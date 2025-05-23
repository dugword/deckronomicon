package game

import (
	"deckronomicon/configs"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// GameState represents the current state of the game.
type GameState struct {
	Battlefield        *Battlefield
	Cheat              bool
	CardPool           string
	CurrentPhase       string
	CurrentPlayer      int
	CurrentStep        string
	Library            *Library
	EventListeners     []EventHandler
	Exile              []*Card
	Graveyard          *Graveyard
	Hand               *Hand
	LandDrop           bool
	LastActionFailed   bool
	Life               int
	ManaPool           *ManaPool
	MaxHandSize        int
	MaxTurns           int
	Message            string
	MessageLog         []string
	Mulligans          int
	PotentialMana      *ManaPool
	SpellsCastThisTurn []string
	StormCount         int
	// TODO: I don't like this, need to rethink how to handle this
	StartingHand   []string
	Turn           int
	TurnCount      int
	TurnMessageLog []string
}

// NewGameState creates a new GameState instance.
func NewGameState() *GameState {
	gameState := GameState{
		Battlefield:        NewBattlefield(),
		Library:            NewLibrary(),
		EventListeners:     []EventHandler{},
		Exile:              []*Card{}, // TODO: make this a struct
		Graveyard:          NewGraveyard(),
		Hand:               NewHand(),
		ManaPool:           NewManaPool(),
		Mulligans:          0,
		PotentialMana:      NewManaPool(),
		SpellsCastThisTurn: []string{}, // TODO: Rethink how this is managed
		TurnMessageLog:     []string{}, // TODO: this sucks, make better
	}
	return &gameState
}

func (g *GameState) DrawStartingHand(startingHand []string) error {
	for _, cardName := range startingHand {
		card, err := g.Library.FindByName(cardName)
		if err != nil {
			return fmt.Errorf("failed to find card %s in library: %w", cardName, err)
		}
		g.Hand.Add(card)
	}
	result, err := g.ResolveAction(GameAction{
		Type:   ActionDraw,
		Target: strconv.Itoa(g.MaxHandSize - g.Hand.Size()),
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to draw starting hand: %w", err)
	}
	g.Log(result.Message)
	return nil
}

// InitializeNewGame initializes a new game with the given configuration.
// TODO: DOn't pass agent here
func (g *GameState) InitializeNewGame(config *configs.Config) error {
	// TODO: Consolidate config file with other configs so they can all be set
	// by either envars or cli flags or config file.
	data, err := os.ReadFile(config.ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	var configFile struct {
		StartingHand []string `json:"StartingHand"`
	}
	if err := json.Unmarshal(data, &configFile); err != nil {
		return fmt.Errorf("failed to unmarshal config file: %w", err)
	}
	g.MaxTurns = config.MaxTurns
	g.Life = config.StartingLife
	g.MaxHandSize = 7
	g.CardPool = config.CardPool
	library, err := importDeck(config.DeckList, g.CardPool)
	if err != nil {
		return err
	}
	g.Library = library
	g.Library.Shuffle()
	g.Hand = &Hand{}
	if err := g.DrawStartingHand(configFile.StartingHand); err != nil {
		return fmt.Errorf("failed to draw starting hand: %w", err)
	}
	g.ManaPool = NewManaPool()
	g.Battlefield = &Battlefield{}
	return nil
}

// Discard discards n cards from the player's hand.
func (g *GameState) Discard(n int, source ChoiceSource, resolver ChoiceResolver) error {
	if n > g.Hand.Size() {
		n = g.Hand.Size()
	}
	for range n {
		choices := CreateObjectChoices(g.Hand.GetAll(), ZoneHand)
		choice, err := resolver.ChooseOne(
			"Which card to discard from hand",
			source,
			choices,
		)
		if err != nil {
			return fmt.Errorf("failed to choose card to discard: %w", err)
		}
		card, err := g.Hand.Get(choice.ID)
		if err != nil {
			return fmt.Errorf("failed to get card from hand: %w", err)
		}
		g.Hand.Remove(card.ID())
		g.Graveyard.Add(card)
	}
	return nil
}

// Draw draws n cards from the library into the player's hand.
func (g *GameState) Draw(n int) error {
	var drawn []GameObject
	var names []string
	for range n {
		card, err := g.Library.TakeTop()
		if err != nil {
			if errors.Is(err, ErrLibraryEmpty) {
				return PlayerLostError{
					Reason: DeckedOut,
				}
			}
			return err
		}
		drawn = append(drawn, card)
		names = append(names, card.Name())
	}
	for _, card := range drawn {
		// TODO: rethink if I want to add one at a time or accept a slice
		g.Hand.Add(card)
	}
	g.Log(fmt.Sprintf("drew: %s", strings.Join(names, ", ")))
	return nil
}

func (g *GameState) Zones() []Zone {
	return []Zone{
		g.Battlefield,
		g.Library,
		g.Hand,
		g.Graveyard,
	}
}

func (g *GameState) GetZone(zone string) (Zone, error) {
	for _, z := range g.Zones() {
		if z.ZoneType() == zone {
			return z, nil
		}
	}
	return nil, fmt.Errorf("zone %s not found", zone)
}

/*
func (g *GameState) CastSpell(card *Card) error {
	// Step 1: Pay costs (mana, sacrifice, etc.)
	// Step 2: Run cast triggers (e.g., storm)
	for _, trigger := range card.CastTriggers {
		if err := trigger(card, g); err != nil {
			return err
		}
	}

	// Step 3: Put spell on stack
	//g.Stack = append(g.Stack, card)
	g.Log = append(g.Log, fmt.Sprintf("Cast %s", card.Name))
	return nil
}
*/

// TODO: Need to do something smarter here, this doesn't account for
// additional mana effects like high tide.
func GetPotentialMana(state *GameState) *ManaPool {
	tempGameState := GameState{
		ManaPool: NewManaPool(),
	}
	for _, permanent := range state.Battlefield.permanents {
		// TODO change this to canpay
		if permanent.IsTapped() {
			continue
		}
		for _, ability := range permanent.ActivatedAbilities() {
			if ability.IsManaAbility() {
				ability.Resolve(&tempGameState, nil) // pass mock resolver
			}
		}
	}
	return tempGameState.ManaPool
}

// TODO Revist this
func PutNBackOnTop(state *GameState, n int, source ChoiceSource, resolver ChoiceResolver) error {
	for range n {
		choices := CreateObjectChoices(state.Hand.GetAll(), ZoneHand)
		choice, err := resolver.ChooseOne(
			"Which card to put back on top",
			source,
			choices,
		)
		if err != nil {
			return fmt.Errorf("failed to choose card to put back on top: %w", err)
		}
		card, err := state.Hand.Take(choice.ID)
		if err != nil {
			return fmt.Errorf("failed to take card from hand: %w", err)
		}
		state.Library.AddTop(card)
	}
	return nil
}

func CanPotentiallyPayFor(state *GameState, manaCost *ManaCost) bool {
	simulated := GetPotentialMana(state)
	for color, need := range manaCost.Colors {
		if simulated.Has(color, need) {
			return false
		}
		simulated.Use(color, need)
	}
	return simulated.HasGeneric(manaCost.Generic)
}

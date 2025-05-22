package game

import (
	"deckronomicon/configs"
	"errors"
	"fmt"
	"strings"
)

// GameState represents the current state of the game.
type GameState struct {
	Battlefield        *Battlefield
	Cheat              bool
	CurrentPhase       string
	CurrentPlayer      int
	CurrentStep        string
	Library            *Library
	EventListeners     []EventHandler
	Exile              []*Card
	Graveyard          []*Card
	Hand               *Hand
	LandDrop           bool
	LastActionFailed   bool
	Life               int
	ManaPool           *ManaPool
	MaxHandSize        int
	MaxTurns           int
	Message            string
	MessageLog         []string
	PotentialMana      *ManaPool
	SpellsCastThisTurn []string
	StormCount         int
	Turn               int
	TurnCount          int
	TurnMessageLog     []string
}

// NewGameState creates a new GameState instance.
func NewGameState() *GameState {
	gameState := GameState{
		Battlefield:        NewBattlefield(),
		Library:            NewLibrary(),
		EventListeners:     []EventHandler{},
		Exile:              []*Card{}, // TODO: make this a struct
		Graveyard:          []*Card{}, // TODO: Make this a struct
		Hand:               NewHand(),
		ManaPool:           NewManaPool(),
		PotentialMana:      NewManaPool(),
		SpellsCastThisTurn: []string{}, // TODO: Rethink how this is managed
		TurnMessageLog:     []string{}, // TODO: this sucks, make better
	}
	return &gameState
}

// InitializeNewGame initializes a new game with the given configuration.
func (g *GameState) InitializeNewGame(config *configs.Config) error {
	g.MaxTurns = config.MaxTurns
	g.Life = config.StartingLife
	g.MaxHandSize = 7
	library, err := importDeck(config.DeckList, config.CardPool)
	if err != nil {
		return err
	}
	g.Library = library
	g.Library.Shuffle()
	g.Hand = &Hand{}
	result, err := g.ResolveAction(GameAction{
		Type:   ActionDraw,
		Target: "7",
	}, nil)
	g.Log(result.Message)
	g.ManaPool = NewManaPool()
	g.Battlefield = &Battlefield{}
	return nil
}

// Discard discards n cards from the player's hand.
func (g *GameState) Discard(n int, source ChoiceSource, resolver ChoiceResolver) error {
	if n > len(g.Hand.Cards()) {
		n = len(g.Hand.Cards())
	}
	for range n {
		choices := g.Hand.CardChoices()
		choice, err := resolver.ChooseOne(
			"Which card to discard from hand",
			source,
			choices,
		)
		if err != nil {
			return fmt.Errorf("failed to choose card to discard: %w", err)
		}
		card, err := g.Hand.GetCard(choice.ID)
		if err != nil {
			return fmt.Errorf("failed to get card from hand: %w", err)
		}
		g.Hand.RemoveCard(card)
		g.Graveyard = append(g.Graveyard, card)
	}
	return nil
}

// Draw draws n cards from the library into the player's hand.
func (g *GameState) Draw(n int) error {
	drawn, err := g.Library.TakeCards(n)
	if errors.Is(err, ErrLibraryEmpty) {
		return PlayerLostError{
			Reason: DeckedOut,
		}
	}
	if err != nil {
		return err
	}
	var names []string
	for _, card := range drawn {
		names = append(names, card.Name())
	}
	g.Hand.Add(drawn...)
	g.Log(fmt.Sprintf("drew: %s", strings.Join(names, ", ")))
	return nil
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
func PutNBackOnTop(state *GameState, n int, source GameObject, resolver ChoiceResolver) error {
	for range n {
		choices := state.Hand.CardChoices()
		choice, err := resolver.ChooseOne(
			"Which card to put back on top",
			source,
			choices,
		)
		if err != nil {
			return fmt.Errorf("failed to choose card to put back on top: %w", err)
		}
		card, err := state.Hand.GetCard(choice.ID)
		if err != nil {
			return fmt.Errorf("failed to get card from hand: %w", err)
		}
		state.Hand.RemoveCard(card)
		state.Library.PutOnTop(card)
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

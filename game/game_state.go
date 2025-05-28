package game

import (
	"deckronomicon/configs"
	"errors"
	"fmt"
	"strings"
)

// GameState represents the current state of the game.
type GameState struct {
	ActivePlayer       *Player
	NonActivePlayer    *Player
	Cheat              bool
	CardPool           string
	CurrentPhase       string
	CurrentPlayer      int
	CurrentStep        string
	EventListeners     []EventHandler
	LastActionFailed   bool
	Players            []*Player
	MaxTurns           int
	Message            string
	MessageLog         []string
	SpellsCastThisTurn []string
	Stack              *Stack
	StormCount         int
	TurnMessageLog     []string
}

// NewGameState creates a new GameState instance.
func NewGameState() *GameState {
	gameState := GameState{
		EventListeners:     []EventHandler{},
		Players:            []*Player{},
		SpellsCastThisTurn: []string{}, // TODO: Rethink how this is managed
		Stack:              NewStack(),
		TurnMessageLog:     []string{}, // TODO: this sucks, make better
	}
	return &gameState
}

func (g *GameState) DrawStartingHand(player *Player) error {
	for _, cardName := range player.StartingHand {
		card, err := player.Library.TakeByName(cardName)
		if err != nil {
			return fmt.Errorf("failed to take card %s from library: %w", cardName, err)
		}
		player.Hand.Add(card)
	}
	if err := g.Draw(player.MaxHandSize-len(player.StartingHand), player); err != nil {
		return fmt.Errorf("failed to draw starting hand: %w", err)
	}
	g.Log(fmt.Sprintf("drawn starting hand for %s: %s", player.ID, strings.Join(player.StartingHand, ", ")))
	return nil
}

// InitializeNewGame initializes a new game with the given configuration.
// TODO: DOn't pass agent here
func (g *GameState) InitializeNewGame(
	scenario *configs.Scenerio,
	player *Player,
	opponent *Player,
	config *configs.Config,
) error {
	g.Cheat = config.Cheat
	g.MaxTurns = scenario.Setup.MaxTurns
	g.CardPool = config.CardPool
	playerLibrary, err := importDeck(scenario.PlayerDeck, g.CardPool)
	if err != nil {
		return err
	}
	player.Library = playerLibrary
	player.Library.Shuffle()
	opponentLibrary, err := importDeck(scenario.OpponentDeck, g.CardPool)
	if err != nil {
		return err
	}
	opponent.Library = opponentLibrary
	opponent.Library.Shuffle()
	// TODO: Do something with the scenario setup for on the play.
	g.Players = []*Player{player, opponent}
	g.ActivePlayer = player
	g.NonActivePlayer = opponent
	// TODO Move this to game engine/mulligan
	if err := g.DrawStartingHand(player); err != nil {
		return fmt.Errorf("failed to draw starting hand: %w", err)
	}
	if err := g.DrawStartingHand(opponent); err != nil {
		return fmt.Errorf("failed to draw starting hand: %w", err)
	}

	return nil
}

// Discard discards n cards from the player's hand.
func (g *GameState) Discard(n int, source ChoiceSource, player *Player) error {
	if n > player.Hand.Size() {
		n = player.Hand.Size()
	}
	for range n {
		choices := CreateChoices(player.Hand.GetAll(), ZoneHand)
		choice, err := player.Agent.ChooseOne(
			"Which card to discard from hand",
			source,
			choices,
		)
		if err != nil {
			return fmt.Errorf("failed to choose card to discard: %w", err)
		}
		card, err := player.Hand.Get(choice.ID)
		if err != nil {
			return fmt.Errorf("failed to get card from hand: %w", err)
		}
		player.Hand.Remove(card.ID())
		player.Graveyard.Add(card)
	}
	return nil
}

// Draw draws n cards from the library into the player's hand.
// TODO: Probably just needs *PLayer
func (g *GameState) Draw(n int, player *Player) error {
	var names []string
	for range n {
		card, err := player.Library.TakeTop()
		if err != nil {
			if errors.Is(err, ErrLibraryEmpty) {
				return PlayerLostError{
					Reason: DeckedOut,
				}
			}
			return err
		}
		player.Hand.Add(card)
		names = append(names, card.Name())
	}
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
func GetPotentialMana(state *GameState, player *Player) *ManaPool {
	// TODO FIX THIS
	return NewManaPool()
	/*
		tempPlayer := Player{
			ManaPool: NewManaPool(),
		}
		for _, permanent := range agent.Player().Battlefield.permanents {
			// TODO change this to canpay
			if permanent.IsTapped() {
				continue
			}
			for _, ability := range permanent.ActivatedAbilities() {
				if ability.IsManaAbility() {
					// TODO
					ability.Resolve(state, agent)
				}
			}
		}
		return tempPlayer.ManaPool
	*/
}

// TODO Revist this
func PutNBackOnTop(state *GameState, n int, source ChoiceSource, player *Player) error {
	for range n {
		choices := CreateChoices(player.Hand.GetAll(), ZoneHand)
		choice, err := player.Agent.ChooseOne(
			"Which card to put back on top",
			source,
			choices,
		)
		if err != nil {
			return fmt.Errorf("failed to choose card to put back on top: %w", err)
		}
		card, err := player.Hand.Take(choice.ID)
		if err != nil {
			return fmt.Errorf("failed to take card from hand: %w", err)
		}
		player.Library.AddTop(card)
	}
	return nil
}

// TODO Rethink this function
// Calculating mana by activating all abilities seems like a bad idea
// because I could impact other game state
func CanPotentiallyPayFor(state *GameState, manaCost *ManaCost) bool {
	return false
	// TODO:
	/*
		simulated := GetPotentialMana(state)
		for color, need := range manaCost.Colors {
			if simulated.Has(color, need) {
				return false
			}
			simulated.Use(color, need)
		}
		return simulated.HasGeneric(manaCost.Generic)
	*/
}

func (g *GameState) CanCastSorcery() bool {
	if g.IsMainPhase() && g.StackIsEmpty() && g.IsPlayerTurn(g.CurrentPlayer) {
		return true
	}
	return false
}

func (g *GameState) IsMainPhase() bool {
	return g.CurrentPhase == PhasePrecombatMain || g.CurrentPhase == PhasePostcombatMain
}

func (g *GameState) StackIsEmpty() bool {
	return g.Stack.Size() == 0
}

func (g *GameState) IsPlayerTurn(playerID int) bool {
	return g.CurrentPlayer == playerID
}

func (g *GameState) GetPlayer(id string) (*Player, error) {
	for _, player := range g.Players {
		if player.ID == id {
			return player, nil
		}
	}
	return nil, fmt.Errorf("player with ID %s not found", id)
}

// TODO This would faile with more than 1 player
func (g *GameState) GetOpponent(id string) (*Player, error) {
	for _, player := range g.Players {
		if player.ID != id {
			return player, nil
		}
	}
	return nil, errors.New("opponent not found")
}

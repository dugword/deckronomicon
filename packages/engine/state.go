package engine

import (
	// 	"deckronomicon/configs"

	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/permanent"
	"deckronomicon/packages/game/player"
	"deckronomicon/packages/game/zone"
	"deckronomicon/packages/query"
	"errors"
	"fmt"
)

// GameState represents the current state of the game.
type GameState struct {
	activePlayer    int
	battlefield     *zone.Battlefield
	CheatsEnabled   bool
	CardPool        string
	CardDefinitions map[string]definition.Card
	CurrentPhase    mtg.Phase
	CurrentStep     mtg.Step
	// EventListeners     []EventHandler
	LastActionFailed   bool
	players            []*player.Player
	MaxTurns           int
	Message            string
	MessageLog         []string
	SpellsCastThisTurn []string
	stack              *zone.Stack
	StormCount         int
	TurnMessageLog     []string
	nextID             int
}

func NewGameState() *GameState {
	// TODO hand on the play
	gameState := GameState{
		battlefield: zone.NewBattlefield(),
		players:     []*player.Player{},
		// EventListeners:     []EventHandler{},
		SpellsCastThisTurn: []string{}, // TODO: Rethink how this is managed
		stack:              zone.NewStack(),
		TurnMessageLog:     []string{}, // TODO: this sucks, make better
	}
	return &gameState
}

func (g *GameState) AddToBattlefield(perm *permanent.Permanent) {
	g.battlefield.Add(perm)
}

func (g *GameState) AddToStack(item zone.Resolvable) {
	g.stack.Add(item)
}

func (g *GameState) Battlefield() query.View {
	return query.NewView(g.battlefield.Name(), g.battlefield.GetAll())
}

func (g *GameState) CanCastSorcery(playerID string) bool {
	if g.IsMainPhase() && g.StackIsEmpty() && g.IsPlayerActive(playerID) {
		return true
	}
	return false
}

func (g *GameState) GetActivePlayer() *player.Player {
	return g.players[g.activePlayer]
}

func (g *GameState) GetPlayer(id string) (*player.Player, error) {
	for _, player := range g.players {
		if player.ID() == id {
			return player, nil
		}
	}
	return nil, fmt.Errorf("player with ID %s not found", id)
}

// TODO: This will fail with more than 2 players
func (g *GameState) GetOpponent(id string) (*player.Player, error) {
	if len(g.players) != 2 {
		panic("GetOpponent called with more than 2 players")
	}
	for _, player := range g.players {
		if player.ID() != id {
			return player, nil
		}
	}
	return nil, errors.New("opponent not found")
}

func (g *GameState) IsMainPhase() bool {
	return g.CurrentPhase == mtg.PhasePrecombatMain ||
		g.CurrentPhase == mtg.PhasePostcombatMain
}

func (g *GameState) IsPlayerActive(playerID string) bool {
	return g.GetActivePlayer().ID() == playerID
}

func (g *GameState) NextPlayerTurn() {
	g.activePlayer++
	if g.activePlayer >= len(g.players) {
		g.activePlayer = 0
	}
}

func (g *GameState) Players() []*player.Player {
	return g.players
}

func (g *GameState) PopFromStack() (zone.Resolvable, error) {
	if g.StackIsEmpty() {
		return nil, errors.New("stack is empty")
	}
	item, err := g.stack.Pop()
	if err != nil {
		return nil, fmt.Errorf("failed to pop from stack: %w", err)
	}
	return item, nil
}

func (g *GameState) StackIsEmpty() bool {
	return g.stack.Size() == 0
}

func (g *GameState) Stack() query.View {
	return query.NewView(g.stack.Name(), g.stack.GetAll())
}

// TODO: Need to do something smarter here, this doesn't account for
// additional mana effects like high tide.
/*
func GetPotentialMana(state *GameState, player *Player) *ManaPool {
	// TODO FIX THIS
	return NewManaPool()
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
}
*/

/*
// TODO Rethink this function
// Calculating mana by activating all abilities seems like a bad idea
// because I could impact other game state
func CanPotentiallyPayFor(state *GameState, manaCost *ManaCost) bool {
	return false
	// TODO:
		simulated := GetPotentialMana(state)
		for color, need := range manaCost.Colors {
			if simulated.Has(color, need) {
				return false
			}
			simulated.Use(color, need)
		}
		return simulated.HasGeneric(manaCost.Generic)
}
*/

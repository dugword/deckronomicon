package engine

import (
	"deckronomicon/packages/configs"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/player"
	"deckronomicon/packages/game/zone"
	"fmt"
	"strings"
)

type Engine struct {
	Scenario        *configs.Scenario
	GameState       *GameState
	CardDefinitions map[string]definition.Card
}

// TODO: Consolidate with top level log package.
// Log logs a message to the game state message log.
func (e *Engine) Log(message ...string) {
	// TODO: Something like this so we know which player to which action
	// outMessage := fmt.Sprintf("%s: %s", g.ActivatePlayer.ID )
	// TODO: There's probably a more elegant way to do this
	e.GameState.TurnMessageLog = append(
		e.GameState.TurnMessageLog,
		strings.Join(message, " "),
	)
	e.GameState.MessageLog = append(
		e.GameState.MessageLog,
		strings.Join(message, " "),
	)
}

// Error logs an error message to the game state message log.
// TODO: There's probably a more elegant way to do this
func (e *Engine) Error(err error) {
	e.GameState.TurnMessageLog = append(
		e.GameState.TurnMessageLog,
		"ERROR: "+err.Error(),
	)
	e.GameState.MessageLog = append(
		e.GameState.MessageLog,
		"ERROR: "+err.Error(),
	)
}

// TODO: Figure out better naming so I don't have plyr
// InitializeNewGame initializes a new game with the given configuration.
func InitializeNewGame(
	scenario *configs.Scenario,
	players []*player.Player,
	cardDefinitions map[string]definition.Card,
) *Engine {
	gameState := NewGameState()
	gameState.players = players
	for i, plyer := range players {
		if plyer.ID() == scenario.Setup.OnThePlay {
			gameState.activePlayer = i
			break
		}
	}
	gameState.CheatsEnabled = scenario.Setup.CheatsEnabled
	engine := Engine{
		CardDefinitions: cardDefinitions,
		GameState:       gameState,
		Scenario:        scenario,
	}
	return &engine
}

// TODO: Rename this to something more descriptive, like get players ready to
// play the game with library and hand, or something more sucient.
// PlayerSetup sets up the players for the game. It builds their libraries,
// draws their starting hands, and handles mulligans.
func (e *Engine) PlayerSetup() error {
	for _, plyer := range e.GameState.players {
		// TODO: Should this be a method on Engine?
		playerLibrary, err := zone.BuildLibrary(
			e.GameState,
			e.Scenario.Players[plyer.ID()].DeckList,
			e.CardDefinitions,
		)
		if err != nil {
			return err
		}
		plyer.AssignLibrary(playerLibrary)
		plyer.ShuffleLibrary()
		startingHand := e.Scenario.Players[plyer.ID()].StartingHand
		cardNames, err := DrawStartingHand(
			plyer,
			startingHand,
		)
		if err != nil {
			return fmt.Errorf(
				"%s failed to draw starting hand: %w",
				plyer.ID(),
				err,
			)
		}
		e.Log(fmt.Sprintf(
			"% drew: %s",
			plyer.ID(),
			strings.Join(cardNames, ", "),
		))
		if err := Mulligan(e.GameState, plyer, startingHand); err != nil {
			return fmt.Errorf(
				"player %s failed to mulligan: %w",
				plyer.ID(),
				err,
			)
		}
		// TODO: Implement Logf
		e.Log(fmt.Sprintf(
			"Player %s mulliganned %d times",
			plyer.ID(),
			plyer.Mulligans,
		))
	}
	return nil
}

// RunGameLoop drives the turn-based game.
func (e *Engine) RunGameLoop() error {
	if err := e.PlayerSetup(); err != nil {
		return fmt.Errorf("failed to setup players: %w", err)
	}
	for {
		if e.GameState.GetActivePlayer().Turn > e.Scenario.MaxTurns {
			return ErrMaxTurnsExceeded
		}
		if err := e.RunTurn(e.GameState.GetActivePlayer()); err != nil {
			return err
		}
		e.GameState.NextPlayerTurn()
	}
}

// RunTurn executes a single game turn using the provided player.
// Maybe split into ResolvePhase and ResolveStep functions
// RunTurn executes a single game turn using the provided player
func (e *Engine) RunTurn(player *player.Player) error {
	player.Turn++
	player.LandDrop = false // TODO Move this
	// player.Battlefield.RemoveSummoningSickness()
	e.GameState.TurnMessageLog = []string{} // TODO: this sucks, make better
	e.Log(fmt.Sprintf("%s's turn %d", player.ID(), player.Turn))
	if err := e.RunPhase(e.beginningPhase(), player); err != nil {
		return fmt.Errorf("failed to run beginning phase: %w", err)
	}
	if err := e.RunPhase(e.precombatMainPhase(), player); err != nil {
		return fmt.Errorf("failed to run precombat main phase: %w", err)
	}
	if err := e.RunPhase(e.combatPhase(), player); err != nil {
		return fmt.Errorf("failed to run combat phase: %w", err)
	}
	if err := e.RunPhase(e.postcombatMainPhase(), player); err != nil {
		return fmt.Errorf("failed to run postcombat main phase: %w", err)
	}
	if err := e.RunPhase(e.endingPhase(), player); err != nil {
		return fmt.Errorf("failed to run ending phase: %w", err)
	}
	return nil
}

func (e *Engine) RunPhase(phase *GamePhase, player *player.Player) error {
	e.GameState.CurrentPhase = phase.Name
	e.Log(fmt.Sprintf("Running phase: %s", phase.Name))
	for _, step := range phase.Steps {
		if err := e.RunStep(step, player); err != nil {
			return fmt.Errorf("failed to run phase %s step %s: %w", phase.Name, step.Name, err)
		}
	}
	e.Log(fmt.Sprintf("Phase completed: %s", phase.Name))
	return nil
}

func (e *Engine) RunStep(step GameStep, player *player.Player) error {
	e.GameState.CurrentStep = step.Name
	e.Log(fmt.Sprintf("Running step: %s", step.Name))
	// g.EmitEvent(Event{Type: step.EventEvent}, player)
	if err := step.Handler(e.GameState, player); err != nil {
		return fmt.Errorf("failed to run step %s: %w", step.Name, err)
	}
	player.EmptyManaPool()
	e.Log(fmt.Sprintf("Step completed: %s", step.Name))
	return nil
}

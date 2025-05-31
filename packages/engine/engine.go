package engine

import (
	"deckronomicon/packages/configs"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/player"
	"deckronomicon/packages/game/zone"
	"deckronomicon/packages/query/has"
	"fmt"
	"strings"
)

type Engine struct {
	Scenario        *configs.Scenario
	GameState       *GameState
	CardDefinitions map[string]definition.Card
}

// TODO: Figure out better naming so I don't have plyr
// InitializeNewGame initializes a new game with the given configuration.
func InitializeNewGame(
	scenario *configs.Scenario,
	players []*player.Player,
	cardDefinitions map[string]definition.Card,
) *Engine {
	gameState := NewGameState()
	gameState.Players = players
	for i, plyer := range players {
		if plyer.ID() == scenario.Setup.OnThePlay {
			gameState.ActivePlayer = i
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

// Log logs a message to the game state message log.
func (e *Engine) Log(message ...string) {
	// TODO: Something like this so we know which player to which action
	// outMessage := fmt.Sprintf("%s: %s", g.ActivatePlayer.ID )
	// TODO: There's probably a more elegant way to do this
	e.GameState.TurnMessageLog = append(e.GameState.TurnMessageLog, strings.Join(message, " "))
	e.GameState.MessageLog = append(e.GameState.MessageLog, strings.Join(message, " "))
}

// Error logs an error message to the game state message log.
// TODO: There's probably a more elegant way to do this
func (e *Engine) Error(err error) {
	e.GameState.TurnMessageLog = append(e.GameState.TurnMessageLog, "ERROR: "+err.Error())
	e.GameState.MessageLog = append(e.GameState.MessageLog, "ERROR: "+err.Error())
}

// Mulligan allows the player to mulligan their hand. The player draws 7 new
// cards but then needs to put 1 card back on the bottom of their library for
// each time they've mulliganed.
func Mulligan(state *GameState, player *player.Player, startingHand []string) error {
	var accept bool
	for (player.Mulligans < player.MaxHandSize) || !accept {
		player.Agent.ReportState(state)
		accept, err := player.Agent.Confirm("Keep Hand? (y/n)", nil)
		if err != nil {
			return fmt.Errorf("failed to confirm mulligan: %w", err)
		}
		if accept {
			break
		}
		for _, card := range player.Hand.GetAll() {
			player.Hand.Remove(card.ID())
			player.Library.Add(card)
		}
		player.Library.Shuffle()
		DrawStartingHand(player, startingHand)
		player.Mulligans++
	}
	if player.Mulligans != 0 {
		/*
			if err := PutNBackOnTop(g, player.Mulligans, NewChoiceSource(ChoiceMulligan, ChoiceMulligan), player); err != nil {
				return fmt.Errorf("failed to put cards back on top: %w", err)
			}
		*/
	}
	return nil
}

func DrawStartingHand(plyer *player.Player, startingHand []string) ([]string, error) {
	var cardsDrawn []string
	for _, cardName := range startingHand {
		if err := plyer.Tutor(has.Name(cardName)); err != nil {
			return nil, fmt.Errorf("failed to tutor card %s: %w", cardName, err)
		}
		cardsDrawn = append(cardsDrawn, cardName)
	}
	remainingDraws := plyer.MaxHandSize - plyer.Hand.Size()
	for range remainingDraws {
		cardName, err := plyer.Draw()
		if err != nil {
			return nil, fmt.Errorf("failed to draw card: %w", err)
		}
		cardsDrawn = append(cardsDrawn, cardName)
	}
	return cardsDrawn, nil
}

// RunGameLoop drives the turn-based game.
func (e *Engine) RunGameLoop() error {
	for _, plyer := range e.GameState.Players {
		// TODO: Should this be a method on Engine?
		playerLibrary, err := zone.BuildLibrary(
			e.GameState,
			e.Scenario.Players[plyer.ID()].DeckList,
			e.CardDefinitions,
		)
		if err != nil {
			return err
		}
		plyer.Library = playerLibrary
		plyer.Library.Shuffle()
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
	}
	for _, player := range e.GameState.Players {
		// TODO: Implement Logf
		e.Log(fmt.Sprintf(
			"Player %s mulliganned %d times",
			player.ID(),
			player.Mulligans,
		))
	}
	if e.GameState.GetActivePlayer().Turn > e.Scenario.MaxTurns {
		return ErrMaxTurnsExceeded
	}
	if err := e.RunTurn(
		e.GameState.GetActivePlayer(),
	); err != nil {
		return err
	}
	e.GameState.NextPlayerTurn()
	return nil
}

func (e *Engine) RunStep(step GameStep, player *player.Player) error {
	e.GameState.CurrentStep = step.Name
	e.Log(fmt.Sprintf("Running step: %s", step.Name))
	// g.EmitEvent(Event{Type: step.EventEvent}, player)
	if err := step.Handler(e.GameState, player); err != nil {
		return fmt.Errorf("failed to run step %s: %w", step.Name, err)
	}
	player.ManaPool.Empty()
	e.Log(fmt.Sprintf("Step completed: %s", step.Name))
	return nil
}

func (e *Engine) RunPhase(phase GamePhase, player *player.Player) error {
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

var ChoiceMulligan = "Mulligan"

// var ChoiceSourceMulligan = NewChoiceSource(ChoiceMulligan, ChoiceMulligan)

// RunTurn executes a single game turn using the provided player.
// Maybe split into ResolvePhase and ResolveStep functions
// RunTurn executes a single game turn using the provided player
func (e *Engine) RunTurn(player *player.Player) error {
	// Beginning Phase
	player.Turn++
	player.LandDrop = false
	// player.Battlefield.RemoveSummoningSickness()
	e.GameState.TurnMessageLog = []string{} // TODO: this sucks, make better
	e.Log(fmt.Sprintf("Turn %d", player.Turn))
	if err := e.RunPhase(beginningPhase, player); err != nil {
		return fmt.Errorf("failed to run beginning phase: %w", err)
	}
	if err := e.RunPhase(precombatMainPhase, player); err != nil {
		return fmt.Errorf("failed to run precombat main phase: %w", err)
	}
	if err := e.RunPhase(combatPhase, player); err != nil {
		return fmt.Errorf("failed to run combat phase: %w", err)
	}
	if err := e.RunPhase(postcombatMainPhase, player); err != nil {
		return fmt.Errorf("failed to run postcombat main phase: %w", err)
	}
	if err := e.RunPhase(endingPhase, player); err != nil {
		return fmt.Errorf("failed to run ending phase: %w", err)
	}
	return nil
}

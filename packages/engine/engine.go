package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/logger"
	"deckronomicon/packages/state"
	"fmt"
)

type Engine struct {
	agents map[string]PlayerAgent
	game   state.Game
	record *GameRecord
	rng    *RNG
	log    *logger.Logger
}

type EngineConfig struct {
	Players []state.Player
	Agents  map[string]PlayerAgent
	Seed    int64
	// Cards are just strings for now, but will be a Card type later
	DeckLists map[string][]gob.Card
}

func NewEngine(config EngineConfig) *Engine {
	agents := map[string]PlayerAgent{}
	for id, agent := range config.Agents {
		agents[id] = agent
	}
	return &Engine{
		agents: agents,
		game:   state.Game{}.WithPlayers(config.Players),
		log:    &logger.Logger{},
		record: NewGameRecord(config.Seed),
		rng:    NewRNG(config.Seed),
	}
}

func (e *Engine) RunGame() error {
	e.log.IncludeContext = true
	e.log.ContextFunc = func() any {
		return nil
	}
	e.log.Debug("Running game")
	if err := e.Apply(event.NewBeginGameEvent()); err != nil {
		return fmt.Errorf("error starting game: %w", err)
	}
	// shuffle all player decks
	// draw initial hands for players
	// resovle mulligans
	for !e.game.IsGameOver() {
		err := e.RunTurn()
		if err != nil {
			return err
		}
		if err := e.Apply(event.SetNextPlayerEvent{}); err != nil {
			return fmt.Errorf("error completing turn: %w", err)
		}
	}
	// e.Apply(event.NewGameOverEvent())
	return nil
}

func (e *Engine) RunTurn() error {
	activePlayerID := e.game.ActivePlayerID()
	e.log.Debug("Running turn for player: ", activePlayerID)
	if err := e.Apply(event.NewBeginTurnEvent(activePlayerID)); err != nil {
		return fmt.Errorf("error starting turn: %w", err)
	}
	for _, phase := range e.GamePhases() {
		if err := e.RunPhase(phase); err != nil {
			return fmt.Errorf("error running phase %s: %w", phase.name, err)
		}
	}
	if err := e.Apply(event.NewEndTurnEvent()); err != nil {
		return fmt.Errorf("error ending turn: %w", err)
	}
	return nil
}

func (e *Engine) RunPhase(phase GamePhase) error {
	e.log.Debug("Running phase:", phase.name)
	for _, step := range phase.steps {
		if err := e.RunStep(step); err != nil {
			return fmt.Errorf("error running step %s: %w", step.name, err)
		}
	}
	return nil
}

func (e *Engine) RunStep(step GameStep) error {
	e.log.Debug("Running step:", step.name)
	if err := e.Apply(event.NewBeginStepEvent(step.name)); err != nil {
		return fmt.Errorf("error starting step %s: %w", step.name, err)
	}
	for _, action := range step.actions {
		e.log.Debug("Completing action:", action.Name())
		event, err := action.Complete(e.game, choose.Choice{})
		if err != nil {
			return fmt.Errorf(
				"error completing action %s: %w",
				action.Name(),
				err,
			)
		}
		if err := e.Apply(event); err != nil {
			return fmt.Errorf("error applying event: %w", err)
		}
	}
	if err := e.RunPriority(); err != nil {
		return fmt.Errorf("error running priority: %w", err)
	}
	if err := e.Apply(event.NewEndStepEvent(step.name)); err != nil {
		return fmt.Errorf("error ending step %s: %w", step.name, err)
	}
	return nil
}

func (e *Engine) RunPriority() error {
	priorityPlayerID := e.game.ActivePlayerID()
	if err := e.Apply(
		event.ReceivePriorityEvent{PlayerID: priorityPlayerID},
	); err != nil {
		return fmt.Errorf("error applying receive priority event: %w", err)
	}
	for !e.game.AllPlayersPassedPriority() {
		priorityPlayerID = e.game.PriorityPlayerID()
		e.log.Debug("Player %s received priority", priorityPlayerID)
		agent := e.agents[priorityPlayerID]
		action, err := agent.GetNextAction()
		if err != nil {
			return fmt.Errorf(
				"error getting next action for player %s: %w",
				priorityPlayerID,
				err,
			)
		}
		evnt, err := action.Complete(e.game, choose.Choice{})
		if err != nil {
			e.log.Error("Error completing action:", err)
		}
		if err := e.Apply(evnt); err != nil {
			return fmt.Errorf("error applying event: %w", err)
		}
		if !e.game.PlayerPassedPriority(priorityPlayerID) {
			if err := e.Apply(event.ResetPriorityPassesEvent{}); err != nil {
				return fmt.Errorf("error resetting priority passes: %w", err)
			}
		}
		if e.game.AllPlayersPassedPriority() {
			if err := e.Apply(event.AllPlayersPassedPriorityEvent{}); err != nil {
				return fmt.Errorf("error applying all players passed priority event: %w", err)
			}
		}
	}
	return nil
}

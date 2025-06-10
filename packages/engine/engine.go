package engine

// Engine manages the control flow of the game, including running turns, phases, and steps.

import (
	"deckronomicon/packages/configs"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/logger"
	"deckronomicon/packages/state"
	"fmt"
)

type ResolutionEnvironment struct {
	EffectRegistry *EffectRegistry

	// definitions maybe should live here too
}

type Engine struct {
	agents                map[string]PlayerAgent
	definitions           map[string]definition.Card
	deckLists             map[string]configs.DeckList
	game                  state.Game
	record                *GameRecord
	rng                   *RNG
	log                   *logger.Logger
	resolutionEnvironment *ResolutionEnvironment
}

type EngineConfig struct {
	Players []state.Player
	Agents  map[string]PlayerAgent
	Seed    int64
	// Cards are just strings for now, but will be a Card type later
	DeckLists   map[string]configs.DeckList
	Definitions map[string]definition.Card
}

func NewEngine(config EngineConfig) *Engine {
	agents := map[string]PlayerAgent{}
	for id, agent := range config.Agents {
		agents[id] = agent
	}
	resolutionEnvironment := ResolutionEnvironment{
		EffectRegistry: NewEffectRegistry(),
	}
	return &Engine{
		agents:                agents,
		definitions:           config.Definitions,
		deckLists:             config.DeckLists,
		game:                  state.Game{}.WithPlayers(config.Players),
		log:                   &logger.Logger{},
		record:                NewGameRecord(config.Seed),
		rng:                   NewRNG(config.Seed),
		resolutionEnvironment: &resolutionEnvironment,
	}
}

func (e *Engine) RunGame() error {
	e.log.IncludeContext = true
	e.log.ContextFunc = func() any {
		return nil
	}
	// TODO This shouldn't live here. It should be in the Apply reducers and managed via events.
	// Or it needs to be created prior to the game starting and passed in.
	for _, playerID := range e.game.PlayerIDsInTurnOrder() {
		e.log.Debug("Setting up player deck:", playerID)
		deckList, ok := e.deckLists[playerID]
		if !ok {
			return fmt.Errorf("deck list for player %q not found", playerID)
		}
		newGame, deck, err := e.game.WithBuildDeck(
			deckList,
			e.definitions,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to build deck for player %q: %w",
				playerID,
				err,
			)
		}
		e.game = newGame
		player, ok := e.game.GetPlayer(playerID)
		if !ok {
			return fmt.Errorf("player %q not found", playerID)
		}
		newPlayer := player.WithLibrary(state.NewLibrary(deck))
		e.game = e.game.WithUpdatedPlayer(newPlayer)
	}

	e.log.Debug("Running game")
	if err := e.Apply(event.BeginGameEvent{}); err != nil {
		return fmt.Errorf("failed to start game: %w", err)
	}
	e.log.Debug("Shuffling decks")
	for _, playerID := range e.game.PlayerIDsInTurnOrder() {
		if err := e.Apply(event.ShuffleDeckEvent{PlayerID: playerID}); err != nil {
			return fmt.Errorf("failed to shuffle decks for player %q: %w", playerID, err)
		}
	}
	for _, playerID := range e.game.PlayerIDsInTurnOrder() {
		e.log.Debug("Drawing starting hand for player:", playerID)
		startingHandAction := DrawStartingHandAction{playerID: playerID}
		if err := e.CompleteAction(startingHandAction); err != nil {
			return fmt.Errorf(
				"failed to draw starting hand for player %q: %w",
				playerID,
				err,
			)
		}
	}
	// resolve mulligans
	for !e.game.IsGameOver() {
		err := e.RunTurn()
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) RunTurn() error {
	activePlayerID := e.game.ActivePlayerID()
	e.log.Debug("Running turn for player: ", activePlayerID)
	if err := e.Apply(event.BeginTurnEvent{PlayerID: activePlayerID}); err != nil {
		return fmt.Errorf("failed to start turn: %w", err)
	}
	for _, phase := range e.GamePhases() {
		if err := e.RunPhase(phase); err != nil {
			return fmt.Errorf("failed to run phase %s: %w", phase.name, err)
		}
	}
	if err := e.Apply(event.EndTurnEvent{PlayerID: activePlayerID}); err != nil {
		return fmt.Errorf("failed to end turn: %w", err)
	}
	// TODO: Not sure if I like this here
	activePlayerID = e.game.NextPlayerID(activePlayerID)
	// TODO: Move this to the start of the turn
	if err := e.Apply(event.SetActivePlayerEvent{PlayerID: activePlayerID}); err != nil {
		return fmt.Errorf("failed to complete turn: %w", err)
	}
	return nil
}

func (e *Engine) RunPhase(phase GamePhase) error {
	e.log.Debug("Running phase:", phase.name)
	if err := e.Apply(event.NewBeginPhaseEvent(phase.name, e.game.ActivePlayerID())); err != nil {
		return fmt.Errorf("failed to start phase %s: %w", phase.name, err)
	}
	for _, step := range phase.steps {
		if err := e.RunStep(step); err != nil {
			return fmt.Errorf("failed to run step %s: %w", step.name, err)
		}
	}
	if err := e.Apply(event.NewEndPhaseEvent(phase.name, e.game.ActivePlayerID())); err != nil {
		return fmt.Errorf("failed to end phase %s: %w", phase.name, err)
	}
	return nil
}

func (e *Engine) RunStep(step GameStep) error {
	e.log.Debug("Running step:", step.name)
	if err := e.Apply(event.NewBeginStepEvent(step.name, e.game.ActivePlayerID())); err != nil {
		return fmt.Errorf("failed to start step %s: %w", step.name, err)
	}
	for _, action := range step.actions {
		e.log.Debug("Completing action:", action.Name())

		if err := e.CompleteAction(action); err != nil {
			return fmt.Errorf(
				"failed to apply action %s: %w",
				action.Name(),
				err,
			)
		}
	}
	if err := e.RunPriority(); err != nil {
		return fmt.Errorf("failed to run priority: %w", err)
	}
	if err := e.Apply(event.NewEndStepEvent(step.name, e.game.ActivePlayerID())); err != nil {
		return fmt.Errorf("failed to end step %s: %w", step.name, err)
	}
	return nil
}

func (e *Engine) RunPriority() error {
	priorityPlayerID := e.game.ActivePlayerID()
	if err := e.Apply(
		event.ReceivePriorityEvent{PlayerID: priorityPlayerID},
	); err != nil {
		return fmt.Errorf("failed to apply receive priority event: %w", err)
	}
	for !e.game.AllPlayersPassedPriority() {
		priorityPlayerID = e.game.PriorityPlayerID()
		e.log.Debugf("Player %q received priority", priorityPlayerID)
		agent := e.agents[priorityPlayerID]
		action, err := agent.GetNextAction(e.game)
		if err != nil {
			return fmt.Errorf(
				"failed to get next action for player %q: %w",
				priorityPlayerID,
				err,
			)
		}
		if err := e.CompleteAction(action); err != nil {
			return fmt.Errorf(
				"failed to apply action for player %q: %w",
				priorityPlayerID,
				err,
			)
		}
		if e.game.PlayerPassedPriority(priorityPlayerID) {
			nextPlayerIDWithPriority := e.game.NextPlayerID(priorityPlayerID)
			if err := e.Apply(event.ReceivePriorityEvent{
				PlayerID: nextPlayerIDWithPriority,
			}); err != nil {
				return fmt.Errorf("failed to apply receive priority event: %w", err)
			}
		} else {
			if err := e.Apply(event.ResetPriorityPassesEvent{}); err != nil {
				return fmt.Errorf("failed to reset priority passes: %w", err)
			}
		}
		if e.game.AllPlayersPassedPriority() {
			if err := e.Apply(event.AllPlayersPassedPriorityEvent{}); err != nil {
				return fmt.Errorf("failed to apply all players passed priority event: %w", err)
			}
		}
	}
	return nil
}

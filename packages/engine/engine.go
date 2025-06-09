package engine

import (
	"deckronomicon/packages/choose"
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

func (e *Engine) ApplyAction(action Action) error {
	choicePrompt, err := action.GetPrompt(e.game)
	if err != nil {
		return fmt.Errorf(
			"error getting choice prompt for action %s: %w",
			action.Name(),
			err,
		)
	}
	choices := []choose.Choice{}
	if len(choicePrompt.Choices) != 0 {
		cs, err := e.agents[action.PlayerID()].Choose(choicePrompt)
		if err != nil {
			return fmt.Errorf(
				"error getting choices for action %s: %w",
				action.Name(),
				err,
			)
		}
		choices = cs
	}
	evnts, err := action.Complete(e.game, e.resolutionEnvironment, choices)
	if err != nil {
		return fmt.Errorf(
			"error completing action %s: %w",
			action.Name(),
			err,
		)
	}
	for _, evnt := range evnts {
		if err := e.Apply(evnt); err != nil {
			return fmt.Errorf(
				"error applying event %s: %w",
				evnt.EventType(),
				err,
			)
		}
	}
	return nil
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
			return fmt.Errorf("deck list for player %s not found", playerID)
		}
		newGame, deck, err := e.game.WithBuildDeck(
			deckList,
			e.definitions,
		)
		if err != nil {
			return fmt.Errorf(
				"error building deck for player %s: %w",
				playerID,
				err,
			)
		}
		e.game = newGame
		player, ok := e.game.GetPlayer(playerID)
		if !ok {
			return fmt.Errorf("player '%s' not found", playerID)
		}
		newPlayer := player.WithLibrary(state.NewLibrary(deck))
		e.game = e.game.WithUpdatedPlayer(newPlayer)
	}

	e.log.Debug("Running game")
	if err := e.Apply(event.NewBeginGameEvent()); err != nil {
		return fmt.Errorf("error starting game: %w", err)
	}
	e.log.Debug("Shuffling decks")
	for _, playerID := range e.game.PlayerIDsInTurnOrder() {
		if err := e.Apply(event.NewShuffDeckEvent(playerID)); err != nil {
			return fmt.Errorf("error shuffling decks: %w", err)
		}
	}
	for _, playerID := range e.game.PlayerIDsInTurnOrder() {
		e.log.Debug("Drawing starting hand for player:", playerID)
		startingHandAction := DrawStartingHandAction{playerID: playerID}
		if err := e.ApplyAction(startingHandAction); err != nil {
			return fmt.Errorf(
				"error drawing starting hand for player %s: %w",
				playerID,
				err,
			)
		}
	}
	// resovle mulligans
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
	if err := e.Apply(event.SetNextPlayerEvent{}); err != nil {
		return fmt.Errorf("error completing turn: %w", err)
	}
	return nil
}

func (e *Engine) RunPhase(phase GamePhase) error {
	e.log.Debug("Running phase:", phase.name)
	if err := e.Apply(event.NewBeginPhaseEvent(phase.name)); err != nil {
		return fmt.Errorf("error starting phase %s: %w", phase.name, err)
	}
	for _, step := range phase.steps {
		if err := e.RunStep(step); err != nil {
			return fmt.Errorf("error running step %s: %w", step.name, err)
		}
	}
	if err := e.Apply(event.NewEndPhaseEvent(phase.name)); err != nil {
		return fmt.Errorf("error ending phase %s: %w", phase.name, err)
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

		if err := e.ApplyAction(action); err != nil {
			return fmt.Errorf(
				"error applying action %s: %w",
				action.Name(),
				err,
			)
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
		e.log.Debugf("Player '%s' received priority", priorityPlayerID)
		agent := e.agents[priorityPlayerID]
		action, err := agent.GetNextAction(e.game)
		if err != nil {
			return fmt.Errorf(
				"error getting next action for player '%s': %w",
				priorityPlayerID,
				err,
			)
		}
		if err := e.ApplyAction(action); err != nil {
			return fmt.Errorf(
				"error applying action for player '%s': %w",
				priorityPlayerID,
				err,
			)
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

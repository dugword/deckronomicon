package engine

// Engine manages the control flow of the game, including running turns, phases, and steps.

import (
	"deckronomicon/packages/configs"
	"deckronomicon/packages/engine/effect"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/engine/rng"
	"deckronomicon/packages/engine/turnaction"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/logger"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

type Engine struct {
	agents      map[string]PlayerAgent
	deckLists   map[string]configs.DeckList
	game        state.Game
	record      *GameRecord
	rng         *rng.RNG
	log         *logger.Logger
	definitions map[string]definition.Card
	resEnv      *resenv.ResEnv
	// effectRegistry *effect.EffectRegistry
}

/*
type ResolutionEnvironment struct {
	RNG         *RNG
	Log         *logger.Logger
	Definitions map[string]definition.Card
	//effectRegistry *effect.EffectRegistry
}
*/

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
	rng := rng.NewRNG(config.Seed)
	log := &logger.Logger{}
	definitions := config.Definitions
	return &Engine{
		agents:    agents,
		deckLists: config.DeckLists,
		game:      state.Game{}.WithPlayers(config.Players),
		log:       log,
		record:    NewGameRecord(config.Seed),
		// rng:         rng,
		definitions: definitions,
		resEnv: &resenv.ResEnv{
			RNG:         rng,
			Log:         log,
			Definitions: definitions,
		},
		//effectRegistry: effect.NewEffectRegistry(),
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
			playerID,
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
	if err := e.ApplyEvent(event.BeginGameEvent{}); err != nil {
		return fmt.Errorf("failed to start game: %w", err)
	}
	e.log.Debug("Shuffling decks")
	for _, playerID := range e.game.PlayerIDsInTurnOrder() {
		player, ok := e.game.GetPlayer(playerID)
		if !ok {
			return fmt.Errorf("player %q not found", playerID)
		}
		shuffledCardsIDs := e.resEnv.RNG.ShuffleCardsIDs(player.Library().GetAll())
		if err := e.ApplyEvent(event.ShuffleLibraryEvent{PlayerID: playerID, ShuffledCardsIDs: shuffledCardsIDs}); err != nil {
			return fmt.Errorf("failed to shuffle decks for player %q: %w", playerID, err)
		}
	}
	for _, playerID := range e.game.PlayerIDsInTurnOrder() {
		e.log.Debug("Drawing starting hand for player:", playerID)
		startingHandAction := turnaction.NewDrawStartingHandAction(playerID)
		// TODO: This could probably just be an event, or maybe a separate type than "action"
		if err := e.CompleteTurnAction(startingHandAction); err != nil {
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
	if err := e.ApplyEvent(event.BeginTurnEvent{PlayerID: activePlayerID}); err != nil {
		return fmt.Errorf("failed to start turn: %w", err)
	}
	for _, phase := range e.GamePhases() {
		if err := e.RunPhase(phase); err != nil {
			return fmt.Errorf("failed to run phase %s: %w", phase.name, err)
		}
	}
	if err := e.ApplyEvent(event.EndTurnEvent{PlayerID: activePlayerID}); err != nil {
		return fmt.Errorf("failed to end turn: %w", err)
	}
	// TODO: Not sure if I like this here
	activePlayerID = e.game.NextPlayerID(activePlayerID)
	// TODO: Move this to the start of the turn
	if err := e.ApplyEvent(event.SetActivePlayerEvent{PlayerID: activePlayerID}); err != nil {
		return fmt.Errorf("failed to complete turn: %w", err)
	}
	return nil
}

func (e *Engine) RunPhase(phase GamePhase) error {
	e.log.Debug("Running phase:", phase.name)
	if err := e.ApplyEvent(event.NewBeginPhaseEvent(phase.name, e.game.ActivePlayerID())); err != nil {
		return fmt.Errorf("failed to start phase %s: %w", phase.name, err)
	}
	for _, step := range phase.steps {
		if err := e.RunStep(step); err != nil {
			return fmt.Errorf("failed to run step %s: %w", step.name, err)
		}
	}
	if err := e.ApplyEvent(event.NewEndPhaseEvent(phase.name, e.game.ActivePlayerID())); err != nil {
		return fmt.Errorf("failed to end phase %s: %w", phase.name, err)
	}
	return nil
}

func (e *Engine) RunStep(step GameStep) error {
	e.log.Debug("Running step:", step.name)
	if err := e.ApplyEvent(event.NewBeginStepEvent(step.name, e.game.ActivePlayerID())); err != nil {
		return fmt.Errorf("failed to start step %s: %w", step.name, err)
	}
	for _, action := range step.actions {
		e.log.Debug("Completing action:", action.Name())
		// TODO: This should be a separate event, not an action.
		if err := e.CompleteTurnAction(action); err != nil {
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
	if err := e.ApplyEvent(event.NewEndStepEvent(step.name, e.game.ActivePlayerID())); err != nil {
		return fmt.Errorf("failed to end step %s: %w", step.name, err)
	}
	return nil
}

func (e *Engine) RunPriority() error {
	priorityPlayerID := e.game.ActivePlayerID()
	if err := e.ApplyEvent(
		event.ReceivePriorityEvent{PlayerID: priorityPlayerID},
	); err != nil {
		return fmt.Errorf("failed to apply receive priority event: %w", err)
	}
	for !e.game.AllPlayersPassedPriority() {
		priorityPlayerID = e.game.PriorityPlayerID()
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
			if errors.Is(err, ErrInvalidUserAction) {
				e.log.Debugf("Error completing action for player %q: %s", priorityPlayerID, err)
				continue
			}
			return fmt.Errorf(
				"failed to apply action for player %q: %w",
				priorityPlayerID,
				err,
			)
		}
		if e.game.PlayerPassedPriority(priorityPlayerID) {
			nextPlayerIDWithPriority := e.game.NextPlayerID(priorityPlayerID)
			if err := e.ApplyEvent(event.ReceivePriorityEvent{
				PlayerID: nextPlayerIDWithPriority,
			}); err != nil {
				return fmt.Errorf("failed to apply receive priority event: %w", err)
			}
		} else {
			if err := e.ApplyEvent(event.ResetPriorityPassesEvent{}); err != nil {
				return fmt.Errorf("failed to reset priority passes: %w", err)
			}
		}
		if e.game.AllPlayersPassedPriority() {
			if err := e.ApplyEvent(event.AllPlayersPassedPriorityEvent{}); err != nil {
				return fmt.Errorf("failed to apply all players passed priority event: %w", err)
			}
			if e.game.Stack().Size() == 0 {
				continue
			}
			// GetNextStackItem?
			resolvable, _, ok := e.game.Stack().TakeTop()
			if !ok {
				return fmt.Errorf("failed to take top from stack: %w", err)
			}
			if err := e.ResolveResolvable(resolvable); err != nil {
				return fmt.Errorf("failed to resolve resolvable: %w", err)
			}
			if err := e.ApplyEvent(event.ResetPriorityPassesEvent{}); err != nil {
				return fmt.Errorf("failed to reset priority passes: %w", err)
			}
		}
	}
	return nil
}

func (e *Engine) ResolveResolvable(resolvable state.Resolvable) error {
	player, ok := e.game.GetPlayer(resolvable.Controller())
	if !ok {
		return fmt.Errorf("player %q not found", resolvable.Controller())
	}
	events := []event.GameEvent{
		event.ResolveTopObjectOnStackEvent{
			Name: resolvable.Name(),
			ID:   resolvable.ID(),
		},
	}
	// TODO: This could probably be more elegent, instead of having the top level effect array be different than the chain it starts,
	// maybe it could all be one thing.
	// TODO: Rename to EffectSpecs()
	targets := resolvable.Targets()
	for _, effectSpec := range resolvable.Effects() {
		e.log.Debug("Resolving effect:", effectSpec.Name)
		efct, err := effect.Build(effectSpec)
		if err != nil {
			return fmt.Errorf("effect %q not found: %w", effectSpec.Name, err)
		}
		//effectResult := effect.EffectResult{}
		effectResult, err := efct.Resolve(e.game, player, resolvable, targets[effectSpec.Name])
		if err != nil {
			return fmt.Errorf("failed to apply effect %q: %w", effectSpec.Name, err)
		}
		// apply events as we go, instead of at the end.
		// later effects might depend on earlier ones.
		// TODO: This whole section feels clunky.
		for _, evnt := range effectResult.Events {
			if err := e.ApplyEvent(evnt); err != nil {
				return fmt.Errorf("failed to apply event %T: %w", evnt, err)
			}
		}
		// TODO: This needs to be a loop, right now we only handle a depth of 1
		if effectResult.ChoicePrompt.ChoiceOpts != nil {
			agent := e.agents[player.ID()]
			choiceResults, err := agent.Choose(effectResult.ChoicePrompt)
			if err != nil {
				return fmt.Errorf("failed to choose action for player %q: %w", player.ID(), err)
			}
			if effectResult.ResumeFunc == nil {
				return fmt.Errorf("effect %q requires choices but has no resume function", efct.Name())
			}
			effectResult, err = effectResult.ResumeFunc(choiceResults)
			if err != nil {
				return fmt.Errorf("failed to resume function for effect %s: %w", efct.Name(), err)
			}
			if effectResult.ChoicePrompt.ChoiceOpts != nil {
				panic("only one level of recursion supported")
			}
			events = append(events, effectResult.Events...)
		}
	}
	if spell, ok := resolvable.(gob.Spell); ok {
		if spell.Flashback() {
			events = append(events, event.PutSpellInExileEvent{
				PlayerID: spell.Owner(),
				SpellID:  resolvable.ID(),
			})
		} else {
			if spell.Match(is.Permanent()) {
				events = append(events, event.PutPermanentOnBattlefieldEvent{
					PlayerID: spell.Owner(),
					CardID:   spell.ID(),
					FromZone: mtg.ZoneStack,
				})
			} else {
				events = append(events, event.PutSpellInGraveyardEvent{
					PlayerID: spell.Owner(),
					SpellID:  resolvable.ID(),
				})
			}
		}
	}
	if ability, ok := resolvable.(gob.AbilityOnStack); ok {
		events = append(events, event.RemoveAbilityFromStackEvent{
			PlayerID:  ability.Owner(),
			AbilityID: ability.ID(),
		})
	}
	for _, evnt := range events {
		if err := e.ApplyEvent(evnt); err != nil {
			return fmt.Errorf("failed to apply event %T: %w", evnt, err)
		}
	}
	return nil
}

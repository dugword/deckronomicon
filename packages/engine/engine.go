package engine

// Engine manages the control flow of the game, including running turns, phases, and steps.

import (
	"deckronomicon/packages/configs"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/engine/resolver"
	"deckronomicon/packages/engine/rng"
	"deckronomicon/packages/engine/turnaction"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

type Action interface {
	Name() string
	Complete(state.Game, state.Player, *resenv.ResEnv) ([]event.GameEvent, error)
}

type Logger interface {
	Debug(...any)
	Debugf(format string, args ...any)
	Info(...any)
	Infof(format string, args ...any)
	Warn(...any)
	Warnf(format string, args ...any)
	Error(...any)
	Errorf(format string, args ...any)
	Critical(...any)
	Criticalf(format string, args ...any)
}

// TODO: I don't think this is the right place for this,
// not sure if the action interface above is in the right place either
var ErrInvalidUserAction = errors.New("invalid user action")

type Engine struct {
	agents      map[string]PlayerAgent
	deckLists   map[string]configs.DeckList
	game        state.Game
	record      *GameRecord
	log         Logger
	definitions map[string]definition.Card
	resEnv      *resenv.ResEnv
}

type EngineConfig struct {
	Players     []state.Player
	Agents      map[string]PlayerAgent
	Seed        int64
	DeckLists   map[string]configs.DeckList
	Definitions map[string]definition.Card
	Log         Logger
}

func NewEngine(config EngineConfig) *Engine {
	agents := map[string]PlayerAgent{}
	for id, agent := range config.Agents {
		agents[id] = agent
	}
	rng := rng.NewRNG(config.Seed)
	definitions := config.Definitions
	return &Engine{
		agents:      agents,
		deckLists:   config.DeckLists,
		game:        state.Game{}.WithPlayers(config.Players),
		log:         config.Log,
		record:      NewGameRecord(config.Seed),
		definitions: definitions,
		resEnv: &resenv.ResEnv{
			RNG:         rng,
			Definitions: definitions,
		},
	}
}

func (e *Engine) Game() state.Game {
	return e.game
}

func (e *Engine) Record() *GameRecord {
	return e.record
}

func (e *Engine) RunGame() error {
	// TODO This shouldn't live here. It should be in the Apply reducers and managed via events.
	// Or it needs to be created prior to the game starting and passed in.
	for _, playerID := range e.game.PlayerIDsInTurnOrder() {
		e.log.Debug("Setting up player deck:", playerID)
		deckList, ok := e.deckLists[playerID]
		if !ok {
			return fmt.Errorf("deck list for player %q not found", playerID)
		}
		game, deck, err := e.game.WithBuildDeck(
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
		e.game = game
		player := e.game.GetPlayer(playerID)
		newPlayer := player.WithLibrary(state.NewLibrary(deck))
		e.game = e.game.WithUpdatedPlayer(newPlayer)
	}

	e.log.Debug("Running game")
	if err := e.ApplyEvent(event.BeginGameEvent{}); err != nil {
		return fmt.Errorf("failed to start game: %w", err)
	}
	e.log.Debug("Shuffling decks")
	for _, playerID := range e.game.PlayerIDsInTurnOrder() {
		player := e.game.GetPlayer(playerID)
		var cardIDs []string
		for _, card := range player.Library().GetAll() {
			cardIDs = append(cardIDs, card.ID())
		}
		shuffledCardsIDs := e.resEnv.RNG.ShuffleIDs(cardIDs)
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
		if err := e.RunTurn(); err != nil {
			return fmt.Errorf("failed to run turn: %w", err)
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
		if err := e.CompleteTurnAction(action); err != nil {
			return fmt.Errorf(
				"failed to apply action %s: %w", action.Name(), err,
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
	for {
		if e.game.DidAllPlayersPassPriority() && e.game.Stack().Size() == 0 {
			return nil
		}
		priorityPlayerID := e.game.PriorityPlayerID()
		e.log.Debug("Running priority for player:", priorityPlayerID)
		if err := e.ApplyEvent(
			event.ReceivePriorityEvent{PlayerID: priorityPlayerID},
		); err != nil {
			return fmt.Errorf("failed to apply receive priority event: %w", err)
		}
		if err := e.RunPlayerActions(priorityPlayerID); err != nil {
			return fmt.Errorf(
				"failed to run player actions for %q: %w", priorityPlayerID, err,
			)
		}
		spellOrAbility, ok := e.game.Stack().GetTop()
		if !ok {
			continue
		}
		if err := e.ResolveSpellOrAbility(spellOrAbility); err != nil {
			return fmt.Errorf(
				"failed to resolve spell or ability %q: %w", spellOrAbility.Name(), err,
			)
		}
		if err := e.ApplyEvent(event.ResetPriorityPassesEvent{}); err != nil {
			return fmt.Errorf("failed to reset priority passes: %w", err)
		}
	}
}

func (e *Engine) RunPlayerActions(playerID string) error {
	for {
		if e.game.DidPlayerPassPriority(playerID) {
			return nil
		}
		e.log.Debugf("Running player actions for %q", playerID)
		action, err := e.agents[playerID].GetNextAction(e.game)
		if err != nil {
			return fmt.Errorf(
				"failed to get next action for player %q: %w", playerID, err,
			)
		}
		player := e.game.GetPlayer(playerID)
		evnts, err := action.Complete(e.game, player, e.resEnv)
		if err != nil {
			// TODO: Actually return this error in action.Complete
			// Right now I don't, and for now I'm going to ignore all errors...
			if errors.Is(err, ErrInvalidUserAction) {
				e.log.Debugf("Invalid player action for %q: %s", playerID, err)
				continue
			}
			e.log.Warnf("Hopefully this is just an invalid action, continue on everything for now: %s", err)
			e.log.Debugf("Invalid player action for %q: %s", playerID, err)
			continue
			/*
				return fmt.Errorf(
					"failed to complete action %q: %w", action.Name(), errors.Join(err, ErrInvalidUserAction),
				)
			*/
		}
		for _, evnt := range evnts {
			if err := e.ApplyEvent(evnt); err != nil {
				return fmt.Errorf(
					"failed to apply event %q: %w", evnt.EventType(), err,
				)
			}
		}
	}
}

func (e *Engine) ResolveSpellOrAbility(resolvable state.Resolvable) error {
	e.log.Debugf("Resolving spell or ability %q <%s> for %q", resolvable.Name(), resolvable.ID(), resolvable.Controller())
	e.ApplyEvent(event.ResolveTopObjectOnStackEvent{
		Name: resolvable.Name(),
		ID:   resolvable.ID(),
	})
	for _, effectWithTarget := range resolvable.EffectWithTargets() {
		player := e.game.GetPlayer(resolvable.Controller())
		e.log.Debugf("Resolving effect: %T", effectWithTarget.Effect)
		effectResult, err := resolver.Resolve(e.game, player.ID(), resolvable, effectWithTarget, e.resEnv)
		if err != nil {
			return fmt.Errorf("failed to resolve effect %q: %w", effectWithTarget.Effect.Name(), err)
		}
		if err := e.ResolveEffectResult(player.ID(), effectResult); err != nil {
			return fmt.Errorf("failed to resolve effect result for effect %q: %w", effectWithTarget.Effect.Name(), err)
		}
	}
	if resolvable.Match(query.And(is.Spell(), is.PermanentCardType())) {
		// TODO: Maybe permanents should have an effect that applies them to the battlefield
		// instead of this being a special case.
		if err := e.ApplyEvent(event.PutPermanentOnBattlefieldEvent{
			PlayerID: resolvable.Owner(),
			CardID:   resolvable.SourceID(),
			FromZone: mtg.ZoneStack,
		}); err != nil {
			return fmt.Errorf("failed to apply event PutPermanentOnBattlefieldEvent: %w", err)
		}
	} else {
		if err := e.ApplyEvent(event.RemoveSpellOrAbilityFromStackEvent{
			PlayerID: resolvable.Owner(),
			ObjectID: resolvable.ID(),
		}); err != nil {
			return fmt.Errorf("failed to apply event RemoveSpellOrAbilityFromStackEvent: %w", err)
		}
		// TODO: I don't like this, I moved having the reducer manage putting flashback
		// cards in the exile zone because I wanted it to be managed in one spot, this
		// feels like it should be managed there too, but I need the rng from engine...
		// Dunno, think about this more.
		if resolvable.Match(query.And(is.Spell(), has.Subtype(mtg.SubtypeOmen))) {
			var cardIDs []string
			player := e.game.GetPlayer(resolvable.Owner())
			for _, card := range player.Library().GetAll() {
				cardIDs = append(cardIDs, card.ID())
			}
			shuffledCardsIDs := e.resEnv.RNG.ShuffleIDs(cardIDs)
			if err := e.ApplyEvent(event.ShuffleLibraryEvent{
				PlayerID:         player.ID(),
				ShuffledCardsIDs: shuffledCardsIDs,
			}); err != nil {
				return fmt.Errorf("failed to apply event ShuffleLibraryEvent for omen: %w", err)
			}
		}
	}
	return nil
}

func (e *Engine) ResolveEffectResult(
	playerID string,
	result resolver.Result,
) error {
	agent := e.agents[playerID]
	for {
		e.log.Debugf("Resolving effect result for player %q", playerID)
		for _, evnt := range result.Events {
			if err := e.ApplyEvent(evnt); err != nil {
				return fmt.Errorf("failed to apply event %T: %w", evnt, err)
			}
		}
		if result.ChoicePrompt.ChoiceOpts == nil {
			return nil
		}
		choiceResults, err := agent.Choose(result.ChoicePrompt)
		if err != nil {
			return fmt.Errorf("failed to get choice from player agent %q: %w", agent.PlayerID(), err)
		}
		if result.Resume == nil {
			return fmt.Errorf("missing resume function")
		}
		result, err = result.Resume(choiceResults)
		if err != nil {
			return fmt.Errorf("failed to resume effect result: %w", err)
		}
	}
}

package store

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resolver"
	"deckronomicon/packages/engine/store/reducer"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/state"
	"fmt"
	"slices"
)

type ApplyEventFunc func(game *state.Game, event event.GameEvent) (*state.Game, error)

type Middleware func(next ApplyEventFunc) ApplyEventFunc

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

type Store struct {
	game            *state.Game
	middlewareChain ApplyEventFunc
}

func (s *Store) DONOTUSESetGameState(newGameState *state.Game) {
	s.game = newGameState
}

func NewStore(initialState *state.Game, middlewareBuilder func(store *Store) []Middleware) *Store {
	// Copy the initial state to avoid modifying the original state
	copiedState := *initialState
	store := Store{
		game: &copiedState,
	}
	// middleware := []Middleware{TriggeredAbilitiesMiddleware(&store)}
	middleware := []Middleware{}
	if middlewareBuilder != nil {
		middleware = append(middleware, middlewareBuilder(&store)...)
	}
	store.middlewareChain = ChainMiddleware(CoreApplyEvent, middleware...)
	return &store
}

func (s *Store) Apply(evnt event.GameEvent) error {
	eventQueue := []event.GameEvent{evnt}
	for len(eventQueue) > 0 {
		currentEvent := eventQueue[0]
		eventQueue = eventQueue[1:]
		fmt.Println("Current Event:", currentEvent.EventType())

		newGame, err := s.middlewareChain(s.game, currentEvent)
		if err != nil {
			return fmt.Errorf("failed to apply event %q: %w", currentEvent.EventType(), err)
		}

		triggerEvents := GenerateTriggerEvents(s.game, newGame, currentEvent)
		eventQueue = append(eventQueue, triggerEvents...)

		s.game = newGame
		triggeredAbilities, err := CheckTriggeredAbilities(s.game, currentEvent)
		if err != nil {
			return fmt.Errorf("failed to check triggered abilities: %w", err)
		}
		for _, triggeredAbility := range triggeredAbilities {
			triggeredEvents, err := HandleTriggeredAbility(s.game, triggeredAbility, currentEvent)
			if err != nil {
				return fmt.Errorf("failed to handle triggered ability: %w", err)
			}
			eventQueue = append(eventQueue, triggeredEvents...)
		}
		// TODO: Make this more elegant and generic
		if currentEvent.EventType() == "BeginEndStep" {
			// Remove all triggered abilities that are set to end at the end of the turn.
			for _, ta := range s.game.RegisteredTriggeredAbilities() {
				if ta.Duration == mtg.DurationEndOfTurn {
					eventQueue = append(eventQueue, &event.RemoveTriggeredAbilityEvent{
						ID: ta.ID,
					})
				}
			}
		}
	}
	return nil
}

func (s *Store) Game() *state.Game {
	// TODO: Consider returning a copy or immutable version of the state
	return s.game
}

func CoreApplyEvent(game *state.Game, ev event.GameEvent) (*state.Game, error) {
	game, err := reducer.Apply(game, ev)
	if err != nil {
		return nil, fmt.Errorf("failed to apply event %q: %w", ev.EventType(), err)
	}
	return game, nil
}

func ChainMiddleware(final ApplyEventFunc, middleware ...Middleware) ApplyEventFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		final = middleware[i](final)
	}
	return final
}

func CheckTriggeredAbilities(game *state.Game, evnt event.GameEvent) ([]gob.RegisteredTriggeredAbility, error) {
	var triggeredAbilities []gob.RegisteredTriggeredAbility
	for _, ta := range game.RegisteredTriggeredAbilities() {
		if MatchesTrigger(ta, evnt, game, ta.PlayerID) {
			triggeredAbilities = append(triggeredAbilities, ta)
		}
	}
	return triggeredAbilities, nil
}

func HandleTriggeredAbility(game *state.Game, triggeredAbility gob.RegisteredTriggeredAbility, evnt event.GameEvent) ([]event.GameEvent, error) {
	var events []event.GameEvent
	var effectWithTargets []*effect.EffectWithTarget
	for _, efct := range triggeredAbility.Effects {
		if addManaEffect, ok := efct.(*effect.AddMana); ok {
			result, err := resolver.ResolveAddMana(game, triggeredAbility.PlayerID, addManaEffect)
			if err != nil {
				return nil, fmt.Errorf("failed to apply effect %q: %w", "AddMana", err)
			}
			events = append(events, result.Events...)
		} else {
			effectWithTargets = append(effectWithTargets, &effect.EffectWithTarget{
				Effect: efct,
				Target: target.Target{
					Type: mtg.TargetTypeNone,
				},
			})
		}
	}
	if len(effectWithTargets) > 0 {
		events = append(events, &event.PutAbilityOnStackEvent{
			PlayerID:          triggeredAbility.PlayerID,
			SourceID:          triggeredAbility.SourceID,
			AbilityID:         triggeredAbility.ID,
			AbilityName:       "Triggered Effect",
			EffectWithTargets: effectWithTargets,
		})
	}
	return events, nil
}

func MatchesTrigger(triggeredAbility gob.RegisteredTriggeredAbility, evnt event.GameEvent, game *state.Game, playerID string) bool {
	// TODO: This match logic should live in the trigger itself I think, otherwise this is going to get out of hand.
	// Or maybe not because we have a generic "filter" in the trigger that is applied differently based on the event type.
	// Maybe this needs to be applied in a dispatching reducer pattern like the apply events function.
	// Maybe this should be in the judge package.
	// TODO: Yeah probably should be in the judge package.
	fmt.Printf("Triggered Ability => %+v\n", triggeredAbility)
	fmt.Printf("Event => %s ::  %+v\n", evnt.EventType(), evnt)
	if triggeredAbility.Trigger.EventType != evnt.EventType() {
		return false
	}
	switch e := evnt.(type) {
	case *event.DeathEvent:
		if triggeredAbility.Trigger.SelfTrigger {
			if triggeredAbility.SourceID == e.CardID {
				return true
			}
		}
		return false
	case *event.EnteredTheBattlefieldEvent:
		fmt.Printf("Triggered Ability => %+v\n", triggeredAbility)
		if triggeredAbility.Trigger.SelfTrigger {
			if triggeredAbility.SourceID == e.PermanentID {
				return true
			}
		}
		return false
	case *event.LandTappedForManaEvent:
		if e.PlayerID != playerID {
			return false
		}
		if triggeredAbility.Trigger.Filter.Subtypes != nil {
			for _, subtype := range triggeredAbility.Trigger.Filter.Subtypes {
				if !slices.Contains(e.Subtypes, subtype) {
					return false
				}
			}
		}
		return true
	case *event.BeginEndStepEvent:
		if e.PlayerID != playerID {
			return false
		}
		return true
	}
	return false
}

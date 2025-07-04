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
		game:            &copiedState,
		middlewareChain: ChainMiddleware(CoreApplyEvent),
	}
	if middlewareBuilder != nil {
		store.middlewareChain = ChainMiddleware(CoreApplyEvent, middlewareBuilder(&store)...)
	}
	return &store
}

func (s *Store) Apply(evnt event.GameEvent) error {
	newGame, err := s.middlewareChain(s.game, evnt)
	if err != nil {
		return fmt.Errorf("failed to apply event %q: %w", evnt.EventType(), err)
	}
	s.game = newGame
	triggeredAbilities, err := CheckTriggeredAbilities(s.game, evnt)
	if err != nil {
		return fmt.Errorf("failed to check triggered abilities: %w", err)
	}
	for _, triggeredAbility := range triggeredAbilities {
		triggeredEvents, err := HandleTriggeredAbility(s.game, triggeredAbility, evnt)
		if err != nil {
			return fmt.Errorf("failed to handle triggered ability: %w", err)
		}
		for _, triggeredEvent := range triggeredEvents {
			err := s.Apply(triggeredEvent)
			if err != nil {
				return fmt.Errorf("failed to apply triggered event %q: %w", triggeredEvent.EventType(), err)
			}
		}
	}
	// TODO: Make this more elegant and generic
	if evnt.EventType() == "BeginEndStep" {
		// Remove all triggered abilities that are set to end at the end of the turn.
		for _, ta := range s.game.TriggeredAbilities() {
			if ta.Duration == mtg.DurationEndOfTurn {
				if err := s.Apply(&event.RemoveTriggeredAbilityEvent{
					ID: ta.ID,
				}); err != nil {
					return fmt.Errorf("failed to remove triggered ability %q: %w", ta.ID, err)
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

func CheckTriggeredAbilities(game *state.Game, evnt event.GameEvent) ([]gob.TriggeredAbility, error) {
	var triggeredAbilities []gob.TriggeredAbility
	for _, ta := range game.TriggeredAbilities() {
		if MatchesTrigger(ta.Trigger, evnt, game, ta.PlayerID) {
			triggeredAbilities = append(triggeredAbilities, ta)
		}
	}
	return triggeredAbilities, nil
}

func HandleTriggeredAbility(game *state.Game, triggeredAbility gob.TriggeredAbility, evnt event.GameEvent) ([]event.GameEvent, error) {
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

func MatchesTrigger(trigger gob.Trigger, evnt event.GameEvent, game *state.Game, playerID string) bool {
	// TODO: This match logic should live in the trigger itself I think, otherwise this is going to get out of hand.
	// Or maybe not because we have a generic "filter" in the trigger that is applied differently based on the event type.
	// Maybe this needs to be applied in a dispatching reducer pattern like the apply events function.
	// Maybe this should be in the judge package.
	// TODO: Yeah probably should be in the judge package.
	switch trigger.EventType {
	case "LandTappedForMana":
		LandTappedForManaEvent, ok := evnt.(*event.LandTappedForManaEvent)
		if !ok {
			return false
		}
		if LandTappedForManaEvent.PlayerID != playerID {
			return false
		}
		if trigger.Filter.Subtypes != nil {
			for _, subtype := range trigger.Filter.Subtypes {
				if !slices.Contains(LandTappedForManaEvent.Subtypes, subtype) {
					return false
				}
			}
		}
		return true
	case "BeginEndStep":
		BeginEndStepEvent, ok := evnt.(*event.BeginEndStepEvent)
		if !ok {
			return false
		}
		if BeginEndStepEvent.PlayerID != playerID {
			return false
		}
		return true
	}
	return false
}

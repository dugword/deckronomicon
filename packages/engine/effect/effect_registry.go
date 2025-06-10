package effect

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
)

type EffectHandler func(
	game state.Game,
	player state.Player,
	source query.Object,
	modifiers []gob.Tag,
) ([]event.GameEvent, error)

type EffectRegistry struct {
	handlers map[string]EffectHandler
}

func NewEffectRegistry() *EffectRegistry {
	registry := EffectRegistry{
		handlers: map[string]EffectHandler{},
	}
	// TODO: This should probably be moved elsewhere
	registry.RegisterEffect("AddMana", AddManaEffectHandler)
	return &registry
}

func (r *EffectRegistry) RegisterEffect(name string, handler EffectHandler) {
	if _, exists := r.handlers[name]; exists {
		panic("Effect already registered: " + name)
	}
	r.handlers[name] = handler
}

func (r *EffectRegistry) Get(name string) (EffectHandler, bool) {
	handler, exists := r.handlers[name]
	if !exists {
		return nil, false
	}
	return handler, true
}

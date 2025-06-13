package effect

type EffectRegistry struct {
	handlers map[string]EffectHandler
}

// TODO: Maybe effect modifiers should live here, they can be interfaces in the definitions
// and then they can be cast here so they can be managed with types specifically for effects

func NewEffectRegistry() *EffectRegistry {
	registry := EffectRegistry{
		handlers: map[string]EffectHandler{},
	}
	// TODO: This should probably be moved elsewhere
	registry.RegisterEffect("AddMana", AddManaEffectHandler)
	registry.RegisterEffect("Draw", DrawEffectHandler)
	registry.RegisterEffect("Counterspell", CounterspellEffectHandler)
	registry.RegisterEffect("Scry", ScryEffectHandler)
	registry.RegisterEffect("PutBackOnTop", PutBackOnTopHandler)
	registry.RegisterEffect("Search", SearchEffectHandler)
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

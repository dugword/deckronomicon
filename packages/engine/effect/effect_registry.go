package effect

type EffectRegistry struct {
	handlers map[string]EffectHandler
}

func NewEffectRegistry() *EffectRegistry {
	registry := EffectRegistry{
		handlers: map[string]EffectHandler{},
	}
	// TODO: This should probably be moved elsewhere
	registry.RegisterEffect("AddMana", AddManaEffectHandler)
	registry.RegisterEffect("Draw", DrawEffectHandler)
	registry.RegisterEffect("Typecycling", TypecyclingEffectHandler)
	registry.RegisterEffect("Scry", ScryEffectHandler)
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

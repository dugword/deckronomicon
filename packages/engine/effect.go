package engine

type EffectStep struct {
	Prompt  *string
	Resolve func([]string) error
}

type EffectChain []EffectStep

type Effect interface {
	Name() string
	GetResolutionSteps() (EffectChain, error)
}

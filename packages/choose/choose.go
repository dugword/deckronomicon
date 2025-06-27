package choose

type ChoiceType string

type Choice interface {
	Name() string
	ID() string
}

type Source interface {
	Name() string
}

type ChoicePrompt struct {
	ChoiceOpts ChoiceOpts
	Message    string
	Source     Source
	Optional   bool
}

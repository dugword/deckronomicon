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
}

// TODO: Use an interface for ChoicePrompt, so it can pass a chooseOne, chooseN, etc.

type ChoicePromptOld struct {
	Choices    []Choice
	MaxChoices int
	Message    string
	MinChoices int
	Optional   bool
	Source     Source
}

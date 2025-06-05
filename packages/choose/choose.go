package choose

type Source interface {
	Name() string
}

type Choice struct {
	Name string
	ID   string
}

type ChoicePrompt struct {
	Choices    []Choice
	MaxChoices int
	Message    string
	MinChoices int
	Optional   bool
	Source     Source
}

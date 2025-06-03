package choose

type Choice struct {
}

type ChoicePrompt struct {
	Message  string
	Choices  []Choice
	Optional bool
}

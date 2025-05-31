package choose

type GenericChoice struct {
	name string
	id   string
}

func NewGenericChoice(name, id string) GenericChoice {
	return GenericChoice{
		name: name,
		id:   id,
	}
}

func (c GenericChoice) Name() string {
	return c.name
}

func (c GenericChoice) ID() string {
	return c.id
}

func NewChoices[T Choice](items []T) []Choice {
	var choices []Choice
	for _, item := range items {
		choices = append(choices, item)
	}
	return choices
}

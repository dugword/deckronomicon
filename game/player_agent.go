package game

type ChoiceSource interface {
	Name() string
	ID() string
}

// TODO: Don't have a generic choice source, have a specific one for each type
// of choice, e.g. MenuChoiceSource....
type choiceSource struct {
	id   string
	name string
}

func (c *choiceSource) Name() string {
	return c.name
}

func (c *choiceSource) ID() string {
	return c.id
}

func NewChoiceSource(name, id string) ChoiceSource {
	choiceSource := choiceSource{
		id:   id,
		name: name,
	}
	return &choiceSource
}

// CreateObjectChoices generates a list of choices from a slice of GameObjects.
func CreateObjectChoices(objects []GameObject, zone string) []Choice {
	var choices []Choice
	for _, object := range objects {
		choices = append(choices, Choice{
			ID:   object.ID(),
			Name: object.Name(),
			Zone: zone,
		})
	}
	return choices
}

// CreateActivatedAbilityChoices generates a list of choices from a slice of
// ActivatedAbilities.
func CreateActivatedAbilityChoices(abilities []*ActivatedAbility) []Choice {
	var choices []Choice
	for _, ability := range abilities {
		choices = append(choices, Choice{
			ID:     ability.ID(),
			Name:   ability.Description(),
			Source: ability.source.Name(),
			Zone:   ability.Zone,
		})
	}
	return choices
}

// PlayerAgent defines how player decisions are made.
type PlayerAgent interface {
	// TODO: Not sure I love this here
	//ChooseAny(prompt string, choices []Choice) []Choice
	//ChooseN(prompt string, choices []Choice, n int) []Choice
	//ChooseUpToN(prompt string, choices []Choice, n int) []Choice
	ChooseOne(prompt string, source ChoiceSource, choices []Choice) (Choice, error)
	Confirm(prompt string, source ChoiceSource) (bool, error)
	EnterNumber(prompt string, source ChoiceSource) (int, error)
	GetNextAction(state *GameState) *GameAction
	PlayerID() string
	ReportState(state *GameState)
}

const ChoiceNone = "None"

// OptionalChoice returns adds an optional choice to the list of choices.
func AddOptionalChoice(choices []Choice) []Choice {
	choices = append([]Choice{{
		// TODO: Make this a constant, maybe special character to prevent
		// collision with other IDs
		ID:     ChoiceNone,
		Name:   ChoiceNone,
		Source: ChoiceNone,
		Zone:   ChoiceNone,
	}}, choices...)
	return choices
}

// Choice represents a choice made by the player.
type Choice struct {
	ID     string
	Name   string
	Source string
	Zone   string
}

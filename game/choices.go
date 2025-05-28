package game

import "sort"

const ChoiceNone = "None"

// Choice represents a choice made by the player.
type Choice struct {
	ID     string
	Name   string
	Source string
	Zone   string
}

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

type actionSource struct {
	action string
}

func (c *actionSource) Name() string {
	return c.action
}

func (c *actionSource) ID() string {
	return c.action
}

type gameObjectSource struct {
	source GameObject
}

func (c *gameObjectSource) Name() string {
	return c.source.Name()
}

func (c *gameObjectSource) ID() string {
	return c.source.ID()
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

// CreateChoices generates a list of choices from a slice of GameObjects.
func CreateChoices(objects []GameObject, zone string) []Choice {
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

func CreateGroupedChoices(gourpedObjects map[string][]GameObject) []Choice {
	var choices []Choice
	var groupNames []string
	for groupName := range gourpedObjects {
		groupNames = append(groupNames, groupName)
	}
	sort.Strings(groupNames)
	for _, groupName := range groupNames {
		objects := gourpedObjects[groupName]
		for _, object := range objects {
			choices = append(choices, Choice{
				ID:   object.ID(),
				Name: object.Name(),
				// Source: groupName,
				Zone: groupName,
			})
		}
	}
	return choices
}

// OptionalChoice returns adds an optional choice to the list of choices.
func AddOptionalChoice(choices []Choice) []Choice {
	choices = append([]Choice{{
		// TODO: Make this a constant, maybe special character to prevent
		// collision with other IDs
		ID:   ChoiceNone,
		Name: ChoiceNone,
	}}, choices...)
	return choices
}

/*
func GroupedChoicesData(title string, choices []game.Choice) (BoxData, []game.Choice) {
	grouped := make(map[string][]game.Choice)
	for _, choice := range choices {
		if choice.ID == game.ChoiceNone {
			grouped[game.ChoiceNone] = append(grouped[game.ChoiceNone], choice)
		} else if choice.Zone == "" {
			// TODO: Handle this different
			grouped[""] = append(grouped[""], choice)
		} else {
			grouped[choice.Zone] = append(grouped[choice.Zone], choice)
		}
	}
	var groupNames []string
	for groupName := range grouped {
		if groupName == game.ChoiceNone {
			continue
		}
		groupNames = append(groupNames, groupName)
	}
	sort.Strings(groupNames)
	// TODO: Make 0 always none, and have it be not an option if the choice is
	// not optional
	var orderedChoices []game.Choice
	var i = 0
	var lines []string
	for _, groupName := range groupNames {
		lines = append(lines, fmt.Sprintf("--- %s ---", groupName))
		for _, choice := range grouped[groupName] {
			var line string
			if choice.Source == "" {
				line = fmt.Sprintf("%d: %s", i, choice.Name)
			} else {
				line = fmt.Sprintf("%d: %s - %s", i, choice.Source, choice.Name)
			}
			lines = append(lines, line)
			orderedChoices = append(orderedChoices, choice)
			i++
		}
	}
	if grouped[game.ChoiceNone] != nil {
		lines = append(lines, "---------")
		lines = append(lines, fmt.Sprintf("%d: %s", i, game.ChoiceNone))
		orderedChoices = append(orderedChoices, game.Choice{ID: game.ChoiceNone})
	}
	return BoxData{
		Title:   title,
		Content: lines,
	}, orderedChoices
}
*/

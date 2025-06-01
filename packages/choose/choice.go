package choose

import (
	"deckronomicon/packages/query"
	"errors"
	"sort"
)

var ErrNoChoices = errors.New("no choices available")

// Choice represents a choice made by the player.
type Choice struct {
	ID          string
	Name        string
	Description string // Prompt
	Source      Source
}

type Choosable interface {
	ID() string
	Name() string
	Description() string
}

type Source interface {
	Name() string
}

type choiceSource struct {
	name string
}

func (c *choiceSource) Name() string {
	return c.name
}

func NewChoiceSource(name string) Source {
	return &choiceSource{name: name}
}

var ChoiceNone = Choice{
	ID:          "none",
	Name:        "None",
	Description: "No choice made",
}

// OptionalChoice returns adds an optional choice to the list of choices.
func AddOptionalChoice(choices []Choice) []Choice {
	choices = append([]Choice{ChoiceNone}, choices...)
	return choices
}

func CreateChoices[T Choosable](objects []T, source Source) []Choice {
	var choices []Choice
	for _, object := range objects {
		choice := Choice{
			ID:     object.ID(),
			Name:   object.Name(),
			Source: source,
		}
		choices = append(choices, choice)
	}
	return choices
}

// This was a good idea, but it needs to be reworked after the refactor
func CreateGroupedChoices(gourpedObjects map[string][]query.Object) []Choice {
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
				ID:     object.ID(),
				Name:   object.Name(),
				Source: NewChoiceSource(groupName),
			})
		}
	}
	return choices
}

/*
func By[T query.Object](objects []T, predicate query.Predicate) []T {
	var result []T
	for _, object := range objects {
		if predicate(object) {
			result = append(result, object)
		}
	}
	return result
}
*/

/*
 */

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

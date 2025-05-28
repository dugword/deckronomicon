package interactive

import (
	"deckronomicon/game"
	"fmt"
)

var choiceBoxTmpl = `{{.TopLine}}
{{.Title}}
{{.MiddleLine}}
{{range .Lines}}{{.}}
{{end}}{{.BottomLine}}
`

// Confirm prompts the user to confirm an action.
func (a *InteractivePlayerAgent) Confirm(prompt string, source game.ChoiceSource) (bool, error) {
	for {
		a.Prompt(prompt)
		accept, err := a.ReadInputConfirm()
		if err != nil {
			fmt.Println("Invalid choice. Please enter '(y)es' or '(n)o'")
			continue
		}
		return accept, nil
	}
}

// EnterNumber prompts the user to enter a number.
func (a *InteractivePlayerAgent) EnterNumber(prompt string, source game.ChoiceSource) (int, error) {
	for {
		a.Prompt(prompt)
		number, err := a.ReadInputNumber(-1)
		if err != nil {
			fmt.Println("Invalid choice. Please enter a number")
			continue
		}
		return number, nil
	}
}

// ChooseMany prompts the user to choose many of the given choices.
func (a *InteractivePlayerAgent) ChooseMany(prompt string, source game.ChoiceSource, choices []game.Choice) ([]game.Choice, error) {
	if len(choices) == 0 {
		return nil, fmt.Errorf("no choices available")
	}
	title := fmt.Sprintf("%s requires a choice", source.Name())
	a.uiBuffer.UpdateChoices(title, choices)
	if err := a.uiBuffer.Render(); err != nil {
		return nil, fmt.Errorf("error rendering UI buffer: %w", err)
	}
	for {
		a.Prompt(prompt)
		max := len(choices) - 1 // 0 based
		selected, err := a.ReadNumberMany(max)
		if err != nil {
			fmt.Printf("Invalid choice. Please enter numbers: %d - %d\n", 0, max)
			continue
		}
		if len(selected) == 0 {
			fmt.Println("You must choose at least one option.")
			continue
		}
		var selectedChoices []game.Choice
		for _, choice := range selected {
			selectedChoices = append(selectedChoices, choices[choice])
		}
		return selectedChoices, nil
	}
}

// ChoseOne prompts the user to choose one of the given choices.
// TODO: Need to enable a way to cancel
func (a *InteractivePlayerAgent) ChooseOne(prompt string, source game.ChoiceSource, choices []game.Choice) (game.Choice, error) {
	if len(choices) == 0 {
		return game.Choice{}, fmt.Errorf("no choices available")
	}
	title := fmt.Sprintf("%s requires a choice", source.Name())
	a.uiBuffer.UpdateChoices(title, choices)
	if err := a.uiBuffer.Render(); err != nil {
		return game.Choice{}, fmt.Errorf("error rendering UI buffer: %w", err)
	}

	for {
		a.Prompt(prompt)
		max := len(choices) - 1 // 0 based
		choice, err := a.ReadInputNumber(max)
		if err != nil {
			fmt.Printf("Invalid choice. Please enter a number: %d - %d\n", 0, max)
			continue
		}
		return choices[choice], nil
	}
}

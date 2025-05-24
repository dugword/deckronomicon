package interactive

import (
	"deckronomicon/game"
	"fmt"
	"os"
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

// ChoseOne prompts the user to choose one of the given choices.
// TODO: Need to enable a way to cancel
func (a *InteractivePlayerAgent) ChooseOne(prompt string, source game.ChoiceSource, choices []game.Choice) (game.Choice, error) {
	if len(choices) == 0 {
		return game.Choice{}, fmt.Errorf("no choices available")
	}
	title := fmt.Sprintf("%s requires a choice", source.Name())
	groupedChoicesData, orderedChoices := GroupedChoicesData(title, choices)
	a.UpdateChoiceData(groupedChoicesData)
	displayBoxes := a.BuildDisplayBoxes()
	ClearScreen()
	if err := a.DisplayTemplate.ExecuteTemplate(
		// TODO: use passed in stdout from Run
		os.Stdout,
		"display.tmpl",
		displayBoxes,
	); err != nil {
		fmt.Println("Error executing template:", err)
		// TODO: handle error return to main
		os.Exit(1)
	}
	for {
		a.Prompt(prompt)
		max := len(orderedChoices) - 1 // 0 based
		choice, err := a.ReadInputNumber(max)
		if err != nil {
			fmt.Printf("Invalid choice. Please enter a number: %d - %d\n", 0, max)
			continue
		}
		return orderedChoices[choice], nil
	}
}

package interactive

import (
	"deckronomicon/game"
	"fmt"
	"os"
	"text/template"
)

var choiceBoxTmpl = `{{.TopLine}}
{{.Title}}
{{.MiddleLine}}
{{range .Lines}}{{.}}
{{end}}{{.BottomLine}}
`

// ChoseOne prompts the user to choose one of the given choices.
// TODO: Need to enable a way to cancel
func (a *InteractivePlayerAgent) ChooseOne(prompt string, source game.ChoiceSource, choices []game.Choice) (game.Choice, error) {
	if len(choices) == 0 {
		return game.Choice{}, fmt.Errorf("no choices available")
	}
	title := fmt.Sprintf("%s requires a choice", source.Name())
	tmpl := template.Must(template.New("choices").Parse(choiceBoxTmpl))
	if err := tmpl.ExecuteTemplate(
		// TODO: use passed in stdout from Run
		os.Stdout,
		"choices",
		ChoicesBox(title, choices),
	); err != nil {
		fmt.Println("Error executing template:", err)
		// TODO: handle error return to main
		os.Exit(1)
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

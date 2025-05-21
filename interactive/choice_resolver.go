package interactive

import (
	"deckronomicon/game"
	"fmt"
)

// ChoseOne prompts the user to choose one of the given choices.
// TODO: Need to enable a way to cancel
func (a *InteractivePlayerAgent) ChooseOne(prompt string, source string, choices []game.Choice) (game.Choice, error) {
	if len(choices) == 0 {
		return game.Choice{}, fmt.Errorf("no choices available")
	}
	fmt.Println("SOURCE: ", source)
	for i, opt := range choices {
		fmt.Printf("%d: %s\n", i, opt.Name)
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

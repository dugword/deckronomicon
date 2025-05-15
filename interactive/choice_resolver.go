package interactive

import (
	"deckronomicon/game"
	"fmt"
)

func (a *InteractivePlayerAgent) ChooseOne(prompt string, choices []game.Choice) game.Choice {
	for i, opt := range choices {
		fmt.Printf("%d: %s\n", i, opt.Name)
	}
	for {
		a.Prompt(prompt)
		choice, err := a.ReadInputNumber(len(choices))
		if err != nil {
			fmt.Printf("Invalid choice. Please enter a number: %d - %d\n", 0, len(choices))
			continue
		}
		return choices[choice]
	}
}

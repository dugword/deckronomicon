package interactive

import (
	"deckronomicon/game"
	"fmt"
)

func (a *InteractivePlayerAgent) ChooseOne(prompt string, options []game.Choice) game.Choice {
	for i, opt := range options {
		fmt.Printf("%d: %s\n", i, opt.Name)
	}
	for {
		a.Prompt(prompt)
		choice, err := a.ReadInputNumber(len(options))
		if err != nil {
			fmt.Printf("Invalid choice. Please enter a number: %d - %d\n", 0, len(options))
			continue
		}
		return options[choice]
	}
}

package interactive

import (
	"deckronomicon/packages/choose"
	"fmt"
)

// Confirm prompts the user to confirm an action.
func (a *Agent) Confirm(prompt string, source choose.Source) (bool, error) {
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
func (a *Agent) EnterNumber(prompt string, source choose.Source) (int, error) {
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
func (a *Agent) ChooseMany(prompt string, source choose.Source, choices []choose.Choice) ([]choose.Choice, error) {
	if len(choices) == 0 {
		return nil, choose.ErrNoChoices
	}
	title := fmt.Sprintf("%s requires a choice", source.Name())
	a.uiBuffer.UpdateChoices(title, choices)
	if err := a.uiBuffer.Render(); err != nil {
		return nil, fmt.Errorf("failed to render UI Buffer: %w", err)
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
		var selectedChoices []choose.Choice
		for _, choice := range selected {
			selectedChoices = append(selectedChoices, choices[choice])
		}
		return selectedChoices, nil
	}
}

func (a *Agent) ChooseOne(prompt choose.ChoicePrompt) (choose.Choice, error) {
	prompt.MinChoices = 1
	prompt.MaxChoices = 1
	choices, err := a.Choose(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to choose one: %w", err)
	}
	if len(choices) == 0 {
		// TODO is this an error?
		return nil, nil
	}
	return choices[0], nil
}

func (a *Agent) Choose(prompt choose.ChoicePrompt) ([]choose.Choice, error) {
	if len(prompt.Choices) == 0 {
		// TODO is this an error?
		return []choose.Choice{}, choose.ErrNoChoices
	}
	title := fmt.Sprintf("%s requires a choice", prompt.Source.Name())
	if prompt.MinChoices > 1 {
		if prompt.MinChoices == prompt.MaxChoices {
			title = fmt.Sprintf(
				"%s requires %d choices",
				prompt.Source.Name(),
				prompt.MinChoices,
			)
		} else {
			title = fmt.Sprintf(
				"%s requires %d - %d choices",
				prompt.Source.Name(),
				prompt.MinChoices,
				prompt.MaxChoices,
			)
		}
	}
	var selected []choose.Choice
	var message []string
	if prompt.Optional {
		message = append(message, "Optional: Choose 0 to skip")
	}
	if prompt.MinChoices != prompt.MaxChoices {
		message = append(message, "Choose 0 to complete selection")
	}
	for len(selected) < prompt.MinChoices || len(selected) > prompt.MaxChoices {
		a.uiBuffer.UpdateMessage(message)
		a.uiBuffer.UpdateChoices(title, prompt.Choices)
		if err := a.uiBuffer.Render(); err != nil {
			return []choose.Choice{}, fmt.Errorf("failed to render UI Buffer: %w", err)
		}
		a.Prompt(prompt.Message)
		max := len(prompt.Choices)
		choicePlus1, err := a.ReadInputNumber(max)
		if err != nil {
			fmt.Printf("Invalid choice. Please enter a number: %d - %d\n", 0, max)
			continue
		}
		if choicePlus1 == 0 {
			if prompt.Optional || len(selected) > prompt.MinChoices {
				break
			} else {
				continue
			}
		}
		choice := choicePlus1 - 1 // Convert to 0 based index
		selected = append(selected, prompt.Choices[choice])
		message = append(message, fmt.Sprintf(
			"Selected: %s <id:%s>",
			prompt.Choices[choice].Name(),
			prompt.Choices[choice].ID(),
		))
		// TODO Should this follow the same immutable stuff
		prompt.Choices = append(prompt.Choices[:choice], prompt.Choices[choice+1:]...)
	}
	return selected, nil
}

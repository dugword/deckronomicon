package interactive

import (
	"deckronomicon/packages/choose"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// TODO I Don't like this here, or at all really...
func parseIndexes(line string) []int {
	var indexes []int
	parts := strings.Fields(line)
	for _, part := range parts {
		var idx int
		_, err := fmt.Sscanf(part, "%d", &idx)
		if err == nil {
			indexes = append(indexes, idx)
		}
	}
	return indexes
}

// Confirm prompts the user to confirm an action.
func (a *Agent) confirm(prompt string, source choose.Source) (bool, error) {
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
func (a *Agent) enterNumber(prompt string, source choose.Source) (int, error) {
	for {
		a.Prompt(prompt)
		input := a.ReadInput()
		number, err := a.inputToNumber(input, -1)
		if err != nil {
			fmt.Println("Invalid choice. Please enter a number")
			continue
		}
		return number, nil
	}
}

/*
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
*/

/*
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
*/

func (a *Agent) Choose(prompt choose.ChoicePrompt) (choose.ChoiceResults, error) {
	switch opts := prompt.ChoiceOpts.(type) {
	case choose.ChooseOneOpts:
		choices, err := a.chooseMany(opts.Choices, 1, 1, prompt.Message, prompt.Source, prompt.Optional)
		if errors.Is(err, choose.ErrNoChoicesAvailable) && prompt.Optional {
			return choose.ChooseOneResults{}, nil
		}
		if err != nil {
			return nil, fmt.Errorf("failed to choose one: %w", err)
		}
		if len(choices) == 0 {
			return choose.ChooseOneResults{}, nil
		}
		return choose.ChooseOneResults{Choice: choices[0]}, nil
	case choose.ChooseManyOpts:
		choices, err := a.chooseMany(opts.Choices, opts.Min, opts.Max, prompt.Message, prompt.Source, prompt.Optional)
		if errors.Is(err, choose.ErrNoChoicesAvailable) && prompt.Optional {
			return choose.ChooseManyResults{}, nil
		}
		return choose.ChooseManyResults{Choices: choices}, err
	case choose.MapChoicesToBucketsOpts:
		return a.chooseMapChoicesToBuckets(prompt.Message, opts, prompt.Source)
	case choose.ChooseNumberOpts:
		return a.chooseNumber(prompt.Message, opts, prompt.Source)
	default:
		return nil, fmt.Errorf("unsupported choice type: %T", opts)
	}
}

func (a *Agent) chooseNumber(
	message string,
	opts choose.ChooseNumberOpts,
	source choose.Source,
) (choose.ChoiceResults, error) {
	number, err := a.enterNumber(message, source)
	if err != nil {
		return nil, fmt.Errorf("failed to enter number: %w", err)
	}
	return choose.ChooseNumberResults{Number: number}, nil
}

func (a *Agent) chooseMapChoicesToBuckets(message string, opts choose.MapChoicesToBucketsOpts, source choose.Source) (choose.ChoiceResults, error) {
	assignments := map[choose.Bucket][]choose.Choice{}
	title := fmt.Sprintf("%s requires a choice", source.Name())
	userMessage := []string{}
	userMessage = append(userMessage, "Press Enter to complete selection")
	userMessage = append(userMessage, "Unassigned choices will be assigned to the last bucket by default")
	choices := slices.Clone(opts.Choices)
	for _, bucket := range opts.Buckets {
		if len(choices) == 0 {
			break
		}
		var remaining []choose.Choice
		a.uiBuffer.UpdateMessage(userMessage)
		a.uiBuffer.UpdateChoices(title, choices)
		if err := a.uiBuffer.Render(); err != nil {
			return nil, fmt.Errorf("failed to render UI Buffer: %w", err)
		}
		a.Prompt(message)
		fmt.Printf("assign cards to %s (space separated numbers): ", bucket)
		line := a.ReadInput()
		indexes := parseIndexes(line)
		for _, idx := range indexes {
			if idx < 1 || idx > len(choices) {
				fmt.Printf("Invalid card number: %d\n", idx)
				continue
			}
			assignments[bucket] = append(assignments[bucket], choices[idx-1])
		}
		for i, choice := range choices {
			idxFound := false
			for _, idx := range indexes {
				if i+1 == idx { // +1 because user input is 1-based
					idxFound = true
				}
			}
			if !idxFound {
				remaining = append(remaining, choice)
			}
		}
		choices = remaining
	}
	lastBucket := opts.Buckets[len(opts.Buckets)-1]
	for _, choice := range choices {
		fmt.Printf("Warning: Card %s was not assigned, defaulting to %s.\n", choice.Name(), lastBucket)
		assignments[lastBucket] = append(assignments[lastBucket], choice)
	}
	return choose.MapChoicesToBucketsResults{
		Assignments: assignments,
	}, nil
}

func (a *Agent) chooseMany(
	choices []choose.Choice,
	min int,
	max int,
	message string,
	source choose.Source,
	optional bool,
) ([]choose.Choice, error) {
	if len(choices) == 0 {
		// TODO is this an error?
		return nil, choose.ErrNoChoicesAvailable
	}
	title := fmt.Sprintf("%s requires a choice", source.Name())
	if min > 1 {
		if min == max {
			title = fmt.Sprintf(
				"%s requires %d choices",
				source.Name(),
				min,
			)
		} else {
			title = fmt.Sprintf(
				"%s requires %d - %d choices",
				source.Name(),
				min,
				max,
			)
		}
	}
	var selected []choose.Choice
	var userMessage []string
	if optional {
		userMessage = append(userMessage, "Optional: Press Enter to skip selection")
	}
	if min != max {
		userMessage = append(userMessage, "Press Enter to complete selection")
	}
	for len(selected) < max {
		a.uiBuffer.UpdateMessage(userMessage)
		a.uiBuffer.UpdateChoices(title, choices)
		if err := a.uiBuffer.Render(); err != nil {
			return nil, fmt.Errorf("failed to render UI Buffer: %w", err)
		}
		a.Prompt(message)
		max := len(choices)
		input := a.ReadInput()
		if input == "" {
			if optional || len(selected) >= min {
				break // User chose to skip selection
			}
			continue
		}
		choicePlus1, err := a.inputToNumber(input, max)
		if err != nil {
			fmt.Printf("Invalid choice. Please enter a number: %d - %d\n", 0, max)
			continue
		}
		choice := choicePlus1 - 1 // Convert to 0 based index
		selected = append(selected, choices[choice])
		userMessage = append(userMessage, fmt.Sprintf(
			"Selected: %s <%s>",
			choices[choice].Name(),
			choices[choice].ID(),
		))
		// TODO Should this follow the same immutable stuff
		choices = append(choices[:choice], choices[choice+1:]...)
	}
	return selected, nil
}

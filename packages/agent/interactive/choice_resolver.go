package interactive

import (
	"deckronomicon/packages/choose"
	"errors"
	"fmt"
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
		choices, err := a.ChooseMany(opts.Choices, 1, 1, prompt.Message, prompt.Source, prompt.Optional)
		if errors.Is(err, choose.ErrNoChoices) && prompt.Optional {
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
		choices, err := a.ChooseMany(opts.Choices, opts.Min, opts.Max, prompt.Message, prompt.Source, prompt.Optional)
		if errors.Is(err, choose.ErrNoChoices) && prompt.Optional {
			return choose.ChooseManyResults{}, nil
		}
		return choose.ChooseManyResults{Choices: choices}, err
	case choose.MapChoicesToBucketsOpts:
		return a.ChooseMapChoicesToBuckets(prompt.Message, opts, prompt.Source)
	default:
		return nil, fmt.Errorf("unsupported choice type: %T", opts)
	}
}

func (a *Agent) ChooseMapChoicesToBuckets(message string, opts choose.MapChoicesToBucketsOpts, source choose.Source) (choose.ChoiceResults, error) {
	assignments := make(map[choose.Bucket][]choose.Choice)
	assigned := make(map[int]bool) // track which cards have been assigned
	// For each bucket, prompt the user
	title := fmt.Sprintf("%s requires a choice", source.Name())
	userMessage := []string{}
	fmt.Println("choices:", opts.Choices)
	for _, bucket := range opts.Buckets {
		a.uiBuffer.UpdateMessage(userMessage)
		a.uiBuffer.UpdateChoices(title, opts.Choices)
		if err := a.uiBuffer.Render(); err != nil {
			return nil, fmt.Errorf("failed to render UI Buffer: %w", err)
		}
		a.Prompt(message)
		fmt.Printf("%s: ", bucket)
		line := a.ReadInput()
		indexes := parseIndexes(line)
		for _, idx := range indexes {
			if idx < 1 || idx > len(opts.Choices) {
				fmt.Printf("Invalid card number: %d\n", idx)
				continue
			}
			if assigned[idx] {
				fmt.Printf("Card %d already assigned, skipping.\n", idx)
				continue
			}
			assigned[idx] = true
			assignments[bucket] = append(assignments[bucket], opts.Choices[idx-1])
		}
	}
	// Handle unassigned cards → mama recommends putting them in last bucket by default!
	lastBucket := opts.Buckets[len(opts.Buckets)-1]
	for i := 1; i <= len(opts.Choices); i++ {
		if !assigned[i] {
			fmt.Printf("Warning: Card %d was not assigned, defaulting to %s.\n", i, lastBucket)
			assignments[lastBucket] = append(assignments[lastBucket], opts.Choices[i-1])
		}
	}
	return choose.MapChoicesToBucketsResults{
		Assignments: assignments,
	}, nil
}

func (a *Agent) ChooseMany(
	choices []choose.Choice,
	min int,
	max int,
	message string,
	source choose.Source,
	optional bool,
) ([]choose.Choice, error) {
	if len(choices) == 0 {
		// TODO is this an error?
		return nil, choose.ErrNoChoices
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
		userMessage = append(userMessage, "Optional: Choose 0 to skip")
	}
	if min != max {
		userMessage = append(userMessage, "Choose 0 to complete selection")
	}
	for len(selected) < max {
		a.uiBuffer.UpdateMessage(userMessage)
		a.uiBuffer.UpdateChoices(title, choices)
		if err := a.uiBuffer.Render(); err != nil {
			return nil, fmt.Errorf("failed to render UI Buffer: %w", err)
		}
		a.Prompt(message)
		max := len(choices)
		choicePlus1, err := a.ReadInputNumber(max)
		if err != nil {
			fmt.Printf("Invalid choice. Please enter a number: %d - %d\n", 0, max)
			continue
		}
		if choicePlus1 == 0 {
			if optional || len(selected) > min {
				break
			} else {
				continue
			}
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

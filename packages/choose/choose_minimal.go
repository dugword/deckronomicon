package choose

import (
	"fmt"
)

func ChooseMinimal(prompt ChoicePrompt) (ChoiceResults, error) {
	switch opts := prompt.ChoiceOpts.(type) {
	case ChooseOneOpts:
		if prompt.Optional {
			return nil, nil
		}
		if len(opts.Choices) == 0 {
			return nil, fmt.Errorf("no choices available")
		}
		return ChooseOneResults{Choice: opts.Choices[0]}, nil
	case ChooseManyOpts:
		if prompt.Optional || opts.Min == 0 {
			return nil, nil
		}
		if len(opts.Choices) < opts.Min {
			return nil, fmt.Errorf("not enough choices available")
		}
		return ChooseManyResults{Choices: opts.Choices[:opts.Min]}, nil
	case MapChoicesToBucketsOpts:
		if len(opts.Buckets) == 0 {
			return nil, fmt.Errorf("no buckets available")
		}
		return MapChoicesToBucketsResults{
			Assignments: map[Bucket][]Choice{
				opts.Buckets[0]: opts.Choices,
			},
		}, nil
	case ChooseNumberOpts:
		if prompt.Optional || opts.Min == 0 {
			return nil, nil
		}
		return ChooseNumberResults{Number: opts.Min}, nil
	default:
		return nil, fmt.Errorf("unknown choice options type: %T", opts)
	}
}

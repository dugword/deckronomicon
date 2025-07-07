package auto

import (
	"deckronomicon/packages/agent/auto/strategy/evalstate"
	"deckronomicon/packages/choose"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"fmt"
)

func (a *RuleBasedAgent) Choose(game *state.Game, prompt choose.ChoicePrompt) (choose.ChoiceResults, error) {
	if a.interactive {
		a.ReportState(game)
	}
	ctx := evalstate.EvalState{
		Game:     game,
		PlayerID: a.playerID,
		Mode:     a.mode,
	}
	for _, choice := range a.strategy.Choices[a.mode] {
		if choice.Name != prompt.Source.Name() {
			continue
		}
		if !choice.When.Evaluate(&ctx) {
			continue
		}
		switch opts := prompt.ChoiceOpts.(type) {
		case choose.ChooseOneOpts:
			if a.interactive {
				a.EnterToContinueOnChoices(fmt.Sprintf("Matched choice: %s", choice.Name), prompt.Message, opts.Choices)
			}
			chosen := choice.Choose.Select(NewQueryObjectsFromChoices(opts.Choices))
			if len(chosen) == 0 {
				return choose.ChooseMinimal(prompt)
			}
			return choose.ChooseOneResults{
				Choice: chosen[0],
			}, nil
		case choose.ChooseManyOpts:
			if a.interactive {
				a.EnterToContinueOnChoices(fmt.Sprintf("Matched choice: %s", choice.Name), prompt.Message, opts.Choices)
			}
			chosen := choice.Choose.Select(NewQueryObjectsFromChoices(opts.Choices))
			if len(chosen) < opts.Min { // todo get diff and union
				return choose.ChooseMinimal(prompt)
			}
			return choose.ChooseManyResults{
				Choices: choose.NewChoices(chosen[:opts.Min]), // Return the minimum number of choices
			}, nil
		case choose.MapChoicesToBucketsOpts:
			if a.interactive {
				a.EnterToContinueOnChoices(fmt.Sprintf("Matched choice: %s", choice.Name), prompt.Message, opts.Choices)
			}
			chosen := choice.Choose.Select(NewQueryObjectsFromChoices(opts.Choices))
			if len(chosen) == 0 {
				return choose.MapChoicesToBucketsResults{
					Assignments: map[choose.Bucket][]choose.Choice{
						opts.Buckets[len(opts.Buckets)-1]: opts.Choices,
					},
				}, nil
			}
			return choose.MapChoicesToBucketsResults{
				Assignments: map[choose.Bucket][]choose.Choice{
					opts.Buckets[0]: choose.NewChoices(chosen),
				},
			}, nil
		case choose.ChooseNumberOpts:
		default:
			return nil, fmt.Errorf("unsupported choice options type: %T", opts)
		}
	}
	return choose.ChooseMinimal(prompt)
}

func NewQueryObjectsFromChoices(choices []choose.Choice) []query.Object {
	var objs []query.Object
	for _, choice := range choices {
		if obj, ok := choice.(query.Object); ok {
			objs = append(objs, obj)
		}
	}
	return objs
}

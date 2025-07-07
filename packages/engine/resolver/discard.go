package resolver

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/state"
	"errors"
)

func ResolveDiscard(
	game *state.Game,
	playerID string,
	discard *effect.Discard,
	target target.Target,
	resolvable state.Resolvable,
) (Result, error) {
	player := game.GetPlayer(playerID)
	cards := player.Hand().GetAll()
	choicePrompt := choose.ChoicePrompt{
		Message: "Chose cards to discard",
		Source:  resolvable,
		ChoiceOpts: choose.ChooseManyOpts{
			Choices: choose.NewChoices(cards),
			Min:     discard.Count,
			Max:     discard.Count,
		},
	}
	resumeFunc := func(choiceResults choose.ChoiceResults) (Result, error) {
		selected, ok := choiceResults.(choose.ChooseManyResults)
		if !ok {
			return Result{}, errors.New("invalid choice results for Discarding")
		}
		var events []event.GameEvent
		for _, choice := range selected.Choices {
			// Create the discard event
			discardEvent := &event.DiscardCardEvent{
				PlayerID: player.ID(),
				CardID:   choice.ID(),
			}
			// Apply the discard event to the game state
			events = append(events, discardEvent)
		}
		return Result{
			Events: events,
		}, nil
	}
	return Result{
		ChoicePrompt: choicePrompt,
		Resume:       resumeFunc,
	}, nil
}

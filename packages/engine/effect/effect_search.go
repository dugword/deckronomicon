package effect

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

func SearchEffectHandler(
	game state.Game,
	player state.Player,
	source query.Object,
	modifiers []definition.EffectModifier,
) (EffectResult, error) {
	query, err := buildQuery(modifiers)
	if err != nil {
		return EffectResult{}, fmt.Errorf("failed to build query for Search effect: %w", err)
	}
	cards := player.Library().FindAll(query)
	choicePrompt := choose.ChoicePrompt{
		// TODO: provide more detail on what kind of card to choose
		Message: "Choose a card to put into your hand",
		Source:  source,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(cards),
		},
	}
	resumeFunc := func(choiceResults choose.ChoiceResults) (EffectResult, error) {
		selected, ok := choiceResults.(choose.ChooseOneResults)
		if !ok {
			return EffectResult{}, fmt.Errorf("expected a single choice result")
		}
		card, ok := selected.Choice.(gob.Card)
		if !ok {
			return EffectResult{}, errors.New("choice is not a card")
		}
		events := []event.GameEvent{
			event.PutCardInHandEvent{
				PlayerID: player.ID(),
				CardID:   card.ID(),
				FromZone: mtg.ZoneLibrary,
			},
		}
		return EffectResult{
			Events: events,
		}, nil
	}
	// Need to get choices
	return EffectResult{
		ChoicePrompt: choicePrompt,
		ResumeFunc:   resumeFunc,
	}, nil
}

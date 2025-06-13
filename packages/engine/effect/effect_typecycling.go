package effect

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

func TypecyclingEffectHandler(
	game state.Game,
	player state.Player,
	source query.Object,
	modifiers []definition.EffectModifier,
) (EffectResult, error) {
	var subtypeValue string
	for _, modifier := range modifiers {
		if modifier.Key == "Subtype" {
			subtypeValue = modifier.Value
			break
		}
	}
	if subtypeValue == "" {
		return EffectResult{}, fmt.Errorf("missing required modifier %q for Typecycling effect", "Subtype")
	}
	subtype, ok := mtg.StringToSubtype(subtypeValue)
	if !ok {
		return EffectResult{}, fmt.Errorf("invalid subtype %q for Typecycling effect", subtypeValue)
	}
	cards := player.Library().FindAll(
		has.Subtype(subtype),
	)
	choicePrompt := choose.ChoicePrompt{
		Message: fmt.Sprintf("Choose a card with subtype %q to put into your hand", subtype),
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

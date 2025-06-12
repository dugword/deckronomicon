package effect

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
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
	modifiers []gob.Tag,
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
		Choices:    choose.NewChoices(cards),
		MinChoices: 1,
		MaxChoices: 1,
		Message:    fmt.Sprintf("Choose a card with subtype %q to put into your hand", subtype),
		Source:     source,
	}
	resumeFunc := func(choices []choose.Choice) (EffectResult, error) {
		if len(choices) == 0 {
			return EffectResult{}, errors.New("no choices selected for Typecyling")
		}
		card, ok := choices[0].(gob.Card)
		if !ok {
			return EffectResult{}, errors.New("choice is not a card")
		}
		events := []event.GameEvent{
			event.MoveCardEvent{
				PlayerID: player.ID(),
				CardID:   card.ID(),
				FromZone: mtg.ZoneLibrary,
				ToZone:   mtg.ZoneHand,
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

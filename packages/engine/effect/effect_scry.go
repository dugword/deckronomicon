package effect

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"errors"
)

func ScryEffectHandler(
	game state.Game,
	player state.Player,
	source query.Object,
	modifiers []gob.Tag,
) (EffectResult, error) {
	return EffectResult{}, errors.New("ScryEffectHandler not implemented")
	/*
		scryCount := 1
		for _, modifier := range modifiers {
			if modifier.Key == "Count" {
				count, err := strconv.Atoi(modifier.Value)
				if err != nil {
					return EffectResult{}, fmt.Errorf("invalid modifier %q for Scry effect: %w", modifier.Key, err)
				}
				scryCount = count
			}
		}
		if scryCount == 0 {
			return EffectResult{}, fmt.Errorf("missing required modifier %q for Scry effect", "Count")
		}
		choicePrompt := choose.ChoicePrompt{
			Choices:    choose.NewChoices(cards),
			MinChoices: 1,
			MaxChoices: 1,
			Message:    fmt.Sprintf("Choose a card with subtype %q to put into your hand", subtype),
			Source:     source,
		}
		resumeFunc := func(choices []choose.Choice) (EffectResult, error) {
			if len(choices) == 0 {
				return EffectResult{}, errors.New("no choices selected for Scrying")
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
	*/
}

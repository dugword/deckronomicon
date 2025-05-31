package effect

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/target"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/take"
	"deckronomicon/packages/state"
	"encoding/json"
	"errors"
	"fmt"
)

type LookAndChooseEffect struct {
	Look      int            `json:"Look"`
	Choose    int            `json:"Choose"`
	CardTypes []mtg.CardType `json:"CardTypes"`
	Rest      mtg.Zone       `json:"Rest"`
	Order     string         `json:"Order"`
}

func NewLookAndChooseEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var lookAndChooseEffect LookAndChooseEffect
	if err := json.Unmarshal(effectSpec.Modifiers, &lookAndChooseEffect); err != nil {
		return nil, fmt.Errorf("failed to unmarshal LookAndChooseEffectModifiers: %w", err)
	}
	return lookAndChooseEffect, nil
}

func (e LookAndChooseEffect) Name() string {
	return "LookAndChoose"
}

func (e LookAndChooseEffect) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}

func (e LookAndChooseEffect) Resolve(
	game state.Game,
	player state.Player,
	source query.Object,
	target target.TargetValue,
) (EffectResult, error) {
	if e.Look <= 0 {
		return EffectResult{}, fmt.Errorf("invalid required modifier %q for Scry effect", "Count")
	}
	cards, _ := player.Library().TakeN(e.Look)
	var events []event.GameEvent
	for _, card := range cards {
		events = append(events, event.RevealCardEvent{
			PlayerID: player.ID(),
			CardID:   card.ID(),
			FromZone: mtg.ZoneLibrary,
		})
	}
	// TODO: Rename build predicate
	predicate, err := buildQuery(QueryOpts{
		CardTypes: e.CardTypes,
	})
	if err != nil {
		return EffectResult{}, fmt.Errorf("failed to build query for LookAndChoose: %w", err)
	}
	choiceCards := query.FindAll(cards, predicate)
	fmt.Println("LookAndChoose found cards:", len(choiceCards))
	choicePrompt := choose.ChoicePrompt{
		// TODO: Add type information to message
		Message:  fmt.Sprintf("Look at the top %d cards of your library. Choose %d of them", e.Look, e.Choose),
		Source:   source,
		Optional: true,
		ChoiceOpts: choose.ChooseManyOpts{
			Choices: choose.NewChoices(choiceCards),
			Max:     e.Choose,
		},
	}
	resumeFunc := func(choiceResults choose.ChoiceResults) (EffectResult, error) {
		selected, ok := choiceResults.(choose.ChooseManyResults)
		if !ok {
			return EffectResult{}, errors.New("invalid choice results for LookAndChoose")
		}
		var selectedCards []gob.Card
		for _, choice := range selected.Choices {
			taken, remaining, ok := take.By(cards, has.ID(choice.ID()))
			if !ok {
				return EffectResult{}, fmt.Errorf("selected card %q not found in looked at cards", choice.ID())
			}
			selectedCards = append(selectedCards, taken)
			cards = remaining
		}
		var events []event.GameEvent
		for _, card := range selectedCards {
			events = append(events, event.PutCardInHandEvent{
				PlayerID: player.ID(),
				CardID:   card.ID(),
				FromZone: mtg.ZoneLibrary,
			})
		}
		if e.Rest == mtg.ZoneLibrary {
			// Put the rest on the bottom of the library in a random order
			// TODO: Support ordering
			for _, card := range cards {
				events = append(events, event.PutCardOnBottomOfLibraryEvent{
					PlayerID: player.ID(),
					CardID:   card.ID(),
					FromZone: mtg.ZoneLibrary,
				})
			}
		} else if e.Rest == mtg.ZoneGraveyard {
			// Put the rest in the graveyard
			for _, card := range cards {
				events = append(events, event.PutCardInGraveyardEvent{
					PlayerID: player.ID(),
					CardID:   card.ID(),
					FromZone: mtg.ZoneLibrary,
				})
			}
		} else {
			return EffectResult{}, fmt.Errorf("invalid Rest zone %q for LookAndChoose", e.Rest)
		}
		return EffectResult{
			Events: events,
		}, nil
	}
	return EffectResult{
		Events:       events,
		ChoicePrompt: choicePrompt,
		ResumeFunc:   resumeFunc,
	}, nil
}

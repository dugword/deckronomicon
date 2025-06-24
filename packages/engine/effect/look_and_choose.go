package effect

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/take"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

type LookAndChooseEffect struct {
	Look      int
	Choose    int
	CardTypes []mtg.CardType
	Rest      mtg.Zone
	Order     string
}

func NewLookAndChooseEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var lookAndChooseEffect LookAndChooseEffect
	look, ok := effectSpec.Modifiers["Look"].(int)
	if !ok || look <= 0 {
		return nil, fmt.Errorf("LookAndChooseEffect requires a 'Look' modifier of type int greater than 0, got %T", effectSpec.Modifiers["Look"])
	}
	lookAndChooseEffect.Look = look
	choose, ok := effectSpec.Modifiers["Choose"].(int)
	if !ok || choose <= 0 {
		return nil, fmt.Errorf("LookAndChooseEffect requires a 'Choose' modifier of type int greater than 0, got %T", effectSpec.Modifiers["Choose"])
	}
	lookAndChooseEffect.Choose = choose
	cardTypeVals, ok := effectSpec.Modifiers["CardTypes"].([]any)
	if !ok || len(cardTypeVals) == 0 {
		return nil, fmt.Errorf("LookAndChooseEffect requires a non-empty 'CardTypes' modifier of type []mtg.CardType, got %T", effectSpec.Modifiers["CardTypes"])
	}
	var cardTypes []mtg.CardType
	for _, cardTypeVal := range cardTypeVals {
		cardTypeString, ok := cardTypeVal.(string)
		if !ok {
			return nil, fmt.Errorf("LookAndChooseEffect requires 'CardTypes' modifier values to be of type string, got %T", cardTypeVal)
		}
		cardType, ok := mtg.StringToCardType(cardTypeString)
		if !ok {
			return nil, fmt.Errorf("LookAndChooseEffect requires valid 'CardTypes' modifier with values like 'Creature', 'Instant', etc., got %q", cardTypeString)
		}
		cardTypes = append(cardTypes, cardType)
	}
	lookAndChooseEffect.CardTypes = cardTypes
	restString, ok := effectSpec.Modifiers["Rest"].(string)
	if !ok {
		return nil, fmt.Errorf("LookAndChooseEffect requires a 'Rest' modifier of type string, got %T", effectSpec.Modifiers["Rest"])
	}
	rest, ok := mtg.StringToZone(restString)
	if !ok && (rest != mtg.ZoneLibrary && rest != mtg.ZoneGraveyard) {
		return nil, fmt.Errorf("LookAndChooseEffect requires a 'Rest' modifier of type mtg.Zone with value 'Library' or 'Graveyard', got %T", effectSpec.Modifiers["Rest"])
	}
	lookAndChooseEffect.Rest = rest
	order, ok := effectSpec.Modifiers["Order"].(string)
	if ok && rest == mtg.ZoneLibrary && order != "Any" && order != "Random" {
		return nil, fmt.Errorf("LookAndChooseEffect requires an 'Order' modifier of type string with value 'Any' or 'Random' when rest is Library, got %T", effectSpec.Modifiers["Order"])
	}
	lookAndChooseEffect.Order = order
	// Validate that Look is greater than or equal to Choose
	if lookAndChooseEffect.Look < lookAndChooseEffect.Choose {
		return nil, fmt.Errorf("LookAndChooseEffect requires 'Look' (%d) to be greater than or equal to 'Choose' (%d)", lookAndChooseEffect.Look, lookAndChooseEffect.Choose)
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
	resEnv *resenv.ResEnv,
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

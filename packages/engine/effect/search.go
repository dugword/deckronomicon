package effect

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/target"

	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

type SearchEffect struct {
	CardTypes  []mtg.CardType
	Colors     []mtg.Color
	Subtypes   []mtg.Subtype
	ManaValues []int
}

func NewSearchEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var searchEffect SearchEffect
	cardTypesRaw, ok := effectSpec.Modifiers["CardTypes"].([]any)
	if ok {
		for _, cardTypeRaw := range cardTypesRaw {
			cardType, ok := mtg.StringToCardType(fmt.Sprintf("%v", cardTypeRaw))
			if !ok {
				return nil, fmt.Errorf("SearchEffect requires a 'CardTypes' modifier of type []mtg.CardType, got %T", cardTypeRaw)
			}
			searchEffect.CardTypes = append(searchEffect.CardTypes, cardType)
		}
	}
	colorsRaw, ok := effectSpec.Modifiers["Colors"].([]any)
	if ok {
		for _, colorRaw := range colorsRaw {
			color, ok := mtg.StringToColor(fmt.Sprintf("%v", colorRaw))
			if !ok {
				return nil, fmt.Errorf("SearchEffect requires a 'Colors' modifier of type []mtg.Color, got %T", colorRaw)
			}
			searchEffect.Colors = append(searchEffect.Colors, color)
		}
	}
	subtypesRaw, ok := effectSpec.Modifiers["Subtypes"].([]any)
	if ok {
		for _, subtypeRaw := range subtypesRaw {
			subtype, ok := mtg.StringToSubtype(fmt.Sprintf("%v", subtypeRaw))
			if !ok {
				return nil, fmt.Errorf("SearchEffect requires a 'Subtypes' modifier of type []mtg.Subtype, got %T", subtypeRaw)
			}
			searchEffect.Subtypes = append(searchEffect.Subtypes, subtype)
		}
	}
	manaValuesRaw, ok := effectSpec.Modifiers["ManaValues"].([]any)
	if ok {
		for _, manaValueRaw := range manaValuesRaw {
			manaValue, ok := manaValueRaw.(int)
			if !ok {
				return nil, fmt.Errorf("SearchEffect requires a 'ManaValues' modifier of type []int, got %T", manaValueRaw)
			}
			searchEffect.ManaValues = append(searchEffect.ManaValues, manaValue)
		}
	}
	return searchEffect, nil
}

func (e SearchEffect) Name() string {
	return "Search"
}

func (e SearchEffect) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}

func (e SearchEffect) Resolve(
	game state.Game,
	player state.Player,
	source query.Object,
	target target.TargetValue,
	resEnv *resenv.ResEnv,
) (EffectResult, error) {
	query, err := buildQuery(QueryOpts(e))
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

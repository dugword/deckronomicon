package effect

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

type CounterspellEffect struct {
	CardTypes  []mtg.CardType `json:"CardTypes,omitempty"`
	Colors     []mtg.Color    `json:"Colors,omitempty"`
	Subtypes   []mtg.Subtype  `json:"Subtypes,omitempty"`
	ManaValues []int          `json:"ManaValues,omitempty"`
}

func NewCounterspellEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var counterspellEffect CounterspellEffect
	cardTypesRaw, ok := effectSpec.Modifiers["CardTypes"].([]any)
	if ok {
		for _, cardTypeRaw := range cardTypesRaw {
			cardType, ok := mtg.StringToCardType(fmt.Sprintf("%v", cardTypeRaw))
			if !ok {
				return nil, fmt.Errorf("CounterspellEffect requires a 'CardTypes' modifier of type []mtg.CardType, got %T", cardTypeRaw)
			}
			counterspellEffect.CardTypes = append(counterspellEffect.CardTypes, cardType)
		}
	}
	colorsRaw, ok := effectSpec.Modifiers["Colors"].([]any)
	if ok {
		for _, colorRaw := range colorsRaw {
			color, ok := mtg.StringToColor(fmt.Sprintf("%v", colorRaw))
			if !ok {
				return nil, fmt.Errorf("CounterspellEffect requires a 'Colors' modifier of type []mtg.Color, got %T", colorRaw)
			}
			counterspellEffect.Colors = append(counterspellEffect.Colors, color)
		}
	}
	subtypesRaw, ok := effectSpec.Modifiers["Subtypes"].([]any)
	if ok {
		for _, subtypeRaw := range subtypesRaw {
			subtype, ok := mtg.StringToSubtype(fmt.Sprintf("%v", subtypeRaw))
			if !ok {
				return nil, fmt.Errorf("CounterspellEffect requires a 'Subtypes' modifier of type []mtg.Subtype, got %T", subtypeRaw)
			}
			counterspellEffect.Subtypes = append(counterspellEffect.Subtypes, subtype)
		}
	}
	manaValuesRaw, ok := effectSpec.Modifiers["ManaValues"].([]any)
	if ok {
		for _, manaValueRaw := range manaValuesRaw {
			manaValue, ok := manaValueRaw.(int)
			if !ok {
				return nil, fmt.Errorf("CounterspellEffect requires a 'ManaValues' modifier of type []int, got %T", manaValueRaw)
			}
			counterspellEffect.ManaValues = append(counterspellEffect.ManaValues, manaValue)
		}
	}
	return counterspellEffect, nil
}

func (e CounterspellEffect) Name() string {
	return "Counterspell"
}

func (e CounterspellEffect) TargetSpec() target.TargetSpec {
	return target.SpellTargetSpec{
		CardTypes:  e.CardTypes,
		Colors:     e.Colors,
		Subtypes:   e.Subtypes,
		ManaValues: e.ManaValues,
	}
}

// TODO: Does Resolve need to check that the target is valid?
func (e CounterspellEffect) Resolve(
	game state.Game,
	player state.Player,
	source query.Object,
	target target.TargetValue,
	resEnv *resenv.ResEnv,
) (EffectResult, error) {
	resolvable, ok := game.Stack().Find(has.ID(target.TargetID))
	if !ok {
		return EffectResult{
			Events: []event.GameEvent{
				event.SpellOrAbilityFizzlesEvent{
					PlayerID: player.ID(),
					ObjectID: target.TargetID,
				},
			},
		}, nil
	}
	spell, ok := resolvable.(gob.Spell)
	if !ok {
		return EffectResult{}, errors.New("choice is not a spell")
	}
	events := []event.GameEvent{
		event.RemoveSpellOrAbilityFromStackEvent{
			PlayerID: spell.Owner(),
			ObjectID: spell.ID(),
		},
	}
	return EffectResult{
		Events: events,
	}, nil
}

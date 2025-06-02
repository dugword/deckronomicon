package effect

import (
	"deckronomicon/packages/game/core"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"errors"
	"fmt"
)

// BuildEffectAdditionalMana creates an effect that adds additional mana when
// a trigger happens, like tapping an island for mana.
// Supported Modifier Keys:
//   - Mana: <ManaString>
//   - Target: <subtype>
//   - Duration: <eventType>
func BuildEffectAdditionalMana(source query.Object, spec definition.EffectSpec) (*Effect, error) {
	var mana string
	var target string
	var duration string
	for _, modifier := range spec.Modifiers {
		if modifier.Key == "Mana" {
			mana += modifier.Value
		}
		if modifier.Key == "Target" {
			target = modifier.Value
		}
		if modifier.Key == "Duration" {
			duration = modifier.Value
		}
	}
	if mana == "" {
		return nil, errors.New("no mana string provided")
	}
	if !mtg.IsMana(mana) {
		return nil, fmt.Errorf("invalid mana string: %s", mana)
	}
	if target == "" {
		return nil, errors.New("no target provided")
	}
	if duration != "EndOfTurn" {
		// return nil, errors.New("no duration provided")
		return nil, errors.New("only EndOfTurn duration is supported")
	}
	_, err := mtg.StringToSubtype(target)
	if err != nil {
		return nil, fmt.Errorf("invalid target subtype: %s", target)
	}
	effect := Effect{id: spec.ID}
	// id := getNextEventID()
	/*
		eventHandler := EventHandler{
			ID: id,
			Callback: func(event Event, state *Gamecore.State, player *core.Player) {
				// Move this into the register so I don't have to check for
				// it.
				if event.Type != EventTapForMana {
					return
				}
				if !event.Source.HasSubtype(subtype) {
					return
				}
				if err := player.ManaPool.AddMana(mana); err != nil {
					// TODO: Handle this better
					panic("failed to add mana: " + err.Error())
				}
				return
			},
		}
	*/
	var tags []core.Tag
	/*
		for _, symbol := range mtg.ManaStringToManaSymbols(mana) {
			tags = append(tags, Tag{Key: AbilityTagManaAbility, Value: symbol})
		}
	*/
	effect.tags = tags
	effect.description = fmt.Sprintf("add additional %s when you tap an %s for mana", mana, target)
	effect.Apply = func(state core.State, player core.Player) error {
		/*
			state.RegisterListenerUntil(
				eventHandler,
				EventEndStep,
			)
		*/
		return nil
	}
	return &effect, nil
}

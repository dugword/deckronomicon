package effect

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/target"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"encoding/json"
	"fmt"
)

type AddManaEffect struct {
	Mana string `json:"Mana"`
}

func NewAddManaEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var addManaEffect AddManaEffect
	if err := json.Unmarshal(effectSpec.Modifiers, &addManaEffect); err != nil {
		return nil, fmt.Errorf("failed to unmarshal AddManaEffect: %w", err)
	}
	return addManaEffect, nil
}

func (e AddManaEffect) Name() string {
	return "AddMana"
}

func (e AddManaEffect) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}

func (e AddManaEffect) Resolve(
	game state.Game,
	player state.Player,
	source query.Object,
	target target.TargetValue,
) (EffectResult, error) {
	amount, err := mana.ParseManaString(e.Mana)
	if err != nil {
		return EffectResult{}, fmt.Errorf("failed to parse mana string %q: %w", e.Mana, err)
	}

	// Think through how to best handle this and how the events will be represented in JSON.
	// I could have the mana struct pretty print to a string like "2{R}{G}" or and then reparse it when I apply the event.
	// Or I could generate multiple events for each color of mana like I am doing here.
	var events []event.GameEvent
	for color, amount := range amount.Colors() {
		if amount <= 0 {
			continue // Skip colors with no mana
		}
		events = append(events, event.AddManaEvent{
			PlayerID: player.ID(),
			Amount:   amount,
			ManaType: color,
		})
	}
	return EffectResult{
		Events: events,
	}, nil
}

package effect

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"fmt"
)

type AddManaEffect struct {
	Mana string
}

func NewAddManaEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var addManaEffect AddManaEffect
	manaString, ok := effectSpec.Modifiers["Mana"].(string)
	if !ok {
		return nil, fmt.Errorf("AddManaEffect requires a 'Mana' modifier of type string, got %T", effectSpec.Modifiers["Mana"])
	}
	addManaEffect.Mana = manaString
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
	resEnv *resenv.ResEnv,
) (EffectResult, error) {
	amount, err := mana.ParseManaString(e.Mana)
	if err != nil {
		return EffectResult{}, fmt.Errorf("failed to parse mana string %q: %w", e.Mana, err)
	}
	var events []event.GameEvent
	if amount.Generic() > 0 {
		events = append(events, event.AddManaEvent{
			PlayerID: player.ID(),
			Amount:   amount.Generic(),
			Color:    mana.Colorless,
		})
	}
	if amount.Colorless() > 0 {
		events = append(events, event.AddManaEvent{
			PlayerID: player.ID(),
			Amount:   amount.Colorless(),
			Color:    mana.Colorless,
		})
	}
	if amount.White() > 0 {
		events = append(events, event.AddManaEvent{
			PlayerID: player.ID(),
			Amount:   amount.White(),
			Color:    mana.White,
		})
	}
	if amount.Blue() > 0 {
		events = append(events, event.AddManaEvent{
			PlayerID: player.ID(),
			Amount:   amount.Blue(),
			Color:    mana.Blue,
		})
	}
	if amount.Black() > 0 {
		events = append(events, event.AddManaEvent{
			PlayerID: player.ID(),
			Amount:   amount.Black(),
			Color:    mana.Black,
		})
	}
	if amount.Red() > 0 {
		events = append(events, event.AddManaEvent{
			PlayerID: player.ID(),
			Amount:   amount.Red(),
			Color:    mana.Red,
		})
	}
	if amount.Green() > 0 {
		events = append(events, event.AddManaEvent{
			PlayerID: player.ID(),
			Amount:   amount.Green(),
			Color:    mana.Green,
		})
	}
	return EffectResult{
		Events: events,
	}, nil
}

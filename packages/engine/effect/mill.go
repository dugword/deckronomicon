package effect

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"encoding/json"
	"fmt"
)

type MillEffect struct {
	Count  int    `json:"Count"`
	Target string `json:"Target"`
}

func NewMillEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var MillEffect MillEffect
	if err := json.Unmarshal(effectSpec.Modifiers, &MillEffect); err != nil {
		return nil, fmt.Errorf("failed to unmarshal MillEffectModifiers: %w", err)
	}
	return MillEffect, nil
}

func (d MillEffect) Name() string {
	return "Mill"
}

func (d MillEffect) TargetSpec() target.TargetSpec {
	switch d.Target {
	case "":
		return target.NoneTargetSpec{}
	case "Player":
		return target.PlayerTargetSpec{}
	default:
		panic(fmt.Sprintf("unknown target spec %q for MillEffect", d.Target))
		return target.NoneTargetSpec{}
	}
}

func (e MillEffect) Resolve(
	game state.Game,
	player state.Player,
	source query.Object,
	target target.TargetValue,
) (EffectResult, error) {
	cards := player.Library().GetN(e.Count)
	var events []event.GameEvent
	for _, card := range cards {
		events = append(events, event.PutCardInGraveyardEvent{
			PlayerID: player.ID(),
			CardID:   card.ID(),
			FromZone: mtg.ZoneLibrary,
		})
	}
	return EffectResult{Events: events}, nil
}

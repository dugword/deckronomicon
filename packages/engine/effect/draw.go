package effect

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/target"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"encoding/json"
	"fmt"
)

type DrawEffect struct {
	Count  int    `json:"Count"`
	Target string `json:"Target"`
}

func NewDrawEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var drawEffect DrawEffect
	if err := json.Unmarshal(effectSpec.Modifiers, &drawEffect); err != nil {
		return nil, fmt.Errorf("failed to unmarshal DrawEffectModifiers: %w", err)
	}
	return drawEffect, nil
}

func (d DrawEffect) Name() string {
	return "Draw"
}

func (d DrawEffect) TargetSpec() target.TargetSpec {
	switch d.Target {
	case "":
		fmt.Println("Returning NoneTargetSpec for empty target in DrawEffect")
		return target.NoneTargetSpec{}
	case "Player":
		return target.PlayerTargetSpec{}
	default:
		panic(fmt.Sprintf("unknown target spec %q for DrawEffect", d.Target))
		return target.NoneTargetSpec{}
	}
}

func (e DrawEffect) Resolve(
	game state.Game,
	player state.Player,
	source query.Object,
	trgt target.TargetValue,
) (EffectResult, error) {
	switch trgt.TargetType {
	case target.TargetTypeNone:
		return e.resolveForPlayer(game, player.ID())
	case target.TargetTypePlayer:
		return e.resolveForPlayer(game, trgt.PlayerID)
	default:
		panic(fmt.Sprintf("unexpected target type %s for DrawEffect", trgt.TargetType))
		return EffectResult{}, nil
	}
}

func (e DrawEffect) resolveForPlayer(game state.Game, playerID string) (EffectResult, error) {
	var events []event.GameEvent
	for range e.Count {
		events = append(events, event.DrawCardEvent{
			PlayerID: playerID,
		})
	}
	return EffectResult{
		Events: events,
	}, nil
}

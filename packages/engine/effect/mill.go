package effect

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"fmt"
)

type MillEffect struct {
	Count  int
	Target string
}

func NewMillEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var MillEffect MillEffect
	count, ok := effectSpec.Modifiers["Count"].(int)
	if !ok || count <= 0 {
		return nil, fmt.Errorf("MillEffect requires a 'Count' modifier of type int greater than 0, got %T", effectSpec.Modifiers["Count"])
	}
	MillEffect.Count = count
	targetStr, ok := effectSpec.Modifiers["Target"].(string)
	if ok && targetStr != "" && targetStr != "Player" {
		return nil, fmt.Errorf("MillEffect requires a 'Target' modifier of either empty or 'Player', got %q", targetStr)
	}
	MillEffect.Target = targetStr
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
	resEnv *resenv.ResEnv,
) (EffectResult, error) {
	targetPlayerID := player.ID()
	if target.TargetID != "" {
		targetPlayerID = target.TargetID
	}
	// TODO: This is getting a player from user input,
	// That means this could fail if the player ID is invalid.
	// We need to check that before proceeding.
	targetPlayer := game.GetPlayer(targetPlayerID)
	cards := targetPlayer.Library().GetN(e.Count)
	var events []event.GameEvent
	for _, card := range cards {
		events = append(events, event.PutCardInGraveyardEvent{
			PlayerID: targetPlayerID,
			CardID:   card.ID(),
			FromZone: mtg.ZoneLibrary,
		})
	}
	return EffectResult{Events: events}, nil
}

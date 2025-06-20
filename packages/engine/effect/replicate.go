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
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"encoding/json"
	"fmt"
)

type ReplicateEffect struct {
	Count int `json:"Count,omitempty"`
}

func NewReplicateEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var replicateEffect ReplicateEffect
	if err := json.Unmarshal(effectSpec.Modifiers, &replicateEffect); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ReplicateEffect: %w", err)
	}
	return replicateEffect, nil
}

func (e ReplicateEffect) Name() string {
	return "Replicate"
}

func (e ReplicateEffect) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}

// TODO: Does Resolve need to check that the target is valid?
func (e ReplicateEffect) Resolve(
	game state.Game,
	player state.Player,
	source query.Object,
	targetValue target.TargetValue,
	resEnv *resenv.ResEnv,
) (EffectResult, error) {
	resolvable, ok := game.Stack().Find(query.And(is.Spell(), has.SourceID(targetValue.ObjectID)))
	if !ok {
		return EffectResult{
			Events: []event.GameEvent{
				event.SpellOrAbilityFizzlesEvent{
					PlayerID: player.ID(),
					ObjectID: targetValue.ObjectID,
				},
			},
		}, nil
	}
	spell, ok := resolvable.(gob.Spell)
	if !ok {
		return EffectResult{}, fmt.Errorf("resolvable with ID %q is %T not a spell", targetValue.ObjectID, resolvable)
	}
	var effectWithNewTargets []gob.EffectWithTarget
	for _, effectWithTarget := range spell.EffectWithTargets() {
		if effectWithTarget.Target.TargetType != target.TargetTypeNone {
			effectWithNewTargets = append(effectWithNewTargets, effectWithTarget)
			continue
		}
		// TODO: Do something where I return a chain of EffectResult for each target so I can prompt the
		// player for a new target for each effect.
		effectWithNewTargets = append(effectWithNewTargets, effectWithTarget)
	}
	events := []event.GameEvent{}
	for range e.Count {
		events = append(events, event.PutCopiedSpellOnStackEvent{
			PlayerID:          player.ID(),
			SpellID:           spell.ID(),
			EffectWithTargets: effectWithNewTargets,
			FromZone:          mtg.ZoneHand,
		})
	}
	return EffectResult{
		Events: events,
	}, nil
}

package resolver

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"fmt"
)

func ResolveReplicate(
	game *state.Game,
	playerID string,
	replicate *effect.Replicate,
	source gob.Object,
	resEnv *resenv.ResEnv,
) (Result, error) {
	resolvable, ok := game.Stack().Find(query.And(is.Spell(), has.SourceID(source.ID())))
	if !ok {
		return Result{
			Events: []event.GameEvent{
				&event.SpellOrAbilityFizzlesEvent{
					PlayerID: playerID,
					ObjectID: source.ID(),
				},
			},
		}, nil
	}
	spell, ok := resolvable.(*gob.Spell)
	if !ok {
		return Result{}, fmt.Errorf("resolvable with ID %q is %T not a spell", source.ID(), resolvable)
	}
	var effectWithNewTargets []*effect.EffectWithTarget
	for _, effectWithTarget := range spell.EffectWithTargets() {
		if effectWithTarget.Target.Type != mtg.TargetTypeNone {
			effectWithNewTargets = append(effectWithNewTargets, effectWithTarget)
			continue
		}
		// TODO: Do something where I return a chain of EffectResult for each target so I can prompt the
		// player for a new target for each effect.
		effectWithNewTargets = append(effectWithNewTargets, effectWithTarget)
	}
	events := []event.GameEvent{}
	for range replicate.Count {
		events = append(events, &event.PutCopiedSpellOnStackEvent{
			PlayerID:          playerID,
			SpellID:           spell.ID(),
			EffectWithTargets: effectWithNewTargets,
			FromZone:          mtg.ZoneHand,
		})
	}
	return Result{
		Events: events,
	}, nil
}

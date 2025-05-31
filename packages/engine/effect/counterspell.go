package effect

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/target"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"

	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"encoding/json"
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
	if err := json.Unmarshal(effectSpec.Modifiers, &counterspellEffect); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CounterspellEffect: %w", err)
	}
	return counterspellEffect, nil
}

func (e CounterspellEffect) Name() string {
	return "Counterspell"
}

func (e CounterspellEffect) TargetSpec() target.TargetSpec {
	query, err := buildQuery(QueryOpts(e))
	if err != nil {
		panic(fmt.Errorf("failed to build query for Search effect: %w", err))
	}
	return target.SpellTargetSpec{
		Predicate: query,
	}
}

// TODO: Does Resolve need to check that the target is valid?
func (e CounterspellEffect) Resolve(
	game state.Game,
	player state.Player,
	source query.Object,
	target target.TargetValue,
) (EffectResult, error) {
	resolvable, ok := game.Stack().Find(has.ID(target.ObjectID))
	if !ok {
		return EffectResult{
			Events: []event.GameEvent{
				event.SpellOrAbilityFizzlesEvent{
					PlayerID: player.ID(),
					ObjectID: target.ObjectID,
				},
			},
		}, nil
	}
	spell, ok := resolvable.(gob.Spell)
	if !ok {
		return EffectResult{}, errors.New("choice is not a spell")
	}
	var events []event.GameEvent
	if spell.Flashback() {
		events = append(events, event.PutSpellInExileEvent{
			PlayerID: player.ID(),
			SpellID:  spell.ID(),
		})
	} else {
		events = append(events, event.PutSpellInGraveyardEvent{
			PlayerID: player.ID(),
			SpellID:  spell.ID(),
		})
	}
	return EffectResult{
		Events: events,
	}, nil
}

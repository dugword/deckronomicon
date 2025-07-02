package resolver

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"errors"
)

func ResolveCounterspell(
	game *state.Game,
	playerID string,
	counterspell *effect.Counterspell,
	target target.Target,
) (Result, error) {
	resolvable, ok := game.Stack().Find(has.ID(target.ID))
	if !ok {
		return Result{
			Events: []event.GameEvent{
				&event.SpellOrAbilityFizzlesEvent{
					PlayerID: playerID,
					ObjectID: target.ID,
				},
			},
		}, nil
	}
	spell, ok := resolvable.(*gob.Spell)
	if !ok {
		return Result{}, errors.New("choice is not a spell")
	}
	events := []event.GameEvent{
		&event.RemoveSpellOrAbilityFromStackEvent{
			PlayerID: spell.Owner(),
			ObjectID: spell.ID(),
		},
	}
	return Result{
		Events: events,
	}, nil
}

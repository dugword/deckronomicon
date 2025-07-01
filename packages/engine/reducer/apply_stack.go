package reducer

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"fmt"
)

// These are events that manage the priority system in the game.

func applyStackEvent(game state.Game, stackEvent event.StackEvent) (state.Game, error) {
	switch evnt := stackEvent.(type) {
	case event.ResolveTopObjectOnStackEvent:
		return game, nil
	case event.PutAbilityOnStackEvent:
		return applyPutAbilityOnStackEvent(game, evnt)
	case event.PutCopiedSpellOnStackEvent:
		return applyPutCopiedSpellOnStackEvent(game, evnt)
	case event.PutSpellOnStackEvent:
		return applyPutSpellOnStackEvent(game, evnt)
	case event.RemoveSpellOrAbilityFromStackEvent:
		return applyRemoveSpellOrAbilityFromStackEvent(game, evnt)
	default:
		return game, fmt.Errorf("unknown stack event type '%T'", evnt)
	}
}

func applyPutCopiedSpellOnStackEvent(
	game state.Game,
	evnt event.PutCopiedSpellOnStackEvent,
) (state.Game, error) {
	resolvable, ok := game.Stack().Get(evnt.SpellID)
	if !ok {
		return game, fmt.Errorf("spell %q not found on stack", evnt.SpellID)
	}
	spell, ok := resolvable.(gob.Spell)
	if !ok {
		return game, fmt.Errorf("object %q is not a spell", resolvable.ID())
	}
	game, err := game.WithPutCopiedSpellOnStack(spell, evnt.PlayerID, evnt.EffectWithTargets)
	if err != nil {
		return game, fmt.Errorf("failed to put copy of spell %q on stack: %w", evnt.SpellID, err)
	}
	return game, nil
}

func applyPutSpellOnStackEvent(
	game state.Game,
	evnt event.PutSpellOnStackEvent,
) (state.Game, error) {
	player := game.GetPlayer(evnt.PlayerID)
	card, player, ok := player.TakeCardFromZone(evnt.CardID, evnt.FromZone)
	if !ok {
		return game, fmt.Errorf("card %q not in zone %q", evnt.CardID, evnt.FromZone)
	}
	game = game.WithUpdatedPlayer(player)
	game, err := game.WithPutSpellOnStack(card, evnt.PlayerID, evnt.EffectWithTargets, evnt.Flashback)
	if err != nil {
		return game, fmt.Errorf("failed to put card %q on stack: %w", evnt.CardID, err)
	}
	return game, nil
}

func applyPutAbilityOnStackEvent(
	game state.Game,
	evnt event.PutAbilityOnStackEvent,
) (state.Game, error) {
	game, err := game.WithPutAbilityOnStack(
		evnt.PlayerID,
		evnt.SourceID,
		evnt.AbilityID,
		evnt.AbilityName,
		evnt.EffectWithTargets,
	)
	if err != nil {
		return game, fmt.Errorf("failed to put ability %q on stack: %w", evnt.AbilityID, err)
	}
	return game, nil
}

func applyRemoveSpellOrAbilityFromStackEvent(
	game state.Game,
	evnt event.RemoveSpellOrAbilityFromStackEvent,
) (state.Game, error) {
	player := game.GetPlayer(evnt.PlayerID)
	resolvable, stack, ok := game.Stack().Take(evnt.ObjectID)
	if !ok {
		return game, fmt.Errorf("object %q not found on stack", evnt.ObjectID)
	}
	switch obj := resolvable.(type) {
	case gob.Spell:
		if obj.IsCopy() {
			break
		}
		zone := mtg.ZoneGraveyard
		if obj.Flashback() {
			zone = mtg.ZoneExile
		}
		if obj.Match(has.Subtype(mtg.SubtypeOmen)) {
			zone = mtg.ZoneLibrary
		}
		player, ok = player.WithAddCardToZone(obj.Card(), zone)
		if !ok {
			return game, fmt.Errorf("failed to move card %q to %q", obj.Card().ID(), zone)
		}
		game = game.WithUpdatedPlayer(player)
	}
	game = game.WithStack(stack)
	return game, nil
}

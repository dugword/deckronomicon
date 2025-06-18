package reducer

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
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
	case event.PutSpellOnStackEvent:
		return applyPutSpellOnStackEvent(game, evnt)
	case event.RemoveAbilityFromStackEvent:
		return applyRemoveAbilityFromStackEvent(game, evnt)
	case event.PutSpellInExileEvent:
		return applyPutSpellInZoneEvent(game, evnt.PlayerID, evnt.SpellID, mtg.ZoneExile)
	case event.PutSpellInGraveyardEvent:
		return applyPutSpellInZoneEvent(game, evnt.PlayerID, evnt.SpellID, mtg.ZoneGraveyard)
	default:
		return game, fmt.Errorf("unknown stack event type '%T'", evnt)
	}
}

func applyPutSpellOnStackEvent(
	game state.Game,
	evnt event.PutSpellOnStackEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(evnt.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", evnt.PlayerID)
	}
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

func applyPutSpellInZoneEvent(
	game state.Game,
	playerID string,
	spellID string,
	zone mtg.Zone,
) (state.Game, error) {
	player, ok := game.GetPlayer(playerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", playerID)
	}
	object, stack, ok := game.Stack().Take(spellID)
	if !ok {
		return game, fmt.Errorf("object %q not found on stack", spellID)
	}
	game = game.WithStack(stack)
	spell, ok := object.(gob.Spell)
	if !ok {
		return game, fmt.Errorf("object %q is not a spell", object.ID())
	}
	player, ok = player.WithAddCardToZone(spell.Card(), zone)
	if !ok {
		return game, fmt.Errorf("failed to move card %q to %q", spell.Card().ID(), zone)
	}
	game = game.WithUpdatedPlayer(player)
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

func applyRemoveAbilityFromStackEvent(
	game state.Game,
	evnt event.RemoveAbilityFromStackEvent,
) (state.Game, error) {
	_, stack, ok := game.Stack().Take(evnt.AbilityID)
	if !ok {
		return game, fmt.Errorf("ability %q not found on stack", evnt.AbilityID)
	}
	game = game.WithStack(stack)
	return game, nil
}

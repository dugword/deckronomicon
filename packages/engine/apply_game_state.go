package engine

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/mana"
	"deckronomicon/packages/state"
	"fmt"
)

// These are game state change events that modify the game state directly. The should be resuable and not
// specific to a single action.

func (e *Engine) applyGameStateChangeEvent(game state.Game, gameStateChangeEvent event.GameStateChangeEvent) (state.Game, error) {
	switch evnt := gameStateChangeEvent.(type) {
	case event.AddManaEvent:
		return e.applyAddManaEvent(game, evnt)
	case event.CheatEnabledEvent:
		game = game.WithCheatsEnabled(true)
		e.log.Info("Cheats enabled")
		return game, nil
	case event.DiscardCardEvent:
		return e.applyDiscardCardEvent(game, evnt)
	case event.DrawCardEvent:
		return e.applyDrawCardEvent(game, evnt)
	case event.MoveCardEvent:
		return e.applyMoveCardEvent(game, evnt)
	case event.PutAbilityOnStackEvent:
		return e.applyPutAbilityOnStackEvent(game, evnt)
	case event.PutPermanentOnBattlefieldEvent:
		return e.applyPutPermanentOnBattlefieldEvent(game, evnt)
	case event.PutSpellInGraveyardEvent:
		return e.applyPutSpellInGraveyardEvent(game, evnt)
	case event.PutSpellOnStackEvent:
		return e.applyPutSpellOnStackEvent(game, evnt)
	case event.SetActivePlayerEvent:
		game = game.WithActivePlayer(evnt.PlayerID)
		return game, nil
	case event.RemoveAbilityFromStackEvent:
		return e.applyRemoveAbilityFromStackEvent(game, evnt)
	case event.ShuffleDeckEvent:
		player, ok := game.GetPlayer(evnt.PlayerID)
		if !ok {
			return game, fmt.Errorf("player %q not found", evnt.PlayerID)
		}
		player = player.WithShuffleDeck(e.rng.DeckShuffler())
		game := game.WithUpdatedPlayer(player)
		return game, nil
	case event.SpendManaEvent:
		player, ok := game.GetPlayer(evnt.PlayerID)
		if !ok {
			return game, fmt.Errorf("player %q not found", evnt.PlayerID)
		}
		amount, err := mana.ParseManaString(evnt.ManaString)
		if err != nil {
			return game, fmt.Errorf("failed to parse mana string %q: %w", evnt.ManaString, err)
		}
		manaPool, err := player.ManaPool().WithSpentFromManaAmount(amount)
		if err != nil {
			return game, fmt.Errorf("failed to update mana pool for player %q: %w", player.ID(), err)
		}
		player = player.WithManaPool(manaPool)
		game = game.WithUpdatedPlayer(player)
		return game, nil

	case event.TapPermanentEvent:
		return e.applyTapPermanentEvent(game, evnt)
	case event.UntapPermanentEvent:
		return e.applyUntapPermanentEvent(game, evnt)
	default:
		return game, fmt.Errorf("unknown game state change event type '%T'", evnt)
	}
}

func (e *Engine) applyAddManaEvent(
	game state.Game,
	addManaEvent event.AddManaEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(addManaEvent.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", addManaEvent.PlayerID)
	}
	player = player.WithAddMana(addManaEvent.ManaType, addManaEvent.Amount)
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

func (e *Engine) applyDiscardCardEvent(
	game state.Game,
	event event.DiscardCardEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(event.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", event.PlayerID)
	}
	player, err := player.WithDiscardCard(event.CardID)
	if err != nil {
		return game, fmt.Errorf("failed to discard card %q: %w", event.CardID, err)
	}
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

func (e *Engine) applyDrawCardEvent(
	game state.Game,
	event event.DrawCardEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(event.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", event.PlayerID)
	}
	player, _, err := player.WithDrawCard()
	if err != nil {
		return game, fmt.Errorf("failed to draw card for %q: %w", player.ID(), err)
	}
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

func (e *Engine) applyMoveCardEvent(
	game state.Game,
	evnt event.MoveCardEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(evnt.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", evnt.PlayerID)
	}
	card, player, ok := player.TakeCardFromZone(evnt.CardID, evnt.FromZone)
	if !ok {
		return game, fmt.Errorf("card %q not in zone %q", evnt.CardID, evnt.FromZone)
	}
	player, ok = player.WithAddCardToZone(card, evnt.ToZone)
	if !ok {
		return game, fmt.Errorf("failed to move card %q to zone %q", evnt.CardID, evnt.ToZone)
	}
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

func (e *Engine) applyPutPermanentOnBattlefieldEvent(
	game state.Game,
	evnt event.PutPermanentOnBattlefieldEvent,
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
	game, err := game.WithPutPermanentOnBattlefield(card, evnt.PlayerID)
	if err != nil {
		return game, fmt.Errorf("failed to put card %q on battlefield: %w", card.ID(), err)
	}
	return game, nil
}

func (e *Engine) applyPutSpellInGraveyardEvent(
	game state.Game,
	evnt event.PutSpellInGraveyardEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(evnt.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", evnt.PlayerID)
	}
	object, stack, ok := game.Stack().Take(evnt.SpellID)
	if !ok {
		return game, fmt.Errorf("object %q not found on stack", evnt.SpellID)
	}
	game = game.WithStack(stack)
	spell, ok := object.(gob.Spell)
	if !ok {
		return game, fmt.Errorf("object %q is not a spell", object.ID())
	}
	player, ok = player.WithAddCardToZone(spell.Card(), mtg.ZoneGraveyard)
	if !ok {
		return game, fmt.Errorf("failed to move card %q to graveyard", spell.Card().ID())
	}
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

func (e *Engine) applyPutSpellOnStackEvent(
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
	game, err := game.WithPutSpellOnStack(card, evnt.PlayerID)
	if err != nil {
		return game, fmt.Errorf("failed to put card %q on stack: %w", evnt.CardID, err)
	}
	return game, nil
}

func (e *Engine) applyPutAbilityOnStackEvent(
	game state.Game,
	evnt event.PutAbilityOnStackEvent,
) (state.Game, error) {
	game, err := game.WithPutAbilityOnStack(
		evnt.PlayerID,
		evnt.SourceID,
		evnt.AbilityID,
		evnt.AbilityName,
		evnt.Effects,
	)
	if err != nil {
		return game, fmt.Errorf("failed to put ability %q on stack: %w", evnt.AbilityID, err)
	}
	return game, nil
}

func (e *Engine) applyRemoveAbilityFromStackEvent(
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

func (e *Engine) applyTapPermanentEvent(
	game state.Game,
	evnt event.TapPermanentEvent,
) (state.Game, error) {
	permanent, ok := game.Battlefield().Get(evnt.PermanentID)
	if !ok {
		return game, fmt.Errorf("permanent %q not found", evnt.PermanentID)
	}
	permanent, err := permanent.Tap()
	if err != nil {
		return game, fmt.Errorf("failed to tap permanent %q: %w", evnt.PermanentID, err)
	}
	battlefield := game.Battlefield().WithUpdatedPermanent(permanent)
	game = game.WithBattlefield(battlefield)
	return game, nil
}

func (e *Engine) applyUntapPermanentEvent(
	game state.Game,
	evnt event.UntapPermanentEvent,
) (state.Game, error) {
	permanent, ok := game.Battlefield().Get(evnt.PermanentID)
	if !ok {
		return game, fmt.Errorf("permanent %q not found", evnt.PermanentID)
	}
	permanent = permanent.Untap()
	battlefield := game.Battlefield().WithUpdatedPermanent(permanent)
	game = game.WithBattlefield(battlefield)
	return game, nil
}

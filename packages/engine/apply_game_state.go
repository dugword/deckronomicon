package engine

import (
	"deckronomicon/packages/engine/event"
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
	case event.PutCardOnBattlefieldEvent:
		return e.applyPutCardOnBattlefieldEvent(game, evnt)
	case event.PutCardOnStackEvent:
		return e.applyPutCardOnStackEvent(game, evnt)
	case event.SetActivePlayerEvent:
		game = game.WithActivePlayer(evnt.PlayerID)
		return game, nil
	case event.ShuffleDeckEvent:
		player, ok := game.GetPlayer(evnt.PlayerID)
		if !ok {
			return game, fmt.Errorf("player %q not found", evnt.PlayerID)
		}
		player = player.WithShuffleDeck(e.rng.DeckShuffler())
		game := game.WithUpdatedPlayer(player)
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

func (e *Engine) applyPutCardOnBattlefieldEvent(
	game state.Game,
	evnt event.PutCardOnBattlefieldEvent,
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
	game, err := game.WithPutCardOnBattlefield(card, evnt.PlayerID)
	if err != nil {
		return game, fmt.Errorf("failed to put card %q on battlefield: %w", card.ID(), err)
	}
	return game, nil
}

func (e *Engine) applyPutCardOnStackEvent(
	game state.Game,
	evnt event.PutCardOnStackEvent,
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
	game, err := game.WithPutCardOnStack(card, evnt.PlayerID)
	if err != nil {
		return game, fmt.Errorf("failed to put card %q on stack: %w", evnt.CardID, err)
	}
	game = game.WithUpdatedPlayer(player)
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

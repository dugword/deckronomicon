package reducer

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"fmt"
)

// These are game state change events that modify the game state directly. The should be resuable and not
// specific to a single action.

func applyGameStateChangeEvent(game state.Game, gameStateChangeEvent event.GameStateChangeEvent) (state.Game, error) {
	switch evnt := gameStateChangeEvent.(type) {
	case event.AddManaEvent:
		return applyAddManaEvent(game, evnt)
	case event.CheatEnabledEvent:
		game = game.WithCheatsEnabled(true)
		return game, nil
	case event.DiscardCardEvent:
		return applyDiscardCardEvent(game, evnt)
	case event.DrawCardEvent:
		return applyDrawCardEvent(game, evnt)
	case event.GainLifeEvent:
		return applyGainLifeEvent(game, evnt)
	case event.LoseLifeEvent:
		return applyeLoseLifeEvent(game, evnt)
	//case event.MoveCardEvent:
	//return applyMoveCardEvent(game, evnt)
	case event.PutCardInHandEvent:
		return applyPutCardInHandEvent(game, evnt)
	case event.PutCardInGraveyardEvent:
		return applyPutCardInGraveyardEvent(game, evnt)
	case event.PutCardOnBottomOfLibraryEvent:
		return applyPutCardOnBottomOfLibraryEvent(game, evnt)
	case event.PutCardOnTopOfLibraryEvent:
		return applyPutCardOnTopOfLibraryEvent(game, evnt)
	case event.PutPermanentOnBattlefieldEvent:
		return applyPutPermanentOnBattlefieldEvent(game, evnt)
	case event.RevealCardEvent:
		player := game.GetPlayer(evnt.PlayerID)
		card, ok := player.GetCardFromZone(evnt.CardID, evnt.FromZone)
		if !ok {
			return game, fmt.Errorf("card %q not found in zone %q", evnt.CardID, evnt.FromZone)
		}
		revealed := player.Revealed().Add(card)
		player = player.WithRevealed(revealed)
		game = game.WithUpdatedPlayer(player)
		return game, nil
	case event.SetActivePlayerEvent:
		game = game.WithActivePlayer(evnt.PlayerID)
		return game, nil
	case event.ShuffleLibraryEvent:
		return applyShuffleLibraryEvent(game, evnt)
	case event.SpendManaEvent:
		player := game.GetPlayer(evnt.PlayerID)
		amount, err := mana.ParseManaString(evnt.ManaString)
		if err != nil {
			return game, fmt.Errorf("failed to parse mana string %q: %w", evnt.ManaString, err)
		}
		// TODO: Let the user choose which colors to spend by providing a list of colors.
		manaPool, err := player.ManaPool().WithSpendFromManaAmount(amount, mana.Colors())
		if err != nil {
			return game, fmt.Errorf("failed to update mana pool for player %q: %w", player.ID(), err)
		}
		player = player.WithManaPool(manaPool)
		game = game.WithUpdatedPlayer(player)
		return game, nil

	case event.TapPermanentEvent:
		return applyTapPermanentEvent(game, evnt)
	case event.UntapPermanentEvent:
		return applyUntapPermanentEvent(game, evnt)
	default:
		return game, fmt.Errorf("unknown game state change event type '%T'", evnt)
	}
}

func applyAddManaEvent(
	game state.Game,
	addManaEvent event.AddManaEvent,
) (state.Game, error) {
	player := game.GetPlayer(addManaEvent.PlayerID)
	player = player.WithAddMana(addManaEvent.Color, addManaEvent.Amount)
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

func applyDiscardCardEvent(
	game state.Game,
	event event.DiscardCardEvent,
) (state.Game, error) {
	player := game.GetPlayer(event.PlayerID)
	player, err := player.WithDiscardCard(event.CardID)
	if err != nil {
		return game, fmt.Errorf("failed to discard card %q: %w", event.CardID, err)
	}
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

func applyDrawCardEvent(
	game state.Game,
	event event.DrawCardEvent,
) (state.Game, error) {
	player := game.GetPlayer(event.PlayerID)
	player, _, err := player.WithDrawCard()
	if err != nil {
		if err == mtg.ErrLibraryEmpty {
			return game, mtg.PlayerLostError{
				Reason: mtg.DeckedOut,
			}
		}
		return game, fmt.Errorf("failed to draw card for %q: %w", player.ID(), err)
	}
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

func applyGainLifeEvent(
	game state.Game,
	evnt event.GainLifeEvent,
) (state.Game, error) {
	player := game.GetPlayer(evnt.PlayerID)
	player = player.WithGainLife(evnt.Amount)
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

func applyeLoseLifeEvent(
	game state.Game,
	evnt event.LoseLifeEvent,
) (state.Game, error) {
	player := game.GetPlayer(evnt.PlayerID)
	player = player.WithLoseLife(evnt.Amount)
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

func applyPutCardInHandEvent(
	game state.Game,
	evnt event.PutCardInHandEvent,
) (state.Game, error) {
	player := game.GetPlayer(evnt.PlayerID)
	card, player, ok := player.TakeCardFromZone(evnt.CardID, evnt.FromZone)
	if !ok {
		return game, fmt.Errorf("card %q not in zone %q", evnt.CardID, evnt.FromZone)
	}
	player, ok = player.WithAddCardToZone(card, mtg.ZoneHand)
	if !ok {
		return game, fmt.Errorf("failed to move card %q to hand", evnt.CardID)
	}
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

func applyPutCardInGraveyardEvent(
	game state.Game,
	evnt event.PutCardInGraveyardEvent,
) (state.Game, error) {
	player := game.GetPlayer(evnt.PlayerID)
	card, player, ok := player.TakeCardFromZone(evnt.CardID, evnt.FromZone)
	if !ok {
		return game, fmt.Errorf("card %q not in zone %q", evnt.CardID, evnt.FromZone)
	}
	player, ok = player.WithAddCardToZone(card, mtg.ZoneGraveyard)
	if !ok {
		return game, fmt.Errorf("failed to move card %q to graveyard", evnt.CardID)
	}
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

func applyPutCardOnBottomOfLibraryEvent(
	game state.Game,
	evnt event.PutCardOnBottomOfLibraryEvent,
) (state.Game, error) {
	player := game.GetPlayer(evnt.PlayerID)
	card, player, ok := player.TakeCardFromZone(evnt.CardID, evnt.FromZone)
	if !ok {
		return game, fmt.Errorf("card %q not in zone %q", evnt.CardID, evnt.FromZone)
	}
	player, ok = player.WithAddCardToZone(card, mtg.ZoneLibrary)
	if !ok {
		return game, fmt.Errorf("failed to move card %q to library", evnt.CardID)
	}
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

func applyPutCardOnTopOfLibraryEvent(
	game state.Game,
	evnt event.PutCardOnTopOfLibraryEvent,
) (state.Game, error) {
	player := game.GetPlayer(evnt.PlayerID)
	card, player, ok := player.TakeCardFromZone(evnt.CardID, evnt.FromZone)
	if !ok {
		return game, fmt.Errorf("card %q not in zone %q", evnt.CardID, evnt.FromZone)
	}
	player, ok = player.WithAddCardToTopOfZone(card, mtg.ZoneLibrary)
	if !ok {
		return game, fmt.Errorf("failed to move card %q to library", evnt.CardID)
	}
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

func applyPutPermanentOnBattlefieldEvent(
	game state.Game,
	evnt event.PutPermanentOnBattlefieldEvent,
) (state.Game, error) {
	player := game.GetPlayer(evnt.PlayerID)
	switch evnt.FromZone {
	case mtg.ZoneHand, mtg.ZoneLibrary, mtg.ZoneGraveyard, mtg.ZoneExile:
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
	// TODO: This probably needs to be handled differently.
	// I think I want cards that are permanents to create spells
	// that have a "put on battlefield" effect, and then
	// this event should be used to resolve that effect.
	case mtg.ZoneStack:
		resolvable, stack, ok := game.Stack().Take(evnt.CardID)
		if !ok {
			return game, fmt.Errorf("card %q not found in stack", evnt.CardID)
		}
		spell, ok := resolvable.(gob.Spell)
		if !ok {
			return game, fmt.Errorf("resolvable %q is not a spell", evnt.CardID)
		}
		card := spell.Card()
		game = game.WithStack(stack)
		game, err := game.WithPutPermanentOnBattlefield(card, evnt.PlayerID)
		if err != nil {
			return game, fmt.Errorf("failed to put card %q on battlefield: %w", card.ID(), err)
		}
		return game, nil
	default:
		return game, fmt.Errorf("invalid zone %q for putting permanent on battlefield", evnt.FromZone)
	}

}

func applyShuffleLibraryEvent(
	game state.Game,
	evnt event.ShuffleLibraryEvent,
) (state.Game, error) {
	player := game.GetPlayer(evnt.PlayerID)
	player, err := player.WithShuffledLibrary(evnt.ShuffledCardsIDs)
	if err != nil {
		return game, fmt.Errorf(
			"failed to shuffle library for player %q: %w",
			evnt.PlayerID, err,
		)
	}
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

func applyTapPermanentEvent(
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

func applyUntapPermanentEvent(
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

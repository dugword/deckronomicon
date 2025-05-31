package engine

import (
	"deckronomicon/packages/engine/effect"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/target"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
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
	case event.GainLifeEvent:
		return e.applyGainLifeEvent(game, evnt)
	case event.LoseLifeEvent:
		return e.applyeLoseLifeEvent(game, evnt)
	//case event.MoveCardEvent:
	//return e.applyMoveCardEvent(game, evnt)
	case event.PutCardInHandEvent:
		return e.applyPutCardInHandEvent(game, evnt)
	case event.PutCardInGraveyardEvent:
		return e.applyPutCardInGraveyardEvent(game, evnt)
	case event.PutCardOnBottomOfLibraryEvent:
		return e.applyPutCardOnBottomOfLibraryEvent(game, evnt)
	case event.PutCardOnTopOfLibraryEvent:
		return e.applyPutCardOnTopOfLibraryEvent(game, evnt)
	case event.PutPermanentOnBattlefieldEvent:
		return e.applyPutPermanentOnBattlefieldEvent(game, evnt)
	case event.RevealCardEvent:
		player, ok := game.GetPlayer(evnt.PlayerID)
		if !ok {
			return game, fmt.Errorf("player %q not found", evnt.PlayerID)
		}
		card, ok := player.GetCardFromZone(evnt.CardID, evnt.FromZone)
		if !ok {
			return game, fmt.Errorf("card %q not found in zone %q", evnt.CardID, evnt.FromZone)
		}
		revealed := player.Revealed().Add(card)
		player = player.WithRevealed(revealed)
		game = game.WithUpdatedPlayer(player)
		e.log.Infof("Player %q reveals card %q from %q", player.ID(), card.Name(), evnt.FromZone)
		return game, nil
	case event.SetActivePlayerEvent:
		game = game.WithActivePlayer(evnt.PlayerID)
		return game, nil
	case event.ResolveManaAbilityEvent:
		return e.applyResolveManaAbilityEvent(game, evnt)
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

func (e *Engine) applyGainLifeEvent(
	game state.Game,
	evnt event.GainLifeEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(evnt.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", evnt.PlayerID)
	}
	player = player.WithGainLife(evnt.Amount)
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

func (e *Engine) applyeLoseLifeEvent(
	game state.Game,
	evnt event.LoseLifeEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(evnt.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", evnt.PlayerID)
	}
	player = player.WithLoseLife(evnt.Amount)
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

/*
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
*/

func (e *Engine) applyPutCardInHandEvent(
	game state.Game,
	evnt event.PutCardInHandEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(evnt.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", evnt.PlayerID)
	}
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

func (e *Engine) applyPutCardInGraveyardEvent(
	game state.Game,
	evnt event.PutCardInGraveyardEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(evnt.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", evnt.PlayerID)
	}
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

func (e *Engine) applyPutCardOnBottomOfLibraryEvent(
	game state.Game,
	evnt event.PutCardOnBottomOfLibraryEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(evnt.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", evnt.PlayerID)
	}
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

func (e *Engine) applyPutCardOnTopOfLibraryEvent(
	game state.Game,
	evnt event.PutCardOnTopOfLibraryEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(evnt.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", evnt.PlayerID)
	}
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

func (e *Engine) applyPutPermanentOnBattlefieldEvent(
	game state.Game,
	evnt event.PutPermanentOnBattlefieldEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(evnt.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", evnt.PlayerID)
	}
	if evnt.FromZone == mtg.ZoneStack {
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

// TODO: Not sure I like how this is handled, think through how I want to have
// events generated in the reducer.
func (e *Engine) applyResolveManaAbilityEvent(
	game state.Game,
	evnt event.ResolveManaAbilityEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(evnt.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", evnt.PlayerID)
	}
	var events []event.GameEvent
	for _, effectSpec := range evnt.EffectSpecs {
		efct, err := effect.Build(effectSpec)
		if err != nil {
			return game, fmt.Errorf("effect %q not found: %w", effectSpec.Name, err)
		}
		effectResults, err := efct.Resolve(game, player, nil, target.TargetValue{})
		if err != nil {
			return game, fmt.Errorf("failed to apply effect %q: %w", effectSpec.Name, err)
		}
		events = append(events, effectResults.Events...)
	}
	for _, evnt := range events {
		var err error
		game, err = e.applyEvent(game, evnt)
		if err != nil {
			return game, fmt.Errorf("failed to apply event %T: %w", evnt, err)
		}
	}
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

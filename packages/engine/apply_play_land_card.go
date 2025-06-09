package engine

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/judge"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

// TODO: Need to document at what level stuff should happen.
// maybe something with if cards move from player zone to player zone, vs if a
// card moves from player to battlefield.... like a method should only operate
// on it's own fields? But wht about sub fields. Probably no, because
// otherwise everything would happen on game.
func (e *Engine) applyPlayLandEvent(
	game state.Game,
	evnt event.PlayLandEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(evnt.PlayerID)
	if !ok {
		return game, fmt.Errorf("player '%s' not found", evnt.PlayerID)
	}
	if evnt.Zone != mtg.ZoneHand {
		return game, errors.New("card not played from hand")
	}
	card, player, err := player.TakeCardFromHand(evnt.CardID)
	if err != nil {
		return game, fmt.Errorf("failed to take card '%s' from hand: %w", evnt.CardID, err)
	}
	if !judge.CanPlayCard(game, player, mtg.ZoneHand, card) {
		return game, errors.New("card cannot be played")
	}
	newGame, err := game.WithPutCardOnBattlefield(card, player.ID())
	if err != nil {
		return game, fmt.Errorf("failed to put card '%s' on battlefield: %w", evnt.CardID, err)
	}
	newerGame := newGame.WithUpdatedPlayer(player.WithLandPlayedThisTurn())
	return newerGame, nil
}

package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/judge"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"fmt"
)

type PlayLandAction struct {
	player          state.Player
	cardInZone      gob.CardInZone
	ManaPayment     mana.Pool
	AdditionalCosts []string
	Targets         []string
}

func NewPlayLandAction(player state.Player, cardInZone gob.CardInZone) PlayLandAction {
	return PlayLandAction{
		player:     player,
		cardInZone: cardInZone,
	}
}

func (a PlayLandAction) PlayerID() string {
	return a.player.ID()
}

func (a PlayLandAction) Name() string {
	return "Play Land"
}

func (a PlayLandAction) Description() string {
	return "The active player plays a land from their hand."
}

func (a PlayLandAction) Complete(game state.Game, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	// TODO: Maybe this should happen in judge.CanPlayLand?
	if !a.player.ZoneContains(a.cardInZone.Zone(), has.ID(a.cardInZone.ID())) {
		return nil, fmt.Errorf(
			"player %q does not have card %q in zone %q",
			a.player.ID(),
			a.cardInZone.ID(),
			a.cardInZone.Zone(),
		)
	}
	ruling := judge.Ruling{Explain: true}
	if !judge.CanPlayLand(game, a.player, a.cardInZone.Zone(), a.cardInZone.Card(), &ruling) {
		return nil, fmt.Errorf(
			"player %q cannot play card %q, %s",
			a.player.ID(),
			a.cardInZone.ID(),
			ruling.Why(),
		)
	}
	if a.cardInZone.Card().Match(is.Land()) {
		return []event.GameEvent{
			event.PlayLandEvent{
				PlayerID: a.player.ID(),
				CardID:   a.cardInZone.ID(),
				Zone:     a.cardInZone.Zone(),
			},
			event.PutPermanentOnBattlefieldEvent{
				PlayerID: a.player.ID(),
				CardID:   a.cardInZone.ID(),
				FromZone: a.cardInZone.Zone(),
			},
		}, nil
	}
	return nil, fmt.Errorf(
		"card %q is not a land, cannot play",
		a.cardInZone.ID(),
	)
}

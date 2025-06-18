package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/judge"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"fmt"
)

type PlayLandRequest struct {
	CardID string
}

type PlayLandAction struct {
	playerID string
	cardID   string
}

func NewPlayLandAction(player state.Player, request PlayLandRequest) PlayLandAction {
	return PlayLandAction{
		playerID: player.ID(),
		cardID:   request.CardID,
	}
}

func (a PlayLandAction) Name() string {
	return "Play Land"
}

func (a PlayLandAction) Complete(game state.Game, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	player, ok := game.GetPlayer(a.playerID)
	if !ok {
		return nil, fmt.Errorf("player %q not found in game", a.playerID)
	}
	landToPlay, ok := player.GetCardFromZone(a.cardID, mtg.ZoneHand)
	if !ok {
		return nil, fmt.Errorf("player %q does not have card %q in hand", a.playerID, a.cardID)
	}
	ruling := judge.Ruling{Explain: true}
	if !judge.CanPlayLand(game, player, mtg.ZoneHand, landToPlay, &ruling) {
		return nil, fmt.Errorf(
			"player %q cannot play land %q, %s",
			a.playerID,
			landToPlay.ID(),
			ruling.Why(),
		)
	}
	return []event.GameEvent{
		event.PlayLandEvent{
			PlayerID: a.playerID,
			CardID:   landToPlay.ID(),
			Zone:     mtg.ZoneHand,
		},
		event.PutPermanentOnBattlefieldEvent{
			PlayerID: a.playerID,
			CardID:   landToPlay.ID(),
			FromZone: mtg.ZoneHand,
		},
	}, nil
}

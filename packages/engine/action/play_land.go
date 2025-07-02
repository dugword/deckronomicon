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
	cardID string
}

func NewPlayLandAction(request PlayLandRequest) PlayLandAction {
	return PlayLandAction{
		cardID: request.CardID,
	}
}

func (a PlayLandAction) Name() string {
	return "Play Land"
}

func (a PlayLandAction) Complete(game *state.Game, playerID string, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	player := game.GetPlayer(playerID)
	landToPlay, ok := player.GetCardFromZone(a.cardID, mtg.ZoneHand)
	if !ok {
		return nil, fmt.Errorf("player %q does not have card %q in hand", playerID, a.cardID)
	}
	ruling := judge.Ruling{Explain: true}
	if !judge.CanPlayLand(game, playerID, mtg.ZoneHand, landToPlay, &ruling) {
		return nil, fmt.Errorf(
			"player %q cannot play land %q, %s",
			playerID,
			landToPlay.ID(),
			ruling.Why(),
		)
	}
	return []event.GameEvent{
		&event.PlayLandEvent{
			PlayerID: playerID,
			CardID:   landToPlay.ID(),
			Zone:     mtg.ZoneHand,
		},
		&event.PutPermanentOnBattlefieldEvent{
			PlayerID: playerID,
			CardID:   landToPlay.ID(),
			FromZone: mtg.ZoneHand,
		},
	}, nil
}

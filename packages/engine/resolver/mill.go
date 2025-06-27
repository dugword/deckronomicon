package resolver

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
)

func ResolveMill(
	game state.Game,
	playerID string,
	mill effect.Mill,
	target effect.Target,
) (Result, error) {
	targetPlayerID := playerID
	if target.ID != "" {
		targetPlayerID = target.ID
	}
	targetPlayer := game.GetPlayer(targetPlayerID)
	cards := targetPlayer.Library().GetN(mill.Count)
	var events []event.GameEvent
	for _, card := range cards {
		events = append(events, event.PutCardInGraveyardEvent{
			PlayerID: targetPlayerID,
			CardID:   card.ID(),
			FromZone: mtg.ZoneLibrary,
		})
	}
	return Result{Events: events}, nil
}

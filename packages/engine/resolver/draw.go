package resolver

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/state"
	"fmt"
)

func ResolveDraw(
	game state.Game,
	playerID string,
	draw effect.Draw,
	target target.Target,
) (Result, error) {
	switch target.Type {
	case mtg.TargetTypeNone:
		return resolveDrawForPlayer(game, playerID, draw.Count)
	case mtg.TargetTypePlayer:
		return resolveDrawForPlayer(game, target.ID, draw.Count)
	default:
		panic(fmt.Sprintf("unexpected target type %s for DrawEffect", target.Type))
		return Result{}, nil
	}
}

func resolveDrawForPlayer(game state.Game, playerID string, count int) (Result, error) {
	var events []event.GameEvent
	for range count {
		events = append(events, event.DrawCardEvent{
			PlayerID: playerID,
		})
	}
	return Result{
		Events: events,
	}, nil
}

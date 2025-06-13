package effect

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"strconv"
)

func DrawEffectHandler(
	game state.Game,
	player state.Player,
	source query.Object,
	modifiers []definition.EffectModifier,
) (EffectResult, error) {
	var events []event.GameEvent
	drawCount := 1
	for _, modifier := range modifiers {
		if modifier.Key == "Type" && modifier.Value == "Cycle" {
			events = append(events, event.CycleCardEvent{
				PlayerID: player.ID(),
			})
		}
		if modifier.Key == "Count" {
			var err error
			drawCount, err = strconv.Atoi(modifier.Value)
			if err != nil {
				return EffectResult{}, err
			}
		}
	}
	for range drawCount {
		events = append(events, event.DrawCardEvent{
			PlayerID: player.ID(),
		})
	}
	return EffectResult{
		Events: events,
	}, nil
}

package reducer

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
)

// These are events that are triggered by the game engine at specific points in the game lifecycle.
// They are not triggered by player actions, but rather by the game state itself. These events are used to
// manage the game state, such as untapping permanents, progressing sagas, or checking day/night status.

func applyTurnBasedActionEvent(game *state.Game, turnBasedActionEvent event.TurnBasedActionEvent) (*state.Game, error) {
	switch evnt := turnBasedActionEvent.(type) {
	case *event.CheckDayNightEvent:
		return game, nil
	case *event.DiscardToHandSizeEvent:
		return game, nil
	case *event.DrawStartingHandEvent:
		return game, nil
	case *event.PhaseInPhaseOutEvent:
		return game, nil
	case *event.ProgressSagaEvent:
		return game, nil
	case *event.RemoveDamageEvent:
		return game, nil
	case *event.UntapAllEvent:
		return applyUntapAllEvent(game, evnt)
	case *event.UpkeepEvent:
		return game, nil
	default:
		return game, fmt.Errorf("unknown turn-based action event type '%T'", evnt)
	}
}

func applyUntapAllEvent(
	game *state.Game,
	evnt *event.UntapAllEvent,
) (*state.Game, error) {
	battlefield := game.Battlefield().UntapAll(evnt.PlayerID)
	game = game.WithBattlefield(battlefield)
	return game, nil
}

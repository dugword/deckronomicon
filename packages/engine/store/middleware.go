package store

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type RunResult struct {
	RunID       int              `json:"RunID"`
	Turns       int              `json:"Turns"`
	Totals      map[string]int   `json:"Totals"`
	Cumulatives map[string][]int `json:"Cumulatives"`
}

func AnalyticsMiddleware(runResult *RunResult) Middleware {
	playerID := "You"
	runResult.Totals["CardDrawn"] = 0
	return func(next ApplyEventFunc) ApplyEventFunc {
		return func(game *state.Game, ev event.GameEvent) (*state.Game, error) {
			game, err := next(game, ev)
			if err != nil {
				return nil, err
			}
			switch myEVent := ev.(type) {
			case *event.DrawCardEvent:
				if myEVent.PlayerID != playerID {
					return game, nil
				}
				runResult.Totals["CardDrawn"]++
				hand := game.GetPlayer(myEVent.PlayerID).Hand().GetAll()
				cardName := hand[len(hand)-1].Name()
				if cardName == "" {

				}
			case *event.EndTurnEvent:
				if myEVent.PlayerID != playerID {
					return game, nil
				}
				runResult.Turns++
				runResult.Cumulatives["CardDrawn"] = append(
					runResult.Cumulatives["CardDrawn"],
					runResult.Totals["CardDrawn"],
				)
			}
			return game, nil
		}
	}
}

func LoggingMiddleware(log Logger) Middleware {
	return func(next ApplyEventFunc) ApplyEventFunc {
		return func(game *state.Game, ev event.GameEvent) (*state.Game, error) {
			log.Debug("Applying event:", ev.EventType())
			return next(game, ev)
		}
	}
}

func RecordEventMiddleware(recordFunc func(event event.GameEvent)) Middleware {
	return func(next ApplyEventFunc) ApplyEventFunc {
		return func(game *state.Game, ev event.GameEvent) (*state.Game, error) {
			newGame, err := next(game, ev)
			if err == nil {
				recordFunc(ev)
			}
			return newGame, err
		}
	}
}

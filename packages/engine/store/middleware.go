package store

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

func AnalyticsMiddleware(results map[string]int) Middleware {
	var totalTimesSeen int
	results["CardDrawn"] = 0
	results["FlushOutDrawn"] = 0
	results["ThrillOfPossibilityDrawn"] = 0
	results["PearledUnicornDrawn"] = 0
	return func(next ApplyEventFunc) ApplyEventFunc {
		return func(game *state.Game, ev event.GameEvent) (*state.Game, error) {
			game, err := next(game, ev)
			if err != nil {
				return nil, err
			}
			switch myEVent := ev.(type) {
			case *event.DrawCardEvent:
				results["CardDrawn"]++
				hand := game.GetPlayer(myEVent.PlayerID).Hand().GetAll()
				card := hand[len(hand)-1].Name()
				if card == "Flush Out" {
					totalTimesSeen++
					results["FlushOutDrawn"]++
				}
				if card == "Thrill of Possibility" {
					results["ThrillOfPossibilityDrawn"]++
				}
				if card == "Pearled Unicorn" {
					results["PearledUnicornDrawn"]++
				}
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

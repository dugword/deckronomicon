package store

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

func AnalyticsMiddleware() Middleware {
	return func(next ApplyEventFunc) ApplyEventFunc {
		return func(game *state.Game, ev event.GameEvent) (*state.Game, error) {
			game, err := next(game, ev)
			if err != nil {
				return nil, err
			}
			switch myEVent := ev.(type) {
			case *event.DrawCardEvent:
				hand := game.GetPlayer(myEVent.PlayerID).Hand().GetAll()
				card := hand[len(hand)-1].Name()
				if card == "Flush Out" {
					// fmt.Println("Flush Out drawn!")
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

package store

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func AnalyticsMiddleware() Middleware {
	var totalTimesSeen int
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
					totalTimesSeen++
					filename := filepath.Join("results/", fmt.Sprintf("%s_%s", game.RunID(), "flush_out_seen.txt"))
					f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						log.Fatal(err)
					}
					message := fmt.Sprintf("Flush Out seen %d times in run %s\n", totalTimesSeen, game.RunID())
					if _, err := f.Write([]byte(message)); err != nil {
						f.Close() // ignore error; Write error takes precedence
						log.Fatal(err)
					}
					if err := f.Close(); err != nil {
						log.Fatal(err)
					}
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

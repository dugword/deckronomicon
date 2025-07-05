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
	var totalTimesSeen int
	playerID := "You"
	runResult.Totals["CardDrawn"] = 0
	runResult.Totals["FlushOutDrawn"] = 0
	runResult.Totals["ThrillOfPossibilityDrawn"] = 0
	runResult.Totals["PearledUnicornDrawn"] = 0
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
				card := hand[len(hand)-1].Name()
				if card == "Flush Out" {
					totalTimesSeen++
					runResult.Totals["FlushOutDrawn"]++
				}
				if card == "Thrill of Possibility" {
					runResult.Totals["ThrillOfPossibilityDrawn"]++
				}
				if card == "Pearled Unicorn" {
					runResult.Totals["PearledUnicornDrawn"]++
				}
			case *event.EndTurnEvent:
				if myEVent.PlayerID != playerID {
					return game, nil
				}
				//fmt.Println("Turns =>", runResult.Turns)
				runResult.Turns++
				runResult.Cumulatives["CardDrawn"] = append(runResult.Cumulatives["CardDrawn"], runResult.Totals["CardDrawn"])
				runResult.Cumulatives["FlushOutDrawn"] = append(runResult.Cumulatives["FlushOutDrawn"], runResult.Totals["FlushOutDrawn"])
				runResult.Cumulatives["ThrillOfPossibilityDrawn"] = append(runResult.Cumulatives["ThrillOfPossibilityDrawn"], runResult.Totals["ThrillOfPossibilityDrawn"])
				runResult.Cumulatives["PearledUnicornDrawn"] = append(runResult.Cumulatives["PearledUnicornDrawn"], runResult.Totals["PearledUnicornDrawn"])
				// fmt.Println("End of turn:", myEVent.PlayerID)
				// fmt.Println("Total cumulative cards drawn:", runResult.Cumulatives["CardDrawn"])
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

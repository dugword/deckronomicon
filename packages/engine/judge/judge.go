package judge

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"strings"
)

type Ruling struct {
	Explain bool
	Reasons []string
}

func (r *Ruling) Why() string {
	if r == nil || !r.Explain {
		return ""
	}
	if len(r.Reasons) == 0 {
		return "No reasons provided."
	}
	return "reasons: " + strings.Join(r.Reasons, ", ")
}

// TODO should this live on engine?
func CanPlaySorcerySpeed(game state.Game, playerID string, ruling *Ruling) bool {
	can := true
	if !game.IsStackEmtpy() {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "stack is not empty")
		}
		can = false
	}
	if game.ActivePlayerID() != playerID {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "player is not active player")
		}
		can = false
	}
	if !(game.Step() == mtg.StepPrecombatMain ||
		game.Step() == mtg.StepPostcombatMain) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "not in a main phase")
		}
		can = false
	}
	return can
}

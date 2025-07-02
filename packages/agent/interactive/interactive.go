package interactive

// TODO: Rethink and refactor this package and how it handles Player.
// Did a bunch of quick and dirty changes to get it working with the new
// game state and player system, but it could be cleaner and more consistent.

import (
	"bufio"
	"deckronomicon/packages/agent/interactive/actionparser"
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"deckronomicon/packages/ui"
	"deckronomicon/packages/view"
	"fmt"
	"slices"
)

// InteractivePlayerAgent implements the PlayerAgent interface for interactive
// play.
type Agent struct {
	uiBuffer      *ui.Buffer
	scanner       *bufio.Scanner
	playerID      string
	inputError    string
	prompt        string
	stops         []mtg.Step
	verbose       bool
	message       []string
	autopay       bool
	autopayColors []mana.Color
}

func NewAgent(
	scanner *bufio.Scanner,
	playerID string,
	stops []mtg.Step,
	displayFile string,
	autopay bool,
	autopayColors []mana.Color,
	verbose bool,
) *Agent {
	agent := Agent{
		playerID:      playerID,
		prompt:        ">> ",
		scanner:       scanner,
		stops:         stops,
		uiBuffer:      ui.NewBuffer(displayFile),
		verbose:       verbose,
		autopay:       autopay,
		autopayColors: autopayColors,
	}
	return &agent
}

func (a *Agent) PlayerID() string {
	return a.playerID
}

func (a *Agent) ReportState(game *state.Game) {
	opponentID := game.GetOpponent(a.playerID).ID()
	a.uiBuffer.Update(
		view.NewGameViewFromState(game),
		view.NewPlayerViewFromState(game, a.playerID, ""),
		view.NewPlayerViewFromState(game, opponentID, ""),
	)
}

func (a *Agent) GetNextAction(game *state.Game) (engine.Action, error) {
	for {
		pass := true
		if slices.Contains(a.stops, game.Step()) {
			if game.ActivePlayerID() == a.playerID {
				pass = false
			}
		}
		if pass {
			return action.NewPassPriorityAction(), nil
		}
		// TODO Don't call this here, run update or something
		// PrintCommands(s.CheatsEnabled)
		// TODO: maybe move the prompt to the read functions?
		if a.inputError != "" {
			a.uiBuffer.UpdateMessage([]string{a.inputError})
			a.inputError = ""
		}
		if err := a.uiBuffer.Render(); err != nil {
			return nil, fmt.Errorf("failed to render UI Buffer: %w", err)
		}
		a.Prompt("take action")
		input := a.ReadInput()
		action, err := actionparser.ParseInput(
			input,
			a,
			game,
			a.playerID,
			a.autopay,
			a.autopayColors,
		)
		if err != nil {
			a.inputError = err.Error()
			a.uiBuffer.UpdateChoices("", []choose.Choice{})
			a.uiBuffer.UpdateMessage([]string{})
			continue
		}
		return action, nil
	}
}

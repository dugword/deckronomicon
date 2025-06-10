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
	uiBuffer   *ui.Buffer
	scanner    *bufio.Scanner
	playerID   string
	inputError string
	prompt     string
	stops      []mtg.Step
	verbose    bool
	message    []string
}

func NewAgent(
	scanner *bufio.Scanner,
	playerID string,
	stops []mtg.Step,
	displayFile string,
	verbose bool,
) *Agent {
	agent := Agent{
		playerID: playerID,
		prompt:   ">> ",
		scanner:  scanner,
		stops:    stops,
		uiBuffer: ui.NewBuffer(displayFile),
		verbose:  verbose,
	}
	return &agent
}

func (a *Agent) PlayerID() string {
	return a.playerID
}

func (a *Agent) ReportState(game state.Game) error {
	player, ok := game.GetPlayer(a.playerID)
	if !ok {
		return fmt.Errorf("player %q not found", a.playerID)
	}
	opponent, ok := game.GetOpponent(a.playerID)
	if !ok {
		return fmt.Errorf("opponent for player %q not found", a.playerID)
	}
	a.uiBuffer.Update(
		view.NewGameViewFromState(game),
		view.NewPlayerViewFromState(player),
		view.NewPlayerViewFromState(opponent),
	)
	return nil
}

func (a *Agent) GetNextAction(game state.Game) (engine.Action, error) {
	player, ok := game.GetPlayer(a.playerID)
	if !ok {
		return nil, fmt.Errorf("player %q not found", a.playerID)
	}
	for {
		pass := true
		if slices.Contains(a.stops, game.Step()) {
			if game.ActivePlayerID() == player.ID() {
				pass = false
			}
		}
		if pass {
			return action.NewPassPriorityAction(player), nil
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
		commandParser := actionparser.CommandParser{}
		command, err := commandParser.ParseInput(
			input,
			a.Choose,
			game,
			player,
		)
		if err != nil {
			a.inputError = err.Error()
			a.uiBuffer.UpdateChoices("", []choose.Choice{})
			a.uiBuffer.UpdateMessage([]string{})
			continue
		}
		action, err := command.Build(game, player)
		if err != nil {
			a.inputError = err.Error()
			continue
		}
		return action, nil
	}
}

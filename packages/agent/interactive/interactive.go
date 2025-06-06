package interactive

// TODO: Rethink and refactor this package and how it handles Player.
// Did a bunch of quick and dirty changes to get it working with the new
// game state and player system, but it could be cleaner and more consistent.

import (
	"bufio"
	"deckronomicon/packages/agent/actionparser"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"deckronomicon/packages/ui"
	"fmt"
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
	verbose bool,
) *Agent {
	agent := Agent{
		playerID: playerID,
		prompt:   ">> ",
		scanner:  scanner,
		stops:    stops,
		uiBuffer: ui.NewBuffer(),
		verbose:  verbose,
	}
	return &agent
}

// ReportState displays the current game state in a terminal UI.
// TODO maybe get player from state.CurrentPlayer?
func (a *Agent) ReportState(game state.Game) error {
	// TODO Don't panic
	player, err := game.GetPlayer(a.playerID)
	if err != nil {
		return fmt.Errorf("Player %s not found in game state", a.playerID)
	}
	// Update the UI buffer with the current game state.
	opponent, err := game.GetOpponent(a.playerID)
	if err != nil {
		return fmt.Errorf("Opponent for player %s not found in game state: %v", a.playerID, err)
	}
	a.uiBuffer.UpdateFromState(game, player, opponent)
	/// TODO: This should return an error instead of panicking
	/*
		if err := a.uiBuffer.Render(); err != nil {
			return fmt.Errorf("failed to render UI buffer: %v", err)
		}
	*/
	return nil
}

func (a *Agent) GetNextAction(game state.Game) (engine.Action, error) {
	for {
		// TODO Don't call this here, run update or something
		// PrintCommands(s.CheatsEnabled)
		// TODO: maybe move the prompt to the read functions?
		if a.inputError != "" {
			fmt.Println("Error:", a.inputError)
			a.inputError = ""
		}
		if err := a.uiBuffer.Render(); err != nil {
			return nil, fmt.Errorf("error rendering UI buffer: %w", err)
		}
		a.Prompt("take action")
		commandParser := actionparser.CommandParser{}
		command, err := commandParser.ParseInput(
			a.ReadInput,
			a.Choose,
			game,
			a.playerID,
		)
		if err != nil {
			a.inputError = err.Error()
			continue
		}
		fmt.Println("HERE =>", command)
		action, err := command.Build(game, a.playerID)
		if err != nil {
			a.inputError = err.Error()
			continue
		}
		return action, nil
	}
}

package interactive

// TODO: Rethink and refactor this package and how it handles Player.
// Did a bunch of quick and dirty changes to get it working with the new
// game state and player system, but it could be cleaner and more consistent.

import (
	"bufio"
	"fmt"

	"deckronomicon/game"
	"deckronomicon/ui"
)

// InteractivePlayerAgent implements the PlayerAgent interface for interactive
// play.
type InteractivePlayerAgent struct {
	uiBuffer   *ui.Buffer
	Scanner    *bufio.Scanner
	playerID   string
	inputError string
	prompt     string
}

// Player returns the player ID controlled by this agent.
func (a *InteractivePlayerAgent) PlayerID() string {
	return a.playerID
}

// NewInteractivePlayerAgent creates a new InteractivePlayerAgent with the
// given scanner.
func NewInteractivePlayerAgent(scanner *bufio.Scanner, playerID string) *InteractivePlayerAgent {
	return &InteractivePlayerAgent{
		playerID: playerID,
		uiBuffer: ui.NewBuffer(),
		Scanner:  scanner,
		prompt:   ">> ",
	}
}

// ReportState displays the current game state in a terminal UI.
// TODO maybe get player from state.CurrentPlayer?
func (a *InteractivePlayerAgent) ReportState(state *game.GameState) {
	a.uiBuffer.UpdateFromState(state, a.playerID)
	if err := a.uiBuffer.Render(); err != nil {
		panic(fmt.Sprintf("Error rendering UI buffer: %v", err))
	}

}

func (a *InteractivePlayerAgent) GetNextAction(state *game.GameState) *game.GameAction {
	for {
		// TODO Don't call this here, run update or something
		a.ReportState(state)
		PrintCommands(state.Cheat)
		// TODO: maybe move the prompt to the read functions?
		if a.inputError != "" {
			fmt.Println("Error:", a.inputError)
			a.inputError = ""
		}
		a.Prompt("take action")
		userInput, arg := a.ReadInput()
		if alias, ok := game.CommandAliases[userInput]; ok {
			userInput = alias
		}
		if userInput == "help" {
			PrintHelp(state.Cheat)
			continue
		}
		command, ok := game.Commands[userInput]
		if !ok {
			a.inputError = "Invalid input. Try again."
			continue
		}
		command.Action.Target = game.ActionTarget{Name: arg}
		return &command.Action
	}
}

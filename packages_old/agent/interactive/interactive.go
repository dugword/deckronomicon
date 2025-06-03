package interactive

// TODO: Rethink and refactor this package and how it handles Player.
// Did a bunch of quick and dirty changes to get it working with the new
// game state and player system, but it could be cleaner and more consistent.

import (
	"bufio"
	"fmt"

	"deckronomicon/packages/engine"
	"deckronomicon/packages/game"
	"deckronomicon/packages/game/player"
	"deckronomicon/packages/ui"
)

// InteractivePlayerAgent implements the PlayerAgent interface for interactive
// play.
type InteractivePlayerAgent struct {
	uiBuffer   *ui.Buffer
	Scanner    *bufio.Scanner
	player     *player.Player
	inputError string
	prompt     string
}

// Player returns the player ID controlled by this agent.
func (a *InteractivePlayerAgent) RegisterPlayer(player *player.Player) {
	a.player = player
}

// NewInteractivePlayerAgent creates a new InteractivePlayerAgent with the
// given scanner.
func NewInteractivePlayerAgent(scanner *bufio.Scanner) *InteractivePlayerAgent {
	return &InteractivePlayerAgent{
		uiBuffer: ui.NewBuffer(),
		Scanner:  scanner,
		prompt:   ">> ",
	}
}

// ReportState displays the current game state in a terminal UI.
// TODO maybe get player from state.CurrentPlayer?
func (a *InteractivePlayerAgent) ReportState(state player.GameState) {
	s, ok := state.(*engine.GameState)
	if !ok {
		panic(fmt.Sprintf("InteractivePlayerAgent.ReportState: state is not of type *engine.GameState, got %T", state))
		return
	}
	a.uiBuffer.UpdateFromState(s, a.player)
	/// TODO: This should return an error instead of panicking
	if err := a.uiBuffer.Render(); err != nil {
		panic(fmt.Sprintf("Error rendering UI buffer: %v", err))
	}
}

func (a *InteractivePlayerAgent) GetNextAction(state player.GameState) (game.Action, error) {
	s, ok := state.(*engine.GameState)
	if !ok {
		return game.Action{}, fmt.Errorf("InteractivePlayerAgent.GetNextAction: state is not of type *engine.GameState")
	}
	for {
		// TODO Don't call this here, run update or something
		a.ReportState(s)
		PrintCommands(s.CheatsEnabled)
		// TODO: maybe move the prompt to the read functions?
		if a.inputError != "" {
			fmt.Println("Error:", a.inputError)
			a.inputError = ""
		}
		a.Prompt("take action")
		userInput, arg := a.ReadInput()
		if alias, ok := engine.CommandAliases[userInput]; ok {
			userInput = alias
		}
		if userInput == "help" {
			PrintHelp(s.CheatsEnabled)
			continue
		}
		command, ok := engine.Commands[userInput]
		if !ok {
			a.inputError = "Invalid input. Try again."
			continue
		}
		command.Action.Target = game.ActionTarget{Name: arg}
		return command.Action, nil
	}
}

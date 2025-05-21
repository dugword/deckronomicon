package interactive

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"text/template"

	game "deckronomicon/game"
)

// InteractivePlayerAgent implements the PlayerAgent interface for interactive
// play.
type InteractivePlayerAgent struct {
	Scanner *bufio.Scanner
	prompt  string
}

// NewInteractivePlayerAgent creates a new InteractivePlayerAgent with the
// given scanner.
func NewInteractivePlayerAgent(scanner *bufio.Scanner) *InteractivePlayerAgent {
	return &InteractivePlayerAgent{
		Scanner: scanner,
		prompt:  ">> ",
	}
}

// ClearScreen clears the terminal screen.
func ClearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// ReportState displays the current game state in a terminal UI.
func (a *InteractivePlayerAgent) ReportState(state *game.GameState) {
	// ClearScreen()
	tmpl := template.New("display")
	tmpl = template.Must(tmpl.ParseFiles("./interactive/display.tmpl"))
	displayData := struct {
		BattlefieldBox Box
		GameStatusBox  Box
		GraveyardBox   Box
		HandBox        Box
		ManaPoolBox    Box
		MessageBox     Box
	}{
		BattlefieldBox: BattlefieldBox(state),
		GameStatusBox:  GameStatusBox(state),
		GraveyardBox:   GraveyardBox(state),
		HandBox:        HandBox(state),
		ManaPoolBox:    ManaPoolBox(state),
		MessageBox:     MessageBox(state),
	}
	if err := tmpl.ExecuteTemplate(
		// TODO: use passed in stdout from Run
		os.Stdout,
		"display.tmpl",
		displayData,
	); err != nil {
		fmt.Println("Error executing template:", err)
		// TODO: handle error return to main
		os.Exit(1)
	}
}

func (a *InteractivePlayerAgent) GetNextAction(state *game.GameState) game.GameAction {
	for {
		PrintCommands(state.Cheat)
		// TODO: maybe move the prompt to the read functions?
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
			fmt.Println("Invalid input. Try again.")
			continue
		}
		command.Action.Target = arg
		return command.Action
	}
}

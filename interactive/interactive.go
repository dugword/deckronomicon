package interactive

import (
	"bufio"
	"fmt"

	game "deckronomicon/game"
)

type InteractivePlayerAgent struct {
	Scanner *bufio.Scanner
	prompt  string
}

func NewInteractivePlayerAgent(scanner *bufio.Scanner) *InteractivePlayerAgent {
	return &InteractivePlayerAgent{
		Scanner: scanner,
		prompt:  ">> ",
	}
}

func (a *InteractivePlayerAgent) ReportState(state *game.GameState) {
	ClearScreen()
	DisplayGameState(state)
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

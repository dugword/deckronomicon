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

type DisplayData struct {
	BattlefieldData BoxData
	ChoiceData      BoxData
	GameStatusData  BoxData
	GraveyardData   BoxData
	HandData        BoxData
	ManaPoolData    BoxData
	MessageData     BoxData
	StackData       BoxData
}

type DisplayBoxes struct {
	GameStatusBox Box
	MessageBox    Box
	PlaySpaceBox  Box
	PlayerBox     Box
}

// InteractivePlayerAgent implements the PlayerAgent interface for interactive
// play.
type InteractivePlayerAgent struct {
	DisplayData     DisplayData
	DisplayTemplate *template.Template
	Scanner         *bufio.Scanner
	inputError      string
	prompt          string
}

func (a *InteractivePlayerAgent) UpdateDisplayData(state *game.GameState) {
	a.DisplayData = DisplayData{
		BattlefieldData: BattlefieldData(state),
		ChoiceData:      BoxData{Title: "Nothing to choose"},
		GameStatusData:  GameStatusData(state),
		GraveyardData:   GraveyardData(state),
		HandData:        HandData(state),
		ManaPoolData:    ManaPoolData(state),
		MessageData:     MessageData(state),
		StackData:       StackData(state),
	}
}

func (a *InteractivePlayerAgent) UpdateChoiceData(choiceData BoxData) {
	a.DisplayData.ChoiceData = choiceData
}

func (a *InteractivePlayerAgent) BuildDisplayBoxes() DisplayBoxes {
	return DisplayBoxes{
		GameStatusBox: CombineBoxesSideBySide(
			CreateBox(a.DisplayData.GameStatusData),
			CreateBox(a.DisplayData.ManaPoolData),
		),
		PlaySpaceBox: CombineBoxesSideBySide(
			CombineBoxesSideBySide(
				CreateBox(a.DisplayData.GraveyardData),
				CreateBox(a.DisplayData.BattlefieldData),
			),
			CreateBox(a.DisplayData.StackData),
		),
		PlayerBox: CombineBoxesSideBySide(
			CreateBox(a.DisplayData.HandData),
			CreateBox(a.DisplayData.ChoiceData),
		),
		MessageBox: CreateBox(a.DisplayData.MessageData),
	}
}

// NewInteractivePlayerAgent creates a new InteractivePlayerAgent with the
// given scanner.
func NewInteractivePlayerAgent(scanner *bufio.Scanner) *InteractivePlayerAgent {
	tmpl := template.New("display")
	tmpl = template.Must(tmpl.ParseFiles(
		"./interactive/display.tmpl",
	))
	return &InteractivePlayerAgent{
		DisplayData:     DisplayData{},
		Scanner:         scanner,
		DisplayTemplate: tmpl,
		prompt:          ">> ",
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
	a.UpdateDisplayData(state)
	displayBoxes := a.BuildDisplayBoxes()
	ClearScreen()
	if err := a.DisplayTemplate.ExecuteTemplate(
		// TODO: use passed in stdout from Run
		os.Stdout,
		"display.tmpl",
		displayBoxes,
	); err != nil {
		fmt.Println("Error executing template:", err)
		// TODO: handle error return to main
		os.Exit(1)
	}
}

func (a *InteractivePlayerAgent) GetNextAction(state *game.GameState) *game.GameAction {
	for {
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
		command.Action.Target = arg
		return &command.Action
	}
}

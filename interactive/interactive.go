package interactive

// TODO: Rethink and refactor this package and how it handles Player.
// Did a bunch of quick and dirty changes to get it working with the new
// game state and player system, but it could be cleaner and more consistent.

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"text/template"

	"deckronomicon/game"
)

type DisplayData struct {
	BattlefieldData BoxData
	ChoiceData      BoxData
	GameStatusData  BoxData
	GraveyardData   BoxData
	HandData        BoxData
	MessageData     BoxData
	OpponentData    BoxData
	RevealedData    BoxData
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
	playerID        string
	inputError      string
	prompt          string
}

// Player returns the player ID controlled by this agent.
func (a *InteractivePlayerAgent) PlayerID() string {
	return a.playerID
}

// TODO maybe get from state.CurrentPlayer?
func (a *InteractivePlayerAgent) UpdateDisplayData(state *game.GameState) {
	// TODO: Handle these errors
	player, err := state.GetPlayer(a.playerID)
	if err != nil {
		panic("InteractivePlayerAgent: Player not found in state: " + a.playerID)
	}
	opponent, err := state.GetOpponent(player.ID)
	if err != nil {
		panic("InteractivePlayerAgent: Opponent not found for player: " + player.ID)
	}
	a.DisplayData = DisplayData{
		BattlefieldData: BattlefieldData(player),
		ChoiceData:      BoxData{Title: "Nothing to choose"},
		GameStatusData:  GameStatusData(state, player),
		GraveyardData:   GraveyardData(player),
		HandData:        HandData(player),
		OpponentData:    OpponentData(state, opponent),
		MessageData:     MessageData(state),
		RevealedData:    RevealedData(player),
		StackData:       StackData(state),
	}
}

func (a *InteractivePlayerAgent) UpdateChoiceData(choiceData BoxData) {
	a.DisplayData.ChoiceData = choiceData
}

func (a *InteractivePlayerAgent) BuildDisplayBoxes() DisplayBoxes {
	playerBox := CombineBoxesSideBySide(
		CreateBox(a.DisplayData.HandData),
		CreateBox(a.DisplayData.ChoiceData),
	)
	if len(a.DisplayData.RevealedData.Content) > 0 {
		playerBox = CombineBoxesSideBySide(
			playerBox,
			CreateBox(a.DisplayData.RevealedData),
		)
	}
	return DisplayBoxes{
		GameStatusBox: CombineBoxesSideBySide(
			CreateBox(a.DisplayData.GameStatusData),
			CreateBox(a.DisplayData.OpponentData),
		),
		PlaySpaceBox: CombineBoxesSideBySide(
			CombineBoxesSideBySide(
				CreateBox(a.DisplayData.GraveyardData),
				CreateBox(a.DisplayData.BattlefieldData),
			),
			CreateBox(a.DisplayData.StackData),
		),
		PlayerBox:  playerBox,
		MessageBox: CreateBox(a.DisplayData.MessageData),
	}
}

// NewInteractivePlayerAgent creates a new InteractivePlayerAgent with the
// given scanner.
func NewInteractivePlayerAgent(scanner *bufio.Scanner, playerID string) *InteractivePlayerAgent {
	tmpl := template.New("display")
	tmpl = template.Must(tmpl.ParseFiles(
		"./interactive/display.tmpl",
	))
	return &InteractivePlayerAgent{
		playerID:        playerID,
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
// TODO maybe get player from state.CurrentPlayer?
func (a *InteractivePlayerAgent) ReportState(state *game.GameState) {
	a.UpdateDisplayData(state)
	displayBoxes := a.BuildDisplayBoxes()
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
		command.Action.Target = arg
		return &command.Action
	}
}

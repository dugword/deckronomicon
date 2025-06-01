package ui

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/game/permanent"
	"deckronomicon/packages/game/player"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"text/template"
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

type BoxData struct {
	Title   string
	Content []string
}

type Buffer struct {
	displayData     *DisplayData
	displayTemplate *template.Template
}

// NewBuffer creates a new Buffer with empty display data.
func NewBuffer() *Buffer {
	tmpl := template.New("display")
	tmpl = template.Must(tmpl.ParseFiles(
		"./ui/term/display.tmpl",
	))
	buffer := Buffer{
		displayData: &DisplayData{
			BattlefieldData: BoxData{Title: "Battlefield", Content: []string{}},
			ChoiceData:      BoxData{Title: "Nothing to choose"},
			GameStatusData:  BoxData{Title: "Game Status", Content: []string{}},
			GraveyardData:   BoxData{Title: "Graveyard", Content: []string{}},
			HandData:        BoxData{Title: "Hand", Content: []string{}},
			MessageData:     BoxData{Title: "Message", Content: []string{}},
			OpponentData:    BoxData{Title: "Opponent Status", Content: []string{}},
			RevealedData:    BoxData{Title: "Revealed Cards", Content: []string{}},
			StackData:       BoxData{Title: "Stack", Content: []string{}},
		},
		displayTemplate: tmpl,
	}
	return &buffer
}

func (b *Buffer) UpdateFromState(state *engine.GameState, player *player.Player) error {
	opponent, err := state.GetOpponent(player.ID())
	if err != nil {
		return fmt.Errorf("opponent not found in state for player: %s", player.ID())
	}
	displayData := DisplayData{
		BattlefieldData: BattlefieldData(state),
		ChoiceData:      BoxData{Title: "Nothing to choose"},
		GameStatusData:  GameStatusData(state, player),
		GraveyardData:   GraveyardData(player),
		HandData:        HandData(player),
		OpponentData:    OpponentData(state, opponent),
		MessageData:     MessageData(state),
		RevealedData:    RevealedData(player),
		StackData:       StackData(state),
	}
	b.displayData = &displayData
	return nil
}

func (b *Buffer) UpdateChoices(title string, choices []choose.Choice) {
	if len(choices) == 0 {
		b.displayData.ChoiceData = BoxData{Title: "Nothing to choose"}
		return
	}
	choiceData := BoxData{
		Title:   title,
		Content: []string{},
	}
	// var subgroupTitle string
	for i, choice := range choices {
		/*
			if choice.Zone != "" && choice.Zone != subgroupTitle {
				choiceData.Content = append(choiceData.Content, fmt.Sprintf("--- %s ---", choice.Zone))
				subgroupTitle = choice.Zone
			}
		*/
		/*
			if choice.Source != "" {
				line = fmt.Sprintf("(%s) - %s", choice.Source, line)
			}
		*/
		var line = choice.Name
		line = fmt.Sprintf("%d: %s", i, line)
		choiceData.Content = append(choiceData.Content, line)
	}
	b.displayData.ChoiceData = choiceData
	//b.displayData.ChoiceData.Content = append(b.displayData.ChoiceData.Content, fmt.Sprintf("Total Choices: %d", len(orderedChoices)))
}

// GameStatusData creates the box data for game status information.
// TODO: Could probably get from state.CurrentPlayer
func GameStatusData(state *engine.GameState, player *player.Player) BoxData {
	// potentialMana := "(empty)"
	manaPool := "(empty)"
	/*
		// TODO: Potential mana is broken
		if player.PotentialMana.AvailableGeneric() > 0 {
			potentialMana = player.PotentialMana.Describe()
		}
	*/
	if player.ManaPool().AvailableGeneric() > 0 {
		manaPool = player.ManaPool().Describe()
	}
	return BoxData{
		Title: "Game Status",
		Content: []string{
			fmt.Sprintf("Current Player: %s", state.GetActivePlayer().ID()),
			fmt.Sprintf("Player Mode: %s", player.Mode),
			fmt.Sprintf("Player 1 Life: %d", player.Life),
			fmt.Sprintf("Turn: %d", player.Turn),
			fmt.Sprintf("Phase: %s", state.CurrentPhase),
			fmt.Sprintf("Step: %s", state.CurrentStep),
			fmt.Sprintf("Library: %d cards", player.Library().Size()),
			fmt.Sprintf("Graveyard: %d cards", player.Graveyard().Size()),
			fmt.Sprintf("Exile: %d cards", player.Exile().Size()),
			fmt.Sprintf("Hand: %d cards", player.Hand().Size()),
			// fmt.Sprintf("Potential Mana: %s", potentialMana),
			fmt.Sprintf("Mana Pool: %s", manaPool),
		},
	}
}

func OpponentData(state *engine.GameState, player *player.Player) BoxData {
	return BoxData{
		Title: "Opponent Status",
		Content: []string{
			fmt.Sprintf("Opponent Life: %d", player.Life()),
			fmt.Sprintf("Library: %d cards", player.Library().Size()),
			fmt.Sprintf("Graveyard: %d cards", player.Graveyard().Size()),
			fmt.Sprintf("Battlefield: %d permanents", state.Battlefield().Size()),
			fmt.Sprintf("Hand: %d cards", player.Hand().Size()),
		},
	}
}

// MessageData creates the box data for displaying messages.
func MessageData(state *engine.GameState) BoxData {
	lines := splitLines(state.Message, 70)
	return BoxData{
		Title:   "Message",
		Content: lines,
	}
}

// BattlefieldData creates the box data for displaying permanents on the
// battlefield.
func BattlefieldData(state *engine.GameState) BoxData {
	var lines []string
	for _, obj := range state.Battlefield().GetAll() {
		perm, ok := obj.(*permanent.Permanent)
		if !ok {
			continue // Skip if not a permanent
		}
		line := perm.Name()
		if perm.IsTapped() {
			line += " (tapped)"
		}
		if perm.HasSummoningSickness() {
			line += " (summoning sick)"
		}
		lines = append(lines, line)
	}
	return BoxData{
		Title:   "Battlefield",
		Content: lines,
	}
}

// GraveyardData creates the box data for displaying cards in the graveyard.
func GraveyardData(player *player.Player) BoxData {
	var lines []string
	for _, card := range player.Graveyard().GetAll() {
		line := card.Name()
		lines = append(lines, line)
	}
	return BoxData{
		Title:   "Graveyard",
		Content: lines,
	}
}

// HandData creates the box data for displaying cards in the player's hand.
func HandData(player *player.Player) BoxData {
	var lines []string
	for _, card := range player.Hand().GetAll() {
		line := card.Name()
		lines = append(lines, line)
	}
	return BoxData{
		Title:   "Hand",
		Content: lines,
	}
}

// StackData creates the box data for displaying cards in the player's hand.
func StackData(state *engine.GameState) BoxData {
	var lines []string
	for _, spell := range state.Stack.GetAll() {
		lines = append(lines, spell.Name())
	}
	return BoxData{
		Title:   "Stack",
		Content: lines,
	}
}

// RevealedData creates the box data for displaying cards in the player's hand.
func RevealedData(player *player.Player) BoxData {
	var lines []string
	for _, spell := range player.Revealed().GetAll() {
		lines = append(lines, spell.Name())
	}
	return BoxData{
		Title:   "Revealed Cards",
		Content: lines,
	}
}

func (b *Buffer) BuildDisplayBoxes() DisplayBoxes {
	playerBox := CombineBoxesSideBySide(
		CreateBox(b.displayData.HandData),
		CreateBox(b.displayData.ChoiceData),
	)
	if len(b.displayData.RevealedData.Content) > 0 {
		playerBox = CombineBoxesSideBySide(
			playerBox,
			CreateBox(b.displayData.RevealedData),
		)
	}
	return DisplayBoxes{
		GameStatusBox: CombineBoxesSideBySide(
			CreateBox(b.displayData.GameStatusData),
			CreateBox(b.displayData.OpponentData),
		),
		PlaySpaceBox: CombineBoxesSideBySide(
			CombineBoxesSideBySide(
				CreateBox(b.displayData.GraveyardData),
				CreateBox(b.displayData.BattlefieldData),
			),
			CreateBox(b.displayData.StackData),
		),
		PlayerBox:  playerBox,
		MessageBox: CreateBox(b.displayData.MessageData),
	}
}

func (b *Buffer) Render() error {
	ClearScreen()
	displayBoxes := b.BuildDisplayBoxes()
	if err := b.displayTemplate.ExecuteTemplate(
		// TODO: use passed in stdout from Run
		os.Stdout,
		"display.tmpl",
		displayBoxes,
	); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}
	return nil
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

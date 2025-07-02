package ui

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/view"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"text/template"
)

const templateName = "display.tmpl"

type DisplayData struct {
	BattlefieldData BoxData
	ChoiceData      BoxData
	GameStatusData  BoxData
	GraveyardData   BoxData
	HandData        BoxData
	MessageData     BoxData
	OpponentData    BoxData
	PlayerData      BoxData
	RevealedData    BoxData
	StackData       BoxData
}

type DisplayBoxes struct {
	PlayersStatusBox Box
	GameStatusBox    Box
	MessageBox       Box
	PlaySpaceBox     Box
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
func NewBuffer(displayFile string) *Buffer {
	tmpl := template.New("display")
	tmpl = template.Must(tmpl.ParseFiles(
		displayFile,
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
			PlayerData:      BoxData{Title: "Player Status", Content: []string{}},
			RevealedData:    BoxData{Title: "Revealed Cards", Content: []string{}},
			StackData:       BoxData{Title: "Stack", Content: []string{}},
		},
		displayTemplate: tmpl,
	}
	return &buffer
}

func (b *Buffer) Update(
	game *view.Game, player *view.Player, opponent *view.Player,
) {
	displayData := DisplayData{
		BattlefieldData: BattlefieldData(game.Battlefield),
		ChoiceData:      BoxData{Title: "Nothing to choose"},
		GameStatusData:  GameStatusData(game),
		GraveyardData:   CardListData("Graveyard", player.Graveyard),
		HandData:        CardListData("Hand", player.Hand),
		OpponentData:    PlayerData(opponent),
		PlayerData:      PlayerData(player),
		MessageData:     MessageData([]string{}),
		RevealedData:    CardListData("Revealed Cards", player.Revealed),
		StackData:       StackData(game.Stack),
	}
	b.displayData = &displayData
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
	for i, choice := range choices {
		var line = fmt.Sprintf("%s <%s>", choice.Name(), choice.ID())
		line = fmt.Sprintf("%d: %s", i+1, line)
		choiceData.Content = append(choiceData.Content, line)
	}
	b.displayData.ChoiceData = choiceData
}

func (b *Buffer) UpdateMessage(lines []string) {
	b.displayData.MessageData = MessageData(lines)
}

func GameStatusData(game *view.Game) BoxData {
	return BoxData{
		Title: "Game Status",
		Content: []string{
			fmt.Sprintf("Active Player: %s", game.ActivePlayerID),
			fmt.Sprintf("Phase: %s", game.Phase),
			fmt.Sprintf("Step: %s", game.Step),
			fmt.Sprintf("Stack: %d", len(game.Stack)),
			fmt.Sprintf("Battlefield: %d", len(game.Battlefield)),
		},
	}
}

// BattlefieldData creates the box data for displaying permanents on the
// battlefield.
func BattlefieldData(permanents []*view.Permanent) BoxData {
	var lines []string
	for _, permanent := range permanents {
		line := fmt.Sprintf(
			"[%s] %s <%s>",
			permanent.Controller,
			permanent.Name,
			permanent.ID,
		)
		if permanent.Tapped {
			line += " (tapped)"
		}
		if permanent.SummoningSick {
			line += " (summoning sick)"
		}
		lines = append(lines, line)
	}
	return BoxData{
		Title:   "Battlefield",
		Content: lines,
	}
}

// BattlefieldData creates the box data for displaying permanents on the
// battlefield.
func StackData(Resolvable []*view.Resolvable) BoxData {
	var lines []string
	for _, resolvable := range Resolvable {
		line := fmt.Sprintf(
			"[%s] %s <%s>",
			resolvable.Controller,
			resolvable.Name,
			resolvable.ID,
		)
		lines = append(lines, line)
	}
	return BoxData{
		Title:   "Stack",
		Content: lines,
	}
}

func PlayerData(player *view.Player) BoxData {
	content := []string{}
	if player.Mode != "" {
		content = append(content, fmt.Sprintf("Mode: %s", player.Mode))
	}
	content = append(content,
		fmt.Sprintf("Turn: %d", player.Turn),
		fmt.Sprintf("Life: %d", player.Life),
		fmt.Sprintf("Hand: %d cards", len(player.Hand)),
		fmt.Sprintf("Library: %d cards", player.LibrarySize),
		fmt.Sprintf("Graveyard: %d cards", len(player.Graveyard)),
		fmt.Sprintf("Exile: %d cards", len(player.Exile)),
		fmt.Sprintf("Mana Pool: %s", player.ManaPool),
		fmt.Sprintf("Potential Mana Pool: %s", player.PotentialManaPool),
	)
	return BoxData{
		Title:   fmt.Sprintf("Status for %s", player.ID),
		Content: content,
	}
}

// GraveyardData creates the box data for displaying cards in the graveyard.
func CardListData(title string, cards []*view.Card) BoxData {
	var lines []string
	for _, card := range cards {
		line := fmt.Sprintf("%s <%s>", card.Name, card.ID)
		lines = append(lines, line)
	}
	return BoxData{
		Title:   title,
		Content: lines,
	}
}

func MessageData(lines []string) BoxData {
	var lns []string
	for _, line := range lines {
		// TODO: These var names suck
		x := splitLines(line, 22)
		lns = append(lns, x...)
	}
	return BoxData{
		Title:   "Message",
		Content: lns,
	}
}

func (b *Buffer) BuildDisplayBoxes() DisplayBoxes {
	return DisplayBoxes{
		PlayersStatusBox: CombineBoxesSideBySide(
			CombineBoxesSideBySide(
				CreateBox(b.displayData.PlayerData),
				CreateBox(b.displayData.OpponentData),
			),
			CreateBox(b.displayData.RevealedData),
		),
		GameStatusBox: CombineBoxesSideBySide(
			CreateBox(b.displayData.GameStatusData),
			CreateBox(b.displayData.BattlefieldData),
		),
		PlaySpaceBox: CombineBoxesSideBySide(
			CombineBoxesSideBySide(
				CreateBox(b.displayData.GraveyardData),
				CreateBox(b.displayData.HandData),
			),
			CreateBox(b.displayData.StackData),
		),
		MessageBox: CombineBoxesSideBySide(
			CreateBox(b.displayData.MessageData),
			CreateBox(b.displayData.ChoiceData),
		),
	}
}

// ClearScreen clears the terminal screen.
func clearScreen() {
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

func (b *Buffer) Render() error {
	// clearScreen()
	displayBoxes := b.BuildDisplayBoxes()
	if err := b.displayTemplate.ExecuteTemplate(
		// TODO: use passed in stdout from Run
		os.Stdout,
		templateName,
		displayBoxes,
	); err != nil {
		return fmt.Errorf("failed to execute template %q: %w", templateName, err)
	}
	return nil
}

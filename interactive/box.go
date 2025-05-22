package interactive

import (
	"deckronomicon/game"
	"fmt"
	"strings"
)

// Box represents a terminal UI box with a title and content.
type Box struct {
	Title      string
	TopLine    string
	MiddleLine string
	BottomLine string
	Lines      []string
}

// splitLines splits a string into lines of a specified maximum length.
func splitLines(text string, maxLength int) []string {
	var lines []string
	for len(text) > maxLength {
		line := text[:maxLength]
		lines = append(lines, line)
		text = text[maxLength:]
	}
	if len(text) > 0 {
		lines = append(lines, text)
	}
	return lines
}

// ChoicesBox creates a Box for displaying choices.
func ChoicesBox(title string, choices []game.Choice) Box {
	var lines []string
	for i, choice := range choices {
		lines = append(lines, fmt.Sprintf("%d: %s", i, choice.Name))
	}
	return CreateBox(title, lines)
}

// CreateMessageBox creates a Box for displaying messages.
func MessageBox(state *game.GameState) Box {
	lines := splitLines(state.Message, 70)
	return CreateBox("MESSAGE", lines)
}

// GameStatusBox creates a Box for displaying game status.
func GameStatusBox(state *game.GameState) Box {
	return CreateBox("GAME STATUS", []string{
		fmt.Sprintf("Current Player: Player %d", state.CurrentPlayer),
		fmt.Sprintf("Player 1 Life: %d", state.Life),
		fmt.Sprintf("Turn: %d", state.Turn),
		fmt.Sprintf("Phase: %s", state.CurrentPhase),
		fmt.Sprintf("Library: %d cards", state.Library.Size()),
		fmt.Sprintf("Graveyard: %d cards", len(state.Graveyard)),
		fmt.Sprintf("Hand: %d cards", state.Hand.Size()),
	})
}

// BattlefieldBox creates a Box for displaying permanents on the battlefield.
func BattlefieldBox(state *game.GameState) Box {
	var lines []string
	for _, permanent := range state.Battlefield.Permanents() {
		line := permanent.Name()
		if permanent.IsTapped() {
			line += " (tapped)"
		}
		if permanent.HasSummoningSickness() {
			line += " (summoning sick)"
		}
		lines = append(lines, line)
	}
	return CreateBox("BATTLEFIELD", lines)
}

// GraveyardBox creates a Box for displaying cards in the graveyard.
func GraveyardBox(state *game.GameState) Box {
	var lines []string
	for _, card := range state.Graveyard {
		line := card.Name()
		lines = append(lines, line)
	}
	return CreateBox("GRAVEYARD", lines)
}

// HandBox creates a Box for displaying cards in the player's hand.
func HandBox(state *game.GameState) Box {
	var lines []string
	for _, card := range state.Hand.Cards() {
		line := card.Name()
		lines = append(lines, line)
	}
	return CreateBox("Hand", lines)
}

// ManaPoolBox creates a Box for displaying the player's mana pool.
func ManaPoolBox(state *game.GameState) Box {
	potentialMana := "(empty)"
	manaPool := "(empty)"
	if state.PotentialMana.AvailableGeneric() > 0 {
		potentialMana = state.PotentialMana.Describe()
	}
	if state.ManaPool.AvailableGeneric() > 0 {
		manaPool = state.ManaPool.Describe()
	}
	lines := []string{
		"Potential Mana: " + potentialMana,
		"Mana Pool: " + manaPool,
	}
	return CreateBox("Mana Pool", lines)
}

// CreateBox creates a Box with a title and content, ensuring the title and
// content are padded to the same width.
func CreateBox(title string, content []string) Box {
	maxLength := len(title)
	for _, line := range content {
		if len(line) > maxLength {
			maxLength = len(line)
		}
	}
	var padded []string
	for _, line := range content {
		paddedLine := fmt.Sprintf("║ %-*s ║", maxLength, line)
		padded = append(padded, paddedLine)
	}
	width := maxLength + 2 // 2 for 1 space on each side
	border := strings.Repeat("═", width)
	return Box{
		Title:      fmt.Sprintf("║ %-*s ║", maxLength, title),
		TopLine:    fmt.Sprintf("╔%s╗", border),
		MiddleLine: fmt.Sprintf("╠%s╣", border),
		BottomLine: fmt.Sprintf("╚%s╝", border),
		Lines:      padded,
	}
}

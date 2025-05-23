package interactive

import (
	"deckronomicon/game"
	"fmt"
	"sort"
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

func GroupedChoicesBox(title string, choices []game.Choice) (Box, []game.Choice) {
	grouped := make(map[string][]game.Choice)
	for _, choice := range choices {
		if choice.ID == game.ChoiceNone {
			grouped[game.ChoiceNone] = append(grouped[game.ChoiceNone], choice)
		} else if choice.Zone == "" {
			// TODO: Handle this different
			grouped[""] = append(grouped[""], choice)
		} else {
			grouped[choice.Zone] = append(grouped[choice.Zone], choice)
		}
	}
	var groupNames []string
	for groupName := range grouped {
		if groupName == game.ChoiceNone {
			continue
		}
		groupNames = append(groupNames, groupName)
	}
	sort.Strings(groupNames)
	// TODO: Make 0 always none, and have it be not an option if the choice is
	// not optional
	var orderedChoices []game.Choice
	var i = 0
	var lines []string
	for _, groupName := range groupNames {
		lines = append(lines, fmt.Sprintf("--- %s ---", groupName))
		for _, choice := range grouped[groupName] {
			var line string
			if choice.Source == "" {
				line = fmt.Sprintf("%d: %s", i, choice.Name)
			} else {
				line = fmt.Sprintf("%d: %s - %s", i, choice.Source, choice.Name)
			}
			lines = append(lines, line)
			orderedChoices = append(orderedChoices, choice)
			i++
		}
	}
	if grouped[game.ChoiceNone] != nil {
		lines = append(lines, "---------")
		lines = append(lines, fmt.Sprintf("%d: %s", i, game.ChoiceNone))
		orderedChoices = append(orderedChoices, game.Choice{ID: game.ChoiceNone})
	}
	return CreateBox(title, lines), orderedChoices
}

// ChoicesBox creates a Box for displaying choices.
func ChoicesBox(title string, choices []game.Choice) Box {
	var lines []string
	for i, choice := range choices {
		var line string
		if choice.ID == game.ChoiceNone {
			line = fmt.Sprintf("%d: %s", i, choice.Name)
		} else if choice.Source == "" {
			line = fmt.Sprintf("%d: %s - %s", i, choice.Zone, choice.Name)
		} else {
			line = fmt.Sprintf("%d: %s - %s - %s", i, choice.Zone, choice.Source, choice.Name)
		}
		lines = append(lines, line)
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
		fmt.Sprintf("Graveyard: %d cards", len(state.Graveyard.Cards())),
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
	for _, card := range state.Graveyard.Cards() {
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

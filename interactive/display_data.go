package interactive

import (
	"deckronomicon/game"
	"fmt"
	"sort"
)

type BoxData struct {
	Title   string
	Content []string
}

// GameStatusData creates the box data for game status information.
func GameStatusData(state *game.GameState) BoxData {
	return BoxData{
		Title: "Game Status",
		Content: []string{
			fmt.Sprintf("Current Player: Player %d", state.CurrentPlayer),
			fmt.Sprintf("Player 1 Life: %d", state.Life),
			fmt.Sprintf("Turn: %d", state.Turn),
			fmt.Sprintf("Phase: %s", state.CurrentPhase),
			fmt.Sprintf("Step: %s", state.CurrentStep),
			fmt.Sprintf("Library: %d cards", state.Library.Size()),
			fmt.Sprintf("Graveyard: %d cards", len(state.Graveyard.Cards())),
			fmt.Sprintf("Hand: %d cards", state.Hand.Size()),
		},
	}
}

// ManaPoolData creates the box data for displaying the player's mana pool.
func ManaPoolData(state *game.GameState) BoxData {
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
	return BoxData{
		Title:   "Mana Pool",
		Content: lines,
	}
}

// MessageData creates the box data for displaying messages.
func MessageData(state *game.GameState) BoxData {
	lines := splitLines(state.Message, 70)
	return BoxData{
		Title:   "Message",
		Content: lines,
	}
}

// BattlefieldData creates the box data for displaying permanents on the
// battlefield.
func BattlefieldData(state *game.GameState) BoxData {
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
	return BoxData{
		Title:   "Battlefield",
		Content: lines,
	}
}

// GraveyardData creates the box data for displaying cards in the graveyard.
func GraveyardData(state *game.GameState) BoxData {
	var lines []string
	for _, card := range state.Graveyard.Cards() {
		line := card.Name()
		lines = append(lines, line)
	}
	return BoxData{
		Title:   "Graveyard",
		Content: lines,
	}
}

// HandData creates the box data for displaying cards in the player's hand.
func HandData(state *game.GameState) BoxData {
	var lines []string
	for _, card := range state.Hand.Cards() {
		line := card.Name()
		lines = append(lines, line)
	}
	return BoxData{
		Title:   "Hand",
		Content: lines,
	}
}

// StackData creates the box data for displaying cards in the player's hand.
func StackData(state *game.GameState) BoxData {
	var lines []string
	for _, spell := range state.Stack.GetAll() {
		lines = append(lines, spell.Name())
	}
	return BoxData{
		Title:   "Stack",
		Content: lines,
	}
}

func GroupedChoicesData(title string, choices []game.Choice) (BoxData, []game.Choice) {
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
	return BoxData{
		Title:   title,
		Content: lines,
	}, orderedChoices
}

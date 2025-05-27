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
// TODO: Could probably get from state.CurrentPlayer
func GameStatusData(state *game.GameState, player *game.Player) BoxData {
	potentialMana := "(empty)"
	manaPool := "(empty)"
	// TODO: Potential mana is broken
	if player.PotentialMana.AvailableGeneric() > 0 {
		potentialMana = player.PotentialMana.Describe()
	}
	if player.ManaPool.AvailableGeneric() > 0 {
		manaPool = player.ManaPool.Describe()
	}
	return BoxData{
		Title: "Game Status",
		Content: []string{
			fmt.Sprintf("Current Player: Player %d", state.CurrentPlayer),
			fmt.Sprintf("Player 1 Life: %d", player.Life),
			fmt.Sprintf("Turn: %d", player.Turn),
			fmt.Sprintf("Phase: %s", state.CurrentPhase),
			fmt.Sprintf("Step: %s", state.CurrentStep),
			fmt.Sprintf("Library: %d cards", player.Library.Size()),
			fmt.Sprintf("Graveyard: %d cards", player.Graveyard.Size()),
			fmt.Sprintf("Exile: %d cards", player.Exile.Size()),
			fmt.Sprintf("Hand: %d cards", player.Hand.Size()),
			fmt.Sprintf("Potential Mana: %s", potentialMana),
			fmt.Sprintf("Mana Pool: %s", manaPool),
		},
	}
}

func OpponentData(state *game.GameState, player *game.Player) BoxData {
	return BoxData{
		Title: "Opponent Status",
		Content: []string{
			fmt.Sprintf("Opponent Life: %d", player.Life),
			fmt.Sprintf("Library: %d cards", player.Library.Size()),
			fmt.Sprintf("Graveyard: %d cards", player.Graveyard.Size()),
			fmt.Sprintf("Battlefield: %d permanents", player.Battlefield.Size()),
			fmt.Sprintf("Hand: %d cards", player.Hand.Size()),
		},
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
func BattlefieldData(player *game.Player) BoxData {
	var lines []string
	for _, permanent := range player.Battlefield.Permanents() {
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
func GraveyardData(player *game.Player) BoxData {
	var lines []string
	for _, card := range player.Graveyard.GetAll() {
		line := card.Name()
		lines = append(lines, line)
	}
	return BoxData{
		Title:   "Graveyard",
		Content: lines,
	}
}

// HandData creates the box data for displaying cards in the player's hand.
func HandData(player *game.Player) BoxData {
	var lines []string
	for _, card := range player.Hand.GetAll() {
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

// RevealedData creates the box data for displaying cards in the player's hand.
func RevealedData(player *game.Player) BoxData {
	var lines []string
	for _, spell := range player.Revealed.GetAll() {
		lines = append(lines, spell.Name())
	}
	return BoxData{
		Title:   "Revealed Cards",
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

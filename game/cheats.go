package game

// GameCheatType represents the type of player cheat
type GameCheatType string

var ChoiceCheat = "Cheat"
var ChoiceSourceCheat = NewChoiceSource(ChoiceCheat, ChoiceCheat)

const (
	CheatAddMana   GameCheatType = "CheatAddMana"
	CheatConjure   GameCheatType = "CheatConjure"
	CheatDiscard   GameCheatType = "CheatDiscard"
	CheatDraw      GameCheatType = "CheatDraw"
	CheatFind      GameCheatType = "CheatFind"
	CheatLandDrop  GameCheatType = "CheatLandDrop"
	CheatPeek      GameCheatType = "CheatPeek"
	CheatPrintDeck GameCheatType = "CheatPrintDeck"
	CheatShuffle   GameCheatType = "CheatShuffle"
)

// CheatResult represents the result of a cheat action.
type CheatResult struct {
	Message string
}

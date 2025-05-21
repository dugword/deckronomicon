package game

// GameCheatType represents the type of player cheat
type GameCheatType string

const (
	CheatDiscard   GameCheatType = "CheatDiscard"
	CheatDraw      GameCheatType = "CheatDraw"
	CheatFind      GameCheatType = "CheatFind"
	CheatPeek      GameCheatType = "CheatPeek"
	CheatPrintDeck GameCheatType = "CheatPrintDeck"
	CheatShuffle   GameCheatType = "CheatShuffle"
)

// CheatResult represents the result of a cheat action.
type CheatResult struct {
	Message string
}

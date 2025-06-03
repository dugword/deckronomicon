package game

// Not sure I love this here, maybe it should be in engine?

type ActionTarget struct {
	Name string
	ID   string
}

// Action represents an action a player can take.
type Action struct {
	Target ActionTarget
	Type   ActionType
}

// ActionType represents the type of player action.
type ActionType string

// Constants for game action types.
const (
	ActionActivate ActionType = "Activate"
	ActionDiscard  ActionType = "Discard"
	ActionDraw     ActionType = "Draw"
	ActionCheat    ActionType = "Cheat"
	ActionConcede  ActionType = "Concede"
	ActionPass     ActionType = "Pass"
	ActionPlay     ActionType = "Play"
	ActionUntap    ActionType = "Untap"
	ActionView     ActionType = "View"

	// Cheat actions
	CheatAddMana   ActionType = "CheatAddMana"
	CheatConjure   ActionType = "CheatConjure"
	CheatDiscard   ActionType = "CheatDiscard"
	CheatDraw      ActionType = "CheatDraw"
	CheatFind      ActionType = "CheatFind"
	CheatLandDrop  ActionType = "CheatLandDrop"
	CheatPeek      ActionType = "CheatPeek"
	CheatPrintDeck ActionType = "CheatPrintDeck"
	CheatShuffle   ActionType = "CheatShuffle"
)

// TODO: Find some way to enforce this, maybe an action map or lookup function
// TODO: Function signature should also be func(state, player.Player, string) (*ActionResult, error)
// type ActionFunc func(state, player.Player, Target) (*ActionResult, error)

// TODO: This is a bit of a hack. We should probably have a better way to
// to create the object for ChoiceSource.
// This is for the ChoiceSorce interface
func (a ActionType) Name() string {
	return string(a)
}

// This is for the ChoiceSorce interface
func (a ActionType) ID() string {
	return string(a)
}

// PlayerActions is a map of game action types to booleans indicating if the
// action is a player action.
// TODO make this a function that takes an ActionType and returns a bool
var PlayerActions = map[ActionType]bool{
	ActionActivate: true,
	ActionCheat:    true,
	ActionConcede:  true,
	ActionPass:     true,
	ActionPlay:     true,
	ActionView:     true,
}

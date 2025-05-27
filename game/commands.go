package game

// Command represents a command that can be executed in the game.
type Command struct {
	Action      GameAction
	Cheat       bool
	Description string
}

// Commands is a map of command names to their corresponding Command structs.
var Commands = map[string]Command{
	"addmana": {
		Action:      GameAction{Type: CheatAddMana},
		Cheat:       true,
		Description: "Add mana to your mana pool",
	},
	"activate": {
		Action:      GameAction{Type: ActionActivate},
		Description: "Activate an ability",
	},
	"cheat": {
		Action:      GameAction{Type: ActionCheat},
		Cheat:       true,
		Description: "Enable cheat actions",
	},
	"concede": {
		Action:      GameAction{Type: ActionConcede},
		Description: "Concede the game",
	},
	"conjure": {
		Action:      GameAction{Type: CheatConjure},
		Cheat:       true,
		Description: "Conjure a card",
	},
	"draw": {
		Action:      GameAction{Type: CheatDraw},
		Cheat:       true,
		Description: "Draw a card",
	},
	"discard": {
		Action:      GameAction{Type: CheatDiscard},
		Cheat:       true,
		Description: "Discard a card",
	},
	"find": {
		Action:      GameAction{Type: CheatFind},
		Cheat:       true,
		Description: "Find a card in the library",
	},
	"help": {
		// TODO: maybe have something here incase it some how gets passed to the game loop....
		Action:      GameAction{},
		Description: "Print more information about a command",
	},
	"landdrop": {
		Action:      GameAction{Type: CheatLandDrop},
		Cheat:       true,
		Description: "Reset land drop check",
	},
	"pass": {
		Action:      GameAction{Type: ActionPass},
		Description: "Pass the turn",
	},
	"peek": {
		Action:      GameAction{Type: CheatPeek},
		Cheat:       true,
		Description: "Peek at the top of the library",
	},
	"play": {
		Action:      GameAction{Type: ActionPlay},
		Description: "Play a land or cast a spell",
	},
	"shuffle": {
		Action:      GameAction{Type: CheatShuffle},
		Cheat:       true,
		Description: "Shuffle the library",
	},
	"untap": {
		Action:      GameAction{Type: ActionUntap},
		Description: "Untap a card",
	},
	"view": {
		Action:      GameAction{Type: ActionView},
		Description: "View an object's description",
	},
}

// CommandAliases is a map of command aliases to their corresponding command
// names.
var CommandAliases = map[string]string{
	"cast": "play",
	"exit": "concede",
	"quit": "concede",
	"tap":  "activate",
}

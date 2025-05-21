package game

// Command represents a command that can be executed in the game.
type Command struct {
	Action      GameAction
	Description string
}

// Commands is a map of command names to their corresponding Command structs.
var Commands = map[string]Command{
	"activate": {
		Description: "Activate an ability",
		Action:      GameAction{Type: ActionActivate},
	},
	"cheat": {
		Description: "Enable cheat actions",
		Action:      GameAction{Type: ActionCheat},
	},
	"concede": {
		Description: "Concede the game",
		Action:      GameAction{Type: ActionConcede},
	},
	"draw": {
		Description: "Draw a card",
		Action:      GameAction{Cheat: CheatDraw},
	},
	"discard": {
		Description: "Discard a card",
		Action:      GameAction{Cheat: CheatDiscard},
	},
	"help": {
		Description: "Print more information about a command",
		// TODO: maybe have something here incase it some how gets passed to the game loop....
		Action: GameAction{},
	},
	"pass": {
		Description: "Pass the turn",
		Action:      GameAction{Type: ActionPass},
	},
	"peek": {
		Description: "Peek at the top of the library",
		Action:      GameAction{Cheat: CheatPeek},
	},
	"play": {
		Description: "Play a land or cast a spell",
		Action:      GameAction{Type: ActionPlay},
	},
	"shuffle": {
		Description: "Shuffle the library",
		Action:      GameAction{Cheat: CheatShuffle},
	},
	"view": {
		Description: "View an object's description",
		Action:      GameAction{Type: ActionView},
	},
}

// CommandAliases is a map of command aliases to their corresponding command
// names.
var CommandAliases = map[string]string{
	"exit": "concede",
	"quit": "concede",
	"tap":  "activate",
}

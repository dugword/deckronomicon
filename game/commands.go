package game

type Command struct {
	Action      GameAction
	Description string
}

var Commands = map[string]Command{
	"activate": {
		Description: "Activate an ability",
		Action:      GameAction{Type: ActionActivate},
	},
	"draw": {
		Description: "Draw a card",
		Action:      GameAction{Cheat: CheatDraw},
	},
	"play": {
		Description: "Play a land or cast a spell",
		Action:      GameAction{Type: ActionPlay},
	},
	"discard": {
		Description: "Discard a card",
		Action:      GameAction{Cheat: CheatDiscard},
	},
	"pass": {
		Description: "Pass the turn",
		Action:      GameAction{Type: ActionPass},
	},
	"concede": {
		Description: "Concede the game",
		Action:      GameAction{Type: ActionConcede},
	},
	"cheat": {
		Description: "Enable cheat actions",
		Action:      GameAction{Type: ActionCheat},
	},
	"shuffle": {
		Description: "Shuffle the library",
		Action:      GameAction{Cheat: CheatShuffle},
	},
	"peek": {
		Description: "Peek at the top of the library",
		Action:      GameAction{Cheat: CheatPeek},
	},
	"help": {
		Description: "Print more information about a command",
		// TODO: maybe have something here incase it some how gets passed to the game loop....
		Action: GameAction{},
	},
	"view": {
		Description: "View an object's description",
		Action:      GameAction{Type: ActionView},
	},
}

var CommandAliases = map[string]string{
	"exit": "concede",
	"quit": "concede",
	"tap":  "activate",
}

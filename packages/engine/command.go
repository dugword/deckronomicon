package engine

// Command represents a command that can be executed in the game.
type Command struct {
	NewAction   Action
	Cheat       bool
	Description string
}

// TODO: Add more commands, add stop, remove stop, set mode,

// Commands is a map of command names to their corresponding Command structs.
var Commands = map[string]Command{
	/*
		"addmana": {
			Action:      game.Action{Type: game.CheatAddMana},
			Cheat:       true,
			Description: "Add mana to your mana pool",
		},
		"activate": {
			Action:      game.Action{Type: game.ActionActivate},
			Description: "Activate an ability",
		},
		"cheat": {
			Action:      game.Action{Type: game.ActionCheat},
			Cheat:       true,
			Description: "Enable cheat actions",
		},
		"concede": {
			Action:      game.Action{Type: game.ActionConcede},
			Description: "Concede the game",
		},
		"conjure": {
			Action:      game.Action{Type: game.CheatConjure},
			Cheat:       true,
			Description: "Conjure a card",
		},
		"draw": {
			Action:      game.Action{Type: game.CheatDraw},
			Cheat:       true,
			Description: "Draw a card",
		},
		"discard": {
			Action:      game.Action{Type: game.CheatDiscard},
			Cheat:       true,
			Description: "Discard a card",
		},
		"find": {
			Action:      game.Action{Type: game.CheatFind},
			Cheat:       true,
			Description: "Find a card in the library",
		},
		"help": {
			// TODO: maybe have something here incase it some how gets passed to the game loop....
			Action:      game.Action{},
			Description: "Print more information about a command",
		},
		"landdrop": {
			Action:      game.Action{Type: game.CheatLandDrop},
			Cheat:       true,
			Description: "Reset land drop check",
		},
		"pass": {
			Action:      game.Action{Type: game.ActionPass},
			Description: "Pass the turn",
		},
		"peek": {
			Action:      game.Action{Type: game.CheatPeek},
			Cheat:       true,
			Description: "Peek at the top of the library",
		},
	*/
	"play": {
		// Action: PlayAction{},
	},
	/*
		"shuffle": {
			Action:      game.Action{Type: game.CheatShuffle},
			Cheat:       true,
			Description: "Shuffle the library",
		},
		"untap": {
			Action:      game.Action{Type: game.ActionUntap},
			Description: "Untap a card",
		},
		"view": {
			Action:      game.Action{Type: game.ActionView},
			Description: "View an object's description",
		},
	*/
}

// CommandAliases is a map of command aliases to their corresponding command
// names.
var CommandAliases = map[string]string{
	"cast": "play",
	"exit": "concede",
	"quit": "concede",
	"tap":  "activate",
}

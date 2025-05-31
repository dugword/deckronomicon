package engine

import "deckronomicon/packages/game/action"

// Command represents a command that can be executed in the game.
type Command struct {
	Action      action.Action
	Cheat       bool
	Description string
}

// TODO: Add more commands, add stop, remove stop, set mode,

// Commands is a map of command names to their corresponding Command structs.
var Commands = map[string]Command{
	"addmana": {
		Action:      action.Action{Type: action.CheatAddMana},
		Cheat:       true,
		Description: "Add mana to your mana pool",
	},
	"activate": {
		Action:      action.Action{Type: action.ActionActivate},
		Description: "Activate an ability",
	},
	"cheat": {
		Action:      action.Action{Type: action.ActionCheat},
		Cheat:       true,
		Description: "Enable cheat actions",
	},
	"concede": {
		Action:      action.Action{Type: action.ActionConcede},
		Description: "Concede the game",
	},
	"conjure": {
		Action:      action.Action{Type: action.CheatConjure},
		Cheat:       true,
		Description: "Conjure a card",
	},
	"draw": {
		Action:      action.Action{Type: action.CheatDraw},
		Cheat:       true,
		Description: "Draw a card",
	},
	"discard": {
		Action:      action.Action{Type: action.CheatDiscard},
		Cheat:       true,
		Description: "Discard a card",
	},
	"find": {
		Action:      action.Action{Type: action.CheatFind},
		Cheat:       true,
		Description: "Find a card in the library",
	},
	"help": {
		// TODO: maybe have something here incase it some how gets passed to the game loop....
		Action:      action.Action{},
		Description: "Print more information about a command",
	},
	"landdrop": {
		Action:      action.Action{Type: action.CheatLandDrop},
		Cheat:       true,
		Description: "Reset land drop check",
	},
	"pass": {
		Action:      action.Action{Type: action.ActionPass},
		Description: "Pass the turn",
	},
	"peek": {
		Action:      action.Action{Type: action.CheatPeek},
		Cheat:       true,
		Description: "Peek at the top of the library",
	},
	"play": {
		Action:      action.Action{Type: action.ActionPlay},
		Description: "Play a land or cast a spell",
	},
	"shuffle": {
		Action:      action.Action{Type: action.CheatShuffle},
		Cheat:       true,
		Description: "Shuffle the library",
	},
	"untap": {
		Action:      action.Action{Type: action.ActionUntap},
		Description: "Untap a card",
	},
	"view": {
		Action:      action.Action{Type: action.ActionView},
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

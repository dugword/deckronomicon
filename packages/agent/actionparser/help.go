package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

// Command represents a command that can be executed in the game.
type CommandHelp struct {
	Cheat       bool
	Description string
}

// TODO: Add more commands, add stop, remove stop, set mode,

// Commands is a map of command names to their corresponding Command structs.
var Commands = map[string]CommandHelp{
	"addmana": {
		Cheat:       true,
		Description: "Add mana to your mana pool",
	},
	"activate": {
		Description: "Activate an ability",
	},
	"cheat": {
		Cheat:       true,
		Description: "Enable cheat actions",
	},
	"cast": {
		Description: engine.PlayCardAction{}.Description(),
	},
	"concede": {
		Description: "Concede the game",
	},
	"conjure": {
		Cheat:       true,
		Description: "Conjure a card",
	},
	"draw": {
		Cheat:       true,
		Description: "Draw a card",
	},
	"discard": {
		Cheat:       true,
		Description: "Discard a card",
	},
	"exit": {
		Description: engine.ConcedeAction{}.Description(),
	},
	"find": {
		Cheat:       true,
		Description: "Find a card in the library",
	},
	"help": {
		// TODO: maybe have something here incase it some how gets passed to the game loop....
		Description: "Print more information about a command",
	},
	"landdrop": {
		Cheat:       true,
		Description: "Reset land drop check",
	},
	"pass": {
		Description: engine.PassPriorityAction{}.Description(),
	},
	"peek": {
		Cheat:       true,
		Description: "Peek at the top of the library",
	},
	"play": {
		Description: engine.PlayCardAction{}.Description(),
	},
	"shuffle": {
		Cheat:       true,
		Description: "Shuffle the library",
	},
	"untap": {
		Description: "Untap a card",
	},
	"view": {
		Description: "View an object's description",
	},
}

type HelpCommand struct {
}

func (p *HelpCommand) IsComplete() bool {
	return true
}

func (p *HelpCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return nil, nil
}

package actionparser

// TODO: Document what level things live at. This package is for parsing user
// input or user actions from configuration files. It tries to accurately
// generate requests based on the game state, to be helpful and provide quick
// feed back to the user by only letting them make valid plays, but these are
// just requests to be sent to the game engine. THe game engine is responsible
// for actually verifying things work according to the rules.

// TODO Use a expression tree parser like we do for the JSON parser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
	"strings"
)

type CommandSource struct {
	name string
}

func (c CommandSource) Name() string {
	return c.name
}

type Command interface {
	IsComplete() bool
	Build(game state.Game, player state.Player) (engine.Action, error)
	// PromptNext(game state.Game, player state.Player) (choose.ChoicePrompt, error)
}

type CommandParser struct {
	Command Command
}

func (p *CommandParser) ParseInput(
	input string,
	choose func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
	game state.Game,
	player state.Player,
) (Command, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil, errors.New("no command provided")
	}
	command, arg := parts[0], strings.Join(parts[1:], " ")
	command = strings.ToLower(command)
	switch command {
	case "activate", "tap":
		return parseActivateAbilityCommand(arg, choose, game, player)
	case "cheat":
		return &CheatCommand{Player: player}, nil
	case "clear":
		return &ClearCommand{Player: player}, nil
	case "concede", "exit", "quit":
		return &ConcedeCommand{Player: player}, nil
	case "help":
		return &HelpCommand{}, nil
	case "pass", "next", "done":
		return &PassPriorityCommand{Player: player}, nil
	case "play":
		return parsePlayLandCommand(arg, choose, game, player)
	case "cast":
		return parseCastSpellCommand(arg, choose, game, player)
	case "view":
		return parseViewCommand(arg, choose, game, player)
	default:
		if game.CheatsEnabled() {
			return p.ParseCheatCommand(command, arg, choose, game, player)
		}
		return nil, fmt.Errorf("unknown command %q", command)
	}
}

func (p *CommandParser) ParseCheatCommand(
	command string,
	arg string,
	choose func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
	game state.Game,
	player state.Player,
) (Command, error) {
	switch command {
	case "addmana":
		return parseAddManaCheatCommand(arg, player)
	case "conjure":
		return parseConjureCardCheatCommand(arg, player)
	case "draw":
		return &DrawCheatCommand{Player: player}, nil
	case "discard":
		return parseDiscardCheatCommand(arg, choose, game, player)
	case "find", "tutor":
		return parseFindCardCheatCommand(arg, choose, player)
	case "landdrop":
		return &ResetLandDropCommand{Player: player}, nil
	case "peek":
		return &PeekCheatCommand{Player: player}, nil
	case "shuffle":
		return &ShuffleCheatCommand{Player: player}, nil
	case "untap":
		return parseUntapCheatCommand(arg, choose, game, player)
	default:
		return nil, fmt.Errorf("unknown cheat command %q", command)
	}
}

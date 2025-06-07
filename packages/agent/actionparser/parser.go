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
	Build(game state.Game, playerID string) (engine.Action, error)
	// PromptNext(game state.Game, playerID string) (choose.ChoicePrompt, error)
}

type CommandParser struct {
	Command Command
}

func (p *CommandParser) ParseInput(
	input string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	playerID string,
) (Command, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil, errors.New("no command provided")
	}
	command, args := parts[0], parts[1:]
	command = strings.ToLower(command)
	switch command {
	case "addmana":
		return parseAddManaCheatCommand(command, args, getChoices, game, playerID)
	case "activate", "tap":
		return parseActivateAbilityCommand(command, args, getChoices, game, playerID)
	case "cheat":
		return &CheatCommand{PlayerID: playerID}, nil
	case "concede", "exit", "quit":
		return &ConcedeCommand{playerID}, nil
	case "conjure":
		return parseConjureCardCheatCommand(command, args, getChoices, game, playerID)
	case "draw":
		return parseDrawCheatCommand(command, args, getChoices, game, playerID)
	case "discard":
		return parseDiscardCheatCommand(command, args, getChoices, game, playerID)
	case "find", "tutor":
		return parseFindCardCheatCommand(command, args, getChoices, game, playerID)
	case "help":
		return &HelpCommand{}, nil
	case "landdrop":
		return &ResetLandDropCommand{PlayerID: playerID}, nil
	case "pass", "next", "done":
		return &PassPriorityCommand{playerID}, nil
	case "peek":
		return parsePeekCheatCommand(command, args, getChoices, game, playerID)
	case "play", "cast":
		return parsePlayCardCommand(command, args, getChoices, game, playerID)
	case "shuffle":
		return &ShuffleCheatCommand{PlayerID: playerID}, nil
	case "untap":
		return parseUntapCheatCommand(command, args, getChoices, game, playerID)
	case "view":
		return parseViewCommand(command, args, getChoices, game, playerID)
	default:
		return nil, fmt.Errorf("unrecognized command '%s'", command)
	}
}

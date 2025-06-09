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
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	player state.Player,
) (Command, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil, errors.New("no command provided")
	}
	command, args := parts[0], parts[1:]
	command = strings.ToLower(command)
	switch command {
	case "addmana":
		return parseAddManaCheatCommand(command, args, getChoices, game, player)
	case "activate", "tap":
		return parseActivateAbilityCommand(command, args, getChoices, game, player)
	case "cheat":
		return &CheatCommand{Player: player}, nil
	case "concede", "exit", "quit":
		return &ConcedeCommand{Player: player}, nil
	case "conjure":
		return parseConjureCardCheatCommand(command, args, getChoices, game, player)
	case "draw":
		return parseDrawCheatCommand(command, args, getChoices, game, player)
	case "discard":
		return parseDiscardCheatCommand(command, args, getChoices, game, player)
	case "find", "tutor":
		return parseFindCardCheatCommand(command, args, getChoices, game, player)
	case "help":
		return &HelpCommand{}, nil
	case "landdrop":
		return &ResetLandDropCommand{Player: player}, nil
	case "pass", "next", "done":
		return &PassPriorityCommand{Player: player}, nil
	case "peek":
		return parsePeekCheatCommand(command, args, getChoices, game, player)
	case "play", "cast":
		return parsePlayCardCommand(command, args, getChoices, game, player)
	case "shuffle":
		return &ShuffleCheatCommand{Player: player}, nil
	case "untap":
		return parseUntapCheatCommand(command, args, getChoices, game, player)
	case "view":
		return parseViewCommand(command, args, getChoices, game, player)
	default:
		return nil, fmt.Errorf("unrecognized command '%s'", command)
	}
}

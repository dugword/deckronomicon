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
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
	"strings"
)

func ParseInput(
	input string,
	choose func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
	game state.Game,
	player state.Player,
) (engine.Action, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil, errors.New("no command provided")
	}
	command, arg := parts[0], strings.Join(parts[1:], " ")
	command = strings.ToLower(command)
	switch command {
	case "activate", "tap":
		return parseActivateAbilityCommand(arg, game, player, choose)
	case "cheat":
		return action.NewCheatAction(player), nil
	case "clear":
		return action.NewClearRevealedAction(player), nil
	case "concede", "exit", "quit":
		return action.NewConcedeAction(player), nil
	case "help":
		fmt.Println("Need to implement help command")
		return nil, nil
	case "pass", "next", "done":
		return action.NewPassPriorityAction(player), nil
	case "play":
		return parsePlayLandCommand(arg, game, player, choose)
	case "cast":
		return parseCastSpellCommand(arg, game, player, choose)
	case "view":
		return parseViewCommand(arg, game, player, choose)
	default:
		if game.CheatsEnabled() {
			return parseCheatCommand(command, arg, choose, game, player)
		}
		return nil, fmt.Errorf("unknown command %q", command)
	}
}

func parseCheatCommand(
	command string,
	arg string,
	choose func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
	game state.Game,
	player state.Player,
) (engine.Action, error) {
	switch command {
	case "addmana":
		return parseAddManaCheatCommand(arg, player)
	case "conjure":
		return parseConjureCardCheatCommand(arg, player)
	case "draw":
		return action.NewDrawCheatAction(player), nil
	case "discard":
		return parseDiscardCheatCommand(arg, game, player, choose)
	case "find", "tutor":
		return parseFindCardCheatCommand(arg, player, choose)
	case "landdrop":
		return action.NewResetLandDropCheatAction(player), nil
	case "peek":
		return action.NewPeekCheatAction(player), nil
	case "shuffle":
		return action.NewShuffleCheatAction(player), nil
	case "untap":
		return parseUntapCheatCommand(arg, game, player, choose)
	default:
		return nil, fmt.Errorf("unknown cheat command %q", command)
	}
}

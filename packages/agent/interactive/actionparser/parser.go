package actionparser

// TODO: Document what level things live at. This package is for parsing user
// input or user actions from configuration files. It tries to accurately
// generate requests based on the game state, to be helpful and provide quick
// feed back to the user by only letting them make valid plays, but these are
// just requests to be sent to the game engine. THe game engine is responsible
// for actually verifying things work according to the rules.

// TODO Use a expression tree parser like we do for the strategy parser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
	"strings"
)

// TODO: Think through having this return a request interface instead of an
// action. The request interface would be a common interface for all requests
// that can be built into an action.
// Need to decide how request.Build() is different from NewAction()
// Also need to decide if it validates and can return an error, or
// if it just builds the action and the action is responsible for validation.
func ParseInput(
	input string,
	agent engine.PlayerAgent,
	game *state.Game,
	playerID string,
	autoPay bool,
	autoPayColors []mana.Color,
) (engine.Action, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil, errors.New("no command provided")
	}
	command, arg := parts[0], strings.Join(parts[1:], " ")
	command = strings.ToLower(command)
	switch command {
	case "activate", "tap":
		request, err := parseActivateAbilityCommand(arg, game, playerID, agent, autoPay, autoPayColors)
		if err != nil {
			return nil, fmt.Errorf("failed to parse activate command: %w", err)
		}
		return request.Build(), nil
	case "cheat":
		return action.NewCheatAction(), nil
	case "clear":
		return action.NewClearRevealedAction(), nil
	case "concede", "exit", "quit":
		return action.NewConcedeAction(), nil
	case "emit":
		return parseEmitMetric(arg)
	case "help":
		fmt.Println("Need to implement help command")
		return nil, nil
	case "log":
		return parseLogMessage(arg)
	case "pass", "next", "done":
		return action.NewPassPriorityAction(), nil
	case "play":
		return parsePlayLandCommand(arg, game, playerID, agent)
	case "cast":
		request, err := parseCastSpellCommand(arg, game, playerID, agent, autoPay, autoPayColors)
		if err != nil {
			return nil, fmt.Errorf("failed to parse cast command: %w", err)
		}
		return request.Build(), nil
	case "view":
		return parseViewCommand(arg, game, playerID, agent)
	default:
		if game.CheatsEnabled() {
			return parseCheatCommand(command, arg, game, playerID, agent)
		}
		return nil, fmt.Errorf("unknown command %q", command)
	}
}

func parseCheatCommand(
	command string,
	arg string,
	game *state.Game,
	playerID string,
	agent engine.PlayerAgent,
) (engine.Action, error) {
	switch command {
	case "addmana":
		return parseAddManaCheatCommand(arg, playerID)
	case "conjure":
		return parseConjureCardCheatCommand(arg)
	case "draw":
		return action.NewDrawCheatAction(), nil
	case "discard":
		return parseDiscardCheatCommand(arg, game, playerID, agent)
	case "effect":
		return parseEffectCheatCommand(arg)
	case "find", "tutor":
		return parseFindCardCheatCommand(arg, game, playerID, agent)
	case "landdrop":
		return action.NewResetLandDropCheatAction(), nil
	case "peek":
		return action.NewPeekCheatAction(playerID), nil
	case "shuffle":
		return action.NewShuffleCheatAction(playerID), nil
	case "untap":
		return parseUntapCheatCommand(arg, game, playerID, agent)
	default:
		return nil, fmt.Errorf("unknown cheat command %q", command)
	}
}

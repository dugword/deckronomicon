package actionparser

import (
	"deckronomicon/packages/engine/action"
	"fmt"
	"strings"
)

func parseEffectCheatCommand(
	input string,
) (action.EffectCheatAction, error) {
	if input == "" {
		return action.EffectCheatAction{}, fmt.Errorf("effect command requires an effect name string")
	}
	parts := strings.SplitN(input, " ", 2)
	if len(parts) != 2 {
		return action.EffectCheatAction{}, fmt.Errorf("effect command requires an effect name and modifiers")
	}
	effectName := parts[0]
	modifiers := parts[1]
	return action.NewEffectCheatAction(effectName, modifiers), nil
}

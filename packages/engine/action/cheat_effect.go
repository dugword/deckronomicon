package action

import (
	"deckronomicon/packages/engine/effect"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/state"
	"encoding/json"
	"fmt"
)

type EffectCheatAction struct {
	Player     state.Player
	EffectName string
	Modifiers  string
}

func NewEffectCheatAction(player state.Player, effectName string, modifers string) EffectCheatAction {
	return EffectCheatAction{
		Player:     player,
		EffectName: effectName,
		Modifiers:  modifers,
	}
}

func (a EffectCheatAction) Name() string {
	return fmt.Sprintf("CHEAT: %s", a.EffectName)
}

func (a EffectCheatAction) Complete(game state.Game, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	player := game.GetPlayer(a.Player.ID())
	if a.EffectName == "" {
		return nil, fmt.Errorf("effect action is missing effect name")
	}
	effectSpec := definition.EffectSpec{
		Name:      a.EffectName,
		Modifiers: json.RawMessage(a.Modifiers),
	}
	efct, err := effect.Build(effectSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to get effect %q: %w", a.EffectName, err)
	}
	effectResults, err := efct.Resolve(game, player, nil, target.TargetValue{}, resEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve effect %q: %w", a.EffectName, err)
	}
	if effectResults.ChoicePrompt.ChoiceOpts != nil {
		return nil, fmt.Errorf("effect %q has unresolved choices", a.EffectName)
	}
	return effectResults.Events, nil
}

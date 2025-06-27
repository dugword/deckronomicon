package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
	"fmt"
)

type EffectCheatAction struct {
	EffectName string
	Modifiers  string
}

func NewEffectCheatAction(effectName string, modifers string) EffectCheatAction {
	return EffectCheatAction{
		EffectName: effectName,
		Modifiers:  modifers,
	}
}

func (a EffectCheatAction) Name() string {
	return fmt.Sprintf("CHEAT: %s", a.EffectName)
}

func (a EffectCheatAction) Complete(game state.Game, player state.Player, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	panic("EffectCheatAction is not implemented yet")
	/*
		if !game.CheatsEnabled() {
			return nil, fmt.Errorf("no cheating you cheater")
		}
		if a.EffectName == "" {
			return nil, fmt.Errorf("effect action is missing effect name")
		}
		var modifiers map[string]any
		if err := json.Unmarshal([]byte(a.Modifiers), &modifiers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal modifiers for effect %q: %w", a.EffectName, err)
		}
		effect := gob.Effect{
			Name:      a.EffectName,
			Modifiers: modifiers,
		}
		effectResolver, err := resolver.Build(effectSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to get effect %q: %w", a.EffectName, err)
		}
		effectResults, err := effectResolver.Resolve(game, player, nil, gob.Target{}, resEnv)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve effect %q: %w", a.EffectName, err)
		}
		if effectResults.ChoicePrompt.ChoiceOpts != nil {
			return nil, fmt.Errorf("effect %q has unresolved choices", a.EffectName)
		}
		return effectResults.Events, nil
	*/
}

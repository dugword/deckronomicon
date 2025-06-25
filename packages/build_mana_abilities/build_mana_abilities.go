package buildmanaabilities

import (
	"deckronomicon/packages/engine/effect"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/state"
	"fmt"
)

// TODO: Putting this here to fix a circular dependency
// but this is a really silly place for it.
// It should be somewhere in a package that does
// other stuff

func BuildManaAbilityEvents(
	game state.Game,
	player state.Player,
	// effectSpecs []definition.EffectSpec,
	effectWithTargets []target.EffectWithTarget,
	resEnv *resenv.ResEnv,
) ([]event.GameEvent, error) {
	var events []event.GameEvent
	// TODO: This could fail if the effect has a resume func
	for _, effectWithTarget := range effectWithTargets {
		efct, err := effect.Build(effectWithTarget.EffectSpec)
		if err != nil {
			return nil, fmt.Errorf("effect %q not found: %w", effectWithTarget.EffectSpec.Name, err)
		}
		effectResults, err := efct.Resolve(game, player, nil, effectWithTarget.Target, resEnv)
		if err != nil {
			return nil, fmt.Errorf("failed to apply effect %q: %w", effectWithTarget.EffectSpec.Name, err)
		}
		events = append(events, effectResults.Events...)
	}
	return events, nil
}

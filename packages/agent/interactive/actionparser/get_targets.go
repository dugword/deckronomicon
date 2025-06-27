package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"fmt"
)

// TODO: Move to judge
// TODO: This probably should return 3 values instead of a map,
// key, value, error
// and then build the map in the caller.
func getTargetsForEffects(
	object gob.Object,
	Effects []effect.Effect,
	game state.Game,
	agent engine.PlayerAgent,
) (map[effect.EffectTargetKey]effect.Target, error) {
	targetsForEffects := map[effect.EffectTargetKey]effect.Target{}
	for i, efct := range Effects {
		effectTargetKey := effect.EffectTargetKey{
			SourceID:    object.ID(),
			EffectIndex: i,
		}
		switch targetSpec := efct.TargetSpec().(type) {
		case nil, effect.NoneTargetSpec:
			targetsForEffects[effectTargetKey] = effect.Target{
				Type: mtg.TargetTypeNone,
			}
		case effect.PlayerTargetSpec:
			playerTarget, err := getPlayerTarget(
				object,
				game,
				agent,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to get player target: %w", err)
			}
			targetsForEffects[effectTargetKey] = playerTarget
		case effect.SpellTargetSpec:
			spellTarget, err := getSpellTarget(
				targetSpec,
				object,
				game,
				agent,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to get spell target: %w", err)
			}
			targetsForEffects[effectTargetKey] = spellTarget
		case effect.PermanentTargetSpec:
			permanentTarget, err := getPermanentTarget(
				targetSpec,
				object,
				game,
				agent,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to get permanent target: %w", err)
			}
			targetsForEffects[effectTargetKey] = permanentTarget
		default:
			return nil, fmt.Errorf("unsupported target spec type %T", targetSpec)
		}

	}
	return targetsForEffects, nil
}

func getPlayerTarget(
	object gob.Object,
	game state.Game,
	agent engine.PlayerAgent,
) (effect.Target, error) {
	prompt := choose.ChoicePrompt{
		Message: "Choose a player to target",
		Source:  object,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(game.Players()),
		},
	}
	choiceResults, err := agent.Choose(prompt)
	if err != nil {
		return effect.Target{}, fmt.Errorf("failed to get choice results: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return effect.Target{}, fmt.Errorf("expected a single choice result")
	}
	selectedPlayer, ok := selected.Choice.(state.Player)
	if !ok {
		return effect.Target{}, fmt.Errorf("selected choice is not a player")
	}
	return effect.Target{
		Type: mtg.TargetTypePlayer,
		ID:   selectedPlayer.ID(),
	}, nil
}

// TODO: Double check this works, I moved it from the counterspell effect
// and it was not tested here.
func getSpellTarget(
	targetSpec effect.SpellTargetSpec,
	object gob.Object,
	game state.Game,
	agent engine.PlayerAgent,
) (effect.Target, error) {
	query, err := buildQuery(QueryOpts(targetSpec))
	if err != nil {
		panic(fmt.Errorf("failed to build query for Search effect: %w", err))
	}
	spells := game.Stack().FindAll(query)
	prompt := choose.ChoicePrompt{
		Message: "Choose a spell to target",
		Source:  object,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(spells),
		},
	}
	choiceResults, err := agent.Choose(prompt)
	if err != nil {
		return effect.Target{}, fmt.Errorf("failed to get choice results: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return effect.Target{}, fmt.Errorf("expected a single choice result")
	}
	selectedSpell, ok := selected.Choice.(gob.Spell)
	if !ok {
		return effect.Target{}, fmt.Errorf("selected choice is not a spell")
	}
	return effect.Target{
		Type: mtg.TargetTypeSpell,
		ID:   selectedSpell.ID(),
	}, nil
}

func getPermanentTarget(
	targetSpec effect.PermanentTargetSpec,
	object gob.Object,
	game state.Game,
	agent engine.PlayerAgent,
) (effect.Target, error) {
	permanents := game.Battlefield().GetAll()
	prompt := choose.ChoicePrompt{
		Message: "Choose a permanent to target",
		Source:  object,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(permanents),
		},
	}
	choiceResults, err := agent.Choose(prompt)
	if err != nil {
		return effect.Target{}, fmt.Errorf("failed to get choice results: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return effect.Target{}, fmt.Errorf("expected a single choice result")
	}
	selectedPermanent, ok := selected.Choice.(gob.Permanent)
	if !ok {
		return effect.Target{}, fmt.Errorf("selected choice is not a permanent")
	}
	return effect.Target{
		Type: mtg.TargetTypePermanent,
		ID:   selectedPermanent.ID(),
	}, nil
}

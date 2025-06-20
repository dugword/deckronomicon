package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/engine/effect"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"fmt"
)

// TODO: This probably should return 3 values instead of a map,
// key, value, error
// and then build the map in the caller.
func getTargetsForEffects(
	object query.Object,
	effectSpecs []definition.EffectSpec,
	game state.Game,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
) (map[action.EffectTargetKey]target.TargetValue, error) {
	targetsForEffects := map[action.EffectTargetKey]target.TargetValue{}
	for i, effectSpec := range effectSpecs {
		effectTargetKey := action.EffectTargetKey{
			SourceID:    object.ID(),
			EffectIndex: i,
		}
		efct, err := effect.Build(effectSpec)
		if err != nil {
			return nil, fmt.Errorf("effect %q not found: %w", effectSpec.Name, err)
		}
		switch targetSpec := efct.TargetSpec().(type) {
		case nil, target.NoneTargetSpec:
			fmt.Println("No target spec for effect", effectSpec.Name, "on object", object.Name())
			targetsForEffects[effectTargetKey] = target.TargetValue{
				TargetType: target.TargetTypeNone,
			}
		case target.PlayerTargetSpec:
			playerTarget, err := getPlayerTarget(
				object,
				game,
				chooseFunc,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to get player target: %w", err)
			}
			targetsForEffects[effectTargetKey] = playerTarget
		case target.SpellTargetSpec:
			spellTarget, err := getSpellTarget(
				targetSpec,
				object,
				game,
				chooseFunc,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to get spell target: %w", err)
			}
			targetsForEffects[effectTargetKey] = spellTarget
		case target.PermanentTargetSpec:
			permanentTarget, err := getPermanentTarget(
				targetSpec,
				object,
				game,
				chooseFunc,
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
	object query.Object,
	game state.Game,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
) (target.TargetValue, error) {
	prompt := choose.ChoicePrompt{
		Message: "Choose a player to target",
		Source:  object,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(game.Players()),
		},
	}
	choiceResults, err := chooseFunc(prompt)
	if err != nil {
		return target.TargetValue{}, fmt.Errorf("failed to get choice results: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return target.TargetValue{}, fmt.Errorf("expected a single choice result")
	}
	selectedPlayer, ok := selected.Choice.(state.Player)
	if !ok {
		return target.TargetValue{}, fmt.Errorf("selected choice is not a player")
	}
	return target.TargetValue{
		TargetType: target.TargetTypePlayer,
		PlayerID:   selectedPlayer.ID(),
	}, nil
}

func getSpellTarget(
	targetSpec target.SpellTargetSpec,
	object query.Object,
	game state.Game,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
) (target.TargetValue, error) {
	spells := game.Stack().FindAll(targetSpec.Predicate)
	prompt := choose.ChoicePrompt{
		Message: "Choose a spell to target",
		Source:  object,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(spells),
		},
	}
	choiceResults, err := chooseFunc(prompt)
	if err != nil {
		return target.TargetValue{}, fmt.Errorf("failed to get choice results: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return target.TargetValue{}, fmt.Errorf("expected a single choice result")
	}
	selectedSpell, ok := selected.Choice.(gob.Spell)
	if !ok {
		return target.TargetValue{}, fmt.Errorf("selected choice is not a spell")
	}
	return target.TargetValue{
		TargetType: target.TargetTypeSpell,
		ObjectID:   selectedSpell.ID(),
	}, nil
}

func getPermanentTarget(
	targetSpec target.PermanentTargetSpec,
	object query.Object,
	game state.Game,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
) (target.TargetValue, error) {
	permanents := game.Battlefield().GetAll()
	prompt := choose.ChoicePrompt{
		Message: "Choose a permanent to target",
		Source:  object,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(permanents),
		},
	}
	choiceResults, err := chooseFunc(prompt)
	if err != nil {
		return target.TargetValue{}, fmt.Errorf("failed to get choice results: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return target.TargetValue{}, fmt.Errorf("expected a single choice result")
	}
	selectedPermanent, ok := selected.Choice.(gob.Permanent)
	if !ok {
		return target.TargetValue{}, fmt.Errorf("selected choice is not a permanent")
	}
	return target.TargetValue{
		TargetType: target.TargetTypePermanent,
		ObjectID:   selectedPermanent.ID(),
	}, nil
}

package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/query/querybuilder"
	"deckronomicon/packages/state"
	"fmt"
)

func getTargetsForCost(
	playerID string,
	card *gob.Card,
	game *state.Game,
	agent engine.PlayerAgent,
) (target.Target, error) {
	if card.AdditionalCost() == nil {
		return target.Target{}, nil // No additional costs to pay
	}
	additionalCost := card.AdditionalCost()
	costWithTargets, ok := additionalCost.(cost.CostWithTarget)
	if !ok {
		fmt.Println("WARN: If composite cost has a cost with targets, this will miss it")
		return target.Target{}, nil
	}
	switch targetSpec := costWithTargets.TargetSpec().(type) {
	case nil, target.NoneTargetSpec:
		return target.Target{}, nil // No targets needed for NoneTargetSpec
	case target.CardTargetSpec:
		return getCardTarget(playerID, targetSpec, card, game, agent)
	case target.PermanentTargetSpec:
		return getPermanentTarget(targetSpec, card, game, agent)
	default:
		return target.Target{}, fmt.Errorf("unsupported target spec type %T", costWithTargets.TargetSpec())
	}
}

// TODO: Move to judge
// TODO: This probably should return 3 values instead of a map,
// key, value, error
// and then build the map in the caller.
func getTargetsForEffects(
	playerID string,
	object gob.Object,
	Effects []effect.Effect,
	game *state.Game,
	agent engine.PlayerAgent,
) (map[effect.EffectTargetKey]target.Target, error) {
	targetsForEffects := map[effect.EffectTargetKey]target.Target{}
	for i, efct := range Effects {
		effectTargetKey := effect.EffectTargetKey{
			SourceID:    object.ID(),
			EffectIndex: i,
		}
		switch targetSpec := efct.TargetSpec().(type) {
		case nil, target.NoneTargetSpec:
			targetsForEffects[effectTargetKey] = target.Target{
				Type: mtg.TargetTypeNone,
			}
		case target.PlayerTargetSpec:
			playerTarget, err := getPlayerTarget(
				object,
				game,
				agent,
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
				agent,
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
	game *state.Game,
	agent engine.PlayerAgent,
) (target.Target, error) {
	prompt := choose.ChoicePrompt{
		Message: "Choose a player to target",
		Source:  object,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(game.Players()),
		},
	}
	choiceResults, err := agent.Choose(prompt)
	if err != nil {
		return target.Target{}, fmt.Errorf("failed to get choice results: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return target.Target{}, fmt.Errorf("expected a single choice result")
	}
	selectedPlayer, ok := selected.Choice.(*state.Player)
	if !ok {
		return target.Target{}, fmt.Errorf("selected choice is not a player")
	}
	return target.Target{
		Type: mtg.TargetTypePlayer,
		ID:   selectedPlayer.ID(),
	}, nil
}

// TODO: Double check this works, I moved it from the counterspell effect
// and it was not tested here.
func getSpellTarget(
	targetSpec target.SpellTargetSpec,
	object gob.Object,
	game *state.Game,
	agent engine.PlayerAgent,
) (target.Target, error) {
	query, err := querybuilder.Build(querybuilder.Opts(targetSpec))
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
		return target.Target{}, fmt.Errorf("failed to get choice results: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return target.Target{}, fmt.Errorf("expected a single choice result")
	}
	selectedSpell, ok := selected.Choice.(*gob.Spell)
	if !ok {
		return target.Target{}, fmt.Errorf("selected choice is not a spell")
	}
	return target.Target{
		Type: mtg.TargetTypeSpell,
		ID:   selectedSpell.ID(),
	}, nil
}

func getPermanentTarget(
	targetSpec target.PermanentTargetSpec,
	object gob.Object,
	game *state.Game,
	agent engine.PlayerAgent,
) (target.Target, error) {
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
		return target.Target{}, fmt.Errorf("failed to get choice results: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return target.Target{}, fmt.Errorf("expected a single choice result")
	}
	selectedPermanent, ok := selected.Choice.(*gob.Permanent)
	if !ok {
		return target.Target{}, fmt.Errorf("selected choice is not a permanent")
	}
	return target.Target{
		Type: mtg.TargetTypePermanent,
		ID:   selectedPermanent.ID(),
	}, nil
}

func getCardTarget(
	playerID string,
	targetSpec target.CardTargetSpec,
	object gob.Object,
	game *state.Game,
	agent engine.PlayerAgent,
) (target.Target, error) {
	player := game.GetPlayer(playerID)
	predicate, err := querybuilder.Build(querybuilder.Opts(targetSpec))
	if err != nil {
		panic(fmt.Errorf("failed to build query for Search effect: %w", err))
	}
	cards := player.Hand().FindAll(query.And(predicate, is.Not(has.ID(object.ID()))))
	prompt := choose.ChoicePrompt{
		Message: "Choose a card to target",
		Source:  object,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(cards),
		},
	}
	choiceResults, err := agent.Choose(prompt)
	if err != nil {
		return target.Target{}, fmt.Errorf("failed to get choice results: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return target.Target{}, fmt.Errorf("expected a single choice result")
	}
	selectedCard, ok := selected.Choice.(*gob.Card)
	if !ok {
		return target.Target{}, fmt.Errorf("selected choice is not a card")
	}
	return target.Target{
		Type: mtg.TargetTypeCard,
		ID:   selectedCard.ID(),
	}, nil
}

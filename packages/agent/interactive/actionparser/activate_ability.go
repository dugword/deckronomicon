package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/engine/judge"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"fmt"
)

func parseActivateAbilityCommand(
	idOrName string,
	game state.Game,
	player state.Player,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
) (action.ActivateAbilityAction, error) {
	var abilityInZone gob.AbilityInZone
	var err error
	ruling := judge.Ruling{Explain: true}
	abilityInZones := judge.GetAbilitiesAvailableToActivate(game, player, &ruling)
	if idOrName == "" {
		abilityInZone, err = buildActivateAbilityCommandByChoice(abilityInZones, chooseFunc, player)
		if err != nil {
			return action.ActivateAbilityAction{}, fmt.Errorf("failed to choose an ability to activate: %w", err)
		}
	} else {
		found, ok := query.Find(abilityInZones, query.Or(has.ID(idOrName), has.Name(idOrName)))
		if !ok {
			return action.ActivateAbilityAction{}, fmt.Errorf("failed to find ability with ID or name %q: %w", idOrName, err)
		}
		abilityInZone = found
	}
	targetsForEffects, err := getTargetsForEffects(
		abilityInZone.Ability(),
		abilityInZone.Ability().EffectSpecs(),
		game,
		chooseFunc,
	)
	if err != nil {
		return action.ActivateAbilityAction{}, fmt.Errorf("failed to get targets for spell: %w", err)
	}
	fmt.Println("TargetsForEffects")
	for key, value := range targetsForEffects {
		fmt.Printf("  %v: %v\n", key, value)
	}
	fmt.Println("Length of targetsForEffects:", len(targetsForEffects))
	request := action.ActivateAbilityRequest{
		AbilityID:         abilityInZone.Ability().ID(),
		SourceID:          abilityInZone.Source().ID(),
		Zone:              abilityInZone.Zone(),
		TargetsForEffects: targetsForEffects,
	}
	return action.NewActivateAbilityAction(player.ID(), request), nil
}

func buildActivateAbilityCommandByChoice(
	abilities []gob.AbilityInZone,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
	player state.Player,
) (gob.AbilityInZone, error) {
	prompt := choose.ChoicePrompt{
		Message:  "Choose an ability to activate",
		Source:   player,
		Optional: true,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(abilities),
		},
	}
	choiceResults, err := chooseFunc(prompt)
	if err != nil {
		return gob.AbilityInZone{}, fmt.Errorf("failed to get choices: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return gob.AbilityInZone{}, fmt.Errorf("expected choose one results, got %T", selected)
	}
	ability, ok := selected.Choice.(gob.AbilityInZone)
	if !ok {
		return gob.AbilityInZone{}, fmt.Errorf("selected choice is not an ability on an object in a zone")
	}
	return ability, nil
}

package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/engine/judge"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"fmt"
)

func parseActivateAbilityCommand(
	idOrName string,
	game *state.Game,
	playerID string,
	agent engine.PlayerAgent,
	autoPayCost bool,
	autoPayColors []mana.Color, // Colors to prioritize when auto-paying costs, if applicable
) (action.ActivateAbilityRequest, error) {
	var abilityInZone *gob.AbilityInZone
	var err error
	ruling := judge.Ruling{Explain: true}
	abilityInZones := judge.GetAbilitiesAvailableToActivate(game, playerID, autoPayCost, autoPayColors, &ruling)
	if idOrName == "" {
		abilityInZone, err = buildActivateAbilityCommandByChoice(abilityInZones, agent, playerID, game)
		if err != nil {
			return action.ActivateAbilityRequest{}, fmt.Errorf("failed to choose an ability to activate: %w", err)
		}
	} else {
		found, ok := query.Find(abilityInZones, query.Or(has.ID(idOrName), has.Name(idOrName)))
		if !ok {
			return action.ActivateAbilityRequest{}, fmt.Errorf("failed to find ability with ID or name %q: %w", idOrName, ErrAbilityNotFound)
		}
		abilityInZone = found
	}
	targetsForEffects, err := getTargetsForEffects(
		playerID,
		abilityInZone.Ability(),
		abilityInZone.Ability().Effects(),
		game,
		agent,
	)
	if err != nil {
		return action.ActivateAbilityRequest{}, fmt.Errorf("failed to get targets for spell: %w", err)
	}
	request := action.ActivateAbilityRequest{
		AbilityID:         abilityInZone.Ability().ID(),
		SourceID:          abilityInZone.Source().ID(),
		Zone:              abilityInZone.Zone(),
		TargetsForEffects: targetsForEffects,
		AutoPayCost:       autoPayCost,
		AutoPayColors:     autoPayColors,
	}
	return request, nil
}

func buildActivateAbilityCommandByChoice(
	abilities []*gob.AbilityInZone,
	agent engine.PlayerAgent,
	playerID string,
	game *state.Game,
) (*gob.AbilityInZone, error) {
	prompt := choose.ChoicePrompt{
		Message:  "Choose an ability to activate",
		Source:   nil, // TODO: Make this better
		Optional: true,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(abilities),
		},
	}
	choiceResults, err := agent.Choose(game, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get choices: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return nil, fmt.Errorf("expected choose one results, got %T", selected)
	}
	if selected.Choice == nil {
		return nil, fmt.Errorf("no card selected: %w", choose.ErrNoChoiceSelected)
	}
	ability, ok := selected.Choice.(*gob.AbilityInZone)
	if !ok {
		return nil, fmt.Errorf("selected choice is not an ability on an object in a zone")
	}
	return ability, nil
}

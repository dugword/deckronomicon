package action

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/judge"
	"deckronomicon/packages/state"
	"fmt"
)

type ActivateAbilityAction struct {
	abilityOnObjectInZone gob.AbilityInZone
	player                state.Player
}

func NewActivateAbilityAction(player state.Player, abilityOnObjectInZone gob.AbilityInZone) ActivateAbilityAction {
	return ActivateAbilityAction{
		abilityOnObjectInZone: abilityOnObjectInZone,
		player:                player,
	}
}

func (a ActivateAbilityAction) PlayerID() string {
	return a.player.ID()
}

func (a ActivateAbilityAction) Name() string {
	return "Activate Ability"
}

func (a ActivateAbilityAction) Description() string {
	return "The active player activates an ability."
}

func (a ActivateAbilityAction) GetPrompt(game state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{}, nil
}

func (a ActivateAbilityAction) Complete(
	game state.Game,
	choiceResults choose.ChoiceResults,
) ([]event.GameEvent, error) {
	ability := a.abilityOnObjectInZone.Ability()
	cost, err := cost.ParseCost(ability.Cost(), a.abilityOnObjectInZone.Ability().Source())
	if err != nil {
		return nil, fmt.Errorf("failed to parse cost %q: %w", ability.Cost(), err)
	}
	if !judge.CanPayCost(cost, ability.Source(), game, a.player) {
		return nil, fmt.Errorf("player %q cannot pay cost %q", a.player.ID(), cost.Description())
	}
	events := []event.GameEvent{
		event.ActivateAbilityEvent{
			PlayerID:  a.player.ID(),
			SourceID:  ability.Source().ID(),
			AbilityID: ability.Name(),
			Zone:      a.abilityOnObjectInZone.Zone(),
		},
	}
	costEvents, err := PayCost(cost, ability.Source(), a.player)
	if err != nil {
		return nil, fmt.Errorf("failed to pay cost %q: %w", cost.Description(), err)
	}
	events = append(events, costEvents...)
	// THIS SHOULD MOVE TO THE APPLY EFFECTS PHASE
	if ability.IsManaAbility() {
		events = append(events, event.ResolveManaAbilityEvent{
			PlayerID:    a.player.ID(),
			SourceID:    ability.Source().ID(),
			AbilityID:   ability.ID(),
			FromZone:    a.abilityOnObjectInZone.Zone(),
			AbilityName: ability.Name(),
			Effects:     ability.Effects(),
		})
		return events, nil
	}
	events = append(events, event.PutAbilityOnStackEvent{
		PlayerID:    a.player.ID(),
		SourceID:    ability.Source().ID(),
		AbilityID:   ability.ID(),
		FromZone:    a.abilityOnObjectInZone.Zone(),
		AbilityName: ability.Name(),
		Effects:     ability.Effects(),
	})
	return events, nil
}

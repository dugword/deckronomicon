package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/judge"
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

func (a ActivateAbilityAction) Complete(game state.Game) ([]event.GameEvent, error) {
	ability := a.abilityOnObjectInZone.Ability()
	ruling := judge.Ruling{Explain: true}
	// TODO: Should it be abilityOnObjectInZone.Object() instead of Source()?
	if !judge.CanActivateAbility(game, a.player, a.abilityOnObjectInZone.Source(), ability, &ruling) {
		return nil, fmt.Errorf(
			"player %q cannot activate ability %q on %q in %q: %v",
			a.player.ID(),
			ability.Name(),
			a.abilityOnObjectInZone.Source().Name(),
			a.abilityOnObjectInZone.Zone(),
			ruling.Reasons,
		)
	}
	events := []event.GameEvent{
		event.ActivateAbilityEvent{
			PlayerID:  a.player.ID(),
			SourceID:  ability.Source().ID(),
			AbilityID: ability.Name(),
			Zone:      a.abilityOnObjectInZone.Zone(),
		},
	}
	cost, err := cost.ParseCost(ability.Cost(), a.abilityOnObjectInZone.Ability().Source())
	if err != nil {
		return nil, fmt.Errorf("failed to parse cost %q: %w", ability.Cost(), err)
	}
	costEvents, err := PayCost(cost, ability.Source(), a.player)
	if err != nil {
		return nil, fmt.Errorf("failed to pay cost %q: %w", cost.Description(), err)
	}
	events = append(events, costEvents...)
	// TODO: I think I need to do this on a per effect basis, like I think for
	// chromatic star I get the mana fixing buy the draw a card effect still
	// goes on the stack.
	if ability.IsManaAbility() {
		events = append(events, event.ResolveManaAbilityEvent{
			PlayerID:    a.player.ID(),
			SourceID:    ability.Source().ID(),
			AbilityID:   ability.ID(),
			FromZone:    a.abilityOnObjectInZone.Zone(),
			AbilityName: ability.Name(),
			EffectSpecs: ability.EffectSpecs(),
		})
		return events, nil
	}
	events = append(events, event.PutAbilityOnStackEvent{
		PlayerID:    a.player.ID(),
		SourceID:    ability.Source().ID(),
		AbilityID:   ability.ID(),
		FromZone:    a.abilityOnObjectInZone.Zone(),
		AbilityName: ability.Name(),
		EffectSpecs: ability.EffectSpecs(),
	})
	return events, nil
}

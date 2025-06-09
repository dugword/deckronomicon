package engine

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
	return choose.ChoicePrompt{
		Message:  "Activating an ability",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a ActivateAbilityAction) Complete(
	game state.Game,
	env *ResolutionEnvironment,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	ability := a.abilityOnObjectInZone.Ability()

	cost, err := cost.ParseCost(ability.Cost(), a.abilityOnObjectInZone.Ability().Source())
	if err != nil {
		return nil, fmt.Errorf("failed to parse cost: %w", err)
	}
	if !judge.CanPayCost(cost, game, a.player) {
		return nil, fmt.Errorf("player '%s' cannot pay cost '%s'", a.player.ID(), cost.Description())
	}
	costEvents, err := PayCost(cost, game, a.player)
	if err != nil {
		return nil, fmt.Errorf("failed to pay cost: %w", err)
	}
	var effectsEvents []event.GameEvent
	fmt.Printf("Generating %d effects for ability '%s'\n", len(ability.Effects()), ability.Name())
	for _, effect := range ability.Effects() {
		handler, ok := env.EffectRegistry.Get(effect.Name())
		if !ok {
			return nil, fmt.Errorf("effect '%s' not found in registry", effect.Name())
		}
		effectEvents, err := handler(game, a.player, ability.Source(), effect.Modifiers())
		if err != nil {
			return nil, fmt.Errorf("failed to apply effect '%s': %w", effect.Name(), err)
		}
		fmt.Println("Effect events generated:", len(effectEvents))
		effectsEvents = append(effectsEvents, effectEvents...)
	}

	return append(costEvents, effectsEvents...), nil
}

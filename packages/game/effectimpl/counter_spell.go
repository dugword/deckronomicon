package effectimpl

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/game/core"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/spell"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"fmt"
	"strings"
)

// BuildEffectCounterSpell creates an effect that counters a spell.
// Supported Modifier Keys (last applies):
//   - Target: <target> CardType
//
// Multiple targets can be specified and will be OR'd together.
// If no target is specified, the effect will counter any spell.
func BuildEffectCounterSpell(source query.Object, spec definition.EffectSpec) (*effect.Effect, error) {

	var targetTypes []string
	for _, modifier := range spec.Modifiers {
		if modifier.Key == "Target" {
			targetTypes = append(targetTypes, modifier.Value)
		}
	}
	var cardTypes []mtg.CardType
	if len(targetTypes) != 0 {
		for _, target := range targetTypes {
			cardType, err := mtg.StringToCardType(target)
			if err != nil {
				return nil, fmt.Errorf("invalid target card type: %s", target)
			}
			cardTypes = append(cardTypes, cardType)
		}
	}
	apply := func(state core.State, player core.Player) error {
		s, ok := source.(State)
		if !ok {
			return fmt.Errorf("source is not a valid state: %T", source)
		}
		resolvables := s.Stack().GetAll()
		var spells []query.Object
		for _, resolvable := range resolvables {

			spell, ok := resolvable.(*spell.Spell)
			if !ok {
				continue
			}
			if len(cardTypes) == 0 {
				spells = append(spells, spell)
			}
			for _, cardType := range cardTypes {
				if spell.Match(has.CardType(cardType)) {
					spells = append(spells, spell)
					break // No need to check other types if one matches
				}
			}

		}
		choices := choose.CreateChoices(spells, mtg.ZoneStack)
		if len(choices) == 0 {
			return fmt.Errorf("no spells to counter")
		}
		/*
			chosen, err := player.Agent.ChooseOne(
				"Choose a spell to counter",
				source,
				choices,
			)
			if err != nil {
				return fmt.Errorf("failed to choose spell: %w", err)
			}
			// Ensure the spell is on the stack
			if _, err := state.Stack.Get(chosen.ID); err != nil {
				// TODO: Handle fizzling consistently
				state.Log("spell fizzled - no targets")
				return nil
			}
			object, err := state.Stack.Take(chosen.ID)
			if err != nil {
				return fmt.Errorf("failed to remove spell from stack: %w", err)
			}
			spell, ok := object.(*Spell)
			if !ok {
				return fmt.Errorf("object is not a spell: %s", object.ID())
			}
			player.Graveyard.Add(spell.Card())
		*/
		return nil
	}
	var tags []core.Tag
	for _, target := range targetTypes {
		tags = append(tags, core.Tag{Key: "CounterSpell", Value: target})
	}

	eff := effect.New(
		spec.ID,
		fmt.Sprintf("counter a spell of type %s", strings.Join(targetTypes, ", ")),
		tags,

		apply,
	)
	return eff, nil
}

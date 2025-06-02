package effectimpl

import (
	"deckronomicon/packages/game/core"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"errors"
	"fmt"
)

// BuildEffectAddMana creates an effect that adds mana to the player's mana
// pool.
// Supported Modifier Keys (concats multiple modifiers):
//   - Mana: <ManaString>
func BuildEffectAddMana(source query.Object, spec definition.EffectSpec) (*effect.Effect, error) {

	var mana string
	for _, modifier := range spec.Modifiers {
		if modifier.Key == "Mana" {
			mana += modifier.Value
		}
	}
	if mana == "" {
		return nil, errors.New("no mana string provided")
	}
	if !mtg.IsMana(mana) {
		return nil, fmt.Errorf("invalid mana string: %s", mana)
	}
	var tags []core.Tag
	for _, symbol := range mtg.ManaStringToManaSymbols(mana) {
		tags = append(tags, core.Tag{Key: effect.TagManaAbility, Value: symbol})
	}
	eff := effect.New(
		spec.ID,
		fmt.Sprintf("add %s", mana),
		tags,

		func(state core.State, player core.Player) error {
			if err := player.AddMana(mana); err != nil {
				return fmt.Errorf("failed to add mana: %w", err)
			}
			return nil
		},
	)

	return eff, nil
}

package effect

import (
	"deckronomicon/packages/game/core"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/query"
	"fmt"
	"strconv"
)

// BuildEffectDraw creates a draw effect based on the provided modifiers.
// Keys: Count, Type
// Default: Count: 1
func BuildEffectDraw(source query.Object, spec definition.EffectSpec) (*Effect, error) {
	effect := Effect{id: spec.ID}
	count := "1"
	var drawType string
	for _, modifier := range spec.Modifiers {
		if modifier.Key == "Count" {
			count = modifier.Value
		}
		if modifier.Key == "Type" {
			drawType = modifier.Value
		}
	}
	_, err := strconv.Atoi(count)
	if err != nil {
		return nil, fmt.Errorf("invalid count: %s", count)
	}
	tags := []core.Tag{{Key: "Draw", Value: count}}
	if drawType != "" {
		tags = append(tags, core.Tag{Key: "Type", Value: drawType})
	}
	effect.description = fmt.Sprintf("draw %d cards", count)
	effect.tags = tags
	effect.Apply = func(state core.State, player core.Player) error {
		p, ok := player.(Player)
		if !ok {
			return fmt.Errorf("player does not implement Player interface: %T", player)
		}
		_, err := p.DrawCard()
		if err != nil {
			return fmt.Errorf("failed to draw card: %w", err)
		}
		return nil
	}
	return &effect, nil
}

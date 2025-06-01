package engine

/*
// BuildSpellAbility builds a spell ability from the given specification.
func BuildSpellAbility(spec *SpellAbilitySpec, source game.Object) (*SpellAbility, error) {
	ability := SpellAbility{}
	for _, effectSpec := range spec.EffectSpecs {
		effect, err := BuildEffect(source, effectSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to create effect: %w", err)
		}
		ability.Effects = append(ability.Effects, effect)
	}
	return &ability, nil
}


*/
/*
func BuildReplicateAbility(card *Card, replicateCount int) *TriggeredAbility {
	handler := func(state *GameState, player *Player) error {
		for range replicateCount {
			spell, err := NewSpell(card)
			if err != nil {
				return fmt.Errorf("failed to create spell from %s: %w", card.Name(), err)
			}
			state.Stack.Add(spell)
		}
		return nil
	}
	tags := []EffectTag{{Key: "Count", Value: strconv.Itoa(replicateCount)}}
	description := fmt.Sprintf("Replicate %s %d times", card.Name(), replicateCount)
	replicateAbility := TriggeredAbility{
		name: "Replicate Trigger",
		id:   GetNextID(),
		Effects: []Effect{
			{
				Apply:       handler,
				Tags:        tags,
				Description: description,
			},
		},
	}
	return &replicateAbility
}
*/

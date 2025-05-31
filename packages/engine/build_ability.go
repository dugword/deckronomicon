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

// BuildStaticAbility builds a static ability from the given specification.
func BuildStaticAbility(spec StaticAbilitySpec, source Object) (*StaticAbility, error) {
	ability := StaticAbility{
		ID: spec.ID,
	}
	for _, modifer := range spec.Modifiers {
		efectTag := EffectTag{
			Key:   modifer.Key,
			Value: modifer.Value,
		}
		ability.Modifiers = append(ability.Modifiers, efectTag)
	}
	return &ability, nil
}

// BuildActivatedAbility builds an activated ability from the given
// specification.
func BuildActivatedAbility(spec ActivatedAbilitySpec, source Object) (*ActivatedAbility, error) {
	// ZoneBattlefield is the default zone for activated abilities.
	zone := ZoneBattlefield
	if spec.Zone != "" {
		zone = spec.Zone
	}
	ability := ActivatedAbility{
		name:   spec.Name,
		id:     GetNextID(),
		source: source,
		Zone:   zone,
		Speed:  spec.Speed,
	}
	cost, err := NewCost(spec.Cost, source)
	if err != nil {
		return nil, fmt.Errorf("failed to create cost: %w", err)
	}
	ability.Cost = cost
	for _, effectSpec := range spec.EffectSpecs {
		effect, err := BuildEffect(source, effectSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to create effect: %w", err)
		}
		ability.Effects = append(ability.Effects, effect)
	}
	return &ability, nil
}

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

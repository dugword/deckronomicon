package card

func (c *Card) BuildActivatedAbility() error {
	// TODO: Maybe have engine call this by passing in a function so it can do
	// all the look ups when that data is available :shrug:
	/*
		for _, spec := range card.activatedAbilitySpecs {
			if spec.Zone == ZoneHand || spec.Zone == ZoneGraveyard {
				ability, err := BuildActivatedAbility(*spec, &card)
				if err != nil {
					return nil, fmt.Errorf("failed to build activated ability: %w", err)
				}
				card.activatedAbilities = append(card.activatedAbilities, ability)
			}
		}
	*/
	return nil
}

func (c *Card) BuildStaticAbility() error {
	/*
		for _, spec := range cardData.StaticAbilitySpecs {
			staticAbility, err := BuildStaticAbility(*spec, &card)
			if err != nil {
				return nil, fmt.Errorf("failed to build static ability: %w", err)
			}
			card.staticAbilities = append(card.staticAbilities, staticAbility)
		}
	*/
	return nil
}

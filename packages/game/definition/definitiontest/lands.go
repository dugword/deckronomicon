package definitiontest

import (
	"deckronomicon/packages/game/definition"
)

func PlainsDefinition(id string, playerID string) *definition.Permanent {
	plains := definition.Permanent{
		ID:         id,
		Name:       "Plains",
		CardTypes:  []string{"Land"},
		Supertypes: []string{"Basic"},
		Subtypes:   []string{"Plains"},
		Controller: playerID,
		Owner:      playerID,
		ActivatedAbilities: []*definition.Ability{{
			Cost: "{T}",
			Effects: []*definition.Effect{{
				Name: "AddMana",
				Modifiers: map[string]any{
					"Mana": "{W}",
				},
			}},
		}},
	}
	return &plains
}

func IslandDefinition(id string, playerID string) *definition.Permanent {
	island := definition.Permanent{
		ID:         id,
		Name:       "Island",
		CardTypes:  []string{"Land"},
		Supertypes: []string{"Basic"},
		Subtypes:   []string{"Island"},
		Controller: playerID,
		Owner:      playerID,
		ActivatedAbilities: []*definition.Ability{{
			Cost: "{T}",
			Effects: []*definition.Effect{{
				Name: "AddMana",
				Modifiers: map[string]any{
					"Mana": "{U}",
				},
			}},
		}},
	}
	return &island
}

func SwampDefinition(id string, playerID string) *definition.Permanent {
	swamp := definition.Permanent{
		ID:         id,
		Name:       "Swamp",
		CardTypes:  []string{"Land"},
		Supertypes: []string{"Basic"},
		Subtypes:   []string{"Swamp"},
		Controller: playerID,
		Owner:      playerID,
		ActivatedAbilities: []*definition.Ability{{
			Cost: "{T}",
			Effects: []*definition.Effect{{
				Name: "AddMana",
				Modifiers: map[string]any{
					"Mana": "{B}",
				},
			}},
		}},
	}
	return &swamp
}

func MountainDefinition(id string, playerID string) *definition.Permanent {
	mountain := definition.Permanent{
		ID:         id,
		Name:       "Mountain",
		CardTypes:  []string{"Land"},
		Supertypes: []string{"Basic"},
		Subtypes:   []string{"Mountain"},
		Controller: playerID,
		Owner:      playerID,
		ActivatedAbilities: []*definition.Ability{{
			Cost: "{T}",
			Effects: []*definition.Effect{{
				Name: "AddMana",
				Modifiers: map[string]any{
					"Mana": "{R}",
				},
			}},
		}},
	}
	return &mountain
}

func ForestDefinition(id string, playerID string) *definition.Permanent {
	forest := definition.Permanent{
		ID:         id,
		Name:       "Forest",
		CardTypes:  []string{"Land"},
		Supertypes: []string{"Basic"},
		Subtypes:   []string{"Forest"},
		Controller: playerID,
		Owner:      playerID,
		ActivatedAbilities: []*definition.Ability{{
			Cost: "{T}",
			Effects: []*definition.Effect{{
				Name: "AddMana",
				Modifiers: map[string]any{
					"Mana": "{G}",
				},
			}},
		}},
	}
	return &forest
}

func WastesDefinition(id string, playerID string) *definition.Permanent {
	wastes := definition.Permanent{
		ID:         id,
		Name:       "Wastes",
		CardTypes:  []string{"Land"},
		Supertypes: []string{"Basic"},
		Controller: playerID,
		Owner:      playerID,
		ActivatedAbilities: []*definition.Ability{{
			Cost: "{T}",
			Effects: []*definition.Effect{{
				Name: "AddMana",
				Modifiers: map[string]any{
					"Mana": "{C}",
				},
			}},
		}},
	}
	return &wastes
}

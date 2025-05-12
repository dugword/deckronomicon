package game

type CardTagMap map[string][]string

// e.g. {"Candy Trail": ["egg"], "Chromatic Star": ["egg", "draw_spell"], ...}

var CardPool = map[string]Card{
	"Island": {
		Object: Object{
			Name:               "Island",
			CardTypes:          []CardType{CardTypeLand},
			Subtypes:           []Subtype{SubtypeIsland},
			Supertypes:         []Supertype{SupertypeBasic},
			ActivatedAbilities: []ActivatedAbility{AbilityIsland},
			RulesText:          "{T}: Add {U}.",
		},
	},
	"Swamp": {
		Object: Object{
			Name:               "Swamp",
			CardTypes:          []CardType{CardTypeLand},
			Subtypes:           []Subtype{SubtypeSwamp},
			Supertypes:         []Supertype{SupertypeBasic},
			ActivatedAbilities: []ActivatedAbility{AbilitySwamp},
			RulesText:          "{T}: Add {B}.",
		},
	},
	"Mountain": {
		Object: Object{
			Name:               "Mountain",
			CardTypes:          []CardType{CardTypeLand},
			Subtypes:           []Subtype{SubtypeMountain},
			Supertypes:         []Supertype{SupertypeBasic},
			ActivatedAbilities: []ActivatedAbility{AbilityMountain},
			RulesText:          "{T}: Add {R}.",
		},
	},
	"Forest": {
		Object: Object{
			Name:               "Forest",
			CardTypes:          []CardType{CardTypeLand},
			Subtypes:           []Subtype{SubtypeForest},
			Supertypes:         []Supertype{SupertypeBasic},
			ActivatedAbilities: []ActivatedAbility{AbilityForest},
			RulesText:          "{T}: Add {G}.",
		},
	},
	"Plains": {
		Object: Object{
			Name:               "Plains",
			CardTypes:          []CardType{CardTypeLand},
			Subtypes:           []Subtype{SubtypePlains},
			Supertypes:         []Supertype{SupertypeBasic},
			ActivatedAbilities: []ActivatedAbility{AbilityPlains},
			RulesText:          "{T}: Add {W}.",
		},
	},
	"Pearled Unicorn": {
		Object: Object{
			Name:      "Pearled Unicorn",
			CardTypes: []CardType{CardTypeCreature},
			Subtypes:  []Subtype{SubtypeUnicorn},
			Power:     2,
			Toughness: 2,
			ManaCost:  ManaCost{Generic: 2, Colors: map[string]int{"W": 1}},
		},
	},
	"Preordain": {
		Object: Object{
			Name:         "Preordain",
			CardTypes:    []CardType{CardTypeSorcery},
			ManaCost:     ManaCost{Colors: map[string]int{"U": 1}},
			SpellAbility: &AbilityScry2, // Neds to draw
			RulesText:    "Scry 2",      // Should grab dynamically from abilities
		},
	},
	"Brainstorm": {
		Object: Object{
			Name:         "Brainstorm",
			CardTypes:    []CardType{CardTypeInstant},
			ManaCost:     ManaCost{Colors: map[string]int{"U": 1}},
			SpellAbility: &AbilityBrainstorm,
		},
	},
	"Candy Trail": {
		Object: Object{
			Name:      "Candy Trail",
			CardTypes: []CardType{CardTypeArtifact},
			ManaCost:  ManaCost{Generic: 1},
		},
	},
	"Foundry Inspector": {
		Object: Object{
			Name:      "Foundry Inspector",
			CardTypes: []CardType{CardTypeArtifact, CardTypeCreature},
			Subtypes:  []Subtype{SubtypeConstruct},
			ManaCost:  ManaCost{Generic: 3},
		},
	},
}

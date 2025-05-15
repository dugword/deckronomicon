package game

type CardTagMap map[string][]string

// e.g. {"Candy Trail": ["egg"], "Chromatic Star": ["egg", "draw_spell"], ...}

var CardPool = map[string]Card{
	"Island": {
		object: object{
			Name:               "Island",
			CardTypes:          []CardType{CardTypeLand},
			Subtypes:           []Subtype{SubtypeIsland},
			Supertypes:         []Supertype{SupertypeBasic},
			ActivatedAbilities: []ActivatedAbility{AbilityIsland},
			RulesText:          "{T}: Add {U}.",
		},
	},
	"Swamp": {
		object: object{
			Name:               "Swamp",
			CardTypes:          []CardType{CardTypeLand},
			Subtypes:           []Subtype{SubtypeSwamp},
			Supertypes:         []Supertype{SupertypeBasic},
			ActivatedAbilities: []ActivatedAbility{AbilitySwamp},
			RulesText:          "{T}: Add {B}.",
		},
	},
	"Mountain": {
		object: object{
			Name:               "Mountain",
			CardTypes:          []CardType{CardTypeLand},
			Subtypes:           []Subtype{SubtypeMountain},
			Supertypes:         []Supertype{SupertypeBasic},
			ActivatedAbilities: []ActivatedAbility{AbilityMountain},
			RulesText:          "{T}: Add {R}.",
		},
	},
	"Forest": {
		object: object{
			Name:               "Forest",
			CardTypes:          []CardType{CardTypeLand},
			Subtypes:           []Subtype{SubtypeForest},
			Supertypes:         []Supertype{SupertypeBasic},
			ActivatedAbilities: []ActivatedAbility{AbilityForest},
			RulesText:          "{T}: Add {G}.",
		},
	},
	"Plains": {
		object: object{
			Name:               "Plains",
			CardTypes:          []CardType{CardTypeLand},
			Subtypes:           []Subtype{SubtypePlains},
			Supertypes:         []Supertype{SupertypeBasic},
			ActivatedAbilities: []ActivatedAbility{AbilityPlains},
			RulesText:          "{T}: Add {W}.",
		},
	},
	"Pearled Unicorn": {
		object: object{
			Name:      "Pearled Unicorn",
			CardTypes: []CardType{CardTypeCreature},
			Subtypes:  []Subtype{SubtypeUnicorn},
			Power:     2,
			Toughness: 2,
			ManaCost:  ManaCost{Generic: 2, Colors: map[string]int{"W": 1}},
		},
	},
	"Preordain": {
		object: object{
			Name:         "Preordain",
			CardTypes:    []CardType{CardTypeSorcery},
			ManaCost:     ManaCost{Colors: map[string]int{"U": 1}},
			SpellAbility: &AbilityScry2, // Neds to draw
			RulesText:    "Scry 2",      // Should grab dynamically from abilities
		},
	},
	"Brainstorm": {
		object: object{
			Name:         "Brainstorm",
			CardTypes:    []CardType{CardTypeInstant},
			ManaCost:     ManaCost{Colors: map[string]int{"U": 1}},
			SpellAbility: &AbilityBrainstorm,
		},
	},
	"Candy Trail": {
		object: object{
			Name:      "Candy Trail",
			CardTypes: []CardType{CardTypeArtifact},
			ManaCost:  ManaCost{Generic: 1},
		},
	},
	"Foundry Inspector": {
		object: object{
			Name:      "Foundry Inspector",
			CardTypes: []CardType{CardTypeArtifact, CardTypeCreature},
			Subtypes:  []Subtype{SubtypeConstruct},
			ManaCost:  ManaCost{Generic: 3},
		},
	},
}

package game

type GameObject interface {
	ActivatedAbilities() []ActivatedAbility
	Card() *Card
	CardTypes() []CardType
	Color() string
	ColorIdicator() string
	Defense() int
	HasType(CardType) bool
	ID() string // unique instance ID
	Loyalty() int
	ManaCost() ManaCost
	Name() string
	Power() int
	RulesText() string
	SpellAbility() *SpellAbility
	StaticAbilities() []StaticAbility
	Subtypes() []Subtype
	Supertypes() []Supertype
	Toughness() int
	TriggeredAbilities() []TriggeredAbility
}

type object struct {
	ActivatedAbilities []ActivatedAbility
	CardTypes          []CardType
	Color              string
	ColorIdicator      string
	Defense            int
	ID                 string
	Loyalty            int
	ManaCost           ManaCost
	Name               string
	Power              int
	RulesText          string
	SpellAbility       *SpellAbility
	StaticAbilities    []StaticAbility
	Subtypes           []Subtype
	Supertypes         []Supertype
	Toughness          int
	TriggeredAbilities []TriggeredAbility
}

func (o *object) HasType(cardType CardType) bool {
	hasType := false
	for _, t := range o.CardTypes {
		if t == cardType {
			hasType = true
			break
		}
	}
	return hasType
}

func (o *object) HasSubtype(subtype Subtype) bool {
	hasSubtype := false
	for _, t := range o.Subtypes {
		if t == subtype {
			hasSubtype = true
			break
		}
	}
	return hasSubtype
}

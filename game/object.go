package game

// spell abilities, activated abilities, triggered abilities, and static abilitie
type Object struct {
	Name               string
	ManaCost           ManaCost
	Color              string
	ColorIdicator      string
	CardTypes          []CardType
	Subtypes           []Subtype
	Supertypes         []Supertype
	RulesText          string
	StaticAbilities    []StaticAbility
	ActivatedAbilities []ActivatedAbility
	TriggeredAbilities []TriggeredAbility
	SpellAbility       *SpellAbility
	Power              int
	Toughness          int
	Loyalty            int
	Defense            int
}

func (o *Object) HasType(cardType CardType) bool {
	hasType := false
	for _, t := range o.CardTypes {
		if t == cardType {
			hasType = true
			break
		}
	}
	return hasType
}

func (o *Object) HasSubtype(subtype Subtype) bool {
	hasSubtype := false
	for _, t := range o.Subtypes {
		if t == subtype {
			hasSubtype = true
			break
		}
	}
	return hasSubtype
}

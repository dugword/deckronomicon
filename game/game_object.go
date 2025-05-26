package game

// GameObject is an interface that represents any object in the game, a card,
// a spell, a permanent, or a token. It provides methods to access the
// object's properties and abilities.
type GameObject interface {
	HasSubtype(Subtype) bool
	HasCardType(CardType) bool
	HasColor(color Color) bool
	Name() string
	ActivatedAbilities() []*ActivatedAbility
	StaticAbilities() []*StaticAbility
	HasStaticAbility(string) bool
	ID() string
	ManaValue() int
	/*
		Card() *Card
		CardTypes() []CardType
		Colors() []string
		HasCardType(CardType) bool
		ID() string // unique instance ID
		Loyalty() int
		ManaCost() string
		Power() int
		RulesText() string
		StaticAbilities() []*StaticAbility
		Subtypes() []Subtype
		Supertypes() []Supertype
		Toughness() int
		TriggeredAbilities() []*TriggeredAbility
	*/
}

package game

// GameObject is an interface that represents any object in the game, a card,
// a spell, a permanent, or a token. It provides methods to access the
// object's properties and abilities.
type GameObject interface {
	ActivatedAbilities() []*ActivatedAbility
	Card() *Card
	CardTypes() []CardType
	Colors() *Colors
	Defense() int
	HasType(CardType) bool
	ID() string // unique instance ID
	Loyalty() int
	ManaCost() *ManaCost
	Name() string
	Power() int
	RulesText() string
	StaticAbilities() []*StaticAbility
	Subtypes() []Subtype
	Supertypes() []Supertype
	Toughness() int
	TriggeredAbilities() []*TriggeredAbility
}

// Colors represents the colors of a card or object in the game.
type Colors struct {
	Black bool
	Blue  bool
	Green bool
	Red   bool
	White bool
}

type EffectSpec struct {
	ID        string
	Modifiers []EffectModifier
}

type EffectModifier struct {
	Key   string
	Value string
}

// ActivatedAbilitySpec represents the specification of an activated ability.
type ActivatedAbilitySpec struct {
	CostExpression string       `json:"Cost,omitempty"`
	EffectSpecs    []EffectSpec `json:"Effects,omitempty"`
	Zone           string       `json:"Zone,omitempty"`
}

// SpellAbilitySpec represents the specification of a spell ability.
type SpellAbilitySpec struct {
	// CostExpression // TODO: AdditionalCosts?
	EffectSpecs []EffectSpec `json:"Effects,omitempty"`
	Zone        string       `json:"Zone,omitempty"`
}

// StaticAbilitySpec represents the specification of static ability.
type StaticAbilitySpec struct {
	EffectSpecs []EffectSpec `json:"Effects,omitempty"`
	Zone        string       `json:"Zone,omitempty"`
}

// object is the base struct for all game objects. It contains common fields
// and methods that are shared among different types of game objects.
// Is this always the underlying object and never accessed directly?
// E.g. game objects are always wrapped in a Card or Permanent or a Spell or
// something else?
type object struct {
	// If so these shouldn't exist here and only contain the spec
	// when the object is transformed into a different thing then
	// the spec should be used to create the actual abilities list
	ActivatedAbilities []ActivatedAbility
	// TODO: Make all the references to this name plural when an array
	ActivatedAbilitiesSpec []ActivatedAbilitySpec
	CardTypes              []CardType
	Colors                 *Colors
	Defense                int
	ID                     string
	Loyalty                int
	ManaCost               *ManaCost
	Name                   string
	Power                  int
	RulesText              string
	SpellAbility           *SpellAbility
	SpellAbilitySpec       *SpellAbilitySpec
	StaticAbilities        []StaticAbility
	Subtypes               []Subtype
	Supertypes             []Supertype
	Toughness              int
	TriggeredAbilities     []TriggeredAbility
}

// HasType checks if the object has the specified card type.
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

// HasSubtype checks if the object has the specified subtype.
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

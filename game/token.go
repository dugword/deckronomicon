package game

// Token represents a token object in the game.
type Token struct {
	object
}

// NewToken creates a new Token instance from an object.
func NewToken(object object) *Token {
	token := Token{
		object: object,
	}
	return &token
}

// ActivatedAbilities returns the activated abilities of the token.
func (t *Token) ActivatedAbilities() []ActivatedAbility {
	return t.object.ActivatedAbilities
}

// Card returns the card associated with the token.
// Tokens do not have a card associated with them, so this returns nil.
func (t *Token) Card() *Card {
	return nil
}

// CardTypes returns the card types of the token.
func (t *Token) CardTypes() []CardType {
	return t.object.CardTypes
}

// Colors returns the colors of the token.
func (t *Token) Colors() *Colors {
	return t.object.Colors
}

// Defense returns the defense of the token.
func (t *Token) Defense() int {
	return t.object.Defense
}

// HasType checks if the token has the specified card type.
func (t *Token) HasType(cardType CardType) bool {
	return t.object.HasType(cardType)
}

// ID returns the ID of the token.
func (t *Token) ID() string {
	return t.object.ID
}

// Loyalty returns the loyalty of the token.
func (t *Token) Loyalty() int {
	return t.object.Loyalty
}

// ManaCost returns the mana cost of the token.
func (t *Token) ManaCost() *ManaCost {
	return t.object.ManaCost
}

// Name returns the name of the token.
func (t *Token) Name() string {
	return t.object.Name
}

// Power returns the power of the token.
func (t *Token) Power() int {
	return t.object.Power
}

// RulesText returns the rules text of the token.
func (t *Token) RulesText() string {
	return t.object.RulesText
}

// StaticAbilities returns the static abilities of the token.
func (t *Token) StaticAbilities() []StaticAbility {
	return t.object.StaticAbilities
}

// Subtypes returns the subtypes of the token.
func (t *Token) Subtypes() []Subtype {
	return t.object.Subtypes
}

// Supertypes returns the supertypes of the token.
func (t *Token) Supertypes() []Supertype {
	return t.object.Supertypes
}

// Toughness returns the toughness of the token.
func (t *Token) Toughness() int {
	return t.object.Toughness
}

// TriggeredAbilities returns the triggered abilities of the token.
func (t *Token) TriggeredAbilities() []TriggeredAbility {
	return t.object.TriggeredAbilities
}

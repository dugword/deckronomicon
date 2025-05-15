package game

// spell abilities, activated abilities, triggered abilities, and static abilitie
type Token struct {
	object
}

func (t *Token) ActivatedAbilities() []ActivatedAbility {
	return t.object.ActivatedAbilities
}
func (t *Token) Card() *Card {
	return nil
}
func (t *Token) CardTypes() []CardType {
	return t.object.CardTypes
}
func (t *Token) Color() string {
	return t.object.Color
}
func (t *Token) ColorIdicator() string {
	return t.object.ColorIdicator
}
func (t *Token) Defense() int {
	return t.object.Defense
}
func (t *Token) HasType(cardType CardType) bool {
	return t.object.HasType(cardType)
}
func (t *Token) ID() string {
	return t.object.ID
}
func (t *Token) Loyalty() int {
	return t.object.Loyalty
}
func (t *Token) ManaCost() ManaCost {
	return t.object.ManaCost
}
func (t *Token) Name() string {
	return t.object.Name
}
func (t *Token) Power() int {
	return t.object.Power
}
func (t *Token) RulesText() string {
	// TODO: Need to build this from abilities....
	return t.object.RulesText
}
func (t *Token) SpellAbility() *SpellAbility {
	return t.object.SpellAbility
}
func (t *Token) StaticAbilities() []StaticAbility {
	return t.object.StaticAbilities
}
func (t *Token) Subtypes() []Subtype {
	return t.object.Subtypes
}
func (t *Token) Supertypes() []Supertype {
	return t.object.Supertypes
}
func (t *Token) Toughness() int {
	return t.object.Toughness
}
func (t *Token) TriggeredAbilities() []TriggeredAbility {
	return t.object.TriggeredAbilities
}

package game

// spell abilities, activated abilities, triggered abilities, and static abilitie
type Card struct {
	object
}

func (c *Card) IsPermanent() bool {
	if c.HasType(CardTypeLand) ||
		c.HasType(CardTypeArtifact) ||
		c.HasType(CardTypeBattle) ||
		c.HasType(CardTypeCreature) ||
		c.HasType(CardTypePlaneswalker) {
		return true
	}
	return false
}

func (c *Card) ActivatedAbilities() []ActivatedAbility {
	return c.object.ActivatedAbilities
}
func (c *Card) Card() *Card {
	return c
}
func (c *Card) CardTypes() []CardType {
	return c.object.CardTypes
}
func (c *Card) Color() string {
	return c.object.Color
}
func (c *Card) ColorIdicator() string {
	return c.object.ColorIdicator
}
func (c *Card) Defense() int {
	return c.object.Defense
}
func (c *Card) HasType(cardType CardType) bool {
	return c.object.HasType(cardType)
}
func (c *Card) ID() string {
	return c.object.ID
}
func (c *Card) Loyalty() int {
	return c.object.Loyalty
}
func (c *Card) ManaCost() ManaCost {
	return c.object.ManaCost
}
func (c *Card) Name() string {
	return c.object.Name
}
func (c *Card) Power() int {
	return c.object.Power
}
func (c *Card) RulesText() string {
	// TODO: Need to build this from abilities....
	return c.object.RulesText
}
func (c *Card) SpellAbility() *SpellAbility {
	return c.object.SpellAbility
}
func (c *Card) StaticAbilities() []StaticAbility {
	return c.object.StaticAbilities
}
func (c *Card) Subtypes() []Subtype {
	return c.object.Subtypes
}
func (c *Card) Supertypes() []Supertype {
	return c.object.Supertypes
}
func (c *Card) Toughness() int {
	return c.object.Toughness
}
func (c *Card) TriggeredAbilities() []TriggeredAbility {
	return c.object.TriggeredAbilities
}

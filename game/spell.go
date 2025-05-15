package game

type Spell struct {
	card *Card
}

func (s *Spell) ActivatedAbilities() []ActivatedAbility {
	return s.card.object.ActivatedAbilities
}
func (s *Spell) Card() *Card {
	return s.card
}
func (s *Spell) CardTypes() []CardType {
	return s.card.object.CardTypes
}
func (s *Spell) Color() string {
	return s.card.object.Color
}
func (s *Spell) ColorIdicator() string {
	return s.card.object.ColorIdicator
}
func (s *Spell) Defense() int {
	return s.card.object.Defense
}
func (s *Spell) HasType(cardType CardType) bool {
	return s.card.object.HasType(cardType)
}
func (s *Spell) ID() string {
	return s.card.object.ID
}
func (s *Spell) Loyalty() int {
	return s.card.object.Loyalty
}
func (s *Spell) ManaCost() ManaCost {
	return s.card.object.ManaCost
}
func (s *Spell) Name() string {
	return s.card.object.Name
}
func (s *Spell) Power() int {
	return s.card.object.Power
}
func (s *Spell) RulesText() string {
	// TODO: Need to build this from abilities....
	return s.card.object.RulesText
}
func (s *Spell) SpellAbility() *SpellAbility {
	return s.card.object.SpellAbility
}
func (s *Spell) StaticAbilities() []StaticAbility {
	return s.card.object.StaticAbilities
}
func (s *Spell) Subtypes() []Subtype {
	return s.card.object.Subtypes
}
func (s *Spell) Supertypes() []Supertype {
	return s.card.object.Supertypes
}
func (s *Spell) Toughness() int {
	return s.card.object.Toughness
}
func (s *Spell) TriggeredAbilities() []TriggeredAbility {
	return s.card.object.TriggeredAbilities
}

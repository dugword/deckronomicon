package game

type ManaPool map[string]int

type Card struct {
	Object
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

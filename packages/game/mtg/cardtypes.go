package mtg

import "fmt"

// CardType is a type of card in Magic: The Gathering.
type CardType string

// CardType constants represent the different types of cards.
const (
	CardTypeArtifact     CardType = "Artifact"
	CardTypeBattle       CardType = "Battle"
	CardTypeCreature     CardType = "Creature"
	CardTypeEnchantment  CardType = "Enchantment"
	CardTypeInstant      CardType = "Instant"
	CardTypeLand         CardType = "Land"
	CardTypePlaneswalker CardType = "Planeswalker"
	CardTypeSorcery      CardType = "Sorcery"
)

// PermanentCardTypes is a set of permanent card types.
var PermanentCardTypes = map[CardType]struct{}{
	CardTypeArtifact:     {},
	CardTypeBattle:       {},
	CardTypeCreature:     {},
	CardTypeEnchantment:  {},
	CardTypePlaneswalker: {},
}

// SpellCardTypes is a set of spell card types.
// NOTE: Technical all non-land cards are spells, but common usage
// is to refer to non-permanent cards as spells.
// TODO: This may be removed later as we implement an actual stack system
// and convert all cards to Spells as they are cast.
var SpellCardTypes = map[CardType]struct{}{
	CardTypeInstant: {},
	CardTypeSorcery: {},
}

// IsLand checks if the CardType is a Land type.
func (t CardType) IsLand() bool {
	if t == CardTypeLand {
		return true
	}
	return false
}

// IsPermanent checks if the CardType is a permanent type.
func (t CardType) IsPermanent() bool {
	if _, ok := PermanentCardTypes[t]; ok {
		return true
	}
	return false
}

// IsSpell checks if the CardType is a Spell type.
func (t CardType) IsSpell() bool {
	if _, ok := SpellCardTypes[t]; ok {
		return true
	}
	return false
}

// StringToCardType converts a string to a CardType.
func StringToCardType(s string) (CardType, error) {
	stringToCardType := map[string]CardType{
		"Artifact":     CardTypeArtifact,
		"Battle":       CardTypeBattle,
		"Creature":     CardTypeCreature,
		"Enchantment":  CardTypeEnchantment,
		"Instant":      CardTypeInstant,
		"Land":         CardTypeLand,
		"Planeswalker": CardTypePlaneswalker,
		"Sorcery":      CardTypeSorcery,
	}
	cardType, ok := stringToCardType[s]
	if !ok {
		return "", fmt.Errorf("unknown CardType '%s'", s)
	}
	return cardType, nil
}

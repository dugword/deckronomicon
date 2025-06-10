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

// IsLand checks if the CardType is a Land type.
func (t CardType) IsLand() bool {
	return t == CardTypeLand
}

// IsPermanent checks if the CardType is a permanent type.
func (t CardType) IsPermanent() bool {
	// PermanentCardTypes is a set of permanent card types.
	permanentCardTypes := map[CardType]struct{}{
		CardTypeArtifact:     {},
		CardTypeBattle:       {},
		CardTypeCreature:     {},
		CardTypeEnchantment:  {},
		CardTypePlaneswalker: {},
		CardTypeLand:         {},
	}
	if _, ok := permanentCardTypes[t]; ok {
		return true
	}
	return false
}

// IsSpell checks if the CardType is a Spell type.
func (t CardType) IsSpell() bool {
	// SpellCardTypes is a set of spell card types.
	// NOTE: Technical all non-land cards are spells, but common usage
	// is to refer to non-permanent cards as spells.
	// TODO: This may be removed later as we implement an actual stack system
	// and convert all cards to Spells as they are cast.
	spellCardTypes := map[CardType]struct{}{
		CardTypeInstant: {},
		CardTypeSorcery: {},
	}
	if _, ok := spellCardTypes[t]; ok {
		return true
	}
	return false
}

// TODO: Have all the StringTo* Functions return a bool indicating success
// and remove the error return type this will  be consistent with the rest
// of the codebase
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
		return "", fmt.Errorf("unknown CardType %q", s)
	}
	return cardType, nil
}

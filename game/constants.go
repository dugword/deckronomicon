package game

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
		return "", fmt.Errorf("unknown card type: %s", s)
	}
	return cardType, nil
}

// Subtype is a subtype of card in Magic: The Gathering.
type Subtype string

// Artifact Subtypes
const (
	SubtypeBlood         Subtype = "Blood"
	SubtypeClue          Subtype = "Clue"
	SubtypeFood          Subtype = "Food"
	SubtypeGold          Subtype = "Gold"
	SubtypeIncubator     Subtype = "Incubator"
	SubtypeJunk          Subtype = "Junk"
	SubtypeMap           Subtype = "Map"
	SubtypePowerstone    Subtype = "Powerstone"
	SubtypeTreasure      Subtype = "Treasure"
	SubtypeEquipment     Subtype = "Equipment"
	SubtypeFortification Subtype = "Fortification"
	SubtypeVehicle       Subtype = "Vehicle"
)

// StringToArtifactSubtype converts a string to a Subtype.
var StringToArtifactSubtype = map[string]Subtype{
	"Blood":         SubtypeBlood,
	"Clue":          SubtypeClue,
	"Food":          SubtypeFood,
	"Gold":          SubtypeGold,
	"Incubator":     SubtypeIncubator,
	"Junk":          SubtypeJunk,
	"Map":           SubtypeMap,
	"Powerstone":    SubtypePowerstone,
	"Treasure":      SubtypeTreasure,
	"Equipment":     SubtypeEquipment,
	"Fortification": SubtypeFortification,
	"Vehicle":       SubtypeVehicle,
}

func (t Subtype) IsArtifactSubtype() bool {
	if _, ok := StringToArtifactSubtype[string(t)]; ok {
		return true
	}
	return false
}

// Battle Subtypes
const (
	SubtypeSiege Subtype = "Siege"
)

// StringToBattleSubtype converts a string to a Subtype.
var StringToBattleSubtype = map[string]Subtype{
	"Siege": SubtypeSiege,
}

func (t Subtype) IsBattleSubtype() bool {
	if _, ok := StringToBattleSubtype[string(t)]; ok {
		return true
	}
	return false
}

// Enchantment Subtypes
const (
	SubtypeAura       Subtype = "Aura"
	SubtypeBackground Subtype = "Background"
	SubtypeCartouche  Subtype = "Cartouche"
	SubtypeCase       Subtype = "Case"
	SubtypeClass      Subtype = "Class"
	SubtypeCurse      Subtype = "Curse"
	SubtypeRole       Subtype = "Role"
	SubtypeRune       Subtype = "Rune"
	SubtypeSaga       Subtype = "Saga"
	SubtypeShard      Subtype = "Shard"
	SubtypeShrine     Subtype = "Shrine"
)

// StringToEnchantmentSubtype converts a string to a Subtype.
var StringToEnchantmentSubtype = map[string]Subtype{
	"Aura":       SubtypeAura,
	"Background": SubtypeBackground,
	"Cartouche":  SubtypeCartouche,
	"Case":       SubtypeCase,
	"Class":      SubtypeClass,
	"Curse":      SubtypeCurse,
	"Role":       SubtypeRole,
	"Rune":       SubtypeRune,
	"Saga":       SubtypeSaga,
	"Shard":      SubtypeShard,
	"Shrine":     SubtypeShrine,
}

func (t Subtype) IsEnchantmentSubtype() bool {
	if _, ok := StringToEnchantmentSubtype[string(t)]; ok {
		return true
	}
	return false
}

// Basic Land Subtypes
const (
	SubtypeForest   Subtype = "Forest"
	SubtypeIsland   Subtype = "Island"
	SubtypeMountain Subtype = "Mountain"
	SubtypePlains   Subtype = "Plains"
	SubtypeSwamp    Subtype = "Swamp"
)

// StringToBasicLandSubtype converts a string to a Subtype.
var StringToBasicLandSubtype = map[string]Subtype{
	"Forest":   SubtypeForest,
	"Island":   SubtypeIsland,
	"Mountain": SubtypeMountain,
	"Plains":   SubtypePlains,
	"Swamp":    SubtypeSwamp,
}

// IsBasicLandSubtype checks if the Subtype is a basic land subtype.
func (t Subtype) IsBasicLandSubtype() bool {
	if _, ok := StringToBasicLandSubtype[string(t)]; ok {
		return true
	}
	return false
}

// Nonbasic Land Subtypes
const (
	SubtypeCave       Subtype = "Cave"
	SubtypeDesert     Subtype = "Desert"
	SubtypeGate       Subtype = "Gate"
	SubtypeLair       Subtype = "Lair"
	SubtypeLocus      Subtype = "Locus"
	SubtypeMine       Subtype = "Mine"
	SubtypePowerPlant Subtype = "Power-Plant"
	SubtypeSphere     Subtype = "Sphere"
	SubtypeTower      Subtype = "Tower"
	SubtypeUrzas      Subtype = "Urza's"
)

// StringToNonbasicLandSubtype converts a string to a Subtype.
var StringToNonbasicLandSubtype = map[string]Subtype{
	"Cave":        SubtypeCave,
	"Desert":      SubtypeDesert,
	"Gate":        SubtypeGate,
	"Lair":        SubtypeLair,
	"Locus":       SubtypeLocus,
	"Mine":        SubtypeMine,
	"Power-Plant": SubtypePowerPlant,
	"Sphere":      SubtypeSphere,
	"Tower":       SubtypeTower,
	"Urza's":      SubtypeUrzas,
}

// IsNonbasicLandSubtype checks if the Subtype is a nonbasic land subtype.
func (t Subtype) IsNonbasicLandSubtype() bool {
	if _, ok := StringToNonbasicLandSubtype[string(t)]; ok {
		return true
	}
	return false
}

// Instant/Sorcery Subtypes
const (
	SubtypeAdventure Subtype = "Adventure"
	SubtypeArcane    Subtype = "Arcane"
	SubtypeChorus    Subtype = "Chorus"
	SubtypeLesson    Subtype = "Lesson"
	SubtypeOmen      Subtype = "Omen"
	SubtypeTrap      Subtype = "Trap" // Instant only
)

// StringToInstantSorcerySubtype converts a string to a Subtype.
var StringToInstantSorcerySubtype = map[string]Subtype{
	"Adventure": SubtypeAdventure,
	"Arcane":    SubtypeArcane,
	"Chorus":    SubtypeChorus,
	"Lesson":    SubtypeLesson,
	"Omen":      SubtypeOmen,
	"Trap":      SubtypeTrap,
}

// IsInstantSorcerySubtype checks if the Subtype is an Instant or Sorcery
// subtype.
func (t Subtype) IsInstantSorcerySubtype() bool {
	if _, ok := StringToInstantSorcerySubtype[string(t)]; ok {
		return true
	}
	return false
}

func StringToSubtype(s string) (Subtype, error) {
	if t, ok := StringToArtifactSubtype[s]; ok {
		return t, nil
	}
	if t, ok := StringToBattleSubtype[s]; ok {
		return t, nil
	}
	if t, ok := StringToEnchantmentSubtype[s]; ok {
		return t, nil
	}
	if t, ok := StringToBasicLandSubtype[s]; ok {
		return t, nil
	}
	if t, ok := StringToNonbasicLandSubtype[s]; ok {
		return t, nil
	}
	if t, ok := StringToInstantSorcerySubtype[s]; ok {
		return t, nil
	}
	return "", fmt.Errorf("unknown Subtype: %s")
}

type Supertype string

const (
	SupertypeBasic     Supertype = "Basic"
	SupertypeLegendary Supertype = "Legendary"
	SupertypeSnow      Supertype = "Snow"
	SupertypeWorld     Supertype = "World"
)

// StringToSupertype converts a string to a Supertype.
func StringToSupertype(s string) (Supertype, error) {
	stringToSupertype := map[string]Supertype{
		"Basic":     SupertypeBasic,
		"Legendary": SupertypeLegendary,
		"Snow":      SupertypeSnow,
		"World":     SupertypeWorld,
	}
	supertype, ok := stringToSupertype[s]
	if !ok {
		return "", fmt.Errorf("unknown supertype: %s", s)
	}
	return supertype, nil
}

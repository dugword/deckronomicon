package mtg

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

func StringToSubtype(s string) (Subtype, bool) {
	if t, ok := StringToArtifactSubtype[s]; ok {
		return t, true
	}
	if t, ok := StringToBattleSubtype[s]; ok {
		return t, true
	}
	if t, ok := StringToEnchantmentSubtype[s]; ok {
		return t, true
	}
	if t, ok := StringToBasicLandSubtype[s]; ok {
		return t, true
	}
	if t, ok := StringToNonbasicLandSubtype[s]; ok {
		return t, true
	}
	if t, ok := StringToInstantSorcerySubtype[s]; ok {
		return t, true
	}
	if t, ok := StringToCreatureSubtype[s]; ok {
		return t, true
	}
	return "", false
}

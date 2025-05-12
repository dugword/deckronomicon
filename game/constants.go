package game

type CardType string

const (
	CardTypeArtifact     CardType = "Artifact"
	CardTypeCreature     CardType = "Creature"
	CardTypeEnchantment  CardType = "Enchantment"
	CardTypePlaneswalker CardType = "Planeswalker"
	CardTypeBattle       CardType = "Battle"
	CardTypeInstant      CardType = "Instant"
	CardTypeSorcery      CardType = "Sorcery"
	CardTypeLand         CardType = "Land"
)

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

// Battle Subtypes
const (
	SubtypeSiege Subtype = "Siege"
)

// Enchantment Subtypes
const (
	SubtypeAura       Subtype = "Aura"
	SubtypeBackground Subtype = "Background"
	SubtypeSaga       Subtype = "Saga"
	SubtypeRole       Subtype = "Role"
	SubtypeShard      Subtype = "Shard"
	SubtypeCartouche  Subtype = "Cartouche"
	SubtypeCase       Subtype = "Case"
	SubtypeClass      Subtype = "Class"
	SubtypeCurse      Subtype = "Curse"
	SubtypeRune       Subtype = "Rune"
	SubtypeShrine     Subtype = "Shrine"
)

// Basic Land Subtypes
const (
	SubtypePlains   Subtype = "Plains"
	SubtypeIsland   Subtype = "Island"
	SubtypeSwamp    Subtype = "Swamp"
	SubtypeMountain Subtype = "Mountain"
	SubtypeForest   Subtype = "Forest"
)

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

// Instant/Sorcery Subtypes
const (
	SubtypeAdventure Subtype = "Adventure"
	SubtypeArcane    Subtype = "Arcane"
	SubtypeChorus    Subtype = "Chorus"
	SubtypeLesson    Subtype = "Lesson"
	SubtypeOmen      Subtype = "Omen"
	SubtypeTrap      Subtype = "Trap" // Instant only
)

type Supertype string

const (
	SupertypeBasic     Supertype = "Basic"
	SupertypeLegendary Supertype = "Legendary"
	SupertypeSnow      Supertype = "Snow"
	SupertypeWorld     Supertype = "World"
)

type CardColor string

const (
	ColorWhite CardColor = "W"
	ColorBlue  CardColor = "U"
	ColorBlack CardColor = "B"
	ColorRed   CardColor = "R"
	ColorGreen CardColor = "G"
	Colorless  CardColor = "C"
)

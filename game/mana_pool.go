package game

import (
	"fmt"
	"maps"
	"regexp"
	"strings"
)

// ManaPool represents a pool of mana available to the player.
type ManaPool struct {
	pool map[Color]int
}

// ManaStringToManaSymbols converts a mana string to a slice of mana symbols.
// The string should be in the format "{W}{U}{B}{R}{G}{C}".
func ManaStringToManaSymbols(mana string) []string {
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(mana, -1)
	var manaSymbols []string
	for _, match := range matches {
		symbol := fmt.Sprintf("{%s}", strings.ToUpper(match[1]))
		manaSymbols = append(manaSymbols, symbol)
	}
	return manaSymbols
}

// NewManaPool creates a new empty mana pool.
func NewManaPool() *ManaPool {
	manaPool := ManaPool{
		pool: map[Color]int{},
	}
	return &manaPool
}

// Add adds mana to the mana pool.
// TODO: Do I need to validate color? Or should I just assume it's valid?
func (m *ManaPool) Add(color Color, amount int) {
	m.pool[color] += amount
}

// AddMana adds mana to the mana pool from a string representation of a
// mana pool. The string should be in the format "{W}{U}{B}{R}{G}{C}".
// TODO: Centralize this with the mana cost parser.
func (m *ManaPool) AddMana(mana string) error {
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(mana, -1)
	for _, match := range matches {
		symbol := strings.ToUpper(match[1])
		color, err := StringToColor(symbol)
		if err != nil {
			return err
		}
		m.pool[color]++
	}
	return nil
}

func (m *ManaPool) Available(color Color) int {
	return m.pool[color]
}

func (m *ManaPool) AvailableGeneric() int {
	total := 0
	for _, amount := range m.pool {
		total += amount
	}
	return total
}

// ColorsAvailable returns a slice of colors that are available in the mana
// pool.
func (m *ManaPool) ColorsAvailable() []Color {
	var colors []Color
	for color, amount := range m.pool {
		if amount > 0 {
			colors = append(colors, color)
		}
	}
	return colors
}

// Copy creates a copy of the mana pool.
func (m *ManaPool) Copy() *ManaPool {
	newPool := ManaPool{
		pool: map[Color]int{},
	}
	maps.Copy(newPool.pool, m.pool)
	return &newPool
}

// Has checks if the mana pool has the specified amount of mana of the
// specified color.
func (m *ManaPool) Has(color Color, amount int) bool {
	if m.pool[color] < amount {
		return false
	}
	return true
}

// Get returns the amount of mana generic mana available.
func (m *ManaPool) HasGeneric(amount int) bool {
	total := 0
	for _, amount := range m.pool {
		total += amount
	}
	if total < amount {
		return false
	}
	return true
}

// Describe returns a string representation of the mana pool.
func (m *ManaPool) Describe() string {
	var mana []string
	// TODO: Consolidate with Description in ManaCost
	for _, color := range []Color{
		// Do not reorder this list, follows the color wheel standard for mana
		// costs.
		ColorColorless,
		ColorWhite,
		ColorBlue,
		ColorBlack,
		ColorRed,
		ColorGreen,
	} {
		if m.pool[color] > 0 {
			for range m.pool[color] {
				mana = append(mana, fmt.Sprintf("{%s}", color))
			}
		}
	}
	return strings.Join(mana, "")
}

// Empty empties the mana pool.
func (m *ManaPool) Empty() {
	m.pool = map[Color]int{}
}

// Use removes the specified amount of mana from the pool.
func (m *ManaPool) Use(color Color, amount int) error {
	if m.pool[color] < amount {
		return fmt.Errorf("insufficient %s mana in pool", color)
	}
	m.pool[color] -= amount
	return nil
}

package player

import (
	"deckronomicon/packages/game/mtg"
	"fmt"
	"maps"
	"regexp"
	"strings"
)

// pool represents a pool of mana available to the player.
type manaPool struct {
	pool map[mtg.Color]int
}

// NewPool creates a new empty mana pool.
func newManaPool() *manaPool {
	pool := manaPool{
		pool: map[mtg.Color]int{},
	}
	return &pool
}

// Add adds mana to the mana pool.
// TODO: Do I need to validate color? Or should I just assume it's valid?
func (p *manaPool) Add(color mtg.Color, amount int) {
	p.pool[color] += amount
}

// AddMana adds mana to the mana pool from a string representation of a
// mana pool. The string should be in the format "{W}{U}{B}{R}{G}{C}".
// TODO: Centralize this with the mana cost parser.
func (p *manaPool) AddMana(mana string) error {
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(mana, -1)
	for _, match := range matches {
		symbol := strings.ToUpper(match[1])
		color, err := mtg.StringToColor(symbol)
		if err != nil {
			return err
		}
		p.pool[color]++
	}
	return nil
}

func (p *manaPool) Available(color mtg.Color) int {
	return p.pool[color]
}

func (p *manaPool) AvailableGeneric() int {
	total := 0
	for _, amount := range p.pool {
		total += amount
	}
	return total
}

// ColorsAvailable returns a slice of colors that are available in the mana
// pool.
func (p *manaPool) ColorsAvailable() []mtg.Color {
	var colors []mtg.Color
	for color, amount := range p.pool {
		if amount > 0 {
			colors = append(colors, color)
		}
	}
	return colors
}

// Copy creates a copy of the mana pool.
func (p *manaPool) Copy() *manaPool {
	copiedPool := manaPool{
		pool: map[mtg.Color]int{},
	}
	maps.Copy(copiedPool.pool, p.pool)
	return &copiedPool
}

// Has checks if the mana pool has the specified amount of mana of the
// specified color.
func (p *manaPool) Has(color mtg.Color, amount int) bool {
	if p.pool[color] < amount {
		return false
	}
	return true
}

// Get returns the amount of mana generic mana available.
func (p *manaPool) HasGeneric(amount int) bool {
	total := 0
	for _, amount := range p.pool {
		total += amount
	}
	if total < amount {
		return false
	}
	return true
}

// Describe returns a string representation of the mana pool.
func (p *manaPool) Describe() string {
	var mana []string
	// TODO: Consolidate with Description in ManaCost
	for _, color := range []mtg.Color{
		// Do not reorder this list, follows the color wheel standard for mana
		// costs.
		mtg.ColorColorless,
		mtg.ColorWhite,
		mtg.ColorBlue,
		mtg.ColorBlack,
		mtg.ColorRed,
		mtg.ColorGreen,
	} {
		if p.pool[color] > 0 {
			for range p.pool[color] {
				mana = append(mana, fmt.Sprintf("{%s}", color))
			}
		}
	}
	return strings.Join(mana, "")
}

// Empty empties the mana pool.
func (p *manaPool) Empty() {
	p.pool = map[mtg.Color]int{}
}

// Use removes the specified amount of mana from the pool.
func (p *manaPool) Use(color mtg.Color, amount int) error {
	if p.pool[color] < amount {
		return fmt.Errorf("insufficient %s mana in pool", color)
	}
	p.pool[color] -= amount
	return nil
}

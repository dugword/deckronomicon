package mana

import (
	"fmt"
	"maps"
	"regexp"
	"strconv"
	"strings"
)

type ManaType string

// This might be beter as a const in mtg package
const (
	White     ManaType = "W"
	Blue      ManaType = "U"
	Black     ManaType = "B"
	Red       ManaType = "R"
	Green     ManaType = "G"
	Colorless ManaType = "C"
)

var AllManaTypes = []ManaType{
	White,
	Blue,
	Black,
	Red,
	Green,
	Colorless,
}

func (m ManaType) String() string {
	return fmt.Sprintf("{%s}", string(m))
}

type Pool struct {
	mana map[ManaType]int
}

func (p Pool) Describe() string {
	descriptions := []string{}
	for _, manaType := range AllManaTypes {
		for range p.mana[manaType] {
			descriptions = append(descriptions, manaType.String())
		}
	}
	if p.Total() == 0 {
		return "(empty)"
	}
	return strings.Join(descriptions, "")
}

func NewManaPool() Pool {
	m := Pool{
		mana: map[ManaType]int{},
	}
	for _, manaType := range AllManaTypes {
		m.mana[manaType] = 0
	}
	return m
}

func (p Pool) Total() int {
	total := 0
	for _, amount := range p.mana {
		total += amount
	}
	return total
}

func (p Pool) WithAddedAmount(amount Amount) Pool {
	newPool := NewManaPool()
	maps.Copy(newPool.mana, p.mana)
	newPool.mana[Colorless] += amount.generic
	for mt, amt := range amount.colors {
		newPool.mana[mt] += amt
	}
	return newPool
}

func (p Pool) WithAddedMana(manaType ManaType, amount int) Pool {
	newPool := NewManaPool()
	maps.Copy(newPool.mana, p.mana)
	newPool.mana[manaType] += amount
	return newPool
}

// This auto pays for mana, spending colored mana first, then generic mana.
// This needs to be made to be smarter or allow for users to specify how they want to pay of each color.
// But this is a good start for now.
func (p Pool) WithSpentFromManaAmount(amount Amount) (Pool, error) {
	// First check if we can pay the colored part
	for mt, required := range amount.colors {
		if p.mana[mt] < required {
			return p, fmt.Errorf("not enough %s mana", mt)
		}
	}
	availableForGeneric := p.Total() - colorSpecificUsed(amount)
	if availableForGeneric < amount.generic {
		return p, fmt.Errorf("not enough total mana to pay generic cost (%d required)", amount.generic)
	}
	newPool := NewManaPool()
	maps.Copy(newPool.mana, p.mana)
	// Spend colored
	for mt, required := range amount.colors {
		newPool.mana[mt] -= required
	}
	// Spend generic
	genericLeft := amount.generic
	for genericLeft > 0 {
		for _, mt := range AllManaTypes {
			if newPool.mana[mt] > 0 {
				newPool.mana[mt]--
				genericLeft--
				if genericLeft == 0 {
					break
				}
			}
		}
	}
	return newPool, nil
}

// How much total mana is used that requires a color
// This is a little bit of a hack, but it works for now
func colorSpecificUsed(amount Amount) int {
	used := 0
	for _, required := range amount.colors {
		used += required
	}
	return used
}

type Amount struct {
	generic int
	colors  map[ManaType]int
}

func (a Amount) Total() int {
	total := a.generic
	for _, amount := range a.colors {
		total += amount
	}
	return total
}

func (a Amount) Generic() int {
	return a.generic
}

func (a Amount) Colors() map[ManaType]int {
	newColors := map[ManaType]int{}
	maps.Copy(newColors, a.colors)
	return newColors
}

func NewAmount() Amount {
	amount := Amount{
		generic: 0,
		colors:  map[ManaType]int{},
	}
	return amount
}

var manaStringRegexp = regexp.MustCompile(`\{(\d+|[WUBRGC])\}`)

func ParseManaString(manaString string) (Amount, error) {
	amount := NewAmount()
	if manaString == "" {
		return amount, nil
	}
	matches := manaStringRegexp.FindAllStringSubmatch(manaString, -1)
	if matches == nil {
		return amount, fmt.Errorf("invalid mana string: %s", manaString)
	}
	for _, match := range matches {
		symbol := match[1]
		switch symbol {
		case "W":
			amount.colors[White]++
		case "U":
			amount.colors[Blue]++
		case "B":
			amount.colors[Black]++
		case "R":
			amount.colors[Red]++
		case "G":
			amount.colors[Green]++
		case "C":
			amount.colors[Colorless]++
		default:
			n, err := strconv.Atoi(symbol)
			if err != nil {
				return amount, fmt.Errorf("invalid mana symbol %q", symbol)
			}
			amount.generic += n
		}
	}
	return amount, nil
}

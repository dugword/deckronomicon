package cost

import (
	"deckronomicon/packages/mana"
	"fmt"
	"regexp"
)

// ParseManaCost parses a mana cost string and returns a ManaCost.
func ParseManaCost(costStr string) (ManaCost, error) {
	amount, err := mana.ParseManaString(costStr)
	if err != nil {
		return ManaCost{}, fmt.Errorf("failed to parse mana cost %q: %w", costStr, err)
	}
	return NewManaCost(amount), nil
}

// ManaPattern is a regex pattern that matches valid mana costs.
// TODO: Support X costs and other special cases.
var manaPattern = regexp.MustCompile(`^(?:\{[0-9WUBRGC]+\})*$`)

// isManaCost checks if the input string is a valid mana cost.
// THIS IS ISMANACOST!
// TODO
func IsManaCost(input string) bool {
	return manaPattern.MatchString(input)
}

type ManaCost struct {
	amount mana.Amount
}

func NewManaCost(amount mana.Amount) ManaCost {
	return ManaCost{
		amount: amount,
	}
}

func (c ManaCost) isCost() {}

func (c ManaCost) Amount() mana.Amount {
	return c.amount
}

func (c ManaCost) ManaString() string {
	return c.amount.ManaString()
}

// Description returns a string representation of the mana cost.
// TODO Use a mana cost String method, this feels redundant.
func (c ManaCost) Description() string {
	return fmt.Sprintf("Pay %s", c.ManaString())
}

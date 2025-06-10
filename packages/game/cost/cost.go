package cost

import (
	"deckronomicon/packages/mana"
	"deckronomicon/packages/query"
	"fmt"
	"regexp"
	"strings"
)

// ManaPattern is a regex pattern that matches valid mana costs.
// TODO: Support X costs and other special cases.
var manaPattern = regexp.MustCompile(`^(?:\{[0-9WUBRGC]+\})*$`)

// isManaCost checks if the input string is a valid mana cost.
// THIS IS ISMANACOST!
// TODO
func IsManaCost(input string) bool {
	return manaPattern.MatchString(input)
}

type Cost interface {
	isCost()
	Description() string
}

type TapCost struct{}

// Description returns a string representation of the tap cost.
func (c TapCost) Description() string {
	return "{T}"
}

// GetChoices returns the choices available for this cost.
/*
func (c *TapCost) GetChoices() []gob.Choice {
 // not sure what this will look like yet, but probably have action.complete
 // get input from the player about what to target/sac/or whatever
}


*/

func (c TapCost) isCost() {}

type CompositeCost struct {
	costs []Cost
}

func (c CompositeCost) isCost() {}

func (c CompositeCost) Costs() []Cost {
	return c.costs
}

// Description returns a string representation of the composite cost.
func (c CompositeCost) Description() string {
	// Cost ordered as Mana, Tap, Sacrifice
	var costStrings []string
	// Mana
	for _, cost := range c.costs {
		if _, ok := cost.(ManaCost); ok {
			costStrings = append(costStrings, cost.Description())
		}
	}
	// Tap
	for _, cost := range c.costs {
		if _, ok := cost.(TapCost); ok {
			costStrings = append(costStrings, cost.Description())
		}
	}
	//Sacrifice
	/*
		for _, cost := range c.costs {
			if _, ok := cost.(*SacrificeCost); ok {
				costStrings = append(costStrings, cost.Description())
			}
		}
	*/
	return strings.Join(costStrings, ", ")
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

// Maybe instead of returning a cost we return a "CostSpec" or "CostTemplate" or "CostBuilding"
// that action can then use to provide the cost when the action is executed. Kinda like how we do with command
// parsing with a "IsComplete" method that returns true when the cost is fully specified.

func ParseCost(costString string, source query.Object) (Cost, error) {
	parts := strings.Split(costString, ",")
	var costs []Cost
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		switch {
		case IsManaCost(trimmed):
			manaCost, err := ParseManaCost(trimmed)
			if err != nil {
				return nil, fmt.Errorf("failed to parse mana cost %q: %w", trimmed, err)
			}
			costs = append(costs, manaCost)
		case isTapCost(trimmed):
			costs = append(costs, TapCost{})
		default:
			return nil, fmt.Errorf("unknown cost %q", trimmed)
		}
	}
	if len(costs) == 1 {
		return costs[0], nil
	}
	return CompositeCost{costs: costs}, nil
}

// ParseManaCost parses a mana cost string and returns a ManaCost.
func ParseManaCost(costStr string) (ManaCost, error) {
	amount, err := mana.ParseManaString(costStr)
	if err != nil {
		return ManaCost{}, fmt.Errorf("failed to parse mana cost %q: %w", costStr, err)
	}
	return NewManaCost(amount), nil
}

// isTapCost checks if the input string is a tap cost.
func isTapCost(input string) bool {
	return input == "{T}"
}

// Description returns a string representation of the mana cost.
// TODO Use a mana cost String method, this feels redundant.
func (c ManaCost) Description() string {
	// TODO: colors might not be the right name, costs? Symbols?
	var cs []string
	if c.amount.Generic() > 0 {
		cs = append(cs, fmt.Sprintf("{%d}", c.amount.Generic()))
	}
	for _, color := range mana.AllManaTypes {
		if c.amount.Colors()[color] > 0 {
			for range c.amount.Colors()[color] {
				cs = append(cs, fmt.Sprintf("{%s}", color))
			}
		}
	}
	return strings.Join(cs, "")
}

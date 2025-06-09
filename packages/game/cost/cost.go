package cost

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/mana"
	"deckronomicon/packages/query"
	"fmt"
	"strings"
)

type Cost interface {
	isCost()
	Description() string
}

type TapCost struct {
	permanent gob.Permanent
}

// Description returns a string representation of the tap cost.
func (c TapCost) Description() string {
	return "{T}"
}

func (c TapCost) Permanent() gob.Permanent {
	return c.permanent
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
		case isTapCost(trimmed):
			fmt.Println("Found tap cost:", trimmed)
			costs = append(costs, TapCost{permanent: source.(gob.Permanent)})
		default:
			return nil, fmt.Errorf("unknown cost '%s'", trimmed)
		}
	}
	if len(costs) == 1 {
		fmt.Println("Returning single cost:", costs[0])
		return costs[0], nil
	}
	return CompositeCost{costs: costs}, nil
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

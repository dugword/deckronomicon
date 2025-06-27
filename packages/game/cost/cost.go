package cost

import (
	"fmt"
	"strings"
)

// CombineCosts takes a cost and a list of costs and combines them into a single new composite cost.
// If either the cost or the costs are composite costs, it flattens them into a single composite cost.
// If any cost is nil or if any composite cost is empty they are skipped and the new composite cost
// will not include them. If the cost is nil, it returns a new composite cost with the costs.
func CombineCosts(cost Cost, costs ...Cost) CompositeCost {
	if cost == nil && len(costs) == 0 {
		return CompositeCost{}
	}
	var combinedCosts []Cost
	if cost != nil {
		if compositeCost, ok := cost.(CompositeCost); ok {
			combinedCosts = append(combinedCosts, compositeCost.Costs()...)
		} else {
			combinedCosts = append(combinedCosts, cost)
		}
	}
	for _, c := range costs {
		if c == nil {
			continue
		}
		if compositeCost, ok := c.(CompositeCost); ok {
			for _, subCost := range compositeCost.Costs() {
				if subCost == nil {
					continue
				}
				if _, ok := subCost.(CompositeCost); ok {
					// If the subCost is a composite cost, we need to flatten it.
					// TODO: This should be recursive.
					combinedCosts = append(combinedCosts, subCost.(CompositeCost).Costs()...)
				} else {
					combinedCosts = append(combinedCosts, subCost)
				}
			}
		} else {
			combinedCosts = append(combinedCosts, c)
		}
	}
	if len(combinedCosts) == 0 {
		return CompositeCost{}
	}
	return NewCompositeCost(combinedCosts...)
}

// Has checks if the cost has a specific cost type. E.g.
// Has(cost, ManaCost{}) returns true if the cost has a mana cost.
// or Has(cost, TapThisCost{}) returns true if the cost has a tap this cost.
// If the cost is a composite cost, it checks if any of the sub-costs match the cost type.
func HasCostType(cost Cost, costType Cost) bool {
	switch c := cost.(type) {
	case CompositeCost:
		for _, subCost := range c.Costs() {
			if HasCostType(subCost, costType) {
				return true
			}
		}
	case ManaCost:
		_, ok := costType.(ManaCost)
		return ok
	case TapThisCost:
		_, ok := costType.(TapThisCost)
		return ok
	case DiscardThisCost:
		_, ok := costType.(DiscardThisCost)
		return ok
	case LifeCost:
		_, ok := costType.(LifeCost)
		return ok
	default:
		return false
	}
	return false
}

type Cost interface {
	isCost()
	Description() string
}

type CompositeCost struct {
	costs []Cost
}

func NewCompositeCost(costs ...Cost) CompositeCost {
	return CompositeCost{costs: costs}
}

func (c CompositeCost) isCost() {}

func (c CompositeCost) Costs() []Cost {
	return c.costs
}

// Description returns a string representation of the composite cost.
func (c CompositeCost) Description() string {
	// Cost ordered as Mana, Tap, Sacrifice, Discard
	var costStrings []string
	// Mana
	for _, cost := range c.costs {
		if _, ok := cost.(ManaCost); ok {
			costStrings = append(costStrings, cost.Description())
		}
	}
	// Tap
	for _, cost := range c.costs {
		if _, ok := cost.(TapThisCost); ok {
			costStrings = append(costStrings, cost.Description())
		}
	}
	// Discard
	for _, cost := range c.costs {
		if _, ok := cost.(DiscardThisCost); ok {
			costStrings = append(costStrings, cost.Description())
		}
	}
	// Life
	for _, cost := range c.costs {
		if _, ok := cost.(LifeCost); ok {
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

// Maybe instead of returning a cost we return a "CostSpec" or "CostTemplate" or "CostBuilding"
// that action can then use to provide the cost when the action is executed. Kinda like how we do with command
// parsing with a "IsComplete" method that returns true when the cost is fully specified.

func Parse(costString string) (Cost, error) {
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
		case isTapThisCost(trimmed):
			costs = append(costs, TapThisCost{})
		case isDiscardThisCost(trimmed):
			costs = append(costs, DiscardThisCost{})
		case isLifeCost(trimmed):
			lifeCost, err := ParseLifeCost(trimmed)
			if err != nil {
				return nil, fmt.Errorf("failed to parse life cost %q: %w", trimmed, err)
			}
			costs = append(costs, lifeCost)
		default:
			return nil, fmt.Errorf("unknown cost %q", trimmed)
		}
	}
	if len(costs) == 1 {
		return costs[0], nil
	}
	return CompositeCost{costs: costs}, nil
}

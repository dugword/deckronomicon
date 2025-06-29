package cost

import "strings"

type CompositeCost struct {
	costs []Cost
}

// NewComposite takes a list of costs and combines them into a single flat CompositeCost.
// Any nested CompositeCosts are flattened into the new CompositeCost.
func NewComposite(costs ...Cost) CompositeCost {
	var combined []Cost
	for _, c := range costs {
		if c == nil {
			continue
		}
		switch cc := c.(type) {
		case CompositeCost:
			combined = append(combined, NewComposite(cc.Costs()...).Costs()...)
		default:
			combined = append(combined, c)
		}
	}
	return CompositeCost{costs: combined}
}

func (c CompositeCost) isCost() {}

func (c CompositeCost) Costs() []Cost {
	return c.costs
}

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

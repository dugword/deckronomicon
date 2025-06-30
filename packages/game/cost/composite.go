package cost

import "strings"

type Composite struct {
	costs []Cost
}

// NewComposite takes a list of costs and combines them into a single flat Composite.
// Any nested Composites are flattened into the new Composite.
func NewComposite(costs ...Cost) Composite {
	var combined []Cost
	for _, c := range costs {
		if c == nil {
			continue
		}
		switch cc := c.(type) {
		case Composite:
			combined = append(combined, NewComposite(cc.Costs()...).Costs()...)
		default:
			combined = append(combined, c)
		}
	}
	return Composite{costs: combined}
}

func (c Composite) isCost() {}

func (c Composite) Costs() []Cost {
	return c.costs
}

func (c Composite) Description() string {
	// Cost ordered as Mana, Tap, Sacrifice, Discard
	var costStrings []string
	// Mana
	for _, cost := range c.costs {
		if _, ok := cost.(Mana); ok {
			costStrings = append(costStrings, cost.Description())
		}
	}
	// Tap
	for _, cost := range c.costs {
		if _, ok := cost.(TapThis); ok {
			costStrings = append(costStrings, cost.Description())
		}
	}
	// Discard this card
	for _, cost := range c.costs {
		if _, ok := cost.(DiscardThis); ok {
			costStrings = append(costStrings, cost.Description())
		}
	}
	// Discard a card
	for _, cost := range c.costs {
		if _, ok := cost.(DiscardACard); ok {
			costStrings = append(costStrings, cost.Description())
		}
	}
	// Life
	for _, cost := range c.costs {
		if _, ok := cost.(Life); ok {
			costStrings = append(costStrings, cost.Description())
		}
	}
	// Sacrifice this permanent
	for _, cost := range c.costs {
		if _, ok := cost.(SacrificeThis); ok {
			costStrings = append(costStrings, cost.Description())
		}
	}
	// Sacrifice target permanent
	for _, cost := range c.costs {
		if st, ok := cost.(SacrificeTarget); ok {
			costStrings = append(costStrings, st.Description())
		}
	}
	return strings.Join(costStrings, ", ")
}

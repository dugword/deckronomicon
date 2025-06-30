package cost

import "deckronomicon/packages/game/target"

type Cost interface {
	isCost()
	Description() string
}

// TODO: I think I like this pattern a lot,
// Are there other places I can use this?
// Object With Zone? Effect With Target?
type CostWithTarget interface {
	Cost
	Target() target.Target
	WithTarget(target target.Target) CostWithTarget
	TargetSpec() target.TargetSpec
}

// Has checks if the cost has a specific cost type. E.g.
// Has(cost, Mana{}) returns true if the cost has a mana cost.
// or Has(cost, TapThis{}) returns true if the cost has a tap this cost.
// If the cost is a composite cost, it checks if any of the sub-costs match the cost type.
func HasType(cost Cost, costType Cost) bool {
	switch c := cost.(type) {
	case Composite:
		_, ok := costType.(Composite)
		if ok {
			return true
		}
		for _, subCost := range c.Costs() {
			if HasType(subCost, costType) {
				return true
			}
		}
	case DiscardThis:
		_, ok := costType.(DiscardThis)
		return ok
	case DiscardACard:
		_, ok := costType.(DiscardACard)
		return ok
	case Life:
		_, ok := costType.(Life)
		return ok
	case Mana:
		_, ok := costType.(Mana)
		return ok
	case TapThis:
		_, ok := costType.(TapThis)
		return ok
	default:
		return false
	}
	return false
}

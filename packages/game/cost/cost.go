package cost

type Cost interface {
	isCost()
	Description() string
}

// TODO: I think I like this pattern a lot,
// Are there other places I can use this?
// Object With Zone? Effect With Target?
type CostWithTarget interface {
	Cost
	TargetID() string
	WithTargetID(targetID string) CostWithTarget
}

// Has checks if the cost has a specific cost type. E.g.
// Has(cost, ManaCost{}) returns true if the cost has a mana cost.
// or Has(cost, TapThisCost{}) returns true if the cost has a tap this cost.
// If the cost is a composite cost, it checks if any of the sub-costs match the cost type.
func HasType(cost Cost, costType Cost) bool {
	switch c := cost.(type) {
	case CompositeCost:
		_, ok := costType.(CompositeCost)
		if ok {
			return true
		}
		for _, subCost := range c.Costs() {
			if HasType(subCost, costType) {
				return true
			}
		}
	case DiscardThisCost:
		_, ok := costType.(DiscardThisCost)
		return ok
	case DiscardACardCost:
		_, ok := costType.(DiscardACardCost)
		return ok
	case LifeCost:
		_, ok := costType.(LifeCost)
		return ok
	case ManaCost:
		_, ok := costType.(ManaCost)
		return ok
	case TapThisCost:
		_, ok := costType.(TapThisCost)
		return ok
	default:
		return false
	}
	return false
}

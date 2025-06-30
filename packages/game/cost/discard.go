package cost

type DiscardThisCost struct{}

func (c DiscardThisCost) isCost() {}

func (c DiscardThisCost) Description() string {
	return "Discard this card"
}

type DiscardACardCost struct {
	targetID string
}

func (c DiscardACardCost) isCost() {}

func (c DiscardACardCost) Description() string {
	return "Discard a card"
}

func (c DiscardACardCost) TargetID() string {
	return c.targetID
}

func (c DiscardACardCost) WithTargetID(targetID string) CostWithTarget {
	c.targetID = targetID
	return c
}

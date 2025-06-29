package cost

type DiscardThisCost struct{}

func (c DiscardThisCost) isCost() {}

func (c DiscardThisCost) Description() string {
	return "Discard this card"
}

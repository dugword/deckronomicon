package cost

type DiscardThisCost struct{}

func (c DiscardThisCost) Description() string {
	return "Discard this card"
}

func (c DiscardThisCost) isCost() {}

func isDiscardThisCost(input string) bool {
	return input == "Discard this card"
}

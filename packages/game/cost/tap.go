package cost

type TapThisCost struct{}

// Description returns a string representation of the tap cost.
func (c TapThisCost) Description() string {
	return "Tap this permanent"
}

func (c TapThisCost) isCost() {}

// isTapCost checks if the input string is a tap cost.
func isTapThisCost(input string) bool {
	return input == "{T}"
}

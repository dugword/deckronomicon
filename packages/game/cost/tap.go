package cost

type TapThisCost struct{}

func (c TapThisCost) Description() string {
	return "Tap this permanent"
}

func (c TapThisCost) isCost() {}

package cost

type TapThis struct{}

func (c TapThis) Description() string {
	return "Tap this permanent"
}

func (c TapThis) isCost() {}

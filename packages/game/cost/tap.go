package cost

type TapThisCost struct{}

// Description returns a string representation of the tap cost.
func (c TapThisCost) Description() string {
	return "Tap this permanent"
}

// GetChoices returns the choices available for this cost.
/*
func (c *TapCost) GetChoices() []gob.Choice {
 // not sure what this will look like yet, but probably have action.complete
 // get input from the player about what to target/sac/or whatever
}


*/

func (c TapThisCost) isCost() {}

// isTapCost checks if the input string is a tap cost.
func isTapThisCost(input string) bool {
	return input == "{T}"
}

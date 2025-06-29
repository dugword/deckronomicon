package cost

import (
	"fmt"
)

type LifeCost struct {
	amount int
}

func (l LifeCost) isCost() {}

func (l LifeCost) Description() string {
	// Return a string representation of the life cost
	return fmt.Sprintf("Pay %d life", l.amount)
}

func (l LifeCost) Amount() int {
	return l.amount
}

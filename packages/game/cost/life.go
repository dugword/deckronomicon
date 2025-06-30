package cost

import (
	"fmt"
)

type Life struct {
	amount int
}

func (l Life) isCost() {}

func (l Life) Description() string {
	// Return a string representation of the life cost
	return fmt.Sprintf("Pay %d life", l.amount)
}

func (l Life) Amount() int {
	return l.amount
}

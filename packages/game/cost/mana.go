package cost

import (
	"deckronomicon/packages/game/mana"
	"fmt"
)

type Mana struct {
	amount mana.Amount
}

func (c Mana) isCost() {}

func (c Mana) Amount() mana.Amount {
	return c.amount
}

func (c Mana) ManaString() string {
	return c.amount.ManaString()
}

// Description returns a string representation of the mana cost.
func (c Mana) Description() string {
	return fmt.Sprintf("Pay %s", c.ManaString())
}

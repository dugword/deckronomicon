package cost

import (
	"deckronomicon/packages/game/mana"
	"fmt"
)

type ManaCost struct {
	amount mana.Amount
}

func (c ManaCost) isCost() {}

func (c ManaCost) Amount() mana.Amount {
	return c.amount
}

func (c ManaCost) ManaString() string {
	return c.amount.ManaString()
}

// Description returns a string representation of the mana cost.
func (c ManaCost) Description() string {
	return fmt.Sprintf("Pay %s", c.ManaString())
}

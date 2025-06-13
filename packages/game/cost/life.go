package cost

import (
	"fmt"
	"regexp"
	"strconv"
)

var LifeCostPattern = regexp.MustCompile(`^Pay \d+ life$`)

func isLifeCost(input string) bool {
	return LifeCostPattern.MatchString(input)
}

type LifeCost struct {
	amount int
}

func (l LifeCost) Description() string {
	// Return a string representation of the life cost
	return fmt.Sprintf("Pay %d life", l.amount)
}

func (l LifeCost) isCost() {}

func (l LifeCost) Amount() int {
	return l.amount
}

func ParseLifeCost(input string) (LifeCost, error) {
	// Example input: "Pay 3 life"
	re := regexp.MustCompile(`^Pay (\d+) life$`)
	matches := re.FindStringSubmatch(input)
	if len(matches) != 2 {
		return LifeCost{}, fmt.Errorf("invalid life cost format: %s", input)
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return LifeCost{}, fmt.Errorf("invalid life amount: %s", matches[1])
	}
	return LifeCost{amount: amount}, nil
}

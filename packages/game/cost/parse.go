package cost

import (
	"deckronomicon/packages/game/mana"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func Parse(costString string) (Cost, error) {
	parts := strings.Split(costString, ",")
	var costs []Cost
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		switch {
		case isDiscardThisCost(trimmed):
			costs = append(costs, DiscardThisCost{})
		case isLifeCost(trimmed):
			lifeCost, err := parseLifeCost(trimmed)
			if err != nil {
				return nil, fmt.Errorf("failed to parse life cost %q: %w", trimmed, err)
			}
			costs = append(costs, lifeCost)
		case IsManaCost(trimmed):
			manaCost, err := ParseManaCost(trimmed)
			if err != nil {
				return nil, fmt.Errorf("failed to parse mana cost %q: %w", trimmed, err)
			}
			costs = append(costs, manaCost)
		case isTapThisCost(trimmed):
			costs = append(costs, TapThisCost{})
		default:
			return nil, fmt.Errorf("unknown cost %q", trimmed)
		}
	}
	if len(costs) == 1 {
		return costs[0], nil
	}
	return CompositeCost{costs: costs}, nil
}

func ParseManaCost(costStr string) (ManaCost, error) {
	amount, err := mana.ParseManaString(costStr)
	if err != nil {
		return ManaCost{}, fmt.Errorf("failed to parse mana cost %q: %w", costStr, err)
	}
	return ManaCost{
		amount: amount,
	}, nil
}

func parseLifeCost(input string) (LifeCost, error) {
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

func isDiscardThisCost(input string) bool {
	return input == "Discard this card"
}

var manaPattern = regexp.MustCompile(`^(?:\{[0-9WUBRGC]+\})*$`)

func IsManaCost(input string) bool {
	return manaPattern.MatchString(input)
}

func isTapThisCost(input string) bool {
	return input == "{T}"
}

var LifeCostPattern = regexp.MustCompile(`^Pay \d+ life$`)

func isLifeCost(input string) bool {
	return LifeCostPattern.MatchString(input)
}

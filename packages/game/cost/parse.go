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
		case isDiscardThis(trimmed):
			costs = append(costs, DiscardThis{})
		case isDiscardACard(trimmed):
			costs = append(costs, DiscardACard{})
		case isLife(trimmed):
			life, err := parseLife(trimmed)
			if err != nil {
				return nil, fmt.Errorf("failed to parse life cost %q: %w", trimmed, err)
			}
			costs = append(costs, life)
		case IsMana(trimmed):
			mana, err := ParseMana(trimmed)
			if err != nil {
				return nil, fmt.Errorf("failed to parse mana cost %q: %w", trimmed, err)
			}
			costs = append(costs, mana)
		case isSacrificeThis(trimmed):
			costs = append(costs, SacrificeThis{})
		case isSacrificeTarget(trimmed):
			sacrifice, err := parseSacrificeTarget(trimmed)
			if err != nil {
				return nil, fmt.Errorf("failed to parse sacrifice target cost %q: %w", trimmed, err)
			}
			costs = append(costs, sacrifice)
		case isTapThis(trimmed):
			costs = append(costs, TapThis{})
		default:
			return nil, fmt.Errorf("unknown cost %q", trimmed)
		}
	}
	if len(costs) == 1 {
		return costs[0], nil
	}
	return Composite{costs: costs}, nil
}

func ParseMana(costStr string) (Mana, error) {
	amount, err := mana.ParseManaString(costStr)
	if err != nil {
		return Mana{}, fmt.Errorf("failed to parse mana cost %q: %w", costStr, err)
	}
	return Mana{
		amount: amount,
	}, nil
}

func parseLife(input string) (Life, error) {
	// Example input: "Pay 3 life"
	re := regexp.MustCompile(`^Pay (\d+) life$`)
	matches := re.FindStringSubmatch(input)
	if len(matches) != 2 {
		return Life{}, fmt.Errorf("invalid life cost format: %s", input)
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return Life{}, fmt.Errorf("invalid life amount: %s", matches[1])
	}
	return Life{amount: amount}, nil
}

func isDiscardThis(input string) bool {
	return input == "Discard this card"
}

func isDiscardACard(input string) bool {
	return input == "Discard a card"
}

var manaPattern = regexp.MustCompile(`^(?:\{[0-9WUBRGC]+\})*$`)

func IsMana(input string) bool {
	return manaPattern.MatchString(input)
}

func isTapThis(input string) bool {
	return input == "{T}"
}

var LifePattern = regexp.MustCompile(`^Pay \d+ life$`)

func isLife(input string) bool {
	return LifePattern.MatchString(input)
}

func isSacrificeThis(input string) bool {
	return input == "Sacrifice this permanent"
}

func isSacrificeTarget(input string) bool {
	return strings.HasPrefix(input, "Sacrifice target ")
}

func parseSacrificeTarget(input string) (SacrificeTarget, error) {
	// Example input: "Sacrifice target creature"
	re := regexp.MustCompile(`^Sacrifice target (\w+)$`)
	matches := re.FindStringSubmatch(input)
	if len(matches) != 2 {
		return SacrificeTarget{}, fmt.Errorf("invalid sacrifice target format: %s", input)
	}
	targetType := matches[1]
	return SacrificeTarget{cardTypes: []string{targetType}}, nil
}

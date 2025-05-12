package auto

import (
	game "deckronomicon/game"
	"strconv"
	"strings"
)

func hasCardNamed(cards []*game.Card, name string) bool {
	for _, c := range cards {
		if strings.EqualFold(c.Name, name) {
			return true
		}
	}
	return false
}

func hasPermanentNamed(perms []*game.Permanent, name string) bool {
	for _, p := range perms {
		if strings.EqualFold(p.Name, name) {
			return true
		}
	}
	return false
}

func cardsFromPerms(perms []*game.Permanent) []*game.Card {
	cards := make([]*game.Card, 0, len(perms))
	for _, p := range perms {
		cards = append(cards, &game.Card{Object: p.Object})
	}
	return cards
}

func allCardsPresent(names []string, cards []*game.Card) bool {
	for _, name := range names {
		if !hasCardNamed(cards, name) {
			return false
		}
	}
	return true
}

func anyCardPresent(names []string, cards []*game.Card) bool {
	for _, name := range names {
		if hasCardNamed(cards, name) {
			return true
		}
	}
	return false
}

func allCardsAbsent(names []string, cards []*game.Card) bool {
	for _, name := range names {
		if hasCardNamed(cards, name) {
			return false
		}
	}
	return true
}

func anyCardAbsent(names []string, cards []*game.Card) bool {
	for _, name := range names {
		if !hasCardNamed(cards, name) {
			return true
		}
	}
	return false
}

func allGroupsSatisfied(groups [][]string, cards []*game.Card) bool {
	for _, group := range groups {
		found := false
		for _, name := range group {
			if hasCardNamed(cards, name) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func anyGroupSatisfied(groups [][]string, cards []*game.Card) bool {
	for _, group := range groups {
		for _, name := range group {
			if hasCardNamed(cards, name) {
				return true
			}
		}
	}
	return false
}

func allGroupsAbsent(groups [][]string, cards []*game.Card) bool {
	for _, group := range groups {
		allPresent := true
		for _, name := range group {
			if !hasCardNamed(cards, name) {
				allPresent = false
				break
			}
		}
		if allPresent {
			return false
		}
	}
	return true
}

func noGroupFullyPresent(groups [][]string, cards []*game.Card) bool {
	for _, group := range groups {
		groupPresent := true
		for _, name := range group {
			if !hasCardNamed(cards, name) {
				groupPresent = false
				break
			}
		}
		if groupPresent {
			return false
		}
	}
	return true
}

// evaluateIntComparison checks if an integer value satisfies a comparison string, e.g. ">=3", "<5", "==2"
func evaluateIntComparison(actual int, expr string) bool {
	expr = strings.TrimSpace(expr)
	if strings.HasPrefix(expr, ">=") {
		if val, err := strconv.Atoi(strings.TrimPrefix(expr, ">=")); err == nil {
			return actual >= val
		}
	} else if strings.HasPrefix(expr, "<=") {
		if val, err := strconv.Atoi(strings.TrimPrefix(expr, "<=")); err == nil {
			return actual <= val
		}
	} else if strings.HasPrefix(expr, ">") {
		if val, err := strconv.Atoi(strings.TrimPrefix(expr, ">")); err == nil {
			return actual > val
		}
	} else if strings.HasPrefix(expr, "<") {
		if val, err := strconv.Atoi(strings.TrimPrefix(expr, "<")); err == nil {
			return actual < val
		}
	} else if strings.HasPrefix(expr, "==") {
		if val, err := strconv.Atoi(strings.TrimPrefix(expr, "==")); err == nil {
			return actual == val
		}
	} else if strings.HasPrefix(expr, "!=") {
		if val, err := strconv.Atoi(strings.TrimPrefix(expr, "!=")); err == nil {
			return actual != val
		}
	} else if val, err := strconv.Atoi(expr); err == nil {
		return actual == val
	}
	return false // Fallback if the expression is malformed
}

// parseManaCost parses a string like "2UG" into a ManaCost struct
func parseManaCost(input string) game.ManaCost {
	cost := game.ManaCost{
		Generic: 0,
		Colors:  make(map[string]int),
	}

	digits := ""
	for i := 0; i < len(input); i++ {
		c := string(input[i])
		if c >= "0" && c <= "9" {
			digits += c
		} else {
			amount := 1
			if digits != "" {
				amount, _ = strconv.Atoi(digits)
				digits = ""
			}
			cost.Colors[c] += amount
		}
	}

	// Any leftover digits are generic mana
	if digits != "" {
		amount, _ := strconv.Atoi(digits)
		cost.Generic += amount
	}

	return cost
}

func PlanManaActivation(battlefield []*game.Permanent, cost game.ManaCost) ([]game.GameAction, bool) {
	needed := make(map[string]int)
	for color, amt := range cost.Colors {
		needed[color] = amt
	}
	genericNeeded := cost.Generic
	actions := []game.GameAction{}

	used := make(map[int]bool) // index of tapped permanents

	// First pass: fulfill colored mana
	for color, required := range needed {
		for i, perm := range battlefield {
			if used[i] || perm.Tapped {
				continue
			}
			for _, ability := range perm.Object.ActivatedAbilities {
				if !ability.IsManaAbility {
					continue
				}
				// TODO This sucks
			Foo:
				for _, tag := range ability.Tags {
					if tag.Key == "ManaSource" && tag.Value == color {
						actions = append(actions, game.GameAction{
							Type:   game.ActionActivate,
							Target: perm.Object.Name,
						})
						used[i] = true
						required--
						needed[color]--
						if required == 0 {
							break Foo
						}
					}
				}
			}
			if required == 0 {
				break
			}
		}
		if needed[color] > 0 {
			return nil, false // can't satisfy colored mana
		}
	}

	// Second pass: satisfy generic mana
	for i, perm := range battlefield {
		if used[i] || perm.Tapped {
			continue
		}
		for _, ability := range perm.Object.ActivatedAbilities {
			if ability.IsManaAbility {
				actions = append(actions, game.GameAction{
					Type:   game.ActionActivate,
					Target: perm.Object.Name,
				})
				used[i] = true
				genericNeeded--
				break
			}
		}
		if genericNeeded <= 0 {
			break
		}
	}

	if genericNeeded > 0 {
		return nil, false // not enough mana
	}

	return actions, true
}

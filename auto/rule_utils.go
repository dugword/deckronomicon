package auto

import (
	game "deckronomicon/game"
	"strconv"
	"strings"
)

// hasCardNamed checks if a card with the given name exists in the slice of
// cards.
func hasCardNamed(cards []*game.Card, name string) bool {
	for _, c := range cards {
		if strings.EqualFold(c.Name(), name) {
			return true
		}
	}
	return false
}

// hasPermanentNamed checks if a permanent with the given name exists in the
// slice of permanents.
func hasPermanentNamed(perms []*game.Permanent, name string) bool {
	for _, p := range perms {
		if strings.EqualFold(p.Name(), name) {
			return true
		}
	}
	return false
}

// allCardsPresent checks if all cards with the given names exist in the slice
// of cards.
func allCardsPresent(names []string, cards []*game.Card) bool {
	for _, name := range names {
		if !hasCardNamed(cards, name) {
			return false
		}
	}
	return true
}

// anyCardPresent checks if any card with the given names exists in the slice
// of cards.
func anyCardPresent(names []string, cards []*game.Card) bool {
	for _, name := range names {
		if hasCardNamed(cards, name) {
			return true
		}
	}
	return false
}

// anyCardAbsent checks if any card with the given names does not exist in
// the slice of cards.
func allCardsAbsent(names []string, cards []*game.Card) bool {
	for _, name := range names {
		if hasCardNamed(cards, name) {
			return false
		}
	}
	return true
}

// anyCardAbsent checks if any card with the given names does not exist in
// the slice of cards.
func anyCardAbsent(names []string, cards []*game.Card) bool {
	for _, name := range names {
		if !hasCardNamed(cards, name) {
			return true
		}
	}
	return false
}

// allGroupsSatisfied checks if all groups of cards are satisfied by the
// given cards. A group is satisfied if at least one card in the group is
// present in the slice of cards.
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

// anyGroupSatisfied checks if any group of cards is satisfied by the given
// cards. A group is satisfied if at least one card in the group is present
// in the slice of cards.
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

// anyGroupAbsent checks if any group of cards is not satisfied by the
// given cards. A group is not satisfied if none of the cards in the group
// are present in the slice of cards.
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

// noGroupFullyPresent checks if no group of cards is fully present in the
// given cards. A group is fully present if all cards in the group are
// present in the slice of cards.
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

// evaluateIntComparison checks if an integer value satisfies a comparison
// string, e.g. ">=3", "<5", "==2"
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

// PlanManaActivation generates a list of game actions to activate mana
// abilities.
func PlanManaActivation(battlefield []*game.Permanent, cost *game.ManaCost) ([]game.GameAction, bool) {
	// TODO: Handle error
	needed := make(map[game.Color]int)
	for color, amt := range cost.Colors {
		needed[color] = amt
	}
	genericNeeded := cost.Generic
	actions := []game.GameAction{}

	used := make(map[int]bool) // index of tapped permanents

	// First pass: fulfill colored mana
	for color, required := range needed {
		for i, perm := range battlefield {
			if used[i] || perm.IsTapped() {
				continue
			}
			for _, ability := range perm.ActivatedAbilities() {
				if !ability.IsManaAbility() {
					continue
				}
				// TODO This sucks - do I need it? does break break the inner loop or the outer? Is it obvious?
			Foo:
				for _, tag := range ability.Tags() {
					// TODO Handl this error
					tagValue, _ := game.StringToColor(tag.Value)
					if tag.Key == game.AbilityTagManaAbility && tagValue == color {
						actions = append(actions, game.GameAction{
							Type:   game.ActionActivate,
							Target: perm.Name(),
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
		if used[i] || perm.IsTapped() {
			continue
		}
		for _, ability := range perm.ActivatedAbilities() {
			if ability.IsManaAbility() {
				actions = append(actions, game.GameAction{
					Type:   game.ActionActivate,
					Target: perm.Name(),
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

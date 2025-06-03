package mtg

import (
	"fmt"
	"regexp"
	"strings"
)

// ManaPattern is a regex pattern that matches valid mana costs.
// TODO: Support X costs and other special cases.
var manaPattern = regexp.MustCompile(`^(?:\{[0-9WUBRGC]+\})*$`)

// isManaCost checks if the input string is a valid mana cost.
// THIS IS ISMANACOST!
// TODO
func IsMana(input string) bool {
	return manaPattern.MatchString(input)
}

// ManaStringToManaSymbols converts a mana string to a slice of mana symbols.
// The string should be in the format "{W}{U}{B}{R}{G}{C}".
func ManaStringToManaSymbols(mana string) []string {
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(mana, -1)
	var manaSymbols []string
	for _, match := range matches {
		symbol := fmt.Sprintf("{%s}", strings.ToUpper(match[1]))
		manaSymbols = append(manaSymbols, symbol)
	}
	return manaSymbols
}

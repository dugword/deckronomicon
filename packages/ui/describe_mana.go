package ui

// TODO: Maybe this should be in a different package with other "Pretty Print" functions?

// Maybe all UI pretty print functions should be in this file.

// I think a "pretry print" package is a really good idea.

import (
	"fmt"
	"strings"
)

// TODO: this is redundant with the one in packages/mana/mana.go
// But I think I want that one for handing printing mana in the game state,
// like in errors and logs, and this one for the UI.
// I don't want this to be in the packages/mana package because it is specific to the UI,
// and I don't want to import the UI package in the packages/mana package.
// Maybe it should be in a separate package like "prettyprint" or "uiutils"?

var AllManaTypes = []string{
	"W",
	"U",
	"B",
	"R",
	"B",
	"G",
}

func describeManaPool(manaPool map[string]int) string {
	descriptions := []string{}
	for _, manaType := range AllManaTypes {
		for range manaPool[manaType] {
			descriptions = append(descriptions, fmt.Sprintf("{%s}", string(manaType)))
		}
	}
	if len(descriptions) == 0 {
		return "(empty)"
	}
	return strings.Join(descriptions, "")
}

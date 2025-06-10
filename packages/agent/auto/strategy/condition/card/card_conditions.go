package card

import (
	"deckronomicon/packages/game/mtg"
	"slices"
	"strings"
)

// TODO: This should be centralized probably
type GameObject interface {
	Name() string
	HasCardType(cardType mtg.CardType) bool
	HasSubtype(subtype mtg.Subtype) bool
}

// TODO Maybe make this the same as Condition with an evaluate method
type CardCondition interface {
	// TODO This shold return an error,
	Matches(objs []GameObject, defs map[string][]string) bool
}

// TODO Should be CardNameCondition or something like that
// --- Leaf node: matches a specific card name
type NameMatchCondition struct {
	Name string
}

// TODO: Move the expansion logic to parse time
func (c *NameMatchCondition) Matches(objs []GameObject, defs map[string][]string) bool {
	for _, obj := range objs {
		if obj.Name() == c.Name {
			return true
		}
	}
	return false
}

type CardTypeCondition struct {
	CardType string // e.g., "Creature", "Instant", etc.
}

func (c *CardTypeCondition) Matches(objs []GameObject, defs map[string][]string) bool {
	cardType, err := mtg.StringToCardType(c.CardType)
	if err != nil {
		panic(err) // or handle error appropriately
	}
	for _, obj := range objs {
		if obj.HasCardType(cardType) {
			return true
		}
	}
	return false
}

type CardSubtypeCondition struct {
	Subtype string // e.g., "Island", "Elf", etc.
}

func (c *CardSubtypeCondition) Matches(objs []GameObject, defs map[string][]string) bool {
	subtype, err := mtg.StringToSubtype(c.Subtype)
	if err != nil {
		panic(err) // or handle error appropriately
	}
	for _, obj := range objs {
		if obj.HasSubtype(subtype) {
			return true
		}
	}
	return false
}

// --- Logical AND: all conditions must match
type AndCardCondition struct {
	Conditions []CardCondition
}

func (c *AndCardCondition) Matches(objs []GameObject, defs map[string][]string) bool {
	for _, cond := range c.Conditions {
		if !cond.Matches(objs, defs) {
			return false
		}
	}
	return true
}

// --- Logical OR: any condition must match
type OrCardCondition struct {
	Conditions []CardCondition
}

func (c *OrCardCondition) Matches(objs []GameObject, defs map[string][]string) bool {
	for _, cond := range c.Conditions {
		if cond.Matches(objs, defs) {
			return true
		}
	}
	return false
}

// --- Logical NOT: negates a condition
type NotCardCondition struct {
	Condition CardCondition
}

func (c *NotCardCondition) Matches(objs []GameObject, defs map[string][]string) bool {
	return !c.Condition.Matches(objs, defs)
}

// --- Group reference (e.g., $ComboPiece)
type GroupRefCondition struct {
	Group string
}

func (c *GroupRefCondition) Matches(objs []GameObject, defs map[string][]string) bool {
	names := defs[c.Group]
	for _, obj := range objs {
		if slices.Contains(names, obj.Name()) {
			return true
		}
	}
	return false
}

// TODO Make this a thing that takes a thing and wraps it with the group stuff
// --- Factory: make a CardCondition from raw string or group reference
func NewNameConditionOrGroupRef(token string) CardCondition {
	if strings.HasPrefix(token, "$") {
		return &GroupRefCondition{Group: strings.TrimPrefix(token, "$")}
	}
	return &NameMatchCondition{Name: token}
}

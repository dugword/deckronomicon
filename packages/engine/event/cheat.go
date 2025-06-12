package event

const (
	EventTypeCheatAddMana       = "CheatAddMana"
	EventTypeCheatConjureCard   = "CheatConjureCard"
	EventTypeCheatDiscard       = "CheatDiscard"
	EventTypeCheatDraw          = "CheatDraw"
	EventTypeCheatEnabled       = "CheatEnabled"
	EventTypeCheatFindCard      = "CheatFindCard"
	EventTypeCheatPeek          = "CheatPeek"
	EventTypeCheatResetLandDrop = "CheatResetLandDrop"
	EventTypeCheatShuffleDeck   = "CheatShuffleDeck"
	EventTypeCheatUntap         = "CheatUntap"
)

type CheatEvent interface{ isCheatEvent() }

type CheatEventBase struct{}

func (e CheatEventBase) isCheatEvent() {}

type CheatAddManaEvent struct {
	CheatEventBase
	Player string // Player ID who added mana
}

func (e CheatAddManaEvent) EventType() string {
	return EventTypeCheatAddMana
}

type CheatConjureCardEvent struct {
	CheatEventBase
	PlayerID string // Player ID who conjured the card
	CardName string // Name of the card conjured
}

func (e CheatConjureCardEvent) EventType() string {
	return "ConjureCard"
}

type CheatDiscardEvent struct {
	CheatEventBase
	PlayerID string // Player ID who discarded the card
}

func (e CheatDiscardEvent) EventType() string {
	return EventTypeCheatDiscard
}

type CheatDrawEvent struct {
	CheatEventBase
	PlayerID string // Player ID who drew the card
}

func (e CheatDrawEvent) EventType() string {
	return EventTypeCheatDraw
}

type CheatPeekEvent struct {
	CheatEventBase
	PlayerID string
}

func (e CheatPeekEvent) EventType() string {
	return EventTypeCheatPeek
}

type CheatFindCardEvent struct {
	CheatEventBase
	PlayerID string
	CardID   string
}

func (e CheatFindCardEvent) EventType() string {
	return EventTypeCheatFindCard
}

type CheatResetLandDropEvent struct {
	CheatEventBase
	PlayerID string
}

func (e CheatResetLandDropEvent) EventType() string {
	return EventTypeCheatResetLandDrop
}

type CheatShuffleDeckEvent struct {
	CheatEventBase
	PlayerID string
}

func (e CheatShuffleDeckEvent) EventType() string {
	return EventTypeCheatShuffleDeck
}

type CheatUntapEvent struct {
	CheatEventBase
	PlayerID string
}

func (e CheatUntapEvent) EventType() string {
	return EventTypeCheatUntap
}

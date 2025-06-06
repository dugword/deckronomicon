package event

// Events go in the event record log so they can be replayed later. That means
// the properties should be serializable and public.

const (
	EventTypeDrawCard     = "DrawCard"
	EventTypeShuffleDeck  = "ShuffleDeck"
	EventDrawStartingHand = "DrawStartingHand"
	EVentTypePlayLand     = "PlayLand"
)

// TODO: maybe use typed constants for event types
const (
	EventTypeSetNextPlayer = "SetNextPlayer"
)

type GameEvent interface {
	EventType() string
}

type Source interface {
	Name() string
}

type CastSpellEvent struct {
	PlayerID string
	CardID   string
}

func (e CastSpellEvent) EventType() string {
	return "CastSpell"
}

type PlayLandEvent struct {
	PlayerID string
	CardID   string
}

func (e PlayLandEvent) EventType() string {
	return EVentTypePlayLand
}

type DrawStartingHandEvent struct {
	PlayerID string
}

func (e DrawStartingHandEvent) EventType() string {
	return EventDrawStartingHand
}

func NewDrawStartingHandEvent(playerID string) DrawStartingHandEvent {
	return DrawStartingHandEvent{
		PlayerID: playerID,
	}
}

type ShuffleDeckEvent struct {
	PlayerID string
}

func (e ShuffleDeckEvent) EventType() string {
	return EventTypeShuffleDeck
}

func NewShuffDeckEvent(playerID string) ShuffleDeckEvent {
	return ShuffleDeckEvent{
		PlayerID: playerID,
	}
}

type SetNextPlayerEvent struct {
}

func (e SetNextPlayerEvent) EventType() string {
	return EventTypeSetNextPlayer
}

type DiscardCardEvent struct {
	PlayerID string
	Source   Source
	CardID   string
}

func (e DiscardCardEvent) EventType() string {
	return "DiscardCard"
}

type DrawCardEvent struct {
	PlayerID string
	Source   Source
}

func (e DrawCardEvent) EventType() string {
	return EventTypeDrawCard
}

func NewDrawCardEvent(playerID string) DrawCardEvent {
	return DrawCardEvent{
		PlayerID: playerID,
		Source:   nil, // Source can be set later if needed
	}
}

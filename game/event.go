package game

import "fmt"

type Event struct {
	Type   EventType
	Source *Object
	// TODO: Better this
	Data map[string]interface{} // anything extra (like how it died, what it targeted, etc.)
}

type EventType string

const (
	EventCastSpell        EventType = "CastSpell"
	EventEnterBattlefield EventType = "EnterBattlefield"
	EventLeaveBattlefield EventType = "LeaveBattlefield"
	EventDrawCard         EventType = "DrawCard"
	EventTapPermanent     EventType = "TapPermanent"
	EventActivateAbility  EventType = "ActivateAbility"
	EventEndStep          EventType = "EndStep"
	EventUpkeep           EventType = "Upkeep"
)

type EventHandler func(Event, *GameState, ChoiceResolver)

// GameState additions for event system
func (g *GameState) RegisterListener(listener EventHandler) {
	g.EventListeners = append(g.EventListeners, listener)
}

func (g *GameState) EmitEvent(evt Event, resolver ChoiceResolver) {
	for _, listener := range g.EventListeners {
		listener(evt, g, resolver)
	}
}

// Hook for triggered abilities
func NewTriggeredListener(ta TriggeredAbility) EventHandler {
	return func(event Event, state *GameState, resolver ChoiceResolver) {
		if ta.TriggerCondition(event) {
			fmt.Println("Triggering:", ta.Description)
			ta.Effect(state, resolver)
		}
	}
}

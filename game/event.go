package game

// TODO: Rework this file with better handling for wrapper functions and ID
// managment for cleanup.

// Event represents an event in the game.
type Event struct {
	Type   EventType
	Source GameObject
	// TODO: Better this
	Data map[string]interface{} // anything extra (like how it died, what it targeted, etc.)
}

// EventType represents the type of event.
type EventType string

const (
	EventActivateAbility  EventType = "ActivateAbility"
	EventCastSpell        EventType = "CastSpell"
	EventDrawCard         EventType = "DrawCard"
	EventEnterBattlefield EventType = "EnterBattlefield"
	EventLeaveBattlefield EventType = "LeaveBattlefield"
	EventTapForMana       EventType = "TapForMana"
	EventTapPermanent     EventType = "TapPermanent"
	//Phase
	EventPostcombatMainPhase EventType = "PrecombatMainPhase"
	EventPrecombatMainPhase  EventType = "PrecombatMainPhase"
	// Step
	EventBeginningOfCombatStep EventType = "BeginningOfCombatStep"
	EventCleanupStep           EventType = "CleanupStep"
	EventCombatDamageStep      EventType = "CombatDamageStep"
	EventDeclareAttackersStep  EventType = "DeclareAttackersStep"
	EventDeclareBlockersStep   EventType = "DeclareBlockersStep"
	EventDrawStep              EventType = "DrawStep"
	EventEndOfCombatStep       EventType = "EndOfCombatStep"
	EventEndStep               EventType = "EndStep"
	EventUntapStep             EventType = "UntapStep"
	EventUpkeepStep            EventType = "UpkeepStep"
)

var nextEventID int

// TODO This probably isn't goroutine safe
// also should use UUID like value
func getNextEventID() int {
	nextEventID++
	return nextEventID
}

// EventListener is a function that handles events.
type EventHandler struct {
	// TODO: UUID
	ID       int
	Callback func(Event, *GameState, *Player)
}

// RegisterListener registers a new event listener.
func (g *GameState) RegisterListener(listener EventHandler) {
	g.EventListeners = append(g.EventListeners, listener)
}

// DeregisterListener
func (g *GameState) DeregisterListener(id int) {
	for i, l := range g.EventListeners {
		if l.ID == id {
			g.EventListeners = append(g.EventListeners[:i], g.EventListeners[i+1:]...)
			break
		}
	}
}

// RegisterOneShotListener registers a one-shot event listener.
func (g *GameState) RegisterOneShotListener(listener EventHandler) {
	wrappedListener := EventHandler{
		ID: listener.ID,
		Callback: func(event Event, state *GameState, player *Player) {
			// Call the original listener's callback
			listener.Callback(event, state, player)
			// Deregister the listener after it has been called
			state.DeregisterListener(listener.ID)
		},
	}

	g.RegisterListener(wrappedListener)
}

// RegisterListenenerUntil registers a listener until a specific event occurs.
func (g *GameState) RegisterListenerUntil(listener EventHandler, untilEvent EventType) {
	cleanUpID := getNextEventID()
	cleanUpHandler := EventHandler{
		ID: cleanUpID,
		Callback: func(event Event, state *GameState, player *Player) {
			if event.Type == untilEvent {
				state.DeregisterListener(listener.ID)
				state.DeregisterListener(cleanUpID)
			}
		},
	}
	g.RegisterListener(listener)
	g.RegisterListener(cleanUpHandler)
}

// EmitEvent emits an event to all registered listeners.
func (g *GameState) EmitEvent(evt Event, player *Player) {
	for _, listener := range g.EventListeners {
		listener.Callback(evt, g, player)
	}
}

// TODO This probably isn't goroutine safe
// Hook for triggered abilities
func NewTriggeredListener(ta TriggeredAbility) (EventHandler, int) {
	id := getNextEventID()
	return EventHandler{
		ID: id,
		Callback: func(event Event, state *GameState, player *Player) {
			if ta.TriggerCondition(event) {
				ta.Resolve(state, player)
			}
		},
	}, id
}

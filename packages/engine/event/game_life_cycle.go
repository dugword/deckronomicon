package event

import "deckronomicon/packages/game/mtg"

const (
	EventTypeBeginGame = "BeginGame"
	EventTypeEndGame   = "EndGame"
)

const (
	EventTypeBeginTurn = "BeginTurn"
	EventTypeEndTurn   = "EndTurn"
)

const (
	// Beginning Phase
	EventTypeBeginBeginningPhase = "BeginBeginningPhase"
	EventTypeEndBeginningPhase   = "EndBeginningPhase"
	EventTypeBeginUntapStep      = "BeginUntapStep"
	EventTypeEndUntapStep        = "EndUntapStep"
	EventTypeBeginUpkeepStep     = "BeginUpkeepStep"
	EventTypeEndUpkeepStep       = "EndUpkeepStep"
	EventTypeBeginDrawStep       = "BeginDrawStep"
	EventTypeEndDrawStep         = "EndDrawStep"
)

const (
	// Precombat Main Phase
	EventTypeBeginPrecombatMainPhase = "BeginPrecombatMainPhase"
	EventTypeEndPrecombatMainPhase   = "EndPrecombatMainPhase"
	EventTypeBeginPrecombatMainStep  = "BeginPrecombatMainStep"
	EventTypeEndPrecombatMainStep    = "EndPrecombatMainStep"
)

const (
	// Combat Phase
	EventTypeBeginCombatPhase           = "BeginCombatPhase"
	EventTypeEndCombatPhase             = "EndCombatPhase"
	EventTypeBeginBeginningOfCombatStep = "BeginBeginningOfCombatStep"
	EventTypeEndBeginningOfCombatStep   = "EndBeginningOfCombatStep"
	EventTypeBeginDeclareAttackersStep  = "BeginDeclareAttackersStep"
	EventTypeEndDeclareAttackersStep    = "EndDeclareAttackersStep"
	EventTypeBeginDeclareBlockersStep   = "BeginDeclareBlockersStep"
	EventTypeEndDeclareBlockersStep     = "EndDeclareBlockersStep"
	EventTypeBeginCombatDamageStep      = "BeginCombatDamageStep"
	EentTypeEndCombatDamageStep         = "EndCombatDamageStep"
	EventTypeBeginEndOfCombatStep       = "BeginEndOfCombatStep"
	EventTypeEndEndOfCombatStep         = "EndEndOfCombatStep"
)

const (
	// Postcombat Main Phase
	EventTypeBeginPostcombatMainPhase = "BeginPostcombatMainPhase"
	EventTypeEndPostcombatMainPhase   = "EndPostcombatMainPhase"
	EventTypeBeginPostcombatMainStep  = "BeginPostcombatMainStep"
	EventTypeEndPostcombatMainStep    = "EndPostcombatMainStep"
)

const (
	// Ending Phase
	EventTypeBeginEndingPhase = "BeginEndingPhase"
	EventTypeEndEndingPhase   = "EndEndingPhase"
	EventTypeBeginEndStep     = "BeginEndStep"
	EventTypeEndEndStep       = "EndEndStep"
	EventTypeBeginCleanupStep = "BeginCleanupStep"
	EventTypeEndCleanupStep   = "EndCleanupStep"
)

type GameLifecycleEvent interface{ isGameLifecycleEvent() }

type GameLifecycleBaseEvent struct{}

func (GameLifecycleBaseEvent) isGameLifecycleEvent() {}

type TurnEvent interface{ isTurnEvent() }

type TurnBaseEvent struct {
	GameLifecycleBaseEvent
}

func (e TurnBaseEvent) isTurnEvent() {}

type BeginPhaseEvent interface{ isBeginPhaseEvent() }

type BeginPhaseBaseEvent struct{ TurnBaseEvent }

func (e BeginPhaseBaseEvent) isBeginPhaseEvent() {}

type EndPhaseEvent interface{ isEndPhaseEvent() }

type EndPhaseBaseEvent struct{ TurnBaseEvent }

func (e EndPhaseBaseEvent) isEndPhaseEvent() {}

type BeginStepEvent interface{ isBeginStepEvent() }

type BeginStepBaseEvent struct{ TurnBaseEvent }

func (e BeginStepBaseEvent) isBeginStepEvent() {}

type EndStepEvent interface{ isEndStepEvent() }

type EndStepBaseEvent struct{ TurnBaseEvent }

func (e EndStepBaseEvent) isEndStepEvent() {}

type BeginGameEvent struct {
	GameLifecycleBaseEvent
}

func (e BeginGameEvent) EventType() string {
	return EventTypeBeginGame
}

type EndGameEvent struct {
	GameLifecycleBaseEvent
	WinnerID string
}

func (e EndGameEvent) EventType() string {
	return EventTypeEndGame
}

type BeginTurnEvent struct {
	TurnBaseEvent
	PlayerID string
}

func (e BeginTurnEvent) EventType() string {
	return EventTypeBeginTurn
}

type EndTurnEvent struct {
	TurnBaseEvent
	PlayerID string
}

func (e EndTurnEvent) EventType() string {
	return EventTypeEndTurn
}

type BeginBeginningPhaseEvent struct {
	BeginPhaseBaseEvent
	PlayerID string
}

func (e BeginBeginningPhaseEvent) EventType() string {
	return EventTypeBeginBeginningPhase
}

type BeginPrecombatMainPhaseEvent struct {
	BeginPhaseBaseEvent
	PlayerID string
}

func (e BeginPrecombatMainPhaseEvent) EventType() string {
	return EventTypeBeginPrecombatMainPhase
}

type BeginCombatPhaseEvent struct {
	BeginPhaseBaseEvent
	PlayerID string
}

func (e BeginCombatPhaseEvent) EventType() string {
	return EventTypeBeginCombatPhase
}

type BeginPostcombatMainPhaseEvent struct {
	BeginPhaseBaseEvent
	PlayerID string
}

func (e BeginPostcombatMainPhaseEvent) EventType() string {
	return EventTypeBeginPostcombatMainPhase
}

type BeginEndingPhaseEvent struct {
	BeginPhaseBaseEvent
	PlayerID string
}

func (e BeginEndingPhaseEvent) EventType() string {
	return EventTypeBeginEndingPhase
}

type EndBeginningPhaseEvent struct {
	EndPhaseBaseEvent
	PlayerID string
}

func (e EndBeginningPhaseEvent) EventType() string {
	return EventTypeEndBeginningPhase
}

type EndPrecombatMainPhaseEvent struct {
	EndPhaseBaseEvent
	PlayerID string
}

func (e EndPrecombatMainPhaseEvent) EventType() string {
	return EventTypeEndPrecombatMainPhase
}

type EndCombatPhaseEvent struct {
	EndPhaseBaseEvent
	PlayerID string
}

func (e EndCombatPhaseEvent) EventType() string {
	return EventTypeEndCombatPhase
}

type EndPostcombatMainPhaseEvent struct {
	EndPhaseBaseEvent
	PlayerID string
}

func (e EndPostcombatMainPhaseEvent) EventType() string {
	return EventTypeEndPostcombatMainPhase
}

type EndEndingPhaseEvent struct {
	EndPhaseBaseEvent
	PlayerID string
}

func (e EndEndingPhaseEvent) EventType() string {
	return EventTypeEndEndingPhase
}

type BeginUntapStepEvent struct {
	BeginStepBaseEvent
	PlayerID string
}

func (e BeginUntapStepEvent) EventType() string {
	return EventTypeBeginUntapStep
}

type EndUntapStepEvent struct {
	EndStepBaseEvent
	PlayerID string
}

func (e EndUntapStepEvent) EventType() string {
	return EventTypeEndUntapStep
}

type BeginUpkeepStepEvent struct {
	BeginStepBaseEvent
	PlayerID string
}

func (e BeginUpkeepStepEvent) EventType() string {
	return EventTypeBeginUpkeepStep
}

type EndUpkeepStepEvent struct {
	EndStepBaseEvent
	PlayerID string
}

func (e EndUpkeepStepEvent) EventType() string {
	return EventTypeEndUpkeepStep
}

type BeginDrawStepEvent struct {
	BeginStepBaseEvent
	PlayerID string
}

func (e BeginDrawStepEvent) EventType() string {
	return EventTypeBeginDrawStep
}

type EndDrawStepEvent struct {
	EndStepBaseEvent
	PlayerID string
}

func (e EndDrawStepEvent) EventType() string {
	return EventTypeEndDrawStep
}

type BeginPrecombatMainStepEvent struct {
	BeginStepBaseEvent
	PlayerID string
}

func (e BeginPrecombatMainStepEvent) EventType() string {
	return EventTypeBeginPrecombatMainStep
}

type EndPrecombatMainStepEvent struct {
	EndStepBaseEvent
	PlayerID string
}

func (e EndPrecombatMainStepEvent) EventType() string {
	return EventTypeEndPrecombatMainStep
}

type BeginBeginningOfCombatStepEvent struct {
	BeginStepBaseEvent
	PlayerID string
}

func (e BeginBeginningOfCombatStepEvent) EventType() string {
	return EventTypeBeginBeginningOfCombatStep
}

type EndBeginningOfCombatStepEvent struct {
	EndStepBaseEvent
	PlayerID string
}

func (e EndBeginningOfCombatStepEvent) EventType() string {
	return EventTypeEndBeginningOfCombatStep
}

type BeginDeclareAttackersStepEvent struct {
	BeginStepBaseEvent
	PlayerID string
}

func (e BeginDeclareAttackersStepEvent) EventType() string {
	return EventTypeBeginDeclareAttackersStep
}

type EndDeclareAttackersStepEvent struct {
	EndStepBaseEvent
	PlayerID string
}

func (e EndDeclareAttackersStepEvent) EventType() string {
	return EventTypeEndDeclareAttackersStep
}

type BeginDeclareBlockersStepEvent struct {
	BeginStepBaseEvent
	PlayerID string
}

func (e BeginDeclareBlockersStepEvent) EventType() string {
	return EventTypeBeginDeclareBlockersStep
}

type EndDeclareBlockersStepEvent struct {
	EndStepBaseEvent
	PlayerID string
}

func (e EndDeclareBlockersStepEvent) EventType() string {
	return EventTypeEndDeclareBlockersStep
}

type BeginCombatDamageStepEvent struct {
	BeginStepBaseEvent
	PlayerID string
}

func (e BeginCombatDamageStepEvent) EventType() string {
	return EventTypeBeginCombatDamageStep
}

type EndCombatDamageStepEvent struct {
	EndStepBaseEvent
	PlayerID string
}

func (e EndCombatDamageStepEvent) EventType() string {
	return EentTypeEndCombatDamageStep
}

type BeginEndOfCombatStepEvent struct {
	BeginStepBaseEvent
	PlayerID string
}

func (e BeginEndOfCombatStepEvent) EventType() string {
	return EventTypeBeginEndOfCombatStep
}

type EndEndOfCombatStepEvent struct {
	EndStepBaseEvent
	PlayerID string
}

func (e EndEndOfCombatStepEvent) EventType() string {
	return EventTypeEndEndOfCombatStep
}

type BeginPostcombatMainStepEvent struct {
	BeginStepBaseEvent
	PlayerID string
}

func (e BeginPostcombatMainStepEvent) EventType() string {
	return EventTypeBeginPostcombatMainStep
}

type EndPostcombatMainStepEvent struct {
	EndStepBaseEvent
	PlayerID string
}

func (e EndPostcombatMainStepEvent) EventType() string {
	return EventTypeEndPostcombatMainStep
}

type BeginEndStepEvent struct {
	BeginStepBaseEvent
	PlayerID string
}

func (e BeginEndStepEvent) EventType() string {
	return EventTypeBeginEndStep
}

type EndEndStepEvent struct {
	EndStepBaseEvent
	PlayerID string
}

func (e EndEndStepEvent) EventType() string {
	return EventTypeEndEndStep
}

type BeginCleanupStepEvent struct {
	BeginStepBaseEvent
	PlayerID string
}

func (e BeginCleanupStepEvent) EventType() string {
	return EventTypeBeginCleanupStep
}

type EndCleanupStepEvent struct {
	EndStepBaseEvent
	PlayerID string
}

func (e EndCleanupStepEvent) EventType() string {
	return EventTypeEndCleanupStep
}

func NewBeginStepEvent(step mtg.Step, playerID string) GameEvent {
	switch step {
	case mtg.StepUntap:
		return BeginUntapStepEvent{PlayerID: playerID}
	case mtg.StepUpkeep:
		return BeginUpkeepStepEvent{PlayerID: playerID}
	case mtg.StepDraw:
		return BeginDrawStepEvent{PlayerID: playerID}
	case mtg.StepPrecombatMain:
		return BeginPrecombatMainStepEvent{PlayerID: playerID}
	case mtg.StepBeginningOfCombat:
		return BeginBeginningOfCombatStepEvent{PlayerID: playerID}
	case mtg.StepDeclareAttackers:
		return BeginDeclareAttackersStepEvent{PlayerID: playerID}
	case mtg.StepDeclareBlockers:
		return BeginDeclareBlockersStepEvent{PlayerID: playerID}
	case mtg.StepCombatDamage:
		return BeginCombatDamageStepEvent{PlayerID: playerID}
	case mtg.StepEndOfCombat:
		return BeginEndOfCombatStepEvent{PlayerID: playerID}
	case mtg.StepPostcombatMain:
		return BeginPostcombatMainStepEvent{PlayerID: playerID}
	case mtg.StepEnd:
		return BeginEndStepEvent{PlayerID: playerID}
	case mtg.StepCleanup:
		return BeginCleanupStepEvent{PlayerID: playerID}
	default:
		panic("unknown step")
	}
}

func NewEndStepEvent(step mtg.Step, playerID string) GameEvent {
	switch step {
	case mtg.StepUntap:
		return EndUntapStepEvent{PlayerID: playerID}
	case mtg.StepUpkeep:
		return EndUpkeepStepEvent{PlayerID: playerID}
	case mtg.StepDraw:
		return EndDrawStepEvent{PlayerID: playerID}
	case mtg.StepPrecombatMain:
		return EndPrecombatMainStepEvent{PlayerID: playerID}
	case mtg.StepBeginningOfCombat:
		return EndBeginningOfCombatStepEvent{PlayerID: playerID}
	case mtg.StepDeclareAttackers:
		return EndDeclareAttackersStepEvent{PlayerID: playerID}
	case mtg.StepDeclareBlockers:
		return EndDeclareBlockersStepEvent{PlayerID: playerID}
	case mtg.StepCombatDamage:
		return EndCombatDamageStepEvent{PlayerID: playerID}
	case mtg.StepEndOfCombat:
		return EndEndOfCombatStepEvent{PlayerID: playerID}
	case mtg.StepPostcombatMain:
		return EndPostcombatMainStepEvent{PlayerID: playerID}
	case mtg.StepEnd:
		return EndEndStepEvent{PlayerID: playerID}
	case mtg.StepCleanup:
		return EndCleanupStepEvent{PlayerID: playerID}
	default:
		panic("unknown step")
	}
}

// TODO Should this return a BeginPhaseEvent?
func NewBeginPhaseEvent(phase mtg.Phase, playerID string) GameEvent {
	switch phase {
	case mtg.PhaseBeginning:
		return BeginBeginningPhaseEvent{PlayerID: playerID}
	case mtg.PhasePrecombatMain:
		return BeginPrecombatMainPhaseEvent{PlayerID: playerID}
	case mtg.PhaseCombat:
		return BeginCombatPhaseEvent{PlayerID: playerID}
	case mtg.PhasePostcombatMain:
		return BeginPostcombatMainPhaseEvent{PlayerID: playerID}
	case mtg.PhaseEnding:
		return BeginEndingPhaseEvent{PlayerID: playerID}
	default:
		panic("unknown phase")
	}
}

func NewEndPhaseEvent(phase mtg.Phase, playerID string) GameEvent {
	switch phase {
	case mtg.PhaseBeginning:
		return EndBeginningPhaseEvent{PlayerID: playerID}
	case mtg.PhasePrecombatMain:
		return EndPrecombatMainPhaseEvent{PlayerID: playerID}
	case mtg.PhaseCombat:
		return EndCombatPhaseEvent{PlayerID: playerID}
	case mtg.PhasePostcombatMain:
		return EndPostcombatMainPhaseEvent{PlayerID: playerID}
	case mtg.PhaseEnding:
		return EndEndingPhaseEvent{PlayerID: playerID}
	default:
		panic("unknown phase")
	}
}

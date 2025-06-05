package event

import "deckronomicon/packages/game/mtg"

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

type TurnEvent interface {
	isTurnEvent()
}

type BeginPhaseEvent interface {
	isBeginPhaseEvent()
}

type EndPhaseEvent interface {
	isEndPhaseEvent()
}

type BeginStepEvent interface {
	isBeginStepEvent()
}

type EndStepEvent interface {
	isEndStepEvent()
}

type TurnEventBase struct {
	PlayerID string
}

func (e TurnEventBase) isTurnEvent() {}

type BeginPhaseEventBase struct {
	TurnEventBase
}

func (e BeginPhaseEventBase) isBeginPhaseEvent() {}

type BeginBeginningPhaseEvent struct {
	BeginPhaseEventBase
}

func (e BeginBeginningPhaseEvent) EventType() string {
	return EventTypeBeginBeginningPhase
}

type BeginPrecombatMainPhaseEvent struct {
	BeginPhaseEventBase
}

func (e BeginPrecombatMainPhaseEvent) EventType() string {
	return EventTypeBeginPrecombatMainPhase
}

type BeginCombatPhaseEvent struct {
	BeginPhaseEventBase
}

func (e BeginCombatPhaseEvent) EventType() string {
	return EventTypeBeginCombatPhase
}

type BeginPostcombatMainPhaseEvent struct {
	BeginPhaseEventBase
}

func (e BeginPostcombatMainPhaseEvent) EventType() string {
	return EventTypeBeginPostcombatMainPhase
}

type BeginEndingPhaseEvent struct {
	BeginPhaseEventBase
}

func (e BeginEndingPhaseEvent) EventType() string {
	return EventTypeBeginEndingPhase
}

type EndPhaseEventBase struct {
	TurnEventBase
}

func (e EndPhaseEventBase) isEndPhaseEvent() {}

type EndBeginningPhaseEvent struct {
	EndPhaseEventBase
}

func (e EndBeginningPhaseEvent) EventType() string {
	return EventTypeEndBeginningPhase
}

type EndPrecombatMainPhaseEvent struct {
	EndPhaseEventBase
}

func (e EndPrecombatMainPhaseEvent) EventType() string {
	return EventTypeEndPrecombatMainPhase
}

type EndCombatPhaseEvent struct {
	EndPhaseEventBase
}

func (e EndCombatPhaseEvent) EventType() string {
	return EventTypeEndCombatPhase
}

type EndPostcombatMainPhaseEvent struct {
	EndPhaseEventBase
}

func (e EndPostcombatMainPhaseEvent) EventType() string {
	return EventTypeEndPostcombatMainPhase
}

type EndEndingPhaseEvent struct {
	EndPhaseEventBase
}

func (e EndEndingPhaseEvent) EventType() string {
	return EventTypeEndEndingPhase
}

type BeginStepEventBase struct {
	TurnEventBase
}

func (e BeginStepEventBase) isBeginStepEvent() {}

type EndStepEventBase struct {
	TurnEventBase
}

func (e EndStepEventBase) isEndStepEvent() {}

type BeginUntapStepEvent struct {
	BeginStepEventBase
}

func (e BeginUntapStepEvent) EventType() string {
	return EventTypeBeginUntapStep
}

type EndUntapStepEvent struct {
	EndStepEventBase
}

func (e EndUntapStepEvent) EventType() string {
	return EventTypeEndUntapStep
}

type BeginUpkeepStepEvent struct {
	BeginStepEventBase
}

func (e BeginUpkeepStepEvent) EventType() string {
	return EventTypeBeginUpkeepStep
}

type EndUpkeepStepEvent struct {
	EndStepEventBase
}

func (e EndUpkeepStepEvent) EventType() string {
	return EventTypeEndUpkeepStep
}

type BeginDrawStepEvent struct {
	BeginStepEventBase
}

func (e BeginDrawStepEvent) EventType() string {
	return EventTypeBeginDrawStep
}

type EndDrawStepEvent struct {
	EndStepEventBase
}

func (e EndDrawStepEvent) EventType() string {
	return EventTypeEndDrawStep
}

type BeginPrecombatMainStepEvent struct {
	BeginStepEventBase
}

func (e BeginPrecombatMainStepEvent) EventType() string {
	return EventTypeBeginPrecombatMainStep
}

type EndPrecombatMainStepEvent struct {
	EndStepEventBase
}

func (e EndPrecombatMainStepEvent) EventType() string {
	return EventTypeEndPrecombatMainStep
}

type BeginBeginningOfCombatStepEvent struct {
	BeginStepEventBase
}

func (e BeginBeginningOfCombatStepEvent) EventType() string {
	return EventTypeBeginBeginningOfCombatStep
}

type EndBeginningOfCombatStepEvent struct {
	EndStepEventBase
}

func (e EndBeginningOfCombatStepEvent) EventType() string {
	return EventTypeEndBeginningOfCombatStep
}

type BeginDeclareAttackersStepEvent struct {
	BeginStepEventBase
}

func (e BeginDeclareAttackersStepEvent) EventType() string {
	return EventTypeBeginDeclareAttackersStep
}

type EndDeclareAttackersStepEvent struct {
	EndStepEventBase
}

func (e EndDeclareAttackersStepEvent) EventType() string {
	return EventTypeEndDeclareAttackersStep
}

type BeginDeclareBlockersStepEvent struct {
	BeginStepEventBase
}

func (e BeginDeclareBlockersStepEvent) EventType() string {
	return EventTypeBeginDeclareBlockersStep
}

type EndDeclareBlockersStepEvent struct {
	EndStepEventBase
}

func (e EndDeclareBlockersStepEvent) EventType() string {
	return EventTypeEndDeclareBlockersStep
}

type BeginCombatDamageStepEvent struct {
	BeginStepEventBase
}

func (e BeginCombatDamageStepEvent) EventType() string {
	return EventTypeBeginCombatDamageStep
}

type EndCombatDamageStepEvent struct {
	EndStepEventBase
}

func (e EndCombatDamageStepEvent) EventType() string {
	return EentTypeEndCombatDamageStep
}

type BeginEndOfCombatStepEvent struct {
	BeginStepEventBase
}

func (e BeginEndOfCombatStepEvent) EventType() string {
	return EventTypeBeginEndOfCombatStep
}

type EndEndOfCombatStepEvent struct {
	EndStepEventBase
}

func (e EndEndOfCombatStepEvent) EventType() string {
	return EventTypeEndEndOfCombatStep
}

type BeginPostcombatMainStepEvent struct {
	BeginStepEventBase
}

func (e BeginPostcombatMainStepEvent) EventType() string {
	return EventTypeBeginPostcombatMainStep
}

type EndPostcombatMainStepEvent struct {
	EndStepEventBase
}

func (e EndPostcombatMainStepEvent) EventType() string {
	return EventTypeEndPostcombatMainStep
}

type BeginEndStepEvent struct {
	BeginStepEventBase
}

func (e BeginEndStepEvent) EventType() string {
	return EventTypeBeginEndStep
}

type EndEndStepEvent struct {
	EndStepEventBase
}

func (e EndEndStepEvent) EventType() string {
	return EventTypeEndEndStep
}

type BeginCleanupStepEvent struct {
	BeginStepEventBase
}

func (e BeginCleanupStepEvent) EventType() string {
	return EventTypeBeginCleanupStep
}

type EndCleanupStepEvent struct {
	EndStepEventBase
}

func (e EndCleanupStepEvent) EventType() string {
	return EventTypeEndCleanupStep
}

func NewBeginStepEvent(step mtg.Step) GameEvent {
	switch step {
	case mtg.StepUntap:
		return BeginUntapStepEvent{}
	case mtg.StepUpkeep:
		return BeginUpkeepStepEvent{}
	case mtg.StepDraw:
		return BeginDrawStepEvent{}
	case mtg.StepPrecombatMain:
		return BeginPrecombatMainStepEvent{}
	case mtg.StepBeginningOfCombat:
		return BeginBeginningOfCombatStepEvent{}
	case mtg.StepDeclareAttackers:
		return BeginDeclareAttackersStepEvent{}
	case mtg.StepDeclareBlockers:
		return BeginDeclareBlockersStepEvent{}
	case mtg.StepCombatDamage:
		return BeginCombatDamageStepEvent{}
	case mtg.StepEndOfCombat:
		return BeginEndOfCombatStepEvent{}
	case mtg.StepPostcombatMain:
		return BeginPostcombatMainStepEvent{}
	case mtg.StepEnd:
		return BeginEndStepEvent{}
	case mtg.StepCleanup:
		return BeginCleanupStepEvent{}
	default:
		panic("unknown step")
	}
}

func NewEndStepEvent(step mtg.Step) GameEvent {
	switch step {
	case mtg.StepUntap:
		return EndUntapStepEvent{}
	case mtg.StepUpkeep:
		return EndUpkeepStepEvent{}
	case mtg.StepDraw:
		return EndDrawStepEvent{}
	case mtg.StepPrecombatMain:
		return EndPrecombatMainStepEvent{}
	case mtg.StepBeginningOfCombat:
		return EndBeginningOfCombatStepEvent{}
	case mtg.StepDeclareAttackers:
		return EndDeclareAttackersStepEvent{}
	case mtg.StepDeclareBlockers:
		return EndDeclareBlockersStepEvent{}
	case mtg.StepCombatDamage:
		return EndCombatDamageStepEvent{}
	case mtg.StepEndOfCombat:
		return EndEndOfCombatStepEvent{}
	case mtg.StepPostcombatMain:
		return EndPostcombatMainStepEvent{}
	case mtg.StepEnd:
		return EndEndStepEvent{}
	case mtg.StepCleanup:
		return EndCleanupStepEvent{}
	default:
		panic("unknown step")
	}
}

// TODO Should this return a BeginPhaseEvent?
func NewBeginPhaseEvent(phase mtg.Phase) GameEvent {
	switch phase {
	case mtg.PhaseBeginning:
		return BeginBeginningPhaseEvent{}
	case mtg.PhasePrecombatMain:
		return BeginPrecombatMainPhaseEvent{}
	case mtg.PhaseCombat:
		return BeginCombatPhaseEvent{}
	case mtg.PhasePostcombatMain:
		return BeginPostcombatMainPhaseEvent{}
	case mtg.PhaseEnding:
		return BeginEndingPhaseEvent{}
	default:
		panic("unknown phase")
	}
}

func NewEndPhaseEvent(phase mtg.Phase) GameEvent {
	switch phase {
	case mtg.PhaseBeginning:
		return EndBeginningPhaseEvent{}
	case mtg.PhasePrecombatMain:
		return EndPrecombatMainPhaseEvent{}
	case mtg.PhaseCombat:
		return EndCombatPhaseEvent{}
	case mtg.PhasePostcombatMain:
		return EndPostcombatMainPhaseEvent{}
	case mtg.PhaseEnding:
		return EndEndingPhaseEvent{}
	default:
		panic("unknown phase")
	}
}

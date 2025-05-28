package game

import (
	"fmt"
	"strconv"
)

const (
	PhaseBeginning      = "Beginning"
	PhasePrecombatMain  = "PrecombatMain"
	PhaseCombat         = "Combat"
	PhasePostcombatMain = "PostcombatMain"
	PhaseEnding         = "Ending"
)

const (
	StepUntap             = "Untap"
	StepUpkeep            = "Upkeep"
	StepDraw              = "Draw"
	StepPrecombatMain     = "PrecombatMain"
	StepBeginningOfCombat = "BeginningOfCombat"
	StepDeclareAttackers  = "DeclareAttackers"
	StepDeclareBlockers   = "DeclareBlockers"
	StepCombatDamage      = "CombatDamage"
	StepEndOfCombat       = "EndOfCombat"
	StepPostcombatMain    = "PostcombatMain"
	StepEnd               = "End"
	StepCleanup           = "Cleanup"
)

type GamePhase struct {
	Name  string
	Steps []GameStep
}

type GameStep struct {
	Name    string
	Handler func(g *GameState, player *Player) error
	// TODO This name sucks
	EventEvent EventType
}

var beginningPhase = GamePhase{
	Name: PhaseBeginning,
	Steps: []GameStep{
		{
			Name:       StepUntap,
			EventEvent: EventUntapStep,
			Handler: func(g *GameState, player *Player) error {
				g.Log("Untapping all permanents")
				actionResult, err := ActionUntapFunc(g, player, UntapAll)
				if err != nil {
					return fmt.Errorf("failed to untap: %w", err)
				}
				g.Message = actionResult.Message
				return nil
			},
		},
		{
			Name:       StepUpkeep,
			EventEvent: EventUpkeepStep,
			Handler: wrapStep(func(g *GameState, player *Player) error {
				return nil
			}),
		},
		{
			Name:       StepDraw,
			EventEvent: EventDrawStep,
			Handler: wrapStep(func(g *GameState, player *Player) error {
				g.Log("Drawing a card")
				actionResult, err := ActionDrawFunc(g, player, "1")
				if err != nil {
					return fmt.Errorf("failed to draw: %w", err)
				}
				g.Message = actionResult.Message
				return nil
			}),
		},
	},
}

var precombatMainPhase = GamePhase{
	Name: PhasePrecombatMain,
	Steps: []GameStep{
		{
			Name:       StepPrecombatMain,
			EventEvent: EventPrecombatMainPhase,
			Handler: wrapStep(func(g *GameState, player *Player) error {
				return nil
			}),
		},
	},
}

var combatPhase = GamePhase{
	Name: PhaseCombat,
	Steps: []GameStep{
		{
			Name:       StepBeginningOfCombat,
			EventEvent: EventBeginningOfCombatStep,

			Handler: wrapStep(func(g *GameState, player *Player) error {
				return nil
			}),
		},
		{
			Name:       StepDeclareAttackers,
			EventEvent: EventDeclareAttackersStep,
			Handler: wrapStep(func(g *GameState, player *Player) error {
				return nil
			}),
		},
		{
			Name:       StepDeclareBlockers,
			EventEvent: EventDeclareBlockersStep,
			Handler: wrapStep(func(g *GameState, player *Player) error {
				return nil
			}),
		},
		{
			Name:       StepCombatDamage,
			EventEvent: EventCombatDamageStep,
			Handler: wrapStep(func(g *GameState, player *Player) error {
				return nil
			}),
		},
		{
			Name: StepEndOfCombat,

			EventEvent: EventEndOfCombatStep,
			Handler: wrapStep(func(g *GameState, player *Player) error {
				return nil
			}),
		},
	},
}

var postcombatMainPhase = GamePhase{
	Name: PhasePostcombatMain,
	Steps: []GameStep{
		{
			Name:       StepPostcombatMain,
			EventEvent: EventPostcombatMainPhase,
			Handler: wrapStep(func(g *GameState, player *Player) error {
				return nil
			}),
		},
	},
}

var endingPhase = GamePhase{
	Name: PhaseEnding,
	Steps: []GameStep{
		{
			Name:       StepEnd,
			EventEvent: EventEndStep,
			Handler: wrapStep(func(g *GameState, player *Player) error {
				return nil
			}),
		},
		{
			Name:       StepCleanup,
			EventEvent: EventCleanupStep,
			Handler: func(g *GameState, player *Player) error {
				toDiscard := player.Hand.Size() - player.MaxHandSize
				if toDiscard > 0 {
					g.Log(fmt.Sprintf("Discarding %d cards to maintain max hand size", toDiscard))
					actionResult, err := ActionDiscardFunc(g, player, strconv.Itoa(toDiscard))
					if err != nil {
						return fmt.Errorf("failed to discard cards: %w", err)
					}
					g.Message = actionResult.Message
					return nil
				}
				actionResult := &ActionResult{
					Message: "Cleanup step completed",
				}
				g.Message = actionResult.Message
				return nil
			},
		},
	},
}

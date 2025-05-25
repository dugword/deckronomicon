package game

import (
	"fmt"
	"strconv"
)

const (
	PhaseBeginning      = "Beginning"
	PhasePreCombatMain  = "PreCombatMain"
	PhaseCombat         = "Combat"
	PhasePostCombatMain = "PostCombatMain"
	PhaseEnding         = "Ending"
)

const (
	StepUntap             = "Untap"
	StepUpkeep            = "Upkeep"
	StepDraw              = "Draw"
	StepPreCombatMain     = "PreCombatMain"
	StepBeginningOfCombat = "BeginningOfCombat"
	StepDeclareAttackers  = "DeclareAttackers"
	StepDeclareBlockers   = "DeclareBlockers"
	StepCombatDamage      = "CombatDamage"
	StepEndOfCombat       = "EndOfCombat"
	StepPostCombatMain    = "PostCombatMain"
	StepEnd               = "End"
	StepCleanup           = "Cleanup"
)

var beginningPhase = GamePhase{
	Name: PhaseBeginning,
	Steps: []GameStep{
		{
			Name:       StepUntap,
			EventEvent: EventUntapStep,
			Handler: func(g *GameState, agent PlayerAgent) error {
				g.Log("Untapping all permanents")
				actionResult, err := ActionUntapFunc(g, UntapAll, agent)
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
			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
				return nil
			}),
		},
		{
			Name:       StepDraw,
			EventEvent: EventDrawStep,
			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
				g.Log("Drawing a card")
				actionResult, err := ActionDrawFunc(g, "1", agent)
				if err != nil {
					return fmt.Errorf("failed to draw: %w", err)
				}
				g.Message = actionResult.Message
				return nil
			}),
		},
	},
}

var preCombatMainPhase = GamePhase{
	Name: PhasePreCombatMain,
	Steps: []GameStep{
		{
			Name:       StepPreCombatMain,
			EventEvent: EventPrecombatMainPhase,
			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
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

			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
				return nil
			}),
		},
		{
			Name:       StepDeclareAttackers,
			EventEvent: EventDeclareAttackersStep,
			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
				return nil
			}),
		},
		{
			Name:       StepDeclareBlockers,
			EventEvent: EventDeclareBlockersStep,
			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
				return nil
			}),
		},
		{
			Name:       StepCombatDamage,
			EventEvent: EventCombatDamageStep,
			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
				return nil
			}),
		},
		{
			Name: StepEndOfCombat,

			EventEvent: EventEndOfCombatStep,
			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
				return nil
			}),
		},
	},
}

var postCombatMainPhase = GamePhase{
	Name: PhasePostCombatMain,
	Steps: []GameStep{
		{
			Name:       StepPostCombatMain,
			EventEvent: EventPostcombatMainPhase,
			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
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
			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
				return nil
			}),
		},
		{
			Name:       StepCleanup,
			EventEvent: EventCleanupStep,
			Handler: func(g *GameState, agent PlayerAgent) error {
				toDiscard := g.Hand.Size() - g.MaxHandSize
				if toDiscard > 0 {
					g.Log(fmt.Sprintf("Discarding %d cards to maintain max hand size", toDiscard))
					actionResult, err := ActionDiscardFunc(g, strconv.Itoa(toDiscard), agent)
					if err != nil {
						return fmt.Errorf("failed to discard cards: %w", err)
					}
					g.Message = actionResult.Message
					return nil
				}
				g.ManaPool.Empty()
				actionResult := &ActionResult{
					Message: "Cleanup step completed",
				}
				g.Message = actionResult.Message
				return nil
			},
		},
	},
}

package game

import (
	"fmt"
	"strconv"
)

const (
	PhaseBeginning      = "Beginning"
	PhasePreCombatMain  = "Pre-combat Main"
	PhaseCombat         = "Combat"
	PhasePostCombatMain = "Post-combat Main"
	PhaseEnding         = "Ending"
)

var beginningPhase = GamePhase{
	Name: "Beginning",
	Steps: []GameStep{
		{
			Name:       "Untap",
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
			Name:       "Upkeep",
			EventEvent: EventUpkeepStep,
			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
				return nil
			}),
		},
		{
			Name:       "Draw",
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
	Name: "Pre-combat Main",
	Steps: []GameStep{
		{
			Name:       "Pre-combat Main",
			EventEvent: EventPrecombatMainPhase,
			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
				return nil
			}),
		},
	},
}

var combatPhase = GamePhase{
	Name: "Combat",
	Steps: []GameStep{
		{
			Name:       "Beginning of Combat",
			EventEvent: EventBeginningOfCombatStep,

			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
				return nil
			}),
		},
		{
			Name:       "Declare Attackers",
			EventEvent: EventDeclareAttackersStep,
			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
				return nil
			}),
		},
		{
			Name:       "Declare Blockers",
			EventEvent: EventDeclareBlockersStep,
			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
				return nil
			}),
		},
		{
			Name:       "Combat Damage",
			EventEvent: EventCombatDamageStep,
			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
				return nil
			}),
		},
		{
			Name: "End of Combat",

			EventEvent: EventEndOfCombatStep,
			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
				return nil
			}),
		},
	},
}

var postCombatMainPhase = GamePhase{
	Name: "Post-combat Main",
	Steps: []GameStep{
		{
			Name:       "Post-combat Main",
			EventEvent: EventPostcombatMainPhase,
			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
				return nil
			}),
		},
	},
}

var endingPhase = GamePhase{
	Name: "Ending",
	Steps: []GameStep{
		{
			Name:       "End Step",
			EventEvent: EventEndStep,
			Handler: wrapStep(func(g *GameState, agent PlayerAgent) error {
				return nil
			}),
		},
		{
			Name:       "Cleanup",
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

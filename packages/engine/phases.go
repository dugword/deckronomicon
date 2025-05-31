package engine

import (
	"deckronomicon/packages/game/action"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/player"
	"fmt"
	"strconv"
)

// TODO: move this to game/phases load the handler in engine

type GamePhase struct {
	Name  mtg.Phase
	Steps []GameStep
}

type GameStep struct {
	Name    mtg.Step
	Handler func(state *GameState, player *player.Player) error
	// TODO This name sucks
	// // EventEvent EventType
}

var beginningPhase = GamePhase{
	Name: mtg.PhaseBeginning,
	Steps: []GameStep{
		{
			Name: mtg.StepUntap,
			// EventEvent: EventUntapStep,
			Handler: func(state *GameState, player *player.Player) error {
				state.Log("Untapping all permanents")
				actionResult, err := ActionUntapFunc(state, player, action.ActionTarget{Name: UntapAll})
				if err != nil {
					return fmt.Errorf("failed to untap: %w", err)
				}
				state.Message = actionResult.Message
				return nil
			},
		},
		{
			Name: mtg.StepUpkeep,
			// EventEvent: EventUpkeepStep,
			Handler: wrapStep(func(state *GameState, player *player.Player) error {
				return nil
			}),
		},
		{
			Name: mtg.StepDraw,
			// EventEvent: EventDrawStep,
			Handler: wrapStep(func(state *GameState, player *player.Player) error {
				state.Log("Drawing a card")
				actionResult, err := ActionDrawFunc(state, player, action.ActionTarget{Name: "1"})
				if err != nil {
					return fmt.Errorf("failed to draw: %w", err)
				}
				state.Message = actionResult.Message
				return nil
			}),
		},
	},
}

var precombatMainPhase = GamePhase{
	Name: mtg.PhasePrecombatMain,
	Steps: []GameStep{
		{
			Name: mtg.StepPrecombatMain,
			// EventEvent: EventPrecombatMainPhase,
			Handler: wrapStep(func(state *GameState, player *player.Player) error {
				return nil
			}),
		},
	},
}

var combatPhase = GamePhase{
	Name: mtg.PhaseCombat,
	Steps: []GameStep{
		{
			Name: mtg.StepBeginningOfCombat,
			// EventEvent: EventBeginningOfCombatStep,

			Handler: wrapStep(func(state *GameState, player *player.Player) error {
				return nil
			}),
		},
		{
			Name: mtg.StepDeclareAttackers,
			// EventEvent: EventDeclareAttackersStep,
			Handler: wrapStep(func(state *GameState, player *player.Player) error {
				return nil
			}),
		},
		{
			Name: mtg.StepDeclareBlockers,
			// EventEvent: EventDeclareBlockersStep,
			Handler: wrapStep(func(state *GameState, player *player.Player) error {
				return nil
			}),
		},
		{
			Name: mtg.StepCombatDamage,
			// EventEvent: EventCombatDamageStep,
			Handler: wrapStep(func(state *GameState, player *player.Player) error {
				return nil
			}),
		},
		{
			Name: mtg.StepEndOfCombat,

			// EventEvent: EventEndOfCombatStep,
			Handler: wrapStep(func(state *GameState, player *player.Player) error {
				return nil
			}),
		},
	},
}

var postcombatMainPhase = GamePhase{
	Name: mtg.PhasePostcombatMain,
	Steps: []GameStep{
		{
			Name: mtg.StepPostcombatMain,
			// EventEvent: EventPostcombatMainPhase,
			Handler: wrapStep(func(state *GameState, player *player.Player) error {
				return nil
			}),
		},
	},
}

var endingPhase = GamePhase{
	Name: mtg.PhaseEnding,
	Steps: []GameStep{
		{
			Name: mtg.StepEnd,
			// EventEvent: EventEndStep,
			Handler: wrapStep(func(state *GameState, player *player.Player) error {
				return nil
			}),
		},
		{
			Name: mtg.StepCleanup,
			// EventEvent: EventCleanupStep,
			Handler: func(state *GameState, player *player.Player) error {
				toDiscard := player.Hand.Size() - player.MaxHandSize
				if toDiscard > 0 {
					state.Log(fmt.Sprintf("Discarding %d cards to maintain max hand size", toDiscard))
					actionResult, err := ActionDiscardFunc(state, player, action.ActionTarget{Name: strconv.Itoa(toDiscard)})
					if err != nil {
						return fmt.Errorf("failed to discard cards: %w", err)
					}
					state.Message = actionResult.Message
					return nil
				}
				actionResult := &ActionResult{
					Message: "Cleanup step completed",
				}
				state.Message = actionResult.Message
				return nil
			},
		},
	},
}

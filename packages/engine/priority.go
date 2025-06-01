package engine

import (
	"deckronomicon/packages/game/action"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/player"
	"errors"
	"fmt"
)

// TODO Handle Preactions maybe?
/*
// TODO: Not sure if I like this here, it's primarily for
// automatically tapping lands for mana before casting spells.
if len(action.Preactions) > 0 {
	for _, preaction := range action.Preactions {
		g.Log("Preaction: " + string(preaction.Type))
		// Should consolidate this with ResolveAction below
		_, err := g.ResolveAction(preaction, agent)
		if err != nil {
			g.Error(err)
			g.LastActionFailed = true
			// TODO: Think through if we want to put errors here, or
			// if we want to get the error from the action result
			// struct.
			// or if errors should go to an g.ErrorMessage or
			// something
			g.Message = "Error: " + err.Error()
			// TODO: Should this break the outer loop?
			// Probably...
			continue
		}
	}
}
*/

func (e *Engine) wrapStep(handler func(*GameState, *player.Player) error) func(*GameState, *player.Player) error {
	return func(state *GameState, player *player.Player) error {
		if err := handler(state, player); err != nil {
			return err
		}
		if err := e.HandlePriority(player); err != nil {
			return err
		}
		return nil
	}
}

func (e *Engine) RunPriorityLoop(player *player.Player) error {
	if !player.HasStop(e.GameState.CurrentStep) {
		player.Agent.ReportState(e.GameState)
		return nil
	}
	for {
		// player.PotentialMana = GetPotentialMana(g, player)
		player.Agent.ReportState(e.GameState)
		// TODO: Better package name so this doesn't conflict
		act, err := player.Agent.GetNextAction(e.GameState)
		if err != nil {
			return fmt.Errorf(
				"failed to get next action for %s: %w",
				player.ID,
				err,
			)
		}
		if !action.PlayerActions[act.Type] {
			if !e.GameState.CheatsEnabled {
				e.GameState.Message = fmt.Sprintf("Invalid player action: %s", act.Type)
				continue
			}
		}
		result, err := e.ResolveAction(act, player)
		if err != nil {
			if errors.Is(err, mtg.ErrGameOver) {
				return err
			}
			errString := fmt.Sprintf("ERROR: %s", err)
			e.Log(errString)
			e.GameState.Message = errString
			continue
		}
		e.GameState.Message = result.Message
		if result.Pass {
			return nil
		}
	}
}

func (e *Engine) HandlePriority(player *player.Player) error {
	for {
		if err := e.RunPriorityLoop(player); err != nil {
			return err
		}
		if e.GameState.Stack.Size() == 0 {
			break
		}
		e.Log("Resolving stack...")
		object, err := e.GameState.Stack.Pop()
		if err != nil {
			return fmt.Errorf("failed to pop resolvable from stack: %w", err)
		}
		// Resolve should live in this pacakge: engine
		resolvable, ok := object.(Resolvable)
		if !ok {
			return fmt.Errorf("object is not resolvable: %T\n", object)
		}
		e.Log("Resolving: " + resolvable.Name())
		if err := resolvable.Resolve(e.GameState, player); err != nil {
			return fmt.Errorf("failed to resolve: %w", err)
		}
		e.Log("Rsolved: " + resolvable.Name())
	}
	return nil
}

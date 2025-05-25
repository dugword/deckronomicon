package game

import (
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

func wrapStep(handler func(*GameState, *Player) error) func(*GameState, *Player) error {
	return func(g *GameState, player *Player) error {
		if err := handler(g, player); err != nil {
			return err
		}
		if err := g.HandlePriority(player); err != nil {
			return err
		}
		return nil
	}
}

func (g *GameState) RunPriorityLoop(player *Player) error {
	if player.ShouldAutoPass(g.CurrentStep) {
		player.Agent.ReportState(g)
		return nil
	}
	for {
		player.PotentialMana = GetPotentialMana(g, player)
		player.Agent.ReportState(g)
		action := player.Agent.GetNextAction(g)
		if !PlayerActions[action.Type] {
			if !g.Cheat {
				g.Message = fmt.Sprintf("Invalid player action: %s", action.Type)
				continue
			}
		}
		if action.Type == ActionPass {
			g.Log("Player passed priority")
			return nil
		}
		result, err := g.ResolveAction(action, player)
		if err != nil {
			if errors.Is(err, ErrGameOver) {
				return err
			}
			errString := fmt.Sprintf("ERROR: %s", err)
			g.Log(errString)
			g.Message = errString
			continue
		}
		g.Message = result.Message
	}
}

func (g *GameState) HandlePriority(player *Player) error {
	for {
		if err := g.RunPriorityLoop(player); err != nil {
			return err
		}
		if g.Stack.Size() == 0 {
			break
		}
		g.Log("Resolving stack...")
		object, err := g.Stack.Pop()
		if err != nil {
			return fmt.Errorf("failed to pop spell from stack: %w", err)
		}
		spell, ok := object.(*Spell)
		if !ok {
			return fmt.Errorf("object is not a spell: %T\n", object)
		}
		g.Log("Resolving spell: " + spell.Name())
		if err := spell.Resolve(g, player); err != nil {
			return fmt.Errorf("failed to resolve spell: %w", err)
		}
		g.Log("Spell resolved: " + spell.Name())
	}
	return nil
}

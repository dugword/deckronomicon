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

func wrapStep(handler func(*GameState, PlayerAgent) error) func(*GameState, PlayerAgent) error {
	return func(g *GameState, agent PlayerAgent) error {
		if err := handler(g, agent); err != nil {
			return err
		}
		if err := g.HandlePriority(agent); err != nil {
			return err
		}
		return nil
	}
}

func (g *GameState) RunPriorityLoop(agent PlayerAgent) error {
	for {
		g.PotentialMana = GetPotentialMana(g)
		agent.ReportState(g)
		action := agent.GetNextAction(g)
		if !PlayerActions[action.Type] {
			if !g.Cheat {
				return errors.New("invalid player action: " + string(action.Type))
			}
		}
		if action.Type == ActionPass {
			g.Log("Player passed priority")
			return nil
		}
		result, err := g.ResolveAction(action, agent)
		if err != nil {
			if errors.Is(err, ErrGameOver) {
				return err
			}
			errString := fmt.Sprintf("ERROR: %s" + err.Error())
			g.Log(errString)
			g.Message = errString
			continue
		}
		g.Message = result.Message
	}
}

func (g *GameState) HandlePriority(agent PlayerAgent) error {
	for {
		if err := g.RunPriorityLoop(agent); err != nil {
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
		if err := spell.Resolve(g, agent); err != nil {
			return fmt.Errorf("failed to resolve spell: %w", err)
		}
		g.Log("Spell resolved: " + spell.Name())
	}
	return nil
}

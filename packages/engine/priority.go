package engine

import (
	"deckronomicon/packages/game/player"
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

func wrapStep(handler func(*GameState, *player.Player) error) func(*GameState, *player.Player) error {
	return func(g *GameState, player *player.Player) error {
		if err := handler(g, player); err != nil {
			return err
		}
		if err := g.HandlePriority(player); err != nil {
			return err
		}
		return nil
	}
}

func (g *GameState) RunPriorityLoop(player *player.Player) error {
	if !player.HasStop(g.CurrentStep) {
		// player.Agent.ReportState(g)
		return nil
	}
	for {
		// player.PotentialMana = GetPotentialMana(g, player)
		// player.Agent.ReportState(g)
		/*
			action, err := player.Agent.GetNextAction(g)
			if err != nil {
				return fmt.Errorf("failed to get next action from agent: %w", err)
			}
			if !PlayerActions[action.Type] {
				if !g.Cheat {
					g.Message = fmt.Sprintf("Invalid player action: %s", action.Type)
					continue
				}
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
			if result.Pass {
				return nil
			}
		*/
	}
}

func (g *GameState) HandlePriority(player *player.Player) error {
	for {
		/*
			if err := g.RunPriorityLoop(player); err != nil {
				return err
			}
			if g.Stack.Size() == 0 {
				break
			}
			g.Log("Resolving stack...")
			object, err := g.Stack.Pop()
			if err != nil {
				return fmt.Errorf("failed to pop resolvable from stack: %w", err)
			}
			// Resolve should live in this pacakge: engine
			resolvable, ok := object.(Resolvable)
			if !ok {
				return fmt.Errorf("object is not resolvable: %T\n", object)
			}
			g.Log("Resolving: " + resolvable.Name())
			if err := resolvable.Resolve(g, player); err != nil {
				return fmt.Errorf("failed to resolve: %w", err)
			}
			g.Log("Rsolved: " + resolvable.Name())
		*/
	}
	return nil
}

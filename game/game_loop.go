package game

import (
	"fmt"
)

// RunGameLoop drives the turn-based game for the given agent
func (g *GameState) RunGameLoop(agent PlayerAgent, maxTurns int) {
	for !g.GameLost && g.Turn < maxTurns {
		g.RunTurn(agent)
		if g.GameLost {
			return
		}
		for len(g.Hand) > 7 {
			ActionDiscardFunc(g, agent)
		}
	}
}

// RunTurn executes a single game turn using the provided agent
func (g *GameState) RunTurn(agent PlayerAgent) {
	g.Turn++
	g.ManaPool = make(ManaPool)
	g.LandDrop = false
	UntapPermanents(g.Battlefield)
	g.DrawCards(1)
	g.CurrentPhase = "First Main"
	g.TurnMessageLog = []string{}
	for {
		g.PotentialMana = GetPotentialMana(g)
		agent.ReportState(g)
		action := agent.GetNextAction(g)
		if action.Cheat != "" && g.Cheat {
			result, err := g.ResolveCheat(action, agent)
			if err != nil {
				g.Error(err)
			}
			// TODO: Hand printing somewhere more cohesive
			if result != nil {
				fmt.Println("CHEAT RESULT:", result.Message)
			}
		}
		result, err := g.ResolveAction(action, agent)
		if result != nil {
			g.Message = result.Message
		} else {
			g.Message = ""
		}
		if err != nil {
			g.Error(err)
			fmt.Println("ERROR: " + err.Error())
			g.LastActionFailed = true
		} else {
			g.LastActionFailed = false
		}
		if result != nil && result.EndTurn {
			return
		}
	}
}

func (g *GameState) ResolveCheat(action GameAction, agent PlayerAgent) (*ActionResult, error) {
	switch action.Cheat {
	case CheatDraw:
		g.DrawCards(1)
	case CheatPeek:
		return &ActionResult{
			Message: "Top Card: " + g.Deck[0].Name,
		}, nil
	case CheatShuffle:
		g.ShuffleDeck()
	case CheatPrintDeck:
		g.PrintDeck()
	case CheatDiscard:
		return ActionDiscardFunc(g, agent)
	default:
		g.Log("Unknown cheat: " + string(action.Type))
	}
	return nil, nil
}

// TODO: Rethink this error handling, are invalid actions errors?
func (g *GameState) ResolveAction(action GameAction, agent PlayerAgent) (result *ActionResult, err error) {
	switch action.Type {
	case ActionCheat:
		g.Cheat = true
		g.Log("Cheats enabled... you cheater")
	case ActionPlay:
		// todo: make this nicer
		g.Log("Action: play " + action.Target)
		return ActionPlayFunc(g, action.Target, agent)
	case ActionActivate:
		g.Log("Action: activate")
		return ActionActivateFunc(g, action.Target, agent)
	case ActionPass:
		g.Log("Action: pass")
		return &ActionResult{EndTurn: true}, nil
	case ActionConcede:
		g.Log("Action: concede")
		g.GameLost = true
		return &ActionResult{EndTurn: true}, nil
	case ActionView:
		g.Log("Action: view")
		return ActionViewFunc(g, agent)
	default:
		g.Log("Unknown action: " + string(action.Type))
	}
	return nil, nil
}

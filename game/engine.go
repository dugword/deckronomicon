package game

import (
	"fmt"
)

// GameAction represents an action a player can take
type GameAction struct {
	Type       GameActionType
	Cheat      GameCheatType
	Target     string
	Preactions []GameAction
}

// RunGameLoop drives the turn-based game for the given agent
func (g *GameState) RunGameLoop(agent PlayerAgent) error {
	for {
		if g.Turn > g.MaxTurns {
			return ErrMaxTurnsExceeded
		}
		if err := g.RunTurn(agent); err != nil {
			return err
		}
	}
}

/*
type Step func(state *GameState, agent PlayerAgent, stepName string) error

func ResolvePhase(state *GameState, agent PlayerAgent, steps []*Step) error {

}

func ResolveStep(state *GameState, agent PlayerAgent, stepName string) error {

}
*/

// Maybe split into ResolvePhase and ResolveStep functions
// RunTurn executes a single game turn using the provided agent
func (g *GameState) RunTurn(agent PlayerAgent) error {
	// Beginning Phase
	g.Turn++
	g.LandDrop = false
	g.Battlefield.RemoveSummoningSickness()
	g.TurnMessageLog = []string{} // TODO: this sucks, make better
	g.CurrentPhase = "Beginning"
	// Untap Step
	g.CurrentStep = "Untap"
	actionResult, err := g.ResolveAction(GameAction{Type: ActionUntap}, agent)
	if err != nil {
		return err
	}
	g.Log(actionResult.Message)

	// Upkeep Step
	g.CurrentStep = "Upkeep"

	// Draw Step
	g.CurrentStep = "Draw"
	actionResult, err = g.ResolveAction(GameAction{Type: ActionDraw}, agent)
	if err != nil {
		return err
	}
	g.Log(actionResult.Message)

	// Pre-combat Main Phase
	g.CurrentPhase = "Pre-combat Main"
	g.CurrentStep = ""
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
		if len(action.Preactions) > 0 {
			fmt.Println("Doing preaction stuff, make this more inline with actions")
			for _, preaction := range action.Preactions {
				fmt.Println("Preaction =>", preaction.Type)
				_, err := g.ResolveAction(preaction, agent)
				// TODO: need to handle messages and last actions and results
				if err != nil {
					g.Error(err)
					fmt.Println("ERROR: " + err.Error())
					g.LastActionFailed = true
				}
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
		if result != nil && result.Pass {
			break
		}
	}

	// Combat Phase
	g.CurrentPhase = "Combat Phase"
	g.CurrentStep = "Beginning of Combat"
	g.CurrentStep = "Declare Attackers"
	g.CurrentStep = "Declare Blockers"
	g.CurrentStep = "Combat Damage"
	g.CurrentStep = "End of Combat"

	// Post-combat Main Phase
	g.CurrentPhase = "Post-combat Main Phase"
	g.CurrentStep = ""

	// Ending Phase
	g.CurrentPhase = "Ending"
	g.CurrentStep = "End"
	g.CurrentStep = "Cleanup"
	for g.Hand.Size() > g.MaxHandSize {
		ActionDiscardFunc(g, agent)
	}
	g.ManaPool = make(ManaPool)

	return nil
}

func (g *GameState) ResolveCheat(action GameAction, agent PlayerAgent) (*ActionResult, error) {
	switch action.Cheat {
	case CheatDraw:
		g.Log("CHEAT! Action: draw")
		return ActionDrawFunc(g, agent)
	case CheatPeek:
		g.Log("CHEAT! Action: peek")
		return &ActionResult{
			Message: "Top Card: " + g.Deck.cards[0].Name(),
		}, nil
	case CheatShuffle:
		g.Log("CHEAT! Action: shuffle")
		g.Deck.Shuffle()
	case CheatDiscard:
		g.Log("CHEAT! Action: discard")
		return ActionDiscardFunc(g, agent)
	default:
		g.Log("Unknown cheat: " + string(action.Type))
	}
	return nil, nil
}

// TODO: Rethink this error handling, are invalid actions errors?
func (g *GameState) ResolveAction(action GameAction, agent PlayerAgent) (result *ActionResult, err error) {
	switch action.Type {
	case ActionActivate:
		g.Log("Action: activate")
		return ActionActivateFunc(g, action.Target, agent)
	case ActionDraw:
		g.Log("Action: draw")
		return ActionDrawFunc(g, agent)
	case ActionCheat:
		g.Log("Action: cheat... you cheater")
		g.Cheat = true
	case ActionConcede:
		g.Log("Action: concede")
		return &ActionResult{Pass: true}, PlayerLostError{Reason: Conceded}
	case ActionPass:
		g.Log("Action: pass")
		return &ActionResult{Pass: true}, nil
	case ActionPlay:
		// todo: make this nicer
		g.Log("Action: play " + action.Target)
		return ActionPlayFunc(g, action.Target, agent)
	case ActionUntap:
		g.Log("Action: untap")
		return ActionUntapFunc(g, agent)
	case ActionView:
		g.Log("Action: view")
		return ActionViewFunc(g, agent)
	default:
		g.Log("Unknown action: " + string(action.Type))
	}
	return nil, nil
}

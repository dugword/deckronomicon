package game

import (
	"errors"
	"fmt"
	"strconv"
)

// RunGameLoop drives the turn-based game for the given agent.
func (g *GameState) RunGameLoop(agent PlayerAgent) error {
	if err := g.Mulligan(agent); err != nil {
		return fmt.Errorf("failed to mulligan: %w", err)
	}
	for {
		if g.Turn > g.MaxTurns {
			return ErrMaxTurnsExceeded
		}
		if err := g.RunTurn(agent); err != nil {
			return err
		}
	}
}

var ChoiceMulligan = "Mulligan"
var ChoiceSourceMulligan = NewChoiceSource(ChoiceMulligan, ChoiceMulligan)

// Mulligan allows the player to mulligan their hand. The player draws 7 new
// cards but then needs to put 1 card back on the bottom of their library for
// each time they've mulliganed.
func (g *GameState) Mulligan(agent PlayerAgent) error {
	var accept bool
	var err error
	for (g.Mulligans < g.MaxHandSize) || !accept {
		agent.ReportState(g)
		accept, err = agent.Confirm("Keep Hand? (y/n)", ChoiceSourceMulligan)
		if err != nil {
			return fmt.Errorf("failed to confirm mulligan: %w", err)
		}
		if accept {
			break
		}
		g.Log("Mulliganing...")
		for _, card := range g.Hand.GetAll() {
			g.Hand.Remove(card.ID())
			g.Library.Add(card)
		}
		g.Library.Shuffle()
		g.DrawStartingHand(g.StartingHand)
		g.Mulligans++
	}
	if g.Mulligans != 0 {
		if err := PutNBackOnTop(g, g.Mulligans, NewChoiceSource(ChoiceMulligan, ChoiceMulligan), agent); err != nil {
			return fmt.Errorf("failed to put cards back on top: %w", err)
		}
	}
	return nil
}

/*
type Step func(state *GameState, agent PlayerAgent, stepName string) error

func ResolvePhase(state *GameState, agent PlayerAgent, steps []*Step) error {

}

func ResolveStep(state *GameState, agent PlayerAgent, stepName string) error {

}
*/

// RunTurn executes a single game turn using the provided agent.
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
	// agent.ReportState(g)
	g.EmitEvent(Event{Type: EventUntapStep}, agent)
	actionResult, err := g.ResolveAction(
		GameAction{
			Type:   ActionUntap,
			Target: UntapAll,
		}, agent)
	if err != nil {
		return err
	}
	g.Log(actionResult.Message)

	// Upkeep Step
	g.CurrentStep = "Upkeep"
	// agent.ReportState(g)
	g.EmitEvent(Event{Type: EventUpkeepStep}, agent)

	// Draw Step
	g.CurrentStep = "Draw"
	// agent.ReportState(g)
	g.EmitEvent(Event{Type: EventDrawStep}, agent)
	actionResult, err = g.ResolveAction(GameAction{
		Type:   ActionDraw,
		Target: "1",
	}, agent)
	if err != nil {
		return err
	}
	g.Log(actionResult.Message)

	// Pre-combat Main Phase
	g.CurrentPhase = "Pre-combat Main"
	g.EmitEvent(Event{Type: EventPrecombatMainPhase}, agent)
	g.CurrentStep = "Pre-combat Main"
	for {
		g.PotentialMana = GetPotentialMana(g)
		agent.ReportState(g)
		action := agent.GetNextAction(g)
		// TODO Revisit this, it's a bit messy
		if !g.Cheat {
			if !PlayerActions[action.Type] {
				g.Log("Invalid player action: " + string(action.Type))
				continue
			}
		}
		if action.Cheat != "" && g.Cheat {
			result, err := g.ResolveCheat(action, agent)
			if err != nil {
				g.Error(err)
				g.Message = "Error: " + err.Error()
				continue
			}
			// TODO: Hand printing somewhere more cohesive
			if result != nil {
				g.Log("CHEAT RESULT:", result.Message)
				g.Message = result.Message
			}
			continue
		}
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
		result, err := g.ResolveAction(action, agent)
		if err != nil {
			if errors.Is(err, ErrGameOver) {
				g.Log("Game Over: " + err.Error())
				return err
			}
			if result != nil {
				g.Message = result.Message
			} else {
				g.Message = "Error: " + err.Error()
			}
			g.Error(err)
			g.LastActionFailed = true
			continue
		}
		g.LastActionFailed = false
		// TODO: Maybe always pass back an ActionResult?
		if result != nil {
			g.Message = result.Message
			if result.Pass {
				break
			}
		}
	}

	// Combat Phase
	g.CurrentPhase = "Combat Phase"
	g.CurrentStep = "Beginning of Combat"
	// agent.ReportState(g)
	g.EmitEvent(Event{Type: EventBeginningOfCombatStep}, agent)
	g.CurrentStep = "Declare Attackers"
	// agent.ReportState(g)
	g.EmitEvent(Event{Type: EventDeclareAttackersStep}, agent)
	g.CurrentStep = "Declare Blockers"
	// agent.ReportState(g)
	g.EmitEvent(Event{Type: EventDeclareBlockersStep}, agent)
	g.CurrentStep = "Combat Damage"
	// agent.ReportState(g)
	g.EmitEvent(Event{Type: EventCombatDamageStep}, agent)
	g.CurrentStep = "End of Combat"
	// agent.ReportState(g)
	g.EmitEvent(Event{Type: EventEndOfCombatStep}, agent)

	// Post-combat Main Phase
	g.CurrentPhase = "Post-combat Main Phase"
	g.EmitEvent(Event{Type: EventPostcombatMainPhase}, agent)
	g.CurrentStep = "Post-combat Main Phase"
	//agent.ReportState(g)

	// Ending Phase
	g.CurrentPhase = "Ending"
	g.CurrentStep = "End"
	// agent.ReportState(g)
	g.EmitEvent(Event{Type: EventEndStep}, agent)
	g.CurrentStep = "Cleanup"
	// agent.ReportState(g)
	toDiscard := g.Hand.Size() - g.MaxHandSize
	if toDiscard > 0 {
		ActionDiscardFunc(g, strconv.Itoa(toDiscard), agent)
	}
	g.ManaPool.Empty()

	return nil
}

// ResolveCheat handles the resolution of cheat actions.
func (g *GameState) ResolveCheat(action GameAction, agent PlayerAgent) (*ActionResult, error) {
	switch action.Cheat {
	case CheatAddMana:
		g.Log("CHEAT! Action: add mana")
		return ActionAddManaFunc(g, action.Target, agent)
	case CheatConjure:
		g.Log("Action: conjure")
		return ActionConjureFunc(g, action.Target, agent)
	case CheatDraw:
		g.Log("CHEAT! Action: draw")
		return ActionDrawFunc(g, action.Target, agent)
	case CheatFind:
		g.Log("CHEAT! Action: find")
		return ActionFindFunc(g, action.Target, agent)
	case CheatLandDrop:
		g.Log("CHEAT! Action: land drop")
		return ActionLandDropFunc(g, action.Target, agent)
	case CheatPeek:
		g.Log("CHEAT! Action: peek")
		return &ActionResult{
			// TODO: No .cards access
			Message: "Top Card: " + g.Library.Peek().Name(),
		}, nil
	case CheatShuffle:
		g.Log("CHEAT! Action: shuffle")
		g.Library.Shuffle()
	case CheatDiscard:
		g.Log("CHEAT! Action: discard")
		return ActionDiscardFunc(g, "1", agent)
	default:
		g.Log("Unknown cheat: " + string(action.Type))
	}
	return nil, nil
}

// ResolveAction handles the resolution of game actions.
func (g *GameState) ResolveAction(action GameAction, agent PlayerAgent) (result *ActionResult, err error) {
	switch action.Type {
	case ActionActivate:
		g.Log("Action: activate")
		return ActionActivateFunc(g, action.Target, agent)
	case ActionDraw:
		g.Log("Action: draw")
		return ActionDrawFunc(g, action.Target, agent)
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
		return ActionUntapFunc(g, action.Target, agent)
	case ActionView:
		g.Log("Action: view")
		return ActionViewFunc(g, action.Target, agent)
	default:
		g.Log("Unknown action: " + string(action.Type))
	}
	return nil, nil
}

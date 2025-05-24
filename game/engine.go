package game

import (
	"fmt"
)

type GamePhase struct {
	Name  string
	Steps []GameStep
}

type GameStep struct {
	Name string
	// TODO: Rename action, it's overloaded as a term
	Handler func(g *GameState, agent PlayerAgent) error
	// TODO This name sucks
	EventEvent EventType
}

func (g *GameState) RunStep(step GameStep, agent PlayerAgent) error {
	g.CurrentPhase = step.Name
	g.Log("Running step: %s", step.Name)
	g.EmitEvent(Event{Type: step.EventEvent}, agent)
	if err := step.Handler(g, agent); err != nil {
		return fmt.Errorf("failed to run step %s: %w", step.Name, err)
	}
	g.Log("Step completed: %s", step.Name)
	return nil
}

func (g *GameState) RunPhase(phase GamePhase, agent PlayerAgent) error {
	g.CurrentPhase = phase.Name
	g.Log("Running phase: %s", phase.Name)
	for _, step := range phase.Steps {
		if err := g.RunStep(step, agent); err != nil {
			return fmt.Errorf("failed to run phase %s step %s: %w", phase.Name, step.Name, err)
		}
	}
	g.Log("Phase completed: %s", phase.Name)
	return nil
}

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

// RunTurn executes a single game turn using the provided agent.
// Maybe split into ResolvePhase and ResolveStep functions
// RunTurn executes a single game turn using the provided agent
func (g *GameState) RunTurn(agent PlayerAgent) error {
	// Beginning Phase
	g.Turn++
	g.LandDrop = false
	g.Battlefield.RemoveSummoningSickness()
	g.TurnMessageLog = []string{} // TODO: this sucks, make better
	g.Log(fmt.Sprintf("Turn %d", g.Turn))

	if err := g.RunPhase(beginningPhase, agent); err != nil {
		return fmt.Errorf("failed to run beginning phase: %w", err)
	}

	if err := g.RunPhase(preCombatMainPhase, agent); err != nil {
		return fmt.Errorf("failed to run pre-combat main phase: %w", err)
	}

	if err := g.RunPhase(combatPhase, agent); err != nil {
		return fmt.Errorf("failed to run combat phase: %w", err)
	}

	if err := g.RunPhase(postCombatMainPhase, agent); err != nil {
		return fmt.Errorf("failed to run post-combat main phase: %w", err)
	}

	if err := g.RunPhase(endingPhase, agent); err != nil {
		return fmt.Errorf("failed to run ending phase: %w", err)
	}

	// Pre-combat Main Phase
	g.CurrentPhase = "Pre-combat Main"
	g.EmitEvent(Event{Type: EventPrecombatMainPhase}, agent)
	g.CurrentStep = "Pre-combat Main"

	return nil
}

// ResolveAction handles the resolution of game actions.
func (g *GameState) ResolveAction(action *GameAction, agent PlayerAgent) (result *ActionResult, err error) {
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

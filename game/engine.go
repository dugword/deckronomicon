package game

import (
	"fmt"
)

func (g *GameState) RunStep(step GameStep, player *Player) error {
	g.CurrentStep = step.Name
	g.Log(fmt.Sprintf("Running step: %s", step.Name))
	g.EmitEvent(Event{Type: step.EventEvent}, player)
	if err := step.Handler(g, player); err != nil {
		return fmt.Errorf("failed to run step %s: %w", step.Name, err)
	}
	player.ManaPool.Empty()
	g.Log(fmt.Sprintf("Step completed: %s", step.Name))
	return nil
}

func (g *GameState) RunPhase(phase GamePhase, player *Player) error {
	g.CurrentPhase = phase.Name
	g.Log(fmt.Sprintf("Running phase: %s", phase.Name))
	for _, step := range phase.Steps {
		if err := g.RunStep(step, player); err != nil {
			return fmt.Errorf("failed to run phase %s step %s: %w", phase.Name, step.Name, err)
		}
	}
	g.Log(fmt.Sprintf("Phase completed: %s", phase.Name))
	return nil
}

// RunGameLoop drives the turn-based game.
func (g *GameState) RunGameLoop() error {
	if err := g.Mulligan(g.ActivePlayer); err != nil {
		return fmt.Errorf("active player failed to mulligan: %w", err)
	}
	if err := g.Mulligan(g.NonActivePlayer); err != nil {
		return fmt.Errorf("non-active player failed to mulligan: %w", err)
	}
	for {
		g.ActivePlayer, g.NonActivePlayer = g.NonActivePlayer, g.ActivePlayer
		if g.ActivePlayer.Turn > g.MaxTurns {
			return ErrMaxTurnsExceeded
		}
		if err := g.RunTurn(g.ActivePlayer); err != nil {
			return err
		}
	}
}

var ChoiceMulligan = "Mulligan"
var ChoiceSourceMulligan = NewChoiceSource(ChoiceMulligan, ChoiceMulligan)

// Mulligan allows the player to mulligan their hand. The player draws 7 new
// cards but then needs to put 1 card back on the bottom of their library for
// each time they've mulliganed.
func (g *GameState) Mulligan(player *Player) error {
	var accept bool
	for (player.Mulligans < player.MaxHandSize) || !accept {
		player.Agent.ReportState(g)
		accept, err := player.Agent.Confirm("Keep Hand? (y/n)", ChoiceSourceMulligan)
		if err != nil {
			return fmt.Errorf("failed to confirm mulligan: %w", err)
		}
		if accept {
			break
		}
		g.Log("Mulliganing...")
		for _, card := range player.Hand.GetAll() {
			player.Hand.Remove(card.ID())
			player.Library.Add(card)
		}
		player.Library.Shuffle()
		g.DrawStartingHand(player)
		player.Mulligans++
	}
	if player.Mulligans != 0 {
		if err := PutNBackOnTop(g, player.Mulligans, NewChoiceSource(ChoiceMulligan, ChoiceMulligan), player); err != nil {
			return fmt.Errorf("failed to put cards back on top: %w", err)
		}
	}
	return nil
}

// RunTurn executes a single game turn using the provided player.
// Maybe split into ResolvePhase and ResolveStep functions
// RunTurn executes a single game turn using the provided player
func (g *GameState) RunTurn(player *Player) error {
	// Beginning Phase
	player.Turn++
	player.LandDrop = false
	player.Battlefield.RemoveSummoningSickness()
	g.TurnMessageLog = []string{} // TODO: this sucks, make better
	g.Log(fmt.Sprintf("Turn %d", player.Turn))
	if err := g.RunPhase(beginningPhase, player); err != nil {
		return fmt.Errorf("failed to run beginning phase: %w", err)
	}
	if err := g.RunPhase(precombatMainPhase, player); err != nil {
		return fmt.Errorf("failed to run precombat main phase: %w", err)
	}
	if err := g.RunPhase(combatPhase, player); err != nil {
		return fmt.Errorf("failed to run combat phase: %w", err)
	}
	if err := g.RunPhase(postcombatMainPhase, player); err != nil {
		return fmt.Errorf("failed to run postcombat main phase: %w", err)
	}
	if err := g.RunPhase(endingPhase, player); err != nil {
		return fmt.Errorf("failed to run ending phase: %w", err)
	}
	return nil
}

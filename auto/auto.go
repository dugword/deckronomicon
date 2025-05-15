// Package auto provides a simple automated agent for playing Magic: The Gathering decks
// without requiring a user-defined strategy. It is primarily used for simulations where
// statistical sampling of deck behavior is more important than specific in-game decision logic.
//
// The auto agent operates using basic heuristics and randomized decision-making, allowing it to
// play a full game turn-by-turn with minimal configuration. It is designed to be plug-and-play:
// you can load a deck, run the simulation, and observe how the deck performs over many games.
//
// Core Behavior:
//
//   - Automatically makes a land drop each turn, if one is available.
//   - Attempts to cast spells and activate abilities based on simple prioritization rules.
//   - Selects valid targets and choices at random, unless otherwise constrained.
//   - Plays out entire turns until no further actions are detected or until a pass condition is met.
//
// Modes:
//
// The auto agent supports multiple operational modes:
//
//   - Random: Purely random decision-making across all valid options.
//   - Standard: Applies default heuristics such as always making a land drop, prioritizing card draw,
//     and sequencing known best-practice plays.
//
// Unlike the rules agent, the auto package does not require a strategy file or rule tree.
// It is intended to provide a baseline for performance testing, deck validation, and behavior sampling.
//
// Use the auto package when you want to:
//
//   - Quickly test how a deck plays without tuning a strategy
//   - Run large-scale simulations with minimal setup
//   - Benchmark new decks or mechanics under a baseline autopilot

package auto

import (
	game "deckronomicon/game"
)

// AutoAgent implements PlayerAgent for automatic play
type AutoPlayerAgent struct{}

func NewAutoPlayerAgent() *AutoPlayerAgent {
	return &AutoPlayerAgent{}
}

func (a *AutoPlayerAgent) ReportState(state *game.GameState) {

}

func (a *AutoPlayerAgent) GetNextAction(state *game.GameState) game.GameAction {
	if state.Hand.Size() > 0 {
		// Just play the first card if any
		return game.GameAction{
			Type: game.ActionPlay,
		}
	}
	return game.GameAction{Type: game.ActionPass}
}

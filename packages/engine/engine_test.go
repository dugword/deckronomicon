package engine

import (
	"testing"
)

func TestEngine(t *testing.T) {
	/*
		// TODO: Need to implement a way to pass in or configure
		// phases and steps so this doesn't decend through all
		// the subsequent phases and steps.
		t.Run("with RunTurn", func(t *testing.T) {
			player1ID := "Player 1"
			player2ID := "Player 2"
			engine := NewEngine(EngineConfig{
				Players: []*state.Player{
					state.LoadPlayerFromConfig(statetest.PlayerConfig{ID: player1ID}),
					state.LoadPlayerFromConfig(statetest.PlayerConfig{ID: player2ID}),
				},
				Agents: map[string]PlayerAgent{
					player1ID: &mockPlayerAgent{playerID: player1ID},
					player2ID: &mockPlayerAgent{playerID: player2ID},
				},
				Seed: 13,
				Log:  &mockLogger{},
			})
			engine.game = state.LoadGameFromConfig(statetest.GameConfig{
				ActivePlayerIdx: 0,
				Players: []statetest.PlayerConfig{
					{ID: player1ID},
					{ID: player2ID},
				},
			})
			if err := engine.RunTurn(); err != nil {
				t.Fatalf("RunTurn(); err = %v; want %v", err, nil)
			}
		})
	*/
}

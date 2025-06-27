package state

import (
	"deckronomicon/packages/game/definition"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPriorityPlayerID(t *testing.T) {
	const (
		player1 = "Player 1"
		player2 = "Player 2"
		player3 = "Player 3"
		player4 = "Player 4"
	)
	testCases := []struct {
		name                  string
		activePlayerID        string
		playersPassedPriority []string
		want                  string
	}{
		{
			name:           "with active player 1, no players passed priority",
			activePlayerID: player1,
			want:           player1,
		},
		{
			name:           "with active player 2, no players passed priority",
			activePlayerID: player2,
			want:           player2,
		},
		{
			name:           "with active player 3, no players passed priority",
			activePlayerID: player3,
			want:           player3,
		},
		{
			name:           "with active player 4, no players passed priority",
			activePlayerID: player4,
			want:           player4,
		},
		{
			name:                  "with active player 1, player 1 passed priority",
			activePlayerID:        player1,
			playersPassedPriority: []string{player1},
			want:                  player2,
		},
		{
			name:                  "with active player 2, player 2 passed priority",
			activePlayerID:        player2,
			playersPassedPriority: []string{player2},
			want:                  player3,
		},
		{
			name:                  "with active player 3, player 3 passed priority",
			activePlayerID:        player3,
			playersPassedPriority: []string{player3},
			want:                  player4,
		},
		{
			name:                  "with active player 4, player 4 passed priority",
			activePlayerID:        player4,
			playersPassedPriority: []string{player4},
			want:                  player1,
		},
		{
			name:                  "with active player 1, players 1 and 2 passed priority",
			activePlayerID:        player1,
			playersPassedPriority: []string{player1, player2},
			want:                  player3,
		},
		{
			name:                  "with active player 1, players 1, 2 and 3 passed priority",
			activePlayerID:        player1,
			playersPassedPriority: []string{player1, player2, player3},
			want:                  player4,
		},
		{
			name:                  "with active player 1, all players passed priority",
			activePlayerID:        player1,
			playersPassedPriority: []string{player1, player2, player3, player4},
			want:                  player1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			playersPassedPriority := map[string]bool{}
			for _, playerID := range tc.playersPassedPriority {
				playersPassedPriority[playerID] = true
			}
			game := NewGameFromDefinition(definition.Game{
				Players: []definition.Player{
					{ID: player1},
					{ID: player2},
					{ID: player3},
					{ID: player4},
				},
				PlayersPassedPriority: playersPassedPriority,
			})
			game = game.WithActivePlayer(tc.activePlayerID)
			got := game.PriorityPlayerID()
			if diff := cmp.Diff(tc.want, got, AllowAllUnexported); diff != "" {
				t.Errorf("PriorityPlayerID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

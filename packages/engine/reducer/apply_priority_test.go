package reducer

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"deckronomicon/packages/state/statetest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestApplyPriorityEvent(t *testing.T) {
	const playerID = "Test Player"
	testCases := []struct {
		name string
		evnt event.PriorityEvent
		game state.Game
		want state.Game
	}{
		{
			name: "with AllPlayersPassedPriorityEvent",
			evnt: event.AllPlayersPassedPriorityEvent{},
			game: state.Game{},
			want: state.Game{},
		},
		{
			name: "with PassPriorityEvent",
			evnt: event.PassPriorityEvent{
				PlayerID: playerID,
			},
			game: state.LoadGameFromConfig(statetest.GameConfig{
				Players: []statetest.PlayerConfig{{
					ID: playerID,
				}},
			}),
			want: state.LoadGameFromConfig(statetest.GameConfig{
				PlayersPassedPriority: map[string]bool{
					playerID: true,
				},
				Players: []statetest.PlayerConfig{{
					ID: playerID,
				}},
			}),
		},
		{
			name: "with ReceivePriorityEvent",
			evnt: event.ReceivePriorityEvent{
				PlayerID: playerID,
			},
			game: state.Game{},
			want: state.Game{},
		},
		{
			name: "with ResetPriorityPassesEvent",
			evnt: event.ResetPriorityPassesEvent{},
			game: state.LoadGameFromConfig(statetest.GameConfig{
				PlayersPassedPriority: map[string]bool{
					playerID: true,
				},
				Players: []statetest.PlayerConfig{{
					ID: playerID,
				}},
			}),
			want: state.LoadGameFromConfig(statetest.GameConfig{
				Players: []statetest.PlayerConfig{{
					ID: playerID,
				}},
			}),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := applyPriorityEvent(tc.game, tc.evnt)
			if err != nil {
				t.Fatalf("applyPriorityEvent(game, %T); err = %v; want %v", tc.evnt, err, nil)
			}
			if diff := cmp.Diff(tc.want, got, AllowAllUnexported, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("applyPriorityEvent(game, %T) mismatch (-want +got):\n%s", tc.evnt, diff)
			}
		})
	}
}

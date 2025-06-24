package node_test

import (
	"deckronomicon/packages/agent/auto/strategy/evalstate"
	"deckronomicon/packages/agent/auto/strategy/node"
	"deckronomicon/packages/agent/auto/strategy/predicate"
	"deckronomicon/packages/game/gob/gobtest"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"deckronomicon/packages/state/statetest"
	"testing"
)

func TestInZone(t *testing.T) {
	playerID := "Test Player"
	tests := []struct {
		name string
		when node.InZone
		is   state.Game
		want bool
	}{
		{
			name: "with when Island in Hand is Island in Hand",
			when: node.InZone{
				Cards: &predicate.Name{Name: "Island"},
				Zone:  mtg.ZoneHand,
			},
			is: state.LoadGameFromConfig(statetest.GameConfig{
				Players: []statetest.PlayerConfig{{
					ID: playerID,
					Hand: statetest.HandConfig{
						Cards: []gobtest.CardConfig{{Name: "Island"}},
					},
				}},
			}),
			want: true,
		},
		{
			name: "with when Island in Hand is Swamp in Hand",
			when: node.InZone{
				Cards: &predicate.Name{Name: "Island"},
				Zone:  mtg.ZoneHand,
			},
			is: state.LoadGameFromConfig(statetest.GameConfig{
				Players: []statetest.PlayerConfig{{
					ID: playerID,
					Hand: statetest.HandConfig{
						Cards: []gobtest.CardConfig{{Name: "Swamp"}},
					},
				}},
			}),
			want: false,
		},
		{
			name: "with when Island in Hand is Swamp and Island in Hand",
			when: node.InZone{
				Cards: &predicate.Name{Name: "Island"},
				Zone:  mtg.ZoneHand,
			},
			is: state.LoadGameFromConfig(statetest.GameConfig{
				Players: []statetest.PlayerConfig{{
					ID: playerID,
					Hand: statetest.HandConfig{
						Cards: []gobtest.CardConfig{
							{Name: "Swamp"},
							{Name: "Island"},
						},
					},
				}},
			}),
			want: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := evalstate.EvalState{
				Game:     test.is,
				PlayerID: playerID,
			}
			got := test.when.Evaluate(&ctx)
			if got != test.want {
				t.Errorf("Evaluate(...) = %t; want %v", got, test.want)
			}
		})
	}
}

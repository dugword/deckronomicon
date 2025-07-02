package resolver

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/state"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func addManaTestGame(playerID string) *state.Game {
	game := state.NewGameFromDefinition(&definition.Game{
		Players: []*definition.Player{{
			ID: playerID,
		}},
	})
	return game
}

func TestAddManaEffect(t *testing.T) {
	playerID := "Test Player"
	testCases := []struct {
		name       string
		effect     *effect.AddMana
		wantEvents []event.GameEvent
	}{
		{
			name:   "with WUBRG mana",
			effect: &effect.AddMana{Mana: "{W}{U}{B}{R}{G}"},
			wantEvents: []event.GameEvent{
				&event.AddManaEvent{PlayerID: playerID, Amount: 1, Color: mana.White},
				&event.AddManaEvent{PlayerID: playerID, Amount: 1, Color: mana.Blue},
				&event.AddManaEvent{PlayerID: playerID, Amount: 1, Color: mana.Black},
				&event.AddManaEvent{PlayerID: playerID, Amount: 1, Color: mana.Red},
				&event.AddManaEvent{PlayerID: playerID, Amount: 1, Color: mana.Green},
			},
		},
		{
			name:   "with 2UU mana",
			effect: &effect.AddMana{Mana: "{2}{U}{U}"},
			wantEvents: []event.GameEvent{
				&event.AddManaEvent{PlayerID: playerID, Amount: 2, Color: mana.Colorless},
				&event.AddManaEvent{PlayerID: playerID, Amount: 2, Color: mana.Blue},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game := addManaTestGame(playerID)
			got, err := ResolveAddMana(game, playerID, tc.effect)
			if err != nil {
				t.Fatalf("AddMana.Resolve(...); err = %v; want = %v", err, nil)
			}
			if diff := cmp.Diff(tc.wantEvents, got.Events); diff != "" {
				t.Errorf("AddMana.Resolve(...) mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

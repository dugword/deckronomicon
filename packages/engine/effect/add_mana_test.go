package effect

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/state"
	"deckronomicon/packages/state/statetest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func addManaTestGame(playerID string) state.Game {
	game := state.LoadGameFromConfig(statetest.GameConfig{
		Players: []statetest.PlayerConfig{{
			ID: playerID,
		}},
	})
	return game
}

func TestAddManaEffectName(t *testing.T) {
	addManaEffect, err := NewAddManaEffect(definition.EffectSpec{
		Name:      "AddMana",
		Modifiers: map[string]any{"Mana": "{W}{U}{B}{R}{G}"},
	})
	if err != nil {
		t.Fatalf("NewAddManaEffect(EffectSpec); err = %v; want = %v", err, nil)
	}
	if addManaEffect.Name() != "AddMana" {
		t.Errorf("AddManaEffect.Name() = %q; want = %q", addManaEffect.Name(), "AddMana")
	}
}

func TestAddManaEffect(t *testing.T) {
	playerID := "Test Player"
	testCases := []struct {
		name       string
		modifiers  map[string]any
		wantEvents []event.GameEvent
	}{
		{
			name:      "with WUBRG mana",
			modifiers: map[string]any{"Mana": "{W}{U}{B}{R}{G}"},
			wantEvents: []event.GameEvent{
				event.AddManaEvent{PlayerID: playerID, Amount: 1, Color: mana.White},
				event.AddManaEvent{PlayerID: playerID, Amount: 1, Color: mana.Blue},
				event.AddManaEvent{PlayerID: playerID, Amount: 1, Color: mana.Black},
				event.AddManaEvent{PlayerID: playerID, Amount: 1, Color: mana.Red},
				event.AddManaEvent{PlayerID: playerID, Amount: 1, Color: mana.Green},
			},
		},
		{
			name:      "with 2UU mana",
			modifiers: map[string]any{"Mana": "{2}{U}{U}"},
			wantEvents: []event.GameEvent{
				event.AddManaEvent{PlayerID: playerID, Amount: 2, Color: mana.Colorless},
				event.AddManaEvent{PlayerID: playerID, Amount: 2, Color: mana.Blue},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game := addManaTestGame(playerID)
			player := game.GetPlayer(playerID)
			addManaEffect, err := NewAddManaEffect(definition.EffectSpec{
				Name:      "AddMana",
				Modifiers: tc.modifiers,
			})
			if err != nil {
				t.Fatalf("NewAddManaEffect(EffectSpec); err = %v; want = %v", err, nil)
			}
			resEnv := &resenv.ResEnv{}
			got, err := addManaEffect.Resolve(game, player, nil, target.TargetValue{}, resEnv)
			if err != nil {
				t.Fatalf("AddManaEffect.Resolve(...); err = %v; want = %v", err, nil)
			}
			if diff := cmp.Diff(tc.wantEvents, got.Events); diff != "" {
				t.Errorf("AddManaEffect.Resolve(...) mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

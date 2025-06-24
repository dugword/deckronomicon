package effect

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/state"
	"deckronomicon/packages/state/statetest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func additionalManaTestGame(playerID string) state.Game {
	game := state.LoadGameFromConfig(statetest.GameConfig{
		Players: []statetest.PlayerConfig{{
			ID: playerID,
		}},
	})
	return game
}

func TestAdditionalManaEffectName(t *testing.T) {
	additionalManaEffect, err := NewAdditionalManaEffect(definition.EffectSpec{
		Name: "AdditionalMana",
		Modifiers: map[string]any{
			"Subtype":  "Island",
			"Mana":     "{U}",
			"Duration": "EndOfTurn",
		},
	})
	if err != nil {
		t.Fatalf("NewAdditionalManaEffect(EffectSpec); err = %v; want = %v", err, nil)
	}
	if additionalManaEffect.Name() != "AdditionalMana" {
		t.Errorf("AdditionalManaEffect.Name() = %q; want = %q", additionalManaEffect.Name(), "AdditionalMana")
	}
}

func TestAdditionalManaEffect(t *testing.T) {
	playerID := "Test Player"

	testCases := []struct {
		name       string
		modifiers  map[string]any
		wantEvents []event.GameEvent
	}{
		{
			name: "with WUBRG mana from Island",
			modifiers: map[string]any{
				"Subtype":  "Island",
				"Mana":     "{W}{U}{B}{R}{G}",
				"Duration": "EndOfTurn",
			},
			wantEvents: []event.GameEvent{
				event.RegisterTriggeredEffectEvent{
					PlayerID: "Test Player",
					Trigger: state.Trigger{
						EventType: "LandTappedForMana",
						Filter: state.Filter{
							Subtypes: []mtg.Subtype{"Island"},
						},
					},
					EffectSpecs: []definition.EffectSpec{{
						Name: "AddMana",
						Modifiers: map[string]any{
							"Mana": "{W}{U}{B}{R}{G}",
						},
					}},
					Duration: "EndOfTurn",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game := additionalManaTestGame(playerID)
			player := game.GetPlayer(playerID)
			additionalManaEffect, err := NewAdditionalManaEffect(definition.EffectSpec{
				Name:      "AdditionalMana",
				Modifiers: tc.modifiers,
			})
			if err != nil {
				t.Fatalf("NewAdditionalManaEffect(EffectSpec); err = %v; want = %v", err, nil)
			}
			resEnv := &resenv.ResEnv{}
			got, err := additionalManaEffect.Resolve(game, player, nil, target.TargetValue{}, resEnv)
			if err != nil {
				t.Fatalf("AdditionalManaEffect.Resolve(...); err = %v; want = %v", err, nil)
			}
			if diff := cmp.Diff(tc.wantEvents, got.Events); diff != "" {
				t.Errorf("AdditionalManaEffect.Resolve(...) mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

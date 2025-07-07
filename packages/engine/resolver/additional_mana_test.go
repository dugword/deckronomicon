package resolver

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func additionalManaTestGame(playerID string) *state.Game {
	game := state.NewGameFromDefinition(&definition.Game{
		Players: []*definition.Player{{
			ID: playerID,
		}},
	})
	return game
}

func TestResolveAdditionalMana(t *testing.T) {
	playerID := "Test Player"
	testCases := []struct {
		name       string
		effect     *effect.AdditionalMana
		wantEvents []event.GameEvent
	}{
		{
			name: "with WUBRG mana from Island",
			effect: &effect.AdditionalMana{
				Subtype:  "Island",
				Mana:     "{W}{U}{B}{R}{G}",
				Duration: mtg.DurationEndOfTurn,
			},
			wantEvents: []event.GameEvent{
				&event.RegisterTriggeredAbilityEvent{
					PlayerID: "Test Player",
					Trigger: gob.Trigger{
						EventType: "LandTappedForMana",
						Filter: query.Opts{
							Subtypes: []mtg.Subtype{"Island"},
						},
					},
					Effects: []effect.Effect{&effect.AddMana{
						Mana: "{W}{U}{B}{R}{G}",
					}},
					Duration: "EndOfTurn",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ResolveAdditionalMana(playerID, tc.effect)
			if err != nil {
				t.Fatalf("AdditionalMana.Resolve(...); err = %v; want = %v", err, nil)
			}
			if diff := cmp.Diff(tc.wantEvents, got.Events); diff != "" {
				t.Errorf("AdditionalMana.Resolve(...) mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

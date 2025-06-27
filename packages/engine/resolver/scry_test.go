package resolver

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/state"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func scryTestGame(playerID string) state.Game {
	game := state.NewGameFromDefinition(definition.Game{
		Players: []definition.Player{{
			ID: playerID,
			Library: definition.Library{
				Cards: []definition.Card{
					{ID: "Card 1 ID"},
					{ID: "Card 2 ID"},
					{ID: "Card 3 ID"},
					{ID: "Card 4 ID"},
					{ID: "Card 5 ID"},
				},
			},
		}},
	})
	return game
}

// TODO: Maps in the Choose package make this non-deterministic
func TestScryEffect(t *testing.T) {
	playerID := "Test Player"
	testCases := []struct {
		name       string
		count      int
		wantEvents []event.GameEvent
	}{
		{
			name:  "with count 2",
			count: 2,
			wantEvents: []event.GameEvent{
				event.PutCardOnTopOfLibraryEvent{PlayerID: playerID, CardID: "Card 1 ID", FromZone: "Library"},
				event.PutCardOnBottomOfLibraryEvent{PlayerID: playerID, CardID: "Card 2 ID", FromZone: "Library"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game := scryTestGame(playerID)
			efct := effect.Scry{
				Count: tc.count,
			}
			effectResult, err := ResolveScry(game, playerID, efct, nil)
			if err != nil {
				t.Fatalf("ScryEffect.Resolve(...); err = %v; want = %v", err, nil)
			}
			mapChoicesToBucketsOpts, ok := effectResult.ChoicePrompt.ChoiceOpts.(choose.MapChoicesToBucketsOpts)
			if !ok {
				t.Fatalf("ChoiceOpts is not MapChoicesToBucketsOpts; got = %T; want = %T", effectResult.ChoicePrompt.ChoiceOpts, choose.MapChoicesToBucketsOpts{})
			}
			if len(mapChoicesToBucketsOpts.Choices) != tc.count {
				t.Errorf("mapChoicesToBucketsOpts.Count = %d; want = %d", len(mapChoicesToBucketsOpts.Choices), tc.count)
			}
			chooseResults := choose.MapChoicesToBucketsResults{
				Assignments: map[choose.Bucket][]choose.Choice{
					mapChoicesToBucketsOpts.Buckets[0]: mapChoicesToBucketsOpts.Choices[:1],
					mapChoicesToBucketsOpts.Buckets[1]: mapChoicesToBucketsOpts.Choices[1:2],
				},
			}
			got, err := effectResult.Resume(chooseResults)
			if err != nil {
				t.Fatalf("effectResult.ResumeFunc(...) failed; err = %v; want = %v", err, nil)
			}
			if diff := cmp.Diff(tc.wantEvents, got.Events); diff != "" {
				t.Errorf("ScryEffect.Resolve(...) mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

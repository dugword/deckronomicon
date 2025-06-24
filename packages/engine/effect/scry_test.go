package effect

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/gob/gobtest"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/state"
	"deckronomicon/packages/state/statetest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func scryTestGame(playerID string) state.Game {
	game := state.LoadGameFromConfig(statetest.GameConfig{
		Players: []statetest.PlayerConfig{{
			ID: playerID,
			Library: statetest.LibraryConfig{
				Cards: []gobtest.CardConfig{
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

func TestScryEffectName(t *testing.T) {
	scryEffect, err := NewScryEffect(definition.EffectSpec{
		Name:      "Scry",
		Modifiers: map[string]any{"Count": 4},
	})
	if err != nil {
		t.Fatalf("NewScryEffect(EffectSpec); err = %v; want = %v", err, nil)
	}
	if scryEffect.Name() != "Scry" {
		t.Errorf("ScryEffect.Name() = %q; want = %q", scryEffect.Name(), "Scry")
	}
}

func TestScryEffect(t *testing.T) {
	playerID := "Test Player"
	testCases := []struct {
		name        string
		count       int
		targetValue target.TargetValue
		wantEvents  []event.GameEvent
	}{
		{
			name:        "with count 2",
			count:       2,
			targetValue: target.TargetValue{},
			wantEvents: []event.GameEvent{
				event.PutCardOnTopOfLibraryEvent{PlayerID: playerID, CardID: "Card 1 ID", FromZone: "Library"},
				event.PutCardOnBottomOfLibraryEvent{PlayerID: playerID, CardID: "Card 2 ID", FromZone: "Library"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game := scryTestGame(playerID)
			player := game.GetPlayer(playerID)
			scryEffect, err := NewScryEffect(definition.EffectSpec{
				Name:      "Scry",
				Modifiers: map[string]any{"Count": tc.count},
			})
			if err != nil {
				t.Fatalf("NewScryEffect(EffectSpec); err = %v; want = %v", err, nil)
			}
			source := gob.LoadSpellFromConfig(gobtest.SpellConfig{ID: "Test Spell"})
			resEnv := &resenv.ResEnv{}
			effectResult, err := scryEffect.Resolve(game, player, source, tc.targetValue, resEnv)
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
			// TODO: Rename to "Resume" no reason to have "Func" in the name
			got, err := effectResult.ResumeFunc(chooseResults)
			if err != nil {
				t.Fatalf("effectResult.ResumeFunc(...) failed; err = %v; want = %v", err, nil)
			}
			if diff := cmp.Diff(tc.wantEvents, got.Events); diff != "" {
				t.Errorf("ScryEffect.Resolve(...) mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

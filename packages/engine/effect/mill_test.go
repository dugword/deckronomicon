package effect

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/gob/gobtest"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/state"
	"deckronomicon/packages/state/statetest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func millTestGame(playerID string) state.Game {
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
		}, {
			ID: "Test Opponent",
			Library: statetest.LibraryConfig{
				Cards: []gobtest.CardConfig{
					{ID: "Opponent Card 1 ID"},
					{ID: "Opponent Card 2 ID"},
					{ID: "Opponent Card 3 ID"},
					{ID: "Opponent Card 4 ID"},
					{ID: "Opponent Card 5 ID"},
				},
			},
		}},
	})
	return game
}

func TestMillEffectName(t *testing.T) {
	millEffect, err := NewMillEffect(definition.EffectSpec{
		Name:      "Mill",
		Modifiers: map[string]any{"Count": 4, "Target": "Player"},
	})
	if err != nil {
		t.Fatalf("NewMillEffect(EffectSpec); err = %v; want = %v", err, nil)
	}
	if millEffect.Name() != "Mill" {
		t.Errorf("MillEffect.Name() = %q; want = %q", millEffect.Name(), "Mill")
	}
}

func TestMillEffect(t *testing.T) {
	playerID := "Test Player"
	testCases := []struct {
		name        string
		Modifers    map[string]any
		targetValue target.TargetValue
		wantEvents  []event.GameEvent
	}{
		{
			name:        "with count 1",
			Modifers:    map[string]any{"Count": 1},
			targetValue: target.TargetValue{},
			wantEvents: []event.GameEvent{
				event.PutCardInGraveyardEvent{PlayerID: playerID, CardID: "Card 1 ID", FromZone: mtg.ZoneLibrary},
			},
		},
		{
			name:     "with count 1 target player",
			Modifers: map[string]any{"Count": 1, "Target": "Player"},
			targetValue: target.TargetValue{
				TargetType: target.TargetTypePlayer,
				TargetID:   "Test Opponent",
			},
			wantEvents: []event.GameEvent{
				event.PutCardInGraveyardEvent{PlayerID: "Test Opponent", CardID: "Opponent Card 1 ID", FromZone: mtg.ZoneLibrary},
			},
		},
		{
			name:     "with count 4 target player",
			Modifers: map[string]any{"Count": 4, "Target": "Player"},
			targetValue: target.TargetValue{
				TargetType: target.TargetTypePlayer,
				TargetID:   playerID,
			},
			wantEvents: []event.GameEvent{
				event.PutCardInGraveyardEvent{PlayerID: playerID, CardID: "Card 1 ID", FromZone: mtg.ZoneLibrary},
				event.PutCardInGraveyardEvent{PlayerID: playerID, CardID: "Card 2 ID", FromZone: mtg.ZoneLibrary},
				event.PutCardInGraveyardEvent{PlayerID: playerID, CardID: "Card 3 ID", FromZone: mtg.ZoneLibrary},
				event.PutCardInGraveyardEvent{PlayerID: playerID, CardID: "Card 4 ID", FromZone: mtg.ZoneLibrary},
			},
		},
		{
			name:        "with count higher than library size",
			Modifers:    map[string]any{"Count": 10},
			targetValue: target.TargetValue{},
			wantEvents: []event.GameEvent{
				event.PutCardInGraveyardEvent{PlayerID: playerID, CardID: "Card 1 ID", FromZone: mtg.ZoneLibrary},
				event.PutCardInGraveyardEvent{PlayerID: playerID, CardID: "Card 2 ID", FromZone: mtg.ZoneLibrary},
				event.PutCardInGraveyardEvent{PlayerID: playerID, CardID: "Card 3 ID", FromZone: mtg.ZoneLibrary},
				event.PutCardInGraveyardEvent{PlayerID: playerID, CardID: "Card 4 ID", FromZone: mtg.ZoneLibrary},
				event.PutCardInGraveyardEvent{PlayerID: playerID, CardID: "Card 5 ID", FromZone: mtg.ZoneLibrary},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game := millTestGame(playerID)
			player := game.GetPlayer(playerID)
			millEffect, err := NewMillEffect(definition.EffectSpec{
				Name:      "Mill",
				Modifiers: tc.Modifers,
			})
			if err != nil {
				t.Fatalf("NewMillEffect(EffectSpec); err = %v; want = %v", err, nil)
			}
			source := gob.LoadSpellFromConfig(gobtest.SpellConfig{ID: "Test Spell"})
			resEnv := &resenv.ResEnv{}
			got, err := millEffect.Resolve(game, player, source, tc.targetValue, resEnv)
			if err != nil {
				t.Fatalf("MillEffect.Resolve(...); err = %v; want = %v", err, nil)
			}
			if diff := cmp.Diff(tc.wantEvents, got.Events); diff != "" {
				t.Errorf("MillEffect.Resolve(...) mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

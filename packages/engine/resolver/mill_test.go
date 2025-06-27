package resolver

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func millTestGame(playerID string) state.Game {
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
		}, {
			ID: "Test Opponent",
			Library: definition.Library{
				Cards: []definition.Card{
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

func TestMill(t *testing.T) {
	playerID := "Test Player"
	testCases := []struct {
		name       string
		effect     effect.Mill
		target     effect.Target
		wantEvents []event.GameEvent
	}{
		{
			name:   "with count 1",
			effect: effect.Mill{Count: 1},
			wantEvents: []event.GameEvent{
				event.PutCardInGraveyardEvent{PlayerID: playerID, CardID: "Card 1 ID", FromZone: mtg.ZoneLibrary},
			},
		},
		{
			name:   "with count 1 target player",
			effect: effect.Mill{Count: 1, Target: mtg.TargetTypePlayer},
			target: effect.Target{
				Type: mtg.TargetTypePlayer,
				ID:   "Test Opponent",
			},
			wantEvents: []event.GameEvent{
				event.PutCardInGraveyardEvent{PlayerID: "Test Opponent", CardID: "Opponent Card 1 ID", FromZone: mtg.ZoneLibrary},
			},
		},
		{
			name:   "with count 4 target player",
			effect: effect.Mill{Count: 4, Target: mtg.TargetTypePlayer},
			target: effect.Target{
				Type: mtg.TargetTypePlayer,
				ID:   playerID,
			},
			wantEvents: []event.GameEvent{
				event.PutCardInGraveyardEvent{PlayerID: playerID, CardID: "Card 1 ID", FromZone: mtg.ZoneLibrary},
				event.PutCardInGraveyardEvent{PlayerID: playerID, CardID: "Card 2 ID", FromZone: mtg.ZoneLibrary},
				event.PutCardInGraveyardEvent{PlayerID: playerID, CardID: "Card 3 ID", FromZone: mtg.ZoneLibrary},
				event.PutCardInGraveyardEvent{PlayerID: playerID, CardID: "Card 4 ID", FromZone: mtg.ZoneLibrary},
			},
		},
		{
			name:   "with count higher than library size",
			effect: effect.Mill{Count: 10},
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
			got, err := ResolveMill(game, playerID, tc.effect, tc.target)
			if err != nil {
				t.Fatalf("MillEffect.Resolve(...); err = %v; want = %v", err, nil)
			}
			if diff := cmp.Diff(tc.wantEvents, got.Events); diff != "" {
				t.Errorf("MillEffect.Resolve(...) mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

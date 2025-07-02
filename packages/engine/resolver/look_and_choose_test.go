package resolver

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func lookAndChooseTestGame(playerID string) *state.Game {
	game := state.NewGameFromDefinition(&definition.Game{
		Players: []*definition.Player{{
			ID: playerID,
			Library: &definition.Library{
				Cards: []*definition.Card{
					{ID: "Card 1 ID"},
					{
						ID:        "Card 2 ID",
						CardTypes: []string{string(mtg.CardTypeSorcery)},
					},
					{ID: "Card 3 ID"},
					{
						ID:        "Card 4 ID",
						CardTypes: []string{string(mtg.CardTypeInstant)},
					},
					{ID: "Card 5 ID"},
				},
			},
		}},
	})
	return game
}

func TestLookAndChoose(t *testing.T) {
	playerID := "Test Player"
	testCases := []struct {
		name       string
		effect     *effect.LookAndChoose
		wantEvents []event.GameEvent
	}{
		{
			name: "with look 5 choose 2 Sorcery or Instant rest to Graveyard",
			effect: &effect.LookAndChoose{
				Look:      5,
				Choose:    2,
				CardTypes: []mtg.CardType{mtg.CardTypeSorcery, mtg.CardTypeInstant},
				Rest:      "Graveyard",
			},
			// TODO: I think the maps make this non-deterministic
			wantEvents: []event.GameEvent{
				&event.PutCardInHandEvent{PlayerID: "Test Player", CardID: "Card 2 ID", FromZone: "Library"},
				&event.PutCardInHandEvent{PlayerID: "Test Player", CardID: "Card 4 ID", FromZone: "Library"},
				&event.PutCardInGraveyardEvent{PlayerID: "Test Player", CardID: "Card 1 ID", FromZone: "Library"},
				&event.PutCardInGraveyardEvent{PlayerID: "Test Player", CardID: "Card 3 ID", FromZone: "Library"},
				&event.PutCardInGraveyardEvent{PlayerID: "Test Player", CardID: "Card 5 ID", FromZone: "Library"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game := lookAndChooseTestGame(playerID)
			result, err := ResolveLookAndChoose(game, playerID, tc.effect, nil)
			if err != nil {
				t.Fatalf("LookAndChooseEffect.Resolve(...); err = %v; want = %v", err, nil)
			}
			chooseManyOpts, ok := result.ChoicePrompt.ChoiceOpts.(choose.ChooseManyOpts)
			if !ok {
				t.Fatalf("ChoiceOpts is not ChooseManyOpts; got = %T; want = %T", result.ChoicePrompt.ChoiceOpts, choose.MapChoicesToBucketsOpts{})
			}
			chooseResults := choose.ChooseManyResults{
				Choices: chooseManyOpts.Choices,
			}
			got, err := result.Resume(chooseResults)
			if err != nil {
				t.Fatalf("result.Resume(...) failed; err = %v; want = %v", err, nil)
			}
			if diff := cmp.Diff(tc.wantEvents, got.Events); diff != "" {
				t.Errorf("ResolveLookAndChoose(...) mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

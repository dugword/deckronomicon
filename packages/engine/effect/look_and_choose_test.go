package effect

import (
	"deckronomicon/packages/choose"
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

func lookAndChooseTestGame(playerID string) state.Game {
	game := state.LoadGameFromConfig(statetest.GameConfig{
		Players: []statetest.PlayerConfig{{
			ID: playerID,
			Library: statetest.LibraryConfig{
				Cards: []gobtest.CardConfig{
					{ID: "Card 1 ID"},
					{
						ID:        "Card 2 ID",
						CardTypes: []mtg.CardType{mtg.CardTypeSorcery},
					},
					{ID: "Card 3 ID"},
					{
						ID:        "Card 4 ID",
						CardTypes: []mtg.CardType{mtg.CardTypeInstant},
					},
					{ID: "Card 5 ID"},
				},
			},
		}},
	})
	return game
}

func TestLookAndChooseEffectName(t *testing.T) {
	LookAndChooseEffect, err := NewLookAndChooseEffect(definition.EffectSpec{
		Name: "LookAndChoose",
		Modifiers: map[string]any{
			"Look":      5,
			"Choose":    2,
			"CardTypes": []any{"Sorcery", "Instant"},
			"Rest":      "Graveyard",
		},
	})
	if err != nil {
		t.Fatalf("NewLookAndChooseEffect(EffectSpec); err = %v; want = %v", err, nil)
	}
	if LookAndChooseEffect.Name() != "LookAndChoose" {
		t.Errorf("LookAndChooseEffect.Name() = %q; want = %q", LookAndChooseEffect.Name(), "LookAndChoose")
	}
}

func TestLookAndChooseEffect(t *testing.T) {
	playerID := "Test Player"
	testCases := []struct {
		name       string
		modifiers  map[string]any
		wantEvents []event.GameEvent
	}{
		{
			name: "with look 5 choose 2 Sorcery or Instant rest to Graveyard",
			modifiers: map[string]any{
				"Look":      5,
				"Choose":    2,
				"CardTypes": []any{"Sorcery", "Instant"},
				"Rest":      "Graveyard",
			},
			// TODO: I think the maps make this non-deterministic
			wantEvents: []event.GameEvent{
				event.PutCardInGraveyardEvent{PlayerID: "Test Player", CardID: "Card 1 ID", FromZone: "Library"},
				event.PutCardInGraveyardEvent{PlayerID: "Test Player", CardID: "Card 2 ID", FromZone: "Library"},
				event.PutCardInGraveyardEvent{PlayerID: "Test Player", CardID: "Card 3 ID", FromZone: "Library"},
				event.PutCardInGraveyardEvent{PlayerID: "Test Player", CardID: "Card 4 ID", FromZone: "Library"},
				event.PutCardInGraveyardEvent{PlayerID: "Test Player", CardID: "Card 5 ID", FromZone: "Library"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game := lookAndChooseTestGame(playerID)
			player := game.GetPlayer(playerID)
			lookAndChooseEffect, err := NewLookAndChooseEffect(definition.EffectSpec{
				Name:      "LookAndChoose",
				Modifiers: tc.modifiers,
			})
			if err != nil {
				t.Fatalf("NewLookAndChooseEffect(EffectSpec); err = %v; want = %v", err, nil)
			}
			source := gob.LoadSpellFromConfig(gobtest.SpellConfig{ID: "Test Spell"})
			resEnv := &resenv.ResEnv{}
			effectResult, err := lookAndChooseEffect.Resolve(game, player, source, target.TargetValue{}, resEnv)
			if err != nil {
				t.Fatalf("LookAndChooseEffect.Resolve(...); err = %v; want = %v", err, nil)
			}
			chooseManyOpts, ok := effectResult.ChoicePrompt.ChoiceOpts.(choose.ChooseManyOpts)
			if !ok {
				t.Fatalf("ChoiceOpts is not ChooseManyOpts; got = %T; want = %T", effectResult.ChoicePrompt.ChoiceOpts, choose.MapChoicesToBucketsOpts{})
			}
			chooseResults := choose.ChooseManyResults{
				Choices: chooseManyOpts.Choices,
			}
			got, err := effectResult.ResumeFunc(chooseResults)
			if err != nil {
				t.Fatalf("effectResult.ResumeFunc(...) failed; err = %v; want = %v", err, nil)
			}
			if diff := cmp.Diff(tc.wantEvents, got.Events); diff != "" {
				t.Errorf("LookAndChooseEffect.Resolve(...) mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

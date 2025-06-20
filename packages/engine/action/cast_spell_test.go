package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCastSpellActionComplete(t *testing.T) {
	const playerID = "Test Player"
	testCases := []struct {
		name   string
		action CastSpellAction
		want   []event.GameEvent
	}{
		{
			name: "with basic spell",
			action: CastSpellAction{
				playerID: playerID,
				cardID:   "Test Card ID",
			},
			want: []event.GameEvent{
				event.CastSpellEvent{PlayerID: "Test Player", CardID: "Test Card ID", FromZone: mtg.ZoneHand},
				event.PutSpellOnStackEvent{PlayerID: "Test Player", CardID: "Test Card ID", FromZone: mtg.ZoneHand},
			},
		},
		{
			name: "with effects",
			action: CastSpellAction{
				playerID: playerID,
				cardID:   "Card with Effects ID",
				targetsForEffects: map[EffectTargetKey]target.TargetValue{
					{SourceID: "Card with Effects ID", EffectIndex: 0}: target.TargetValue{
						ObjectID: "Target Object ID",
					},
					{SourceID: "Card with Effects ID", EffectIndex: 1}: target.TargetValue{
						PlayerID: "Target Player ID",
					},
				},
			},
			want: []event.GameEvent{
				event.CastSpellEvent{
					PlayerID: "Test Player",
					CardID:   "Card with Effects ID",
					FromZone: mtg.ZoneHand,
				},
				event.PutSpellOnStackEvent{
					PlayerID: "Test Player",
					CardID:   "Card with Effects ID",
					FromZone: mtg.ZoneHand,
					EffectWithTargets: []gob.EffectWithTarget{
						{
							EffectSpec: definition.EffectSpec{
								Name:      "Effect 1",
								Modifiers: json.RawMessage(`{"Target": "Permanent"}`),
							},
							Target:   target.TargetValue{ObjectID: "Target Object ID"},
							SourceID: "Card with Effects ID",
						},
						{
							EffectSpec: definition.EffectSpec{
								Name:      "Effect 2",
								Modifiers: json.RawMessage(`{"Target": "Player"}`),
							},
							Target:   target.TargetValue{PlayerID: "Target Player ID"},
							SourceID: "Card with Effects ID",
						},
					},
				},
			},
		},
		{
			name: "with flashback spell",
			action: CastSpellAction{
				playerID:  playerID,
				cardID:    "Card with Flashback ID",
				flashback: true,
			},
			want: []event.GameEvent{
				event.CastSpellEvent{PlayerID: "Test Player", CardID: "Card with Flashback ID", FromZone: mtg.ZoneGraveyard},
				event.PutSpellOnStackEvent{PlayerID: "Test Player", CardID: "Card with Flashback ID", FromZone: mtg.ZoneGraveyard, Flashback: true},
			},
		},
		{
			name: "with spell that has targets",
			action: CastSpellAction{
				playerID: playerID,
				cardID:   "Card with Target ID",
				targetsForEffects: map[EffectTargetKey]target.TargetValue{
					{SourceID: "Card with Target ID", EffectIndex: 0}: target.TargetValue{
						PlayerID: playerID,
					},
				},
			},
			want: []event.GameEvent{
				event.CastSpellEvent{
					PlayerID: "Test Player",
					CardID:   "Card with Target ID",
					FromZone: mtg.ZoneHand,
				},
				event.PutSpellOnStackEvent{
					PlayerID: "Test Player",
					CardID:   "Card with Target ID",
					FromZone: mtg.ZoneHand,
					EffectWithTargets: []gob.EffectWithTarget{
						{
							EffectSpec: definition.EffectSpec{
								Name:      "Target",
								Modifiers: json.RawMessage(`{"Target": "Player"}`),
							},
							Target:   target.TargetValue{PlayerID: "Test Player"},
							SourceID: "Card with Target ID",
						},
					},
				},
			},
		},
		{
			name: "with replicated spell",
			action: CastSpellAction{
				playerID:       playerID,
				cardID:         "Card with Replicate ID",
				replicateCount: 3,
			},
			want: []event.GameEvent{
				event.CastSpellEvent{PlayerID: "Test Player", CardID: "Card with Replicate ID", FromZone: mtg.ZoneHand},
				event.PutSpellOnStackEvent{PlayerID: "Test Player", CardID: "Card with Replicate ID", FromZone: mtg.ZoneHand},
				event.PutAbilityOnStackEvent{
					PlayerID:    "Test Player",
					SourceID:    "Card with Replicate ID",
					FromZone:    "Hand",
					AbilityName: "Replicate",
					EffectWithTargets: []gob.EffectWithTarget{
						{
							EffectSpec: definition.EffectSpec{
								Name:      string(mtg.StaticKeywordReplicate),
								Modifiers: json.RawMessage(`{"Count": 3}`),
							},
							Target: target.TargetValue{
								ObjectID: "Card with Replicate ID",
							},
							SourceID: "Card with Replicate ID",
						},
					},
				},
			},
		},
		{
			name: "with splice spell",
			action: CastSpellAction{
				playerID:      playerID,
				cardID:        "Acane Card ID",
				spliceCardIDs: []string{"Card with Splice ID"},
			},
			want: []event.GameEvent{
				event.CastSpellEvent{PlayerID: "Test Player", CardID: "Acane Card ID", FromZone: mtg.ZoneHand},
				event.PutSpellOnStackEvent{
					PlayerID: "Test Player",
					CardID:   "Acane Card ID",
					FromZone: mtg.ZoneHand,
				},
			},
		},
	}
	for _, tc := range testCases {
		game := newTestGame(playerID)
		resEnv := resenv.ResEnv{}
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.action.Complete(game, &resEnv)
			if err != nil {
				t.Fatalf("action.Complete(...); err = %v; want %v", err, nil)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("action.Complete(...); mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

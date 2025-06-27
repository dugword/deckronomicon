package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
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
				cardID: "Test Card ID",
			},
			want: []event.GameEvent{
				event.CastSpellEvent{PlayerID: "Test Player", CardID: "Test Card ID", FromZone: mtg.ZoneHand},
				event.PutSpellOnStackEvent{PlayerID: "Test Player", CardID: "Test Card ID", FromZone: mtg.ZoneHand},
			},
		},
		{
			name: "with effects",
			action: CastSpellAction{
				cardID: "Card with Effects ID",
				targetsForEffects: map[effect.EffectTargetKey]effect.Target{
					{SourceID: "Card with Effects ID", EffectIndex: 0}: effect.Target{
						ID: "Target Object ID",
					},
					{SourceID: "Card with Effects ID", EffectIndex: 1}: effect.Target{
						ID: "Another Target Object ID",
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
					EffectWithTargets: []effect.EffectWithTarget{
						{
							Effect: effect.TargetEffect{
								Target: "Permanent",
							},
							Target:   effect.Target{ID: "Target Object ID"},
							SourceID: "Card with Effects ID",
						},
						{
							Effect: effect.TargetEffect{
								Target: "Permanent",
							},
							Target:   effect.Target{ID: "Another Target Object ID"},
							SourceID: "Card with Effects ID",
						},
					},
				},
			},
		},
		{
			name: "with flashback spell",
			action: CastSpellAction{
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
				cardID: "Card with Target ID",
				targetsForEffects: map[effect.EffectTargetKey]effect.Target{
					{SourceID: "Card with Target ID", EffectIndex: 0}: effect.Target{
						ID: playerID,
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
					EffectWithTargets: []effect.EffectWithTarget{
						{
							Effect: effect.TargetEffect{
								Target: "Permanent",
							},
							Target:   effect.Target{ID: "Test Player"},
							SourceID: "Card with Target ID",
						},
					},
				},
			},
		},
		{
			name: "with replicated spell",
			action: CastSpellAction{
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
					EffectWithTargets: []effect.EffectWithTarget{
						{
							Effect: effect.Replicate{
								Count: 3,
							},
							Target: effect.Target{
								ID: "Card with Replicate ID",
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
			player := game.GetPlayer(playerID)
			got, err := tc.action.Complete(game, player, &resEnv)
			if err != nil {
				t.Fatalf("action.Complete(...); err = %v; want %v", err, nil)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("action.Complete(...); mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

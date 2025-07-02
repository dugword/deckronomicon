package pay_test

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/pay"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPayCost(t *testing.T) {
	const playerID = "Test Player"
	tests := []struct {
		name   string
		cost   string
		object gob.Object
		want   []event.GameEvent
		target target.Target
	}{
		{
			name: "with white mana cost",
			cost: "{W}",
			want: []event.GameEvent{
				&event.SpendManaEvent{
					ManaString: "{W}",
					PlayerID:   playerID,
				},
			},
		},
		{
			name: "with wubrg mana cost",
			cost: "{W}{U}{B}{R}{G}",
			want: []event.GameEvent{
				&event.SpendManaEvent{
					ManaString: "{W}{U}{B}{R}{G}",
					PlayerID:   playerID,
				},
			},
		},
		{
			name: "with pay 3 life cost",
			cost: "Pay 3 life",
			want: []event.GameEvent{
				&event.LoseLifeEvent{
					PlayerID: playerID,
					Amount:   3,
				},
			},
		},
		{
			name: "with tap this cost",
			cost: "{T}",
			object: gob.NewPermanentFromDefinition((&definition.Permanent{
				ID: "Object ID",
			})),
			want: []event.GameEvent{
				&event.TapPermanentEvent{
					PlayerID:    playerID,
					PermanentID: "Object ID",
				},
			},
		},
		{
			name: "with discard cost",
			cost: "Discard this card",
			object: gob.NewCardFromDefinition(&definition.Card{
				ID: "Object ID",
			}),
			want: []event.GameEvent{
				&event.DiscardCardEvent{
					PlayerID: playerID,
					CardID:   "Object ID",
				},
			},
		},
		{
			name: "with discard a card cost",
			cost: "Discard a card",
			target: target.Target{
				ID:   "Card ID",
				Type: mtg.TargetTypeCard,
			},
			want: []event.GameEvent{
				&event.DiscardCardEvent{
					PlayerID: playerID,
					CardID:   "Card ID",
				},
			},
		},
		{
			name: "with composite cost of mana and life",
			cost: `{W}, Pay 2 life`,
			want: []event.GameEvent{
				&event.SpendManaEvent{
					ManaString: "{W}",
					PlayerID:   playerID,
				},
				&event.LoseLifeEvent{
					PlayerID: playerID,
					Amount:   2,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testCost, err := cost.Parse(test.cost)
			if err != nil {
				t.Fatalf("cost.Parse(%s); err = %v; want %v", testCost, err, nil)
			}
			if costWithTarget, ok := testCost.(cost.CostWithTarget); ok && test.target.ID != "" {
				testCost = costWithTarget.WithTarget(test.target)
			}
			got, err := pay.Cost(testCost, test.object, playerID)
			if err != nil {
				t.Fatalf("pay.Cost(%s) returned error: %v", testCost.Description(), err)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("pay.Cost(%s) mismatch (-want +got):\n%s", testCost.Description(), diff)
			}
		})
	}
}

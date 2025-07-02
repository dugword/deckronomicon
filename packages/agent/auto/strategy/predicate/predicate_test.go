package predicate

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/query"
	"testing"
)

func TestNameMatchCondition(t *testing.T) {
	tests := []struct {
		name string
		when Name
		is   []string
		want bool
	}{
		{
			name: "with when Island is Island",
			when: Name{Name: "Island"},
			is:   []string{"Island"},
			want: true,
		},
		{
			name: "with when Island is Swamp",
			when: Name{Name: "Island"},
			is:   []string{"Swamp"},
			want: false,
		},
		{
			name: "with when Island is Island and Swamp",
			when: Name{Name: "Island"},
			is:   []string{"Island", "Swamp"},
			want: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var objs []query.Object
			for _, name := range test.is {
				card := gob.NewCardFromDefinition(&definition.Card{Name: name})
				objs = append(objs, card)
			}
			got := test.when.Matches(objs)
			if got != test.want {
				t.Errorf("Matches(%v) = %v; want %v", test.is, got, test.want)
			}
		})
	}
}

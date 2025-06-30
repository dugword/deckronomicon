package cost

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var AllowAllUnexported = cmp.Exporter(func(reflect.Type) bool { return true })

func TestHasCostType(t *testing.T) {
	tests := []struct {
		name     string
		cost     Cost
		costType Cost
		want     bool
	}{
		{
			name: "with has discard this cost",
			cost: Composite{
				costs: []Cost{
					Mana{},
					DiscardThis{},
				},
			},
			costType: DiscardThis{},
			want:     true,
		},
		{
			name: "with no discard this cost",
			cost: Composite{
				costs: []Cost{Mana{}},
			},
			costType: DiscardThis{},
			want:     false,
		},
		{
			name: "with has discard a card cost",
			cost: Composite{
				costs: []Cost{
					Mana{},
					DiscardACard{},
				},
			},
			costType: DiscardACard{},
			want:     true,
		},
		{
			name: "with no discard a card cost",
			cost: Composite{
				costs: []Cost{Mana{}},
			},
			costType: DiscardACard{},
			want:     false,
		},
		{
			name: "with has life cost",
			cost: Composite{
				costs: []Cost{
					Mana{},
					Life{},
				},
			},
			costType: Life{},
			want:     true,
		},
		{
			name: "with no life cost",
			cost: Composite{
				costs: []Cost{
					Mana{},
				},
			},
			costType: Life{},
			want:     false,
		},
		{
			name: "with has mana cost",
			cost: Composite{
				costs: []Cost{
					Mana{},
					TapThis{},
				},
			},
			costType: Mana{},
			want:     true,
		},
		{
			name: "with no mana cost",
			cost: Composite{
				costs: []Cost{
					TapThis{},
				},
			},
			costType: Mana{},
			want:     false,
		},
		{
			name: "with has tap this cost",
			cost: Composite{
				costs: []Cost{
					TapThis{},
					Mana{},
				},
			},
			costType: TapThis{},
			want:     true,
		},
		{
			name: "with no tap this cost",
			cost: Composite{
				costs: []Cost{
					Mana{},
				},
			},
			costType: TapThis{},
			want:     false,
		},
		{
			name: "with has composite cost",
			cost: Composite{
				costs: []Cost{
					TapThis{},
					Mana{},
					DiscardThis{},
				},
			},
			costType: Composite{},
			want:     true,
		},
		{
			name:     "with no composite cost",
			cost:     Mana{},
			costType: Composite{},
			want:     false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := HasType(test.cost, test.costType); got != test.want {
				t.Errorf("HasCostType() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestNewComposite(t *testing.T) {
	tests := []struct {
		name  string
		costs []Cost
		want  Composite
	}{
		{
			name:  "with no costs",
			costs: []Cost{},
			want:  Composite{},
		},
		{
			name:  "with single cost",
			costs: []Cost{Mana{}},
			want:  Composite{costs: []Cost{Mana{}}},
		},
		{
			name:  "with multiple costs",
			costs: []Cost{Mana{}, Mana{}},
			want:  Composite{costs: []Cost{Mana{}, Mana{}}},
		},
		{
			name:  "with composite costs",
			costs: []Cost{Composite{costs: []Cost{Mana{}}}, Mana{}},
			want:  Composite{costs: []Cost{Mana{}, Mana{}}},
		},
		{
			name: "with nested composite costs",
			costs: []Cost{
				Composite{
					costs: []Cost{
						Composite{
							costs: []Cost{
								Mana{},
								Composite{
									costs: []Cost{
										Mana{},
									},
								},
							},
						},
					},
				},
				Composite{
					costs: []Cost{
						Composite{
							costs: []Cost{
								Mana{},
							},
						},
						Mana{},
					},
				},
			},
			want: Composite{costs: []Cost{Mana{}, Mana{}, Mana{}, Mana{}}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewComposite(test.costs...)
			if diff := cmp.Diff(test.want, got, AllowAllUnexported); diff != "" {
				t.Errorf("NewComposite() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

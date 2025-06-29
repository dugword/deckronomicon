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
			cost: CompositeCost{
				costs: []Cost{
					ManaCost{},
					DiscardThisCost{},
				},
			},
			costType: DiscardThisCost{},
			want:     true,
		},
		{
			name: "with no discard this cost",
			cost: CompositeCost{
				costs: []Cost{ManaCost{}},
			},
			costType: DiscardThisCost{},
			want:     false,
		},
		{
			name: "with has life cost",
			cost: CompositeCost{
				costs: []Cost{
					ManaCost{},
					LifeCost{},
				},
			},
			costType: LifeCost{},
			want:     true,
		},
		{
			name: "with no life cost",
			cost: CompositeCost{
				costs: []Cost{
					ManaCost{},
				},
			},
			costType: LifeCost{},
			want:     false,
		},
		{
			name: "with has mana cost",
			cost: CompositeCost{
				costs: []Cost{
					ManaCost{},
					TapThisCost{},
				},
			},
			costType: ManaCost{},
			want:     true,
		},
		{
			name: "with no mana cost",
			cost: CompositeCost{
				costs: []Cost{
					TapThisCost{},
				},
			},
			costType: ManaCost{},
			want:     false,
		},
		{
			name: "with has tap this cost",
			cost: CompositeCost{
				costs: []Cost{
					TapThisCost{},
					ManaCost{},
				},
			},
			costType: TapThisCost{},
			want:     true,
		},
		{
			name: "with no tap this cost",
			cost: CompositeCost{
				costs: []Cost{
					ManaCost{},
				},
			},
			costType: TapThisCost{},
			want:     false,
		},
		{
			name: "with has composite cost",
			cost: CompositeCost{
				costs: []Cost{
					TapThisCost{},
					ManaCost{},
					DiscardThisCost{},
				},
			},
			costType: CompositeCost{},
			want:     true,
		},
		{
			name:     "with no composite cost",
			cost:     ManaCost{},
			costType: CompositeCost{},
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
		want  CompositeCost
	}{
		{
			name:  "with no costs",
			costs: []Cost{},
			want:  CompositeCost{},
		},
		{
			name:  "with single cost",
			costs: []Cost{ManaCost{}},
			want:  CompositeCost{costs: []Cost{ManaCost{}}},
		},
		{
			name:  "with multiple costs",
			costs: []Cost{ManaCost{}, ManaCost{}},
			want:  CompositeCost{costs: []Cost{ManaCost{}, ManaCost{}}},
		},
		{
			name:  "with composite costs",
			costs: []Cost{CompositeCost{costs: []Cost{ManaCost{}}}, ManaCost{}},
			want:  CompositeCost{costs: []Cost{ManaCost{}, ManaCost{}}},
		},
		{
			name: "with nested composite costs",
			costs: []Cost{
				CompositeCost{
					costs: []Cost{
						CompositeCost{
							costs: []Cost{
								ManaCost{},
								CompositeCost{
									costs: []Cost{
										ManaCost{},
									},
								},
							},
						},
					},
				},
				CompositeCost{
					costs: []Cost{
						CompositeCost{
							costs: []Cost{
								ManaCost{},
							},
						},
						ManaCost{},
					},
				},
			},
			want: CompositeCost{costs: []Cost{ManaCost{}, ManaCost{}, ManaCost{}, ManaCost{}}},
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

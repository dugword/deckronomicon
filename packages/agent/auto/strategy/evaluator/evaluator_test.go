package evaluator_test

import (
	"deckronomicon/packages/agent/auto/strategy/evalstate"
	"deckronomicon/packages/agent/auto/strategy/evaluator"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"testing"
)

func TestAnd(t *testing.T) {
	tests := []struct {
		name string
		when evaluator.And
		want bool
	}{
		{
			name: "with all true",
			when: evaluator.And{
				Evaluators: []evaluator.Evaluator{
					&evaluator.True{},
					&evaluator.True{},
				},
			},
			want: true,
		},
		{
			name: "with when one true one false",
			when: evaluator.And{
				Evaluators: []evaluator.Evaluator{
					&evaluator.True{},
					&evaluator.False{},
				},
			},
			want: false,
		},
		{
			name: "with when all false",
			when: evaluator.And{
				Evaluators: []evaluator.Evaluator{
					&evaluator.False{},
					&evaluator.False{},
				},
			},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := evalstate.EvalState{}
			got := test.when.Evaluate(&ctx)
			if got != test.want {
				t.Errorf("Evaluate(...) = %t; want %v", got, test.want)
			}
		})
	}
}

func TestOr(t *testing.T) {
	tests := []struct {
		name string
		when evaluator.Or
		want bool
	}{
		{
			name: "with all true",
			when: evaluator.Or{
				Evaluators: []evaluator.Evaluator{
					&evaluator.True{},
					&evaluator.True{},
				},
			},
			want: true,
		},
		{
			name: "with when one true one false",
			when: evaluator.Or{
				Evaluators: []evaluator.Evaluator{
					&evaluator.True{},
					&evaluator.False{},
				},
			},
			want: true,
		},
		{
			name: "with when all false",
			when: evaluator.Or{
				Evaluators: []evaluator.Evaluator{
					&evaluator.False{},
					&evaluator.False{},
				},
			},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := evalstate.EvalState{}
			got := test.when.Evaluate(&ctx)
			if got != test.want {
				t.Errorf("Evaluate(...) = %t; want %v", got, test.want)
			}
		})
	}
}

func TestMode(t *testing.T) {
	tests := []struct {
		name string
		when evaluator.Mode
		is   string
		want bool
	}{
		{
			name: "with when Setup is Setup",
			when: evaluator.Mode{Mode: "Setup"},
			is:   "Setup",
			want: true,
		},
		{
			name: "with when Setup is Combo",
			when: evaluator.Mode{Mode: "Setup"},
			is:   "Combo",
			want: false,
		},
		{
			name: "with when Combo is Combo",
			when: evaluator.Mode{Mode: "Combo"},
			is:   "Combo",
			want: true,
		},
		{
			name: "with when Setup is empty",
			when: evaluator.Mode{Mode: "Setup"},
			is:   "",
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.when.Evaluate(&evalstate.EvalState{
				Mode: test.is,
			})
			if got != test.want {
				t.Errorf("Evaluate(%q) = %t; want %v", test.is, got, test.want)
			}
		})
	}
}

func TestStep(t *testing.T) {
	tests := []struct {
		name string
		when evaluator.Step
		is   mtg.Step
		want bool
	}{
		{
			name: "with when PrecombaMain is PrecombatMain",
			when: evaluator.Step{Step: mtg.StepPrecombatMain},
			is:   mtg.StepPrecombatMain,
			want: true,
		},
		{
			name: "with when PrecombaMain is Upkeep",
			when: evaluator.Step{Step: mtg.StepPrecombatMain},
			is:   mtg.StepUpkeep,
			want: false,
		},
		{
			name: "with when PrecombaMain is empty",
			when: evaluator.Step{Step: mtg.StepUpkeep},
			is:   "",
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			playerID := "Test Player"
			ctx := evalstate.EvalState{
				Game: state.NewGameFromDefinition(definition.Game{
					ActivePlayerID: playerID,
					Players: []definition.Player{{
						ID: playerID,
					}},
					Step: string(test.is),
				}),
				PlayerID: playerID,
			}
			got := test.when.Evaluate(&ctx)
			if got != test.want {
				t.Errorf("Evaluate(%q) = %t; want %v", test.is, got, test.want)
			}
		})
	}
}

func TestLandPlayedThisTurn(t *testing.T) {
	tests := []struct {
		name string
		when evaluator.LandPlayedThisTurn
		is   bool
		want bool
	}{
		{
			name: "with when false is false",
			when: evaluator.LandPlayedThisTurn{LandPlayedThisTurn: false},
			is:   false,
			want: true,
		},
		{
			name: "with when true is true",
			when: evaluator.LandPlayedThisTurn{LandPlayedThisTurn: true},
			is:   true,
			want: true,
		},
		{
			name: "with when false is true",
			when: evaluator.LandPlayedThisTurn{LandPlayedThisTurn: true},
			is:   false,
			want: false,
		},
		{
			name: "with true is false",
			when: evaluator.LandPlayedThisTurn{LandPlayedThisTurn: false},
			is:   true,
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := evalstate.EvalState{
				Game: state.NewGameFromDefinition(definition.Game{
					Players: []definition.Player{{
						LandPlayedThisTurn: test.is,
					}},
				}),
			}
			got := test.when.Evaluate(&ctx)
			if got != test.want {
				t.Errorf("Evaluate(%t) = %t; want %v", test.is, got, test.want)
			}
		})
	}
}

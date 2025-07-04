package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/evaluator"
	"testing"
)

func TestParseSorcerySpeedEvaluator(t *testing.T) {

	parser, err := NewStrategyParser(StrategyParserData{})
	if err != nil {
		t.Fatalf("Failed to create strategy parser: %v", err)
	}
	eval, err := parser.parseEvaluator(map[string]any{
		"SorcerySpeed": true,
	})
	if err != nil {
		t.Fatalf("Failed to parse SorcerySpeed evaluator: %v", err)
	}

	if eval == nil {
		t.Fatal("Expected evaluator to be parsed successfully, got nil")
	}

	if _, ok := eval.(*evaluator.SorcerySpeed); !ok {
		t.Fatalf("Expected SorcerySpeed evaluator, got %T", eval)
	}

}

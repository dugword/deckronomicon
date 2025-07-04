package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/evaluator"
	"deckronomicon/packages/agent/auto/strategy/node"
	"deckronomicon/packages/game/mtg"
	"fmt"
	"strconv"
	"strings"
)

func (p *StrategyParser) parseEvaluator(raw any) (evaluator.Evaluator, error) {
	switch val := raw.(type) {
	case string:
		eval, ok := p.conditions[val]
		if !ok {
			return nil, fmt.Errorf("unknown condition %q", val)
		}
		return eval, nil
	case map[string]any:
		var evaluators []evaluator.Evaluator
		for k, v := range val {
			switch k {
			case "And", "All", "AllOf":
				eval, err := p.parseLogicalEvaluator("And", v)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'And' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "Or", "Any", "AnyOf":
				eval, err := p.parseLogicalEvaluator("Or", v)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Or' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "Not":
				eval, err := p.parseNotEvaluator(v)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Not' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "True", "Yes":
				evaluators = append(evaluators, &evaluator.True{})
			case "False", "No":
				evaluators = append(evaluators, &evaluator.False{})
			case "Step":
				eval, err := p.parseStepEvaluator(v)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Step' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "Mode":
				eval, err := p.parseModeEvaluator(v)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Mode' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "Player":
				eval, err := p.parsePlayerEvaluator(v)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Player' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "ManaPool":
				eval, err := p.parseManaPoolEvaluator(v)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'ManaPool' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "AvailableMana":
				eval, err := p.parseAvailableManaEvaluator(v)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'AvailableMana' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "InHand":
				eval, err := p.parseInZoneAliasEvaluator("Hand", v)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'InHand' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "InZone":
				eval, err := p.parseInZoneEvaluator(v)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'InZone' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "StackEmpty":
				eval, err := p.parseStackEmptyEvaluator(v)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'StackEmpty' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "SorcerySpeed":
				eval, err := p.parseSorcerySpeedEvaluator(v)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'SorcerySpeed' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "LandPlayedThisTurn":
				eval, err := p.parseLandPlayedThisTurnEvaluator(v)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'LandPlayedThisTurn' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			default:
				return nil, fmt.Errorf("unknown key %q in evaluator", k)
			}
		}
		if len(evaluators) == 0 {
			return nil, fmt.Errorf("no evaluators found in %v", val)
		}
		if len(evaluators) == 1 {
			return evaluators[0], nil
		}
		return &evaluator.And{Evaluators: evaluators}, nil
	default:
		return nil, fmt.Errorf("expected object, got %T", raw)
	}
}

func (p *StrategyParser) parseSorcerySpeedEvaluator(value any) (evaluator.Evaluator, error) {
	sorcerySpeed, ok := value.(bool)
	if !ok {
		return nil, fmt.Errorf("expected boolean for 'SorcerySpeed', got %T", value)
	}
	return &evaluator.SorcerySpeed{
		SorcerySpeed: sorcerySpeed,
	}, nil
}

func (p *StrategyParser) parseStackEmptyEvaluator(value any) (evaluator.Evaluator, error) {
	stackEmpty, ok := value.(bool)
	if !ok {
		return nil, fmt.Errorf("expected boolean for 'StackEmpty', got %T", value)
	}
	return &evaluator.StackEmpty{
		StackEmpty: stackEmpty,
	}, nil
}

func (p *StrategyParser) parseLandPlayedThisTurnEvaluator(value any) (evaluator.Evaluator, error) {
	landPlayedThisTurn, ok := value.(bool)
	if !ok {
		return nil, fmt.Errorf("expected boolean for 'LandPlayedThisTurn', got %T", value)
	}
	return &evaluator.LandPlayedThisTurn{
		LandPlayedThisTurn: landPlayedThisTurn,
	}, nil
}

func (p *StrategyParser) parseNotEvaluator(value any) (evaluator.Evaluator, error) {
	eval, err := p.parseEvaluator(value)
	if err != nil {
		return nil, fmt.Errorf("error parsing 'Not' evaluator: %w", err)
	}
	return &evaluator.Not{Evaluator: eval}, nil
}

func (p *StrategyParser) parseLogicalEvaluator(op string, value any) (evaluator.Evaluator, error) {
	items, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("expected array for '%s', got %T", op, value)
	}
	var evaluators []evaluator.Evaluator
	for _, item := range items {
		eval, err := p.parseEvaluator(item)
		if err != nil {
			return nil, fmt.Errorf("error parsing item in %q: %w", op, err)
		}
		evaluators = append(evaluators, eval)
	}
	switch op {
	case "And":
		return &evaluator.And{Evaluators: evaluators}, nil
	case "Or":
		return &evaluator.Or{Evaluators: evaluators}, nil
	}
	return nil, fmt.Errorf("unknown logical operator %q", op)
}

// TODO: Support a slice of steps
func (p *StrategyParser) parseStepEvaluator(value any) (evaluator.Evaluator, error) {
	switch v := value.(type) {
	case string:
		step, ok := mtg.StringToStep(v)
		if !ok {
			return nil, fmt.Errorf("invalid step %q", v)
		}
		return &evaluator.Step{Step: step}, nil
	case map[string]any:
		s, ok := v["Step"].(string)
		if !ok {
			return nil, fmt.Errorf("expected 'Step' to be a string, got %T", v["Step"])
		}
		return p.parseStepEvaluator(s)
	}
	return nil, fmt.Errorf("step must be a string or an object with a 'Step' key")
}

// TODO: Support a slice of modes
func (p *StrategyParser) parseModeEvaluator(value any) (evaluator.Evaluator, error) {
	switch v := value.(type) {
	case string:
		return &evaluator.Mode{Mode: v}, nil
	case map[string]any:
		if m, ok := v["Mode"].(string); ok {
			return &evaluator.Mode{Mode: m}, nil
		}
	}
	return nil, fmt.Errorf("mode must be a string or an object with a 'Mode' key")
}

func (p *StrategyParser) parseInZoneAliasEvaluator(alias string, value any) (evaluator.Evaluator, error) {
	zoneWrapper := map[string]any{
		"Zone":  alias,
		"Cards": value,
	}
	return p.parseInZoneEvaluator(zoneWrapper)
}

func (p *StrategyParser) parseInZoneEvaluator(value any) (evaluator.Evaluator, error) {
	obj, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected object for 'InZone', got %T", value)
	}
	zoneName, ok := obj["Zone"].(string)
	if !ok {
		return nil, fmt.Errorf("expected 'Zone' key to be a string in 'InZone', got %T", obj["Zone"])
	}
	zone, ok := mtg.StringToZone(zoneName)
	if !ok {
		return nil, fmt.Errorf("invalid zone name in 'InZone': %s", zoneName)
	}
	cardsRaw, ok := obj["Cards"]
	if !ok {
		return nil, fmt.Errorf("missing 'Cards' in 'InZone'")
	}
	pred, err := p.parsePredicate(cardsRaw)
	if err != nil {
		return nil, fmt.Errorf("error parsing 'Cards' in 'InZone': %w", err)
	}
	return &node.InZone{
		Zone:  zone,
		Cards: pred,
	}, nil
}

func (p *StrategyParser) parseNumericOpValue(raw any) (evaluator.Op, any, error) {
	switch value := raw.(type) {
	case int:
		return evaluator.OpEqual, value, nil
	case string:
		return p.parseNumericOpValueString(value)
	case map[string]any:
		for k, v := range value {
			switch k {
			case "AtLeast":
				return evaluator.OpGreaterThanOrEqual, v, nil
			case "AtMost":
				return evaluator.OpLessThanOrEqual, v, nil
			case "Equal":
				return evaluator.OpEqual, v, nil
			case "Op":
				opString, ok := v.(string)
				if !ok {
					return "", nil, fmt.Errorf("expected 'Op' to be a string got %T", v)
				}
				opValue, ok := value["Value"].(int)
				if !ok {
					return "", nil, fmt.Errorf("expected 'Value' to be an int got %T", value["Value"])
				}
				op, err := p.parseOpString(opString)
				if err != nil {
					return "", nil, fmt.Errorf("error parsing operator %q: %w", opString, err)
				}
				return op, opValue, nil
			}
		}
		return "", nil, fmt.Errorf("no valid numeric operator found in %v", value)
	default:
		return "", nil, fmt.Errorf("expected int, string, or map for numeric value, got %T", raw)
	}
}

func (p *StrategyParser) parseNumericOpValueString(value string) (evaluator.Op, any, error) {
	parts := strings.Fields(value)
	s := strings.Join(parts, "")
	for _, op := range evaluator.Operators() {
		if strings.HasPrefix(s, string(op)) {
			valuePart := strings.TrimPrefix(s, string(op))
			value, err := strconv.Atoi(valuePart)
			if err != nil {
				return "", nil, fmt.Errorf("could not parse numeric value from string %q: %v", valuePart, err)
			}
			return op, value, nil
		}
	}
	return "", nil, fmt.Errorf("could not parse numeric operator from string %q", value)
}

func (p *StrategyParser) parseOpString(opString string) (evaluator.Op, error) {
	switch opString {
	case "==":
		return evaluator.OpEqual, nil
	case ">":
		return evaluator.OpGreaterThan, nil
	case ">=":
		return evaluator.OpGreaterThanOrEqual, nil
	case "<":
		return evaluator.OpLessThan, nil
	case "<=":
		return evaluator.OpLessThanOrEqual, nil
	case "!=":
		return evaluator.OpNotEqual, nil
	default:
		return "", fmt.Errorf("unknown operator %q", opString)
	}
}

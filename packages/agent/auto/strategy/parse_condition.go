package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/evaluator"
	"deckronomicon/packages/agent/auto/strategy/node"
	"deckronomicon/packages/game/mtg"
	"fmt"
)

func (p *StrategyParser) parseEvaluator(raw any) evaluator.Evaluator {
	switch val := raw.(type) {
	case map[string]any:
		var evaluators []evaluator.Evaluator
		for k, v := range val {
			switch k {
			case "And", "All", "AllOf":
				evaluators = append(evaluators, p.parseLogicalEvaluator("And", v))
			case "Or", "Any", "AnyOf":
				evaluators = append(evaluators, p.parseLogicalEvaluator("Or", v))
			case "Not":
				evaluators = append(evaluators, p.parseNotEvaluator(v))
			case "True", "Yes":
				evaluators = append(evaluators, &evaluator.True{})
			case "False", "No":
				evaluators = append(evaluators, &evaluator.False{})
			case "Step":
				evaluators = append(evaluators, p.parseStepEvaluator(v))
			case "Mode":
				evaluators = append(evaluators, p.parseModeEvaluator(v))
			case "InHand":
				evaluators = append(evaluators, p.parseInZoneAliasEvaluator("Hand", v))
			case "InZone":
				evaluators = append(evaluators, p.parseInZoneEvaluator(v))
			case "LandPlayedThisTurn":
				evaluators = append(evaluators, p.parseLandPlayedThisTurnEvaluator(v))
			default:
				p.errors.Add(fmt.Errorf("unknown key: %s", k))
				return nil
			}
		}
		if len(evaluators) == 0 {
			p.errors.Add(fmt.Errorf("no evaluators found in %v", val))
		}
		if len(evaluators) == 1 {
			return evaluators[0]
		}
		return &evaluator.And{Evaluators: evaluators}
	default:
		p.errors.Add(fmt.Errorf("expected object, got %T", raw))
		return nil
	}
}

func (p *StrategyParser) parseLandPlayedThisTurnEvaluator(value interface{}) evaluator.Evaluator {
	landPlayedThisTurn, ok := value.(bool)
	if !ok {
		p.errors.Add(fmt.Errorf("expected boolean for 'LandPlayedThisTurn', got %T", value))
		return nil
	}
	return &evaluator.LandPlayedThisTurn{
		LandPlayedThisTurn: landPlayedThisTurn,
	}
}

func (p *StrategyParser) parseNotEvaluator(value interface{}) evaluator.Evaluator {
	return &evaluator.Not{Evaluator: p.parseEvaluator(value)}
}

func (p *StrategyParser) parseLogicalEvaluator(op string, value interface{}) evaluator.Evaluator {
	items, ok := value.([]interface{})
	if !ok {
		p.errors.Add(fmt.Errorf("expected array for '%s', got %T", op, value))
		return nil
	}
	var evaluators []evaluator.Evaluator
	for _, item := range items {
		evaluators = append(evaluators, p.parseEvaluator(item))
	}
	switch op {
	case "And":
		return &evaluator.And{Evaluators: evaluators}
	case "Or":
		return &evaluator.Or{Evaluators: evaluators}
	}
	p.errors.Add(fmt.Errorf("unknown logical operator: %s", op))
	return nil
}

// TODO: Support a slice of steps
func (p *StrategyParser) parseStepEvaluator(value interface{}) evaluator.Evaluator {
	switch v := value.(type) {
	case string:
		step, err := mtg.StringToStep(v)
		if err != nil {
			p.errors.Add(fmt.Errorf("invalid step: %s", v))
			return nil
		}
		return &evaluator.Step{Step: step}
	case map[string]any:
		s, ok := v["Step"].(string)
		if !ok {
			p.errors.Add(fmt.Errorf("expected 'Step' to be a string, got %T", v["Step"]))
		}
		step, err := mtg.StringToStep(s)
		if err != nil {
			p.errors.Add(fmt.Errorf("invalid step: %s", v))
			return nil
		}
		return &evaluator.Step{Step: step}
	}
	p.errors.Add(fmt.Errorf("step must be a string or an object with a 'Step' key"))
	return nil
}

// TODO: Support a slice of modes
func (p *StrategyParser) parseModeEvaluator(value interface{}) evaluator.Evaluator {
	switch v := value.(type) {
	case string:
		return &evaluator.Mode{Mode: v}
	case map[string]any:
		if m, ok := v["Mode"].(string); ok {
			return &evaluator.Mode{Mode: m}
		}
	}
	p.errors.Add(fmt.Errorf("mode must be a string or an object with a 'Mode' key"))
	return nil
}

func (p *StrategyParser) parseInZoneAliasEvaluator(alias string, value interface{}) evaluator.Evaluator {
	zoneWrapper := map[string]any{
		"Zone":  alias,
		"Cards": value,
	}
	return p.parseInZoneEvaluator(zoneWrapper)
}

func (p *StrategyParser) parseInZoneEvaluator(value interface{}) evaluator.Evaluator {
	obj, ok := value.(map[string]any)
	if !ok {
		p.errors.Add(fmt.Errorf("expected object for 'InZone', got %T", value))
		return nil
	}
	zoneName, ok := obj["Zone"].(string)
	if !ok {
		p.errors.Add(fmt.Errorf("expected 'Zone' key to be a string in 'InZone', got %T", obj["Zone"]))
		return nil
	}
	zone, ok := mtg.StringToZone(zoneName)
	if !ok {
		p.errors.Add(fmt.Errorf("invalid zone name in 'InZone': %s", zoneName))
		return nil
	}
	cardsRaw, ok := obj["Cards"]
	if !ok {
		p.errors.Add(fmt.Errorf("missing cards in 'InZone'"))
		return nil
	}
	return &node.InZone{
		Zone:  zone,
		Cards: p.parsePredicate(cardsRaw),
	}
}

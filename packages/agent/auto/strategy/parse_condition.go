package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/evaluator"
	"deckronomicon/packages/agent/auto/strategy/matcher"
	"deckronomicon/packages/game/mtg"
	"fmt"
)

func (p *StrategyParser) parseEvaluator(raw any) evaluator.Evaluator {
	switch val := raw.(type) {
	case map[string]any:
		var evaluators []evaluator.Evaluator
		for k, v := range val {
			switch k {
			case "And":
				evaluators = append(evaluators, p.parseLogicalEvaluator("And", v))
			case "Or":
				evaluators = append(evaluators, p.parseLogicalEvaluator("Or", v))
			case "Not":
				evaluators = append(evaluators, p.parseNotEvaluator(v))
			case "True":
				evaluators = append(evaluators, &evaluator.True{})
			case "False":
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
	return &evaluator.InZone{
		Zone:  zone,
		Cards: p.parseMatcher(cardsRaw),
	}
}

func (p *StrategyParser) parseMatcher(data any) matcher.Matcher {
	switch node := data.(type) {
	case string:
		return &matcher.Name{Name: node}
	case []any:
		var matchers []matcher.Matcher
		for _, item := range node {
			matchers = append(matchers, p.parseMatcher(item))
		}
		return &matcher.And{Matchers: matchers}
	case map[string]any:
		if len(node) != 1 {
			// TODO: add context path to p.errors.
			p.errors.Add(fmt.Errorf("expected exactly one key, got %d", len(node)))
			return nil
		}
		for key, val := range node {
			switch key {
			// todo move to a separate function
			case "And":
				items, ok := val.([]any)
				if !ok {
					p.errors.Add(fmt.Errorf("expected list for 'and', got %T", val))
					return nil
				}
				var matchers []matcher.Matcher
				for _, item := range items {
					matchers = append(matchers, p.parseMatcher(item))
				}
				return &matcher.And{Matchers: matchers}
				// todo move to a separate function
			case "Or":
				items, ok := val.([]any)
				if !ok {
					p.errors.Add(fmt.Errorf("expected list for 'or', got %T", val))
					return nil
				}
				var matchers []matcher.Matcher
				for _, item := range items {
					matchers = append(matchers, p.parseMatcher(item))
				}
				return &matcher.Or{Matchers: matchers}
			case "Not":
				return p.parseNotMatcher(val)
			case "CardType":
				return p.parseCardTypeMatcher(val)
			case "Subtype":
				return p.parseSubtypeMatcher(val)
			case "Supertype":
			case "Color":
			case "ManaCost":
			case "Power":
			case "Toughness":
			default:
				p.errors.Add(fmt.Errorf("unknown key: %s", key))
				return nil
			}
		}
	default:
		p.errors.Add(fmt.Errorf("unexpected type: %T", data))
		return nil
	}
	// TODO handle the switch statement to always return so we don't have the compiler hit this
	panic("unreachable")
}

func (p *StrategyParser) parseNotMatcher(value any) matcher.Matcher {
	return &matcher.Not{Matcher: p.parseMatcher(value)}
}

func (p *StrategyParser) parseCardTypeMatcher(value any) matcher.Matcher {
	typeStr, ok := value.(string)
	if !ok {
		p.errors.Add(fmt.Errorf("expected string for 'CardType', got %T", value))
		return nil
	}
	cardType, ok := mtg.StringToCardType(typeStr)
	if !ok {
		p.errors.Add(fmt.Errorf("invalid card type: %s", typeStr))
		return nil
	}
	return &matcher.CardType{CardType: cardType}
}

func (p *StrategyParser) parseSubtypeMatcher(value any) matcher.Matcher {
	typeStr, ok := value.(string)
	if !ok {
		p.errors.Add(fmt.Errorf("expected string for 'Subtype', got %T", value))
		return nil
	}
	subtype, ok := mtg.StringToSubtype(typeStr)
	if !ok {
		p.errors.Add(fmt.Errorf("invalid subtype: %s", typeStr))
		return nil
	}
	return &matcher.Subtype{Subtype: subtype}
}

package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/evaluator"
	"fmt"
)

func (p *StrategyParser) parseManaPoolEvaluator(raw any) (evaluator.Evaluator, error) {
	switch value := raw.(type) {
	case map[string]any:
		var evaluators []evaluator.Evaluator
		for k, v := range value {
			switch k {
			case "Any":
				eval, err := p.parseManaPoolNumericEvaluator(v, evaluator.ManaPoolAny)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Any' in 'ManaPool' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "White":
				eval, err := p.parseManaPoolNumericEvaluator(v, evaluator.ManaPoolWhite)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'White' in 'ManaPool' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "Blue":
				eval, err := p.parseManaPoolNumericEvaluator(v, evaluator.ManaPoolBlue)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Blue' in 'ManaPool' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "Black":
				eval, err := p.parseManaPoolNumericEvaluator(v, evaluator.ManaPoolBlack)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Black' in 'ManaPool' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "Red":
				eval, err := p.parseManaPoolNumericEvaluator(v, evaluator.ManaPoolRed)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Red' in 'ManaPool' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "Green":
				eval, err := p.parseManaPoolNumericEvaluator(v, evaluator.ManaPoolGreen)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Green' in 'ManaPool' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "Colorless":
				eval, err := p.parseManaPoolNumericEvaluator(v, evaluator.ManaPoolColorless)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Colorless' in 'ManaPool' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			default:
				return nil, fmt.Errorf("unknown key %q in 'ManaPool' evaluator", k)
			}
		}
		if len(evaluators) == 0 {
			return evaluators[0], nil
		}
		return &evaluator.And{Evaluators: evaluators}, nil
	default:
		return nil, fmt.Errorf("expected  object for 'ManaPool', got %T", value)
	}
}

func (p *StrategyParser) parseManaPoolNumericEvaluator(value any, stat evaluator.ManaPoolStat) (evaluator.Evaluator, error) {
	op, amount, err := p.parseNumericOpValue(value)
	if err != nil {
		return nil, fmt.Errorf("error parsing numeric stat for %q in 'ManaPool': %w", stat, err)
	}
	return &evaluator.ManaPool{
		Stat:  stat,
		Op:    op,
		Value: amount,
	}, nil
}

func (p *StrategyParser) parseAvailableManaEvaluator(raw any) (evaluator.Evaluator, error) {
	switch value := raw.(type) {
	case map[string]any:
		var evaluators []evaluator.Evaluator
		for k, v := range value {
			switch k {
			case "Any":
				eval, err := p.parseAvailableManaNumericEvaluator(v, evaluator.AvailableManaAny)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Any' in 'AvailableMana' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "White":
				eval, err := p.parseAvailableManaNumericEvaluator(v, evaluator.AvailableManaWhite)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'White' in 'AvailableMana' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "Blue":
				eval, err := p.parseAvailableManaNumericEvaluator(v, evaluator.AvailableManaBlue)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Blue' in 'AvailableMana' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "Black":
				eval, err := p.parseAvailableManaNumericEvaluator(v, evaluator.AvailableManaBlack)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Black' in 'AvailableMana' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "Red":
				eval, err := p.parseAvailableManaNumericEvaluator(v, evaluator.AvailableManaRed)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Red' in 'AvailableMana' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "Green":
				eval, err := p.parseAvailableManaNumericEvaluator(v, evaluator.AvailableManaGreen)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Green' in 'AvailableMana' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "Colorless":
				eval, err := p.parseAvailableManaNumericEvaluator(v, evaluator.AvailableManaColorless)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Colorless' in 'AvailableMana' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			default:
				return nil, fmt.Errorf("unknown key %q in 'AvailableMana' evaluator", k)
			}
		}
		if len(evaluators) == 0 {
			return evaluators[0], nil
		}
		return &evaluator.And{Evaluators: evaluators}, nil
	default:
		return nil, fmt.Errorf("expected object for 'AvailableMana', got %T", value)
	}
}

func (p *StrategyParser) parseAvailableManaNumericEvaluator(value any, stat evaluator.AvailableManaStat) (evaluator.Evaluator, error) {
	op, amount, err := p.parseNumericOpValue(value)
	if err != nil {
		return nil, fmt.Errorf("error parsing numeric stat for %q in 'AvailableMana': %w", stat, err)
	}
	return &evaluator.AvailableMana{
		Stat:  stat,
		Op:    op,
		Value: amount,
	}, nil
}

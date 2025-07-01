package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/evaluator"
	"fmt"
)

func (p *StrategyParser) parseManaPoolEvaluator(raw any) evaluator.Evaluator {
	switch value := raw.(type) {
	case map[string]any:
		var evaluators []evaluator.Evaluator
		for k, v := range value {
			switch k {
			case "Any":
				evaluators = append(evaluators, p.parseManaPoolNumericEvaluator(v, evaluator.ManaPoolAny))
			case "White":
				evaluators = append(evaluators, p.parseManaPoolNumericEvaluator(v, evaluator.ManaPoolWhite))
			case "Blue":
				evaluators = append(evaluators, p.parseManaPoolNumericEvaluator(v, evaluator.ManaPoolBlue))
			case "Black":
				evaluators = append(evaluators, p.parseManaPoolNumericEvaluator(v, evaluator.ManaPoolBlack))
			case "Red":
				evaluators = append(evaluators, p.parseManaPoolNumericEvaluator(v, evaluator.ManaPoolRed))
			case "Green":
				evaluators = append(evaluators, p.parseManaPoolNumericEvaluator(v, evaluator.ManaPoolGreen))
			case "Colorless":
				evaluators = append(evaluators, p.parseManaPoolNumericEvaluator(v, evaluator.ManaPoolColorless))
			default:
				p.errors.Add(fmt.Errorf("unknown key '%s' in 'ManaPool'", k))
				return nil
			}
		}
		if len(evaluators) == 0 {
			return evaluators[0]
		}
		return &evaluator.And{Evaluators: evaluators}
	default:
		p.errors.Add(fmt.Errorf("expected  object for 'ManaPool', got %T", value))
		return nil
	}
}

func (p *StrategyParser) parseManaPoolNumericEvaluator(value any, stat evaluator.ManaPoolStat) evaluator.Evaluator {
	op, amount := p.parseNumericOpValue(value)
	return &evaluator.ManaPool{
		Stat:  stat,
		Op:    op,
		Value: amount,
	}
}

func (p *StrategyParser) parseAvailableManaEvaluator(raw any) evaluator.Evaluator {
	switch value := raw.(type) {
	case map[string]any:
		var evaluators []evaluator.Evaluator
		for k, v := range value {
			switch k {
			case "Any":
				evaluators = append(evaluators, p.parseAvailableManaNumericEvaluator(v, evaluator.AvailableManaAny))
			case "White":
				evaluators = append(evaluators, p.parseAvailableManaNumericEvaluator(v, evaluator.AvailableManaWhite))
			case "Blue":
				evaluators = append(evaluators, p.parseAvailableManaNumericEvaluator(v, evaluator.AvailableManaBlue))
			case "Black":
				evaluators = append(evaluators, p.parseAvailableManaNumericEvaluator(v, evaluator.AvailableManaBlack))
			case "Red":
				evaluators = append(evaluators, p.parseAvailableManaNumericEvaluator(v, evaluator.AvailableManaRed))
			case "Green":
				evaluators = append(evaluators, p.parseAvailableManaNumericEvaluator(v, evaluator.AvailableManaGreen))
			case "Colorless":
				evaluators = append(evaluators, p.parseAvailableManaNumericEvaluator(v, evaluator.AvailableManaColorless))
			default:
				p.errors.Add(fmt.Errorf("unknown key '%s' in 'AvailableMana'", k))
				return nil
			}
		}
		if len(evaluators) == 0 {
			return evaluators[0]
		}
		return &evaluator.And{Evaluators: evaluators}
	default:
		p.errors.Add(fmt.Errorf("expected object for 'AvailableMana', got %T", value))
		return nil
	}
}

func (p *StrategyParser) parseAvailableManaNumericEvaluator(value any, stat evaluator.AvailableManaStat) evaluator.Evaluator {
	op, amount := p.parseNumericOpValue(value)
	return &evaluator.AvailableMana{
		Stat:  stat,
		Op:    op,
		Value: amount,
	}
}

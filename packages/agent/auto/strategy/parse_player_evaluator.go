package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/evaluator"
	"fmt"
)

func (p *StrategyParser) parsePlayerEvaluator(raw any) (evaluator.Evaluator, error) {
	switch value := raw.(type) {
	case map[string]any:
		var evaluators []evaluator.Evaluator
		for k, v := range value {
			switch k {
			case "ID":
				eval, err := p.parsePlayerStringStatEvaluator(v, evaluator.PlayerStatID)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'ID' in 'Player' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "Name":
				eval, err := p.parsePlayerStringStatEvaluator(v, evaluator.PlayerStatName)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Name' in 'Player' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "Life":
				eval, err := p.parsePlayerNumericStatEvaluator(v, evaluator.PlayerStatLife)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Life' in 'Player' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "Turn":
				eval, err := p.parsePlayerNumericStatEvaluator(v, evaluator.PlayerStatTurn)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Turn' in 'Player' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "HandSize":
				eval, err := p.parsePlayerNumericStatEvaluator(v, evaluator.PlayerStatHandSize)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'HandSize' in 'Player' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "LibrarySize":
				eval, err := p.parsePlayerNumericStatEvaluator(v, evaluator.PlayerStatLibrarySize)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'LibrarySize' in 'Player' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			case "GraveyardSize":
				eval, err := p.parsePlayerNumericStatEvaluator(v, evaluator.PlayerStatGraveyardSize)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'GraveyardSize' in 'Player' evaluator: %w", err)
				}
				evaluators = append(evaluators, eval)
			default:
				return nil, fmt.Errorf("unknown key %q in 'Player' evaluator", k)
			}
		}
		if len(evaluators) == 0 {
			return evaluators[0], nil
		}
		return &evaluator.And{Evaluators: evaluators}, nil
	default:
		return nil, fmt.Errorf("expected  object for 'Player', got %T", value)
	}
}

func (p *StrategyParser) parsePlayerStringStatEvaluator(raw any, stat evaluator.PlayerStat) (evaluator.Evaluator, error) {
	value, ok := raw.(string)
	if !ok {
		return nil, fmt.Errorf("expected %q to be a string in 'Player', got %T", stat, raw)
	}
	return &evaluator.Player{
		Stat:  stat,
		Op:    evaluator.OpEqual,
		Value: value,
	}, nil
}

func (p *StrategyParser) parsePlayerNumericStatEvaluator(raw any, stat evaluator.PlayerStat) (evaluator.Evaluator, error) {
	op, value, err := p.parseNumericOpValue(raw)
	if err != nil {
		return nil, fmt.Errorf("error parsing numeric stat for %q in 'Player': %w", stat, err)
	}
	return &evaluator.Player{
		Stat:  stat,
		Op:    op,
		Value: value,
	}, nil
}

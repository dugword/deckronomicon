package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/evaluator"
	"fmt"
)

func (p *StrategyParser) parsePlayerEvaluator(raw any) evaluator.Evaluator {
	switch value := raw.(type) {
	case map[string]any:
		var evaluators []evaluator.Evaluator
		for k, v := range value {
			switch k {
			case "ID":
				evaluators = append(evaluators, p.parsePlayerStringStatEvaluator(v, evaluator.PlayerStatID))
			case "Name":
				evaluators = append(evaluators, p.parsePlayerStringStatEvaluator(v, evaluator.PlayerStatName))
			case "Life":
				evaluators = append(evaluators, p.parsePlayerNumericStatEvaluator(v, evaluator.PlayerStatLife))
			case "Turn":
				evaluators = append(evaluators, p.parsePlayerNumericStatEvaluator(v, evaluator.PlayerStatTurn))
			case "HandSize":
				evaluators = append(evaluators, p.parsePlayerNumericStatEvaluator(v, evaluator.PlayerStatHandSize))
			case "LibrarySize":
				evaluators = append(evaluators, p.parsePlayerNumericStatEvaluator(v, evaluator.PlayerStatLibrarySize))
			case "GraveyardSize":
				evaluators = append(evaluators, p.parsePlayerNumericStatEvaluator(v, evaluator.PlayerStatGraveyardSize))
			default:
				p.errors.Add(fmt.Errorf("unknown key '%s' in 'Player'", k))
				return nil
			}
		}
		if len(evaluators) == 0 {
			return evaluators[0]
		}
		return &evaluator.And{Evaluators: evaluators}
	default:
		p.errors.Add(fmt.Errorf("expected  object for 'Player', got %T", value))
		return nil
	}
}

func (p *StrategyParser) parsePlayerStringStatEvaluator(raw any, stat evaluator.PlayerStat) evaluator.Evaluator {
	value, ok := raw.(string)
	if !ok {
		p.errors.Add(fmt.Errorf("expected %q to be a string in 'Player', got %T", stat, raw))
		return nil
	}
	return &evaluator.Player{
		Stat:  stat,
		Op:    evaluator.OpEqual,
		Value: value,
	}
}

func (p *StrategyParser) parsePlayerNumericStatEvaluator(raw any, stat evaluator.PlayerStat) evaluator.Evaluator {
	op, value := p.parseNumericOpValue(raw)
	return &evaluator.Player{
		Stat:  stat,
		Op:    op,
		Value: value,
	}
}

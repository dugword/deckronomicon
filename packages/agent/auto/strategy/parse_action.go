package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/action"
	"deckronomicon/packages/agent/auto/strategy/predicate"
	"fmt"
)

func (p *StrategyParser) parseActionNode(raw any) (action.ActionNode, error) {
	switch val := raw.(type) {
	case string:
		act, ok := p.actions[val]
		if !ok {
			return nil, fmt.Errorf("unknown action %q", val)
		}
		return act, nil
	case map[string]any:
		// TODO: Make this check for 1, keeping zero ok for now
		// for some testing
		if len(val) > 1 {
			return nil, fmt.Errorf("expected a single action key, got %d keys", len(val))
		}
		for k, v := range val {
			switch k {
			case "Play":
				return p.parsePlayAction(v)
			case "Cast":
				return p.parseCastAction(v)
			case "Log":
				return p.parseLogMessageAction(v)
			case "Emit":
				return p.parseEmitMetricAction(v)
			case "Concede":
				return p.parseConcedeAction(v)
			case "Activate", "Tap":
				return p.parseActivateAction(v)
			case "Pass":
				return p.parsePassAction(v)
			}
			return nil, fmt.Errorf("unknown action key %q", k)
		}
		return nil, fmt.Errorf("no action key found in %v", val)
	default:
		return nil, fmt.Errorf("expected action object, got %T", raw)
	}
}

func (p *StrategyParser) parseEmitMetricAction(value any) (action.ActionNode, error) {
	switch val := value.(type) {
	case map[string]any:
		name, ok := val["Metric"].(string)
		if !ok {
			return nil, fmt.Errorf("expected 'Metric' key to be a string in 'Emit' action, got %T", val["Metric"])
		}
		// TODO: Should this default to 1 if not specified?
		value, ok := val["Value"].(int)
		if !ok {
			return nil, fmt.Errorf("expected 'Value' key to be an int in 'Emit' action, got %T", val["Value"])
		}
		return &action.EmitMetricActionNode{Name: name, Value: value}, nil
	default:
		return nil, fmt.Errorf("expected string or object for 'Emit' action, got %T", value)
	}
}

func (p *StrategyParser) parseLogMessageAction(value any) (action.ActionNode, error) {
	switch val := value.(type) {
	case string:
		return &action.LogMessageActionNode{Message: val}, nil
	case map[string]any:
		message, ok := val["Message"].(string)
		if !ok {
			return nil, fmt.Errorf("expected 'Message' key to be a string in 'Log' action, got %T", val["Message"])
		}
		return &action.LogMessageActionNode{Message: message}, nil
	default:
		return nil, fmt.Errorf("expected string or object for 'Log' action, got %T", value)
	}
}

func (p *StrategyParser) parsePlayAction(value any) (action.ActionNode, error) {
	switch val := value.(type) {
	case string:
		pred, err := p.parsePredicate(val)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'Play' action: %w", err)
		}
		return &action.PlayLandCardActionNode{Cards: pred}, nil
	case []any:
		var predicates []predicate.Predicate
		for _, item := range val {
			pred, err := p.parsePredicate(item)
			if err != nil {
				return nil, fmt.Errorf("error parsing 'Play' action item %v: %w)", item, err)
			}
			predicates = append(predicates, pred)
		}
		return &action.PlayLandCardActionNode{Cards: &predicate.Or{Predicates: predicates}}, nil
	case map[string]any:
		cardName, ok := val["Card"].(string)
		if !ok {
			return nil, fmt.Errorf("expected 'Card' key to be a string in 'Play' action, got %T", val["Card"])
		}
		pred, err := p.parsePredicate(cardName)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'Card' in 'Play' action: %w", err)
		}
		return &action.PlayLandCardActionNode{Cards: pred}, nil
	default:
		return nil, fmt.Errorf("expected string or object for 'Play' action, got %T", value)
	}
}

// TODO: Look this over, wrote it quickly late at night
func (p *StrategyParser) parseCastAction(value any) (action.ActionNode, error) {
	switch val := value.(type) {
	case string:
		pred, err := p.parsePredicate(val)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'Cast' action: %w", err)
		}
		return &action.CastSpellActionNode{Cards: pred}, nil
	case []any:
		var predicates []predicate.Predicate
		for _, item := range val {
			pred, err := p.parsePredicate(item)
			if err != nil {
				return nil, fmt.Errorf("error parsing 'Cast' action item %v: %w", item, err)
			}
			predicates = append(predicates, pred)
		}
		return &action.CastSpellActionNode{Cards: &predicate.Or{Predicates: predicates}}, nil
	case map[string]any:
		cardName, ok := val["Card"].(string)
		if !ok {
			return nil, fmt.Errorf("expected 'Card' key to be a string in 'Cast' action, got %T", val["Card"])
		}
		pred, err := p.parsePredicate(cardName)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'Card' in 'Cast' action: %w", err)
		}
		var additionalCost predicate.Selector
		if additionalCostRaw, ok := val["AdditionalCost"].(map[string]any); ok {
			additionalCost, err = p.parseAdditionalCost(additionalCostRaw)
			if err != nil {
				return nil, fmt.Errorf("error parsing 'AdditionalCost' in 'Cast' action: %w", err)
			}
		}
		return &action.CastSpellActionNode{
			Cards:          pred,
			AdditionalCost: additionalCost,
		}, nil
	default:
		return nil, fmt.Errorf("expected string or object for 'Cast' action, got %T", value)
	}
}

func (p *StrategyParser) parseAdditionalCost(additionalCost map[string]any) (predicate.Predicate, error) {
	for k, v := range additionalCost {
		switch k {
		case "Discard":
			return p.parsePredicate(v)
		}
	}
	return nil, fmt.Errorf("unknown additional cost key")
}

func (p *StrategyParser) parseConcedeAction(value any) (action.ActionNode, error) {
	if value != nil {
		return nil, fmt.Errorf("expected 'Concede' action to have no value, got %T", value)
	}
	return &action.ConcedeActionNode{}, nil
}

func (p *StrategyParser) parsePassAction(value any) (action.ActionNode, error) {
	if value != nil {
		return nil, fmt.Errorf("expected 'Pass' action to have no value, got %T", value)
	}
	return &action.PassPriorityActionNode{}, nil
}

func (p *StrategyParser) parseActivateAction(raw any) (action.ActionNode, error) {
	switch val := raw.(type) {
	case string:
		pred, err := p.parsePredicate(val)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'Activate' action: %w", err)
		}
		return &action.ActivateActionNode{Cards: pred}, nil
	case []any:
		var predicates []predicate.Predicate
		for _, item := range val {
			pred, err := p.parsePredicate(item)
			if err != nil {
				return nil, fmt.Errorf("error parsing 'Cast' action item %v: %w", item, err)
			}
			predicates = append(predicates, pred)
		}
		return &action.ActivateActionNode{Cards: &predicate.Or{Predicates: predicates}}, nil
	case map[string]any:
		cardName, ok := val["Card"].(string)
		if !ok {
			return nil, fmt.Errorf("expected 'Card' key to be a string in 'Cast' action, got %T", val["Card"])
		}
		pred, err := p.parsePredicate(cardName)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'Card' in 'Cast' action: %w", err)
		}
		return &action.ActivateActionNode{
			Cards: pred,
		}, nil
	default:
		return nil, fmt.Errorf("expected string or object for 'Cast' action, got %T", val)
	}
}

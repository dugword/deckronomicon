package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/action"
	"deckronomicon/packages/agent/auto/strategy/predicate"
	"fmt"
)

func (p *StrategyParser) parseActionNode(raw any) action.ActionNode {
	switch val := raw.(type) {
	case map[string]any:
		// TODO: Make this check for 1, keeping zero ok for now
		// for some testing
		if len(val) > 1 {
			p.errors.Add(fmt.Errorf("expected a single action key, got %d keys", len(val)))
			return nil
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
			default:
				p.errors.Add(fmt.Errorf("unknown action key: %s", k))
				return nil
			}
		}
		// TODO: Think through this, maybe configure pass or fail as default?
		return &action.PassPriorityActionNode{} // No action specified, default to pass
	default:
		p.errors.Add(fmt.Errorf("expected action object, got %T", raw))
		return nil
	}
}

func (p *StrategyParser) parseEmitMetricAction(value any) action.ActionNode {
	fmt.Println("Parsing Emit action with value:", value)
	switch val := value.(type) {
	case map[string]any:
		name, ok := val["Metric"].(string)
		if !ok {
			p.errors.Add(fmt.Errorf("expected 'Metroc' key to be a string in 'Emit' action, got %T", val["Metric"]))
			return nil
		}
		// TODO: Should this default to 1 if not specified?
		value, ok := val["Value"].(int)
		if !ok {
			p.errors.Add(fmt.Errorf("expected 'Value' key to be an int in 'Emit' action, got %T", val["Value"]))
			return nil
		}
		return &action.EmitMetricActionNode{Name: name, Value: value}
	default:
		p.errors.Add(fmt.Errorf("expected string or object for 'Emit' action, got %T", value))
		return nil
	}
}

func (p *StrategyParser) parseLogMessageAction(value any) action.ActionNode {
	fmt.Println("Parsing Log action with value:", value)
	switch val := value.(type) {
	case string:
		return &action.LogMessageActionNode{Message: val}
	case map[string]any:
		message, ok := val["Message"].(string)
		if !ok {
			p.errors.Add(fmt.Errorf("expected 'Message' key to be a string in 'Log' action, got %T", val["Message"]))
			return nil
		}
		return &action.LogMessageActionNode{Message: message}
	default:
		p.errors.Add(fmt.Errorf("expected string or object for 'Log' action, got %T", value))
		return nil
	}
}

func (p *StrategyParser) parsePlayAction(value any) action.ActionNode {
	fmt.Println("Parsing Play action with value:", value)
	switch val := value.(type) {
	case string:
		predicate := p.parsePredicate(val)
		return &action.PlayLandCardActionNode{Cards: predicate}
	case []any:
		var predicates []predicate.Predicate
		for _, item := range val {
			predicate := p.parsePredicate(item)
			predicates = append(predicates, predicate)
		}
		return &action.PlayLandCardActionNode{Cards: &predicate.Or{Predicates: predicates}}
	case map[string]any:
		cardName, ok := val["Card"].(string)
		if !ok {
			p.errors.Add(fmt.Errorf("expected 'Card' key to be a string in 'Play' action, got %T", val["Card"]))
			return nil
		}
		predicate := p.parsePredicate(cardName)
		return &action.PlayLandCardActionNode{Cards: predicate}
	default:
		p.errors.Add(fmt.Errorf("expected string or object for 'Play' action, got %T", value))
		return nil
	}
}

// TODO: Look this over, wrote it quickly late at night
func (p *StrategyParser) parseCastAction(value any) action.ActionNode {
	fmt.Println("Parsing Cast action with value:", value)
	switch val := value.(type) {
	case string:
		predicate := p.parsePredicate(val)
		return &action.CastSpellActionNode{Cards: predicate}
	case []any:
		var predicates []predicate.Predicate
		for _, item := range val {
			predicate := p.parsePredicate(item)
			predicates = append(predicates, predicate)
		}
		return &action.CastSpellActionNode{Cards: &predicate.Or{Predicates: predicates}}
	case map[string]any:
		cardName, ok := val["Card"].(string)
		if !ok {
			p.errors.Add(fmt.Errorf("expected 'Card' key to be a string in 'Cast' action, got %T", val["Card"]))
			return nil
		}
		predicate := p.parsePredicate(cardName)
		if additionalCostRaw, ok := val["AdditionalCost"].(map[string]any); ok {
			cardsRaw, ok := additionalCostRaw["Cards"]
			if !ok {
				p.errors.Add(fmt.Errorf("missing cards in 'AdditionalCost'"))
				return nil
			}
			additionalCost := p.parsePredicate(cardsRaw)
			return &action.CastSpellActionNode{
				Cards:          predicate,
				AdditionalCost: additionalCost,
			}
		}
		return &action.CastSpellActionNode{
			Cards: predicate,
		}
	default:
		p.errors.Add(fmt.Errorf("expected string or object for 'Cast' action, got %T", value))
		return nil
	}
}

func (p *StrategyParser) parseConcedeAction(value any) action.ActionNode {
	if value != nil {
		p.errors.Add(fmt.Errorf("expected 'Concede' action to have no value, got %T", value))
		return nil
	}
	return &action.ConcedeActionNode{}
}

func (p *StrategyParser) parsePassAction(value any) action.ActionNode {
	if value != nil {
		p.errors.Add(fmt.Errorf("expected 'Pass' action to have no value, got %T", value))
	}
	return &action.PassPriorityActionNode{}
}

// TODO: This is broken
func (p *StrategyParser) parseActivateAction(_ any) action.ActionNode {
	/*
		switch val := raw.(type) {
		case string:
			return &action.ActivateActionNode{AbilityInZone: gob.AbilityInZone{}}
		case []any:
			var abilities []string
			for _, item := range val {
				ability, ok := item.(string)
				if !ok {
					p.errors.Add(fmt.Errorf("expected string in 'Activate' action array, got %T", item))
					return nil
				}
				abilities = append(abilities, ability)
			}
			return &action.ActivateActionNode{AbilityInZone: gob.AbilityInZone{}}
		case map[string]any:
			_, ok := val["Ability"].(string)
			if !ok {
				p.errors.Add(fmt.Errorf("expected 'Ability' key to be a string in 'Activate' action, got %T", val["Ability"]))
				return nil
			}
			return &action.ActivateActionNode{AbilityInZone: gob.AbilityInZone{}}
		default:
			p.errors.Add(fmt.Errorf("expected string or object for 'Activate' action, got %T", raw))
			return nil
		}
	*/
	return nil
}

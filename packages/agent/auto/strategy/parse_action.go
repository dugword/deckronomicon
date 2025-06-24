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
			case "Play", "Cast":
				return p.parsePlayAction(v)
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

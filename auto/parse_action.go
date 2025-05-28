package auto

import "fmt"

func (p *StrategyParser) parseActionNode(raw any) ActionNode {
	switch val := raw.(type) {
	case map[string]any:
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
		return &PassAction{} // No action specified, default to pass
	default:
		p.errors.Add(fmt.Errorf("expected action object, got %T", raw))
		return nil
	}
}

func (p *StrategyParser) parsePlayAction(value any) ActionNode {
	switch val := value.(type) {
	case string:
		return &PlayAction{CardNames: []string{val}}
	case []any:
		var cards []string
		for _, item := range val {
			card, ok := item.(string)
			if !ok {
				p.errors.Add(fmt.Errorf("expected string in 'Play' action array, got %T", item))
				return nil
			}
			cards = append(cards, card)
		}
		return &PlayAction{CardNames: cards}
	case map[string]any:
		cardName, ok := val["Card"].(string)
		if !ok {
			p.errors.Add(fmt.Errorf("expected 'Card' key to be a string in 'Play' action, got %T", val["Card"]))
			return nil
		}
		return &PlayAction{CardNames: []string{cardName}}
	default:
		p.errors.Add(fmt.Errorf("expected string or object for 'Play' action, got %T", value))
		return nil
	}
}

func (p *StrategyParser) parseConcedeAction(value any) ActionNode {
	if value != nil {
		p.errors.Add(fmt.Errorf("expected 'Concede' action to have no value, got %T", value))
		return nil
	}
	return &ConcedeAction{}
}

func (p *StrategyParser) parsePassAction(value any) ActionNode {
	if value != nil {
		p.errors.Add(fmt.Errorf("expected 'Pass' action to have no value, got %T", value))
	}
	return &PassAction{}
}

func (p *StrategyParser) parseActivateAction(raw any) ActionNode {
	switch val := raw.(type) {
	case string:
		return &ActivateAction{AbilityNames: []string{val}}
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
		return &ActivateAction{AbilityNames: abilities}
	case map[string]any:
		ability, ok := val["Ability"].(string)
		if !ok {
			p.errors.Add(fmt.Errorf("expected 'Ability' key to be a string in 'Activate' action, got %T", val["Ability"]))
			return nil
		}
		return &ActivateAction{AbilityNames: []string{ability}}
	default:
		p.errors.Add(fmt.Errorf("expected string or object for 'Activate' action, got %T", raw))
		return nil
	}
}

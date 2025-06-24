package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/action"
	"fmt"
	"strings"
)

func (p *StrategyParser) parseActionNode(raw any) action.ActionNode {
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
		return &action.PassPriorityActionNode{} // No action specified, default to pass
	default:
		p.errors.Add(fmt.Errorf("expected action object, got %T", raw))
		return nil
	}
}

func (p *StrategyParser) parseGroupThingy(input string) []string {
	if !strings.HasPrefix(input, "$") {
		fmt.Println("Missing prefix")
		p.errors.Add(fmt.Errorf("group name must start with '$', got '%s'", input))
		return nil
	}
	groupName := strings.TrimPrefix(input, "$")
	group, ok := p.groups[groupName]
	if !ok {
		fmt.Println("Group not found", groupName)
		p.errors.Add(fmt.Errorf("group '%s' not found", groupName))
		return nil
	}
	var cardNames []string
	for _, item := range group {
		cardName, ok := item.(string)
		if !ok {
			fmt.Println("Expected string in group", groupName, "but got", item)
			p.errors.Add(fmt.Errorf("expected string in group '%s', got %T", groupName, item))
			return nil
		}
		cardNames = append(cardNames, cardName)
	}
	fmt.Println("HELLO FROM PARSE GROUP THINGY", cardNames)
	return cardNames
}

func (p *StrategyParser) parsePlayAction(value any) action.ActionNode {
	fmt.Println("Parsing Play action with value:", value)
	switch val := value.(type) {
	case string:
		if strings.HasPrefix(val, "$") {
			cardNames := p.parseGroupThingy(val)
			return &action.PlayLandCardActionNode{CardNames: cardNames}
		}
		return &action.PlayLandCardActionNode{CardNames: []string{val}}
	case []any:
		var cardNames []string
		for _, item := range val {
			card, ok := item.(string)
			if !ok {
				p.errors.Add(fmt.Errorf("expected string in 'Play' action array, got %T", item))
				return nil
			}
			if strings.HasPrefix(card, "$") {
				cardNamesFromGroup := p.parseGroupThingy(card)
				if cardNamesFromGroup != nil {
					cardNames = append(cardNames, cardNamesFromGroup...)
				}
			}
			cardNames = append(cardNames, card)
		}
		return &action.PlayLandCardActionNode{CardNames: cardNames}
	case map[string]any:
		_, ok := val["Card"].(string)
		if !ok {
			p.errors.Add(fmt.Errorf("expected 'Card' key to be a string in 'Play' action, got %T", val["Card"]))
			return nil
		}
		return &action.PlayLandCardActionNode{CardNames: []string{val["Card"].(string)}}
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

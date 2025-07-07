package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/action"
	"deckronomicon/packages/agent/auto/strategy/evaluator"
	"deckronomicon/packages/agent/auto/strategy/predicate"
	"errors"
	"fmt"
)

// region Node Parsers

func (p *StrategyParser) ParseChoiceNode(raw any) (*Choice, error) {
	switch node := raw.(type) {
	case string:
		choice, ok := p.choices[node]
		if !ok {
			return nil, fmt.Errorf("choice %q not found in strategy", node)
		}
		return choice, nil
	case map[string]any:
		var choice Choice
		for key, value := range node {
			switch key {
			case "Name":
				choice.Name = value.(string)
			case "Description":
				choice.Description = value.(string)
			case "Source":
				choice.Source = value.(string)
				if choice.Source == "" {
					return nil, errors.New("choice has no 'Source' defined")
				}
			case "When":
				eval, err := p.parseEvaluator(value)
				if err != nil {
					return nil, fmt.Errorf("failed to parse 'When' condition: %w", err)
				}
				choice.When = eval
			case "Choose":
				selector, err := p.parsePredicate(value)
				if err != nil {
					return nil, fmt.Errorf("failed to parse 'Choose' selector: %w", err)
				}
				choice.Choose = selector
			}
		}
		return &choice, nil
	default:
		return nil, fmt.Errorf("invalid choice node type: %T", node)
	}
}

func (p *StrategyParser) ParseGroupNode(raw any) (*Group, error) {
	switch node := raw.(type) {
	case map[string]any:
		var group Group
		for key, value := range node {
			switch key {
			case "Name":
				group.Name = value.(string)
				if group.Name == "" {
					return nil, fmt.Errorf("group has no name defined")
				}
			case "Description":
				group.Description = value.(string)
			case "Members":
				members, ok := value.([]any)
				if !ok {
					return nil, fmt.Errorf("group has invalid 'Members' type: %T", value)
				}
				if len(members) == 0 {
					return nil, fmt.Errorf("group has no members defined")
				}
				for _, member := range members {
					memberStr, ok := member.(string)
					if !ok {
						return nil, fmt.Errorf("group member is not a string: %T", member)
					}
					if memberStr == "" {
						return nil, fmt.Errorf("group member cannot be an empty string")
					}
					group.Members = append(group.Members, memberStr)
				}
				if len(group.Members) == 0 {
					return nil, fmt.Errorf("group has no members defined")
				}
			default:
				return nil, fmt.Errorf("unknown key %q in group node", key)
			}
		}
		return &group, nil
	default:
		return nil, fmt.Errorf("invalid group node type: %T", node)
	}
}

func (p *StrategyParser) ParseModeNode(raw any) (*Mode, error) {
	switch node := raw.(type) {
	case string:
		mode, ok := p.modes[node]
		if !ok {
			return nil, fmt.Errorf("mode %q not found in strategy", node)
		}
		return mode, nil
	case map[string]any:
		var mode Mode
		for key, value := range node {
			switch key {
			case "Name":
				mode.Name = value.(string)
			case "Description":
				mode.Description = value.(string)
			case "When":
				eval, err := p.parseEvaluator(value)
				if err != nil {
					return nil, fmt.Errorf("failed to parse 'When' evaluator: %w", err)
				}
				mode.When = eval
			}
		}
		return &mode, nil
	default:
		return nil, fmt.Errorf("invalid mode node type: %T", node)
	}
}

func (p *StrategyParser) ParseNamedActionNode(raw any) (action.ActionNode, error) {
	switch node := raw.(type) {
	case string:
		actionNode, ok := p.actions[node]
		if !ok {
			return nil, fmt.Errorf("action %q not found in strategy", node)
		}
		return actionNode, nil
	case map[string]any:
		for key, value := range node {
			switch key {
			case "Then":
				return p.parseActionNode(value)
			}
		}
		return nil, fmt.Errorf("named action node has no 'Action' key: %v", node)
	default:
		return nil, fmt.Errorf("invalid named action node type: %T", node)
	}
}

func (p *StrategyParser) ParseNamedConditionNode(raw any) (evaluator.Evaluator, error) {
	switch node := raw.(type) {
	case string:
		condition, ok := p.conditions[node]
		if !ok {
			return nil, fmt.Errorf("condition %q not found in strategy", node)
		}
		return condition, nil
	case map[string]any:
		for key, value := range node {
			switch key {
			case "When":
				return p.parseEvaluator(value)
			}
		}
		return nil, fmt.Errorf("named condition node has no 'When' key: %v", node)
	default:
		return nil, fmt.Errorf("invalid named condition node type: %T", node)
	}
}

func (p *StrategyParser) ParseNamedSelectorNode(raw any) (predicate.Selector, error) {
	switch node := raw.(type) {
	case string:
		selector, ok := p.selectors[node]
		if !ok {
			return nil, fmt.Errorf("selector %q not found in strategy", node)
		}
		return selector, nil
	case map[string]any:
		for key, value := range node {
			switch key {
			case "Choose":
				return p.parsePredicate(value)
			}
		}
		return nil, fmt.Errorf("named selector node has no 'Choose' key: %v", node)
	default:
		return nil, fmt.Errorf("invalid named selector node type: %T", node)
	}
}

func (p *StrategyParser) ParseRuleNode(raw any) (*Rule, error) {
	switch node := raw.(type) {
	case string:
		rule, ok := p.rules[node]
		if !ok {
			return nil, fmt.Errorf("rule %q not found in strategy", node)
		}
		return rule, nil
	case map[string]any:
		var rule Rule
		for key, value := range node {
			switch key {
			case "Name":
				rule.Name = value.(string)
			case "Description":
				rule.Description = value.(string)
			case "When":
				eval, err := p.parseEvaluator(value)
				if err != nil {
					return nil, fmt.Errorf("failed to parse 'When' condition: %w", err)
				}
				rule.When = eval
			case "Then":
				eval, err := p.parseActionNode(value)
				if err != nil {
					return nil, fmt.Errorf("failed to parse 'Then' action: %w", err)
				}
				rule.Then = eval
			}
		}
		return &rule, nil
	default:
		return nil, fmt.Errorf("invalid rule node type: %T", node)
	}
}

// endregion

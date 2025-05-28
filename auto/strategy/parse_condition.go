package strategy

import (
	"fmt"
)

// TODO Maybe something like this
// type ErrorCollector struct {
// 	errors []error
// }
// func shouldConditionNode(ConditionNode, error) ConditionNode {
// 	if err != nil {
// 		errors = append(errors, err)
// 		return nil
// 	}
// 	return conditionNode
// }

func (p *StrategyParser) parseConditionNode(raw any) ConditionNode {
	switch val := raw.(type) {
	case map[string]any:
		var conditionNodes []ConditionNode
		for k, v := range val {
			switch k {
			case "And":
				conditionNodes = append(conditionNodes, p.parseLogicalNode("And", v))
			case "Or":
				conditionNodes = append(conditionNodes, p.parseLogicalNode("Or", v))
			case "Not":
				conditionNodes = append(conditionNodes, p.parseNotConditionNode(v))
			//case "LifeTotal":
			//return parseNumericCondition("LifeTotal", v)
			case "Step":
				conditionNodes = append(conditionNodes, p.parseStepCondition(v))
			case "Mode":
				conditionNodes = append(conditionNodes, p.parseModeCondition(v))
			case "InHand":
				conditionNodes = append(conditionNodes, p.parseInZoneAliasCondition("Hand", v))
			case "InZone":
				conditionNodes = append(conditionNodes, p.parseInZoneCondition(v))
			case "LandDrop":
				conditionNodes = append(conditionNodes, p.parseLandDropCondition(v))
			default:
				p.errors.Add(fmt.Errorf("unknown condition key: %s", k))
				return nil
			}
		}
		if len(conditionNodes) == 0 {
			p.errors.Add(fmt.Errorf("no valid conditions found in %v", val))
		}
		// If there's only one condition, return it directly
		if len(conditionNodes) == 1 {
			return conditionNodes[0]
		}
		return &AndCondition{Conditions: conditionNodes} // If multiple conditions, return an AND condition
	default:
		p.errors.Add(fmt.Errorf("expected condition object, got %T", raw))
		return nil
	}
}

func (p *StrategyParser) parseLandDropCondition(value interface{}) ConditionNode {
	landDrop, ok := value.(bool)
	if !ok {
		p.errors.Add(fmt.Errorf("expected boolean for 'LandDrop' condition, got %T", value))
		return nil
	}
	return &LandDropCondition{
		LandDrop: landDrop,
	}
}

func (p *StrategyParser) parseNotConditionNode(value interface{}) ConditionNode {
	return &NotCondition{Condition: p.parseConditionNode(value)}
}

func (p *StrategyParser) parseLogicalNode(op string, value interface{}) ConditionNode {
	items, ok := value.([]interface{})
	if !ok {
		p.errors.Add(fmt.Errorf("expected array for '%s' condition, got %T", op, value))
		return nil
	}
	var conditions []ConditionNode
	for _, item := range items {
		conditions = append(conditions, p.parseConditionNode(item))
	}
	switch op {
	case "And":
		return &AndCondition{Conditions: conditions}
	case "Or":
		return &OrCondition{Conditions: conditions}
	}
	p.errors.Add(fmt.Errorf("unknown logical operator: %s", op))
	return nil
}

// TODO: Support a slice of steps
func (p *StrategyParser) parseStepCondition(value interface{}) ConditionNode {
	switch v := value.(type) {
	case string:
		return &StepCondition{Step: v}
	case map[string]any:
		if s, ok := v["Step"].(string); ok {
			return &StepCondition{Step: s}
		}
	}
	p.errors.Add(fmt.Errorf("step condition must be a string or an object with a 'Step' key"))
	return nil
}

// TODO: Support a slice of modes
func (p *StrategyParser) parseModeCondition(value interface{}) ConditionNode {
	switch v := value.(type) {
	case string:
		return &ModeCondition{Mode: v}
	case map[string]any:
		if m, ok := v["Mode"].(string); ok {
			return &ModeCondition{Mode: m}
		}
	}
	p.errors.Add(fmt.Errorf("mode condition must be a string or an object with a 'Mode' key"))
	return nil
}

func (p *StrategyParser) parseInZoneAliasCondition(alias string, value interface{}) ConditionNode {
	zoneWrapper := map[string]any{
		"Zone":  alias,
		"Cards": value,
	}
	return p.parseInZoneCondition(zoneWrapper)
}

func (p *StrategyParser) parseInZoneCondition(value interface{}) ConditionNode {
	obj, ok := value.(map[string]any)
	if !ok {
		p.errors.Add(fmt.Errorf("expected object for 'InZone' condition, got %T", value))
		return nil
	}
	zoneName, ok := obj["Zone"].(string)
	if !ok {
		p.errors.Add(fmt.Errorf("expected 'Zone' key to be a string in 'InZone' condition, got %T", obj["Zone"]))
		return nil
	}
	cardsRaw, ok := obj["Cards"]
	if !ok {
		p.errors.Add(fmt.Errorf("missing cards in 'InZone' condition"))
		return nil
	}
	return &InZoneCondition{
		Zone:  zoneName,
		Cards: p.parseCardConditionNode(cardsRaw),
	}
}

// TODO not used, but I like beinble able to have a card be an object or a string
// objects would let me do {"CardType": "Land"} instead of just Names.
func (p *StrategyParser) parseCardCondition(value interface{}) CardCondition {
	switch v := value.(type) {
	case string:
		return &CardNameCondition{Name: v}
	case map[string]any:
		name, ok := v["Name"].(string)
		if !ok {
			p.errors.Add(fmt.Errorf("expected 'Name' key to be a string in card condition, got %T", v["Name"]))
			return nil
		}
		return &CardNameCondition{Name: name}
	}
	p.errors.Add(fmt.Errorf("expeacted card condition to be a string or an object, got %T", value))
	return nil
}

// Parses a generic card condition from JSON
func (p *StrategyParser) parseCardConditionNode(data any) CardCondition {
	switch node := data.(type) {
	case string:
		return NewNameConditionOrGroupRef(node)
	case []any:
		var conds []CardCondition
		for _, item := range node {
			conds = append(conds, p.parseCardConditionNode(item))
		}
		return &AndCardCondition{Conditions: conds}
	case map[string]any:
		if len(node) != 1 {
			p.errors.Add(fmt.Errorf("invalid card condition node, expected exactly one key, got %d", len(node)))
			return nil
		}
		for key, val := range node {
			switch key {
			// todo move to a separate function
			case "And":
				items, ok := val.([]any)
				if !ok {
					p.errors.Add(fmt.Errorf("expected list for 'and', got %T", val))
					return nil
				}
				var conds []CardCondition
				for _, item := range items {
					conds = append(conds, p.parseCardConditionNode(item))
				}
				return &AndCardCondition{Conditions: conds}
				// todo move to a separate function
			case "Or":
				items, ok := val.([]any)
				if !ok {
					p.errors.Add(fmt.Errorf("expected list for 'or', got %T", val))
					return nil
				}
				var conds []CardCondition
				for _, item := range items {
					conds = append(conds, p.parseCardConditionNode(item))
				}
				return &OrCardCondition{Conditions: conds}
			case "Not":
				return p.parseCardNotConditionNode(val)
			case "CardType":
				return p.parseCardTypeCondition(val)
			case "Subtype":
				return p.parseCardSubtypeCondition(val)
			case "Supertype":
			case "Color":
			case "ManaCost":
			case "Power":
			case "Toughness":
			default:
				p.errors.Add(fmt.Errorf("unknown card condition key: %s", key))
				return nil
			}
		}
	default:
		p.errors.Add(fmt.Errorf("unexpected card condition type: %T", data))
		return nil
	}
	// TODO handle the switch statement to always return so we don't have the compiler hit this
	panic("unreachable")
}

func (p *StrategyParser) parseCardNotConditionNode(value any) CardCondition {
	return &NotCardCondition{Condition: p.parseCardConditionNode(value)}
}

func (p *StrategyParser) parseCardTypeCondition(value any) CardCondition {
	typeStr, ok := value.(string)
	if !ok {
		p.errors.Add(fmt.Errorf("expected string for 'CardType', got %T", value))
		return nil
	}
	return &CardTypeCondition{CardType: typeStr}
}

func (p *StrategyParser) parseCardSubtypeCondition(value any) CardCondition {
	typeStr, ok := value.(string)
	if !ok {
		p.errors.Add(fmt.Errorf("expected string for 'Subtype', got %T", value))
		return nil
	}
	return &CardSubtypeCondition{Subtype: typeStr}
}

/*
func NewNameConditionOrGroupRef(input string) CardCondition {
	if strings.HasPrefix(input, "$") {
		// It’s a group reference like "$ComboPiece"
		return &GroupRefCondition{
			GroupName: input[1:], // strip the "$"
		}
	}
	// It’s a literal card name
	return &NameMatchCondition{
		Name: input,
	}
}

type GroupRefCondition struct {
	GroupName string
}

func (g *GroupRefCondition) Matches(card *Card, strategy *Strategy) bool {
	group, ok := strategy.Definitions[g.GroupName]
	if !ok {
		return false
	}
	for _, name := range group {
		if card.Name == name {
			return true
		}
	}
	return false
}

type NameMatchCondition struct {
	Name string
}

func (n *NameMatchCondition) Matches(card *Card, _ *Strategy) bool {
	return card.Name == n.Name
}
*/

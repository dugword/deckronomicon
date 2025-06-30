package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/predicate"
	"deckronomicon/packages/game/mtg"
	"fmt"
	"strings"
)

func (p *StrategyParser) parsePredicate(data any) predicate.Predicate {
	switch node := data.(type) {
	case string:
		if strings.HasPrefix(node, "$") {
			return p.parseGroupPredicate(node)
		}
		return &predicate.Name{Name: node}
	case []any:
		var predicates []predicate.Predicate
		for _, item := range node {
			predicates = append(predicates, p.parsePredicate(item))
		}
		return &predicate.And{Predicates: predicates}
	case map[string]any:
		/* // TODO Think about this
		if len(node) != 1 {
			// TODO: add context path to p.errors.
			p.errors.Add(fmt.Errorf("expected exactly one key, got %d", len(node)))
			return nil
		}
		*/
		for key, val := range node {
			var predicates []predicate.Predicate
			switch key {
			// todo move to a separate function
			case "And", "All", "AllOf":
				items, ok := val.([]any)
				if !ok {
					p.errors.Add(fmt.Errorf("expected list for 'and', got %T", val))
					return nil
				}
				for _, item := range items {
					predicates = append(predicates, p.parsePredicate(item))
				}
				predicates = append(predicates, &predicate.And{Predicates: predicates})
				// todo move to a separate function
			case "Or", "Any", "AnyOf":
				items, ok := val.([]any)
				if !ok {
					p.errors.Add(fmt.Errorf("expected list for 'or', got %T", val))
					return nil
				}
				for _, item := range items {
					predicates = append(predicates, p.parsePredicate(item))
				}
				predicates = append(predicates, &predicate.Or{Predicates: predicates})
			case "Not":
				predicates = append(predicates, p.parseNotPredicate(val))
			case "CardType":
				predicates = append(predicates, p.parseCardTypePredicate(val))
			case "Subtype":
				predicates = append(predicates, p.parseSubtypePredicate(val))
			case "Supertype":
			case "Color":
			case "Mana":
			case "Power":
			case "Toughness":
			default:
				p.errors.Add(fmt.Errorf("unknown key: %s", key))
				return nil
			}
			if len(predicates) == 0 {
				p.errors.Add(fmt.Errorf("no valid predicates found for key '%s'", key))
				return nil
			}
			if len(predicates) == 1 {
				return predicates[0]
			}
			return &predicate.And{Predicates: predicates}
		}
	default:
		p.errors.Add(fmt.Errorf("unexpected type: %T", data))
		return nil
	}
	// TODO handle the switch statement to always return so we don't have the compiler hit this
	panic("unreachable")
}

func (p *StrategyParser) parseGroupPredicate(name string) predicate.Predicate {
	if !strings.HasPrefix(name, "$") {
		p.errors.Add(fmt.Errorf("group predicate must start with '$', got %s", name))
		return nil
	}
	groupName := strings.TrimPrefix(name, "$")
	if groupName == "" {
		p.errors.Add(fmt.Errorf("group predicate name cannot be empty"))
		return nil
	}
	group, ok := p.groups[groupName]
	if !ok {
		p.errors.Add(fmt.Errorf("group '%s' not found in definitions", groupName))
		return nil
	}
	var predicates []predicate.Predicate
	for _, item := range group {
		predicate := p.parsePredicate(item)
		if predicate != nil {
			predicates = append(predicates, predicate)
		}
	}
	if len(predicates) == 0 {
		p.errors.Add(fmt.Errorf("no valid predicate found in group '%s'", groupName))
		return nil
	}
	if len(predicates) == 1 {
		return predicates[0]
	}
	return &predicate.Or{Predicates: predicates}
}

func (p *StrategyParser) parseNotPredicate(value any) predicate.Predicate {
	return &predicate.Not{Predicate: p.parsePredicate(value)}
}

func (p *StrategyParser) parseCardTypePredicate(value any) predicate.Predicate {
	typeStr, ok := value.(string)
	if !ok {
		p.errors.Add(fmt.Errorf("expected string for 'CardType', got %T", value))
		return nil
	}
	cardType, ok := mtg.StringToCardType(typeStr)
	if !ok {
		p.errors.Add(fmt.Errorf("invalid card type: %s", typeStr))
		return nil
	}
	return &predicate.CardType{CardType: cardType}
}

func (p *StrategyParser) parseSubtypePredicate(value any) predicate.Predicate {
	typeStr, ok := value.(string)
	if !ok {
		p.errors.Add(fmt.Errorf("expected string for 'Subtype', got %T", value))
		return nil
	}
	subtype, ok := mtg.StringToSubtype(typeStr)
	if !ok {
		p.errors.Add(fmt.Errorf("invalid subtype: %s", typeStr))
		return nil
	}
	return &predicate.Subtype{Subtype: subtype}
}

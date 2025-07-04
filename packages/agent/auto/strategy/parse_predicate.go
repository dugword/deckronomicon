package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/predicate"
	"deckronomicon/packages/game/mtg"
	"errors"
	"fmt"
	"strings"
)

func (p *StrategyParser) parsePredicate(raw any) (predicate.Predicate, error) {
	switch node := raw.(type) {
	case string:
		if strings.HasPrefix(node, "$") {
			return p.parseGroupPredicate(node)
		}
		return &predicate.Name{Name: node}, nil
	case []any:
		var predicates []predicate.Predicate
		for _, item := range node {
			pred, err := p.parsePredicate(item)
			if err != nil {
				return nil, fmt.Errorf("error parsing predicate item %v: %w", item, err)
			}
			predicates = append(predicates, pred)
		}
		return &predicate.And{Predicates: predicates}, nil
	case map[string]any:
		var predicates []predicate.Predicate
		for key, val := range node {
			switch key {
			// todo move to a separate function
			case "And", "All", "AllOf":
				items, ok := val.([]any)
				if !ok {
					return nil, fmt.Errorf("expected list for 'And', got %T", val)
				}
				var children []predicate.Predicate
				for _, item := range items {
					pred, err := p.parsePredicate(item)
					if err != nil {
						return nil, fmt.Errorf("error parsing 'And' predicate item %v: %w", item, err)
					}
					children = append(children, pred)
				}
				predicates = append(predicates, &predicate.And{Predicates: children})
			case "Or", "Any", "AnyOf":
				items, ok := val.([]any)
				if !ok {
					return nil, fmt.Errorf("expected list for 'Or', got %T", val)
				}
				var children []predicate.Predicate
				for _, item := range items {
					pred, err := p.parsePredicate(item)
					if err != nil {
						return nil, fmt.Errorf("error parsing 'Or' predicate item %v: %w", item, err)
					}
					children = append(children, pred)
				}
				predicates = append(predicates, &predicate.Or{Predicates: children})
			case "Not":
				pred, err := p.parseNotPredicate(val)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Not' predicate: %w", err)
				}
				predicates = append(predicates, pred)
			case "CardType":
				pred, err := p.parseCardTypePredicate(val)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'CardType' predicate: %w", err)
				}
				predicates = append(predicates, pred)
			case "Subtype":
				pred, err := p.parseSubtypePredicate(val)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'Subtype' predicate: %w", err)
				}
				predicates = append(predicates, pred)
				/*
					case "Supertype":
						pred, err := p.parseSupertypePredicate(val)
						if err != nil {
							return nil, fmt.Errorf("error parsing 'Supertype' predicate: %w", err)
						}
						predicates = append(predicates, pred)
					case "Color":
						pred, err := p.parseColorPredicate(val)
						if err != nil {
							return nil, fmt.Errorf("error parsing 'Color' predicate: %w", err)
						}
						predicates = append(predicates, pred)
					case "Mana":
						pred, err := p.parseManaPredicate(val)
						if err != nil {
							return nil, fmt.Errorf("error parsing 'Mana' predicate: %w", err)
						}
						predicates = append(predicates, pred)
					case "Power":
						pred, err := p.parsePowerPredicate(val)
						if err != nil {
							return nil, fmt.Errorf("error parsing 'Power' predicate: %w", err)
						}
						predicates = append(predicates, pred)
					case "Toughness":
						pred, err := p.parseToughnessPredicate(val)
						if err != nil {
							return nil, fmt.Errorf("error parsing 'Toughness' predicate: %w", err)
						}
						predicates = append(predicates, pred)
				*/
			default:
				return nil, fmt.Errorf("unknown key %q in predicate", key)
			}
		}
		if len(predicates) == 0 {
			return nil, errors.New("no valid predicates found")
		}
		if len(predicates) == 1 {
			return predicates[0], nil
		}
		return &predicate.And{Predicates: predicates}, nil
	default:
		return nil, fmt.Errorf("expected string, list, or object for predicate, got %T", raw)
	}
}

func (p *StrategyParser) parseGroupPredicate(name string) (predicate.Predicate, error) {
	if !strings.HasPrefix(name, "$") {
		return nil, fmt.Errorf("group predicate must start with '$', got %s", name)
	}
	groupName := strings.TrimPrefix(name, "$")
	if groupName == "" {
		return nil, fmt.Errorf("group predicate name cannot be empty")
	}
	group, ok := p.groups[groupName]
	if !ok {
		return nil, fmt.Errorf("group '%s' not found in definitions", groupName)
	}
	var predicates []predicate.Predicate
	for _, item := range group.Members {
		pred, err := p.parsePredicate(item)
		if err != nil {
			return nil, fmt.Errorf("error parsing predicate in group '%s': %w", groupName, err)
		}
		if pred != nil {
			predicates = append(predicates, pred)
		}
	}
	if len(predicates) == 0 {
		return nil, fmt.Errorf("no valid predicate found in group '%s'", groupName)
	}
	if len(predicates) == 1 {
		return predicates[0], nil
	}
	return &predicate.Or{Predicates: predicates}, nil
}

func (p *StrategyParser) parseNotPredicate(value any) (predicate.Predicate, error) {
	pred, err := p.parsePredicate(value)
	if err != nil {
		return nil, fmt.Errorf("error parsing 'Not' predicate: %w", err)
	}
	return &predicate.Not{Predicate: pred}, nil
}

func (p *StrategyParser) parseCardTypePredicate(value any) (predicate.Predicate, error) {
	typeStr, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("expected string for 'CardType', got %T", value)
	}
	cardType, ok := mtg.StringToCardType(typeStr)
	if !ok {
		return nil, fmt.Errorf("invalid card type: %s", typeStr)
	}
	return &predicate.CardType{CardType: cardType}, nil
}

func (p *StrategyParser) parseSubtypePredicate(value any) (predicate.Predicate, error) {
	typeStr, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("expected string for 'Subtype', got %T", value)
	}
	subtype, ok := mtg.StringToSubtype(typeStr)
	if !ok {
		return nil, fmt.Errorf("invalid subtype: %s", typeStr)
	}
	return &predicate.Subtype{Subtype: subtype}, nil
}

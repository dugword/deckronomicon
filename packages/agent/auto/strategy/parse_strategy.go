package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/action"
	"deckronomicon/packages/agent/auto/strategy/evaluator"
	"deckronomicon/packages/agent/auto/strategy/predicate"
	"fmt"
	"path/filepath"
)

// region Strategy Parser Types

type StrategyParser struct {
	actions    map[string]action.ActionNode
	choices    map[string]*Choice
	conditions map[string]evaluator.Evaluator
	groups     map[string]*Group
	modes      map[string]*Mode
	rules      map[string]*Rule
	selectors  map[string]predicate.Selector
}

// StrategyParserData holds a map of file names to their contents
// for each type of fragment file used in the strategy parser.
type StrategyParserData struct {
	ActionFragmentFilesData    map[string][]byte
	ChoiceFragmentFilesData    map[string][]byte
	ConditionFragmentFilesData map[string][]byte
	GroupFragmentFilesData     map[string][]byte
	ModeFragmentFilesData      map[string][]byte
	RuleFragmentFilesData      map[string][]byte
	SelectorFragmentFilesData  map[string][]byte
}

type StrategyFile struct {
	Name        string `json:"Name,omitempty" yaml:"Name,omitempty"`
	Description string `json:"Description,omitempty" yaml:"Description,omitempty"`
	// TODO: Groups cannot contain references to other groups.
	Groups  []any            `json:"Groups,omitempty" yaml:"Groups,omitempty"`
	Modes   []any            `json:"Modes,omitempty" yaml:"Modes,omitempty"`
	Rules   map[string][]any `json:"Rules,omitempty" yaml:"Rules,omitempty"`
	Choices map[string][]any `json:"Choices,omitempty" yaml:"Choices,omitempty"`
}

// endregion

// region New Strategy Parser
func NewStrategyParser(data StrategyParserData) (*StrategyParser, error) {
	parser := StrategyParser{
		actions:    map[string]action.ActionNode{},
		choices:    map[string]*Choice{},
		conditions: map[string]evaluator.Evaluator{},
		groups:     map[string]*Group{},
		modes:      map[string]*Mode{},
		rules:      map[string]*Rule{},
		selectors:  map[string]predicate.Selector{},
	}
	// Groups must be parsed first as they may be referenced in other fragments.
	groupFragments, err := parser.ParseGroupFragmentFiles(data.GroupFragmentFilesData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse group fragments: %w", err)
	}
	parser.groups = groupFragments
	// Actions, Conditions, and Selector fragments must be parsed before Modes, Rules and Choice fragments.
	// This is because they may be referenced in the latter.
	actions, err := parser.ParseActionFragmentFiles(data.ActionFragmentFilesData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse action fragments: %w", err)
	}
	parser.actions = actions
	conditions, err := parser.ParseConditionFragmentFiles(data.ConditionFragmentFilesData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse condition fragments: %w", err)
	}
	parser.conditions = conditions
	selectors, err := parser.ParseSelectorFragmentFiles(data.SelectorFragmentFilesData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse selector fragments: %w", err)
	}
	parser.selectors = selectors

	// Choices, Mode, and Rule fragments are parsed last as they may reference the previously parsed fragments.
	choices, err := parser.ParseChoiceFragmentFiles(data.ChoiceFragmentFilesData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse choice fragments: %w", err)
	}
	parser.choices = choices
	modes, err := parser.ParseModeFragmentFiles(data.ModeFragmentFilesData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mode fragments: %w", err)
	}
	parser.modes = modes
	rules, err := parser.ParseRuleFragmentFiles(data.RuleFragmentFilesData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rule fragments: %w", err)
	}
	parser.rules = rules
	return &parser, nil
}

// endregion

// region Parser

func (p *StrategyParser) Parse(data []byte, filePath string) (*Strategy, error) {
	var strategyFile StrategyFile
	if err := unmarshalByExt(filepath.Ext(filePath), data, &strategyFile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal strategy file: %w", err)
	}
	// Groups must be parsed first as they may be referenced in other fragments.
	for _, groupFragment := range strategyFile.Groups {
		group, err := p.ParseGroupNode(groupFragment)
		if err != nil {
			return nil, fmt.Errorf("failed to parse group fragment: %w", err)
		}
		if _, ok := p.groups[group.Name]; ok {
			return nil, fmt.Errorf("group %q already exists in strategy file %s", group.Name, filePath)
		}
		p.groups[group.Name] = group
	}
	// region Parse Strategy Nodes
	strategy := Strategy{
		Name:        strategyFile.Name,
		Description: strategyFile.Description,
	}
	for _, rawMode := range strategyFile.Modes {
		mode, err := p.ParseModeNode(rawMode)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mode node: %w", err)
		}
		strategy.Modes = append(strategy.Modes, mode)
	}
	strategy.Rules = map[string][]*Rule{}
	for mode, rawRules := range strategyFile.Rules {
		for _, rawRule := range rawRules {
			rule, err := p.ParseRuleNode(rawRule)
			if err != nil {
				return nil, fmt.Errorf("failed to parse rule node: %w", err)
			}
			strategy.Rules[mode] = append(strategy.Rules[mode], rule)
		}
	}
	strategy.Choices = map[string][]*Choice{}
	for mode, rawChoices := range strategyFile.Choices {
		for _, rawChoice := range rawChoices {
			choice, err := p.ParseChoiceNode(rawChoice)
			if err != nil {
				return nil, fmt.Errorf("failed to parse choice node: %w", err)
			}
			strategy.Choices[mode] = append(strategy.Choices[mode], choice)
		}
	}
	// endregion

	return &strategy, nil
}

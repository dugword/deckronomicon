package strategy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type StrategyFile struct {
	Name        string           `json:"Name,omitempty" yaml:"Name,omitempty"`
	Description string           `json:"Description,omitempty" yaml:"Description,omitempty"`
	Groups      map[string][]any `json:"Groups,omitempty" yaml:"Groups,omitempty"`
	Modes       []map[string]any `json:"Modes,omitempty" yaml:"Modes,omitempty"`
	Rules       map[string]any   `json:"Rules,omitempty" yaml:"Rules,omitempty"`
}

type StrategyParser struct {
	errors         *ParseErrors
	sourceFile     string
	sourceFileType string // ".json", ".yaml"
	groups         map[string][]any
	rules          map[string]Rule
}

type ParseErrors struct {
	errors []error
}

func (e *ParseErrors) Add(err error) {
	if err != nil {
		e.errors = append(e.errors, err)
	}
}

func (e *ParseErrors) HasErrors() bool {
	return len(e.errors) > 0
}

func (e *ParseErrors) Error() string {
	if len(e.errors) == 0 {
		return ""
	}
	var errStrings []string
	for _, err := range e.errors {
		errStrings = append(errStrings, err.Error())
	}
	return "Parse errors: " + strings.Join(errStrings, ", ")
}

func LoadStrategy(path string) (*Strategy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	// New Parser - TODO make this a construtor
	parser := &StrategyParser{
		errors:         &ParseErrors{},
		sourceFile:     path,
		sourceFileType: filepath.Ext(path),
	}
	// Maybe just pass in path?
	strategy, err := parser.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("errors encountered while parsing strategy: %s, %w", parser.errors.Error(), err)
	}
	return strategy, nil
}

func (p *StrategyParser) Parse(data []byte) (*Strategy, error) {
	var strategyFile StrategyFile
	switch p.sourceFileType {
	case ".json":
		if err := json.Unmarshal(data, &strategyFile); err != nil {
			p.errors.Add(fmt.Errorf("failed to unmarshal strategy JSON: %w", err))
			return nil, p.errors
		}
	case ".yaml":
		if err := yaml.Unmarshal(data, &strategyFile); err != nil {
			p.errors.Add(fmt.Errorf("failed to unmarshal strategy YAML: %w", err))
			return nil, p.errors
		}
	default:
		p.errors.Add(fmt.Errorf("unsupported file type: %s", p.sourceFileType))
		return nil, p.errors
	}
	p.groups = strategyFile.Groups
	if err := p.ParseRuleFiles(); err != nil {
		p.errors.Add(fmt.Errorf("failed to parse rule files: %w", err))
		return nil, p.errors
	}
	var strategy Strategy
	var modes []*Rule
	for _, rawMode := range strategyFile.Modes {
		mode, err := p.ParseModeNode(rawMode)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mode: %w", err)
		}
		modes = append(modes, mode)
	}
	strategy.Modes = modes
	strategy.Rules = map[string][]*Rule{}
	for name, rawRules := range strategyFile.Rules {
		var rules []*Rule
		for _, rawRule := range rawRules.([]any) {
			rule, err := p.ParseRuleNode(rawRule)
			if err != nil {
				return nil, fmt.Errorf("failed to parse rule in %s: %w", name, err)
			}
			rules = append(rules, rule)
		}
		strategy.Rules[name] = rules
	}
	return &strategy, nil
}

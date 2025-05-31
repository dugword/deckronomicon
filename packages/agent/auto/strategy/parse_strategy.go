package strategy

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type StrategyParser struct {
	errors      *ParseErrors
	sourceFile  string
	definitions map[string][]string
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
		errors:     &ParseErrors{},
		sourceFile: path,
	}

	// Maybe just pass in path?
	strategy, err := parser.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("errors encountered while parsing strategy: %s", parser.errors.Error())
	}

	return strategy, nil
}

func (p *StrategyParser) Parse(data []byte) (*Strategy, error) {
	var strategy Strategy
	if err := json.Unmarshal(data, &strategy); err != nil {
		p.errors.Add(fmt.Errorf("failed to unmarshal strategy JSON: %w", err))
		return nil, p.errors
	}
	p.definitions = strategy.Definitions
	var outRules []Rule
	for _, mode := range strategy.Modes {
		outRules = append(outRules, p.parseRule(mode))
	}
	strategy.Modes = outRules
	for name, rules := range strategy.Rules {
		var outRules []Rule
		for _, rule := range rules {
			outRules = append(outRules, p.parseRule(rule))
		}
		strategy.Rules[name] = outRules
	}
	if p.errors.HasErrors() {
		return nil, fmt.Errorf("errors encountered while parsing strategy: %w", p.errors)
	}
	return &strategy, nil
}

func (p *StrategyParser) parseRule(rule Rule) Rule {
	var r Rule
	r.Name = rule.Name
	r.Description = rule.Description
	r.Then = p.parseActionNode(rule.RawThen)
	fmt.Println("Parsed action:", r.Then)
	r.When = p.parseConditionNode(rule.RawWhen)
	fmt.Println("Parsed condition:", r.When)
	return r
}

package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/action"
	"deckronomicon/packages/agent/auto/strategy/evaluator"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type RuleFile struct {
	Name        string              `json:"Name" yaml:"Name"`
	Description string              `json:"Description" yaml:"Description"`
	RawWhen     map[string]any      `json:"When" yaml:"When"`
	When        evaluator.Evaluator `json:"-" yaml:"-"`
	RawThen     map[string]any      `json:"Then" yaml:"Then"`
	Then        action.ActionNode   `json:"-" yaml:"-"`
}

func (p *StrategyParser) ParseRuleNode(ruleNode any) (*Rule, error) {
	var r Rule
	switch node := ruleNode.(type) {
	case string:
		rule, ok := p.rules[node]
		if !ok {
			return nil, fmt.Errorf("rule %s not found in rules map", node)
		}
		r = rule
	case map[string]any:
		for key, value := range node {
			switch key {
			case "Name":
				if name, ok := value.(string); ok {
					r.Name = name
				} else {
					return nil, fmt.Errorf("invalid type for Name: expected string, got %T", value)
				}
			case "Description":
				if desc, ok := value.(string); ok {
					r.Description = desc
				}
			case "When":
				r.When = p.parseEvaluator(value)
			case "Then":
				r.Then = p.parseActionNode(value)
			default:
				return nil, fmt.Errorf("unknown key in rule node: %s", key)
			}
		}
	}
	if r.Then == nil {
		return nil, fmt.Errorf("rule %s has no action defined", r.Name)
	}
	if p.errors.HasErrors() {
		return nil, fmt.Errorf("errors encountered while parsing rule node: %s", p.errors.Error())
	}
	return &r, nil
}

func (p *StrategyParser) ParseModeNode(ruleNode map[string]any) (*Rule, error) {
	var r Rule
	for key, value := range ruleNode {
		switch key {
		case "Name":
			if name, ok := value.(string); ok {
				r.Name = name
			} else {
				return nil, fmt.Errorf("invalid type for Name: expected string, got %T", value)
			}
		case "Description":
			if desc, ok := value.(string); ok {
				r.Description = desc
			}
		case "When":
			r.When = p.parseEvaluator(value)
		default:
			return nil, fmt.Errorf("unknown key in rule node: %s", key)
		}
	}
	if p.errors.HasErrors() {
		return nil, fmt.Errorf("errors encountered while parsing rule node: %s", p.errors.Error())
	}
	return &r, nil
}

func (p *StrategyParser) ParseRuleFile(ruleFile RuleFile) Rule {
	var r Rule
	r.Name = ruleFile.Name
	r.Description = ruleFile.Description
	r.Then = p.parseActionNode(ruleFile.RawThen)
	r.When = p.parseEvaluator(ruleFile.RawWhen)
	return r
}

/*
func (p *StrategyParser) parseRule(rule Rule) Rule {

}
*/

// TODO: Think about how to manage directories and files for rule fragments,
// I would like it to be less magical

func (p *StrategyParser) ParseRuleFiles() error {
	ruleDiretory := filepath.Join(filepath.Dir(p.sourceFile), "rules")
	files, err := os.ReadDir(ruleDiretory)
	if err != nil {
		return fmt.Errorf("failed to read rules directory: %w", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}
		filePath := filepath.Join(ruleDiretory, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			p.errors.Add(fmt.Errorf("failed to read rule file %s: %w", file.Name(), err))
			continue
		}
		var ruleFile RuleFile
		switch filepath.Ext(file.Name()) {
		case ".json":
			if err := json.Unmarshal(data, &ruleFile); err != nil {
				p.errors.Add(fmt.Errorf("failed to unmarshal rule JSON from %s: %w", file.Name(), err))
				continue
			}
		case ".yaml":
			if err := yaml.Unmarshal(data, &ruleFile); err != nil {
				p.errors.Add(fmt.Errorf("failed to unmarshal rule YAML from %s: %w", file.Name(), err))
				continue
			}
		default:
			p.errors.Add(fmt.Errorf("unsupported rule file type: %s", file.Name()))
			continue
		}
		if ruleFile.Name == "" {
			p.errors.Add(fmt.Errorf("rule in file %s has no name", file.Name()))
			continue
		}
		rule := p.ParseRuleFile(ruleFile)
		if p.errors.HasErrors() {
			return fmt.Errorf("errors encountered while parsing rule from %s: %s", file.Name(), p.errors.Error())
		}
		if p.rules == nil {
			p.rules = map[string]Rule{}
		}
		if _, ok := p.rules[ruleFile.Name]; ok {
			p.errors.Add(fmt.Errorf("duplicate rule name found: %s in file %s", ruleFile.Name, file.Name()))
		}
		p.rules[rule.Name] = rule
	}
	return nil
}

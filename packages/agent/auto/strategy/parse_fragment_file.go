package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/action"
	"deckronomicon/packages/agent/auto/strategy/evaluator"
	"deckronomicon/packages/agent/auto/strategy/predicate"
	"encoding/json"
	"fmt"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// region Fragment File Parsers

func (p *StrategyParser) ParseActionFragmentFiles(fragmentFilesData map[string][]byte) (map[string]action.ActionNode, error) {
	return parseNamedFragments(
		fragmentFilesData,
		"Actions",
		p.ParseNamedActionNode,
	)
}

func (p *StrategyParser) ParseChoiceFragmentFiles(fragmentFilesData map[string][]byte) (map[string]*Choice, error) {
	return parseNamedFragments(
		fragmentFilesData,
		"Choices",
		p.ParseChoiceNode,
	)
}

func (p *StrategyParser) ParseConditionFragmentFiles(fragmentFilesData map[string][]byte) (map[string]evaluator.Evaluator, error) {
	return parseNamedFragments(
		fragmentFilesData,
		"Conditions",
		p.ParseNamedConditionNode,
	)
}

func (p *StrategyParser) ParseGroupFragmentFiles(fragmentFilesData map[string][]byte) (map[string]*Group, error) {
	return parseNamedFragments(
		fragmentFilesData,
		"Groups",
		p.ParseGroupNode,
	)
}

func (p *StrategyParser) ParseModeFragmentFiles(fragmentFilesData map[string][]byte) (map[string]*Mode, error) {
	return parseNamedFragments(
		fragmentFilesData,
		"Modes",
		p.ParseModeNode,
	)
}

func (p *StrategyParser) ParseRuleFragmentFiles(fragmentFilesData map[string][]byte) (map[string]*Rule, error) {
	return parseNamedFragments(
		fragmentFilesData,
		"Rules",
		p.ParseRuleNode,
	)
}

func (p *StrategyParser) ParseSelectorFragmentFiles(fragmentFilesData map[string][]byte) (map[string]predicate.Selector, error) {
	return parseNamedFragments(
		fragmentFilesData,
		"Selectors",
		p.ParseNamedSelectorNode,
	)
}

func parseNamedFragments[ResultT any](
	fragmentFilesData map[string][]byte,
	kind string,
	parse func(any) (ResultT, error),
) (map[string]ResultT, error) {
	result := map[string]ResultT{}
	fragments, err := unmarshalFragmentFileData(fragmentFilesData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal fragment files: %w", err)
	}
	for fileName, fragment := range fragments {
		for _, itemRaw := range fragment[kind].([]any) {
			item, ok := itemRaw.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("invalid fragment format in %s expected map[string]any for %q, got %T", fileName, kind, itemRaw)
			}
			name := item["Name"].(string)
			if name == "" {
				return nil, fmt.Errorf("item in %q has no name in %q", kind, fileName)
			}
			result[name], err = parse(itemRaw)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %q fragment %s in %q: %w", kind, name, fileName, err)
			}
		}
	}
	return result, nil
}

// endregion

// region Fragment File Data Unmarshallers

func unmarshalFragmentFileData(fragmentFilesData map[string][]byte) (map[string]map[string]any, error) {
	fragments := map[string]map[string]any{}
	for fileName, data := range fragmentFilesData {
		var fragment map[string]any
		if err := unmarshalByExt(filepath.Ext(fileName), data, &fragment); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file %q: %w", fileName, err)
		}
		fragments[fileName] = fragment
	}
	return fragments, nil
}

func unmarshalByExt(ext string, data []byte, v any) error {
	switch ext {
	case ".json":
		return json.Unmarshal(data, v)
	case ".yaml":
		return yaml.Unmarshal(data, v)
	default:
		return fmt.Errorf("unsupported file extension %q", ext)
	}
}

// endregion

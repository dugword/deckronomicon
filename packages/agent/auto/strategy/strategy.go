package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/action"
	"deckronomicon/packages/agent/auto/strategy/evaluator"
	"deckronomicon/packages/agent/auto/strategy/predicate"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type EvaluatorContext struct {
	Game     *state.Game
	PlayerID string
	Strategy *Strategy
}

type Choice struct {
	Name        string
	Description string
	Source      string
	When        evaluator.Evaluator
	Choose      predicate.Selector
}

type Group struct {
	Name        string
	Description string
	Members     []string
}

type Mode struct {
	Name        string
	Description string
	When        evaluator.Evaluator
}

type Rule struct {
	Name        string
	Description string
	When        evaluator.Evaluator
	Then        action.ActionNode
}

type Strategy struct {
	Name        string
	Description string
	Modes       []*Mode
	Rules       map[string][]*Rule
	Choices     map[string][]*Choice
}

type StrategyDirectories struct {
	Actions    string
	Choices    string
	Conditions string
	Groups     string
	Modes      string
	Rules      string
	Selectors  string
}

func LoadStrategy(
	scenarioPath string,
	strategyFileName string,
	strategyDirectories StrategyDirectories,
) (*Strategy, error) {
	if !strings.HasSuffix(strategyFileName, ".yaml") {
		strategyFileName += ".yaml"
	}
	strategyFilePath := filepath.Join(scenarioPath, strategyFileName)
	strategyFileData, err := os.ReadFile(strategyFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", strategyFilePath, err)
	}
	actionsDirPath := filepath.Join(scenarioPath, strategyDirectories.Actions)
	actionFragmentFilesData, err := loadFragmentFiles(actionsDirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load action fragment files in %s: %w", actionsDirPath, err)
	}
	choicesDirPath := filepath.Join(scenarioPath, strategyDirectories.Choices)
	choiceFragmentFilesData, err := loadFragmentFiles(choicesDirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load choice fragment files in %s: %w", choicesDirPath, err)
	}
	conditionsDirPath := filepath.Join(scenarioPath, strategyDirectories.Conditions)
	conditionFragmentFilesData, err := loadFragmentFiles(conditionsDirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load condition fragment files in %s: %w", conditionsDirPath, err)
	}
	groupsDirPath := filepath.Join(scenarioPath, strategyDirectories.Groups)
	groupFragmentFilesData, err := loadFragmentFiles(groupsDirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load group fragment files in %s: %w", groupsDirPath, err)
	}
	modesDirPath := filepath.Join(scenarioPath, strategyDirectories.Modes)
	modeFragmentFilesData, err := loadFragmentFiles(modesDirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load mode fragment files in %s: %w", modesDirPath, err)
	}
	rulesDirPath := filepath.Join(scenarioPath, strategyDirectories.Rules)
	ruleFragmentFilesData, err := loadFragmentFiles(rulesDirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load rule fragment files in %s: %w", rulesDirPath, err)
	}
	selectorDirPath := filepath.Join(scenarioPath, "selectors")
	selectorFragmentFilesData, err := loadFragmentFiles(selectorDirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load selector fragment files in %s: %w", selectorDirPath, err)
	}
	parser, err := NewStrategyParser(StrategyParserData{
		ActionFragmentFilesData:    actionFragmentFilesData,
		ChoiceFragmentFilesData:    choiceFragmentFilesData,
		ConditionFragmentFilesData: conditionFragmentFilesData,
		GroupFragmentFilesData:     groupFragmentFilesData,
		ModeFragmentFilesData:      modeFragmentFilesData,
		RuleFragmentFilesData:      ruleFragmentFilesData,
		SelectorFragmentFilesData:  selectorFragmentFilesData,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create strategy parser: %w", err)
	}
	strategy, err := parser.Parse(
		strategyFileData,
		strategyFilePath,
	)
	if err != nil {
		return nil, fmt.Errorf("errors encountered while parsing strategy %s: %w", strategyFileName, err)
	}
	return strategy, nil
}

func loadFragmentFiles(dirName string) (map[string][]byte, error) {
	filesData := map[string][]byte{}
	if err := filepath.WalkDir(dirName, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk directory %s: %w", dirName, err)
		}
		if d.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}
		filesData[path] = data
		return nil
	}); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to walk directory %s: %w", dirName, err)
	}
	return filesData, nil
}

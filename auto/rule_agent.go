package auto

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	game "deckronomicon/game"
)

type RuleBasedAgent struct {
	Rules       []Rule
	ChoiceRules []ChoiceRule
}

func NewRuleBasedAgent(ruleFile string) *RuleBasedAgent {
	agent := &RuleBasedAgent{}
	agent.LoadRules(ruleFile)
	return agent
}

func (a *RuleBasedAgent) LoadRules(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Failed to load rule file: %v", err))
	}
	type rawRules struct {
		Rules       []Rule       `json:"rules"`
		ChoiceRules []ChoiceRule `json:"choice_rules"`
	}
	var parsed rawRules
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse rule file: %v", err))
	}
	// Do I need this? Or Should the rules just be in a ordered array already?
	sort.SliceStable(parsed.Rules, func(i, j int) bool {
		return parsed.Rules[i].Priority > parsed.Rules[j].Priority
	})
	a.Rules = parsed.Rules
	a.ChoiceRules = parsed.ChoiceRules
}

func (a *RuleBasedAgent) ReportState(state *game.GameState) {}

func (a *RuleBasedAgent) GetNextAction(state *game.GameState) game.GameAction {
	// TODO: make this configurable in the rules
	if state.LastActionFailed {
		return game.GameAction{Type: game.ActionConcede}
	}
	for _, rule := range a.Rules {
		fmt.Println("Rule =>", rule.Name)
		// TODO: change to !disabled so we don't need it in every rule
		if rule.Enabled && MatchesConditionSet(state, rule.When) {
			gameAction := rule.ToGameAction()
			if gameAction.Type == game.ActionPlay {

			}
		}
	}
	return game.GameAction{Type: game.ActionPass}
}

// TODO: Something better
func (a *RuleBasedAgent) ChooseOne(prompt string, options []game.Choice) game.Choice {
	// always chose the first option for now
	return options[0]
}

// --- Rule definitions ---

type Rule struct {
	Name     string       `json:"name"`
	Phase    string       `json:"phase,omitempty"` // this isn't main combat, it's combo not combo, rename
	When     ConditionSet `json:"when"`
	Then     Action       `json:"then"`
	Priority int          `json:"priority,omitempty"`
	Enabled  bool         `json:"enabled"`
}

type ConditionSet struct {
	// Zone-Based Card Presence
	HandContains          []string   `json:"hand_contains,omitempty"`
	HandContainsAny       []string   `json:"hand_contains_any,omitempty"`
	HandContainsAllGroups [][]string `json:"hand_contains_all_groups,omitempty"`
	HandContainsAnyGroups [][]string `json:"hand_contains_any_groups,omitempty"`
	HandLacks             []string   `json:"hand_lacks,omitempty"`
	HandLacksAny          []string   `json:"hand_lacks_any,omitempty"`
	HandLacksAllGroups    [][]string `json:"hand_lacks_all_groups,omitempty"`
	HandLacksAnyGroups    [][]string `json:"hand_lacks_any_groups,omitempty"`

	BattlefieldContains          []string   `json:"battlefield_contains,omitempty"`
	BattlefieldContainsAny       []string   `json:"battlefield_contains_any,omitempty"`
	BattlefieldContainsAllGroups [][]string `json:"battlefield_contains_all_groups,omitempty"`
	BattlefieldContainsAnyGroups [][]string `json:"battlefield_contains_any_groups,omitempty"`
	BattlefieldLacks             []string   `json:"battlefield_lacks,omitempty"`
	BattlefieldLacksAny          []string   `json:"battlefield_lacks_any,omitempty"`
	BattlefieldLacksAllGroups    [][]string `json:"battlefield_lacks_all_groups,omitempty"`
	BattlefieldLacksAnyGroups    [][]string `json:"battlefield_lacks_any_groups,omitempty"`

	GraveyardContains          []string   `json:"graveyard_contains,omitempty"`
	GraveyardContainsAny       []string   `json:"graveyard_contains_any,omitempty"`
	GraveyardContainsAllGroups [][]string `json:"graveyard_contains_all_groups,omitempty"`
	GraveyardContainsAnyGroups [][]string `json:"graveyard_contains_any_groups,omitempty"`
	GraveyardLacks             []string   `json:"graveyard_lacks,omitempty"`
	GraveyardLacksAny          []string   `json:"graveyard_lacks_any,omitempty"`
	GraveyardLacksAllGroups    [][]string `json:"graveyard_lacks_all_groups,omitempty"`
	GraveyardLacksAnyGroups    [][]string `json:"graveyard_lacks_any_groups,omitempty"`

	// Tags for flexible card classification
	Tags map[string]string `json:"tags,omitempty" yaml:"tags,omitempty"`

	// Turn & Game State Conditions
	TurnRange          []int    `json:"turn_range,omitempty"`
	Storm              string   `json:"storm,omitempty"`
	ManaAvailable      string   `json:"mana_available,omitempty"`
	LibrarySize        string   `json:"library_size,omitempty"`
	CardsInHand        string   `json:"cards_in_hand,omitempty"`
	GraveyardSize      string   `json:"graveyard_size,omitempty"`
	HasCastThisTurn    []string `json:"has_cast_this_turn,omitempty"`
	SpellCountThisTurn string   `json:"spell_count_this_turn,omitempty"`
	HasPlayedLand      *bool    `json:"has_played_land,omitempty"`
}

type Action struct {
	ActionType string `json:"action" yaml:"action"`
	Target     string `json:"target" yaml:"target"`
}

func (r Rule) ToGameAction() game.GameAction {
	command, ok := game.Commands[r.Then.ActionType]
	if ok {
		command.Action.Target = r.Then.Target
		fmt.Printf("Action %s => target %s\n", command.Action.Type, command.Action.Target)
		return command.Action
	}
	fmt.Println("unknown command =>", r.Then.ActionType)
	panic("need to handle an error here")
}

// --- Choice Rule System ---

type ChoiceRule struct {
	When    ChoiceCondition `json:"when"`
	Then    ChoiceAction    `json:"then"`
	Enabled bool            `json:"enabled"`
}

type ChoiceCondition struct {
	PromptContains string   `json:"prompt_contains,omitempty"`
	CardNames      []string `json:"card_names,omitempty"`
	CardTags       []string `json:"card_tags,omitempty"`
}

type ChoiceAction struct {
	Choose string `json:"choose"` // "top" or "bottom" or name match
}

func (cr ChoiceRule) Applies(prompt game.OptionPrompt) bool {
	if cr.When.PromptContains != "" && !strings.Contains(prompt.Message, cr.When.PromptContains) {
		return false
	}
	return true // can extend later with card checks
}

func (cr ChoiceRule) Resolve(prompt game.OptionPrompt) game.Choice {
	for i, option := range prompt.Options {
		if strings.EqualFold(cr.Then.Choose, option) {
			return game.Choice{Name: option, Index: i}
		}
	}
	return game.Choice{Name: prompt.Options[0], Index: 0}
}

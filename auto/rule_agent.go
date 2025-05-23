// Package rulesagent provides a configurable, strategy-driven player agent for simulating
// Magic: The Gathering gameplay using deterministic decision-making based on rule definitions.
//
// The rules agent reads a user-supplied strategy configuration file (referred to as a RuleSet),
// which defines a prioritized list of strategic behaviors. During each decision point in a game,
// the agent evaluates these strategies against the current GameState, including factors such as:
//
//   - Cards in hand
//   - Battlefield state (creatures, permanents, lands, etc.)
//   - Graveyard contents
//   - Life totals
//   - Mana available or potential mana generation
//
// Once a strategy condition is matched, the associated action is executed — such as casting a
// spell, activating an ability, or targeting a permanent — using defined or derived targets.
//
// Unlike the auto agent, which plays randomly or heuristically, the rules agent simulates how a
// skilled player might execute an idealized game plan. This enables precise, data-driven evaluations
// of deck performance and strategy efficacy.
//
// Key Features:
//
//   - Supports custom strategy files (RuleSets) defining prioritized, conditional logic.
//   - Includes a rich set of conditionals (ConditionSet) to express complex game state checks.
//   - Deterministic execution of matched rules ensures reproducible simulations.
//   - Useful for analyzing win rates, fizzle rates, and sequencing decisions across deck variants.
//
// Example Use Cases:
//
//   - Compare two strategies for the same deck by changing rule priority or structure.
//   - Measure fizzle rate when the agent always follows an optimal strategy.
//   - Validate if a combo line is consistently reachable under specific game conditions.
//
// Use the rulesagent package when you want to:
//
//   - Model advanced player behavior through conditional logic
//   - Evaluate strategic differences between deck variants
//   - Simulate high-skill-level performance for testing or optimization

package auto

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"

	game "deckronomicon/game"
)

// RuleBasedAgent implements the game.Agent interface, providing a rule-based
// decision-making system for Magic: The Gathering gameplay. It evaluates
// game states against a set of user-defined rules and conditions to determine
// the best course of action. The agent can also handle choice prompts and
// make decisions based on the defined rules.
type RuleBasedAgent struct {
	Rules       []Rule
	ChoiceRules []ChoiceRule
}

// NewRuleBasedAgent creates a new RuleBasedAgent instance and loads the rules
// from the specified JSON file. The rules are parsed and stored in the agent
// for use during gameplay. The rules are stored by priority, with higher
// priority rules evaluated first.
func NewRuleBasedAgent(ruleFile string) *RuleBasedAgent {
	agent := &RuleBasedAgent{}
	agent.LoadRules(ruleFile)
	return agent
}

// LoadRules loads the rules from a JSON file and parses them into the
// RuleBasedAgent instance.
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

// ReportState is a placeholder method for reporting the game state. It is not
// used in the the rule agent but is required by the game.Agent interface.
// This method can be implemented to log or analyze the game state at any
// point in the game.
func (a *RuleBasedAgent) ReportState(state *game.GameState) {}

// Confirm is a placeholder method for confirming actions during gameplay.
// TODO: Implement this
func (a *RuleBasedAgent) Confirm(prompt string, source game.ChoiceSource) (bool, error) {
	return true, nil
}

// GetNextAction evaluates the current game state against the defined rules
// and returns the next action to be taken. If a rule's condition is met, the
// corresponding action is executed. If no rules match, a default action
// (ActionPass) is returned. The method also handles special cases, such as
// conceding if the last action failed.
func (a *RuleBasedAgent) GetNextAction(state *game.GameState) game.GameAction {
	// TODO: make this configurable in the rules
	if state.LastActionFailed {
		return game.GameAction{Type: game.ActionConcede}
	}
	for _, rule := range a.Rules {
		fmt.Println("Rule =>", rule.Name)
		// TODO: change to !disabled so we don't need it in every rule
		conditionSetMatched, err := MatchesConditionSet(state, rule.When)
		if err != nil {
			// TODO: handle this better
			panic("need to handle this as an error")
		}
		if rule.Enabled && conditionSetMatched {
			gameAction := rule.ToGameAction()
			// TODO This could be more elegant
			if gameAction.Type == game.ActionPlay {
				if rule.Then.Target != "" {
					object, err := state.Hand.FindByName(rule.Then.Target)
					if err != nil {
						fmt.Println("ERROR: could not find card in hand =>", rule.Then.Target)
						os.Exit(1)
						return game.GameAction{Type: ""}
					}
					card, ok := object.(*game.Card)
					if !ok {
						fmt.Println("ERROR: could not cast card =>", rule.Then.Target)
						os.Exit(1)
					}

					// TODO canCast should probably be an error?
					preactions, canCast := PlanManaActivation(
						state.Battlefield.Permanents(),
						card.ManaCost(),
					)
					if !canCast {
						fmt.Println("ERROR (HANDLE THIS BETTER): can not cast even if I tap all my lands")
						for _, p := range state.Battlefield.Permanents() {
							fmt.Println("Permanent =>", p.Name())
						}
						os.Exit(0)
						return game.GameAction{Type: ""}
					}
					gameAction.Preactions = preactions
					return gameAction
				}
			}
		}
	}
	return game.GameAction{Type: game.ActionPass}
}

// ChooseOne handles choice prompts during gameplay. It evaluates the prompt
// against the defined choice rules and selects the best option based on the
// rules. If no choice rules match, the first option is selected by default.
// This method can be extended to include more complex decision-making logic
// based on the game state and the available choices.
func (a *RuleBasedAgent) ChooseOne(prompt string, source game.ChoiceSource, choices []game.Choice) (game.Choice, error) {
	if len(choices) == 0 {
		return game.Choice{}, errors.New("no choices available")
	}
	// always chose the first option for now
	return choices[0], nil
}

// Rule represents a single rule in the rules engine. Each rule has a name,
// a phase (optional), a set of conditions to check against the game state,
// an action to perform if the conditions are met, a priority for rule
// evaluation, and an enabled flag to determine if the rule is active.
// The priority determines the order in which rules are evaluated, with
// higher priority rules being checked first. The phase indicates the
// specific phase of the game in which the rule should be applied, such as
// "main" or "combat". The conditions are defined in a ConditionSet,
// which includes various checks for card presence, game state, and other
// factors. The action specifies what to do if the conditions are met,
// such as casting a spell or activating an ability. The enabled flag
// indicates whether the rule is currently active and should be evaluated.
type Rule struct {
	Name     string       `json:"name"`
	Phase    string       `json:"phase,omitempty"` // this isn't main combat, it's combo not combo, rename
	When     ConditionSet `json:"when"`
	Then     Action       `json:"then"`
	Priority int          `json:"priority,omitempty"`
	Enabled  bool         `json:"enabled"`
}

// MatchesConditionSet checks if the current game state matches the conditions
// defined in the ConditionSet. It evaluates the game state against the
// conditions specified in the rule, such as card presence in hand,
// battlefield, or graveyard, and game state conditions like storm count,
// mana available, and other factors. If all conditions are met, it returns
// true, indicating that the rule should be applied. If any condition is not
// met, it returns false, indicating that the rule does not apply in the
// current game state.
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

// Action represents the action to be taken when a rule's conditions are met.
// It includes the action type (e.g., "cast", "activate", "target") and the
// target of the action.
type Action struct {
	ActionType string `json:"action" yaml:"action"`
	Target     string `json:"target" yaml:"target"`
}

// ToGameAction converts the Rule's action into a game.GameAction.
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

// ChoiceRule represents a choice rule in the rules engine. Each choice rule
// has a condition that specifies when the rule applies, an action that
// defines what to do when the condition is met, and an enabled flag to
// determine if the rule is active. The condition includes checks for
// specific prompts, card names, and card tags. The action specifies how to
// choose between options, such as selecting the top or bottom option or
// matching a specific name. The enabled flag indicates whether the rule is
// currently active and should be evaluated.
type ChoiceRule struct {
	When    ChoiceCondition `json:"when"`
	Then    ChoiceAction    `json:"then"`
	Enabled bool            `json:"enabled"`
}

// ChoiceCondition represents the conditions under which a choice rule
// applies. It includes checks for specific prompts, card names, and card
// tags. The prompt is checked against the provided string, and the card
// names and tags are checked against the cards in the game state. The
// conditions are used to determine if the choice rule should be applied
// during gameplay. The conditions are defined in a way that allows for
// flexible matching, including partial matches and exact matches.
type ChoiceCondition struct {
	PromptContains string   `json:"prompt_contains,omitempty"`
	CardNames      []string `json:"card_names,omitempty"`
	CardTags       []string `json:"card_tags,omitempty"`
}

// ChoiceAction represents the action to be taken when a choice rule's
// conditions are met.
type ChoiceAction struct {
	Choose string `json:"choose"` // "top" or "bottom" or name match
}

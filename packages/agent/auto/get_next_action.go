package auto

import (
	"deckronomicon/packages/agent/auto/strategy/evalstate"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
	"fmt"
	"slices"
)

func (a *RuleBasedAgent) GetNextAction(game *state.Game) (engine.Action, error) {
	if a.interactive {
		a.ReportState(game)
	}
	ctx := evalstate.EvalState{
		Game:     game,
		PlayerID: a.playerID,
		Mode:     a.mode,
	}
	for _, mode := range a.strategy.Modes {
		if mode.Name == a.mode {
			continue
		}
		if mode.When.Evaluate(&ctx) {
			a.mode = mode.Name
			break
		}
	}
	var act engine.Action = action.PassPriorityAction{}
	for _, rule := range a.strategy.Rules[a.mode] {
		if !rule.When.Evaluate(&ctx) {
			continue
		}
		if a.interactive {
			a.EnterToContinue(fmt.Sprintf("Matched rule: %s", rule.Name))
		}
		var err error
		act, err = rule.Then.Resolve(&ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve action for rule %s: %w", rule.Name, err)
		}
		break
	}
	if slices.Contains(a.stops, game.Step()) && a.interactive {
		a.EnterToContinue(fmt.Sprintf("Action to apply: %s", act.Name()))
	}
	return act, nil
}

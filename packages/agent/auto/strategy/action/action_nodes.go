package action

// TODO: Rename this package so it doesn't conflict with engine/action

import (
	"deckronomicon/packages/agent/auto/strategy/evalstate"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/state"
)

// TODO: Rename nodes to align with actions/commands/events

type ActionNode interface {
	Resolve(ctx *evalstate.EvalState) (engine.Action, error)
}

type PassPriorityActionNode struct {
	Player state.Player
}

func (n *PassPriorityActionNode) Resolve(ctx *evalstate.EvalState) (engine.Action, error) {
	return action.NewPassPriorityAction(n.Player), nil
}

type ConcedeActionNode struct {
	Player state.Player
}

func (n *ConcedeActionNode) Resolve(ctx *evalstate.EvalState) (engine.Action, error) {
	return action.NewConcedeAction(n.Player), nil
}

type ActivateActionNode struct {
	AbilityInZone gob.AbilityInZone
	Player        state.Player
}

func (n *ActivateActionNode) Resolve(ctx *evalstate.EvalState) (engine.Action, error) {
	// var available []GameObject
	// TODO: This should be a method of Player so I don't have to iterate through a map. Think through managing this,
	// maybe have player return choices?
	/*
		for _, objects := range ctx.Player.GetAvailableToActivate(ctx.State) {
			available = append(available, objects...)
		}
		var filters []engine.FilterFunc
		for _, name := range a.AbilityNames {
			filters = append(filters, engine.HasName(name))
		}
		object, err := engine.FindFirstBy(available, engine.Or(filters...))
		if err != nil {
			return nil, errors.New("abilites not found: " + strings.Join(a.AbilityNames, ", "))
		}
		return &engine.GameAction{
			Type:   engine.ActionActivate,
			Target: engine.ActionTarget{ID: object.ID()},
		}, nil
	*/
	request := action.ActivateAbilityRequest{}
	return action.NewActivateAbilityAction(n.Player.ID(), request), nil
}

type PlayCardActionNode struct {
	CardInZone gob.CardInZone
	Player     state.Player
}

func (n *PlayCardActionNode) Resolve(ctx *evalstate.EvalState) (engine.Action, error) {
	/*
		var available []GameObject
		for _, objects := range ctx.Player.GetAvailableToPlay(ctx.State) {
			available = append(available, objects...)
		}
		var filters []engine.FilterFunc
		for _, name := range p.CardNames {
			filters = append(filters, engine.HasName(name))
		}
		object, err := engine.FindFirstBy(available, engine.Or(filters...))
		if err != nil {
			return nil, errors.New("cards not found: " + strings.Join(p.CardNames, ", "))
		}
		return &engine.GameAction{
			Type:   engine.ActionPlay,
			Target: engine.ActionTarget{ID: object.ID()},
		}, nil
	*/
	request := action.PlayLandRequest{
		CardID: n.CardInZone.Card().ID(),
	}
	return action.NewPlayLandAction(n.Player, request), nil
}

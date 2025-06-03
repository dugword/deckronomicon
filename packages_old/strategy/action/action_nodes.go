package action

import (
	"deckronomicon/packages/game"
	"deckronomicon/packages/strategy/evalstate"
)

type ActionNode interface {
	Resolve(ctx *evalstate.EvalState) (game.Action, error)
}

type PassAction struct{}

func (p *PassAction) Resolve(ctx *evalstate.EvalState) (game.Action, error) {
	return game.Action{
		Type: game.ActionPass,
	}, nil
}

type ConcedeAction struct{}

func (c *ConcedeAction) Resolve(ctx *evalstate.EvalState) (game.Action, error) {
	return game.Action{
		Type: game.ActionConcede,
	}, nil
}

type ActivateAction struct {
	AbilityNames []string `json:"card_name"`
}

func (a *ActivateAction) Resolve(ctx *evalstate.EvalState) (game.Action, error) {
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
	return game.Action{}, nil
}

type PlayAction struct {
	// TODO: Should this be an object for more complex play actions?
	// e.g. play a CardType: "Creature"
	CardNames []string `json:"card_name"`
}

func (p *PlayAction) Resolve(ctx *evalstate.EvalState) (game.Action, error) {
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
	return game.Action{}, nil
}

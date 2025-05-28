package auto

import (
	"deckronomicon/game"
	"errors"
	"strings"
)

type ActionNode interface {
	Resolve(ctx *EvaluatorContext) (*game.GameAction, error)
}

type PassAction struct{}

func (p *PassAction) Resolve(ctx *EvaluatorContext) (*game.GameAction, error) {
	return &game.GameAction{
		Type: game.ActionPass,
	}, nil
}

type ConcedeAction struct{}

func (c *ConcedeAction) Resolve(ctx *EvaluatorContext) (*game.GameAction, error) {
	return &game.GameAction{
		Type: game.ActionConcede,
	}, nil
}

type ActivateAction struct {
	AbilityNames []string `json:"card_name"`
}

func (a *ActivateAction) Resolve(ctx *EvaluatorContext) (*game.GameAction, error) {
	var available []game.GameObject
	// TODO: This should be a method of Player so I don't have to iterate through a map. Think through managing this,
	// maybe have player return choices?
	for _, objects := range ctx.player.GetAvailableToActivate(ctx.state) {
		available = append(available, objects...)
	}

	expanded, err := expandDefinitions(a.AbilityNames, ctx.strategy.Definitions)
	if err != nil {
		return nil, err
	}
	var filters []game.FilterFunc
	for _, name := range expanded {
		filters = append(filters, game.HasName(name))
	}
	object, err := game.FindFirstBy(available, game.Or(filters...))
	if err != nil {
		return nil, errors.New("abilites not found: " + strings.Join(a.AbilityNames, ", "))
	}
	return &game.GameAction{
		Type:   game.ActionActivate,
		Target: game.ActionTarget{ID: object.ID()},
	}, nil
}

type PlayAction struct {
	// TODO: Should this be an object for more complex play actions?
	// e.g. play a CardType: "Creature"
	CardNames []string `json:"card_name"`
}

func (p *PlayAction) Resolve(ctx *EvaluatorContext) (*game.GameAction, error) {
	var available []game.GameObject
	for _, objects := range ctx.player.GetAvailableToPlay(ctx.state) {
		available = append(available, objects...)
	}
	expaned, err := expandDefinitions(p.CardNames, ctx.strategy.Definitions)
	if err != nil {
		return nil, err
	}
	var filters []game.FilterFunc
	for _, name := range expaned {
		filters = append(filters, game.HasName(name))
	}
	object, err := game.FindFirstBy(available, game.Or(filters...))
	if err != nil {
		return nil, errors.New("cards not found: " + strings.Join(expaned, ", "))
	}
	return &game.GameAction{
		Type:   game.ActionPlay,
		Target: game.ActionTarget{ID: object.ID()},
	}, nil
}

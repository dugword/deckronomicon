package auto

import (
	"deckronomicon/game"
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
	panic("fix me")
	/*
		available := ctx.player.GetAvailableToActivate(ctx.state)
		expanded, err := expandDefinitions(a.AbilityNames, ctx.strategy.Definitions)
		if err != nil {
			return nil, err
		}
		var filters []game.FilterFunc
		for _, name := range expanded {
			filters = append(filters, game.HasName(name))
		}
		object, err := game.FindFirstBy(objects, game.Or(filters...))
		if err != nil {
			return nil, errors.New("abilites not found: " + strings.Join(a.AbilityNames, ", "))
		}
		return &game.GameAction{
			Type:   game.ActionActivate,
			Target: object.Name(),
		}, nil
	*/
}

type PlayAction struct {
	// TODO: Should this be an object for more complex play actions?
	// e.g. play a CardType: "Creature"
	CardNames []string `json:"card_name"`
}

func (p *PlayAction) Resolve(ctx *EvaluatorContext) (*game.GameAction, error) {
	panic("fix me")
	/*
		objects := ctx.player.GetAvailableToPlay(ctx.state)
		expaned, err := expandDefinitions(p.CardNames, ctx.strategy.Definitions)
		if err != nil {
			return nil, err
		}
		var filters []game.FilterFunc
		for _, name := range expaned {
			filters = append(filters, game.HasName(name))
		}
		object, err := game.FindFirstBy(objects, game.Or(filters...))
		if err != nil {
			return nil, errors.New("cards not found: " + strings.Join(expaned, ", "))
		}
		return &game.GameAction{
			Type:   game.ActionPlay,
			Target: object.Name(),
		}, nil
	*/
}

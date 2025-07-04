package action

// TODO: Rename this package so it doesn't conflict with engine/action

import (
	"deckronomicon/packages/agent/auto/strategy/evalstate"
	"deckronomicon/packages/agent/auto/strategy/predicate"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/engine/judge"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"fmt"
)

// TODO: Rename nodes to align with actions/commands/events

type ActionNode interface {
	Resolve(ctx *evalstate.EvalState) (engine.Action, error)
}

type PassPriorityActionNode struct{}

func (n *PassPriorityActionNode) Resolve(ctx *evalstate.EvalState) (engine.Action, error) {
	return action.NewPassPriorityAction(), nil
}

type ConcedeActionNode struct{}

func (n *ConcedeActionNode) Resolve(ctx *evalstate.EvalState) (engine.Action, error) {
	return action.NewConcedeAction(), nil
}

type ActivateActionNode struct {
	AbilityInZone *gob.AbilityInZone
}

func (n *ActivateActionNode) Resolve(ctx *evalstate.EvalState) (engine.Action, error) {
	request := action.ActivateAbilityRequest{
		AbilityID: n.AbilityInZone.Ability().ID(),
		SourceID:  n.AbilityInZone.Source().ID(),
		Zone:      n.AbilityInZone.Zone(),
		// TODO: Need to handle targets for effects properly
		// TargetsForEffects: n.AbilityInZone.Ability().TargetsForEffects(),
	}
	return action.NewActivateAbilityAction(request), nil
}

type LogMessageActionNode struct {
	Message string
}

func (n *LogMessageActionNode) Resolve(ctx *evalstate.EvalState) (engine.Action, error) {
	return action.NewLogMessageAction(n.Message), nil
}

type EmitMetricActionNode struct {
	Name  string
	Value int
}

func (n *EmitMetricActionNode) Resolve(ctx *evalstate.EvalState) (engine.Action, error) {
	return action.NewEmitMetricAction(n.Name, n.Value), nil
}

type PlayLandCardActionNode struct {
	Cards predicate.Predicate
}

func (n *PlayLandCardActionNode) Resolve(ctx *evalstate.EvalState) (engine.Action, error) {
	player := ctx.Game.GetPlayer(ctx.PlayerID)
	found := n.Cards.Select(query.NewQueryObjects(player.Hand().GetAll()))
	if len(found) == 0 {
		return nil, fmt.Errorf("no land cards found in hand for player %s", player.ID())
	}
	var playable []*gob.Card
	ruling := judge.Ruling{Explain: true}
	for _, obj := range found {
		card, ok := obj.(*gob.Card)
		if !ok {
			return nil, fmt.Errorf("object %s is not a card", obj.ID())
		}
		if !judge.CanPlayLand(ctx.Game, player.ID(), mtg.ZoneHand, card, &ruling) {
			continue
		}
		playable = append(playable, card)
	}
	if len(playable) == 0 {
		return nil, fmt.Errorf("no playable land cards found in hand for player %s", player.ID())
	}
	request := action.PlayLandRequest{
		CardID: playable[0].ID(), // Assuming we play the first found land card
	}
	return action.NewPlayLandAction(request), nil
}

type CastSpellActionNode struct {
	Cards          predicate.Predicate
	AdditionalCost predicate.Selector
}

// TODO: Look this over, wrote it quickly late at night
func (n *CastSpellActionNode) Resolve(ctx *evalstate.EvalState) (engine.Action, error) {
	player := ctx.Game.GetPlayer(ctx.PlayerID)
	found := n.Cards.Select(query.NewQueryObjects(player.Hand().GetAll()))
	if len(found) == 0 {
		return nil, fmt.Errorf("no spell cards found in hand for player %s", player.ID())
	}
	var castable []*gob.Card
	ruling := judge.Ruling{Explain: true}
	var costTarget target.Target
	for _, obj := range found {
		card, ok := obj.(*gob.Card)
		if !ok {
			return nil, fmt.Errorf("object %s is not a card", obj.ID())
		}
		autoPay := true
		additionalCost := card.AdditionalCost()
		if additionalCost != nil {
			if _, ok := additionalCost.(cost.Composite); ok {
				panic("Composite costs are not supported as additional costs yet")
			}
			if costWithTarget, ok := additionalCost.(cost.CostWithTarget); ok {
				cards := n.AdditionalCost.Select(query.NewQueryObjects(player.Hand().GetAll()))
				if len(cards) == 0 {
					return nil, fmt.Errorf("no additional cost cards found in hand for card %s", card.ID())
				}
				costTarget = target.Target{
					ID: cards[0].ID(),
				}
				additionalCost = costWithTarget.WithTarget(costTarget)
			}
		}
		totalCost := cost.NewComposite(card.ManaCost(), additionalCost)
		if !judge.CanCastSpellFromHand(ctx.Game, player.ID(), card, totalCost, autoPay, mana.Colors(), &ruling) {
			continue
		}
		castable = append(castable, card)
	}
	if len(castable) == 0 {
		return nil, fmt.Errorf("no castable spell cards found in hand for player %s", player.ID())
	}
	request := action.CastSpellRequest{
		CardID:         castable[0].ID(), // Assuming we cast the first found spell card
		ReplicateCount: 0,                // TODO: Handle replicate count
		// TargetsForEffects: n.CardInZone.Card().TargetsForEffects(),
		SpliceCardIDs: nil,           // TODO: Handle splice cards
		Flashback:     false,         // TODO: Handle flashback
		AutoPayCost:   true,          // TODO: Handle auto pay cost
		AutoPayColors: mana.Colors(), // TODO: Handle auto pay colors
		CostTarget:    costTarget,
	}
	return action.NewCastSpellAction(request), nil
}

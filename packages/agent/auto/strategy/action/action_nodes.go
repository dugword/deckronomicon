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
	AbilityInZone gob.AbilityInZone
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

type PlayLandCardActionNode struct {
	Cards predicate.Predicate
}

func (n *PlayLandCardActionNode) Resolve(ctx *evalstate.EvalState) (engine.Action, error) {
	player := ctx.Game.GetPlayer(ctx.PlayerID)
	found := n.Cards.Select(query.NewQueryObjects(player.Hand().GetAll()))
	if len(found) == 0 {
		return nil, fmt.Errorf("no land cards found in hand for player %s", player.ID())
	}
	var playable []gob.Card
	ruling := judge.Ruling{Explain: true}
	for _, obj := range found {
		card, ok := obj.(gob.Card)
		if !ok {
			return nil, fmt.Errorf("object %s is not a card", obj.ID())
		}
		if !judge.CanPlayLand(ctx.Game, player, mtg.ZoneHand, card, &ruling) {
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
	Cards predicate.Predicate
}

// TODO: Look this over, wrote it quickly late at night
func (n *CastSpellActionNode) Resolve(ctx *evalstate.EvalState) (engine.Action, error) {
	player := ctx.Game.GetPlayer(ctx.PlayerID)
	found := n.Cards.Select(query.NewQueryObjects(player.Hand().GetAll()))
	if len(found) == 0 {
		return nil, fmt.Errorf("no spell cards found in hand for player %s", player.ID())
	}
	var castable []gob.Card
	ruling := judge.Ruling{Explain: true}
	for _, obj := range found {
		card, ok := obj.(gob.Card)
		if !ok {
			return nil, fmt.Errorf("object %s is not a card", obj.ID())
		}
		autoPay := true
		totalCost := cost.NewComposite(card.ManaCost(), card.AdditionalCost())
		if !judge.CanCastSpellFromHand(ctx.Game, player, card, totalCost, autoPay, mana.Colors(), &ruling) {
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
	}
	return action.NewCastSpellAction(request), nil
}

package action

// TODO: Rename this package so it doesn't conflict with engine/action

import (
	"deckronomicon/packages/agent/auto/strategy/evalstate"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/engine/judge"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query/has"
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
	CardNames []string
}

func (n *PlayLandCardActionNode) Resolve(ctx *evalstate.EvalState) (engine.Action, error) {
	player := ctx.Game.GetPlayer(ctx.PlayerID)
	card, ok := player.Hand().Find(has.AnyName(n.CardNames...))
	if !ok {
		return nil, fmt.Errorf("failed to find playable land card in hand: %v", n.CardNames)
	}
	ruling := judge.Ruling{Explain: true}
	if !judge.CanPlayLand(ctx.Game, player, mtg.ZoneHand, card, &ruling) {
		return nil, fmt.Errorf("cannot play land card %q <%s>: %s", card.Name(), card.ID(), ruling.Why())
	}
	request := action.PlayLandRequest{
		CardID: card.ID(),
	}
	return action.NewPlayLandAction(request), nil
}

type CastSpellActionNode struct {
	CardInZone gob.CardInZone
}

func (n *CastSpellActionNode) Resolve(ctx *evalstate.EvalState) (engine.Action, error) {
	request := action.CastSpellRequest{
		CardID:         n.CardInZone.Card().ID(),
		ReplicateCount: 0, // TODO: Handle replicate count
		// TargetsForEffects: n.CardInZone.Card().TargetsForEffects(),
		SpliceCardIDs: nil,   // TODO: Handle splice cards
		Flashback:     false, // TODO: Handle flashback
	}
	return action.NewCastSpellAction(request), nil
}

package node

import (
	"deckronomicon/packages/agent/auto/strategy/evalstate"
	"deckronomicon/packages/agent/auto/strategy/predicate"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
)

type Node interface {
	Evaluate(ctx *evalstate.EvalState) bool
	Collect(ctx *evalstate.EvalState) []query.Object
}

type InZone struct {
	Zone  mtg.Zone
	Cards predicate.Predicate
}

func (n *InZone) Evaluate(ctx *evalstate.EvalState) bool {
	player := ctx.Game.GetPlayer(ctx.PlayerID)
	if n.Cards == nil {
		return true
	}
	switch n.Zone {
	case mtg.ZoneHand, mtg.ZoneGraveyard, mtg.ZoneLibrary, mtg.ZoneExile:
		cards, ok := player.GetCardsInZone(n.Zone)
		if !ok {
			// TODO: Think through how to handle invalid cases
			// Maybe something like the error record in the strategy
			// parser, or the ruling passed into the judge methods.
			return false
		}
		return n.Cards.Matches(query.NewQueryObjects(cards))
	case mtg.ZoneBattlefield:
		permanents := ctx.Game.Battlefield().GetAll()
		return n.Cards.Matches(query.NewQueryObjects(permanents))
	default:
		return false
	}
}

func (n *InZone) Collect(ctx *evalstate.EvalState) []query.Object {
	player := ctx.Game.GetPlayer(ctx.PlayerID)
	if n.Cards == nil {
		return nil
	}
	switch n.Zone {
	case mtg.ZoneHand, mtg.ZoneGraveyard, mtg.ZoneLibrary, mtg.ZoneExile:
		cardsInZone, ok := player.GetCardsInZone(n.Zone)
		if !ok {
			return nil
		}
		found := n.Cards.Select(query.NewQueryObjects(cardsInZone))
		return found
	case mtg.ZoneBattlefield:
		permanents := ctx.Game.Battlefield().GetAll()
		found := n.Cards.Select(query.NewQueryObjects(permanents))
		return found
	default:
		return nil
	}
}

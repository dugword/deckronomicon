package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/game/ability/activated"
	"deckronomicon/packages/game/action"
	"deckronomicon/packages/game/card"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/permanent"
	"deckronomicon/packages/game/player"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"fmt"
)

func GetAvailableToActivate(state *GameState, p *player.Player) (map[string][]query.Object, error) {
	available := map[string][]query.Object{}
	for _, object := range state.Battlefield().GetAll() {
		const ZoneBattlefield = string(mtg.ZoneBattlefield)
		perm, ok := object.(*permanent.Permanent)
		if !ok {
			return nil, ErrObjectNotPermanent
		}
		for _, ability := range perm.ActivatedAbilities() {
			if !ability.CanActivate(state, p.ID()) {
				continue
			}
			if !ability.Cost.CanPay(state, p) {
				continue
			}
			available[ZoneBattlefield] = append(
				available[ZoneBattlefield],
				ability,
			)
		}
	}
	for _, object := range p.Hand().GetAll() {
		const ZoneHand = string(mtg.ZoneHand)
		card, ok := object.(*card.Card)
		if !ok {
			return nil, ErrObjectNotCard
		}
		for _, ability := range card.ActivatedAbilities() {
			if !ability.CanActivate(state, p.ID()) {
				continue
			}
			if !ability.Cost.CanPay(state, p) {
				continue
			}
			available[ZoneHand] = append(
				available[ZoneHand],
				ability,
			)
		}
	}
	return available, nil
	/*
		available := map[string][]query.Object{}
		for _, object := range p.Hand().GetAll() {
			const ZoneHand = string(mtg.ZoneHand)
			c, ok := object.(*card.Card)
			if !ok {
				return nil, ErrObjectNotCard
			}
			if state.CanCastSorcery(p.ID()) {
				if c.Match(is.Land()) {
					if p.LandDrop == false {
						available[ZoneHand] = append(
							available[ZoneHand],
							c,
						)
						continue
					}
				} else if c.ManaCost().CanPay(state, p) {
					available[ZoneHand] = append(available[ZoneHand], c)
					continue
				}
			}
			if c.ManaCost().CanPay(state, p) && c.Match(has.CardType(mtg.CardTypeInstant)) {
				available[ZoneHand] = append(available[ZoneHand], c)
				continue
			}
		}
		return available, nil
	*/
}

// ActionActivateFunc handles the activate action. This is performed by the
// player to activate an ability of a permanent on the battlefield, or a card
// in hand or in the graveyard. The target is the name of the permanent or
// card.
// TODO: Support more than one activated ability
// TODO: Support activated abilities in hand and graveyard
func ActionActivateFunc(state *GameState, player *player.Player, target action.ActionTarget) (ActionResult, error) {
	available, err := GetAvailableToActivate(state, player)
	if err != nil {
		return ActionResult{}, fmt.Errorf(
			"failed to get available activated abilities: %w",
			err,
		)
	}
	choices := choose.CreateGroupedChoices(available)
	if len(choices) == 0 {
		return ActionResult{}, choose.ErrNoChoices
	}
	choice, err := player.Agent.ChooseOne(
		"Which ability to activate",
		choose.NewChoiceSource(string(action.ActionActivate)),
		choose.AddOptionalChoice(choices),
	)
	if err != nil {
		return ActionResult{}, fmt.Errorf("failed to choose ability: %w", err)
	}
	if choice == choose.ChoiceNone {
		return ActionResult{Message: "No choice made"}, nil
	}
	choiceZone := choice.Source.Name()
	targetObj, ok := query.Find(available[choiceZone], has.ID(choice.ID))
	if !ok {
		return ActionResult{}, query.ErrNotFound
	}
	var ability *activated.Ability
	ability, ok = targetObj.(*activated.Ability)
	if !ok {
		return ActionResult{}, fmt.Errorf("object is not an activated ability: %w", err)
	}
	if err := ability.Cost.Pay(state, player); err != nil {
		return ActionResult{}, fmt.Errorf("cannot pay activated ability cost: %w", err)
	}
	// Mana abilities
	if ability.IsManaAbility() {
		if err := ability.Resolve(state, player); err != nil {
			return ActionResult{}, fmt.Errorf("failed to resolve mana ability: %w", err)
		}
		_, ok := ability.Cost.(*cost.TapCost)
		if ok {
			/*
				state.EmitEvent(Event{
					Type:   EventTapForMana,
					Source: ability.source,
				}, player)
			*/
		}
	} else {
		// TODO: Figure out the interface for this
		state.Stack.Add(ability)
	}
	return ActionResult{
		Message: fmt.Sprintf(
			"ability activated: %s (%s)",
			ability.Name(),
			ability.Description(),
		),
	}, nil
	return ActionResult{}, nil
}

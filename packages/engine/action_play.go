package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/game/action"
	"deckronomicon/packages/game/card"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/permanent"
	"deckronomicon/packages/game/player"
	"deckronomicon/packages/game/spell"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"fmt"
)

func GetAvailableToPlay(state *GameState, p *player.Player) (map[string][]query.Object, error) {
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
}

// TODO This whole look up system is a bit of a hack. Should extract it out to a
// function and share it with action activate.
// ActionPlayFunc handles the play action. This is performed by the player to
// play a card from their hand. The target is the name of the card to play.
func ActionPlayFunc(state *GameState, player *player.Player, target action.ActionTarget) (result ActionResult, err error) {
	available, err := GetAvailableToPlay(state, player)
	if err != nil {
		return ActionResult{}, fmt.Errorf(
			"failed to get available cards to play: %w",
			err,
		)
	}
	var choiceZone string
	var targetObj query.Object
	if target.ID != "" {
	outer:
		for zone, objects := range available {
			for _, obj := range objects {
				if obj.ID() != target.ID {
					continue
				}
				choiceZone = zone
				targetObj = obj
				break outer
			}
		}
	}
	if targetObj == nil {
		choices := choose.CreateGroupedChoices(available)
		if len(choices) == 0 {
			return ActionResult{}, choose.ErrNoChoices
		}
		// TODO: Rethink this
		if target.Name != "" {
			var filteredChoices []choose.Choice
			for _, choice := range choices {
				if choice.Name == target.Name {
					filteredChoices = append(filteredChoices, choice)
				}
			}
			choices = filteredChoices
		}

		var choice choose.Choice
		if target.Name != "" && len(choices) == 1 {
			choice = choices[0]
		} else {
			choice, err = player.Agent.ChooseOne(
				"Which card to play",
				choose.NewChoiceSource(string(action.ActionPlay)),
				choose.AddOptionalChoice(choices),
			)
			if err != nil {
				return ActionResult{}, fmt.Errorf("failed to choose card: %w", err)
			}
			if choice == choose.ChoiceNone {
				return ActionResult{Message: "No choice made"}, nil
			}
		}
		choiceZone = choice.Source.Name()
		var ok bool
		targetObj, ok = query.Find(available[choiceZone], has.ID(choice.ID))
		if !ok {
			return ActionResult{}, fmt.Errorf("failed to find card to play")
		}
	}
	c, ok := targetObj.(*card.Card)
	if !ok {
		return ActionResult{}, ErrObjectNotCard
	}
	if c.Match(has.CardType(mtg.CardTypeLand)) {
		if err := PlayLand(state, player, c); err != nil {
			return ActionResult{}, fmt.Errorf("failed to play land: %w", err)
		}
		return ActionResult{
			Message: fmt.Sprintf("Played land %s", c.Name()),
		}, nil
	}
	if err := CastSpell(state, player, c, choiceZone); err != nil {
		return ActionResult{}, fmt.Errorf("failed to cast spell: %w", err)
	}
	return ActionResult{
		Message: fmt.Sprintf("Played spell %s", c.Name()),
	}, nil
}

// TODO: Should this be a method? Standardize on Logging
func PlayLand(state *GameState, player *player.Player, card *card.Card) error {
	if card.Match(has.CardType(mtg.CardTypeLand)) {
		if player.LandDrop {
			return mtg.ErrLandAlreadyPlayed
		}
		player.LandDrop = true
	}
	if _, err := player.TakeCard(card.ID(), mtg.ZoneHand); err != nil {
		return fmt.Errorf(
			"failed to remove card %s from player %s: %w",
			card.Name(),
			player.ID(),
			err,
		)
	}
	perm, err := permanent.NewPermanent(card, state, player)
	if err != nil {
		return fmt.Errorf(
			"failed to create permanent from %s: %w",
			card.Name(),
			err,
		)
	}
	state.battlefield.Add(perm)
	return nil
}

// TODO: Maybe this should be a method off of GameState
// or maybe a method off of Card, e.g. card.Cast() like Ability.Resolve()
func CastSpell(state *GameState, player *player.Player, c *card.Card, fromZone string) error {
	spell, err := spell.New(state, c)
	if err != nil {
		return fmt.Errorf("failed to create spell from %s: %w", c.Name(), err)
	}
	var spellCost cost.Cost
	var isFlashback bool
	// var isSplice bool
	// var spliceCost Cost
	if fromZone == string(mtg.ZoneGraveyard) {
		for _, ability := range spell.StaticAbilities() {
			if ability.ID() != string(mtg.StaticKeywordFlashback) {
				continue
			}
			isFlashback = true
			for _, modifier := range ability.Modifiers {
				if modifier.Key != "Cost" {
					continue
				}
				spellCost, err = cost.New(modifier.Value, spell)
				if err != nil {
					return fmt.Errorf("failed to create cost: %w", err)
				}
			}
		}
	} else {
		spellCost = spell.ManaCost()
	}
	var toSplice []query.Object
	if spell.Match(has.Subtype(mtg.SubtypeArcane)) {
		spliceCards := player.Hand().FindAll(
			query.And(
				has.StaticKeyword(mtg.StaticKeywordSplice),
				/*
					has.StaticAbilityModifier(
						game.AbilityKeywordSplice,
						EffectTag{Key: "Onto", Value: "Arcane"},
					),
				*/
			),
		)
		// var spliceCosts []Cost
		choices := choose.CreateChoices(spliceCards, mtg.ZoneHand)

		chosen, err := player.Agent.ChooseMany(
			"Choose cards to splice onto the spell",
			mtg.StaticKeywordSplice,
			choose.AddOptionalChoice(choices),
		)
		if err != nil {
			return fmt.Errorf("failed to choose card: %w", err)
		}
		for _, c := range chosen {
			if c == choose.ChoiceNone {
				break
			}
			obj, ok := player.Hand().Get(c.ID)
			if !ok {
				return fmt.Errorf("failed to get splice card from hand: %w", err)
			}
			spliceCard, ok := obj.(*card.Card)
			if !ok {
				return ErrObjectNotCard
			}
			for _, ability := range spliceCard.StaticAbilities() {
				if ability.ID() != mtg.StaticKeywordSplice.Name() {
					continue
				}
				for _, modifier := range ability.Modifiers {
					if modifier.Key != "Cost" {
						continue
					}
					spliceCost, err := cost.New(modifier.Value, spell)
					if err != nil {
						return fmt.Errorf("failed to create splice cost: %w", err)
					}
					accept, err := player.Agent.Confirm(
						fmt.Sprintf("Splice %s onto %s for %s?",
							spliceCard.Name(),
							spell.Name(),
							spliceCost.Description(),
						),
						mtg.StaticKeywordSplice,
					)
					if err != nil {
						return fmt.Errorf("failed to confirm splice cost: %w", err)
					}
					if !accept {
						continue
					}
					spellCost = cost.AddCosts(spellCost, spliceCost)
					toSplice = append(toSplice, spliceCard)
				}
			}
		}
	}

	var isReplicate bool
	var replicateCost cost.Cost
	var replicateCount int
	for _, ability := range spell.StaticAbilities() {
		if ability.ID() != mtg.StaticKeywordReplicate.Name() {
			continue
		}
		for _, modifier := range ability.Modifiers {
			if modifier.Key != "Cost" {
				continue
			}
			isReplicate = true
			replicateCost, err = cost.New(modifier.Value, spell)
			if err != nil {
				return fmt.Errorf("failed to create cost: %w", err)
			}
		}
	}
	if isReplicate {
		replicateCount, err = player.Agent.EnterNumber(
			fmt.Sprintf("Replicate how many times for %s", replicateCost.Description()),
			mtg.StaticKeywordReplicate,
		)
		if err != nil {
			return fmt.Errorf("failed to confirm replicate cost: %w", err)
		}
	}
	for range replicateCount {
		spellCost = cost.AddCosts(spellCost, replicateCost)
	}
	if err := spellCost.Pay(state, player); err != nil {
		return fmt.Errorf("cannot pay spell cost: %w", err)
	}
	for _, spliceCard := range toSplice {
		card, ok := spliceCard.(*card.Card)
		if !ok {
			return ErrObjectNotCard
		}
		spell.Splice(state, card)
	}
	z, err := mtg.StringToZone(fromZone)
	if err != nil {
		return fmt.Errorf("invalid zone: %s", fromZone)
	}
	c2, err := player.TakeCard(c.ID(), z)
	if err != nil {
		return fmt.Errorf("failed to take card %s from hand: %w", c2.Name(), err)
	}
	if isFlashback {
		spell.Flashback()
	}
	for _, effect := range spell.Effects() {
		fmt.Println("Effect ID: ", effect.ID)
	}
	state.AddToStack(spell)
	if replicateCount > 0 {
		/*
			replicateAbility := BuildReplicateAbility(card, replicateCount)
			state.Stack.Add(replicateAbility)
		*/
	}
	return nil
}

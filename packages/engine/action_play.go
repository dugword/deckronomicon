package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/game/action"
	"deckronomicon/packages/game/card"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/permanent"
	"deckronomicon/packages/game/player"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/find"
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
		targetObj, err = find.FirstBy(available[choiceZone], has.ID(choice.ID))
		if err != nil {
			return ActionResult{}, fmt.Errorf("failed to find card to play: %w", err)
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
	return actionCastSpellFunc(state, player, c, choiceZone)
}

// TODO: Should this be a method? Standardize on Logging
func PlayLand(state *GameState, player *player.Player, card *card.Card) error {
	if card.Match(has.CardType(mtg.CardTypeLand)) {
		if player.LandDrop {
			return mtg.ErrLandAlreadyPlayed
		}
		player.LandDrop = true
	}
	if err := player.RemoveCardFromHand(card.ID()); err != nil {
		return fmt.Errorf(
			"failed to remove card %s from player %s: %w",
			card.Name(),
			player.ID(),
			err,
		)
	}
	perm, err := permanent.NewPermanent(card, state)
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
func actionCastSpellFunc(state *GameState, player *player.Player, card *card.Card, zone string) (result ActionResult, err error) {
	/*
		spell, err := NewSpell(card)
		if err != nil {
			return nil, fmt.Errorf("failed to create spell from %s: %w", card.Name(), err)
		}
		var spellCost Cost
		var isFlashback bool
		// var isSplice bool
		// var spliceCost Cost
		if zone == ZoneGraveyard {
			for _, ability := range spell.StaticAbilities() {
				if ability.ID != game.AbilityKeywordFlashback {
					continue
				}
				isFlashback = true
				for _, modifier := range ability.Modifiers {
					if modifier.Key != "Cost" {
						continue
					}
					spellCost, err = NewCost(modifier.Value, spell)
					if err != nil {
						return nil, fmt.Errorf("failed to create cost: %w", err)
					}
				}
			}
		} else {
			spellCost = spell.ManaCost()
		}
		var toSplice []game.Object
		if spell.HasSubtype(game.SubtypeArcane) {
			spliceCards := FindInZoneBy(
				player.Hand,
				And(
					HasStaticAbility(game.AbilityKeywordSplice),
					HasStaticAbilityModifier(
						game.AbilityKeywordSplice,
						EffectTag{Key: "Onto", Value: "Arcane"},
					),
				),
			)
			// var spliceCosts []Cost
			choices := CreateChoices(spliceCards, ZoneHand)
			chosen, err := player.Agent.ChooseMany(
				"Choose cards to splice onto the spell",
				NewChoiceSource(game.AbilityKeywordSplice, game.AbilityKeywordSplice),
				AddOptionalChoice(choices),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to choose card: %w", err)
			}
			for _, c := range chosen {
				if c.ID == ChoiceNone {
					break
				}
				spliceCard, err := player.Hand.Get(c.ID)
				if err != nil {
					return nil, fmt.Errorf("failed to get splice card from hand: %w", err)
				}
				for _, ability := range spliceCard.StaticAbilities() {
					if ability.ID != game.AbilityKeywordSplice {
						continue
					}
					for _, modifier := range ability.Modifiers {
						if modifier.Key != "Cost" {
							continue
						}
						spliceCost, err := NewCost(modifier.Value, spell)
						if err != nil {
							return nil, fmt.Errorf("failed to create splice cost: %w", err)
						}
						accept, err := player.Agent.Confirm(
							fmt.Sprintf("Splice %s onto %s for %s?",
								spliceCard.Name(),
								spell.Name(),
								spliceCost.Description(),
							),
							NewChoiceSource(game.AbilityKeywordSplice, game.AbilityKeywordSplice),
						)
						if err != nil {
							return nil, fmt.Errorf("failed to confirm splice cost: %w", err)
						}
						if !accept {
							continue
						}
						state.Log("adding splice to toSplice...")
						spellCost = AddCosts(spellCost, spliceCost)
						toSplice = append(toSplice, spliceCard)
					}
				}
			}
		}

		var isReplicate bool
		var replicateCost Cost
		var replicateCount int
		for _, ability := range spell.StaticAbilities() {
			if ability.ID != game.AbilityKeywordReplicate {
				continue
			}
			for _, modifier := range ability.Modifiers {
				if modifier.Key != "Cost" {
					continue
				}
				isReplicate = true
				replicateCost, err = NewCost(modifier.Value, spell)
				if err != nil {
					return nil, fmt.Errorf("failed to create cost: %w", err)
				}
			}
		}
		if isReplicate {
			replicateCount, err = player.Agent.EnterNumber(
				fmt.Sprintf("Replicate how many times for %s", replicateCost.Description()),
				NewChoiceSource(mtg.AbilityKeywordReplicate, mtg.AbilityKeywordReplicate),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to confirm replicate cost: %w", err)
			}
		}
		for range replicateCount {
			spellCost = AddCosts(spellCost, replicateCost)
		}
		fmt.Println("Spell cost: ", spellCost.Description())
		if err := spellCost.Pay(state, player); err != nil {
			return nil, err
		}
		state.Log("toSplice size: " + strconv.Itoa(len(toSplice)))
		for _, spliceCard := range toSplice {
			card, ok := spliceCard.(*card.Card)
			if !ok {
				return nil, fmt.Errorf("object is not a card: %w", err)
			}
			state.Log("Splicing " + card.Name() + " onto " + spell.Name())
			spell.Splice(card)
		}
		cardZone, err := player.GetZone(zone)
		if err != nil {
			return nil, fmt.Errorf("failed to get zone %s: %w", zone, err)
		}
		if err := cardZone.Remove(card.ID()); err != nil {
			return nil, fmt.Errorf("failed to remove card from hand: %w", err)
		}
		state.Log("Casing spell: " + card.Name())
		if isFlashback {
			spell.Flashback()
		}
		for _, effect := range spell.SpellAbility().Effects {
			fmt.Println("Effect ID: ", effect.ID)
		}
		if err := state.Stack.Add(spell); err != nil {
			return nil, fmt.Errorf("failed to add spell to stack: %w", err)
		}
		if replicateCount > 0 {
			replicateAbility := BuildReplicateAbility(card, replicateCount)
			if err := state.Stack.Add(replicateAbility); err != nil {
				return nil, fmt.Errorf("failed to add replicate trigger to stack: %w", err)
			}
		}
		return ActionResult{
			Message: "played card: " + card.Name(),
		}, nil
	*/

	return ActionResult{}, nil
}

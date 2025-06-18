package judge

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"encoding/json"
	"fmt"
)

// TODO: For function signatures think about when to use little > big vs big > little
// e.g. game, playerID, zone, card vs card, zone playerID, game

// TODO: This probably shouldn't be in Judge, maybe judge is just boolean functions
// that return true or false, and this should be where it's used, like in the agent functions
func GetLandsAvailableToPlay(game state.Game, player state.Player, ruling *Ruling) []gob.CardInZone {
	var availableCards []gob.CardInZone
	for _, card := range player.Hand().GetAll() {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("[card %q]: ", card.Name()))
		}
		if CanPlayLand(game, player, mtg.ZoneHand, card, ruling) {

			availableCards = append(availableCards, gob.NewCardInZone(card, mtg.ZoneHand))
		}
	}
	for _, card := range player.Graveyard().GetAll() {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("[card %q]: ", card.Name()))
		}
		if CanPlayLand(game, player, mtg.ZoneGraveyard, card, ruling) {
			availableCards = append(availableCards, gob.NewCardInZone(card, mtg.ZoneGraveyard))
		}
	}
	return availableCards
}

func GetSpellsAvailableToCast(game state.Game, player state.Player, ruling *Ruling) []gob.CardInZone {
	var availableCards []gob.CardInZone
	for _, card := range player.Hand().GetAll() {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("[card %q]: ", card.Name()))
		}
		if CanCastSpellFromHand(game, player, card, card.ManaCost(), ruling) {
			availableCards = append(availableCards, gob.NewCardInZone(card, mtg.ZoneHand))
		}
	}
	for _, card := range player.Graveyard().GetAll() {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("[card %q]: ", card.Name()))
		}
		flashbackCost, ok := card.GetStaticAbilityCost(mtg.StaticKeywordFlashback)
		if !ok {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(
					ruling.Reasons,
					fmt.Sprintf("card %s does not have flashback", card.Name()),
				)
			}
			continue
		}
		if CanCastSpellWithFlashback(game, player, card, flashbackCost, ruling) {
			availableCards = append(availableCards, gob.NewCardInZone(card, mtg.ZoneGraveyard))
		}
	}
	return availableCards
}

// TODO: This probably shouldn't be in Judge, maybe judge is just boolean functions
// that return true or false, and this should be where it's used, like in the agent functions
func GetAbilitiesAvailableToActivate(game state.Game, player state.Player, ruling *Ruling) []gob.AbilityInZone {
	var availableAbilities []gob.AbilityInZone
	for _, permanent := range game.Battlefield().GetAll() {
		for _, ability := range permanent.ActivatedAbilities() {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("[ability %q]: ", ability.Name()))
			}
			if CanActivateAbility(game, player, permanent, ability, ruling) {
				availableAbilities = append(availableAbilities, gob.NewAbilityInZone(ability, permanent, mtg.ZoneBattlefield))
			}
		}
	}
	for _, card := range player.Hand().GetAll() {
		for _, ability := range card.ActivatedAbilities() {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("[ability %q]: ", ability.Name()))
			}
			if ability.Zone() != mtg.ZoneHand {
				if ruling != nil && ruling.Explain {
					ruling.Reasons = append(ruling.Reasons, "ability not available in hand")
				}
				continue
			}
			if CanActivateAbility(game, player, card, ability, ruling) {
				availableAbilities = append(availableAbilities, gob.NewAbilityInZone(ability, card, mtg.ZoneHand))
			}
		}
	}
	return availableAbilities
}

// TODO: Account for Cost
func GetSplicableCards(
	game state.Game,
	player state.Player,
	cardToCast gob.CardInZone,
	ruling *Ruling,
) ([]gob.CardInZone, error) {
	var splicableCards []gob.CardInZone
	for _, card := range player.Hand().GetAll() {
		if card.ID() == cardToCast.ID() {
			continue
		}
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("[card %q]: ", card.Name()))
		}
		var spliceModifiers gob.SpliceModifiers
		rawModifiers, ok := card.GetStaticAbilityModifiers(mtg.StaticKeywordSplice)
		if !ok {
			continue
		}
		if err := json.Unmarshal(rawModifiers, &spliceModifiers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal splice modifiers for card %s: %w", card.Name(), err)
		}
		if !cardToCast.Match(has.Subtype(spliceModifiers.Subtype)) {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(
					ruling.Reasons,
					fmt.Sprintf("cannot splice card %s onto %s: card does not have subtype %s",
						card.Name(),
						cardToCast.Card().Name(),
						spliceModifiers.Subtype,
					),
				)
			}
			continue
		}
		spliceCost, ok := card.GetStaticAbilityCost(mtg.StaticKeywordSplice)
		if !ok {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(
					ruling.Reasons,
					fmt.Sprintf("cannot splice card %s onto %s: card does not have a splice cost",
						card.Name(),
						cardToCast.Card().Name(),
					),
				)
			}
			continue
		}
		totalCost := cost.CombineCosts(
			cardToCast.Card().ManaCost(),
			spliceCost,
		)
		if !CanPayCost(totalCost, card, game, player, ruling) {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(
					ruling.Reasons,
					fmt.Sprintf("cannot splice card %s onto %s: insufficient resources to pay cost %s",
						card.Name(),
						cardToCast.Card().Name(),
						totalCost.Description(),
					),
				)
			}
			continue
		}
		splicableCards = append(splicableCards, gob.NewCardInZone(card, mtg.ZoneHand))
	}
	return splicableCards, nil
}

package judge

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/pay"
	"deckronomicon/packages/engine/reducer"
	"deckronomicon/packages/engine/resolver"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/staticability"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"fmt"
)

// TODO: I really like this idea,
// It feels line with Get Available Cards to play, cast, activate
func GetTargetsForEffect() bool {
	return false
}

func GetLandsAvailableToPlay(game *state.Game, playerID string, ruling *Ruling) []*gob.CardInZone {
	player := game.GetPlayer(playerID)
	var availableCards []*gob.CardInZone
	for _, card := range player.Hand().GetAll() {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("[card %q]: ", card.Name()))
		}
		if CanPlayLand(game, playerID, mtg.ZoneHand, card, ruling) {

			availableCards = append(availableCards, gob.NewCardInZone(card, mtg.ZoneHand))
		}
	}
	for _, card := range player.Graveyard().GetAll() {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("[card %q]: ", card.Name()))
		}
		if CanPlayLand(game, playerID, mtg.ZoneGraveyard, card, ruling) {
			availableCards = append(availableCards, gob.NewCardInZone(card, mtg.ZoneGraveyard))
		}
	}
	return availableCards
}

func GetSpellsAvailableToCast(
	game *state.Game,
	playerID string,
	autoPayCost bool,
	autoPayColors []mana.Color,
	ruling *Ruling,
	apply func(game *state.Game, event event.GameEvent) (*state.Game, error),
) []*gob.CardInZone {
	player := game.GetPlayer(playerID)
	var availableCards []*gob.CardInZone
	for _, card := range player.Hand().GetAll() {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("[card %q]: ", card.Name()))
		}
		if CanCastSpellFromHand(
			game,
			playerID,
			card,
			card.ManaCost(),
			autoPayCost,
			autoPayColors,
			ruling,
			apply,
		) {
			availableCards = append(availableCards, gob.NewCardInZone(card, mtg.ZoneHand))
		}
	}
	for _, card := range player.Graveyard().GetAll() {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("[card %q]: ", card.Name()))
		}
		staticAbility, ok := card.StaticAbility(mtg.StaticKeywordFlashback)
		if !ok {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(
					ruling.Reasons,
					fmt.Sprintf("card %s does not have flashback", card.Name()),
				)
			}
			continue
		}
		flashbackAbility, ok := staticAbility.(*staticability.Flashback)
		if !ok {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(
					ruling.Reasons,
					fmt.Sprintf("card %s does not have flashback ability", card.Name()),
				)
			}
			continue
		}
		if CanCastSpellWithFlashback(game, playerID, card, flashbackAbility.Cost, autoPayCost, autoPayColors, ruling) {
			availableCards = append(availableCards, gob.NewCardInZone(card, mtg.ZoneGraveyard))
		}
	}
	return availableCards
}

// TODO: This probably shouldn't be in Judge, maybe judge is just boolean functions
// that return true or false, and this should be where it's used, like in the agent functions
func GetAbilitiesAvailableToActivate(game *state.Game, playerID string, ruling *Ruling) []*gob.AbilityInZone {
	player := game.GetPlayer(playerID)
	var availableAbilities []*gob.AbilityInZone
	for _, permanent := range game.Battlefield().GetAll() {
		for _, ability := range permanent.ActivatedAbilities() {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("[ability %q]: ", ability.Name()))
			}
			if CanActivateAbility(game, playerID, permanent, ability, ruling) {
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
			if CanActivateAbility(game, playerID, card, ability, ruling) {
				availableAbilities = append(availableAbilities, gob.NewAbilityInZone(ability, card, mtg.ZoneHand))
			}
		}
	}
	return availableAbilities
}

// TODO: Account for Cost
func GetSplicableCards(
	game *state.Game,
	playerID string,
	cardToCast *gob.CardInZone,
	ruling *Ruling,
) ([]*gob.CardInZone, error) {
	player := game.GetPlayer(playerID)
	var splicableCards []*gob.CardInZone
	for _, card := range player.Hand().GetAll() {
		if card.ID() == cardToCast.ID() {
			continue
		}
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("[card %q]: ", card.Name()))
		}
		spliceAbilityRaw, ok := card.StaticAbility(mtg.StaticKeywordSplice)
		if !ok {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(
					ruling.Reasons,
					fmt.Sprintf("cannot splice card %s onto %s: card does not have splice ability",
						card.Name(),
						cardToCast.Card().Name(),
					),
				)
			}
			continue
		}
		spliceAbility, ok := spliceAbilityRaw.(*staticability.Splice)
		if !ok {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(
					ruling.Reasons,
					fmt.Sprintf("cannot splice card %s onto %s: splice ability is not a Splice ability",
						card.Name(),
						cardToCast.Card().Name(),
					),
				)
			}
			continue
		}
		if !cardToCast.Match(has.Subtype(spliceAbility.Subtype)) {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(
					ruling.Reasons,
					fmt.Sprintf("cannot splice card %s onto %s: card does not have subtype %s",
						card.Name(),
						cardToCast.Card().Name(),
						spliceAbility.Subtype,
					),
				)
			}
			continue
		}
		totalCost := cost.NewComposite(
			cardToCast.Card().ManaCost(),
			spliceAbility.Cost,
		)
		if !CanPayCost(totalCost, card, game, playerID, ruling) {
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

// TODO: Redundate with pay automatic activation
func GetAvailableMana(game *state.Game, playerID string) mana.Pool {
	for _, untappedLand := range game.Battlefield().FindAll(
		query.And(has.Controller(playerID), is.Land(), is.Untapped())) {
		for _, ability := range untappedLand.ActivatedAbilities() {
			if !ability.Match(is.ManaAbility()) {
				continue
			}
			events := []event.GameEvent{
				&event.ActivateAbilityEvent{
					PlayerID:  playerID,
					SourceID:  untappedLand.ID(),
					AbilityID: ability.Name(),
					Zone:      mtg.ZoneBattlefield,
				},
			}
			costEvents, err := pay.Cost(
				ability.Cost(),
				untappedLand,
				playerID,
			)
			if err != nil {
				// TODO: Don't panic
				panic(fmt.Errorf("failed to pay cost for ability %q: %w", ability.Name(), err))
			}
			events = append(events, costEvents...)
			events = append(events, &event.LandTappedForManaEvent{
				PlayerID: playerID,
				ObjectID: untappedLand.ID(),
				Subtypes: untappedLand.Subtypes(),
			})
			for _, event := range events {
				var err error
				game, err = reducer.ApplyEventAndTriggers(game, event)
				if err != nil {
					panic(fmt.Errorf("failed to apply event %q: %w", event.EventType(), err))
				}
			}
			for _, efct := range ability.Effects() {
				addMana, ok := efct.(*effect.AddMana)
				if !ok {
					continue
				}
				result, err := resolver.ResolveAddMana(game, playerID, addMana)
				if err != nil {
					panic(fmt.Errorf("failed to resolve add mana effect: %w", err))
				}
				for _, event := range result.Events {
					game, err = reducer.ApplyEventAndTriggers(game, event)
					if err != nil {
						panic(fmt.Errorf("failed to apply event %q: %w", event.EventType(), err))
					}
				}
				events = append(events, result.Events...)
			}
		}
	}
	player := game.GetPlayer(playerID)
	return player.ManaPool()
}

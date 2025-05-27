package auto

import (
	game "deckronomicon/game"
)

// MatchesConditionSet determines if a given ConditionSet is satisfied by the game state.
// It checks various conditions such as land drops, mana availability, card
// presence in hand, battlefield, and graveyard, as well as game state scalars
// like storm count, library size, and cards in hand.
func MatchesConditionSet(state *game.GameState, player *game.Player, cond ConditionSet) (bool, error) {
	// --- Land Drop Checks ---
	if cond.HasPlayedLand != nil && *cond.HasPlayedLand != player.LandDrop {
		return false, nil
	}
	// --- Mana Cost ---
	if cond.ManaAvailable != "" {
		manaNeeded, err := game.ParseManaCost(cond.ManaAvailable)
		if err != nil {
			return false, err
		}
		if !game.CanPotentiallyPayFor(state, manaNeeded) {
			return false, nil
		}
	}
	// --- Hand Checks ---
	if !allCardsPresent(cond.HandContains, player.Hand.GetAll()) {
		return false, nil
	}
	if len(cond.HandContainsAny) > 0 && !anyCardPresent(cond.HandContainsAny, player.Hand.GetAll()) {
		return false, nil
	}
	if !allGroupsSatisfied(cond.HandContainsAllGroups, player.Hand.GetAll()) {
		return false, nil
	}
	if len(cond.HandContainsAnyGroups) > 0 && !anyGroupSatisfied(cond.HandContainsAnyGroups, player.Hand.GetAll()) {
		return false, nil
	}
	if !allCardsAbsent(cond.HandLacks, player.Hand.GetAll()) {
		return false, nil
	}
	if len(cond.HandLacksAny) > 0 && !anyCardAbsent(cond.HandLacksAny, player.Hand.GetAll()) {
		return false, nil
	}
	if !allGroupsAbsent(cond.HandLacksAllGroups, player.Hand.GetAll()) {
		return false, nil
	}
	if !noGroupFullyPresent(cond.HandLacksAnyGroups, player.Hand.GetAll()) {
		return false, nil
	}

	/*
		// --- Battlefield Checks ---
		if !allCardsPresent(cond.BattlefieldContains, cardsFromPerms(player.Battlefield)) {
			return false, nil
		}

		if len(cond.BattlefieldContainsAny) > 0 && !anyCardPresent(cond.BattlefieldContainsAny, cardsFromPerms(player.Battlefield)) {
			return false, nil
		}
		if !allGroupsSatisfied(cond.BattlefieldContainsAllGroups, cardsFromPerms(player.Battlefield)) {
			return false, nil
		}
		if len(cond.BattlefieldContainsAnyGroups) > 0 && !anyGroupSatisfied(cond.BattlefieldContainsAnyGroups, cardsFromPerms(player.Battlefield)) {
			return false, nil
		}
		if !allCardsAbsent(cond.BattlefieldLacks, cardsFromPerms(player.Battlefield)) {
			return false, nil
		}
		if len(cond.BattlefieldLacksAny) > 0 && !anyCardAbsent(cond.BattlefieldLacksAny, cardsFromPerms(player.Battlefield)) {
			return false, nil
		}
		if !allGroupsAbsent(cond.BattlefieldLacksAllGroups, cardsFromPerms(player.Battlefield)) {
			return false, nil
		}
		if !noGroupFullyPresent(cond.BattlefieldLacksAnyGroups, cardsFromPerms(player.Battlefield)) {
			return false, nil
		}
	*/

	// --- Graveyard Checks ---
	if !allCardsPresent(cond.GraveyardContains, player.Graveyard.GetAll()) {
		return false, nil
	}
	if len(cond.GraveyardContainsAny) > 0 && !anyCardPresent(cond.GraveyardContainsAny, player.Graveyard.GetAll()) {
		return false, nil
	}
	if !allGroupsSatisfied(cond.GraveyardContainsAllGroups, player.Graveyard.GetAll()) {
		return false, nil
	}
	if len(cond.GraveyardContainsAnyGroups) > 0 && !anyGroupSatisfied(cond.GraveyardContainsAnyGroups, player.Graveyard.GetAll()) {
		return false, nil
	}
	if !allCardsAbsent(cond.GraveyardLacks, player.Graveyard.GetAll()) {
		return false, nil
	}
	if len(cond.GraveyardLacksAny) > 0 && !anyCardAbsent(cond.GraveyardLacksAny, player.Graveyard.GetAll()) {
		return false, nil
	}
	if !allGroupsAbsent(cond.GraveyardLacksAllGroups, player.Graveyard.GetAll()) {
		return false, nil
	}
	if !noGroupFullyPresent(cond.GraveyardLacksAnyGroups, player.Graveyard.GetAll()) {
		return false, nil
	}

	// --- Game State Scalar Conditions ---
	if cond.Storm != "" && !evaluateIntComparison(state.StormCount, cond.Storm) {
		return false, nil
	}
	if cond.LibrarySize != "" && !evaluateIntComparison(player.Library.Size(), cond.LibrarySize) {
		return false, nil
	}
	if cond.CardsInHand != "" && !evaluateIntComparison(player.Hand.Size(), cond.CardsInHand) {
		return false, nil
	}
	if cond.GraveyardSize != "" && !evaluateIntComparison(player.Graveyard.Size(), cond.GraveyardSize) {
		return false, nil
	}
	if cond.SpellCountThisTurn != "" && !evaluateIntComparison(state.StormCount, cond.SpellCountThisTurn) {
		return false, nil
	}
	/*
		for _, name := range cond.HasCastThisTurn {
			if !state.SpellsCastThisTurn[name] {
				return false, nil
			}
		}
	*/
	return true, nil
}

// TODO: Future Enhancements for Rule Matching
//
// - Mana Curve Logic:
//   Allow conditions like `has_card_costing: 1`, `max_card_cost_in_hand: 3`, or
//   average CMC thresholds to support deck tempo and combo-timing logic.
//
// - Spell History:
//   Track which spells were cast this turn (or past turns), e.g.:
//   `has_cast_this_turn: ["High Tide"]`
//   or more complex logic like `cast_count_this_turn: {"High Tide": ">=2"}`.
//
// - Tag-Based Filtering:
//   Add support for card tags like `"arcane"`, `"mana_engine"`, etc. in rules.
//   Example: `hand_contains_tags: ["arcane"]`
//   Useful for generalizing across similar card types.
//
// - User-Defined Groups:
//   Let rules JSON define reusable named groups to simplify repeated conditions.
//   Example:
//     "groups": {
//       "combo_core": ["High Tide", "Psychic Puppetry"],
//       "arcane_package": ["Ideas Unbound", "Peer through Depths"]
//     }
//   Then use: `hand_contains_all_groups: ["combo_core", "arcane_package"]`
//   This avoids copy/paste and makes the rule config more maintainable.

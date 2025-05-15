package auto

import (
	game "deckronomicon/game"
)

// MatchesConditionSet determines if a given ConditionSet is satisfied by the game state.
func MatchesConditionSet(state *game.GameState, cond ConditionSet) bool {
	// --- Land Drop Checks ---
	if cond.HasPlayedLand != nil && *cond.HasPlayedLand != state.LandDrop {
		return false
	}
	// --- Mana Cost ---
	if cond.ManaAvailable != "" {
		//manaNeeded := parseManaCost(cond.ManaAvailable)
		/*
			if !game.CanPotentiallyPayFor(state, manaNeeded) {
				return false
			}
		*/
	}
	// --- Hand Checks ---
	if !allCardsPresent(cond.HandContains, state.Hand.Cards()) {
		return false
	}
	if len(cond.HandContainsAny) > 0 && !anyCardPresent(cond.HandContainsAny, state.Hand.Cards()) {
		return false
	}
	if !allGroupsSatisfied(cond.HandContainsAllGroups, state.Hand.Cards()) {
		return false
	}
	if len(cond.HandContainsAnyGroups) > 0 && !anyGroupSatisfied(cond.HandContainsAnyGroups, state.Hand.Cards()) {
		return false
	}
	if !allCardsAbsent(cond.HandLacks, state.Hand.Cards()) {
		return false
	}
	if len(cond.HandLacksAny) > 0 && !anyCardAbsent(cond.HandLacksAny, state.Hand.Cards()) {
		return false
	}
	if !allGroupsAbsent(cond.HandLacksAllGroups, state.Hand.Cards()) {
		return false
	}
	if !noGroupFullyPresent(cond.HandLacksAnyGroups, state.Hand.Cards()) {
		return false
	}

	/*
		// --- Battlefield Checks ---
		if !allCardsPresent(cond.BattlefieldContains, cardsFromPerms(state.Battlefield)) {
			return false
		}

		if len(cond.BattlefieldContainsAny) > 0 && !anyCardPresent(cond.BattlefieldContainsAny, cardsFromPerms(state.Battlefield)) {
			return false
		}
		if !allGroupsSatisfied(cond.BattlefieldContainsAllGroups, cardsFromPerms(state.Battlefield)) {
			return false
		}
		if len(cond.BattlefieldContainsAnyGroups) > 0 && !anyGroupSatisfied(cond.BattlefieldContainsAnyGroups, cardsFromPerms(state.Battlefield)) {
			return false
		}
		if !allCardsAbsent(cond.BattlefieldLacks, cardsFromPerms(state.Battlefield)) {
			return false
		}
		if len(cond.BattlefieldLacksAny) > 0 && !anyCardAbsent(cond.BattlefieldLacksAny, cardsFromPerms(state.Battlefield)) {
			return false
		}
		if !allGroupsAbsent(cond.BattlefieldLacksAllGroups, cardsFromPerms(state.Battlefield)) {
			return false
		}
		if !noGroupFullyPresent(cond.BattlefieldLacksAnyGroups, cardsFromPerms(state.Battlefield)) {
			return false
		}
	*/

	// --- Graveyard Checks ---
	if !allCardsPresent(cond.GraveyardContains, state.Graveyard) {
		return false
	}
	if len(cond.GraveyardContainsAny) > 0 && !anyCardPresent(cond.GraveyardContainsAny, state.Graveyard) {
		return false
	}
	if !allGroupsSatisfied(cond.GraveyardContainsAllGroups, state.Graveyard) {
		return false
	}
	if len(cond.GraveyardContainsAnyGroups) > 0 && !anyGroupSatisfied(cond.GraveyardContainsAnyGroups, state.Graveyard) {
		return false
	}
	if !allCardsAbsent(cond.GraveyardLacks, state.Graveyard) {
		return false
	}
	if len(cond.GraveyardLacksAny) > 0 && !anyCardAbsent(cond.GraveyardLacksAny, state.Graveyard) {
		return false
	}
	if !allGroupsAbsent(cond.GraveyardLacksAllGroups, state.Graveyard) {
		return false
	}
	if !noGroupFullyPresent(cond.GraveyardLacksAnyGroups, state.Graveyard) {
		return false
	}

	// --- Game State Scalar Conditions ---
	if cond.Storm != "" && !evaluateIntComparison(state.StormCount, cond.Storm) {
		return false
	}
	if cond.LibrarySize != "" && !evaluateIntComparison(state.Deck.Size(), cond.LibrarySize) {
		return false
	}
	if cond.CardsInHand != "" && !evaluateIntComparison(state.Hand.Size(), cond.CardsInHand) {
		return false
	}
	if cond.GraveyardSize != "" && !evaluateIntComparison(len(state.Graveyard), cond.GraveyardSize) {
		return false
	}
	if cond.SpellCountThisTurn != "" && !evaluateIntComparison(state.StormCount, cond.SpellCountThisTurn) {
		return false
	}
	/*
		for _, name := range cond.HasCastThisTurn {
			if !state.SpellsCastThisTurn[name] {
				return false
			}
		}
	*/
	return true
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

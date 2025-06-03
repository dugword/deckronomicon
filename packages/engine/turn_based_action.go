package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
)

/*
From the Comprehensive Rules (November 8, 2024—Magic: The Gathering Foundations)
703. Turn-Based Actions
703.1. Turn-based actions are game actions that happen automatically when certain steps or phases begin, or when each step and phase ends. Turn-based actions don’t use the stack.
703.1a Abilities that watch for a specified step or phase to begin are triggered abilities, not turn-based actions. (See rule 603, “Handling Triggered Abilities.”)
703.2. Turn-based actions are not controlled by any player.
703.3. Whenever a step or phase begins, if it’s a step or phase that has any turn-based action associated with it, those turn-based actions are automatically dealt with first. This happens before state-based actions are checked, before triggered abilities are put on the stack, and before players receive priority.
703.4. The turn-based actions are as follows:
703.4a Immediately after the untap step begins, all phased-in permanents with phasing that the active player controls phase out, and all phased-out permanents that the active player controlled when they phased out phase in. This all happens simultaneously. See rule 502.1.
703.4b Immediately after the phasing action has been completed during the untap step, if the game has either the day or night designation, it checks to see whether that designation should change. If it’s neither day nor night, this check doesn’t happen. See rule 502.2.
703.4c Immediately after the game checks to see if its day or night designation should change during the untap step or, if the game doesn’t have a day or night designation, immediately after the phasing action has been completed during the untap step, the active player determines which permanents they control will untap. Then they untap them all simultaneously. See rule 502.3.
703.4d Immediately after the draw step begins, the active player draws a card. See rule 504.1.
703.4e In an Archenemy game (see rule 904), immediately after the archenemy’s precombat main phase begins, that player sets the top card of their scheme deck in motion. See rule 701.25.
703.4f Immediately after a player’s precombat main phase begins, that player puts a lore counter on each Saga enchantment they control. In an Archenemy game, this happens after the archenemy’s scheme action. See rule 714, “Saga Cards.”
703.4g Immediately after the action of placing lore counters has been completed, if the active player controls any Attractions, that player rolls to visit their Attractions. See rule 701.49, “Roll to Visit Your Attractions.”
703.4h Immediately after the beginning of combat step begins, if the game being played is a multiplayer game in which the active player’s opponents don’t all automatically become defending players, the active player chooses one of their opponents. That player becomes the defending player. See rule 507.1.
703.4i Immediately after the declare attackers step begins, the active player declares attackers. See rule 508.1.
703.4j Immediately after the declare blockers step begins, the defending player declares blockers. See rule 509.1.
703.4k Immediately after the combat damage step begins, each player in APNAP order announces how each attacking or blocking creature they control assigns its combat damage. See rule 510.1.
703.4m Immediately after combat damage has been assigned during the combat damage step, all combat damage is dealt simultaneously. See rule 510.2.
703.4n Immediately after the cleanup step begins, if the active player’s hand contains more cards than their maximum hand size (normally seven), they discard enough cards to reduce their hand size to that number. See rule 514.1.
703.4p Immediately after the active player has discarded cards (if necessary) during the cleanup step, all damage is removed from permanents and all “until end of turn” and “this turn” effects end. These actions happen simultaneously. See rule 514.2.
703.4q When each step or phase ends, any unused mana left in a player’s mana pool empties. See rule 500.4
*/

type TurnBasedAction interface {
	Name() string
	Description() string
	Complete(choose.Choice) (event.GameEvent, error)
}

// DrawAction represents the action of drawing a card during a player's turn.
type DrawAction struct {
	PlayerID string
}

func (a DrawAction) Name() string {
	return "Draw Card"
}

func (a DrawAction) Description() string {
	return "The active player draws a card."
}

func (a DrawAction) GetPrompt(state GameState, player Player) choose.ChoicePrompt {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Drawing a card",
		Choices:  nil,
		Optional: false,
	}
}
func (a DrawAction) Complete(choose.Choice) (event.GameEvent, error) {
	return event.GameEvent{}, nil
	// return event.NewDrawCardEvent(a.PlayerID), nil
}

// UntapAction represents the action of untapping permanents during the untap
// step.
type UntapAction struct {
	PlayerID string
}

func (a UntapAction) Name() string {
	return "Untap Permanents"
}

func (a UntapAction) Description() string {
	return "The active player untaps all permanents they control."
}

func (a UntapAction) GetPrompt(state GameState, player Player) choose.ChoicePrompt {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Untapping permanents",
		Choices:  nil,
		Optional: false,
	}
}

func (a UntapAction) Complete(choose.Choice) (event.GameEvent, error) {
	return event.GameEvent{}, nil
}

type DeclareAttackersAction struct {
	PlayerID string
}

func (a DeclareAttackersAction) Name() string {
	return "Declare Attackers"
}
func (a DeclareAttackersAction) Description() string {
	return "The active player declares attackers."
}
func (a DeclareAttackersAction) GetPrompt(state GameState, player Player) choose.ChoicePrompt {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Declaring attackers",
		Choices:  nil,
		Optional: false,
	}
}
func (a DeclareAttackersAction) Complete(choose.Choice) (event.GameEvent, error) {
	// This action would typically involve the player choosing which creatures to attack with.
	// For now, we return an empty event as a placeholder.
	return event.GameEvent{}, nil
}

type DeclareBlockersAction struct {
	PlayerID string
}

func (a DeclareBlockersAction) Name() string {
	return "Declare Blockers"
}
func (a DeclareBlockersAction) Description() string {
	return "The defending player declares blockers."
}
func (a DeclareBlockersAction) GetPrompt(state GameState, player Player) choose.ChoicePrompt {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Declaring blockers",
		Choices:  nil,
		Optional: false,
	}
}

func (a DeclareBlockersAction) Complete(choose.Choice) (event.GameEvent, error) {
	// This action would typically involve the player choosing which creatures to block with.
	// For now, we return an empty event as a placeholder.
	return event.GameEvent{}, nil
}

// CombatDamageAction represents the action of assigning combat damage during
// the combat damage step.
type CombatDamageAction struct {
	PlayerID string
}

func (a CombatDamageAction) Name() string {
	return "Assign Combat Damage"
}
func (a CombatDamageAction) Description() string {
	return "The active player assigns combat damage for attacking and blocking creatures."
}
func (a CombatDamageAction) GetPrompt(state GameState, player Player) choose.ChoicePrompt {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Assigning combat damage",
		Choices:  nil,
		Optional: false,
	}
}
func (a CombatDamageAction) Complete(choose.Choice) (event.GameEvent, error) {
	// This action would typically involve the player assigning combat damage.
	// For now, we return an empty event as a placeholder.
	return event.GameEvent{}, nil
}

// CleanupAction represents the action of cleaning up at the end of a turn.
type CleanupAction struct {
	PlayerID string
}

func (a CleanupAction) Name() string {
	return "Cleanup"
}
func (a CleanupAction) Description() string {
	return "The active player discards down to their maximum hand size and removes all damage from permanents."
}
func (a CleanupAction) GetPrompt(state GameState, player Player) choose.ChoicePrompt {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Cleaning up",
		Choices:  nil,
		Optional: false,
	}
}
func (a CleanupAction) Complete(choose.Choice) (event.GameEvent, error) {
	// This action would typically involve the player discarding cards and removing damage.
	// For now, we return an empty event as a placeholder.
	return event.GameEvent{}, nil
}

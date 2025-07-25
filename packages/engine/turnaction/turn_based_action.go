package turnaction

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"fmt"
)

// TODO: Not sure if these actually should be "actions", or if turn based actions should just be
// events, and "action" package is just for player-initiated actions.

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

// technically draw starting hand isn't a turn-based action, but it is a game
// lifecycle event

type DrawStartingHandAction struct {
	playerID string
}

func NewDrawStartingHandAction(playerID string) DrawStartingHandAction {
	return DrawStartingHandAction{
		playerID: playerID,
	}
}

func (a DrawStartingHandAction) PlayerID() string {
	return a.playerID
}

func (a DrawStartingHandAction) Name() string {
	return "Draw Starting Hand"
}
func (a DrawStartingHandAction) Description() string {
	return "The active player draws their starting hand."
}

func (a DrawStartingHandAction) GetPrompt(
	game *state.Game,
) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{}, nil
}

func (a DrawStartingHandAction) Complete(
	game *state.Game,
	choiceResults choose.ChoiceResults,
) ([]event.GameEvent, error) {
	// This action would typically involve the player drawing their starting hand.
	// For now, we return an empty event as a placeholder.
	player := game.GetPlayer(a.PlayerID())
	// Draw the starting hand, which is typically 7 cards, but can be
	// different.
	var drawEvents []event.GameEvent
	for range player.MaxHandSize() {
		drawEvents = append(drawEvents, &event.DrawCardEvent{PlayerID: a.PlayerID()})
	}
	return append([]event.GameEvent{
			&event.DrawStartingHandEvent{PlayerID: a.PlayerID()},
		}, drawEvents...),
		nil
}

type PhaseInPhaseOutAction struct {
	playerID string
}

func NewPhaseInPhaseOutAction(playerID string) PhaseInPhaseOutAction {
	return PhaseInPhaseOutAction{
		playerID: playerID,
	}
}

func (a PhaseInPhaseOutAction) PlayerID() string {
	return a.playerID
}

func (a PhaseInPhaseOutAction) Name() string {
	return "Phase In/Out"
}

func (a PhaseInPhaseOutAction) Description() string {
	return "Phased-in permanents with phasing phase out, and phased-out permanents phase in."
}

func (a PhaseInPhaseOutAction) GetPrompt(game *state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{}, nil
}

func (a PhaseInPhaseOutAction) Complete(game *state.Game, choiceResults choose.ChoiceResults) ([]event.GameEvent, error) {
	return []event.GameEvent{
		&event.PhaseInPhaseOutEvent{PlayerID: a.PlayerID()},
	}, nil
}

type CheckDayNightAction struct {
	playerID string
}

func NewCheckDayNightAction(playerID string) CheckDayNightAction {
	return CheckDayNightAction{
		playerID: playerID,
	}
}

func (a CheckDayNightAction) PlayerID() string {
	return a.playerID
}

func (a CheckDayNightAction) Name() string {
	return "Check Day/Night"
}

func (a CheckDayNightAction) Description() string {
	return "Check if the day/night designation should change."
}

func (a CheckDayNightAction) GetPrompt(game *state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{}, nil
}

func (a CheckDayNightAction) Complete(game *state.Game, choiceResults choose.ChoiceResults) ([]event.GameEvent, error) {
	// This action would typically involve checking the game state to see if the day/night designation should change.
	// For now, we return an empty event as a placeholder.
	return []event.GameEvent{
		&event.CheckDayNightEvent{PlayerID: a.PlayerID()},
	}, nil
}

// UntapAction represents the action of untapping permanents during the untap
// step.
type UntapAction struct {
	playerID string
}

func NewUntapAction(playerID string) UntapAction {
	return UntapAction{
		playerID: playerID,
	}
}

func (a UntapAction) PlayerID() string {
	return a.playerID
}

func (a UntapAction) Name() string {
	return "Untap Permanents"
}

func (a UntapAction) Description() string {
	return "The active player untaps all permanents they control."
}

func (a UntapAction) GetPrompt(game *state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{}, nil
}

func (a UntapAction) Complete(*state.Game, choose.ChoiceResults) ([]event.GameEvent, error) {
	return []event.GameEvent{
		&event.UntapAllEvent{PlayerID: a.PlayerID()},
	}, nil
}

type UpkeepAction struct {
	playerID string
}

func NewUpkeepAction(playerID string) UpkeepAction {
	return UpkeepAction{
		playerID: playerID,
	}
}

func (a UpkeepAction) PlayerID() string {
	return a.playerID
}

func (a UpkeepAction) Name() string {
	return "Upkeep"
}

func (a UpkeepAction) Description() string {
	return "The active player performs any upkeep actions."
}

func (a UpkeepAction) GetPrompt(game *state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{}, nil
}

func (a UpkeepAction) Complete(game *state.Game, choiceResults choose.ChoiceResults) ([]event.GameEvent, error) {
	// This action would typically involve the player performing any upkeep actions.
	// For now, we return an empty event as a placeholder.
	return []event.GameEvent{
		&event.UpkeepEvent{PlayerID: a.PlayerID()},
	}, nil
}

// DrawAction represents the action of drawing a card during a player's turn.
type DrawAction struct {
	playerID string
}

func NewDrawAction(playerID string) DrawAction {
	return DrawAction{
		playerID: playerID,
	}
}

func (a DrawAction) PlayerID() string {
	return a.playerID
}

func (a DrawAction) Name() string {
	return "Draw Card"
}

func (a DrawAction) Description() string {
	return "The active player draws a card."
}

func (a DrawAction) GetPrompt(game *state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{}, nil
}
func (a DrawAction) Complete(game *state.Game, choose choose.ChoiceResults) ([]event.GameEvent, error) {
	// TODO: Make this more better
	// Player skips draw on their first turn if they go first
	player := game.GetPlayer(a.PlayerID())
	opponent := game.GetOpponent(a.PlayerID())
	if player.Turn() == 1 && opponent.Turn() == 0 {
		return nil, nil
	}
	return []event.GameEvent{
		&event.DrawCardEvent{
			PlayerID: a.PlayerID(),
		}}, nil
	// return event.NewDrawCardEvent(a.PlayerID), nil
}

type DeclareAttackersAction struct {
	playerID string
}

func NewDeclareAttackersAction(playerID string) DeclareAttackersAction {
	return DeclareAttackersAction{
		playerID: playerID,
	}
}

func (a DeclareAttackersAction) PlayerID() string {
	return a.playerID
}

func (a DeclareAttackersAction) Name() string {
	return "Declare Attackers"
}
func (a DeclareAttackersAction) Description() string {
	return "The active player declares attackers."
}
func (a DeclareAttackersAction) GetPrompt(game *state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{}, nil
}
func (a DeclareAttackersAction) Complete(*state.Game, choose.ChoiceResults) ([]event.GameEvent, error) {
	// This action would typically involve the player choosing which creatures to attack with.
	// For now, we return an empty event as a placeholder.
	return []event.GameEvent{
		&event.DeclareAttackersEvent{},
	}, nil
}

type DeclareBlockersAction struct {
	playerID string
}

func NewDeclareBlockersAction(playerID string) DeclareBlockersAction {
	return DeclareBlockersAction{
		playerID: playerID,
	}
}

func (a DeclareBlockersAction) PlayerID() string {
	return a.playerID
}

func (a DeclareBlockersAction) Name() string {
	return "Declare Blockers"
}
func (a DeclareBlockersAction) Description() string {
	return "The defending player declares blockers."
}
func (a DeclareBlockersAction) GetPrompt(game *state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{}, nil
}

func (a DeclareBlockersAction) Complete(*state.Game, choose.ChoiceResults) ([]event.GameEvent, error) {
	// This action would typically involve the player choosing which creatures to block with.
	// For now, we return an empty event as a placeholder.
	return []event.GameEvent{
		&event.DeclareBlockersEvent{},
	}, nil
}

// CombatDamageAction represents the action of assigning combat damage during
// the combat damage step.
type CombatDamageAction struct {
	playerID string
}

func NewCombatDamageAction(playerID string) CombatDamageAction {
	return CombatDamageAction{
		playerID: playerID,
	}
}

func (a CombatDamageAction) PlayerID() string {
	return a.playerID
}

func (a CombatDamageAction) Name() string {
	return "Assign Combat Damage"
}
func (a CombatDamageAction) Description() string {
	return "The active player assigns combat damage for attacking and blocking creatures."
}
func (a CombatDamageAction) GetPrompt(game *state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{}, nil
}
func (a CombatDamageAction) Complete(*state.Game, choose.ChoiceResults) ([]event.GameEvent, error) {
	// This action would typically involve the player assigning combat damage.
	// For now, we return an empty event as a placeholder.
	return []event.GameEvent{
		&event.AssignCombatDamageEvent{PlayerID: a.PlayerID()},
	}, nil
}

type DiscardToHandSizeAction struct {
	playerID string
}

func NewDiscardToHandSizeAction(playerID string) DiscardToHandSizeAction {
	return DiscardToHandSizeAction{
		playerID: playerID,
	}
}

func (a DiscardToHandSizeAction) PlayerID() string {
	return a.playerID
}

func (a DiscardToHandSizeAction) Name() string {
	return "Discard to Hand Size"
}

func (a DiscardToHandSizeAction) Description() string {
	return "The active player discards down to their maximum hand size."
}

// TODO This should move to the player agent
func (a DiscardToHandSizeAction) GetPrompt(game *state.Game) (choose.ChoicePrompt, error) {
	player := game.GetPlayer(a.PlayerID())
	hand := player.Hand()
	excess := hand.Size() - player.MaxHandSize()
	if excess <= 0 {
		return choose.ChoicePrompt{}, nil
	}
	choices := []choose.Choice{}
	for _, card := range hand.GetAll() {
		choices = append(choices, gob.NewCardInZone(card, mtg.ZoneHand))
	}
	return choose.ChoicePrompt{
		Message: fmt.Sprintf("%d cards in hand, discard %d", hand.Size(), excess),
		Source:  a,
		ChoiceOpts: choose.ChooseManyOpts{
			Choices: choices,
			Min:     excess,
			Max:     excess,
		},
	}, nil
}

func (a DiscardToHandSizeAction) Complete(
	game *state.Game,
	choiceResults choose.ChoiceResults,
) ([]event.GameEvent, error) {
	var discardEvents []event.GameEvent
	if choiceResults == nil {
		return nil, nil
	}
	selected, ok := choiceResults.(choose.ChooseManyResults)
	if !ok {
		return nil, fmt.Errorf("expected multiple choice results")
	}
	// TODO: Verify that the choices are valid cards in the player's hand
	for _, choice := range selected.Choices {
		if choice.ID() == "" {
			continue // Skip empty choices
		}
		discardEvents = append(discardEvents, &event.DiscardCardEvent{
			PlayerID: a.PlayerID(),
			CardID:   choice.ID(),
		})
	}
	return discardEvents, nil
}

type RemoveDamageAction struct {
	playerID string
}

func NewRemoveDamageAction(playerID string) RemoveDamageAction {
	return RemoveDamageAction{
		playerID: playerID,
	}
}
func (a RemoveDamageAction) PlayerID() string {
	return a.playerID
}

func (a RemoveDamageAction) Name() string {
	return "Remove Damage"
}

func (a RemoveDamageAction) Description() string {
	return "Remove all damage from permanents and end all 'until end of turn' effects."
}

func (a RemoveDamageAction) GetPrompt(game *state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{}, nil
}

func (a RemoveDamageAction) Complete(game *state.Game, choiceResults choose.ChoiceResults) ([]event.GameEvent, error) {
	// This action would typically involve removing all damage from permanents and ending all "until end of turn" effects.
	// For now, we return an empty event as a placeholder.
	return []event.GameEvent{
		&event.RemoveDamageEvent{PlayerID: a.PlayerID()},
	}, nil
}

type ProgressSagaAction struct {
	playerID string
}

func NewProgressSagaAction(playerID string) ProgressSagaAction {
	return ProgressSagaAction{
		playerID: playerID,
	}
}
func (a ProgressSagaAction) PlayerID() string {
	return a.playerID
}

func (a ProgressSagaAction) Name() string {
	return "Progress Saga"
}

func (a ProgressSagaAction) Description() string {
	return "The active player progresses each Saga they control."
}

func (a ProgressSagaAction) GetPrompt(game *state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{}, nil
}

func (a ProgressSagaAction) Complete(*state.Game, choose.ChoiceResults) ([]event.GameEvent, error) {
	return []event.GameEvent{
		&event.ProgressSagaEvent{PlayerID: a.PlayerID()},
	}, nil
}

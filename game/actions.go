package game

import (
	"errors"
	"fmt"
	"strconv"
)

// TODO: This could theoretically conflict with a card name.
// Perhaps a standard special character like "!" or "@" could be used.
const UntapAll = "UntapAll"

// GameAction represents an action a player can take.
type GameAction struct {
	Cheat      GameCheatType
	Preactions []GameAction
	Target     string
	Type       GameActionType
}

// GameActionType represents the type of player action.
type GameActionType string

// Constants for game action types.
const (
	ActionActivate GameActionType = "Activate"
	ActionDiscard  GameActionType = "Discard"
	ActionDraw     GameActionType = "Draw"
	ActionCheat    GameActionType = "Cheat"
	ActionConcede  GameActionType = "Concede"
	ActionPass     GameActionType = "Pass"
	ActionPlay     GameActionType = "Play"
	ActionUntap    GameActionType = "Untap"
	ActionView     GameActionType = "View"
)

// PlayerActions is a map of game action types to booleans indicating if the
// action is a player action.
var PlayerActions = map[GameActionType]bool{
	ActionActivate: true,
	ActionCheat:    true,
	ActionConcede:  true,
	ActionPass:     true,
	ActionPlay:     true,
	ActionView:     true,
}

// ActionResult represents the result of an action.
type ActionResult struct {
	Message string
	Pass    bool
}

// ActionActivateFunc handles the activate action. This is performed by the
// player to activate an ability of a permanent on the battlefield, or a card
// in hand or in the graveyard. The target is the name of the permanent or
// card.
// TODO: Support more than one activated ability
// TODO: Support activated abilities in hand and graveyard
func ActionActivateFunc(state *GameState, target string, resolver ChoiceResolver) (*ActionResult, error) {
	var selectedPermanent *Permanent
	if target != "" {
		selectedPermanent = state.Battlefield.FindPermanentWithAvailableActivatedAbility(target, state)
	}
	if selectedPermanent == nil {
		choices := state.Battlefield.GetPermanentsWithActivatedAbilities(state)
		choice, err := resolver.ChooseOne(
			"Which permanent to activate",
			string(ActionActivate),
			choices,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose permanent: %w", err)
		}
		selectedPermanent = state.Battlefield.GetPermanent(choice.Index)
	}
	// TODO: Support more than 1
	if len(selectedPermanent.ActivatedAbilities()) > 1 {
		return nil, fmt.Errorf("no support for multiple activated abilities")
	}
	activatedAbility := selectedPermanent.ActivatedAbilities()[0]
	if err := activatedAbility.Cost.Pay(state, resolver); err != nil {
		return nil, fmt.Errorf("cannot pay activated ability cost: %w", err)
	}
	if err := activatedAbility.Resolve(state, resolver); err != nil {
		return nil, fmt.Errorf("cannot resolve ability: %w", err)
	}
	// TODO: Need to validate stack order of triggers and abilities
	// Emit trigger events
	// TODO: Make this a function or something so all the conditional logic
	// isn't in the action function.
	// TODO: Make this more generic
	{
		// Tap for Mana
		fmt.Println("Checking for tap for mana")
		_, ok := activatedAbility.Cost.(*TapCost)
		if ok {
			fmt.Println("is tap cost")
			if activatedAbility.IsManaAbility() {
				fmt.Println("is mana ability")
				state.EmitEvent(Event{
					Type:   EventTapForMana,
					Source: selectedPermanent,
				}, resolver)
				fmt.Println("Emitted event")
			}
		}
	}
	state.Log(fmt.Sprintf(
		"%s, paid %s to %s",
		selectedPermanent.Name(),
		activatedAbility.Cost.Description(),
		activatedAbility.Description()),
	)
	return &ActionResult{
		Message: fmt.Sprintf(
			"ability resolved: %s (%s)",
			selectedPermanent.Name(),
			activatedAbility.Description(),
		),
	}, nil
}

// ActionDiscardFunc handles the discard action. This is performed
// automatically at the end of turn by during the clean up state. It can also
// be performed manually by the player if Cheat is enabled in the game state.
// The target is the number of cards to discard.
func ActionDiscardFunc(state *GameState, target string, resolver ChoiceResolver) (*ActionResult, error) {
	n, err := strconv.Atoi(target)
	if err != nil {
		return nil, fmt.Errorf("failed to convert %s to int: %w", target, err)
	}
	state.Discard(n, string(ActionDiscard), resolver)
	return &ActionResult{
		Message: "card discarded",
	}, nil
}

// ActionDrawFunc handles the draw action. This is performed automatically at
// the beginning of turn during the draw step. It can also be performed
// manually by the player if Cheat is enabled in the game state. The target is
// the number of cards to draw.
func ActionDrawFunc(state *GameState, target string, resolver ChoiceResolver) (*ActionResult, error) {
	n, err := strconv.Atoi(target)
	if err != nil {
		return nil, fmt.Errorf("failed to convert %s to int: %w", target, err)
	}
	if err := state.Draw(n); err != nil {
		return nil, err
	}
	return &ActionResult{
		Message: "card drawn",
	}, nil
}

// ActionPlayFunc handles the play action. This is performed by the player to
// play a card from their hand. The target is the name of the card to play.
func ActionPlayFunc(state *GameState, target string, resolver ChoiceResolver) (result *ActionResult, err error) {
	var card *Card
	if target != "" {
		card = state.Hand.FindCard(target)
	}
	if card == nil {
		choices := state.Hand.CardChoices()
		choice, err := resolver.ChooseOne(
			"Which card to play from hand",
			string(ActionPlay),
			choices,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose card: %w", err)
		}
		card = state.Hand.GetCard(choice.Index)
	}
	if card != nil {
		if card.HasType(CardTypeLand) {
			return actionPlayLandFunc(state, resolver, card)
		}
		return actionCastSpellFunc(state, resolver, card)
	}
	return nil, nil
}

// ActionUntapFunc handles the untap action. This is performed automatically
// at the beginning of turn during the untap step. It can also be performed
// manually by the player if Cheat is enabled in the game state. The target
// is the name of the card to untap. If target == const(UntapAll), all
// permanents are untapped.
func ActionUntapFunc(state *GameState, target string, resolver PlayerAgent) (*ActionResult, error) {
	if target == UntapAll {
		state.Battlefield.UntapPermanents()
		return &ActionResult{Message: "all permanents untapped"}, nil
	}
	var selectedPermanent *Permanent
	if target != "" {
		selectedPermanent = state.Battlefield.FindTappedPermanent(target)
	}
	if selectedPermanent == nil {
		choices := state.Battlefield.GetTappedPermanents(state)
		choice, err := resolver.ChooseOne(
			"Which permanent to activate",
			string(ActionActivate),
			choices,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose permanent: %w", err)
		}
		selectedPermanent = state.Battlefield.GetPermanent(choice.Index)
	}
	return &ActionResult{
		Message: fmt.Sprintf("%s untapped", selectedPermanent.Name()),
	}, nil
}

// ActionViewFunc handles the view action. This is performed by the player
// to view a card in their hand, battlefield, or graveyard. The target is
// the zone to view. The player can choose to view their hand, battlefield,
// or graveyard. They will then be prompted to view a specific card.
func ActionViewFunc(state *GameState, target string, resolver PlayerAgent) (*ActionResult, error) {
	var choice Choice
	// TODO: Don't like this
	var err error
	if target == "" {
		choices := []Choice{
			{
				Name: "Hand",
			}, {
				Name: "Battlefield",
			}, {
				Name: "Graveyard",
			},
		}
		choice, err = resolver.ChooseOne("Which zone", "view zone", choices)
		if err != nil {
			return nil, fmt.Errorf("failed to choose zone: %w", err)
		}
	} else {
		choice = Choice{Name: target}
	}
	if choice.Name == "Hand" {
		return viewHand(state, resolver)
	}
	if choice.Name == "Battlefield" {
		return viewBattlefield(state, resolver)
	}
	if choice.Name == "Graveyard" {
		return viewGraveyard(state, resolver)
	}
	return nil, fmt.Errorf("unknown zone or not yet implemented")
}

// TODO: There's probably an abstraction we can set up for viewHand, viewBattlefield, and viewGraveyard
// probably a general abstraction for cards/permanents and zones
func viewHand(state *GameState, resolver ChoiceResolver) (result *ActionResult, err error) {
	choices := state.Hand.CardChoices()
	choice, err := resolver.ChooseOne("Which card", "view hand", choices)
	if err != nil {
		return nil, fmt.Errorf("failed to choose card: %w", err)
	}
	card := state.Hand.GetCard(choice.Index)
	state.Log("viewed " + card.Name())
	result = &ActionResult{Message: fmt.Sprintf("CARD: %s :: %s :: %s", card.Name(), card.CardTypes(), card.RulesText())}
	return result, nil
}

func viewBattlefield(state *GameState, resolver ChoiceResolver) (result *ActionResult, err error) {
	var choices []Choice
	permanents := state.Battlefield.Permanents()
	for i, card := range permanents {
		choices = append(choices, Choice{Name: card.Name(), Index: i})
	}
	choice, err := resolver.ChooseOne("Which card", "view battlefield", choices)
	if err != nil {
		return nil, fmt.Errorf("failed to choose card: %w", err)
	}
	card := permanents[choice.Index]
	state.Log("viewed " + card.Name())
	result = &ActionResult{Message: fmt.Sprintf("CARD: %s :: %s", card.Name, card.RulesText)}
	return result, nil
}

func viewGraveyard(state *GameState, resolver ChoiceResolver) (result *ActionResult, err error) {
	var choices []Choice
	for i, card := range state.Graveyard {
		choices = append(choices, Choice{Name: card.Name(), Index: i})
	}
	choice, err := resolver.ChooseOne("Which card", "view graveyard", choices)
	if err != nil {
		return nil, fmt.Errorf("failed to choose card: %w", err)
	}
	card := state.Graveyard[choice.Index]
	state.Log("viewed " + card.Name())
	result = &ActionResult{Message: fmt.Sprintf("CARD: %s :: %s", card.Name, card.RulesText)}
	return result, nil
}

// TODO: Maybe this should be a method off of GameState
func actionPlayLandFunc(state *GameState, resolver ChoiceResolver, card *Card) (result *ActionResult, err error) {
	if card.HasType(CardTypeLand) {
		if state.LandDrop {
			return nil, errors.New("land already played this turn")
		}
		state.LandDrop = true
	}
	state.Log("Played land: " + card.Name())
	state.Hand.RemoveCard(card)
	permanent, err := NewPermanent(card)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create permanent from %s: %w",
			card.Name(),
			err,
		)
	}
	state.Battlefield.AddPermanent(permanent)
	return &ActionResult{
		Message: "played land: " + card.Name(),
	}, nil
}

// TODO: Maybe this should be a method off of GameState
// or maybe a method off of Card, e.g. card.Cast() like Ability.Resolve()
func actionCastSpellFunc(state *GameState, resolver ChoiceResolver, card *Card) (result *ActionResult, err error) {
	if err := card.ManaCost().Pay(state, resolver); err != nil {
		return nil, err
	}
	state.Log("Casing spell: " + card.Name())
	state.Hand.RemoveCard(card)
	if card.IsSpell() {
		state.Log("Card: " + card.Name() + " is a spell")
		spell, err := NewSpell(card)
		if err != nil {
			return nil, fmt.Errorf("failed to create spell from %s: %w", card.Name(), err)
		}
		if err := spell.SpellAbility().Resolve(state, resolver); err != nil {
			return nil, fmt.Errorf("failed to resolve spell: %w", err)
		}
	} else if card.IsPermanent() {
		state.Log("Card: " + card.Name() + " is permanent")
		permanent, err := NewPermanent(card)
		if err != nil {
			return nil, fmt.Errorf("failed to create permanent from %s: %w", card.Name(), err)
		}
		state.Battlefield.AddPermanent(permanent)
	}
	return &ActionResult{
		Message: "played card: " + card.Name(),
	}, nil
}

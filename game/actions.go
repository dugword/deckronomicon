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

// TODO: Find some way to enforce this, maybe an action map or lookup function
// type ActionFunc func(*GameState, string, ChoiceResolver) (*ActionResult, error)

// TODO: This is a bit of a hack. We should probably have a better way to
// to create the object for ChoiceSource.
// This is for the ChoiceSorce interface
func (a GameActionType) Name() string {
	return string(a)
}

// This is for the ChoiceSorce interface
func (a GameActionType) ID() string {
	return string(a)
}

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
	/*
		var selectedObject GameObject
		if target != "" {
			selectedObject = state.Battlefield.FindPermanentWithAvailableActivatedAbility(target, state)
		}
	*/
	var abilities []*ActivatedAbility
	for _, zone := range state.Zones() {
		abilities = append(abilities, zone.AvailableActivatedAbilities(state)...)
	}
	if len(abilities) == 0 {
		// TODO: Do I need to check this?
		return nil, fmt.Errorf("no activated abilities available")
	}
	choices := CreateActivatedAbilityChoices(abilities)
	choice, err := resolver.ChooseOne(
		"Which ability to activate",
		ActionActivate,
		AddOptionalChoice(choices),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to choose ability: %w", err)
	}
	if choice.ID == ChoiceNone {
		return nil, nil
	}
	var ability *ActivatedAbility
	for _, a := range abilities {
		if a.ID == choice.ID {
			ability = a
			break
		}
	}
	if ability == nil {
		return nil, fmt.Errorf("failed to find activated ability: %w", err)
	}
	if err := ability.Cost.Pay(state, resolver); err != nil {
		return nil, fmt.Errorf("cannot pay activated ability cost: %w", err)
	}
	if err := ability.Resolve(state, resolver); err != nil {
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
		_, ok := ability.Cost.(*TapCost)
		if ok {
			if ability.IsManaAbility() {
				state.EmitEvent(Event{
					Type:   EventTapForMana,
					Source: ability.source,
				}, resolver)
			}
		}
	}
	state.Log(fmt.Sprintf(
		"%s, paid %s to %s",
		ability.source.Name(),
		ability.Cost.Description(),
		ability.Description(),
	),
	)
	return &ActionResult{
		Message: fmt.Sprintf(
			"ability resolved: %s (%s)",
			ability.source.Name(),
			ability.Description(),
		),
	}, nil
	return nil, nil
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
	// TODO:
	state.Discard(n, ActionDiscard, resolver)
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
	/*
		var card GameObject
			if target != "" {
				card, err = state.Hand.FindByName(target)
				if err != nil {
					return nil, fmt.Errorf("failed to find card in hand: %w", err)
				}
			}
	*/
	var choices []Choice
	for _, zone := range state.Zones() {
		cards := zone.AvailableToPlay(state)
		cs := CreateObjectChoices(cards, zone.ZoneType())
		choices = append(choices, cs...)
	}
	if len(choices) == 0 {
		return nil, fmt.Errorf("no cards available to play")
	}
	choice, err := resolver.ChooseOne(
		"Which card to play",
		ActionPlay,
		AddOptionalChoice(choices),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to choose card: %w", err)
	}
	if choice.ID == ChoiceNone {
		return nil, nil
	}
	zone, err := state.GetZone(choice.Zone)
	if err != nil {
		return nil, fmt.Errorf("failed to get zone: %w", err)
	}
	card, err := zone.Get(choice.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to take card from zone: %w", err)
	}
	c, ok := card.(*Card)
	if !ok {
		return nil, fmt.Errorf("object is not a card: %w", err)
	}
	if c.HasCardType(CardTypeLand) {
		return actionPlayLandFunc(state, resolver, c)
	}
	return actionCastSpellFunc(state, resolver, c)
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
	var selectedObject GameObject
	var err error
	if target != "" {
		selectedObject, err = state.Battlefield.FindTappedPermanent(target)
		if err != nil {
			return nil, fmt.Errorf("failed to find permanent: %w", err)
		}
	}
	if selectedObject == nil {
		objects := state.Battlefield.GetTappedPermanents()
		choices := CreateObjectChoices(objects, ZoneBattlefield)
		choice, err := resolver.ChooseOne(
			"Which permanent to untap",
			ActionUntap,
			AddOptionalChoice(choices),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose permanent: %w", err)
		}
		if choice.ID == ChoiceNone {
			return nil, nil
		}
		selectedObject, err = state.Battlefield.Get(choice.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get permanent from battlefield: %w", err)
		}
	}
	selectedPermanent, ok := selectedObject.(*Permanent)
	if !ok {
		return nil, fmt.Errorf("object is not a permanent: %w", err)
	}
	selectedPermanent.Untap()
	return &ActionResult{
		Message: fmt.Sprintf("%s untapped", selectedObject.Name()),
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
				Name: ZoneHand,
			}, {
				Name: ZoneBattlefield,
			}, {
				Name: ZoneGraveyard,
			},
		}
		choice, err = resolver.ChooseOne(
			"Which zone",
			NewChoiceSource("View Zone", "View Zone"),
			AddOptionalChoice(choices),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose zone: %w", err)
		}
		if choice.ID == ChoiceNone {
			return nil, nil
		}
	} else {
		choice = Choice{Name: target}
	}
	if choice.Name == ZoneHand {
		return viewHand(state, resolver)
	}
	if choice.Name == ZoneBattlefield {
		return viewBattlefield(state, resolver)
	}
	if choice.Name == ZoneGraveyard {
		return viewGraveyard(state, resolver)
	}
	return nil, fmt.Errorf("unknown zone or not yet implemented")
}

// TODO: There's probably an abstraction we can set up for viewHand, viewBattlefield, and viewGraveyard
// probably a general abstraction for cards/permanents and zones
func viewHand(state *GameState, resolver ChoiceResolver) (result *ActionResult, err error) {
	/*
		choices := state.Hand.CardChoices()
		choice, err := resolver.ChooseOne(
			"Which card",
			NewChoiceSource("View Hand", "View Hand"),
			AddOptionalChoice(choices),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose card: %w", err)
		}
		if choice.ID == ChoiceNone {
			return nil, nil
		}
		card, err := state.Hand.GetCard(choice.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get card from hand: %w", err)
		}
		state.Log("viewed " + card.Name())
		result = &ActionResult{Message: fmt.Sprintf("CARD: %s :: %s :: %s", card.Name(), card.CardTypes(), card.RulesText())}
		return result, nil
	*/
	return nil, nil
}

func viewBattlefield(state *GameState, resolver ChoiceResolver) (result *ActionResult, err error) {
	/*
		var choices []Choice
		permanents := state.Battlefield.Permanents()
		for _, permanent := range permanents {
			choices = append(choices, Choice{Name: permanent.Name(), ID: permanent.ID()})
		}
		choice, err := resolver.ChooseOne(
			"Which card",
			NewChoiceSource("View Battlefield", "View Battlefield"),
			AddOptionalChoice(choices),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose card: %w", err)
		}
		if choice.ID == ChoiceNone {
			return nil, nil
		}
		card, err := state.Battlefield.GetPermanent(choice.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get permanent from battlefield: %w", err)
		}
		state.Log("viewed " + card.Name())
		result = &ActionResult{Message: fmt.Sprintf("CARD: %s :: %s", card.Name, card.RulesText)}
		return result, nil
	*/
	return nil, nil
}

func viewGraveyard(state *GameState, resolver ChoiceResolver) (result *ActionResult, err error) {
	var choices []Choice
	// TODO Remove the .cards access
	for _, card := range state.Graveyard.cards {
		choices = append(choices, Choice{Name: card.Name(), ID: card.ID()})
	}
	choice, err := resolver.ChooseOne(
		"Which card",
		NewChoiceSource("View Graveyard", "View Graveyard"),
		AddOptionalChoice(choices),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to choose card: %w", err)
	}
	if choice.ID == ChoiceNone {
		return nil, nil
	}
	var selectedCard *Card
	// TODO remove the .cards access
	for _, card := range state.Graveyard.cards {
		if card.ID() == choice.ID {
			selectedCard = card
			break
		}
	}
	if selectedCard == nil {
		return nil, fmt.Errorf("failed to get card from graveyard: %w", err)
	}
	state.Log("viewed " + selectedCard.Name())
	result = &ActionResult{Message: fmt.Sprintf("CARD: %s :: %s", selectedCard.Name(), selectedCard.RulesText())}
	return result, nil
}

// TODO: Maybe this should be a method off of GameState
func actionPlayLandFunc(state *GameState, resolver ChoiceResolver, card *Card) (result *ActionResult, err error) {
	if card.HasCardType(CardTypeLand) {
		if state.LandDrop {
			return nil, errors.New("land already played this turn")
		}
		state.LandDrop = true
	}
	state.Log("Played land: " + card.Name())
	state.Hand.Remove(card.ID())
	permanent, err := NewPermanent(card)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create permanent from %s: %w",
			card.Name(),
			err,
		)
	}
	state.Battlefield.Add(permanent)
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
	state.Hand.Remove(card.ID())
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
		state.Battlefield.Add(permanent)
	}
	return &ActionResult{
		Message: "played card: " + card.Name(),
	}, nil
}

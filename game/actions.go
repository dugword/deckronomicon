package game

import (
	"errors"
	"fmt"
)

// GameActionType represents the type of player action
type GameActionType string

// GameCheatType represents the type of player cheat
type GameCheatType string

const (
	ActionActivate    GameActionType = "Activate"
	ActionBattlefield GameActionType = "Battlefield"
	ActionCheat       GameActionType = "Cheat"
	ActionConcede     GameActionType = "Concede"
	ActionGraveyard   GameActionType = "Graveyard"
	ActionPass        GameActionType = "Pass"
	ActionPlay        GameActionType = "Play"
	ActionView        GameActionType = "View"
)

const (
	CheatDiscard   GameCheatType = "CheatDiscard"
	CheatDraw      GameCheatType = "CheatDraw"
	CheatFind      GameCheatType = "CheatFind"
	CheatPeek      GameCheatType = "CheatPeek"
	CheatPrintDeck GameCheatType = "CheatPrintDeck"
	CheatShuffle   GameCheatType = "CheatShuffle"
)

// GameAction represents an action a player can take
type GameAction struct {
	Type   GameActionType
	Cheat  GameCheatType
	Target string
}

func cardToPermanent(card *Card) *Permanent {
	return &Permanent{Object: card.Object}
}

func ActionDiscardFunc(state *GameState, resolver ChoiceResolver) (result *ActionResult, err error) {
	options := getCardOptions(state.Hand)
	choice := resolver.ChooseOne("Which card to discard from hand", options)
	card := state.Hand[choice.Index]
	if card != nil {
		card := state.Hand[choice.Index]
		if card != nil {
			state.Hand = removeCard(state.Hand, card)
			state.Graveyard = append(state.Graveyard, card)
		}
	}
	// TODO: Maybe every action resultion should log and put a status message
	return nil, nil
}

// TODO: Need to check bounds on arrays and stuff, probably seems like there's a more elegent solution there too
// TODO: Only supports 1 activated ability and only tap costs
func ActionActivateFunc(state *GameState, target string, resolver ChoiceResolver) (result *ActionResult, err error) {
	var selectedPermanent *Permanent
	if target != "" {
		selectedPermanent = findUntappedPermanent(state.Battlefield, target)
	}
	if selectedPermanent == nil {
		options := FindPermanentsWithActivatableAbilities(state.Battlefield)
		choice := resolver.ChooseOne("Which permanent to activate", options)
		selectedPermanent = state.Battlefield[choice.Index]
	}
	if selectedPermanent.Tapped {
		return nil, fmt.Errorf("permanent already tapped")
	}
	selectedPermanent.Tapped = true
	if len(selectedPermanent.ActivatedAbilities) > 1 {
		return nil, fmt.Errorf("no support for multiple activated abilities")
	}
	// TODO: maybe check to see if the ability exists so this doesn't explode?
	activatedAbility := selectedPermanent.ActivatedAbilities[0]
	activatedAbility.Effect(state, resolver)
	state.Log(fmt.Sprintf("tapping %s for %s", selectedPermanent.Name, activatedAbility.Description))
	return nil, nil
}

func ActionPlayFunc(state *GameState, target string, resolver ChoiceResolver) (result *ActionResult, err error) {
	var card *Card
	if target != "" {
		card = findCard(state.Hand, target)
	}
	if card == nil {
		options := getCardOptions(state.Hand)
		choice := resolver.ChooseOne("Which card to play from hand", options)
		card = state.Hand[choice.Index]
	}
	if card != nil {
		if card.HasType(CardTypeLand) {
			return actionPlayLandFunc(state, resolver, card)
		}
		return actionCastSpellFunc(state, resolver, card)
	}
	return nil, nil
}

func ActionViewFunc(state *GameState, resolver PlayerAgent) (result *ActionResult, err error) {
	options := []Choice{
		{
			Name: "Hand",
		}, {
			Name: "Battlefield",
		}, {
			Name: "Graveyard",
		},
	}
	choice := resolver.ChooseOne("Which zone", options)
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
	options := getCardOptions(state.Hand)
	choice := resolver.ChooseOne("Which card", options)
	card := state.Hand[choice.Index]
	state.Log("viewed " + card.Name)
	result = &ActionResult{Message: fmt.Sprintf("CARD: %s :: %s", card.Name, card.RulesText)}
	return result, nil
}

func viewBattlefield(state *GameState, resolver ChoiceResolver) (result *ActionResult, err error) {
	var options []Choice
	for i, card := range state.Battlefield {
		options = append(options, Choice{Name: card.Name, Index: i})
	}
	choice := resolver.ChooseOne("Which card", options)
	card := state.Battlefield[choice.Index]
	state.Log("viewed " + card.Name)
	result = &ActionResult{Message: fmt.Sprintf("CARD: %s :: %s", card.Name, card.RulesText)}
	return result, nil
}

func viewGraveyard(state *GameState, resolver ChoiceResolver) (result *ActionResult, err error) {
	var options []Choice
	for i, card := range state.Graveyard {
		options = append(options, Choice{Name: card.Name, Index: i})
	}
	choice := resolver.ChooseOne("Which card", options)
	card := state.Graveyard[choice.Index]
	state.Log("viewed " + card.Name)
	result = &ActionResult{Message: fmt.Sprintf("CARD: %s :: %s", card.Name, card.RulesText)}
	return result, nil
}

func actionPlayLandFunc(state *GameState, resolver ChoiceResolver, card *Card) (result *ActionResult, err error) {
	if card.HasType(CardTypeLand) {
		if state.LandDrop {
			return nil, errors.New("land already played this turn")
		}
		state.LandDrop = true
	}
	state.Log("Played card: " + card.Name)
	state.Hand = removeCard(state.Hand, card)
	state.Battlefield = append(state.Battlefield, cardToPermanent(card))
	return nil, nil
}

func actionCastSpellFunc(state *GameState, resolver ChoiceResolver, card *Card) (result *ActionResult, err error) {
	if err := card.ManaCost.Pay(state, resolver, &card.Object); err != nil {
		return nil, err
	}
	state.Log("Played card: " + card.Name)
	state.Hand = removeCard(state.Hand, card)
	if card.SpellAbility != nil {
		card.SpellAbility.Effect(state, resolver)
	}
	if card.IsPermanent() {
		state.Log("Card: " + card.Name + " is permanent")
		state.Battlefield = append(state.Battlefield, cardToPermanent(card))
	}

	return nil, nil
}

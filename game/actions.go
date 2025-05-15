package game

import (
	"errors"
	"fmt"
	"strings"
)

// GameActionType represents the type of player action
type GameActionType string

const (
	ActionActivate    GameActionType = "Activate"
	ActionDraw        GameActionType = "Draw"
	ActionBattlefield GameActionType = "Battlefield"
	ActionCheat       GameActionType = "Cheat"
	ActionConcede     GameActionType = "Concede"
	ActionGraveyard   GameActionType = "Graveyard"
	ActionPass        GameActionType = "Pass"
	ActionPlay        GameActionType = "Play"
	ActionUntap       GameActionType = "Untap"
	ActionView        GameActionType = "View"
)

type ActionResult struct {
	Message string
	Pass    bool
}

func ActionDrawFunc(state *GameState, resolver ChoiceResolver) (*ActionResult, error) {
	drawn, err := state.Deck.DrawCards(1)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, card := range drawn {
		names = append(names, card.Name())
	}
	state.Hand.Add(drawn...)
	return &ActionResult{
		Message: fmt.Sprintf("drew: %s", strings.Join(names, ", ")),
	}, nil
}

// TODO: Need to check bounds on arrays and stuff, probably seems like there's a more elegent solution there too
// TODO: Only supports 1 activated ability and only tap costs
func ActionActivateFunc(state *GameState, target string, resolver ChoiceResolver) (*ActionResult, error) {
	var selectedPermanent *Permanent
	if target != "" {
		selectedPermanent = state.Battlefield.FindPermanentWithAvailableActivatedAbility(target, state)
	}
	if selectedPermanent == nil {
		choices := state.Battlefield.GetPermanentsWithActivatedAbilities(state)
		choice := resolver.ChooseOne("Which permanent to activate", choices)
		selectedPermanent = state.Battlefield.GetPermanent(choice.Index)
	}
	// TODO: Support more than 1
	if len(selectedPermanent.ActivatedAbilities()) > 1 {
		return nil, fmt.Errorf("no support for multiple activated abilities")
	}
	activatedAbility := selectedPermanent.ActivatedAbilities()[0]
	if err := activatedAbility.Cost.Pay(state, selectedPermanent); err != nil {
		// TODO wrap this
		return nil, err
	}
	activatedAbility.Effect(state, resolver)
	state.Log(fmt.Sprintf("tapping %s for %s", selectedPermanent.Name(), activatedAbility.Description))
	return nil, nil
}

func ActionPlayFunc(state *GameState, target string, resolver ChoiceResolver) (result *ActionResult, err error) {
	var card *Card
	if target != "" {
		card = state.Hand.FindCard(target)
	}
	if card == nil {
		choices := state.Hand.CardChoices()
		choice := resolver.ChooseOne("Which card to play from hand", choices)
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

func ActionUntapFunc(state *GameState, resolver PlayerAgent) (*ActionResult, error) {
	state.Battlefield.UntapPermanents()
	return &ActionResult{Message: "permanents untapped"}, nil
}

func ActionViewFunc(state *GameState, resolver PlayerAgent) (*ActionResult, error) {
	_ = []Choice{
		{
			Name: "Hand",
		}, {
			Name: "Battlefield",
		}, {
			Name: "Graveyard",
		},
	}
	/* // TODO: might not need this
	choice := resolver.ChooseOne("Which zone", choices)
	if choice.Name == "Hand" {
		return viewHand(state, resolver)
	}
	if choice.Name == "Battlefield" {
		return viewBattlefield(state, resolver)
	}
	if choice.Name == "Graveyard" {
		return viewGraveyard(state, resolver)
	}
	*/
	return nil, fmt.Errorf("unknown zone or not yet implemented")
}

/*
// TODO: There's probably an abstraction we can set up for viewHand, viewBattlefield, and viewGraveyard
// probably a general abstraction for cards/permanents and zones
func viewHand(state *GameState, resolver ChoiceResolver) (result *ActionResult, err error) {
	choices := state.Hand.CardChoices()
	choice := resolver.ChooseOne("Which card", choices)
	card := state.Hand.GetCard(choice.Index)
	state.Log("viewed " + card.Name())
	result = &ActionResult{Message: fmt.Sprintf("CARD: %s :: %s", card.Name, card.RulesText)}
	return result, nil
}

func viewBattlefield(state *GameState, resolver ChoiceResolver) (result *ActionResult, err error) {
	var choices []Choice
	for i, card := range state.Battlefield {
		choices = append(choices, Choice{Name: card.Name(), Index: i})
	}
	choice := resolver.ChooseOne("Which card", choices)
	card := state.Battlefield[choice.Index]
	state.Log("viewed " + card.Name())
	result = &ActionResult{Message: fmt.Sprintf("CARD: %s :: %s", card.Name, card.RulesText)}
	return result, nil
}

func viewGraveyard(state *GameState, resolver ChoiceResolver) (result *ActionResult, err error) {
	var choices []Choice
	for i, card := range state.Graveyard {
		choices = append(choices, Choice{Name: card.Name(), Index: i})
	}
	choice := resolver.ChooseOne("Which card", choices)
	card := state.Graveyard[choice.Index]
	state.Log("viewed " + card.Name())
	result = &ActionResult{Message: fmt.Sprintf("CARD: %s :: %s", card.Name, card.RulesText)}
	return result, nil
}
*/

func actionPlayLandFunc(state *GameState, resolver ChoiceResolver, card *Card) (result *ActionResult, err error) {
	if card.HasType(CardTypeLand) {
		if state.LandDrop {
			return nil, errors.New("land already played this turn")
		}
		state.LandDrop = true
	}
	state.Log("Played card: " + card.Name())
	state.Hand.RemoveCard(card)
	state.Battlefield.AddPermanent(NewPermanent(card))
	return nil, nil
}

func actionCastSpellFunc(state *GameState, resolver ChoiceResolver, card *Card) (result *ActionResult, err error) {
	if err := card.ManaCost().Pay(state, resolver, card); err != nil {
		return nil, err
	}
	state.Log("Played card: " + card.Name())
	state.Hand.RemoveCard(card)
	if card.SpellAbility() != nil {
		card.SpellAbility().Effect(state, resolver)
	}
	if card.IsPermanent() {
		state.Log("Card: " + card.Name() + " is permanent")
		state.Battlefield.AddPermanent(NewPermanent(card))
	}

	return nil, nil
}

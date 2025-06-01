package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/game/action"
	"deckronomicon/packages/game/card"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/permanent"
	"deckronomicon/packages/game/player"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"errors"
	"fmt"
	"strconv"
)

// TODO: This could theoretically conflict with a card name.
// Perhaps a standard special character like "!" or "@" could be used.
const UntapAll = "UntapAll"

// ActionResult represents the result of an action.
type ActionResult struct {
	Message string
	Pass    bool
}

// TODO: ActionAction sucks, maybe move it to the game package and have
// game.Action

// ResolveAction handles the resolution of game actions.
// TODO Maybe a map, that would enforce the action func signature
func (e *Engine) ResolveAction(act action.Action, player *player.Player) (result ActionResult, err error) {
	switch act.Type {
	case action.ActionActivate:
		e.Log("Action: activate")
		return ActionActivateFunc(e.GameState, player, act.Target)
	case action.ActionDraw:
		e.Log("Action: draw")
		return ActionDrawFunc(e.GameState, player, act.Target)
	case action.ActionCheat:
		e.Log("Action: cheat... you cheater")
		e.GameState.CheatsEnabled = true
		return ActionResult{Message: "Cheat mode enabled"}, nil
	case action.ActionConcede:
		e.Log("Action: concede")
		// TODO: Should this be an error?
		return ActionResult{}, mtg.PlayerLostError{Reason: mtg.Conceded}
	case action.ActionPass:
		e.Log("Action: pass")
		return ActionResult{
			Message: "Player passed",
			Pass:    true,
		}, nil
	case action.ActionPlay:
		// todo: make this nicer
		e.Log(fmt.Sprintf("Action: %+v", act))
		return ActionPlayFunc(e.GameState, player, act.Target)
	case action.ActionUntap:
		e.Log("Action: untap")
		return ActionUntapFunc(e.GameState, player, act.Target)
	case action.ActionView:
		e.Log("Action: view")
		return ActionViewFunc(e.GameState, player, act.Target)
	case action.CheatAddMana:
		e.Log("CHEAT! Action: add mana")
		return ActionAddManaFunc(e.GameState, player, act.Target)
	case action.CheatConjure:
		e.Log("Action: conjure")
		return ActionConjureFunc(e.GameState, player, act.Target)
	case action.CheatDraw:
		e.Log("CHEAT! Action: draw")
		return ActionDrawFunc(e.GameState, player, act.Target)
	case action.CheatFind:
		e.Log("CHEAT! Action: find")
		return ActionFindFunc(e.GameState, player, act.Target)
	case action.CheatLandDrop:
		e.Log("CHEAT! Action: land drop")
		return ActionLandDropFunc(e.GameState, player, act.Target)
	case action.CheatPeek:
		e.Log("CHEAT! Action: peek")
		/*
			return ActionResult{
				// TODO: No .cards access
				Message: "Top Card: " + player.Library.Peek().Name(),
			}, nil
		*/
		// TODO: fix this
		return ActionResult{Message: "Top Card: <hidden>"}, nil
	case action.CheatShuffle:
		e.Log("CHEAT! Action: shuffle")
		player.ShuffleLibrary()
		return ActionResult{Message: "Library shuffled"}, nil
	case action.CheatDiscard:
		e.Log("CHEAT! Action: discard")
		// TODO: Name is a bad field for this, but it works for now
		return ActionDiscardFunc(e.GameState, player, action.ActionTarget{Name: "1"})
	default:
		e.Log("Unknown Action: " + string(act.Type))
		// TODO: Should this be an error? Need to think through error handling
		return ActionResult{Message: "Unknown Action"}, nil
	}
}

// ActionAddManaFunc handles the add mana action. This is performed by the
// player to add mana to their mana pool. The target is the amount of mana
// to add. This is a cheat.
func ActionAddManaFunc(state *GameState, player *player.Player, target action.ActionTarget) (ActionResult, error) {
	// TODO: This is a hack, we should probably have a better way to
	mana := target.Name
	if mana == "" {
		choices := []choose.Choice{
			{
				Name: "White",
				ID:   "{W}",
				// Source: ChoiceCheat,
			}, {
				Name: "Blue",
				ID:   "{U}",
				// Source: ChoiceCheat,
			}, {
				Name: "Black",
				ID:   "{B}",
				// Source: ChoiceCheat,
			}, {
				Name: "Red",
				ID:   "{R}",
				// Source: ChoiceCheat,
			}, {
				Name: "Green",
				ID:   "{G}",
				// Source: ChoiceCheat,
			}, {
				Name: "Colorless",
				ID:   "{C}",
				// Source: ChoiceCheat,
			},
		}
		choice, err := player.Agent.ChooseOne(
			"Add mana to mana pool",
			// TODO: This should be a constant
			choose.NewChoiceSource("Cheat"),
			choices,
		)
		if err != nil {
			return ActionResult{}, fmt.Errorf("failed to choose mana: %w", err)
		}
		mana = choice.ID
	}
	player.AddMana(mana)
	return ActionResult{
		Message: fmt.Sprintf("%s mana added to pool", mana),
	}, nil
	return ActionResult{}, nil
}

// ActionConjureFunc handles the conjure action. This is performed by the
// player to conjure a card. The target is the name of the card to conjure.
// This is a cheat.
func ActionConjureFunc(state *GameState, player *player.Player, target action.ActionTarget) (ActionResult, error) {
	cardName := target.Name
	if cardName == "" {
		return ActionResult{}, errors.New("no card name provided")
	}
	cardDefinition, ok := state.CardDefinitions[cardName]
	if !ok {
		return ActionResult{}, fmt.Errorf(
			"card %s not found in card definitions",
			target,
		)
	}
	card, err := card.NewCardFromCardDefinition(state, cardDefinition)
	if err != nil {
		return ActionResult{}, fmt.Errorf(
			"failed to create card %s: %w",
			target,
			err,
		)
	}
	player.CheatAddCard(card)
	return ActionResult{
		Message: "conjured card: " + card.Name(),
	}, nil
	return ActionResult{}, nil
}

// ActionDiscardFunc handles the discard action. This is performed
// automatically at the end of turn by during the clean up state. It can also
// be performed manually by the player if Cheat is enabled in the game state.
// The target is the number of cards to discard.
func ActionDiscardFunc(state *GameState, player *player.Player, target action.ActionTarget) (ActionResult, error) {
	_, err := strconv.Atoi(target.Name)
	if err != nil {
		return ActionResult{}, fmt.Errorf(
			"failed to convert %s to int: %w",
			target,
			err,
		)
	}
	// TODO: Fix this
	/*
		Discard(n, ActionDiscard, player)
		return ActionResult{
			Message: "card discarded",
		}, nil
	*/
	return ActionResult{}, nil
}

// ActionDrawFunc handles the draw action. This is performed automatically at
// the beginning of turn during the draw step. It can also be performed
// manually by the player if Cheat is enabled in the game state. The target is
// the number of cards to draw.
func ActionDrawFunc(state *GameState, player *player.Player, target action.ActionTarget) (ActionResult, error) {
	count := target.Name
	n := 1
	var err error
	if count != "" {
		n, err = strconv.Atoi(count)
		if err != nil {
			return ActionResult{}, fmt.Errorf("failed to convert %s to int: %w", count, err)
		}
	}
	for range n {
		if _, err := player.DrawCard(); err != nil {
			return ActionResult{}, fmt.Errorf("failed to draw card: %w", err)
		}
	}
	return ActionResult{
		Message: "card drawn",
	}, nil
	return ActionResult{}, nil
}

// ActionFindFunc handles the find action. This is performed by the player to
// find a card in their library. The target is the name of the card to find.
// This is a cheat.
func ActionFindFunc(state *GameState, player *player.Player, target action.ActionTarget) (ActionResult, error) {
	var card query.Object
	cardName := target.Name
	if cardName != "" {
		if err := player.Tutor(has.Name(target.Name)); err != nil {
			return ActionResult{}, fmt.Errorf("failed to find card in library: %w", err)
		}
	} else {
		choices := choose.CreateChoices(
			player.Library().GetAll(),
			choose.NewChoiceSource(string(mtg.ZoneLibrary)),
		)
		choice, err := player.Agent.ChooseOne(
			"Which card to find",
			choose.NewChoiceSource("Cheat"),
			choose.AddOptionalChoice(choices),
		)
		if err != nil {
			return ActionResult{}, fmt.Errorf("failed to choose card: %w", err)
		}
		if err = player.Tutor(has.ID(choice.ID)); err != nil {
			return ActionResult{}, fmt.Errorf(
				"failed to tutor card: %w",
				err,
			)
		}
	}
	return ActionResult{
		Message: "found card: " + card.Name(),
	}, nil
	return ActionResult{}, nil
}

// ActionLandDropFunc handles the land drop action. This is performed by the
// player. Is resets the land drop flag for the player.
// This is a cheat.
func ActionLandDropFunc(state *GameState, player *player.Player, target action.ActionTarget) (ActionResult, error) {
	player.LandDrop = false
	return ActionResult{
		Message: "land drop reset",
	}, nil
}

// ActionUntapFunc handles the untap action. This is performed automatically
// at the beginning of turn during the untap step. It can also be performed
// manually by the player if Cheat is enabled in the game state. The target
// is the name of the card to untap. If target == const(UntapAll), all
// permanents are untapped.
func ActionUntapFunc(state *GameState, player *player.Player, target action.ActionTarget) (ActionResult, error) {
	permanentName := target.Name
	if permanentName == UntapAll {
		for _, obj := range state.Battlefield().GetAll() {
			perm, ok := obj.(*permanent.Permanent)
			if !ok {
				return ActionResult{}, ErrObjectNotPermanent
			}
			perm.Untap()
		}
		return ActionResult{Message: "all permanents untapped"}, nil
	}
	var err error
	var selectedObject query.Object
	if permanentName != "" {
		found, ok := state.Battlefield().Find(
			query.And(has.Name(permanentName), is.Tapped()),
		)
		if !ok {
			return ActionResult{}, fmt.Errorf("failed to find permanent: %w", err)
		}
		perm, ok := found.(*permanent.Permanent)
		if !ok {
			return ActionResult{}, fmt.Errorf("object is not a permanent: %w", err)
		}
		perm.Untap()
		return ActionResult{
			Message: fmt.Sprintf("%s untapped", perm.Name()),
		}, nil
	}
	if selectedObject == nil {
		objects := state.Battlefield().FindAll(is.Tapped())
		choices := choose.CreateChoices(
			objects,
			choose.NewChoiceSource(string(mtg.ZoneBattlefield)),
		)
		choice, err := player.Agent.ChooseOne(
			"Which permanent to untap",
			action.ActionUntap,
			choose.AddOptionalChoice(choices),
		)
		if err != nil {
			return ActionResult{}, fmt.Errorf("failed to choose permanent: %w", err)
		}
		if choice == choose.ChoiceNone {
			// TODO: Make this a constant
			return ActionResult{Message: "no choice made"}, nil
		}
		var ok bool
		selectedObject, ok = state.Battlefield().Get(choice.ID)
		if !ok {
			return ActionResult{}, fmt.Errorf("failed to get permanent from battlefield: %w", err)
		}
	}
	selectedPermanent, ok := selectedObject.(*permanent.Permanent)
	if !ok {
		return ActionResult{}, fmt.Errorf("object is not a permanent: %w", err)
	}
	selectedPermanent.Untap()
	return ActionResult{
		Message: fmt.Sprintf("%s untapped", selectedObject.Name()),
	}, nil
	return ActionResult{}, nil
}

// ActionViewFunc handles the view action. This is performed by the player
// to view a card in their hand, battlefield, or graveyard. The target is
// the zone to view. The player can choose to view their hand, battlefield,
// or graveyard. They will then be prompted to view a specific card.
func ActionViewFunc(state *GameState, player *player.Player, target action.ActionTarget) (ActionResult, error) {
	/*
		var choice Choice
		// TODO: Don't like this
		var err error
		if target.Name == "" {
			choices := []Choice{
				{
					Name: ZoneHand,
				}, {
					Name: ZoneBattlefield,
				}, {
					Name: ZoneGraveyard,
				},
			}
			choice, err = player.Agent.ChooseOne(
				"Which zone",
				NewChoiceSource("View Zone", "View Zone"),
				AddOptionalChoice(choices),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to choose zone: %w", err)
			}
			if choice.ID == ChoiceNone {
				return ActionResult{Message: "No choice made"}, nil
			}
		} else {
			choice = Choice{Name: target.Name}
		}
		if choice.Name == ZoneHand {
			return viewHand(state, player)
		}
		if choice.Name == ZoneBattlefield {
			return viewBattlefield(state, player)
		}
		if choice.Name == ZoneGraveyard {
			return viewGraveyard(state, player)
		}
		return nil, errors.New("unknown zone or not yet implemented")
	*/
	return ActionResult{}, nil
}

// TODO: There's probably an abstraction we can set up for viewHand, viewBattlefield, and viewGraveyard
// probably a general abstraction for cards/permanents and zones
func viewHand(state *GameState, player *player.Player) (result ActionResult, err error) {
	/*
		choices := state.Hand.CardChoices()
		choice, err := player.PlayerAgent.ChooseOne(
			"Which card",
			NewChoiceSource("View Hand", "View Hand"),
			AddOptionalChoice(choices),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose card: %w", err)
		}
		if choice.ID == ChoiceNone {
			return ActionResult{Message: "No choice made"}, nil
		}
		card, err := state.Hand.GetCard(choice.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get card from hand: %w", err)
		}
		state.Log("viewed " + card.Name())
		result = ActionResult{Message: fmt.Sprintf("CARD: %s :: %s :: %s", card.Name(), card.CardTypes(), card.RulesText())}
		return result, nil
	*/
	return ActionResult{Message: "No choice made"}, nil
}

func viewBattlefield(state *GameState, player *player.Player) (result ActionResult, err error) {
	/*
		var choices []Choice
		permanents := state.Battlefield.Permanents()
		for _, permanent := range permanents {
			choices = append(choices, Choice{Name: permanent.Name(), ID: permanent.ID()})
		}
		choice, err := player.PlayerAgent.ChooseOne(
			"Which card",
			NewChoiceSource("View Battlefield", "View Battlefield"),
			AddOptionalChoice(choices),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose card: %w", err)
		}
		if choice.ID == ChoiceNone {
			return ActionResult{Message: "No choice made"}, nil
		}
		card, err := state.Battlefield.GetPermanent(choice.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get permanent from battlefield: %w", err)
		}
		state.Log("viewed " + card.Name())
		result = ActionResult{Message: fmt.Sprintf("CARD: %s :: %s", card.Name, card.RulesText)}
		return result, nil
	*/
	return ActionResult{Message: "No choice made"}, nil
}

func viewGraveyard(state *GameState, player *player.Player) (result ActionResult, err error) {
	/*
		var choices []Choice
		// TODO Remove the .cards access
		for _, card := range player.Graveyard.cards {
			choices = append(choices, Choice{Name: card.Name(), ID: card.ID()})
		}
		choice, err := player.Agent.ChooseOne(
			"Which card",
			NewChoiceSource("View Graveyard", "View Graveyard"),
			AddOptionalChoice(choices),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose card: %w", err)
		}
		if choice.ID == ChoiceNone {
			return ActionResult{Message: "No choice made"}, nil
		}
		var selectedCard *card.Card
		// TODO remove the .cards access
		for _, card := range player.Graveyard.cards {
			if card.ID() == choice.ID {
				selectedCard = card
				break
			}
		}
		if selectedCard == nil {
			return nil, fmt.Errorf("failed to get card from graveyard: %w", err)
		}
		state.Log("viewed " + selectedCard.Name())
		result = ActionResult{Message: fmt.Sprintf("CARD: %s :: %s", selectedCard.Name(), selectedCard.RulesText())}
		return result, nil
	*/
	return ActionResult{}, nil
}

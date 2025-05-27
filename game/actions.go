package game

import (
	"errors"
	"fmt"
	"strconv"
)

// TODO: This could theoretically conflict with a card name.
// Perhaps a standard special character like "!" or "@" could be used.
const UntapAll = "UntapAll"

var ChoiceCheat = "Cheat"
var ChoiceSourceCheat = NewChoiceSource(ChoiceCheat, ChoiceCheat)

// GameAction represents an action a player can take.
type GameAction struct {
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

	// Cheat actions
	CheatAddMana   GameActionType = "CheatAddMana"
	CheatConjure   GameActionType = "CheatConjure"
	CheatDiscard   GameActionType = "CheatDiscard"
	CheatDraw      GameActionType = "CheatDraw"
	CheatFind      GameActionType = "CheatFind"
	CheatLandDrop  GameActionType = "CheatLandDrop"
	CheatPeek      GameActionType = "CheatPeek"
	CheatPrintDeck GameActionType = "CheatPrintDeck"
	CheatShuffle   GameActionType = "CheatShuffle"
)

// TODO: Find some way to enforce this, maybe an action map or lookup function
// TODO: Function signature should also be func(*GameState, Player, string) (*ActionResult, error)
// type ActionFunc func(*GameState, Player, Target) (*ActionResult, error)

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

// ResolveAction handles the resolution of game actions.
func (g *GameState) ResolveAction(action *GameAction, player *Player) (result *ActionResult, err error) {
	switch action.Type {
	case ActionActivate:
		g.Log("Action: activate")
		return ActionActivateFunc(g, player, action.Target)
	case ActionDraw:
		g.Log("Action: draw")
		return ActionDrawFunc(g, player, action.Target)
	case ActionCheat:
		g.Log("Action: cheat... you cheater")
		g.Cheat = true
		return &ActionResult{Message: "Cheat mode enabled"}, nil
	case ActionConcede:
		g.Log("Action: concede")
		// TODO: Should this be an error?
		return nil, PlayerLostError{Reason: Conceded}
	case ActionPass:
		g.Log("Action: pass")
		return &ActionResult{Pass: true}, nil
	case ActionPlay:
		// todo: make this nicer
		g.Log("Action: play " + action.Target)
		return ActionPlayFunc(g, player, action.Target)
	case ActionUntap:
		g.Log("Action: untap")
		return ActionUntapFunc(g, player, action.Target)
	case ActionView:
		g.Log("Action: view")
		return ActionViewFunc(g, player, action.Target)
	case CheatAddMana:
		g.Log("CHEAT! Action: add mana")
		return ActionAddManaFunc(g, player, action.Target)
	case CheatConjure:
		g.Log("Action: conjure")
		return ActionConjureFunc(g, player, action.Target)
	case CheatDraw:
		g.Log("CHEAT! Action: draw")
		return ActionDrawFunc(g, player, action.Target)
	case CheatFind:
		g.Log("CHEAT! Action: find")
		return ActionFindFunc(g, player, action.Target)
	case CheatLandDrop:
		g.Log("CHEAT! Action: land drop")
		return ActionLandDropFunc(g, player, action.Target)
	case CheatPeek:
		g.Log("CHEAT! Action: peek")
		return &ActionResult{
			// TODO: No .cards access
			Message: "Top Card: " + player.Library.Peek().Name(),
		}, nil
	case CheatShuffle:
		g.Log("CHEAT! Action: shuffle")
		player.Library.Shuffle()
		return &ActionResult{Message: "Library shuffled"}, nil
	case CheatDiscard:
		g.Log("CHEAT! Action: discard")
		return ActionDiscardFunc(g, player, "1")
	default:
		g.Log("Unknown Action: " + string(action.Type))
		// TODO: Should this be an error? Need to think through error handling
		return &ActionResult{Message: "Unknown Action"}, nil
	}
}

// ActionActivateFunc handles the activate action. This is performed by the
// player to activate an ability of a permanent on the battlefield, or a card
// in hand or in the graveyard. The target is the name of the permanent or
// card.
// TODO: Support more than one activated ability
// TODO: Support activated abilities in hand and graveyard
func ActionActivateFunc(state *GameState, player *Player, target string) (*ActionResult, error) {
	var abilities []*ActivatedAbility
	for _, zone := range player.Zones() {
		abilities = append(abilities, zone.AvailableActivatedAbilities(state, player)...)
	}
	if len(abilities) == 0 {
		// TODO: Do I need to check this?
		return nil, errors.New("no activated abilities available")
	}
	choices := CreateActivatedAbilityChoices(abilities)
	choice, err := player.Agent.ChooseOne(
		"Which ability to activate",
		ActionActivate,
		AddOptionalChoice(choices),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to choose ability: %w", err)
	}
	if choice.ID == ChoiceNone {
		return &ActionResult{Message: "No choice made"}, nil
	}
	var ability *ActivatedAbility
	for _, a := range abilities {
		if a.ID() == choice.ID {
			ability = a
			break
		}
	}
	if ability == nil {
		return nil, fmt.Errorf("failed to find activated ability: %w", err)
	}
	if err := ability.Cost.Pay(state, player); err != nil {
		return nil, fmt.Errorf("cannot pay activated ability cost: %w", err)
	}
	state.Log("Activating ability: " + ability.Description())
	// Mana abilities
	if ability.IsManaAbility() {
		if err := ability.Resolve(state, player); err != nil {
			return nil, fmt.Errorf("failed to resolve mana ability: %w", err)
		}
		_, ok := ability.Cost.(*TapCost)
		if ok {
			state.EmitEvent(Event{
				Type:   EventTapForMana,
				Source: ability.source,
			}, player)
		}
	} else {
		if err := state.Stack.Add(ability); err != nil {
			return nil, fmt.Errorf("failed to add ability to stack: %w", err)
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
			"ability activated: %s (%s)",
			ability.source.Name(),
			ability.Description(),
		),
	}, nil
}

// ActionAddManaFunc handles the add mana action. This is performed by the
// player to add mana to their mana pool. The target is the amount of mana
// to add. This is a cheat.
func ActionAddManaFunc(state *GameState, player *Player, target string) (*ActionResult, error) {
	mana := target
	if mana == "" {
		choices := []Choice{
			{
				Name:   "White",
				ID:     "{W}",
				Source: ChoiceCheat,
			}, {
				Name:   "Blue",
				ID:     "{U}",
				Source: ChoiceCheat,
			}, {
				Name:   "Black",
				ID:     "{B}",
				Source: ChoiceCheat,
			}, {
				Name:   "Red",
				ID:     "{R}",
				Source: ChoiceCheat,
			}, {
				Name:   "Green",
				ID:     "{G}",
				Source: ChoiceCheat,
			}, {
				Name:   "Colorless",
				ID:     "{C}",
				Source: ChoiceCheat,
			},
		}
		choice, err := player.Agent.ChooseOne(
			"Add mana to mana pool",
			ChoiceSourceCheat,
			choices,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose mana: %w", err)
		}
		mana = choice.ID
	}
	player.ManaPool.AddMana(mana)
	return &ActionResult{
		Message: fmt.Sprintf("%s mana added to pool", mana),
	}, nil
}

// ActionConjureFunc handles the conjure action. This is performed by the
// player to conjure a card. The target is the name of the card to conjure.
// This is a cheat.
func ActionConjureFunc(state *GameState, player *Player, target string) (*ActionResult, error) {
	if target == "" {
		return nil, errors.New("no card name provided")
	}
	cardPoolData, err := LoadCardPoolData(state.CardPool)
	if err != nil {
		return nil, fmt.Errorf("failed to load card pool data: %w", err)
	}
	cardData, ok := cardPoolData[target]
	if !ok {
		return nil, fmt.Errorf("card %s not found in card pool data", target)
	}
	card, err := NewCardFromCardData(cardData)
	if err != nil {
		return nil, fmt.Errorf("failed to create card %s: %w", target, err)
	}
	if err := player.Hand.Add(card); err != nil {
		return nil, fmt.Errorf("failed to add card to hand: %w", err)
	}
	state.Log("Conjured card: " + card.Name())
	return &ActionResult{
		Message: "conjured card: " + card.Name(),
	}, nil
}

// ActionDiscardFunc handles the discard action. This is performed
// automatically at the end of turn by during the clean up state. It can also
// be performed manually by the player if Cheat is enabled in the game state.
// The target is the number of cards to discard.
func ActionDiscardFunc(state *GameState, player *Player, target string) (*ActionResult, error) {
	n, err := strconv.Atoi(target)
	if err != nil {
		return nil, fmt.Errorf("failed to convert %s to int: %w", target, err)
	}
	// TODO:
	state.Discard(n, ActionDiscard, player)
	return &ActionResult{
		Message: "card discarded",
	}, nil
}

// ActionDrawFunc handles the draw action. This is performed automatically at
// the beginning of turn during the draw step. It can also be performed
// manually by the player if Cheat is enabled in the game state. The target is
// the number of cards to draw.
func ActionDrawFunc(state *GameState, player *Player, target string) (*ActionResult, error) {
	n := 1
	var err error
	if target != "" {
		n, err = strconv.Atoi(target)
		if err != nil {
			return nil, fmt.Errorf("failed to convert %s to int: %w", target, err)
		}
	}
	if err := state.Draw(n, player); err != nil {
		return nil, err
	}
	return &ActionResult{
		Message: "card drawn",
	}, nil
}

// ActionFindFunc handles the find action. This is performed by the player to
// find a card in their library. The target is the name of the card to find.
// This is a cheat.
func ActionFindFunc(state *GameState, player *Player, target string) (*ActionResult, error) {
	var card GameObject
	var err error
	if target != "" {
		card, err = FindFirstInZoneBy(player.Library, HasName(target))
		if err != nil {
			return nil, fmt.Errorf("failed to find card in library: %w", err)
		}
	} else {
		choices := CreateObjectChoices(
			player.Library.GetAll(),
			ZoneLibrary,
		)
		choice, err := player.Agent.ChooseOne(
			"Which card to find",
			ChoiceSourceCheat,
			AddOptionalChoice(choices),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose card: %w", err)
		}
		card, err = player.Library.Take(choice.ID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to take card from library: %w", err)
	}
	if err := player.Hand.Add(card); err != nil {
		return nil, fmt.Errorf("failed to add card to hand: %w", err)
	}
	return &ActionResult{
		Message: "found card: " + card.Name(),
	}, nil
}

// ActionLandDropFunc handles the land drop action. This is performed by the
// player. Is resets the land drop flag for the player.
// This is a cheat.
func ActionLandDropFunc(state *GameState, player *Player, target string) (*ActionResult, error) {
	player.LandDrop = false
	state.Log("Land drop reset")
	return &ActionResult{
		Message: "land drop reset",
	}, nil
}

// ActionPlayFunc handles the play action. This is performed by the player to
// play a card from their hand. The target is the name of the card to play.
func ActionPlayFunc(state *GameState, player *Player, target string) (result *ActionResult, err error) {
	var choices []Choice
	for _, zone := range player.Zones() {
		cards := zone.AvailableToPlay(state, player)
		cs := CreateObjectChoices(cards, zone.ZoneType())
		choices = append(choices, cs...)
	}
	if len(choices) == 0 {
		return nil, errors.New("no cards available to play")
	}
	// TODO: Rethink this
	if target != "" {
		var filteredChoices []Choice
		for _, choice := range choices {
			if choice.Name == target {
				filteredChoices = append(filteredChoices, choice)
			}
		}
		choices = filteredChoices
	}

	var choice Choice
	if target != "" && len(choices) == 1 {
		choice = choices[0]
	} else {

		choice, err = player.Agent.ChooseOne(
			"Which card to play",
			ActionPlay,
			AddOptionalChoice(choices),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose card: %w", err)
		}
		if choice.ID == ChoiceNone {
			return &ActionResult{Message: "No choice made"}, nil
		}
	}
	zone, err := player.GetZone(choice.Zone)
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
		return actionPlayLandFunc(state, player, c)
	}
	return actionCastSpellFunc(state, player, c, choice.Zone)
}

// ActionUntapFunc handles the untap action. This is performed automatically
// at the beginning of turn during the untap step. It can also be performed
// manually by the player if Cheat is enabled in the game state. The target
// is the name of the card to untap. If target == const(UntapAll), all
// permanents are untapped.
func ActionUntapFunc(state *GameState, player *Player, target string) (*ActionResult, error) {
	if target == UntapAll {
		player.Battlefield.UntapPermanents()
		return &ActionResult{Message: "all permanents untapped"}, nil
	}
	var err error
	var selectedObject GameObject
	if target != "" {
		selectedObject, err = FindFirstInZoneBy(
			player.Battlefield,
			And(IsTapped(), HasName(target)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to find permanent: %w", err)
		}
	}
	if selectedObject == nil {
		objects := FindInZoneBy(player.Battlefield, IsTapped())
		choices := CreateObjectChoices(objects, ZoneBattlefield)
		choice, err := player.Agent.ChooseOne(
			"Which permanent to untap",
			ActionUntap,
			AddOptionalChoice(choices),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose permanent: %w", err)
		}
		if choice.ID == ChoiceNone {
			// TODO: Make this a constant
			return &ActionResult{Message: "no choice made"}, nil
		}
		selectedObject, err = player.Battlefield.Get(choice.ID)
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
func ActionViewFunc(state *GameState, player *Player, target string) (*ActionResult, error) {
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
		choice, err = player.Agent.ChooseOne(
			"Which zone",
			NewChoiceSource("View Zone", "View Zone"),
			AddOptionalChoice(choices),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose zone: %w", err)
		}
		if choice.ID == ChoiceNone {
			return &ActionResult{Message: "No choice made"}, nil
		}
	} else {
		choice = Choice{Name: target}
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
}

// TODO: There's probably an abstraction we can set up for viewHand, viewBattlefield, and viewGraveyard
// probably a general abstraction for cards/permanents and zones
func viewHand(state *GameState, player *Player) (result *ActionResult, err error) {
	/*
		choices := state.Hand.CardChoices()
		choice, err := PlayerAgent.ChooseOne(
			"Which card",
			NewChoiceSource("View Hand", "View Hand"),
			AddOptionalChoice(choices),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose card: %w", err)
		}
		if choice.ID == ChoiceNone {
			return &ActionResult{Message: "No choice made"}, nil
		}
		card, err := state.Hand.GetCard(choice.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get card from hand: %w", err)
		}
		state.Log("viewed " + card.Name())
		result = &ActionResult{Message: fmt.Sprintf("CARD: %s :: %s :: %s", card.Name(), card.CardTypes(), card.RulesText())}
		return result, nil
	*/
	return &ActionResult{Message: "No choice made"}, nil
}

func viewBattlefield(state *GameState, player *Player) (result *ActionResult, err error) {
	/*
		var choices []Choice
		permanents := state.Battlefield.Permanents()
		for _, permanent := range permanents {
			choices = append(choices, Choice{Name: permanent.Name(), ID: permanent.ID()})
		}
		choice, err := PlayerAgent.ChooseOne(
			"Which card",
			NewChoiceSource("View Battlefield", "View Battlefield"),
			AddOptionalChoice(choices),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose card: %w", err)
		}
		if choice.ID == ChoiceNone {
			return &ActionResult{Message: "No choice made"}, nil
		}
		card, err := state.Battlefield.GetPermanent(choice.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get permanent from battlefield: %w", err)
		}
		state.Log("viewed " + card.Name())
		result = &ActionResult{Message: fmt.Sprintf("CARD: %s :: %s", card.Name, card.RulesText)}
		return result, nil
	*/
	return &ActionResult{Message: "No choice made"}, nil
}

func viewGraveyard(state *GameState, player *Player) (result *ActionResult, err error) {
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
		return &ActionResult{Message: "No choice made"}, nil
	}
	var selectedCard *Card
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
	result = &ActionResult{Message: fmt.Sprintf("CARD: %s :: %s", selectedCard.Name(), selectedCard.RulesText())}
	return result, nil
}

// TODO: Maybe this should be a method off of GameState
func actionPlayLandFunc(state *GameState, player *Player, card *Card) (result *ActionResult, err error) {
	if card.HasCardType(CardTypeLand) {
		if player.LandDrop {
			return nil, errors.New("land already played this turn")
		}
		player.LandDrop = true
	}
	state.Log("Played land: " + card.Name())
	player.Hand.Remove(card.ID())
	permanent, err := NewPermanent(card)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create permanent from %s: %w",
			card.Name(),
			err,
		)
	}
	player.Battlefield.Add(permanent)
	return &ActionResult{
		Message: "played land: " + card.Name(),
	}, nil
}

// TODO: Maybe this should be a method off of GameState
// or maybe a method off of Card, e.g. card.Cast() like Ability.Resolve()
func actionCastSpellFunc(state *GameState, player *Player, card *Card, zone string) (result *ActionResult, err error) {
	spell, err := NewSpell(card)
	if err != nil {
		return nil, fmt.Errorf("failed to create spell from %s: %w", card.Name(), err)
	}
	var spellCost Cost
	var isFlashback bool
	if zone == ZoneGraveyard {
		for _, ability := range spell.StaticAbilities() {
			if ability.ID != AbilityKeywordFlashback {
				continue
			}
			isFlashback = true
			for _, modifier := range ability.Modifiers {
				if modifier.Key != "Cost" {
					continue
				}
				spellCost, err = NewCost(modifier.Value, spell)
				if err != nil {
					return nil, fmt.Errorf("failed to create cost: %w", err)
				}
			}
		}
	} else {
		spellCost = spell.ManaCost()
	}
	var isReplicate bool
	var replicateCost Cost
	var replicateCount int
	for _, ability := range spell.StaticAbilities() {
		if ability.ID != AbilityKeywordReplicate {
			continue
		}
		for _, modifier := range ability.Modifiers {
			if modifier.Key != "Cost" {
				continue
			}
			isReplicate = true
			replicateCost, err = NewCost(modifier.Value, spell)
			if err != nil {
				return nil, fmt.Errorf("failed to create cost: %w", err)
			}
		}
	}
	if isReplicate {
		replicateCount, err = player.Agent.EnterNumber(
			fmt.Sprintf("Replicate how many times for %s", replicateCost.Description()),
			NewChoiceSource(AbilityKeywordReplicate, AbilityKeywordReplicate),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to confirm replicate cost: %w", err)
		}
	}
	for range replicateCount {
		spellCost = spellCost.Add(replicateCost)
	}
	if err := spellCost.Pay(state, player); err != nil {
		return nil, err
	}
	cardZone, err := player.GetZone(zone)
	if err != nil {
		return nil, fmt.Errorf("failed to get zone %s: %w", zone, err)
	}
	if err := cardZone.Remove(card.ID()); err != nil {
		return nil, fmt.Errorf("failed to remove card from hand: %w", err)
	}
	state.Log("Casing spell: " + card.Name())
	if isFlashback {
		spell.Flashback()
	}
	if err := state.Stack.Add(spell); err != nil {
		return nil, fmt.Errorf("failed to add spell to stack: %w", err)
	}
	if replicateCount > 0 {
		replicateAbility := BuildReplicateAbility(card, replicateCount)
		if err := state.Stack.Add(replicateAbility); err != nil {
			return nil, fmt.Errorf("failed to add replicate trigger to stack: %w", err)
		}
	}
	return &ActionResult{
		Message: "played card: " + card.Name(),
	}, nil
}

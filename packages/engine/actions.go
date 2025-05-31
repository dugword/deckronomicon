package engine

import (
	"deckronomicon/packages/game/action"
	"deckronomicon/packages/game/card"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/player"
	"fmt"
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
func (e *Engine) ResolveAction(act *action.Action, player *player.Player) (result ActionResult, err error) {
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
			Message: "player.Player passed",
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
		return ActionResult{
			// TODO: No .cards access
			Message: "Top Card: " + player.Library.Peek().Name(),
		}, nil
	case action.CheatShuffle:
		e.Log("CHEAT! Action: shuffle")
		player.Library.Shuffle()
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

// ActionActivateFunc handles the activate action. This is performed by the
// player to activate an ability of a permanent on the battlefield, or a card
// in hand or in the graveyard. The target is the name of the permanent or
// card.
// TODO: Support more than one activated ability
// TODO: Support activated abilities in hand and graveyard
func ActionActivateFunc(state *GameState, player *player.Player, target action.ActionTarget) (ActionResult, error) {
	/*
		available := player.GetAvailableToActivate(state)
		choices := CreateGroupedChoices(available)
		if len(choices) == 0 {
			// TODO: Do I need to check this?
			return nil, errors.New("no activated abilities available")
		}
		choice, err := player.Agent.ChooseOne(
			"Which ability to activate",
			ActionActivate,
			AddOptionalChoice(choices),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose ability: %w", err)
		}
		if choice.ID == ChoiceNone {
			return ActionResult{Message: "No choice made"}, nil
		}
		var ability *ActivatedAbility
		object, err := FindFirstBy(available[choice.Zone], HasID(choice.ID))
		if err != nil {
			return nil, fmt.Errorf("failed to find activated ability: %w", err)
		}
		ability, ok := object.(*ActivatedAbility)
		if !ok {
			return nil, fmt.Errorf("object is not an activated ability: %w", err)
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
		return ActionResult{
			Message: fmt.Sprintf(
				"ability activated: %s (%s)",
				ability.source.Name(),
				ability.Description(),
			),
		}, nil
	*/
	return ActionResult{}, nil
}

// ActionAddManaFunc handles the add mana action. This is performed by the
// player to add mana to their mana pool. The target is the amount of mana
// to add. This is a cheat.
func ActionAddManaFunc(state *GameState, player *player.Player, target action.ActionTarget) (ActionResult, error) {
	// TODO: This is a hack, we should probably have a better way to
	/*
		mana := target.Name
		if mana == "" {
			choices := []choice.Choice{
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
				ChoiceSourceCheat,
				choices,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to choose mana: %w", err)
			}
			mana = choice.ID
		}
		player.ManaPool.AddMana(mana)
		return ActionResult{
			Message: fmt.Sprintf("%s mana added to pool", mana),
		}, nil
	*/

	return ActionResult{}, nil
}

// ActionConjureFunc handles the conjure action. This is performed by the
// player to conjure a card. The target is the name of the card to conjure.
// This is a cheat.
func ActionConjureFunc(state *GameState, player *player.Player, target action.ActionTarget) (ActionResult, error) {
	/*
		cardName := target.Name
		if cardName == "" {
			return nil, errors.New("no card name provided")
		}
		cardPoolData, err := LoadCardPoolData(state.CardPool)
		if err != nil {
			return nil, fmt.Errorf("failed to load card pool data: %w", err)
		}
		cardData, ok := cardPoolData[cardName]
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
		return ActionResult{
			Message: "conjured card: " + card.Name(),
		}, nil
	*/
	return ActionResult{}, nil
}

// ActionDiscardFunc handles the discard action. This is performed
// automatically at the end of turn by during the clean up state. It can also
// be performed manually by the player if Cheat is enabled in the game state.
// The target is the number of cards to discard.
func ActionDiscardFunc(state *GameState, player *player.Player, target action.ActionTarget) (ActionResult, error) {
	/*
		n, err := strconv.Atoi(target.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to convert %s to int: %w", target, err)
		}
		// TODO:
		state.Discard(n, ActionDiscard, player)
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
	/*
		count := target.Name
		n := 1
		var err error
		if count != "" {
			n, err = strconv.Atoi(count)
			if err != nil {
				return nil, fmt.Errorf("failed to convert %s to int: %w", count, err)
			}
		}
		if err := state.Draw(n, player); err != nil {
			return nil, err
		}
		return ActionResult{
			Message: "card drawn",
		}, nil
	*/
	return ActionResult{}, nil
}

// ActionFindFunc handles the find action. This is performed by the player to
// find a card in their library. The target is the name of the card to find.
// This is a cheat.
func ActionFindFunc(state *GameState, player *player.Player, target action.ActionTarget) (ActionResult, error) {
	/*
		var card game.Object
		cardName := target.Name
		var err error
		if cardName != "" {
			card, err = FindFirstInZoneBy(player.Library, HasName(cardName))
			if err != nil {
				return nil, fmt.Errorf("failed to find card in library: %w", err)
			}
		} else {
			choices := CreateChoices(
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
		return ActionResult{
			Message: "found card: " + card.Name(),
		}, nil
	*/
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

// TODO This whole look up system is a bit of a hack. Should extract it out to a
// function and share it with action activate.
// ActionPlayFunc handles the play action. This is performed by the player to
// play a card from their hand. The target is the name of the card to play.
func ActionPlayFunc(state *GameState, player *player.Player, target action.ActionTarget) (result ActionResult, err error) {
	/*
		available := player.GetAvailableToPlay(state)
		choiceZone := ""
		var obj game.Object
		if target.ID != "" {
			for z, objects := range available {
				c, err := FindFirstBy(objects, HasID(target.ID))
				if errors.Is(err, ErrObjectNotFound) {
					continue
				}
				obj = c
				choiceZone = z
			}
		}

		if obj == nil {
			choices := CreateGroupedChoices(available)
			if len(choices) == 0 {
				return nil, errors.New("no cards available to play")
			}
			// TODO: Rethink this
			if target.Name != "" {
				var filteredChoices []Choice
				for _, choice := range choices {
					if choice.Name == target.Name {
						filteredChoices = append(filteredChoices, choice)
					}
				}
				choices = filteredChoices
			}

			var choice Choice
			if target.Name != "" && len(choices) == 1 {
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
					return ActionResult{Message: "No choice made"}, nil
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
			obj = card
			choiceZone = choice.Zone
		}
		c, ok := obj.(*card.Card)
		if !ok {
			return nil, fmt.Errorf("object is not a card: %w", err)
		}
		if c.HasCardType(game.CardTypeLand) {
			return actionPlayLandFunc(state, player, c)
		}
		return actionCastSpellFunc(state, player, c, choiceZone)
	*/
	return ActionResult{}, nil
}

// ActionUntapFunc handles the untap action. This is performed automatically
// at the beginning of turn during the untap step. It can also be performed
// manually by the player if Cheat is enabled in the game state. The target
// is the name of the card to untap. If target == const(UntapAll), all
// permanents are untapped.
func ActionUntapFunc(state *GameState, player *player.Player, target action.ActionTarget) (ActionResult, error) {
	/*
		permanentName := target.Name
		if permanentName == UntapAll {
			player.Battlefield.UntapPermanents()
			return ActionResult{Message: "all permanents untapped"}, nil
		}
		var err error
		var selectedObject game.Object
		if permanentName != "" {
			selectedObject, err = FindFirstInZoneBy(
				player.Battlefield,
				And(IsTapped(), HasName(permanentName)),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to find permanent: %w", err)
			}
		}
		if selectedObject == nil {
			objects := FindInZoneBy(player.Battlefield, IsTapped())
			choices := CreateChoices(objects, ZoneBattlefield)
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
				return ActionResult{Message: "no choice made"}, nil
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
		return ActionResult{
			Message: fmt.Sprintf("%s untapped", selectedObject.Name()),
		}, nil
	*/
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

// TODO: Maybe this should be a method off of GameState
func actionPlayLandFunc(state *GameState, player *player.Player, card *card.Card) (result ActionResult, err error) {
	/*
		if card.HasCardType(game.CardTypeLand) {
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
		return ActionResult{
			Message: "played land: " + card.Name(),
		}, nil
	*/
	return ActionResult{}, nil
}

// TODO: Maybe this should be a method off of GameState
// or maybe a method off of Card, e.g. card.Cast() like Ability.Resolve()
func actionCastSpellFunc(state *GameState, player *player.Player, card *card.Card, zone string) (result ActionResult, err error) {
	/*
		spell, err := NewSpell(card)
		if err != nil {
			return nil, fmt.Errorf("failed to create spell from %s: %w", card.Name(), err)
		}
		var spellCost Cost
		var isFlashback bool
		// var isSplice bool
		// var spliceCost Cost
		if zone == ZoneGraveyard {
			for _, ability := range spell.StaticAbilities() {
				if ability.ID != game.AbilityKeywordFlashback {
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
		var toSplice []game.Object
		if spell.HasSubtype(game.SubtypeArcane) {
			spliceCards := FindInZoneBy(
				player.Hand,
				And(
					HasStaticAbility(game.AbilityKeywordSplice),
					HasStaticAbilityModifier(
						game.AbilityKeywordSplice,
						EffectTag{Key: "Onto", Value: "Arcane"},
					),
				),
			)
			// var spliceCosts []Cost
			choices := CreateChoices(spliceCards, ZoneHand)
			chosen, err := player.Agent.ChooseMany(
				"Choose cards to splice onto the spell",
				NewChoiceSource(game.AbilityKeywordSplice, game.AbilityKeywordSplice),
				AddOptionalChoice(choices),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to choose card: %w", err)
			}
			for _, c := range chosen {
				if c.ID == ChoiceNone {
					break
				}
				spliceCard, err := player.Hand.Get(c.ID)
				if err != nil {
					return nil, fmt.Errorf("failed to get splice card from hand: %w", err)
				}
				for _, ability := range spliceCard.StaticAbilities() {
					if ability.ID != game.AbilityKeywordSplice {
						continue
					}
					for _, modifier := range ability.Modifiers {
						if modifier.Key != "Cost" {
							continue
						}
						spliceCost, err := NewCost(modifier.Value, spell)
						if err != nil {
							return nil, fmt.Errorf("failed to create splice cost: %w", err)
						}
						accept, err := player.Agent.Confirm(
							fmt.Sprintf("Splice %s onto %s for %s?",
								spliceCard.Name(),
								spell.Name(),
								spliceCost.Description(),
							),
							NewChoiceSource(game.AbilityKeywordSplice, game.AbilityKeywordSplice),
						)
						if err != nil {
							return nil, fmt.Errorf("failed to confirm splice cost: %w", err)
						}
						if !accept {
							continue
						}
						state.Log("adding splice to toSplice...")
						spellCost = AddCosts(spellCost, spliceCost)
						toSplice = append(toSplice, spliceCard)
					}
				}
			}
		}

		var isReplicate bool
		var replicateCost Cost
		var replicateCount int
		for _, ability := range spell.StaticAbilities() {
			if ability.ID != game.AbilityKeywordReplicate {
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
				NewChoiceSource(mtg.AbilityKeywordReplicate, mtg.AbilityKeywordReplicate),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to confirm replicate cost: %w", err)
			}
		}
		for range replicateCount {
			spellCost = AddCosts(spellCost, replicateCost)
		}
		fmt.Println("Spell cost: ", spellCost.Description())
		if err := spellCost.Pay(state, player); err != nil {
			return nil, err
		}
		state.Log("toSplice size: " + strconv.Itoa(len(toSplice)))
		for _, spliceCard := range toSplice {
			card, ok := spliceCard.(*card.Card)
			if !ok {
				return nil, fmt.Errorf("object is not a card: %w", err)
			}
			state.Log("Splicing " + card.Name() + " onto " + spell.Name())
			spell.Splice(card)
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
		for _, effect := range spell.SpellAbility().Effects {
			fmt.Println("Effect ID: ", effect.ID)
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
		return ActionResult{
			Message: "played card: " + card.Name(),
		}, nil
	*/

	return ActionResult{}, nil
}

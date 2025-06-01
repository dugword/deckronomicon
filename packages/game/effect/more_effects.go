package effect

import (
	"deckronomicon/packages/game/core"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const TagManaAbility = "ManaAbility"

// TODO  would be good to ensure all BuildEffect functions return the same
// type EffectBuilder func(source GameObject, spec definition.EffectSpec) (*Effect, error)

// BuildEffect creates an effect based on the provided definition.EffectSpec.
func BuildEffect(source query.Object, spec definition.EffectSpec) (*Effect, error) {
	switch spec.ID {
	case "AdditionalMana":
		return BuildEffectAdditionalMana(source, spec)
	case "AddMana":
		return BuildEffectAddMana(source, spec)
	case "CounterSpell":
		return BuildEffectCounterSpell(source, spec)
	case "Discard":
		return BuildEffectDiscard(source, spec)
	case "Draw":
		return BuildEffectDraw(source, spec)
	case "PutBackOnTop":
		return BuildEffectPutBackOnTop(source, spec)
	case "LookAndChoose":
		return BuildEffectLookAndChoose(source, spec)
	case "Mill":
		return BuildEffectMill(source, spec)
	case "Scry":
		return BuildEffectScry(source, spec)
	case "Search":
		return BuildEffectSearch(source, spec)
	case "Transmute":
		return BuildEffectTransmute(source, spec)
	case "Typecycling":
		return BuildEffectTypecycling(source, spec)
	case "ShuffleFromGraveyard":
		return BuildEffectShuffleFromGraveyard(source, spec)
	case "Tap":
		return BuildEffectTap(source, spec)
	case "TapOrUntap":
		return BuildEffectTapOrUntap(source, spec)
	default:
		panic("effect not implemented: " + spec.ID)
		return &Effect{
			id:          "UnknownEffect",
			description: fmt.Sprintf("unknown effect: %s", spec.ID),
			tags:        []Tag{{Key: "Unknown", Value: spec.ID}},
			Apply: func(state core.State, player core.Player) error {
				return nil
			},
		}, nil
	}
}

// BuildEffectDraw creates a draw effect based on the provided modifiers.
// Keys: Count, Type
// Default: Count: 1
func BuildEffectDraw(source query.Object, spec definition.EffectSpec) (*Effect, error) {
	effect := Effect{id: spec.ID}
	count := "1"
	var drawType string
	for _, modifier := range spec.Modifiers {
		if modifier.Key == "Count" {
			count = modifier.Value
		}
		if modifier.Key == "Type" {
			drawType = modifier.Value
		}
	}
	_, err := strconv.Atoi(count)
	if err != nil {
		return nil, fmt.Errorf("invalid count: %s", count)
	}
	tags := []Tag{{Key: "Draw", Value: count}}
	if drawType != "" {
		tags = append(tags, Tag{Key: "Type", Value: drawType})
	}
	effect.description = fmt.Sprintf("draw %d cards", count)
	effect.tags = tags
	effect.Apply = func(state core.State, player core.Player) error {
		/*
			if err := state.Draw(n, player); err != nil {
				fmt.Errorf("failed to draw %d cards: %w", n, err)
			}
		*/
		return nil
	}
	return &effect, nil
}

// BuildEffectAddMana creates an effect that adds mana to the player's mana
// pool.
// Supported Modifier Keys (concats multiple modifiers):
//   - Mana: <ManaString>
func BuildEffectAddMana(source query.Object, spec definition.EffectSpec) (*Effect, error) {
	effect := Effect{id: spec.ID}
	var mana string
	for _, modifier := range spec.Modifiers {
		if modifier.Key == "Mana" {
			mana += modifier.Value
		}
	}
	if mana == "" {
		return nil, errors.New("no mana string provided")
	}
	if !mtg.IsMana(mana) {
		return nil, fmt.Errorf("invalid mana string: %s", mana)
	}
	var tags []Tag
	for _, symbol := range mtg.ManaStringToManaSymbols(mana) {
		tags = append(tags, Tag{Key: TagManaAbility, Value: symbol})
	}
	effect.description = fmt.Sprintf("add %s", mana)
	effect.tags = tags
	effect.Apply = func(state core.State, player core.Player) error {
		if err := player.AddMana(mana); err != nil {
			return fmt.Errorf("failed to add mana: %w", err)
		}
		return nil
	}
	return &effect, nil
}

// BuildEffectAdditionalMana creates an effect that adds additional mana when
// a trigger happens, like tapping an island for mana.
// Supported Modifier Keys:
//   - Mana: <ManaString>
//   - Target: <subtype>
//   - Duration: <eventType>
func BuildEffectAdditionalMana(source query.Object, spec definition.EffectSpec) (*Effect, error) {
	var mana string
	var target string
	var duration string
	for _, modifier := range spec.Modifiers {
		if modifier.Key == "Mana" {
			mana += modifier.Value
		}
		if modifier.Key == "Target" {
			target = modifier.Value
		}
		if modifier.Key == "Duration" {
			duration = modifier.Value
		}
	}
	if mana == "" {
		return nil, errors.New("no mana string provided")
	}
	/*
		if !isMana(mana) {
			return nil, fmt.Errorf("invalid mana string: %s", mana)
		}
	*/
	if target == "" {
		return nil, errors.New("no target provided")
	}
	if duration != "EndOfTurn" {
		// return nil, errors.New("no duration provided")
		return nil, errors.New("only EndOfTurn duration is supported")
	}
	_, err := mtg.StringToSubtype(target)
	if err != nil {
		return nil, fmt.Errorf("invalid target subtype: %s", target)
	}
	effect := Effect{id: spec.ID}
	// id := getNextEventID()
	/*
		eventHandler := EventHandler{
			ID: id,
			Callback: func(event Event, state *Gamecore.State, player *core.Player) {
				// Move this into the register so I don't have to check for
				// it.
				if event.Type != EventTapForMana {
					return
				}
				if !event.Source.HasSubtype(subtype) {
					return
				}
				if err := player.ManaPool.AddMana(mana); err != nil {
					// TODO: Handle this better
					panic("failed to add mana: " + err.Error())
				}
				return
			},
		}
	*/
	var tags []Tag
	/*
		for _, symbol := range mtg.ManaStringToManaSymbols(mana) {
			tags = append(tags, Tag{Key: AbilityTagManaAbility, Value: symbol})
		}
	*/
	effect.tags = tags
	effect.description = fmt.Sprintf("add additional %s when you tap an %s for mana", mana, target)
	effect.Apply = func(state core.State, player core.Player) error {
		/*
			state.RegisterListenerUntil(
				eventHandler,
				EventEndStep,
			)
		*/
		return nil
	}
	return &effect, nil
}

// BuildEffectCounterSpell creates an effect that counters a spell.
// Supported Modifier Keys (last applies):
//   - Target: <target> CardType
//
// Multiple targets can be specified and will be OR'd together.
// If no target is specified, the effect will counter any spell.
func BuildEffectCounterSpell(source query.Object, spec definition.EffectSpec) (*Effect, error) {
	effect := Effect{id: spec.ID}
	var targetTypes []string
	for _, modifier := range spec.Modifiers {
		if modifier.Key == "Target" {
			targetTypes = append(targetTypes, modifier.Value)
		}
	}
	var cardTypes []mtg.CardType
	if len(targetTypes) != 0 {
		for _, target := range targetTypes {
			cardType, err := mtg.StringToCardType(target)
			if err != nil {
				return nil, fmt.Errorf("invalid target card type: %s", target)
			}
			cardTypes = append(cardTypes, cardType)
		}
	}
	effect.Apply = func(state core.State, player core.Player) error {
		/*
			resolvables := state.Stack.GetAll()
			var spells []query.Object
			for _, resolvable := range resolvables {
				spell, ok := resolvable.(*Spell)
				if !ok {
					continue
				}
				if len(cardTypes) == 0 {
					spells = append(spells, spell)
				}
				for _, cardType := range cardTypes {
					if spell.HasCardType(cardType) {
						spells = append(spells, spell)
						break // No need to check other types if one matches
					}
				}
			}
			choices := CreateChoices(spells, ZoneStack)
			if len(choices) == 0 {
				return fmt.Errorf("no spells to counter")
			}
			chosen, err := player.Agent.ChooseOne(
				"Choose a spell to counter",
				source,
				choices,
			)
			if err != nil {
				return fmt.Errorf("failed to choose spell: %w", err)
			}
			// Ensure the spell is on the stack
			if _, err := state.Stack.Get(chosen.ID); err != nil {
				// TODO: Handle fizzling consistently
				state.Log("spell fizzled - no targets")
				return nil
			}
			object, err := state.Stack.Take(chosen.ID)
			if err != nil {
				return fmt.Errorf("failed to remove spell from stack: %w", err)
			}
			spell, ok := object.(*Spell)
			if !ok {
				return fmt.Errorf("object is not a spell: %s", object.ID())
			}
			player.Graveyard.Add(spell.Card())
			return nil
		*/
		return nil
	}
	effect.description = fmt.Sprintf("counter a spell of type %s", strings.Join(targetTypes, ", "))
	var tags []Tag
	for _, target := range targetTypes {
		tags = append(tags, Tag{Key: "CounterSpell", Value: target})
	}
	effect.tags = tags
	return &effect, nil
}

// BuildEffectLookAndChoose creates an effect that allows the player to look
// at cards and choose one or more to put into their hand.
// Supported Modifier Keys (last applies):
//   - Look: <Cards to look at> Default: 1
//   - Choose: <Cards to choose> Default: 1
//   - Target<Type>: <target>
//   - Rest <Zone> Default: Library (Bottom) | Graveyard
//   - Order <order> Default: Any
func BuildEffectLookAndChoose(source query.Object, spec definition.EffectSpec) (*Effect, error) {
	effect := Effect{id: spec.ID}
	lookCount := "1"
	chooseCount := "1"
	restZone := "Library"
	// var filters []FilterFunc
	order := "Any"
	for _, modifier := range spec.Modifiers {
		if modifier.Key == "Look" {
			lookCount = modifier.Value
		}
		if modifier.Key == "Choose" {
			chooseCount = modifier.Value
		}
		if modifier.Key == "Rest" {
			restZone = modifier.Value
		}
		// TODO: Could probably just have this be "type" and figure out if
		// it's a card type or color
		if modifier.Key == "TargetCardType" {
			cardType, err := mtg.StringToCardType(modifier.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid target card type: %s", modifier.Value)
			}
			fmt.Println("Adding filter for card type:", cardType)
			//	filters = append(filters, HasCardType(cardType))
		}
		if modifier.Key == "Rest" {
			restZone = modifier.Value
		}
		if modifier.Key == "Order" {
			order = modifier.Value
		}
	}
	nLook, err := strconv.Atoi(lookCount)
	if err != nil {
		return nil, fmt.Errorf("invalid look count: %s", lookCount)
	}
	nChoose, err := strconv.Atoi(chooseCount)
	if err != nil {
		return nil, fmt.Errorf("invalid choose count: %s", chooseCount)
	}
	effect.Apply = func(state core.State, player core.Player) error {
		/*
			var objects []query.Object
			for range nLook {
				card, err := player.Library.TakeTop()
				if err != nil {
					// Not an error to look at an empty library
					if errors.Is(err, ErrLibraryEmpty) {
						break
					}
					return fmt.Errorf("failed to take top cards: %w", err)
				}
				objects = append(objects, card)
			}
			for _, object := range objects {
				if err := state.ActivePlayer.Revealed.Add(object); err != nil {
					return fmt.Errorf("failed to add card to revealed zone: %w", err)
				}
			}
			// TODO: I'm not sure I like this
			state.ActivePlayer.Agent.ReportState(state)
			defer func() {
				state.ActivePlayer.Revealed.Clear()
				state.ActivePlayer.Agent.ReportState(state)
			}()
			for range nChoose {
				choices := CreateChoices(
					FindBy(objects, Or(filters...)),
					ZoneRevealed,
				)
				choice, err := player.Agent.ChooseOne(
					fmt.Sprintf("Choose %d cards to put into your hand", nChoose),
					source,
					AddOptionalChoice(choices),
				)
				if err != nil {
					return fmt.Errorf("failed to choose cards: %w", err)
				}
				if choice.ID == ChoiceNone {
					state.Log("Player chose not to take any cards")
					return nil
				}

				taken, remaining, err := TakeFirstBy(objects, Or(HasID(choice.ID), HasID(choice.ID)))
				if err != nil {
					return fmt.Errorf("failed to find chosen card: %w", err)
				}
				if err := player.Hand.Add(taken); err != nil {
					return fmt.Errorf("failed to add card to hand: %w", err)
				}
				objects = remaining
			}
			switch restZone {
			case "Library":
				if order == "Any" {
					for len(objects) > 0 {
						choices := CreateChoices(objects, ZoneLibrary)
						chosen, err := player.Agent.ChooseOne(
							"Choose a card to put on the bottom of your library",
							source,
							choices,
						)
						if err != nil {
							return fmt.Errorf("failed to choose card: %w", err)
						}
						object, remaining, err := TakeFirstBy(objects, HasID(chosen.ID))
						if err != nil {
							return fmt.Errorf("failed to take card: %w", err)
						}
						objects = remaining
						if err := player.Library.Add(object); err != nil {
							return fmt.Errorf("failed to add card to library: %w", err)
						}
					}
				} else {
					for _, card := range objects {
						if err := player.Library.Add(card); err != nil {
							return fmt.Errorf("failed to add card to library: %w", err)
						}
					}
				}
			case "Graveyard":
				for _, card := range objects {
					if err := player.Graveyard.Add(card); err != nil {
						return fmt.Errorf("failed to add card to graveyard: %w", err)
					}
				}
			default:
				return fmt.Errorf("invalid rest zone: %s", restZone)
			}
		*/
		return nil
	}
	effect.tags = []Tag{
		{Key: "Look", Value: lookCount},
		{Key: "Choose", Value: chooseCount},
		{Key: "Rest", Value: restZone},
		{Key: "Order", Value: order},
	}
	effect.description = fmt.Sprintf("look at the top %d cards of your library, choose %d to put into your hand, and put the rest on the %s of your library in %s order", nLook, nChoose, restZone, order)
	return &effect, nil
}

// BuildEffectMill creates an effect that mills cards from the top of the
// library.
// Supported Modifier Keys (last applies):
//   - Count: <Cards to mill> Default: 1
//   - Target <target> Player | Self | Opponent
//
// TODO: Target needs to be selected on cast, not on resolution.
func BuildEffectMill(source query.Object, spec definition.EffectSpec) (*Effect, error) {
	effect := Effect{id: spec.ID}
	count := "1"
	var target string
	for _, modifier := range spec.Modifiers {
		if modifier.Key == "Count" {
			count = modifier.Value
		}
		if modifier.Key == "Target" {
			target = modifier.Value
		}
	}
	if target != "Player" && target != "Self" && target != "Opponent" {
		// TODO Support more targets
		// return nil, fmt.Errorf("invalid target: %s, must be Player, Self, or Opponent", target)
		return nil, fmt.Errorf("invalid target: %s, must be Player")
	}
	n, err := strconv.Atoi(count)
	if err != nil {
		return nil, fmt.Errorf("invalid count: %s", count)
	}
	effect.Apply = func(state core.State, player core.Player) error {
		// Get the target player
		// choices := []choose.Choice{}
		/*
			for _, player := range state.Players {
				choices = append(choices, Choice{
					Name: player.ID,
					ID:   player.ID,
				})
			}
			choice, err := player.Agent.ChooseOne("Choose a player to mill", source, choices)
			if err != nil {
				return fmt.Errorf("failed to choose player: %w", err)
			}
			var targetPlayer *Player
			// TODO Implement a DetPlayer(ID) function
			for _, player := range state.Players {
				if player.ID == choice.ID {
					targetPlayer = player
					break
				}
			}
			if targetPlayer == nil {
				return fmt.Errorf("failed to find player: %s", choice.ID)
			}
			for range n {
				taken, err := targetPlayer.Library.TakeTop()
				if err != nil {
					// Not an error to mill on an empty library
					if errors.Is(err, ErrLibraryEmpty) {
						return nil
					}
					return fmt.Errorf("failed to take top cards: %w", err)
				}
				targetPlayer.Graveyard.Add(taken)
			}
		*/
		return nil
	}
	effect.description = fmt.Sprintf("mill %d cards from your library", n)
	effect.tags = []Tag{{Key: "Mill", Value: count}}
	return &effect, nil
}

// BuildEffectPutBackOnTop creates an effect that puts cards back on top of
// the library.
// Supported Modifier Keys (last applies):
//   - Count: <Cards to put back> Default: 1
func BuildEffectPutBackOnTop(source query.Object, spec definition.EffectSpec) (*Effect, error) {
	effect := Effect{id: spec.ID}
	count := "1"
	for _, modifier := range spec.Modifiers {
		if modifier.Key != "Count" {
			continue
		}
		count = modifier.Value
	}
	n, err := strconv.Atoi(count)
	if err != nil {
		return nil, fmt.Errorf("invalid count: %s", count)
	}
	effect.Apply = func(state core.State, player core.Player) error {
		/*
			if err := PutNBackOnTop(state, 2, source, player); err != nil {
				return err
			}
		*/
		return nil
	}
	effect.description = fmt.Sprintf("put %d cards from your hand on top of your library in any order", n)
	effect.tags = []Tag{{Key: "PutBackOnTop", Value: count}}
	return &effect, nil
}

// BuildEffectScry creates an effect that allows the player to scry.
// Supported Modifier Keys (last applies):
//   - Count: <Cards to scry> Default: 1
func BuildEffectScry(source query.Object, spec definition.EffectSpec) (*Effect, error) {
	effect := Effect{id: spec.ID}
	count := "1"
	for _, modifier := range spec.Modifiers {
		if modifier.Key != "Count" {
			continue
		}
		count = modifier.Value
	}
	n, err := strconv.Atoi(count)
	if err != nil {
		return nil, fmt.Errorf("invalid count: %s", count)
	}
	effect.description = fmt.Sprintf("look at the top %d cards of your library, then put them back on top or bottom of your library in any order.", n)
	effect.Apply = func(state core.State, player core.Player) error {
		/*
			if err := Scry(state, source, n, player); err != nil {
				return err
			}
		*/
		return nil
	}
	effect.tags = []Tag{{Key: "Scry", Value: count}}
	return &effect, nil
}

// BuildEffectDiscard creates an effect that discards cards from the
// player's hand.
// Supported Modifier Keys (last applies):
//   - Count: <Cards to discard> Default: 1
//   - Delay: <Delay until> EndStep
func BuildEffectDiscard(source query.Object, spec definition.EffectSpec) (*Effect, error) {
	effect := Effect{id: spec.ID}
	count := "1"
	var delay string
	for _, modifier := range spec.Modifiers {
		if modifier.Key == "Count" {
			count = modifier.Value
		}
		if modifier.Key == "Delay" {
			delay = modifier.Value
		}
	}
	n, err := strconv.Atoi(count)
	if err != nil {
		return nil, fmt.Errorf("invalid count: %s", count)
	}
	// TODO: This could be more elegant
	var effectFunc func(state core.State, player core.Player) error

	fn := func(state core.State, player core.Player) error {
		/*
			if err := state.Discard(n, source, player); err != nil {
				return err
			}
		*/
		return nil
	}

	effectFunc = fn

	if delay == "EndStep" {
		/*
			var eventHandler EventHandler
			id := getNextEventID()
			eventHandler = EventHandler{
				ID: id,
				Callback: func(event Event, state *Gamecore.State, player *core.Player) {
					if event.Type != EventEndStep {
						return
					}
					// TODO Handle errors some how...
					_ = fn(state, player)
					state.DeregisterListener(id)
					return
				}}
			effectFunc = func(state *Gamecore.State, player *core.Player) error {
				state.RegisterListener(eventHandler)
				return nil
			}
		*/
	}
	tags := []Tag{{Key: "Discard", Value: count}}
	if delay != "" {
		tags = append(tags, Tag{Key: "Delay", Value: delay})
	}
	description := fmt.Sprintf("discard %d cards", n)
	if delay != "EndStep" {
		description += " at the beginning of your next end step"
	}
	effect.tags = tags
	effect.description = description
	effect.Apply = effectFunc
	return &effect, nil
}

// BuildEffectSearch creates an effect that allows the player to search
// their library for a card and put it into their hand.
// Supported Modifier Keys
//   - Target: <target> CardType | Color
//
// Multiple targets can be specified and will be AND'd together.gt
func BuildEffectSearch(source query.Object, spec definition.EffectSpec) (*Effect, error) {
	effect := Effect{id: spec.ID}
	/*
		var filters []FilterFunc
		for _, modifier := range spec.Modifiers {
			// TODO: Could probably just have this be "type" and figure out if
			// it's a card type or color
			if modifier.Key == "TargetCardType" {
				cardType, err := game.StringToCardType(modifier.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid target card type: %s", modifier.Value)
				}
				filters = append(filters, HasCardType(cardType))
			}
			if modifier.Key == "TargetColor" {
				color, err := game.StringToColor(modifier.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid target color: %s", modifier.Value)
				}
				filters = append(filters, HasColor(color))
			}
		}
		effect.Apply = func(state *Gamecore.State, player *core.Player) error {
			objects := player.Library.FindBy(And(filters...))
			if len(objects) == 0 {
				return fmt.Errorf("no cards found matching the specified targets")
			}
			choices := CreateChoices(objects, ZoneLibrary)
			chosen, err := player.Agent.ChooseOne(
				fmt.Sprintf("Choose a card to put into your hand"),
				source,
				choices,
			)
			if err != nil {
				return fmt.Errorf("failed to choose card: %w", err)
			}
			card, err := player.Library.Take(chosen.ID)
			if err != nil {
				return fmt.Errorf("failed to take card: %w", err)
			}
			player.Library.Shuffle()
			player.Hand.Add(card)

			return nil
		}
		effect.Description = fmt.Sprintf("search library for a card matching the specified targets")
		var tags []EffectTag
		for _, modifier := range spec.Modifiers {
			if strings.HasPrefix(modifier.Key, "Target") {
				tags = append(tags, EffectTag{Key: "Search", Value: modifier.Value})
			}
		}
	*/
	// effect.tags = tags
	return &effect, nil
}

// BuildEffectTransmute creates an effect that allows the player to transmute
// a card from their hand.
// Supported Modifier Keys (last applies):
func BuildEffectTransmute(source query.Object, spec definition.EffectSpec) (*Effect, error) {
	/*
			effect := Effect{ID: spec.ID}
			card, ok := source.(*Card)
			if !ok {
				return nil, fmt.Errorf("source is not a card: %T", source)
			}
			effect.Apply = func(state *Gamecore.State, player *core.Player) error {
				objects := FindInZoneBy(player.Library, HasManaValue(card.ManaValue()))
				if len(objects) == 0 {
					return fmt.Errorf(
						"no cards with mana value %s found", card.ManaValue(),
					)
				}
				choices := CreateChoices(objects, ZoneLibrary)
				chosen, err := player.Agent.ChooseOne(
					fmt.Sprintf("Choose a card to put into your hand"),
					source,
					choices,
				)
				if err != nil {
					return fmt.Errorf("failed to choose card: %w", err)
				}
				card, err := player.Library.Take(chosen.ID)
				if err != nil {
					return fmt.Errorf("failed to take card: %w", err)
				}
				player.Library.Shuffle()
				player.Hand.Add(card)
				return nil
			}
			effect.Description = fmt.Sprintf("search library for a card of mana value %d", card.ManaValue())
			effect.Tags = []EffectTag{{Key: "Transmute", Value: strconv.Itoa(card.ManaValue())}}
		return &effect, nil
	*/
	return nil, nil
}

// BuildEffectTypecycling creates an effect that search cards from the library.
// Supported Modifier Keys (last applies):
//   - Subtype <subtype>
func BuildEffectTypecycling(source query.Object, spec definition.EffectSpec) (*Effect, error) {
	/*
			effect := Effect{ID: spec.ID}
			var subtype string
			for _, modifier := range spec.Modifiers {
				if modifier.Key == "Subtype" {
					subtype = modifier.Value
				}
			}
			if subtype == "" {
				return nil, errors.New("no subtype provided")
			}
			subtypeEnum, err := game.StringToSubtype(subtype)
			if err != nil {
				return nil, fmt.Errorf("invalid subtype: %s", subtype)
			}
			effect.Apply = func(state *Gamecore.State, player *core.Player) error {
				objects := FindInZoneBy(player.Library, HasSubtype(subtypeEnum))
				if len(objects) == 0 {
					return fmt.Errorf("no cards of subtype %s found", subtype)
				}
				choices := CreateChoices(objects, ZoneLibrary)
				chosen, err := player.Agent.ChooseOne(
					fmt.Sprintf("Choose a card to put into your hand"),
					source,
					choices,
				)
				if err != nil {
					return fmt.Errorf("failed to choose card: %w", err)
				}
				card, err := player.Library.Take(chosen.ID)
				if err != nil {
					return fmt.Errorf("failed to take card: %w", err)
				}
				player.Library.Shuffle()
				player.Hand.Add(card)
				return nil
			}
			effect.Description = fmt.Sprintf("search library for a card of subtype %s", subtype)
			effect.Tags = []EffectTag{{Key: "Typecycling", Value: subtype}}
		return &effect, nil
	*/
	return nil, nil
}

// BuildEffectShuffleFromGraveyard creates an effect that shuffles cards from
// the graveyard back into the library.
// Supported Modifier Keys (last applies):
//   - Count: <Cards to shuffle> Default: 1
func BuildEffectShuffleFromGraveyard(source query.Object, spec definition.EffectSpec) (*Effect, error) {
	/*
			effect := Effect{ID: spec.ID}
			count := "1"
			for _, modifier := range spec.Modifiers {
				if modifier.Key == "Count" {
					count = modifier.Value
				}
			}
			n, err := strconv.Atoi(count)
			if err != nil {
				return nil, fmt.Errorf("invalid count: %s", count)
			}
			effect.Apply = func(state *Gamecore.State, player *core.Player) error {
				for range n {
					choices := CreateChoices(player.Graveyard.GetAll(), ZoneGraveyard)
					if len(choices) == 0 {
						break
					}
					chosen, err := player.Agent.ChooseOne(
						"Choose cards to shuffle into your library",
						source,
						AddOptionalChoice(choices),
					)
					if err != nil {
						return fmt.Errorf("failed to choose cards: %w", err)
					}
					if chosen.ID == ChoiceNone {
						break
					}
					card, err := player.Graveyard.Take(chosen.ID)
					if err != nil {
						return fmt.Errorf("failed to take card from graveyard: %w", err)
					}
					player.Library.Add(card)
					player.Library.Shuffle()
				}
				return nil
			}
			effect.Description = fmt.Sprintf("shuffle %d cards from your graveyard into your library", n)
			effect.Tags = []EffectTag{{Key: "ShuffleFromGraveyard", Value: count}}
		return &effect, nil
	*/
	return nil, nil
}

// BuildEffectTap creates an effect that taps a card.
// Supported Modifier Keys (last applies):
//   - Target: Permanent
//
// TODO: Support other targets
func BuildEffectTap(source query.Object, spec definition.EffectSpec) (*Effect, error) {
	/*
		effect := Effect{ID: spec.ID}
		var target string
		for _, modifier := range spec.Modifiers {
			if modifier.Key == "Target" {
				target = modifier.Value
			}
		}
		if target == "" {
			return nil, errors.New("no target provided")
		}
		if target != "Permanent" {
			return nil, fmt.Errorf("only Permanent target is supported: %s", target)
		}
		effect.Apply = func(state *Gamecore.State, player *core.Player) error {
			cards := player.Battlefield.GetAll()
			if len(cards) == 0 {
				// TODO: Spells can't be cast without targets
				return errors.New("no available targets")
			}
			choices := CreateChoices(cards, ZoneBattlefield)
			chosen, err := player.Agent.ChooseOne(
				fmt.Sprintf("Choose a card to tap"),
				source,
				choices,
			)
			if err != nil {
				return fmt.Errorf("failed to choose card: %w", err)
			}
			permanent, err := player.Battlefield.Get(chosen.ID)
			if err != nil {
				return fmt.Errorf("failed to get permanent: %w", err)
			}
			p, ok := permanent.(*Permanent)
			if !ok {
				return fmt.Errorf("object is not a permanent: %s", chosen.ID)
			}
			if err := p.Tap(); err != nil {
				if errors.Is(err, ErrAlreadyTapped) {
					// It's not an error to tap a card that's already tapped.
					return nil
				}
				return fmt.Errorf("failed to tap card: %w", err)
			}
			return nil
		}
		effect.Description = fmt.Sprintf("tap a card of type %s", target)
		effect.Tags = []EffectTag{{Key: "Tap", Value: target}}
	*/
	// return &effect, nil
	return nil, nil
}

// BuildEffectTapOrUntap creates an effect that taps or untaps a card.
// Supported Modifier Keys (last applies):
// TODO
//   - Target: Permanent // Permanent should be the default and only need to
//     specify if there's a more specific choice
func BuildEffectTapOrUntap(source query.Object, spec definition.EffectSpec) (*Effect, error) {
	/*
		effect := Effect{ID: spec.ID}
		var target string
		for _, modifier := range spec.Modifiers {
			if modifier.Key == "Target" {
				target = modifier.Value
			}
		}
		if target == "" {
			return nil, errors.New("no target provided")
		}
		if target != "Permanent" {
			return nil, fmt.Errorf("only Permanent target is supported: %s", target)
		}
		effect.Apply = func(state *Gamecore.State, player *core.Player) error {
			cards := player.Battlefield.GetAll()
			if len(cards) == 0 {
				return errors.New("no available targets")
			}
			choices := CreateChoices(cards, ZoneBattlefield)
			chosen, err := player.Agent.ChooseOne(
				fmt.Sprintf("Choose a card to tap or untap"),
				source,
				choices,
			)
			if err != nil {
				return fmt.Errorf("failed to choose card: %w", err)
			}
			permanent, err := player.Battlefield.Get(chosen.ID)
			if err != nil {
				return fmt.Errorf("failed to get permanent: %w", err)
			}
			p, ok := permanent.(*Permanent)
			if !ok {
				return fmt.Errorf("object is not a permanent: %s", chosen.ID)
			}
			var action string
			if p.IsTapped() {
				action = "untap"
				accept, err := player.Agent.Confirm(
					fmt.Sprintf("Do you want to untap %s?", p.Name()),
					source,
				)
				if err != nil {
					return fmt.Errorf("failed to confirm untap: %w", err)
				}
				if accept {
					p.Untap()
				}
			} else {
				action = "tap"
				accept, err := player.Agent.Confirm(
					fmt.Sprintf("Do you want to tap %s?", p.Name()),
					source,
				)
				if err != nil {
					return fmt.Errorf("failed to confirm tap: %w", err)
				}
				if accept {
					if err := p.Tap(); err != nil {
						if errors.Is(err, ErrAlreadyTapped) {
							return nil // Not an error to tap a card that's already tapped.
						}
						return fmt.Errorf("failed to tap card: %w", err)
					}
				}
			}
			state.Log(fmt.Sprintf("%s %s", action, p.Name()))
			return nil
		}
		effect.Description = fmt.Sprintf("tap or untap a card of type %s", target)
		effect.Tags = []EffectTag{{Key: "TapOrUntap", Value: target}}
	*/
	// return &effect, nil
	return nil, nil
}

func Scry(state core.State, source query.Object, n int, player core.Player) error {
	/*
		var taken []query.Object
		for range n {
			card, err := player.Library.TakeTop()
			if err != nil {
				// Not an error to scry on an empty library
				if errors.Is(err, ErrLibraryEmpty) {
					break
				}
				return fmt.Errorf("failed to take top card: %w", err)
			}
			taken = append(taken, card)
		}
		used := map[string]bool{}

		for range len(taken) {
			// Build option list from unplaced cards
			var choices []Choice
			for _, card := range taken {
				if !used[card.ID()] {
					choices = append(choices, Choice{
						Name: card.Name(),
						ID:   card.ID(),
					})
				}
			}
			chosen, err := player.Agent.ChooseOne(
				"Choose a card to place",
				source,
				choices,
			)
			if err != nil {
				return fmt.Errorf("failed to choose card: %w", err)
			}
			var chosenCard query.Object
			for _, card := range taken {
				if card.ID() == chosen.ID {
					chosenCard = card
					break
				}
			}
			if chosenCard == nil {
				return fmt.Errorf("failed to find chosen card: %s", chosen.ID)
			}
			used[chosen.ID] = true
			// TODO: Maybe have a global set of constants for choices like this
			const ChoiceTop = "Top"
			const ChoiceBottom = "Bottom"
			topBottomchoices := []Choice{
				{
					Name: ChoiceTop,
					ID:   ChoiceTop,
				},
				{
					Name: ChoiceBottom,
					ID:   ChoiceBottom,
				},
			}
			placement, err := player.Agent.ChooseOne(
				fmt.Sprintf("Place %s on top or bottom of your library?", chosenCard.Name()),
				source,
				topBottomchoices,
			)
			if err != nil {
				return fmt.Errorf("failed to choose placement: %w", err)
			}
			if placement.ID == ChoiceTop {
				player.Library.AddTop(chosenCard)
			} else {
				player.Library.Add(chosenCard)
			}
		}
	*/
	return nil
}

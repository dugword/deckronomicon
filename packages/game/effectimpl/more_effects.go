package effectimpl

import (
	"deckronomicon/packages/game/core"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/query"
	"fmt"
	"strconv"
)

type State interface {
	Stack() query.View
}

type Player interface {
	DrawCard() (string, error)
	// TakeTopCard() (*object.Card, error)
}

// TODO  would be good to ensure all BuildEffect functions return the same
// type EffectBuilder func(source GameObject, spec definition.EffectSpec) (*Effect, error)

// BuildEffect creates an effect based on the provided definition.EffectSpec.
func BuildEffect(source query.Object, spec definition.EffectSpec) (*effect.Effect, error) {
	switch spec.ID {
	/*
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
	*/
	default:
		panic("effect not implemented: " + spec.ID)
		/*
			return &Effect{
				id:          "UnknownEffect",
				description: fmt.Sprintf("unknown effect: %s", spec.ID),
				tags:        []core.Tag{{Key: "Unknown", Value: spec.ID}},
				Apply: func(state core.State, player core.Player) error {
					return nil
				},
			}, nil
		*/
	}
}

// BuildEffectMill creates an effect that mills cards from the top of the
// library.
// Supported Modifier Keys (last applies):
//   - Count: <Cards to mill> Default: 1
//   - Target <target> Player | Self | Opponent
//
// TODO: Target needs to be selected on cast, not on resolution.
func BuildEffectMill(source query.Object, spec definition.EffectSpec) (*effect.Effect, error) {
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
	_ = func(state core.State, player core.Player) error {
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
	eff := effect.New(
		spec.ID,
		fmt.Sprintf("mill %d cards from your library", n),
		[]core.Tag{{Key: "Mill", Value: count}},
		func(core.State, core.Player) error { return nil },
	)
	return eff, nil
}

// //
// // BuildEffectPutBackOnTop creates an effect that puts cards back on top of
// // the library.
// // Supported Modifier Keys (last applies):
// //   - Count: <Cards to put back> Default: 1
// func BuildEffectPutBackOnTop(source query.Object, spec definition.EffectSpec) (*effect.Effect, error) {
// 	count := "1"
// 	for _, modifier := range spec.Modifiers {
// 		if modifier.Key != "Count" {
// 			continue
// 		}
// 		count = modifier.Value
// 	}
// 	n, err := strconv.Atoi(count)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid count: %s", count)
// 	}
// 	effect.Apply = func(state core.State, player core.Player) error {
// 		/*
// 			if err := PutNBackOnTop(state, 2, source, player); err != nil {
// 				return err
// 			}
// 		*/
// 		return nil
// 	}
// 	effect.description = fmt.Sprintf("put %d cards from your hand on top of your library in any order", n)
// 	effect.tags = []core.Tag{{Key: "PutBackOnTop", Value: count}}
// 	return &effect, nil
// }

// // BuildEffectScry creates an effect that allows the player to scry.
// // Supported Modifier Keys (last applies):
// //   - Count: <Cards to scry> Default: 1
// func BuildEffectScry(source query.Object, spec definition.EffectSpec) (*Effect, error) {
// 	effect := Effect{id: spec.ID}
// 	count := "1"
// 	for _, modifier := range spec.Modifiers {
// 		if modifier.Key != "Count" {
// 			continue
// 		}
// 		count = modifier.Value
// 	}
// 	n, err := strconv.Atoi(count)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid count: %s", count)
// 	}
// 	effect.description = fmt.Sprintf("look at the top %d cards of your library, then put them back on top or bottom of your library in any order.", n)
// 	effect.Apply = func(state core.State, player core.Player) error {
// 		/*
// 			if err := Scry(state, source, n, player); err != nil {
// 				return err
// 			}
// 		*/
// 		return nil
// 	}
// 	effect.tags = []core.Tag{{Key: "Scry", Value: count}}
// 	return &effect, nil
// }

// // BuildEffectDiscard creates an effect that discards cards from the
// // player's hand.
// // Supported Modifier Keys (last applies):
// //   - Count: <Cards to discard> Default: 1
// //   - Delay: <Delay until> EndStep
// func BuildEffectDiscard(source query.Object, spec definition.EffectSpec) (*Effect, error) {
// 	effect := Effect{id: spec.ID}
// 	count := "1"
// 	var delay string
// 	for _, modifier := range spec.Modifiers {
// 		if modifier.Key == "Count" {
// 			count = modifier.Value
// 		}
// 		if modifier.Key == "Delay" {
// 			delay = modifier.Value
// 		}
// 	}
// 	n, err := strconv.Atoi(count)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid count: %s", count)
// 	}
// 	// TODO: This could be more elegant
// 	var effectFunc func(state core.State, player core.Player) error

// 	fn := func(state core.State, player core.Player) error {
// 		/*
// 			if err := state.Discard(n, source, player); err != nil {
// 				return err
// 			}
// 		*/
// 		return nil
// 	}

// 	effectFunc = fn

// 	if delay == "EndStep" {
// 		/*
// 			var eventHandler EventHandler
// 			id := getNextEventID()
// 			eventHandler = EventHandler{
// 				ID: id,
// 				Callback: func(event Event, state *Gamecore.State, player *core.Player) {
// 					if event.Type != EventEndStep {
// 						return
// 					}
// 					// TODO Handle errors some how...
// 					_ = fn(state, player)
// 					state.DeregisterListener(id)
// 					return
// 				}}
// 			effectFunc = func(state *Gamecore.State, player *core.Player) error {
// 				state.RegisterListener(eventHandler)
// 				return nil
// 			}
// 		*/
// 	}
// 	tags := []core.Tag{{Key: "Discard", Value: count}}
// 	if delay != "" {
// 		tags = append(tags, core.Tag{Key: "Delay", Value: delay})
// 	}
// 	description := fmt.Sprintf("discard %d cards", n)
// 	if delay != "EndStep" {
// 		description += " at the beginning of your next end step"
// 	}
// 	effect.tags = tags
// 	effect.description = description
// 	effect.Apply = effectFunc
// 	return &effect, nil
// }

// // BuildEffectSearch creates an effect that allows the player to search
// // their library for a card and put it into their hand.
// // Supported Modifier Keys
// //   - Target: <target> CardType | Color
// //
// // Multiple targets can be specified and will be AND'd together.gt
// func BuildEffectSearch(source query.Object, spec definition.EffectSpec) (*Effect, error) {
// 	effect := Effect{id: spec.ID}
// 	/*
// 		var filters []FilterFunc
// 		for _, modifier := range spec.Modifiers {
// 			// TODO: Could probably just have this be "type" and figure out if
// 			// it's a card type or color
// 			if modifier.Key == "TargetCardType" {
// 				cardType, err := game.StringToCardType(modifier.Value)
// 				if err != nil {
// 					return nil, fmt.Errorf("invalid target card type: %s", modifier.Value)
// 				}
// 				filters = append(filters, HasCardType(cardType))
// 			}
// 			if modifier.Key == "TargetColor" {
// 				color, err := game.StringToColor(modifier.Value)
// 				if err != nil {
// 					return nil, fmt.Errorf("invalid target color: %s", modifier.Value)
// 				}
// 				filters = append(filters, HasColor(color))
// 			}
// 		}
// 		effect.Apply = func(state *Gamecore.State, player *core.Player) error {
// 			objects := player.Library.FindBy(And(filters...))
// 			if len(objects) == 0 {
// 				return fmt.Errorf("no cards found matching the specified targets")
// 			}
// 			choices := CreateChoices(objects, ZoneLibrary)
// 			chosen, err := player.Agent.ChooseOne(
// 				fmt.Sprintf("Choose a card to put into your hand"),
// 				source,
// 				choices,
// 			)
// 			if err != nil {
// 				return fmt.Errorf("failed to choose card: %w", err)
// 			}
// 			card, err := player.Library.Take(chosen.ID)
// 			if err != nil {
// 				return fmt.Errorf("failed to take card: %w", err)
// 			}
// 			player.Library.Shuffle()
// 			player.Hand.Add(card)

// 			return nil
// 		}
// 		effect.Description = fmt.Sprintf("search library for a card matching the specified targets")
// 		var tags []EffectTag
// 		for _, modifier := range spec.Modifiers {
// 			if strings.HasPrefix(modifier.Key, "Target") {
// 				tags = append(tags, EffectTag{Key: "Search", Value: modifier.Value})
// 			}
// 		}
// 	*/
// 	// effect.tags = tags
// 	return &effect, nil
// }

// // BuildEffectTransmute creates an effect that allows the player to transmute
// // a card from their hand.
// // Supported Modifier Keys (last applies):
// func BuildEffectTransmute(source query.Object, spec definition.EffectSpec) (*Effect, error) {
// 	/*
// 			effect := Effect{ID: spec.ID}
// 			card, ok := source.(*Card)
// 			if !ok {
// 				return nil, fmt.Errorf("source is not a card: %T", source)
// 			}
// 			effect.Apply = func(state *Gamecore.State, player *core.Player) error {
// 				objects := FindInZoneBy(player.Library, HasManaValue(card.ManaValue()))
// 				if len(objects) == 0 {
// 					return fmt.Errorf(
// 						"no cards with mana value %s found", card.ManaValue(),
// 					)
// 				}
// 				choices := CreateChoices(objects, ZoneLibrary)
// 				chosen, err := player.Agent.ChooseOne(
// 					fmt.Sprintf("Choose a card to put into your hand"),
// 					source,
// 					choices,
// 				)
// 				if err != nil {
// 					return fmt.Errorf("failed to choose card: %w", err)
// 				}
// 				card, err := player.Library.Take(chosen.ID)
// 				if err != nil {
// 					return fmt.Errorf("failed to take card: %w", err)
// 				}
// 				player.Library.Shuffle()
// 				player.Hand.Add(card)
// 				return nil
// 			}
// 			effect.Description = fmt.Sprintf("search library for a card of mana value %d", card.ManaValue())
// 			effect.Tags = []EffectTag{{Key: "Transmute", Value: strconv.Itoa(card.ManaValue())}}
// 		return &effect, nil
// 	*/
// 	return nil, nil
// }

// // BuildEffectTypecycling creates an effect that search cards from the library.
// // Supported Modifier Keys (last applies):
// //   - Subtype <subtype>
// func BuildEffectTypecycling(source query.Object, spec definition.EffectSpec) (*Effect, error) {
// 	/*
// 			effect := Effect{ID: spec.ID}
// 			var subtype string
// 			for _, modifier := range spec.Modifiers {
// 				if modifier.Key == "Subtype" {
// 					subtype = modifier.Value
// 				}
// 			}
// 			if subtype == "" {
// 				return nil, errors.New("no subtype provided")
// 			}
// 			subtypeEnum, err := game.StringToSubtype(subtype)
// 			if err != nil {
// 				return nil, fmt.Errorf("invalid subtype: %s", subtype)
// 			}
// 			effect.Apply = func(state *Gamecore.State, player *core.Player) error {
// 				objects := FindInZoneBy(player.Library, HasSubtype(subtypeEnum))
// 				if len(objects) == 0 {
// 					return fmt.Errorf("no cards of subtype %s found", subtype)
// 				}
// 				choices := CreateChoices(objects, ZoneLibrary)
// 				chosen, err := player.Agent.ChooseOne(
// 					fmt.Sprintf("Choose a card to put into your hand"),
// 					source,
// 					choices,
// 				)
// 				if err != nil {
// 					return fmt.Errorf("failed to choose card: %w", err)
// 				}
// 				card, err := player.Library.Take(chosen.ID)
// 				if err != nil {
// 					return fmt.Errorf("failed to take card: %w", err)
// 				}
// 				player.Library.Shuffle()
// 				player.Hand.Add(card)
// 				return nil
// 			}
// 			effect.Description = fmt.Sprintf("search library for a card of subtype %s", subtype)
// 			effect.Tags = []EffectTag{{Key: "Typecycling", Value: subtype}}
// 		return &effect, nil
// 	*/
// 	return nil, nil
// }

// // BuildEffectShuffleFromGraveyard creates an effect that shuffles cards from
// // the graveyard back into the library.
// // Supported Modifier Keys (last applies):
// //   - Count: <Cards to shuffle> Default: 1
// func BuildEffectShuffleFromGraveyard(source query.Object, spec definition.EffectSpec) (*Effect, error) {
// 	/*
// 			effect := Effect{ID: spec.ID}
// 			count := "1"
// 			for _, modifier := range spec.Modifiers {
// 				if modifier.Key == "Count" {
// 					count = modifier.Value
// 				}
// 			}
// 			n, err := strconv.Atoi(count)
// 			if err != nil {
// 				return nil, fmt.Errorf("invalid count: %s", count)
// 			}
// 			effect.Apply = func(state *Gamecore.State, player *core.Player) error {
// 				for range n {
// 					choices := CreateChoices(player.Graveyard.GetAll(), ZoneGraveyard)
// 					if len(choices) == 0 {
// 						break
// 					}
// 					chosen, err := player.Agent.ChooseOne(
// 						"Choose cards to shuffle into your library",
// 						source,
// 						AddOptionalChoice(choices),
// 					)
// 					if err != nil {
// 						return fmt.Errorf("failed to choose cards: %w", err)
// 					}
// 					if chosen.ID == ChoiceNone {
// 						break
// 					}
// 					card, err := player.Graveyard.Take(chosen.ID)
// 					if err != nil {
// 						return fmt.Errorf("failed to take card from graveyard: %w", err)
// 					}
// 					player.Library.Add(card)
// 					player.Library.Shuffle()
// 				}
// 				return nil
// 			}
// 			effect.Description = fmt.Sprintf("shuffle %d cards from your graveyard into your library", n)
// 			effect.Tags = []EffectTag{{Key: "ShuffleFromGraveyard", Value: count}}
// 		return &effect, nil
// 	*/
// 	return nil, nil
// }

// // BuildEffectTap creates an effect that taps a card.
// // Supported Modifier Keys (last applies):
// //   - Target: Permanent
// //
// // TODO: Support other targets
// func BuildEffectTap(source query.Object, spec definition.EffectSpec) (*Effect, error) {
// 	/*
// 		effect := Effect{ID: spec.ID}
// 		var target string
// 		for _, modifier := range spec.Modifiers {
// 			if modifier.Key == "Target" {
// 				target = modifier.Value
// 			}
// 		}
// 		if target == "" {
// 			return nil, errors.New("no target provided")
// 		}
// 		if target != "Permanent" {
// 			return nil, fmt.Errorf("only Permanent target is supported: %s", target)
// 		}
// 		effect.Apply = func(state *Gamecore.State, player *core.Player) error {
// 			cards := player.Battlefield.GetAll()
// 			if len(cards) == 0 {
// 				// TODO: Spells can't be cast without targets
// 				return errors.New("no available targets")
// 			}
// 			choices := CreateChoices(cards, ZoneBattlefield)
// 			chosen, err := player.Agent.ChooseOne(
// 				fmt.Sprintf("Choose a card to tap"),
// 				source,
// 				choices,
// 			)
// 			if err != nil {
// 				return fmt.Errorf("failed to choose card: %w", err)
// 			}
// 			permanent, err := player.Battlefield.Get(chosen.ID)
// 			if err != nil {
// 				return fmt.Errorf("failed to get permanent: %w", err)
// 			}
// 			p, ok := permanent.(*Permanent)
// 			if !ok {
// 				return fmt.Errorf("object is not a permanent: %s", chosen.ID)
// 			}
// 			if err := p.Tap(); err != nil {
// 				if errors.Is(err, ErrAlreadyTapped) {
// 					// It's not an error to tap a card that's already tapped.
// 					return nil
// 				}
// 				return fmt.Errorf("failed to tap card: %w", err)
// 			}
// 			return nil
// 		}
// 		effect.Description = fmt.Sprintf("tap a card of type %s", target)
// 		effect.Tags = []EffectTag{{Key: "Tap", Value: target}}
// 	*/
// 	// return &effect, nil
// 	return nil, nil
// }

// // BuildEffectTapOrUntap creates an effect that taps or untaps a card.
// // Supported Modifier Keys (last applies):
// // TODO
// //   - Target: Permanent // Permanent should be the default and only need to
// //     specify if there's a more specific choice
// func BuildEffectTapOrUntap(source query.Object, spec definition.EffectSpec) (*Effect, error) {
// 	/*
// 		effect := Effect{ID: spec.ID}
// 		var target string
// 		for _, modifier := range spec.Modifiers {
// 			if modifier.Key == "Target" {
// 				target = modifier.Value
// 			}
// 		}
// 		if target == "" {
// 			return nil, errors.New("no target provided")
// 		}
// 		if target != "Permanent" {
// 			return nil, fmt.Errorf("only Permanent target is supported: %s", target)
// 		}
// 		effect.Apply = func(state *Gamecore.State, player *core.Player) error {
// 			cards := player.Battlefield.GetAll()
// 			if len(cards) == 0 {
// 				return errors.New("no available targets")
// 			}
// 			choices := CreateChoices(cards, ZoneBattlefield)
// 			chosen, err := player.Agent.ChooseOne(
// 				fmt.Sprintf("Choose a card to tap or untap"),
// 				source,
// 				choices,
// 			)
// 			if err != nil {
// 				return fmt.Errorf("failed to choose card: %w", err)
// 			}
// 			permanent, err := player.Battlefield.Get(chosen.ID)
// 			if err != nil {
// 				return fmt.Errorf("failed to get permanent: %w", err)
// 			}
// 			p, ok := permanent.(*Permanent)
// 			if !ok {
// 				return fmt.Errorf("object is not a permanent: %s", chosen.ID)
// 			}
// 			var action string
// 			if p.IsTapped() {
// 				action = "untap"
// 				accept, err := player.Agent.Confirm(
// 					fmt.Sprintf("Do you want to untap %s?", p.Name()),
// 					source,
// 				)
// 				if err != nil {
// 					return fmt.Errorf("failed to confirm untap: %w", err)
// 				}
// 				if accept {
// 					p.Untap()
// 				}
// 			} else {
// 				action = "tap"
// 				accept, err := player.Agent.Confirm(
// 					fmt.Sprintf("Do you want to tap %s?", p.Name()),
// 					source,
// 				)
// 				if err != nil {
// 					return fmt.Errorf("failed to confirm tap: %w", err)
// 				}
// 				if accept {
// 					if err := p.Tap(); err != nil {
// 						if errors.Is(err, ErrAlreadyTapped) {
// 							return nil // Not an error to tap a card that's already tapped.
// 						}
// 						return fmt.Errorf("failed to tap card: %w", err)
// 					}
// 				}
// 			}
// 			state.Log(fmt.Sprintf("%s %s", action, p.Name()))
// 			return nil
// 		}
// 		effect.Description = fmt.Sprintf("tap or untap a card of type %s", target)
// 		effect.Tags = []EffectTag{{Key: "TapOrUntap", Value: target}}
// 	*/
// 	// return &effect, nil
// 	return nil, nil
// }

// func Scry(state core.State, source query.Object, n int, player core.Player) error {
// 	/*
// 		var taken []query.Object
// 		for range n {
// 			card, err := player.Library.TakeTop()
// 			if err != nil {
// 				// Not an error to scry on an empty library
// 				if errors.Is(err, ErrLibraryEmpty) {
// 					break
// 				}
// 				return fmt.Errorf("failed to take top card: %w", err)
// 			}
// 			taken = append(taken, card)
// 		}
// 		used := map[string]bool{}

// 		for range len(taken) {
// 			// Build option list from unplaced cards
// 			var choices []Choice
// 			for _, card := range taken {
// 				if !used[card.ID()] {
// 					choices = append(choices, Choice{
// 						Name: card.Name(),
// 						ID:   card.ID(),
// 					})
// 				}
// 			}
// 			chosen, err := player.Agent.ChooseOne(
// 				"Choose a card to place",
// 				source,
// 				choices,
// 			)
// 			if err != nil {
// 				return fmt.Errorf("failed to choose card: %w", err)
// 			}
// 			var chosenCard query.Object
// 			for _, card := range taken {
// 				if card.ID() == chosen.ID {
// 					chosenCard = card
// 					break
// 				}
// 			}
// 			if chosenCard == nil {
// 				return fmt.Errorf("failed to find chosen card: %s", chosen.ID)
// 			}
// 			used[chosen.ID] = true
// 			// TODO: Maybe have a global set of constants for choices like this
// 			const ChoiceTop = "Top"
// 			const ChoiceBottom = "Bottom"
// 			topBottomchoices := []Choice{
// 				{
// 					Name: ChoiceTop,
// 					ID:   ChoiceTop,
// 				},
// 				{
// 					Name: ChoiceBottom,
// 					ID:   ChoiceBottom,
// 				},
// 			}
// 			placement, err := player.Agent.ChooseOne(
// 				fmt.Sprintf("Place %s on top or bottom of your library?", chosenCard.Name()),
// 				source,
// 				topBottomchoices,
// 			)
// 			if err != nil {
// 				return fmt.Errorf("failed to choose placement: %w", err)
// 			}
// 			if placement.ID == ChoiceTop {
// 				player.Library.AddTop(chosenCard)
// 			} else {
// 				player.Library.Add(chosenCard)
// 			}
// 		}
// 	*/
// 	return nil
// }

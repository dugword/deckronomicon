package game

import (
	"errors"
	"fmt"
	"strconv"
)

// Effect represents an effect that can be applied to a game state.
type Effect struct {
	Apply       func(*GameState, *Player) error
	Description string
	Tags        []EffectTag
}

// EffectTag represents a tag associated with an effect. These are used to
// define the effect and its properties. E.g. a Draw effect will have a tag
// Key of "Count" and a Value of the number of cards to draw.
type EffectTag struct {
	Key   string
	Value string
}

// TODO  would be good to ensure all BuildEffect functions return the same
// type EffectBuilder func(source GameObject, effectModifiers []EffectModifier) (*Effect, error)

// BuildEffect creates an effect based on the provided EffectSpec.
func BuildEffect(source GameObject, spec EffectSpec) (*Effect, error) {
	switch spec.ID {
	case "AdditionalMana":
		return BuildEffectAdditionalMana(source, spec.Modifiers)
	case "AddMana":
		return BuildEffectAddMana(source, spec.Modifiers)
	case "Discard":
		return BuildEffectDiscard(source, spec.Modifiers)
	case "Draw":
		return BuildEffectDraw(source, spec.Modifiers)
	case "PutBackOnTop":
		return BuildEffectPutBackOnTop(source, spec.Modifiers)
	case "Mill":
		return BuildEffectMill(source, spec.Modifiers)
	case "Scry":
		return BuildEffectScry(source, spec.Modifiers)
	case "Search":
		return BuildEffectSearch(source, spec.Modifiers)
	case "SuffleFromGraveyard":
		return BuildEffectShuffleFromGraveyard(source, spec.Modifiers)
	case "Replicate":
		return BuildEffectReplicate(source, spec.Modifiers)
	case "Tap":
		return BuildEffectTap(source, spec.Modifiers)
	default:
		return &Effect{
			Description: fmt.Sprintf("unknown effect: %s", spec.ID),
			Tags:        []EffectTag{{Key: "Unknown", Value: spec.ID}},
			Apply: func(state *GameState, player *Player) error {
				return nil
			},
		}, nil
	}
}

// BuildEffectDraw creates a draw effect based on the provided modifiers.
// Keys: Count, Type
// Default: Count: 1
func BuildEffectDraw(source GameObject, modifiers []EffectModifier) (*Effect, error) {
	effect := Effect{}
	count := "1"
	var drawType string
	for _, modifier := range modifiers {
		if modifier.Key == "Count" {
			count = modifier.Value
		}
		if modifier.Key == "Type" {
			drawType = modifier.Value
		}
	}
	n, err := strconv.Atoi(count)
	if err != nil {
		return nil, fmt.Errorf("invalid count: %s", count)
	}
	tags := []EffectTag{{Key: "Draw", Value: count}}
	if drawType != "" {
		tags = append(tags, EffectTag{Key: "Type", Value: drawType})
	}
	effect.Description = fmt.Sprintf("draw %d cards", count)
	effect.Tags = tags
	effect.Apply = func(state *GameState, player *Player) error {
		if err := state.Draw(n, player); err != nil {
			fmt.Errorf("failed to draw %d cards: %w", n, err)
		}
		return nil
	}
	return &effect, nil
}

// BuildEffectAddMana creates an effect that adds mana to the player's mana
// pool.
// Supported Modifier Keys (concats multiple modifiers):
//   - Mana: <ManaString>
func BuildEffectAddMana(source GameObject, modifiers []EffectModifier) (*Effect, error) {
	effect := Effect{}
	var mana string
	for _, modifier := range modifiers {
		if modifier.Key == "Mana" {
			mana += modifier.Value
		}
	}
	if mana == "" {
		return nil, errors.New("no mana string provided")
	}
	if !isMana(mana) {
		return nil, fmt.Errorf("invalid mana string: %s", mana)
	}
	var tags []EffectTag
	for _, symbol := range ManaStringToManaSymbols(mana) {
		tags = append(tags, EffectTag{Key: AbilityTagManaAbility, Value: symbol})
	}
	effect.Description = fmt.Sprintf("add %s", mana)
	effect.Tags = tags
	effect.Apply = func(state *GameState, player *Player) error {
		if err := player.ManaPool.AddMana(mana); err != nil {
			return err
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
func BuildEffectAdditionalMana(source GameObject, modifiers []EffectModifier) (*Effect, error) {
	var mana string
	var target string
	var duration string
	for _, modifier := range modifiers {
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
	if !isMana(mana) {
		return nil, fmt.Errorf("invalid mana string: %s", mana)
	}
	if target == "" {
		return nil, errors.New("no target provided")
	}
	if duration != "EndOfTurn" {
		// return nil, errors.New("no duration provided")
		return nil, errors.New("only EndOfTurn duration is supported")
	}
	subtype, err := StringToSubtype(target)
	if err != nil {
		return nil, fmt.Errorf("invalid target subtype: %s", target)
	}
	effect := Effect{}
	id := getNextEventID()
	eventHandler := EventHandler{
		ID: id,
		Callback: func(event Event, state *GameState, player *Player) {
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
	var tags []EffectTag
	for _, symbol := range ManaStringToManaSymbols(mana) {
		tags = append(tags, EffectTag{Key: AbilityTagManaAbility, Value: symbol})
	}
	effect.Tags = tags
	effect.Description = fmt.Sprintf("add additional %s when you tap an %s for mana", mana, target)
	effect.Apply = func(state *GameState, player *Player) error {
		state.RegisterListenerUntil(
			eventHandler,
			EventEndStep,
		)
		return nil
	}
	return &effect, nil
}

// BuildEffectMill creates an effect that mills cards from the top of the
// library.
// Supported Modifier Keys (last applies):
//   - Count: <Cards to mill> Default: 1
//   - Target <target> Player | Self | Opponent
func BuildEffectMill(source GameObject, effectModifiers []EffectModifier) (*Effect, error) {
	effect := Effect{}
	count := "1"
	var target string
	for _, modifier := range effectModifiers {
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
	effect.Apply = func(state *GameState, player *Player) error {
		// Get the target player
		choices := []Choice{}
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
		return nil
	}
	effect.Description = fmt.Sprintf("mill %d cards from your library", n)
	effect.Tags = []EffectTag{{Key: "Mill", Value: count}}
	return &effect, nil
}

// BuildEffectPutBackOnTop creates an effect that puts cards back on top of
// the library.
// Supported Modifier Keys (last applies):
//   - Count: <Cards to put back> Default: 1
func BuildEffectPutBackOnTop(source GameObject, effectModifiers []EffectModifier) (*Effect, error) {
	effect := Effect{}
	count := "1"
	for _, modifier := range effectModifiers {
		if modifier.Key != "Count" {
			continue
		}
		count = modifier.Value
	}
	n, err := strconv.Atoi(count)
	if err != nil {
		return nil, fmt.Errorf("invalid count: %s", count)
	}
	effect.Apply = func(state *GameState, player *Player) error {
		if err := PutNBackOnTop(state, 2, source, player); err != nil {
			return err
		}
		return nil
	}
	effect.Description = fmt.Sprintf("put %d cards from your hand on top of your library in any order", n)
	effect.Tags = []EffectTag{{Key: "PutBackOnTop", Value: count}}
	return &effect, nil
}

// BuildEffectScry creates an effect that allows the player to scry.
// Supported Modifier Keys (last applies):
//   - Count: <Cards to scry> Default: 1
func BuildEffectScry(source GameObject, effectModifiers []EffectModifier) (*Effect, error) {
	effect := Effect{}
	count := "1"
	for _, modifier := range effectModifiers {
		if modifier.Key != "Count" {
			continue
		}
		count = modifier.Value
	}
	n, err := strconv.Atoi(count)
	if err != nil {
		return nil, fmt.Errorf("invalid count: %s", count)
	}
	effect.Description = fmt.Sprintf("look at the top %d cards of your library, then put them back on top or bottom of your library in any order.", n)
	effect.Apply = func(state *GameState, player *Player) error {
		if err := Scry(state, source, n, player); err != nil {
			return err
		}
		return nil
	}
	effect.Tags = []EffectTag{{Key: "Scry", Value: count}}
	return &effect, nil
}

// BuildEffectDiscard creates an effect that discards cards from the
// player's hand.
// Supported Modifier Keys (last applies):
//   - Count: <Cards to discard> Default: 1
//   - Delay: <Delay until> EndStep
func BuildEffectDiscard(source GameObject, effectModifiers []EffectModifier) (*Effect, error) {
	effect := Effect{}
	count := "1"
	var delay string
	for _, modifier := range effectModifiers {
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
	var effectFunc func(state *GameState, player *Player) error

	fn := func(state *GameState, player *Player) error {
		if err := state.Discard(n, source, player); err != nil {
			return err
		}
		return nil
	}

	effectFunc = fn

	if delay == "EndStep" {
		var eventHandler EventHandler
		id := getNextEventID()
		eventHandler = EventHandler{
			ID: id,
			Callback: func(event Event, state *GameState, player *Player) {
				if event.Type != EventEndStep {
					return
				}
				// TODO Handle errors some how...
				_ = fn(state, player)
				state.DeregisterListener(id)
				return
			}}
		effectFunc = func(state *GameState, player *Player) error {
			state.RegisterListener(eventHandler)
			return nil
		}
	}
	tags := []EffectTag{{Key: "Discard", Value: count}}
	if delay != "" {
		tags = append(tags, EffectTag{Key: "Delay", Value: delay})
	}
	description := fmt.Sprintf("discard %d cards", n)
	if delay != "EndStep" {
		description += " at the beginning of your next end step"
	}
	effect.Tags = tags
	effect.Description = description
	effect.Apply = effectFunc
	return &effect, nil
}

// BuildEffectSearch creates an effect that search cards from the library.
// Supported Modifier Keys (last applies):
//   - Subtype <subtype>
func BuildEffectSearch(source GameObject, effectModifiers []EffectModifier) (*Effect, error) {
	effect := Effect{}
	var subtype string
	for _, modifier := range effectModifiers {
		if modifier.Key == "Subtype" {
			subtype = modifier.Value
		}
	}
	if subtype == "" {
		return nil, errors.New("no subtype provided")
	}
	subtypeEnum, err := StringToSubtype(subtype)
	if err != nil {
		return nil, fmt.Errorf("invalid subtype: %s", subtype)
	}
	effect.Apply = func(state *GameState, player *Player) error {
		objects := player.Library.FindAllBySubtype(subtypeEnum)
		if len(objects) == 0 {
			return fmt.Errorf("no cards of subtype %s found", subtype)
		}
		choices := CreateObjectChoices(objects, ZoneLibrary)
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
	effect.Tags = []EffectTag{{Key: "Tutor", Value: subtype}}
	return &effect, nil
}

// BuildEffectShuffleFromGraveyard creates an effect that shuffles cards from
// the graveyard back into the library.
// Supported Modifier Keys (last applies):
//   - Count: <Cards to shuffle> Default: 1
func BuildEffectShuffleFromGraveyard(source GameObject, effectModifiers []EffectModifier) (*Effect, error) {
	effect := Effect{}
	count := "1"
	for _, modifier := range effectModifiers {
		if modifier.Key == "Count" {
			count = modifier.Value
		}
	}
	n, err := strconv.Atoi(count)
	if err != nil {
		return nil, fmt.Errorf("invalid count: %s", count)
	}
	effect.Apply = func(state *GameState, player *Player) error {
		for range n {
			choices := CreateObjectChoices(player.Graveyard.GetAll(), ZoneGraveyard)
			if len(choices) == 0 {
				break
			}
			chosen, err := player.Agent.ChooseOne("Choose cards to shuffle into your library", source, choices)
			if err != nil {
				return fmt.Errorf("failed to choose cards: %w", err)
			}
			card, err := player.Graveyard.Get(chosen.ID)
			if err != nil {
				return fmt.Errorf("failed to get card from graveyard: %w", err)
			}
			player.Library.Add(card)
			player.Library.Shuffle()
		}
		return nil
	}
	effect.Description = fmt.Sprintf("shuffle %d cards from your graveyard into your library", n)
	effect.Tags = []EffectTag{{Key: "ShuffleFromGraveyard", Value: count}}
	return &effect, nil
}

// BuildEffectReplicate creates an effect that replicates a spell for each
// tiem the replicate cost is paid.
// Supported Modifier Keys (last applies):
//   - Cost: <cost>
func BuildEffectReplicate(source GameObject, effectModifiers []EffectModifier) (*Effect, error) {
	var costString string
	for _, modifier := range effectModifiers {
		if modifier.Key == "Cost" {
			costString += modifier.Value
		}
	}
	card, ok := source.(*Card)
	if !ok {
		return nil, fmt.Errorf("source is not a card: %s", source.ID())
	}
	effect := Effect{}
	effect.Apply = func(state *GameState, player *Player) error {
		spell, err := NewSpell(card)
		if err != nil {
			return fmt.Errorf("failed to create spell: %w", err)
		}
		state.Stack.Add(spell)
		return nil
	}
	effect.Description = fmt.Sprintf("replicate for %s", costString)
	effect.Tags = []EffectTag{{Key: "Replicate", Value: costString}}
	return &effect, nil
}

// BuildEffectTap creates an effect that taps a card.
// Supported Modifier Keys (last applies):
//   - Target: Permanent
//
// TODO: Support other targets
func BuildEffectTap(source GameObject, effectModifiers []EffectModifier) (*Effect, error) {
	effect := Effect{}
	var target string
	for _, modifier := range effectModifiers {
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
	effect.Apply = func(state *GameState, player *Player) error {
		cards := player.Battlefield.GetAll()
		if len(cards) == 0 {
			// TODO: Spells can't be cast without targets
			return errors.New("no available targets")
		}
		choices := CreateObjectChoices(cards, ZoneBattlefield)
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
	return &effect, nil
}

func Scry(state *GameState, source GameObject, n int, player *Player) error {
	var taken []GameObject
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
		var chosenCard GameObject
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
	return nil
}

package game

import (
	"errors"
	"fmt"
	"strconv"
)

// Effect represents an effect that can be applied to a game state.
type Effect struct {
	Apply       func(*GameState, ChoiceResolver) error
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
	case "Scry":
		return BuildEffectScry(source, spec.Modifiers)
	case "Search":
		return BuildEffectSearch(source, spec.Modifiers)
	default:
		return &Effect{
			Description: fmt.Sprintf("unknown effect: %s", spec.ID),
			Tags:        []EffectTag{{Key: "Unknown", Value: spec.ID}},
			Apply: func(state *GameState, resolver ChoiceResolver) error {
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
	effect.Apply = func(state *GameState, resolver ChoiceResolver) error {
		if err := state.Draw(n); err != nil {
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
	effect.Apply = func(state *GameState, resolver ChoiceResolver) error {
		if err := state.ManaPool.AddMana(mana); err != nil {
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
		Callback: func(event Event, state *GameState, resolver ChoiceResolver) {
			// Move this into the register so I don't have to check for
			// it.
			if event.Type != EventTapForMana {
				return
			}
			if !event.Source.HasSubtype(subtype) {
				return
			}
			if err := state.ManaPool.AddMana(mana); err != nil {
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
	effect.Apply = func(state *GameState, resolver ChoiceResolver) error {
		state.RegisterListenerUntil(
			eventHandler,
			EventEndStep,
		)
		return nil
	}
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
	effect.Apply = func(state *GameState, resolver ChoiceResolver) error {
		if err := PutNBackOnTop(state, 2, source, resolver); err != nil {
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
	effect.Apply = func(state *GameState, resolver ChoiceResolver) error {
		if err := Scry(state, source, n, resolver); err != nil {
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
	var effectFunc func(state *GameState, resolver ChoiceResolver) error

	fn := func(state *GameState, resolver ChoiceResolver) error {
		if err := state.Discard(n, source, resolver); err != nil {
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
			Callback: func(event Event, state *GameState, resolver ChoiceResolver) {
				if event.Type != EventEndStep {
					return
				}
				// TODO Handle errors some how...
				_ = fn(state, resolver)
				state.DeregisterListener(id)
				return
			}}
		effectFunc = func(state *GameState, resolver ChoiceResolver) error {
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
	effect.Apply = func(state *GameState, resolver ChoiceResolver) error {
		choices := state.Library.ChooseCardsBySubtype(subtypeEnum)
		if len(choices) == 0 {
			return fmt.Errorf("no cards of subtype %s found", subtype)
		}
		chosen, err := resolver.ChooseOne(
			fmt.Sprintf("Choose a card to put into your hand"),
			source,
			choices,
		)
		if err != nil {
			return fmt.Errorf("failed to choose card: %w", err)
		}
		card, err := state.Library.TakeCardByIndex(chosen.Index)
		if err != nil {
			return fmt.Errorf("failed to take card: %w", err)
		}
		state.Library.Shuffle()
		state.Hand.Add(card)
		return nil
	}
	effect.Description = fmt.Sprintf("search library for a card of subtype %s", subtype)
	effect.Tags = []EffectTag{{Key: "Tutor", Value: subtype}}
	return &effect, nil
}

func Scry(state *GameState, source GameObject, n int, resolver ChoiceResolver) error {
	taken, _ := state.Library.TakeCards(n)
	used := make([]bool, len(taken))
	for range len(taken) {
		// Build option list from unplaced cards
		var choices []Choice
		for index, card := range taken {
			if !used[index] {
				choices = append(choices, Choice{
					Name:  card.Name(),
					Index: index,
				})
			}
		}
		chosen, err := resolver.ChooseOne(
			"Choose a card to place",
			source,
			choices,
		)
		if err != nil {
			return fmt.Errorf("failed to choose card: %w", err)
		}
		chosenCard := taken[chosen.Index]
		used[chosen.Index] = true
		topBottomchoices := []Choice{
			{Name: "Top"},
			{Name: "Bottom"},
		}
		placement, err := resolver.ChooseOne(
			fmt.Sprintf("Place %s on top or bottom of your library?", chosenCard.Name()),
			source,
			topBottomchoices,
		)
		if err != nil {
			return fmt.Errorf("failed to choose placement: %w", err)
		}
		if placement.Index == 0 {
			state.Library.PutOnTop(chosenCard)
		} else {
			state.Library.PutOnBottom(chosenCard)
		}
	}
	return nil
}

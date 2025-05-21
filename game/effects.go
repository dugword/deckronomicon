package game

import (
	"fmt"
	"strconv"
)

// Effect represents an effect that can be applied to a game state.
type Effect struct {
	Description string
	Apply       func(*GameState, ChoiceResolver) error
	Tags        []AbilityTag
}

type EffectBuilder func(source string, effectModifiers []EffectModifier) (*Effect, error)

// EffectMap is a map of effect keys to their corresponding Effect instances.
// TODO: Maybe have this be a function and wrap things in a consistent way for
// delayed triggers and tags and stuff.
var EffectMap = map[string]EffectBuilder{
	"AddMana": func(source string, effectModifiers []EffectModifier) (*Effect, error) {
		/* Supported Modifier Keys (concats multiple modifiers):
		- Mana: <ManaString>
		*/
		var mana string
		for _, modifier := range effectModifiers {
			if modifier.Key != "Mana" {
				continue
			}
			mana += modifier.Value
		}
		if !isMana(mana) {
			return nil, fmt.Errorf("invalid mana string: %s", mana)
		}
		if mana == "" {
			panic("no mana string provided")

		}
		return &Effect{
			Description: fmt.Sprintf("add %s", mana),
			Apply: func(state *GameState, resolver ChoiceResolver) error {
				if err := state.ManaPool.AddMana(mana); err != nil {
					return err
				}
				return nil
			},
			// TODO: Split the mana string into multiple tags
			Tags: []AbilityTag{{Key: "ManaSource", Value: mana}},
		}, nil
	},
	"AdditionalMana": func(source string, effectModifiers []EffectModifier) (*Effect, error) {
		id := getNextEventID()
		eventHandler := EventHandler{
			ID: id,
			Callback: func(event Event, state *GameState, resolver ChoiceResolver) {
				// Move this into the register so I don't have to check for
				// it.
				fmt.Println("AdditionalMana event handler called")
				if event.Type != EventTapForMana {
					return
				}
				fmt.Println("AdditionalMana event handler called 2")
				fmt.Printf("%+v\n", event.Source.Name())
				fmt.Printf("%+v\n", event.Source.Subtypes())
				fmt.Printf("%+v\n", event.Source.Card())
				fmt.Printf("%+v\n", event.Source.Card().object)
				if !event.Source.Card().HasSubtype(SubtypeIsland) {
					return
				}
				fmt.Println("AdditionalMana event handler called 3")
				state.ManaPool.Add(ColorBlue, 1)
				fmt.Println("AdditionalMana event handler called 4")
				return
			}}
		return &Effect{
			Description: fmt.Sprintf("add additional {U} when you tap an island for mana"),
			Apply: func(state *GameState, resolver ChoiceResolver) error {
				state.RegisterListenerUntil(
					eventHandler,
					EventEndStep,
				)
				return nil
			},
			// TODO: Don't hard code U
			Tags: []AbilityTag{{Key: "AdditionalMana", Value: "{U}"}},
		}, nil
	},
	"Draw": func(source string, effectModifiers []EffectModifier) (*Effect, error) {
		/* Supported Modifier Keys (last applies):
		- Count: <Cards to draw> Default 1
		*/
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
		return &Effect{
			Description: fmt.Sprintf("draw %d cards", n),
			Apply: func(state *GameState, resolver ChoiceResolver) error {
				if err := state.Draw(n); err != nil {
					return err
				}
				return nil
			},
			Tags: []AbilityTag{{Key: "Draw", Value: count}},
		}, nil
	},
	"PutBackOnTop": func(source string, effectModifiers []EffectModifier) (*Effect, error) {
		/* Supported Modifier Keys (last applies):
		- Count: <Cards to draw> Default: 1
		*/
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
		return &Effect{
			Description: fmt.Sprintf("put %d cards from your hand on top of your library in any order", n),
			Apply: func(state *GameState, resolver ChoiceResolver) error {
				if err := PutNBackOnTop(state, 2, source, resolver); err != nil {
					return err
				}
				return nil
			},
			Tags: []AbilityTag{{Key: "PutBackOnTop", Value: "2"}},
		}, nil
	},
	"Scry": func(source string, effectModifiers []EffectModifier) (*Effect, error) {
		/* Supported Modifier Keys (last applies):
		- Count: <Cards to scry> Default 1
		*/
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
		return &Effect{
			Description: fmt.Sprintf("Look at the top %d cards of your library, then put them back on top or bottom of your library in any order.", n),
			Apply: func(state *GameState, resolver ChoiceResolver) error {
				if err := Scry(state, n, resolver); err != nil {
					return err
				}
				return nil
			},
			Tags: []AbilityTag{{Key: "Scry", Value: count}},
		}, nil
	},
	"Discard": func(source string, effectModifiers []EffectModifier) (*Effect, error) {
		/* Supported Modifier Keys (last applies):
		- Count: <Cards to discard> default 1
		- Delay: <Delay until> EndStep
		*/
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
					fmt.Println("The event was called")
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
		tags := []AbilityTag{{Key: "Discard", Value: count}}
		if delay != "" {
			tags = append(tags, AbilityTag{Key: "Delay", Value: delay})
		}
		description := fmt.Sprintf("discard %d cards", n)
		if delay != "EndStep" {
			description += " at the beginning of your next end step"
		}
		return &Effect{
			Description: description,
			// TODO: Standardize this
			//"at the beginning of your next end step.",
			Apply: effectFunc,
			Tags:  tags,
		}, nil
	},
}

func Scry(state *GameState, n int, resolver ChoiceResolver) error {
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
			"ScryChoseCard",
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
			"ScryPlaceCard",

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

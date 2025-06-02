package effect

import (
	"deckronomicon/packages/game/core"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"fmt"
	"strconv"
)

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
				card, err := player.TakeTopCard()
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
	effect.tags = []core.Tag{
		{Key: "Look", Value: lookCount},
		{Key: "Choose", Value: chooseCount},
		{Key: "Rest", Value: restZone},
		{Key: "Order", Value: order},
	}
	effect.description = fmt.Sprintf("look at the top %d cards of your library, choose %d to put into your hand, and put the rest on the %s of your library in %s order", nLook, nChoose, restZone, order)
	return &effect, nil
}

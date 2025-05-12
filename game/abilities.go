package game

// TODO: Implement mana abilities to not use the stack

// Ability is the general interface for all abilities.
type Ability interface {
	Resolve(game *GameState, choiceResolver ChoiceResolver) error
	Description() string
}

type AbilityTag struct {
	Key   string
	Value string
}

// ActivatedAbility represents abilities that require activation costs.
type ActivatedAbility struct {
	Cost          Cost
	Effect        func(*GameState, ChoiceResolver)
	Description   string
	IsManaAbility bool // maybe make this a generic set of tags
	Tags          []AbilityTag
}

// TriggeredAbility represents abilities that trigger on specific events.
type TriggeredAbility struct {
	TriggerCondition func(event Event) bool
	Effect           func(*GameState, ChoiceResolver)
	Description      string
}

// StaticAbility represents continuous effects.
type StaticAbility struct {
	Effect      func(*GameState, ChoiceResolver)
	Description string
}

// SpellAbility represents abilities on instant or sorcery spells.
type SpellAbility struct {
	Effect      func(*GameState, ChoiceResolver)
	Description string
}

func CreateClueTrigger() TriggeredAbility {
	return TriggeredAbility{
		TriggerCondition: func(event Event) bool {
			return event.Type == EventLeaveBattlefield &&
				event.Source != nil && event.Source.HasType("Artifact") && event.Source.HasSubtype("Clue")
		},
		Effect: func(gs *GameState, resolver ChoiceResolver) {
			// maybe draw a card, or create a new token
		},
		Description: "When a Clue is sacrificed, do X",
	}
}

/*
func Scry(n int) Ability {
	return func(card *Card, state *GameState) error {
		if len(state.Deck) < n {
			n = len(state.Deck)
		}
		// topCards
		_ = state.Deck[:n]

		// Naive: Keep all on top for now
		// Later: Add logic to reorder or bottom some

		state.Log = append(state.Log, fmt.Sprintf("%s scried %d cards", card.Name, n))
		return nil
	}
}

func FanaticalOfferingSacrifice() Ability {
	return func(card *Card, state *GameState) error {
		// Find a sacrificial artifact
		for i, target := range state.Battlefield {
			if target.Types["Artifact"] {
				state.Graveyard = append(state.Graveyard, target)
				state.Battlefield = append(state.Battlefield[:i], state.Battlefield[i+1:]...)
				state.DrawCards(2)
				state.Log = append(state.Log, fmt.Sprintf("%s was sacrificed for Fanatical Offering", target.Name))
				return nil
			}
		}
		return fmt.Errorf("no valid artifact to sacrifice for Fanatical Offering")
	}
}

func CandyTrailScry2() Ability {
	return func(card *Card, state *GameState) error {
		n := 2
		if len(state.Deck) < n {
			n = len(state.Deck)
		}
		// For now, just leave them on top — simple sim logic
		state.Log = append(state.Log, fmt.Sprintf("%s scries %d cards", card.Name, n))
		return nil
	}
}

func CandyTrailSacrificeDraw() Ability {
	return func(card *Card, state *GameState) error {
		// Remove from battlefield, add to graveyard
		for i, c := range state.Battlefield {
			if c == card {
				state.Battlefield = append(state.Battlefield[:i], state.Battlefield[i+1:]...)
				state.Graveyard = append(state.Graveyard, card)
				state.DrawCards(1)
				state.Log = append(state.Log, fmt.Sprintf("%s was sacrificed to draw a card", card.Name))
				return nil
			}
		}
		return fmt.Errorf("could not sacrifice %s", card.Name)
	}
}

func ConjurersBaubleEffect() Ability {
	return func(card *Card, state *GameState) error {
		if len(state.Graveyard) == 0 {
			return fmt.Errorf("graveyard is empty, bauble fizzles")
		}

		// Prioritize draw spells
		for i, target := range state.Graveyard {
			if strings.Contains(target.Text, "draw") {
				state.Deck = append(state.Deck, target) // bottom of deck
				state.Graveyard = append(state.Graveyard[:i], state.Graveyard[i+1:]...)
				state.Log = append(state.Log, fmt.Sprintf("Bauble put %s on bottom of deck", target.Name))
				break
			}
		}

		state.DrawCards(1)
		return nil
	}
}

func StreamOfThoughtShuffle() Ability {
	return func(card *Card, state *GameState) error {
		var toShuffle []*Card

		for i := 0; i < len(state.Graveyard) && len(toShuffle) < 4; i++ {
			c := state.Graveyard[i]
			if strings.Contains(c.Text, "draw") {
				toShuffle = append(toShuffle, c)
			}
		}

		if len(toShuffle) == 0 {
			return fmt.Errorf("nothing valuable to shuffle with Stream of Thought")
		}

		for _, c := range toShuffle {
			state.Deck = append(state.Deck, c)
			state.Graveyard = removeCard(state.Graveyard, c)
			state.Log = append(state.Log, fmt.Sprintf("Stream of Thought shuffled in %s", c.Name))
		}

		state.ShuffleDeck()
		return nil
	}
}
*/

package game

type Hand struct {
	cards []*Card
}

// TODO: for display remove later
func (h *Hand) Cards() []*Card {
	return h.cards
}

func (h *Hand) Add(cards ...*Card) {
	h.cards = append(h.cards, cards...)
}

func (h *Hand) CardChoices() []Choice {
	var choices []Choice
	for i, card := range h.cards {
		choices = append(choices, Choice{Name: card.Name(), Index: i})
	}
	return choices
}

func (h *Hand) FindCard(name string) *Card {
	for _, card := range h.cards {
		if card.Name() == name {
			return card
		}
	}
	return nil
}

// TODO: need to figure out how to handle duplicates so this can be better, maybe by an ID?
func (h *Hand) GetCard(i int) *Card {
	return h.cards[i]
}

// Does it matter if there are duplicates? Do I need to track ID? Probably someday, but not now
// TODO: Should this return an error if the card is not found
func (h *Hand) RemoveCard(target *Card) {
	newHand := h.cards
	for i, card := range h.cards {
		if card == target {
			newHand = append(newHand[:i], newHand[i+1:]...)
			break
		}
	}
	h.cards = newHand
}

func (h *Hand) Size() int {
	return len(h.cards)
}

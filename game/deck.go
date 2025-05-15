package game

type Deck struct {
	cards []*Card
}

func (d *Deck) Shuffle() {
	shuffleCards(d.cards)
}

// TODO Maybe return error?
func (d *Deck) DrawCards(n int) ([]*Card, error) {
	if len(d.cards) < n {
		return nil, PlayerLostError{
			Reason: DeckedOut,
		}
	}
	drawn, remaining := d.cards[:n], d.cards[n:]
	d.cards = remaining
	return drawn, nil
}

func (d *Deck) Size() int {
	return len(d.cards)
}

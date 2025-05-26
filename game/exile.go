package game

type Exile struct {
	cards []*Card
}

func NewExile() *Exile {
	return &Exile{
		cards: []*Card{},
	}
}

func (e *Exile) Add(object GameObject) error {
	return nil
}

// This probably makes more sense as a method of Player
func (e *Exile) AvailableActivatedAbilities(*GameState, PlayerAgent) []*ActivatedAbility {
	return nil
}

// This probably makes more sense as a method of Player
func (e *Exile) AvailableToPlay(*GameState, PlayerAgent) []GameObject {
	return nil
}
func (e *Exile) Find(id string) (GameObject, error) {
	return nil, nil
}
func (e *Exile) Get(id string) (GameObject, error) {
	return nil, nil
}
func (e *Exile) GetAll() []GameObject {
	return nil
}
func (e *Exile) Remove(id string) error {
	return nil
}
func (e *Exile) Take(id string) (GameObject, error) {
	return nil, nil
}
func (e *Exile) Size() int {
	return 0
}
func (e *Exile) ZoneType() string {
	return ZoneExile
}

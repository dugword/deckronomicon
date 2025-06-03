package gob

type Permanent struct {
	Card               *Card // The card representing the permanent
	Name               string
	ActivatedAbilities []*Ability // List of abilities that can be activated
}

package gob

// gob - (G)ame (Ob)jects is a package that provides structures and methods
// for creating and managing game objects in Magic the Gathering.

type Card struct {
	Name               string
	ActivatedAbilities []*Ability // List of abilities that can be activated
}

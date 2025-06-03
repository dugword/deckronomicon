package engine

type PlayerAgent interface {
	// TODO Will be a complex type in the future, string works for now
	GetNextAction() string
}

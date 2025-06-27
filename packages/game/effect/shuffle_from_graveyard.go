package effect

type ShuffleFromGraveyard struct {
	Count int
}

func NewShuffleFromGraveyard(modifiers map[string]any) (ShuffleFromGraveyard, error) {
	countModifier, err := parseCount(modifiers)
	if err != nil {
		return ShuffleFromGraveyard{}, err
	}
	return ShuffleFromGraveyard{
		Count: countModifier,
	}, nil
}

func (e ShuffleFromGraveyard) Name() string {
	return "ShuffleFromGraveyard"
}

func (e ShuffleFromGraveyard) TargetSpec() TargetSpec {
	return NoneTargetSpec{}
}

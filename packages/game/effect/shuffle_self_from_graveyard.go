package effect

import "deckronomicon/packages/game/target"

type ShuffleSelfFromGraveyard struct {
}

func NewShuffleSelfFromGraveyard(_ map[string]any) (*ShuffleSelfFromGraveyard, error) {
	return &ShuffleSelfFromGraveyard{}, nil
}

func (e *ShuffleSelfFromGraveyard) Name() string {
	return "ShuffleSelfFromGraveyard"
}

func (e *ShuffleSelfFromGraveyard) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}

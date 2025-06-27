package effect

type PutBackOnTop struct {
	Count int
}

func NewPutBackOnTop(modifiers map[string]any) (PutBackOnTop, error) {
	countModifier, err := parseCount(modifiers)
	if err != nil {
		return PutBackOnTop{}, err
	}
	return PutBackOnTop{
		Count: countModifier,
	}, nil
}

func (e PutBackOnTop) Name() string {
	return "PutBackOnTop"
}

func (e PutBackOnTop) TargetSpec() TargetSpec {
	return NoneTargetSpec{}
}

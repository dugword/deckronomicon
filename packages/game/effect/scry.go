package effect

type Scry struct {
	Count int `json:"Count"`
}

func NewScry(modifiers map[string]any) (Scry, error) {
	countModifier, err := parseCount(modifiers)
	if err != nil {
		return Scry{}, err
	}
	return Scry{
		Count: countModifier,
	}, nil
}

func (e Scry) Name() string {
	return "Scry"
}

func (e Scry) TargetSpec() TargetSpec {
	return NoneTargetSpec{}
}

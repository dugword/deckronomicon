package choose

type ChoiceResults interface {
	ChoiceType() ChoiceType
}

type ChooseOneResults struct {
	Choice Choice
}

func (r ChooseOneResults) ChoiceType() ChoiceType {
	return ChoiceTypeChooseOne
}

type ChooseManyResults struct {
	Choices []Choice
}

func (r ChooseManyResults) ChoiceType() ChoiceType {
	return ChoiceTypeChooseMany
}

type MapChoicesToBucketsResults struct {
	Assignments map[Bucket][]Choice
}

func (r MapChoicesToBucketsResults) ChoiceType() ChoiceType {
	return ChoiceTypeMapChoicesToBuckets
}

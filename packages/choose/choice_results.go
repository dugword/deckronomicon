package choose

// TODO: I really should just be returning IDs and not trusting player agents to return the correct
// objects. This could lead to inaccuracies where players can choose cards that they shouldn't be able to choose,
// or choose cards that are not in the correct zone.

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

type ChooseNumberResults struct {
	Number int
}

func (r ChooseNumberResults) ChoiceType() ChoiceType {
	return ChoiceTypeChooseNumber
}

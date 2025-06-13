package choose

const (
	ChoiceTypeChooseOne           ChoiceType = "ChooseOne"
	ChoiceTypeChooseMany          ChoiceType = "ChooseMany"
	ChoiceTypeMapChoicesToBuckets ChoiceType = "MapChoicesToBuckets"
)

type ChoiceOpts interface {
	ChoiceType() ChoiceType
}

type ChooseOneOpts struct {
	Choices  []Choice
	Optional bool
}

func (o ChooseOneOpts) ChoiceType() ChoiceType {
	return ChoiceTypeChooseOne
}

type ChooseManyOpts struct {
	Choices  []Choice
	Min      int
	Max      int
	Optional bool
}

func (o ChooseManyOpts) ChoiceType() ChoiceType {
	return ChoiceTypeChooseMany
}

type MapChoicesToBucketsOpts struct{}

func (o MapChoicesToBucketsOpts) ChoiceType() ChoiceType {
	return ChoiceTypeMapChoicesToBuckets
}

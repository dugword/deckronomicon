package choose

// TODO: This should probably live under agent

const (
	ChoiceTypeChooseOne           ChoiceType = "ChooseOne"
	ChoiceTypeChooseMany          ChoiceType = "ChooseMany"
	ChoiceTypeMapChoicesToBuckets ChoiceType = "MapChoicesToBuckets"
	ChoiceTypeChooseNumber        ChoiceType = "ChooseNumber"
)

const (
	BucketTop    Bucket = "Top"
	BucketBottom Bucket = "Bottom"
)

type ChoiceOpts interface {
	ChoiceType() ChoiceType
}

type ChooseOneOpts struct {
	Choices []Choice
}

func (o ChooseOneOpts) ChoiceType() ChoiceType {
	return ChoiceTypeChooseOne
}

type ChooseManyOpts struct {
	Choices []Choice
	Min     int
	Max     int
}

func (o ChooseManyOpts) ChoiceType() ChoiceType {
	return ChoiceTypeChooseMany
}

type Bucket string

type MapChoicesToBucketsOpts struct {
	Buckets []Bucket
	Choices []Choice
}

func (o MapChoicesToBucketsOpts) ChoiceType() ChoiceType {
	return ChoiceTypeMapChoicesToBuckets
}

type ChooseNumberOpts struct {
	Min int
	Max int
}

func (o ChooseNumberOpts) ChoiceType() ChoiceType {
	return ChoiceTypeChooseNumber
}

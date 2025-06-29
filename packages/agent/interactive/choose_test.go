package interactive

import (
	"bufio"
	"deckronomicon/packages/choose"
	"deckronomicon/packages/game/mtg"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TODO: Figure out how to disable the ">>" prompt for agent in tests.

var AllowAllUnexported = cmp.Exporter(func(reflect.Type) bool { return true })

type mockSource struct{}

func (m mockSource) Name() string {
	return "MockSource"
}

type mockChoice struct {
	name string
	id   string
}

func (m mockChoice) Name() string {
	return m.name
}

func (m mockChoice) ID() string {
	return m.id
}

func getMockChoices() []mockChoice {
	return []mockChoice{
		{"Ness", "Ness ID"},
		{"Paula", "Paula ID"},
		{"Jeff", "Jeff ID"},
		{"Poo", "Poo ID"},
	}
}

func TestChoose(t *testing.T) {
	testCases := []struct {
		name         string
		inputs       []string
		choicePrompt choose.ChoicePrompt
		want         choose.ChoiceResults
	}{
		{
			name:   "with ChooseOneOpts input 1",
			inputs: []string{"1"},
			choicePrompt: choose.ChoicePrompt{
				ChoiceOpts: choose.ChooseOneOpts{
					Choices: choose.NewChoices(getMockChoices()),
				},
			},
			want: choose.ChooseOneResults{
				Choice: getMockChoices()[0],
			},
		},
		{
			name:   "with ChooseOneOpts input 1 item 2",
			inputs: []string{"2"},
			choicePrompt: choose.ChoicePrompt{
				ChoiceOpts: choose.ChooseOneOpts{
					Choices: choose.NewChoices(getMockChoices()),
				},
			},
			want: choose.ChooseOneResults{
				Choice: getMockChoices()[1],
			},
		},
		{
			name:   "with ChooseOneOpts no input optional",
			inputs: []string{""},
			choicePrompt: choose.ChoicePrompt{
				Optional: true,
				ChoiceOpts: choose.ChooseOneOpts{
					Choices: choose.NewChoices(getMockChoices()),
				},
			},
			want: choose.ChooseOneResults{},
		},
		{
			name:   "with ChooseManyOpts min 1 max 1 input 1",
			inputs: []string{"1"},
			choicePrompt: choose.ChoicePrompt{
				ChoiceOpts: choose.ChooseManyOpts{
					Choices: choose.NewChoices(getMockChoices()),
					Min:     1,
					Max:     1,
				},
			},
			want: choose.ChooseManyResults{
				Choices: choose.NewChoices(getMockChoices()[:1]),
			},
		},
		{
			name:   "with ChooseManyOpts min 1 max 1 input 1 item 3",
			inputs: []string{"3"},
			choicePrompt: choose.ChoicePrompt{
				ChoiceOpts: choose.ChooseManyOpts{
					Choices: choose.NewChoices(getMockChoices()),
					Min:     1,
					Max:     1,
				},
			},
			want: choose.ChooseManyResults{
				Choices: choose.NewChoices(getMockChoices()[2:3]),
			},
		},
		{
			name:   "with ChooseManyOpts min 2 max 2 input 2",
			inputs: []string{"1", "1"},
			choicePrompt: choose.ChoicePrompt{
				ChoiceOpts: choose.ChooseManyOpts{
					Choices: choose.NewChoices(getMockChoices()),
					Min:     2,
					Max:     2,
				},
			},
			want: choose.ChooseManyResults{
				Choices: choose.NewChoices(getMockChoices()[:2]),
			},
		},
		{
			name:   "with ChooseManyOpts min 1 max 2 input 1",
			inputs: []string{"1", ""},
			choicePrompt: choose.ChoicePrompt{
				ChoiceOpts: choose.ChooseManyOpts{
					Choices: choose.NewChoices(getMockChoices()),
					Min:     1,
					Max:     2,
				},
			},
			want: choose.ChooseManyResults{
				Choices: choose.NewChoices(getMockChoices()[:1]),
			},
		},
		{
			name:   "with ChooseManyOpts min 1 max 2 input 2",
			inputs: []string{"1", "1"},
			choicePrompt: choose.ChoicePrompt{
				ChoiceOpts: choose.ChooseManyOpts{
					Choices: choose.NewChoices(getMockChoices()),
					Min:     1,
					Max:     2,
				},
			},
			want: choose.ChooseManyResults{
				Choices: choose.NewChoices(getMockChoices()[:2]),
			},
		},
		{
			name:   "with ChooseManyOpts min 0 max 1 no input",
			inputs: []string{""},
			choicePrompt: choose.ChoicePrompt{
				ChoiceOpts: choose.ChooseManyOpts{
					Choices: choose.NewChoices(getMockChoices()),
					Min:     0,
					Max:     1,
				},
			},
			want: choose.ChooseManyResults{},
		},
		{
			name:   "with ChooseManyOpts min 2 max 2 input 3",
			inputs: []string{"1", "1", "1"},
			choicePrompt: choose.ChoicePrompt{
				ChoiceOpts: choose.ChooseManyOpts{
					Choices: choose.NewChoices(getMockChoices()),
					Min:     2,
					Max:     2,
				},
			},
			want: choose.ChooseManyResults{
				Choices: choose.NewChoices(getMockChoices()[:2]),
			},
		},
		{
			name:   "with ChooseManyOpts min 2 max 2 no input optional",
			inputs: []string{""},
			choicePrompt: choose.ChoicePrompt{
				Optional: true,
				ChoiceOpts: choose.ChooseManyOpts{
					Choices: choose.NewChoices(getMockChoices()),
					Min:     2,
					Max:     2,
				},
			},
			want: choose.ChooseManyResults{},
		},
		{
			name:   "with MapChoicesToBucketsOpts 1 to first",
			inputs: []string{"1", ""},
			choicePrompt: choose.ChoicePrompt{
				ChoiceOpts: choose.MapChoicesToBucketsOpts{
					Buckets: []choose.Bucket{choose.BucketTop, choose.BucketBottom},
					Choices: choose.NewChoices(getMockChoices()),
				},
			},
			want: choose.MapChoicesToBucketsResults{
				Assignments: map[choose.Bucket][]choose.Choice{
					choose.BucketTop:    {getMockChoices()[0]},
					choose.BucketBottom: {getMockChoices()[1], getMockChoices()[2], getMockChoices()[3]},
				},
			},
		},
		{
			name:   "with MapChoicesToBucketsOpts assign 1 to first 1 to second",
			inputs: []string{"1", "1"},
			choicePrompt: choose.ChoicePrompt{
				ChoiceOpts: choose.MapChoicesToBucketsOpts{
					Buckets: []choose.Bucket{choose.BucketTop, choose.BucketBottom},
					Choices: choose.NewChoices(getMockChoices()),
				},
			},
			want: choose.MapChoicesToBucketsResults{
				Assignments: map[choose.Bucket][]choose.Choice{
					choose.BucketTop:    {getMockChoices()[0]},
					choose.BucketBottom: {getMockChoices()[1], getMockChoices()[2], getMockChoices()[3]},
				},
			},
		},
		{
			name:   "with MapChoicesToBucketsOpts no input",
			inputs: []string{"", ""},
			choicePrompt: choose.ChoicePrompt{
				ChoiceOpts: choose.MapChoicesToBucketsOpts{
					Buckets: []choose.Bucket{choose.BucketTop, choose.BucketBottom},
					Choices: choose.NewChoices(getMockChoices()),
				},
			},
			want: choose.MapChoicesToBucketsResults{
				Assignments: map[choose.Bucket][]choose.Choice{
					choose.BucketBottom: choose.NewChoices(getMockChoices()),
				},
			},
		},
		{
			name:   "with MapChoicesToBucketsOpts all to top",
			inputs: []string{"1 2 3 4"},
			choicePrompt: choose.ChoicePrompt{
				ChoiceOpts: choose.MapChoicesToBucketsOpts{
					Buckets: []choose.Bucket{choose.BucketTop, choose.BucketBottom},
					Choices: choose.NewChoices(getMockChoices()),
				},
			},
			want: choose.MapChoicesToBucketsResults{
				Assignments: map[choose.Bucket][]choose.Choice{
					choose.BucketTop: choose.NewChoices(getMockChoices()),
				},
			},
		},
		{
			name:   "with MapChoicesToBucketsOpts all to bottom",
			inputs: []string{"", "1 2 3 4"},
			choicePrompt: choose.ChoicePrompt{
				ChoiceOpts: choose.MapChoicesToBucketsOpts{
					Buckets: []choose.Bucket{choose.BucketTop, choose.BucketBottom},
					Choices: choose.NewChoices(getMockChoices()),
				},
			},
			want: choose.MapChoicesToBucketsResults{
				Assignments: map[choose.Bucket][]choose.Choice{
					choose.BucketBottom: choose.NewChoices(getMockChoices()),
				},
			},
		},
		{
			name:   "with MapChoicesToBucketsOpts all to top reverse order",
			inputs: []string{"4 3 2 1"},
			choicePrompt: choose.ChoicePrompt{
				ChoiceOpts: choose.MapChoicesToBucketsOpts{
					Buckets: []choose.Bucket{choose.BucketTop, choose.BucketBottom},
					Choices: choose.NewChoices(getMockChoices()),
				},
			},
			want: choose.MapChoicesToBucketsResults{
				Assignments: map[choose.Bucket][]choose.Choice{
					choose.BucketTop: {getMockChoices()[3], getMockChoices()[2], getMockChoices()[1], getMockChoices()[0]},
				},
			},
		},
		{
			name:   "with MapChoicesToBucketsOpts mix up order",
			inputs: []string{"3 1", "2 1"},
			choicePrompt: choose.ChoicePrompt{
				ChoiceOpts: choose.MapChoicesToBucketsOpts{
					Buckets: []choose.Bucket{choose.BucketTop, choose.BucketBottom},
					Choices: choose.NewChoices(getMockChoices()),
				},
			},
			want: choose.MapChoicesToBucketsResults{
				Assignments: map[choose.Bucket][]choose.Choice{
					choose.BucketTop:    {getMockChoices()[2], getMockChoices()[0]},
					choose.BucketBottom: {getMockChoices()[3], getMockChoices()[1]},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockScanner := bufio.NewScanner(
				strings.NewReader(strings.Join(tc.inputs, "\n") + "\n"),
			)
			agent := NewAgent(
				mockScanner,
				"Test Player",
				[]mtg.Step{},
				"./testdata/display.tmpl",
				false, // No autopay for these tests
				nil,   // No autopay for these tests
				false,
			)
			tc.choicePrompt.Source = mockSource{}
			got, err := agent.Choose(tc.choicePrompt)
			if err != nil {
				t.Fatalf("Choose(...); err = %v; want %v", err, nil)
			}
			if diff := cmp.Diff(tc.want, got, AllowAllUnexported); diff != "" {
				t.Errorf("Choose(...) mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

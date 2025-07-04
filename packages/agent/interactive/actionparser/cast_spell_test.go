package actionparser

import (
	"deckronomicon/packages/agent/dummy"
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/state"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func mockMaybeApplyEvent(game *state.Game, event event.GameEvent) (*state.Game, error) {
	// Mock implementation for testing purposes
	return game, nil
}

func TestParseCastSpellCommand(t *testing.T) {
	const playerID = "Test Player"
	// #region Success Cases
	testCases := []struct {
		name  string
		arg   string
		agent engine.PlayerAgent
		want  action.CastSpellRequest
	}{
		{
			name: "with name",
			arg:  "Test Card",
			want: action.CastSpellRequest{
				CardID: "Test Card ID",
			},
		},
		{
			name: "with id",
			arg:  "Test Card ID",
			want: action.CastSpellRequest{
				CardID: "Test Card ID",
			},
		},
		{
			name: "with choice",
			arg:  "",
			want: action.CastSpellRequest{
				CardID: "Test Card ID",
			},
		},
		{
			name: "with selected card replicated",
			arg:  "Card with Replicate",
			want: action.CastSpellRequest{
				CardID:         "Card with Replicate ID",
				ReplicateCount: 1,
			},
		},
		{
			name: "with Arcane card and spliced card",
			arg:  "Arcane Card",
			want: action.CastSpellRequest{
				CardID:        "Acane Card ID",
				SpliceCardIDs: []string{"Card with Splice ID"},
			},
		},
		{
			name: "with selected card requiring target",
			arg:  "Card with Target",
			agent: dummy.NewChooseProvided("Test Player", dummy.ChooseProvidedConfig{
				OneChoiceIDs: []string{"Test Permanent ID"},
			}),
			want: action.CastSpellRequest{
				CardID: "Card with Target ID",
				TargetsForEffects: map[effect.EffectTargetKey]target.Target{
					{SourceID: "Card with Target ID", EffectIndex: 0}: {
						Type: mtg.TargetTypePermanent,
						ID:   "Test Permanent ID",
					},
				},
			},
		},
		{
			name: "with selected card using Flashback",
			arg:  "Card with Flashback",
			agent: dummy.NewChooseProvided("Test Player", dummy.ChooseProvidedConfig{
				OneChoiceIDs: []string{"Card with Flashback ID"},
			}),
			want: action.CastSpellRequest{
				CardID:    "Card with Flashback ID",
				Flashback: true,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game := newTestGame(playerID)
			var agent engine.PlayerAgent = dummy.NewChooseOneAgent(playerID)
			if tc.agent != nil {
				agent = tc.agent
			}
			got, err := parseCastSpellCommand(tc.arg, game, playerID, agent, false, nil, mockMaybeApplyEvent)
			if err != nil {
				t.Fatalf("parseCastSpellCommand(%q ...); err = %v; want %v", tc.arg, err, nil)
			}
			if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("parseCastSpellCommand(%q ...); mismatch (-want +got):\n%s", tc.arg, diff)
			}
		})
	}
	// #endregion

	// #region Failure Cases
	failCases := []struct {
		name  string
		arg   string
		agent engine.PlayerAgent
		want  error
	}{
		{
			name:  "with invalid name or id",
			arg:   "Invalid Card",
			agent: dummy.NewChooseOneAgent(playerID),
			want:  ErrCardNotFound,
		},
		{
			name:  "with no choice",
			arg:   "",
			agent: dummy.NewChooseNoneAgent(playerID),
			want:  choose.ErrNoChoiceSelected,
		},
	}
	for _, tc := range failCases {
		t.Run(tc.name, func(t *testing.T) {
			game := newTestGame(playerID)
			_, err := parseCastSpellCommand(tc.arg, game, playerID, tc.agent, false, nil, mockMaybeApplyEvent)
			if !errors.Is(err, tc.want) {
				t.Errorf("parseCastSpellCommand(%q ...); err = %v; want %v", tc.arg, err, tc.want)
			}
		})
	}
	// #endregion
}

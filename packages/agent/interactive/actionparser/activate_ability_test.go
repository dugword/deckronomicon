package actionparser

import (
	"deckronomicon/packages/agent/dummy"
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseActivateAbilityCommand(t *testing.T) {
	const playerID = "Test Player"
	// #region Success Cases
	testCases := []struct {
		name  string
		arg   string
		agent engine.PlayerAgent
		want  action.ActivateAbilityRequest
	}{
		{
			name: "with name of ability on card",
			arg:  "Ability on Card",
			want: action.ActivateAbilityRequest{
				AbilityID:         "Card with Ability ID-1",
				SourceID:          "Card with Ability ID",
				Zone:              mtg.ZoneHand,
				TargetsForEffects: map[effect.EffectTargetKey]effect.Target{},
			},
		},
		{
			name: "with id of ability on card",
			arg:  "Card with Ability ID-1",
			want: action.ActivateAbilityRequest{
				AbilityID:         "Card with Ability ID-1",
				SourceID:          "Card with Ability ID",
				Zone:              mtg.ZoneHand,
				TargetsForEffects: map[effect.EffectTargetKey]effect.Target{},
			},
		},
		{
			name: "with name of ability on permanent",
			arg:  "Ability on Permanent",
			want: action.ActivateAbilityRequest{
				AbilityID:         "Test Permanent ID-1",
				SourceID:          "Test Permanent ID",
				Zone:              mtg.ZoneBattlefield,
				TargetsForEffects: map[effect.EffectTargetKey]effect.Target{},
			},
		},
		{
			name: "with id of ability on permanent",
			arg:  "Test Permanent ID-1",
			want: action.ActivateAbilityRequest{
				AbilityID:         "Test Permanent ID-1",
				SourceID:          "Test Permanent ID",
				Zone:              mtg.ZoneBattlefield,
				TargetsForEffects: map[effect.EffectTargetKey]effect.Target{},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game := newTestGame(playerID)
			player := game.GetPlayer(playerID)
			var agent engine.PlayerAgent = dummy.NewChooseOneAgent(playerID)
			if tc.agent != nil {
				agent = tc.agent
			}
			got, err := parseActivateAbilityCommand(tc.arg, game, player, agent)
			if err != nil {
				t.Fatalf("parseActivateAbilityCommand(%q ...); err = %v; want %v", tc.arg, err, nil)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("parseActivateAbilityCommand(%q ...); mismatch (-want +got):\n%s", tc.arg, diff)
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
			arg:   "Invalid Ability",
			agent: dummy.NewChooseOneAgent(playerID),
			want:  ErrAbilityNotFound,
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
			player := game.GetPlayer(playerID)
			_, err := parseActivateAbilityCommand(tc.arg, game, player, tc.agent)
			if !errors.Is(err, tc.want) {
				t.Errorf("parseActivateAbilityCommand(%q ...); err = %v; want %v", tc.arg, err, tc.want)
			}
		})
	}
	// #endregion
}

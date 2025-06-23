package card_test

import (
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/engine/reducer"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/state"
	"deckronomicon/packages/state/statetest"
	"testing"
)

// TODO: Update tests to use standard test failure patterns

func TestStreamOfThought(t *testing.T) {
	resEnv := resenv.ResEnv{}
	const playerID = "Test Player"
	definitions, err := definition.LoadCardDefinitions("../../definitions/cards")
	if err != nil {
		t.Fatalf("Failed to load card definitions: %v", err)
	}
	streamOfThoughtCard, err := gob.NewCardFromCardDefinition("Stream of Thought ID", playerID, definitions["Stream of Thought"])
	if err != nil {
		t.Fatalf("error creating card from definition: %v", err)
	}
	game := state.LoadGameFromConfig(statetest.GameConfig{
		Step: mtg.StepPrecombatMain,
		Players: []statetest.PlayerConfig{
			{
				ID: playerID,
			},
		},
	})
	player := game.GetPlayer(playerID)
	hand := player.Hand().Add(streamOfThoughtCard)
	player = player.WithHand(hand)
	player = player.WithAddMana(mana.Blue, 9)
	game = game.WithUpdatedPlayer(player)
	action := action.NewCastSpellAction(
		action.CastSpellRequest{
			CardID: streamOfThoughtCard.ID(),
			TargetsForEffects: map[action.EffectTargetKey]target.TargetValue{
				{SourceID: "Stream of Thought ID", EffectIndex: 0}: {
					TargetType: target.TargetTypePlayer,
					// TODO: Maybe this could just be TargetID and we'd know based on the target type.
					PlayerID: playerID,
				},
				{SourceID: "Stream of Thought ID", EffectIndex: 1}: target.TargetValue{
					TargetType: target.TargetTypeNone,
				},
			},
		},
	)
	events, err := action.Complete(game, player, &resEnv)
	if err != nil {
		t.Fatalf("Failed to complete action: %v", err)
	}
	for _, event := range events {
		game, err = reducer.ApplyEvent(game, event)
		if err != nil {
			t.Fatalf("Failed to apply event: %v", err)
		}
	}
}

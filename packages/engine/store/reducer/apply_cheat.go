package reducer

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
)

// These are events that manage the priority system in the game.

func applyCheatEvent(game *state.Game, cheatEvent event.CheatEvent) (*state.Game, error) {
	switch evnt := cheatEvent.(type) {
	case *event.CheatAddManaEvent:
		return game, nil
	case *event.CheatConjureCardEvent:
		return applyConjureCardCheatEvent(game, evnt)
	case *event.CheatDiscardEvent:
		return game, nil
	case *event.CheatDrawEvent:
		return game, nil
	case *event.CheatFindCardEvent:
		return game, nil
	case *event.CheatPeekEvent:
		return applyCheatPeekEvent(game, evnt)
	case *event.CheatResetLandDropEvent:
		return applyResetLandDropCheatEvent(game, evnt)
	case *event.CheatShuffleDeckEvent:
		return game, nil
	case *event.CheatUntapEvent:
		return game, nil
	default:
		return game, fmt.Errorf("unknown cheat event type '%T'", evnt)
	}
}

func applyConjureCardCheatEvent(game *state.Game, evnt *event.CheatConjureCardEvent) (*state.Game, error) {
	return game, nil
	/*
		cardDef, ok := e.definitions[evnt.CardName]
		if !ok {
			return game, fmt.Errorf("card definition for %q not found", evnt.CardName)
		}
		id, game := game.GetNextID()
		card, err := gob.NewCardFromCardDefinition(id, evnt.PlayerID, cardDef)
		if err != nil {
			return game, fmt.Errorf("failed to create card from definition: %w", err)
		}
		player, ok := game.GetPlayer(evnt.PlayerID)
		if !ok {
			return game, fmt.Errorf("player %q not found", evnt.PlayerID)
		}
		game = game.WithUpdatedPlayer(
			player.WithHand(
				player.Hand().Add(card),
			),
		)
		return game, nil
	*/
}

func applyCheatPeekEvent(game *state.Game, evnt *event.CheatPeekEvent) (*state.Game, error) {
	player := game.GetPlayer(evnt.PlayerID)
	card := player.Library().Peek()
	revealed := player.Revealed().Add(card)
	game = game.WithUpdatedPlayer(player.WithRevealed(revealed))
	return game, nil
}

func applyResetLandDropCheatEvent(
	game *state.Game,
	evnt *event.CheatResetLandDropEvent,
) (*state.Game, error) {
	player := game.GetPlayer(evnt.PlayerID)
	player = player.WithClearLandPlayedThisTurn()
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

package resolver

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
)

func ResolveShuffleSelfFromGraveyard(
	game *state.Game,
	playerID string,
	shuffleSelfFromGraveyard *effect.ShuffleSelfFromGraveyard,
	// TODO: Update all the functions to take state.Resolvable since that's what it is
	resolvable state.Resolvable,
	resEnv *resenv.ResEnv,
) (Result, error) {
	player := game.GetPlayer(playerID)
	card, ok := player.Graveyard().Get(resolvable.SourceID())
	if !ok {
		panic("handle this")
	}
	var cardIDs []string
	for _, card := range player.Library().GetAll() {
		cardIDs = append(cardIDs, card.ID())
	}
	var events []event.GameEvent
	events = append(events, &event.PutCardOnBottomOfLibraryEvent{
		PlayerID: player.ID(),
		CardID:   card.ID(),
		FromZone: mtg.ZoneGraveyard,
	})
	// TODO: This doesn't feel super clear,
	// but we need the card ID order to shuffle the library.
	// I think I'd rather have the the cards added to the library
	// and then shuffled.
	// Maybe that could be a separate effect?
	// If we do that I'd need to update how the effectResult chanins
	// are handled. Like just a bool "should continue" or something.
	// This could be useful for other effects where I need to
	// apply state changes before the next part of the effect
	// is resolved.
	cardIDs = append(cardIDs, card.ID())

	shuffledCardsIDs := resEnv.RNG.ShuffleIDs(cardIDs)
	events = append(events, &event.ShuffleLibraryEvent{
		PlayerID:         player.ID(),
		ShuffledCardsIDs: shuffledCardsIDs,
	})
	return Result{
		Events: events,
	}, nil
}

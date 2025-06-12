package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/judge"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"fmt"
)

type PlayLandCommand struct {
	CardInZone gob.CardInZone
	Player     state.Player
}

func (p *PlayLandCommand) IsComplete() bool {
	return p.CardInZone.ID() == "" && p.Player.ID() != ""
}

func (p *PlayLandCommand) Build(
	game state.Game,
	player state.Player,
) (engine.Action, error) {
	return action.NewPlayLandAction(player, p.CardInZone), nil
}

func parsePlayLandCommand(
	idOrName string,
	chooseOne func(prompt choose.ChoicePrompt) (choose.Choice, error),
	game state.Game,
	player state.Player,
) (*PlayLandCommand, error) {
	cards := judge.GetLandsAvailableToPlay(game, player)
	if idOrName == "" {
		return buildPlayLandCommandByChoice(cards, chooseOne, player)
	}
	return buildPlayLandCommandByIDOrName(cards, idOrName, player)
}

func buildPlayLandCommandByChoice(
	cards []gob.CardInZone,
	chooseOne func(prompt choose.ChoicePrompt) (choose.Choice, error), player state.Player) (*PlayLandCommand, error) {
	prompt := choose.ChoicePrompt{
		Message:  "Choose a land to play",
		Choices:  choose.NewChoices(cards),
		Source:   CommandSource{"Play a land"},
		Optional: true,
	}
	selected, err := chooseOne(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get choices: %w", err)
	}
	card, ok := selected.(gob.CardInZone)
	if !ok {
		return nil, fmt.Errorf("selected choice is not a card in a zone")
	}
	return &PlayLandCommand{
		CardInZone: card,
		Player:     player,
	}, nil
}

func buildPlayLandCommandByIDOrName(
	cards []gob.CardInZone,
	idOrName string,
	player state.Player,
) (*PlayLandCommand, error) {
	if card, ok := query.Find(cards, query.Or(has.ID(idOrName), has.Name(idOrName))); ok {
		return &PlayLandCommand{
			CardInZone: card,
			Player:     player,
		}, nil
	}
	return nil, fmt.Errorf("no land found with ID or name %q", idOrName)
}

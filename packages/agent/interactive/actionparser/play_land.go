package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/judge"
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
	choose func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
	game state.Game,
	player state.Player,
) (*PlayLandCommand, error) {
	cards := judge.GetLandsAvailableToPlay(game, player)
	if idOrName == "" {
		return buildPlayLandCommandByChoice(cards, choose, player)
	}
	return buildPlayLandCommandByIDOrName(cards, idOrName, player)
}

func buildPlayLandCommandByChoice(
	cards []gob.CardInZone,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
	player state.Player) (*PlayLandCommand, error) {
	prompt := choose.ChoicePrompt{
		Message:  "Choose a land to play",
		Source:   CommandSource{"Play a land"},
		Optional: true,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(cards),
		},
	}
	choiceResults, err := chooseFunc(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get choices: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return nil, fmt.Errorf("expected a single choice result")
	}
	card, ok := selected.Choice.(gob.CardInZone)
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

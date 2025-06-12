package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"fmt"
)

type UntapCheatCommand struct {
	Player    state.Player
	Permanent gob.Permanent
}

func (p *UntapCheatCommand) IsComplete() bool {
	return p.Player.ID() != "" && p.Permanent.ID() != ""
}

func (p *UntapCheatCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return action.NewUntapCheatAction(p.Player, p.Permanent), nil
}

func parseUntapCheatCommand(
	idOrName string,
	chooseOne func(prompt choose.ChoicePrompt) (choose.Choice, error),
	game state.Game,
	player state.Player,
) (*UntapCheatCommand, error) {
	permanents := game.Battlefield().FindAll(
		query.And(has.Controller(player.ID()), is.Tapped()),
	)
	if idOrName == "" {
		return buildUntapCommandByChoice(permanents, chooseOne, player)
	}
	return buildUntapCommandByIDOrName(permanents, idOrName, player)
}

func buildUntapCommandByChoice(
	permanents []gob.Permanent,
	chooseOne func(prompt choose.ChoicePrompt) (choose.Choice, error),
	player state.Player,
) (*UntapCheatCommand, error) {
	prompt := choose.ChoicePrompt{
		Message:  "Choose a permanent to untap",
		Choices:  choose.NewChoices(permanents),
		Source:   CommandSource{"Untap a permanent"},
		Optional: true,
	}
	selected, err := chooseOne(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get choices: %w", err)
	}
	// TODO: Do something where I can select this without having to index an slice with a magic 0.
	// Maybe that choice type that's an interface or something.
	permanent, ok := selected.(gob.Permanent)
	if !ok {
		return nil, fmt.Errorf("selected choice is not a permanent")
	}
	return &UntapCheatCommand{
		Player:    player,
		Permanent: permanent,
	}, nil
}

func buildUntapCommandByIDOrName(
	permanents []gob.Permanent,
	idOrName string,
	player state.Player,
) (*UntapCheatCommand, error) {
	if permanent, ok := query.Find(permanents, query.Or(has.ID(idOrName), has.Name(idOrName))); ok {
		return &UntapCheatCommand{
			Player:    player,
			Permanent: permanent,
		}, nil
	}
	return nil, fmt.Errorf("no permanent found with ID or name %q", idOrName)
}

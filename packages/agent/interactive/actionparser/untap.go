package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/state"
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
	arg string,
	chooseOne func(prompt choose.ChoicePrompt) (choose.Choice, error),
	game state.Game,
	player state.Player,
) (*UntapCheatCommand, error) {
	return nil, nil
}

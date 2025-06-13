package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
)

type ViewCommand struct {
	Player state.Player
	Card   string
	Zone   string
}

func (p *ViewCommand) IsComplete() bool {
	return p.Player.ID() != "" && p.Card != ""
}

func (p *ViewCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return action.NewViewAction(player, p.Zone, p.Card), nil
}

func parseViewCommand(
	arg string,
	choose func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
	game state.Game,
	player state.Player,
) (*ViewCommand, error) {
	return nil, nil
}

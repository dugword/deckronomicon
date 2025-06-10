package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/judge"
	"deckronomicon/packages/mana"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"fmt"
)

type CastSpellCommand struct {
	CardInZone      gob.CardInZone
	Player          state.Player
	SpendMana       mana.Pool
	Targets         []gob.CardInZone
	AdditionalCosts []string
}

func (p *CastSpellCommand) IsComplete() bool {
	return p.CardInZone.ID() == "" && p.Player.ID() != ""
}

func (p *CastSpellCommand) Build(
	game state.Game,
	player state.Player,
) (engine.Action, error) {
	return action.NewCastSpellAction(player, p.CardInZone), nil
}

func parseCastSpellCommand(
	idOrName string,
	chooseOne func(prompt choose.ChoicePrompt) (choose.Choice, error),
	game state.Game,
	player state.Player,
) (*CastSpellCommand, error) {
	cards := judge.GetSpellsAvailableToCast(game, player)
	if idOrName == "" {
		return buildCastSpellCommandByChoice(cards, chooseOne, player)
	}
	return buildCastSpellCommandByIDOrName(cards, idOrName, player)
}

func buildCastSpellCommandByChoice(
	cards []gob.CardInZone,
	chooseOne func(prompt choose.ChoicePrompt) (choose.Choice, error),
	player state.Player,
) (*CastSpellCommand, error) {
	var choices []choose.Choice
	for _, card := range cards {
		choices = append(choices, card)
	}
	prompt := choose.ChoicePrompt{
		Message:  "Choose a card to play",
		Choices:  choices,
		Source:   CommandSource{"Play a card"},
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
	return &CastSpellCommand{
		CardInZone: card,
		Player:     player,
	}, nil
}

func buildCastSpellCommandByIDOrName(
	cards []gob.CardInZone,
	idOrName string,
	player state.Player,
) (*CastSpellCommand, error) {
	if card, ok := query.Find(cards, query.Or(has.ID(idOrName), has.Name(idOrName))); ok {
		return &CastSpellCommand{
			CardInZone: card,
			Player:     player,
		}, nil
	}
	return nil, fmt.Errorf("no land found with ID or name %q", idOrName)
}

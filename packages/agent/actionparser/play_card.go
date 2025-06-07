package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"fmt"
)

type PlayCardCommand struct {
	CardID   string
	PlayerID string
	Zone     mtg.Zone
}

func (p *PlayCardCommand) IsComplete() bool {
	return p.CardID != "" && p.PlayerID != "" && p.Zone == mtg.ZoneHand
}

func (p *PlayCardCommand) Build(
	game state.Game,
	playerID string,
) (engine.Action, error) {
	return engine.NewPlayCardAction(p.PlayerID, p.Zone, p.CardID), nil
}

func parsePlayCardCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	playerID string,
) (*PlayCardCommand, error) {
	if len(args) == 0 {
		var choices []choose.Choice
		cards, err := game.GetCardsAvailableToPlay(playerID)
		if err != nil {
			return nil, fmt.Errorf("failed to get cards available to play: %w", err)
		}
		for _, card := range cards {
			choices = append(choices, choose.Choice{
				Name: card.Name(),
				ID:   card.ID(),
			})
		}
		prompt := choose.ChoicePrompt{
			Message:    "Choose a card to play",
			Choices:    choices,
			Source:     CommandSource{"Play a card"},
			MinChoices: 1,
			MaxChoices: 1,
			Optional:   true,
		}
		selected, err := getChoices(prompt)
		if err != nil {
			return nil, fmt.Errorf("failed to get choices: %w", err)
		}
		if len(selected) == 0 {
			return nil, fmt.Errorf("no card selected")
		}
		return &PlayCardCommand{
			CardID:   selected[0].ID,
			PlayerID: playerID,
			Zone:     mtg.ZoneHand,
		}, nil
	}
	player, err := game.GetPlayer(playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player '%s': %w", playerID, err)
	}
	if player.Hand().Contains(has.ID(args[0])) {
		return &PlayCardCommand{
			CardID:   args[0],
			PlayerID: playerID,
			Zone:     mtg.ZoneHand,
		}, nil
	}
	card, ok := player.Hand().Find(has.Name(args[0]))
	if !ok {
		return parsePlayCardCommand(
			command,
			args,
			getChoices,
			game,
			playerID,
		)
	}
	return &PlayCardCommand{
		CardID:   card.ID(),
		PlayerID: playerID,
		// TODO: Determine the zone based on where the card was found
		Zone: mtg.ZoneHand,
	}, nil
}

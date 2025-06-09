package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/judge"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"fmt"
)

type PlayCardCommand struct {
	CardInZone gob.CardInZone
	Player     state.Player
}

func (p *PlayCardCommand) IsComplete() bool {
	return p.CardInZone.ID() == "" && p.Player.ID() != ""
}

func (p *PlayCardCommand) Build(
	game state.Game,
	player state.Player,
) (engine.Action, error) {
	return engine.NewPlayCardAction(player, p.CardInZone), nil
}

func parsePlayCardCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	player state.Player,
) (*PlayCardCommand, error) {
	if len(args) == 0 {
		var choices []choose.Choice
		cards := judge.GetCardsAvailableToPlay(game, player)
		for _, card := range cards {
			choices = append(choices, card)
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
		// TODO: Do something where I can select this without having to index an slice with a magic 0.
		// Maybe that choice type that's an interface or something.
		card, ok := selected[0].(gob.CardInZone)
		if !ok {
			return nil, fmt.Errorf("selected choice is not a card in a zone")
		}
		return &PlayCardCommand{
			CardInZone: card,
			Player:     player,
		}, nil
	}
	if card, ok := player.Hand().Get(args[0]); ok {
		// If the card is found in the hand, we can return it directly
		return &PlayCardCommand{
			CardInZone: gob.NewCardInZone(card, mtg.ZoneHand),
			Player:     player,
		}, nil
	}
	card, ok := player.Hand().Find(has.Name(args[0]))
	if !ok {
		return parsePlayCardCommand(
			command,
			args,
			getChoices,
			game,
			player,
		)
	}
	return &PlayCardCommand{
		CardInZone: gob.NewCardInZone(card, mtg.ZoneHand),
		Player:     player,
	}, nil
}

package actionparser

// TODO: Document what level things live at. This package is for parsing user
// input or user actions from configuration files. It tries to accurately
// generate requests based on the game state, to be helpful and provide quick
// feed back to the user by only letting them make valid plays, but these are
// just requests to be sent to the game engine. THe game engine is responsible
// for actually verifying things work according to the rules.

// TODO Use a expression tree parser like we do for the JSON parser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"fmt"
	"strings"
)

type CommandSource struct {
	name string
}

func (c CommandSource) Name() string {
	return c.name
}

type Command interface {
	IsComplete() bool
	Build(game state.Game, playerID string) (engine.Action, error)
	// PromptNext(game state.Game, playerID string) (choose.ChoicePrompt, error)
}

type CommandParser struct {
	Command Command
}

type PassPriorityCommand struct {
	PlayerID string
}

func (p *PassPriorityCommand) IsComplete() bool {
	return p.PlayerID != ""
}

func (p *PassPriorityCommand) Build(game state.Game, playerID string) (engine.Action, error) {
	return engine.NewPassPriorityAction(p.PlayerID), nil
}

type PlayCardCommand struct {
	CardID   string
	PlayerID string
}

func (p *PlayCardCommand) IsComplete() bool {
	return p.CardID != "" && p.PlayerID != ""
}

func (p *PlayCardCommand) Build(
	game state.Game,
	playerID string,
) (engine.Action, error) {
	return engine.NewPlayCardAction(p.PlayerID, p.CardID), nil
}

func (p *CommandParser) ParseInput(
	getInput func() string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	playerID string,
) (Command, error) {
	for {
		input := getInput()
		fmt.Println("Input received:", input)
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}
		command, args := parts[0], parts[1:]
		switch command {
		case "pass":
			return &PassPriorityCommand{playerID}, nil
		case "play":
			return parsePlayCardCommand(command, args, getChoices, game, playerID)
		default:
			return nil, fmt.Errorf("unrecognized command '%s'", command)
		}
	}
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
		}
		selected, err := getChoices(prompt)
		if err != nil {
			return nil, fmt.Errorf("failed to get choices: %w", err)
		}
		return &PlayCardCommand{
			CardID:   selected[0].ID,
			PlayerID: playerID,
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
	}, nil
}

/*


type PlayCommand struct {
	CardID string
}

func (c *PlayCommand) IsComplete() bool {
	return c.CardID != ""
}

func (c *PlayCommand) PromptNext(game state.Game, playerID string) (choose.ChoicePrompt, error) {
	if c.CardID == "" {
		player, err := game.GetPlayer(playerID)
		if err != nil {
			return choose.ChoicePrompt{}, fmt.Errorf("failed to get player %s: %w", playerID, err)
		}
		var choices []choose.Choice
		for _, card := range player.Hand().GetAll() {
			if game.CanPlayCard(card, playerID) {
				choices = append(choices, choose.Choice{
					Name: card.Name(),
					ID:   card.ID(),
				})
			}
		}
		return choose.ChoicePrompt{
			Message: "Choose a card to play",
			Choices: choices,
			Source:  CommandHelper{"Choose a card to play"},
		}, nil
	}
	return choose.ChoicePrompt{}, nil
}

func (c *PlayCommand) AddInput(input string) error {
	c.CardID = input
	return nil
}

func (c *PlayCommand) Build(game state.Game, playerID string) (engine.Action, error) {
	/*
		if c.CardID == "" {
			return nil, fmt.Errorf("no card selected to play")
		}
		player, err := game.GetPlayer(playerID)
		if err != nil {
			return nil, fmt.Errorf("failed to get player %s: %w", playerID, err)
		}
		card := player.Hand().GetByID(c.CardID)
		if card == nil {
			return nil, fmt.Errorf("card with ID %s not found in hand", c.CardID)
		}
		return engine.NewPlayCardAction(playerID, card.ID()), nil
	return nil, nil
}

type BuildActionRequest struct {
	Command string
	Args    []string
}

type BuildActionResult struct {
	Action             engine.Action
	ChoicePrompt       choose.ChoicePrompt
	BuildActionRequest BuildActionRequest
}

type Parser interface {
	Parse(input string) (BuildActionRequest, error)
}

func ParseCommand(input string) BuildActionRequest {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return BuildActionRequest{}
	}
	return BuildActionRequest{
		Command: strings.ToLower(parts[0]),
		Args:    parts[1:],
	}
}

func BuildAction(req BuildActionRequest, game state.Game, playerID string) (BuildActionResult, error) {
	switch req.Command {
	case "pass":
		return buildPassAction(playerID), nil
	case "play":
		return buildPlayAction(req, game, playerID)
	default:
		return BuildActionResult{}, fmt.Errorf("unrecognized command: %s", req.Command)
	}
}

func buildPassAction(playerID string) BuildActionResult {
	return BuildActionResult{
		Action: engine.NewPassPriorityAction(playerID),
	}
}

func buildPlayAction(req BuildActionRequest, game state.Game, playerID string) (BuildActionResult, error) {
	player, err := game.GetPlayer(playerID)
	if err != nil {
		return BuildActionResult{}, fmt.Errorf("failed to get hand for player %s: %w", playerID, err)
	}
	if len(req.Args) == 0 {
	}
	arg := req.Args[0]
	var candidates []gob.Card
	// TODO use filter function
	for _, card := range player.Hand().GetAll() {
		if card.ID() == arg || card.Name() == arg {
			candidates = append(candidates, card)
		}
	}

	switch len(candidates) {
	case 0:
		return BuildActionResult{}, fmt.Errorf("no playable card found for '%s'", arg)
	case 1:
		return BuildActionResult{
			Action: engine.NewPlayCardAction(playerID, candidates[0].ID()),
		}, nil
	default:
		choices := make([]choose.Choice, len(candidates))
		for i, card := range candidates {
			choices[i] = choose.Choice{
				Name: card.Name(),
				ID:   card.ID(),
			}
		}
		return BuildActionResult{
			ChoicePrompt: choose.ChoicePrompt{
				Message: "Choose a card to play",
				Choices: choices,
				Source:  CommandHelper{"Choose a card to play"},
			},
		}, nil
	}
}

/*
// ParseCommand parses a simple string command like "pass", "play", or "activate"
// and returns the appropriate GameAction for the engine to process.
func ParseCommand(
	input string,
	game state.Game,
	playerID string,
) (engine.Action, error) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty command")
	}
	cmd := strings.ToLower(fields[0])
	args := fields[1:]

	switch cmd {
	case "pass":
		return engine.NewPassPriorityAction(playerID), nil

	case "play":
		if len(args) == 0 {
			return engine.NewPlayCardAction(playerID, cardID), nil
		}
		cardArg := strings.Join(args, " ")
		return engine.NewPlayCardAction(playerID, cardArg), nil

	case "activate":
		if len(args) == 0 {
			return engine.NewActivatePromptAction(playerID), nil
		}
		abilityArg := strings.Join(args, " ")
		return engine.NewActivateAction(playerID, abilityArg), nil

	default:
		return nil, fmt.Errorf("unrecognized command: %s", cmd)
	}
}
*/

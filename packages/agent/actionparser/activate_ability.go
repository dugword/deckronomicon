package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/judge"
	"deckronomicon/packages/state"
	"fmt"
)

type ActivateAbilityCommand struct {
	Player                state.Player
	AbilityOnObjectInZone gob.AbilityInZone
}

func (p *ActivateAbilityCommand) IsComplete() bool {
	return p.Player.ID() != "" && p.AbilityOnObjectInZone.ID() != ""
}

// TODO: What is this really even doing? Maybe just return a new action instead of building a command?
func (p *ActivateAbilityCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return engine.NewActivateAbilityAction(p.Player, p.AbilityOnObjectInZone), nil
}

func parseActivateAbilityCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	player state.Player,
) (*ActivateAbilityCommand, error) {
	if len(args) == 0 {
		var choices []choose.Choice
		abilities := judge.GetAbilitiesAvailableToActivate(game, player)
		for _, ability := range abilities {
			choices = append(choices, ability)
		}
		prompt := choose.ChoicePrompt{
			Message:    "Choose an ability to activate",
			Choices:    choices,
			Source:     CommandSource{"Activate an ability"},
			MinChoices: 1,
			MaxChoices: 1,
			Optional:   true,
		}
		selected, err := getChoices(prompt)
		if err != nil {
			return nil, fmt.Errorf("failed to get choices: %w", err)
		}
		if len(selected) == 0 {
			return nil, fmt.Errorf("no ability selected")
		}
		// TODO: Do something where I can select this without having to index an slice with a magic 0.
		// Maybe that choice type that's an interface or something.
		ability, ok := selected[0].(gob.AbilityInZone)
		if !ok {
			return nil, fmt.Errorf("selected choice is not an ability on an object in a zone")
		}
		return &ActivateAbilityCommand{
			Player:                player,
			AbilityOnObjectInZone: ability,
		}, nil
	}
	return parseActivateAbilityCommand(
		command,
		args,
		getChoices,
		game,
		player,
	)
	/* // TODO: This doesn't work, need to get abilities not permanents
	if game.Battlefield().Contains(has.ID(args[0])) {
		return &ActivateAbilityCommand{
			Player: player,
			AbilityOnObjectInZone:
		}, nil
	}

	ability, ok := game.Battlefield().Find(has.Name(args[0]))
	if !ok {
		return parseActivateAbilityCommand(
			command,
			args,
			getChoices,
			game,
			player,
		)
	}
	return &ActivateAbilityCommand{
		Player:  player,
		AbilityID: ability.ID(),
	}, nil
	*/
}

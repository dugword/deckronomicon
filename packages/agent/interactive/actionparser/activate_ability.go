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

type ActivateAbilityCommand struct {
	Player                state.Player
	AbilityOnObjectInZone gob.AbilityInZone
}

func (p *ActivateAbilityCommand) IsComplete() bool {
	return p.Player.ID() != "" && p.AbilityOnObjectInZone.ID() != ""
}

// TODO: What is this really even doing? Maybe just return a new action instead of building a command?
func (p *ActivateAbilityCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return action.NewActivateAbilityAction(p.Player, p.AbilityOnObjectInZone), nil
}

func parseActivateAbilityCommand(
	idOrName string,
	chooseOne func(prompt choose.ChoicePrompt) (choose.Choice, error),
	game state.Game,
	player state.Player,
) (*ActivateAbilityCommand, error) {
	ruling := judge.Ruling{Explain: true}
	abilities := judge.GetAbilitiesAvailableToActivate(game, player, &ruling)
	if idOrName == "" {
		return buildActivateAbilityCommandByChoice(abilities, chooseOne, player)
	}
	return buildActivateAbilityCommandByIDOrName(abilities, idOrName, player)
}

func buildActivateAbilityCommandByChoice(
	abilities []gob.AbilityInZone,
	chooseOne func(prompt choose.ChoicePrompt) (choose.Choice, error),
	player state.Player,
) (*ActivateAbilityCommand, error) {
	prompt := choose.ChoicePrompt{
		Message:  "Choose an ability to activate",
		Choices:  choose.NewChoices(abilities),
		Source:   CommandSource{"Activate an ability"},
		Optional: true,
	}
	selected, err := chooseOne(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get choices: %w", err)
	}
	// TODO: Do something where I can select this without having to index an slice with a magic 0.
	// Maybe that choice type that's an interface or something.
	ability, ok := selected.(gob.AbilityInZone)
	if !ok {
		return nil, fmt.Errorf("selected choice is not an ability on an object in a zone")
	}
	return &ActivateAbilityCommand{
		Player:                player,
		AbilityOnObjectInZone: ability,
	}, nil
}

func buildActivateAbilityCommandByIDOrName(
	abilities []gob.AbilityInZone,
	idOrName string,
	player state.Player,
) (*ActivateAbilityCommand, error) {
	if ability, ok := query.Find(abilities, query.Or(has.ID(idOrName), has.Name(idOrName))); ok {
		return &ActivateAbilityCommand{
			Player:                player,
			AbilityOnObjectInZone: ability,
		}, nil
	}
	return nil, fmt.Errorf("no ability found with ID or name %q", idOrName)
}

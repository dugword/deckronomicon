package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/engine/effect"
	"deckronomicon/packages/engine/target"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/judge"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"fmt"
)

type CastSpellCommand struct {
	CardInZone      gob.CardInZone
	Player          state.Player
	SpendMana       mana.Pool
	Targets         map[string]target.TargetValue
	AdditionalCosts []string
}

func (p *CastSpellCommand) IsComplete() bool {
	return p.CardInZone.ID() == "" && p.Player.ID() != ""
}

func (p *CastSpellCommand) Build(
	game state.Game,
	player state.Player,
) (engine.Action, error) {
	return action.NewCastSpellAction(player, p.CardInZone, p.Targets), nil
}

func parseCastSpellCommand(
	idOrName string,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
	game state.Game,
	player state.Player,
) (*CastSpellCommand, error) {
	cards := judge.GetSpellsAvailableToCast(game, player)
	var card gob.CardInZone
	var err error
	if idOrName == "" {
		card, err = getCardByChoice(cards, chooseFunc, player)
		if err != nil {
			return nil, fmt.Errorf("failed to get card by choice: %w", err)
		}
	} else {
		found, ok := query.Find(cards, query.Or(has.ID(idOrName), has.Name(idOrName)))
		if !ok {
			return nil, fmt.Errorf("no spell found with id or name %q", idOrName)
		}
		card = found
	}
	targets, err := getTargetsForSpell(card, chooseFunc, game, player)
	if err != nil {
		return nil, fmt.Errorf("failed to get targets for spell: %w", err)
	}
	return &CastSpellCommand{
		CardInZone: card,
		Player:     player,
		Targets:    targets,
	}, nil
}

func getTargetsForSpell(
	card gob.CardInZone,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
	game state.Game,
	player state.Player,
) (map[string]target.TargetValue, error) {
	targets := map[string]target.TargetValue{}
	for _, effectSpec := range card.Card().SpellAbility() {
		efct, err := effect.Build(effectSpec)
		if err != nil {
			return nil, fmt.Errorf("effect %q not found: %w", effectSpec.Name, err)
		}
		switch targetSpec := efct.TargetSpec().(type) {
		case nil, target.NoneTargetSpec:
			targets[effectSpec.Name] = target.TargetValue{
				TargetType: target.TargetTypeNone,
			}
		// TODO: Move these to functions
		case target.PlayerTargetSpec:
			prompt := choose.ChoicePrompt{
				Message: "Choose a player to target",
				Source:  CommandSource{"Cast a spell"},
				ChoiceOpts: choose.ChooseOneOpts{
					Choices: choose.NewChoices(game.Players()),
				},
			}
			choiceResults, err := chooseFunc(prompt)
			if err != nil {
				return nil, fmt.Errorf("failed to get choice results: %w", err)
			}
			selected, ok := choiceResults.(choose.ChooseOneResults)
			if !ok {
				return nil, fmt.Errorf("expected a single choice result")
			}
			selectedPlayer, ok := selected.Choice.(state.Player)
			if !ok {
				return nil, fmt.Errorf("selected choice is not a player")
			}
			targets[effectSpec.Name] = target.TargetValue{
				TargetType: target.TargetTypePlayer,
				PlayerID:   selectedPlayer.ID(),
			}
		case target.SpellTargetSpec:
			spells := game.Stack().FindAll(targetSpec.Predicate)
			prompt := choose.ChoicePrompt{
				Message: "Choose a spell to target",
				Source:  CommandSource{"Target a spell"},
				ChoiceOpts: choose.ChooseOneOpts{
					Choices: choose.NewChoices(spells),
				},
			}
			choiceResults, err := chooseFunc(prompt)
			if err != nil {
				return nil, fmt.Errorf("failed to get choice results: %w", err)
			}
			selected, ok := choiceResults.(choose.ChooseOneResults)
			if !ok {
				return nil, fmt.Errorf("expected a single choice result")
			}
			selectedSpell, ok := selected.Choice.(gob.Spell)
			if !ok {
				return nil, fmt.Errorf("selected choice is not a spell")
			}
			targets[effectSpec.Name] = target.TargetValue{
				TargetType: target.TargetTypeSpell,
				ObjectID:   selectedSpell.ID(),
			}
		default:
			return nil, fmt.Errorf("unsupported target spec type %T", targetSpec)
		}
	}
	return targets, nil
}

func getCardByChoice(
	cards []gob.CardInZone,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
	player state.Player,
) (gob.CardInZone, error) {
	prompt := choose.ChoicePrompt{
		Message: "Choose a spell to cast",

		Source:   CommandSource{"Cast a spell"},
		Optional: true,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(cards),
		},
	}
	choiceResults, err := chooseFunc(prompt)
	if err != nil {
		return gob.CardInZone{}, fmt.Errorf("failed to get choices: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return gob.CardInZone{}, fmt.Errorf("selected choice is not a card in a zone")
	}
	card, ok := selected.Choice.(gob.CardInZone)
	if !ok {
		return gob.CardInZone{}, fmt.Errorf("selected choice is not a card in a zone")
	}
	return card, nil
}

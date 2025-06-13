package effect

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

func CounterspellEffectHandler(
	game state.Game,
	player state.Player,
	source query.Object,
	modifiers []definition.EffectModifier,
) (EffectResult, error) {
	query, err := buildQuery(modifiers)
	if err != nil {
		return EffectResult{}, fmt.Errorf("failed to build query for Search effect: %w", err)
	}
	spells := game.Stack().FindAll(query)
	choicePrompt := choose.ChoicePrompt{
		// TODO: provide more detail on what kind of card to choose
		Message: "Choose a spell to counter",
		Source:  source,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(spells),
		},
	}
	resumeFunc := func(choiceResults choose.ChoiceResults) (EffectResult, error) {
		selected, ok := choiceResults.(choose.ChooseOneResults)
		if !ok {
			return EffectResult{}, fmt.Errorf("expected a single choice result")
		}
		spell, ok := selected.Choice.(gob.Spell)
		if !ok {
			return EffectResult{}, errors.New("choice is not a spell")
		}
		var events []event.GameEvent
		if spell.Flashback() {
			events = append(events, event.PutSpellInExileEvent{
				PlayerID: player.ID(),
				SpellID:  spell.ID(),
			})
		} else {
			events = append(events, event.PutSpellInGraveyardEvent{
				PlayerID: player.ID(),
				SpellID:  spell.ID(),
			})
		}
		return EffectResult{
			Events: events,
		}, nil
	}
	// Need to get choices
	return EffectResult{
		ChoicePrompt: choicePrompt,
		ResumeFunc:   resumeFunc,
	}, nil
}

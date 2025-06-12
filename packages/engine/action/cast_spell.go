package action

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/judge"
	"deckronomicon/packages/mana"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"fmt"
)

type CastSpellAction struct {
	player          state.Player
	cardInZone      gob.CardInZone
	ManaPayment     mana.Pool
	AdditionalCosts []string
	Targets         []string
}

func NewCastSpellAction(player state.Player, cardInZone gob.CardInZone) CastSpellAction {
	return CastSpellAction{
		player:     player,
		cardInZone: cardInZone,
	}
}

func (a CastSpellAction) PlayerID() string {
	return a.player.ID()
}

func (a CastSpellAction) Name() string {
	return "Cast Spell"
}
func (a CastSpellAction) Description() string {
	return "The active player casts a spell from their hand."
}
func (a CastSpellAction) GetPrompt(state state.Game) (choose.ChoicePrompt, error) {
	/*
		choices := state.GetPlayableCards(a.playerID)
		if len(choices) == 0 {
			return choose.ChoicePrompt{
				Message:  "No playable cards",
				Choices:  nil,
				Optional: false,
			}, nil
		}
		return choose.ChoicePrompt{
			Message:  "Choose a card to play",
			Choices:  choices,
			Optional: false,
		}, nil
	*/
	return choose.ChoicePrompt{}, nil
}

func (a CastSpellAction) Complete(
	game state.Game,
	env *ResolutionEnvironment,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	// TODO: This should probably be a part of CanCastSpell from judge maybe?
	if !a.player.ZoneContains(a.cardInZone.Zone(), has.ID(a.cardInZone.ID())) {
		return nil, fmt.Errorf(
			"player %q does not have card %q in zone %q",
			a.player.ID(),
			a.cardInZone.ID(),
			a.cardInZone.Zone(),
		)
	}
	ruling := judge.Ruling{Explain: true}
	if !judge.CanCastSpell(
		game,
		a.player,
		a.cardInZone.Zone(),
		a.cardInZone.Card(),
		&ruling,
	) {
		return nil, fmt.Errorf(
			"player %q cannot cast card %q from zone %q: %s",
			a.player.ID(),
			a.cardInZone.ID(),
			a.cardInZone.Zone(),
			ruling.Why(),
		)
	}
	costEvents, err := PayCost(a.cardInZone.Card().ManaCost(), a.cardInZone.Card(), a.player)
	if err != nil {
		return nil, fmt.Errorf("failed to pay cost: %w", err)
	}
	events := append(costEvents,
		event.CastSpellEvent{
			PlayerID: a.player.ID(),
			CardID:   a.cardInZone.ID(),
			Zone:     a.cardInZone.Zone(),
		},
		event.PutSpellOnStackEvent{
			PlayerID: a.player.ID(),
			CardID:   a.cardInZone.ID(),
			FromZone: a.cardInZone.Zone(),
		},
	)
	return events, nil
}

package action

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/judge"
	"deckronomicon/packages/mana"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"fmt"
)

type PlayLandAction struct {
	player          state.Player
	cardInZone      gob.CardInZone
	ManaPayment     mana.Pool
	AdditionalCosts []string
	Targets         []string
}

func NewPlayLandAction(player state.Player, cardInZone gob.CardInZone) PlayLandAction {
	return PlayLandAction{
		player:     player,
		cardInZone: cardInZone,
	}
}

func (a PlayLandAction) PlayerID() string {
	return a.player.ID()
}

func (a PlayLandAction) Name() string {
	return "Play Land"
}

// TODO: Not sure if I like this pattern, I think I want to have the command/strategy parser
// Figureout everything that needs to be done and then just pass the choices directly to the action.
// However, I'm not sure how to handle the turn based actions like declare attackers, declare blockers, and assign combat damage.
// Maybe I need to split out Turn Based Actions from Player Actions.
func (a PlayLandAction) GetPrompt(state state.Game) (choose.ChoicePrompt, error) {
	return choose.ChoicePrompt{}, nil
}

func (a PlayLandAction) Complete(
	game state.Game,
	env *ResolutionEnvironment,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	// TODO: Maybe this should happen in judge.CanPlayLand?
	if !a.player.ZoneContains(a.cardInZone.Zone(), has.ID(a.cardInZone.ID())) {
		return nil, fmt.Errorf(
			"player %q does not have card %q in zone %q",
			a.player.ID(),
			a.cardInZone.ID(),
			a.cardInZone.Zone(),
		)
	}
	ruling := judge.Ruling{Explain: true}
	// TODO: Think through structured errors or "reason" structs for responses from the judge.
	if !judge.CanPlayLand(game, a.player, a.cardInZone.Zone(), a.cardInZone.Card(), &ruling) {
		return nil, fmt.Errorf(
			"player %q cannot play card %q, %s",
			a.player.ID(),
			a.cardInZone.ID(),
			ruling.Why(),
		)
	}
	if a.cardInZone.Card().Match(is.Land()) {
		return []event.GameEvent{
			event.PlayLandEvent{
				PlayerID: a.player.ID(),
				CardID:   a.cardInZone.ID(),
				Zone:     a.cardInZone.Zone(),
			},
			event.PutCardOnBattlefieldEvent{
				PlayerID: a.player.ID(),
				CardID:   a.cardInZone.ID(),
				FromZone: a.cardInZone.Zone(),
			},
		}, nil
	}
	return nil, fmt.Errorf(
		"card %q is not a land, cannot play",
		a.cardInZone.ID(),
	)
}

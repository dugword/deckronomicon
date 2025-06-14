package action

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/judge"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
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
	choiceResults choose.ChoiceResults,
) ([]event.GameEvent, error) {

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
	switch a.cardInZone.Zone() {
	case mtg.ZoneGraveyard:
		return castFromGraveyard(game, a.player, a.cardInZone.Card())
	case mtg.ZoneHand:
		return castFromHand(game, a.player, a.cardInZone.Card())
	default:
		return nil, fmt.Errorf("casting from zone %q not implemented", a.cardInZone.Zone())
	}
}

func castFromGraveyard(game state.Game, player state.Player, card gob.Card) ([]event.GameEvent, error) {
	if !card.Match(has.StaticKeyword(mtg.StaticKeywordFlashback)) {
		// Only flashback is supported for now
		return nil, fmt.Errorf("card %q does not have flashback", card.ID())
	}
	var costString string
	// TODO: Maybe these should be methods on the Card type?
	// Or maybe a helper function?
	for _, ability := range card.StaticAbilities() {
		if ability.StaticKeyword() != mtg.StaticKeywordFlashback {
			continue
		}
		for _, modifier := range ability.Modifiers {
			if modifier.Key != "Cost" {
				continue
			}
			costString = modifier.Value
		}
	}
	if costString == "" {
		return nil, fmt.Errorf("flashback ability on card %q is missing Cost modifier", card.Name())
	}
	costs, err := cost.ParseCost(costString, card)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cost: %w", err)
	}
	costEvents, err := PayCost(costs, card, player)
	if err != nil {
		return nil, fmt.Errorf("failed to pay cost: %w", err)
	}
	events := append(costEvents,
		event.CastSpellEvent{
			PlayerID: player.ID(),
			CardID:   card.ID(),
			FromZone: mtg.ZoneGraveyard,
		},
		event.PutSpellOnStackEvent{
			PlayerID:  player.ID(),
			CardID:    card.ID(),
			FromZone:  mtg.ZoneGraveyard,
			Flashback: true,
		},
	)
	return events, nil
}

func castFromHand(game state.Game, player state.Player, card gob.Card) ([]event.GameEvent, error) {
	costEvents, err := PayCost(card.ManaCost(), card, player)
	if err != nil {
		return nil, fmt.Errorf("failed to pay cost: %w", err)
	}
	events := append(costEvents,
		event.CastSpellEvent{
			PlayerID: player.ID(),
			CardID:   card.ID(),
			FromZone: mtg.ZoneHand,
		},
		event.PutSpellOnStackEvent{
			PlayerID: player.ID(),
			CardID:   card.ID(),
			FromZone: mtg.ZoneHand,
		},
	)
	return events, nil
}

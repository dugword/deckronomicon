package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/static_abilities"
	"deckronomicon/packages/engine/target"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/judge"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"encoding/json"
	"fmt"
)

type CastSpellAction struct {
	player          state.Player
	cardInZone      gob.CardInZone
	manaPayment     mana.Pool
	additionalCosts []string
	targets         map[string]target.TargetValue
}

func NewCastSpellAction(player state.Player, cardInZone gob.CardInZone, targets map[string]target.TargetValue) CastSpellAction {
	return CastSpellAction{
		player:     player,
		cardInZone: cardInZone,
		targets:    targets,
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

func (a CastSpellAction) Complete(game state.Game) ([]event.GameEvent, error) {
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
		return a.castFromGraveyard(game)
	case mtg.ZoneHand:
		return a.castFromHand(game)
	default:
		return nil, fmt.Errorf("casting from zone %q not implemented", a.cardInZone.Zone())
	}
}

func (a CastSpellAction) castFromGraveyard(game state.Game) ([]event.GameEvent, error) {
	if !a.cardInZone.Card().Match(has.StaticKeyword(mtg.StaticKeywordFlashback)) {
		// Only flashback is supported for now
		return nil, fmt.Errorf("card %q does not have flashback", a.cardInZone.ID())
	}
	var costString string
	// TODO: Maybe these should be methods on the Card type?
	// Or maybe a helper function?
	for _, ability := range a.cardInZone.Card().StaticAbilities() {
		if ability.StaticKeyword() != mtg.StaticKeywordFlashback {
			continue
		}
		if ability.StaticKeyword() != mtg.StaticKeywordFlashback {
			continue
		}
		//TODO Have some kind of helper function to get static abilities
		// or maybe some kind of interface so I can have these live on the card type....
		flashbackModifiers := static_abilities.FlashbackModifiers{}
		if err := json.Unmarshal(ability.Modifiers, &flashbackModifiers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal flashback modifiers: %w", err)
		}
		costString = flashbackModifiers.Cost
	}
	if costString == "" {
		return nil, fmt.Errorf("flashback ability on card %q is missing Cost modifier", a.cardInZone.Name())
	}
	costs, err := cost.ParseCost(costString, a.cardInZone.Card())
	if err != nil {
		return nil, fmt.Errorf("failed to parse cost: %w", err)
	}
	costEvents, err := PayCost(costs, a.cardInZone.Card(), a.player)
	if err != nil {
		return nil, fmt.Errorf("failed to pay cost: %w", err)
	}
	events := append(costEvents,
		event.CastSpellEvent{
			PlayerID: a.player.ID(),
			CardID:   a.cardInZone.ID(),
			FromZone: mtg.ZoneGraveyard,
		},
		event.PutSpellOnStackEvent{
			PlayerID:  a.player.ID(),
			CardID:    a.cardInZone.ID(),
			FromZone:  mtg.ZoneGraveyard,
			Targets:   a.targets,
			Flashback: true,
		},
	)
	return events, nil
}

func (a CastSpellAction) castFromHand(game state.Game) ([]event.GameEvent, error) {
	costEvents, err := PayCost(a.cardInZone.Card().ManaCost(), a.cardInZone.Card(), a.player)
	if err != nil {
		return nil, fmt.Errorf("failed to pay cost: %w", err)
	}
	events := append(costEvents,
		event.CastSpellEvent{
			PlayerID: a.player.ID(),
			CardID:   a.cardInZone.ID(),
			FromZone: mtg.ZoneHand,
		},
		event.PutSpellOnStackEvent{
			PlayerID: a.player.ID(),
			CardID:   a.cardInZone.ID(),
			Targets:  a.targets,
			FromZone: mtg.ZoneHand,
		},
	)
	return events, nil
}

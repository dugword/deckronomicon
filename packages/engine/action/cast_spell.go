package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/judge"
	"deckronomicon/packages/engine/pay"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/state"
	"fmt"
)

type CastSpellRequest struct {
	CardID            string
	SpliceCardIDs     []string
	TargetsForEffects map[EffectTargetKey]target.TargetValue
	WithFlashback     bool
}

type CastSpellAction struct {
	playerID          string
	cardID            string
	targetsForEffects map[EffectTargetKey]target.TargetValue
	spliceCardIDs     []string
	withFlashback     bool
}

func NewCastSpellAction(
	playerID string,
	request CastSpellRequest,
) CastSpellAction {
	return CastSpellAction{
		playerID:          playerID,
		cardID:            request.CardID,
		targetsForEffects: request.TargetsForEffects,
		spliceCardIDs:     request.SpliceCardIDs,
		withFlashback:     request.WithFlashback,
	}
}

func (a CastSpellAction) Name() string {
	return "Cast Spell"
}

func (a CastSpellAction) Complete(game state.Game, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	player, ok := game.GetPlayer(a.playerID)
	if !ok {
		return nil, fmt.Errorf("player %q not found in game", a.playerID)
	}
	if a.withFlashback {
		return a.castWithFlashback(game, player)
	}
	return a.castFromHand(game, player)
}

func (a CastSpellAction) castWithFlashback(game state.Game, player state.Player) ([]event.GameEvent, error) {
	cardToCast, ok := player.GetCardFromZone(a.cardID, mtg.ZoneGraveyard)
	if !ok {
		return nil, fmt.Errorf("player %q does not have card %q in graveyard", a.playerID, a.cardID)
	}
	flashbackCost, ok := cardToCast.GetStaticAbilityCost(mtg.StaticKeywordFlashback)
	if !ok {
		return nil, fmt.Errorf("card %q does not have flashback", cardToCast.ID())
	}
	ruling := judge.Ruling{Explain: true}
	if !judge.CanCastSpellWithFlashback(
		game,
		player,
		cardToCast,
		flashbackCost,
		&ruling,
	) {
		return nil, fmt.Errorf(
			"player %q cannot cast card %q with flashback: %s",
			a.playerID,
			cardToCast.ID(),
			ruling.Why(),
		)
	}
	effectWithTargets, err := buildEffectWithTargets(cardToCast.ID(), cardToCast.SpellAbility(), a.targetsForEffects)
	if err != nil {
		return nil, err
	}
	costEvents := pay.PayCost(flashbackCost, cardToCast, player)
	events := append(
		costEvents,
		event.CastSpellEvent{
			PlayerID: a.playerID,
			CardID:   cardToCast.ID(),
			FromZone: mtg.ZoneGraveyard,
		},
		event.PutSpellOnStackEvent{
			PlayerID:          a.playerID,
			CardID:            cardToCast.ID(),
			FromZone:          mtg.ZoneGraveyard,
			EffectWithTargets: effectWithTargets,
			Flashback:         true,
		},
	)
	return events, nil
}

func (a CastSpellAction) castFromHand(game state.Game, player state.Player) ([]event.GameEvent, error) {
	cardToCast, ok := player.GetCardFromZone(a.cardID, mtg.ZoneHand)
	if !ok {
		return nil, fmt.Errorf("player %q does not have card %q in hand", a.playerID, a.cardID)
	}
	var totalCost cost.Cost = cardToCast.ManaCost()
	effectWithTargets, err := buildEffectWithTargets(cardToCast.ID(), cardToCast.SpellAbility(), a.targetsForEffects)
	if err != nil {
		return nil, err
	}
	if len(a.spliceCardIDs) > 0 {
		spliceEffectWithCards, spliceCosts, err := buildSpliceEffectWithCards(
			player,
			cardToCast,
			a.spliceCardIDs,
			a.targetsForEffects,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to build splice effect with cards: %w", err)
		}
		effectWithTargets = append(effectWithTargets, spliceEffectWithCards...)
		totalCost = cost.CombineCosts(totalCost, spliceCosts...)
	}
	ruling := judge.Ruling{Explain: true}
	if !judge.CanCastSpellFromHand(
		game,
		player,
		cardToCast,
		totalCost,
		&ruling,
	) {
		return nil, fmt.Errorf(
			"player %q cannot cast card %q from from hand: %s",
			a.playerID,
			cardToCast.ID(),
			ruling.Why(),
		)
	}
	costEvents := pay.PayCost(totalCost, cardToCast, player)
	events := append(
		costEvents,
		event.CastSpellEvent{
			PlayerID: a.playerID,
			CardID:   cardToCast.ID(),
			FromZone: mtg.ZoneHand,
		},
		event.PutSpellOnStackEvent{
			PlayerID:          a.playerID,
			CardID:            cardToCast.ID(),
			EffectWithTargets: effectWithTargets,
			FromZone:          mtg.ZoneHand,
		},
	)
	return events, nil
}

func buildSpliceEffectWithCards(
	player state.Player,
	cardToCast gob.Card,
	spliceCardIDs []string,
	targetsForEffects map[EffectTargetKey]target.TargetValue,
) ([]gob.EffectWithTarget, []cost.Cost, error) {
	var spliceCosts []cost.Cost
	var effectWithTargets []gob.EffectWithTarget
	for _, spliceCardID := range spliceCardIDs {
		spliceCard, ok := player.GetCardFromZone(spliceCardID, mtg.ZoneHand)
		if !ok {
			return nil, nil, fmt.Errorf("card %q not found in hand for splicing", spliceCardID)
		}
		var ruling judge.Ruling
		if !judge.CanSpliceCard(
			player,
			cardToCast,
			spliceCard,
			&ruling,
		) {
			return nil, nil, fmt.Errorf(
				"player %q cannot splice card %q onto spell %q: %s",
				player.ID(),
				spliceCardID,
				cardToCast.ID(),
				ruling.Why(),
			)
		}
		staticAbilityCost, ok := spliceCard.GetStaticAbilityCost(mtg.StaticKeywordSplice)
		if !ok {
			return nil, nil, fmt.Errorf("card %q does not have splice ability", spliceCardID)
		}
		spliceCosts = append(spliceCosts, staticAbilityCost)
		effectWithSpliceTargets, err := buildEffectWithTargets(
			spliceCard.ID(),
			spliceCard.SpellAbility(),
			targetsForEffects,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to build effect with targets for splice card %q: %w", spliceCardID, err)
		}
		effectWithTargets = append(effectWithTargets, effectWithSpliceTargets...)
	}
	return effectWithTargets, spliceCosts, nil
}

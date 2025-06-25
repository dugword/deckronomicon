package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/judge"
	"deckronomicon/packages/engine/pay"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/state"
	"fmt"
)

// TODO: Probably needs to include zone for cast from exile, graveyard, etc.
type CastSpellRequest struct {
	CardID            string
	ReplicateCount    int
	SpliceCardIDs     []string
	TargetsForEffects map[target.EffectTargetKey]target.TargetValue
	Flashback         bool
	AutoPayCost       bool
	Preactions        []ActivateAbilityRequest // Preactions to be executed before casting the spell
}

func (r CastSpellRequest) Build(string) CastSpellAction {
	return NewCastSpellAction(r)
}

type CastSpellAction struct {
	cardID            string
	replicateCount    int
	targetsForEffects map[target.EffectTargetKey]target.TargetValue
	spliceCardIDs     []string
	flashback         bool
	autoPayCost       bool
	preactions        []ActivateAbilityRequest // Preactions to be executed before casting the spell
}

func NewCastSpellAction(
	request CastSpellRequest,
) CastSpellAction {
	return CastSpellAction{
		cardID:            request.CardID,
		replicateCount:    request.ReplicateCount,
		targetsForEffects: request.TargetsForEffects,
		spliceCardIDs:     request.SpliceCardIDs,
		flashback:         request.Flashback,
		autoPayCost:       request.AutoPayCost,
		preactions:        request.Preactions,
	}
}

func (a CastSpellAction) Name() string {
	return "Cast Spell"
}

func (a CastSpellAction) Complete(game state.Game, player state.Player, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if a.flashback {
		return a.castWithFlashback(game, player)
	}
	return a.castFromHand(game, player)
}

func (a CastSpellAction) castWithFlashback(game state.Game, player state.Player) ([]event.GameEvent, error) {
	cardToCast, ok := player.GetCardFromZone(a.cardID, mtg.ZoneGraveyard)
	if !ok {
		return nil, fmt.Errorf("player %q does not have card %q in graveyard", player.ID(), a.cardID)
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
			player.ID(),
			cardToCast.ID(),
			ruling.Why(),
		)
	}
	effectWithTargets, err := target.BuildEffectWithTargets(cardToCast.ID(), cardToCast.SpellAbility(), a.targetsForEffects)
	if err != nil {
		return nil, err
	}
	costEvents := pay.PayCost(flashbackCost, cardToCast, player)
	events := append(
		costEvents,
		event.CastSpellEvent{
			PlayerID: player.ID(),
			CardID:   cardToCast.ID(),
			FromZone: mtg.ZoneGraveyard,
		},
		event.PutSpellOnStackEvent{
			PlayerID:          player.ID(),
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
		return nil, fmt.Errorf("player %q does not have card %q in hand", player.ID(), a.cardID)
	}
	var totalCost cost.Cost = cardToCast.ManaCost()
	effectWithTargets, err := target.BuildEffectWithTargets(cardToCast.ID(), cardToCast.SpellAbility(), a.targetsForEffects)
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
	if a.replicateCount > 0 {
		if !judge.CanReplicateCard(
			game,
			player,
			cardToCast,
			&judge.Ruling{Explain: true},
		) {
			return nil, fmt.Errorf(
				"player %q cannot replicate card %q",
				player.ID(),
				cardToCast.ID(),
			)
		}
		replicateCost, ok := cardToCast.GetStaticAbilityCost(mtg.StaticKeywordReplicate)
		if !ok {
			return nil, fmt.Errorf("card %q does not have replicate ability", cardToCast.ID())
		}
		for range a.replicateCount {
			totalCost = cost.CombineCosts(totalCost, replicateCost)
		}
	}
	fooCostEvents := []event.GameEvent{}
	if a.autoPayCost {
		autoPayEvents, err := pay.AutoPay(game, player.ID(), cardToCast, totalCost)
		if err != nil {
			return nil, fmt.Errorf("failed to auto pay cost: %w", err)
		}
		fooCostEvents = append(fooCostEvents, autoPayEvents...)
	} else {
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
				player.ID(),
				cardToCast.ID(),
				ruling.Why(),
			)
		}
		costEvents := pay.PayCost(totalCost, cardToCast, player)
		fooCostEvents = append(fooCostEvents, costEvents...)
	}
	events := append(
		fooCostEvents,
		event.CastSpellEvent{
			PlayerID: player.ID(),
			CardID:   cardToCast.ID(),
			FromZone: mtg.ZoneHand,
		},
		event.PutSpellOnStackEvent{
			PlayerID:          player.ID(),
			CardID:            cardToCast.ID(),
			EffectWithTargets: effectWithTargets,
			FromZone:          mtg.ZoneHand,
		},
	)
	if a.replicateCount > 0 {
		effectWithTargets := []target.EffectWithTarget{{
			EffectSpec: definition.EffectSpec{
				Name:      string(mtg.StaticKeywordReplicate),
				Modifiers: map[string]any{"Count": a.replicateCount},
			},
			Target: target.TargetValue{
				TargetID: cardToCast.ID(),
			},
			SourceID: cardToCast.ID(),
		}}
		events = append(
			events,
			event.PutAbilityOnStackEvent{
				PlayerID: player.ID(),
				// TODO: Maybe the engine needs to manage creating new IDs for game objects.
				// There might be a case where the same card is cast multiple times in the same
				// priority loop and we need to ensure that the IDs are unique.
				SourceID:          cardToCast.ID(),
				FromZone:          mtg.ZoneHand,
				AbilityName:       string(mtg.StaticKeywordReplicate),
				EffectWithTargets: effectWithTargets,
			},
		)
	}
	return events, nil
}

func buildSpliceEffectWithCards(
	player state.Player,
	cardToCast gob.Card,
	spliceCardIDs []string,
	targetsForEffects map[target.EffectTargetKey]target.TargetValue,
) ([]target.EffectWithTarget, []cost.Cost, error) {
	var spliceCosts []cost.Cost
	var effectWithTargets []target.EffectWithTarget
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
		effectWithSpliceTargets, err := target.BuildEffectWithTargets(
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

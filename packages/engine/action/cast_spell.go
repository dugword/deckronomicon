package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/judge"
	"deckronomicon/packages/engine/pay"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/staticability"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/state"
	"fmt"
)

// TODO: Probably needs to include zone for cast from exile, graveyard, etc.
type CastSpellRequest struct {
	CardID            string
	ReplicateCount    int
	SpliceCardIDs     []string
	TargetsForEffects map[effect.EffectTargetKey]target.Target
	Flashback         bool
	AutoPayCost       bool
	AutoPayColors     []mana.Color             // Colors to prioritize when auto-paying costs, if applicable
	Preactions        []ActivateAbilityRequest // Preactions to be executed before casting the spell
	CostTarget        target.Target            // Target for additional costs, if applicable
}

func (r CastSpellRequest) Build(string) CastSpellAction {
	return NewCastSpellAction(r)
}

type CastSpellAction struct {
	cardID            string
	replicateCount    int
	targetsForEffects map[effect.EffectTargetKey]target.Target
	spliceCardIDs     []string
	flashback         bool
	autoPayCost       bool
	autoPayColors     []mana.Color             // Colors to prioritize when auto-paying costs, if applicable
	preactions        []ActivateAbilityRequest // Preactions to be executed before casting the spell
	costTarget        target.Target            // Target for additional costs, if applicable
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
		autoPayColors:     request.AutoPayColors,
		preactions:        request.Preactions,
		costTarget:        request.CostTarget,
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
	staticAbility, ok := cardToCast.StaticAbility(mtg.StaticKeywordFlashback)
	if !ok {
		return nil, fmt.Errorf("card %q does not have flashback ability", cardToCast.ID())
	}
	flashback, ok := staticAbility.(staticability.Flashback)
	if !ok {
		return nil, fmt.Errorf("card %q does not have flashback ability", cardToCast.ID())
	}
	ruling := judge.Ruling{Explain: true}
	if !judge.CanCastSpellWithFlashback(
		game,
		player,
		cardToCast,
		flashback.Cost,
		a.autoPayCost,
		a.autoPayColors,
		&ruling,
	) {
		return nil, fmt.Errorf(
			"player %q cannot cast card %q with flashback: %s",
			player.ID(),
			cardToCast.ID(),
			ruling.Why(),
		)
	}
	effectWithTargets, err := effect.BuildEffectWithTargets(cardToCast.ID(), cardToCast.SpellAbility(), a.targetsForEffects)
	if err != nil {
		return nil, err
	}
	var events []event.GameEvent
	if a.autoPayCost {
		var err error
		activateEvents, err := pay.AutoActivateManaSources(game, flashback.Cost, cardToCast, player.ID(), a.autoPayColors)
		if err != nil {
			return nil, fmt.Errorf("failed to auto-pay cost for card %q: %w", cardToCast.ID(), err)
		}
		events = append(events, activateEvents...)
	}
	costEvents, err := pay.Cost(flashback.Cost, cardToCast, player.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to pay cost for card %q: %w", cardToCast.ID(), err)
	}
	events = append(events, costEvents...)
	events = append(
		events,
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
	effectWithTargets, err := effect.BuildEffectWithTargets(cardToCast.ID(), cardToCast.SpellAbility(), a.targetsForEffects)
	if err != nil {
		return nil, err
	}
	if len(a.spliceCardIDs) > 0 {
		spliceEffectWithCards, spliceCost, err := buildSpliceEffectWithCards(
			player,
			cardToCast,
			a.spliceCardIDs,
			a.targetsForEffects,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to build splice effect with cards: %w", err)
		}
		effectWithTargets = append(effectWithTargets, spliceEffectWithCards...)
		totalCost = cost.NewComposite(totalCost, spliceCost)
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
		staticAbility, ok := cardToCast.StaticAbility(mtg.StaticKeywordReplicate)
		if !ok {
			return nil, fmt.Errorf("card %q does not have replicate ability", cardToCast.ID())
		}
		replicate, ok := staticAbility.(staticability.Replicate)
		if !ok {
			return nil, fmt.Errorf("card %q does not have replicate ability", cardToCast.ID())
		}
		for range a.replicateCount {
			totalCost = cost.NewComposite(totalCost, replicate.Cost)
		}
	}
	ruling := judge.Ruling{Explain: true}
	if !judge.CanCastSpellFromHand(
		game,
		player,
		cardToCast,
		totalCost,
		a.autoPayCost,
		a.autoPayColors,
		&ruling,
	) {
		return nil, fmt.Errorf(
			"player %q cannot cast card %q from from hand: %s",
			player.ID(),
			cardToCast.ID(),
			ruling.Why(),
		)
	}
	var events []event.GameEvent
	if a.autoPayCost {
		var err error
		activateEvents, err := pay.AutoActivateManaSources(game, totalCost, cardToCast, player.ID(), a.autoPayColors)
		if err != nil {
			return nil, fmt.Errorf("failed to auto-pay cost for card %q: %w", cardToCast.ID(), err)
		}
		events = append(events, activateEvents...)
	}
	if cardToCast.AdditionalCost() != nil {
		if costWithTarget, ok := cardToCast.AdditionalCost().(cost.CostWithTarget); ok {
			costWithTarget = costWithTarget.WithTarget(a.costTarget)
			totalCost = cost.NewComposite(totalCost, costWithTarget)
		}
	}
	costEvents, err := pay.Cost(totalCost, cardToCast, player.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to pay cost for card %q: %w", cardToCast.ID(), err)
	}
	events = append(events, costEvents...)
	events = append(
		events,
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
		effectWithTargets := []effect.EffectWithTarget{{
			Effect: effect.Replicate{Count: a.replicateCount},
			Target: target.Target{
				ID: cardToCast.ID(),
			},
			SourceID: cardToCast.ID(),
		}}
		events = append(
			events,
			event.PutAbilityOnStackEvent{
				PlayerID:          player.ID(),
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
	targetsForEffects map[effect.EffectTargetKey]target.Target,
) ([]effect.EffectWithTarget, cost.Cost, error) {
	var spliceCost cost.Cost
	var effectWithTargets []effect.EffectWithTarget
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
		staticAbility, ok := spliceCard.StaticAbility(mtg.StaticKeywordSplice)
		if !ok {
			return nil, nil, fmt.Errorf("card %q does not have splice ability", spliceCardID)
		}
		spliceAbility, ok := staticAbility.(staticability.Splice)
		if !ok {
			return nil, nil, fmt.Errorf("card %q does not have splice ability", spliceCardID)
		}
		spliceCost = cost.NewComposite(spliceCost, spliceAbility.Cost)
		effectWithSpliceTargets, err := effect.BuildEffectWithTargets(
			spliceCard.ID(),
			spliceCard.SpellAbility(),
			targetsForEffects,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to build effect with targets for splice card %q: %w", spliceCardID, err)
		}
		effectWithTargets = append(effectWithTargets, effectWithSpliceTargets...)
	}
	return effectWithTargets, spliceCost, nil
}

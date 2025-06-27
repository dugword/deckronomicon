package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/judge"
	"deckronomicon/packages/engine/pay"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/engine/resolver"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"fmt"
)

type ActivateAbilityRequest struct {
	AbilityID         string
	SourceID          string
	Zone              mtg.Zone
	TargetsForEffects map[effect.EffectTargetKey]effect.Target
}

func (r ActivateAbilityRequest) Build(string) ActivateAbilityAction {
	return NewActivateAbilityAction(r)
}

type ActivateAbilityAction struct {
	abilityID         string
	sourceID          string
	zone              mtg.Zone
	targetsForEffects map[effect.EffectTargetKey]effect.Target
}

func NewActivateAbilityAction(
	request ActivateAbilityRequest,
) ActivateAbilityAction {
	return ActivateAbilityAction{
		abilityID:         request.AbilityID,
		sourceID:          request.SourceID,
		zone:              request.Zone,
		targetsForEffects: request.TargetsForEffects,
	}
}

func (a ActivateAbilityAction) Name() string {
	return "Activate Ability"
}

func (a ActivateAbilityAction) Complete(game state.Game, player state.Player, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	ability, ok := getAbilityOnSourceInZone(
		game,
		player,
		a.sourceID,
		a.abilityID,
		a.zone,
	)
	if !ok {
		return nil, fmt.Errorf(
			"player %q does not have ability %q on source %q in zone %q",
			player.ID(),
			a.abilityID,
			a.sourceID,
			a.zone,
		)
	}
	source, ok := getAbilitySource(
		game,
		player,
		a.zone,
		ability.Source().ID(),
	)
	if !ok {
		return nil, fmt.Errorf(
			"source %q for ability %q not found in zone %q",
			ability.Source().ID(),
			ability.Name(),
			a.zone,
		)
	}
	ruling := judge.Ruling{Explain: true}
	// TODO: Should it be abilityOnObjectInZone.Object() instead of Source()?
	if !judge.CanActivateAbility(game, player, source, ability, &ruling) {
		return nil, fmt.Errorf(
			"player %q cannot activate ability %q on %q in %q: %v",
			player.ID(),
			ability.Name(),
			source.Name(),
			a.zone,
			ruling.Reasons,
		)
	}
	events := []event.GameEvent{
		event.ActivateAbilityEvent{
			PlayerID:  player.ID(),
			SourceID:  source.ID(),
			AbilityID: ability.Name(),
			Zone:      a.zone,
		},
	}
	effectWithTargets, err := effect.BuildEffectWithTargets(ability.ID(), ability.Effects(), a.targetsForEffects)
	if err != nil {
		return nil, fmt.Errorf("failed to build effect with targets: %w", err)
	}
	costEvents := pay.Cost(ability.Cost(), source, player)
	events = append(events, costEvents...)
	var nonAddManaEffectsWithTargets []effect.EffectWithTarget
	for _, effectWithTarget := range effectWithTargets {
		addManaEffect, ok := effectWithTarget.Effect.(effect.AddMana)
		if !ok {
			nonAddManaEffectsWithTargets = append(nonAddManaEffectsWithTargets, effectWithTarget)
			continue
		}
		if source.Match(is.Land()) && cost.HasCostType(ability.Cost(), cost.TapThisCost{}) {
			land, ok := source.(gob.Permanent)
			if !ok {
				return nil, fmt.Errorf("source %q is not a permanent", source.ID())
			}
			events = append(events, event.LandTappedForManaEvent{
				PlayerID: player.ID(),
				ObjectID: land.ID(),
				Subtypes: land.Subtypes(),
			})
		}
		result, err := resolver.ResolveAddMana(game, player.ID(), addManaEffect)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve add mana effect: %w", err)
		}
		events = append(events, result.Events...)
	}
	if len(nonAddManaEffectsWithTargets) > 0 {
		events = append(events, event.PutAbilityOnStackEvent{
			PlayerID:          player.ID(),
			SourceID:          source.ID(),
			AbilityID:         ability.ID(),
			FromZone:          a.zone,
			AbilityName:       ability.Name(),
			EffectWithTargets: nonAddManaEffectsWithTargets,
		})
	}
	return events, nil
}

func getAbilityOnSourceInZone(
	game state.Game,
	player state.Player,
	sourceID string,
	abilityID string,
	zone mtg.Zone,
) (gob.Ability, bool) {
	switch zone {
	case mtg.ZoneBattlefield:
		return getAbilityFromPermanent(game, sourceID, abilityID)
	case mtg.ZoneHand, mtg.ZoneGraveyard, mtg.ZoneExile, mtg.ZoneLibrary:
		return getAbilityFromCardInZone(game, player, sourceID, abilityID, zone)
	default:
		return gob.Ability{}, false
	}
}

func getAbilityFromPermanent(
	game state.Game,
	sourceID string,
	abilityID string,
) (gob.Ability, bool) {
	permanent, ok := game.Battlefield().Get(sourceID)
	if !ok {
		return gob.Ability{}, false
	}
	for _, ability := range permanent.ActivatedAbilities() {
		if ability.ID() == abilityID {
			return ability, true
		}
	}
	return gob.Ability{}, false
}

func getAbilityFromCardInZone(
	game state.Game,
	player state.Player,
	sourceID string,
	abilityID string,
	zone mtg.Zone,
) (gob.Ability, bool) {
	card, ok := player.GetCardFromZone(sourceID, zone)
	if !ok {
		return gob.Ability{}, false
	}
	for _, ability := range card.ActivatedAbilities() {
		if ability.ID() == abilityID {
			return ability, true
		}
	}
	return gob.Ability{}, false
}

func getAbilitySource(
	game state.Game,
	player state.Player,
	zone mtg.Zone,
	sourceID string,
) (gob.Object, bool) {
	switch zone {
	case mtg.ZoneBattlefield:
		return getAbilitySourceFromPermanent(game, sourceID)
	default:
		return getAbilitySourceFromCardInZone(player, sourceID, zone)
	}
}

func getAbilitySourceFromPermanent(
	game state.Game,
	sourceID string,
) (gob.Object, bool) {
	permanent, ok := game.Battlefield().Get(sourceID)
	if !ok {
		return nil, false
	}
	return permanent, true
}

func getAbilitySourceFromCardInZone(
	player state.Player,
	sourceID string,
	zone mtg.Zone,
) (gob.Object, bool) {
	card, ok := player.GetCardFromZone(sourceID, zone)
	if !ok {
		return nil, false
	}
	return card, true
}

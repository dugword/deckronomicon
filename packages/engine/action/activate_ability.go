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
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"fmt"
)

type ActivateAbilityRequest struct {
	AbilityID         string
	SourceID          string
	Zone              mtg.Zone
	TargetsForEffects map[effect.EffectTargetKey]target.Target
	AutoPayCost       bool
	AutoPayColors     []mana.Color // Colors to prioritize when auto-paying costs, if applicable
}

func (r ActivateAbilityRequest) Build() ActivateAbilityAction {
	return NewActivateAbilityAction(r)
}

type ActivateAbilityAction struct {
	abilityID         string
	sourceID          string
	zone              mtg.Zone
	targetsForEffects map[effect.EffectTargetKey]target.Target
	autoPayCost       bool
	autoPayColors     []mana.Color // Colors to prioritize when auto-paying costs, if applicable
}

func NewActivateAbilityAction(
	request ActivateAbilityRequest,
) ActivateAbilityAction {
	return ActivateAbilityAction{
		abilityID:         request.AbilityID,
		sourceID:          request.SourceID,
		zone:              request.Zone,
		targetsForEffects: request.TargetsForEffects,
		autoPayCost:       request.AutoPayCost,
		autoPayColors:     request.AutoPayColors,
	}
}

func (a ActivateAbilityAction) Name() string {
	return "Activate Ability"
}

func (a ActivateAbilityAction) Complete(game *state.Game, playerID string, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	ability, ok := getAbilityOnSourceInZone(
		game,
		playerID,
		a.sourceID,
		a.abilityID,
		a.zone,
	)
	if !ok {
		return nil, fmt.Errorf(
			"player %q does not have ability %q on source %q in zone %q",
			playerID,
			a.abilityID,
			a.sourceID,
			a.zone,
		)
	}
	source, ok := getAbilitySource(
		game,
		playerID,
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
	if !judge.CanActivateAbility(
		game,
		playerID,
		source,
		ability,
		a.autoPayCost,
		a.autoPayColors,
		&ruling,
	) {
		return nil, fmt.Errorf(
			"player %q cannot activate ability %q on %q in %q: %v",
			playerID,
			ability.Name(),
			source.Name(),
			a.zone,
			ruling.Reasons,
		)
	}

	var events []event.GameEvent
	if a.autoPayCost {
		var err error
		activateEvents, err := pay.AutoActivateManaSources(game, ability.Cost(), ability, playerID, a.autoPayColors)
		if err != nil {
			return nil, fmt.Errorf("failed to auto-pay cost for ability %q: %w", ability.ID(), err)
		}
		events = append(events, activateEvents...)
	}
	costEvents, err := pay.Cost(ability.Cost(), source, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to pay ability cost: %w", err)
	}
	events = append(events, costEvents...)
	events = append(events,
		&event.ActivateAbilityEvent{
			PlayerID:  playerID,
			SourceID:  source.ID(),
			AbilityID: ability.Name(),
			Zone:      a.zone,
		},
	)
	var nonAddManaEffectsWithTargets []*effect.EffectWithTarget
	effectWithTargets, err := effect.BuildEffectWithTargets(ability.ID(), ability.Effects(), a.targetsForEffects)
	if err != nil {
		return nil, fmt.Errorf("failed to build effect with targets: %w", err)
	}
	for _, effectWithTarget := range effectWithTargets {
		addManaEffect, ok := effectWithTarget.Effect.(*effect.AddMana)
		if !ok {
			nonAddManaEffectsWithTargets = append(nonAddManaEffectsWithTargets, effectWithTarget)
			continue
		}
		if source.Match(is.Land()) && cost.HasType(ability.Cost(), cost.TapThis{}) {
			land, ok := source.(*gob.Permanent)
			if !ok {
				return nil, fmt.Errorf("source %q is not a permanent", source.ID())
			}
			// TODO: Figure out how to move this to generate trigger events middleware
			events = append(events, &event.LandTappedForManaEvent{
				PlayerID: playerID,
				ObjectID: land.ID(),
				Subtypes: land.Subtypes(),
			})
		}
		result, err := resolver.ResolveAddMana(game, playerID, addManaEffect)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve add mana effect: %w", err)
		}
		events = append(events, result.Events...)
	}
	if len(nonAddManaEffectsWithTargets) > 0 {
		events = append(events, &event.PutAbilityOnStackEvent{
			PlayerID:          playerID,
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
	game *state.Game,
	playerID string,
	sourceID string,
	abilityID string,
	zone mtg.Zone,
) (*gob.Ability, bool) {
	switch zone {
	case mtg.ZoneBattlefield:
		return getAbilityFromPermanent(game, sourceID, abilityID)
	case mtg.ZoneHand, mtg.ZoneGraveyard, mtg.ZoneExile, mtg.ZoneLibrary:
		return getAbilityFromCardInZone(game, playerID, sourceID, abilityID, zone)
	default:
		return nil, false
	}
}

func getAbilityFromPermanent(
	game *state.Game,
	sourceID string,
	abilityID string,
) (*gob.Ability, bool) {
	permanent, ok := game.Battlefield().Get(sourceID)
	if !ok {
		return nil, false
	}
	for _, ability := range permanent.ActivatedAbilities() {
		if ability.ID() == abilityID {
			return ability, true
		}
	}
	return nil, false
}

func getAbilityFromCardInZone(
	game *state.Game,
	playerID string,
	sourceID string,
	abilityID string,
	zone mtg.Zone,
) (*gob.Ability, bool) {
	player := game.GetPlayer(playerID)
	card, ok := player.GetCardFromZone(sourceID, zone)
	if !ok {
		return nil, false
	}
	for _, ability := range card.ActivatedAbilities() {
		if ability.ID() == abilityID {
			return ability, true
		}
	}
	return nil, false
}

func getAbilitySource(
	game *state.Game,
	playerID string,
	zone mtg.Zone,
	sourceID string,
) (gob.Object, bool) {
	switch zone {
	case mtg.ZoneBattlefield:
		return getAbilitySourceFromPermanent(game, sourceID)
	default:
		return getAbilitySourceFromCardInZone(game, playerID, sourceID, zone)
	}
}

func getAbilitySourceFromPermanent(
	game *state.Game,
	sourceID string,
) (gob.Object, bool) {
	permanent, ok := game.Battlefield().Get(sourceID)
	if !ok {
		return nil, false
	}
	return permanent, true
}

func getAbilitySourceFromCardInZone(
	game *state.Game,
	playerID string,
	sourceID string,
	zone mtg.Zone,
) (gob.Object, bool) {
	player := game.GetPlayer(playerID)
	card, ok := player.GetCardFromZone(sourceID, zone)
	if !ok {
		return nil, false
	}
	return card, true
}

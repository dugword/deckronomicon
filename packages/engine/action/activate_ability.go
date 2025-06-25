package action

import (
	buildmanaabilities "deckronomicon/packages/build_mana_abilities"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/judge"
	"deckronomicon/packages/engine/pay"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"fmt"
)

type ActivateAbilityRequest struct {
	AbilityID         string
	SourceID          string
	Zone              mtg.Zone
	TargetsForEffects map[target.EffectTargetKey]target.TargetValue
}

func (r ActivateAbilityRequest) Build(string) ActivateAbilityAction {
	return NewActivateAbilityAction(r)
}

type ActivateAbilityAction struct {
	abilityID         string
	sourceID          string
	zone              mtg.Zone
	targetsForEffects map[target.EffectTargetKey]target.TargetValue
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
	effectWithTargets, err := target.BuildEffectWithTargets(ability.ID(), ability.EffectSpecs(), a.targetsForEffects)
	if err != nil {
		return nil, fmt.Errorf("failed to build effect with targets: %w", err)
	}
	costEvents := pay.PayCost(ability.Cost(), source, player)
	events = append(events, costEvents...)
	// TODO: I think I need to do this on a per effect basis, like I think for
	// chromatic star I get the mana fixing buy the draw a card effect still
	// goes on the stack.
	if ability.IsManaAbility() {
		if permanent, ok := source.(gob.Permanent); ok {
			if permanent.Match(is.Land()) {
				if cost.HasCostType(ability.Cost(), cost.TapThisCost{}) {
					events = append(events, event.LandTappedForManaEvent{
						PlayerID: player.ID(),
						ObjectID: permanent.ID(),
						Subtypes: permanent.Subtypes(),
					})
				}
			}
		}
		manaEvents, err := buildmanaabilities.BuildManaAbilityEvents(game, player, effectWithTargets, resEnv)
		if err != nil {
			return nil, fmt.Errorf("failed to build mana ability events: %w", err)
		}
		events = append(events, manaEvents...)
		return events, nil
	}
	events = append(events, event.PutAbilityOnStackEvent{
		PlayerID:          player.ID(),
		SourceID:          source.ID(),
		AbilityID:         ability.ID(),
		FromZone:          a.zone,
		AbilityName:       ability.Name(),
		EffectWithTargets: effectWithTargets,
	})
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
) (query.Object, bool) {
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
) (query.Object, bool) {
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
) (query.Object, bool) {
	card, ok := player.GetCardFromZone(sourceID, zone)
	if !ok {
		return nil, false
	}
	return card, true
}

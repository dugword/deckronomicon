package gob

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"fmt"
	"strconv"
)

type Trigger struct {
	EventType   string
	SelfTrigger bool
	Filter      query.Opts
}

type TriggeredAbility struct {
	id      string
	name    string
	trigger Trigger
	effects []effect.Effect
	source  Object
	zone    mtg.Zone
}

func NewTriggeredAbility(source Object, idx int, abilityDefinition *definition.TriggeredAbility) (*TriggeredAbility, error) {
	zone := mtg.ZoneBattlefield
	specZone, ok := mtg.StringToZone(abilityDefinition.Zone)
	if ok {
		zone = specZone
	}
	var effects []effect.Effect
	for _, effectDefinition := range abilityDefinition.Effects {
		effect, err := effect.New(effectDefinition)
		if err != nil {
			return nil, fmt.Errorf("failed to build effect %s: %w", effectDefinition.Name, err)
		}
		effects = append(effects, effect)
	}

	ability := TriggeredAbility{
		name:    abilityDefinition.Name,
		effects: effects,
		id:      fmt.Sprintf("%s-%d", source.ID(), idx+1),
		zone:    zone,
		source:  source,
		trigger: Trigger{
			EventType: abilityDefinition.Trigger.Event,
			Filter:    buildOpts(abilityDefinition.Trigger.Filter),
		},
	}
	if len(abilityDefinition.Trigger.Filter) == 0 {
		ability.trigger.SelfTrigger = true
	}
	return &ability, nil
}

func (a *TriggeredAbility) Effects() []effect.Effect {
	return a.effects
}

func (a *TriggeredAbility) ID() string {
	return a.id
}

func (a *TriggeredAbility) Name() string {
	return a.name
}

func (a *TriggeredAbility) Trigger() Trigger {
	return a.trigger
}

func (a *TriggeredAbility) Zone() mtg.Zone {
	return a.zone
}

func (a *TriggeredAbility) Source() Object {
	return a.source
}

func buildOpts(
	input map[string][]string,
) query.Opts {
	var opts query.Opts
	for _, v := range input["CardTypes"] {
		cardType, ok := mtg.StringToCardType(v)
		if !ok {
			panic("invalid card type: " + v)
		}
		opts.CardTypes = append(opts.CardTypes, cardType)
	}
	for _, v := range input["Colors"] {
		color, ok := mtg.StringToColor(v)
		if !ok {
			panic("invalid color: " + v)
		}
		opts.Colors = append(opts.Colors, color)
	}
	for _, v := range input["Subtypes"] {
		subtype, ok := mtg.StringToSubtype(v)
		if !ok {
			panic("invalid subtype: " + v)
		}
		opts.Subtypes = append(opts.Subtypes, subtype)
		for _, v := range input["ManaValues"] {
			manaValue, err := strconv.Atoi(v)
			if err != nil {
				panic("invalid mana value: " + v)
			}
			opts.ManaValues = append(opts.ManaValues, manaValue)
		}
	}
	return opts
}

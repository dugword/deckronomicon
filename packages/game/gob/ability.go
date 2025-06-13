package gob

import (
	//	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"fmt"
	"strings"
)

// Ability represents abilities that require activation costs.
type Ability struct {
	name string
	// TODO Maybe parse this into a cost type?
	cost string
	// effects     []Effect
	// effectSpecs []definition.EffectSpec
	effects []definition.EffectSpec
	id      string
	zone    mtg.Zone
	source  query.Object
	speed   mtg.Speed
}

func NewAbility(id string) Ability {
	ability := Ability{
		id: id,
	}
	return ability
}

func (a Ability) Controller() string {
	return a.source.Controller()
}

func (a Ability) Owner() string {
	return a.source.Owner()
}

func (a Ability) Cost() string {
	return a.cost
}

func (a Ability) Effects() []definition.EffectSpec {
	return a.effects
}

/*
func (a Ability) EffectSpecs() []definition.EffectSpec {
	return a.effectSpecs
}
*/

func (a Ability) Speed() mtg.Speed {
	return a.speed
}

func (a Ability) Source() query.Object {
	return a.source
}

func (a Ability) Zone() mtg.Zone {
	return a.zone
}

// BuildActivatedAbility builds an activated ability from the given
// specification.
// TODO use a better interface
/*
func BuildActivatedAbility(state core.State, spec definition.ActivatedAbilitySpec, source query.Object) (*Ability, error) {
	ability := Ability{
		name:   spec.Name,
		id:     state.GetNextID(),
		source: source,
		Zone:   "",
		Speed:  spec.Speed,
	}
	cost, err := cost.New(spec.Cost, source)
	if err != nil {
		return nil, fmt.Errorf("failed to create cost: %w", err)
	}
	ability.Cost = cost
	for _, _ = range spec.EffectSpecs {
			effect, err := effectimpl.BuildEffect(source, effectSpec)
			if err != nil {
				return nil, fmt.Errorf("failed to create effect: %w", err)
			}
			ability.Effects = append(ability.Effects, effect)
	}
	return &ability, nil
}
*/

func (a Ability) Name() string {
	//return fmt.Sprintf("%s - %s", a.source.Name(), a.name)
	return a.name
}

func (a Ability) ID() string {
	return a.id
}

/*
func (a Ability) CanActivate(state core.State, playerID string) bool {
	if a.Speed == mtg.SpeedSorcery && !state.CanCastSorcery(playerID) {
		return false
	}
	return true
}
*/

// Description returns a string representation of the activated ability.
func (a Ability) Description() string {
	var descriptions []string
	for _, effect := range a.effects {
		// TODO: Come up with a better way to handle descriptions
		descriptions = append(descriptions, effect.Name)
		//descriptions = append(descriptions, effect.Description())
	}
	// return fmt.Sprintf("%s: %s", a.Cost.Description(), strings.Join(descriptions, ", "))
	return fmt.Sprintf("%s: %s", "<cost>", strings.Join(descriptions, ", "))
}

// IsManaAbility checks if the activated ability is a mana ability.
func (a Ability) IsManaAbility() bool {
	fmt.Println("Checking if ability is mana ability:", a.Name())
	for _, tag := range a.Tags() {
		fmt.Printf("Tag: %s=%s\n", tag.Key, tag.Value)
		if tag.Key == TagManaAbility {
			return true
		}
	}
	return false
}

func (a Ability) Match(predicate query.Predicate) bool {
	return predicate(a)
}

// Resolve resolves the activated ability. Any costs must be paid before
// resolving the ability.
/*
func (a Ability) Resolve(state core.State, player core.Player) error {
	for _, effect := range a.Effects {
		if err := effect.Apply(state, player); err != nil {
			return fmt.Errorf("cannot resolve effect: %w", err)
		}
	}
	return nil
}
*/

// Tags returns the tags associated with the activated ability.
func (a Ability) Tags() []Tag {
	var tags []Tag
	for _, effect := range a.effects {
		for _, modifier := range effect.Modifiers {
			tags = append(tags, Tag{
				Key:   modifier.Key,
				Value: modifier.Value,
			})
		}
	}
	return tags
}

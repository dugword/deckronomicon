package game

import (
	"fmt"
	"strings"
)

// ActivatedAbility represents abilities that require activation costs.
type ActivatedAbility struct {
	Cost    Cost
	Effects []*Effect
	id      string
	Zone    string
	source  GameObject
	Speed   string
}

func (a *ActivatedAbility) HasSubtype(Subtype) bool {
	// Activated abilities do not have subtypes.
	return false
}
func (a *ActivatedAbility) Name() string {
	// Activated abilities do not have a name.
	// TODO: Maybe make this the source name + something
	// or the EffectID, or maybe Name: join(EfectID, ", ")
	var ids []string
	for _, effect := range a.Effects {
		ids = append(ids, effect.ID)
	}
	return fmt.Sprintf("%s: %s", a.source.Name(), strings.Join(ids, ", "))
}
func (a *ActivatedAbility) ActivatedAbilities() []*ActivatedAbility {
	// Activated abilities do not have activated abilities.
	return nil
}
func (a *ActivatedAbility) StaticAbilities() []*StaticAbility {
	// Activated abilities do not have static abilities.
	return nil
}
func (a *ActivatedAbility) HasStaticAbility(string) bool {
	// Activated abilities do not have static abilities.
	return false
}
func (a *ActivatedAbility) ID() string {
	// Activated abilities have an ID.
	return a.id
}
func (a *ActivatedAbility) ManaValue() int {
	// Activated abilities do not have a mana value.
	return 0
}

func (a *ActivatedAbility) CanPlay(state *GameState) bool {
	if a.Speed == SpeedSorcery && !state.CanCastSorcery() {
		return false
	}
	return true
}

// BuildActivatedAbility builds an activated ability from the given
// specification.
func BuildActivatedAbility(spec ActivatedAbilitySpec, source GameObject) (*ActivatedAbility, error) {
	// ZoneBattlefield is the default zone for activated abilities.
	zone := ZoneBattlefield
	if spec.Zone != "" {
		zone = spec.Zone
	}
	ability := ActivatedAbility{
		id:     GetNextID(),
		source: source,
		Zone:   zone,
		Speed:  spec.Speed,
	}
	cost, err := NewCost(spec.Cost, source)
	if err != nil {
		return nil, fmt.Errorf("failed to create cost: %w", err)
	}
	ability.Cost = cost
	for _, effectSpec := range spec.EffectSpecs {
		effect, err := BuildEffect(source, effectSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to create effect: %w", err)
		}
		ability.Effects = append(ability.Effects, effect)
	}
	return &ability, nil
}

// Description returns a string representation of the activated ability.
func (a *ActivatedAbility) Description() string {
	var descriptions []string
	for _, effect := range a.Effects {
		descriptions = append(descriptions, effect.Description)
	}
	return fmt.Sprintf("%s: %s", a.Cost.Description(), strings.Join(descriptions, ", "))
}

// IsManaAbility checks if the activated ability is a mana ability.
func (a *ActivatedAbility) IsManaAbility() bool {
	for _, tag := range a.Tags() {
		if tag.Key == AbilityTagManaAbility {
			return true
		}
	}
	return false
}

// Resolve resolves the activated ability. Any costs must be paid before
// resolving the ability.
func (a *ActivatedAbility) Resolve(state *GameState, player *Player) error {
	for _, effect := range a.Effects {
		if err := effect.Apply(state, player); err != nil {
			return fmt.Errorf("cannot resolve effect: %w", err)
		}
	}
	return nil
}

// Tags returns the tags associated with the activated ability.
func (a *ActivatedAbility) Tags() []EffectTag {
	var tags []EffectTag
	for _, effect := range a.Effects {
		tags = append(tags, effect.Tags...)
	}
	return tags
}

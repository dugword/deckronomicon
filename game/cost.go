package game

import (
	"errors"
	"fmt"
)

type ManaCost struct {
	Generic int            // The {1}, {2}, etc.object.
	Colors  map[string]int // {"G": 1, "U": 1}
}

type Cost interface {
	CanPay(game *GameState, source *Permanent) bool
	Pay(game *GameState, source *Permanent) error
	//Description() string
}

// --- CompositeCost ---

type CompositeCost struct {
	Costs []Cost
}

// --- TapCost ---

type TapCost struct{}

// this feels odd since the cost is tied to the ability, and the ability is tied to a source... so weird you have to pass it in and can't read up the chain
func (c TapCost) CanPay(game *GameState, source *Permanent) bool {
	return !source.IsTapped()
}

func (c TapCost) Pay(game *GameState, source *Permanent) error {
	// TODO: Either return this directly, or wrap the function, probably wrap with a pay thing
	if err := source.Tap(); err != nil {
		return err
	}
	return nil
}

/*
// TODO Needs to be abstract interface or something
// Also the logic is hard, maybe it should be its own type with methods
func (c ManaCost) CanPay(game *GameState, source *Object) bool {
	tempPool := map[string]int{}
	maps.Copy(tempPool, game.ManaPool)
	fmt.Println("checking colors")
	for color, amount := range source.ManaCost.Colors {
		if tempPool[color] >= amount {
			tempPool[color] -= amount
		} else {
			return false
		}
	}
	tempGenericMana := 0
	fmt.Println("checking generic")
	for _, amount := range tempPool {
		tempGenericMana += amount
	}
	if tempGenericMana < source.ManaCost.Generic {
		return false
	}
	return true
}

// TODO: I think this is managed in a few different places, probably should centralize it
func CanPotentiallyPayFor(state *GameState, manaCost ManaCost) bool {
	simulated := GetPotentialMana(state)

	for color, need := range manaCost.Colors {
		if simulated[color] < need {
			return false
		}
		simulated[color] -= need
	}

	genericAvailable := 0
	for _, v := range simulated {
		genericAvailable += v
	}

	return genericAvailable >= manaCost.Generic
}
*/

func ChooseManaForGeneric(genericCost int, pool map[string]int, resolver ChoiceResolver) (map[string]int, error) {
	choice := make(map[string]int)
	remaining := genericCost

	for remaining > 0 {
		choices := []Choice{}
		for color, count := range pool {
			if count > 0 {
				choices = append(choices, Choice{
					Name: color,
				})
			}
		}
		if len(choices) == 0 {
			return nil, errors.New("not enough mana to pay for generic cost")
		}

		selected := resolver.ChooseOne("Choose mana to use for generic cost", choices)
		pool[selected.Name]--
		choice[selected.Name]++
		remaining--
	}

	return choice, nil
}

func (c ManaCost) Pay(game *GameState, resolver ChoiceResolver, source GameObject) error {
	// Step 1: Pay colored mana first
	for color, amount := range c.Colors {
		if game.ManaPool[color] < amount {
			return fmt.Errorf("not enough %s mana to pay colored cost", color)
		}
		game.ManaPool[color] -= amount
	}

	// Step 2: Pay generic mana
	if c.Generic > 0 {
		choice, err := ChooseManaForGeneric(c.Generic, game.ManaPool, resolver)
		if err != nil {
			return err
		}
		// Subtract the chosen mana from the pool
		for color, amt := range choice {
			game.ManaPool[color] -= amt
		}
	}
	return nil
}

/*
func (c ManaCost) Description() string {
	return "Pay mana"
}






func (c TapCost) Description() string {
	return "Tap this permanent"
}

// --- SacrificeCost ---

type SacrificeCost struct {
	Filter func(*Permanent) bool
}

func (c SacrificeCost) CanPay(game *GameState, source *Permanent) bool {
	for _, p := range game.Battlefield {
		if c.Filter(p) {
			return true
		}
	}
	return false
}

func (c SacrificeCost) Pay(game *GameState, source *Permanent) error {
	for i, p := range game.Battlefield {
		if c.Filter(p) {
			game.Battlefield = append(game.Battlefield[:i], game.Battlefield[i+1:]...)
			// TODO

				//if p.Card != nil {
				//	game.Graveyard = append(game.Graveyard, *p.Card)
				//}
			//
			return nil
		}
	}
	return fmt.Errorf("no valid permanent to sacrifice")
}

func (c SacrificeCost) Description() string {
	return "Sacrifice a permanent"
}



func (c CompositeCost) CanPay(game *GameState, source *Permanent) bool {
	for _, cost := range c.Costs {
		if !cost.CanPay(game, source) {
			return false
		}
	}
	return true
}

func (c CompositeCost) Pay(game *GameState, source *Permanent) error {
	for _, cost := range c.Costs {
		if err := cost.Pay(game, source); err != nil {
			return err
		}
	}
	return nil
}

func (c CompositeCost) Description() string {
	return "Multiple costs"
}

// --- ParseCost ---

// TODO: This needs more thought
/*
func ParseCost(input string) Cost {
	parts := strings.Split(input, ",")
	var costs []Cost
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		switch trimmed {
		case "{T}":
			costs = append(costs, TapCost{})
		case "Sacrifice a Clue":
			costs = append(costs, SacrificeCost{
				Filter: func(p *Permanent) bool {
					return p.HasType("Artifact") && p.HasSubtype("Clue")
				},
			})
		default:
			if strings.HasPrefix(trimmed, "{") {
				mana := make(map[string]int)
				tokens := strings.Split(trimmed, "}{")
				for _, tok := range tokens {
					tok = strings.Trim(tok, "{}")
					mana[tok]++
				}
				costs = append(costs, ManaCost{Mana: mana})
			}
		}
	}

	if len(costs) == 1 {
		return costs[0]
	}
	return CompositeCost{Costs: costs}
}
*/

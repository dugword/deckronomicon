package game

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Cost represents a cost that can be paid in the game.
type Cost interface {
	CanPay(*GameState, *Player) bool
	Description() string
	Pay(*GameState, *Player) error
	Add(Cost) Cost
}

// ManaPattern is a regex pattern that matches valid mana costs.
// TODO: Support X costs and other special cases.
var ManaPattern = regexp.MustCompile(`^(?:\{[0-9WUBRGC]+\})*$`)

// NewCost creates a new cost based on the input string and the source.
// TODO: Maybe rename to NewCost, only return a composit cost when there's
// more than one
func NewCost(input string, source GameObject) (Cost, error) {
	parts := strings.Split(input, ",")
	var costs []Cost
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		switch {
		case isTapCost(trimmed):
			permanent, ok := source.(*Permanent)
			if !ok {
				return nil, fmt.Errorf("source %s not permanent", source.Name())
			}
			costs = append(costs, &TapCost{Permanent: permanent})
		case isMana(trimmed):
			manaCost, err := ParseManaCost(trimmed)
			if err != nil {
				return nil, fmt.Errorf("failed to parse mana cost %s: %w", trimmed, err)
			}
			costs = append(costs, manaCost)
		default:
			return nil, fmt.Errorf("unknown cost %s", trimmed)
		}
	}
	if len(costs) == 1 {
		return costs[0], nil
	}
	return &CompositeCost{Costs: costs}, nil
}

// CompositeCost represents a cost that is a combination of multiple costs.
type CompositeCost struct {
	Costs []Cost
}

func (c *CompositeCost) Add(cost Cost) Cost {
	c.Costs = append(c.Costs, cost)
	return c
}

// CanPay checks if all costs in the composite cost can be paid with the
// current game state.
func (c *CompositeCost) CanPay(state *GameState, player *Player) bool {
	for _, cost := range c.Costs {
		if !cost.CanPay(state, player) {
			return false
		}
	}
	return true
}

// Description returns a string representation of the composite cost.
func (c *CompositeCost) Description() string {
	// Cost ordered as Mana, Tap, Sacrifice
	var costStrings []string
	// Mana
	for _, cost := range c.Costs {
		if _, ok := cost.(*ManaCost); ok {
			costStrings = append(costStrings, cost.Description())
		}
	}
	// Tap
	for _, cost := range c.Costs {
		if _, ok := cost.(*TapCost); ok {
			costStrings = append(costStrings, cost.Description())
		}
	}
	//Sacrifice
	for _, cost := range c.Costs {
		if _, ok := cost.(*SacrificeCost); ok {
			costStrings = append(costStrings, cost.Description())
		}
	}
	return strings.Join(costStrings, ", ")
}

// Pay pays all costs in the composite cost.
// TODO: If one cost fails, we need to roll back the others.
func (c *CompositeCost) Pay(state *GameState, player *Player) error {
	// TODO: Maybe there's a better way to do this, but this helps with the
	// needing to roll back thing.
	for _, cost := range c.Costs {
		if !cost.CanPay(state, player) {
			return fmt.Errorf(
				"failed to pay composite cost",
			)
		}
	}
	for _, cost := range c.Costs {
		if err := cost.Pay(state, player); err != nil {
			return fmt.Errorf(
				"failed to pay composite cost: %w",
				err,
			)
		}
	}
	return nil
}

// ManaCost represents the mana cost of a card or ability.
type ManaCost struct {
	Colors  map[Color]int
	Generic int
}

// AddCost adds a cost to the mana cost.
func (c *ManaCost) Add(cost Cost) Cost {
	cc := CompositeCost{
		Costs: []Cost{},
	}
	cc.Costs = append(cc.Costs, c)
	cc.Costs = append(cc.Costs, cost)
	return &cc
}

// CanPay checks if the cost can be paid with the current game state.
// TODO Maybe this should just be *Player
func (c *ManaCost) CanPay(state *GameState, player *Player) bool {
	tempPool := player.ManaPool.Copy()
	for color, amount := range c.Colors {
		if !tempPool.Has(color, amount) {
			return false
		}
		tempPool.Use(color, amount)
	}
	if !tempPool.HasGeneric(c.Generic) {
		return false
	}
	return true
}

// Description returns a string representation of the mana cost.
func (c *ManaCost) Description() string {
	// TODO: colors might not be the right name, costs? Symbols?
	var colors []string
	if c.Generic > 0 {
		colors = append(colors, fmt.Sprintf("{%d}", c.Generic))
	}
	for _, color := range []Color{
		// Do not reorder this list, follows the color wheel standard for mana
		// costs.
		ColorColorless,
		ColorWhite,
		ColorBlue,
		ColorBlack,
		ColorRed,
		ColorGreen,
	} {
		if c.Colors[color] > 0 {
			for range c.Colors[color] {
				colors = append(colors, fmt.Sprintf("{%s}", color))
			}
		}
	}
	return strings.Join(colors, "")
}

// Pay pays the mana cost by using the mana from the mana pool.
// TODO: Need to roll back if the cost is partially paid and fails.
func (c *ManaCost) Pay(game *GameState, player *Player) error {
	// Pay colored mana
	for color, amount := range c.Colors {
		if err := player.ManaPool.Use(color, amount); err != nil {
			return err
		}
	}
	// Pay generic mana
	if c.Generic > 0 {
		choice, err := chooseManaForGeneric(c.Generic, player)
		if err != nil {
			return err
		}
		for colorChoice, amount := range choice {
			color, err := StringToColor(colorChoice)
			if err != nil {
				return err
			}
			player.ManaPool.Use(color, amount)
		}
	}
	return nil
}

// SacrificeCost represents a cost that requires sacrificing the permanent.
// TODO: Support sacrificing other permanents.
type SacrificeCost struct {
	Permanent *Permanent
}

// AddCost adds a cost to the mana cost.
func (c *SacrificeCost) Add(cost Cost) Cost {
	cc := CompositeCost{
		Costs: []Cost{},
	}
	cc.Costs = append(cc.Costs, c)
	cc.Costs = append(cc.Costs, cost)
	return &cc
}

// CanPay checks if the sacrifice cost can be paid with the current game
// state.
func (c *SacrificeCost) CanPay(game *GameState, player *Player) bool {
	// TODO: Pretty much always true unless there is some state that says the
	// permanent can't be sacrificed.
	return true
}

// Description returns a string representation of the sacrifice cost.
func (c *SacrificeCost) Description() string {
	return "sacrifice this permanent"
}

// Pay pays the sacrifice cost by sacrificing the permanent.
func (c *SacrificeCost) Pay(game *GameState, player *Player) error {
	// TODO: Implement this
	return errors.New("not implemented")
}

// TapCost represents a cost that requires tapping the permanent.
type TapCost struct {
	Permanent *Permanent
}

// AddCost adds a cost to the mana cost.
func (c *TapCost) Add(cost Cost) Cost {
	cc := CompositeCost{
		Costs: []Cost{},
	}
	cc.Costs = append(cc.Costs, c)
	cc.Costs = append(cc.Costs, cost)
	return &cc
}

// CanPay checks if the tap cost can be paid with the current game state.
func (c *TapCost) CanPay(game *GameState, player *Player) bool {
	return !c.Permanent.IsTapped()
}

// Description returns a string representation of the tap cost.
func (c *TapCost) Description() string {
	return "{T}"
}

// Pay pays the tap cost by tapping the permanent.
func (c *TapCost) Pay(game *GameState, player *Player) error {
	if err := c.Permanent.Tap(); err != nil {
		return fmt.Errorf("failed to pay {T} cost: %w", err)
	}
	return nil
}

// ParseManaCost parses a mana cost string and returns a ManaCost.
func ParseManaCost(costStr string) (*ManaCost, error) {
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(costStr, -1)
	manaCost := ManaCost{
		Generic: 0,
		Colors:  map[Color]int{},
	}
	for _, match := range matches {
		symbol := strings.ToUpper(match[1])
		if val, err := strconv.Atoi(symbol); err == nil {
			manaCost.Generic += val
		} else {
			color, err := StringToColor(symbol)
			if err != nil {
				return nil, err
			}
			manaCost.Colors[color]++
		}
	}
	return &manaCost, nil
}

// chooseManaForGeneric prompts the user to choose mana for a generic cost.
func chooseManaForGeneric(genericCost int, player *Player) (map[string]int, error) {
	choice := make(map[string]int)
	remaining := genericCost
	for remaining > 0 {
		choices := []Choice{}
		for _, color := range player.ManaPool.ColorsAvailable() {
			choices = append(choices, Choice{
				Name: string(color),
			})
		}
		if len(choices) == 0 {
			return nil, errors.New("insufficient mana")
		}
		selected, err := player.Agent.ChooseOne(
			"Choose mana to use for generic cost",
			NewChoiceSource("Pay Generic Cost", "Pay Generic Cost"),
			choices,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to choose mana for generic cost: %w", err)
		}
		// TODO: maybe make this a function
		color, err := StringToColor(selected.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to apply choice for generic mana: %w", err)
		}
		player.ManaPool.Use(color, 1)
		choice[selected.Name]++
		remaining--
	}
	return choice, nil
}

// isManaCost checks if the input string is a valid mana cost.
func isMana(input string) bool {
	return ManaPattern.MatchString(input)
}

// isTapCost checks if the input string is a tap cost.
func isTapCost(input string) bool {
	return input == "{T}"
}

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
}

func AddCosts(costs ...Cost) Cost {
	// If there's only one cost, return it directly.
	if len(costs) == 1 {
		return costs[0]
	}
	// Otherwise, return a composite cost.
	return &CompositeCost{Costs: costs}
}

// ManaPattern is a regex pattern that matches valid mana costs.
// TODO: Support X costs and other special cases.
var ManaPattern = regexp.MustCompile(`^(?:\{[0-9WUBRGC]+\})*$`)

var LifeCostPattern = regexp.MustCompile(`^Pay \d+ life$`)

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
		case isLifeCost(trimmed):
			lifeCost, err := ParseLifeCost(trimmed)
			if err != nil {
				return nil, fmt.Errorf("failed to parse life cost %s: %w", trimmed, err)
			}
			costs = append(costs, lifeCost)
		case isDiscardCost(trimmed):
			card, ok := source.(*Card)
			if !ok {
				return nil, fmt.Errorf("source %s not a card", source.Name())
			}
			costs = append(costs, &DiscardCost{Card: card})
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

type LifeCost struct {
	Amount int
	Player *Player
}

func ParseLifeCost(input string) (*LifeCost, error) {
	// Example input: "Pay 3 life"
	re := regexp.MustCompile(`^Pay (\d+) life$`)
	matches := re.FindStringSubmatch(input)
	if len(matches) != 2 {
		return nil, fmt.Errorf("invalid life cost format: %s", input)
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid life amount: %s", matches[1])
	}
	return &LifeCost{Amount: amount}, nil
}

func (l *LifeCost) CanPay(state *GameState, player *Player) bool {
	// Check if the player has enough life to pay the cost
	return player.Life >= l.Amount
}

func (l *LifeCost) Description() string {
	// Return a string representation of the life cost
	return fmt.Sprintf("Pay %d life", l.Amount)
}

func (l *LifeCost) Pay(state *GameState, player *Player) error {
	// Check if the player can pay the cost
	if !l.CanPay(state, player) {
		return fmt.Errorf("not enough life to pay cost: %d", l.Amount)
	}
	// Subtract the life cost from the player's life total
	player.Life -= l.Amount
	return nil
}

// ManaCost represents the mana cost of a card or ability.
type ManaCost struct {
	Colors  map[Color]int
	Generic int
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

func (c *ManaCost) ManaValue() int {
	// The mana value is the total of generic mana and colored mana.
	// Each colored mana counts as 1, and generic mana counts as its value.
	manaValue := c.Generic
	for _, amount := range c.Colors {
		manaValue += amount
	}
	return manaValue
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

type DiscardCost struct {
	Card *Card
}

func (c *DiscardCost) CanPay(game *GameState, player *Player) bool {
	if _, err := player.Hand.Get(c.Card.ID()); err != nil {
		return false
	}
	return true
}

func (c *DiscardCost) Description() string {
	return fmt.Sprintf("discard %s", c.Card.Name())
}

func (c *DiscardCost) Pay(game *GameState, player *Player) error {
	// Check if the player can pay the discard cost
	if !c.CanPay(game, player) {
		return fmt.Errorf("cannot discard card %s", c.Card.Name())
	}
	// Remove the card from the player's hand
	card, err := player.Hand.Take(c.Card.ID())
	if err != nil {
		return fmt.Errorf("failed to discard card %s: %w", c.Card.Name(), err)
	}
	if err := player.Graveyard.Add(card); err != nil {
		return fmt.Errorf("failed to move discarded card %s to graveyard: %w", c.Card.Name(), err)
	}
	return nil
}

// TapCost represents a cost that requires tapping the permanent.
type TapCost struct {
	Permanent *Permanent
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
	tempPool := player.ManaPool.Copy()
	for remaining > 0 {
		choices := []Choice{}
		for _, color := range tempPool.ColorsAvailable() {
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
		tempPool.Use(color, 1)
		choice[selected.Name]++
		remaining--
	}
	return choice, nil
}

// isManaCost checks if the input string is a valid mana cost.
func isMana(input string) bool {
	return ManaPattern.MatchString(input)
}

func isLifeCost(input string) bool {
	return LifeCostPattern.MatchString(input)
}

// isTapCost checks if the input string is a tap cost.
func isTapCost(input string) bool {
	return input == "{T}"
}

func isDiscardCost(input string) bool {
	return input == "Discard this card"
}

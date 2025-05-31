package cost

import (
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// TODO: Make this so only GameState can be passed in.
type State any

type Card interface {
	ID() string
	Name() string
}

type Zone interface {
	Get(string) (Card, error)
	Take(string) (Card, error)
	Add(Card) error
}

// TODO: Make this so only Player can be passed in.
type Player interface {
	ChooseManaForGeneric(int) (map[string]int, error)
	Graveyard() Zone
	Hand() Zone
	Life() int
	LoseLife(int) error
	ManaPool() *mana.Pool
}

// TODO: Make this so only Object can be passed in.
type Object interface {
	Name() string
}

// Cost represents a cost that can be paid in the
type Cost interface {
	CanPay(State, Player) bool
	Description() string
	Pay(State, Player) error
}

type Permanent interface {
	IsTapped() bool
	Tap() error
}

func AddCosts(costs ...Cost) Cost {
	// If there's only one cost, return it directly.
	if len(costs) == 1 {
		return costs[0]
	}
	// Otherwise, return a composite cost.
	return &CompositeCost{costs: costs}
}

// ManaPattern is a regex pattern that matches valid mana costs.
// TODO: Support X costs and other special cases.
var ManaPattern = regexp.MustCompile(`^(?:\{[0-9WUBRGC]+\})*$`)

var LifeCostPattern = regexp.MustCompile(`^Pay \d+ life$`)

// NewCost creates a new cost based on the input string and the source.
// TODO: Maybe rename to NewCost, only return a composit cost when there's
// more than one
func NewCost(input string, source Object) (Cost, error) {
	parts := strings.Split(input, ",")
	var costs []Cost
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		switch {
		case isTapCost(trimmed):
			permanent, ok := source.(Permanent)
			if !ok {
				return nil, fmt.Errorf("source %s not permanent", source.Name())
			}
			costs = append(costs, &TapCost{permanent: permanent})
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
			card, ok := source.(Card)
			if !ok {
				return nil, fmt.Errorf("source %s not a card", source.Name())
			}
			costs = append(costs, &DiscardCost{card: card})
		default:
			return nil, fmt.Errorf("unknown cost %s", trimmed)
		}
	}
	if len(costs) == 1 {
		return costs[0], nil
	}
	return &CompositeCost{costs: costs}, nil
}

// CompositeCost represents a cost that is a combination of multiple costs.
type CompositeCost struct {
	costs []Cost
}

// CanPay checks if all costs in the composite cost can be paid with the
// current game state.
func (c *CompositeCost) CanPay(state State, player Player) bool {
	for _, cost := range c.costs {
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
	for _, cost := range c.costs {
		if _, ok := cost.(*ManaCost); ok {
			costStrings = append(costStrings, cost.Description())
		}
	}
	// Tap
	for _, cost := range c.costs {
		if _, ok := cost.(*TapCost); ok {
			costStrings = append(costStrings, cost.Description())
		}
	}
	//Sacrifice
	for _, cost := range c.costs {
		if _, ok := cost.(*SacrificeCost); ok {
			costStrings = append(costStrings, cost.Description())
		}
	}
	return strings.Join(costStrings, ", ")
}

// Pay pays all costs in the composite cost.
// TODO: If one cost fails, we need to roll back the others.
func (c *CompositeCost) Pay(state State, player Player) error {
	// TODO: Maybe there's a better way to do this, but this helps with the
	// needing to roll back thing.
	for _, cost := range c.costs {
		if !cost.CanPay(state, player) {
			return fmt.Errorf(
				"failed to pay composite cost",
			)
		}
	}
	for _, cost := range c.costs {
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
	amount int
	player *Player
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
	return &LifeCost{amount: amount}, nil
}

func (l *LifeCost) CanPay(state State, player Player) bool {
	// Check if the player has enough life to pay the cost
	return player.Life() >= l.amount
}

func (l *LifeCost) Description() string {
	// Return a string representation of the life cost
	return fmt.Sprintf("Pay %d life", l.amount)
}

func (l *LifeCost) Pay(state State, player Player) error {
	// Check if the player can pay the cost
	if !l.CanPay(state, player) {
		return fmt.Errorf("not enough life to pay cost: %d", l.amount)
	}
	// Subtract the life cost from the player's life total
	if err := player.LoseLife(l.amount); err != nil {
		return fmt.Errorf("failed to pay life cost: %w", err)
	}
	return nil
}

// ManaCost represents the mana cost of a card or ability.
type ManaCost struct {
	colors  map[mtg.Color]int
	generic int
}

// CanPay checks if the cost can be paid with the current game state.
// TODO Maybe this should just be *Player
func (c *ManaCost) CanPay(state State, player Player) bool {
	tempPool := player.ManaPool().Copy()
	for color, amount := range c.colors {
		if !tempPool.Has(color, amount) {
			return false
		}
		tempPool.Use(color, amount)
	}
	if !tempPool.HasGeneric(c.generic) {
		return false
	}
	return true
}

// Description returns a string representation of the mana cost.
func (c *ManaCost) Description() string {
	// TODO: colors might not be the right name, costs? Symbols?
	var cs []string
	if c.generic > 0 {
		cs = append(cs, fmt.Sprintf("{%d}", c.generic))
	}
	for _, color := range []mtg.Color{
		// Do not reorder this list, follows the color wheel standard for mana
		// costs.
		mtg.ColorColorless,
		mtg.ColorWhite,
		mtg.ColorBlue,
		mtg.ColorBlack,
		mtg.ColorRed,
		mtg.ColorGreen,
	} {
		if c.colors[color] > 0 {
			for range c.colors[color] {
				cs = append(cs, fmt.Sprintf("{%s}", color))
			}
		}
	}
	return strings.Join(cs, "")
}

func (c *ManaCost) ManaValue() int {
	// The mana value is the total of generic mana and colored mana.
	// Each colored mana counts as 1, and generic mana counts as its value.
	manaValue := c.generic
	for _, amount := range c.colors {
		manaValue += amount
	}
	return manaValue
}

// Pay pays the mana cost by using the mana from the mana pool.
// TODO: Need to roll back if the cost is partially paid and fails.
func (c *ManaCost) Pay(game State, player Player) error {
	// Pay colored mana
	for color, amount := range c.colors {
		if err := player.ManaPool().Use(color, amount); err != nil {
			return err
		}
	}
	// Pay generic mana
	if c.generic > 0 {
		choice, err := player.ChooseManaForGeneric(c.generic)
		if err != nil {
			return err
		}
		for colorChoice, amount := range choice {
			color, err := mtg.StringToColor(colorChoice)
			if err != nil {
				return err
			}
			player.ManaPool().Use(color, amount)
		}
	}
	return nil
}

// SacrificeCost represents a cost that requires sacrificing the permanent.
// TODO: Support sacrificing other permanents.
type SacrificeCost struct {
	permanent *Permanent
}

// CanPay checks if the sacrifice cost can be paid with the current game
// state.
func (c *SacrificeCost) CanPay(game State, player Player) bool {
	// TODO: Pretty much always true unless there is some state that says the
	// permanent can't be sacrificed.
	return true
}

// Description returns a string representation of the sacrifice cost.
func (c *SacrificeCost) Description() string {
	return "sacrifice this permanent"
}

// Pay pays the sacrifice cost by sacrificing the permanent.
func (c *SacrificeCost) Pay(game State, player Player) error {
	// TODO: Implement this
	return errors.New("not implemented")
}

type DiscardCost struct {
	card Card
}

func (c *DiscardCost) CanPay(state State, player Player) bool {
	if _, err := player.Hand().Get(c.card.ID()); err != nil {
		return false
	}
	return true
}

func (c *DiscardCost) Description() string {
	return fmt.Sprintf("discard %s", c.card.Name())
}

func (c *DiscardCost) Pay(game State, player Player) error {
	// Check if the player can pay the discard cost
	if !c.CanPay(game, player) {
		return fmt.Errorf("cannot discard card %s", c.card.Name())
	}
	// Remove the card from the player's hand
	card, err := player.Hand().Take(c.card.ID())
	if err != nil {
		return fmt.Errorf("failed to discard card %s: %w", c.card.Name(), err)
	}
	if err := player.Graveyard().Add(card); err != nil {
		return fmt.Errorf("failed to move discarded card %s to graveyard: %w", c.card.Name(), err)
	}
	return nil
}

// TapCost represents a cost that requires tapping the permanent.
type TapCost struct {
	permanent Permanent
}

// CanPay checks if the tap cost can be paid with the current game state.
func (c *TapCost) CanPay(game State, player Player) bool {
	return !c.permanent.IsTapped()
}

// Description returns a string representation of the tap cost.
func (c *TapCost) Description() string {
	return "{T}"
}

// Pay pays the tap cost by tapping the permanent.
func (c *TapCost) Pay(game State, player Player) error {
	if err := c.permanent.Tap(); err != nil {
		return fmt.Errorf("failed to pay {T} cost: %w", err)
	}
	return nil
}

// ParseManaCost parses a mana cost string and returns a ManaCost.
func ParseManaCost(costStr string) (*ManaCost, error) {
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(costStr, -1)
	manaCost := ManaCost{
		generic: 0,
		colors:  map[mtg.Color]int{},
	}
	for _, match := range matches {
		symbol := strings.ToUpper(match[1])
		if val, err := strconv.Atoi(symbol); err == nil {
			manaCost.generic += val
		} else {
			color, err := mtg.StringToColor(symbol)
			if err != nil {
				return nil, err
			}
			manaCost.colors[color]++
		}
	}
	return &manaCost, nil
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

package game

import (
	"fmt"
	"strings"
)

type Deck []*Card

type GameState struct {
	LastActionFailed   bool
	CurrentPlayer      int
	CurrentPhase       string
	Cheat              bool
	Life               int
	GameLost           bool
	Deck               Deck
	Turn               int
	SpellsCastThisTurn []string
	Hand               []*Card
	Battlefield        []*Permanent
	Graveyard          []*Card
	Exile              []*Card
	ManaPool           ManaPool
	PotentialMana      ManaPool
	StormCount         int
	TurnCount          int
	MessageLog         []string
	TurnMessageLog     []string
	LandDrop           bool
	EventListeners     []EventHandler
	Message            string
}

type CheatResult struct {
	Message string
}

type ActionResult struct {
	EndTurn bool
	Message string
}

// Maybe do something where I can pass in "play Island" and it'll take the second param as the Choice and only prompt if it is missing
// maybe support typing in the number or the name of the card
type ChoiceResolver interface {
	ChooseOne(prompt string, options []Choice) Choice
	//ChooseN(prompt string, options []Choice, n int) []Choice
	//ChooseUpToN(prompt string, options []Choice, n int) []Choice
	//ChooseAny(prompt string, options []Choice) []Choice
	//Confirm(prompt string) bool // For simple yes/no prompts
}

// PlayerAgent defines how player decisions are made
type PlayerAgent interface {
	ReportState(state *GameState)
	GetNextAction(state *GameState) GameAction
	ChoiceResolver
}

type OptionPrompt struct {
	Message string
	Options []string
	Min     int
	Max     int
}

type Choice struct {
	Name   string
	Index  int
	Source string
}

// TODO: make this more better
func (g *GameState) PrintDeck() {
	for _, card := range g.Deck {
		fmt.Println(card.Name)
	}
}

var startingLife = 20

func NewGameState() *GameState {
	return &GameState{}
}

func (g *GameState) InitializeNewGame(filename string) error {
	g.Life = startingLife
	deck, err := importDeck(filename)
	if err != nil {
		return err
	}
	g.Deck = deck
	g.ShuffleDeck()
	// TODO: Start with 6 so we draw into 7 to simulate on the play
	// this should be configurable
	g.DrawCards(6)
	return nil
}

func (g *GameState) Log(message string) {
	// TODO: There's probably a more elegant way to do this
	g.TurnMessageLog = append(g.TurnMessageLog, message)
	g.MessageLog = append(g.MessageLog, message)
}

// TODO: There's probably a more elegant way to do this
func (g *GameState) Error(err error) {
	g.TurnMessageLog = append(g.TurnMessageLog, "ERROR: "+err.Error())
	g.MessageLog = append(g.MessageLog, "ERROR: "+err.Error())
}

func (g *GameState) ShuffleDeck() {
	shuffleCards(g.Deck)
}

// TODO Maybe return error?
func (g *GameState) DrawCards(n int) {
	taken, remaining, err := takeNCards(g.Deck, n)
	if err != nil {
		g.GameLost = true
		g.Log("Game Lost: Decked")
		return
	}
	var names []string
	for _, card := range taken {
		names = append(names, card.Name)
	}
	g.Deck = remaining
	g.Hand = append(g.Hand, taken...)
	g.Log(fmt.Sprintf("Drew %d cards: %s", n, strings.Join(names, ", ")))
}

/*
func (g *GameState) CastSpell(card *Card) error {
	// Step 1: Pay costs (mana, sacrifice, etc.)
	// Step 2: Run cast triggers (e.g., storm)
	for _, trigger := range card.CastTriggers {
		if err := trigger(card, g); err != nil {
			return err
		}
	}

	// Step 3: Put spell on stack
	//g.Stack = append(g.Stack, card)
	g.Log = append(g.Log, fmt.Sprintf("Cast %s", card.Name))
	return nil
}

func (g *GameState) ResolveSpell(card *Card) error {
	for _, effect := range card.OnResolutionEffects {
		if err := effect(card, g); err != nil {
			return err
		}
	}

	if card.Types["Creature"] || card.Types["Artifact"] || card.Types["Enchantment"] {
		g.Battlefield = append(g.Battlefield, card)
		for _, etb := range card.ETBTriggers {
			if err := etb(card, g); err != nil {
				return err
			}
		}
	} else {
		// Instants, sorceries go straight to graveyard
		g.Graveyard = append(g.Graveyard, card)
	}

	//g.Stack = g.Stack[:len(g.Stack)-1]
	return nil
}

func (g *GameState) Activate(card *Card, abilityIndex int) error {
	if abilityIndex < 0 || abilityIndex >= len(card.ActivatedAbilities) {
		return fmt.Errorf("invalid ability index for %s", card.Name)
	}
	return card.ActivatedAbilities[abilityIndex](card, g)
}

func (c *Card) Resolve(state *GameState) error {
	for _, ability := range c.OnResolutionEffects {
		if err := ability(c, state); err != nil {
			return err
		}
	}
	return nil
}

func (gs *GameState) EmitEvent(evt GameEvent) {
	// Add it to a queue if needed or process it immediately
	for _, listener := range gs.EventListeners {
		listener(evt, gs)
	}
}
*/

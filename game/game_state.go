package game

type ManaPool map[string]int

type GameState struct {
	Battlefield        *Battlefield
	Cheat              bool
	CurrentPhase       string
	CurrentPlayer      int
	CurrentStep        string
	Deck               *Deck
	EventListeners     []EventHandler
	Exile              []*Card
	Graveyard          []*Card
	Hand               *Hand
	LandDrop           bool
	LastActionFailed   bool
	Life               int
	ManaPool           ManaPool
	MaxHandSize        int
	MaxTurns           int
	Message            string
	MessageLog         []string
	PotentialMana      ManaPool
	SpellsCastThisTurn []string
	StormCount         int
	Turn               int
	TurnCount          int
	TurnMessageLog     []string
}

// Maybe do something where I can pass in "play Island" and it'll take the second param as the Choice and only prompt if it is missing
// maybe support typing in the number or the name of the card
type ChoiceResolver interface {
	ChooseOne(prompt string, choices []Choice) Choice
	//ChooseN(prompt string, choices []Choice, n int) []Choice
	//ChooseUpToN(prompt string, choices []Choice, n int) []Choice
	//ChooseAny(prompt string, choices []Choice) []Choice
	//Confirm(prompt string) bool // For simple yes/no prompts
}

// PlayerAgent defines how player decisions are made
type PlayerAgent interface {
	ReportState(state *GameState)
	GetNextAction(state *GameState) GameAction
	ChoiceResolver
}

type Choice struct {
	Name   string
	Index  int
	Source string
}

type GameStateConfig struct {
	StartingLife int
	MaxTurns     int
	DeckList     string
}

func NewGameState() *GameState {
	return &GameState{}
}

func (g *GameState) InitializeNewGame(config GameStateConfig) error {
	g.MaxTurns = config.MaxTurns
	g.Life = config.StartingLife
	g.MaxHandSize = 7
	deck, err := importDeck(config.DeckList)
	if err != nil {
		return err
	}
	g.Deck = deck
	g.Deck.Shuffle()
	g.Hand = &Hand{}
	for range 7 { // TODO: This sucks, figure this out along will mulls
		_, err := g.ResolveAction(GameAction{
			Type: ActionDraw,
		}, nil)
		if err != nil {
			return err
		}
	}
	g.ManaPool = make(ManaPool)
	g.Battlefield = &Battlefield{}
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

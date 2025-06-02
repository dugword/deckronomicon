package player

import (
	"deckronomicon/packages/game/card"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/zone"
	"deckronomicon/packages/query"
	"errors"
	"fmt"
)

// TODO Some places I pass PlayerAgent where I really should just pass Player
// Think about how to refactor this so that PlayerAgent is not needed
// everywhere, unless decisions or input are needed from the player.

// TODO Maybe PlayerAgent should be a field of Player.

// TODO Not sure how I want to handle the "Modes" of a player, standardize them here or let them be defiend ad-hoc,
// in the strategy rules.
const (
	ModeSetup = "Setup"
)

type Player struct {
	Agent Agent
	// TODO: This should be a field of GameState, there is only one
	// battlefield per game.
	exile       *zone.Exile
	graveyard   *zone.Graveyard
	hand        *zone.Hand
	id          string
	LandDrop    bool
	library     *zone.Library
	life        int
	manaPool    *manaPool
	MaxHandSize int
	Mode        string
	Mulligans   int
	// TODO: This is broken and maybe a bad idea
	// PotentialMana *ManaPool
	revealed *zone.Revealed
	Stops    []mtg.Step
	// TODO This should probably be used in game engine
	// StartingHand []string
	Turn int
}

// NewPlayer creates a new Player instance.
// TODO: Make this a constructor that takes a config file or parameters.
func New(agent Agent, id string, life int, mode string) *Player {
	player := Player{
		Agent:     agent,
		exile:     zone.NewExile(),
		graveyard: zone.NewGraveyard(),
		hand:      zone.NewHand(),
		id:        id,
		library:   zone.NewLibrary(),
		life:      life,
		manaPool:  newManaPool(),
		// I don't like this, but it works for now
		// TODO: pass this from somewhere
		MaxHandSize: 7,
		Mode:        mode,
		// PotentialMana: NewManaPool(),
		revealed: zone.NewRevealed(),
		// TODO: Make this configurable
		Stops: []mtg.Step{mtg.StepDraw, mtg.StepPrecombatMain},
	}
	agent.RegisterPlayer(&player)
	return &player
}

func (p *Player) AddMana(mana string) error {
	return p.manaPool.AddMana(mana)
}

func (p *Player) AssignLibrary(library *zone.Library) {
	p.library = library
}

func (p *Player) BottomCard(cardID string) error {
	card, err := p.hand.Take(cardID)
	if err != nil {
		return fmt.Errorf("failed to take card %s from hand: %w", cardID, err)
	}
	p.library.Add(card)
	return nil
}

func (p *Player) CheatAddCard(card *card.Card) {
	p.hand.Add(card)
}

func (p *Player) DiscardCard(cardID string) error {
	card, err := p.hand.Take(cardID)
	if err != nil {
		return fmt.Errorf("failed to take card %s from hand: %w", cardID, err)
	}
	p.graveyard.Add(card)
	return nil
}

// Draw draws a card from the library into the player's hand.
func (p *Player) DrawCard() (string, error) {
	card, err := p.library.TakeTop()
	if err != nil {
		if errors.Is(err, mtg.ErrLibraryEmpty) {
			return "", mtg.PlayerLostError{
				Reason: mtg.DeckedOut,
			}
		}
		return "", err
	}
	p.hand.Add(card)
	return card.Name(), nil
}

func (p *Player) EmptyManaPool() {
	p.manaPool.Empty()
}

func (p *Player) Exile() query.View {
	return query.NewView(p.exile.Name(), p.exile.GetAll())
}

func (p *Player) GainLife(amount int) {
	if amount < 0 {
		amount = 0
	}
	p.life += amount
}

func (p *Player) Graveyard() query.View {
	return query.NewView(p.graveyard.Name(), p.graveyard.GetAll())
}

func (p *Player) HasStop(step mtg.Step) bool {
	for _, stop := range p.Stops {
		if stop == step {
			return true
		}
	}
	return false
}

func (p *Player) Hand() query.View {
	return query.NewView(p.hand.Name(), p.hand.GetAll())
}

func (p *Player) ID() string {
	return p.id
}

func (p *Player) Library() query.View {
	return query.NewView(p.library.Name(), p.library.GetAll())
}

func (p *Player) Life() int {
	return p.life
}

func (p *Player) LoseLife(amount int) error {
	// TODO: Confirm this is how the game works per the rules.
	if amount < 0 {
		amount = 0
	}
	p.life -= amount
	if p.life <= 0 {
		return (mtg.PlayerLostError{
			Reason: mtg.LifeTotalZero,
		})
	}
	return nil
}

// ManaPool returns a copy of the player's mana pool.
func (p *Player) ManaPool() *manaPool {
	return p.manaPool.Copy()
}

func (p *Player) PlaceCard(card *card.Card, zone mtg.Zone) error {
	switch zone {
	case mtg.ZoneHand:
		p.hand.Add(card)
	case mtg.ZoneGraveyard:
		p.graveyard.Add(card)
	case mtg.ZoneExile:
		p.exile.Add(card)
	case mtg.ZoneLibrary:
		p.library.Add(card)
	default:
		return fmt.Errorf("unknown zone: %s", zone)
	}
	return nil
}

func (p *Player) Revealed() query.View {
	return query.NewView(p.revealed.Name(), p.revealed.GetAll())
}

func (p *Player) ShuffleLibrary() {
	p.library.Shuffle()
}

func (p *Player) TakeCard(cardID string, zone mtg.Zone) (*card.Card, error) {
	switch zone {
	case mtg.ZoneHand:
		return p.hand.Take(cardID)
	case mtg.ZoneGraveyard:
		return p.graveyard.Take(cardID)
	case mtg.ZoneExile:
		return p.exile.Take(cardID)
	case mtg.ZoneLibrary:
		return p.library.Take(cardID)
	default:
		return nil, fmt.Errorf("unknown zone: %s", zone)
	}
}

func (p *Player) TakeTopCard() (*card.Card, error) {
	card, err := p.library.TakeTop()
	if err != nil {
		return nil, nil
	}
	return card, nil
}

func (p *Player) Tutor(query query.Predicate) error {
	card, err := p.library.TakeBy(query)
	if err != nil {
		return fmt.Errorf("failed to take card from library: %w", err)
	}
	p.hand.Add(card)
	return nil
}

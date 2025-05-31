package player

import (
	"deckronomicon/packages/game/mana"
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
	Agent       Agent
	Battlefield *zone.Battlefield
	Exile       *zone.Exile
	Graveyard   *zone.Graveyard
	Hand        *zone.Hand
	id          string
	LandDrop    bool
	Library     *zone.Library
	Life        int
	ManaPool    *mana.Pool
	MaxHandSize int
	Mode        string
	Mulligans   int
	// TODO: This is broken and maybe a bad idea
	// PotentialMana *ManaPool
	Revealed *zone.Revealed
	Stops    []mtg.Step
	// TODO This should probably be used in game engine
	// StartingHand []string
	Turn int
}

// NewPlayer creates a new Player instance.
// TODO: Make this a constructor that takes a config file or parameters.
func New(agent Agent, id string, life int, mode string) *Player {
	player := Player{
		Agent:       agent,
		Battlefield: zone.NewBattlefield(),
		Exile:       zone.NewExile(),
		Graveyard:   zone.NewGraveyard(),
		Hand:        zone.NewHand(),
		id:          id,
		Library:     zone.NewLibrary(),
		Life:        life,
		ManaPool:    mana.NewPool(),
		// I don't like this, but it works for now
		// TODO: pass this from somewhere
		MaxHandSize: 7,
		Mode:        mode,
		// PotentialMana: NewManaPool(),
		Revealed: zone.NewRevealed(),
		// TODO: Make this configurable
		Stops: []mtg.Step{mtg.StepDraw, mtg.StepPrecombatMain},
	}
	agent.RegisterPlayer(&player)
	return &player
}

func (p *Player) HasStop(step mtg.Step) bool {
	for _, stop := range p.Stops {
		if stop == step {
			return true
		}
	}
	return false
}

func (p *Player) ID() string {
	return p.id
}

func (p *Player) Tutor(query query.Predicate) error {
	card, err := p.Library.TakeBy(query)
	if err != nil {
		return fmt.Errorf("failed to take card from library: %w", err)
	}
	p.Hand.Add(card)
	return nil
}

// Draw draws a card from the library into the player's hand.
func (p *Player) Draw() (string, error) {
	card, err := p.Library.TakeTop()
	if err != nil {
		if errors.Is(err, mtg.ErrLibraryEmpty) {
			return "", mtg.PlayerLostError{
				Reason: mtg.DeckedOut,
			}
		}
		return "", err
	}
	p.Hand.Add(card)
	return card.Name(), nil
}

func (p *Player) Discard(player Player) error {
	/*
		choices := CreateChoices(player.Hand.GetAll(), ZoneHand)
		choice, err := player.Agent.ChooseOne(
			"Which card to discard from hand",
			source,
			choices,
		)
		if err != nil {
			return fmt.Errorf("failed to choose card to discard: %w", err)
		}
		card, err := player.Hand.Get(choice.ID)
		if err != nil {
			return fmt.Errorf("failed to get card from hand: %w", err)
		}
		player.Hand.Remove(card.ID())
		player.Graveyard.Add(card)
		return nil
	*/
	return nil
}

/*
func (p *Player) Zones() []Zone {
	return []Zone{
		p.Battlefield,
		p.Library,
		p.Hand,
		p.Graveyard,
	}
}
*/

/*
func (p *Player) GetAvailableToPlay(state *GameState) map[string][]Object {
	available := map[string][]Object{}
	for _, zone := range p.Zones() {
		available[zone.ZoneType()] = append(available[zone.ZoneType()], zone.AvailableToPlay(state, p)...)
	}
	return available
}

func (p *Player) GetAvailableToActivate(state *GameState) map[string][]Object {
	available := map[string][]Object{}
	for _, zone := range p.Zones() {
		// TODO: Make this reutrn objects
		for _, ability := range zone.AvailableActivatedAbilities(state, p) {
			available[zone.ZoneType()] = append(available[zone.ZoneType()], ability)
		}
	}
	return available

}
*/

/*
func (p *Player) GetZone(zone string) (Zone, error) {
	for _, z := range p.Zones() {
		if z.ZoneType() == zone {
			return z, nil
		}
	}
	return nil, fmt.Errorf("zone %s not found", zone)
}
*/

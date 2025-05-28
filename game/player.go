package game

import "fmt"

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
	Agent         PlayerAgent
	Battlefield   *Battlefield
	Exile         *Exile
	Graveyard     *Graveyard
	Hand          *Hand
	ID            string
	LandDrop      bool
	Library       *Library
	Life          int
	ManaPool      *ManaPool
	MaxHandSize   int
	Mode          string
	Mulligans     int
	PotentialMana *ManaPool
	Revealed      *Revealed
	Stops         []string
	StartingHand  []string
	Turn          int
}

// NewPlayer creates a new Player instance.
// TODO: Make this a constructor that takes a config file or parameters.
func NewPlayer(agent PlayerAgent) *Player {
	player := Player{
		Agent:       agent,
		Battlefield: NewBattlefield(),
		Exile:       NewExile(),
		Graveyard:   NewGraveyard(),
		Hand:        NewHand(),
		ID:          agent.PlayerID(),
		Library:     NewLibrary(),
		// I don't like this, but it works for now
		Life:     20,
		ManaPool: NewManaPool(),
		// I don't like this, but it works for now
		MaxHandSize:   7,
		Mode:          ModeSetup,
		PotentialMana: NewManaPool(),
		Revealed:      NewRevealed(),
		// TODO: Make this configurable
		Stops: []string{StepDraw, StepPrecombatMain},
	}
	return &player
}

func (p *Player) Zones() []Zone {
	return []Zone{
		p.Battlefield,
		p.Library,
		p.Hand,
		p.Graveyard,
	}
}

func (p *Player) GetAvailableToPlay(state *GameState) map[string][]GameObject {
	available := map[string][]GameObject{}
	for _, zone := range p.Zones() {
		available[zone.ZoneType()] = append(available[zone.ZoneType()], zone.AvailableToPlay(state, p)...)
	}
	return available
}

func (p *Player) GetAvailableToActivate(state *GameState) map[string][]GameObject {
	available := map[string][]GameObject{}
	for _, zone := range p.Zones() {
		// TODO: Make this reutrn objects
		for _, ability := range zone.AvailableActivatedAbilities(state, p) {
			available[zone.ZoneType()] = append(available[zone.ZoneType()], ability)
		}
	}
	return available

}

func (p *Player) GetZone(zone string) (Zone, error) {
	for _, z := range p.Zones() {
		if z.ZoneType() == zone {
			return z, nil
		}
	}
	return nil, fmt.Errorf("zone %s not found", zone)
}

func (p *Player) ShouldAutoPass(currentStep string) bool {
	for _, stop := range p.Stops {
		if stop == currentStep {
			return false
		}
	}
	return true
}

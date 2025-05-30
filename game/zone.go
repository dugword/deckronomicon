package game

// library, hand, battlefield, graveyard, stack, exile, and command
const (
	ZoneBattlefield = "Battlefield"
	ZoneCommand     = "Command"
	ZoneExile       = "Exile"
	ZoneGraveyard   = "Graveyard"
	ZoneHand        = "Hand"
	ZoneLibrary     = "Library"
	ZoneStack       = "Stack"
	ZoneRevealed    = "Revealed"
)

type Zone interface {
	// TODO make this accept multiple GameObjects
	Add(object GameObject) error
	// This probably makes more sense as a method of Player
	AvailableActivatedAbilities(*GameState, *Player) []GameObject
	// This probably makes more sense as a method of Player
	AvailableToPlay(*GameState, *Player) []GameObject
	Get(id string) (GameObject, error)
	GetAll() []GameObject
	Remove(id string) error
	Take(id string) (GameObject, error)
	Size() int
	ZoneType() string
}

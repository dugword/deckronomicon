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
)

type Zone interface {
	Add(object GameObject) error
	AvailableActivatedAbilities(*GameState) []*ActivatedAbility
	AvailableToPlay(*GameState) []GameObject
	Find(id string) (GameObject, error)
	FindByName(name string) (GameObject, error)
	FindAllBySubtype(subtype Subtype) []GameObject
	Get(id string) (GameObject, error)
	Remove(id string) error
	Take(id string) (GameObject, error)
	Size() int
	ZoneType() string
}

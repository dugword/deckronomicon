package mtg

type Zone string

const (
	ZoneBattlefield Zone = "Battlefield"
	ZoneCommand     Zone = "Command"
	ZoneExile       Zone = "Exile"
	ZoneGraveyard   Zone = "Graveyard"
	ZoneHand        Zone = "Hand"
	ZoneLibrary     Zone = "Library"
	ZoneRevealed    Zone = "Revealed"
	ZoneSideboard   Zone = "Sideboard"
	ZoneStack       Zone = "Stack"
)

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

func (z Zone) Name() string {
	return string(z)
}

func StringToZone(s string) (Zone, bool) {
	stringToZone := map[string]Zone{
		"Battlefield": ZoneBattlefield,
		"Command":     ZoneCommand,
		"Exile":       ZoneExile,
		"Graveyard":   ZoneGraveyard,
		"Hand":        ZoneHand,
		"Library":     ZoneLibrary,
		"Revealed":    ZoneRevealed,
		"Sideboard":   ZoneSideboard,
		"Stack":       ZoneStack,
	}
	zone, ok := stringToZone[s]
	return zone, ok
}

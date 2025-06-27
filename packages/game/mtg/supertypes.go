package mtg

type Supertype string

const (
	SupertypeBasic     Supertype = "Basic"
	SupertypeLegendary Supertype = "Legendary"
	SupertypeSnow      Supertype = "Snow"
	SupertypeWorld     Supertype = "World"
)

// StringToSupertype converts a string to a Supertype.
func StringToSupertype(s string) (Supertype, bool) {
	stringToSupertype := map[string]Supertype{
		"Basic":     SupertypeBasic,
		"Legendary": SupertypeLegendary,
		"Snow":      SupertypeSnow,
		"World":     SupertypeWorld,
	}
	supertype, ok := stringToSupertype[s]
	if !ok {
		return "", false
	}
	return supertype, true
}

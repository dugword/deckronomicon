package mtg

// TODO This package name is gunna be confusing. rename it so we don't have
// color := colors.color maybe put it in a core package or something
import "fmt"

// Color is a color in Magic: The Gathering.
type Color string

const (
	ColorBlack     Color = "B"
	ColorBlue      Color = "U"
	ColorColorless Color = "C"
	ColorGreen     Color = "G"
	ColorRed       Color = "R"
	ColorWhite     Color = "W"
)

// StringToColor converts a string to a Color.
func StringToColor(s string) (Color, error) {
	colors := map[string]Color{
		"B": ColorBlack,
		"C": ColorColorless,
		"G": ColorGreen,
		"R": ColorRed,
		"U": ColorBlue,
		"W": ColorWhite,
	}
	color, ok := colors[s]
	if !ok {
		return "", fmt.Errorf("unknown color: %s", s)
	}
	return color, nil
}

// Colors represents the colors of a card or object in the game.
type Colors struct {
	Black bool
	Blue  bool
	Green bool
	Red   bool
	White bool
}

// StringToColor converts a string to a Color.
func StringsToColors(ss []string) (Colors, error) {
	colors := Colors{}
	for _, s := range ss {
		color, err := StringToColor(s)
		if err != nil {
			return Colors{}, fmt.Errorf("failed to parse color %q: %w", s, err)
		}
		switch color {
		case ColorBlack:
			colors.Black = true
		case ColorBlue:
			colors.Blue = true
		case ColorGreen:
			colors.Green = true
		case ColorRed:
			colors.Red = true
		case ColorWhite:
			colors.White = true
		}
	}
	return colors, nil
}

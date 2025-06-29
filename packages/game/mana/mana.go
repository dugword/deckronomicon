package mana

import (
	"fmt"
	"regexp"
	"strconv"
)

type Color string

// This might be beter as a const in mtg package
const (
	White     Color = "W"
	Blue      Color = "U"
	Black     Color = "B"
	Red       Color = "R"
	Green     Color = "G"
	Colorless Color = "C"
)

func Colors() []Color {
	return []Color{
		Colorless,
		White,
		Blue,
		Black,
		Red,
		Green,
	}
}

func (c Color) String() string {
	return fmt.Sprintf("{%s}", string(c))
}

func StringToColor(s string) (Color, bool) {
	switch s {
	case "W", "{W}":
		return White, true
	case "U", "{U}":
		return Blue, true
	case "B", "{B}":
		return Black, true
	case "R", "{R}":
		return Red, true
	case "G", "{G}":
		return Green, true
	case "C", "{C}":
		return Colorless, true
	default:
		return "", false
	}
}

var manaStringRegexp = regexp.MustCompile(`\{(\d+|[WUBRGC])\}`)

func ParseManaString(manaString string) (Amount, error) {
	amount := Amount{}
	if manaString == "" {
		return amount, nil
	}
	matches := manaStringRegexp.FindAllStringSubmatch(manaString, -1)
	if matches == nil {
		return amount, fmt.Errorf("invalid mana string: %s", manaString)
	}
	for _, match := range matches {
		symbol := match[1]
		switch symbol {
		case "W":
			amount.white++
		case "U":
			amount.blue++
		case "B":
			amount.black++
		case "R":
			amount.red++
		case "G":
			amount.green++
		case "C":
			amount.colorless++
		default:
			n, err := strconv.Atoi(symbol)
			if err != nil {
				return amount, fmt.Errorf("invalid mana symbol %q", symbol)
			}
			amount.generic += n
		}
	}
	return amount, nil
}

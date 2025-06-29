package mana

import (
	"fmt"
	"strings"
)

type Amount struct {
	generic   int
	white     int
	blue      int
	black     int
	red       int
	green     int
	colorless int
}

func (a Amount) Has(color Color) bool {
	switch color {
	case White:
		return a.white > 0
	case Blue:
		return a.blue > 0
	case Black:
		return a.black > 0
	case Red:
		return a.red > 0
	case Green:
		return a.green > 0
	case Colorless:
		return a.colorless > 0
	default:
		panic(fmt.Sprintf("unknown color %s", color))
	}
}

func (a Amount) AmountOf(color Color) int {
	switch color {
	case White:
		return a.white
	case Blue:
		return a.blue
	case Black:
		return a.black
	case Red:
		return a.red
	case Green:
		return a.green
	case Colorless:
		return a.colorless
	default:
		panic(fmt.Sprintf("unknown color %s", color))
	}
}

func (a Amount) White() int {
	return a.white
}

func (a Amount) Blue() int {
	return a.blue
}

func (a Amount) Black() int {
	return a.black
}

func (a Amount) Red() int {
	return a.red
}

func (a Amount) Green() int {
	return a.green
}

func (a Amount) Colorless() int {
	return a.colorless
}

func (a Amount) Generic() int {
	return a.generic
}

func (a Amount) Total() int {
	total := a.generic
	total += a.white
	total += a.blue
	total += a.black
	total += a.red
	total += a.green
	total += a.colorless
	return total
}

func (a Amount) WithAddAmount(amount Amount) Amount {
	a.white += amount.white
	a.blue += amount.blue
	a.black += amount.black
	a.red += amount.red
	a.green += amount.green
	a.colorless += amount.colorless
	a.generic += amount.generic
	return a
}

func (a Amount) WithAddMana(amount int, color Color) Amount {
	switch color {
	case White:
		a.white += amount
	case Blue:
		a.blue += amount
	case Black:
		a.black += amount
	case Red:
		a.red += amount
	case Green:
		a.green += amount
	case Colorless:
		a.colorless += amount
	default:
		panic(fmt.Sprintf("unknown color %s", color))
	}
	return a
}

func (a Amount) ManaString() string {
	descriptions := []string{}
	if a.generic > 0 {
		descriptions = append(descriptions, fmt.Sprintf("{%d}", a.generic))
	}
	for range a.colorless {
		descriptions = append(descriptions, Colorless.String())
	}
	for range a.white {
		descriptions = append(descriptions, White.String())
	}
	for range a.blue {
		descriptions = append(descriptions, Blue.String())
	}
	for range a.black {
		descriptions = append(descriptions, Black.String())
	}
	for range a.red {
		descriptions = append(descriptions, Red.String())
	}
	for range a.green {
		descriptions = append(descriptions, Green.String())
	}
	return strings.Join(descriptions, "")
}

package mana

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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
	return []Color{White,
		Blue,
		Black,
		Red,
		Green,
		Colorless,
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

type Pool struct {
	white     int
	blue      int
	black     int
	red       int
	green     int
	colorless int
}

func (p Pool) White() int {
	return p.white
}

func (p Pool) Blue() int {
	return p.blue
}

func (p Pool) Black() int {
	return p.black
}

func (p Pool) Red() int {
	return p.red
}

func (p Pool) Green() int {
	return p.green
}

func (p Pool) Colorless() int {
	return p.colorless
}

func (p Pool) ManaString() string {
	descriptions := []string{}
	for range p.colorless {
		descriptions = append(descriptions, Colorless.String())
	}
	for range p.white {
		descriptions = append(descriptions, White.String())
	}
	for range p.blue {
		descriptions = append(descriptions, Blue.String())
	}
	for range p.black {
		descriptions = append(descriptions, Black.String())
	}
	for range p.red {
		descriptions = append(descriptions, Red.String())
	}
	for range p.green {
		descriptions = append(descriptions, Green.String())
	}
	return strings.Join(descriptions, "")
}

func (p Pool) Copy() Pool {
	return p
}

func (p Pool) Total() int {
	total := 0
	total += p.white
	total += p.blue
	total += p.black
	total += p.red
	total += p.green
	total += p.colorless
	return total
}

func (p Pool) WithAddedAmount(amount Amount) Pool {
	p.white += amount.white
	p.blue += amount.blue
	p.black += amount.black
	p.red += amount.red
	p.green += amount.green
	p.colorless += amount.colorless
	p.colorless += amount.generic
	return p
}

func (p Pool) WithAddedMana(color Color, amount int) Pool {
	switch color {
	case White:
		p.white += amount
	case Blue:
		p.blue += amount
	case Black:
		p.black += amount
	case Red:
		p.red += amount
	case Green:
		p.green += amount
	case Colorless:
		p.colorless += amount
	default:
		panic(fmt.Sprintf("unknown color %s", color))
	}
	return p
}

// WithPayDownFromAmount returns a new Pool with the specified amount of mana paid down.
// It will spend as much mana as possible from the pool, and return the remaining amount
// that could not be paid down.
func (p Pool) WithPayDownFromAmount(amount Amount, colorsForGeneric []Color) (Pool, Amount) {
	if amount.white > p.white {
		amount.white -= p.white
		p.white = 0
	} else {
		p.white -= amount.white
		amount.white = 0
	}
	if amount.blue > p.blue {
		amount.blue -= p.blue
		p.blue = 0
	} else {
		p.blue -= amount.blue
		amount.blue = 0
	}
	if amount.black > p.black {
		amount.black -= p.black
		p.black = 0
	} else {
		p.black -= amount.black
		amount.black = 0
	}
	if amount.red > p.red {
		amount.red -= p.red
		p.red = 0
	} else {
		p.red -= amount.red
		amount.red = 0
	}
	if amount.green > p.green {
		amount.green -= p.green
		p.green = 0
	} else {
		p.green -= amount.green
		amount.green = 0
	}
	if amount.colorless > p.colorless {
		amount.colorless -= p.colorless
		p.colorless = 0
	} else {
		p.colorless -= amount.colorless
		amount.colorless = 0
	}
	if amount.generic > p.Total() {
		genericDeficit := amount.generic - (p.Total())
		amount.generic = genericDeficit
		p.white = 0
		p.blue = 0
		p.black = 0
		p.red = 0
		p.green = 0
		p.colorless = 0
	}
	genericLeft := amount.generic
	if p.white > 0 && genericLeft > 0 {
		if genericLeft >= p.white {
			genericLeft -= p.white
			p.white = 0
		} else {
			p.white -= genericLeft
			genericLeft = 0
		}
	}
	for _, color := range colorsForGeneric {
		p, genericLeft = SpendGenericManaByColor(p, genericLeft, color)
	}
	if genericLeft > 0 {
		amount.generic = genericLeft
	} else {
		amount.generic = 0
	}
	return p, amount
}

func (p Pool) WithSpendFromManaAmount(amount Amount, colorsForGeneric []Color) (Pool, error) {
	manaDeficit := []string{}
	if amount.white > p.white {
		whiteDeficit := amount.white - p.white
		manaDeficit = append(manaDeficit, strings.Repeat(string(White), whiteDeficit))
	}
	p.white -= amount.white
	if amount.blue > p.blue {
		blueDeficit := amount.blue - p.blue
		manaDeficit = append(manaDeficit, strings.Repeat(string(Blue), blueDeficit))
	}
	p.blue -= amount.blue
	if amount.black > p.black {
		blackDeficit := amount.black - p.black
		manaDeficit = append(manaDeficit, strings.Repeat(string(Black), blackDeficit))
	}
	p.black -= amount.black
	if amount.red > p.red {
		redDeficit := amount.red - p.red
		manaDeficit = append(manaDeficit, strings.Repeat(string(Red), redDeficit))
	}
	p.red -= amount.red
	if amount.green > p.green {
		greenDeficit := amount.green - p.green
		manaDeficit = append(manaDeficit, strings.Repeat(string(Green), greenDeficit))
	}
	p.green -= amount.green
	if amount.colorless > p.colorless {
		colorlessDeficit := amount.colorless - p.colorless
		manaDeficit = append(manaDeficit, strings.Repeat(string(Colorless), colorlessDeficit))
	}
	p.colorless -= amount.colorless
	if amount.generic > p.Total() {
		genericDeficit := amount.generic - (p.Total())
		manaDeficit = append(manaDeficit, fmt.Sprintf("{%d}", genericDeficit))
	}
	if len(manaDeficit) > 0 {
		return p, fmt.Errorf("cannot pay %s", strings.Join(manaDeficit, ""))
	}
	// Spend generic
	genericLeft := amount.generic
	if p.white > 0 && genericLeft > 0 {
		if genericLeft >= p.white {
			genericLeft -= p.white
			p.white = 0
		} else {
			p.white -= genericLeft
			genericLeft = 0
		}
	}
	for _, color := range colorsForGeneric {
		p, genericLeft = SpendGenericManaByColor(p, genericLeft, color)
	}
	if genericLeft > 0 {
		return p, fmt.Errorf("cannot pay {%d}", genericLeft)
	}
	return p, nil
}

func SpendGenericManaByColor(p Pool, amount int, color Color) (Pool, int) {
	if amount <= 0 {
		return p, amount
	}
	switch color {
	case White:
		if p.white > 0 && amount > 0 {
			if amount >= p.white {
				amount -= p.white
				p.white = 0
			} else {
				p.white -= amount
				amount = 0
			}
		}
	case Blue:
		if p.blue > 0 && amount > 0 {
			if amount >= p.blue {
				amount -= p.blue
				p.blue = 0
			} else {
				p.blue -= amount
				amount = 0
			}
		}
	case Black:
		if p.black > 0 && amount > 0 {
			if amount >= p.black {
				amount -= p.black
				p.black = 0
			} else {
				p.black -= amount
				amount = 0
			}
		}
	case Red:
		if p.red > 0 && amount > 0 {
			if amount >= p.red {
				amount -= p.red
				p.red = 0
			} else {
				p.red -= amount
				amount = 0
			}
		}
	case Green:
		if p.green > 0 && amount > 0 {
			if amount >= p.green {
				amount -= p.green
				p.green = 0
			} else {
				p.green -= amount
				amount = 0
			}
		}
	case Colorless:
		if p.colorless > 0 && amount > 0 {
			if amount >= p.colorless {
				amount -= p.colorless
				p.colorless = 0
			} else {
				p.colorless -= amount
				amount = 0
			}
		}
	default:
		panic(fmt.Sprintf("unknown color %s", color))
	}
	return p, amount
}

type Amount struct {
	generic   int
	white     int
	blue      int
	black     int
	red       int
	green     int
	colorless int
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

func (a Amount) WithAddedAmount(amount Amount) Amount {
	a.white += amount.white
	a.blue += amount.blue
	a.black += amount.black
	a.red += amount.red
	a.green += amount.green
	a.colorless += amount.colorless
	a.generic += amount.generic
	return a
}

func (a Amount) ManaString() string {
	descriptions := []string{}
	if a.colorless > 0 {
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

package mana

import (
	"fmt"
	"strings"
)

type Pool struct {
	white     int
	blue      int
	black     int
	red       int
	green     int
	colorless int
}

func (p Pool) Has(color Color) bool {
	switch color {
	case White:
		return p.white > 0
	case Blue:
		return p.blue > 0
	case Black:
		return p.black > 0
	case Red:
		return p.red > 0
	case Green:
		return p.green > 0
	case Colorless:
		return p.colorless > 0
	default:
		panic(fmt.Sprintf("unknown color %s", color))
	}
}

func (p Pool) AmountOf(color Color) int {
	switch color {
	case White:
		return p.white
	case Blue:
		return p.blue
	case Black:
		return p.black
	case Red:
		return p.red
	case Green:
		return p.green
	case Colorless:
		return p.colorless
	default:
		panic(fmt.Sprintf("unknown color %s", color))
	}
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

func (p Pool) WithAddAmount(amount Amount) Pool {
	p.white += amount.white
	p.blue += amount.blue
	p.black += amount.black
	p.red += amount.red
	p.green += amount.green
	p.colorless += amount.colorless
	p.colorless += amount.generic
	return p
}

func (p Pool) WithAddMana(amount int, color Color) Pool {
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

func (p Pool) WithSpendAmount(amount Amount, colorsForGeneric []Color) (Pool, Amount) {
	spendFromPool(&p.white, &amount.white)
	spendFromPool(&p.blue, &amount.blue)
	spendFromPool(&p.black, &amount.black)
	spendFromPool(&p.red, &amount.red)
	spendFromPool(&p.green, &amount.green)
	spendFromPool(&p.colorless, &amount.colorless)
	if amount.generic > p.Total() {
		amount.generic = amount.generic - (p.Total())
		p = Pool{}
	}
	for _, color := range colorsForGeneric {
		p, amount.generic = p.WithSpendMana(amount.generic, color)
	}
	return p, amount
}

func (p Pool) WithSpendMana(amount int, color Color) (Pool, int) {
	if amount <= 0 {
		return p, amount
	}
	switch color {
	case White:
		spendFromPool(&p.white, &amount)
	case Blue:
		spendFromPool(&p.blue, &amount)
	case Black:
		spendFromPool(&p.black, &amount)
	case Red:
		spendFromPool(&p.red, &amount)
	case Green:
		spendFromPool(&p.green, &amount)
	case Colorless:
		spendFromPool(&p.colorless, &amount)
	default:
		panic(fmt.Sprintf("unknown color %s", color))
	}
	return p, amount
}

func spendFromPool(poolColor *int, amount *int) {
	if *amount > *poolColor {
		*amount -= *poolColor
		*poolColor = 0
	} else {
		*poolColor -= *amount
		*amount = 0
	}
}

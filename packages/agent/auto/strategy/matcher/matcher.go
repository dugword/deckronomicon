package matcher

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
)

type Matcher interface {
	Matches(objs []query.Object) bool
}

type Name struct {
	Name string
}

func (m *Name) Matches(objs []query.Object) bool {
	return query.Contains(objs, has.Name(m.Name))
}

type CardType struct {
	CardType mtg.CardType
}

func (m *CardType) Matches(objs []query.Object) bool {
	return query.Contains(objs, has.CardType(m.CardType))
}

type Subtype struct {
	Subtype mtg.Subtype
}

func (m *Subtype) Matches(objs []query.Object) bool {
	return query.Contains(objs, has.Subtype(m.Subtype))
}

type And struct {
	Matchers []Matcher
}

func (m *And) Matches(objs []query.Object) bool {
	for _, matcher := range m.Matchers {
		if !matcher.Matches(objs) {
			return false
		}
	}
	return true
}

type Or struct {
	Matchers []Matcher
}

func (m *Or) Matches(objs []query.Object) bool {
	for _, matcher := range m.Matchers {
		if matcher.Matches(objs) {
			return true
		}
	}
	return false
}

type Not struct {
	Matcher Matcher
}

func (m *Not) Matches(objs []query.Object) bool {
	return !m.Matcher.Matches(objs)
}

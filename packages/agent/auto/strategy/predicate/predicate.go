package predicate

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
)

type Predicate interface {
	Matches(objs []query.Object) bool
	Select(objs []query.Object) []query.Object
}

type Matcher interface {
	Matches(objs []query.Object) bool
}

type Selector interface {
	Select(objs []query.Object) []query.Object
}

type Name struct {
	Name string
}

func (p *Name) Matches(objs []query.Object) bool {
	return query.Contains(objs, has.Name(p.Name))
}

func (p *Name) Select(objs []query.Object) []query.Object {
	return query.FindAll(objs, has.Name(p.Name))
}

type CardType struct {
	CardType mtg.CardType
}

func (p *CardType) Matches(objs []query.Object) bool {
	return query.Contains(objs, has.CardType(p.CardType))
}

func (p *CardType) Select(objs []query.Object) []query.Object {
	return query.FindAll(objs, has.CardType(p.CardType))
}

type Subtype struct {
	Subtype mtg.Subtype
}

func (p *Subtype) Matches(objs []query.Object) bool {
	return query.Contains(objs, has.Subtype(p.Subtype))
}

func (p *Subtype) Select(objs []query.Object) []query.Object {
	return query.FindAll(objs, has.Subtype(p.Subtype))
}

type And struct {
	Predicates []Predicate
}

func (p *And) Matches(objs []query.Object) bool {
	for _, predicate := range p.Predicates {
		if !predicate.Matches(objs) {
			return false
		}
	}
	return true
}

// TODO: Preserve order of predicates in selection
func (p *And) Select(objs []query.Object) []query.Object {
	if len(p.Predicates) == 0 {
		return nil
	}
	base := p.Predicates[0].Select(objs)
	idMap := map[string]query.Object{}
	for _, obj := range base {
		idMap[obj.ID()] = obj
	}
	for _, predicate := range p.Predicates[1:] {
		nextIDs := map[string]bool{}
		for _, obj := range predicate.Select(objs) {
			nextIDs[obj.ID()] = true
		}
		for id := range idMap {
			if !nextIDs[id] {
				delete(idMap, id)
			}
		}
	}
	var selected []query.Object
	for _, obj := range idMap {
		selected = append(selected, obj)
	}
	return selected
}

type Or struct {
	Predicates []Predicate
}

func (p *Or) Matches(objs []query.Object) bool {
	for _, predicate := range p.Predicates {
		if predicate.Matches(objs) {
			return true
		}
	}
	return false
}

func (p *Or) Select(objs []query.Object) []query.Object {
	seenIDs := map[string]bool{}
	var selected []query.Object
	for _, predicate := range p.Predicates {
		for _, obj := range predicate.Select(objs) {
			if !seenIDs[obj.ID()] {
				seenIDs[obj.ID()] = true
				selected = append(selected, obj)
			}
		}
	}
	return selected
}

type Not struct {
	Predicate Predicate
}

func (p *Not) Matches(objs []query.Object) bool {
	return !p.Predicate.Matches(objs)
}

func (p *Not) Select(objs []query.Object) []query.Object {
	excluded := map[string]bool{}
	for _, obj := range p.Predicate.Select(objs) {
		excluded[obj.ID()] = true
	}

	var selected []query.Object
	for _, obj := range objs {
		if !excluded[obj.ID()] {
			selected = append(selected, obj)
		}
	}
	return selected
}

package state

import (
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/add"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/remove"
	"deckronomicon/packages/query/take"
)

type Resolvable interface {
	Description() string
	ID() string
	Name() string
	EffectWithTargets() []*effect.EffectWithTarget
	Match(p query.Predicate) bool
	Controller() string
	Owner() string
	SourceID() string
}

type Stack struct {
	resolvables []Resolvable
}

func NewStack() *Stack {
	stack := &Stack{
		resolvables: []Resolvable{},
	}
	return stack
}

func (s *Stack) Add(resolvable Resolvable) *Stack {
	return &Stack{resolvables: add.Item(s.resolvables, resolvable)}
}
func (s *Stack) AddTop(resolvable Resolvable) *Stack {
	return &Stack{resolvables: add.Item([]Resolvable{resolvable}, s.resolvables...)}
}

func (s *Stack) Find(predicate query.Predicate) (Resolvable, bool) {
	return query.Find(s.resolvables, predicate)
}

func (s *Stack) FindAll(predicate query.Predicate) []Resolvable {
	return query.FindAll(s.resolvables, predicate)
}

func (s *Stack) Get(id string) (Resolvable, bool) {
	return query.Get(s.resolvables, id)
}

func (s *Stack) GetTop() (Resolvable, bool) {
	return query.GetTop(s.resolvables)
}

func (s *Stack) GetAll() []Resolvable {
	return s.resolvables
}

func (s *Stack) Name() string {
	return string(mtg.ZoneStack)
}

func (s *Stack) Remove(id string) (*Stack, bool) {
	resolvables, ok := remove.By(s.resolvables, has.ID(id))
	if !ok {
		return nil, false
	}
	return &Stack{resolvables: resolvables}, true

}

func (s *Stack) Size() int {
	return len(s.resolvables)
}

func (s *Stack) Take(id string) (Resolvable, *Stack, bool) {
	resolvable, resolvables, ok := take.By(s.resolvables, has.ID(id))
	if !ok {
		return nil, nil, false
	}
	return resolvable, &Stack{resolvables: resolvables}, true
}

func (s *Stack) TakeBy(predicate query.Predicate) (Resolvable, *Stack, bool) {
	resolvable, resolvables, ok := take.By(s.resolvables, predicate)
	if !ok {
		return nil, nil, false
	}
	return resolvable, &Stack{resolvables: resolvables}, true
}

func (s *Stack) TakeTop() (Resolvable, *Stack, bool) {
	resolvable, resolvables, ok := take.Top(s.resolvables)
	if !ok {
		return nil, nil, false
	}
	return resolvable, &Stack{resolvables: resolvables}, true
}

func (s *Stack) ZoneType() string {
	return "Stack"
}

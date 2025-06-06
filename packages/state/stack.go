package state

import (
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
	Resolve(game Game, player Player) error
	Match(p query.Predicate) bool
}

type Stack struct {
	resolvables []Resolvable
}

func NewStack() Stack {
	stack := Stack{
		resolvables: []Resolvable{},
	}
	return stack
}

func (s Stack) Add(resolvable Resolvable) (Stack, bool) {
	return Stack{resolvables: add.Item(s.resolvables, resolvable)}, true
}

func (s Stack) Get(id string) (Resolvable, bool) {
	for _, resolvable := range s.resolvables {
		if resolvable.ID() == id {
			return resolvable, true
		}
	}
	return nil, false
}

func (s Stack) GetAll() []Resolvable {
	return s.resolvables
}

func (s Stack) Name() string {
	return string(mtg.ZoneStack)
}

func (s Stack) Remove(id string) (Stack, bool) {
	resolvables, ok := remove.By(s.resolvables, has.ID(id))
	if !ok {
		return s, false
	}
	return Stack{resolvables: resolvables}, true

}

func (s Stack) Take(id string) (Resolvable, Stack, bool) {
	resolvable, resolvables, ok := take.By(s.resolvables, has.ID(id))
	if !ok {
		return nil, s, false
	}
	return resolvable, Stack{resolvables: resolvables}, true
}

func (s Stack) Size() int {
	return len(s.resolvables)
}

func (s Stack) ZoneType() string {
	return "Stack"
}

func (s Stack) TakeTop() (Resolvable, Stack, bool) {
	resolvable, resolvables, ok := take.Top(s.resolvables)
	if !ok {
		return nil, s, false
	}
	return resolvable, Stack{resolvables: resolvables}, true
}

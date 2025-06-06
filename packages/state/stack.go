package state

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"errors"
	"fmt"
)

type Resolvable interface {
	Description() string
	ID() string
	Name() string
	Resolve(game Game, player Player) error
	Match(p query.Predicate) bool
}

type Stack struct {
	stack []Resolvable
}

func NewStack() Stack {
	stack := Stack{
		stack: []Resolvable{},
	}
	return stack
}

func (s Stack) Add(resolvable Resolvable) Stack {
	s.stack = append(s.stack, resolvable)
	return s
}

func (s Stack) Get(id string) (Resolvable, error) {
	for _, resolvable := range s.stack {
		if resolvable.ID() == id {
			return resolvable, nil
		}
	}
	return nil, fmt.Errorf("object with ID %s not found in stack", id)
}

func (s Stack) GetAll() []Resolvable {
	return s.stack
}

func (s Stack) Name() string {
	return string(mtg.ZoneStack)
}

func (s Stack) Remove(id string) error {
	for i, resolvable := range s.stack {
		if resolvable.ID() == id {
			s.stack = append(s.stack[:i], s.stack[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("object with ID %s not found in stack", id)
}

func (s Stack) Take(id string) (Resolvable, error) {
	for i, resolvable := range s.stack {
		if resolvable.ID() == id {
			s.stack = append(s.stack[:i], s.stack[i+1:]...)
			return resolvable, nil
		}
	}
	return nil, fmt.Errorf("object with ID %s not found in stack", id)
}

func (s Stack) Size() int {
	return len(s.stack)
}

func (s Stack) ZoneType() string {
	return "Stack"
}

func (s Stack) Pop() (Resolvable, Stack, error) {
	if len(s.stack) == 0 {
		return nil, s, errors.New("stack is empty")
	}
	top := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return top, s, nil
}

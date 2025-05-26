package game

import (
	"errors"
	"fmt"
)

type Resolvable interface {
	Name() string
	Description() string
	Resolve(*GameState, *Player) error
}

type Stack struct {
	stack []Resolvable
}

func NewStack() *Stack {
	return &Stack{
		stack: []Resolvable{},
	}
}

func (s *Stack) Add(object GameObject) error {
	resolvable, ok := object.(Resolvable)
	if !ok {
		return fmt.Errorf("object %s is not resolvable", object.ID())
	}
	s.stack = append(s.stack, resolvable)
	return nil
}

func (s *Stack) AvailableActivatedAbilities(*GameState) []*ActivatedAbility {
	return nil
}

func (s *Stack) AvailableToPlay(*GameState) []GameObject {
	return nil
}

func (s *Stack) Find(id string) (GameObject, error) {
	return nil, nil
}

func (s *Stack) FindByName(name string) (GameObject, error) {
	return nil, nil
}

func (s *Stack) FindAllBySubtype(subtype Subtype) []GameObject {
	return nil
}

func (s *Stack) Get(id string) (GameObject, error) {
	return nil, nil
}

func (s *Stack) GetAll() []GameObject {
	var objects []GameObject
	for _, obj := range s.stack {
		gameObject, ok := obj.(GameObject)
		if !ok {
			fmt.Printf("Object %s is not a GameObject\n", obj.Description())
		}
		objects = append(objects, gameObject)
	}
	return objects
}

func (s *Stack) Remove(id string) error {
	return nil
}

func (s *Stack) Take(id string) (GameObject, error) {
	return nil, nil
}

func (s *Stack) Size() int {
	return len(s.stack)
}

func (s *Stack) ZoneType() string {
	return "Stack"
}

func (s *Stack) Pop() (GameObject, error) {
	if len(s.stack) == 0 {
		return nil, errors.New("stack is empty")
	}
	top := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	gameObject, ok := top.(GameObject)
	if !ok {
		return nil, fmt.Errorf("top of stack is not a GameObject: %s", top.Description())
	}
	return gameObject, nil
}

package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/gob"
	"fmt"
)

type Stack struct {
	items []StackItem
}

type StackItem struct {
}

func (s *StackItem) Resolve(state *GameState) error {
	return nil
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) Pop() (*StackItem, error) {
	if len(s.items) == 0 {
		return nil, fmt.Errorf("stack is empty")
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return &item, nil
}

func (s *Stack) Push(item StackItem) {
	s.items = append(s.items, item)
}

type Engine struct {
	agents map[string]PlayerAgent
	state  *GameState
	record *GameRecord
	rng    *RNG
}

type EngineConfig struct {
	Players []*Player
	Agents  map[string]PlayerAgent
	Seed    int64
	// Cards are just strings for now, but will be a Card type later
	DeckLists map[string][]*gob.Card
}

func NewEngine(config EngineConfig) *Engine {
	agents := map[string]PlayerAgent{}
	for id, agent := range config.Agents {
		agents[id] = agent
	}
	return &Engine{
		agents: agents,
		state: NewGameState(
			GameStateConfig{
				Players: config.Players,
			},
		),
		record: NewGameRecord(config.Seed),
		rng:    NewRNG(config.Seed),
	}
}

func (e *Engine) RunGame() error {
	fmt.Println("Running game...")
	// shuffle all player decks
	// draw initial hands for players
	// resovle mulligans
	for !e.state.IsGameOver() {
		fmt.Println("Starting turn for player:", e.state.ActivePlayer())
		err := e.RunTurn()
		if err != nil {
			return err
		}
		nextPlayer := e.state.NextPlayer()
		e.Apply(event.NewChangeActivePlayerEvent(nextPlayer.id))
	}
	e.Apply(event.NewGameOverEvent("?"))
	return nil
}

func (e *Engine) RunTurn() error {
	fmt.Println("Running turn for player:", e.state.ActivePlayer())
	for _, phase := range GamePhases() {
		if err := e.RunPhase(phase); err != nil {
			return fmt.Errorf("error running phase %s: %w", phase.name, err)
		}
	}
	// This is a place holder, phases will be a struct later
	return nil
}

func (e *Engine) RunPhase(phase GamePhase) error {
	fmt.Println("Running phase:", phase.name)
	for _, step := range phase.steps {
		if err := e.RunStep(step); err != nil {
			return fmt.Errorf("error running step %s: %w", step.name, err)
		}
	}
	return nil
}

func (e *Engine) RunStep(step GameStep) error {
	fmt.Println("Running step:", step.name)
	for _, action := range step.actions {
		fmt.Println("Completing action:", action.Name())
		action.Complete(choose.Choice{})
	}
	for _, playerID := range e.state.PlayersInTurnOrder() {
		fmt.Println("Starting priority for player:", playerID)
		e.RunPriority(playerID)
		for {
			if e.state.stack.IsEmpty() {
				break
			}
			if stackItem, err := e.state.stack.Pop(); err != nil {
				_ = stackItem.Resolve(e.state)
			}
			e.RunPriority(playerID)
		}
	}
	return nil
}

func (e *Engine) RunPriority(playerID string) {
	for {
		fmt.Println("Running priority for player:", playerID)
		agent := e.agents[playerID]
		action := agent.GetNextAction()
		if action == "pass" {
			fmt.Println("Player", playerID, "passed priority.")
			e.PassPriority(playerID)
			break
		}
	}
}

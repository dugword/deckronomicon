package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
	"fmt"
)

type mockLogger struct{}

func (m *mockLogger) Debug(...any)                         {}
func (m *mockLogger) Debugf(format string, args ...any)    {}
func (m *mockLogger) Info(...any)                          {}
func (m *mockLogger) Infof(format string, args ...any)     {}
func (m *mockLogger) Warn(...any)                          {}
func (m *mockLogger) Warnf(format string, args ...any)     {}
func (m *mockLogger) Error(...any)                         {}
func (m *mockLogger) Errorf(format string, args ...any)    {}
func (m *mockLogger) Critical(...any)                      {}
func (m *mockLogger) Criticalf(format string, args ...any) {}

type mockPlayerAgent struct {
	playerID string
}

type mockAction struct {
	name string
}

func (a mockAction) Name() string {
	return a.name
}

func (a mockAction) Complete(game state.Game, player state.Player, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	return []event.GameEvent{
		event.PassPriorityEvent{
			PlayerID: player.ID(),
		},
	}, nil
}

func (m *mockPlayerAgent) GetNextAction(state.Game) (Action, error) {
	return &mockAction{
		name: "Mock Pass Priority Action",
	}, nil
}

func (m *mockPlayerAgent) ReportState(state.Game) error {
	return nil
}

func (m *mockPlayerAgent) Choose(prompt choose.ChoicePrompt) (choose.ChoiceResults, error) {
	switch prompt.ChoiceOpts.(type) {
	case choose.ChooseOneOpts:
		return choose.ChooseOneResults{}, nil
	case choose.ChooseManyOpts:
		return choose.ChooseManyResults{}, nil
	case choose.ChooseNumberOpts:
		return choose.ChooseNumberResults{}, nil
	case choose.MapChoicesToBucketsOpts:
		return choose.MapChoicesToBucketsResults{}, nil
	default:
		return nil, fmt.Errorf("unknown choice prompt type %T", prompt.ChoiceOpts)
	}
}

func (m *mockPlayerAgent) PlayerID() string {
	return m.playerID
}

package ui

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/view"
	"slices"
	"testing"
)

func TestNewBuffer(t *testing.T) {
	buffer := NewBuffer("./testdata/display.tmpl")
	if buffer == nil {
		t.Errorf("NewBuffer() = nil, want non-nil")
	}
}

func TestBufferUpdate(t *testing.T) {
	buffer := NewBuffer("./testdata/display.tmpl")
	buffer.Update(
		&view.Game{ActivePlayerID: "A"},
		&view.Player{ID: "A"},
		&view.Player{ID: "B"},
	)
	if buffer.displayData == nil &&
		!slices.Contains(buffer.displayData.GameStatusData.Content, "Active Player: A") &&
		!slices.Contains(buffer.displayData.PlayerData.Content, "ID: A") &&
		!slices.Contains(buffer.displayData.OpponentData.Content, "ID: B") {
		t.Errorf("Buffer.Update() did not set display data correctly")
	}
}

func TestBufferUpdateChoices(t *testing.T) {
	buffer := NewBuffer("./testdata/display.tmpl")
	t.Run("No choices", func(t *testing.T) {
		buffer.UpdateChoices("Make a Choice", nil)
		if buffer.displayData == nil &&
			buffer.displayData.ChoiceData.Title != "Nothing to choose" {
			t.Errorf("Buffer.UpdateChoices() did not set display data correctly for no choices")
		}
	})
	t.Run("Multiple choices", func(t *testing.T) {
		buffer.UpdateChoices("Make a Choice", []choose.Choice{
			choose.NewGenericChoice("Choice 1", "c1"),
			choose.NewGenericChoice("Choice 2", "c2"),
			choose.NewGenericChoice("Choice 3", "c3"),
		})
		if buffer.displayData == nil &&
			!slices.Contains(buffer.displayData.ChoiceData.Content, "Choice 1") &&
			!slices.Contains(buffer.displayData.ChoiceData.Content, "Choice 2") &&
			!slices.Contains(buffer.displayData.ChoiceData.Content, "Choice 3") {
			t.Errorf("Buffer.UpdateChoices() did not set display data correctly")
		}
	})
}

func TestUpdateMessage(t *testing.T) {
	buffer := NewBuffer("./testdata/display.tmpl")
	buffer.UpdateMessage([]string{"New message"})
	if buffer.displayData == nil &&
		!slices.Contains(buffer.displayData.MessageData.Content, "New message") {
		t.Errorf("Buffer.UpdateMessage() did not set display data correctly")
	}
}

func TestBattlefieldData(t *testing.T) {
	permanents := []*view.Permanent{
		{ID: "p1", Name: "Elf", Controller: "A", Tapped: false, SummoningSick: false},
		{ID: "p2", Name: "Goblin", Controller: "B", Tapped: true, SummoningSick: true},
	}
	box := BattlefieldData(permanents)
	if box.Title != "Battlefield" {
		t.Errorf("expected title 'Battlefield', got %q", box.Title)
	}
	if len(box.Content) != 2 {
		t.Errorf("expected 2 lines, got %d", len(box.Content))
	}
	if got := box.Content[0]; got == "" {
		t.Errorf("expected non-empty line for first permanent")
	}
	want := "[A] Elf <p1>"
	if box.Content[0] != want {
		t.Errorf("expected first line %q, got %q", want, box.Content[0])
	}
}

func TestStackData(t *testing.T) {
	stack := []*view.Resolvable{
		{ID: "s1", Name: "Lightning Bolt", Controller: "A"},
		{ID: "s2", Name: "Counterspell", Controller: "B"},
	}
	box := StackData(stack)
	if box.Title != "Stack" {
		t.Errorf("expected title 'Stack', got %q", box.Title)
	}
	if len(box.Content) != 2 {
		t.Errorf("expected 2 lines, got %d", len(box.Content))
	}
	want := "[A] Lightning Bolt <s1>"
	if box.Content[0] != want {
		t.Errorf("expected first line %q, got %q", want, box.Content[0])
	}
}

func TestPlayerData(t *testing.T) {
	player := view.Player{
		ID:          "p1",
		Life:        20,
		Hand:        []*view.Card{{ID: "c1", Name: "Island"}},
		LibrarySize: 40,
		Graveyard:   []*view.Card{{ID: "c2", Name: "Mountain"}},
	}
	box := PlayerData(&player)
	if box.Title != "Status for p1" {
		t.Errorf("expected title 'Status for p1', got %q", box.Title)
	}
	if len(box.Content) != 8 {
		t.Errorf("expected 8 lines, got %d", len(box.Content))
	}
	if got := box.Content[0]; got != "Turn: 0" {
		t.Errorf("expected first line 'Turn: 0', got %q", got)
	}
}

func TestGraveyardData(t *testing.T) {
	cards := []*view.Card{
		{ID: "g1", Name: "Swamp"},
		{ID: "g2", Name: "Forest"},
	}
	box := CardListData("Graveyard", cards)
	if box.Title != "Graveyard" {
		t.Errorf("expected title 'Graveyard', got %q", box.Title)
	}
	if len(box.Content) != 2 {
		t.Errorf("expected 2 lines, got %d", len(box.Content))
	}
	want := "Swamp <g1>"
	if box.Content[0] != want {
		t.Errorf("expected first line %q, got %q", want, box.Content[0])
	}
}

func TestHandData(t *testing.T) {
	cards := []*view.Card{
		{ID: "h1", Name: "Plains"},
		{ID: "h2", Name: "Island"},
	}
	box := CardListData("Hand", cards)
	if box.Title != "Hand" {
		t.Errorf("expected title 'Hand', got %q", box.Title)
	}
	if len(box.Content) != 2 {
		t.Errorf("expected 2 lines, got %d", len(box.Content))
	}
	want := "Plains <h1>"
	if box.Content[0] != want {
		t.Errorf("expected first line %q, got %q", want, box.Content[0])
	}
}

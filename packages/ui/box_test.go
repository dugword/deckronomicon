package ui

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSplitLines(t *testing.T) {
	want := []string{"abc", "def", "ghi", "j"}
	text := "abcdefghij"
	got := splitLines(text, 3)
	if len(got) != len(want) {
		t.Errorf("len(text) = %d, want %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("text[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestCreateBox(t *testing.T) {
	want := Box{
		TopLine:    "╔═══════╗",
		Title:      "║ Title ║",
		MiddleLine: "╠═══════╣",
		Lines: []string{
			"║ One   ║",
			"║ Two   ║",
		},
		BottomLine: "╚═══════╝",
	}
	data := BoxData{
		Title:   "Title",
		Content: []string{"One", "Two"},
	}
	got := CreateBox(data)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("CreateBox() mismatch (-want +got):\n%s", diff)
	}
}

func TestCombineBoxesSideBySide(t *testing.T) {
	t.Run("Equal height boxes", func(t *testing.T) {
		want := Box{
			TopLine:    "╔═══════╗  ╔═══════╗",
			Title:      "║ Left  ║  ║ Right ║",
			MiddleLine: "╠═══════╣  ╠═══════╣",
			Lines: []string{
				"║ Line1 ║  ║ Line2 ║",
				"║ Line3 ║  ║ Line4 ║",
			},
			BottomLine: "╚═══════╝  ╚═══════╝",
		}
		left := CreateBox(BoxData{Title: "Left", Content: []string{"Line1", "Line3"}})
		right := CreateBox(BoxData{Title: "Right", Content: []string{"Line2", "Line4"}})
		got := CombineBoxesSideBySide(left, right)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("CombineBoxesSideBySide() mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("Left box taller", func(t *testing.T) {
		want := Box{
			TopLine:    "╔═══╗  ╔═══╗",
			Title:      "║ L ║  ║ R ║",
			MiddleLine: "╠═══╣  ╠═══╣",
			Lines: []string{
				"║ A ║  ║ B ║",
				"║ C ║  ╚═══╝",
			},
			BottomLine: "╚═══╝       ",
		}
		left := CreateBox(BoxData{Title: "L", Content: []string{"A", "C"}})
		right := CreateBox(BoxData{Title: "R", Content: []string{"B"}})
		got := CombineBoxesSideBySide(left, right)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("CombineBoxesSideBySide() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Right box taller", func(t *testing.T) {
		want := Box{
			TopLine:    "╔═══╗  ╔═══╗",
			Title:      "║ L ║  ║ R ║",
			MiddleLine: "╠═══╣  ╠═══╣",
			Lines: []string{
				"║ A ║  ║ B ║",
				"╚═══╝  ║ C ║",
			},
			BottomLine: "       ╚═══╝",
		}
		left := CreateBox(BoxData{Title: "L", Content: []string{"A"}})
		right := CreateBox(BoxData{Title: "R", Content: []string{"B", "C"}})
		got := CombineBoxesSideBySide(left, right)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("CombineBoxesSideBySide() mismatch (-want +got):\n%s", diff)
		}
	})
}

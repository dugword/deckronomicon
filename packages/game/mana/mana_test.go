package mana

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStringToColor(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Color
	}{
		{
			name:  "with  W",
			input: "W",
			want:  White,
		},
		{
			name:  "with {W}",
			input: "{W}",
			want:  White,
		},
		{
			name:  "with  U",
			input: "U",
			want:  Blue,
		},
		{
			name:  "with {U}",
			input: "{U}",
			want:  Blue,
		},
		{
			name:  "with  B",
			input: "B",
			want:  Black,
		},
		{
			name:  "with {B}",
			input: "{B}",
			want:  Black,
		},
		{
			name:  "with  R",
			input: "R",
			want:  Red,
		},
		{
			name:  "with {R}",
			input: "{R}",
			want:  Red,
		},
		{
			name:  "with  G",
			input: "G",
			want:  Green,
		},
		{
			name:  "with {G}",
			input: "{G}",
			want:  Green,
		},
		{
			name:  "with  C",
			input: "C",
			want:  Colorless,
		},
		{
			name:  "with {C}",
			input: "{C}",
			want:  Colorless,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := StringToColor(test.input)
			if !ok {
				t.Fatalf("StringToColor(%q) returned false", test.input)
			}
			if got != test.want {
				t.Errorf("StringToColor(%q) = %v; want %v", test.input, got, test.want)
			}
		})
	}
}

func TestParseManaString(t *testing.T) {
	tests := []struct {
		name       string
		manaString string
		want       Amount
	}{
		{
			name:       "with color",
			manaString: "{U}",
			want:       Amount{blue: 1},
		},
		{
			name:       "with colorless",
			manaString: "{C}",
			want:       Amount{colorless: 1},
		},
		{
			name:       "with WUBRG",
			manaString: "{W}{U}{B}{R}{G}",
			want:       Amount{white: 1, blue: 1, black: 1, red: 1, green: 1},
		},
		{
			name:       "with double WUBRG",
			manaString: "{W}{U}{B}{R}{G}{W}{U}{B}{R}{G}",
			want:       Amount{white: 2, blue: 2, black: 2, red: 2, green: 2},
		},
		{
			name:       "with generic and WUBRGC",
			manaString: "{2}{C}{W}{U}{B}{R}{G}",
			want:       Amount{generic: 2, colorless: 1, white: 1, blue: 1, black: 1, red: 1, green: 1},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseManaString(test.manaString)
			if err != nil {
				t.Fatalf("ParseManaString(%q) returned error: %v", test.manaString, err)
			}
			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Amount{})); diff != "" {
				t.Errorf("ParseManaString(%q) mismatch (-want +got):\n%s", test.manaString, diff)
			}
		})
	}
}

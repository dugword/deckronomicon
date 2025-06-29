package mana

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAmountWithAddAmount(t *testing.T) {
	tests := []struct {
		name       string
		initial    Amount
		additional Amount
		want       Amount
	}{
		{
			name: "with color",
			initial: Amount{
				blue: 1,
			},
			additional: Amount{
				blue: 1,
			},
			want: Amount{
				blue: 2,
			},
		},
		{
			name: "with colorless",
			initial: Amount{
				colorless: 1,
			},
			additional: Amount{
				colorless: 1,
			},
			want: Amount{
				colorless: 2,
			},
		},
		{
			name: "with generic",
			initial: Amount{
				generic: 1,
			},
			additional: Amount{
				generic: 1,
			},
			want: Amount{
				generic: 2,
			},
		},
		{
			name: "with generic colorless and wuburg",
			initial: Amount{
				generic:   1,
				colorless: 1,
				white:     1,
				blue:      1,
				black:     1,
				red:       1,
				green:     1,
			},
			additional: Amount{
				generic:   1,
				colorless: 1,
				white:     1,
				blue:      1,
				black:     1,
				red:       1,
				green:     1,
			},
			want: Amount{
				generic:   2,
				colorless: 2,
				white:     2,
				blue:      2,
				black:     2,
				red:       2,
				green:     2,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.initial.WithAddAmount(test.additional)
			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Amount{})); diff != "" {
				t.Errorf("{%s}.WithAddAmount(%s) mismatch (-want +got):\n%s", test.initial.ManaString(), test.additional.ManaString(), diff)
			}
		})
	}
}

func TestAmountWithAddMana(t *testing.T) {
	tests := []struct {
		name   string
		amount Amount
		color  Color
		amt    int
		want   Amount
	}{
		{
			name:   "with 1 colorless",
			amount: Amount{},
			color:  Colorless,
			amt:    1,
			want:   Amount{colorless: 1},
		},
		{
			name:   "with 1 white",
			amount: Amount{},
			color:  White,
			amt:    1,
			want:   Amount{white: 1},
		},
		{
			name:   "with 1 blue",
			amount: Amount{},
			color:  Blue,
			amt:    1,
			want:   Amount{blue: 1},
		},
		{
			name:   "with 1 black",
			amount: Amount{},
			color:  Black,
			amt:    1,
			want:   Amount{black: 1},
		},
		{
			name:   "with 1 red",
			amount: Amount{},
			color:  Red,
			amt:    1,
			want:   Amount{red: 1},
		},
		{
			name:   "with 1 green",
			amount: Amount{},
			color:  Green,
			amt:    1,
			want:   Amount{green: 1},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.amount.WithAddMana(test.amt, test.color)
			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Amount{})); diff != "" {
				t.Errorf("WithAddMana(%d, %s) mismatch (-want +got):\n%s", test.amt, test.color, diff)
			}
		})
	}
}

func TestAmountManaString(t *testing.T) {
	tests := []struct {
		name string
		amt  Amount
		want string
	}{
		{
			name: "with color",
			amt:  Amount{blue: 1},
			want: "{U}",
		},
		{
			name: "with colorless",
			amt:  Amount{colorless: 1},
			want: "{C}",
		},
		{
			name: "with generic colorless and wubrg",
			amt:  Amount{generic: 1, colorless: 1, white: 1, blue: 1, black: 1, red: 1, green: 1},
			want: "{1}{C}{W}{U}{B}{R}{G}",
		},
		{
			name: "with 2 generic 2 colorless and double wubrg",
			amt:  Amount{generic: 2, colorless: 2, white: 2, blue: 2, black: 2, red: 2, green: 2},
			want: "{2}{C}{C}{W}{W}{U}{U}{B}{B}{R}{R}{G}{G}",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.amt.ManaString()
			if got != test.want {
				t.Errorf("Amount.ManaString() = %q; want %q", got, test.want)
			}
		})
	}
}

func TestAmountTotal(t *testing.T) {
	tests := []struct {
		name   string
		amount Amount
		want   int
	}{
		{
			name:   "with color",
			amount: Amount{blue: 1},
			want:   1,
		},
		{
			name:   "with colorless",
			amount: Amount{colorless: 1},
			want:   1,
		},
		{
			name:   "with wubrg",
			amount: Amount{white: 1, blue: 1, black: 1, red: 1, green: 1},
			want:   5,
		},
		{
			name:   "with 2 generic colorless and wubrg",
			amount: Amount{generic: 2, colorless: 1, white: 1, blue: 1, black: 1, red: 1, green: 1},
			want:   8,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.amount.Total()
			if got != test.want {
				t.Errorf("Amount.Total() = %d; want %d", got, test.want)
			}
		})
	}
}

func TestAmountGetters(t *testing.T) {
	t.Run("with generic", func(t *testing.T) {
		got := Amount{generic: 1}.Generic()
		want := 1
		if got != want {
			t.Errorf("Amount.Generic() = %d; want %d", got, want)
		}
	})
	t.Run("with colorless", func(t *testing.T) {
		got := Amount{colorless: 1}.Colorless()
		want := 1
		if got != want {
			t.Errorf("Amount.Colorless() = %d; want %d", got, want)
		}
	})
	t.Run("with white", func(t *testing.T) {
		got := Amount{white: 1}.White()
		want := 1
		if got != want {
			t.Errorf("Amount.White() = %d; want %d", got, want)
		}
	})
	t.Run("with blue", func(t *testing.T) {
		got := Amount{blue: 1}.Blue()
		want := 1
		if got != want {
			t.Errorf("Amount.Blue() = %d; want %d", got, want)
		}
	})
	t.Run("with black", func(t *testing.T) {
		got := Amount{black: 1}.Black()
		want := 1
		if got != want {
			t.Errorf("Amount.Black() = %d; want %d", got, want)
		}
	})
	t.Run("with red", func(t *testing.T) {
		got := Amount{red: 1}.Red()
		want := 1
		if got != want {
			t.Errorf("Amount.Red() = %d; want %d", got, want)
		}
	})
	t.Run("with green", func(t *testing.T) {
		got := Amount{green: 1}.Green()
		want := 1
		if got != want {
			t.Errorf("Amount.Green() = %d; want %d", got, want)
		}
	})
}

func TestAmountHas(t *testing.T) {
	tests := []struct {
		name   string
		amount Amount
		color  Color
		want   bool
	}{
		{
			name:   "with has white",
			amount: Amount{white: 1},
			color:  White,
			want:   true,
		},
		{
			name:   "with has blue",
			amount: Amount{blue: 1},
			color:  Blue,
			want:   true,
		},
		{
			name:   "with has black",
			amount: Amount{black: 1},
			color:  Black,
			want:   true,
		},
		{
			name:   "with has red",
			amount: Amount{red: 1},
			color:  Red,
			want:   true,
		},
		{
			name:   "with has green",
			amount: Amount{green: 1},
			color:  Green,
			want:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.amount.Has(test.color)
			if got != test.want {
				t.Errorf("Amount.Has(%s) = %v; want %v", test.color, got, test.want)
			}
		})
	}
}

func TestAmountAmountOf(t *testing.T) {
	tests := []struct {
		name   string
		amount Amount
		color  Color
		want   int
	}{
		{
			name:   "with white",
			amount: Amount{white: 1},
			color:  White,
			want:   1,
		},
		{
			name:   "with blue",
			amount: Amount{blue: 1},
			color:  Blue,
			want:   1,
		},
		{
			name:   "with black",
			amount: Amount{black: 1},
			color:  Black,
			want:   1,
		},
		{
			name:   "with red",
			amount: Amount{red: 1},
			color:  Red,
			want:   1,
		},
		{
			name:   "with green",
			amount: Amount{green: 1},
			color:  Green,
			want:   1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.amount.AmountOf(test.color)
			if got != test.want {
				t.Errorf("Amount.AmountOf(%s) = %d; want %d", test.color, got, test.want)
			}
		})
	}
}

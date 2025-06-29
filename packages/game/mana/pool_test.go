package mana

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPoolTotal(t *testing.T) {
	tests := []struct {
		name string
		Pool Pool
		want int
	}{
		{
			name: "with color",
			Pool: Pool{blue: 1},
			want: 1,
		},
		{
			name: "with colorless",
			Pool: Pool{colorless: 1},
			want: 1,
		},
		{
			name: "with wurbrg",
			Pool: Pool{white: 1, blue: 1, black: 1, red: 1, green: 1},
			want: 5,
		},
		{
			name: "with 2 colorless and double wurbrg",
			Pool: Pool{colorless: 2, white: 2, blue: 2, black: 2, red: 2, green: 2},
			want: 12,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.Pool.Total()
			if got != test.want {
				t.Errorf("Pool.Total() = %d; want %d", got, test.want)
			}
		})
	}
}

func TestPoolGetters(t *testing.T) {
	t.Run("with colorless", func(t *testing.T) {
		got := Pool{colorless: 1}.Colorless()
		want := 1
		if got != want {
			t.Errorf("Pool.Colorless() = %d; want %d", got, want)
		}
	})
	t.Run("with white", func(t *testing.T) {
		got := Pool{white: 1}.White()
		want := 1
		if got != want {
			t.Errorf("Pool.White() = %d; want %d", got, want)
		}
	})
	t.Run("with blue", func(t *testing.T) {
		got := Pool{blue: 1}.Blue()
		want := 1
		if got != want {
			t.Errorf("Pool.Blue() = %d; want %d", got, want)
		}
	})
	t.Run("with black", func(t *testing.T) {
		got := Pool{black: 1}.Black()
		want := 1
		if got != want {
			t.Errorf("Pool.Black() = %d; want %d", got, want)
		}
	})
	t.Run("with red", func(t *testing.T) {
		got := Pool{red: 1}.Red()
		want := 1
		if got != want {
			t.Errorf("Pool.Red() = %d; want %d", got, want)
		}
	})
	t.Run("with green", func(t *testing.T) {
		got := Pool{green: 1}.Green()
		want := 1
		if got != want {
			t.Errorf("Pool.Green() = %d; want %d", got, want)
		}
	})
}

func TestPoolManaString(t *testing.T) {
	tests := []struct {
		name string
		pool Pool
		want string
	}{
		{
			name: "with empty pool",
			pool: Pool{},
			want: "",
		},
		{
			name: "with 1 colorless",
			pool: Pool{colorless: 1},
			want: "{C}",
		},
		{
			name: "with 1 white",
			pool: Pool{white: 1},
			want: "{W}",
		},
		{
			name: "with 1 blue",
			pool: Pool{blue: 1},
			want: "{U}",
		},
		{
			name: "with 1 black",
			pool: Pool{black: 1},
			want: "{B}",
		},
		{
			name: "with 1 red",
			pool: Pool{red: 1},
			want: "{R}",
		},
		{
			name: "with 1 green",
			pool: Pool{green: 1},
			want: "{G}",
		},
		{
			name: "with 2 colorless double wubrg",
			pool: Pool{colorless: 2, white: 2, blue: 2, black: 2, red: 2, green: 2},
			want: "{C}{C}{W}{W}{U}{U}{B}{B}{R}{R}{G}{G}",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.pool.ManaString()
			if got != test.want {
				t.Errorf("ManaString() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestPoolWithAddMana(t *testing.T) {
	tests := []struct {
		name   string
		pool   Pool
		amount int
		color  Color
		want   Pool
	}{
		{
			name:   "with 1 colorless",
			pool:   Pool{},
			amount: 1,
			color:  Colorless,
			want:   Pool{colorless: 1},
		},
		{
			name:   "with 1 white",
			pool:   Pool{},
			amount: 1,
			color:  White,
			want:   Pool{white: 1},
		},
		{
			name:   "with 1 blue",
			pool:   Pool{},
			amount: 1,
			color:  Blue,
			want:   Pool{blue: 1},
		},
		{
			name:   "with 1 black",
			pool:   Pool{},
			amount: 1,
			color:  Black,
			want:   Pool{black: 1},
		},
		{
			name:   "with 1 red",
			pool:   Pool{},
			amount: 1,
			color:  Red,
			want:   Pool{red: 1},
		},
		{
			name:   "with 1 green",
			pool:   Pool{},
			amount: 1,
			color:  Green,
			want:   Pool{green: 1},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.pool.WithAddMana(test.amount, test.color)
			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Pool{})); diff != "" {
				t.Errorf("WithAddMana() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPoolWithAddAmount(t *testing.T) {
	tests := []struct {
		name   string
		pool   Pool
		amount Amount
		want   Pool
	}{
		{
			name:   "with 1 colorless",
			pool:   Pool{},
			amount: Amount{colorless: 1},
			want:   Pool{colorless: 1},
		},
		{
			name:   "with 1 white",
			pool:   Pool{},
			amount: Amount{white: 1},
			want:   Pool{white: 1},
		},
		{
			name:   "with 1 blue",
			pool:   Pool{},
			amount: Amount{blue: 1},
			want:   Pool{blue: 1},
		},
		{
			name:   "with 1 black",
			pool:   Pool{},
			amount: Amount{black: 1},
			want:   Pool{black: 1},
		},
		{
			name:   "with 1 red",
			pool:   Pool{},
			amount: Amount{red: 1},
			want:   Pool{red: 1},
		},
		{
			name:   "with 1 green",
			pool:   Pool{},
			amount: Amount{green: 1},
			want:   Pool{green: 1},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.pool.WithAddAmount(test.amount)
			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Pool{})); diff != "" {
				t.Errorf("WithAddAmount() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPoolWithSpendMana(t *testing.T) {
	tests := []struct {
		name        string
		pool        Pool
		spendAmount int
		spendColor  Color
		wantPool    Pool
		wantAmount  int
		colors      []Color
	}{
		{
			name:        "with 1 colorless spend 1 colorless",
			pool:        Pool{colorless: 1},
			spendAmount: 1,
			spendColor:  Colorless,
			wantPool:    Pool{},
			wantAmount:  0,
		},
		{
			name:        "with 1 white spend 1 white",
			pool:        Pool{white: 1},
			spendAmount: 1,
			spendColor:  White,
			wantPool:    Pool{},
			wantAmount:  0,
		},
		{
			name:        "with has 1 blue spend 1 blue",
			pool:        Pool{blue: 1},
			spendAmount: 1,
			spendColor:  Blue,
			wantPool:    Pool{},
			wantAmount:  0,
		},
		{
			name:        "with 1 black spend 1 black",
			pool:        Pool{black: 1},
			spendAmount: 1,
			spendColor:  Black,
			wantPool:    Pool{},
			wantAmount:  0,
		},
		{
			name:        "with 1 red spend 1 red",
			pool:        Pool{red: 1},
			spendAmount: 1,
			spendColor:  Red,
			wantPool:    Pool{},
			wantAmount:  0,
		},
		{
			name:        "with 1 green spend 1 green",
			pool:        Pool{green: 1},
			spendAmount: 1,
			spendColor:  Green,
			wantPool:    Pool{},
			wantAmount:  0,
		},
		{
			name:        "with 0 blue spend 1 blue",
			pool:        Pool{blue: 0},
			spendAmount: 1,
			spendColor:  Blue,
			wantPool:    Pool{blue: 0},
			wantAmount:  1,
		},
		{
			name:        "with 1 blue spend 2 blue",
			pool:        Pool{blue: 1},
			spendAmount: 2,
			spendColor:  Blue,
			wantPool:    Pool{blue: 0},
			wantAmount:  1,
		},
		{
			name: "with 2 blue spend 1 blue",
			pool: Pool{
				blue: 2,
			},
			spendAmount: 1,
			spendColor:  Blue,
			wantPool:    Pool{blue: 1},
			wantAmount:  0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotPool, gotAmount := test.pool.WithSpendMana(test.spendAmount, test.spendColor)
			if gotAmount != test.wantAmount {
				t.Errorf("SpendGenericManaByColor(%d, %s); gotAmount = %v, want %v", test.spendAmount, test.spendColor, gotAmount, test.wantAmount)
			}
			if diff := cmp.Diff(test.wantPool, gotPool, cmp.AllowUnexported(Pool{})); diff != "" {
				t.Errorf("SpendGenericManaByColor() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPoolWithSpendAmount(t *testing.T) {
	tests := []struct {
		name             string
		pool             Pool
		spendAmount      Amount
		colorsForGeneric []Color
		wantPool         Pool
		wantAmount       Amount
	}{
		{
			name:        "with 1 white spend 1 white",
			pool:        Pool{white: 1},
			spendAmount: Amount{white: 1},
			wantPool:    Pool{},
			wantAmount:  Amount{},
		},
		{
			name:        "with 1 blue spend 1 blue",
			pool:        Pool{blue: 1},
			spendAmount: Amount{blue: 1},
			wantPool:    Pool{},
			wantAmount:  Amount{},
		},
		{
			name:        "with 1 black spend 1 black",
			pool:        Pool{black: 1},
			spendAmount: Amount{black: 1},
			wantPool:    Pool{},
			wantAmount:  Amount{},
		},
		{
			name:        "with 1 red spend 1 red",
			pool:        Pool{red: 1},
			spendAmount: Amount{red: 1},
			wantPool:    Pool{},
			wantAmount:  Amount{},
		},
		{
			name:        "with 1 green spend 1 green",
			pool:        Pool{green: 1},
			spendAmount: Amount{green: 1},
			wantPool:    Pool{},
			wantAmount:  Amount{},
		},
		{
			name:        "with 1 colorless spend 1 colorless",
			pool:        Pool{colorless: 1},
			spendAmount: Amount{colorless: 1},
			wantPool:    Pool{},
			wantAmount:  Amount{},
		},
		{
			name: "with 1 colorless and wubrg spend 6 generic",
			pool: Pool{
				colorless: 1,
				white:     1,
				blue:      1,
				black:     1,
				red:       1,
				green:     1,
			},
			spendAmount: Amount{
				generic: 6,
			},
			colorsForGeneric: Colors(),
			wantPool:         Pool{},
			wantAmount:       Amount{},
		},
		{
			name: "with empty pool spend 1 generic 1 colorless and wubrg",
			pool: Pool{},
			spendAmount: Amount{
				generic:   1,
				colorless: 1,
				white:     1,
				blue:      1,
				black:     1,
				red:       1,
				green:     1,
			},
			colorsForGeneric: []Color{White, Blue, Black, Red, Green},
			wantPool:         Pool{},
			wantAmount: Amount{
				generic:   1,
				colorless: 1,
				white:     1,
				blue:      1,
				black:     1,
				red:       1,
				green:     1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotPool, gotAmount := test.pool.WithSpendAmount(test.spendAmount, test.colorsForGeneric)
			if diff := cmp.Diff(test.wantAmount, gotAmount, cmp.AllowUnexported(Amount{})); diff != "" {
				t.Errorf("WithSpendManaAmount(); Amount mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(test.wantPool, gotPool, cmp.AllowUnexported(Pool{})); diff != "" {
				t.Errorf("WithSpendManaAmount(); Pool mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAmountOf(t *testing.T) {
	tests := []struct {
		name  string
		pool  Pool
		color Color
		want  int
	}{
		{
			name:  "with 1 colorless",
			pool:  Pool{colorless: 1},
			color: Colorless,
			want:  1,
		},
		{
			name:  "with 1 white",
			pool:  Pool{white: 1},
			color: White,
			want:  1,
		},
		{
			name:  "with 1 blue",
			pool:  Pool{blue: 1},
			color: Blue,
			want:  1,
		},
		{
			name:  "with 1 black",
			pool:  Pool{black: 1},
			color: Black,
			want:  1,
		},
		{
			name:  "with 1 red",
			pool:  Pool{red: 1},
			color: Red,
			want:  1,
		},
		{
			name:  "with 1 green",
			pool:  Pool{green: 1},
			color: Green,
			want:  1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.pool.AmountOf(test.color)
			if got != test.want {
				t.Errorf("AmountOf(%s) = %d; want %d", test.color, got, test.want)
			}
		})
	}
}

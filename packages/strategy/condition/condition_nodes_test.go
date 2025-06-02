package condition_test

/*
func TestLandDropCondition(t *testing.T) {
	tests := []struct {
		name   string
		player *game.Player
		node   strategy.LandDropCondition
		want   bool
	}{
		{
			name:   "Land drop is avaible, and it's wanted",
			player: &game.Player{LandDrop: false},
			node:   strategy.LandDropCondition{LandDrop: false},
			want:   true,
		},
		{
			name:   "Land drop is not available, and it's not wanted",
			player: &game.Player{LandDrop: true},
			node:   strategy.LandDropCondition{LandDrop: true},
			want:   true,
		},
		{
			name:   "Land drop is avaible, and it's not wanted",
			player: &game.Player{LandDrop: false},
			node:   strategy.LandDropCondition{LandDrop: true},
			want:   false,
		},
		{
			name:   "Land drop is not available, and it's wanted",
			player: &game.Player{LandDrop: true},
			node:   strategy.LandDropCondition{LandDrop: false},
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := &strategy.EvaluatorContext{
				Player: test.player,
			}
			got, err := test.node.Evaluate(ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != test.want {
				t.Errorf("got %v, want %v", test.want, got)
			}
		})
	}
}
*/

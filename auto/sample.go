package auto

var SampleDSLStrategy = []Rule{
	/*
		{
			Name:        "Land Drop",
			Description: "Play a land if you can and have one in hand",
			When: When(
				InHand("Island"),
			),
			Then: Cast("Land"),
		},
		/*
			{
				Name:        "Early Setup",
				Description: "Play Lorien Revealed when setup pieces are in hand and opponent is tapped out.",
				When: When(
					InHand("Lorien Revealed"),
					Or(
						InHand("Combo Piece A"),
						InHand("Combo Piece B"),
					),
					Not(OnBattlefield("Disruption Card")),
					HandSize(">=", 5),
				),
				Then: Cast("Lorien Revealed"),
			},
			{
				Name:        "Emergency Defense",
				Description: "Switch to panic mode if life total is too low.",
				When: When(
					LifeTotal("<", 6),
				),
				Then: Cast("Weather the Storm"),
			},
	*/
}

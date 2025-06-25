package judge

import (
	"testing"
)

func TestGetAvailableMana(t *testing.T) {
	// I can't test this because I need to have the card
	// actualy have the ability, and the config file doesn't give me that.
	// I need to load cards from the defintions.
	/*
		playerID := "Test Player"
		game := state.LoadGameFromConfig(statetest.GameConfig{
			Players: []statetest.PlayerConfig{{ID: playerID}},
			Battlefield: statetest.BattlefieldConfig{
				Permanents: []gobtest.PermanentConfig{
					{ID: "Forest ID", Name: "Forest", CardTypes: []mtg.CardType{mtg.CardTypeLand}, Controller: playerID, Tapped: false},
					{ID: "Island ID", Name: "Island", CardTypes: []mtg.CardType{mtg.CardTypeLand}, Controller: playerID, Tapped: true},
				},
			},
		})
		player := game.GetPlayer(playerID)
		got := GetAvailableMana(game, player).ManaString()
		want := "{G}"
		if got != want {
			t.Errorf("GetAvailableMana() = %v; want %v", got, want)
		}
	*/
}

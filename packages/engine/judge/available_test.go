package judge

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"testing"
)

func TestGetAvailableMana(t *testing.T) {
	// I can't test this because I need to have the card
	// actualy have the ability, and the config file doesn't give me that.
	// I need to load cards from the defintions.
	playerID := "Test Player"
	game := state.NewGameFromDefinition(definition.Game{
		Players: []definition.Player{{ID: playerID}},
		Battlefield: definition.Battlefield{
			Permanents: []definition.Permanent{
				{
					ID:         "Forest ID",
					Name:       "Forest",
					CardTypes:  []string{string(mtg.CardTypeLand)},
					Controller: playerID,
					Tapped:     false,
					ActivatedAbilities: []definition.Ability{{
						Effects: []definition.Effect{{
							Name: "AddMana",
							Modifiers: map[string]any{
								"Mana": "{G}",
							},
						}},
						Zone: string(mtg.ZoneBattlefield),
					}},
				},
				{
					ID:         "Island ID",
					Name:       "Island",
					CardTypes:  []string{string(mtg.CardTypeLand)},
					Controller: playerID,
					Tapped:     true,
					ActivatedAbilities: []definition.Ability{{
						Effects: []definition.Effect{{
							Name: "AddMana",
							Modifiers: map[string]any{
								"Mana": "{U}",
							},
						}},
						Zone: string(mtg.ZoneBattlefield),
					}},
				},
			},
		},
	})
	player := game.GetPlayer(playerID)
	got := GetAvailableMana(game, player).ManaString()
	want := "{G}"
	if got != want {
		t.Errorf("GetAvailableMana() = %v; want %v", got, want)
	}
}

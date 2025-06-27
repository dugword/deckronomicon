package actionparser

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
)

func newTestGame(playerID string) state.Game {
	game := state.NewGameFromDefinition(definition.Game{
		Step: string(mtg.StepPrecombatMain),
		Players: []definition.Player{
			{
				ID: playerID,
				Hand: definition.Hand{Cards: []definition.Card{
					{
						ID:   "Test Card ID",
						Name: "Test Card",
					},
					{
						ID:       "Acane Card ID",
						Name:     "Arcane Card",
						Subtypes: []string{string(mtg.SubtypeArcane)},
					},
					{
						ID:   "Card with Splice ID",
						Name: "Card with Splice",
						StaticAbilities: []definition.StaticAbility{{
							Name:      string(mtg.StaticKeywordSplice),
							Modifiers: map[string]any{"Subtype": "Arcane"},
						}},
					},
					{
						ID:   "Card with Replicate ID",
						Name: "Card with Replicate",
						StaticAbilities: []definition.StaticAbility{{
							Name: string(mtg.StaticKeywordReplicate),
						}},
					},
					{
						ID:   "Card with Target ID",
						Name: "Card with Target",
						SpellAbility: []definition.Effect{{
							Name:      "Target",
							Modifiers: map[string]any{"Target": "Permanent"},
						}},
					},
					{
						ID:         "Card with Ability ID",
						Name:       "Card with Ability",
						Controller: playerID,
						ActivatedAbilities: []definition.Ability{{
							Name: "Ability on Card",
							Zone: string(mtg.ZoneHand),
						}},
					},
					{
						ID:   "Card with Effects ID",
						Name: "Card with Effects",
						SpellAbility: []definition.Effect{
							{
								Name:      "Target",
								Modifiers: map[string]any{"Target": "Permanent"},
							},
							{
								Name:      "Target",
								Modifiers: map[string]any{"Target": "Permanent"},
							},
						},
					},
				}},
				Graveyard: definition.Graveyard{
					Cards: []definition.Card{
						{
							ID:   "Card with Flashback ID",
							Name: "Card with Flashback",
							StaticAbilities: []definition.StaticAbility{{
								Name: string(mtg.StaticKeywordFlashback),
							}},
						},
					},
				},
			},
		},
		Battlefield: definition.Battlefield{
			Permanents: []definition.Permanent{
				{
					ID:         "Test Permanent ID",
					Name:       "Test Permanent",
					Controller: playerID,
					Owner:      playerID,
					ActivatedAbilities: []definition.Ability{{
						Name: "Ability on Permanent",
					}},
				},
			},
		},
	})
	return game
}

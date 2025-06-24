package action

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob/gobtest"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"deckronomicon/packages/state/statetest"
)

func newTestGame(playerID string) state.Game {
	game := state.LoadGameFromConfig(statetest.GameConfig{
		Step: mtg.StepPrecombatMain,
		Players: []statetest.PlayerConfig{
			{
				ID: playerID,
				Hand: statetest.HandConfig{Cards: []gobtest.CardConfig{
					{
						ID:   "Test Card ID",
						Name: "Test Card",
					},
					{
						ID:       "Acane Card ID",
						Name:     "Arcane Card",
						Subtypes: []mtg.Subtype{mtg.SubtypeArcane},
					},
					{
						ID:   "Card with Splice ID",
						Name: "Card with Splice",
						StaticAbilities: []definition.StaticAbilitySpec{{
							Name:      mtg.StaticKeywordSplice,
							Modifiers: map[string]any{"Subtype": "Arcane"},
						}},
					},
					{
						ID:   "Card with Replicate ID",
						Name: "Card with Replicate",
						StaticAbilities: []definition.StaticAbilitySpec{{
							Name: mtg.StaticKeywordReplicate,
						}},
					},
					{
						ID:   "Card with Target ID",
						Name: "Card with Target",
						SpellAbility: []definition.EffectSpec{{
							Name:      "Target",
							Modifiers: map[string]any{"Target": "Player"},
						}},
					},
					{
						ID:   "Card with Ability ID",
						Name: "Card with Ability",
						ActivatedAbilities: []definition.ActivatedAbilitySpec{{
							Name: "Ability on Card",
							Zone: mtg.ZoneHand,
						}},
					},
					{
						ID:   "Card with Effects ID",
						Name: "Card with Effects",
						SpellAbility: []definition.EffectSpec{
							{
								Name:      "Effect 1",
								Modifiers: map[string]any{"Target": "Permanent"},
							},
							{
								Name:      "Effect 2",
								Modifiers: map[string]any{"Target": "Player"},
							},
						},
					},
				}},
				Graveyard: statetest.GraveyardConfig{
					Cards: []gobtest.CardConfig{
						{
							ID:   "Card with Flashback ID",
							Name: "Card with Flashback",
							StaticAbilities: []definition.StaticAbilitySpec{{
								Name: mtg.StaticKeywordFlashback,
							}},
						},
					},
				},
			},
		},
		Battlefield: statetest.BattlefieldConfig{
			Permanents: []gobtest.PermanentConfig{
				{
					ID:         "Test Permanent ID",
					Name:       "Test Permanent",
					Controller: playerID,
					Owner:      playerID,
					ActivatedAbilities: []definition.ActivatedAbilitySpec{{
						Name: "Ability on Permanent",
					}},
				},
			},
		},
	})
	return game
}

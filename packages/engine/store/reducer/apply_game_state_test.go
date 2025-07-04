package reducer

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/state"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestApplyGameStateChangeEvent(t *testing.T) {
	const playerID = "Test Player"
	testCases := []struct {
		name string
		evnt event.GameStateChangeEvent
		game *state.Game
		want *state.Game
	}{
		{
			name: "with AddManaEvent",
			evnt: &event.AddManaEvent{
				PlayerID: playerID,
				Color:    mana.Blue,
				Amount:   2,
			},
			game: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{ID: playerID}},
			}),
			want: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID:       playerID,
					ManaPool: "{U}{U}",
				}},
			}),
		},
		{
			name: "with CheatEnabledEvent",
			evnt: &event.CheatEnabledEvent{
				PlayerID: playerID,
			},
			game: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{ID: playerID}},
			}),
			want: state.NewGameFromDefinition(&definition.Game{
				CheatsEnabled: true,
				Players: []*definition.Player{{
					ID: playerID,
				}},
			}),
		},
		{
			name: "with DiscardCardEvent",
			evnt: &event.DiscardCardEvent{
				PlayerID: playerID,
				CardID:   "Discarded Card ID",
			},
			game: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
					Hand: &definition.Hand{
						Cards: []*definition.Card{{ID: "Discarded Card ID"}},
					},
				}},
			}),
			want: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
					Graveyard: &definition.Graveyard{
						Cards: []*definition.Card{{ID: "Discarded Card ID"}},
					},
				}},
			}),
		},
		{
			name: "with DrawCardEvent",
			evnt: &event.DrawCardEvent{
				PlayerID: playerID,
			},
			game: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
					Library: &definition.Library{
						Cards: []*definition.Card{{ID: "Drawn Card ID"}},
					},
				}},
			}),
			want: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
					Hand: &definition.Hand{
						Cards: []*definition.Card{{ID: "Drawn Card ID"}},
					},
				}},
			}),
		},
		{
			name: "with GainLifeEvent",
			evnt: &event.GainLifeEvent{
				PlayerID: playerID,
				Amount:   5,
			},
			game: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{ID: playerID, Life: 10}},
			}),
			want: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{ID: playerID, Life: 15}},
			}),
		},
		{
			name: "with LoseLifeEvent",
			evnt: &event.LoseLifeEvent{
				PlayerID: playerID,
				Amount:   3,
			},
			game: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{ID: playerID, Life: 10}},
			}),
			want: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{ID: playerID, Life: 7}},
			}),
		},
		{
			name: "with PutCardInHandEvent",
			evnt: &event.PutCardInHandEvent{
				PlayerID: playerID,
				CardID:   "Card ID to Hand",
				FromZone: "Library",
			},
			game: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
					Library: &definition.Library{
						Cards: []*definition.Card{{ID: "Card ID to Hand"}},
					},
				}},
			}),
			want: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
					Hand: &definition.Hand{
						Cards: []*definition.Card{{ID: "Card ID to Hand"}},
					},
				}},
			}),
		},
		{
			name: "with PutCardInGraveyardEvent",
			evnt: &event.PutCardInGraveyardEvent{
				PlayerID: playerID,
				CardID:   "Card ID to Graveyard",
				FromZone: "Hand",
			},
			game: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
					Hand: &definition.Hand{
						Cards: []*definition.Card{{ID: "Card ID to Graveyard"}},
					},
				}},
			}),
			want: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
					Graveyard: &definition.Graveyard{
						Cards: []*definition.Card{{ID: "Card ID to Graveyard"}},
					},
				}},
			}),
		},
		{
			name: "with PutCardOnBottomOfLibraryEvent",
			evnt: &event.PutCardOnBottomOfLibraryEvent{
				PlayerID: playerID,
				CardID:   "Card ID to Bottom of Library",
				FromZone: "Hand",
			},
			game: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
					Hand: &definition.Hand{
						Cards: []*definition.Card{{ID: "Card ID to Bottom of Library"}},
					},
				}},
			}),
			want: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
					Library: &definition.Library{
						Cards: []*definition.Card{{ID: "Card ID to Bottom of Library"}},
					},
				}},
			}),
		},
		{
			name: "with PutCardOnTopOfLibraryEvent",
			evnt: &event.PutCardOnTopOfLibraryEvent{
				PlayerID: playerID,
				CardID:   "Card ID to Top of Library",
				FromZone: "Hand",
			},
			game: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
					Hand: &definition.Hand{
						Cards: []*definition.Card{{ID: "Card ID to Top of Library"}},
					},
				}},
			}),
			want: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
					Library: &definition.Library{
						Cards: []*definition.Card{{ID: "Card ID to Top of Library"}},
					},
				}},
			}),
		},
		{
			name: "with PutPermanentOnBattlefieldEvent",
			evnt: &event.PutPermanentOnBattlefieldEvent{
				PlayerID: playerID,
				CardID:   "Card ID to Battlefield",
				FromZone: "Hand",
			},
			game: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
					Hand: &definition.Hand{
						Cards: []*definition.Card{{ID: "Card ID to Battlefield"}},
					},
				}},
			}),
			want: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
				}},
				Battlefield: &definition.Battlefield{
					Permanents: []*definition.Permanent{
						{
							ID:         "1",
							Controller: playerID,
							Owner:      playerID,
							Card:       &definition.Card{ID: "Card ID to Battlefield"},
						},
					},
				},
			}),
		},
		{
			name: "with SetActivePlayerEvent",
			evnt: &event.SetActivePlayerEvent{
				PlayerID: playerID,
			},
			game: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{ID: playerID}},
			}),
			want: state.NewGameFromDefinition(&definition.Game{
				ActivePlayerID: playerID,
				Players:        []*definition.Player{{ID: playerID}},
			}),
		},
		{
			name: "with ShuffleLibraryEvent",
			evnt: &event.ShuffleLibraryEvent{
				PlayerID: playerID,
				ShuffledCardsIDs: []string{
					"Card ID 4",
					"Card ID 3",
					"Card ID 5",
					"Card ID 1",
					"Card ID 2",
				},
			},
			game: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
					Library: &definition.Library{
						Cards: []*definition.Card{
							{ID: "Card ID 1"},
							{ID: "Card ID 2"},
							{ID: "Card ID 3"},
							{ID: "Card ID 4"},
							{ID: "Card ID 5"},
						},
					},
				}},
			}),
			want: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
					Library: &definition.Library{
						Cards: []*definition.Card{
							{ID: "Card ID 4"},
							{ID: "Card ID 3"},
							{ID: "Card ID 5"},
							{ID: "Card ID 1"},
							{ID: "Card ID 2"},
						},
					},
				}},
			}),
		},
		{
			name: "with SpendManaEvent",
			evnt: &event.SpendManaEvent{
				PlayerID:   playerID,
				ManaString: "{U}{U}",
			},

			game: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID:       playerID,
					ManaPool: "{U}{U}{U}",
				}},
			}),
			want: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID:       playerID,
					ManaPool: "{U}",
				}},
			}),
		},
		{
			name: "with TapPermanentEvent",
			evnt: &event.TapPermanentEvent{
				PlayerID:    playerID,
				PermanentID: "Permanent ID to Tap",
			},
			game: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
				}},
				Battlefield: &definition.Battlefield{
					Permanents: []*definition.Permanent{
						{ID: "Permanent ID to Tap"},
					},
				},
			}),
			want: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
				}},
				Battlefield: &definition.Battlefield{
					Permanents: []*definition.Permanent{
						{ID: "Permanent ID to Tap", Tapped: true},
					},
				},
			}),
		},
		{
			name: "with UntapPermanentEvent",
			evnt: &event.UntapPermanentEvent{
				PlayerID:    playerID,
				PermanentID: "Permanent ID to Untap",
			},
			game: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
				}},
				Battlefield: &definition.Battlefield{
					Permanents: []*definition.Permanent{
						{ID: "Permanent ID to Untap", Tapped: true},
					},
				},
			}),
			want: state.NewGameFromDefinition(&definition.Game{
				Players: []*definition.Player{{
					ID: playerID,
				}},
				Battlefield: &definition.Battlefield{
					Permanents: []*definition.Permanent{
						{ID: "Permanent ID to Untap", Tapped: false},
					},
				},
			}),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := applyGameStateChangeEvent(tc.game, tc.evnt)
			if err != nil {
				t.Fatalf("applyGameStateChangeEvent(game, %T); err = %v; want %v", tc.evnt, err, nil)
			}
			if diff := cmp.Diff(tc.want, got,
				AllowAllUnexported,
				cmpopts.EquateEmpty(),
				// nextID is incremented internally to keep track of object IDs.
				// It should not affect the equality of game states.
				// IDs created by nextID are reflected in the state of game objects.
				// We ignore it here to focus on the game state changes.
				cmpopts.IgnoreFields(state.Game{}, "nextID"),
			); diff != "" {
				t.Errorf("applyGameStateChangeEvent(game, %T) mismatch (-want +got):\n%s", tc.evnt, diff)
			}
		})
	}
}

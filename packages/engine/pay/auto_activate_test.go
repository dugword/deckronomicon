package pay

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/definition/definitiontest"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TODO: Maybe have each package provide a AllowAllUnexported variable
// Or each package provide it's own cmp.Diff function.
// That would let us control if we ever want hide fields from being compared
// in tests, without having to change every test that uses cmp.Diff.
var AllowAllUnexported = cmp.Exporter(func(reflect.Type) bool { return true })

func newTestGame(
	playerID string,
	manaPool string,
	battlefield []definition.Permanent,
) state.Game {
	game := state.NewGameFromDefinition(definition.Game{
		Players: []definition.Player{{ID: playerID, ManaPool: manaPool}},
		Battlefield: definition.Battlefield{
			Permanents: battlefield,
		},
	})
	return game
}

func TestWithActivateManaSources(t *testing.T) {
	const (
		playerID   = "Test Player"
		plainsID   = "Test Plains"
		islandID   = "Test Island"
		swampID    = "Test Swamp"
		mountainID = "Test Mountain"
		forestID   = "Test Forest"
		wastesID   = "Test Wastes"
	)
	tests := []struct {
		name       string
		cost       string
		manaPool   string
		permanents []definition.Permanent
		colors     []mana.Color
		want       []event.GameEvent
	}{
		{
			name: "with 2 generic 1 white mana",
			cost: "{2}{W}",
			permanents: []definition.Permanent{
				definitiontest.PlainsDefinition(plainsID, playerID),
				definitiontest.IslandDefinition(islandID, playerID),
				definitiontest.SwampDefinition(swampID, playerID),
				definitiontest.MountainDefinition(mountainID, playerID),
				definitiontest.ForestDefinition(forestID, playerID),
				definitiontest.WastesDefinition(wastesID, playerID),
			},
			colors: mana.Colors(),
			want: []event.GameEvent{
				// Plains
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: plainsID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: plainsID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: plainsID, Subtypes: []mtg.Subtype{mtg.SubtypePlains}},
				event.AddManaEvent{Amount: 1, Color: mana.White, PlayerID: playerID},
				// Wastes
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: wastesID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: wastesID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: wastesID},
				event.AddManaEvent{Amount: 1, Color: mana.Colorless, PlayerID: playerID},
				// Island
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: islandID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: islandID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: islandID, Subtypes: []mtg.Subtype{mtg.SubtypeIsland}},
				event.AddManaEvent{Amount: 1, Color: mana.Blue, PlayerID: playerID},
			},
		},
		{
			name: "with 2 generic 1 white mana prioritize black and red for generic",
			cost: "{2}{W}",
			permanents: []definition.Permanent{
				definitiontest.PlainsDefinition(plainsID, playerID),
				definitiontest.IslandDefinition(islandID, playerID),
				definitiontest.SwampDefinition(swampID, playerID),
				definitiontest.MountainDefinition(mountainID, playerID),
				definitiontest.ForestDefinition(forestID, playerID),
				definitiontest.WastesDefinition(wastesID, playerID),
			},
			colors: []mana.Color{mana.Black, mana.Red},
			want: []event.GameEvent{
				// Plains
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: plainsID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: plainsID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: plainsID, Subtypes: []mtg.Subtype{mtg.SubtypePlains}},
				event.AddManaEvent{Amount: 1, Color: mana.White, PlayerID: playerID},
				// Swamp
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: swampID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: swampID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: swampID, Subtypes: []mtg.Subtype{mtg.SubtypeSwamp}},
				event.AddManaEvent{Amount: 1, Color: mana.Black, PlayerID: playerID},
				// Mountain
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: mountainID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: mountainID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: mountainID, Subtypes: []mtg.Subtype{mtg.SubtypeMountain}},
				event.AddManaEvent{Amount: 1, Color: mana.Red, PlayerID: playerID},
			},
		},
		{
			name:     "with 2 generic 1 white mana with 1 generic 1 white in pool",
			cost:     "{2}{W}",
			manaPool: "{1}{W}",
			permanents: []definition.Permanent{
				definitiontest.PlainsDefinition(plainsID, playerID),
				definitiontest.IslandDefinition(islandID, playerID),
				definitiontest.SwampDefinition(swampID, playerID),
				definitiontest.MountainDefinition(mountainID, playerID),
				definitiontest.ForestDefinition(forestID, playerID),
				definitiontest.WastesDefinition(wastesID, playerID),
			},
			colors: mana.Colors(),
			want: []event.GameEvent{
				// Wastes
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: wastesID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: wastesID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: wastesID},
				event.AddManaEvent{Amount: 1, Color: mana.Colorless, PlayerID: playerID},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			game := newTestGame(
				playerID,
				test.manaPool,
				test.permanents,
			)
			testCost, err := cost.ParseManaCost(test.cost)
			if err != nil {
				t.Fatalf("ParseManaCost(); error = %v", err)
			}
			_, got, err := withActivateManaSources(game, playerID, testCost, test.colors)
			if err != nil {
				t.Fatalf("withActivateManaSources(); error = %v", err)
			}
			if diff := cmp.Diff(test.want, got, AllowAllUnexported); diff != "" {
				t.Errorf("withActivateManaSources(); events mismatch: %s", diff)
			}
		})
	}
}

func TestActivateManaSourcesForColored(t *testing.T) {
	const (
		playerID   = "Test Player"
		plainsID   = "Test Plains"
		islandID   = "Test Island"
		swampID    = "Test Swamp"
		mountainID = "Test Mountain"
		forestID   = "Test Forest"
		wastesID   = "Test Wastes"
	)
	tests := []struct {
		name          string
		amount        string
		manaColor     mana.Color
		permanents    []definition.Permanent
		want          []event.GameEvent
		wantRemaining string
	}{
		{
			name:   "with white mana",
			amount: "{W}",
			permanents: []definition.Permanent{
				definitiontest.PlainsDefinition(plainsID, playerID),
			},
			manaColor: mana.White,
			want: []event.GameEvent{
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: plainsID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: plainsID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: plainsID, Subtypes: []mtg.Subtype{mtg.SubtypePlains}},
				event.AddManaEvent{Amount: 1, Color: mana.White, PlayerID: playerID},
			},
		},
		{
			name:   "with blue mana",
			amount: "{U}",
			permanents: []definition.Permanent{
				definitiontest.IslandDefinition(islandID, playerID),
			},
			manaColor: mana.Blue,
			want: []event.GameEvent{
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: islandID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: islandID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: islandID, Subtypes: []mtg.Subtype{mtg.SubtypeIsland}},
				event.AddManaEvent{Amount: 1, Color: mana.Blue, PlayerID: playerID},
			},
		},
		{
			name:   "with black mana",
			amount: "{B}",
			permanents: []definition.Permanent{
				definitiontest.SwampDefinition(swampID, playerID),
			},
			manaColor: mana.Black,
			want: []event.GameEvent{
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: swampID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: swampID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: swampID, Subtypes: []mtg.Subtype{mtg.SubtypeSwamp}},
				event.AddManaEvent{Amount: 1, Color: mana.Black, PlayerID: playerID},
			},
		},
		{
			name:   "with red mana",
			amount: "{R}",
			permanents: []definition.Permanent{
				definitiontest.MountainDefinition(mountainID, playerID),
			},
			manaColor: mana.Red,
			want: []event.GameEvent{
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: mountainID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: mountainID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: mountainID, Subtypes: []mtg.Subtype{mtg.SubtypeMountain}},
				event.AddManaEvent{Amount: 1, Color: mana.Red, PlayerID: playerID},
			},
		},
		{
			name:   "with green mana",
			amount: "{G}",
			permanents: []definition.Permanent{
				definitiontest.ForestDefinition(forestID, playerID),
			},
			manaColor: mana.Green,
			want: []event.GameEvent{
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: forestID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: forestID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: forestID, Subtypes: []mtg.Subtype{mtg.SubtypeForest}},
				event.AddManaEvent{Amount: 1, Color: mana.Green, PlayerID: playerID},
			},
		},
		{
			name:   "with colorless mana",
			amount: "{C}",
			permanents: []definition.Permanent{
				definitiontest.WastesDefinition(wastesID, playerID),
			},
			manaColor: mana.Colorless,
			want: []event.GameEvent{
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: wastesID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: wastesID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: wastesID},
				event.AddManaEvent{Amount: 1, Color: mana.Colorless, PlayerID: playerID},
			},
		},
		{
			name:   "with need white have none",
			amount: "{W}",
			permanents: []definition.Permanent{
				definitiontest.IslandDefinition(islandID, playerID),
				definitiontest.SwampDefinition(swampID, playerID),
				definitiontest.MountainDefinition(mountainID, playerID),
				definitiontest.ForestDefinition(forestID, playerID),
				definitiontest.WastesDefinition(wastesID, playerID),
			},
			manaColor:     mana.White,
			want:          nil,
			wantRemaining: "{W}",
		},
		{
			name:   "with need 1 white mana have 2",
			amount: "{W}",
			permanents: []definition.Permanent{
				definitiontest.PlainsDefinition(plainsID, playerID),
				definitiontest.PlainsDefinition("Test Plains 2", playerID),
			},
			manaColor: mana.White,
			want: []event.GameEvent{
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: plainsID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: plainsID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: plainsID, Subtypes: []mtg.Subtype{mtg.SubtypePlains}},
				event.AddManaEvent{Amount: 1, Color: mana.White, PlayerID: playerID},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			game := newTestGame(
				playerID,
				"",
				test.permanents,
			)
			remaining, err := mana.ParseManaString(test.amount)
			if err != nil {
				t.Fatalf("ParseManaString(); error = %v", err)
			}
			_, gotRemaining, got, err := activateManaSourcesForColored(game, playerID, remaining, test.manaColor)
			if err != nil {
				t.Fatalf("activateManaSourcesForColored(); error = %v", err)
			}

			if test.wantRemaining == "" {
				if gotRemaining.Total() != 0 {
					t.Errorf("activateManaSourcesForColored(); remaining = %s, want 0", gotRemaining.ManaString())
				}
			} else {
				wantRemaining, err := mana.ParseManaString(test.wantRemaining)
				if err != nil {
					t.Fatalf("ParseManaString(); error = %v", err)
				}
				if gotRemaining != wantRemaining {
					t.Errorf("activateManaSourcesForColored(); remaining = %v, want %v", gotRemaining.ManaString(), wantRemaining.ManaString())
				}
			}
			if diff := cmp.Diff(test.want, got, AllowAllUnexported); diff != "" {
				t.Errorf("activateManaSourcesForColored(); events mismatch: %s", diff)
			}
		})
	}
}

func TestActivateManaSourcesForGeneric(t *testing.T) {
	const (
		playerID   = "Test Player"
		plainsID   = "Test Plains"
		islandID   = "Test Island"
		swampID    = "Test Swamp"
		mountainID = "Test Mountain"
		forestID   = "Test Forest"
		wastesID   = "Test Wastes"
	)
	tests := []struct {
		name          string
		amount        string
		permanents    []definition.Permanent
		want          []event.GameEvent
		wantRemaining string
	}{
		{
			name:   "with need 1 generic mana",
			amount: "{1}",
			permanents: []definition.Permanent{
				definitiontest.WastesDefinition(wastesID, playerID),
			},
			want: []event.GameEvent{
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: wastesID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: wastesID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: wastesID},
				event.AddManaEvent{Amount: 1, Color: mana.Colorless, PlayerID: playerID},
			},
		},
		{
			name:          "with need 1 generic mana 0 available",
			amount:        "{1}",
			permanents:    []definition.Permanent{},
			want:          nil,
			wantRemaining: "{1}",
		},
		{
			name:   "'with need 1 generic mana 2 available",
			amount: "{1}",
			permanents: []definition.Permanent{
				definitiontest.WastesDefinition(wastesID, playerID),
				definitiontest.PlainsDefinition(plainsID, playerID),
			},
			want: []event.GameEvent{
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: wastesID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: wastesID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: wastesID},
				event.AddManaEvent{Amount: 1, Color: mana.Colorless, PlayerID: playerID},
			},
		},
		{
			name:   "with need 2 generic mana 1 available",
			amount: "{2}",
			permanents: []definition.Permanent{
				definitiontest.WastesDefinition(wastesID, playerID),
			},
			want: []event.GameEvent{
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: wastesID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: wastesID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: wastesID},
				event.AddManaEvent{Amount: 1, Color: mana.Colorless, PlayerID: playerID},
			},
			wantRemaining: "{1}",
		},
		{
			name:   "with need 6 generic mana",
			amount: "{6}",
			permanents: []definition.Permanent{
				definitiontest.PlainsDefinition(plainsID, playerID),
				definitiontest.IslandDefinition(islandID, playerID),
				definitiontest.SwampDefinition(swampID, playerID),
				definitiontest.MountainDefinition(mountainID, playerID),
				definitiontest.ForestDefinition(forestID, playerID),
				definitiontest.WastesDefinition(wastesID, playerID),
			},
			want: []event.GameEvent{
				// Wastes
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: wastesID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: wastesID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: wastesID},
				event.AddManaEvent{Amount: 1, Color: mana.Colorless, PlayerID: playerID},
				// Plains
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: plainsID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: plainsID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: plainsID, Subtypes: []mtg.Subtype{mtg.SubtypePlains}},
				event.AddManaEvent{Amount: 1, Color: mana.White, PlayerID: playerID},
				// Island
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: islandID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: islandID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: islandID, Subtypes: []mtg.Subtype{mtg.SubtypeIsland}},
				event.AddManaEvent{Amount: 1, Color: mana.Blue, PlayerID: playerID},
				// Swamp
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: swampID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: swampID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: swampID, Subtypes: []mtg.Subtype{mtg.SubtypeSwamp}},
				event.AddManaEvent{Amount: 1, Color: mana.Black, PlayerID: playerID},
				// Mountain
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: mountainID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: mountainID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: mountainID, Subtypes: []mtg.Subtype{mtg.SubtypeMountain}},
				event.AddManaEvent{Amount: 1, Color: mana.Red, PlayerID: playerID},
				// Forest
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: forestID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: forestID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: forestID, Subtypes: []mtg.Subtype{mtg.SubtypeForest}},
				event.AddManaEvent{Amount: 1, Color: mana.Green, PlayerID: playerID},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			game := newTestGame(
				playerID,
				"",
				test.permanents,
			)
			remaining, err := mana.ParseManaString(test.amount)
			if err != nil {
				t.Fatalf("ParseManaString(); error = %v", err)
			}
			_, gotRemaining, got, err := activateManaSourcesForGeneric(game, playerID, remaining, mana.Colors())
			if err != nil {
				t.Fatalf("activateManaSourcesForGeneric(); error = %v", err)
			}
			if test.wantRemaining == "" {
				if gotRemaining.Total() != 0 {
					t.Errorf("activateManaSourcesForGeneric(); remaining = %s, want 0", gotRemaining.ManaString())
				}
			} else {
				wantRemaining, err := mana.ParseManaString(test.wantRemaining)
				if err != nil {
					t.Fatalf("ParseManaString(); error = %v", err)
				}
				if gotRemaining != wantRemaining {
					t.Errorf("activateManaSourcesForGeneric(); remaining = %s, want %s", gotRemaining.ManaString(), wantRemaining.ManaString())
				}
			}
			if diff := cmp.Diff(test.want, got, AllowAllUnexported); diff != "" {
				t.Errorf("activateManaSourcesForGeneric(); events mismatch: %s", diff)
			}
		})
	}
}

func TestActivateManaSource(t *testing.T) {
	const (
		playerID   = "Test Player"
		plainsID   = "Test Plains"
		islandID   = "Test Island"
		swampID    = "Test Swamp"
		mountainID = "Test Mountain"
		forestID   = "Test Forest"
		wastesID   = "Test Wastes"
	)
	tests := []struct {
		name   string
		landID string
		want   []event.GameEvent
	}{
		{
			name:   "with activate white source",
			landID: plainsID,
			want: []event.GameEvent{
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: plainsID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: plainsID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: plainsID, Subtypes: []mtg.Subtype{mtg.SubtypePlains}},
				event.AddManaEvent{Amount: 1, Color: mana.White, PlayerID: playerID},
			},
		},
		{
			name:   "with activate blue source",
			landID: islandID,
			want: []event.GameEvent{
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: islandID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: islandID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: islandID, Subtypes: []mtg.Subtype{mtg.SubtypeIsland}},
				event.AddManaEvent{Amount: 1, Color: mana.Blue, PlayerID: playerID},
			},
		},
		{
			name:   "with activate black source",
			landID: swampID,
			want: []event.GameEvent{
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: swampID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: swampID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: swampID, Subtypes: []mtg.Subtype{mtg.SubtypeSwamp}},
				event.AddManaEvent{Amount: 1, Color: mana.Black, PlayerID: playerID},
			},
		},
		{
			name:   "with activate red source",
			landID: mountainID,
			want: []event.GameEvent{
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: mountainID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: mountainID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: mountainID, Subtypes: []mtg.Subtype{mtg.SubtypeMountain}},
				event.AddManaEvent{Amount: 1, Color: mana.Red, PlayerID: playerID},
			},
		},
		{
			name:   "with activate green source",
			landID: forestID,
			want: []event.GameEvent{
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: forestID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: forestID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: forestID, Subtypes: []mtg.Subtype{mtg.SubtypeForest}},
				event.AddManaEvent{Amount: 1, Color: mana.Green, PlayerID: playerID},
			},
		},
		{
			name:   "with activate colorless source",
			landID: wastesID,
			want: []event.GameEvent{
				event.ActivateAbilityEvent{PlayerID: playerID, SourceID: wastesID, Zone: mtg.ZoneBattlefield},
				event.TapPermanentEvent{PlayerID: playerID, PermanentID: wastesID},
				event.LandTappedForManaEvent{PlayerID: playerID, ObjectID: wastesID},
				event.AddManaEvent{Amount: 1, Color: mana.Colorless, PlayerID: playerID},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			permanents := []definition.Permanent{
				definitiontest.PlainsDefinition(plainsID, playerID),
				definitiontest.IslandDefinition(islandID, playerID),
				definitiontest.SwampDefinition(swampID, playerID),
				definitiontest.MountainDefinition(mountainID, playerID),
				definitiontest.ForestDefinition(forestID, playerID),
				definitiontest.WastesDefinition(wastesID, playerID),
			}
			game := newTestGame(
				playerID,
				"",
				permanents,
			)
			_, got, err := activateManaSource(game, playerID, test.landID)
			if err != nil {
				t.Fatalf("activateManaSource(); error = %v", err)
			}
			if diff := cmp.Diff(test.want, got, AllowAllUnexported); diff != "" {
				t.Errorf("activateManaSource(); events mismatch: %s", diff)
			}
		})
	}
}

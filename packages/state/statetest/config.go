package statetest

import (
	"deckronomicon/packages/game/gob/gobtest"
	"deckronomicon/packages/game/mtg"
)

// GameConfig holds the configuration for a game, allowing it to be loaded from a structured format.
// This structure is useful for saving and loading game states, or for initializing a game with predefined settings.
// It includes all necessary fields to reconstruct a game instance, such as player states, battlefield, stack, and game phase.
// These configurations should not be used for managing game state during active gameplay, but rather for testing, setup or serialization
// purposes

type GameConfig struct {
	NextID                int
	CheatsEnabled         bool
	ActivePlayerIdx       int
	PlayerWithPriority    string
	PlayersPassedPriority map[string]bool
	Battlefield           BattlefieldConfig
	Phase                 mtg.Phase
	Step                  mtg.Step
	Players               []PlayerConfig
	Stack                 StackConfig
	WinnerID              string
}

type PlayerConfig struct {
	ID                 string
	Life               int
	LandPlayedThisTurn bool
	Hand               HandConfig
	Library            LibraryConfig
	Graveyard          GraveyardConfig
	Exile              ExileConfig
	ManaPool           string
}

type HandConfig struct {
	Cards []gobtest.CardConfig
}

type BattlefieldConfig struct {
	Permanents []gobtest.PermanentConfig
}

type GraveyardConfig struct {
	Cards []gobtest.CardConfig
}

type LibraryConfig struct {
	Cards []gobtest.CardConfig
}

type ExileConfig struct {
	Cards []gobtest.CardConfig
}

type StackConfig struct {
}

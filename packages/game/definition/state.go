package definition

// GameConfig holds the configuration for a game, allowing it to be loaded from a structured format.
// This structure is useful for saving and loading game states, or for initializing a game with predefined settings.
// It includes all necessary fields to reconstruct a game instance, such as player states, battlefield, stack, and game phase.
// These configurations should not be used for managing game state during active gameplay, but rather for testing, setup or serialization
// purposes

type Battlefield struct {
	Permanents []*Permanent `json:"Permanents,omitempty" yaml:"Permanents,omitempty"`
}

type Exile struct {
	Cards []*Card `json:"Cards,omitempty" yaml:"Cards,omitempty"`
}

type Game struct {
	ActivePlayerID        string          `json:"ActivePlayerID,omitempty" yaml:"ActivePlayerID,omitempty"`
	Battlefield           *Battlefield    `json:"Battlefield" yaml:"Battlefield,omitempty"`
	CheatsEnabled         bool            `json:"CheatsEnabled,omitempty" yaml:"CheatsEnabled,omitempty"`
	NextID                int             `json:"NextID,omitempty" yaml:"NextID,omitempty"`
	Phase                 string          `json:"Phase,omitempty" yaml:"Phase,omitempty"`
	Players               []*Player       `json:"Players,omitempty" yaml:"Players,omitempty"`
	PlayersPassedPriority map[string]bool `json:"PlayersPassedPriority,omitempty" yaml:"PlayersPassedPriority,omitempty"`
	PlayerWithPriority    string          `json:"PlayerWithPriority,omitempty" yaml:"PlayerWithPriority,omitempty"`
	Step                  string          `json:"Step,omitempty" yaml:"Step,omitempty"`
	WinnerID              string          `json:"WinnerID,omitempty" yaml:"WinnerID,omitempty"`
}

type Graveyard struct {
	Cards []*Card `json:"Cards,omitempty" yaml:"Cards,omitempty"`
}

type Hand struct {
	Cards []*Card `json:"Cards,omitempty" yaml:"Cards,omitempty"`
}

type Library struct {
	Cards []*Card `json:"Cards,omitempty" yaml:"Cards,omitempty"`
}

type Player struct {
	Exile              *Exile     `json:"Exile" yaml:"Exile,omitempty"`
	Graveyard          *Graveyard `json:"Graveyard" yaml:"Graveyard,omitempty"`
	Hand               *Hand      `json:"Hand" yaml:"Hand,omitempty"`
	ID                 string     `json:"ID,omitempty" yaml:"ID,omitempty"`
	LandPlayedThisTurn bool       `json:"LandPlayedThisTurn,omitempty" yaml:"LandPlayedThisTurn,omitempty"`
	Library            *Library   `json:"Library" yaml:"Library,omitempty"`
	Life               int        `json:"Life,omitempty" yaml:"Life,omitempty"`
	ManaPool           string     `json:"ManaPool,omitempty" yaml:"ManaPool,omitempty"`
	MaxHandSize        int        `json:"MaxHandSize,omitempty" yaml:"MaxHandSize,omitempty"`
	Revealed           *Revealed  `json:"Revealed" yaml:"Revealed,omitempty"`
	SpellsCastThisTurn int        `json:"SpellsCastThisTurn,omitempty" yaml:"SpellsCastThisTurn,omitempty"`
	Turn               int        `json:"Turn,omitempty" yaml:"Turn,omitempty"`
}

type Revealed struct {
	Cards []*Card `json:"Cards,omitempty" yaml:"Cards,omitempty"`
}

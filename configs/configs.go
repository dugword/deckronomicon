package configs

import (
	"flag"
	"strconv"
)

type Config struct {
	CardPool     string
	DeckList     string
	Interactive  bool
	MaxTurns     int
	OnThePlay    bool
	StartingLife int
	StrategyFile string
	Verbose      bool
}

// LoadConfig loads the configuration from command line arguments and
// environment variables.
// TODO: Add support for reading the list of avilable decks if no deck list is
// specified.
// TODO: Maybe make the .json suffix optional and support typing in the name
// of the deck if it exists in the decks directory.
func LoadConfig(args []string, getenv func(string) string) (*Config, error) {
	var config Config
	if v, err := strconv.ParseBool(getenv("VERBOSE")); err == nil {
		config.Verbose = v
	}

	// Configure flags.
	flags := flag.NewFlagSet("deckronomicon", flag.ContinueOnError)
	// TODO: Maybe load both deck and strategy from a specified directory.
	cardPool := flags.String("card-pool", "cards", "card pool directory")
	deckList := flags.String("deck-list", "decks/example/deck.json", "deck list file")
	interactive := flags.Bool("interactive", false, "run as interactive mode")
	maxTurns := flags.Int("max-turns", 100, "maximum number of turns to simulate")
	onThePlay := flags.Bool("on-the-play", true, "player going first")
	startingLife := flags.Int("starting-life", 20, "starting life total")
	strategyFile := flags.String("strategy", "decks/example/strategy.json", "strategy file")
	verbose := flags.Bool("verbose", config.Verbose, "verbose output")

	if err := flags.Parse(args[1:]); err != nil {
		return nil, err
	}
	config.CardPool = *cardPool
	config.DeckList = *deckList
	config.Interactive = *interactive
	config.MaxTurns = *maxTurns
	config.OnThePlay = *onThePlay
	config.StartingLife = *startingLife
	config.StrategyFile = *strategyFile
	config.Verbose = *verbose

	return &config, nil
}

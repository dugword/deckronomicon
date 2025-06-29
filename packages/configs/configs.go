package configs

import (
	"flag"
	"strconv"
)

type Config struct {
	Definitions  string
	Cheat        bool
	Interactive  bool
	Scenario     string
	ScenariosDir string
	AutoPay      bool
	Verbose      bool
}

// LoadConfig loads the configuration from command line arguments and
// environment variables.
// TODO: Add support for reading the list of avilable decks if no deck list is
// specified.
// TODO: Maybe make the .json suffix optional and support typing in the name
// of the deck if it exists in the decks directory.
func LoadConfig(args []string, getenv func(string) string) (Config, error) {
	var config Config
	if v, err := strconv.ParseBool(getenv("VERBOSE")); err == nil {
		config.Verbose = v
	}
	// Configure flags.
	flags := flag.NewFlagSet("deckronomicon", flag.ContinueOnError)
	// TODO: Maybe load both deck and strategy from a specified directory.
	cheat := flags.Bool("cheat", false, "cheat mode")
	definitions := flags.String("definitions", "definitions/cards", "definitions directory")
	interactive := flags.Bool("interactive", false, "run as interactive mode")
	scenario := flags.String("scenario", "example", "scenario to run")
	scenariosDir := flags.String("scenarios", "scenarios", "scenarios directory")
	autoPay := flags.Bool("autopay", false, "automatically pay costs when possible")
	verbose := flags.Bool("verbose", config.Verbose, "verbose output")
	if err := flags.Parse(args[1:]); err != nil {
		return Config{}, err
	}
	config.Cheat = *cheat
	config.Definitions = *definitions
	config.Interactive = *interactive
	config.Scenario = *scenario
	config.ScenariosDir = *scenariosDir
	config.AutoPay = *autoPay
	config.Verbose = *verbose
	return config, nil
}

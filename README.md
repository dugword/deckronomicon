# Deckronomicon - Magic: The Gathering Simulator

![deckronomicon-logo-burst small](https://github.com/user-attachments/assets/64e4353e-bba0-4636-b199-53de3d746a87)

Join the Deckronomicon community on Discord: [Deckronomicon
Discord](https://discord.gg/Jb7Q9w6K)

## ⚠️ Warning

This project is in early development and is not yet ready for production use.

Features are subject to change, and the API is likely to change.

## Overview

Deckronomicon is a _Magic: The Gathering_ simulator and analysis tool designed
to run large numbers of automated games using configurable, rule-based
strategies. The results of the simulations can be used to evaluate and optimize
decklists for storm and combo archetypes like Eggs and High Tide, as well as
aggressive pump/burn strategies like Infect and Madness Burn.

The simulation engine supports both manual strategy exploration via an
interactive player agent as well as large-scale automated simulations.

## Getting Started and Usage

See the [Getting Started](docs/GettingStarted.md) guide for instructions on how to run
Deckronomicon.

See the [Usage](docs/Usage.md) guide for instructions on how to use
Deckronomicon.


## Why Deckronomicon?

Deckronomicon is designed to bridge the gap in decklist analysis between
tooling limited to statistical calculations and manual playtesting. A
Hypergeometric Calculator will give you the odds of drawing a specific set of
cards, but doesn't account for tutors/Scry/Surveil or other in game effects.
Goldfishing and practice games are useful for understanding how a deck plays,
but are time consuming and do not provide high-confidence results due to small
sample sizes.

By simulating large numbers of automated games, Deckronomicon enables
high-confidence results for determining the probability of specific game
outcomes when following a given strategy. This allows players to explore the
impact of deck building decisions, card choices, and play patterns on the
overall performance of a deck.


## Project Goals

* Simulate full-game sequences of _Magic: the Gathering_ decks.
* Manage automation strategy through composable rule-based configuration files.
* Record all game actions, events, and outcomes for analytics.
* Measure deck efficiency, recursion capability, and combo consistency.  
* Log game outcomes, card interactions, and resource usage.
* Prioritization of Pauper legal cards.

## Core Features

### Simulation Engine

Simulations are run using a deterministic event-based architecture enabling
true turn-by-turn replays and rollbacks of player actions.

The engine is designed to handle complex game states and interactions,
allowing for rules accurate simulations of _Magic: The Gathering_ games. Each
player is represented by a configurable strategy that defines how they play
their deck, including card priorities, land drop policies, and how to respond
to in game events.

### Game Record

Each simulation run generates a detailed game record that includes all player
actions, game events, state changes, and outcomes. This record can be used for
deterministic replays, post-game analysis, debugging, and strategy refinement.

### Configurable Rule-Based Strategy Automation

Player automation is driven by a rule-based strategy system that allows for
customizable and composable game plans. Strategies are defined in YAML,
enabling users to specify a prioritized list of actions and choices for the
player agent to follow during simulations.

The rule-based system supports complex decision-making based on the current
game state following a "When X Game State, Do Y Action" model. This allows for
flexible and adaptive strategies that can be tailored to specific deck types.

### Interactive Player Agent

The interactive player agent allows users to manually control a player during
simulations. This is useful for testing specific strategies, exploring card
interactions, and understanding how a deck performs.

## Contributing

If you are interested in contributing, please join the Discord community and
discuss your ideas. Contributions in the form of bug reports, feature requests,
and card implementation requests are highly encouraged.

Contributions to the codebase are mostly not yet accepted as the codebase is
still rapidly evolving. But if you have a specific feature or bug fix in mind,
please join the Discord and we can discuss it.

## Design and Architecture

See the [Design and Architecture](docs/Architecture.md) document for an
overview of how the Deckronomicon project is structured.

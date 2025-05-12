# 📘 Deckronomicon - Magic the Gathering Deck Analysis Engine

Deckronomicon is a Magic: The Gathering deck simulator and analysis tool designed to run thousands of full-game, turn-by-turn simulations using configurable, rule-based strategies. It helps evaluate and optimize decklists for storm and combo archetypes like Eggs and High Tide, as well as aggressive pump strategies like Infect or Kiln Fiend.

It supports both manual strategy exploration and large-scale automated simulations for data-driven analysis and machine learning–based optimization.

---

## 🎯 Project Goals

* Simulate full-game sequences of Magic the Gathering decks.
* Measure deck efficiency, recursion capability, and combo consistency.  
* Log game outcomes, card interactions, and resource usage.
* Enable tuning and evolution of deck strategies.
* Export simulation data for ML/analytics pipelines and visualiztions.
* Priortization of Pauper legal cards.

---

## 🧰 Core Features

### 🔁 Simulation Engine

* Turn-by-turn execution
* Rules aware game engine
    * Land drop enforcement
    * Spell and ability resolution
    * Mana pool management with smart tapping  
    * End-of-turn hand size enforcement (discard to 7)  
    * Game loss detection
* Game state tracking
    * Deck size
    * Life total
    * Storm count
    * Damage done
    * Opponent mill count

### 🧠 Configuraable Rule Based Strategy

* Priorities for cards played
* Priorities for abilities activated
* Land drop policy (greedy vs reserved)  
* Game state aware conditionals

### 🧩 Gamestate State Tracking

* Zones: hand, battlefield, graveyard, library, exile  
* Tapped/untapped states  
* Mana pool management
* Static Ability tracking
* Triggered Ability tracking

## ⚙️ Engine Mechanics

### 🔄 Game Logic

* Simulates games until win/lose condition detected
    * Deck is exhausted
    * Life total hits zero
    * Combo ends
    * Custom user defined win/lose conditions
* Tracks spells cast, lands played, and mana used per turn  
* Capable of running thousands of concurrent simulations
* Captures per-turn data and aggregate summaries  

### 📊 Logging and Metrics

* Logs spells/abilities/land drops per run  
* Tracks storm count, life gain, recursion loops  
* Outputs JSON logs and summaries  
* Saves top 10 runs in categories such as:  
  * Max storm (total & per-turn)  
  * Max life gained
  * Average fizzle rate
  * Average turn count to reach combo win state 

---

## 🖥️ Integration and Extensibility

* JSON/CSV export for data analysis  
* CLI interface for deck config, run control, and result inspection  
* Roadmap includes web UI for visualization  

---

## 🚧 Known Limitations

### 💡 Supported Deck Types
* Optimized for combo/storm decks that rely on recursion, draw chaining, and resource loops

### ⚠️ Not Ideal For:
* Aggro/tempo decks with minimal recursion or sequencing needs  
* Control decks that rely heavily on reactive play patterns  
* Decks where turn-by-turn non-determinism and opponent interaction are crucial

### Game Engine Limitations
* No stack management
* No opponent interaction
* Only cards implemented by this engine are supported

## 📚 Getting Started

```sh
git clone https://github.com/YOUR_USERNAME/deckronomicon.git
cd deckronomicon
go run ./
```
# Deckronomicon DESIGN.md

## :pushpin: Project Summary

**Deckronomicon** is a rules-accurate Magic: The Gathering simulator built in Go. It simulates full games with deterministic, event-driven, immutable game state. The simulator is designed for:
Testing rule-based strategies
Simulating turn-by-turn gameplay
Extracting game metrics
Training and evaluating machine learning agents

---

## :brain: Core Architecture Principles

### Event-Driven Immutable State

All state changes are triggered via `GameEvent` objects.
Events are applied through `engine.Engine.Apply(event)`.
Applying an event:
Records it in the game log
Calls a reducer to return a new `GameState` (pure function)
Enables deterministic replays and debugging
Triggers follow-up abilities/effects recursively

### Central Game Lifecycle

The main loop is driven by methods on the `engine.Engine`:
`RunGame`
`RunTurn`
`RunPhase`
`RunStep`
`RunPriority`

Each handles one slice of the game lifecycle. Turn-based actions and priority passes are all modeled explicitly as events.

---

## :bricks: Packages Overview

`engine`: Orchestrates the full game lifecycle and event application
`state`: Immutable game state and player state
`event`: All event definitions and dispatch logic
`action`: Turn-based player actions (e.g. Untap, Draw)
`player`: Player agent interface and rule-based decision system
`choose`: ChoicePrompt system for structured decision-making
`gameobject`: Card, spell, permanent, ability types
`mtg`: Magic: The Gathering terminology/constants
`logger`: Internal logging system with contextual, emoji-enhanced logs

---

## :black_joker: Cards, Effects, and Prompts

### JSON-Defined Cards

Each card is defined in JSON and includes:
Name, ManaCost, Color, CardTypes
SpellAbility referencing a list of reusable **effect IDs**

Example:
```json
{
  "Name": "Preordain",
  "ManaCost": "{U}",
  "SpellAbility": {
    "Effects": ["scry", "draw"]
  }
}
```

### Reusable Effect Implementations

Effects like `draw`, `scry`, `create_token` are implemented once and reused across cards using effect IDs and modifiers:
```json
{
  "Effect": "scry",
  "Modifiers": {
    "amount": 2
  }
}
```

### Structured Prompts

All choices (targets, card selections, etc.) are modeled using the `ChoicePrompt` system. Each prompt is:
A type implementing `ChoicePrompt`
Serializable for automation or UI
Resolvable by human input or AI agent

---

## :robot_face: Rule-Based Strategy Agent

### DSL Syntax

Strategies are JSON documents declaring **conditions** and **desired actions**:
```json
{
  "when": "<condition>",
  "mode": "combo",
  "actions": [...]
}
```

Within each action:
```json
{
  "when": "<condition>",
  "do": "cast spell",
  "target": "opponent creature with highest power"
}
```

### Functional Evaluation

The agent matches game state to rules and outputs actions or choice answers. It is:
Stateless and deterministic
Purely functional
Fully traceable and auditable

---

## :mag: Filters and Queries

Zone/card operations use functional, composable utilities:
`filter`, `map`, `find`, `take`, `drop`, etc.
Support for expressions like:
"Find all blue instants"
"Find a creature with the lowest toughness"

Used for:
Targeting
Strategy matching
AI rule evaluation

---

## :gear: Event Handling and Effects

Events are categorized (e.g. TurnStepEvents, PriorityEvents, ZoneChangeEvents)
Each has its own reducer logic
Base structs can be embedded to implement interfaces (`PriorityEvent`, etc.)
Effects like "until end of turn" are implemented with context-aware state

---

## :test_tube: Simulation + Machine Learning

The simulator supports:
Game scenarios (deck, agents, state) defined in JSON
High-speed simulations for metrics collection
Event logs with full game trace
Future integration of ML agents using rule or reward optimization

---

## :memo: Design Benefits

:boom: Deterministic and testable
:arrows_counterclockwise: Replayable for audits
:bricks: Composable by design
:technologist: Human and AI support for actions and choices
:brain: Perfect for tuning decks and learning strategy

---

## :white_check_mark: Still to Document

Triggered ability resolution flow
Stack resolution process
Priority pass handling in nested stack layers
Examples: Preordain, Brainstorm, Pieces of the Puzzle
Scenarios folder structure

---



---

## Execution Flow: Where Logic Should Live

To ensure clarity, maintainability, and separation of concerns, the simulator architecture separates logic across several key levels of abstraction. Here's a breakdown of where different types of logic belong:

---

### `engine.RunX` methods (e.g. `RunTurn`, `RunStep`, `RunPriority`)
**Responsibility:** High-level orchestration and flow control.

**What happens here:**
Drives the progression of the game (turns, phases, steps, etc.)
Calls `GetNextAction` from the PlayerAgent and handles its completion
Initiates `Resolve` or `Complete` on Actions, Spells, or Abilities
Applies events using `e.Apply`
Triggers game events such as start/end of turn, or priority passing
Manages choice resolution by calling into the ChoiceResolver
Responsible for the *gameplay flow*, not the logic of individual game mechanics

---

### `.Resolve()` / `.Complete()` methods
**Responsibility:** Generate `GameEvent`s from a game object (e.g. Spell, Ability, Action)

**What happens here:**
Based on current GameState and player choice, returns a list of `GameEvent`s
Resolves the behavior of a card or action (e.g. "Draw two cards")
Converts declarative effect data (from JSON) into actual game events
Is *pure logic* — does not mutate game state directly
May also include prompts for future decisions (e.g. which cards to discard)

---

### `engine.Apply(event)`
**Responsibility:** Central dispatcher that:
Records the event in the `GameRecord`
Applies the event immutably to produce a new GameState

**What happens here:**
Switches on event type/category
Applies relevant reducers such as `ApplyDrawEvent`, `ApplyUntapEvent`, etc.
May recursively apply triggered events
May enforce continuous effects (e.g. modifications to power/toughness)
Ensures *all* game state changes go through this function (for determinism and replayability)

---

### Reducers (e.g. `ApplyDrawEvent`, `ApplyCreateTokenEvent`)
**Responsibility:** The *only* place where GameState is updated

**What happens here:**
Takes the current GameState and returns a new immutable GameState
Does one and only one thing per event type
Responsible for applying continuous effects and checking for triggers
Performs no I/O or external coordination

---

### Strategy Rule Evaluation
**Responsibility:** Declarative rule system used by agents to make decisions

**What happens here:**
Evaluates rules like `when <condition> then <action>`
Each rule block corresponds to a specific mode (e.g. "Combo", "Defense", "DiscardDown")
Each action and choice within a rule may contain its own mini rule set
Not responsible for applying game state — it feeds into actions which will later generate events

---

Let me know if you’d like this structure visualized or want to break out each layer into its own file or section.

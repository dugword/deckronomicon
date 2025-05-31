# Deckronomicon Architecture

## Overview

Deckronomicon is a deterministic, auditable Magic: The Gathering game simulator and analysis engine. It supports both:

- Interactive human play (via command-line interface)
- Automated agent play (via Strategy JSON configurations)

The architecture is designed to be:

- Modular
- Deterministic (for reproducible simulations and ML training)
- Replayable (full game logs can be replayed exactly)
- Extensible (new cards, effects, agents easily added)

---

## Engine Lifecycle

```plaintext
GameEngine.RunGame
    → Loop: RunTurn
        → Loop: RunPhase
            → Loop: RunStep
                → Complete TurnBasedActions
            → Priority Loop: RunPriority
                → For each player:
                    → PlayerAgent.GetNextAction
                        → Action.Complete
                            → If valid → Apply Events
                            → If error → Reject Action → request new Action
    → Game ends when win/loss/draw condition met
```

---

## Core Principles

✅ The GameEngine only accepts **fully-formed Actions**.  
✅ The Engine **does not prompt** or fill in Action data.  
✅ All rule enforcement is centralized via a pure, stateless **Judge** package.  
✅ Action.Complete is atomic — either produces valid GameEvents or rejects the Action.  
✅ Apply applies **all events atomically**, based on immutable state updates.  
✅ Interactive and Automated agents use separate paths to build Actions.  
✅ The stack is resolved separately via Spell.Resolve, which revalidates targets and applies the spell’s effects.

---

## PlayerAgent Responsibilities

The PlayerAgent interface is responsible for proposing **complete Actions** to the Engine:

```go
type PlayerAgent interface {
    GetNextAction(state GameState) (Action, error)
}
```

- If the agent returns an incomplete Action, the Engine will reject it.
- If Action.Complete returns an error, the Engine will request another Action.

---

## Action Lifecycle

```plaintext
Agent builds Action → Engine calls Action.Complete → Apply Events
```

### Action Interface

```go
type Action interface {
    Complete(state GameState) ([]GameEvent, error)
    // Optional: String() for logging/debugging
}
```

- **Complete:** Responsible for final validation and GameEvent creation.
- Must use the Judge package to re-check legality.
- Must handle "game state changed" gracefully (e.g. targets no longer legal).

---

## Stack Resolution

When an Action puts a spell or ability on the stack:

```plaintext
Action.Complete → creates GameEvent_AddToStack

GameEngine stack loop:
    → Player passes priority
    → Engine pops top of stack → calls Spell.Resolve
        → Resolve re-checks targets using Judge
        → If targets valid → Resolve returns GameEvents to Apply
        → If targets invalid → spell/ability fizzles → no events → continue
```

- Target invalidation causing a "fizzle" is not an error — it is correct Magic behavior.
- Stack resolution is fully deterministic and auditable.

---

## Judge Package

The `judge` package contains all centralized rule enforcement logic:

Example API:

```go
func CanPlayCard(state GameState, player PlayerID, card CardID) bool
func CanPayCost(state GameState, player PlayerID, cost ManaCost, additionalCosts ...) bool
func IsTargetStillValid(state GameState, target Target) bool
func GetValidTargets(state GameState, player PlayerID, ability AbilityID) []Target
func CanActivateAbility(state GameState, player PlayerID, ability AbilityID) bool
```

Principles:

- Stateless and pure
- Safe to call from Agents, Action.Complete, and Spell.Resolve
- Single source of truth for Magic rule logic

---

## Apply Behavior

The Engine applies GameEvents in an **atomic, immutable fashion**:

```plaintext
Action.Complete → returns Events → Apply(Events) → GameState is updated
```

- Apply must be all-or-nothing.
- Apply must enforce hard game invariants (e.g. max hand size).
- Apply must not perform additional rule checks — Action.Complete and Judge ensure legality.
- Engine logs all applied Actions and Events for audit/replay.

---

## Agent Types

### Interactive Agent

```plaintext
User Input → Command Parser → Command → Action Builder → Action
    → If incomplete → Agent prompts user for choices → completes Action
    → Fully-formed Action → Engine
```

### Automated Agent

```plaintext
Strategy JSON → Strategy Engine → Action Builder → Fully-formed Action → Engine
```

- Automated agents must submit fully-formed Actions — no interactive prompts allowed.
- If an Action is rejected, Agent must propose a new Action.

---

## Strategy JSON Structure

Example:

```json
{
    "when": { ... },
    "then": {
        "action": "PlayCard",
        "card": "Preordain",
        "targets": [ ... ],
        "pay": { ... },
        "choices": [
            {
                "for": "ChooseMode",
                "when": { ... },
                "then": "Mode1"
            }
        ]
    }
}
```

- The `choices` section allows automated agents to fully resolve multi-step Actions.
- If a Strategy tries an invalid Action, Agent must avoid retry loops (see below).

---

## Automated Agent Retry Handling

Automated agents must take care to avoid infinite retry loops if an Action is rejected:

- Agents may track "last failed Action" and avoid repeating it in the same state.
- Agents may use Judge helpers to pre-check legality before proposing Actions.
- Agents may implement a "cooldown" or "backoff" for failed Actions.

---

## Logging and Replay

The Engine should log:

- All Actions proposed by Agents.
- Whether Actions were accepted or rejected.
- All GameEvents applied.
- Full GameState after each Apply.

This enables:

- Exact game replay for debugging.
- Deterministic ML training datasets.
- Full audit trail of simulation runs.

---

## Summary

Deckronomicon uses a clean, modular architecture with:

✅ Immutable GameState  
✅ Stateless Judge package  
✅ Atomic Action → GameEvents → Apply pipeline  
✅ Clear Agent / Engine contract  
✅ Stack resolution matching MTG rules  
✅ Deterministic logging and replayability  

This architecture supports:

- Accurate Magic simulation
- Large-scale automated testing
- ML-driven strategy tuning
- Interactive play for experimentation and testing

---

# Final Notes

This architecture provides a solid foundation for further enhancements, such as:

- More sophisticated Agent learning
- Multi-player support
- Comprehensive game rules coverage
- UI integration (optional)

Deckronomicon is designed to evolve while maintaining strong architectural clarity and determinism.

---

# End of Document


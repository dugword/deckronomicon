Architecture Overview
Deckronomicon uses a clean, modular architecture to support both interactive play and automated large-scale simulation of Magic: The Gathering decks.
The core engine is fully deterministic and replayable, with an immutable game state and centralized rule enforcement via a pure, stateless judge package.
Players (human or automated agents) propose complete Actions to the engine. The engine validates each Action, applies resulting GameEvents atomically, and manages the game flow and stack resolution.
Key design features:
:white_check_mark: Immutable GameState
:white_check_mark: Stateless Judge package for centralized rule logic
:white_check_mark: Clear separation of Engine, PlayerAgent, and Action responsibilities
:white_check_mark: Accurate Magic stack resolution (targets can change mid-stack)
:white_check_mark: Deterministic logging for ML and reproducibility
:white_check_mark: Flexible support for both interactive input and JSON-driven strategy agents

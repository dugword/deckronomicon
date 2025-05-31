Deckronomicon
Deckronomicon is a deterministic, rules-accurate Magic: The Gathering game simulator and analysis tool.
It is designed to run thousands of turn-by-turn games using configurable rule-based player strategies, enabling advanced analysis and optimization of deck performance. The system supports:
	•	Testing combo decks (e.g. High Tide, Tron Eggs, Storm)
	•	Evaluating control, prison, and tempo strategies
	•	Simulating aggro and pump-based lethal sequences
	•	Performing machine learning-driven strategy optimization
Deckronomicon uses a pure, auditable game engine with full support for stack resolution, triggered abilities, continuous effects, and complex player choices. All game actions are logged in a reproducible format to enable statistical analysis, replay, and tuning.





10:13
Project Goals
Deckronomicon is designed to support deep analysis and optimization of Magic: The Gathering decks and strategies.
Key goals:
	•	:dart: Simulate complex combo turns accurately
Model storm decks, self-mill loops, recursion engines, and other advanced sequences.
	•	:bar_chart: Analyze deck performance statistically
Track metrics such as:
	•	Average turn to combo assemble
	•	Average life gained, damage dealt, cards drawn per turn/game
	•	Fizzle rate of storm/combo turns
	•	Win patterns against opposing decks
	•	:brain: Support machine learning optimization
Provide deterministic, auditable logs for training ML models to optimize decklists and strategies.
	•	:joystick: Support both human and automated agents
Allow interactive exploration of decks as well as automated large-scale simulations.
	•	:recycle: Enable fully reproducible game logs
Guarantee that simulations can be exactly replayed for debugging, regression testing, and ML training.
	•	:mag: Maintain clean, extensible architecture
Use a pure Engine → Judge → Action → Apply pipeline with strong separation of concerns.









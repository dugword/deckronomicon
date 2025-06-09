iEvent Structure Principles
Deckronomicon models all game state changes via an explicit GameEvent system.
:point_right: This ensures that:
:white_check_mark: All state changes are logged and auditable
:white_check_mark: The engine is fully deterministic and replayable
:white_check_mark: The sequence of game state changes matches the Magic rules model
:white_check_mark: The UI (if added) can animate state changes accurately
:white_check_mark: The engine stays flexible and extensible
⸻
Granularity of Events
GameEvents should represent individual state changes, not compound “intention” actions.
	•	The engine should model the actual sequence of state changes that occur, even if multiple events are triggered by a single player action.
Example:
:point_right: A player activates Island’s “{T}: Add {U}” ability.
The resulting GameEvents should be:
7:07
GameEventTapPermanent → represents the cost being paid (tap Island)
GameEventAddMana → represents the effect being applied (add {U} to mana pool)
7:07
:point_right: Do NOT model this as a single GameEventTapLandForMana.
:point_right: This mixes concerns and prevents clean logging, animation, and reuse.
⸻
Why use granular events?
:white_check_mark: Composable → tapping a permanent is used in many contexts, not just adding mana
:white_check_mark: Auditable → logs can clearly show each state change that occurred
:white_check_mark: Rules-accurate → matches Magic’s model of cost payment and effect resolution
:white_check_mark: Testable → each GameEvent can be unit tested in isolation
:white_check_mark: UI-friendly → each GameEvent can trigger distinct animations
⸻
Design guideline:
	•	Paying costs (TapPermanent, SacrificePermanent, PayMana, PayLife, DiscardCard) → always generate explicit GameEvents.
	•	Effect resolution → also generates GameEvents separately.
	•	Do not combine paying cost and resolving effect into a single compound GameEvent.
⸻
Example GameEvents:
7:08
GameEventTapPermanent
GameEventAddMana
GameEventSacrificePermanent
GameEventPayLife
GameEventDiscardCard
GameEventDrawCard
GameEventMoveCardZone
GameEventGainLife
GameEventLoseLife
GameEventMillCards
7:08
Summary principle:
:point_right: In Deckronomicon, the event stream is the source of truth about what happened in the game.
:point_right: Each GameEvent should represent exactly one state change.
:point_right: This ensures that the engine remains:
	•	Deterministic
	•	Replayable
	•	Auditable
	•	Extensible
	•	Rules-accurate
7:10
Matching Magic’s Model: Cost Payment vs Effect Resolution
This event structure follows the rules model of Magic: The Gathering:
	•	Costs are paid first — before a spell or ability is placed on the stack (for spells) or before the ability resolves (for activated/triggered abilities).
	•	Effects are resolved separately, only after costs have been successfully paid.
:point_right: In Deckronomicon:
:white_check_mark: GameEvents that represent paying costs (TapPermanent, PayMana, SacrificePermanent, etc.) are generated first — during Action.Complete.
:white_check_mark: GameEvents that represent resolving effects (AddMana, DrawCard, etc.) are generated separately — during Spell.Resolve or Ability.Resolve.
This matches the Magic model exactly:
	•	If a player cannot pay a cost → the spell/ability is not activated/cast → no effects happen.
	•	If costs are paid → the spell/ability goes on the stack → effects resolve if the spell/ability resolves.
:point_right: By separating cost GameEvents from effect GameEvents in the event stream, Deckronomicon can correctly model this timing and sequencing.
⸻
Triggers and Replacement Effects
Granular GameEvents are also critical for future engine features:
	•	Triggered abilities → will listen to the GameEvent stream to know when to trigger.
	•	Example: A “Whenever you tap a permanent…” trigger will respond to GameEventTapPermanent.
	•	Example: A “Whenever you gain life…” trigger will respond to GameEventGainLife.
	•	Replacement effects → will modify or replace specific GameEvents as they are about to be applied.
	•	Example: “If you would gain life, prevent that life gain.” → intercepts GameEventGainLife.
	•	Example: “If a card would be put into your graveyard, exile it instead.” → intercepts GameEventMoveCardZone.
:point_right: If you model large compound events (like “TapLandForMana”), you make it impossible for triggers and replacement effects to work properly.
:point_right: Granular, composable GameEvents are the foundation that allows Triggers and Replacement Effects to work cleanly and correctly.
⸻
Summary:
:point_right: This design:
:white_check_mark: Matches Magic’s timing model of cost payment vs effect resolution.
:white_check_mark: Enables Triggers to watch for individual GameEvents.
:white_check_mark: Enables Replacement Effects to safely intercept or modify GameEvents.
:white_check_mark: Keeps the event stream fully auditable and replayable.


doug
  7:22 AM
EffectRegistry and EffectMetadata
Deckronomicon uses a data-driven effect system to ensure that:
:white_check_mark: All effects are clearly defined and documented
:white_check_mark: Player Agents can reason about available effects
:white_check_mark: The engine can validate card definitions and agent strategies
:white_check_mark: Effect resolution is deterministic, auditable, and safe
⸻
EffectMetadata JSON files
Every effect used by the engine must be defined in a JSON file in:
7:22
definitions/effects/
7:23
Each effect has an EffectMetadata JSON file, which defines:
	•	The effect name
	•	Required modifiers (key names expected in EffectSpec.Modifiers)
	•	Optional modifiers
	•	Required player input (choices or targets) if applicable
Example:
7:23
{
  "Name": "AddMana",
  "RequiredModifiers": ["Mana"],
  "OptionalModifiers": [],
  "RequiredChoices": [],
  "RequiredTargets": []
}
7:23
:point_right: Only effects that have a corresponding EffectMetadata file are allowed to be used in card definitions or activated abilities.
:point_right: This ensures that Player Agents can also load this metadata and reason about which choices/targets they must provide.
⸻
EffectRegistry
At engine startup:
	•	The engine loads all definitions/effects/*.json into an EffectMetadata map.
	•	The engine also builds an EffectRegistry, registering the corresponding EffectHandler function for each allowed effect.
Example EffectRegistry:
7:23
type EffectHandler func(state *GameState, source *EffectSource, modifiers map[string]string) ([]GameEvent, error)
type EffectRegistry struct {
    handlers map[string]EffectHandler
    metadata map[string]EffectMetadata
}
7:23
type EffectHandler func(state *GameState, source *EffectSource, modifiers map[string]string) ([]GameEvent, error)
type EffectRegistry struct {
    handlers map[string]EffectHandler
    metadata map[string]EffectMetadata
}
7:23
GameEventTapPermanent     // from cost payment
GameEventAddMana          // from effect resolution
7:24
:point_right: This ensures that:
:white_check_mark: Effect resolution is deterministic and auditable
:white_check_mark: Triggers and replacement effects can hook into individual GameEvents
:white_check_mark: The game state evolves in a sequence that matches Magic’s cost payment and effect resolution model
⸻
Summary Principle
:white_check_mark: Every effect used in the engine must be defined in EffectMetadata JSON in definitions/effects.
:white_check_mark: Only effects defined in JSON and registered in the EffectRegistry are allowed.
:white_check_mark: The engine will exit if an undefined effect is referenced.
:white_check_mark: Player Agents can use the EffectMetadata to validate and configure their strategy logic.
:white_check_mark: Effect resolution produces granular GameEvents consistent with the Magic rules model.


doug
  7:30 AM
What High Tide does (rules-wise):
	•	It creates a temporary replacement effect:
:point_right: Until end of turn,
:point_right: Whenever you tap an Island for mana,
:point_right: Add an additional {U}.
⸻
How this works in your model:
:white_check_mark: The key is that you already have granular events → GameEventTapPermanent and GameEventAddMana.
:point_right: But High Tide affects the AddMana effect, not just the tap itself.
:point_right: Specifically, it wants to replace or augment:
7:30
GameEventAddMana{ Mana: X }
→ add {U} if the source is an Island.
7:31
How this is modeled:
:white_check_mark: You need a Replacement Effect system — this is a very standard engine pattern.
:white_check_mark: You will model High Tide by adding a ReplacementEffect to the game state:
7:31
type ReplacementEffect struct {
    Duration Duration // Until end of turn
    AppliesTo func(event GameEvent) bool
    Replace func(event GameEvent) []GameEvent
}
7:31
When you apply GameEvents, your Apply function would look like:
7:31
for each GameEvent:
    for each active ReplacementEffect:
        if effect.AppliesTo(event):
            events := effect.Replace(event)
            process(events...) instead of original event
            break // only one replacement applies
    else:
        apply original event normally
7:31
High Tide’s ReplacementEffect:
When High Tide resolves, it would add this ReplacementEffect to the GameState:
7:31
ReplacementEffect{
    Duration: UntilEndOfTurn,
    AppliesTo: func(event GameEvent) bool {
        addManaEvent, ok := event.(GameEventAddMana)
        if !ok { return false }
        // Check: was the source tapped Island?
        if !IsIsland(addManaEvent.SourcePermanent) {
            return false
        }
        return true
    },
    Replace: func(event GameEvent) []GameEvent {
        addManaEvent := event.(GameEventAddMana)
        // Add additional {U}
        newManaPool := addManaEvent.Mana
        newManaPool.Add("{U}")
        return []GameEvent{
            GameEventAddMana{
                PlayerID: addManaEvent.PlayerID,
                SourcePermanent: addManaEvent.SourcePermanent,
                Mana: newManaPool,
            },
        }
    },
}
7:31
Summary of flow:
:white_check_mark: Player taps Island → pays cost → GameEventTapPermanent.
:white_check_mark: Ability resolves → GameEventAddMana is produced → SourcePermanent is Island.
:white_check_mark: When applying GameEventAddMana, active ReplacementEffects are checked:
	•	High Tide’s ReplacementEffect matches → replaces event → adds additional {U}.
:white_check_mark: The replacement happens during Apply, not during Action.Complete or Effect resolution — this ensures correct ordering and auditability.
⸻
Why this matches the rules perfectly:
:point_right: High Tide’s effect is a continuous replacement effect — it modifies what happens when an Island is tapped for mana.
:point_right: The correct engine model is:
	•	Tapping Island triggers AddMana → GameEventAddMana.
	•	ReplacementEffect modifies the AddMana as it is applied → adding the bonus {U}.
	•	This allows stacking multiple such effects correctly.
⸻
Summary principle:
:white_check_mark: Your granular events make this possible — because GameEventAddMana is a distinct event, you can intercept it.
:white_check_mark: The pattern of:
	•	ReplacementEffect.AppliesTo → Replace → produce new events
	•	During Apply phase
→ is how all major professional Magic engines model replacement effects.
:white_check_mark: High Tide is a textbook example — your current architecture supports this beautifully.
7:32
ReplacementEffect type
7:32
Inengine/replacement_effect.go
package engine
type ReplacementEffect struct {
    Duration Duration
    AppliesTo func(event GameEvent) bool
    Replace func(event GameEvent) []GameEvent
}
7:32
You can also add an ID or Source field if you want for logging/debugging:
7:32
ID string // e.g. "HighTide-123"
Source *Spell // or *Ability, or nil
7:32
:two: GameState holds active ReplacementEffects
You’ll want a field in GameState:
7:33
type GameState struct {
    // existing fields...
    ActiveReplacementEffects []ReplacementEffect
}
7:33
:three: Engine.Apply with ReplacementEffects
Here’s the Apply loop:
7:33
func (e *Engine) Apply(event GameEvent) {
    // Check for replacement effects first
    replaced := false
    for _, re := range e.state.ActiveReplacementEffects {
        if re.AppliesTo(event) {
            newEvents := re.Replace(event)
            // Optionally log:
            e.logReplacement(event, newEvents, re)
            // Process the replacement events
            for _, newEvt := range newEvents {
                e.Apply(newEvt)
            }
            replaced = true
            break // apply at most one replacement effect
        }
    }
    if replaced {
        return // original event not applied directly
    }
    // Normal Apply logic for this event:
    switch evt := event.(type) {
    case GameEventTapPermanent:
        e.applyTapPermanent(evt)
    case GameEventAddMana:
        e.applyAddMana(evt)
    // ... other cases ...
    default:
        panic(fmt.Sprintf("Unknown GameEvent type: %T", evt))
    }
}
7:33
:white_check_mark: This is the standard replacement effect pattern — matches professional engines.
⸻
:four: High Tide’s EffectHandler → installs ReplacementEffect
Now here’s how High Tide would work:
7:33
func HighTideEffectHandler(state *GameState, source *EffectSource, modifiers map[string]string) ([]GameEvent, error) {
    re := ReplacementEffect{
        Duration: UntilEndOfTurn,
        AppliesTo: func(event GameEvent) bool {
            addManaEvt, ok := event.(GameEventAddMana)
            if !ok {
                return false
            }
            // Check if SourcePermanent is Island
            if addManaEvt.SourcePermanent == nil {
                return false
            }
            if !IsIsland(addManaEvt.SourcePermanent) {
                return false
            }
            return true
        },
        Replace: func(event GameEvent) []GameEvent {
            addManaEvt := event.(GameEventAddMana)
            // Add an additional {U}
            newMana := addManaEvt.Mana
            newMana.Add("{U}")
            return []GameEvent{
                GameEventAddMana{
                    PlayerID: addManaEvt.PlayerID,
                    SourcePermanent: addManaEvt.SourcePermanent,
                    Mana: newMana,
                },
            }
        },
    }
    // Install the ReplacementEffect into GameState
    state.ActiveReplacementEffects = append(state.ActiveReplacementEffects, re)
    // No immediate GameEvents from High Tide itself
    return nil, nil
}
7:33
:white_check_mark: This matches the Magic rules perfectly:
	•	The ReplacementEffect is active until end of turn.
	•	While active, it modifies GameEventAddMana for Islands.
	•	Once removed at EOT, no further impact.
7:33
:five: Cleanup at End of Turn
You will want an Engine step like:
7:33
func (e *Engine) CleanupEndOfTurnEffects() {
    newList := e.state.ActiveReplacementEffects[:0]
    for _, re := range e.state.ActiveReplacementEffects {
        if re.Duration == UntilEndOfTurn {
            // Skip → removing this effect
            continue
        }
        newList = append(newList, re)
    }
    e.state.ActiveReplacementEffects = newList
}
7:34
:white_check_mark: You call this at the End Step / Cleanup Step of your turn loop.
⸻
:six: Why this is perfect for your architecture:
:white_check_mark: Uses your granular event stream → ReplacementEffects work naturally.
:white_check_mark: Preserves deterministic, auditable Apply flow.
:white_check_mark: Correctly models Magic’s ReplacementEffect timing.
:white_check_mark: High Tide is exactly this kind of effect — this pattern will work for:
	•	Doubling Season
	•	Leyline of Anticipation
	•	If a card would go to graveyard, exile instead
	•	Prevent life gain
	•	Replace draw with mill
:white_check_mark: Professional engines (MTGO, Arena) use this same model.


doug
  7:39 AM
Summary flow for High Tide:
:one: Player casts High Tide → resolves → adds ReplacementEffect to GameState.
:two: Player taps Island for mana:
GameEventTapPermanent → Apply → taps Island.
GameEventAddMana → Apply → High Tide’s ReplacementEffect triggers → modifies AddMana → adds extra {U}.
:three: At end of turn → CleanupEndOfTurnEffects → removes the ReplacementEffect.
7:42
What Ideas Unbound does (rules-wise):
:point_right: When the spell resolves, its effect is:
:white_check_mark: “Draw three cards.” → this happens immediately → produces 3x GameEventDrawCard.
:white_check_mark: “At the beginning of the next end step, discard three cards.” → this is a delayed triggered ability.
:point_right: The discard does not happen immediately → instead, a trigger is created and registered when the spell resolves.
⸻
How this works in your model:
:white_check_mark: You already have:
	•	Granular GameEvents
	•	ReplacementEffects → used for modifying GameEvents
	•	The same principle can apply to Triggers:
→ You add ActiveTriggeredAbilities to the GameState → they listen for GameEvents or GamePhases.
⸻
New type: TriggeredAbility
7:42
type TriggeredAbility struct {
    Duration Duration // UntilEndOfTurn, UntilNextTurn, Permanent, etc.
    TriggerCondition func(state *GameState, currentPhase Phase, currentStep Step, currentEvent GameEvent) bool
    TriggerEffect func(state *GameState) []GameEvent
}
7:42
Ideas Unbound’s TriggeredAbility
:white_check_mark: When Ideas Unbound resolves → you will:
:one: Produce GameEventDrawCard x3 immediately.
:two: Register a TriggeredAbility that watches for:
7:42
currentPhase == EndStep && currentStep == Beginning
7:42
And when triggered → produces:
7:42
GameEventDiscardCard x3
7:42
Example TriggeredAbility for Ideas Unbound:
7:43
TriggeredAbility{
    Duration: UntilEndOfTurn,
    TriggerCondition: func(state *GameState, phase Phase, step Step, evt GameEvent) bool {
        return phase == PhaseEndStep && step == StepBeginning
    },
    TriggerEffect: func(state *GameState) []GameEvent {
        playerID := ??? // You should carry SourcePlayer info here
        hand := state.GetPlayer(playerID).Hand
        cardsToDiscard := hand.TakeFirstN(3) // You can add Targeting later
        var events []GameEvent
        for _, card := range cardsToDiscard {
            events = append(events, GameEventDiscardCard{
                PlayerID: playerID,
                CardID:   card.ID(),
            })
        }
        return events
    },
}
7:43
:point_right: You add this TriggeredAbility to:
7:43
state.ActiveTriggeredAbilities = append(state.ActiveTriggeredAbilities, theTrigger)
7:43
Engine processing Triggers
:point_right: You want the Engine to do this during RunPhaseStep:
7:43
func (e *Engine) RunPhaseStep(phase Phase, step Step) {
    e.currentPhase = phase
    e.currentStep = step
    // Check triggered abilities:
    var newTriggeredEvents []GameEvent
    for _, trigger := range e.state.ActiveTriggeredAbilities {
        if trigger.TriggerCondition(&e.state, phase, step, nil) {
            events := trigger.TriggerEffect(&e.state)
            newTriggeredEvents = append(newTriggeredEvents, events...)
        }
    }
    // Process new triggered events:
    for _, evt := range newTriggeredEvents {
        e.Apply(evt)
    }
    // Now run normal step logic...
}
7:43
:point_right: You can also extend this to check GameEvents if you want triggers like “Whenever X happens.”
7:43
Cleanup of TriggeredAbilities
:white_check_mark: Same as ReplacementEffects:
7:43
func (e *Engine) CleanupEndOfTurnTriggers() {
    newList := e.state.ActiveTriggeredAbilities[:0]
    for _, trig := range e.state.ActiveTriggeredAbilities {
        if trig.Duration == UntilEndOfTurn {
            continue // remove it
        }
        newList = append(newList, trig)
    }
    e.state.ActiveTriggeredAbilities = newList
}
7:44
Summary flow for Ideas Unbound:
:white_check_mark: Player casts Ideas Unbound:
	•	Draw 3 cards → produces 3x GameEventDrawCard.
	•	Adds TriggeredAbility to GameState → watches for Beginning of End Step → discard 3 cards.
:white_check_mark: Later:
	•	Beginning of End Step triggers → TriggerEffect runs → produces 3x GameEventDiscardCard.
	•	Discard happens then.
:white_check_mark: After End Step → Cleanup → removes the TriggeredAbility.
⸻
Why this is perfect with your architecture:
:white_check_mark: Your granular events and RunPhaseStep loop allow Triggers to fire cleanly.
:white_check_mark: TriggeredAbility is a clean, simple struct — very extensible.
:white_check_mark: Works perfectly for Ideas Unbound and also:
	•	Delayed triggers like Exile at EOT
	•	“At the beginning of your upkeep…”
	•	“Whenever you gain life…”
	•	“Whenever a creature dies…” (listen to GameEvents instead of phase steps)
⸻
Summary principle:
:white_check_mark: Triggers should be modeled as:
	•	ActiveTriggeredAbilities → added to GameState
	•	Duration → controls when they are cleaned up
	•	TriggerCondition → watches for Phase/Step or GameEvents
	•	TriggerEffect → produces new GameEvents when the condition matches
:white_check_mark: Ideas Unbound is a textbook example of a delayed trigger — your engine is perfectly ready for this pattern!


doug
  8:54 AM
// engine/replacement_effect.go
package engine
type ReplacementEffect struct {
    Duration Duration
    AppliesTo func(event GameEvent) bool
    Replace func(event GameEvent) []GameEvent
}


Summary principle:
:white_check_mark: Paying a cost involves two phases:
:one: CanPayCost → pre-check, no state change, purely Judge logic.
:two: GenerateCostPaymentEvents → Action.Complete produces cost GameEvents → these are applied to change state.
8:57
GameEventMoveCardZone
GameEventTapPermanent
GameEventUntapPermanent
GameEventAttachAuraOrEquipment
GameEventDetachAuraOrEquipment
GameEventAddMana
GameEventSpendMana
GameEventGainLife
GameEventLoseLife
GameEventDealDamage
GameEventAddCounter
GameEventRemoveCounter
GameEventDrawCard
GameEventRevealCard
GameEventDiscardCard
GameEventPlayerConcede
GameEventPlayerLose
GameEventPlayerWin
GameEventSpellCast
GameEventSpellResolved
GameEventAbilityActivated
GameEventAbilityResolved
GameEventMillCard
GameEventShuffleLibrary
GameEventSetLifeTotal


doug
  9:12 AM
Cost Handling Design Principle
Deckronomicon models costs (mana costs, additional costs, activation costs) using a clean separation of concerns between the cost and judge packages.
⸻
Cost package
The cost package is a pure data model and parser:
	•	The Cost type represents the structured components of a cost:
	•	Mana cost
	•	Tap self
	•	Pay life
	•	Sacrifice self or other object
	•	Discard cards
	•	Other cost components
	•	The ParseCost function converts a cost string from card definitions into a structured Cost object.
:point_right: The cost package has no knowledge of the game state and performs no game rule logic.
:point_right: It is safe to import anywhere and is fully reusable.
⸻
Judge package
The judge package is the single authority for Magic rules enforcement.
It is responsible for answering all game rule questions, including cost-related ones:
	•	CanPayCost → determines whether the given player is currently able to pay a given Cost, given the current GameState.
	•	GenerateCostPaymentEvents → produces the ordered sequence of GameEvents required to pay the Cost.
:point_right: Only the judge package depends on both the cost package and the GameState.
:point_right: The rest of the engine always queries judge when it needs to reason about costs.
⸻
Rationale
:white_check_mark: This separation prevents circular dependencies.
:white_check_mark: It keeps cost pure and reusable.
:white_check_mark: It ensures that all cost-related rules logic lives in a single place → judge.
:white_check_mark: It matches the broader architecture pattern:
	•	Judge answers rules questions
	•	Actions/Spells use Judge to validate and model payment
	•	EffectHandlers do not process or reason about costs
⸻
Example flow
:one: Action.Complete calls judge.CanPayCost → validates that the player can legally pay the cost.
:two: If valid, Action.Complete calls judge.GenerateCostPaymentEvents → appends these events first in the GameEvent sequence.
:three: Effect resolution proceeds after cost GameEvents.
⸻
Summary
The engine cleanly separates data modeling of costs from rules logic of paying costs, using:





9:12
Concern
Package
Data model for Cost
cost
Parsing Cost from string
cost
Rules logic: CanPayCost
judge
Rules logic: GenerateCostPaymentEvents
judge
9:12
This ensures a professional, scalable architecture for handling Magic costs.











Jot something down




















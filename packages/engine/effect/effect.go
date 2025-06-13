package effect

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
)

// TODO: I have a rough idea and I'm not sure how to fully express or implement it yet, but
// I see there being some kind of loop where I process effects from the effects array of the ability, and each effect
// can return events and then also supply a choice prompt and a resume function to get the next set of events.
// Right now it feels like I have a slice of effects, and then separately a slice of effect results, but I think there
// might be a more elegaent way to express this.

type EffectResult struct {
	Events       []event.GameEvent
	ChoicePrompt choose.ChoicePrompt
	ResumeFunc   func(choose.ChoiceResults) (EffectResult, error)
}

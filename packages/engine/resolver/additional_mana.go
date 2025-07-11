package resolver

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
)

func ResolveAdditionalMana(
	playerID string,
	additionalMana *effect.AdditionalMana,
) (Result, error) {
	evnt := event.RegisterTriggeredAbilityEvent{
		PlayerID: playerID,
		Trigger: gob.Trigger{
			EventType: "LandTappedForMana",
			Filter: query.Opts{
				Subtypes: []mtg.Subtype{additionalMana.Subtype},
			},
		},
		Duration: additionalMana.Duration,
		Effects: []effect.Effect{&effect.AddMana{
			Mana: additionalMana.Mana,
		}},
	}
	return Result{
		Events: []event.GameEvent{&evnt},
	}, nil
}

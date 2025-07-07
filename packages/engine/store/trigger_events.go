package store

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
)

func GenerateTriggerEvents(oldGame *state.Game, newGame *state.Game, evnt event.GameEvent) []event.GameEvent {
	oldBattlefield := oldGame.Battlefield()
	newBattlefield := newGame.Battlefield()
	var triggerEvents []event.GameEvent
	for _, newPermanent := range newBattlefield.GetAll() {
		oldPermanent, found := oldBattlefield.Get(newPermanent.ID())
		if !found {
			triggerEvents = append(triggerEvents, &event.EnteredTheBattlefieldEvent{
				ControllerID: newPermanent.Controller(),
				PermanentID:  newPermanent.ID(),
				CardTypes:    newPermanent.CardTypes(),
				Subtypes:     newPermanent.Subtypes(),
				Supertypes:   newPermanent.Supertypes(),
			})
			continue
		}
		if !oldPermanent.IsTapped() && newPermanent.IsTapped() {
			triggerEvents = append(triggerEvents, event.TappedEvent{
				ControllerID: newPermanent.Controller(),
				PermanentID:  newPermanent.ID(),
				CardTypes:    newPermanent.CardTypes(),
				Subtypes:     newPermanent.Subtypes(),
				Supertypes:   newPermanent.Supertypes(),
			})
		}
		if oldPermanent.IsTapped() && !newPermanent.IsTapped() {
			triggerEvents = append(triggerEvents, event.UntappedEvent{
				ControllerID: newPermanent.Controller(),
				PermanentID:  newPermanent.ID(),
				CardTypes:    newPermanent.CardTypes(),
				Subtypes:     newPermanent.Subtypes(),
				Supertypes:   newPermanent.Supertypes(),
			})
		}
	}
	for _, oldPermanent := range oldBattlefield.GetAll() {
		fmt.Println("Checking old permanent:", oldPermanent.Name())
		_, found := newBattlefield.Get(oldPermanent.ID())
		if !found {
			triggerEvents = append(triggerEvents, &event.LeftTheBattlefieldEvent{
				ControllerID: oldPermanent.Controller(),
				PermanentID:  oldPermanent.ID(),
				CardTypes:    oldPermanent.CardTypes(),
				Subtypes:     oldPermanent.Subtypes(),
				Supertypes:   oldPermanent.Supertypes(),
			})
			for _, player := range newGame.Players() {
				for _, cardInGraveyard := range player.Graveyard().GetAll() {
					if cardInGraveyard.ID() == oldPermanent.Card().ID() {
						triggerEvents = append(triggerEvents, &event.DeathEvent{
							ControllerID: oldPermanent.Controller(),
							OwnerID:      oldPermanent.Owner(),
							PermanentID:  oldPermanent.ID(),
							CardID:       oldPermanent.Card().ID(),
							CardTypes:    oldPermanent.CardTypes(),
							Subtypes:     oldPermanent.Subtypes(),
							Supertypes:   oldPermanent.Supertypes(),
						})
					}
				}
			}
			continue
		}
	}
	return triggerEvents
}

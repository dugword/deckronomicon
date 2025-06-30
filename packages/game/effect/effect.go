package effect

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/target"
	"fmt"
)

type Effect interface {
	TargetSpec() target.TargetSpec
	Name() string
}

func New(definition definition.Effect) (Effect, error) {
	switch definition.Name {
	case "AddMana":
		return NewAddMana(definition.Modifiers)
	case "AdditionalMana":
		return NewAdditionalMana(definition.Modifiers)
	case "Counterspell":
		return NewCounterspell(definition.Modifiers)
	case "Discard":
		return NewDiscard(definition.Modifiers)
	case "Draw":
		return NewDraw(definition.Modifiers)
	case "LookAndChoose":
		return NewLookAndChoose(definition.Modifiers)
	case "Mill":
		return NewMill(definition.Modifiers)
	case "PutBackOnTop":
		return NewPutBackOnTop(definition.Modifiers)
	case "RegisterDelayedTriggeredAbility":
		return NewRegisterDelayedTriggeredAbility(definition.Modifiers)
	case "Replicate":
		return NewReplicate(definition.Modifiers)
	case "Scry":
		return NewScry(definition.Modifiers)
	case "Search":
		return NewSearch(definition.Modifiers)
	case "ShuffleFromGraveyard":
		return NewShuffleFromGraveyard(definition.Modifiers)
	case "Tap":
		return NewTap(definition.Modifiers)
	case "TapOrUntap":
		return NewTapOrUntap(definition.Modifiers)
	case "Target":
		return NewTarget(definition.Modifiers)
	default:
		return nil, fmt.Errorf("unknown effect %s", definition.Name)
	}
}

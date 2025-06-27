package resolver

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/state"
	"fmt"
)

type Result struct {
	Events       []event.GameEvent
	ChoicePrompt choose.ChoicePrompt
	Resume       func(choose.ChoiceResults) (Result, error)
}

func Resolve(
	game state.Game,
	playerID string,
	source gob.Object,
	effectWithTarget effect.EffectWithTarget,
	resEnv *resenv.ResEnv,
) (Result, error) {
	target := effectWithTarget.Target
	switch efct := effectWithTarget.Effect.(type) {
	case effect.AdditionalMana:
		return ResolveAdditionalMana(playerID, efct)
	case effect.AddMana:
		return ResolveAddMana(game, playerID, efct)
	case effect.Counterspell:
		return ResolveCounterspell(game, playerID, efct, target)
	case effect.Discard:
		return ResolveDiscard(game, playerID, efct, target, source)
	case effect.Draw:
		return ResolveDraw(game, playerID, efct, target)
	case effect.LookAndChoose:
		return ResolveLookAndChoose(game, playerID, efct, source)
	case effect.Mill:
		return ResolveMill(game, playerID, efct, target)
	case effect.PutBackOnTop:
		return ResolvePutBackOnTop(game, playerID, efct, source)
	case effect.RegisterDelayedTriggeredAbility:
		return ResolveRegisterDelayedTriggeredAbility(playerID, efct, source)
	case effect.Replicate:
		return ResolveReplicate(game, playerID, efct, source, resEnv)
	case effect.Scry:
		return ResolveScry(game, playerID, efct, source)
	case effect.Search:
		return ResolveSearch(game, playerID, efct, source)
	case effect.ShuffleFromGraveyard:
		return ResolveShuffleFromGraveyard(game, playerID, efct, source, resEnv)
	case effect.Tap:
		return ResolveTap(game, playerID)
	case effect.TapOrUntap:
		return ResolveTapOrUntap(game, playerID, efct, target)
	case effect.TargetEffect:
		return ResolveTarget(game, playerID)
	default:
		return Result{}, fmt.Errorf("unknown effect type: %T", efct)
	}
}

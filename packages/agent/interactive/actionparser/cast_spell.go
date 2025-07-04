package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/engine/judge"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"fmt"
)

func parseCastSpellCommand(
	idOrName string,
	game *state.Game,
	playerID string,
	agent engine.PlayerAgent,
	autoPayCost bool,
	autoPayColors []mana.Color, // Colors to prioritize when auto-paying costs, if applicable
) (action.CastSpellRequest, error) {
	ruling := judge.Ruling{Explain: true}
	cards := judge.GetSpellsAvailableToCast(game, playerID, autoPayCost, autoPayColors, &ruling)
	var cardInZone *gob.CardInZone
	var err error
	if idOrName == "" {
		cardInZone, err = getCardByChoice(cards, agent, game, playerID)
		if err != nil {
			return action.CastSpellRequest{}, fmt.Errorf("failed to get card by choice: %w", err)
		}
	} else {
		found, ok := query.Find(cards, query.Or(has.ID(idOrName), has.Name(idOrName)))
		if !ok {
			return action.CastSpellRequest{}, fmt.Errorf(
				"failed to find card found with id or name %q: %w", idOrName, ErrCardNotFound,
			)
		}
		cardInZone = found
	}
	costTarget, err := getTargetsForCost(
		playerID,
		cardInZone.Card(),
		game,
		agent,
	)
	if err != nil {
		return action.CastSpellRequest{}, fmt.Errorf("failed to get cost target: %w", err)
	}
	targetsForEffects, err := getTargetsForEffects(
		playerID,
		cardInZone.Card(),
		cardInZone.Card().SpellAbility(),
		game,
		agent,
	)
	if err != nil {
		return action.CastSpellRequest{}, fmt.Errorf("failed to get targets for spell: %w", err)
	}
	// TODO: Pass in cost or something and only get cards that can be paid for.
	// TODO: Handle autoPayCost and autoPayColors properly.
	splicableCards, err := judge.GetSplicableCards(game, playerID, cardInZone, &ruling)
	if err != nil {
		return action.CastSpellRequest{}, fmt.Errorf("failed to get splicable cards: %w", err)
	}
	cardsToSplice, err := chooseSpliceCards(splicableCards, cardInZone.Card(), game, agent)
	if err != nil {
		return action.CastSpellRequest{}, fmt.Errorf("failed to choose splice cards: %w", err)
	}
	var spliceCardIDs []string
	for _, cardToSplice := range cardsToSplice {
		spliceTargetsForEffects, err := getTargetsForEffects(
			playerID,
			cardToSplice.Card(),
			cardToSplice.Card().SpellAbility(),
			game,
			agent,
		)
		if err != nil {
			return action.CastSpellRequest{}, fmt.Errorf("failed to get targets for spell: %w", err)
		}
		for k, v := range spliceTargetsForEffects {
			targetsForEffects[k] = v
		}
		spliceCardIDs = append(spliceCardIDs, cardToSplice.Card().ID())
	}
	flashback := false
	if cardInZone.Zone() == mtg.ZoneGraveyard {
		flashback = true
	}
	replicateCount := 0
	if judge.CanReplicateCard(game, playerID, cardInZone.Card(), &ruling) {
		var err error
		replicateCount, err = getReplicateCount(cardInZone.Card(), game, agent)
		if err != nil {
			return action.CastSpellRequest{}, fmt.Errorf("failed to get replicate count: %w", err)
		}
	}
	request := action.CastSpellRequest{
		CardID:            cardInZone.Card().ID(),
		TargetsForEffects: targetsForEffects,
		ReplicateCount:    replicateCount,
		SpliceCardIDs:     spliceCardIDs,
		Flashback:         flashback,
		AutoPayCost:       autoPayCost, // WARNING: This is turned on for testing, it doesn't work properly yet.,
		AutoPayColors:     autoPayColors,
		CostTarget:        costTarget,
	}
	return request, nil
}

func getReplicateCount(
	card *gob.Card,
	game *state.Game,
	agent engine.PlayerAgent,
) (int, error) {
	replicatePrompt := choose.ChoicePrompt{
		Message:    "Choose how many times to replicate the spell",
		Source:     card,
		Optional:   true,
		ChoiceOpts: choose.ChooseNumberOpts{},
	}
	replicateResults, err := agent.Choose(game, replicatePrompt)
	if err != nil {
		return 0, fmt.Errorf("failed to get replicate count: %w", err)
	}
	selected, ok := replicateResults.(choose.ChooseNumberResults)
	if !ok {
		return 0, fmt.Errorf("expected a number choice result for replicate count")
	}
	if selected.Number < 0 {
		return 0, fmt.Errorf("replicate count cannot be negative")
	}
	return selected.Number, nil
}

func chooseSpliceCards(
	splicableCards []*gob.CardInZone,
	card *gob.Card,
	game *state.Game,
	agent engine.PlayerAgent,
) ([]*gob.CardInZone, error) {
	var cardsInHandToSplice []*gob.CardInZone
	if len(splicableCards) == 0 {
		return nil, nil // No splicable cards available
	}
	splicePrompt := choose.ChoicePrompt{
		Message:  "Choose cards to splice onto the spell",
		Source:   card,
		Optional: true,
		ChoiceOpts: choose.ChooseManyOpts{
			Choices: choose.NewChoices(splicableCards),
			Min:     0,
			Max:     len(splicableCards),
		},
	}
	spliceResults, err := agent.Choose(game, splicePrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get splice choices: %w", err)
	}
	selected, ok := spliceResults.(choose.ChooseManyResults)
	if !ok {
		return nil, fmt.Errorf("expected a multiple choice result for splicing")
	}
	for _, choice := range selected.Choices {
		card, ok := choice.(*gob.CardInZone)
		if !ok {
			return nil, fmt.Errorf("selected choice is not a card in a zone")
		}
		cardsInHandToSplice = append(cardsInHandToSplice, card)
	}
	return cardsInHandToSplice, nil
}

func getCardByChoice(
	cards []*gob.CardInZone,
	agent engine.PlayerAgent,
	game *state.Game,
	playerID string,
) (*gob.CardInZone, error) {
	prompt := choose.ChoicePrompt{
		Message:  "Choose a spell to cast",
		Source:   nil, // TODO: Make this better
		Optional: true,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(cards),
		},
	}
	choiceResults, err := agent.Choose(game, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get choices: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return nil, fmt.Errorf("expected ChooseOneResults, got %T", choiceResults)
	}
	if selected.Choice == nil {
		return nil, fmt.Errorf("no card selected: %w", choose.ErrNoChoiceSelected)
	}
	card, ok := selected.Choice.(*gob.CardInZone)
	if !ok {
		return nil, fmt.Errorf("selected choice is not a card in a zone")
	}
	return card, nil
}

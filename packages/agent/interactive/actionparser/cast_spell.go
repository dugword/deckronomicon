package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/engine/judge"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"fmt"
)

func parseCastSpellCommand(
	idOrName string,
	game state.Game,
	player state.Player,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
) (action.CastSpellAction, error) {
	ruling := judge.Ruling{Explain: true}
	cards := judge.GetSpellsAvailableToCast(game, player, &ruling)
	var cardInZone gob.CardInZone
	var err error
	if idOrName == "" {
		cardInZone, err = getCardByChoice(cards, chooseFunc, player)
		if err != nil {
			return action.CastSpellAction{}, fmt.Errorf("failed to get card by choice: %w", err)
		}
	} else {
		found, ok := query.Find(cards, query.Or(has.ID(idOrName), has.Name(idOrName)))
		if !ok {
			return action.CastSpellAction{}, fmt.Errorf("no spell found with id or name %q", idOrName)
		}
		cardInZone = found
	}
	targetsForEffects, err := getTargetsForEffects(
		cardInZone.Card(),
		cardInZone.Card().SpellAbility(),
		game,
		chooseFunc,
	)
	if err != nil {
		return action.CastSpellAction{}, fmt.Errorf("failed to get targets for spell: %w", err)
	}
	// TODO: Pass in cost or something and only get cards that can be paid for.
	splicableCards, err := judge.GetSplicableCards(game, player, cardInZone, &ruling)
	if err != nil {
		return action.CastSpellAction{}, fmt.Errorf("failed to get splicable cards: %w", err)
	}
	fmt.Println("Ruling for splice:", ruling.Why())
	cardsToSplice, err := chooseSpliceCards(splicableCards, cardInZone.Card(), chooseFunc)
	if err != nil {
		return action.CastSpellAction{}, fmt.Errorf("failed to choose splice cards: %w", err)
	}
	var spliceCardIDs []string
	for _, cardToSplice := range cardsToSplice {
		spliceTargetsForEffects, err := getTargetsForEffects(
			cardToSplice.Card(),
			cardToSplice.Card().SpellAbility(),
			game,
			chooseFunc,
		)
		if err != nil {
			return action.CastSpellAction{}, fmt.Errorf("failed to get targets for spell: %w", err)
		}
		for k, v := range spliceTargetsForEffects {
			targetsForEffects[k] = v
		}
		spliceCardIDs = append(spliceCardIDs, cardToSplice.Card().ID())
	}
	withFlashback := false
	if cardInZone.Zone() == mtg.ZoneGraveyard {
		withFlashback = true
	}
	request := action.CastSpellRequest{
		CardID:            cardInZone.Card().ID(),
		TargetsForEffects: targetsForEffects,
		SpliceCardIDs:     spliceCardIDs,
		WithFlashback:     withFlashback,
	}
	return action.NewCastSpellAction(
		player.ID(),
		request,
	), nil
}

func chooseSpliceCards(
	splicableCards []gob.CardInZone,
	card gob.Card,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
) ([]gob.CardInZone, error) {
	var cardsInHandToSplice []gob.CardInZone
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
	spliceResults, err := chooseFunc(splicePrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get splice choices: %w", err)
	}
	selected, ok := spliceResults.(choose.ChooseManyResults)
	if !ok {
		return nil, fmt.Errorf("expected a multiple choice result for splicing")
	}
	for _, choice := range selected.Choices {
		card, ok := choice.(gob.CardInZone)
		if !ok {
			return nil, fmt.Errorf("selected choice is not a card in a zone")
		}
		cardsInHandToSplice = append(cardsInHandToSplice, card)
	}
	return cardsInHandToSplice, nil
}

func getCardByChoice(
	cards []gob.CardInZone,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
	player state.Player,
) (gob.CardInZone, error) {
	prompt := choose.ChoicePrompt{
		Message:  "Choose a spell to cast",
		Source:   player,
		Optional: true,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(cards),
		},
	}
	choiceResults, err := chooseFunc(prompt)
	if err != nil {
		return gob.CardInZone{}, fmt.Errorf("failed to get choices: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return gob.CardInZone{}, fmt.Errorf("selected choice is not a card in a zone")
	}
	card, ok := selected.Choice.(gob.CardInZone)
	if !ok {
		return gob.CardInZone{}, fmt.Errorf("selected choice is not a card in a zone")
	}
	return card, nil
}

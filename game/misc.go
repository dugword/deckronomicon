package game

func ActionDiscardFunc(state *GameState, resolver ChoiceResolver) (result *ActionResult, err error) {
	choices := state.Hand.CardChoices()
	choice := resolver.ChooseOne("Which card to discard from hand", choices)
	card := state.Hand.GetCard(choice.Index)
	if card != nil {
		card := state.Hand.GetCard(choice.Index)
		if card != nil {
			state.Hand.RemoveCard(card)
			state.Graveyard = append(state.Graveyard, card)
		}
	}
	// TODO: Maybe every action resultion should log and put a status message
	return nil, nil
}

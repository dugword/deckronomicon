package engine

import (
	"deckronomicon/packages/game/player"
)

// TODO This sucks here, find a better place

// chooseManaForGeneric prompts the user to choose mana for a generic cost.
func chooseManaForGeneric(genericCost int, player *player.Player) (map[string]int, error) {
	/*
		choice := make(map[string]int)
		remaining := genericCost
		tempPool := player.ManaPool().Copy()
		for remaining > 0 {
			choices := []Choice{}
			for _, color := range tempPool.ColorsAvailable() {
				choices = append(choices, Choice{
					Name: string(color),
				})
			}
			if len(choices) == 0 {
				return nil, errors.New("insufficient mana")
			}
			selected, err := player.Agent.ChooseOne(
				"Choose mana to use for generic cost",
				NewChoiceSource("Pay Generic Cost", "Pay Generic Cost"),
				choices,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to choose mana for generic cost: %w", err)
			}
			// TODO: maybe make this a function
			color, err := StringToColor(selected.Name)
			if err != nil {
				return nil, fmt.Errorf("failed to apply choice for generic mana: %w", err)
			}
			tempPool.Use(color, 1)
			choice[selected.Name]++
			remaining--
		}
		return choice, nil
	*/
	return nil, nil
}

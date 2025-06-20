package main

import (
	"deckronomicon/packages/engine/effect"
	"deckronomicon/packages/game/definition"
	"errors"
	"fmt"
	"os"
)

func main() {
	if err := Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func Run() error {
	definitions, err := definition.LoadCardDefinitions("./definitions")
	if err != nil {
		return err
	}
	var errs []error
	for name, def := range definitions {
		if validationErrors := validateCard(def); len(validationErrors) > 0 {
			errs = append(errs, fmt.Errorf("card %q has errors", name))
			errs = append(errs, validationErrors...)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func validateCard(card definition.Card) []error {
	var errs []error
	if card.Name == "" {
		errs = append(errs, fmt.Errorf("card name is required"))
		fmt.Println("Missing name in card:", card)
	}
	if len(card.CardTypes) == 0 {
		errs = append(errs, fmt.Errorf("at least one card type is required"))
	}
	for _, effectSpec := range card.SpellAbilitySpec {
		if _, err := effect.Build(effectSpec); err != nil {
			errs = append(errs, fmt.Errorf("invalid effect spec in card %q: %w", card.Name, err))
			continue
		}
	}
	return errs
}

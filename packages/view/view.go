package view

import (
	"deckronomicon/packages/engine/judge"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
)

type Game struct {
	ActivePlayerID string
	Phase          mtg.Phase
	Step           mtg.Step
	Battlefield    []Permanent
	Stack          []Resolvable
}

type Player struct {
	ID   string
	Life int
	Mode string
	Turn int
	// TODO: Change ManaPool to a more structured type
	// and do the formatting in the UI.
	ManaPool          string
	PotentialManaPool string
	Hand              []Card
	Graveyard         []Card
	Exile             []Card
	Revealed          []Card
	LibrarySize       int
}

type Permanent struct {
	Name          string
	Controller    string
	ID            string
	Tapped        bool
	SummoningSick bool
}

type Resolvable struct {
	Name       string
	ID         string
	Controller string
}

// Describe returns a string representation of the mana pool
type Card struct {
	Name string
	ID   string
}

func NewGameViewFromState(game state.Game) Game {
	return Game{
		ActivePlayerID: game.ActivePlayerID(),
		Phase:          game.Phase(),
		Step:           game.Step(),
		Battlefield:    permanentsToViewPermanents(game.Battlefield().GetAll()),
		Stack:          resolvablesToViewResolvables(game.Stack().GetAll()),
	}
}

func resolvablesToViewResolvables(resolvables []state.Resolvable) []Resolvable {
	var viewResolvables []Resolvable
	for _, resolvable := range resolvables {
		viewResolvables = append(viewResolvables, Resolvable{
			Name:       resolvable.Name(),
			ID:         resolvable.ID(),
			Controller: resolvable.Controller(),
		})
	}
	return viewResolvables
}

func cardsToViewCards(cards []gob.Card) []Card {
	var viewCards []Card
	for _, card := range cards {
		viewCards = append(viewCards, Card{
			Name: card.Name(),
			ID:   card.ID(),
		})
	}
	return viewCards
}

func permanentsToViewPermanents(permanents []gob.Permanent) []Permanent {
	var viewPermanents []Permanent
	for _, perm := range permanents {
		viewPermanents = append(viewPermanents, Permanent{
			Name:          perm.Name(),
			Controller:    perm.Controller(),
			ID:            perm.ID(),
			Tapped:        perm.IsTapped(),
			SummoningSick: perm.HasSummoningSickness(),
		})
	}
	return viewPermanents
}

func NewPlayerViewFromState(game state.Game, player state.Player, mode string) Player {
	manaPool := player.ManaPool().ManaString()
	if manaPool == "" {
		manaPool = "(empty)"
	}
	potentialManaPool := judge.GetAvailableMana(game, player).ManaString()
	if potentialManaPool == "" {
		potentialManaPool = "(empty)"
	}
	return Player{
		ID:                player.ID(),
		Life:              player.Life(),
		Mode:              mode,
		Turn:              player.Turn(),
		ManaPool:          manaPool,
		PotentialManaPool: potentialManaPool,
		Hand:              cardsToViewCards(player.Hand().GetAll()),
		Graveyard:         cardsToViewCards(player.Graveyard().GetAll()),
		Exile:             cardsToViewCards(player.Exile().GetAll()),
		Revealed:          cardsToViewCards(player.Revealed().GetAll()),
		LibrarySize:       player.Library().Size(),
	}
}

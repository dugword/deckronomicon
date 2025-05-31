package engine

import (
	"deckronomicon/packages/configs"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/player"
	"deckronomicon/packages/query/has"
	"fmt"
	"strings"
)

type Engine struct {
	Scenario        *configs.Scenario
	GameState       *GameState
	CardDefinitions map[string]definition.Card
}

// TODO: Figure out better naming so I don't have plyr
// InitializeNewGame initializes a new game with the given configuration.
func InitializeNewGame(
	scenario *configs.Scenario,
	players []*player.Player,
	cardDefinitions map[string]definition.Card,
) *Engine {
	gameState := NewGameState()
	gameState.MaxTurns = scenario.Setup.MaxTurns
	gameState.Players = players
	for _, plyer := range players {
		if plyer.ID() == scenario.Setup.OnThePlay {
			gameState.ActivePlayer = plyer
			continue
		}
		gameState.NonActivePlayers = append(gameState.NonActivePlayers, plyer)
	}
	gameState.CheatsEnabled = scenario.Setup.CheatsEnabled
	engine := Engine{
		CardDefinitions: cardDefinitions,
		GameState:       gameState,
		Scenario:        scenario,
	}
	return &engine
}

// Log logs a message to the game state message log.
func (e *Engine) Log(message ...string) {
	// TODO: Something like this so we know which player to which action
	// outMessage := fmt.Sprintf("%s: %s", g.ActivatePlayer.ID )
	// TODO: There's probably a more elegant way to do this
	e.GameState.TurnMessageLog = append(e.GameState.TurnMessageLog, strings.Join(message, " "))
	e.GameState.MessageLog = append(e.GameState.MessageLog, strings.Join(message, " "))
}

// Error logs an error message to the game state message log.
// TODO: There's probably a more elegant way to do this
func (e *Engine) Error(err error) {
	e.GameState.TurnMessageLog = append(e.GameState.TurnMessageLog, "ERROR: "+err.Error())
	e.GameState.MessageLog = append(e.GameState.MessageLog, "ERROR: "+err.Error())
}

func DrawStartingHand(plyer *player.Player, startingHand []string) ([]string, error) {
	var cardsDrawn []string
	for _, cardName := range startingHand {
		if err := plyer.Tutor(has.Name(cardName)); err != nil {
			return nil, fmt.Errorf("failed to tutor card %s: %w", cardName, err)
		}
		cardsDrawn = append(cardsDrawn, cardName)
	}
	remainingDraws := plyer.MaxHandSize - plyer.Hand.Size()
	for range remainingDraws {
		cardName, err := plyer.Draw()
		if err != nil {
			return nil, fmt.Errorf("failed to draw card: %w", err)
		}
		cardsDrawn = append(cardsDrawn, cardName)
	}
	return cardsDrawn, nil
}

// RunGameLoop drives the turn-based game.
func (e *Engine) RunGameLoop() error {
	for _, plyer := range e.GameState.Players {
		// TODO: Should this be a method on Engine?
		playerLibrary, err := BuildDeck(
			e.GameState,
			e.Scenario.Players[plyer.ID()].DeckList,
			e.CardDefinitions,
		)
		if err != nil {
			return err
		}
		plyer.Library = playerLibrary
		plyer.Library.Shuffle()

	}
	{
		cardNames, err := DrawStartingHand(plyer, scenario.Setup.Player.StartingHand)
		if err != nil {
			return fmt.Errorf("%s failed to draw starting hand: %w", plyer.ID(), err)
		}
		g.Log(fmt.Sprintf("% drew: %s", plyer.ID(), strings.Join(cardNames, ", ")))
		cardNames, err = DrawStartingHand(opponent, scenario.Setup.Opponent.StartingHand)
		if err != nil {
			return fmt.Errorf("%s failed to draw starting hand: %w", opponent.ID(), err)
		}
		g.Log(fmt.Sprintf("% drew: %s", plyer.ID(), strings.Join(cardNames, ", ")))
	}
	if err := g.Mulligan(g.ActivePlayer); err != nil {
		return fmt.Errorf("active player failed to mulligan: %w", err)
	}
	if err := g.Mulligan(g.NonActivePlayer); err != nil {
		return fmt.Errorf("non-active player failed to mulligan: %w", err)
	}
	for {
		g.ActivePlayer, g.NonActivePlayer = g.NonActivePlayer, g.ActivePlayer
		if g.ActivePlayer.Turn > g.MaxTurns {
			return ErrMaxTurnsExceeded
		}
		if err := g.RunTurn(g.ActivePlayer); err != nil {
			return err
		}
	}
}

func (g *GameState) RunStep(step GameStep, player *player.Player) error {
	g.CurrentStep = step.Name
	g.Log(fmt.Sprintf("Running step: %s", step.Name))
	// g.EmitEvent(Event{Type: step.EventEvent}, player)
	if err := step.Handler(g, player); err != nil {
		return fmt.Errorf("failed to run step %s: %w", step.Name, err)
	}
	player.ManaPool.Empty()
	g.Log(fmt.Sprintf("Step completed: %s", step.Name))
	return nil
}

func (g *GameState) RunPhase(phase GamePhase, player *player.Player) error {
	g.CurrentPhase = phase.Name
	g.Log(fmt.Sprintf("Running phase: %s", phase.Name))
	for _, step := range phase.Steps {
		if err := g.RunStep(step, player); err != nil {
			return fmt.Errorf("failed to run phase %s step %s: %w", phase.Name, step.Name, err)
		}
	}
	g.Log(fmt.Sprintf("Phase completed: %s", phase.Name))
	return nil
}

var ChoiceMulligan = "Mulligan"

// var ChoiceSourceMulligan = NewChoiceSource(ChoiceMulligan, ChoiceMulligan)

// Mulligan allows the player to mulligan their hand. The player draws 7 new
// cards but then needs to put 1 card back on the bottom of their library for
// each time they've mulliganed.
func (g *GameState) Mulligan(player *player.Player) error {
	var accept bool
	for (player.Mulligans < player.MaxHandSize) || !accept {
		// player.Agent.ReportState(g)
		/*
				accept, err := player.Agent.Confirm("Keep Hand? (y/n)", ChoiceSourceMulligan)
				if err != nil {
					return fmt.Errorf("failed to confirm mulligan: %w", err)
				}
				if accept {
					break
				}
				g.Log("Mulliganing...")
				for _, card := range player.Hand.GetAll() {
					player.Hand.Remove(card.ID())
					player.Library.Add(card)
				}
				player.Library.Shuffle()
				g.DrawStartingHand(player)
				player.Mulligans++
			}
			if player.Mulligans != 0 {
				if err := PutNBackOnTop(g, player.Mulligans, NewChoiceSource(ChoiceMulligan, ChoiceMulligan), player); err != nil {
					return fmt.Errorf("failed to put cards back on top: %w", err)
				}
		*/
	}
	return nil
}

// RunTurn executes a single game turn using the provided player.
// Maybe split into ResolvePhase and ResolveStep functions
// RunTurn executes a single game turn using the provided player
func (g *GameState) RunTurn(player *player.Player) error {
	// Beginning Phase
	player.Turn++
	player.LandDrop = false
	// player.Battlefield.RemoveSummoningSickness()
	g.TurnMessageLog = []string{} // TODO: this sucks, make better
	g.Log(fmt.Sprintf("Turn %d", player.Turn))
	if err := g.RunPhase(beginningPhase, player); err != nil {
		return fmt.Errorf("failed to run beginning phase: %w", err)
	}
	if err := g.RunPhase(precombatMainPhase, player); err != nil {
		return fmt.Errorf("failed to run precombat main phase: %w", err)
	}
	if err := g.RunPhase(combatPhase, player); err != nil {
		return fmt.Errorf("failed to run combat phase: %w", err)
	}
	if err := g.RunPhase(postcombatMainPhase, player); err != nil {
		return fmt.Errorf("failed to run postcombat main phase: %w", err)
	}
	if err := g.RunPhase(endingPhase, player); err != nil {
		return fmt.Errorf("failed to run ending phase: %w", err)
	}
	return nil
}

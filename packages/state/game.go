package state

// TODO: Document what kind of logic lives where.

// I think the state should really just control the mechanics of movings
// things around, not checking the game rules. Games rule enforcement should
// happen in the engine

// stuff like can cast might even need to be moved to the engine package

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"errors"
	"fmt"
	"strconv"
)

type Game struct {
	id                    int
	cheatsEnabled         bool
	activePlayerIdx       int
	playerWithPriority    string
	playersPassedPriority map[string]bool
	battlefield           Battlefield
	phase                 mtg.Phase
	step                  mtg.Step
	players               []Player
	stack                 Stack
	winnerID              string
}

func (g Game) GetCardsAvailableToPlay(playerID string) ([]gob.Card, error) {
	player, err := g.GetPlayer(playerID)
	if err != nil {
		return nil, err
	}
	var availableCards []gob.Card
	for _, card := range player.hand.cards {
		if g.CanPlayCard(card, playerID, mtg.ZoneHand) {
			availableCards = append(availableCards, card)
		}
	}
	for _, card := range player.graveyard.cards {
		if g.CanPlayCard(card, playerID, mtg.ZoneGraveyard) {
			availableCards = append(availableCards, card)
		}
	}
	return availableCards, nil
}

// TODO should this live on engine?
func (g Game) CanCastSorcery(playerID string) bool {
	if !g.IsStackEmtpy() {
		return false
	}
	if g.ActivePlayerID() != playerID {
		return false
	}
	if !(g.step == mtg.StepPrecombatMain ||
		g.step == mtg.StepPostcombatMain) {
		return false
	}
	fmt.Println("Can cast sorcery:", g.step, g.ActivePlayerID(), playerID)
	return true
}

// TODO: Should this live on engine?
func (g Game) CanPlayCard(
	card gob.Card,
	playerID string,
	zone mtg.Zone,
) bool {
	if zone == mtg.ZoneGraveyard {
		return false
	}
	fmt.Println("aaa")
	if card.Match(is.Permanent()) {
		fmt.Println("ccc")
	}
	if card.Match(has.CardType(mtg.CardTypeSorcery)) {
		fmt.Println("bbb")
		if !g.CanCastSorcery(playerID) {
			return false
		}
	}
	return true
}

func (g Game) WithCheatsEnabled(enabled bool) Game {
	g.cheatsEnabled = enabled
	return g
}

func (g Game) CheatsEnabled() bool {
	return g.cheatsEnabled
}

func (g Game) WithPhase(phase mtg.Phase) Game {
	g.phase = phase
	return g
}

func (g Game) WithStep(step mtg.Step) Game {
	g.step = step
	return g
}

func (g Game) Phase() mtg.Phase {
	return g.phase
}

func (g Game) Step() mtg.Step {
	return g.step
}

func (g Game) IsStackEmtpy() bool {
	return g.stack.Size() == 0
}

func (g Game) WithGameOver(winnerID string) Game {
	g.winnerID = winnerID
	return g
}

func (g Game) WithPlayers(players []Player) Game {
	g.players = players
	return g
}

func (g Game) WithClearedPriority() Game {
	g.playerWithPriority = ""
	return g
}

func (g Game) WithPlayerWithPriority(playerID string) Game {
	g.playerWithPriority = playerID
	return g
}

func (g Game) WithActivePlayer(playerID string) Game {
	var idx int
	for i, p := range g.players {
		if p.id == playerID {
			idx = i
			break
		}
	}
	g.activePlayerIdx = idx
	return g
}

func (g Game) WithResetPriorityPasses() Game {
	g.playersPassedPriority = map[string]bool{}
	return g
}

func (g Game) WithPlayerPassedPriority(playerID string) Game {
	newPlayersPassedPriority := map[string]bool{}
	for pID := range g.playersPassedPriority {
		newPlayersPassedPriority[pID] = g.playersPassedPriority[pID]
	}
	newPlayersPassedPriority[playerID] = true
	g.playersPassedPriority = newPlayersPassedPriority
	return g
}

func (g Game) WithUpdatedPlayer(player Player) Game {
	var newPlayers []Player
	for _, p := range g.players {
		if p.id == player.id {
			newPlayers = append(newPlayers, player)
			continue
		}
		newPlayers = append(newPlayers, p)
	}
	g.players = newPlayers
	return g
}

func (g Game) WithBattlefield(battlefield Battlefield) Game {
	g.battlefield = battlefield
	return g
}

func (g Game) Battlefield() Battlefield {
	return g.battlefield
}

func (g Game) Stack() Stack {
	return g.stack
}

// Returns the players starting with the active player and going in turn order
func (g Game) PlayerIDsInTurnOrder() []string {
	var n = g.activePlayerIdx
	var playersInTurnOrder []string
	for i := 0; i < len(g.players); i++ {
		playersInTurnOrder = append(playersInTurnOrder, g.players[n].id)
		n = (n + 1) % len(g.players)
	}
	return playersInTurnOrder
}

func (g Game) ActivePlayerID() string {
	return g.players[g.activePlayerIdx].ID()
}

func (g Game) AllPlayersPassedPriority() bool {
	for _, player := range g.players {
		if !g.playersPassedPriority[player.id] {
			return false
		}
	}
	return true
}

// TODO: Think about removing this and using GetPlayerID instead
func (g Game) GetPlayer(id string) (Player, error) {
	fmt.Println("GetPlayer called with ID:", id)
	for _, player := range g.players {
		fmt.Println("Checking player ID:", player.id)
		if player.id == id {
			return player, nil

		}
	}
	return Player{}, errors.New("player not found")
}

// TODO: THIS WILL BREAK WITH MORE THAN 2 PLAYERS
func (g Game) GetOpponent(id string) (Player, error) {
	if len(g.players) > 2 {
		panic("GetOpponent is not implemented for more than 2 players")
	}
	opponentID := g.NextPlayerID(id)
	for _, player := range g.players {
		if player.id == opponentID {
			return player, nil
		}
	}
	return Player{}, errors.New("opponent not found")
}

func (g Game) NextPlayerID(currentPlayerID string) string {
	for i, player := range g.players {
		if player.id == currentPlayerID {
			nextIdx := (i + 1) % len(g.players)
			return g.players[nextIdx].ID()
		}
	}
	return currentPlayerID
}

func (g Game) PriorityPlayerID() string {
	return g.playerWithPriority
}

func (g Game) PlayerPassedPriority(id string) bool {
	return g.playersPassedPriority[id]
}

func NewGame(config GameConfig) Game {
	state := Game{
		players: config.Players,
	}
	return state
}

type GameConfig struct {
	Players []Player
}

func (g Game) IsGameOver() bool {
	if g.winnerID != "" {
		return true
	}
	return false
}

type GameStateSnapshot struct {
	Turn int
}

func (g Game) GetNextID() (id string, game Game) {
	g.id++
	return strconv.Itoa(g.id), g
}

func (g Game) WithPutCardOnBattlefield(card gob.Card, playerID string) (Game, error) {
	id, newGame := g.GetNextID()
	permanent, err := gob.NewPermanent(id, card)
	if err != nil {
		return newGame, err
	}
	newBattlefield := newGame.battlefield.Append(permanent)
	newerGame := newGame.WithBattlefield(newBattlefield)
	return newerGame, nil
}

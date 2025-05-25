package game

import (
	"strings"
)

// Log logs a message to the game state message log.
func (g *GameState) Log(message ...string) {
	// TODO: Something like this so we know which player to which action
	// outMessage := fmt.Sprintf("%s: %s", g.ActivatePlayer.ID )
	// TODO: There's probably a more elegant way to do this
	g.TurnMessageLog = append(g.TurnMessageLog, strings.Join(message, " "))
	g.MessageLog = append(g.MessageLog, strings.Join(message, " "))
}

// Error logs an error message to the game state message log.
// TODO: There's probably a more elegant way to do this
func (g *GameState) Error(err error) {
	g.TurnMessageLog = append(g.TurnMessageLog, "ERROR: "+err.Error())
	g.MessageLog = append(g.MessageLog, "ERROR: "+err.Error())
}

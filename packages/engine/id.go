package engine

import "strconv"

// TODO: Maybe take a prefix that is the source for the ID, e.g. "card name" or
// "player name"
func (g *GameState) GetNextID() string {
	g.nextID++
	return strconv.Itoa(g.nextID)
}

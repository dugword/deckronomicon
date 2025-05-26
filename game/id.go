package game

import "strconv"

var nextId int

// TODO: Maybe take a prefix that is the source for the ID, e.g. "card name" or
// "player name"
func GetNextID() string {
	nextId++
	return strconv.Itoa(nextId)
}

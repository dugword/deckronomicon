package game

import "strconv"

var nextId int

func GetNextID() string {
	nextId++
	return strconv.Itoa(nextId)
}

package query

type Predicate func(obj Object) bool

// TODO add a ToString method to Predicate so we can have a string
// representation of the predicate for debugging purposes

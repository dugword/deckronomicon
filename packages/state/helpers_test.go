package state

import (
	"reflect"

	"github.com/google/go-cmp/cmp"
)

var AllowAllUnexported = cmp.Exporter(func(reflect.Type) bool { return true })

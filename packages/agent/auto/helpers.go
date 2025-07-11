package auto

import (
	"fmt"
	"strconv"
	"strings"
)

/*
func ZoneHasAll(zone engine.Zone, names []string) bool {
	found := map[string]bool{}
	for _, obj := range zone.GetAll() {
		for _, name := range names {
			if obj.Name() == name {
				found[name] = true
			}
		}
	}
	for _, name := range names {
		if !found[name] {
			return false
		}
	}
	return true
}
*/

func parseStatShortcut(stat string, raw string) (string, string, int, error) {
	ops := []string{"<=", ">=", "!=", "==", "<", ">"}
	for _, op := range ops {
		if strings.HasPrefix(raw, op) {
			valStr := strings.TrimPrefix(raw, op)
			value, err := strconv.Atoi(valStr)
			if err != nil {
				return "", "", 0, fmt.Errorf("invalid value for %s: %s", stat, raw)
			}
			return stat, op, value, nil
		}
	}
	return "", "", 0, fmt.Errorf("invalid shortcut format for %s: %s", stat, raw)
}

package auto

import (
	"deckronomicon/game"
	"fmt"
	"strconv"
	"strings"
)

func ZoneHasAll(zone game.Zone, names []string) bool {
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

func expandDefinitions(names []string, definitions map[string][]string) ([]string, error) {
	var expanded []string
	for _, name := range names {
		if !strings.HasPrefix(name, "$") {
			expanded = append(expanded, name)
			continue
		}
		key := strings.TrimPrefix(name, "$")
		values, ok := definitions[key]
		if !ok {
			return nil, fmt.Errorf("definition %s not found", key)
		}
		expanded = append(expanded, values...)
	}
	return expanded, nil
}

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

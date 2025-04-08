package db

import (
	"fmt"
	"math/rand"
	"strings"
)

func ResolveReference(ref string, inserted map[string][]map[string]any) (any, error) {
	parts := strings.Split(ref, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ref format: %s (expected table.field)", ref)
	}
	table, field := parts[0], parts[1]
	rows := inserted[table]
	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows available for reference: %s", ref)
	}
	return rows[rand.Intn(len(rows))][field], nil
}

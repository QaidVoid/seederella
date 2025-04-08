package resolve

import (
	"fmt"
	"strings"

	"github.com/QaidVoid/seederella/internal/config"
	"github.com/QaidVoid/seederella/internal/db"
	"github.com/QaidVoid/seederella/internal/faker"
)

// SortFields returns the order in which fields should be resolved.
// It performs a topological sort based on dependencies (same_as).
// The function returns an ordered slice of field names or an error if
// a cyclic dependency is detected in the field configuration.
func SortFields(fields map[string]config.ColumnConfig) ([]string, error) {
	graph := map[string][]string{}
	for name, field := range fields {
		if field.SameAs != "" {
			graph[name] = append(graph[name], field.SameAs)
		} else {
			graph[name] = []string{}
		}
	}

	visited := map[string]bool{}
	temp := map[string]bool{}
	var result []string

	var visit func(string) error
	visit = func(node string) error {
		if temp[node] {
			return fmt.Errorf("cyclic dependency detected at %s", node)
		}
		if !visited[node] {
			temp[node] = true
			for _, dep := range graph[node] {
				if err := visit(dep); err != nil {
					return err
				}
			}
			visited[node] = true
			temp[node] = false
			result = append(result, node)
		}
		return nil
	}

	for name := range graph {
		if err := visit(name); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// ResolveField generates a value for a field based on its configuration.
// It supports various value generation strategies including direct values,
// same_as references, database references, faker functions, and copy_from.
// It can also apply transformations to the generated values.
//
// Parameters:
//   - fieldName: Name of the field being resolved
//   - fieldCfg: Configuration for the field
//   - rowData: Map of field names to values for the current row
//   - inserted: Map of table names to previously inserted rows
//
// Returns the resolved value or an error if resolution fails.
func ResolveField(
	fieldName string,
	fieldCfg config.ColumnConfig,
	rowData map[string]any,
	inserted map[string][]map[string]any,
) (any, error) {
	var val any
	var err error

	if fieldCfg.Value != nil {
		val = fieldCfg.Value
	} else if fieldCfg.SameAs != "" {
		val, err = resolveSameAs(fieldCfg.SameAs, rowData)
		if err != nil {
			return nil, fmt.Errorf("field %s references same_as: %s, but that field has no value: %w", fieldName, fieldCfg.SameAs, err)
		}
	} else if fieldCfg.Reference != "" {
		val, err = db.ResolveReference(fieldCfg.Reference, inserted)
		if err != nil {
			return nil, fmt.Errorf("field %s: failed to resolve reference: %w", fieldName, err)
		}
	} else if fieldCfg.Faker != "" {
		val, err = faker.Generate(fieldCfg.Faker)
		if err != nil {
			return nil, fmt.Errorf("field %s: failed to generate faker: %w", fieldName, err)
		}
	} else {
		return nil, fmt.Errorf("field %s has no generation strategy defined", fieldName)
	}

	if fieldCfg.Transform != "" {
		val, err = applyTransform(val, fieldCfg.Transform)
		if err != nil {
			return nil, fmt.Errorf("field %s: failed to apply transform %s: %w", fieldName, fieldCfg.Transform, err)
		}
	}

	return val, nil
}

func applyTransform(val any, transform string) (any, error) {
	switch transform {
	case "lower":
		return strings.ToLower(fmt.Sprintf("%v", val)), nil
	case "upper":
		return strings.ToUpper(fmt.Sprintf("%v", val)), nil
	default:
		return val, nil
	}
}

func resolveSameAs(sameAsField string, rowData map[string]any) (any, error) {
	val, ok := rowData[sameAsField]
	if !ok {
		return nil, fmt.Errorf("same_as field '%s' not found", sameAsField)
	}
	return val, nil
}

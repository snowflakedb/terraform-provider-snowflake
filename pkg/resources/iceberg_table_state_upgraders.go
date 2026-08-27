package resources

import "context"

// v2_20_0_IcebergTableConstraintColumnFieldsStateUpgrader renames the singular constraint
// column list attributes used through v2.20.0 to the plural names aligned with hybrid tables:
// column → columns, ref_column → ref_columns.
func v2_20_0_IcebergTableConstraintColumnFieldsStateUpgrader(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}
	renameConstraintListColumnFields(rawState, "primary_key_constraint")
	renameConstraintListColumnFields(rawState, "unique_constraint")
	renameConstraintListColumnFields(rawState, "foreign_key_constraint")
	return rawState, nil
}

func renameConstraintListColumnFields(rawState map[string]any, attr string) {
	items, ok := rawState[attr].([]any)
	if !ok {
		return
	}
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		renameMapKey(m, "column", "columns")
		renameMapKey(m, "ref_column", "ref_columns")
	}
}

func renameMapKey(m map[string]any, from, to string) {
	v, ok := m[from]
	if !ok {
		return
	}
	if _, exists := m[to]; !exists {
		m[to] = v
	}
	delete(m, from)
}

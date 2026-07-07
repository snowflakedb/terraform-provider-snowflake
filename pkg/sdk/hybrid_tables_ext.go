package sdk

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
)

func (r hybridTableDetailsRow) additionalConvert(result *HybridTableDetails) error {
	type_, collation := r.splitTypeAndCollation()
	result.Type = type_
	result.Collation = collation
	return nil
}

// ShowParameters returns the parameters visible at the TABLE level for the given hybrid table.
// Mirrors pkg/sdk/functions_ext.go:155 (ParametersIn.Function) with ParametersIn.Table.
func (v *hybridTables) ShowParameters(ctx context.Context, id SchemaObjectIdentifier) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Table: id,
		},
	})
}

// splitTypeAndCollation extracts the COLLATE 'X' suffix that DESCRIBE TABLE returns for
// collated columns (e.g. "VARCHAR(200) COLLATE 'en-ci'") and returns the bare type and
// the collation as separate values. Mirrors pkg/sdk/tables.go:736 — the same shape exists
// for classic tables but the generator cannot derive a Collation field from the raw Type
// column, so this lives in the _ext.go file and is invoked manually from convert().
func (r hybridTableDetailsRow) splitTypeAndCollation() (string, *string) {
	collateRegexp := regexp.MustCompile(`COLLATE +'([a-zA-Z0-9_-]*)'`)
	matches := collateRegexp.FindStringSubmatch(r.Type)
	if len(matches) == 2 {
		collation := matches[1]
		type_ := strings.TrimSpace(collateRegexp.ReplaceAllString(r.Type, ""))
		return type_, &collation
	}
	return r.Type, nil
}

// HybridTableConstraint is the merged, ordered view of one table constraint.
type HybridTableConstraint struct {
	Name    string
	Kind    ColumnConstraintType
	Columns []string // ordered by key_sequence

	// FK-only; zero-value for PK/UNIQUE.
	ReferencedTable   SchemaObjectIdentifier
	ReferencedColumns []string // ordered by key_sequence
	DeleteRule        string
	UpdateRule        string
}

// keyRow is the shared in-memory shape for PK and UNIQUE rows. The generator emits
// TablePrimaryKey and TableUniqueKey as two distinct types (one per StructPair), so
// both are field-copied into keyRow before being merged by the single mergeKeyRows
// function, which holds the non-trivial sort/group logic.
type keyRow struct {
	ConstraintName string
	ColumnName     string
	KeySequence    int
}

func primaryKeysToKeyRows(in []TablePrimaryKey) []keyRow {
	out := make([]keyRow, len(in))
	for i, r := range in {
		out[i] = keyRow{ConstraintName: r.ConstraintName, ColumnName: r.ColumnName, KeySequence: r.KeySequence}
	}
	return out
}

func uniqueKeysToKeyRows(in []TableUniqueKey) []keyRow {
	out := make([]keyRow, len(in))
	for i, r := range in {
		out[i] = keyRow{ConstraintName: r.ConstraintName, ColumnName: r.ColumnName, KeySequence: r.KeySequence}
	}
	return out
}

// mergeKeyRows groups PK/UNIQUE rows by constraint_name, preserving discovery order
// of the groups and key_sequence order of columns within each group. A table has at
// most one PRIMARY KEY, so a PK call returns 0 or 1 element.
func mergeKeyRows(rows []keyRow, kind ColumnConstraintType) []HybridTableConstraint {
	if len(rows) == 0 {
		return nil
	}
	var order []string
	byName := make(map[string][]keyRow)
	for _, r := range rows {
		if _, seen := byName[r.ConstraintName]; !seen {
			order = append(order, r.ConstraintName)
		}
		byName[r.ConstraintName] = append(byName[r.ConstraintName], r)
	}
	out := make([]HybridTableConstraint, 0, len(order))
	for _, name := range order {
		grp := byName[name]
		sort.SliceStable(grp, func(i, j int) bool { return grp[i].KeySequence < grp[j].KeySequence })
		c := HybridTableConstraint{Name: name, Kind: kind}
		for _, r := range grp {
			c.Columns = append(c.Columns, r.ColumnName)
		}
		out = append(out, c)
	}
	return out
}

// mergeForeignKeyRows groups imported-key rows by fk_name. Each multi-column FK
// produces one row per column; columns and referenced columns are ordered by key_sequence.
func mergeForeignKeyRows(rows []TableImportedKey) []HybridTableConstraint {
	if len(rows) == 0 {
		return nil
	}
	var order []string
	byName := make(map[string][]TableImportedKey)
	for _, r := range rows {
		if _, seen := byName[r.FkName]; !seen {
			order = append(order, r.FkName)
		}
		byName[r.FkName] = append(byName[r.FkName], r)
	}
	out := make([]HybridTableConstraint, 0, len(order))
	for _, name := range order {
		grp := byName[name]
		sort.SliceStable(grp, func(i, j int) bool { return grp[i].KeySequence < grp[j].KeySequence })
		// ReferencedTable, DeleteRule, and UpdateRule are identical across all rows of a
		// given FK, so we take them from the first row.
		first := grp[0]
		c := HybridTableConstraint{
			Name:            name,
			Kind:            ColumnConstraintTypeForeignKey,
			ReferencedTable: NewSchemaObjectIdentifier(first.PkDatabaseName, first.PkSchemaName, first.PkTableName),
			DeleteRule:      first.DeleteRule,
			UpdateRule:      first.UpdateRule,
		}
		for _, r := range grp {
			c.Columns = append(c.Columns, r.FkColumnName)
			c.ReferencedColumns = append(c.ReferencedColumns, r.PkColumnName)
		}
		out = append(out, c)
	}
	return out
}

func (v *hybridTables) GetConstraints(ctx context.Context, id SchemaObjectIdentifier) ([]HybridTableConstraint, error) {
	var result []HybridTableConstraint

	pkRows, err := v.ShowPrimaryKeys(ctx, NewShowPrimaryKeysHybridTableRequest(id))
	if err != nil {
		return nil, fmt.Errorf("showing primary keys for %s: %w", id.FullyQualifiedName(), err)
	}
	result = append(result, mergeKeyRows(primaryKeysToKeyRows(pkRows), ColumnConstraintTypePrimaryKey)...)

	ukRows, err := v.ShowUniqueKeys(ctx, NewShowUniqueKeysHybridTableRequest(id))
	if err != nil {
		return nil, fmt.Errorf("showing unique keys for %s: %w", id.FullyQualifiedName(), err)
	}
	result = append(result, mergeKeyRows(uniqueKeysToKeyRows(ukRows), ColumnConstraintTypeUnique)...)

	fkRows, err := v.ShowImportedKeys(ctx, NewShowImportedKeysHybridTableRequest(id))
	if err != nil {
		log.Printf("[WARN] SHOW IMPORTED KEYS failed for %s (undocumented for hybrid tables); skipping FK read-back: %v", id.FullyQualifiedName(), err)
		return result, nil
	}
	result = append(result, mergeForeignKeyRows(fkRows)...)

	return result, nil
}

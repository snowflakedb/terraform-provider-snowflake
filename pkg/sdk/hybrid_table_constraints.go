package sdk

import (
	"context"
	"fmt"
	"log"
	"sort"
)

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

// Partial row structs — safe because the SDK db handle is in Unsafe() mode
// (pkg/sdk/client.go:141), so unmapped SHOW columns are ignored.
type showKeyRow struct {
	ConstraintName string `db:"constraint_name"`
	ColumnName     string `db:"column_name"`
	KeySequence    int    `db:"key_sequence"`
}

type showImportedKeyRow struct {
	FkName         string `db:"fk_name"`
	FkColumnName   string `db:"fk_column_name"`
	KeySequence    int    `db:"key_sequence"`
	PkDatabaseName string `db:"pk_database_name"`
	PkSchemaName   string `db:"pk_schema_name"`
	PkTableName    string `db:"pk_table_name"`
	PkColumnName   string `db:"pk_column_name"`
	DeleteRule     string `db:"delete_rule"`
	UpdateRule     string `db:"update_rule"`
}

func showPrimaryKeysSQL(id SchemaObjectIdentifier) string {
	return fmt.Sprintf(`SHOW PRIMARY KEYS IN TABLE %s`, id.FullyQualifiedName())
}

func showUniqueKeysSQL(id SchemaObjectIdentifier) string {
	return fmt.Sprintf(`SHOW UNIQUE KEYS IN TABLE %s`, id.FullyQualifiedName())
}

func showImportedKeysSQL(id SchemaObjectIdentifier) string {
	return fmt.Sprintf(`SHOW IMPORTED KEYS IN TABLE %s`, id.FullyQualifiedName())
}

// GetConstraints reads PK, UNIQUE, and FK constraints for a hybrid table via three
// metadata-only SHOW commands. None require a warehouse (verified preprod).
// SHOW IMPORTED KEYS is undocumented for hybrid tables, so its failure is non-fatal:
// we log a warning and return the PK/UNIQUE results without FKs.
func (v *hybridTables) GetConstraints(ctx context.Context, id SchemaObjectIdentifier) ([]HybridTableConstraint, error) {
	var result []HybridTableConstraint

	var pkRows []showKeyRow
	if err := v.client.query(ctx, &pkRows, showPrimaryKeysSQL(id)); err != nil {
		return nil, fmt.Errorf("showing primary keys for %s: %w", id.FullyQualifiedName(), err)
	}
	result = append(result, mergeKeyRows(pkRows, ColumnConstraintTypePrimaryKey)...)

	var ukRows []showKeyRow
	if err := v.client.query(ctx, &ukRows, showUniqueKeysSQL(id)); err != nil {
		return nil, fmt.Errorf("showing unique keys for %s: %w", id.FullyQualifiedName(), err)
	}
	result = append(result, mergeKeyRows(ukRows, ColumnConstraintTypeUnique)...)

	var fkRows []showImportedKeyRow
	if err := v.client.query(ctx, &fkRows, showImportedKeysSQL(id)); err != nil {
		log.Printf("[WARN] SHOW IMPORTED KEYS failed for %s (undocumented for hybrid tables); skipping FK read-back: %v", id.FullyQualifiedName(), err)
		return result, nil
	}
	result = append(result, mergeForeignKeyRows(fkRows)...)

	return result, nil
}

// mergeKeyRows groups PK/UNIQUE rows by constraint_name, preserving discovery order
// of the groups and key_sequence order of columns within each group. A table has at
// most one PRIMARY KEY, so a PK call returns 0 or 1 element.
func mergeKeyRows(rows []showKeyRow, kind ColumnConstraintType) []HybridTableConstraint {
	if len(rows) == 0 {
		return nil
	}
	var order []string
	byName := make(map[string][]showKeyRow)
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
func mergeForeignKeyRows(rows []showImportedKeyRow) []HybridTableConstraint {
	if len(rows) == 0 {
		return nil
	}
	var order []string
	byName := make(map[string][]showImportedKeyRow)
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

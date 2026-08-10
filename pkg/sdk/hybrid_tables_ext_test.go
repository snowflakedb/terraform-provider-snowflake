package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHybridTableDetailsRow_SplitTypeAndCollation(t *testing.T) {
	testCases := []struct {
		Name              string
		Value             string
		ExpectedType      string
		ExpectedCollation *string
	}{
		{
			Name:              "with utf8",
			Value:             "VARCHAR(10) COLLATE 'utf8'",
			ExpectedType:      "VARCHAR(10)",
			ExpectedCollation: new("utf8"),
		},
		{
			Name:              "with locale",
			Value:             "VARCHAR(10) COLLATE 'en_US'",
			ExpectedType:      "VARCHAR(10)",
			ExpectedCollation: new("en_US"),
		},
		{
			Name:              "with multiple specifiers",
			Value:             "VARCHAR(10) COLLATE 'fr_CA-ai-pi-trim'",
			ExpectedType:      "VARCHAR(10)",
			ExpectedCollation: new("fr_CA-ai-pi-trim"),
		},
		{
			Name:              "with empty collation",
			Value:             "VARCHAR(10) COLLATE ''",
			ExpectedType:      "VARCHAR(10)",
			ExpectedCollation: new(""),
		},
		{
			Name:              "without collation",
			Value:             "NUMBER(38, 0)",
			ExpectedType:      "NUMBER(38, 0)",
			ExpectedCollation: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			row := hybridTableDetailsRow{Type: tc.Value}
			actualType, actualCollation := row.splitTypeAndCollation()
			assert.Equal(t, tc.ExpectedType, actualType)
			if tc.ExpectedCollation == nil {
				assert.Nil(t, actualCollation)
			} else {
				assert.Equal(t, *tc.ExpectedCollation, *actualCollation)
			}
		})
	}
}

func TestHybridTableConstraints_mergeKeyRows(t *testing.T) {
	testCases := []struct {
		name     string
		rows     []keyRow
		kind     ColumnConstraintType
		expected []HybridTableConstraint
	}{
		{
			name:     "empty input returns nil",
			rows:     nil,
			kind:     ColumnConstraintTypePrimaryKey,
			expected: nil,
		},
		{
			name: "single-column primary key",
			rows: []keyRow{
				{ConstraintName: "PK_T", ColumnName: "ID", KeySequence: 1},
			},
			kind: ColumnConstraintTypePrimaryKey,
			expected: []HybridTableConstraint{
				{Name: "PK_T", Kind: ColumnConstraintTypePrimaryKey, Columns: []string{"ID"}},
			},
		},
		{
			name: "two distinct unique constraints, columns ordered by key_sequence",
			rows: []keyRow{
				// Intentionally out of key_sequence order to prove sorting.
				{ConstraintName: "UQ_A", ColumnName: "COL_A2", KeySequence: 2},
				{ConstraintName: "UQ_A", ColumnName: "COL_A1", KeySequence: 1},
				{ConstraintName: "UQ_B", ColumnName: "COL_B", KeySequence: 1},
			},
			kind: ColumnConstraintTypeUnique,
			expected: []HybridTableConstraint{
				{Name: "UQ_A", Kind: ColumnConstraintTypeUnique, Columns: []string{"COL_A1", "COL_A2"}},
				{Name: "UQ_B", Kind: ColumnConstraintTypeUnique, Columns: []string{"COL_B"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, mergeKeyRows(tc.rows, tc.kind))
		})
	}
}

func TestHybridTableConstraints_mergeForeignKeyRows(t *testing.T) {
	t.Run("empty input returns nil", func(t *testing.T) {
		assert.Nil(t, mergeForeignKeyRows(nil))
	})

	t.Run("multi-column foreign key, columns and referenced columns ordered by key_sequence", func(t *testing.T) {
		// Two rows of the same FK, intentionally out of key_sequence order.
		rows := []TableImportedKey{
			{
				FkName:         "FK_T",
				FkColumnName:   "PARENT_B",
				KeySequence:    2,
				PkDatabaseName: "DB",
				PkSchemaName:   "SCH",
				PkTableName:    "PARENT",
				PkColumnName:   "B",
				DeleteRule:     "NO ACTION",
				UpdateRule:     "NO ACTION",
			},
			{
				FkName:         "FK_T",
				FkColumnName:   "PARENT_A",
				KeySequence:    1,
				PkDatabaseName: "DB",
				PkSchemaName:   "SCH",
				PkTableName:    "PARENT",
				PkColumnName:   "A",
				DeleteRule:     "NO ACTION",
				UpdateRule:     "NO ACTION",
			},
		}

		expected := []HybridTableConstraint{
			{
				Name:              "FK_T",
				Kind:              ColumnConstraintTypeForeignKey,
				Columns:           []string{"PARENT_A", "PARENT_B"},
				ReferencedTable:   NewSchemaObjectIdentifier("DB", "SCH", "PARENT"),
				ReferencedColumns: []string{"A", "B"},
				DeleteRule:        "NO ACTION",
				UpdateRule:        "NO ACTION",
			},
		}

		assert.Equal(t, expected, mergeForeignKeyRows(rows))
	})
}

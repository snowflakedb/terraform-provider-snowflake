package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_buildHybridAlterColumnTypeAction(t *testing.T) {
	t.Run("valid type produces request with parsed type", func(t *testing.T) {
		col := hybridTableColumn{
			name:     "AGE",
			dataType: "NUMBER(20,5)",
		}

		req, err := buildHybridAlterColumnTypeAction(col)
		require.NoError(t, err)
		require.NotNil(t, req)

		assert.Equal(t, "AGE", req.ColumnName)
		require.NotNil(t, req.DataType)
		assert.Equal(t, sdk.DataType("NUMBER(20, 5)"), *req.DataType)
	})

	t.Run("invalid type returns error mentioning column name", func(t *testing.T) {
		col := hybridTableColumn{
			name:     "WAT",
			dataType: "NOT_A_REAL_TYPE",
		}

		req, err := buildHybridAlterColumnTypeAction(col)
		require.Error(t, err)
		assert.Nil(t, req)
		assert.Contains(t, err.Error(), "WAT")
	})

	t.Run("empty type returns error", func(t *testing.T) {
		col := hybridTableColumn{
			name:     "X",
			dataType: "",
		}

		req, err := buildHybridAlterColumnTypeAction(col)
		require.Error(t, err)
		assert.Nil(t, req)
	})
}

func Test_buildHybridColumnSpec(t *testing.T) {
	t.Run("plain column with no default", func(t *testing.T) {
		col := hybridTableColumn{
			name:     "ID",
			dataType: "NUMBER(38,0)",
			nullable: true,
		}

		spec, err := buildHybridColumnSpec(col)
		require.NoError(t, err)

		require.NotNil(t, spec.dataType)
		assert.Equal(t, "NUMBER(38,0)", spec.dataType.Canonical())
		assert.Nil(t, spec.defaultValue)
		assert.Empty(t, spec.collate)
		assert.Empty(t, spec.comment)
	})

	t.Run("column with constant default and collate", func(t *testing.T) {
		def := "hello"
		col := hybridTableColumn{
			name:     "NAME",
			dataType: "VARCHAR(200)",
			nullable: true,
			_default: &columnDefault{constant: &def},
			collate:  "en-ci",
			comment:  "the name",
		}

		spec, err := buildHybridColumnSpec(col)
		require.NoError(t, err)

		require.NotNil(t, spec.defaultValue)
		require.NotNil(t, spec.defaultValue.Expression)
		assert.Equal(t, "'hello'", *spec.defaultValue.Expression)
		assert.Equal(t, "en-ci", spec.collate)
		assert.Equal(t, "the name", spec.comment)
	})

	t.Run("invalid data type returns error", func(t *testing.T) {
		col := hybridTableColumn{
			name:     "BAD",
			dataType: "NOT_A_REAL_TYPE",
		}

		_, err := buildHybridColumnSpec(col)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "BAD")
	})

	t.Run("default block with multiple fields set returns error", func(t *testing.T) {
		c := "x"
		e := "f()"
		col := hybridTableColumn{
			name:     "X",
			dataType: "VARCHAR(10)",
			nullable: true,
			_default: &columnDefault{constant: &c, expression: &e},
		}
		_, err := buildHybridColumnSpec(col)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "exactly one")
	})
}

func TestUnit_HybridTable_uniqueConstraintHash_ignoresName(t *testing.T) {
	a := map[string]any{"name": "", "columns": []any{"col_a"}}
	b := map[string]any{"name": "SYS_CONSTRAINT_abc", "columns": []any{"col_a"}}
	require.Equal(t, uniqueConstraintHash(a), uniqueConstraintHash(b))

	c := map[string]any{"name": "", "columns": []any{"col_b"}}
	require.NotEqual(t, uniqueConstraintHash(a), uniqueConstraintHash(c))
}

func TestUnit_HybridTable_foreignKeyHash_ignoresName_and_normalizesTableId(t *testing.T) {
	quoted := map[string]any{
		"name":    "",
		"columns": []any{"parent_id"},
		"references": []any{map[string]any{
			"table_id": `"DB"."SCH"."PARENT"`,
			"columns":  []any{"id"},
		}},
	}
	unquoted := map[string]any{
		"name":    "FK_NAMED",
		"columns": []any{"parent_id"},
		"references": []any{map[string]any{
			"table_id": "DB.SCH.PARENT",
			"columns":  []any{"id"},
		}},
	}
	require.Equal(t, foreignKeyHash(quoted), foreignKeyHash(unquoted))

	differentTable := map[string]any{
		"name":    "",
		"columns": []any{"parent_id"},
		"references": []any{map[string]any{
			"table_id": "DB.SCH.OTHER_PARENT",
			"columns":  []any{"id"},
		}},
	}
	require.NotEqual(t, foreignKeyHash(quoted), foreignKeyHash(differentTable))
}

// TestUnit_HybridTable_indexHash_caseInsensitive proves the load-bearing mechanism
// for index case-suppression: a lowercase config element and its uppercase
// SHOW INDEXES read-back must hash identically, so the TypeSet element identity is
// stable and Read does not produce a spurious "remove old / add new" -> ForceNew.
// This matters because set membership is resolved via this hash BEFORE the leaf
// ignoreCaseSuppressFunc runs, so the suppress func alone cannot prevent the churn.
func TestUnit_HybridTable_indexHash_caseInsensitive(t *testing.T) {
	lower := map[string]any{
		"name":            "idx_x",
		"columns":         []any{"status"},
		"include_columns": schema.NewSet(schema.HashString, []any{"score"}),
	}
	upper := map[string]any{
		"name":            "idx_x",
		"columns":         []any{"STATUS"},
		"include_columns": schema.NewSet(schema.HashString, []any{"SCORE"}),
	}
	require.Equal(t, indexHash(lower), indexHash(upper), "lowercase config and uppercase read-back must hash equal")

	// Genuinely different key columns must hash differently.
	different := map[string]any{
		"name":            "idx_x",
		"columns":         []any{"region"},
		"include_columns": schema.NewSet(schema.HashString, []any{"score"}),
	}
	require.NotEqual(t, indexHash(lower), indexHash(different))

	// include_columns is a TypeSet: declared order must not change the hash.
	orderA := map[string]any{
		"name":            "idx_y",
		"columns":         []any{"a"},
		"include_columns": schema.NewSet(schema.HashString, []any{"x", "y"}),
	}
	orderB := map[string]any{
		"name":            "idx_y",
		"columns":         []any{"a"},
		"include_columns": schema.NewSet(schema.HashString, []any{"y", "x"}),
	}
	require.Equal(t, indexHash(orderA), indexHash(orderB))

	// An index with no include_columns must hash stably (no panic on absent set).
	noInclude := map[string]any{
		"name":    "idx_z",
		"columns": []any{"a"},
	}
	require.NotPanics(t, func() { indexHash(noInclude) })
}

func Test_parseIndexColumnsString(t *testing.T) {
	t.Run("single column", func(t *testing.T) {
		require.Equal(t, []string{"STATUS"}, parseIndexColumnsString("[STATUS]"))
	})
	t.Run("multiple columns split on comma-space", func(t *testing.T) {
		require.Equal(t, []string{"SUBMISSION_ID", "SUITE_NAME", "SUT_TYPE"}, parseIndexColumnsString("[SUBMISSION_ID, SUITE_NAME, SUT_TYPE]"))
	})
	t.Run("empty bracket yields empty slice", func(t *testing.T) {
		require.Empty(t, parseIndexColumnsString("[]"))
	})
	t.Run("no leading space on later elements (not a bare-comma split)", func(t *testing.T) {
		got := parseIndexColumnsString("[A, B]")
		require.Equal(t, []string{"A", "B"}, got)
		require.NotContains(t, got, " B")
	})
}

func Test_buildIndexesStateFromShowIndexes(t *testing.T) {
	indexes := []sdk.HybridTableIndex{
		{Name: "SYS_INDEX_T_PRIMARY", IsUnique: sdk.Bool(true), Columns: sdk.String("[ID]"), IncludedColumns: "[]"},
		{Name: "UQ_EMAIL", IsUnique: sdk.Bool(true), Columns: sdk.String("[EMAIL]"), IncludedColumns: "[]"},
		{Name: "IDX_STATUS", IsUnique: sdk.Bool(false), Columns: sdk.String("[STATUS]"), IncludedColumns: "[]"},
		{Name: "IDX_REGION_INC", IsUnique: sdk.Bool(false), Columns: sdk.String("[REGION]"), IncludedColumns: "[SCORE]"},
	}

	got := buildIndexesStateFromShowIndexes(indexes)

	require.Len(t, got, 2, "only user secondary indexes (is_unique=false) are kept")

	byName := map[string]map[string]any{}
	for _, m := range got {
		byName[m["name"].(string)] = m
	}

	status, ok := byName["IDX_STATUS"]
	require.True(t, ok)
	require.Equal(t, []string{"STATUS"}, status["columns"])
	require.Empty(t, status["include_columns"])

	regionInc, ok := byName["IDX_REGION_INC"]
	require.True(t, ok)
	require.Equal(t, []string{"REGION"}, regionInc["columns"])
	require.Equal(t, []string{"SCORE"}, regionInc["include_columns"])
}

func Test_buildIndexesStateFromShowIndexes_skipsNilDiscriminator(t *testing.T) {
	indexes := []sdk.HybridTableIndex{
		{Name: "UNKNOWN", IsUnique: nil, Columns: sdk.String("[X]"), IncludedColumns: "[]"},
	}
	require.Empty(t, buildIndexesStateFromShowIndexes(indexes), "nil IsUnique is excluded (cannot confirm it is a user secondary index)")
}

func Test_buildIndexesStateFromShowIndexes_nilColumns(t *testing.T) {
	indexes := []sdk.HybridTableIndex{
		{Name: "IDX_NO_COL", IsUnique: sdk.Bool(false), Columns: nil, IncludedColumns: "[]"},
	}
	got := buildIndexesStateFromShowIndexes(indexes)
	require.Len(t, got, 1)
	require.Empty(t, got[0]["columns"])
}

func Test_hybridTableColumnTypeStateFunc(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"VARCHAR(256)", "VARCHAR(256)"},
		{"NUMBER(38,0)", "NUMBER(38, 0)"},
		{"TEXT", "TEXT(16777216)"},
		{"BOOLEAN", "BOOLEAN"},
		{"VARCHAR", "VARCHAR(16777216)"},
		{"NUMBER", "NUMBER(38, 0)"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := DataTypeStateFunc(tc.input)
			require.Equal(t, tc.want, got)
			// Idempotency: applying StateFunc to its own output must be a no-op.
			require.Equal(t, got, DataTypeStateFunc(got))
		})
	}
}

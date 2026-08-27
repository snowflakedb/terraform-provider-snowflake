package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestV2_20_0_IcebergTableConstraintColumnFieldsStateUpgrader(t *testing.T) {
	t.Run("renames singular column list fields on all constraint kinds", func(t *testing.T) {
		rawState := map[string]any{
			"primary_key_constraint": []any{
				map[string]any{"name": "PK", "column": []any{"ID"}},
			},
			"unique_constraint": []any{
				map[string]any{"column": []any{"NAME"}},
			},
			"foreign_key_constraint": []any{
				map[string]any{
					"column":     []any{"REF_ID"},
					"table_name": `"DB"."SCHEMA"."OTHER"`,
					"ref_column": []any{"id"},
				},
			},
		}

		result, err := v2_20_0_IcebergTableConstraintColumnFieldsStateUpgrader(context.Background(), rawState, nil)
		require.NoError(t, err)

		pk := result["primary_key_constraint"].([]any)[0].(map[string]any)
		assert.Equal(t, []any{"ID"}, pk["columns"])
		assert.NotContains(t, pk, "column")
		assert.Equal(t, "PK", pk["name"])

		uq := result["unique_constraint"].([]any)[0].(map[string]any)
		assert.Equal(t, []any{"NAME"}, uq["columns"])
		assert.NotContains(t, uq, "column")

		fk := result["foreign_key_constraint"].([]any)[0].(map[string]any)
		assert.Equal(t, []any{"REF_ID"}, fk["columns"])
		assert.Equal(t, []any{"id"}, fk["ref_columns"])
		assert.NotContains(t, fk, "column")
		assert.NotContains(t, fk, "ref_column")
		assert.Equal(t, `"DB"."SCHEMA"."OTHER"`, fk["table_name"])
	})

	t.Run("nil state is a no-op", func(t *testing.T) {
		result, err := v2_20_0_IcebergTableConstraintColumnFieldsStateUpgrader(context.Background(), nil, nil)
		require.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("missing constraint attributes are a no-op", func(t *testing.T) {
		rawState := map[string]any{"name": "T"}
		result, err := v2_20_0_IcebergTableConstraintColumnFieldsStateUpgrader(context.Background(), rawState, nil)
		require.NoError(t, err)
		assert.Equal(t, "T", result["name"])
		assert.NotContains(t, result, "primary_key_constraint")
	})
}

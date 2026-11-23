package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_migrateTableIdAndNewBooleanStringFields(t *testing.T) {
	t.Run("rewrites the id and sets the new fields to default", func(t *testing.T) {
		rawState := map[string]any{"id": "db|schema|table", "comment": "c"}

		got, err := migrateTableIdAndNewBooleanStringFields(context.Background(), rawState, nil)

		require.NoError(t, err)
		assert.Equal(t, `"db"."schema"."table"`, got["id"])
		assert.Equal(t, "c", got["comment"])
		for _, key := range []string{"is_transient", "enable_schema_evolution", "error_logging", "row_timestamp"} {
			assert.Equal(t, BooleanDefault, got[key], key)
		}
	})

	t.Run("nil state", func(t *testing.T) {
		got, err := migrateTableIdAndNewBooleanStringFields(context.Background(), nil, nil)

		require.NoError(t, err)
		assert.Nil(t, got)
	})
}

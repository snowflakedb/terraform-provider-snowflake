package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_v2_23_0_ShareStateUpgrader(t *testing.T) {
	t.Run("sets secure_objects_only to the default when missing and quotes the id", func(t *testing.T) {
		result, err := v2_23_0_ShareStateUpgrader(context.Background(), map[string]any{"id": "SHARE", "name": "SHARE"}, nil)
		require.NoError(t, err)
		assert.Equal(t, map[string]any{"id": `"SHARE"`, "name": "SHARE", "secure_objects_only": BooleanDefault}, result)
	})

	t.Run("keeps an existing secure_objects_only and an already quoted id", func(t *testing.T) {
		result, err := v2_23_0_ShareStateUpgrader(context.Background(), map[string]any{"id": `"SHARE"`, "secure_objects_only": "false"}, nil)
		require.NoError(t, err)
		assert.Equal(t, map[string]any{"id": `"SHARE"`, "secure_objects_only": "false"}, result)
	})

	t.Run("quotes a name with special characters", func(t *testing.T) {
		result, err := v2_23_0_ShareStateUpgrader(context.Background(), map[string]any{"id": "my share"}, nil)
		require.NoError(t, err)
		assert.Equal(t, `"my share"`, result["id"])
	})

	t.Run("nil state", func(t *testing.T) {
		result, err := v2_23_0_ShareStateUpgrader(context.Background(), nil, nil)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_v2_23_0_IsAgenticSecurityIntegrationStateUpgrader(t *testing.T) {
	t.Run("sets the sentinel when the attribute is missing", func(t *testing.T) {
		rawState := map[string]any{"name": "integration"}

		result, err := v2_23_0_IsAgenticSecurityIntegrationStateUpgrader(context.Background(), rawState, nil)
		require.NoError(t, err)
		assert.Equal(t, BooleanDefault, result["is_agentic"])
		assert.Equal(t, "integration", result["name"])
	})

	t.Run("sets the sentinel when the attribute is empty", func(t *testing.T) {
		rawState := map[string]any{"is_agentic": ""}

		result, err := v2_23_0_IsAgenticSecurityIntegrationStateUpgrader(context.Background(), rawState, nil)
		require.NoError(t, err)
		assert.Equal(t, BooleanDefault, result["is_agentic"])
	})

	t.Run("keeps a value already in state", func(t *testing.T) {
		rawState := map[string]any{"is_agentic": BooleanTrue}

		result, err := v2_23_0_IsAgenticSecurityIntegrationStateUpgrader(context.Background(), rawState, nil)
		require.NoError(t, err)
		assert.Equal(t, BooleanTrue, result["is_agentic"])
	})

	t.Run("nil state is a no-op", func(t *testing.T) {
		result, err := v2_23_0_IsAgenticSecurityIntegrationStateUpgrader(context.Background(), nil, nil)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

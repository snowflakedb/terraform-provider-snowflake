package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_v2_23_0_OpenflowDeploymentSnowflakeManagedUsePrivateLinkUpgrader(t *testing.T) {
	t.Run("sets the sentinel when the attribute is missing", func(t *testing.T) {
		rawState := map[string]any{"name": "deployment"}

		result, err := v2_23_0_OpenflowDeploymentSnowflakeManagedUsePrivateLinkUpgrader(context.Background(), rawState, nil)
		require.NoError(t, err)
		assert.Equal(t, BooleanDefault, result["use_private_link"])
		assert.Equal(t, "deployment", result["name"])
	})

	t.Run("sets the sentinel when the attribute is empty", func(t *testing.T) {
		rawState := map[string]any{"use_private_link": ""}

		result, err := v2_23_0_OpenflowDeploymentSnowflakeManagedUsePrivateLinkUpgrader(context.Background(), rawState, nil)
		require.NoError(t, err)
		assert.Equal(t, BooleanDefault, result["use_private_link"])
	})

	t.Run("keeps a value already in state", func(t *testing.T) {
		rawState := map[string]any{"use_private_link": BooleanTrue}

		result, err := v2_23_0_OpenflowDeploymentSnowflakeManagedUsePrivateLinkUpgrader(context.Background(), rawState, nil)
		require.NoError(t, err)
		assert.Equal(t, BooleanTrue, result["use_private_link"])
	})

	t.Run("nil state is a no-op", func(t *testing.T) {
		result, err := v2_23_0_OpenflowDeploymentSnowflakeManagedUsePrivateLinkUpgrader(context.Background(), nil, nil)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

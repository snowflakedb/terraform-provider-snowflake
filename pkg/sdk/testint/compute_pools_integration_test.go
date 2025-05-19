//go:build !account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_ComputePools(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// TODO(next PR): Add more tests.

	t.Run("show: with like", func(t *testing.T) {
		computePool, funcCleanup := testClientHelper().ComputePool.Create(t)
		t.Cleanup(funcCleanup)

		computePools, err := client.ComputePools.Show(ctx, sdk.NewShowComputePoolRequest().WithLike(sdk.Like{Pattern: &computePool.Name}))
		require.NoError(t, err)
		require.Equal(t, 1, len(computePools))
		require.NotNil(t, computePool)
		require.Equal(t, computePool.Name, computePools[0].Name)
		// TODO(next PR): Add more assertions, based on generated builders.
		// Note that the value of updated_on may differ between both SHOW calls.
	})
}

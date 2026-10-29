package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

func Test_DescribeInSdk_computePool(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.ComputePool}, nil, nil)
	require.True(t, model.Hooks.DescribeInSdk.Generate)

	got := renderLifecycle(t, model)
	require.Contains(t, got, `computePoolDetails, err := meta.(*provider.Context).Client.ComputePools.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}`)
	require.NotContains(t, got, "computePoolDescribeInSdkExt")
}

func Test_DescribeInSdk_hookModes(t *testing.T) {
	t.Run("generated", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithDescribeInSdk(ssot.Generated())))
		require.Equal(t, `poolDetails, err := meta.(*provider.Context).Client.Pools.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}
`, renderTemplate(t, "describeInSdk", model))
	})
	t.Run("ext call", func(t *testing.T) {
		model := lifecycleModel(func(obj *ssot.Object) { obj.DescribeFailRead = true })
		require.Equal(t, `poolDetails, err := poolDescribeInSdkExt(ctx, meta, id)
		if err != nil {
			return diag.FromErr(err)
		}
`, renderTemplate(t, "describeInSdk", model))
	})
	t.Run("ext call best effort", func(t *testing.T) {
		model := lifecycleModel(func(obj *ssot.Object) { obj.DescribeUsesShowId = true })
		require.Equal(t, `
		// Describe is best-effort: log + skip Set, do not fail Read. Uses SHOW object's id.
		poolDetails, describeErr := poolDescribeInSdkExt(ctx, meta, pool.ID())
`, renderTemplate(t, "describeInSdk", model))
	})
	t.Run("omitted", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithDescribeInSdk(ssot.Omitted())))
		require.Empty(t, renderTemplate(t, "describeInSdk", model))
	})
}

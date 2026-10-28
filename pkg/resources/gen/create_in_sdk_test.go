package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

func Test_CreateInSdk_computePool(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.ComputePool}, nil, nil)
	require.True(t, model.Hooks.CreateInSdk.Generate)

	got := renderLifecycle(t, model)
	require.Contains(t, got, `if err := meta.(*provider.Context).Client.ComputePools.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}`)
	require.NotContains(t, got, "computePoolCreateInSdkExt")
}

func Test_CreateInSdk_hookModes(t *testing.T) {
	t.Run("generated", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithCreateInSdk(ssot.Generated())))
		require.Equal(t, `if err := meta.(*provider.Context).Client.Pools.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}`, renderTemplate(t, "createInSdk", model))
	})
	t.Run("ext call", func(t *testing.T) {
		model := lifecycleModel()
		require.Equal(t, `if diags := poolCreateInSdkExt(ctx, meta, id, req); diags != nil {
		return diags
	}
`, renderTemplate(t, "createInSdk", model))
	})
	t.Run("omitted", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithCreateInSdk(ssot.Omitted())))
		require.Empty(t, renderTemplate(t, "createInSdk", model))
	})
}

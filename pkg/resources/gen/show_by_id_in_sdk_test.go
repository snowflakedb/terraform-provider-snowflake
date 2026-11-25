package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

// TODO [next PRs]: Production definitions should not be used in tests; replace them with fake objects or remove the test.
func Test_ShowByIdInSdk_computePool(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.ComputePool}, nil, nil)
	require.True(t, model.Hooks.ShowByIdInSdk.Generate)
	require.True(t, model.Hooks.ShowByIdSafelyInSdk.Generate)

	got := renderLifecycle(t, model)
	require.Contains(t, got, `computePool, err := meta.(*provider.Context).Client.ComputePools.ShowByID(ctx, id)`)
	require.Contains(t, got, `computePool, err := meta.(*provider.Context).Client.ComputePools.ShowByIDSafely(ctx, id)
		if err != nil {
			return showByIdErrorDiagnostics(d, sdk.ObjectTypeComputePool, id, err)
		}`)
	require.NotContains(t, got, "computePoolShowByIdInSdkExt")
	require.NotContains(t, got, "computePoolShowByIdSafelyInSdkExt")
}

func Test_ShowByIdInSdk_hookModes(t *testing.T) {
	t.Run("generated", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithShowByIdInSdk(ssot.Generated())))
		require.Equal(t, `pool, err := meta.(*provider.Context).Client.Pools.ShowByID(ctx, id)
`, renderTemplate(t, "showByIdInSdk", model))
	})
	t.Run("ext call", func(t *testing.T) {
		model := lifecycleModel()
		require.Equal(t, `pool, err := poolShowByIdInSdkExt(ctx, meta, id)
`, renderTemplate(t, "showByIdInSdk", model))
	})
	t.Run("omitted", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithShowByIdInSdk(ssot.Omitted())))
		require.Empty(t, renderTemplate(t, "showByIdInSdk", model))
	})
}

func Test_ShowByIdSafelyInSdk_hookModes(t *testing.T) {
	t.Run("generated", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithShowByIdSafelyInSdk(ssot.Generated())))
		require.Equal(t, `pool, err := meta.(*provider.Context).Client.Pools.ShowByIDSafely(ctx, id)
		if err != nil {
			return showByIdErrorDiagnostics(d, sdk.ObjectTypePool, id, err)
		}
`, renderTemplate(t, "showByIdSafelyInSdk", model))
	})
	t.Run("ext call", func(t *testing.T) {
		model := lifecycleModel()
		require.Equal(t, `pool, diags := poolShowByIdSafelyInSdkExt(ctx, d, meta, id)
		if diags != nil {
			return diags
		}
`, renderTemplate(t, "showByIdSafelyInSdk", model))
	})
	t.Run("omitted", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithShowByIdSafelyInSdk(ssot.Omitted())))
		require.Empty(t, renderTemplate(t, "showByIdSafelyInSdk", model))
	})
}

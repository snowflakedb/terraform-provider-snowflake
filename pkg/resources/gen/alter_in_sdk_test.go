package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

func Test_AlterInSdk_computePool(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.ComputePool}, nil, nil)
	require.True(t, model.Hooks.AlterSetInSdk.Generate)
	require.True(t, model.Hooks.AlterUnsetInSdk.Generate)

	got := renderLifecycle(t, model)
	require.Contains(t, got, `if !reflect.DeepEqual(*setUnset.set, *sdk.NewComputePoolSetRequest()) {
		if err := meta.(*provider.Context).Client.ComputePools.Alter(ctx, sdk.NewAlterComputePoolRequest(id).WithSet(*setUnset.set)); err != nil {
			return diag.FromErr(err)
		}
	}`)
	require.Contains(t, got, `if !reflect.DeepEqual(*setUnset.unset, *sdk.NewComputePoolUnsetRequest()) {
		if err := meta.(*provider.Context).Client.ComputePools.Alter(ctx, sdk.NewAlterComputePoolRequest(id).WithUnset(*setUnset.unset)); err != nil {
			return diag.FromErr(err)
		}
	}`)
	require.NotContains(t, got, "computePoolAlterSetInSdkExt")
	require.NotContains(t, got, "computePoolAlterUnsetInSdkExt")
}

func Test_AlterInSdk_hookModes(t *testing.T) {
	tests := []struct {
		name             string
		hook             func(ssot.ResourceHook) ssot.ResourceHookOption
		generated        string
		generatedPartial string
		ext              string
		extPartial       string
	}{
		{
			name: "alterSetInSdk", hook: ssot.WithAlterSetInSdk,
			generated: `if !reflect.DeepEqual(*setUnset.set, *sdk.NewPoolSetRequest()) {
		if err := meta.(*provider.Context).Client.Pools.Alter(ctx, sdk.NewAlterPoolRequest(id).WithSet(*setUnset.set)); err != nil {
			return diag.FromErr(err)
		}
	}
`,
			generatedPartial: `if !reflect.DeepEqual(*setUnset.set, *sdk.NewPoolSetRequest()) {
		if err := meta.(*provider.Context).Client.Pools.Alter(ctx, sdk.NewAlterPoolRequest(id).WithSet(*setUnset.set)); err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}
`,
			ext: `if err := poolAlterSetInSdkExt(ctx, d, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
`,
			extPartial: `if err := poolAlterSetInSdkExt(ctx, d, meta, id, setUnset); err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}
`,
		},
		{
			name: "alterUnsetInSdk", hook: ssot.WithAlterUnsetInSdk,
			generated: `if !reflect.DeepEqual(*setUnset.unset, *sdk.NewPoolUnsetRequest()) {
		if err := meta.(*provider.Context).Client.Pools.Alter(ctx, sdk.NewAlterPoolRequest(id).WithUnset(*setUnset.unset)); err != nil {
			return diag.FromErr(err)
		}
	}
`,
			generatedPartial: `if !reflect.DeepEqual(*setUnset.unset, *sdk.NewPoolUnsetRequest()) {
		if err := meta.(*provider.Context).Client.Pools.Alter(ctx, sdk.NewAlterPoolRequest(id).WithUnset(*setUnset.unset)); err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}
`,
			ext: `if err := poolAlterUnsetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
`,
			extPartial: `if err := poolAlterUnsetInSdkExt(ctx, meta, id, setUnset); err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, mode := range []struct {
				name        string
				hook        ssot.ResourceHook
				want        string
				wantPartial string
			}{
				{"generated", ssot.Generated(), tt.generated, tt.generatedPartial},
				{"ext call", ssot.ExtCall(), tt.ext, tt.extPartial},
			} {
				t.Run(mode.name, func(t *testing.T) {
					model := lifecycleModel(withHooks(tt.hook(mode.hook)))
					require.Equal(t, mode.want, renderTemplate(t, tt.name, model))
					t.Run("partial state", func(t *testing.T) {
						model.HasPartialOnUpdate = true
						require.Equal(t, mode.wantPartial, renderTemplate(t, tt.name, model))
					})
				})
			}
			t.Run("omitted", func(t *testing.T) {
				model := lifecycleModel(withHooks(tt.hook(ssot.Omitted())))
				require.Empty(t, renderTemplate(t, tt.name, model))
			})
		})
	}
}

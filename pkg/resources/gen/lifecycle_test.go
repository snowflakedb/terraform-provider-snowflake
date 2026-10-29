package gen

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func lifecycleModel(opts ...func(*ssot.Object)) ResourceModel {
	obj := ssot.Object{Name: "Pool", NamePlural: "Pools", IdType: "SchemaObjectIdentifier"}
	for _, opt := range opts {
		opt(&obj)
	}
	return ModelFromInputObject(ResourceObject{Object: obj}, nil, nil)
}

func withHooks(hooks ...ssot.ResourceHookOption) func(*ssot.Object) {
	return func(obj *ssot.Object) { obj.ResourceHooks = hooks }
}

func renderLifecycle(t *testing.T, model ResourceModel) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, LifecycleTemplate.Execute(&b, model))
	return b.String()
}

func Test_lifecycle_setIdEncoding(t *testing.T) {
	t.Run("current id", func(t *testing.T) {
		got := renderLifecycle(t, lifecycleModel())
		require.Contains(t, got, "d.SetId(helpers.EncodeResourceIdentifier(id))")
		require.NotContains(t, got, "EncodeSnowflakeID")
	})

	t.Run("legacy id", func(t *testing.T) {
		got := renderLifecycle(t, lifecycleModel(func(obj *ssot.Object) { obj.LegacyIdEncoding = true }))
		require.Contains(t, got, "d.SetId(helpers.EncodeSnowflakeID(id))")
		require.NotContains(t, got, "EncodeResourceIdentifier")
	})
}

func Test_lifecycle_afterCreate(t *testing.T) {
	t.Run("omitted", func(t *testing.T) {
		require.Equal(t, "return ReadPoolFunc(false)(ctx, d, meta)\n", renderTemplate(t, "afterCreate", lifecycleModel()))
	})

	t.Run("ext call", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithAfterCreate(ssot.ExtCall())))
		require.Equal(t, `afterCreateDiags := poolAfterCreateExt(ctx, d, meta, id)
	return append(afterCreateDiags, ReadPoolFunc(false)(ctx, d, meta)...)
`, renderTemplate(t, "afterCreate", model))
	})
}

func Test_lifecycle_afterAlter(t *testing.T) {
	t.Run("omitted", func(t *testing.T) {
		require.Equal(t, "return ReadPoolFunc(false)(ctx, d, meta)\n", renderTemplate(t, "afterAlter", lifecycleModel()))
	})

	t.Run("ext call", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithAfterAlter(ssot.ExtCall())))
		require.Equal(t, `afterAlterDiags := poolAfterAlterExt(ctx, d, meta, id)
	return append(afterAlterDiags, ReadPoolFunc(false)(ctx, d, meta)...)
`, renderTemplate(t, "afterAlter", model))
	})
}

func Test_lifecycle_importPermissiveMode(t *testing.T) {
	got := renderLifecycle(t, lifecycleModel())
	require.Contains(t, got, "importPermissiveMode := meta.(*provider.Context).Experiments.IsEnabled(experimentalfeatures.ImportPermissiveMode)")
	require.Contains(t, got, "poolSetFieldsNotSetByReadExt(d, id, pool, importPermissiveMode)")

	omitted := renderLifecycle(t, lifecycleModel(withHooks(ssot.WithSetFieldsNotSetByRead(ssot.Omitted()))))
	require.NotContains(t, omitted, "importPermissiveMode")
	require.NotContains(t, omitted, "SetFieldsNotSetByReadExt")
}

func Test_lifecycle_importOmittedWithoutShowByIdInSdk(t *testing.T) {
	require.Contains(t, renderLifecycle(t, lifecycleModel()), "func ImportPool(")

	got := renderLifecycle(t, lifecycleModel(withHooks(ssot.WithShowByIdInSdk(ssot.Omitted()))))
	require.NotContains(t, got, "func ImportPool(")
}

func Test_lifecycle_mappingHooksOmitted(t *testing.T) {
	tests := []struct {
		name     string
		hook     ssot.ResourceHookOption
		function string
	}{
		{"parseIdFromConfig", ssot.WithParseIdFromConfig(ssot.Omitted()), "func poolParseIdFromConfig("},
		{"newCreateRequest", ssot.WithNewCreateRequest(ssot.Omitted()), "func poolNewCreateRequest("},
		{"applyCreateOptionals", ssot.WithApplyCreateOptionals(ssot.Omitted()), "func poolApplyCreateOptionals("},
		{"applySetUnset", ssot.WithApplySetUnset(ssot.Omitted()), "func poolApplySetUnset("},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := lifecycleModel(withHooks(tt.hook))
			require.Empty(t, renderTemplate(t, tt.name, model))
			require.NotContains(t, renderLifecycle(t, model), tt.function)
		})
	}
}

func Test_lifecycle_parseIdFromConfigExt(t *testing.T) {
	model := lifecycleModel()
	require.Equal(t, "id, err := poolParseIdFromConfigExt(d)\n", renderTemplate(t, "parseIdFromConfig", model))
	require.NotContains(t, renderLifecycle(t, model), "func poolParseIdFromConfig(")
}

func Test_lifecycle_dropSafelyInSdk(t *testing.T) {
	t.Run("generated", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithDropSafelyInSdk(ssot.Generated())))
		got := renderLifecycle(t, model)
		require.Contains(t, got, "DeleteContext: TrackingDeleteWrapper(resources.Pool, DeletePool),")
		require.NotContains(t, got, "poolDeleteExt")
		require.Equal(t, `func DeletePool(ctx context.Context, d *tfschema.ResourceData, meta any) diag.Diagnostics {
	return ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.Pools.DropSafely
		},
	)(ctx, d, meta)
}
`, renderTemplate(t, "generatedDelete", model))
		require.Contains(t, got, "func DeletePool(")
	})

	t.Run("legacy id", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithDropSafelyInSdk(ssot.Generated())), func(obj *ssot.Object) {
			obj.LegacyIdEncoding = true
		})
		require.Equal(t, `func DeletePool(ctx context.Context, d *tfschema.ResourceData, meta any) diag.Diagnostics {
	return ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErrLegacy[sdk.SchemaObjectIdentifier],
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.Pools.DropSafely
		},
	)(ctx, d, meta)
}
`, renderTemplate(t, "generatedDelete", model))
	})

	t.Run("preview feature", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithDropSafelyInSdk(ssot.Generated())), func(obj *ssot.Object) {
			obj.PreviewFeature = true
		})
		got := renderLifecycle(t, model)
		require.Contains(t, got, "DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.PoolResource), TrackingDeleteWrapper(resources.Pool, DeletePool)),")
		require.NotContains(t, got, "poolDeleteExt")
	})

	for name, hook := range map[string]ssot.ResourceHook{
		"ext call": ssot.ExtCall(),
		"omitted":  ssot.Omitted(),
	} {
		t.Run(name, func(t *testing.T) {
			model := lifecycleModel(withHooks(ssot.WithDropSafelyInSdk(hook)))
			got := renderLifecycle(t, model)
			require.Contains(t, got, "DeleteContext: poolDeleteExt,")
			require.NotContains(t, got, "func DeletePool(")
		})
	}
}

func Test_lifecycle_outputSetHooks(t *testing.T) {
	tests := []struct {
		name      string
		template  string
		withHook  func(ssot.ResourceHook) ssot.ResourceHookOption
		generated string
		extCall   string
	}{
		{
			name:      "show output set",
			template:  "showOutputSet",
			withHook:  ssot.WithShowOutputSet,
			generated: "d.Set(ShowOutputAttributeName, []map[string]any{schemas.PoolToSchema(pool)}),\n",
			extCall:   "poolShowOutputSetExt(d, pool),\n",
		},
		{
			name:      "describe output set",
			template:  "describeOutputSet",
			withHook:  ssot.WithDescribeOutputSet,
			generated: "d.Set(DescribeOutputAttributeName, []map[string]any{schemas.PoolDetailsToSchema(poolDetails)}),\n",
			extCall:   "poolDescribeOutputSetExt(d, poolDetails),\n",
		},
		{
			name:      "fqn set",
			template:  "fqnSet",
			withHook:  ssot.WithFqnSet,
			generated: "d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),\n",
			extCall:   "poolFqnSetExt(d, id),\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name+" generated", func(t *testing.T) {
			model := lifecycleModel(withHooks(tt.withHook(ssot.Generated())))
			require.Equal(t, tt.generated, renderTemplate(t, tt.template, model))
		})

		t.Run(tt.name+" is generated by default", func(t *testing.T) {
			require.Equal(t, tt.generated, renderTemplate(t, tt.template, lifecycleModel()))
		})

		t.Run(tt.name+" ext call", func(t *testing.T) {
			model := lifecycleModel(withHooks(tt.withHook(ssot.ExtCall())))
			require.Equal(t, tt.extCall, renderTemplate(t, tt.template, model))
		})

		t.Run(tt.name+" omitted", func(t *testing.T) {
			model := lifecycleModel(withHooks(tt.withHook(ssot.Omitted())))
			require.Empty(t, renderTemplate(t, tt.template, model))
		})
	}
}

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

func Test_lifecycle_importOmittedWithoutShowByIdInSdk(t *testing.T) {
	require.Contains(t, renderLifecycle(t, lifecycleModel()), "func ImportPool(")

	got := renderLifecycle(t, lifecycleModel(withHooks(ssot.WithShowByIdInSdk(ssot.Omitted()))))
	require.NotContains(t, got, "func ImportPool(")
}

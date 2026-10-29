package gen

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

func Test_ParseId_identifierKindsAndEncoding(t *testing.T) {
	tests := []struct {
		name   string
		idType string
		legacy bool
		want   string
	}{
		{
			name: "account", idType: "AccountObjectIdentifier",
			want: `func poolParseId(id string) (sdk.AccountObjectIdentifier, error) {
	return sdk.ParseAccountObjectIdentifier(id)
}
`,
		},
		{
			name: "database", idType: "DatabaseObjectIdentifier",
			want: `func poolParseId(id string) (sdk.DatabaseObjectIdentifier, error) {
	return sdk.ParseDatabaseObjectIdentifier(id)
}
`,
		},
		{
			name: "schema", idType: "SchemaObjectIdentifier",
			want: `func poolParseId(id string) (sdk.SchemaObjectIdentifier, error) {
	return sdk.ParseSchemaObjectIdentifier(id)
}
`,
		},
		{
			name: "legacy account", idType: "AccountObjectIdentifier", legacy: true,
			want: `func poolParseId(id string) (sdk.AccountObjectIdentifier, error) {
	return helpers.DecodeSnowflakeIDErrLegacy[sdk.AccountObjectIdentifier](id)
}
`,
		},
		{
			name: "legacy database", idType: "DatabaseObjectIdentifier", legacy: true,
			want: `func poolParseId(id string) (sdk.DatabaseObjectIdentifier, error) {
	return helpers.DecodeSnowflakeIDErrLegacy[sdk.DatabaseObjectIdentifier](id)
}
`,
		},
		{
			name: "legacy schema", idType: "SchemaObjectIdentifier", legacy: true,
			want: `func poolParseId(id string) (sdk.SchemaObjectIdentifier, error) {
	return helpers.DecodeSnowflakeIDErrLegacy[sdk.SchemaObjectIdentifier](id)
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := lifecycleModel(func(obj *ssot.Object) {
				obj.IdType = tt.idType
				obj.LegacyIdEncoding = tt.legacy
				obj.ResourceHooks = []ssot.ResourceHookOption{ssot.WithParseId(ssot.Generated())}
			})
			require.Equal(t, tt.want, renderTemplate(t, "generatedParseId", model))
		})
	}
}

func Test_ParseId_computePoolLifecycleUsesGeneratedParser(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.ComputePool}, nil, nil)
	got := renderLifecycle(t, model)
	require.Equal(t, 3, strings.Count(got, "computePoolParseId(d.Id())"))
	require.Contains(t, got, "func computePoolParseId(id string) (sdk.AccountObjectIdentifier, error)")
	require.NotContains(t, got, "computePoolParseIdExt")
}

func Test_ParseId_hookModes(t *testing.T) {
	t.Run("generated", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithParseId(ssot.Generated())))
		require.Equal(t, "poolParseId", model.ParseIdFunc())
		require.Equal(t, "id, err := poolParseId(d.Id())\n", renderTemplate(t, "parseId", model))
		require.Contains(t, renderLifecycle(t, model), "func poolParseId(")
	})
	t.Run("ext call", func(t *testing.T) {
		model := lifecycleModel()
		require.Equal(t, "poolParseIdExt", model.ParseIdFunc())
		require.Equal(t, "id, err := poolParseIdExt(d.Id())\n", renderTemplate(t, "parseId", model))
		require.NotContains(t, renderLifecycle(t, model), "func poolParseId(")
	})
	t.Run("omitted", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithParseId(ssot.Omitted())))
		require.Empty(t, renderTemplate(t, "parseId", model))
		require.NotContains(t, renderLifecycle(t, model), "func poolParseId(")
	})
}

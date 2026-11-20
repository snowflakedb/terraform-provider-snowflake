package telemetry

import (
	"runtime"
	"testing"

	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/stretchr/testify/require"
)

func Test_commonFields(t *testing.T) {
	got := commonFields(typeProviderInit, "span-1")
	require.Equal(t, source, got["source"])
	require.Equal(t, string(typeProviderInit), got["type"])
	require.Equal(t, "1", got["json_schema_version"])
	require.NotEmpty(t, got["version"])
	require.Equal(t, "span-1", got["span_id"])
}

func Test_NewSpanID(t *testing.T) {
	t.Run("has UUID length", func(t *testing.T) {
		id := NewSpanID()
		require.Len(t, id, 36)
	})

	t.Run("subsequent ids are unique", func(t *testing.T) {
		ids := make(map[string]struct{})
		for range 100 {
			id := NewSpanID()
			require.NotContains(t, ids, id)
			ids[id] = struct{}{}
		}
	})
}

func Test_providerInitFields(t *testing.T) {
	providerCtx := &internalprovider.Context{
		Experiments: experimentalfeatures.New(
			[]string{string(experimentalfeatures.HierarchyRenames)},
			[]string{string(experimentalfeatures.UserEnableDefaultWorkloadIdentity)},
		),
		EnabledFeatures: []string{
			string(previewfeatures.AlertResource),
		},
	}
	got := providerInitFields(providerCtx)
	require.Equal(t, string(experimentalfeatures.HierarchyRenames), got["experimental_features_enabled"])
	require.Equal(t, string(experimentalfeatures.UserEnableDefaultWorkloadIdentity), got["experimental_features_disabled"])
	require.Equal(t, "ACCOUNT_PARAMETER_CATALOG_WRITE_PATH,HIERARCHY_RENAMES,OBJECT_PARAMETER_PARSE_ON_SET", got["experimental_features_effective_enabled"])
	require.Equal(t, string(previewfeatures.AlertResource), got["preview_features_enabled"])
	require.Equal(t, runtime.GOOS, got["os"])
	require.Equal(t, runtime.GOARCH, got["arch"])
	require.Contains(t, got, "ci_environment")
	require.Contains(t, got, "agent_environment")
	require.Contains(t, got, "terraform_host")
	require.Contains(t, got, "auth_type")
}

func Test_Emit_skipsNilClient(t *testing.T) {
	require.NotPanics(t, func() {
		emit(t.Context(), nil, typeProviderInit, "span", nil)
	})
}

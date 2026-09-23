package telemetry

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/tracking"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/stretchr/testify/require"
)

func Test_hashID(t *testing.T) {
	require.Empty(t, hashID(""))

	raw := `"MY_DB"|"PUBLIC"`
	got := hashID(raw)
	sum := sha256.Sum256([]byte(raw))
	require.Equal(t, hex.EncodeToString(sum[:]), got)
	require.Len(t, got, 64)
	require.NotContains(t, got, "MY_DB")
}

func Test_datasourceOpFields(t *testing.T) {
	rawID := "databases_read"
	got := datasourceOpFields(datasources.Databases, rawID)
	require.Equal(t, datasources.Databases.String(), got["datasource_name"])
	require.Equal(t, string(tracking.ReadOperation), got["operation_type"])
	require.Equal(t, hashID(rawID), got["datasource_id"])
	require.NotContains(t, got["datasource_id"], "databases_read")
}

func Test_EmitDatasourceOp_skipsNilMeta(t *testing.T) {
	require.NotPanics(t, func() {
		EmitDatasourceOp(t.Context(), nil, datasources.Databases, "id")
	})
}

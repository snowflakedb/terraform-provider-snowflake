package telemetry

import (
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/tracking"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/stretchr/testify/require"
)

func Test_datasourceOpFields(t *testing.T) {
	duration := 1500 * time.Millisecond

	t.Run("success", func(t *testing.T) {
		diags := diag.Diagnostics{{Severity: diag.Warning, Summary: "careful"}}
		got := datasourceOpFields(datasources.Databases, duration, diags, nil)

		require.Equal(t, datasources.Databases.String(), got["datasource_name"])
		require.Equal(t, string(tracking.ReadOperation), got["operation_type"])
		require.Equal(t, string(statusSuccess), got["status"])
		require.Equal(t, "1.5", got["duration_s"])
	})

	t.Run("failure", func(t *testing.T) {
		diags := diag.Diagnostics{
			{Severity: diag.Warning, Summary: "careful"},
			{Severity: diag.Error, Summary: "boom", Detail: "select secret"},
		}
		got := datasourceOpFields(datasources.Databases, duration, diags, nil)

		require.Equal(t, string(statusFailure), got["status"])
		require.Equal(t, "1.5", got["duration_s"])
	})

	t.Run("panic", func(t *testing.T) {
		got := datasourceOpFields(datasources.Databases, duration, diag.Errorf("ignored"), "exploded")

		require.Equal(t, string(statusPanic), got["status"])
	})
}

func Test_formatSeconds(t *testing.T) {
	require.Equal(t, "1.5", formatSeconds(1500*time.Millisecond))
	require.Equal(t, "1", formatSeconds(time.Second))
	require.Equal(t, "0.001", formatSeconds(time.Millisecond))
	require.Equal(t, "0.002", formatSeconds(1500*time.Microsecond))
	require.Equal(t, "1.235", formatSeconds(1234500*time.Microsecond))
	require.Equal(t, "0", formatSeconds(4131*time.Nanosecond))
}

func Test_EmitDatasourceOp_skipsNilMeta(t *testing.T) {
	require.NotPanics(t, func() {
		EmitDatasourceOp(t.Context(), nil, datasources.Databases, time.Second, diag.Errorf("boom"), nil)
	})
}

package telemetry

import (
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/tracking"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/stretchr/testify/require"
)

func Test_resourceOpFields(t *testing.T) {
	duration := 1500 * time.Millisecond

	t.Run("success", func(t *testing.T) {
		diags := diag.Diagnostics{{Severity: diag.Warning, Summary: "careful"}}
		got := resourceOpFields(resources.Warehouse, tracking.CreateOperation, duration, diags, nil)

		require.Equal(t, resources.Warehouse.String(), got["resource_name"])
		require.Equal(t, string(tracking.CreateOperation), got["operation_type"])
		require.Equal(t, string(statusSuccess), got["status"])
		require.Equal(t, "1.5", got["duration_s"])
	})

	t.Run("failure", func(t *testing.T) {
		diags := diag.Diagnostics{
			{Severity: diag.Warning, Summary: "careful"},
			{Severity: diag.Error, Summary: "boom", Detail: "select secret"},
		}
		got := resourceOpFields(resources.Warehouse, tracking.UpdateOperation, duration, diags, nil)

		require.Equal(t, string(statusFailure), got["status"])
		require.Equal(t, "1.5", got["duration_s"])
	})

	t.Run("panic", func(t *testing.T) {
		got := resourceOpFields(resources.Warehouse, tracking.DeleteOperation, duration, diag.Errorf("ignored"), "exploded")

		require.Equal(t, string(statusPanic), got["status"])
	})
}

func Test_EmitResourceOp_skipsNilMeta(t *testing.T) {
	require.NotPanics(t, func() {
		EmitResourceOp(t.Context(), nil, resources.Warehouse, tracking.ReadOperation, time.Second, diag.Errorf("boom"), nil)
	})
}

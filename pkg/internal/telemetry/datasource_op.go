package telemetry

import (
	"context"
	"strconv"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/tracking"
	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func EmitDatasourceOp(ctx context.Context, meta any, datasourceName datasources.Datasource, duration time.Duration, diags diag.Diagnostics, recovered any) {
	providerCtx, ok := meta.(*internalprovider.Context)
	if !ok || providerCtx == nil {
		return
	}
	emit(ctx, providerCtx.Client, typeDatasourceOp, providerCtx.SpanID, datasourceOpFields(datasourceName, duration, diags, recovered))
}

func datasourceOpFields(datasourceName datasources.Datasource, duration time.Duration, diags diag.Diagnostics, recovered any) map[string]string {
	return map[string]string{
		"datasource_name": datasourceName.String(),
		"operation_type":  string(tracking.ReadOperation),
		"status":          string(getOperationStatus(diags, recovered)),
		"duration_s":      formatSeconds(duration),
	}
}

// formatSeconds reports duration_s in fractional seconds.
func formatSeconds(duration time.Duration) string {
	return strconv.FormatFloat(duration.Seconds(), 'f', -1, 64)
}

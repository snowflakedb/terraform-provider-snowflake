package telemetry

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/tracking"
	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
)

func EmitDatasourceOp(ctx context.Context, meta any, datasourceName datasources.Datasource, terraformID string) {
	providerCtx, ok := meta.(*internalprovider.Context)
	if !ok || providerCtx == nil {
		return
	}
	emit(ctx, providerCtx.Client, typeDatasourceOp, providerCtx.SpanID, datasourceOpFields(datasourceName, terraformID))
}

func datasourceOpFields(datasourceName datasources.Datasource, terraformID string) map[string]string {
	return map[string]string{
		"datasource_name": datasourceName.String(),
		"operation_type":  string(tracking.ReadOperation),
		"datasource_id":   hashID(terraformID),
	}
}

func hashID(id string) string {
	if id == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(id))
	return hex.EncodeToString(sum[:])
}

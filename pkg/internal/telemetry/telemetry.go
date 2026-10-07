package telemetry

import (
	"context"
	"log"
	"maps"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/tracking"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

const (
	source = "terraform_provider"
)

type eventType string

const (
	typeProviderInit eventType = "provider_init"
	typeDatasourceOp eventType = "datasource_op"
	typeResourceOp   eventType = "resource_op"
)

type operationStatus string

const (
	statusSuccess operationStatus = "success"
	statusFailure operationStatus = "failure"
	statusPanic   operationStatus = "panic"
)

func getOperationStatus(diags diag.Diagnostics, recovered any) operationStatus {
	switch {
	case recovered != nil:
		return statusPanic
	case diags.HasError():
		return statusFailure
	default:
		return statusSuccess
	}
}

// NewSpanID returns a UUID used to correlate events for one provider configure/run.
func NewSpanID() (string, error) {
	return uuid.GenerateUUID()
}

func commonFields(eventType eventType, spanID string) map[string]string {
	return map[string]string{
		"source":              source,
		"type":                string(eventType),
		"json_schema_version": tracking.CurrentSchemaVersion,
		"version":             tracking.ProviderVersion,
		"span_id":             spanID,
	}
}

// emit queues a custom driver telemetry event. Failures are logged and never
// returned to the caller so Terraform CRUD/configure cannot fail because of telemetry.
func emit(ctx context.Context, client *sdk.Client, eventType eventType, spanID string, extra map[string]string) {
	if client == nil || spanID == "" {
		return
	}
	data := commonFields(eventType, spanID)
	maps.Copy(data, extra)
	if err := client.AddTelemetry(ctx, data); err != nil {
		log.Printf("[DEBUG] failed to emit %s telemetry: %v", eventType, err)
	}
}

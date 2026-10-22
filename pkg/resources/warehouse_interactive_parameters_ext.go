package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Schema maps
// =============================================================================

var warehouseInteractiveParametersOutputSchemaExt = map[string]*schema.Schema{
	ParametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW PARAMETERS IN WAREHOUSE` for the given interactive warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowWarehouseParametersSchemaInteractive,
		},
	},
}

// =============================================================================
// Read handler
// =============================================================================

// Generated set-fields writes AccountObjectIdentifier.FullyQualifiedName(), which quotes the
// identifier. Interactive warehouses store fallback_warehouse with Name().
func warehouseInteractiveSetParametersFieldsExt(d *schema.ResourceData, parameters *sdk.WarehouseParametersDetails) diag.Diagnostics {
	return JoinDiags(
		setResourceData(d, "max_concurrency_level", parameters.MaxConcurrencyLevel.Value),
		setResourceData(d, "statement_queued_timeout_in_seconds", parameters.StatementQueuedTimeoutInSeconds.Value),
		setResourceData(d, "statement_timeout_in_seconds", parameters.StatementTimeoutInSeconds.Value),
		setResourceData(d, "fallback_warehouse", parameters.FallbackWarehouse.Value.Name()),
	)
}

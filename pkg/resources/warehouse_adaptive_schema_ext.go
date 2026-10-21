package resources

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var warehouseAdaptiveIdentitySchemaExt = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Identifier for the adaptive warehouse; must be unique for your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
}

var warehouseAdaptiveAttributesSchemaExt = map[string]*schema.Schema{
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the adaptive warehouse.",
	},
	"warehouse_type": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Specifies the type for the adaptive warehouse. This field is used for checking external changes and recreating the resource if needed.",
	},
	"max_query_performance_level": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToMaxQueryPerformanceLevel),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToMaxQueryPerformanceLevel),
		Description:      fmt.Sprintf("Specifies the maximum query performance level for the adaptive warehouse. Determines the initial compute capacity. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllMaxQueryPerformanceLevels)),
	},
	"query_throughput_multiplier": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		Description:      "Specifies the query throughput multiplier for the adaptive warehouse.",
	},
	"resource_monitor": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInShow("resource_monitor")),
		Description:      relatedResourceDescription("Specifies the name of a resource monitor that is explicitly assigned to the adaptive warehouse.", resources.ResourceMonitor),
	},
}

var warehouseAdaptiveShowOutputSchemaExt = map[string]*schema.Schema{
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW WAREHOUSES` for the given adaptive warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowWarehouseAdaptiveSchema,
		},
	},
}

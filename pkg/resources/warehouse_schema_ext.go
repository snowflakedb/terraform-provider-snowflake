package resources

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var warehouseIdentitySchemaExt = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Identifier for the virtual warehouse; must be unique for your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
}

var warehouseAttributesSchemaExt = map[string]*schema.Schema{
	"warehouse_type": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToWarehouseTypeUserSettable),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToWarehouseTypeUserSettable), IgnoreChangeToCurrentSnowflakeValueInShow("type")),
		Description:      fmt.Sprintf("Specifies warehouse type. Valid values are (case-insensitive): %s. Warehouse needs to be suspended to change its type. Provider will handle automatic suspension and resumption if needed.", possibleValuesListed(sdk.ValidWarehouseTypesRegularString)),
	},
	"warehouse_size": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToWarehouseSize),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToWarehouseSize), IgnoreChangeToCurrentSnowflakeValueInShow("size")),
		Description:      fmt.Sprintf("Specifies the size of the virtual warehouse. Valid values are (case-insensitive): %s. Consult [warehouse documentation](https://docs.snowflake.com/en/sql-reference/sql/create-warehouse#optional-properties-objectproperties) for the details. Note: removing the size from config will result in the resource recreation.", possibleValuesListed(sdk.ValidWarehouseSizesString)),
	},
	"max_cluster_count": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("max_cluster_count"),
		Description:      "Specifies the maximum number of server clusters for the warehouse.",
	},
	"min_cluster_count": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("min_cluster_count"),
		Description:      "Specifies the minimum number of server clusters for the warehouse (only applies to multi-cluster warehouses).",
	},
	"scaling_policy": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToScalingPolicy),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToWarehouseType), IgnoreChangeToCurrentSnowflakeValueInShow("scaling_policy")),
		Description:      fmt.Sprintf("Specifies the policy for automatically starting and shutting down clusters in a multi-cluster warehouse running in Auto-scale mode. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.ValidWarehouseScalingPoliciesString)),
	},
	"auto_suspend": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("auto_suspend"),
		Description:      "Specifies the number of seconds of inactivity after which a warehouse is automatically suspended.",
		Default:          IntDefault,
	},
	"auto_resume": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("auto_resume"),
		Description:      booleanStringFieldDescription("Specifies whether to automatically resume a warehouse when a SQL statement (e.g. query) is submitted to it."),
		Default:          BooleanDefault,
	},
	"initially_suspended": {
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: IgnoreAfterCreation,
		Description:      "Specifies whether the warehouse is created initially in the ‘Suspended’ state.",
	},
	"resource_monitor": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInShow("resource_monitor")),
		Description:      relatedResourceDescription("Specifies the name of a resource monitor that is explicitly assigned to the warehouse.", resources.ResourceMonitor),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the warehouse.",
	},
	"enable_query_acceleration": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("enable_query_acceleration"),
		Description:      booleanStringFieldDescription("Specifies whether to enable the query acceleration service for queries that rely on this warehouse for compute resources."),
		Default:          BooleanDefault,
	},
	"query_acceleration_max_scale_factor": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 100)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("query_acceleration_max_scale_factor"),
		Description:      "Specifies the maximum scale factor for leasing compute resources for query acceleration. The scale factor is used as a multiplier based on warehouse size.",
		Default:          IntDefault,
	},
	"resource_constraint": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToWarehouseResourceConstraint),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToWarehouseResourceConstraint), IgnoreChangeToCurrentSnowflakeValueInShow("resource_constraint")),
		Description: joinWithSpace("Specifies the resource constraint for the warehouse. Only available for snowpark-optimized warehouses. For setting generation please use the `generation` field.",
			"Please check [Snowflake documentation](https://docs.snowflake.com/en/sql-reference/sql/create-warehouse#optional-properties-objectproperties) for required warehouse sizes for each resource constraint.",
			fmt.Sprintf("Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllWarehouseResourceConstraintsWithoutGenerations)),
			"Gen2 warehouses are not available in all regions. Please consult the [Snowflake Gen2 Region Availability documentation](https://docs.snowflake.com/en/user-guide/warehouses-gen2#region-availability) prior to configuration."),
		ConflictsWith: []string{"generation"},
	},
	"generation": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToWarehouseGeneration),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToWarehouseGeneration), IgnoreChangeToCurrentSnowflakeValueInShow("generation")),
		Description: joinWithSpace(fmt.Sprintf("Specifies the generation for the warehouse. Only available for standard warehouses. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllWarehouseGenerations)),
			"Gen2 warehouses are not available in all regions. Please consult the [Snowflake Gen2 Region Availability documentation](https://docs.snowflake.com/en/user-guide/warehouses-gen2#region-availability) prior to configuration."),
		ConflictsWith: []string{"resource_constraint"},
	},
}

var warehouseShowOutputSchemaExt = map[string]*schema.Schema{
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW WAREHOUSES` for the given warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowWarehouseRegularSchema,
		},
	},
}

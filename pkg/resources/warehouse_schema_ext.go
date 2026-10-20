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

func warehouseAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	s["warehouse_type"].DiffSuppressFunc = SuppressIfAny(NormalizeAndCompare(sdk.ToWarehouseTypeUserSettable), IgnoreChangeToCurrentSnowflakeValueInShow("type"))
	s["warehouse_type"].Description = fmt.Sprintf("Specifies warehouse type. Valid values are (case-insensitive): %s. Warehouse needs to be suspended to change its type. Provider will handle automatic suspension and resumption if needed.", possibleValuesListed(sdk.ValidWarehouseTypesRegularString))

	s["warehouse_size"].DiffSuppressFunc = SuppressIfAny(NormalizeAndCompare(sdk.ToWarehouseSize), IgnoreChangeToCurrentSnowflakeValueInShow("size"))
	s["warehouse_size"].Description = fmt.Sprintf("Specifies the size of the virtual warehouse. Valid values are (case-insensitive): %s. Consult [warehouse documentation](https://docs.snowflake.com/en/sql-reference/sql/create-warehouse#optional-properties-objectproperties) for the details. Note: removing the size from config will result in the resource recreation.", possibleValuesListed(sdk.ValidWarehouseSizesString))

	s["max_cluster_count"].DiffSuppressFunc = IgnoreChangeToCurrentSnowflakeValueInShow("max_cluster_count")
	s["min_cluster_count"].DiffSuppressFunc = IgnoreChangeToCurrentSnowflakeValueInShow("min_cluster_count")

	// Pre-existing: this suppressor normalizes with ToWarehouseType, not ToScalingPolicy.
	s["scaling_policy"].DiffSuppressFunc = SuppressIfAny(NormalizeAndCompare(sdk.ToWarehouseType), IgnoreChangeToCurrentSnowflakeValueInShow("scaling_policy"))

	s["auto_suspend"].DiffSuppressFunc = IgnoreChangeToCurrentSnowflakeValueInShow("auto_suspend")
	s["auto_resume"].DiffSuppressFunc = IgnoreChangeToCurrentSnowflakeValueInShow("auto_resume")

	s["resource_monitor"].ValidateDiagFunc = IsValidIdentifier[sdk.AccountObjectIdentifier]()
	s["resource_monitor"].DiffSuppressFunc = SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInShow("resource_monitor"))
	s["resource_monitor"].Description = relatedResourceDescription("Specifies the name of a resource monitor that is explicitly assigned to the warehouse.", resources.ResourceMonitor)

	s["enable_query_acceleration"].DiffSuppressFunc = IgnoreChangeToCurrentSnowflakeValueInShow("enable_query_acceleration")

	s["query_acceleration_max_scale_factor"].ValidateDiagFunc = validation.ToDiagFunc(validation.IntBetween(0, 100))
	s["query_acceleration_max_scale_factor"].DiffSuppressFunc = IgnoreChangeToCurrentSnowflakeValueInShow("query_acceleration_max_scale_factor")

	s["resource_constraint"].DiffSuppressFunc = SuppressIfAny(NormalizeAndCompare(sdk.ToWarehouseResourceConstraint), IgnoreChangeToCurrentSnowflakeValueInShow("resource_constraint"))
	s["resource_constraint"].Description = joinWithSpace("Specifies the resource constraint for the warehouse. Only available for snowpark-optimized warehouses. For setting generation please use the `generation` field.",
		"Please check [Snowflake documentation](https://docs.snowflake.com/en/sql-reference/sql/create-warehouse#optional-properties-objectproperties) for required warehouse sizes for each resource constraint.",
		fmt.Sprintf("Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllWarehouseResourceConstraintsWithoutGenerations)),
		"Gen2 warehouses are not available in all regions. Please consult the [Snowflake Gen2 Region Availability documentation](https://docs.snowflake.com/en/user-guide/warehouses-gen2#region-availability) prior to configuration.")
	s["resource_constraint"].ConflictsWith = []string{"generation"}

	s["generation"].DiffSuppressFunc = SuppressIfAny(NormalizeAndCompare(sdk.ToWarehouseGeneration), IgnoreChangeToCurrentSnowflakeValueInShow("generation"))
	s["generation"].ConflictsWith = []string{"resource_constraint"}

	s["initially_suspended"] = &schema.Schema{
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: IgnoreAfterCreation,
		Description:      "Specifies whether the warehouse is created initially in the ‘Suspended’ state.",
	}
	return s
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

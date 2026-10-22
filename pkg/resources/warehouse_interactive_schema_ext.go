package resources

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var warehouseInteractiveIdentitySchemaExt = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Identifier for the interactive warehouse; must be unique for your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
}

var warehouseInteractiveAttributesSchemaExt = map[string]*schema.Schema{
	"warehouse_size": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToWarehouseSize),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToWarehouseSize), IgnoreChangeToCurrentSnowflakeValueInShow("size")),
		Description:      fmt.Sprintf("Specifies the size of the interactive warehouse. Valid values are (case-insensitive): %s. Note: changing the size briefly suspends and resumes the warehouse to apply the resize (an interactive warehouse cannot be resized while running); removing the size from config will result in the resource recreation.", possibleValuesListed(sdk.ValidWarehouseSizesString)),
	},
	"max_cluster_count": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("max_cluster_count"),
		Description:      "Specifies the maximum number of server clusters for the interactive warehouse.",
	},
	"min_cluster_count": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("min_cluster_count"),
		Description:      "Specifies the minimum number of server clusters for the interactive warehouse (only applies to multi-cluster warehouses).",
	},
	"auto_suspend": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("auto_suspend"),
		Description:      "Specifies the number of seconds of inactivity after which an interactive warehouse is automatically suspended.",
	},
	"auto_resume": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("auto_resume"),
		Description:      booleanStringFieldDescription("Specifies whether to automatically resume an interactive warehouse when a SQL statement (e.g. query) is submitted to it."),
		Default:          BooleanDefault,
	},
	"initially_suspended": {
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: IgnoreAfterCreation,
		Description:      "Specifies whether the interactive warehouse is created initially in the ‘Suspended’ state.",
	},
	"resource_monitor": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInShow("resource_monitor")),
		Description:      relatedResourceDescription("Specifies the name of a resource monitor that is explicitly assigned to the interactive warehouse.", resources.ResourceMonitor),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the interactive warehouse.",
	},
	"tables": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		},
		DiffSuppressFunc: NormalizeAndCompareIdentifiersInSet("tables"),
		Description:      "Specifies the fully qualified names of the tables associated with the interactive warehouse. Changes are applied incrementally (ADD TABLES / DROP TABLES) rather than by full re-association.",
	},
	"warehouse_type": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Specifies the type for the interactive warehouse. This field is used for checking external changes and recreating the resource if needed.",
	},
}

var warehouseInteractiveShowOutputSchemaExt = map[string]*schema.Schema{
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW WAREHOUSES` for the given interactive warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowWarehouseInteractiveSchema,
		},
	},
}

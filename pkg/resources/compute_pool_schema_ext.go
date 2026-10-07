package resources

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var computePoolIdentitySchemaExt = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the compute pool; must be unique for the account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
}

var computePoolAttributesSchemaExt = map[string]*schema.Schema{
	"for_application": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "Specifies the Snowflake Native App name.",
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"min_nodes": {
		Type:             schema.TypeInt,
		Required:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		Description:      "Specifies the minimum number of nodes for the compute pool.",
	},
	"max_nodes": {
		Type:             schema.TypeInt,
		Required:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		Description:      "Specifies the maximum number of nodes for the compute pool.",
	},
	"instance_family": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToComputePoolInstanceFamily),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToComputePoolInstanceFamily)),
		Description: fmt.Sprintf("Identifies the type of machine you want to provision for the nodes in the compute pool. Valid values are (case-insensitive): %s."+
			" Not all instance families are supported in all regions. Run `SHOW COMPUTE POOL INSTANCE FAMILIES` to see the list of supported instance families in your region.", possibleValuesListed(sdk.AllComputePoolInstanceFamilies)),
	},
	"backup_instance_families": {
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: sdkValidation(sdk.ToComputePoolInstanceFamily),
		},
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToComputePoolInstanceFamily),
		Description: fmt.Sprintf("Specifies an ordered list of instance families to fall back on when the primary `instance_family` is unavailable. The order determines the fallback priority."+
			" Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllComputePoolInstanceFamilies)),
	},
	"auto_resume": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("auto_resume"),
		Description:      booleanStringFieldDescription("Specifies whether to automatically resume a compute pool when a service or job is submitted to it."),
		Default:          BooleanDefault,
	},
	"initially_suspended": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreAfterCreation,
		Description:      "Specifies whether the compute pool is created initially in the suspended state. This field is used only when creating a compute pool. Changes on this field are ignored after creation.",
		Default:          BooleanDefault,
	},
	"auto_suspend_secs": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("auto_suspend_secs"),
		Description:      "Number of seconds of inactivity after which you want Snowflake to automatically suspend the compute pool.",
		Default:          IntDefault,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the compute pool.",
	},
}

var computePoolShowOutputSchemaExt = map[string]*schema.Schema{
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW COMPUTE POOLS` for the given compute pool.",
		Elem: &schema.Resource{
			Schema: schemas.ShowComputePoolSchema,
		},
	},
}

var computePoolDescribeOutputSchemaExt = map[string]*schema.Schema{
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE COMPUTE POOL` for the given compute pool.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeComputePoolDetailsSchema,
		},
	},
}

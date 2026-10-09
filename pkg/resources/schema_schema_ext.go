package resources

import (
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var schemaAttributesSchemaExt = map[string]*schema.Schema{
	"with_managed_access": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      booleanStringFieldDescription("Specifies a managed schema. Managed access schemas centralize privilege management with the schema owner."),
		ValidateDiagFunc: validateBooleanString,
		Default:          BooleanDefault,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShowWithMapping("options", func(x any) any {
			return slices.Contains(sdk.ParseCommaSeparatedStringArray(x.(string), false), "MANAGED ACCESS")
		}),
	},
	"is_transient": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      booleanStringFieldDescription("Specifies the schema as transient. Transient schemas do not have a Fail-safe period so they do not incur additional storage costs once they leave Time Travel; however, this means they are also not protected by Fail-safe in the event of a data loss."),
		ValidateDiagFunc: validateBooleanString,
		Default:          BooleanDefault,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShowWithMapping("options", func(x any) any {
			return slices.Contains(sdk.ParseCommaSeparatedStringArray(x.(string), false), "TRANSIENT")
		}),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the schema.",
	},
}

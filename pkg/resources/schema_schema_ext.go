package resources

import (
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var schemaIdentitySchemaExt = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the schema; must be unique for the database in which the schema is created. When the name is `PUBLIC`, during creation the provider checks if this schema has already been created and, in such case, `ALTER` is used to match the desired state."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the schema."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
}

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

var schemaShowOutputSchemaExt = map[string]*schema.Schema{
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW SCHEMA` for the given object.",
		Elem: &schema.Resource{
			Schema: schemas.ShowSchemaSchema,
		},
	},
}

var schemaDescribeOutputSchemaExt = map[string]*schema.Schema{
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE SCHEMA` for the given object. In order to handle this output, one must grant sufficient privileges, e.g. [grant_ownership](./grant_ownership) on all objects in the schema.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeSchemaDetailsSchema,
		},
	},
}

package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ParameterListSchema represents Snowflake parameter object.
var ParameterListSchema = &schema.Schema{
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: ShowParameterSchema,
	},
}

// DescribePropertyListSchema represents Snowflake property object returned by DESCRIBE query.
// TODO [v3]: remove once nested property-list describe_output consumers are flattened to typed details.
var DescribePropertyListSchema = &schema.Schema{
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: SecurityIntegrationPropertySchema,
	},
}

var FullyQualifiedNameSchema = &schema.Schema{
	Type:        schema.TypeString,
	Computed:    true,
	Description: "Fully qualified name of the resource. For more information, see [object name resolution](https://docs.snowflake.com/en/sql-reference/name-resolution).",
}

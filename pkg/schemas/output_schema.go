package schemas

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	fullyQualifiedNameAttributeName = "fully_qualified_name"
	showOutputAttributeName         = "show_output"
	describeOutputAttributeName     = "describe_output"
	parametersAttributeName         = "parameters"
)

// FullyQualifiedNameOutputSchema is the computed fully_qualified_name attribute.
func FullyQualifiedNameOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		fullyQualifiedNameAttributeName: FullyQualifiedNameSchema,
	}
}

// OutputListSchema is the computed list attribute that holds one output block (show_output, describe_output, or parameters).
func OutputListSchema(attributeName, description string, elemSchema map[string]*schema.Schema) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		attributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: description,
			Elem:        &schema.Resource{Schema: elemSchema},
		},
	}
}

func holdsOutputDescription(sql string, objectType sdk.ObjectType, additional ...string) string {
	return strings.TrimSpace(fmt.Sprintf("Holds the output of `%s` for the given %s. %s", sql, strings.ToLower(objectType.String()), strings.Join(additional, " ")))
}

// ShowOutputListSchema is the computed show_output list for SHOW objectType.
func ShowOutputListSchema(objectType sdk.ObjectType, elemSchema map[string]*schema.Schema, additional ...string) map[string]*schema.Schema {
	return OutputListSchema(showOutputAttributeName, holdsOutputDescription("SHOW "+string(objectType.Plural()), objectType, additional...), elemSchema)
}

// DescribeOutputListSchema is the computed describe_output list for DESCRIBE objectType.
func DescribeOutputListSchema(objectType sdk.ObjectType, elemSchema map[string]*schema.Schema, additional ...string) map[string]*schema.Schema {
	return OutputListSchema(describeOutputAttributeName, holdsOutputDescription("DESCRIBE "+objectType.String(), objectType, additional...), elemSchema)
}

// ParametersOutputListSchema is the computed parameters list for SHOW PARAMETERS FOR objectType.
func ParametersOutputListSchema(objectType sdk.ObjectType, elemSchema map[string]*schema.Schema, additional ...string) map[string]*schema.Schema {
	return OutputListSchema(parametersAttributeName, holdsOutputDescription("SHOW PARAMETERS FOR "+objectType.String(), objectType, additional...), elemSchema)
}

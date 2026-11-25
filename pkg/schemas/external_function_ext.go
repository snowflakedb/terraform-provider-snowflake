package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (externalFunctionToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"arguments": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"return_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (externalFunctionToSchemaMapper) additionalToSchema(src *sdk.ExternalFunction, dst map[string]any) {
	dst["arguments"] = src.ArgumentsStrings()
	dst["return_type"] = src.ReturnTypeString()
}

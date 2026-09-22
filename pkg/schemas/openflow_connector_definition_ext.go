package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (openflowConnectorDefinitionToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// categories is a JSON array in SHOW, parsed to []string in the SDK.
		"categories": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
	}
}

func (openflowConnectorDefinitionToSchemaMapper) additionalToSchema(src *sdk.OpenflowConnectorDefinition, dst map[string]any) {
	dst["categories"] = src.Categories
}

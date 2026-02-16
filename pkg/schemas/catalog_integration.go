package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DescribeCatalogIntegrationSchema = map[string]*schema.Schema{
	"property_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"property_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"property_value": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"property_default": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func CatalogIntegrationPropertyToSchema(property sdk.CatalogIntegrationProperty) map[string]any {
	return map[string]any{
		"property_name":    property.Name,
		"property_type":    property.Type,
		"property_value":   property.Value,
		"property_default": property.Default,
	}
}

func CatalogIntegrationPropertiesToSchema(properties []sdk.CatalogIntegrationProperty) []map[string]any {
	result := make([]map[string]any, len(properties))
	for i, p := range properties {
		result[i] = CatalogIntegrationPropertyToSchema(p)
	}
	return result
}

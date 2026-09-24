package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (storageIntegrationAwsDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allowed_locations": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
		"blocked_locations": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
	}
}

func (storageIntegrationAwsDetailsToSchemaMapper) additionalToSchema(src *sdk.StorageIntegrationAwsDetails, dst map[string]any) {
	dst["allowed_locations"] = src.AllowedLocations
	dst["blocked_locations"] = src.BlockedLocations
}

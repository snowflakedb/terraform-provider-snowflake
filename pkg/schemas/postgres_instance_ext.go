package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (postgresInstanceToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"is_ha": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}

func (postgresInstanceToSchemaMapper) additionalToSchema(src *sdk.PostgresInstance, dst map[string]any) {
	dst["is_ha"] = src.IsHighlyAvailable
}

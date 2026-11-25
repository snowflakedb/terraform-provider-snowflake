package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (failoverGroupToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"object_types": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
		"allowed_integration_types": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
	}
}

func (failoverGroupToSchemaMapper) additionalToSchema(src *sdk.FailoverGroup, dst map[string]any) {
	dst["object_types"] = collections.Map(src.ObjectTypes, func(v sdk.FailoverGroupObjectType) string {
		return string(v)
	})
	dst["allowed_integration_types"] = collections.Map(src.AllowedIntegrationTypes, func(v sdk.FailoverGroupIntegrationType) string {
		return string(v)
	})
}

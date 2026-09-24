package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (serviceDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"external_access_integrations": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
	}
}

func (serviceDetailsToSchemaMapper) additionalToSchema(src *sdk.ServiceDetails, dst map[string]any) {
	// Identifier slice: map each integration to Name() (not FullyQualifiedName).
	dst["external_access_integrations"] = collections.Map(src.ExternalAccessIntegrations, sdk.AccountObjectIdentifier.Name)
}

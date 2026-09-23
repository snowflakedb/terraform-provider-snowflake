package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (connectionToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"failover_allowed_to_accounts": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
	}
}

func (connectionToSchemaMapper) additionalToSchema(src *sdk.Connection, dst map[string]any) {
	// Identifier slice: map each account to Name() (not FullyQualifiedName).
	dst["failover_allowed_to_accounts"] = collections.Map(src.FailoverAllowedToAccounts, sdk.AccountIdentifier.Name)
}

package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (secretToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// oauth_scopes is []string in the SDK; public show_output is a TypeSet.
		"oauth_scopes": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
	}
}

func (secretToSchemaMapper) additionalToSchema(src *sdk.Secret, dst map[string]any) {
	dst["oauth_scopes"] = src.OauthScopes
}

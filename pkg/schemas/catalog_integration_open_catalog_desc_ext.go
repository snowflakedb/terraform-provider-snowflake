package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (catalogIntegrationOpenCatalogDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rest_config": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"catalog_uri": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"catalog_api_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"catalog_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"access_delegation_mode": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"rest_authentication": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"oauth_token_uri": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"oauth_client_id": {
						Type:      schema.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"oauth_allowed_scopes": {
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
	}
}

func (catalogIntegrationOpenCatalogDetailsToSchemaMapper) additionalToSchema(src *sdk.CatalogIntegrationOpenCatalogDetails, dst map[string]any) {
	dst["rest_config"] = []map[string]any{
		{
			"catalog_uri":            src.RestConfig.CatalogUri,
			"catalog_api_type":       src.RestConfig.CatalogApiType,
			"catalog_name":           src.RestConfig.CatalogName,
			"access_delegation_mode": src.RestConfig.AccessDelegationMode,
		},
	}
	dst["rest_authentication"] = []map[string]any{
		{
			"oauth_token_uri":      src.RestAuthentication.OauthTokenUri,
			"oauth_client_id":      src.RestAuthentication.OauthClientId,
			"oauth_allowed_scopes": src.RestAuthentication.OauthAllowedScopes,
		},
	}
}

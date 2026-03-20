package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DescribeCatalogIntegrationOpenCatalogDetailsSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"catalog_source": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"table_format": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"enabled": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"refresh_interval_seconds": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"catalog_namespace": {
		Type:     schema.TypeString,
		Computed: true,
	},
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
					Type:     schema.TypeString,
					Computed: true,
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

var _ = DescribeCatalogIntegrationOpenCatalogDetailsSchema

func CatalogIntegrationOpenCatalogDetailsToSchema(catalogIntegrationOpenCatalogDetails *sdk.CatalogIntegrationOpenCatalogDetails) map[string]any {
	catalogIntegrationOpenCatalogDetailsSchema := make(map[string]any)
	catalogIntegrationOpenCatalogDetailsSchema["id"] = catalogIntegrationOpenCatalogDetails.Id.Name()
	catalogIntegrationOpenCatalogDetailsSchema["catalog_source"] = string(catalogIntegrationOpenCatalogDetails.CatalogSource)
	catalogIntegrationOpenCatalogDetailsSchema["table_format"] = string(catalogIntegrationOpenCatalogDetails.TableFormat)
	catalogIntegrationOpenCatalogDetailsSchema["enabled"] = catalogIntegrationOpenCatalogDetails.Enabled
	catalogIntegrationOpenCatalogDetailsSchema["refresh_interval_seconds"] = catalogIntegrationOpenCatalogDetails.RefreshIntervalSeconds
	catalogIntegrationOpenCatalogDetailsSchema["comment"] = catalogIntegrationOpenCatalogDetails.Comment
	catalogIntegrationOpenCatalogDetailsSchema["catalog_namespace"] = catalogIntegrationOpenCatalogDetails.CatalogNamespace
	catalogIntegrationOpenCatalogDetailsSchema["rest_config"] = []map[string]any{
		{
			"catalog_uri":            catalogIntegrationOpenCatalogDetails.RestConfig.CatalogUri,
			"catalog_api_type":       catalogIntegrationOpenCatalogDetails.RestConfig.CatalogApiType,
			"catalog_name":           catalogIntegrationOpenCatalogDetails.RestConfig.CatalogName,
			"access_delegation_mode": catalogIntegrationOpenCatalogDetails.RestConfig.AccessDelegationMode,
		},
	}
	catalogIntegrationOpenCatalogDetailsSchema["rest_authentication"] = []map[string]any{
		{
			"oauth_token_uri":      catalogIntegrationOpenCatalogDetails.RestAuthentication.OauthTokenUri,
			"oauth_client_id":      catalogIntegrationOpenCatalogDetails.RestAuthentication.OauthClientId,
			"oauth_allowed_scopes": catalogIntegrationOpenCatalogDetails.RestAuthentication.OauthAllowedScopes,
		},
	}
	return catalogIntegrationOpenCatalogDetailsSchema
}

var _ = CatalogIntegrationOpenCatalogDetailsToSchema

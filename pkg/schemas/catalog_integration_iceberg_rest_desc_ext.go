package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (catalogIntegrationIcebergRestDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
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
					"prefix": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"catalog_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"catalog_api_type": {
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
		"oauth_rest_authentication": {
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
		"bearer_rest_authentication": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{},
			},
		},
		"sigv4_rest_authentication": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"sigv4_iam_role": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"sigv4_signing_region": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (catalogIntegrationIcebergRestDetailsToSchemaMapper) additionalToSchema(src *sdk.CatalogIntegrationIcebergRestDetails, dst map[string]any) {
	dst["rest_config"] = []map[string]any{
		{
			"catalog_uri":            src.RestConfig.CatalogUri,
			"prefix":                 src.RestConfig.Prefix,
			"catalog_name":           src.RestConfig.CatalogName,
			"catalog_api_type":       string(src.RestConfig.CatalogApiType),
			"access_delegation_mode": string(src.RestConfig.AccessDelegationMode),
		},
	}
	if src.OAuthRestAuthentication != nil {
		dst["oauth_rest_authentication"] = []map[string]any{
			{
				"oauth_token_uri":      src.OAuthRestAuthentication.OauthTokenUri,
				"oauth_client_id":      src.OAuthRestAuthentication.OauthClientId,
				"oauth_allowed_scopes": src.OAuthRestAuthentication.OauthAllowedScopes,
			},
		}
	} else {
		dst["oauth_rest_authentication"] = []map[string]any{}
	}
	dst["bearer_rest_authentication"] = []map[string]any{}
	if src.SigV4RestAuthentication != nil {
		dst["sigv4_rest_authentication"] = []map[string]any{
			{
				"sigv4_iam_role":       src.SigV4RestAuthentication.Sigv4IamRole,
				"sigv4_signing_region": src.SigV4RestAuthentication.Sigv4SigningRegion,
			},
		}
	} else {
		dst["sigv4_rest_authentication"] = []map[string]any{}
	}
}

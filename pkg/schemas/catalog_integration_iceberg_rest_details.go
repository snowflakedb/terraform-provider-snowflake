package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DescribeCatalogIntegrationIcebergRestDetailsSchema = map[string]*schema.Schema{
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

var _ = DescribeCatalogIntegrationIcebergRestDetailsSchema

func CatalogIntegrationIcebergRestDetailsToSchema(details *sdk.CatalogIntegrationIcebergRestDetails) map[string]any {
	out := map[string]any{
		"id":                       details.Id.Name(),
		"catalog_source":           string(details.CatalogSource),
		"table_format":             string(details.TableFormat),
		"enabled":                  details.Enabled,
		"refresh_interval_seconds": details.RefreshIntervalSeconds,
		"comment":                  details.Comment,
		"catalog_namespace":        details.CatalogNamespace,
		"rest_config": []map[string]any{{
			"catalog_uri":            details.RestConfig.CatalogUri,
			"prefix":                 details.RestConfig.Prefix,
			"catalog_name":           details.RestConfig.CatalogName,
			"catalog_api_type":       string(details.RestConfig.CatalogApiType),
			"access_delegation_mode": string(details.RestConfig.AccessDelegationMode),
		}},
	}
	if details.OAuthRestAuthentication != nil {
		out["oauth_rest_authentication"] = []map[string]any{{
			"oauth_token_uri":      details.OAuthRestAuthentication.OauthTokenUri,
			"oauth_client_id":      details.OAuthRestAuthentication.OauthClientId,
			"oauth_allowed_scopes": details.OAuthRestAuthentication.OauthAllowedScopes,
		}}
	} else {
		out["oauth_rest_authentication"] = []map[string]any{}
	}
	out["bearer_rest_authentication"] = []map[string]any{}
	if details.SigV4RestAuthentication != nil {
		out["sigv4_rest_authentication"] = []map[string]any{{
			"sigv4_iam_role":       details.SigV4RestAuthentication.Sigv4IamRole,
			"sigv4_signing_region": details.SigV4RestAuthentication.Sigv4SigningRegion,
		}}
	} else {
		out["sigv4_rest_authentication"] = []map[string]any{}
	}
	return out
}

var _ = CatalogIntegrationIcebergRestDetailsToSchema

package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeSecretSchema represents output of DESCRIBE query for the single secret.
var DescribeSecretSchema = map[string]*schema.Schema{
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"database_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"schema_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"secret_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"username": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"oauth_access_token_expiry_time": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"oauth_refresh_token_expiry_time": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"oauth_scopes": {
		Type:     schema.TypeSet,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Computed: true,
	},
	"integration_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func SecretDescriptionToSchema(details sdk.SecretDetails) map[string]any {
	s := map[string]any{
		"created_on":       details.CreatedOn.String(),
		"name":             details.Name,
		"database_name":    details.DatabaseName,
		"schema_name":      details.SchemaName,
		"owner":            details.Owner,
		"comment":          details.Comment,
		"secret_type":      details.SecretType,
		"username":         details.Username,
		"oauth_scopes":     details.OauthScopes,
		"integration_name": details.IntegrationName,
	}
	if details.OauthAccessTokenExpiryTime != nil {
		s["oauth_access_token_expiry_time"] = details.OauthAccessTokenExpiryTime.String()
	}
	if details.OauthRefreshTokenExpiryTime != nil {
		s["oauth_refresh_token_expiry_time"] = details.OauthRefreshTokenExpiryTime.String()
	}
	return s
}

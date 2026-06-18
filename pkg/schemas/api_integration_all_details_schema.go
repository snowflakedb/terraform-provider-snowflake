package schemas

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeApiIntegrationAllDetailsSchema is the Terraform schema for DESCRIBE output of any API integration type.
// Fields that do not apply to a given provider are zero-valued.
var DescribeApiIntegrationAllDetailsSchema = map[string]*schema.Schema{
	"enabled":      {Type: schema.TypeBool, Computed: true},
	"api_key":      {Type: schema.TypeString, Computed: true, Sensitive: true},
	"api_provider": {Type: schema.TypeString, Computed: true},
	// AWS
	"api_aws_role_arn":     {Type: schema.TypeString, Computed: true},
	"api_aws_iam_user_arn": {Type: schema.TypeString, Computed: true},
	"api_aws_external_id":  {Type: schema.TypeString, Computed: true},
	// Azure
	"azure_tenant_id":             {Type: schema.TypeString, Computed: true},
	"azure_ad_application_id":     {Type: schema.TypeString, Computed: true},
	"azure_multi_tenant_app_name": {Type: schema.TypeString, Computed: true},
	"azure_consent_url":           {Type: schema.TypeString, Computed: true},
	// Google
	"google_audience":            {Type: schema.TypeString, Computed: true},
	"google_api_service_account": {Type: schema.TypeString, Computed: true},
	// Git HTTPS / External MCP
	"allowed_authentication_secrets": {Type: schema.TypeString, Computed: true},
	"user_auth_type":                 {Type: schema.TypeString, Computed: true},
	"oauth_grant":                    {Type: schema.TypeString, Computed: true},
	"oauth_client_id":                {Type: schema.TypeString, Computed: true},
	"oauth_client_auth_method":       {Type: schema.TypeString, Computed: true},
	"oauth_token_endpoint":           {Type: schema.TypeString, Computed: true},
	"oauth_authorization_endpoint":   {Type: schema.TypeString, Computed: true},
	"oauth_access_token_validity":    {Type: schema.TypeInt, Computed: true},
	"oauth_refresh_token_validity":   {Type: schema.TypeInt, Computed: true},
	"oauth_allowed_scopes":           {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}, Computed: true},
	"oauth_username":                 {Type: schema.TypeString, Computed: true},
	"oauth_assertion_issuer":         {Type: schema.TypeString, Computed: true},
	"oauth_resource_url":             {Type: schema.TypeString, Computed: true},
	"use_privatelink_endpoint":       {Type: schema.TypeBool, Computed: true},
	"tls_trusted_certificates":       {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}, Computed: true},
	// Shared
	"allowed_prefixes": {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}, Computed: true},
	"blocked_prefixes":  {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}, Computed: true},
	"comment":           {Type: schema.TypeString, Computed: true},
}

func ApiIntegrationAllDetailsToSchema(d *sdk.ApiIntegrationAllDetails) map[string]any {
	result := make(map[string]any)
	result["enabled"] = d.Enabled
	result["api_key"] = d.ApiKey
	result["api_provider"] = strings.ToLower(d.ApiProvider)
	result["api_aws_role_arn"] = d.ApiAwsRoleArn
	result["api_aws_iam_user_arn"] = d.ApiAwsIamUserArn
	result["api_aws_external_id"] = d.ApiAwsExternalId
	result["azure_tenant_id"] = d.AzureTenantId
	result["azure_ad_application_id"] = d.AzureAdApplicationId
	result["azure_multi_tenant_app_name"] = d.AzureMultiTenantAppName
	result["azure_consent_url"] = d.AzureConsentUrl
	result["google_audience"] = d.GoogleAudience
	result["google_api_service_account"] = d.GoogleApiServiceAccount
	result["allowed_authentication_secrets"] = d.AllowedAuthenticationSecrets
	result["user_auth_type"] = d.UserAuthType
	result["oauth_grant"] = d.OauthGrant
	result["oauth_client_id"] = d.OauthClientId
	result["oauth_client_auth_method"] = d.OauthClientAuthMethod
	result["oauth_token_endpoint"] = d.OauthTokenEndpoint
	result["oauth_authorization_endpoint"] = d.OauthAuthorizationEndpoint
	result["oauth_access_token_validity"] = d.OauthAccessTokenValidity
	result["oauth_refresh_token_validity"] = d.OauthRefreshTokenValidity
	result["oauth_allowed_scopes"] = d.OauthAllowedScopes
	result["oauth_username"] = d.OauthUsername
	result["oauth_assertion_issuer"] = d.OauthAssertionIssuer
	result["oauth_resource_url"] = d.OauthResourceUrl
	result["use_privatelink_endpoint"] = d.UsePrivatelinkEndpoint
	result["tls_trusted_certificates"] = d.TlsTrustedCertificates
	result["allowed_prefixes"] = d.AllowedPrefixes
	result["blocked_prefixes"] = d.BlockedPrefixes
	result["comment"] = d.Comment
	return result
}

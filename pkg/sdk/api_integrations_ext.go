package sdk

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (r *CreateApiIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

// ApiIntegrationAwsDetails holds the structured output of DESCRIBE API INTEGRATION for AWS integrations.
type ApiIntegrationAwsDetails struct {
	Id               AccountObjectIdentifier
	Enabled          bool
	ApiKey           string
	ApiProvider      string
	ApiAwsRoleArn    string
	ApiAwsIamUserArn string
	ApiAwsExternalId string
	AllowedPrefixes  []string
	BlockedPrefixes  []string
	Comment          string
}

func (d *ApiIntegrationAwsDetails) ID() AccountObjectIdentifier {
	return d.Id
}

// ApiIntegrationAzureDetails holds the structured output of DESCRIBE API INTEGRATION for Azure integrations.
type ApiIntegrationAzureDetails struct {
	Id                      AccountObjectIdentifier
	Enabled                 bool
	ApiKey                  string
	ApiProvider             string
	AzureTenantId           string
	AzureAdApplicationId    string
	AzureMultiTenantAppName string
	AzureConsentUrl         string
	AllowedPrefixes         []string
	BlockedPrefixes         []string
	Comment                 string
}

func (d *ApiIntegrationAzureDetails) ID() AccountObjectIdentifier {
	return d.Id
}

// ApiIntegrationGoogleDetails holds the structured output of DESCRIBE API INTEGRATION for Google integrations.
type ApiIntegrationGoogleDetails struct {
	Id                      AccountObjectIdentifier
	Enabled                 bool
	ApiKey                  string
	ApiProvider             string
	GoogleAudience          string
	GoogleApiServiceAccount string
	AllowedPrefixes         []string
	BlockedPrefixes         []string
	Comment                 string
}

func (d *ApiIntegrationGoogleDetails) ID() AccountObjectIdentifier {
	return d.Id
}

// ApiIntegrationGitHttpsApiDetails holds the structured output of DESCRIBE API INTEGRATION for
// git HTTPS API integrations (covers token-based, GitHub App, OAuth2, and private-link variants).
type ApiIntegrationGitHttpsApiDetails struct {
	Id                           AccountObjectIdentifier
	Enabled                      bool
	ApiProvider                  string
	AllowedAuthenticationSecrets string // "ALL", "NONE", or comma-separated secret identifiers
	UserAuthType                 string
	OauthGrant                   string
	OauthClientId                string
	OauthClientAuthMethod        string
	OauthTokenEndpoint           string
	OauthAuthorizationEndpoint   string
	OauthAccessTokenValidity     int
	OauthRefreshTokenValidity    int
	OauthAllowedScopes           []string
	OauthUsername                string
	OauthAssertionIssuer         string
	OauthResourceUrl             string
	UsePrivatelinkEndpoint       bool
	TlsTrustedCertificates       []string
	AllowedPrefixes              []string
	BlockedPrefixes              []string
	Comment                      string
}

func (d *ApiIntegrationGitHttpsApiDetails) ID() AccountObjectIdentifier {
	return d.Id
}

// ApiIntegrationExternalMcpDetails holds the structured output of DESCRIBE API INTEGRATION for
// external MCP integrations (covers OAuth2 and dynamic-client variants).
type ApiIntegrationExternalMcpDetails struct {
	Id                         AccountObjectIdentifier
	Enabled                    bool
	ApiProvider                string
	UserAuthType               string
	OauthGrant                 string
	OauthClientId              string
	OauthClientAuthMethod      string
	OauthTokenEndpoint         string
	OauthAuthorizationEndpoint string
	OauthAccessTokenValidity   int
	OauthRefreshTokenValidity  int
	OauthAllowedScopes         []string
	OauthUsername              string
	OauthAssertionIssuer       string
	OauthResourceUrl           string
	AllowedPrefixes            []string
	BlockedPrefixes            []string
	Comment                    string
}

func (d *ApiIntegrationExternalMcpDetails) ID() AccountObjectIdentifier {
	return d.Id
}

// DescribeAwsDetails fetches and parses describe output for an AWS API integration.
func (v *apiIntegrations) DescribeAwsDetails(ctx context.Context, id AccountObjectIdentifier) (*ApiIntegrationAwsDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseApiIntegrationAwsDetails(properties, id)
}

// DescribeAzureDetails fetches and parses describe output for an Azure API integration.
func (v *apiIntegrations) DescribeAzureDetails(ctx context.Context, id AccountObjectIdentifier) (*ApiIntegrationAzureDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseApiIntegrationAzureDetails(properties, id)
}

// DescribeGoogleDetails fetches and parses describe output for a Google API integration.
func (v *apiIntegrations) DescribeGoogleDetails(ctx context.Context, id AccountObjectIdentifier) (*ApiIntegrationGoogleDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseApiIntegrationGoogleDetails(properties, id)
}

// DescribeGitHttpsApiDetails fetches and parses describe output for a git HTTPS API integration.
func (v *apiIntegrations) DescribeGitHttpsApiDetails(ctx context.Context, id AccountObjectIdentifier) (*ApiIntegrationGitHttpsApiDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseApiIntegrationGitHttpsApiDetails(properties, id)
}

// DescribeExternalMcpDetails fetches and parses describe output for an external MCP API integration.
func (v *apiIntegrations) DescribeExternalMcpDetails(ctx context.Context, id AccountObjectIdentifier) (*ApiIntegrationExternalMcpDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseApiIntegrationExternalMcpDetails(properties, id)
}

func parseApiIntegrationAwsDetails(properties []ApiIntegrationProperty, id AccountObjectIdentifier) (*ApiIntegrationAwsDetails, error) {
	details := &ApiIntegrationAwsDetails{Id: id}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "ENABLED":
			if v, err := strconv.ParseBool(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.Enabled = v
			}
		case "API_KEY":
			details.ApiKey = prop.Value
		case "API_PROVIDER":
			details.ApiProvider = prop.Value
		case "API_AWS_ROLE_ARN":
			details.ApiAwsRoleArn = prop.Value
		case "API_AWS_IAM_USER_ARN":
			details.ApiAwsIamUserArn = prop.Value
		case "API_AWS_EXTERNAL_ID":
			details.ApiAwsExternalId = prop.Value
		case "API_ALLOWED_PREFIXES":
			details.AllowedPrefixes = ParseCommaSeparatedStringArray(prop.Value, false)
		case "API_BLOCKED_PREFIXES":
			details.BlockedPrefixes = ParseCommaSeparatedStringArray(prop.Value, false)
		case "COMMENT":
			details.Comment = prop.Value
		}
	}
	return details, errors.Join(errs...)
}

func parseApiIntegrationAzureDetails(properties []ApiIntegrationProperty, id AccountObjectIdentifier) (*ApiIntegrationAzureDetails, error) {
	details := &ApiIntegrationAzureDetails{Id: id}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "ENABLED":
			if v, err := strconv.ParseBool(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.Enabled = v
			}
		case "API_KEY":
			details.ApiKey = prop.Value
		case "API_PROVIDER":
			details.ApiProvider = prop.Value
		case "AZURE_TENANT_ID":
			details.AzureTenantId = prop.Value
		case "AZURE_AD_APPLICATION_ID":
			details.AzureAdApplicationId = prop.Value
		case "AZURE_MULTI_TENANT_APP_NAME":
			details.AzureMultiTenantAppName = prop.Value
		case "AZURE_CONSENT_URL":
			details.AzureConsentUrl = prop.Value
		case "API_ALLOWED_PREFIXES":
			details.AllowedPrefixes = ParseCommaSeparatedStringArray(prop.Value, false)
		case "API_BLOCKED_PREFIXES":
			details.BlockedPrefixes = ParseCommaSeparatedStringArray(prop.Value, false)
		case "COMMENT":
			details.Comment = prop.Value
		}
	}
	return details, errors.Join(errs...)
}

func parseApiIntegrationGoogleDetails(properties []ApiIntegrationProperty, id AccountObjectIdentifier) (*ApiIntegrationGoogleDetails, error) {
	details := &ApiIntegrationGoogleDetails{Id: id}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "ENABLED":
			if v, err := strconv.ParseBool(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.Enabled = v
			}
		case "API_KEY":
			details.ApiKey = prop.Value
		case "API_PROVIDER":
			details.ApiProvider = prop.Value
		case "GOOGLE_AUDIENCE":
			details.GoogleAudience = prop.Value
		case "API_GCP_SERVICE_ACCOUNT":
			details.GoogleApiServiceAccount = prop.Value
		case "API_ALLOWED_PREFIXES":
			details.AllowedPrefixes = ParseCommaSeparatedStringArray(prop.Value, false)
		case "API_BLOCKED_PREFIXES":
			details.BlockedPrefixes = ParseCommaSeparatedStringArray(prop.Value, false)
		case "COMMENT":
			details.Comment = prop.Value
		}
	}
	return details, errors.Join(errs...)
}

func parseApiIntegrationGitHttpsApiDetails(properties []ApiIntegrationProperty, id AccountObjectIdentifier) (*ApiIntegrationGitHttpsApiDetails, error) {
	details := &ApiIntegrationGitHttpsApiDetails{Id: id}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "ENABLED":
			if v, err := strconv.ParseBool(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.Enabled = v
			}
		case "API_PROVIDER":
			details.ApiProvider = prop.Value
		case "ALLOWED_AUTHENTICATION_SECRETS":
			details.AllowedAuthenticationSecrets = prop.Value
		case "API_USER_AUTHENTICATION":
			if err := parseUserAuthIntoGitHttpsApi(prop.Value, details); err != nil {
				errs = append(errs, err)
			}
		case "USE_PRIVATELINK_ENDPOINT":
			if v, err := strconv.ParseBool(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.UsePrivatelinkEndpoint = v
			}
		case "TLS_TRUSTED_CERTIFICATES":
			details.TlsTrustedCertificates = ParseCommaSeparatedStringArray(prop.Value, false)
		case "API_ALLOWED_PREFIXES":
			details.AllowedPrefixes = ParseCommaSeparatedStringArray(prop.Value, false)
		case "API_BLOCKED_PREFIXES":
			details.BlockedPrefixes = ParseCommaSeparatedStringArray(prop.Value, false)
		case "COMMENT":
			details.Comment = prop.Value
		}
	}
	return details, errors.Join(errs...)
}

func parseApiIntegrationExternalMcpDetails(properties []ApiIntegrationProperty, id AccountObjectIdentifier) (*ApiIntegrationExternalMcpDetails, error) {
	details := &ApiIntegrationExternalMcpDetails{Id: id}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "ENABLED":
			if v, err := strconv.ParseBool(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.Enabled = v
			}
		case "API_PROVIDER":
			details.ApiProvider = prop.Value
		case "API_USER_AUTHENTICATION":
			if err := parseUserAuthIntoExternalMcp(prop.Value, details); err != nil {
				errs = append(errs, err)
			}
		case "API_ALLOWED_PREFIXES":
			details.AllowedPrefixes = ParseCommaSeparatedStringArray(prop.Value, false)
		case "API_BLOCKED_PREFIXES":
			details.BlockedPrefixes = ParseCommaSeparatedStringArray(prop.Value, false)
		case "COMMENT":
			details.Comment = prop.Value
		}
	}
	return details, errors.Join(errs...)
}

func parseUserAuthIntoGitHttpsApi(value string, details *ApiIntegrationGitHttpsApiDetails) error {
	s := strings.TrimPrefix(value, "{")
	s = strings.TrimSuffix(s, "}")
	parts := ParseOuterCommaSeparatedStringArray(fmt.Sprintf("[%s]", s), false)
	var errs []error
	for _, part := range parts {
		k, v, _ := strings.Cut(part, "=")
		switch k {
		case "TYPE":
			details.UserAuthType = v
		case "OAUTH_GRANT":
			details.OauthGrant = emptyIfNull(v)
		case "OAUTH_CLIENT_ID":
			details.OauthClientId = v
		case "OAUTH_CLIENT_AUTH_METHOD":
			details.OauthClientAuthMethod = emptyIfNull(v)
		case "OAUTH_TOKEN_ENDPOINT":
			details.OauthTokenEndpoint = v
		case "OAUTH_AUTHORIZATION_ENDPOINT":
			details.OauthAuthorizationEndpoint = v
		case "OAUTH_ACCESS_TOKEN_VALIDITY":
			if val, err := strconv.ParseInt(v, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.OauthAccessTokenValidity = int(val)
			}
		case "OAUTH_REFRESH_TOKEN_VALIDITY":
			if val, err := strconv.ParseInt(v, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.OauthRefreshTokenValidity = int(val)
			}
		case "OAUTH_ALLOWED_SCOPES":
			details.OauthAllowedScopes = ParseCommaSeparatedStringArray(emptyIfNull(v), false)
		case "OAUTH_USERNAME":
			details.OauthUsername = emptyIfNull(v)
		case "OAUTH_ASSERTION_ISSUER":
			details.OauthAssertionIssuer = emptyIfNull(v)
		case "OAUTH_RESOURCE_URL":
			details.OauthResourceUrl = emptyIfNull(v)
		}
	}
	return errors.Join(errs...)
}

func parseUserAuthIntoExternalMcp(value string, details *ApiIntegrationExternalMcpDetails) error {
	s := strings.TrimPrefix(value, "{")
	s = strings.TrimSuffix(s, "}")
	parts := ParseOuterCommaSeparatedStringArray(fmt.Sprintf("[%s]", s), false)
	var errs []error
	for _, part := range parts {
		k, v, _ := strings.Cut(part, "=")
		switch k {
		case "TYPE":
			details.UserAuthType = v
		case "OAUTH_GRANT":
			details.OauthGrant = emptyIfNull(v)
		case "OAUTH_CLIENT_ID":
			details.OauthClientId = v
		case "OAUTH_CLIENT_AUTH_METHOD":
			details.OauthClientAuthMethod = emptyIfNull(v)
		case "OAUTH_TOKEN_ENDPOINT":
			details.OauthTokenEndpoint = v
		case "OAUTH_AUTHORIZATION_ENDPOINT":
			details.OauthAuthorizationEndpoint = v
		case "OAUTH_ACCESS_TOKEN_VALIDITY":
			if val, err := strconv.ParseInt(v, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.OauthAccessTokenValidity = int(val)
			}
		case "OAUTH_REFRESH_TOKEN_VALIDITY":
			if val, err := strconv.ParseInt(v, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.OauthRefreshTokenValidity = int(val)
			}
		case "OAUTH_ALLOWED_SCOPES":
			details.OauthAllowedScopes = ParseCommaSeparatedStringArray(emptyIfNull(v), false)
		case "OAUTH_USERNAME":
			details.OauthUsername = emptyIfNull(v)
		case "OAUTH_ASSERTION_ISSUER":
			details.OauthAssertionIssuer = emptyIfNull(v)
		case "OAUTH_RESOURCE_URL":
			details.OauthResourceUrl = emptyIfNull(v)
		}
	}
	return errors.Join(errs...)
}

// ApiIntegrationAllDetails holds the output of DESCRIBE API INTEGRATION for any provider type.
// Fields that do not apply to a given provider are zero-valued.
type ApiIntegrationAllDetails struct {
	Id                           AccountObjectIdentifier
	Enabled                      bool
	ApiKey                       string
	ApiProvider                  string
	ApiAwsRoleArn                string
	ApiAwsIamUserArn             string
	ApiAwsExternalId             string
	AzureTenantId                string
	AzureAdApplicationId         string
	AzureMultiTenantAppName      string
	AzureConsentUrl              string
	GoogleAudience               string
	GoogleApiServiceAccount      string
	AllowedAuthenticationSecrets string
	UserAuthType                 string
	OauthGrant                   string
	OauthClientId                string
	OauthClientAuthMethod        string
	OauthTokenEndpoint           string
	OauthAuthorizationEndpoint   string
	OauthAccessTokenValidity     int
	OauthRefreshTokenValidity    int
	OauthAllowedScopes           []string
	OauthUsername                string
	OauthAssertionIssuer         string
	OauthResourceUrl             string
	UsePrivatelinkEndpoint       bool
	TlsTrustedCertificates       []string
	AllowedPrefixes              []string
	BlockedPrefixes              []string
	Comment                      string
}

func (d *ApiIntegrationAllDetails) ID() AccountObjectIdentifier {
	return d.Id
}

// DescribeAllDetails fetches and parses describe output for any API integration type.
func (v *apiIntegrations) DescribeAllDetails(ctx context.Context, id AccountObjectIdentifier) (*ApiIntegrationAllDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseApiIntegrationAllDetails(properties, id)
}

func parseApiIntegrationAllDetails(properties []ApiIntegrationProperty, id AccountObjectIdentifier) (*ApiIntegrationAllDetails, error) {
	details := &ApiIntegrationAllDetails{Id: id}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "ENABLED":
			if v, err := strconv.ParseBool(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.Enabled = v
			}
		case "API_KEY":
			details.ApiKey = prop.Value
		case "API_PROVIDER":
			details.ApiProvider = prop.Value
		case "API_AWS_ROLE_ARN":
			details.ApiAwsRoleArn = prop.Value
		case "API_AWS_IAM_USER_ARN":
			details.ApiAwsIamUserArn = prop.Value
		case "API_AWS_EXTERNAL_ID":
			details.ApiAwsExternalId = prop.Value
		case "AZURE_TENANT_ID":
			details.AzureTenantId = prop.Value
		case "AZURE_AD_APPLICATION_ID":
			details.AzureAdApplicationId = prop.Value
		case "AZURE_MULTI_TENANT_APP_NAME":
			details.AzureMultiTenantAppName = prop.Value
		case "AZURE_CONSENT_URL":
			details.AzureConsentUrl = prop.Value
		case "GOOGLE_AUDIENCE":
			details.GoogleAudience = prop.Value
		case "API_GCP_SERVICE_ACCOUNT":
			details.GoogleApiServiceAccount = prop.Value
		case "ALLOWED_AUTHENTICATION_SECRETS":
			details.AllowedAuthenticationSecrets = prop.Value
		case "API_USER_AUTHENTICATION":
			if err := parseUserAuthIntoAllDetails(prop.Value, details); err != nil {
				errs = append(errs, err)
			}
		case "USE_PRIVATELINK_ENDPOINT":
			if v, err := strconv.ParseBool(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.UsePrivatelinkEndpoint = v
			}
		case "TLS_TRUSTED_CERTIFICATES":
			details.TlsTrustedCertificates = ParseCommaSeparatedStringArray(prop.Value, false)
		case "API_ALLOWED_PREFIXES":
			details.AllowedPrefixes = ParseCommaSeparatedStringArray(prop.Value, false)
		case "API_BLOCKED_PREFIXES":
			details.BlockedPrefixes = ParseCommaSeparatedStringArray(prop.Value, false)
		case "COMMENT":
			details.Comment = prop.Value
		}
	}
	return details, errors.Join(errs...)
}

func parseUserAuthIntoAllDetails(value string, details *ApiIntegrationAllDetails) error {
	s := strings.TrimPrefix(value, "{")
	s = strings.TrimSuffix(s, "}")
	parts := ParseOuterCommaSeparatedStringArray(fmt.Sprintf("[%s]", s), false)
	var errs []error
	for _, part := range parts {
		k, v, _ := strings.Cut(part, "=")
		switch k {
		case "TYPE":
			details.UserAuthType = v
		case "OAUTH_GRANT":
			details.OauthGrant = emptyIfNull(v)
		case "OAUTH_CLIENT_ID":
			details.OauthClientId = v
		case "OAUTH_CLIENT_AUTH_METHOD":
			details.OauthClientAuthMethod = emptyIfNull(v)
		case "OAUTH_TOKEN_ENDPOINT":
			details.OauthTokenEndpoint = v
		case "OAUTH_AUTHORIZATION_ENDPOINT":
			details.OauthAuthorizationEndpoint = v
		case "OAUTH_ACCESS_TOKEN_VALIDITY":
			if val, err := strconv.ParseInt(v, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.OauthAccessTokenValidity = int(val)
			}
		case "OAUTH_REFRESH_TOKEN_VALIDITY":
			if val, err := strconv.ParseInt(v, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.OauthRefreshTokenValidity = int(val)
			}
		case "OAUTH_ALLOWED_SCOPES":
			details.OauthAllowedScopes = ParseCommaSeparatedStringArray(v, false)
		case "OAUTH_USERNAME":
			details.OauthUsername = emptyIfNull(v)
		case "OAUTH_ASSERTION_ISSUER":
			details.OauthAssertionIssuer = emptyIfNull(v)
		case "OAUTH_RESOURCE_URL":
			details.OauthResourceUrl = emptyIfNull(v)
		}
	}
	return errors.Join(errs...)
}

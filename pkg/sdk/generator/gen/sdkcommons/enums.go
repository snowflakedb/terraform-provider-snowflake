package sdkcommons

type (
	AllowedProviderOption                    string
	ApiIntegrationAwsApiProviderType         string
	AuthenticationMethodsOption              string
	AutoEventLogging                         string
	ClientTypesOption                        string
	ComputePoolInstanceFamily                string
	DataMetricFunctionRefEntityDomainOption  string
	DataType                                 string
	EnforceMfaOnExternalAuthenticationOption string
	ListingState                             string
	ListingRevision                          string
	Location                                 string
	LogLevel                                 string
	MetricLevel                              string
	MfaAuthenticationMethodsOption           string
	MfaEnrollmentOption                      string
	MfaPolicyAllowedMethodsOption            string
	NetworkPolicyEvaluationOption            string
	NetworkRuleMode                          string
	NetworkRuleType                          string
	NullInputBehavior                        string
	OrganizationAccountEdition               string
	ReturnNullValues                         string
	ReturnResultsBehavior                    string
	SecretType                               string
	TraceLevel                               string
)

// copied from SDK for now
const (
	SecretTypePassword                     SecretType = "PASSWORD"
	SecretTypeOAuth2                       SecretType = "OAUTH2"
	SecretTypeGenericString                SecretType = "GENERIC_STRING"
	SecretTypeOAuth2ClientCredentials      SecretType = "OAUTH2_CLIENT_CREDENTIALS"       // #nosec G101
	SecretTypeOAuth2AuthorizationCodeGrant SecretType = "OAUTH2_AUTHORIZATION_CODE_GRANT" // #nosec G101
)

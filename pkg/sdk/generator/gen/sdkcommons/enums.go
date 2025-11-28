package sdkcommons

type (
	AllowedProviderOption                                               string
	ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption     string
	ApiIntegrationAwsApiProviderType                                    string
	AuthenticationMethodsOption                                         string
	AutoEventLogging                                                    string
	ClientTypesOption                                                   string
	ComputePoolInstanceFamily                                           string
	DataMetricFunctionRefEntityDomainOption                             string
	DataType                                                            string
	EnforceMfaOnExternalAuthenticationOption                            string
	ExternalOauthSecurityIntegrationAnyRoleModeOption                   string
	ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption string
	ExternalOauthSecurityIntegrationTypeOption                          string
	ExternalStageAzureEncryptionOption                                  string
	ExternalStageGCSEncryptionOption                                    string
	ExternalStageS3EncryptionOption                                     string
	FileFormatType                                                      string
	FileFormatTypeOptions                                               string
	GCSEncryptionType                                                   string
	InternalStageEncryptionOption                                       string
	ListingState                                                        string
	ListingRevision                                                     string
	Location                                                            string
	LogLevel                                                            string
	MetricLevel                                                         string
	MfaAuthenticationMethodsOption                                      string
	MfaEnrollmentOption                                                 string
	MfaPolicyAllowedMethodsOption                                       string
	NetworkPolicyEvaluationOption                                       string
	NetworkRuleMode                                                     string
	NetworkRuleType                                                     string
	NullInputBehavior                                                   string
	OauthSecurityIntegrationClientOption                                string
	OauthSecurityIntegrationClientTypeOption                            string
	OauthSecurityIntegrationUseSecondaryRolesOption                     string
	OrganizationAccountEdition                                          string
	ReturnNullValues                                                    string
	ReturnResultsBehavior                                               string
	S3EncryptionType                                                    string
	S3Protocol                                                          string
	S3StorageProvider                                                   string
	Saml2SecurityIntegrationSaml2ProviderOption                         string
	Saml2SecurityIntegrationSaml2RequestedNameidFormatOption            string
	ScimSecurityIntegrationRunAsRoleOption                              string
	ScimSecurityIntegrationScimClientOption                             string
	SecretType                                                          string
	StageCopyColumnMapOption                                            string
	StorageProvider                                                     string
	TaskState                                                           string
	TraceLevel                                                          string
	ViewDataMetricScheduleStatusOperationOption                         string
)

// copied from SDK for now
const (
	SecretTypePassword                     SecretType = "PASSWORD"
	SecretTypeOAuth2                       SecretType = "OAUTH2"
	SecretTypeGenericString                SecretType = "GENERIC_STRING"
	SecretTypeOAuth2ClientCredentials      SecretType = "OAUTH2_CLIENT_CREDENTIALS"       // #nosec G101
	SecretTypeOAuth2AuthorizationCodeGrant SecretType = "OAUTH2_AUTHORIZATION_CODE_GRANT" // #nosec G101
)

const (
	StorageProviderGCS   StorageProvider = "GCS"
	StorageProviderAzure StorageProvider = "AZURE"
)

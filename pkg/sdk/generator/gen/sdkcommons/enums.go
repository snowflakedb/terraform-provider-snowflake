package sdkcommons

type (
	AllowedProviderOption                                               string
	ColumnConstraintType                                                string
	ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption     string
	ApiIntegrationAwsApiProviderType                                    string
	AuthenticationMethodsOption                                         string
	AutoEventLogging                                                    string
	AvroCompression                                                     string
	BinaryFormat                                                        string
	CatalogIntegrationAccessDelegationMode                              string
	CatalogIntegrationCatalogApiType                                    string
	CatalogIntegrationCatalogSourceType                                 string
	CatalogIntegrationRestAuthenticationType                            string
	CatalogIntegrationTableFormat                                       string
	ClientTypesOption                                                   string
	ComputePoolInstanceFamily                                           string
	CsvCompression                                                      string
	CsvEncoding                                                         string
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
	JsonCompression                                                     string
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
	ParquetCompression                                                  string
	ReclusterState                                                      string
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
	SequenceName                                                        string
	StageCopyColumnMapOption                                            string
	StorageProvider                                                     string
	TaskState                                                           string
	TraceLevel                                                          string
	ViewDataMetricScheduleStatusOperationOption                         string
	XmlCompression                                                      string
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

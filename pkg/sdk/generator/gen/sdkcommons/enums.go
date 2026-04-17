package sdkcommons

type (
	ColumnConstraintType                                     string
	ApiIntegrationAwsApiProviderType                         string
	AutoEventLogging                                         string
	ComputePoolInstanceFamily                                string
	DataMetricFunctionRefEntityDomainOption                  string
	DataType                                                 string
	ImageRepositoryEncryptionType                            string
	ListingState                                             string
	ListingRevision                                          string
	LogLevel                                                 string
	MetricLevel                                              string
	NetworkRuleMode                                          string
	NetworkRuleType                                          string
	NullInputBehavior                                        string
	OrganizationAccountEdition                               string
	ReclusterState                                           string
	ReturnNullValues                                         string
	ReturnResultsBehavior                                    string
	S3Protocol                                               string
	Saml2SecurityIntegrationSaml2RequestedNameidFormatOption string
	SecretType                                               string
	SequenceName                                             string
	TaskState                                                string
	TraceLevel                                               string
	ViewDataMetricScheduleStatusOperationOption              string
)

// copied from SDK for now
const (
	SecretTypePassword                     SecretType = "PASSWORD"
	SecretTypeOAuth2                       SecretType = "OAUTH2"
	SecretTypeGenericString                SecretType = "GENERIC_STRING"
	SecretTypeOAuth2ClientCredentials      SecretType = "OAUTH2_CLIENT_CREDENTIALS"       // #nosec G101
	SecretTypeOAuth2AuthorizationCodeGrant SecretType = "OAUTH2_AUTHORIZATION_CODE_GRANT" // #nosec G101
)

package sdkcommons

type (
	ActivePythonProfiler                                     string
	BinaryInputFormat                                        string
	BinaryOutputFormat                                       string
	ClientTimestampTypeMapping                               string
	ColumnConstraintType                                     string
	DefaultNullOrdering                                      string
	GeographyOutputFormat                                    string
	GeometryOutputFormat                                     string
	WarehouseGeneration                                      string
	WarehouseSize                                            string
	AutoEventLogging                                         string
	ComputePoolInstanceFamily                                string
	DataType                                                 string
	ImageRepositoryEncryptionType                            string
	LogLevel                                                 string
	MetricLevel                                              string
	StorageSerializationPolicy                               string
	NullInputBehavior                                        string
	ReturnNullValues                                         string
	ReturnResultsBehavior                                    string
	S3Protocol                                               string
	Saml2SecurityIntegrationSaml2RequestedNameidFormatOption string
	SecretType                                               string
	SequenceName                                             string
	TaskState                                                string
	TimestampTypeMapping                                     string
	TraceLevel                                               string
	TransactionDefaultIsolationLevel                         string
	UnsupportedDDLAction                                     string
)

// copied from SDK for now
const (
	SecretTypePassword                     SecretType = "PASSWORD"
	SecretTypeOAuth2                       SecretType = "OAUTH2"
	SecretTypeGenericString                SecretType = "GENERIC_STRING"
	SecretTypeOAuth2ClientCredentials      SecretType = "OAUTH2_CLIENT_CREDENTIALS"       // #nosec G101
	SecretTypeOAuth2AuthorizationCodeGrant SecretType = "OAUTH2_AUTHORIZATION_CODE_GRANT" // #nosec G101
)

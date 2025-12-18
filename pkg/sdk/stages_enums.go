package sdk

type InternalStageEncryptionOption string

var (
	InternalStageEncryptionFull InternalStageEncryptionOption = "SNOWFLAKE_FULL"
	InternalStageEncryptionSSE  InternalStageEncryptionOption = "SNOWFLAKE_SSE"
)

type ExternalStageS3EncryptionOption string

var (
	ExternalStageS3EncryptionCSE    ExternalStageS3EncryptionOption = "AWS_CSE"
	ExternalStageS3EncryptionSSES3  ExternalStageS3EncryptionOption = "AWS_SSE_S3"
	ExternalStageS3EncryptionSSEKMS ExternalStageS3EncryptionOption = "AWS_SSE_KMS"
	ExternalStageS3EncryptionNone   ExternalStageS3EncryptionOption = "NONE"
)

type ExternalStageGCSEncryptionOption string

var (
	ExternalStageGCSEncryptionSSEKMS ExternalStageGCSEncryptionOption = "GCS_SSE_KMS"
	ExternalStageGCSEncryptionNone   ExternalStageGCSEncryptionOption = "NONE"
)

type ExternalStageAzureEncryptionOption string

var (
	ExternalStageAzureEncryptionCSE  ExternalStageAzureEncryptionOption = "AZURE_CSE"
	ExternalStageAzureEncryptionNone ExternalStageAzureEncryptionOption = "NONE"
)

type StageCopyColumnMapOption string

var (
	StageCopyColumnMapCaseSensitive   StageCopyColumnMapOption = "CASE_SENSITIVE"
	StageCopyColumnMapCaseInsensitive StageCopyColumnMapOption = "CASE_INSENSITIVE"
	StageCopyColumnMapCaseNone        StageCopyColumnMapOption = "NONE"
)

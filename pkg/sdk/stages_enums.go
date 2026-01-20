package sdk

type InternalStageEncryptionOption string

var (
	InternalStageEncryptionFull InternalStageEncryptionOption = "SNOWFLAKE_FULL"
	InternalStageEncryptionSSE  InternalStageEncryptionOption = "SNOWFLAKE_SSE"
)

type ExternalStageS3EncryptionOption string

var (
	ExternalStageS3EncryptionCSE       ExternalStageS3EncryptionOption = "AWS_CSE"
	ExternalStageS3EncryptionSSES3     ExternalStageS3EncryptionOption = "AWS_SSE_S3"
	ExternalStageS3EncryptionSSEKMS    ExternalStageS3EncryptionOption = "AWS_SSE_KMS"
	ExternalStageS3EncryptionNoneValue ExternalStageS3EncryptionOption = "NONE"
)

type ExternalStageGCSEncryptionOption string

var (
	ExternalStageGCSEncryptionSSEKMS    ExternalStageGCSEncryptionOption = "GCS_SSE_KMS"
	ExternalStageGCSEncryptionNoneValue ExternalStageGCSEncryptionOption = "NONE"
)

type ExternalStageAzureEncryptionOption string

var (
	ExternalStageAzureEncryptionCSE       ExternalStageAzureEncryptionOption = "AZURE_CSE"
	ExternalStageAzureEncryptionNoneValue ExternalStageAzureEncryptionOption = "NONE"
)

// TODO: move to tables
type StageCopyColumnMapOption string

var (
	StageCopyColumnMapCaseSensitive   StageCopyColumnMapOption = "CASE_SENSITIVE"
	StageCopyColumnMapCaseInsensitive StageCopyColumnMapOption = "CASE_INSENSITIVE"
	StageCopyColumnMapCaseNone        StageCopyColumnMapOption = "NONE"
)

type StageFileFormatAvroCompression string
type StageFileFormatJsonCompression string
type StageFileFormatCsvCompression string
type StageFileFormatParquetCompression string
type StageFileFormatXmlCompression string
type StageFileFormatBinaryFormat string

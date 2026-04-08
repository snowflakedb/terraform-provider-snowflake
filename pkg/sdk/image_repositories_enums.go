package sdk

type ImageRepositoryEncryptionType string

var (
	ImageRepositoryEncryptionTypeSnowflakeFull ImageRepositoryEncryptionType = "SNOWFLAKE_FULL"
	ImageRepositoryEncryptionTypeSnowflakeSse  ImageRepositoryEncryptionType = "SNOWFLAKE_SSE"
)

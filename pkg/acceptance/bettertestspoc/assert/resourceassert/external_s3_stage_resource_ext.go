package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// ExternalS3StageDirectoryTableAssert is used for S3 stage directory table assertions.
// S3 stages don't have notification_integration like Azure stages.
type ExternalS3StageDirectoryTableAssert struct {
	Enable          bool
	RefreshOnCreate *bool
	AutoRefresh     *string
	AwsSnsTopic     *string
}

func (e *ExternalS3StageResourceAssert) HasDirectory(opts ExternalS3StageDirectoryTableAssert) *ExternalS3StageResourceAssert {
	var refreshOnCreate string
	if opts.RefreshOnCreate != nil {
		refreshOnCreate = strconv.FormatBool(*opts.RefreshOnCreate)
	}
	var autoRefresh string
	if opts.AutoRefresh != nil {
		autoRefresh = *opts.AutoRefresh
	}
	var awsSnsTopic string
	if opts.AwsSnsTopic != nil {
		awsSnsTopic = *opts.AwsSnsTopic
	}
	e.ValueSet("directory.#", "1")
	e.ValueSet("directory.0.enable", strconv.FormatBool(opts.Enable))
	e.ValueSet("directory.0.auto_refresh", autoRefresh)
	e.ValueSet("directory.0.refresh_on_create", refreshOnCreate)
	e.ValueSet("directory.0.aws_sns_topic", awsSnsTopic)
	return e
}

func (e *ExternalS3StageResourceAssert) HasEncryptionAwsCse() *ExternalS3StageResourceAssert {
	e.ValueSet("encryption.#", "1")
	e.ValueSet("encryption.0.aws_cse.#", "1")
	e.ValueSet("encryption.0.aws_sse_s3.#", "0")
	e.ValueSet("encryption.0.aws_sse_kms.#", "0")
	e.ValueSet("encryption.0.none.#", "0")
	return e
}

func (e *ExternalS3StageResourceAssert) HasEncryptionAwsSseS3() *ExternalS3StageResourceAssert {
	e.ValueSet("encryption.#", "1")
	e.ValueSet("encryption.0.aws_cse.#", "0")
	e.ValueSet("encryption.0.aws_sse_s3.#", "1")
	e.ValueSet("encryption.0.aws_sse_kms.#", "0")
	e.ValueSet("encryption.0.none.#", "0")
	return e
}

func (e *ExternalS3StageResourceAssert) HasEncryptionAwsSseKms() *ExternalS3StageResourceAssert {
	e.ValueSet("encryption.#", "1")
	e.ValueSet("encryption.0.aws_cse.#", "0")
	e.ValueSet("encryption.0.aws_sse_s3.#", "0")
	e.ValueSet("encryption.0.aws_sse_kms.#", "1")
	e.ValueSet("encryption.0.none.#", "0")
	return e
}

func (e *ExternalS3StageResourceAssert) HasEncryptionNone() *ExternalS3StageResourceAssert {
	e.ValueSet("encryption.#", "1")
	e.ValueSet("encryption.0.aws_cse.#", "0")
	e.ValueSet("encryption.0.aws_sse_s3.#", "0")
	e.ValueSet("encryption.0.aws_sse_kms.#", "0")
	e.ValueSet("encryption.0.none.#", "1")
	return e
}

func (e *ExternalS3StageResourceAssert) HasStageTypeEnum(expected sdk.StageType) *ExternalS3StageResourceAssert {
	e.ValueSet("stage_type", string(expected))
	return e
}

func (e *ExternalS3StageResourceAssert) HasCloudEnum(expected sdk.StageCloud) *ExternalS3StageResourceAssert {
	e.ValueSet("cloud", string(expected))
	return e
}

func (e *ExternalS3StageResourceAssert) HasFileFormatFormatName(expected string) *ExternalS3StageResourceAssert {
	stageApplyFileFormatFormatNameChecks(e.ResourceAssert, expected)
	return e
}

func (e *ExternalS3StageResourceAssert) HasFileFormatCsv() *ExternalS3StageResourceAssert {
	stageApplyFileFormatCsvChecks(e.ResourceAssert)
	return e
}

func (e *ExternalS3StageResourceAssert) HasCredentialsAwsKey(keyId, secretKey string) *ExternalS3StageResourceAssert {
	e.ValueSet("credentials.#", "1")
	e.ValueSet("credentials.0.aws_key_id", keyId)
	e.ValueSet("credentials.0.aws_secret_key", secretKey)
	return e
}

func (e *ExternalS3StageResourceAssert) HasCredentialsAwsRole(roleArn string) *ExternalS3StageResourceAssert {
	e.ValueSet("credentials.#", "1")
	e.ValueSet("credentials.0.aws_role", roleArn)
	return e
}

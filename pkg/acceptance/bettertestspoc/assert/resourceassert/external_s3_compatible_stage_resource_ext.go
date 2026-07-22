package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (e *ExternalS3CompatibleStageResourceAssert) HasDirectory(opts sdk.StageS3CompatibleDirectoryTableOptionsRequest) *ExternalS3CompatibleStageResourceAssert {
	var refreshOnCreate string
	if opts.RefreshOnCreate != nil {
		refreshOnCreate = strconv.FormatBool(*opts.RefreshOnCreate)
	}
	var autoRefresh string
	if opts.AutoRefresh != nil {
		autoRefresh = strconv.FormatBool(*opts.AutoRefresh)
	}
	e.ValueSet("directory.#", "1")
	e.ValueSet("directory.0.enable", strconv.FormatBool(opts.Enable))
	e.ValueSet("directory.0.auto_refresh", autoRefresh)
	e.ValueSet("directory.0.refresh_on_create", refreshOnCreate)
	return e
}

func (e *ExternalS3CompatibleStageResourceAssert) HasStageTypeEnum(expected sdk.StageType) *ExternalS3CompatibleStageResourceAssert {
	e.ValueSet("stage_type", string(expected))
	return e
}

func (e *ExternalS3CompatibleStageResourceAssert) HasCloudEnum(expected sdk.StageCloud) *ExternalS3CompatibleStageResourceAssert {
	e.ValueSet("cloud", string(expected))
	return e
}

func (e *ExternalS3CompatibleStageResourceAssert) HasFileFormatFormatName(expected string) *ExternalS3CompatibleStageResourceAssert {
	stageApplyFileFormatFormatNameChecks(e.ResourceAssert, expected)
	return e
}

func (e *ExternalS3CompatibleStageResourceAssert) HasFileFormatCsv() *ExternalS3CompatibleStageResourceAssert {
	stageApplyFileFormatCsvChecks(e.ResourceAssert)
	return e
}

func (e *ExternalS3CompatibleStageResourceAssert) HasCredentials(awsKeyId string, awsSecretKey string) *ExternalS3CompatibleStageResourceAssert {
	e.ValueSet("credentials.#", "1")
	e.ValueSet("credentials.0.aws_key_id", awsKeyId)
	e.ValueSet("credentials.0.aws_secret_key", awsSecretKey)
	return e
}

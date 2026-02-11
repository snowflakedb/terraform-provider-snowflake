package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (e *ExternalS3CompatStageResourceAssert) HasDirectory(opts sdk.StageS3CommonDirectoryTableOptionsRequest) *ExternalS3CompatStageResourceAssert {
	var refreshOnCreate string
	if opts.RefreshOnCreate != nil {
		refreshOnCreate = strconv.FormatBool(*opts.RefreshOnCreate)
	}
	var autoRefresh string
	if opts.AutoRefresh != nil {
		autoRefresh = strconv.FormatBool(*opts.AutoRefresh)
	}
	e.AddAssertion(assert.ValueSet("directory.#", "1"))
	e.AddAssertion(assert.ValueSet("directory.0.enable", strconv.FormatBool(opts.Enable)))
	e.AddAssertion(assert.ValueSet("directory.0.auto_refresh", autoRefresh))
	e.AddAssertion(assert.ValueSet("directory.0.refresh_on_create", refreshOnCreate))
	return e
}

func (e *ExternalS3CompatStageResourceAssert) HasStageTypeEnum(expected sdk.StageType) *ExternalS3CompatStageResourceAssert {
	e.AddAssertion(assert.ValueSet("stage_type", string(expected)))
	return e
}

func (e *ExternalS3CompatStageResourceAssert) HasCloudEnum(expected sdk.StageCloud) *ExternalS3CompatStageResourceAssert {
	e.AddAssertion(assert.ValueSet("cloud", string(expected)))
	return e
}

func (e *ExternalS3CompatStageResourceAssert) HasFileFormatEmpty() *ExternalS3CompatStageResourceAssert {
	e.AddAssertion(assert.ValueSet("file_format.#", "0"))
	return e
}

func (e *ExternalS3CompatStageResourceAssert) HasFileFormatFormatName(expected string) *ExternalS3CompatStageResourceAssert {
	for _, a := range stageHasFileFormatFormatName(expected) {
		e.AddAssertion(a)
	}
	return e
}

func (e *ExternalS3CompatStageResourceAssert) HasFileFormatCsv() *ExternalS3CompatStageResourceAssert {
	for _, a := range stageHasFileFormatCsv() {
		e.AddAssertion(a)
	}
	return e
}

func (e *ExternalS3CompatStageResourceAssert) HasCredentials(awsKeyId string, awsSecretKey string) *ExternalS3CompatStageResourceAssert {
	e.AddAssertion(assert.ValueSet("credentials.#", "1"))
	e.AddAssertion(assert.ValueSet("credentials.0.aws_key_id", awsKeyId))
	e.AddAssertion(assert.ValueSet("credentials.0.aws_secret_key", awsSecretKey))
	return e
}

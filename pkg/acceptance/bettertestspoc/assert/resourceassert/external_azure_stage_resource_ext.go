package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ExternalStageDirectoryTableAssert struct {
	Enable                  bool
	RefreshOnCreate         *bool
	AutoRefresh             *string
	NotificationIntegration *string
}

func (e *ExternalAzureStageResourceAssert) HasDirectory(opts ExternalStageDirectoryTableAssert) *ExternalAzureStageResourceAssert {
	var notificationIntegration string
	if opts.NotificationIntegration != nil {
		notificationIntegration = *opts.NotificationIntegration
	}
	var refreshOnCreate string
	if opts.RefreshOnCreate != nil {
		refreshOnCreate = strconv.FormatBool(*opts.RefreshOnCreate)
	}
	var autoRefresh string
	if opts.AutoRefresh != nil {
		autoRefresh = *opts.AutoRefresh
	}
	e.AddAssertion(assert.ValueSet("directory.#", "1"))
	e.AddAssertion(assert.ValueSet("directory.0.enable", strconv.FormatBool(opts.Enable)))
	e.AddAssertion(assert.ValueSet("directory.0.auto_refresh", autoRefresh))
	e.AddAssertion(assert.ValueSet("directory.0.notification_integration", notificationIntegration))
	e.AddAssertion(assert.ValueSet("directory.0.refresh_on_create", refreshOnCreate))
	return e
}

func (e *ExternalAzureStageResourceAssert) HasEncryptionAzureCse() *ExternalAzureStageResourceAssert {
	e.AddAssertion(assert.ValueSet("encryption.#", "1"))
	e.AddAssertion(assert.ValueSet("encryption.0.azure_cse.#", "1"))
	e.AddAssertion(assert.ValueSet("encryption.0.none.#", "0"))
	return e
}

func (e *ExternalAzureStageResourceAssert) HasEncryptionNone() *ExternalAzureStageResourceAssert {
	e.AddAssertion(assert.ValueSet("encryption.#", "1"))
	e.AddAssertion(assert.ValueSet("encryption.0.azure_cse.#", "0"))
	e.AddAssertion(assert.ValueSet("encryption.0.none.#", "1"))
	return e
}

func (e *ExternalAzureStageResourceAssert) HasStageTypeEnum(expected sdk.StageType) *ExternalAzureStageResourceAssert {
	e.AddAssertion(assert.ValueSet("stage_type", string(expected)))
	return e
}

func (e *ExternalAzureStageResourceAssert) HasCloudEnum(expected sdk.StageCloud) *ExternalAzureStageResourceAssert {
	e.AddAssertion(assert.ValueSet("cloud", string(expected)))
	return e
}

func (e *ExternalAzureStageResourceAssert) HasFileFormatEmpty() *ExternalAzureStageResourceAssert {
	e.AddAssertion(assert.ValueSet("file_format.#", "0"))
	return e
}

func (e *ExternalAzureStageResourceAssert) HasFileFormatFormatName(expected string) *ExternalAzureStageResourceAssert {
	for _, a := range stageHasFileFormatFormatName(expected) {
		e.AddAssertion(a)
	}
	return e
}

func (e *ExternalAzureStageResourceAssert) HasFileFormatCsv() *ExternalAzureStageResourceAssert {
	for _, a := range stageHasFileFormatCsv() {
		e.AddAssertion(a)
	}
	return e
}

func (e *ExternalAzureStageResourceAssert) HasCredentials(token string) *ExternalAzureStageResourceAssert {
	e.AddAssertion(assert.ValueSet("credentials.#", "1"))
	e.AddAssertion(assert.ValueSet("credentials.0.azure_sas_token", token))
	return e
}

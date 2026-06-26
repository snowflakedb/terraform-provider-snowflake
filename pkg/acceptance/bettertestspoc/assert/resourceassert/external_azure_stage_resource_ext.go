package resourceassert

import (
	"strconv"

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
	e.ValueSet("directory.#", "1")
	e.ValueSet("directory.0.enable", strconv.FormatBool(opts.Enable))
	e.ValueSet("directory.0.auto_refresh", autoRefresh)
	e.ValueSet("directory.0.notification_integration", notificationIntegration)
	e.ValueSet("directory.0.refresh_on_create", refreshOnCreate)
	return e
}

func (e *ExternalAzureStageResourceAssert) HasEncryptionAzureCse() *ExternalAzureStageResourceAssert {
	e.ValueSet("encryption.#", "1")
	e.ValueSet("encryption.0.azure_cse.#", "1")
	e.ValueSet("encryption.0.none.#", "0")
	return e
}

func (e *ExternalAzureStageResourceAssert) HasEncryptionNone() *ExternalAzureStageResourceAssert {
	e.ValueSet("encryption.#", "1")
	e.ValueSet("encryption.0.azure_cse.#", "0")
	e.ValueSet("encryption.0.none.#", "1")
	return e
}

func (e *ExternalAzureStageResourceAssert) HasStageTypeEnum(expected sdk.StageType) *ExternalAzureStageResourceAssert {
	e.ValueSet("stage_type", string(expected))
	return e
}

func (e *ExternalAzureStageResourceAssert) HasCloudEnum(expected sdk.StageCloud) *ExternalAzureStageResourceAssert {
	e.ValueSet("cloud", string(expected))
	return e
}

func (e *ExternalAzureStageResourceAssert) HasFileFormatFormatName(expected string) *ExternalAzureStageResourceAssert {
	stageApplyFileFormatFormatNameChecks(e.ResourceAssert, expected)
	return e
}

func (e *ExternalAzureStageResourceAssert) HasFileFormatCsv() *ExternalAzureStageResourceAssert {
	stageApplyFileFormatCsvChecks(e.ResourceAssert)
	return e
}

func (e *ExternalAzureStageResourceAssert) HasCredentials(token string) *ExternalAzureStageResourceAssert {
	e.ValueSet("credentials.#", "1")
	e.ValueSet("credentials.0.azure_sas_token", token)
	return e
}

package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (e *ExternalGcsStageResourceAssert) HasDirectory(opts ExternalStageDirectoryTableAssert) *ExternalGcsStageResourceAssert {
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

func (e *ExternalGcsStageResourceAssert) HasEncryptionGcsSseKms() *ExternalGcsStageResourceAssert {
	e.ValueSet("encryption.#", "1")
	e.ValueSet("encryption.0.gcs_sse_kms.#", "1")
	e.ValueSet("encryption.0.none.#", "0")
	return e
}

func (e *ExternalGcsStageResourceAssert) HasEncryptionNone() *ExternalGcsStageResourceAssert {
	e.ValueSet("encryption.#", "1")
	e.ValueSet("encryption.0.gcs_sse_kms.#", "0")
	e.ValueSet("encryption.0.none.#", "1")
	return e
}

func (e *ExternalGcsStageResourceAssert) HasStageTypeEnum(expected sdk.StageType) *ExternalGcsStageResourceAssert {
	e.ValueSet("stage_type", string(expected))
	return e
}

func (e *ExternalGcsStageResourceAssert) HasCloudEnum(expected sdk.StageCloud) *ExternalGcsStageResourceAssert {
	e.ValueSet("cloud", string(expected))
	return e
}

func (e *ExternalGcsStageResourceAssert) HasFileFormatFormatName(expected string) *ExternalGcsStageResourceAssert {
	stageApplyFileFormatFormatNameChecks(e.ResourceAssert, expected)
	return e
}

func (e *ExternalGcsStageResourceAssert) HasFileFormatCsv() *ExternalGcsStageResourceAssert {
	stageApplyFileFormatCsvChecks(e.ResourceAssert)
	return e
}

package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
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
	e.AddAssertion(assert.ValueSet("directory.#", "1"))
	e.AddAssertion(assert.ValueSet("directory.0.enable", strconv.FormatBool(opts.Enable)))
	e.AddAssertion(assert.ValueSet("directory.0.auto_refresh", autoRefresh))
	e.AddAssertion(assert.ValueSet("directory.0.notification_integration", notificationIntegration))
	e.AddAssertion(assert.ValueSet("directory.0.refresh_on_create", refreshOnCreate))
	return e
}

func (e *ExternalGcsStageResourceAssert) HasEncryptionGcsSseKms() *ExternalGcsStageResourceAssert {
	e.AddAssertion(assert.ValueSet("encryption.#", "1"))
	e.AddAssertion(assert.ValueSet("encryption.0.gcs_sse_kms.#", "1"))
	e.AddAssertion(assert.ValueSet("encryption.0.none.#", "0"))
	return e
}

func (e *ExternalGcsStageResourceAssert) HasEncryptionNone() *ExternalGcsStageResourceAssert {
	e.AddAssertion(assert.ValueSet("encryption.#", "1"))
	e.AddAssertion(assert.ValueSet("encryption.0.gcs_sse_kms.#", "0"))
	e.AddAssertion(assert.ValueSet("encryption.0.none.#", "1"))
	return e
}

func (e *ExternalGcsStageResourceAssert) HasStageTypeEnum(expected sdk.StageType) *ExternalGcsStageResourceAssert {
	e.AddAssertion(assert.ValueSet("stage_type", string(expected)))
	return e
}

func (e *ExternalGcsStageResourceAssert) HasCloudEnum(expected sdk.StageCloud) *ExternalGcsStageResourceAssert {
	e.AddAssertion(assert.ValueSet("cloud", string(expected)))
	return e
}

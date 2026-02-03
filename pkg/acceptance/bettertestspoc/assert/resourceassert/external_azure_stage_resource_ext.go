package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (e *ExternalAzureStageResourceAssert) HasDirectoryEnableString(expected string) *ExternalAzureStageResourceAssert {
	e.AddAssertion(assert.ValueSet("directory.0.enable", expected))
	return e
}

func (e *ExternalAzureStageResourceAssert) HasDirectoryAutoRefreshString(expected string) *ExternalAzureStageResourceAssert {
	e.AddAssertion(assert.ValueSet("directory.0.auto_refresh", expected))
	return e
}

func (e *ExternalAzureStageResourceAssert) HasDirectoryRefreshOnCreateString(expected string) *ExternalAzureStageResourceAssert {
	e.AddAssertion(assert.ValueSet("directory.0.refresh_on_create", expected))
	return e
}

func (e *ExternalAzureStageResourceAssert) HasDirectory(enable bool, autoRefresh bool) *ExternalAzureStageResourceAssert {
	e.AddAssertion(assert.ValueSet("directory.#", "1"))
	e.AddAssertion(assert.ValueSet("directory.0.enable", strconv.FormatBool(enable)))
	e.AddAssertion(assert.ValueSet("directory.0.auto_refresh", strconv.FormatBool(autoRefresh)))
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

func (e *ExternalAzureStageResourceAssert) HasCredentials() *ExternalAzureStageResourceAssert {
	e.AddAssertion(assert.ValueSet("credentials.#", "1"))
	return e
}

package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// Directory block assertions
func (i *InternalStageResourceAssert) HasDirectoryEnableString(expected string) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("directory.0.enable", expected))
	return i
}

func (i *InternalStageResourceAssert) HasDirectoryAutoRefreshString(expected string) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("directory.0.auto_refresh", expected))
	return i
}

func (i *InternalStageResourceAssert) HasDirectory(enable bool, autoRefresh bool) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("directory.#", "1"))
	i.AddAssertion(assert.ValueSet("directory.0.enable", strconv.FormatBool(enable)))
	i.AddAssertion(assert.ValueSet("directory.0.auto_refresh", strconv.FormatBool(autoRefresh)))
	return i
}

// Encryption block assertions
func (i *InternalStageResourceAssert) HasEncryptionSnowflakeFull() *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("encryption.#", "1"))
	i.AddAssertion(assert.ValueSet("encryption.0.snowflake_full.#", "1"))
	i.AddAssertion(assert.ValueSet("encryption.0.snowflake_sse.#", "0"))
	return i
}

func (i *InternalStageResourceAssert) HasEncryptionSnowflakeSse() *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("encryption.#", "1"))
	i.AddAssertion(assert.ValueSet("encryption.0.snowflake_full.#", "0"))
	i.AddAssertion(assert.ValueSet("encryption.0.snowflake_sse.#", "1"))
	return i
}

func (i *InternalStageResourceAssert) HasStageTypeEnum(expected sdk.StageType) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("stage_type", string(expected)))
	return i
}

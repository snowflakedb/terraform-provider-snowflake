package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
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

// Show output assertions
func (i *InternalStageResourceAssert) HasShowOutputDirectoryEnabled(expected bool) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("show_output.0.directory_enabled", strconv.FormatBool(expected)))
	return i
}

func (i *InternalStageResourceAssert) HasShowOutputComment(expected string) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("show_output.0.comment", expected))
	return i
}

func (i *InternalStageResourceAssert) HasShowOutputName(expected string) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("show_output.0.name", expected))
	return i
}

func (i *InternalStageResourceAssert) HasShowOutputType(expected string) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("show_output.0.type", expected))
	return i
}

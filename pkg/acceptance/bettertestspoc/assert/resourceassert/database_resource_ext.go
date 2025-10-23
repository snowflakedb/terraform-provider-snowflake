package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (d *DatabaseResourceAssert) HasAllDefaultParameters() *DatabaseResourceAssert {
	return d.
		HasDataRetentionTimeInDaysString("1").
		HasMaxDataExtensionTimeInDaysString("14").
		HasExternalVolumeEmpty().
		HasCatalogEmpty().
		HasReplaceInvalidCharactersString("false").
		HasDefaultDdlCollationEmpty().
		HasStorageSerializationPolicyString(string(sdk.StorageSerializationPolicyOptimized)).
		HasLogLevelString("OFF").
		HasTraceLevelString("OFF").
		HasSuspendTaskAfterNumFailuresString("10").
		HasTaskAutoRetryAttemptsString("0").
		HasUserTaskManagedInitialWarehouseSizeString("Medium").
		HasUserTaskTimeoutMsString("3600000").
		HasUserTaskMinimumTriggerIntervalInSecondsString("30").
		HasQuotedIdentifiersIgnoreCaseString("false").
		HasEnableConsoleOutputString("false")
}

func (d *DatabaseResourceAssert) HasReplication(accountIdentifier sdk.AccountIdentifier, withFailover bool, ignoreEditionCheck bool) *DatabaseResourceAssert {
	d.AddAssertion(assert.ValueSet("replication.#", "1"))
	d.AddAssertion(assert.ValueSet("replication.0.enable_to_account.#", "1"))
	d.AddAssertion(assert.ValueSet("replication.0.enable_to_account.0.account_identifier", accountIdentifier.FullyQualifiedName()))
	d.AddAssertion(assert.ValueSet("replication.0.enable_to_account.0.with_failover", strconv.FormatBool(withFailover)))
	d.AddAssertion(assert.ValueSet("replication.0.ignore_edition_check", strconv.FormatBool(ignoreEditionCheck)))
	return d
}

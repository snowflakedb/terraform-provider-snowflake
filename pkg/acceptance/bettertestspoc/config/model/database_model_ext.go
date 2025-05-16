package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

// DatabaseWithParametersSet should be used to create database which sets the parameters that can be altered on the account level in other tests; this way, the test is not affected by the changes.
// TODO [this PR]: consider using helper instead
func DatabaseWithParametersSet(
	resourceName string,
	name string,
) *DatabaseModel {
	return Database(resourceName, name).
		WithDataRetentionTimeInDays(1).
		WithMaxDataExtensionTimeInDays(1).
		// according to the docs SNOWFLAKE is a valid value (https://docs.snowflake.com/en/sql-reference/parameters#catalog)
		WithCatalog("SNOWFLAKE")
}

func (d *DatabaseModel) WithReplication(accountIdentifier sdk.AccountIdentifier, withFailover bool, ignoreEditionCheck bool) *DatabaseModel {
	return d.WithReplicationValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"enable_to_account": tfconfig.ObjectVariable(
					map[string]tfconfig.Variable{
						"account_identifier": tfconfig.StringVariable(accountIdentifier.FullyQualifiedName()),
						"with_failover":      tfconfig.BoolVariable(withFailover),
					},
				),
				"ignore_edition_check": tfconfig.BoolVariable(ignoreEditionCheck),
			},
		),
	)
}

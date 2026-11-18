package datasourcemodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (t *TablesModel) WithEmptyIn() *TablesModel {
	return t.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
		}),
	)
}

func (t *TablesModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *TablesModel {
	return t.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (t *TablesModel) WithInSchema(schemaId sdk.DatabaseObjectIdentifier) *TablesModel {
	return t.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"schema": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
		}),
	)
}

func (t *TablesModel) WithInAccount() *TablesModel {
	return t.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"account": tfconfig.BoolVariable(true),
		}),
	)
}

func (t *TablesModel) WithInApplication(applicationId sdk.AccountObjectIdentifier) *TablesModel {
	return t.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"application": tfconfig.StringVariable(applicationId.Name()),
		}),
	)
}

func (t *TablesModel) WithInApplicationPackage(applicationPackageId sdk.AccountObjectIdentifier) *TablesModel {
	return t.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"application_package": tfconfig.StringVariable(applicationPackageId.Name()),
		}),
	)
}

func (t *TablesModel) WithLimit(rows int) *TablesModel {
	return t.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}

package datasourcemodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (a *AlertsModel) WithEmptyIn() *AlertsModel {
	return a.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
		}),
	)
}

func (a *AlertsModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *AlertsModel {
	return a.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (a *AlertsModel) WithInSchema(schemaId sdk.DatabaseObjectIdentifier) *AlertsModel {
	return a.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"schema": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
		}),
	)
}

func (a *AlertsModel) WithInApplication(applicationId sdk.AccountObjectIdentifier) *AlertsModel {
	return a.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"application": tfconfig.StringVariable(applicationId.Name()),
		}),
	)
}

func (a *AlertsModel) WithInApplicationPackage(applicationPackageId sdk.AccountObjectIdentifier) *AlertsModel {
	return a.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"application_package": tfconfig.StringVariable(applicationPackageId.Name()),
		}),
	)
}

func (a *AlertsModel) WithLimit(rows int) *AlertsModel {
	return a.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}

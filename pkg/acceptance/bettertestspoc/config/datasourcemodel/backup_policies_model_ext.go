package datasourcemodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (b *BackupPoliciesModel) WithEmptyIn() *BackupPoliciesModel {
	return b.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
		}),
	)
}

func (b *BackupPoliciesModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *BackupPoliciesModel {
	return b.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (b *BackupPoliciesModel) WithInSchema(schemaId sdk.DatabaseObjectIdentifier) *BackupPoliciesModel {
	return b.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"schema": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
		}),
	)
}

func (b *BackupPoliciesModel) WithLimit(rows int) *BackupPoliciesModel {
	return b.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}

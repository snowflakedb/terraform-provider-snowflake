package datasourcemodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (p *PasswordPoliciesModel) WithRowsAndFrom(rows int, from string) *PasswordPoliciesModel {
	return p.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
			"from": tfconfig.StringVariable(from),
		}),
	)
}

func (p *PasswordPoliciesModel) WithEmptyIn() *PasswordPoliciesModel {
	return p.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
		}),
	)
}

func (p *PasswordPoliciesModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *PasswordPoliciesModel {
	return p.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (p *PasswordPoliciesModel) WithInSchema(schemaId sdk.DatabaseObjectIdentifier) *PasswordPoliciesModel {
	return p.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"schema": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
		}),
	)
}

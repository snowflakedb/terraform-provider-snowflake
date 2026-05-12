package datasourcemodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (a *AuthenticationPoliciesModel) WithRowsAndFrom(rows int, from string) *AuthenticationPoliciesModel {
	return a.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
			"from": tfconfig.StringVariable(from),
		}),
	)
}

func (a *AuthenticationPoliciesModel) WithEmptyIn() *AuthenticationPoliciesModel {
	return a.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
		}),
	)
}

func (a *AuthenticationPoliciesModel) WithEmptyOn() *AuthenticationPoliciesModel {
	return a.WithOnValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
		}),
	)
}

func (a *AuthenticationPoliciesModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *AuthenticationPoliciesModel {
	return a.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (a *AuthenticationPoliciesModel) WithOnUser(userId sdk.AccountObjectIdentifier) *AuthenticationPoliciesModel {
	return a.WithOnValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"user": tfconfig.StringVariable(userId.Name()),
		}),
	)
}

func (a *AuthenticationPoliciesModel) WithOnAccount() *AuthenticationPoliciesModel {
	return a.WithOnValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"account": tfconfig.BoolVariable(true),
		}),
	)
}

package datasourcemodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (s *SessionPoliciesModel) WithRowsAndFrom(rows int, from string) *SessionPoliciesModel {
	return s.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
			"from": tfconfig.StringVariable(from),
		}),
	)
}

func (s *SessionPoliciesModel) WithEmptyIn() *SessionPoliciesModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
		}),
	)
}

func (s *SessionPoliciesModel) WithEmptyOn() *SessionPoliciesModel {
	return s.WithOnValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
		}),
	)
}

func (s *SessionPoliciesModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *SessionPoliciesModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (s *SessionPoliciesModel) WithInSchema(schemaId sdk.DatabaseObjectIdentifier) *SessionPoliciesModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"schema": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
		}),
	)
}

func (s *SessionPoliciesModel) WithOnUser(userId sdk.AccountObjectIdentifier) *SessionPoliciesModel {
	return s.WithOnValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"user": tfconfig.StringVariable(userId.Name()),
		}),
	)
}

func (s *SessionPoliciesModel) WithOnAccount() *SessionPoliciesModel {
	return s.WithOnValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"account": tfconfig.BoolVariable(true),
		}),
	)
}

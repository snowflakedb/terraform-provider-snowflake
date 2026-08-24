package datasourcemodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (s *StreamlitsModel) WithInAccount() *StreamlitsModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"account": tfconfig.BoolVariable(true),
		}),
	)
}

func (s *StreamlitsModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *StreamlitsModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (s *StreamlitsModel) WithInSchema(schemaId sdk.DatabaseObjectIdentifier) *StreamlitsModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"schema": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
		}),
	)
}

func (s *StreamlitsModel) WithInDatabaseAndSchema(databaseName string, schemaName string) *StreamlitsModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseName),
			"schema":   tfconfig.StringVariable(schemaName),
		}),
	)
}

func (s *StreamlitsModel) WithEmptyIn() *StreamlitsModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
		}),
	)
}

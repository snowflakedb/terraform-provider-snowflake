package datasourcemodel

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (s *SecretsModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *SecretsModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (s *SecretsModel) WithInSchema(schemaId sdk.DatabaseObjectIdentifier) *SecretsModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"schema": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
		}),
	)
}

func (s *SecretsModel) WithInAccount() *SecretsModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"account": tfconfig.BoolVariable(true),
		}),
	)
}

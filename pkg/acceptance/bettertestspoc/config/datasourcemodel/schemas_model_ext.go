package datasourcemodel

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *SchemasModel) WithLimit(rows int) *SchemasModel {
	return s.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}

func (s *SchemasModel) WithRowsAndFrom(rows int, from string) *SchemasModel {
	return s.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
			"from": tfconfig.StringVariable(from),
		}),
	)
}

func (s *SchemasModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *SchemasModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (s *SchemasModel) WithInAccount() *SchemasModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"account": tfconfig.BoolVariable(true),
		}),
	)
}

func (s *SchemasModel) WithInApplication(applicationId sdk.AccountObjectIdentifier) *SchemasModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"application": tfconfig.StringVariable(applicationId.Name()),
		}),
	)
}

func (s *SchemasModel) WithInApplicationPackage(applicationPackageId sdk.AccountObjectIdentifier) *SchemasModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"application_package": tfconfig.StringVariable(applicationPackageId.Name()),
		}),
	)
}

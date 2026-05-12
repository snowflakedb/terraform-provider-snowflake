package datasourcemodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (r *RowAccessPoliciesModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *RowAccessPoliciesModel {
	return r.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (r *RowAccessPoliciesModel) WithInSchema(schemaId sdk.DatabaseObjectIdentifier) *RowAccessPoliciesModel {
	return r.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"schema": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
		}),
	)
}

func (r *RowAccessPoliciesModel) WithInAccount() *RowAccessPoliciesModel {
	return r.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"account": tfconfig.BoolVariable(true),
		}),
	)
}

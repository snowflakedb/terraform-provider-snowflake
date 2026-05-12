package datasourcemodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (m *MaskingPoliciesModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *MaskingPoliciesModel {
	return m.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (m *MaskingPoliciesModel) WithInSchema(schemaId sdk.DatabaseObjectIdentifier) *MaskingPoliciesModel {
	return m.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"schema": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
		}),
	)
}

func (m *MaskingPoliciesModel) WithInAccount() *MaskingPoliciesModel {
	return m.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"account": tfconfig.BoolVariable(true),
		}),
	)
}

func (m *MaskingPoliciesModel) WithRows(rows int) *MaskingPoliciesModel {
	return m.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}

func (m *MaskingPoliciesModel) WithRowsAndFrom(rows int, from string) *MaskingPoliciesModel {
	return m.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
			"from": tfconfig.StringVariable(from),
		}),
	)
}

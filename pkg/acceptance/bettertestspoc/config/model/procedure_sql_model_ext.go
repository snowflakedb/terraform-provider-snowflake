package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

func ProcedureSqlBasicInline(
	resourceName string,
	id sdk.SchemaObjectIdentifierWithArguments,
	returnType datatypes.DataType,
	procedureDefinition string,
) *ProcedureSqlModel {
	return ProcedureSql(resourceName, id.DatabaseName(), id.SchemaName(), id.Name(), procedureDefinition, returnType.ToSql())
}

func (f *ProcedureSqlModel) WithArgument(argName string, argDataType datatypes.DataType) *ProcedureSqlModel {
	return f.WithArgumentsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"arg_name":      tfconfig.StringVariable(argName),
				"arg_data_type": tfconfig.StringVariable(argDataType.ToSql()),
			},
		),
	)
}

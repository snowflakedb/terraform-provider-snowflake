package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

func ProcedureJavascriptBasicInline(
	resourceName string,
	id sdk.SchemaObjectIdentifierWithArguments,
	returnType datatypes.DataType,
	procedureDefinition string,
) *ProcedureJavascriptModel {
	return ProcedureJavascript(resourceName, id.DatabaseName(), id.SchemaName(), id.Name(), procedureDefinition, returnType.ToSql())
}

func (f *ProcedureJavascriptModel) WithArgument(argName string, argDataType datatypes.DataType) *ProcedureJavascriptModel {
	return f.WithArgumentsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"arg_name":      tfconfig.StringVariable(argName),
				"arg_data_type": tfconfig.StringVariable(argDataType.ToSql()),
			},
		),
	)
}

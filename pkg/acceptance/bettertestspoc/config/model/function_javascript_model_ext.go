package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

func FunctionJavascriptInline(resourceName string, id sdk.SchemaObjectIdentifierWithArguments, functionDefinition string, returnType string) *FunctionJavascriptModel {
	f := &FunctionJavascriptModel{ResourceModelMeta: config.Meta(resourceName, resources.FunctionJavascript)}
	f.WithDatabase(id.DatabaseName())
	f.WithFunctionDefinition(functionDefinition)
	f.WithName(id.Name())
	f.WithReturnType(returnType)
	f.WithSchema(id.SchemaName())
	return f
}

func (f *FunctionJavascriptModel) WithArgument(argName string, argDataType datatypes.DataType) *FunctionJavascriptModel {
	return f.WithArgumentsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"arg_name":      tfconfig.StringVariable(argName),
				"arg_data_type": tfconfig.StringVariable(argDataType.ToSql()),
			},
		),
	)
}

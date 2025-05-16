package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
)

func FunctionPythonBasicInline(resourceName string, id sdk.SchemaObjectIdentifierWithArguments, runtimeVersion string, returnType datatypes.DataType, handler string, functionDefinition string) *FunctionPythonModel {
	return FunctionPython(resourceName, id.DatabaseName(), id.SchemaName(), id.Name(), handler, returnType.ToSql(), runtimeVersion).WithFunctionDefinition(functionDefinition)
}

func (f *FunctionPythonModel) WithArgument(argName string, argDataType datatypes.DataType) *FunctionPythonModel {
	return f.WithArgumentsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"arg_name":      tfconfig.StringVariable(argName),
				"arg_data_type": tfconfig.StringVariable(argDataType.ToSql()),
			},
		),
	)
}

func (f *FunctionPythonModel) WithImports(imports ...sdk.NormalizedPath) *FunctionPythonModel {
	return f.WithImportsValue(
		tfconfig.SetVariable(
			collections.Map(imports, func(imp sdk.NormalizedPath) tfconfig.Variable {
				return tfconfig.ObjectVariable(
					map[string]tfconfig.Variable{
						"stage_location": tfconfig.StringVariable(imp.StageLocation),
						"path_on_stage":  tfconfig.StringVariable(imp.PathOnStage),
					},
				)
			})...,
		),
	)
}

func (f *FunctionPythonModel) WithExternalAccessIntegrations(ids ...sdk.AccountObjectIdentifier) *FunctionPythonModel {
	return f.WithExternalAccessIntegrationsValue(
		tfconfig.SetVariable(
			collections.Map(ids, func(id sdk.AccountObjectIdentifier) tfconfig.Variable { return tfconfig.StringVariable(id.Name()) })...,
		),
	)
}

func (f *FunctionPythonModel) WithSecrets(secrets map[string]sdk.SchemaObjectIdentifier) *FunctionPythonModel {
	objects := make([]tfconfig.Variable, 0)
	for k, v := range secrets {
		objects = append(objects, tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"secret_variable_name": tfconfig.StringVariable(k),
				"secret_id":            tfconfig.StringVariable(v.FullyQualifiedName()),
			},
		))
	}

	return f.WithSecretsValue(
		tfconfig.SetVariable(
			objects...,
		),
	)
}

func (f *FunctionPythonModel) WithPackages(pkgs ...string) *FunctionPythonModel {
	return f.WithPackagesValue(
		tfconfig.SetVariable(
			collections.Map(pkgs, func(pkg string) tfconfig.Variable { return tfconfig.StringVariable(pkg) })...,
		),
	)
}

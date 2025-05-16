package model

import (
	"strings"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
)

func FunctionJavaBasicInline(
	resourceName string,
	id sdk.SchemaObjectIdentifierWithArguments,
	returnType datatypes.DataType,
	handler string,
	functionDefinition string,
) *FunctionJavaModel {
	return FunctionJava(resourceName, id.DatabaseName(), id.SchemaName(), id.Name(), handler, returnType.ToSql()).WithFunctionDefinition(functionDefinition)
}

func FunctionJavaBasicStaged(
	resourceName string,
	id sdk.SchemaObjectIdentifierWithArguments,
	returnType datatypes.DataType,
	handler string,
	stageLocation string,
	pathOnStage string,
) *FunctionJavaModel {
	return FunctionJava(resourceName, id.DatabaseName(), id.SchemaName(), id.Name(), handler, returnType.ToSql()).
		WithImport(stageLocation, pathOnStage)
}

func (f *FunctionJavaModel) WithArgument(argName string, argDataType datatypes.DataType) *FunctionJavaModel {
	return f.WithArgumentsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"arg_name":      tfconfig.StringVariable(argName),
				"arg_data_type": tfconfig.StringVariable(argDataType.ToSql()),
			},
		),
	)
}

func (f *FunctionJavaModel) WithArgumentWithDefaultValue(argName string, argDataType datatypes.DataType, value string) *FunctionJavaModel {
	return f.WithArgumentsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"arg_name":          tfconfig.StringVariable(argName),
				"arg_data_type":     tfconfig.StringVariable(argDataType.ToSql()),
				"arg_default_value": tfconfig.StringVariable(value),
			},
		),
	)
}

func (f *FunctionJavaModel) WithImport(stageLocation string, pathOnStage string) *FunctionJavaModel {
	return f.WithImportsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"stage_location": tfconfig.StringVariable(strings.TrimPrefix(stageLocation, "@")),
				"path_on_stage":  tfconfig.StringVariable(pathOnStage),
			},
		),
	)
}

func (f *FunctionJavaModel) WithImports(imports ...sdk.NormalizedPath) *FunctionJavaModel {
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

func (f *FunctionJavaModel) WithPackages(pkgs ...string) *FunctionJavaModel {
	return f.WithPackagesValue(
		tfconfig.SetVariable(
			collections.Map(pkgs, func(pkg string) tfconfig.Variable { return tfconfig.StringVariable(pkg) })...,
		),
	)
}

func (f *FunctionJavaModel) WithExternalAccessIntegrations(ids ...sdk.AccountObjectIdentifier) *FunctionJavaModel {
	return f.WithExternalAccessIntegrationsValue(
		tfconfig.SetVariable(
			collections.Map(ids, func(id sdk.AccountObjectIdentifier) tfconfig.Variable { return tfconfig.StringVariable(id.Name()) })...,
		),
	)
}

func (f *FunctionJavaModel) WithSecrets(secrets map[string]sdk.SchemaObjectIdentifier) *FunctionJavaModel {
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

func (f *FunctionJavaModel) WithTargetPathParts(stageLocation string, pathOnStage string) *FunctionJavaModel {
	return f.WithTargetPathValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"stage_location": tfconfig.StringVariable(stageLocation),
				"path_on_stage":  tfconfig.StringVariable(pathOnStage),
			},
		),
	)
}

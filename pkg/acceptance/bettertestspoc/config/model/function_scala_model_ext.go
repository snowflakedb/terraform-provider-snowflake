package model

import (
	"strings"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
)

func FunctionScalaBasicInline(
	resourceName string,
	id sdk.SchemaObjectIdentifierWithArguments,
	runtimeVersion string,
	returnType datatypes.DataType,
	handler string,
	functionDefinition string,
) *FunctionScalaModel {
	return FunctionScala(resourceName, id.DatabaseName(), id.SchemaName(), id.Name(), handler, returnType.ToSql(), runtimeVersion).WithFunctionDefinition(functionDefinition)
}

func (f *FunctionScalaModel) WithArgument(argName string, argDataType datatypes.DataType) *FunctionScalaModel {
	return f.WithArgumentsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"arg_name":      tfconfig.StringVariable(argName),
				"arg_data_type": tfconfig.StringVariable(argDataType.ToSql()),
			},
		),
	)
}

func (f *FunctionScalaModel) WithImport(stageLocation string, pathOnStage string) *FunctionScalaModel {
	return f.WithImportsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"stage_location": tfconfig.StringVariable(strings.TrimPrefix(stageLocation, "@")),
				"path_on_stage":  tfconfig.StringVariable(pathOnStage),
			},
		),
	)
}

func (f *FunctionScalaModel) WithImports(imports ...sdk.NormalizedPath) *FunctionScalaModel {
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

func (f *FunctionScalaModel) WithPackages(pkgs ...string) *FunctionScalaModel {
	return f.WithPackagesValue(
		tfconfig.SetVariable(
			collections.Map(pkgs, func(pkg string) tfconfig.Variable { return tfconfig.StringVariable(pkg) })...,
		),
	)
}

func (f *FunctionScalaModel) WithExternalAccessIntegrations(ids ...sdk.AccountObjectIdentifier) *FunctionScalaModel {
	return f.WithExternalAccessIntegrationsValue(
		tfconfig.SetVariable(
			collections.Map(ids, func(id sdk.AccountObjectIdentifier) tfconfig.Variable { return tfconfig.StringVariable(id.Name()) })...,
		),
	)
}

func (f *FunctionScalaModel) WithSecrets(secrets map[string]sdk.SchemaObjectIdentifier) *FunctionScalaModel {
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

func (f *FunctionScalaModel) WithTargetPathParts(stageLocation string, pathOnStage string) *FunctionScalaModel {
	return f.WithTargetPathValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"stage_location": tfconfig.StringVariable(stageLocation),
				"path_on_stage":  tfconfig.StringVariable(pathOnStage),
			},
		),
	)
}

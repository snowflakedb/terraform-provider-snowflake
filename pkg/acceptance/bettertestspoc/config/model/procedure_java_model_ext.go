package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
)

func ProcedureJavaBasicInline(
	resourceName string,
	id sdk.SchemaObjectIdentifierWithArguments,
	returnType datatypes.DataType,
	handler string,
	procedureDefinition string,
) *ProcedureJavaModel {
	return ProcedureJava(resourceName, id.DatabaseName(), id.SchemaName(), id.Name(), handler, returnType.ToSql(), "11", "1.14.0").
		WithProcedureDefinition(procedureDefinition)
}

func ProcedureJavaBasicStaged(
	resourceName string,
	id sdk.SchemaObjectIdentifierWithArguments,
	returnType datatypes.DataType,
	handler string,
	stageLocation string,
	pathOnStage string,
) *ProcedureJavaModel {
	return ProcedureJava(resourceName, id.DatabaseName(), id.SchemaName(), id.Name(), handler, returnType.ToSql(), "11", "1.14.0").
		WithImport(stageLocation, pathOnStage)
}

func (f *ProcedureJavaModel) WithArgument(argName string, argDataType datatypes.DataType) *ProcedureJavaModel {
	return f.WithArgumentsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"arg_name":      tfconfig.StringVariable(argName),
				"arg_data_type": tfconfig.StringVariable(argDataType.ToSql()),
			},
		),
	)
}

func (f *ProcedureJavaModel) WithArgumentWithDefaultValue(argName string, argDataType datatypes.DataType, value string) *ProcedureJavaModel {
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

func (f *ProcedureJavaModel) WithImport(stageLocation string, pathOnStage string) *ProcedureJavaModel {
	return f.WithImportsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"stage_location": tfconfig.StringVariable(stageLocation),
				"path_on_stage":  tfconfig.StringVariable(pathOnStage),
			},
		),
	)
}

func (f *ProcedureJavaModel) WithImports(imports ...sdk.NormalizedPath) *ProcedureJavaModel {
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

func (f *ProcedureJavaModel) WithPackages(pkgs ...string) *ProcedureJavaModel {
	return f.WithPackagesValue(
		tfconfig.SetVariable(
			collections.Map(pkgs, func(pkg string) tfconfig.Variable { return tfconfig.StringVariable(pkg) })...,
		),
	)
}

func (f *ProcedureJavaModel) WithExternalAccessIntegrations(ids ...sdk.AccountObjectIdentifier) *ProcedureJavaModel {
	return f.WithExternalAccessIntegrationsValue(
		tfconfig.SetVariable(
			collections.Map(ids, func(id sdk.AccountObjectIdentifier) tfconfig.Variable { return tfconfig.StringVariable(id.Name()) })...,
		),
	)
}

func (f *ProcedureJavaModel) WithSecrets(secrets map[string]sdk.SchemaObjectIdentifier) *ProcedureJavaModel {
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

func (f *ProcedureJavaModel) WithTargetPathParts(stageLocation string, pathOnStage string) *ProcedureJavaModel {
	return f.WithTargetPathValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"stage_location": tfconfig.StringVariable(stageLocation),
				"path_on_stage":  tfconfig.StringVariable(pathOnStage),
			},
		),
	)
}

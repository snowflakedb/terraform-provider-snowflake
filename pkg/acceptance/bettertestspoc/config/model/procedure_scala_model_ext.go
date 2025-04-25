package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

func ProcedureScalaBasicInline(
	resourceName string,
	id sdk.SchemaObjectIdentifierWithArguments,
	returnType datatypes.DataType,
	handler string,
	procedureDefinition string,
) *ProcedureScalaModel {
	return ProcedureScala(resourceName, id.DatabaseName(), id.SchemaName(), id.Name(), handler, returnType.ToSql(), "2.12", "1.14.0").
		WithProcedureDefinition(procedureDefinition)
}

func ProcedureScalaBasicStaged(
	resourceName string,
	id sdk.SchemaObjectIdentifierWithArguments,
	returnType datatypes.DataType,
	handler string,
	stageLocation string,
	pathOnStage string,
) *ProcedureScalaModel {
	return ProcedureScala(resourceName, id.DatabaseName(), id.SchemaName(), id.Name(), handler, returnType.ToSql(), "2.12", "1.14.0").
		WithImport(stageLocation, pathOnStage)
}

func (f *ProcedureScalaModel) WithArgument(argName string, argDataType datatypes.DataType) *ProcedureScalaModel {
	return f.WithArgumentsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"arg_name":      tfconfig.StringVariable(argName),
				"arg_data_type": tfconfig.StringVariable(argDataType.ToSql()),
			},
		),
	)
}

func (f *ProcedureScalaModel) WithArgumentWithDefaultValue(argName string, argDataType datatypes.DataType, value string) *ProcedureScalaModel {
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

func (f *ProcedureScalaModel) WithImport(stageLocation string, pathOnStage string) *ProcedureScalaModel {
	return f.WithImportsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"stage_location": tfconfig.StringVariable(stageLocation),
				"path_on_stage":  tfconfig.StringVariable(pathOnStage),
			},
		),
	)
}

func (f *ProcedureScalaModel) WithImports(imports ...sdk.NormalizedPath) *ProcedureScalaModel {
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

func (f *ProcedureScalaModel) WithPackages(pkgs ...string) *ProcedureScalaModel {
	return f.WithPackagesValue(
		tfconfig.SetVariable(
			collections.Map(pkgs, func(pkg string) tfconfig.Variable { return tfconfig.StringVariable(pkg) })...,
		),
	)
}

func (f *ProcedureScalaModel) WithExternalAccessIntegrations(ids ...sdk.AccountObjectIdentifier) *ProcedureScalaModel {
	return f.WithExternalAccessIntegrationsValue(
		tfconfig.SetVariable(
			collections.Map(ids, func(id sdk.AccountObjectIdentifier) tfconfig.Variable { return tfconfig.StringVariable(id.Name()) })...,
		),
	)
}

func (f *ProcedureScalaModel) WithSecrets(secrets map[string]sdk.SchemaObjectIdentifier) *ProcedureScalaModel {
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

func (f *ProcedureScalaModel) WithTargetPathParts(stageLocation string, pathOnStage string) *ProcedureScalaModel {
	return f.WithTargetPathValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"stage_location": tfconfig.StringVariable(stageLocation),
				"path_on_stage":  tfconfig.StringVariable(pathOnStage),
			},
		),
	)
}

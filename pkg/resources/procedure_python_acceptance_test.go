//go:build !account_level_tests

package resources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ProcedurePython_InlineBasic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)
	idWithChangedNameButTheSameDataType := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	definition := acc.TestClient().Procedure.SamplePythonDefinition(t, funcName, argName)

	procedureModel := model.ProcedurePythonBasicInline("w", id, dataType, funcName, definition).
		WithArgument(argName, dataType)
	procedureModelRenamed := model.ProcedurePythonBasicInline("w", idWithChangedNameButTheSameDataType, dataType, funcName, definition).
		WithArgument(argName, dataType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.ProcedurePython),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, procedureModel),
				Check: assertThat(t,
					resourceassert.ProcedurePythonResource(t, procedureModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanDefault).
						HasCommentString(sdk.DefaultProcedureComment).
						HasImportsLength(0).
						HasRuntimeVersionString("3.8").
						HasProcedureDefinitionString(definition).
						HasProcedureLanguageString("PYTHON").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ProcedureShowOutput(t, procedureModel.ResourceReference()).
						HasIsSecure(false),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "arguments.0.arg_name", argName)),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "arguments.0.arg_data_type", dataType.ToSql())),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "arguments.0.arg_default_value", "")),
				),
			},
			// REMOVE EXTERNALLY (CHECK RECREATION)
			{
				PreConfig: func() {
					acc.TestClient().Procedure.DropProcedureFunc(t, id)()
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(procedureModel.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, procedureModel),
				Check: assertThat(t,
					resourceassert.ProcedurePythonResource(t, procedureModel.ResourceReference()).
						HasNameString(id.Name()),
				),
			},
			// IMPORT
			{
				ResourceName:            procedureModel.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_secure", "arguments.0.arg_data_type", "null_input_behavior", "execute_as"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedProcedurePythonResource(t, id.FullyQualifiedName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_name", argName)),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_data_type", "VARCHAR(16777216)")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_default_value", "")),
				),
			},
			// RENAME
			{
				Config: config.FromModels(t, procedureModelRenamed),
				Check: assertThat(t,
					resourceassert.ProcedurePythonResource(t, procedureModelRenamed.ResourceReference()).
						HasNameString(idWithChangedNameButTheSameDataType.Name()).
						HasFullyQualifiedNameString(idWithChangedNameButTheSameDataType.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_ProcedurePython_InlineFull(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	secretId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	secretId2 := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	networkRule, networkRuleCleanup := acc.TestClient().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)

	secret, secretCleanup := acc.TestClient().Secret.CreateWithGenericString(t, secretId, "test_secret_string")
	t.Cleanup(secretCleanup)

	secret2, secret2Cleanup := acc.TestClient().Secret.CreateWithGenericString(t, secretId2, "test_secret_string_2")
	t.Cleanup(secret2Cleanup)

	externalAccessIntegration, externalAccessIntegrationCleanup := acc.TestClient().ExternalAccessIntegration.CreateExternalAccessIntegrationWithNetworkRuleAndSecret(t, networkRule.ID(), secret.ID())
	t.Cleanup(externalAccessIntegrationCleanup)

	externalAccessIntegration2, externalAccessIntegration2Cleanup := acc.TestClient().ExternalAccessIntegration.CreateExternalAccessIntegrationWithNetworkRuleAndSecret(t, networkRule.ID(), secret2.ID())
	t.Cleanup(externalAccessIntegration2Cleanup)

	tmpPythonFunction := acc.TestClient().CreateSamplePythonFunctionAndModuleOnUserStage(t)
	tmpPythonFunction2 := acc.TestClient().CreateSamplePythonFunctionAndModuleOnUserStage(t)

	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	definition := acc.TestClient().Procedure.SamplePythonDefinition(t, funcName, argName)

	procedureModel := model.ProcedurePythonBasicInline("w", id, dataType, funcName, definition).
		WithArgument(argName, dataType).
		WithImports(
			sdk.NormalizedPath{StageLocation: "~", PathOnStage: tmpPythonFunction.PythonFileName()},
			sdk.NormalizedPath{StageLocation: "~", PathOnStage: tmpPythonFunction2.PythonFileName()},
		).
		WithSnowparkPackage("1.14.0").
		WithPackages("absl-py==0.10.0").
		WithExternalAccessIntegrations(externalAccessIntegration, externalAccessIntegration2).
		WithSecrets(map[string]sdk.SchemaObjectIdentifier{
			"abc": secretId,
			"def": secretId2,
		}).
		WithRuntimeVersion("3.8").
		WithIsSecure("false").
		WithNullInputBehavior(string(sdk.NullInputBehaviorCalledOnNullInput)).
		WithExecuteAs(string(sdk.ExecuteAsCaller)).
		WithComment("some comment")

	procedureModelUpdateWithoutRecreation := model.ProcedurePythonBasicInline("w", id, dataType, funcName, definition).
		WithArgument(argName, dataType).
		WithImports(
			sdk.NormalizedPath{StageLocation: "~", PathOnStage: tmpPythonFunction.PythonFileName()},
			sdk.NormalizedPath{StageLocation: "~", PathOnStage: tmpPythonFunction2.PythonFileName()},
		).
		WithSnowparkPackage("1.14.0").
		WithPackages("absl-py==0.10.0").
		WithExternalAccessIntegrations(externalAccessIntegration).
		WithSecrets(map[string]sdk.SchemaObjectIdentifier{
			"def": secretId2,
		}).
		WithRuntimeVersion("3.8").
		WithIsSecure("false").
		WithNullInputBehavior(string(sdk.NullInputBehaviorCalledOnNullInput)).
		WithExecuteAs(string(sdk.ExecuteAsOwner)).
		WithComment("some other comment")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.ProcedurePython),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, procedureModel),
				Check: assertThat(t,
					resourceassert.ProcedurePythonResource(t, procedureModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanFalse).
						HasImportsLength(2).
						HasRuntimeVersionString("3.8").
						HasProcedureDefinitionString(definition).
						HasCommentString("some comment").
						HasProcedureLanguageString("PYTHON").
						HasExecuteAsString(string(sdk.ExecuteAsCaller)).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "secrets.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "external_access_integrations.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "packages.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "packages.0", "absl-py==0.10.0")),
					resourceshowoutputassert.ProcedureShowOutput(t, procedureModel.ResourceReference()).
						HasIsSecure(false),
				),
			},
			// IMPORT
			{
				ResourceName:            procedureModel.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"arguments.0.arg_data_type"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedProcedurePythonResource(t, id.FullyQualifiedName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_name", argName)),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_data_type", "VARCHAR(16777216)")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_default_value", "")),
				),
			},
			// UPDATE WITHOUT RECREATION
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(procedureModelUpdateWithoutRecreation.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, procedureModelUpdateWithoutRecreation),
				Check: assertThat(t,
					resourceassert.ProcedurePythonResource(t, procedureModelUpdateWithoutRecreation.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanFalse).
						HasImportsLength(2).
						HasRuntimeVersionString("3.8").
						HasProcedureDefinitionString(definition).
						HasCommentString("some other comment").
						HasProcedureLanguageString("PYTHON").
						HasExecuteAsString(string(sdk.ExecuteAsOwner)).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(procedureModelUpdateWithoutRecreation.ResourceReference(), "secrets.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(procedureModelUpdateWithoutRecreation.ResourceReference(), "secrets.0.secret_variable_name", "def")),
					assert.Check(resource.TestCheckResourceAttr(procedureModelUpdateWithoutRecreation.ResourceReference(), "secrets.0.secret_id", secretId2.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(procedureModelUpdateWithoutRecreation.ResourceReference(), "external_access_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(procedureModelUpdateWithoutRecreation.ResourceReference(), "external_access_integrations.0", externalAccessIntegration.Name())),
					assert.Check(resource.TestCheckResourceAttr(procedureModelUpdateWithoutRecreation.ResourceReference(), "packages.#", "1")),
					resourceshowoutputassert.ProcedureShowOutput(t, procedureModelUpdateWithoutRecreation.ResourceReference()).
						HasIsSecure(false),
				),
			},
		},
	})
}

func TestAcc_ProcedurePython_ImportsDifSuppression(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	stage, stageCleanup := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	tmpPythonFunction := acc.TestClient().CreateSamplePythonFunctionAndModuleOnStage(t, stage)
	tmpPythonFunction2 := acc.TestClient().CreateSamplePythonFunctionAndModuleOnStage(t, stage)

	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	definition := acc.TestClient().Procedure.SamplePythonDefinition(t, funcName, argName)

	procedureModel := model.ProcedurePythonBasicInline("w", id, dataType, funcName, definition).
		WithArgument(argName, dataType).
		WithImports(
			sdk.NormalizedPath{StageLocation: stage.ID().FullyQualifiedName(), PathOnStage: tmpPythonFunction.PythonFileName()},
			sdk.NormalizedPath{StageLocation: stage.ID().FullyQualifiedName(), PathOnStage: tmpPythonFunction2.PythonFileName()},
		).
		WithSnowparkPackage("1.14.0").
		WithRuntimeVersion("3.8").
		WithIsSecure("false")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.ProcedurePython),
		Steps: []resource.TestStep{
			// CREATE AND EXPECT EMPTY PLAN
			{
				Config: config.FromModels(t, procedureModel),
				Check: assertThat(t,
					resourceassert.ProcedurePythonResource(t, procedureModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanFalse).
						HasImportsLength(2).
						HasRuntimeVersionString("3.8").
						HasProcedureDefinitionString(definition).
						HasProcedureLanguageString("PYTHON").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ProcedureShowOutput(t, procedureModel.ResourceReference()).
						HasIsSecure(false),
				),
			},
		},
	})
}

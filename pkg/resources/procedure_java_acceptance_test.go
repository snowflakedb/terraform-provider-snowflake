//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"testing"
	"time"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO [SNOW-1348103]: test external changes
// TODO [SNOW-1348103]: test changes of attributes separately

func TestAcc_ProcedureJava_InlineBasic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)
	idWithChangedNameButTheSameDataType := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := acc.TestClient().Procedure.SampleJavaDefinition(t, className, funcName, argName)

	procedureModel := model.ProcedureJavaBasicInline("w", id, dataType, handler, definition).
		WithArgument(argName, dataType)
	procedureModelRenamed := model.ProcedureJavaBasicInline("w", idWithChangedNameButTheSameDataType, dataType, handler, definition).
		WithArgument(argName, dataType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.ProcedureJava),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, procedureModel),
				Check: assertThat(t,
					resourceassert.ProcedureJavaResource(t, procedureModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanDefault).
						HasCommentString(sdk.DefaultProcedureComment).
						HasImportsLength(0).
						HasTargetPathEmpty().
						HasRuntimeVersionString("11").
						HasProcedureDefinitionString(definition).
						HasProcedureLanguageString("JAVA").
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
					resourceassert.ProcedureJavaResource(t, procedureModel.ResourceReference()).
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
					resourceassert.ImportedProcedureJavaResource(t, id.FullyQualifiedName()).
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
					resourceassert.ProcedureJavaResource(t, procedureModelRenamed.ResourceReference()).
						HasNameString(idWithChangedNameButTheSameDataType.Name()).
						HasFullyQualifiedNameString(idWithChangedNameButTheSameDataType.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_ProcedureJava_InlineEmptyArgs(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	className := "TestFunc"
	funcName := "echoVarchar"
	returnDataType := testdatatypes.DataTypeVarchar_100

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes()

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := acc.TestClient().Procedure.SampleJavaDefinitionNoArgs(t, className, funcName)

	procedureModel := model.ProcedureJavaBasicInline("w", id, returnDataType, handler, definition)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.ProcedureJava),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, procedureModel),
				Check: assertThat(t,
					resourceassert.ProcedureJavaResource(t, procedureModel.ResourceReference()).
						HasNameString(id.Name()).
						HasProcedureDefinitionString(definition).
						HasProcedureLanguageString("JAVA").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_ProcedureJava_InlineBasicDefaultArg(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100
	defaultValue := "'hello'"

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := acc.TestClient().Procedure.SampleJavaDefinition(t, className, funcName, argName)

	procedureModel := model.ProcedureJavaBasicInline("w", id, dataType, handler, definition).
		WithArgumentWithDefaultValue(argName, dataType, defaultValue)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.ProcedureJava),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, procedureModel),
				Check: assertThat(t,
					resourceassert.ProcedureJavaResource(t, procedureModel.ResourceReference()).
						HasNameString(id.Name()).
						HasProcedureDefinitionString(definition).
						HasProcedureLanguageString("JAVA").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "arguments.0.arg_name", argName)),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "arguments.0.arg_data_type", dataType.ToSql())),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "arguments.0.arg_default_value", defaultValue)),
				),
			},
		},
	})
}

func TestAcc_ProcedureJava_InlineFull(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	stage, stageCleanup := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

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

	tmpJavaProcedure := acc.TestClient().CreateSampleJavaProcedureAndJarOnUserStage(t)
	tmpJavaProcedure2 := acc.TestClient().CreateSampleJavaProcedureAndJarOnUserStage(t)

	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := acc.TestClient().Procedure.SampleJavaDefinition(t, className, funcName, argName)
	// TODO [SNOW-1850370]: extract to helper
	jarName := fmt.Sprintf("tf-%d-%s.jar", time.Now().Unix(), random.AlphaN(5))

	procedureModel := model.ProcedureJavaBasicInline("w", id, dataType, handler, definition).
		WithArgument(argName, dataType).
		WithImports(
			sdk.NormalizedPath{StageLocation: "~", PathOnStage: tmpJavaProcedure.JarName},
			sdk.NormalizedPath{StageLocation: "~", PathOnStage: tmpJavaProcedure2.JarName},
		).
		WithSnowparkPackage("1.14.0").
		WithPackages("com.snowflake:telemetry:0.1.0").
		WithExternalAccessIntegrations(externalAccessIntegration, externalAccessIntegration2).
		WithSecrets(map[string]sdk.SchemaObjectIdentifier{
			"abc": secretId,
			"def": secretId2,
		}).
		WithTargetPathParts(stage.ID().FullyQualifiedName(), jarName).
		WithRuntimeVersion("11").
		WithIsSecure("false").
		WithNullInputBehavior(string(sdk.NullInputBehaviorCalledOnNullInput)).
		WithExecuteAs(string(sdk.ExecuteAsCaller)).
		WithComment("some comment")

	procedureModelUpdateWithoutRecreation := model.ProcedureJavaBasicInline("w", id, dataType, handler, definition).
		WithArgument(argName, dataType).
		WithImports(
			sdk.NormalizedPath{StageLocation: "~", PathOnStage: tmpJavaProcedure.JarName},
			sdk.NormalizedPath{StageLocation: "~", PathOnStage: tmpJavaProcedure2.JarName},
		).
		WithSnowparkPackage("1.14.0").
		WithPackages("com.snowflake:telemetry:0.1.0").
		WithExternalAccessIntegrations(externalAccessIntegration).
		WithSecrets(map[string]sdk.SchemaObjectIdentifier{
			"def": secretId2,
		}).
		WithTargetPathParts(stage.ID().FullyQualifiedName(), jarName).
		WithRuntimeVersion("11").
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
		CheckDestroy: acc.CheckDestroy(t, resources.ProcedureJava),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, procedureModel),
				Check: assertThat(t,
					resourceassert.ProcedureJavaResource(t, procedureModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanFalse).
						HasImportsLength(2).
						HasRuntimeVersionString("11").
						HasProcedureDefinitionString(definition).
						HasCommentString("some comment").
						HasProcedureLanguageString("JAVA").
						HasExecuteAsString(string(sdk.ExecuteAsCaller)).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "target_path.0.stage_location", stage.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "target_path.0.path_on_stage", jarName)),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "secrets.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "external_access_integrations.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "packages.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "packages.0", "com.snowflake:telemetry:0.1.0")),
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
					resourceassert.ImportedProcedureJavaResource(t, id.FullyQualifiedName()).
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
					resourceassert.ProcedureJavaResource(t, procedureModelUpdateWithoutRecreation.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanFalse).
						HasImportsLength(2).
						HasRuntimeVersionString("11").
						HasProcedureDefinitionString(definition).
						HasCommentString("some other comment").
						HasProcedureLanguageString("JAVA").
						HasExecuteAsString(string(sdk.ExecuteAsOwner)).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(procedureModelUpdateWithoutRecreation.ResourceReference(), "target_path.0.stage_location", stage.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(procedureModelUpdateWithoutRecreation.ResourceReference(), "target_path.0.path_on_stage", jarName)),
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

func TestAcc_ProcedureJava_StagedBasic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	stage, stageCleanup := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	tmpJavaProcedure := acc.TestClient().CreateSampleJavaProcedureAndJarOnStage(t, stage)

	dataType := tmpJavaProcedure.ArgType
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	argName := "x"
	handler := tmpJavaProcedure.JavaHandler()

	procedureModel := model.ProcedureJavaBasicStaged("w", id, dataType, handler, stage.ID().FullyQualifiedName(), tmpJavaProcedure.JarName).
		WithArgument(argName, dataType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.ProcedureJava),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, procedureModel),
				Check: assertThat(t,
					resourceassert.ProcedureJavaResource(t, procedureModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanDefault).
						HasCommentString(sdk.DefaultProcedureComment).
						HasImportsLength(1).
						HasNoProcedureDefinition().
						HasProcedureLanguageString("JAVA").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "imports.0.stage_location", stage.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "imports.0.path_on_stage", tmpJavaProcedure.JarName)),
					resourceshowoutputassert.ProcedureShowOutput(t, procedureModel.ResourceReference()).
						HasIsSecure(false),
				),
			},
		},
	})
}

func TestAcc_ProcedureJava_AllParameters(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := acc.TestClient().Procedure.SampleJavaDefinition(t, className, funcName, argName)

	procedureModel := model.ProcedureJavaBasicInline("w", id, dataType, handler, definition).
		WithArgument(argName, dataType)
	procedureModelWithAllParametersSet := model.ProcedureJavaBasicInline("w", id, dataType, handler, definition).
		WithArgument(argName, dataType).
		WithEnableConsoleOutput(true).
		WithLogLevel(string(sdk.LogLevelWarn)).
		WithMetricLevel(string(sdk.MetricLevelAll)).
		WithTraceLevel(string(sdk.TraceLevelAlways))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.ProcedureJava),
		Steps: []resource.TestStep{
			// create with default values for all the parameters
			{
				Config: config.FromModels(t, procedureModel),
				Check: assertThat(t,
					objectparametersassert.ProcedureParameters(t, id).
						HasAllDefaults().
						HasAllDefaultsExplicit(),
					resourceparametersassert.ProcedureResourceParameters(t, procedureModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
			// import when no parameter set
			{
				ResourceName: procedureModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceparametersassert.ImportedProcedureResourceParameters(t, helpers.EncodeResourceIdentifier(id)).
						HasAllDefaults(),
				),
			},
			// set all parameters
			{
				Config: config.FromModels(t, procedureModelWithAllParametersSet),
				Check: assertThat(t,
					objectparametersassert.ProcedureParameters(t, id).
						HasEnableConsoleOutput(true).
						HasLogLevel(sdk.LogLevelWarn).
						HasMetricLevel(sdk.MetricLevelAll).
						HasTraceLevel(sdk.TraceLevelAlways),
					resourceparametersassert.ProcedureResourceParameters(t, procedureModelWithAllParametersSet.ResourceReference()).
						HasEnableConsoleOutput(true).
						HasLogLevel(sdk.LogLevelWarn).
						HasMetricLevel(sdk.MetricLevelAll).
						HasTraceLevel(sdk.TraceLevelAlways),
				),
			},
			// import when all parameters set
			{
				ResourceName: procedureModelWithAllParametersSet.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceparametersassert.ImportedProcedureResourceParameters(t, helpers.EncodeResourceIdentifier(id)).
						HasEnableConsoleOutput(true).
						HasLogLevel(sdk.LogLevelWarn).
						HasMetricLevel(sdk.MetricLevelAll).
						HasTraceLevel(sdk.TraceLevelAlways),
				),
			},
			// unset all the parameters
			{
				Config: config.FromModels(t, procedureModel),
				Check: assertThat(t,
					objectparametersassert.ProcedureParameters(t, id).
						HasAllDefaults().
						HasAllDefaultsExplicit(),
					resourceparametersassert.ProcedureResourceParameters(t, procedureModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
			// destroy
			{
				Config:  config.FromModels(t, procedureModel),
				Destroy: true,
			},
			// create with all parameters set
			{
				Config: config.FromModels(t, procedureModelWithAllParametersSet),
				Check: assertThat(t,
					objectparametersassert.ProcedureParameters(t, id).
						HasEnableConsoleOutput(true).
						HasLogLevel(sdk.LogLevelWarn).
						HasMetricLevel(sdk.MetricLevelAll).
						HasTraceLevel(sdk.TraceLevelAlways),
					resourceparametersassert.ProcedureResourceParameters(t, procedureModelWithAllParametersSet.ResourceReference()).
						HasEnableConsoleOutput(true).
						HasLogLevel(sdk.LogLevelWarn).
						HasMetricLevel(sdk.MetricLevelAll).
						HasTraceLevel(sdk.TraceLevelAlways),
				),
			},
		},
	})
}

func TestAcc_ProcedureJava_handleExternalLanguageChange(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tmpJavaProcedure := acc.TestClient().CreateSampleJavaProcedureAndJarOnUserStage(t)

	dataType := tmpJavaProcedure.ArgType
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	argName := "x"
	handler := tmpJavaProcedure.JavaHandler()

	procedureModel := model.ProcedureJavaBasicStaged("w", id, dataType, handler, "~", tmpJavaProcedure.JarName).
		WithArgument(argName, dataType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.ProcedureJava),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, procedureModel),
				Check: assertThat(t,
					resourceassert.ProcedureJavaResource(t, procedureModel.ResourceReference()).HasNameString(id.Name()).HasProcedureLanguageString("JAVA"),
				),
			},
			// change type externally by creating a new procedure with the exact same id but using different language
			{
				PreConfig: func() {
					acc.TestClient().Procedure.DropProcedureFunc(t, id)()
					acc.TestClient().Procedure.CreateScalaStaged(t, id, dataType, tmpJavaProcedure.JarLocation(), handler)
				},
				Config: config.FromModels(t, procedureModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(procedureModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.ProcedureJavaResource(t, procedureModel.ResourceReference()).HasNameString(id.Name()).HasProcedureLanguageString("JAVA"),
				),
			},
		},
	})
}

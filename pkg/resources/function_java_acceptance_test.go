//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"testing"
	"time"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/objectassert"
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

func TestAcc_FunctionJava_InlineBasic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)
	idWithChangedNameButTheSameDataType := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := acc.TestClient().Function.SampleJavaDefinition(t, className, funcName, argName)

	functionModel := model.FunctionJavaBasicInline("w", id, dataType, handler, definition).
		WithArgument(argName, dataType)
	functionModelRenamed := model.FunctionJavaBasicInline("w", idWithChangedNameButTheSameDataType, dataType, handler, definition).
		WithArgument(argName, dataType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionJava),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, functionModel),
				Check: assertThat(t,
					resourceassert.FunctionJavaResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanDefault).
						HasCommentString(sdk.DefaultFunctionComment).
						HasImportsLength(0).
						HasTargetPathEmpty().
						HasNoRuntimeVersion().
						HasFunctionDefinitionString(definition).
						HasFunctionLanguageString("JAVA").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.FunctionShowOutput(t, functionModel.ResourceReference()).
						HasIsSecure(false),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_name", argName)),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_data_type", dataType.ToSql())),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_default_value", "")),
				),
			},
			// REMOVE EXTERNALLY (CHECK RECREATION)
			{
				PreConfig: func() {
					acc.TestClient().Function.DropFunctionFunc(t, id)()
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(functionModel.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, functionModel),
				Check: assertThat(t,
					resourceassert.FunctionJavaResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()),
				),
			},
			// IMPORT
			{
				ResourceName:            functionModel.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_secure", "arguments.0.arg_data_type", "null_input_behavior", "return_results_behavior"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedFunctionJavaResource(t, id.FullyQualifiedName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_name", argName)),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_data_type", "VARCHAR(16777216)")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_default_value", "")),
				),
			},
			// RENAME
			{
				Config: config.FromModels(t, functionModelRenamed),
				Check: assertThat(t,
					resourceassert.FunctionJavaResource(t, functionModelRenamed.ResourceReference()).
						HasNameString(idWithChangedNameButTheSameDataType.Name()).
						HasFullyQualifiedNameString(idWithChangedNameButTheSameDataType.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_FunctionJava_InlineEmptyArgs(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	className := "TestFunc"
	funcName := "echoVarchar"
	returnDataType := testdatatypes.DataTypeVarchar_100

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes()

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := acc.TestClient().Function.SampleJavaDefinitionNoArgs(t, className, funcName)

	functionModel := model.FunctionJavaBasicInline("w", id, returnDataType, handler, definition)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionJava),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, functionModel),
				Check: assertThat(t,
					resourceassert.FunctionJavaResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasFunctionDefinitionString(definition).
						HasFunctionLanguageString("JAVA").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_FunctionJava_InlineBasicDefaultArg(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100
	defaultValue := "'hello'"

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := acc.TestClient().Function.SampleJavaDefinition(t, className, funcName, argName)

	functionModel := model.FunctionJavaBasicInline("w", id, dataType, handler, definition).
		WithArgumentWithDefaultValue(argName, dataType, defaultValue)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionJava),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, functionModel),
				Check: assertThat(t,
					resourceassert.FunctionJavaResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasFunctionDefinitionString(definition).
						HasFunctionLanguageString("JAVA").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_name", argName)),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_data_type", dataType.ToSql())),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_default_value", defaultValue)),
				),
			},
		},
	})
}

func TestAcc_FunctionJava_InlineFull(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	t.Setenv(string(testenvs.ConfigureClientOnce), "")

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

	tmpJavaFunction := acc.TestClient().CreateSampleJavaFunctionAndJarOnUserStage(t)
	tmpJavaFunction2 := acc.TestClient().CreateSampleJavaFunctionAndJarOnUserStage(t)

	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := acc.TestClient().Function.SampleJavaDefinition(t, className, funcName, argName)
	// TODO [SNOW-1850370]: extract to helper
	jarName := fmt.Sprintf("tf-%d-%s.jar", time.Now().Unix(), random.AlphaN(5))

	functionModel := model.FunctionJavaBasicInline("w", id, dataType, handler, definition).
		WithArgument(argName, dataType).
		WithImports(
			sdk.NormalizedPath{StageLocation: "~", PathOnStage: tmpJavaFunction.JarName},
			sdk.NormalizedPath{StageLocation: "~", PathOnStage: tmpJavaFunction2.JarName},
		).
		WithPackages("com.snowflake:snowpark:1.14.0", "com.snowflake:telemetry:0.1.0").
		WithExternalAccessIntegrations(externalAccessIntegration, externalAccessIntegration2).
		WithSecrets(map[string]sdk.SchemaObjectIdentifier{
			"abc": secretId,
			"def": secretId2,
		}).
		WithTargetPathParts(stage.ID().FullyQualifiedName(), jarName).
		WithRuntimeVersion("11").
		WithIsSecure(r.BooleanFalse).
		WithNullInputBehavior(string(sdk.NullInputBehaviorCalledOnNullInput)).
		WithReturnResultsBehavior(string(sdk.ReturnResultsBehaviorVolatile)).
		WithComment("some comment")

	functionModelUpdateWithoutRecreation := model.FunctionJavaBasicInline("w", id, dataType, handler, definition).
		WithArgument(argName, dataType).
		WithImports(
			sdk.NormalizedPath{StageLocation: "~", PathOnStage: tmpJavaFunction.JarName},
			sdk.NormalizedPath{StageLocation: "~", PathOnStage: tmpJavaFunction2.JarName},
		).
		WithPackages("com.snowflake:snowpark:1.14.0", "com.snowflake:telemetry:0.1.0").
		WithExternalAccessIntegrations(externalAccessIntegration).
		WithSecrets(map[string]sdk.SchemaObjectIdentifier{
			"def": secretId2,
		}).
		WithTargetPathParts(stage.ID().FullyQualifiedName(), jarName).
		WithRuntimeVersion("11").
		WithIsSecure(r.BooleanFalse).
		WithNullInputBehavior(string(sdk.NullInputBehaviorCalledOnNullInput)).
		WithReturnResultsBehavior(string(sdk.ReturnResultsBehaviorVolatile)).
		WithComment("some other comment")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionJava),
		Steps: []resource.TestStep{
			// CREATE WITH ALL
			{
				Config: config.FromModels(t, functionModel),
				Check: assertThat(t,
					resourceassert.FunctionJavaResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanFalse).
						HasImportsLength(2).
						HasRuntimeVersionString("11").
						HasFunctionDefinitionString(definition).
						HasCommentString("some comment").
						HasFunctionLanguageString("JAVA").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "target_path.0.stage_location", stage.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "target_path.0.path_on_stage", jarName)),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "secrets.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "external_access_integrations.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "packages.#", "2")),
					resourceshowoutputassert.FunctionShowOutput(t, functionModel.ResourceReference()).
						HasIsSecure(false),
				),
			},
			// IMPORT
			{
				ResourceName:            functionModel.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"arguments.0.arg_data_type"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedFunctionJavaResource(t, id.FullyQualifiedName()).
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
						plancheck.ExpectResourceAction(functionModelUpdateWithoutRecreation.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, functionModelUpdateWithoutRecreation),
				Check: assertThat(t,
					resourceassert.FunctionJavaResource(t, functionModelUpdateWithoutRecreation.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanFalse).
						HasImportsLength(2).
						HasRuntimeVersionString("11").
						HasFunctionDefinitionString(definition).
						HasCommentString("some other comment").
						HasFunctionLanguageString("JAVA").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(functionModelUpdateWithoutRecreation.ResourceReference(), "target_path.0.stage_location", stage.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(functionModelUpdateWithoutRecreation.ResourceReference(), "target_path.0.path_on_stage", jarName)),
					assert.Check(resource.TestCheckResourceAttr(functionModelUpdateWithoutRecreation.ResourceReference(), "secrets.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(functionModelUpdateWithoutRecreation.ResourceReference(), "secrets.0.secret_variable_name", "def")),
					assert.Check(resource.TestCheckResourceAttr(functionModelUpdateWithoutRecreation.ResourceReference(), "secrets.0.secret_id", secretId2.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(functionModelUpdateWithoutRecreation.ResourceReference(), "external_access_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(functionModelUpdateWithoutRecreation.ResourceReference(), "external_access_integrations.0", externalAccessIntegration.Name())),
					assert.Check(resource.TestCheckResourceAttr(functionModelUpdateWithoutRecreation.ResourceReference(), "packages.#", "2")),
					resourceshowoutputassert.FunctionShowOutput(t, functionModelUpdateWithoutRecreation.ResourceReference()).
						HasIsSecure(false),
				),
			},
		},
	})
}

func TestAcc_FunctionJava_StagedBasic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	stage, stageCleanup := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	tmpJavaFunction := acc.TestClient().CreateSampleJavaFunctionAndJarOnStage(t, stage)

	dataType := tmpJavaFunction.ArgType
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	argName := "x"
	handler := tmpJavaFunction.JavaHandler()

	functionModel := model.FunctionJavaBasicStaged("w", id, dataType, handler, stage.ID().FullyQualifiedName(), tmpJavaFunction.JarName).
		WithArgument(argName, dataType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionJava),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, functionModel),
				Check: assertThat(t,
					resourceassert.FunctionJavaResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanDefault).
						HasCommentString(sdk.DefaultFunctionComment).
						HasImportsLength(1).
						HasNoFunctionDefinition().
						HasFunctionLanguageString("JAVA").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "imports.0.stage_location", stage.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "imports.0.path_on_stage", tmpJavaFunction.JarName)),
					resourceshowoutputassert.FunctionShowOutput(t, functionModel.ResourceReference()).
						HasIsSecure(false),
				),
			},
		},
	})
}

// TODO [SNOW-1991414]: set these parameters on the pre-created db level
func TestAcc_FunctionJava_AllParameters(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := acc.TestClient().Function.SampleJavaDefinition(t, className, funcName, argName)

	functionModel := model.FunctionJavaBasicInline("w", id, dataType, handler, definition).
		WithArgument(argName, dataType)
	functionModelWithAllParametersSet := model.FunctionJavaBasicInline("w", id, dataType, handler, definition).
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
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionJava),
		Steps: []resource.TestStep{
			// create with default values for all the parameters
			{
				Config: config.FromModels(t, functionModel),
				Check: assertThat(t,
					objectparametersassert.FunctionParameters(t, id).
						HasAllDefaults().
						HasAllDefaultsExplicit(),
					resourceparametersassert.FunctionResourceParameters(t, functionModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
			// import when no parameter set
			{
				ResourceName: functionModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceparametersassert.ImportedFunctionResourceParameters(t, helpers.EncodeResourceIdentifier(id)).
						HasAllDefaults(),
				),
			},
			// set all parameters
			{
				Config: config.FromModels(t, functionModelWithAllParametersSet),
				Check: assertThat(t,
					objectparametersassert.FunctionParameters(t, id).
						HasEnableConsoleOutput(true).
						HasLogLevel(sdk.LogLevelWarn).
						HasMetricLevel(sdk.MetricLevelAll).
						HasTraceLevel(sdk.TraceLevelAlways),
					resourceparametersassert.FunctionResourceParameters(t, functionModelWithAllParametersSet.ResourceReference()).
						HasEnableConsoleOutput(true).
						HasLogLevel(sdk.LogLevelWarn).
						HasMetricLevel(sdk.MetricLevelAll).
						HasTraceLevel(sdk.TraceLevelAlways),
				),
			},
			// import when all parameters set
			{
				ResourceName: functionModelWithAllParametersSet.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceparametersassert.ImportedFunctionResourceParameters(t, helpers.EncodeResourceIdentifier(id)).
						HasEnableConsoleOutput(true).
						HasLogLevel(sdk.LogLevelWarn).
						HasMetricLevel(sdk.MetricLevelAll).
						HasTraceLevel(sdk.TraceLevelAlways),
				),
			},
			// unset all the parameters
			{
				Config: config.FromModels(t, functionModel),
				Check: assertThat(t,
					objectparametersassert.FunctionParameters(t, id).
						HasAllDefaults().
						HasAllDefaultsExplicit(),
					resourceparametersassert.FunctionResourceParameters(t, functionModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
			// destroy
			{
				Config:  config.FromModels(t, functionModel),
				Destroy: true,
			},
			// create with all parameters set
			{
				Config: config.FromModels(t, functionModelWithAllParametersSet),
				Check: assertThat(t,
					objectparametersassert.FunctionParameters(t, id).
						HasEnableConsoleOutput(true).
						HasLogLevel(sdk.LogLevelWarn).
						HasMetricLevel(sdk.MetricLevelAll).
						HasTraceLevel(sdk.TraceLevelAlways),
					resourceparametersassert.FunctionResourceParameters(t, functionModelWithAllParametersSet.ResourceReference()).
						HasEnableConsoleOutput(true).
						HasLogLevel(sdk.LogLevelWarn).
						HasMetricLevel(sdk.MetricLevelAll).
						HasTraceLevel(sdk.TraceLevelAlways),
				),
			},
		},
	})
}

func TestAcc_FunctionJava_handleExternalLanguageChange(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tmpJavaFunction := acc.TestClient().CreateSampleJavaFunctionAndJarOnUserStage(t)

	dataType := tmpJavaFunction.ArgType
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	argName := "x"
	handler := tmpJavaFunction.JavaHandler()

	functionModel := model.FunctionJavaBasicStaged("w", id, dataType, handler, "~", tmpJavaFunction.JarName).
		WithArgument(argName, dataType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionJava),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, functionModel),
				Check: assertThat(t,
					objectassert.Function(t, id).HasLanguage("JAVA"),
					resourceassert.FunctionJavaResource(t, functionModel.ResourceReference()).HasNameString(id.Name()).HasFunctionLanguageString("JAVA"),
					resourceshowoutputassert.FunctionShowOutput(t, functionModel.ResourceReference()).HasLanguage("JAVA"),
				),
			},
			// change type externally by creating a new function with the exact same id but using different language
			{
				PreConfig: func() {
					acc.TestClient().Function.DropFunctionFunc(t, id)()
					acc.TestClient().Function.CreateScalaStaged(t, id, dataType, tmpJavaFunction.JarLocation(), handler)
					objectassert.Function(t, id).HasLanguage("SCALA")
				},
				Config: config.FromModels(t, functionModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(functionModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					objectassert.Function(t, id).HasLanguage("JAVA"),
					resourceassert.FunctionJavaResource(t, functionModel.ResourceReference()).HasNameString(id.Name()).HasFunctionLanguageString("JAVA"),
					resourceshowoutputassert.FunctionShowOutput(t, functionModel.ResourceReference()).HasLanguage("JAVA"),
				),
			},
		},
	})
}

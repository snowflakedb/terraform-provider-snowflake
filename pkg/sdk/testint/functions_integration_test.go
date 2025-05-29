//go:build !account_level_tests

package testint

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-1850370]: HasArgumentsRawFrom(functionId, arguments, return)
// TODO [SNOW-1850370]: extract show assertions with commons fields
// TODO [SNOW-1850370]: test confirming that runtime version is required for Scala function
// TODO [SNOW-1850370]: test create or replace with name change, args change
// TODO [SNOW-1850370]: test rename more (arg stays, can't change arg, rename to different schema)
// TODO [SNOW-1850370]: add test documenting that UNSET SECRETS does not work
// TODO [SNOW-1850370]: add test documenting [JAVA]: 391516 (42601): SQL compilation error: Cannot specify TARGET_PATH without a function BODY.
// TODO [SNOW-1850370]: add a test documenting that we can't set parameters in create (and revert adding these parameters directly in object...)
// TODO [SNOW-1850370]: active warehouse vs validations
// TODO [SNOW-1850370]: add a test documenting STRICT behavior
// TODO [SNOW-1348103]: test weird names for arg name - lower/upper if used with double quotes, to upper without quotes, dots, spaces, and both quotes not permitted
// TODO [SNOW-1348103]: test secure
// TODO [SNOW-1348103]: python aggregate func (100357 (P0000): Could not find accumulate method in function CVVEMHIT_06547800_08D6_DBCA_1AC7_5E422AFF8B39 with handler dump)
// TODO [SNOW-1348103]: add test with multiple imports
// TODO [SNOW-1348103]: test with multiple external access integrations and secrets
func TestInt_Functions(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	secretId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

	networkRule, networkRuleCleanup := testClientHelper().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)

	secret, secretCleanup := testClientHelper().Secret.CreateWithGenericString(t, secretId, "test_secret_string")
	t.Cleanup(secretCleanup)

	externalAccessIntegration, externalAccessIntegrationCleanup := testClientHelper().ExternalAccessIntegration.CreateExternalAccessIntegrationWithNetworkRuleAndSecret(t, networkRule.ID(), secret.ID())
	t.Cleanup(externalAccessIntegrationCleanup)

	tmpJavaFunction := testClientHelper().CreateSampleJavaFunctionAndJarOnUserStage(t)
	tmpPythonFunction := testClientHelper().CreateSamplePythonFunctionAndModuleOnUserStage(t)

	assertParametersSet := func(t *testing.T, functionParametersAssert *objectparametersassert.FunctionParametersAssert) {
		t.Helper()
		assertThatObject(t, functionParametersAssert.
			HasEnableConsoleOutput(true).
			HasLogLevel(sdk.LogLevelWarn).
			HasMetricLevel(sdk.MetricLevelAll).
			HasTraceLevel(sdk.TraceLevelAlways),
		)
	}

	t.Run("create function for Java - inline minimal", func(t *testing.T) {
		className := "TestFunc"
		funcName := "echoVarchar"
		argName := "x"
		dataType := testdatatypes.DataTypeVarchar_100

		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))
		argument := sdk.NewFunctionArgumentRequest(argName, dataType)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		handler := fmt.Sprintf("%s.%s", className, funcName)
		definition := testClientHelper().Function.SampleJavaDefinition(t, className, funcName, argName)

		request := sdk.NewCreateForJavaFunctionRequest(id.SchemaObjectId(), *returns, handler).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithFunctionDefinitionWrapped(definition)

		err := client.Functions.CreateForJava(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultFunctionComment).
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrations("").
			HasSecrets("").
			HasIsExternalFunction(false).
			HasLanguage("JAVA").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(dataType.ToSql()).
			HasReturnDataType(dataType).
			HasReturnNotNull(false).
			HasLanguage("JAVA").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder().
			HasSecretsNil().
			HasImports(`[]`).
			HasExactlyImportsNormalizedInAnyOrder().
			HasHandler(handler).
			HasRuntimeVersionNil().
			HasPackages(`[]`).
			HasExactlyPackagesInAnyOrder().
			HasTargetPathNil().
			HasNormalizedTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create function for Java - inline full", func(t *testing.T) {
		className := "TestFunc"
		funcName := "echoVarchar"
		argName := "x"
		dataType := testdatatypes.DataTypeVarchar_100

		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))
		argument := sdk.NewFunctionArgumentRequest(argName, dataType)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		handler := fmt.Sprintf("%s.%s", className, funcName)
		definition := testClientHelper().Function.SampleJavaDefinition(t, className, funcName, argName)
		jarName := fmt.Sprintf("tf-%d-%s.jar", time.Now().Unix(), random.AlphaN(5))
		targetPath := fmt.Sprintf("@~/%s", jarName)

		request := sdk.NewCreateForJavaFunctionRequest(id.SchemaObjectId(), *returns, handler).
			WithOrReplace(true).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnsNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithReturnNullValues(sdk.ReturnNullValuesNotNull).
			WithRuntimeVersion("11").
			WithComment("comment").
			WithImports([]sdk.FunctionImportRequest{*sdk.NewFunctionImportRequest().WithImport(tmpJavaFunction.JarLocation())}).
			WithPackages([]sdk.FunctionPackageRequest{
				*sdk.NewFunctionPackageRequest().WithPackage("com.snowflake:snowpark:1.14.0"),
				*sdk.NewFunctionPackageRequest().WithPackage("com.snowflake:telemetry:0.1.0"),
			}).
			WithExternalAccessIntegrations([]sdk.AccountObjectIdentifier{externalAccessIntegration}).
			WithSecrets([]sdk.SecretReference{{VariableName: "abc", Name: secretId}}).
			WithTargetPath(targetPath).
			WithEnableConsoleOutput(true).
			WithLogLevel(sdk.LogLevelWarn).
			WithMetricLevel(sdk.MetricLevelAll).
			WithTraceLevel(sdk.TraceLevelAlways).
			WithFunctionDefinitionWrapped(definition)

		err := client.Functions.CreateForJava(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))
		t.Cleanup(testClientHelper().Stage.RemoveFromUserStageFunc(t, jarName))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			HasIsExternalFunction(false).
			HasLanguage("JAVA").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(fmt.Sprintf(`%s NOT NULL`, dataType.ToSql())).
			HasReturnDataType(dataType).
			HasReturnNotNull(true).
			HasLanguage("JAVA").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorReturnsNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			ContainsExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			HasImports(fmt.Sprintf(`[%s]`, tmpJavaFunction.JarLocation())).
			HasExactlyImportsNormalizedInAnyOrder(sdk.NormalizedPath{
				StageLocation: "~", PathOnStage: tmpJavaFunction.JarName,
			}).
			HasHandler(handler).
			HasRuntimeVersion("11").
			HasPackages(`[com.snowflake:snowpark:1.14.0,com.snowflake:telemetry:0.1.0]`).
			HasExactlyPackagesInAnyOrder("com.snowflake:snowpark:1.14.0", "com.snowflake:telemetry:0.1.0").
			HasTargetPath(targetPath).
			HasNormalizedTargetPath("~", jarName).
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create function for Java - staged minimal", func(t *testing.T) {
		dataType := tmpJavaFunction.ArgType
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "x"
		argument := sdk.NewFunctionArgumentRequest(argName, dataType)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		handler := tmpJavaFunction.JavaHandler()
		importPath := tmpJavaFunction.JarLocation()

		requestStaged := sdk.NewCreateForJavaFunctionRequest(id.SchemaObjectId(), *returns, handler).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithImports([]sdk.FunctionImportRequest{*sdk.NewFunctionImportRequest().WithImport(importPath)})

		err := client.Functions.CreateForJava(ctx, requestStaged)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultFunctionComment).
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrations("").
			HasSecrets("").
			HasIsExternalFunction(false).
			HasLanguage("JAVA").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(dataType.ToSql()).
			HasReturnDataType(dataType).
			HasReturnNotNull(false).
			HasLanguage("JAVA").
			HasBodyNil().
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder().
			HasSecretsNil().
			HasImports(fmt.Sprintf(`[%s]`, importPath)).
			HasExactlyImportsNormalizedInAnyOrder(sdk.NormalizedPath{
				StageLocation: "~", PathOnStage: tmpJavaFunction.JarName,
			}).
			HasHandler(handler).
			HasRuntimeVersionNil().
			HasPackages(`[]`).
			HasExactlyPackagesInAnyOrder().
			HasTargetPathNil().
			HasNormalizedTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create function for Java - staged full", func(t *testing.T) {
		dataType := tmpJavaFunction.ArgType
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "x"
		argument := sdk.NewFunctionArgumentRequest(argName, dataType)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		handler := tmpJavaFunction.JavaHandler()

		requestStaged := sdk.NewCreateForJavaFunctionRequest(id.SchemaObjectId(), *returns, handler).
			WithOrReplace(true).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnsNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithReturnNullValues(sdk.ReturnNullValuesNotNull).
			WithRuntimeVersion("11").
			WithComment("comment").
			WithImports([]sdk.FunctionImportRequest{*sdk.NewFunctionImportRequest().WithImport(tmpJavaFunction.JarLocation())}).
			WithPackages([]sdk.FunctionPackageRequest{
				*sdk.NewFunctionPackageRequest().WithPackage("com.snowflake:snowpark:1.14.0"),
				*sdk.NewFunctionPackageRequest().WithPackage("com.snowflake:telemetry:0.1.0"),
			}).
			WithExternalAccessIntegrations([]sdk.AccountObjectIdentifier{externalAccessIntegration}).
			WithSecrets([]sdk.SecretReference{{VariableName: "abc", Name: secretId}})

		err := client.Functions.CreateForJava(ctx, requestStaged)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			HasIsExternalFunction(false).
			HasLanguage("JAVA").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(fmt.Sprintf(`%s NOT NULL`, dataType.ToSql())).
			HasReturnDataType(dataType).
			HasReturnNotNull(true).
			HasLanguage("JAVA").
			HasBodyNil().
			HasNullHandling(string(sdk.NullInputBehaviorReturnsNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			ContainsExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			HasImports(fmt.Sprintf(`[%s]`, tmpJavaFunction.JarLocation())).
			HasExactlyImportsNormalizedInAnyOrder(sdk.NormalizedPath{
				StageLocation: "~", PathOnStage: tmpJavaFunction.JarName,
			}).
			HasHandler(handler).
			HasRuntimeVersion("11").
			HasPackages(`[com.snowflake:snowpark:1.14.0,com.snowflake:telemetry:0.1.0]`).
			HasExactlyPackagesInAnyOrder("com.snowflake:snowpark:1.14.0", "com.snowflake:telemetry:0.1.0").
			HasTargetPathNil().
			HasNormalizedTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create function for Java - different stage", func(t *testing.T) {
		stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(stageCleanup)

		tmpJavaFunctionDifferentStage := testClientHelper().CreateSampleJavaFunctionAndJarOnStage(t, stage)

		dataType := tmpJavaFunctionDifferentStage.ArgType
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

		argName := "x"
		argument := sdk.NewFunctionArgumentRequest(argName, dataType)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		handler := tmpJavaFunctionDifferentStage.JavaHandler()
		importPath := tmpJavaFunctionDifferentStage.JarLocation()

		requestStaged := sdk.NewCreateForJavaFunctionRequest(id.SchemaObjectId(), *returns, handler).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithImports([]sdk.FunctionImportRequest{*sdk.NewFunctionImportRequest().WithImport(importPath)})

		err := client.Functions.CreateForJava(ctx, requestStaged)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasImports(fmt.Sprintf(`[@"%s"."%s".%s/%s]`, stage.ID().DatabaseName(), stage.ID().SchemaName(), stage.ID().Name(), tmpJavaFunctionDifferentStage.JarName)).
			HasExactlyImportsNormalizedInAnyOrder(sdk.NormalizedPath{
				StageLocation: stage.ID().FullyQualifiedName(), PathOnStage: tmpJavaFunctionDifferentStage.JarName,
			}).
			HasHandler(handler).
			HasTargetPathNil().
			HasNormalizedTargetPathNil(),
		)
	})

	// proves that we don't get default argument values from SHOW and DESCRIBE
	t.Run("create function for Java - default argument value", func(t *testing.T) {
		className := "TestFunc"
		funcName := "echoVarchar"
		argName := "x"
		dataType := testdatatypes.DataTypeVarchar_100

		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))
		argument := sdk.NewFunctionArgumentRequest(argName, dataType).WithDefaultValue(`'abc'`)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		handler := fmt.Sprintf("%s.%s", className, funcName)
		definition := testClientHelper().Function.SampleJavaDefinition(t, className, funcName, argName)

		request := sdk.NewCreateForJavaFunctionRequest(id.SchemaObjectId(), *returns, handler).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithFunctionDefinitionWrapped(definition)

		err := client.Functions.CreateForJava(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(DEFAULT %[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())),
		)
	})

	t.Run("create function for Javascript - inline minimal", func(t *testing.T) {
		dataType := testdatatypes.DataTypeFloat
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "d"
		definition := testClientHelper().Function.SampleJavascriptDefinition(t, argName)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewFunctionArgumentRequest(argName, dataType)

		request := sdk.NewCreateForJavascriptFunctionRequestDefinitionWrapped(id.SchemaObjectId(), *returns, definition).
			WithArguments([]sdk.FunctionArgumentRequest{*argument})

		err := client.Functions.CreateForJavascript(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultFunctionComment).
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrations("").
			HasSecrets("").
			HasIsExternalFunction(false).
			HasLanguage("JAVASCRIPT").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(dataType.ToSql()).
			HasReturnDataType(dataType).
			HasReturnNotNull(false).
			HasLanguage("JAVASCRIPT").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder().
			HasSecretsNil().
			HasImportsNil().
			HasExactlyImportsNormalizedInAnyOrder().
			HasHandlerNil().
			HasRuntimeVersionNil().
			HasPackagesNil().
			HasExactlyPackagesInAnyOrder().
			HasTargetPathNil().
			HasNormalizedTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create function for Javascript - inline full", func(t *testing.T) {
		dataType := testdatatypes.DataTypeFloat
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "d"
		definition := testClientHelper().Function.SampleJavascriptDefinition(t, argName)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewFunctionArgumentRequest(argName, dataType)
		request := sdk.NewCreateForJavascriptFunctionRequestDefinitionWrapped(id.SchemaObjectId(), *returns, definition).
			WithOrReplace(true).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithReturnNullValues(sdk.ReturnNullValuesNotNull).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnsNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithComment("comment")

		err := client.Functions.CreateForJavascript(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrations("").
			HasSecrets("").
			HasIsExternalFunction(false).
			HasLanguage("JAVASCRIPT").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(fmt.Sprintf(`%s NOT NULL`, dataType.ToSql())).
			HasReturnDataType(dataType).
			HasReturnNotNull(true).
			HasLanguage("JAVASCRIPT").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorReturnsNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExternalAccessIntegrationsNil().
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder().
			HasSecretsNil().
			HasImportsNil().
			HasExactlyImportsNormalizedInAnyOrder().
			HasHandlerNil().
			HasRuntimeVersionNil().
			HasPackagesNil().
			HasExactlyPackagesInAnyOrder().
			HasTargetPathNil().
			HasNormalizedTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create function for Python - inline minimal", func(t *testing.T) {
		dataType := testdatatypes.DataTypeNumber_36_2
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "i"
		funcName := "dump"
		definition := testClientHelper().Function.SamplePythonDefinition(t, funcName, argName)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewFunctionArgumentRequest(argName, dataType)
		request := sdk.NewCreateForPythonFunctionRequest(id.SchemaObjectId(), *returns, testvars.PythonRuntime, funcName).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithFunctionDefinitionWrapped(definition)

		err := client.Functions.CreateForPython(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultFunctionComment).
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrations("").
			HasSecrets("").
			HasIsExternalFunction(false).
			HasLanguage("PYTHON").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(strings.ReplaceAll(dataType.ToSql(), " ", "")).
			HasReturnDataType(dataType).
			HasReturnNotNull(false).
			HasLanguage("PYTHON").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder().
			HasSecretsNil().
			HasImports(`[]`).
			HasExactlyImportsNormalizedInAnyOrder().
			HasHandler(funcName).
			HasRuntimeVersion(testvars.PythonRuntime).
			HasPackages(`[]`).
			HasExactlyPackagesInAnyOrder().
			HasTargetPathNil().
			HasNormalizedTargetPathNil().
			HasInstalledPackagesNotEmpty().
			HasIsAggregate(false),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create function for Python - inline full", func(t *testing.T) {
		dataType := testdatatypes.DataTypeNumber_36_2
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "i"
		funcName := "dump"
		definition := testClientHelper().Function.SamplePythonDefinition(t, funcName, argName)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewFunctionArgumentRequest(argName, dataType)
		request := sdk.NewCreateForPythonFunctionRequest(id.SchemaObjectId(), *returns, testvars.PythonRuntime, funcName).
			WithOrReplace(true).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithReturnNullValues(sdk.ReturnNullValuesNotNull).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnsNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithComment("comment").
			WithImports([]sdk.FunctionImportRequest{*sdk.NewFunctionImportRequest().WithImport(tmpPythonFunction.PythonModuleLocation())}).
			WithPackages([]sdk.FunctionPackageRequest{
				*sdk.NewFunctionPackageRequest().WithPackage("absl-py==0.12.0"),
				*sdk.NewFunctionPackageRequest().WithPackage("about-time==4.2.1"),
			}).
			WithExternalAccessIntegrations([]sdk.AccountObjectIdentifier{externalAccessIntegration}).
			WithSecrets([]sdk.SecretReference{{VariableName: "abc", Name: secretId}}).
			WithFunctionDefinitionWrapped(definition)

		err := client.Functions.CreateForPython(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			HasIsExternalFunction(false).
			HasLanguage("PYTHON").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(strings.ReplaceAll(dataType.ToSql(), " ", "")+" NOT NULL").
			HasReturnDataType(dataType).
			HasReturnNotNull(true).
			HasLanguage("PYTHON").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorReturnsNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			ContainsExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			HasImports(fmt.Sprintf(`[%s]`, tmpPythonFunction.PythonModuleLocation())).
			HasExactlyImportsNormalizedInAnyOrder(sdk.NormalizedPath{
				StageLocation: "~", PathOnStage: tmpPythonFunction.PythonFileName(),
			}).
			HasHandler(funcName).
			HasRuntimeVersion(testvars.PythonRuntime).
			HasPackages(`['absl-py==0.12.0','about-time==4.2.1']`).
			HasExactlyPackagesInAnyOrder("absl-py==0.12.0", "about-time==4.2.1").
			HasTargetPathNil().
			HasNormalizedTargetPathNil().
			HasInstalledPackagesNotEmpty().
			HasIsAggregate(false),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create function for Python - staged minimal", func(t *testing.T) {
		dataType := testdatatypes.DataTypeVarchar_100
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "i"
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewFunctionArgumentRequest(argName, dataType)
		request := sdk.NewCreateForPythonFunctionRequest(id.SchemaObjectId(), *returns, testvars.PythonRuntime, tmpPythonFunction.PythonHandler()).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithImports([]sdk.FunctionImportRequest{*sdk.NewFunctionImportRequest().WithImport(tmpPythonFunction.PythonModuleLocation())})

		err := client.Functions.CreateForPython(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultFunctionComment).
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrations("").
			HasSecrets("").
			HasIsExternalFunction(false).
			HasLanguage("PYTHON").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(strings.ReplaceAll(dataType.ToSql(), " ", "")).
			HasReturnDataType(dataType).
			HasReturnNotNull(false).
			HasLanguage("PYTHON").
			HasBodyNil().
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder().
			HasSecretsNil().
			HasImports(fmt.Sprintf(`[%s]`, tmpPythonFunction.PythonModuleLocation())).
			HasExactlyImportsNormalizedInAnyOrder(sdk.NormalizedPath{
				StageLocation: "~", PathOnStage: tmpPythonFunction.PythonFileName(),
			}).
			HasHandler(tmpPythonFunction.PythonHandler()).
			HasRuntimeVersion(testvars.PythonRuntime).
			HasPackages(`[]`).
			HasExactlyPackagesInAnyOrder().
			HasTargetPathNil().
			HasNormalizedTargetPathNil().
			HasInstalledPackagesNotEmpty().
			HasIsAggregate(false),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create function for Python - staged full", func(t *testing.T) {
		dataType := testdatatypes.DataTypeVarchar_100
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "i"
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewFunctionArgumentRequest(argName, dataType)
		request := sdk.NewCreateForPythonFunctionRequest(id.SchemaObjectId(), *returns, testvars.PythonRuntime, tmpPythonFunction.PythonHandler()).
			WithOrReplace(true).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithReturnNullValues(sdk.ReturnNullValuesNotNull).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnsNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithComment("comment").
			WithPackages([]sdk.FunctionPackageRequest{
				*sdk.NewFunctionPackageRequest().WithPackage("absl-py==0.12.0"),
				*sdk.NewFunctionPackageRequest().WithPackage("about-time==4.2.1"),
			}).
			WithExternalAccessIntegrations([]sdk.AccountObjectIdentifier{externalAccessIntegration}).
			WithSecrets([]sdk.SecretReference{{VariableName: "abc", Name: secretId}}).
			WithImports([]sdk.FunctionImportRequest{*sdk.NewFunctionImportRequest().WithImport(tmpPythonFunction.PythonModuleLocation())})

		err := client.Functions.CreateForPython(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			HasIsExternalFunction(false).
			HasLanguage("PYTHON").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(strings.ReplaceAll(dataType.ToSql(), " ", "")+" NOT NULL").
			HasReturnDataType(dataType).
			HasReturnNotNull(true).
			HasLanguage("PYTHON").
			HasBodyNil().
			HasNullHandling(string(sdk.NullInputBehaviorReturnsNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			ContainsExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			HasImports(fmt.Sprintf(`[%s]`, tmpPythonFunction.PythonModuleLocation())).
			HasExactlyImportsNormalizedInAnyOrder(sdk.NormalizedPath{
				StageLocation: "~", PathOnStage: tmpPythonFunction.PythonFileName(),
			}).
			HasHandler(tmpPythonFunction.PythonHandler()).
			HasRuntimeVersion(testvars.PythonRuntime).
			HasPackages(`['absl-py==0.12.0','about-time==4.2.1']`).
			HasExactlyPackagesInAnyOrder("about-time==4.2.1", "absl-py==0.12.0").
			HasTargetPathNil().
			HasNormalizedTargetPathNil().
			HasInstalledPackagesNotEmpty().
			HasIsAggregate(false),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create function for Scala - inline minimal", func(t *testing.T) {
		className := "TestFunc"
		funcName := "echoVarchar"
		argName := "x"
		dataType := testdatatypes.DataTypeVarchar_100

		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		definition := testClientHelper().Function.SampleScalaDefinition(t, className, funcName, argName)
		argument := sdk.NewFunctionArgumentRequest(argName, dataType)
		handler := fmt.Sprintf("%s.%s", className, funcName)
		request := sdk.NewCreateForScalaFunctionRequest(id.SchemaObjectId(), dataType, handler, "2.12").
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithFunctionDefinitionWrapped(definition)

		err := client.Functions.CreateForScala(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultFunctionComment).
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrations("").
			HasSecrets("").
			HasIsExternalFunction(false).
			HasLanguage("SCALA").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(dataType.ToSql()).
			HasReturnDataType(dataType).
			HasReturnNotNull(false).
			HasLanguage("SCALA").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder().
			HasSecretsNil().
			HasImports(`[]`).
			HasExactlyImportsNormalizedInAnyOrder().
			HasHandler(handler).
			HasRuntimeVersion("2.12").
			HasPackages(`[]`).
			HasExactlyPackagesInAnyOrder().
			HasTargetPathNil().
			HasNormalizedTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create function for Scala - inline full", func(t *testing.T) {
		className := "TestFunc"
		funcName := "echoVarchar"
		argName := "x"
		dataType := testdatatypes.DataTypeVarchar_100

		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		definition := testClientHelper().Function.SampleScalaDefinition(t, className, funcName, argName)
		argument := sdk.NewFunctionArgumentRequest(argName, dataType)
		handler := fmt.Sprintf("%s.%s", className, funcName)
		jarName := fmt.Sprintf("tf-%d-%s.jar", time.Now().Unix(), random.AlphaN(5))
		targetPath := fmt.Sprintf("@~/%s", jarName)
		request := sdk.NewCreateForScalaFunctionRequest(id.SchemaObjectId(), dataType, handler, "2.12").
			WithOrReplace(true).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnsNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithReturnNullValues(sdk.ReturnNullValuesNotNull).
			WithComment("comment").
			WithImports([]sdk.FunctionImportRequest{*sdk.NewFunctionImportRequest().WithImport(tmpJavaFunction.JarLocation())}).
			WithPackages([]sdk.FunctionPackageRequest{
				*sdk.NewFunctionPackageRequest().WithPackage("com.snowflake:snowpark:1.14.0"),
				*sdk.NewFunctionPackageRequest().WithPackage("com.snowflake:telemetry:0.1.0"),
			}).
			WithTargetPath(targetPath).
			WithExternalAccessIntegrations([]sdk.AccountObjectIdentifier{externalAccessIntegration}).
			WithSecrets([]sdk.SecretReference{{VariableName: "abc", Name: secretId}}).
			WithEnableConsoleOutput(true).
			WithLogLevel(sdk.LogLevelWarn).
			WithMetricLevel(sdk.MetricLevelAll).
			WithTraceLevel(sdk.TraceLevelAlways).
			WithFunctionDefinitionWrapped(definition)

		err := client.Functions.CreateForScala(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))
		t.Cleanup(testClientHelper().Stage.RemoveFromUserStageFunc(t, jarName))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			HasIsExternalFunction(false).
			HasLanguage("SCALA").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(fmt.Sprintf(`%s NOT NULL`, dataType.ToSql())).
			HasReturnDataType(dataType).
			HasReturnNotNull(true).
			HasLanguage("SCALA").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorReturnsNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			ContainsExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			HasImports(fmt.Sprintf(`[%s]`, tmpJavaFunction.JarLocation())).
			HasExactlyImportsNormalizedInAnyOrder(sdk.NormalizedPath{
				StageLocation: "~", PathOnStage: tmpJavaFunction.JarName,
			}).
			HasHandler(handler).
			HasRuntimeVersion("2.12").
			HasPackages(`[com.snowflake:snowpark:1.14.0,com.snowflake:telemetry:0.1.0]`).
			HasExactlyPackagesInAnyOrder("com.snowflake:snowpark:1.14.0", "com.snowflake:telemetry:0.1.0").
			HasTargetPath(targetPath).
			HasNormalizedTargetPath("~", jarName).
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create function for Scala - staged minimal", func(t *testing.T) {
		dataType := tmpJavaFunction.ArgType
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "x"
		argument := sdk.NewFunctionArgumentRequest(argName, dataType)
		handler := tmpJavaFunction.JavaHandler()
		importPath := tmpJavaFunction.JarLocation()

		requestStaged := sdk.NewCreateForScalaFunctionRequest(id.SchemaObjectId(), dataType, handler, "2.12").
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithImports([]sdk.FunctionImportRequest{*sdk.NewFunctionImportRequest().WithImport(importPath)})

		err := client.Functions.CreateForScala(ctx, requestStaged)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultFunctionComment).
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrations("").
			HasSecrets("").
			HasIsExternalFunction(false).
			HasLanguage("SCALA").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(dataType.ToSql()).
			HasReturnDataType(dataType).
			HasReturnNotNull(false).
			HasLanguage("SCALA").
			HasBodyNil().
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder().
			HasSecretsNil().
			HasImports(fmt.Sprintf(`[%s]`, importPath)).
			HasExactlyImportsNormalizedInAnyOrder(sdk.NormalizedPath{
				StageLocation: "~", PathOnStage: tmpJavaFunction.JarName,
			}).
			HasHandler(handler).
			HasRuntimeVersion("2.12").
			HasPackages(`[]`).
			HasExactlyPackagesInAnyOrder().
			HasTargetPathNil().
			HasNormalizedTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create function for Scala - staged full", func(t *testing.T) {
		dataType := tmpJavaFunction.ArgType
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "x"
		argument := sdk.NewFunctionArgumentRequest(argName, dataType)
		handler := tmpJavaFunction.JavaHandler()

		requestStaged := sdk.NewCreateForScalaFunctionRequest(id.SchemaObjectId(), dataType, handler, "2.12").
			WithOrReplace(true).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnsNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithReturnNullValues(sdk.ReturnNullValuesNotNull).
			WithComment("comment").
			WithPackages([]sdk.FunctionPackageRequest{
				*sdk.NewFunctionPackageRequest().WithPackage("com.snowflake:snowpark:1.14.0"),
				*sdk.NewFunctionPackageRequest().WithPackage("com.snowflake:telemetry:0.1.0"),
			}).
			WithExternalAccessIntegrations([]sdk.AccountObjectIdentifier{externalAccessIntegration}).
			WithSecrets([]sdk.SecretReference{{VariableName: "abc", Name: secretId}}).
			WithImports([]sdk.FunctionImportRequest{*sdk.NewFunctionImportRequest().WithImport(tmpJavaFunction.JarLocation())})

		err := client.Functions.CreateForScala(ctx, requestStaged)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasIsExternalFunction(false).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			HasLanguage("SCALA").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(fmt.Sprintf(`%s NOT NULL`, dataType.ToSql())).
			HasReturnDataType(dataType).
			HasReturnNotNull(true).
			HasLanguage("SCALA").
			HasBodyNil().
			HasNullHandling(string(sdk.NullInputBehaviorReturnsNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			ContainsExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			HasImports(fmt.Sprintf(`[%s]`, tmpJavaFunction.JarLocation())).
			HasExactlyImportsNormalizedInAnyOrder(sdk.NormalizedPath{
				StageLocation: "~", PathOnStage: tmpJavaFunction.JarName,
			}).
			HasHandler(handler).
			HasRuntimeVersion("2.12").
			HasPackages(`[com.snowflake:snowpark:1.14.0,com.snowflake:telemetry:0.1.0]`).
			HasExactlyPackagesInAnyOrder("com.snowflake:snowpark:1.14.0", "com.snowflake:telemetry:0.1.0").
			HasTargetPathNil().
			HasNormalizedTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create function for SQL - inline minimal", func(t *testing.T) {
		argName := "x"
		dataType := testdatatypes.DataTypeFloat
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		definition := testClientHelper().Function.SampleSqlDefinition(t)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewFunctionArgumentRequest(argName, dataType)
		request := sdk.NewCreateForSQLFunctionRequestDefinitionWrapped(id.SchemaObjectId(), *returns, definition).
			WithArguments([]sdk.FunctionArgumentRequest{*argument})

		err := client.Functions.CreateForSQL(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultFunctionComment).
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrations("").
			HasSecrets("").
			HasIsExternalFunction(false).
			HasLanguage("SQL").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(dataType.ToSql()).
			HasReturnDataType(dataType).
			HasReturnNotNull(false).
			HasLanguage("SQL").
			HasBody(definition).
			HasNullHandlingNil().
			HasVolatilityNil().
			HasExternalAccessIntegrationsNil().
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder().
			HasSecretsNil().
			HasImportsNil().
			HasExactlyImportsNormalizedInAnyOrder().
			HasHandlerNil().
			HasRuntimeVersionNil().
			HasPackagesNil().
			HasExactlyPackagesInAnyOrder().
			HasTargetPathNil().
			HasNormalizedTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	// proves that we don't get default argument values from SHOW and DESCRIBE
	t.Run("create function for SQL - default argument value", func(t *testing.T) {
		argName := "x"
		dataType := testdatatypes.DataTypeFloat
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		definition := testClientHelper().Function.SampleSqlDefinition(t)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewFunctionArgumentRequest(argName, dataType).WithDefaultValue("3.123")
		request := sdk.NewCreateForSQLFunctionRequestDefinitionWrapped(id.SchemaObjectId(), *returns, definition).
			WithArguments([]sdk.FunctionArgumentRequest{*argument})

		err := client.Functions.CreateForSQL(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(DEFAULT %[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())),
		)
	})

	t.Run("create function for SQL - inline full", func(t *testing.T) {
		argName := "x"
		dataType := testdatatypes.DataTypeFloat
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		definition := testClientHelper().Function.SampleSqlDefinition(t)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewFunctionArgumentRequest(argName, dataType)
		request := sdk.NewCreateForSQLFunctionRequestDefinitionWrapped(id.SchemaObjectId(), *returns, definition).
			WithOrReplace(true).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithReturnNullValues(sdk.ReturnNullValuesNotNull).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithMemoizable(true).
			WithComment("comment")

		err := client.Functions.CreateForSQL(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrations("").
			HasSecrets("").
			HasIsExternalFunction(false).
			HasLanguage("SQL").
			HasIsMemoizable(true).
			HasIsDataMetric(false),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(fmt.Sprintf(`%s NOT NULL`, dataType.ToSql())).
			HasReturnDataType(dataType).
			HasReturnNotNull(true).
			HasLanguage("SQL").
			HasBody(definition).
			HasNullHandlingNil().
			// TODO [SNOW-1348103]: volatility is not returned and is present in create syntax
			// HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasVolatilityNil().
			HasExternalAccessIntegrationsNil().
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder().
			HasSecretsNil().
			HasImportsNil().
			HasExactlyImportsNormalizedInAnyOrder().
			HasHandlerNil().
			HasRuntimeVersionNil().
			HasPackagesNil().
			HasExactlyPackagesInAnyOrder().
			HasTargetPathNil().
			HasNormalizedTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create function for SQL - no arguments", func(t *testing.T) {
		dataType := testdatatypes.DataTypeFloat
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments()

		definition := testClientHelper().Function.SampleSqlDefinition(t)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		request := sdk.NewCreateForSQLFunctionRequestDefinitionWrapped(id.SchemaObjectId(), *returns, definition)

		err := client.Functions.CreateForSQL(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(0).
			HasMaxNumArguments(0).
			HasArgumentsOld([]sdk.DataType{}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s() RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultFunctionComment).
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrations("").
			HasSecrets("").
			HasIsExternalFunction(false).
			HasLanguage("SQL").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature("()").
			HasReturns(dataType.ToSql()).
			HasReturnDataType(dataType).
			HasReturnNotNull(false).
			HasLanguage("SQL").
			HasBody(definition).
			HasNullHandlingNil().
			HasVolatilityNil().
			HasExternalAccessIntegrationsNil().
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder().
			HasSecretsNil().
			HasImportsNil().
			HasExactlyImportsNormalizedInAnyOrder().
			HasHandlerNil().
			HasRuntimeVersionNil().
			HasPackagesNil().
			HasExactlyPackagesInAnyOrder().
			HasTargetPathNil().
			HasNormalizedTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("show parameters", func(t *testing.T) {
		f, fCleanup := testClientHelper().Function.CreateSql(t)
		t.Cleanup(fCleanup)
		id := f.ID()

		param, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterLogLevel, sdk.Object{ObjectType: sdk.ObjectTypeFunction, Name: id})
		require.NoError(t, err)
		assert.Equal(t, string(sdk.LogLevelOff), param.Value)

		parameters, err := client.Parameters.ShowParameters(ctx, &sdk.ShowParametersOptions{
			In: &sdk.ParametersIn{
				Function: id,
			},
		})
		require.NoError(t, err)

		assertThatObject(t, objectparametersassert.FunctionParametersPrefetched(t, id, parameters).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)

		// check that ShowParameters on function level works too
		parameters, err = client.Functions.ShowParameters(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectparametersassert.FunctionParametersPrefetched(t, id, parameters).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("alter function: rename", func(t *testing.T) {
		f, fCleanup := testClientHelper().Function.CreateSql(t)
		t.Cleanup(fCleanup)
		id := f.ID()

		nid := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(id.ArgumentDataTypes()...)
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithRenameTo(nid.SchemaObjectId()))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, nid))

		_, err = client.Functions.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		e, err := client.Functions.ShowByID(ctx, nid)
		require.NoError(t, err)
		require.Equal(t, nid.Name(), e.Name)
	})

	t.Run("alter function: set and unset all for Java", func(t *testing.T) {
		f, fCleanup := testClientHelper().Function.CreateJava(t)
		t.Cleanup(fCleanup)
		id := f.ID()

		assertThatObject(t, objectassert.Function(t, id).
			HasName(id.Name()).
			HasDescription(sdk.DefaultFunctionComment),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, id).
			HasExternalAccessIntegrationsNil().
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder().
			HasSecretsNil(),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)

		request := sdk.NewAlterFunctionRequest(id).WithSet(*sdk.NewFunctionSetRequest().
			WithEnableConsoleOutput(true).
			WithExternalAccessIntegrations([]sdk.AccountObjectIdentifier{externalAccessIntegration}).
			WithSecretsList(*sdk.NewSecretsListRequest([]sdk.SecretReference{{VariableName: "abc", Name: secretId}})).
			WithLogLevel(sdk.LogLevelWarn).
			WithMetricLevel(sdk.MetricLevelAll).
			WithTraceLevel(sdk.TraceLevelAlways).
			WithComment("new comment"),
		)

		err := client.Functions.Alter(ctx, request)
		require.NoError(t, err)

		assertThatObject(t, objectassert.Function(t, id).
			HasName(id.Name()).
			HasDescription("new comment"),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, id).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			ContainsExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}),
		)

		assertParametersSet(t, objectparametersassert.FunctionParameters(t, id))

		unsetRequest := sdk.NewAlterFunctionRequest(id).WithUnset(*sdk.NewFunctionUnsetRequest().
			WithEnableConsoleOutput(true).
			WithExternalAccessIntegrations(true).
			WithEnableConsoleOutput(true).
			WithLogLevel(true).
			WithMetricLevel(true).
			WithTraceLevel(true).
			WithComment(true),
		)

		err = client.Functions.Alter(ctx, unsetRequest)
		require.NoError(t, err)

		assertThatObject(t, objectassert.Function(t, id).
			HasName(id.Name()).
			HasDescription(sdk.DefaultFunctionComment).
			HasExactlyExternalAccessIntegrations().
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, id).
			HasExternalAccessIntegrationsNil().
			HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder().
			// TODO [SNOW-1850370]: apparently UNSET external access integrations cleans out secrets in the describe but leaves it in SHOW
			HasSecretsNil(),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)

		unsetSecretsRequest := sdk.NewAlterFunctionRequest(id).WithSet(*sdk.NewFunctionSetRequest().
			WithSecretsList(*sdk.NewSecretsListRequest([]sdk.SecretReference{})),
		)

		err = client.Functions.Alter(ctx, unsetSecretsRequest)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionDetails(t, id).
			HasSecretsNil(),
		)
	})

	t.Run("alter function: set and unset all for SQL", func(t *testing.T) {
		f, fCleanup := testClientHelper().Function.CreateSql(t)
		t.Cleanup(fCleanup)
		id := f.ID()

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)

		request := sdk.NewAlterFunctionRequest(id).WithSet(*sdk.NewFunctionSetRequest().
			WithEnableConsoleOutput(true).
			WithLogLevel(sdk.LogLevelWarn).
			WithMetricLevel(sdk.MetricLevelAll).
			WithTraceLevel(sdk.TraceLevelAlways).
			WithComment("new comment"),
		)

		err := client.Functions.Alter(ctx, request)
		require.NoError(t, err)

		assertThatObject(t, objectassert.Function(t, id).
			HasName(id.Name()).
			HasDescription("new comment"),
		)

		assertParametersSet(t, objectparametersassert.FunctionParameters(t, id))

		unsetRequest := sdk.NewAlterFunctionRequest(id).WithUnset(*sdk.NewFunctionUnsetRequest().
			WithEnableConsoleOutput(true).
			WithLogLevel(true).
			WithMetricLevel(true).
			WithTraceLevel(true).
			WithComment(true),
		)

		err = client.Functions.Alter(ctx, unsetRequest)
		require.NoError(t, err)

		assertThatObject(t, objectassert.Function(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDescription(sdk.DefaultFunctionComment),
		)

		assertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("alter function: set and unset secure", func(t *testing.T) {
		f, fCleanup := testClientHelper().Function.CreateSql(t)
		t.Cleanup(fCleanup)
		id := f.ID()

		assertThatObject(t, objectassert.FunctionFromObject(t, f).
			HasIsSecure(false),
		)

		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSetSecure(true))
		require.NoError(t, err)

		assertThatObject(t, objectassert.Function(t, id).
			HasIsSecure(true),
		)

		err = client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithUnsetSecure(true))
		require.NoError(t, err)

		assertThatObject(t, objectassert.Function(t, id).
			HasIsSecure(false),
		)
	})

	t.Run("show function: without like", func(t *testing.T) {
		f1, fCleanup := testClientHelper().Function.CreateSql(t)
		t.Cleanup(fCleanup)

		f2, fCleanup2 := testClientHelper().Function.CreateSql(t)
		t.Cleanup(fCleanup2)

		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest())
		require.NoError(t, err)

		require.Contains(t, functions, *f1)
		require.Contains(t, functions, *f2)
	})

	t.Run("show function: with like", func(t *testing.T) {
		f1, fCleanup := testClientHelper().Function.CreateSql(t)
		t.Cleanup(fCleanup)

		f2, fCleanup2 := testClientHelper().Function.CreateSql(t)
		t.Cleanup(fCleanup2)

		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest().WithLike(sdk.Like{Pattern: &f1.Name}))
		require.NoError(t, err)

		require.Equal(t, 1, len(functions))
		require.Contains(t, functions, *f1)
		require.NotContains(t, functions, *f2)
	})

	t.Run("show function: no matches", func(t *testing.T) {
		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest().
			WithIn(sdk.ExtendedIn{In: sdk.In{Schema: testClientHelper().Ids.SchemaId()}}).
			WithLike(sdk.Like{Pattern: sdk.String(NonExistingSchemaObjectIdentifier.Name())}))
		require.NoError(t, err)
		require.Equal(t, 0, len(functions))
	})

	t.Run("describe function: for Java - minimal", func(t *testing.T) {
		f, fCleanup := testClientHelper().Function.CreateJava(t)
		t.Cleanup(fCleanup)
		id := f.ID()

		details, err := client.Functions.Describe(ctx, id)
		require.NoError(t, err)
		assert.Len(t, details, 11)

		pairs := make(map[string]*string)
		for _, detail := range details {
			pairs[detail.Property] = detail.Value
		}
		assert.Equal(t, "(x VARCHAR)", *pairs["signature"])
		assert.Equal(t, "VARCHAR(100)", *pairs["returns"])
		assert.Equal(t, "JAVA", *pairs["language"])
		assert.NotEmpty(t, *pairs["body"])
		assert.Equal(t, string(sdk.NullInputBehaviorCalledOnNullInput), *pairs["null handling"])
		assert.Equal(t, string(sdk.VolatileTableKind), *pairs["volatility"])
		assert.Nil(t, pairs["external_access_integration"])
		assert.Nil(t, pairs["secrets"])
		assert.Equal(t, "[]", *pairs["imports"])
		assert.Equal(t, "TestFunc.echoVarchar", *pairs["handler"])
		assert.Nil(t, pairs["runtime_version"])
	})

	t.Run("describe function: for SQL - with arguments", func(t *testing.T) {
		f, fCleanup := testClientHelper().Function.CreateSql(t)
		t.Cleanup(fCleanup)
		id := f.ID()

		details, err := client.Functions.Describe(ctx, id)
		require.NoError(t, err)
		assert.Len(t, details, 4)

		pairs := make(map[string]string)
		for _, detail := range details {
			pairs[detail.Property] = *detail.Value
		}
		assert.Equal(t, "(x FLOAT)", pairs["signature"])
		assert.Equal(t, "FLOAT", pairs["returns"])
		assert.Equal(t, "SQL", pairs["language"])
		assert.Equal(t, "3.141592654::FLOAT", pairs["body"])
	})

	t.Run("describe function: for SQL - no arguments", func(t *testing.T) {
		f, fCleanup := testClientHelper().Function.CreateSqlNoArgs(t)
		t.Cleanup(fCleanup)
		id := f.ID()

		details, err := client.Functions.Describe(ctx, id)
		require.NoError(t, err)
		assert.Len(t, details, 4)

		pairs := make(map[string]string)
		for _, detail := range details {
			pairs[detail.Property] = *detail.Value
		}
		assert.Equal(t, "()", pairs["signature"])
		assert.Equal(t, "FLOAT", pairs["returns"])
		assert.Equal(t, "SQL", pairs["language"])
		assert.Equal(t, "3.141592654::FLOAT", pairs["body"])
	})

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		dataType := testdatatypes.DataTypeFloat
		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierWithArgumentsInSchema(id1.Name(), schema.ID(), sdk.LegacyDataTypeFrom(dataType))

		_, fCleanup1 := testClientHelper().Function.CreateSqlWithIdentifierAndArgument(t, id1.SchemaObjectId(), dataType)
		t.Cleanup(fCleanup1)
		_, fCleanup2 := testClientHelper().Function.CreateSqlWithIdentifierAndArgument(t, id2.SchemaObjectId(), dataType)
		t.Cleanup(fCleanup2)

		e1, err := client.Functions.ShowByID(ctx, id1)
		require.NoError(t, err)

		e1Id := e1.ID()
		require.NoError(t, err)
		require.Equal(t, id1, e1Id)

		e2, err := client.Functions.ShowByID(ctx, id2)
		require.NoError(t, err)

		e2Id := e2.ID()
		require.NoError(t, err)
		require.Equal(t, id2, e2Id)
	})

	t.Run("show function by id - same name, different arguments", func(t *testing.T) {
		dataType := testdatatypes.DataTypeFloat
		name := testClientHelper().Ids.Alpha()

		id1 := testClientHelper().Ids.NewSchemaObjectIdentifierWithArgumentsInSchema(name, testClientHelper().Ids.SchemaId(), sdk.LegacyDataTypeFrom(dataType))
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierWithArgumentsInSchema(name, testClientHelper().Ids.SchemaId(), sdk.DataTypeInt, sdk.DataTypeVARCHAR)

		e := testClientHelper().Function.CreateWithIdentifier(t, id1)
		testClientHelper().Function.CreateWithIdentifier(t, id2)

		es, err := client.Functions.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, *e, *es)
	})

	// TODO [SNOW-1348103]: remove with old function removal for V1
	t.Run("function returns non detailed data types of arguments - old data types", func(t *testing.T) {
		// This test proves that every detailed data types (e.g. VARCHAR(20) and NUMBER(10, 0)) are generalized
		// on Snowflake side (to e.g. VARCHAR and NUMBER) and that sdk.ToDataType mapping function maps detailed types
		// correctly to their generalized counterparts (same as in Snowflake).

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		args := []sdk.FunctionArgumentRequest{
			*sdk.NewFunctionArgumentRequest("A", nil).WithArgDataTypeOld("NUMBER(2, 0)"),
			*sdk.NewFunctionArgumentRequest("B", nil).WithArgDataTypeOld("DECIMAL"),
			*sdk.NewFunctionArgumentRequest("C", nil).WithArgDataTypeOld("INTEGER"),
			*sdk.NewFunctionArgumentRequest("D", nil).WithArgDataTypeOld(sdk.DataTypeFloat),
			*sdk.NewFunctionArgumentRequest("E", nil).WithArgDataTypeOld("DOUBLE"),
			*sdk.NewFunctionArgumentRequest("F", nil).WithArgDataTypeOld("VARCHAR(20)"),
			*sdk.NewFunctionArgumentRequest("G", nil).WithArgDataTypeOld("CHAR"),
			*sdk.NewFunctionArgumentRequest("H", nil).WithArgDataTypeOld(sdk.DataTypeString),
			*sdk.NewFunctionArgumentRequest("I", nil).WithArgDataTypeOld("TEXT"),
			*sdk.NewFunctionArgumentRequest("J", nil).WithArgDataTypeOld(sdk.DataTypeBinary),
			*sdk.NewFunctionArgumentRequest("K", nil).WithArgDataTypeOld("VARBINARY"),
			*sdk.NewFunctionArgumentRequest("L", nil).WithArgDataTypeOld(sdk.DataTypeBoolean),
			*sdk.NewFunctionArgumentRequest("M", nil).WithArgDataTypeOld(sdk.DataTypeDate),
			*sdk.NewFunctionArgumentRequest("N", nil).WithArgDataTypeOld("DATETIME"),
			*sdk.NewFunctionArgumentRequest("O", nil).WithArgDataTypeOld(sdk.DataTypeTime),
			*sdk.NewFunctionArgumentRequest("R", nil).WithArgDataTypeOld(sdk.DataTypeTimestampLTZ),
			*sdk.NewFunctionArgumentRequest("S", nil).WithArgDataTypeOld(sdk.DataTypeTimestampNTZ),
			*sdk.NewFunctionArgumentRequest("T", nil).WithArgDataTypeOld(sdk.DataTypeTimestampTZ),
			*sdk.NewFunctionArgumentRequest("U", nil).WithArgDataTypeOld(sdk.DataTypeVariant),
			*sdk.NewFunctionArgumentRequest("V", nil).WithArgDataTypeOld(sdk.DataTypeObject),
			*sdk.NewFunctionArgumentRequest("W", nil).WithArgDataTypeOld(sdk.DataTypeArray),
			*sdk.NewFunctionArgumentRequest("X", nil).WithArgDataTypeOld(sdk.DataTypeGeography),
			*sdk.NewFunctionArgumentRequest("Y", nil).WithArgDataTypeOld(sdk.DataTypeGeometry),
			*sdk.NewFunctionArgumentRequest("Z", nil).WithArgDataTypeOld("VECTOR(INT, 16)"),
		}
		err := client.Functions.CreateForPython(ctx, sdk.NewCreateForPythonFunctionRequest(
			id,
			*sdk.NewFunctionReturnsRequest().WithResultDataType(*sdk.NewFunctionReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeVariant)),
			testvars.PythonRuntime,
			"add",
		).
			WithArguments(args).
			WithFunctionDefinitionWrapped("def add(A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, R, S, T, U, V, W, X, Y, Z): A + A"),
		)
		require.NoError(t, err)

		dataTypes := make([]sdk.DataType, len(args))
		for i, arg := range args {
			dataType, err := datatypes.ParseDataType(string(arg.ArgDataTypeOld))
			require.NoError(t, err)
			dataTypes[i] = sdk.LegacyDataTypeFrom(dataType)
		}
		idWithArguments := sdk.NewSchemaObjectIdentifierWithArguments(id.DatabaseName(), id.SchemaName(), id.Name(), dataTypes...)

		function, err := client.Functions.ShowByID(ctx, idWithArguments)
		require.NoError(t, err)
		require.Equal(t, dataTypes, function.ArgumentsOld)
	})

	// This test shows behavior of detailed types (e.g. VARCHAR(20) and NUMBER(10, 0)) on Snowflake side for functions.
	// For SHOW, data type is generalized both for argument and return type (to e.g. VARCHAR and NUMBER).
	// FOR DESCRIBE, data type is generalized for argument and works weirdly for the return type: type is generalized to the canonical one, but we also get the attributes.
	for _, tc := range []string{
		"NUMBER(36, 5)",
		"NUMBER(36)",
		"NUMBER",
		"DECIMAL",
		"INTEGER",
		"FLOAT",
		"DOUBLE",
		"VARCHAR",
		"VARCHAR(20)",
		"CHAR",
		"CHAR(10)",
		"TEXT",
		"BINARY",
		"BINARY(1000)",
		"VARBINARY",
		"BOOLEAN",
		"DATE",
		"DATETIME",
		"TIME",
		"TIMESTAMP_LTZ",
		"TIMESTAMP_NTZ",
		"TIMESTAMP_TZ",
		"VARIANT",
		"OBJECT",
		"ARRAY",
		"GEOGRAPHY",
		"GEOMETRY",
		"VECTOR(INT, 16)",
		"VECTOR(FLOAT, 8)",
	} {
		tc := tc
		t.Run(fmt.Sprintf("function returns non detailed data types of arguments for %s", tc), func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			argName := "A"
			funcName := "identity"
			dataType, err := datatypes.ParseDataType(tc)
			require.NoError(t, err)
			args := []sdk.FunctionArgumentRequest{
				*sdk.NewFunctionArgumentRequest(argName, dataType),
			}

			err = client.Functions.CreateForPython(ctx, sdk.NewCreateForPythonFunctionRequest(
				id,
				*sdk.NewFunctionReturnsRequest().WithResultDataType(*sdk.NewFunctionReturnsResultDataTypeRequest(dataType)),
				testvars.PythonRuntime,
				funcName,
			).
				WithArguments(args).
				WithFunctionDefinitionWrapped(testClientHelper().Function.PythonIdentityDefinition(t, funcName, argName)),
			)
			require.NoError(t, err)

			oldDataType := sdk.LegacyDataTypeFrom(dataType)
			idWithArguments := sdk.NewSchemaObjectIdentifierWithArguments(id.DatabaseName(), id.SchemaName(), id.Name(), oldDataType)

			function, err := client.Functions.ShowByID(ctx, idWithArguments)
			require.NoError(t, err)
			assert.Equal(t, []sdk.DataType{oldDataType}, function.ArgumentsOld)
			assert.Equal(t, fmt.Sprintf("%[1]s(%[2]s) RETURN %[2]s", id.Name(), oldDataType), function.ArgumentsRaw)

			details, err := client.Functions.Describe(ctx, idWithArguments)
			require.NoError(t, err)
			pairs := make(map[string]string)
			for _, detail := range details {
				pairs[detail.Property] = *detail.Value
			}
			assert.Equal(t, fmt.Sprintf("(%s %s)", argName, oldDataType), pairs["signature"])
			assert.Equal(t, dataType.Canonical(), pairs["returns"])
		})
	}

	t.Run("create function for SQL - return table data type", func(t *testing.T) {
		argName := "x"

		returnDataType, err := datatypes.ParseDataType(fmt.Sprintf("TABLE(PRICE %s, THIRD %s)", datatypes.FloatLegacyDataType, datatypes.VarcharLegacyDataType))
		require.NoError(t, err)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(datatypes.VarcharLegacyDataType)

		definition := `
SELECT 2.2::float, 'abc');` // the ending parenthesis has to be there (otherwise SQL compilation error is thrown)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(returnDataType)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewFunctionArgumentRequest(argName, nil).WithArgDataTypeOld(datatypes.VarcharLegacyDataType)
		request := sdk.NewCreateForSQLFunctionRequestDefinitionWrapped(id.SchemaObjectId(), *returns, definition).
			WithArguments([]sdk.FunctionArgumentRequest{*argument})

		err = client.Functions.CreateForSQL(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasArgumentsRawContains(strings.ReplaceAll(returnDataType.ToLegacyDataTypeSql(), "TABLE(", "TABLE (")),
		)

		assertThatObject(t, objectassert.FunctionDetails(t, id).
			HasReturnDataType(returnDataType),
		)
	})
}

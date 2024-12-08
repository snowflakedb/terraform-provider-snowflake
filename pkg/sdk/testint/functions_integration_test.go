package testint

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	assertions "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO [next PR]: schemaName and catalog name are quoted (because we use lowercase)
// TODO [next PR]: HasArgumentsRawFrom(functionId, arguments, return)
// TODO [next PR]: extract show assertions with commons fields
// TODO [this PR]: python aggregate func
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

	tmpJavaFunction := testClientHelper().CreateSampleJavaFunctionAndJar(t)
	tmpPythonFunction := testClientHelper().CreateSamplePythonFunctionAndModule(t)

	//assertParametersSet := func(t *testing.T, functionParametersAssert *objectparametersassert.FunctionParametersAssert) {
	//	assertions.AssertThatObject(t, functionParametersAssert.
	//		HasEnableConsoleOutput(true).
	//		HasLogLevel(sdk.LogLevelWarn).
	//		HasMetricLevel(sdk.MetricLevelAll).
	//		HasTraceLevel(sdk.TraceLevelAlways),
	//	)
	//}

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

		assertions.AssertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(fmt.Sprintf(`"%s"`, id.SchemaName())).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultFunctionComment).
			HasCatalogName(fmt.Sprintf(`"%s"`, id.DatabaseName())).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasIsExternalFunction(false).
			HasLanguage("JAVA").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertions.AssertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(dataType.ToSql()).
			HasLanguage("JAVA").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImports(`[]`).
			HasHandler(handler).
			HasRuntimeVersionNil().
			HasPackages(`[]`).
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertions.AssertThatObject(t, objectparametersassert.FunctionParameters(t, id).
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
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnNullInput)).
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

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(fmt.Sprintf(`"%s"`, id.SchemaName())).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(fmt.Sprintf(`"%s"`, id.DatabaseName())).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasIsExternalFunction(false).
			HasLanguage("JAVA").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertions.AssertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(fmt.Sprintf(`%s NOT NULL`, dataType.ToSql())).
			HasLanguage("JAVA").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorReturnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExternalAccessIntegrations(fmt.Sprintf(`[%s]`, externalAccessIntegration.FullyQualifiedName())).
			// TODO [this PR]: parse to identifier list
			// TODO [this PR]: check multiple secrets (to know how to parse)
			HasSecrets(fmt.Sprintf(`{"abc":"\"%s\".\"%s\".%s"}`, secretId.DatabaseName(), secretId.SchemaName(), secretId.Name())).
			HasImports(fmt.Sprintf(`[%s]`, tmpJavaFunction.JarLocation())).
			HasHandler(handler).
			HasRuntimeVersion("11").
			HasPackages(`[com.snowflake:snowpark:1.14.0,com.snowflake:telemetry:0.1.0]`).
			HasTargetPath(targetPath).
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertions.AssertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)

		// TODO [this PR]: will check after alter
		// TODO [this PR]: add a test documenting that we can't set parameters in create (and revert adding these parametrs directly in object...)
		//assertParametersSet(t, objectparametersassert.FunctionParameters(t, id))
		//
		//// check that ShowParameters works too
		//parameters, err := client.Functions.ShowParameters(ctx, id)
		//require.NoError(t, err)
		//assertParametersSet(t, objectparametersassert.FunctionParametersPrefetched(t, id, parameters))
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

		assertions.AssertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(fmt.Sprintf(`"%s"`, id.SchemaName())).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultFunctionComment).
			HasCatalogName(fmt.Sprintf(`"%s"`, id.DatabaseName())).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasIsExternalFunction(false).
			HasLanguage("JAVA").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertions.AssertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(dataType.ToSql()).
			HasLanguage("JAVA").
			HasBodyNil().
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImports(fmt.Sprintf(`[%s]`, importPath)).
			HasHandler(handler).
			HasRuntimeVersionNil().
			HasPackages(`[]`).
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertions.AssertThatObject(t, objectparametersassert.FunctionParameters(t, id).
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
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnNullInput)).
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

		assertions.AssertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(fmt.Sprintf(`"%s"`, id.SchemaName())).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(fmt.Sprintf(`"%s"`, id.DatabaseName())).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasIsExternalFunction(false).
			HasLanguage("JAVA").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertions.AssertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(fmt.Sprintf(`%s NOT NULL`, dataType.ToSql())).
			HasLanguage("JAVA").
			HasBodyNil().
			HasNullHandling(string(sdk.NullInputBehaviorReturnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExternalAccessIntegrations(fmt.Sprintf(`[%s]`, externalAccessIntegration.FullyQualifiedName())).
			HasSecrets(fmt.Sprintf(`{"abc":"\"%s\".\"%s\".%s"}`, secretId.DatabaseName(), secretId.SchemaName(), secretId.Name())).
			HasImports(fmt.Sprintf(`[%s]`, tmpJavaFunction.JarLocation())).
			HasHandler(handler).
			HasRuntimeVersion("11").
			HasPackages(`[com.snowflake:snowpark:1.14.0,com.snowflake:telemetry:0.1.0]`).
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertions.AssertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
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

		assertions.AssertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(fmt.Sprintf(`"%s"`, id.SchemaName())).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultFunctionComment).
			HasCatalogName(fmt.Sprintf(`"%s"`, id.DatabaseName())).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasIsExternalFunction(false).
			HasLanguage("JAVASCRIPT").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertions.AssertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(dataType.ToSql()).
			HasLanguage("JAVASCRIPT").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImportsNil().
			HasHandlerNil().
			HasRuntimeVersionNil().
			HasPackagesNil().
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertions.AssertThatObject(t, objectparametersassert.FunctionParameters(t, id).
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
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithComment("comment")

		err := client.Functions.CreateForJavascript(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(fmt.Sprintf(`"%s"`, id.SchemaName())).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(fmt.Sprintf(`"%s"`, id.DatabaseName())).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasIsExternalFunction(false).
			HasLanguage("JAVASCRIPT").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertions.AssertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(fmt.Sprintf(`%s NOT NULL`, dataType.ToSql())).
			HasLanguage("JAVASCRIPT").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorReturnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImportsNil().
			HasHandlerNil().
			HasRuntimeVersionNil().
			HasPackagesNil().
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertions.AssertThatObject(t, objectparametersassert.FunctionParameters(t, id).
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
		request := sdk.NewCreateForPythonFunctionRequest(id.SchemaObjectId(), *returns, "3.8", funcName).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithFunctionDefinitionWrapped(definition)

		err := client.Functions.CreateForPython(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(fmt.Sprintf(`"%s"`, id.SchemaName())).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultFunctionComment).
			HasCatalogName(fmt.Sprintf(`"%s"`, id.DatabaseName())).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasIsExternalFunction(false).
			HasLanguage("PYTHON").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertions.AssertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(strings.ReplaceAll(dataType.ToSql(), " ", "")). //TODO [this PR]: do we care about this whitespace?
			HasLanguage("PYTHON").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImports(`[]`).
			HasHandler(funcName).
			HasRuntimeVersion("3.8").
			HasPackages(`[]`).
			HasTargetPathNil().
			HasInstalledPackagesNotEmpty().
			HasIsAggregate(false),
		)

		assertions.AssertThatObject(t, objectparametersassert.FunctionParameters(t, id).
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
		request := sdk.NewCreateForPythonFunctionRequest(id.SchemaObjectId(), *returns, "3.8", funcName).
			WithOrReplace(true).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithReturnNullValues(sdk.ReturnNullValuesNotNull).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithComment("comment").
			WithImports([]sdk.FunctionImportRequest{*sdk.NewFunctionImportRequest().WithImport(tmpPythonFunction.PythonModuleLocation())}).
			WithPackages([]sdk.FunctionPackageRequest{
				*sdk.NewFunctionPackageRequest().WithPackage("absl-py==0.10.0"),
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

		assertions.AssertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(fmt.Sprintf(`"%s"`, id.SchemaName())).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(fmt.Sprintf(`"%s"`, id.DatabaseName())).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasIsExternalFunction(false).
			HasLanguage("PYTHON").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertions.AssertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(strings.ReplaceAll(dataType.ToSql(), " ", "")+" NOT NULL"). //TODO [this PR]: do we care about this whitespace?
			HasLanguage("PYTHON").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorReturnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExternalAccessIntegrations(fmt.Sprintf(`[%s]`, externalAccessIntegration.FullyQualifiedName())).
			HasSecrets(fmt.Sprintf(`{"abc":"\"%s\".\"%s\".%s"}`, secretId.DatabaseName(), secretId.SchemaName(), secretId.Name())).
			HasImports(fmt.Sprintf(`[%s]`, tmpPythonFunction.PythonModuleLocation())).
			HasHandler(funcName).
			HasRuntimeVersion("3.8").
			HasPackages(`['absl-py==0.10.0','about-time==4.2.1']`).
			HasTargetPathNil().
			HasInstalledPackagesNotEmpty().
			HasIsAggregate(false),
		)

		assertions.AssertThatObject(t, objectparametersassert.FunctionParameters(t, id).
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
		request := sdk.NewCreateForPythonFunctionRequest(id.SchemaObjectId(), *returns, "3.8", tmpPythonFunction.PythonHandler()).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithImports([]sdk.FunctionImportRequest{*sdk.NewFunctionImportRequest().WithImport(tmpPythonFunction.PythonModuleLocation())})

		err := client.Functions.CreateForPython(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(fmt.Sprintf(`"%s"`, id.SchemaName())).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultFunctionComment).
			HasCatalogName(fmt.Sprintf(`"%s"`, id.DatabaseName())).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasIsExternalFunction(false).
			HasLanguage("PYTHON").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertions.AssertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(strings.ReplaceAll(dataType.ToSql(), " ", "")). //TODO [this PR]: do we care about this whitespace?
			HasLanguage("PYTHON").
			HasBodyNil().
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImports(fmt.Sprintf(`[%s]`, tmpPythonFunction.PythonModuleLocation())).
			HasHandler(tmpPythonFunction.PythonHandler()).
			HasRuntimeVersion("3.8").
			HasPackages(`[]`).
			HasTargetPathNil().
			HasInstalledPackagesNotEmpty().
			HasIsAggregate(false),
		)

		assertions.AssertThatObject(t, objectparametersassert.FunctionParameters(t, id).
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
		request := sdk.NewCreateForPythonFunctionRequest(id.SchemaObjectId(), *returns, "3.8", tmpPythonFunction.PythonHandler()).
			WithOrReplace(true).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithReturnNullValues(sdk.ReturnNullValuesNotNull).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithComment("comment").
			WithPackages([]sdk.FunctionPackageRequest{
				*sdk.NewFunctionPackageRequest().WithPackage("absl-py==0.10.0"),
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

		assertions.AssertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(fmt.Sprintf(`"%s"`, id.SchemaName())).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(fmt.Sprintf(`"%s"`, id.DatabaseName())).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasIsExternalFunction(false).
			HasLanguage("PYTHON").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertions.AssertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(strings.ReplaceAll(dataType.ToSql(), " ", "")+" NOT NULL"). //TODO [this PR]: do we care about this whitespace?
			HasLanguage("PYTHON").
			HasBodyNil().
			HasNullHandling(string(sdk.NullInputBehaviorReturnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExternalAccessIntegrations(fmt.Sprintf(`[%s]`, externalAccessIntegration.FullyQualifiedName())).
			HasSecrets(fmt.Sprintf(`{"abc":"\"%s\".\"%s\".%s"}`, secretId.DatabaseName(), secretId.SchemaName(), secretId.Name())).
			HasImports(fmt.Sprintf(`[%s]`, tmpPythonFunction.PythonModuleLocation())).
			HasHandler(tmpPythonFunction.PythonHandler()).
			HasRuntimeVersion("3.8").
			HasPackages(`['absl-py==0.10.0','about-time==4.2.1']`).
			HasTargetPathNil().
			HasInstalledPackagesNotEmpty().
			HasIsAggregate(false),
		)

		assertions.AssertThatObject(t, objectparametersassert.FunctionParameters(t, id).
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
		request := sdk.NewCreateForScalaFunctionRequest(id.SchemaObjectId(), dataType, handler).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithRuntimeVersion("2.12").
			WithFunctionDefinitionWrapped(definition)

		err := client.Functions.CreateForScala(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(fmt.Sprintf(`"%s"`, id.SchemaName())).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultFunctionComment).
			HasCatalogName(fmt.Sprintf(`"%s"`, id.DatabaseName())).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasIsExternalFunction(false).
			HasLanguage("SCALA").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertions.AssertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(dataType.ToSql()).
			HasLanguage("SCALA").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImports(`[]`).
			HasHandler(handler).
			HasRuntimeVersion("2.12").
			HasPackages(`[]`).
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertions.AssertThatObject(t, objectparametersassert.FunctionParameters(t, id).
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
		request := sdk.NewCreateForScalaFunctionRequest(id.SchemaObjectId(), dataType, handler).
			WithOrReplace(true).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithReturnNullValues(sdk.ReturnNullValuesNotNull).
			WithRuntimeVersion("2.12").
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

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(fmt.Sprintf(`"%s"`, id.SchemaName())).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(fmt.Sprintf(`"%s"`, id.DatabaseName())).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasIsExternalFunction(false).
			HasLanguage("SCALA").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertions.AssertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(fmt.Sprintf(`%s NOT NULL`, dataType.ToSql())).
			HasLanguage("SCALA").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorReturnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExternalAccessIntegrations(fmt.Sprintf(`[%s]`, externalAccessIntegration.FullyQualifiedName())).
			HasSecrets(fmt.Sprintf(`{"abc":"\"%s\".\"%s\".%s"}`, secretId.DatabaseName(), secretId.SchemaName(), secretId.Name())).
			HasImports(fmt.Sprintf(`[%s]`, tmpJavaFunction.JarLocation())).
			HasHandler(handler).
			HasRuntimeVersion("2.12").
			HasPackages(`[com.snowflake:snowpark:1.14.0,com.snowflake:telemetry:0.1.0]`).
			HasTargetPath(targetPath).
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertions.AssertThatObject(t, objectparametersassert.FunctionParameters(t, id).
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

		requestStaged := sdk.NewCreateForScalaFunctionRequest(id.SchemaObjectId(), dataType, handler).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithRuntimeVersion("2.12").
			WithImports([]sdk.FunctionImportRequest{*sdk.NewFunctionImportRequest().WithImport(importPath)})

		err := client.Functions.CreateForScala(ctx, requestStaged)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Function.DropFunctionFunc(t, id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(fmt.Sprintf(`"%s"`, id.SchemaName())).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultFunctionComment).
			HasCatalogName(fmt.Sprintf(`"%s"`, id.DatabaseName())).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasIsExternalFunction(false).
			HasLanguage("SCALA").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertions.AssertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(dataType.ToSql()).
			HasLanguage("SCALA").
			HasBodyNil().
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImports(fmt.Sprintf(`[%s]`, importPath)).
			HasHandler(handler).
			HasRuntimeVersion("2.12").
			HasPackages(`[]`).
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertions.AssertThatObject(t, objectparametersassert.FunctionParameters(t, id).
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

		requestStaged := sdk.NewCreateForScalaFunctionRequest(id.SchemaObjectId(), dataType, handler).
			WithOrReplace(true).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithReturnNullValues(sdk.ReturnNullValuesNotNull).
			WithRuntimeVersion("2.12").
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

		assertions.AssertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(fmt.Sprintf(`"%s"`, id.SchemaName())).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(fmt.Sprintf(`"%s"`, id.DatabaseName())).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasIsExternalFunction(false).
			HasLanguage("SCALA").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertions.AssertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(fmt.Sprintf(`%s NOT NULL`, dataType.ToSql())).
			HasLanguage("SCALA").
			HasBodyNil().
			HasNullHandling(string(sdk.NullInputBehaviorReturnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExternalAccessIntegrations(fmt.Sprintf(`[%s]`, externalAccessIntegration.FullyQualifiedName())).
			HasSecrets(fmt.Sprintf(`{"abc":"\"%s\".\"%s\".%s"}`, secretId.DatabaseName(), secretId.SchemaName(), secretId.Name())).
			HasImports(fmt.Sprintf(`[%s]`, tmpJavaFunction.JarLocation())).
			HasHandler(handler).
			HasRuntimeVersion("2.12").
			HasPackages(`[com.snowflake:snowpark:1.14.0,com.snowflake:telemetry:0.1.0]`).
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertions.AssertThatObject(t, objectparametersassert.FunctionParameters(t, id).
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

		assertions.AssertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(fmt.Sprintf(`"%s"`, id.SchemaName())).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultFunctionComment).
			HasCatalogName(fmt.Sprintf(`"%s"`, id.DatabaseName())).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasIsExternalFunction(false).
			HasLanguage("SQL").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertions.AssertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(dataType.ToSql()).
			HasLanguage("SQL").
			HasBody(definition).
			HasNullHandlingNil().
			HasVolatilityNil().
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImportsNil().
			HasHandlerNil().
			HasRuntimeVersionNil().
			HasPackagesNil().
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertions.AssertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
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

		assertions.AssertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(fmt.Sprintf(`"%s"`, id.SchemaName())).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(fmt.Sprintf(`"%s"`, id.DatabaseName())).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasIsExternalFunction(false).
			HasLanguage("SQL").
			HasIsMemoizable(true).
			HasIsDataMetric(false),
		)

		assertions.AssertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(fmt.Sprintf(`%s NOT NULL`, dataType.ToSql())).
			HasLanguage("SQL").
			HasBody(definition).
			HasNullHandlingNil().
			// TODO [next PR]: volatility is not returned and is present in create syntax
			//HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasVolatilityNil().
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImportsNil().
			HasHandlerNil().
			HasRuntimeVersionNil().
			HasPackagesNil().
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertions.AssertThatObject(t, objectparametersassert.FunctionParameters(t, id).
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

		assertions.AssertThatObject(t, objectassert.FunctionFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(fmt.Sprintf(`"%s"`, id.SchemaName())).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(0).
			HasMaxNumArguments(0).
			HasArgumentsOld([]sdk.DataType{}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s() RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultFunctionComment).
			HasCatalogName(fmt.Sprintf(`"%s"`, id.DatabaseName())).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasIsExternalFunction(false).
			HasLanguage("SQL").
			HasIsMemoizable(false).
			HasIsDataMetric(false),
		)

		assertions.AssertThatObject(t, objectassert.FunctionDetails(t, function.ID()).
			HasSignature("()").
			HasReturns(dataType.ToSql()).
			HasLanguage("SQL").
			HasBody(definition).
			HasNullHandlingNil().
			HasVolatilityNil().
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImportsNil().
			HasHandlerNil().
			HasRuntimeVersionNil().
			HasPackagesNil().
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasIsAggregateNil(),
		)

		assertions.AssertThatObject(t, objectparametersassert.FunctionParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})
}

func TestInt_OtherFunctions(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	assertFunction := func(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments, secure bool, withArguments bool) {
		t.Helper()

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.NotEmpty(t, function.CreatedOn)
		assert.Equal(t, id.Name(), function.Name)
		assert.Equal(t, false, function.IsBuiltin)
		assert.Equal(t, false, function.IsAggregate)
		assert.Equal(t, false, function.IsAnsi)
		if withArguments {
			assert.Equal(t, 1, function.MinNumArguments)
			assert.Equal(t, 1, function.MaxNumArguments)
		} else {
			assert.Equal(t, 0, function.MinNumArguments)
			assert.Equal(t, 0, function.MaxNumArguments)
		}
		assert.NotEmpty(t, function.ArgumentsRaw)
		assert.NotEmpty(t, function.ArgumentsOld)
		assert.NotEmpty(t, function.Description)
		assert.NotEmpty(t, function.CatalogName)
		assert.Equal(t, false, function.IsTableFunction)
		assert.Equal(t, false, function.ValidForClustering)
		assert.Equal(t, secure, function.IsSecure)
		assert.Equal(t, false, function.IsExternalFunction)
		assert.Equal(t, "SQL", function.Language)
		assert.Equal(t, false, function.IsMemoizable)
	}

	cleanupFunctionHandle := func(id sdk.SchemaObjectIdentifierWithArguments) func() {
		return func() {
			err := client.Functions.Drop(ctx, sdk.NewDropFunctionRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createFunctionForSQLHandle := func(t *testing.T, cleanup bool, withArguments bool) *sdk.Function {
		t.Helper()
		var id sdk.SchemaObjectIdentifierWithArguments
		if withArguments {
			id = testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeFloat)
		} else {
			id = testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments()
		}

		definition := testClientHelper().Function.SampleSqlDefinition(t)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeFloat)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		request := sdk.NewCreateForSQLFunctionRequest(id.SchemaObjectId(), *returns, definition).
			WithOrReplace(true)
		if withArguments {
			argument := sdk.NewFunctionArgumentRequest("x", nil).WithArgDataTypeOld(sdk.DataTypeFloat)
			request = request.WithArguments([]sdk.FunctionArgumentRequest{*argument})
		}
		err := client.Functions.CreateForSQL(ctx, request)
		require.NoError(t, err)
		if cleanup {
			t.Cleanup(cleanupFunctionHandle(id))
		}
		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		return function
	}

	t.Run("alter function: rename", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, false, true)

		id := f.ID()
		nid := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeFloat)
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithRenameTo(nid.SchemaObjectId()))
		if err != nil {
			t.Cleanup(cleanupFunctionHandle(id))
		} else {
			t.Cleanup(cleanupFunctionHandle(nid))
		}
		require.NoError(t, err)

		_, err = client.Functions.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		e, err := client.Functions.ShowByID(ctx, nid)
		require.NoError(t, err)
		require.Equal(t, nid.Name(), e.Name)
	})

	t.Run("alter function: set log level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSet(*sdk.NewFunctionSetRequest().WithLogLevel(sdk.LogLevelDebug)))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: unset log level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithUnset(*sdk.NewFunctionUnsetRequest().WithLogLevel(true)))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: set trace level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSet(*sdk.NewFunctionSetRequest().WithTraceLevel(sdk.TraceLevelAlways)))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: unset trace level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithUnset(*sdk.NewFunctionUnsetRequest().WithTraceLevel(true)))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: set comment", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSet(*sdk.NewFunctionSetRequest().WithComment("test comment")))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: unset comment", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithUnset(*sdk.NewFunctionUnsetRequest().WithComment(true)))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: set secure", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSetSecure(true))
		require.NoError(t, err)
		assertFunction(t, id, true, true)
	})

	t.Run("alter function: set secure with no arguments", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)
		id := f.ID()
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSetSecure(true))
		require.NoError(t, err)
		assertFunction(t, id, true, true)
	})

	t.Run("alter function: unset secure", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithUnsetSecure(true))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("show function for SQL: without like", func(t *testing.T) {
		f1 := createFunctionForSQLHandle(t, true, true)
		f2 := createFunctionForSQLHandle(t, true, true)

		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest())
		require.NoError(t, err)

		require.Contains(t, functions, *f1)
		require.Contains(t, functions, *f2)
	})

	t.Run("show function for SQL: with like", func(t *testing.T) {
		f1 := createFunctionForSQLHandle(t, true, true)
		f2 := createFunctionForSQLHandle(t, true, true)

		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest().WithLike(sdk.Like{Pattern: &f1.Name}))
		require.NoError(t, err)

		require.Equal(t, 1, len(functions))
		require.Contains(t, functions, *f1)
		require.NotContains(t, functions, *f2)
	})

	t.Run("show function for SQL: no matches", func(t *testing.T) {
		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest().WithLike(sdk.Like{Pattern: sdk.String("non-existing-id-pattern")}))
		require.NoError(t, err)
		require.Equal(t, 0, len(functions))
	})

	t.Run("describe function for SQL", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		details, err := client.Functions.Describe(ctx, f.ID())
		require.NoError(t, err)
		pairs := make(map[string]string)
		for _, detail := range details {
			pairs[detail.Property] = *detail.Value
		}
		require.Equal(t, "SQL", pairs["language"])
		require.Equal(t, "FLOAT", pairs["returns"])
		require.Equal(t, "3.141592654::FLOAT", pairs["body"])
		require.Equal(t, "(X FLOAT)", pairs["signature"])
	})

	t.Run("describe function for SQL: no arguments", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, false)

		details, err := client.Functions.Describe(ctx, f.ID())
		require.NoError(t, err)
		pairs := make(map[string]string)
		for _, detail := range details {
			pairs[detail.Property] = *detail.Value
		}
		require.Equal(t, "SQL", pairs["language"])
		require.Equal(t, "FLOAT", pairs["returns"])
		require.Equal(t, "3.141592654::FLOAT", pairs["body"])
		require.Equal(t, "()", pairs["signature"])
	})
}

func TestInt_FunctionsShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupFunctionHandle := func(id sdk.SchemaObjectIdentifierWithArguments) func() {
		return func() {
			err := client.Functions.Drop(ctx, sdk.NewDropFunctionRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createFunctionForSQLHandle := func(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) {
		t.Helper()

		definition := testClientHelper().Function.SampleSqlDefinition(t)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeFloat)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		request := sdk.NewCreateForSQLFunctionRequest(id.SchemaObjectId(), *returns, definition).WithOrReplace(true)

		argument := sdk.NewFunctionArgumentRequest("x", nil).WithArgDataTypeOld(sdk.DataTypeFloat)
		request = request.WithArguments([]sdk.FunctionArgumentRequest{*argument})
		err := client.Functions.CreateForSQL(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeFloat)
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierWithArgumentsInSchema(id1.Name(), schema.ID(), sdk.DataTypeFloat)

		createFunctionForSQLHandle(t, id1)
		createFunctionForSQLHandle(t, id2)

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

	t.Run("show function by id - different name, same arguments", func(t *testing.T) {
		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeInt, sdk.DataTypeFloat, sdk.DataTypeVARCHAR)
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeInt, sdk.DataTypeFloat, sdk.DataTypeVARCHAR)
		e := testClientHelper().Function.CreateWithIdentifier(t, id1)
		testClientHelper().Function.CreateWithIdentifier(t, id2)

		es, err := client.Functions.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, *e, *es)
	})

	t.Run("show function by id - same name, different arguments", func(t *testing.T) {
		name := testClientHelper().Ids.Alpha()
		id1 := testClientHelper().Ids.NewSchemaObjectIdentifierWithArgumentsInSchema(name, testClientHelper().Ids.SchemaId(), sdk.DataTypeInt, sdk.DataTypeFloat, sdk.DataTypeVARCHAR)
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
			"3.8",
			"add",
		).
			WithArguments(args).
			WithFunctionDefinition("def add(A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, R, S, T, U, V, W, X, Y, Z): A + A"),
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
				"3.8",
				funcName,
			).
				WithArguments(args).
				WithFunctionDefinition(testClientHelper().Function.PythonIdentityDefinition(t, funcName, argName)),
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
}

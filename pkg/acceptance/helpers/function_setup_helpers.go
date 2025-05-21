package helpers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/stretchr/testify/require"
)

func (c *TestClient) CreateSampleJavaFunctionAndJarOnUserStage(t *testing.T) *TmpFunction {
	t.Helper()

	return c.CreateSampleJavaFunctionAndJarInLocation(t, "@~")
}

func (c *TestClient) CreateSampleJavaFunctionAndJarOnStage(t *testing.T, stage *sdk.Stage) *TmpFunction {
	t.Helper()

	return c.CreateSampleJavaFunctionAndJarInLocation(t, stage.Location())
}

// TODO [SNOW-1827324]: add TestClient ref to each specific client, so that we enhance specific client and not the base one
func (c *TestClient) CreateSampleJavaFunctionAndJarInLocation(t *testing.T, stageLocation string) *TmpFunction {
	t.Helper()
	ctx := context.Background()

	className := fmt.Sprintf("TestClassAbc%s", random.AlphaLowerN(3))
	funcName := fmt.Sprintf("echoVarchar%s", random.AlphaLowerN(3))
	argName := fmt.Sprintf("arg%s", random.AlphaLowerN(3))
	dataType := testdatatypes.DataTypeVarchar_100

	id := c.Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))
	argument := sdk.NewFunctionArgumentRequest(argName, dataType)
	dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
	returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := c.Function.SampleJavaDefinition(t, className, funcName, argName)
	jarName := fmt.Sprintf("tf-%d-%s.jar", time.Now().Unix(), random.AlphaN(5))
	targetPath := fmt.Sprintf("%s/%s", stageLocation, jarName)

	request := sdk.NewCreateForJavaFunctionRequest(id.SchemaObjectId(), *returns, handler).
		WithArguments([]sdk.FunctionArgumentRequest{*argument}).
		WithTargetPath(targetPath).
		WithFunctionDefinitionWrapped(definition)

	err := c.context.client.Functions.CreateForJava(ctx, request)
	require.NoError(t, err)
	t.Cleanup(c.Function.DropFunctionFunc(t, id))
	t.Cleanup(c.Stage.RemoveFromStageFunc(t, stageLocation, jarName))

	return &TmpFunction{
		FunctionId:    id,
		ClassName:     className,
		FuncName:      funcName,
		ArgName:       argName,
		ArgType:       dataType,
		JarName:       jarName,
		StageLocation: stageLocation,
	}
}

func (c *TestClient) CreateSampleJavaProcedureAndJarOnUserStage(t *testing.T) *TmpFunction {
	t.Helper()

	return c.CreateSampleJavaProcedureAndJarInLocation(t, "@~")
}

func (c *TestClient) CreateSampleJavaProcedureAndJarOnStage(t *testing.T, stage *sdk.Stage) *TmpFunction {
	t.Helper()

	return c.CreateSampleJavaProcedureAndJarInLocation(t, stage.Location())
}

func (c *TestClient) CreateSampleJavaProcedureAndJarInLocation(t *testing.T, stageLocation string) *TmpFunction {
	t.Helper()
	ctx := context.Background()

	className := fmt.Sprintf("TestClassAbc%s", random.AlphaLowerN(3))
	funcName := fmt.Sprintf("echoVarchar%s", random.AlphaLowerN(3))
	argName := fmt.Sprintf("arg%s", random.AlphaLowerN(3))
	dataType := testdatatypes.DataTypeVarchar_100

	id := c.Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))
	argument := sdk.NewProcedureArgumentRequest(argName, dataType)
	dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType)
	returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := c.Procedure.SampleJavaDefinition(t, className, funcName, argName)
	jarName := fmt.Sprintf("tf-%d-%s.jar", time.Now().Unix(), random.AlphaN(5))
	targetPath := fmt.Sprintf("%s/%s", stageLocation, jarName)
	packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:1.14.0")}

	request := sdk.NewCreateForJavaProcedureRequest(id.SchemaObjectId(), *returns, "11", packages, handler).
		WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
		WithTargetPath(targetPath).
		WithProcedureDefinitionWrapped(definition)

	err := c.context.client.Procedures.CreateForJava(ctx, request)
	require.NoError(t, err)
	t.Cleanup(c.Procedure.DropProcedureFunc(t, id))
	t.Cleanup(c.Stage.RemoveFromStageFunc(t, stageLocation, jarName))

	return &TmpFunction{
		FunctionId:    id,
		ClassName:     className,
		FuncName:      funcName,
		ArgName:       argName,
		ArgType:       dataType,
		JarName:       jarName,
		StageLocation: stageLocation,
	}
}

func (c *TestClient) CreateSamplePythonFunctionAndModuleOnUserStage(t *testing.T) *TmpFunction {
	t.Helper()

	return c.CreateSamplePythonFunctionAndModuleInLocation(t, "@~")
}

func (c *TestClient) CreateSamplePythonFunctionAndModuleOnStage(t *testing.T, stage *sdk.Stage) *TmpFunction {
	t.Helper()

	return c.CreateSamplePythonFunctionAndModuleInLocation(t, stage.Location())
}

func (c *TestClient) CreateSamplePythonFunctionAndModuleInLocation(t *testing.T, stageLocation string) *TmpFunction {
	t.Helper()
	ctx := context.Background()

	funcName := fmt.Sprintf("echo%s", random.AlphaLowerN(3))
	argName := fmt.Sprintf("arg%s", random.AlphaLowerN(3))
	dataType := testdatatypes.DataTypeVarchar_100

	id := c.Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))
	argument := sdk.NewFunctionArgumentRequest(argName, dataType)
	dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
	returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
	definition := c.Function.SamplePythonDefinition(t, funcName, argName)

	request := sdk.NewCreateForPythonFunctionRequest(id.SchemaObjectId(), *returns, "3.9", funcName).
		WithArguments([]sdk.FunctionArgumentRequest{*argument}).
		WithFunctionDefinitionWrapped(definition)

	err := c.context.client.Functions.CreateForPython(ctx, request)
	require.NoError(t, err)
	t.Cleanup(c.Function.DropFunctionFunc(t, id))

	// using os.CreateTemp underneath - last * in pattern is replaced with random string
	modulePattern := fmt.Sprintf("example*%s.py", random.AlphaLowerN(3))
	modulePath := c.Stage.PutInLocationWithContent(t, stageLocation, modulePattern, definition)
	moduleFileName := filepath.Base(modulePath)

	return &TmpFunction{
		FunctionId:    id,
		ModuleName:    strings.TrimSuffix(moduleFileName, ".py"),
		FuncName:      funcName,
		ArgName:       argName,
		ArgType:       dataType,
		StageLocation: stageLocation,
	}
}

type TmpFunction struct {
	FunctionId    sdk.SchemaObjectIdentifierWithArguments
	ClassName     string
	ModuleName    string
	FuncName      string
	ArgName       string
	ArgType       datatypes.DataType
	JarName       string
	StageLocation string
}

func (f *TmpFunction) JarLocation() string {
	return fmt.Sprintf("%s/%s", f.StageLocation, f.JarName)
}

func (f *TmpFunction) PythonModuleLocation() string {
	return fmt.Sprintf("%s/%s", f.StageLocation, f.PythonFileName())
}

func (f *TmpFunction) PythonFileName() string {
	return fmt.Sprintf("%s.py", f.ModuleName)
}

func (f *TmpFunction) JavaHandler() string {
	return fmt.Sprintf("%s.%s", f.ClassName, f.FuncName)
}

func (f *TmpFunction) PythonHandler() string {
	return fmt.Sprintf("%s.%s", f.ModuleName, f.FuncName)
}

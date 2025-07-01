package helpers

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/stretchr/testify/require"
)

type FunctionClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewFunctionClient(context *TestClientContext, idsGenerator *IdsGenerator) *FunctionClient {
	return &FunctionClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *FunctionClient) client() sdk.Functions {
	return c.context.client.Functions
}

func (c *FunctionClient) Create(t *testing.T, arguments ...sdk.DataType) *sdk.Function {
	t.Helper()
	return c.CreateWithIdentifier(t, c.ids.RandomSchemaObjectIdentifierWithArguments(arguments...))
}

func (c *FunctionClient) CreateWithIdentifier(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) *sdk.Function {
	t.Helper()

	return c.CreateWithRequest(t, id,
		sdk.NewCreateForSQLFunctionRequestDefinitionWrapped(
			id.SchemaObjectId(),
			*sdk.NewFunctionReturnsRequest().WithResultDataType(*sdk.NewFunctionReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeInt)),
			"SELECT 1",
		),
	)
}

// TODO [SNOW-1850370]: improve this helper (all  other types creation)
func (c *FunctionClient) CreateSecure(t *testing.T, arguments ...sdk.DataType) *sdk.Function {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifierWithArguments(arguments...)

	return c.CreateWithRequest(t, id,
		sdk.NewCreateForSQLFunctionRequestDefinitionWrapped(
			id.SchemaObjectId(),
			*sdk.NewFunctionReturnsRequest().WithResultDataType(*sdk.NewFunctionReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeInt)),
			"SELECT 1",
		).WithSecure(true),
	)
}

func (c *FunctionClient) CreateSql(t *testing.T) (*sdk.Function, func()) {
	t.Helper()
	dataType := testdatatypes.DataTypeFloat
	id := c.ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))
	return c.CreateSqlWithIdentifierAndArgument(t, id.SchemaObjectId(), dataType)
}

func (c *FunctionClient) CreateSqlWithIdentifierAndArgument(t *testing.T, id sdk.SchemaObjectIdentifier, dataType datatypes.DataType) (*sdk.Function, func()) {
	t.Helper()
	ctx := context.Background()

	idWithArgs := sdk.NewSchemaObjectIdentifierWithArgumentsInSchema(id.SchemaId(), id.Name(), sdk.LegacyDataTypeFrom(dataType))
	argName := "x"
	definition := c.SampleSqlDefinition(t)
	dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
	returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
	argument := sdk.NewFunctionArgumentRequest(argName, dataType)
	request := sdk.NewCreateForSQLFunctionRequestDefinitionWrapped(id, *returns, definition).
		WithArguments([]sdk.FunctionArgumentRequest{*argument})

	err := c.client().CreateForSQL(ctx, request)
	require.NoError(t, err)

	function, err := c.client().ShowByID(ctx, idWithArgs)
	require.NoError(t, err)

	return function, c.DropFunctionFunc(t, idWithArgs)
}

func (c *FunctionClient) CreateSqlNoArgs(t *testing.T) (*sdk.Function, func()) {
	t.Helper()
	ctx := context.Background()

	dataType := testdatatypes.DataTypeFloat
	id := c.ids.RandomSchemaObjectIdentifierWithArguments()

	definition := c.SampleSqlDefinition(t)
	dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
	returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
	request := sdk.NewCreateForSQLFunctionRequestDefinitionWrapped(id.SchemaObjectId(), *returns, definition)

	err := c.client().CreateForSQL(ctx, request)
	require.NoError(t, err)
	t.Cleanup(c.DropFunctionFunc(t, id))

	function, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return function, c.DropFunctionFunc(t, id)
}

func (c *FunctionClient) CreateJava(t *testing.T) (*sdk.Function, func()) {
	t.Helper()
	ctx := context.Background()

	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	id := c.ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))
	argument := sdk.NewFunctionArgumentRequest(argName, dataType)
	dt := sdk.NewFunctionReturnsResultDataTypeRequest(dataType)
	returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := c.SampleJavaDefinition(t, className, funcName, argName)

	request := sdk.NewCreateForJavaFunctionRequest(id.SchemaObjectId(), *returns, handler).
		WithArguments([]sdk.FunctionArgumentRequest{*argument}).
		WithFunctionDefinitionWrapped(definition)

	err := c.client().CreateForJava(ctx, request)
	require.NoError(t, err)

	function, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return function, c.DropFunctionFunc(t, id)
}

func (c *FunctionClient) CreatePythonWithRequest(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments, req *sdk.CreateForPythonFunctionRequest) *sdk.Function {
	t.Helper()
	ctx := context.Background()
	err := c.client().CreateForPython(ctx, req)
	require.NoError(t, err)

	t.Cleanup(c.DropFunctionFunc(t, id))

	function, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return function
}

func (c *FunctionClient) CreateScalaStaged(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments, dataType datatypes.DataType, importPath string, handler string) (*sdk.Function, func()) {
	t.Helper()
	ctx := context.Background()

	argName := "x"
	argument := sdk.NewFunctionArgumentRequest(argName, dataType)

	request := sdk.NewCreateForScalaFunctionRequest(id.SchemaObjectId(), dataType, handler, "2.12").
		WithArguments([]sdk.FunctionArgumentRequest{*argument}).
		WithImports([]sdk.FunctionImportRequest{*sdk.NewFunctionImportRequest().WithImport(importPath)})

	err := c.client().CreateForScala(ctx, request)
	require.NoError(t, err)

	function, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return function, c.DropFunctionFunc(t, id)
}

func (c *FunctionClient) CreateWithRequest(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments, req *sdk.CreateForSQLFunctionRequest) *sdk.Function {
	t.Helper()
	ctx := context.Background()
	argumentRequests := make([]sdk.FunctionArgumentRequest, len(id.ArgumentDataTypes()))
	for i, argumentDataType := range id.ArgumentDataTypes() {
		argumentRequests[i] = *sdk.NewFunctionArgumentRequest(c.ids.Alpha(), nil).WithArgDataTypeOld(argumentDataType)
	}
	err := c.client().CreateForSQL(ctx, req.WithArguments(argumentRequests))
	require.NoError(t, err)

	t.Cleanup(c.DropFunctionFunc(t, id))

	function, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return function
}

func (c *FunctionClient) DropFunctionFunc(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropFunctionRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *FunctionClient) Show(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) (*sdk.Function, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *FunctionClient) DescribeDetails(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) (*sdk.FunctionDetails, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().DescribeDetails(ctx, id)
}

// formatFunctionDefinition removes the first newline and replaces all tabs with spaces.
func (c *FunctionClient) formatFunctionDefinition(definition string) string {
	return strings.Replace(strings.ReplaceAll(definition, "\t", "  "), "\n", "", 1)
}

func (c *FunctionClient) SampleJavaDefinition(t *testing.T, className string, funcName string, argName string) string {
	t.Helper()

	return c.formatFunctionDefinition(fmt.Sprintf(`
	class %[1]s {
		public static String %[2]s(String %[3]s) {
			return %[3]s;
		}
	}
`, className, funcName, argName))
}

func (c *FunctionClient) SampleJavaDefinitionNoArgs(t *testing.T, className string, funcName string) string {
	t.Helper()

	return c.formatFunctionDefinition(fmt.Sprintf(`
	class %[1]s {
		public static String %[2]s() {
			return null;
		}
	}
`, className, funcName))
}

func (c *FunctionClient) SampleJavascriptDefinition(t *testing.T, argName string) string {
	t.Helper()

	return c.formatFunctionDefinition(fmt.Sprintf(`
	if (%[1]s == 0) {
		return 1;
	} else {
		return 2;
	}
`, argName))
}

func (c *FunctionClient) SampleJavascriptDefinitionNoArgs(t *testing.T) string {
	t.Helper()
	return c.formatFunctionDefinition(`
return 1;
`)
}

func (c *FunctionClient) SamplePythonDefinition(t *testing.T, funcName string, argName string) string {
	t.Helper()

	return c.formatFunctionDefinition(fmt.Sprintf(`
def %[1]s(%[2]s):
	result = ''
	for a in range(5):
		result += %[2]s
	return result
`, funcName, argName))
}

func (c *FunctionClient) SampleScalaDefinition(t *testing.T, className string, funcName string, argName string) string {
	t.Helper()

	return c.formatFunctionDefinition(fmt.Sprintf(`
	class %[1]s {
		def %[2]s(%[3]s : String): String = {
			return %[3]s
		}
	}
`, className, funcName, argName))
}

// TODO [SNOW-1850370]: use input argument like in other samples
func (c *FunctionClient) SampleSqlDefinition(t *testing.T) string {
	t.Helper()

	return c.formatFunctionDefinition("3.141592654::FLOAT")
}

func (c *FunctionClient) SampleSqlDefinitionWithArgument(t *testing.T, argName string) string {
	t.Helper()

	return c.formatFunctionDefinition(fmt.Sprintf(`
%s
`, argName))
}

func (c *FunctionClient) PythonIdentityDefinition(t *testing.T, funcName string, argName string) string {
	t.Helper()

	return c.formatFunctionDefinition(fmt.Sprintf("def %[1]s(%[2]s): %[2]s", funcName, argName))
}

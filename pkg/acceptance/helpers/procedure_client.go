package helpers

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/stretchr/testify/require"
)

type ProcedureClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewProcedureClient(context *TestClientContext, idsGenerator *IdsGenerator) *ProcedureClient {
	return &ProcedureClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ProcedureClient) client() sdk.Procedures {
	return c.context.client.Procedures
}

func (c *ProcedureClient) CreateSql(t *testing.T) (*sdk.Procedure, func()) {
	t.Helper()
	dataType := testdatatypes.DataTypeFloat
	id := c.ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))
	definition := c.SampleSqlDefinition(t)
	return c.CreateSqlWithIdentifierAndArgument(t, id.SchemaObjectId(), dataType, definition)
}

func (c *ProcedureClient) CreateSqlWithIdentifierAndArgument(t *testing.T, id sdk.SchemaObjectIdentifier, dataType datatypes.DataType, definition string) (*sdk.Procedure, func()) {
	t.Helper()
	ctx := context.Background()

	idWithArgs := sdk.NewSchemaObjectIdentifierWithArgumentsInSchema(id.SchemaId(), id.Name(), sdk.LegacyDataTypeFrom(dataType))
	argName := "x"
	dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType)
	returns := sdk.NewProcedureSQLReturnsRequest().WithResultDataType(*dt)
	argument := sdk.NewProcedureArgumentRequest(argName, dataType)

	request := sdk.NewCreateForSQLProcedureRequestDefinitionWrapped(id, *returns, definition).
		WithArguments([]sdk.ProcedureArgumentRequest{*argument})

	err := c.client().CreateForSQL(ctx, request)
	require.NoError(t, err)

	procedure, err := c.client().ShowByID(ctx, idWithArgs)
	require.NoError(t, err)

	return procedure, c.DropProcedureFunc(t, idWithArgs)
}

func (c *ProcedureClient) CreateJava(t *testing.T) (*sdk.Procedure, func()) {
	t.Helper()
	ctx := context.Background()

	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	id := c.ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))
	argument := sdk.NewProcedureArgumentRequest(argName, dataType)
	dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType)
	returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := c.SampleJavaDefinition(t, className, funcName, argName)
	packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:1.14.0")}

	request := sdk.NewCreateForJavaProcedureRequest(id.SchemaObjectId(), *returns, "11", packages, handler).
		WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
		WithProcedureDefinitionWrapped(definition)

	err := c.client().CreateForJava(ctx, request)
	require.NoError(t, err)

	procedure, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return procedure, c.DropProcedureFunc(t, id)
}

func (c *ProcedureClient) CreatePythonInSchema(t *testing.T, schemaId sdk.DatabaseObjectIdentifier) (*sdk.Procedure, func()) {
	t.Helper()
	ctx := context.Background()

	dataType := testdatatypes.DataTypeNumber_36_2
	id := c.ids.RandomSchemaObjectIdentifierWithArgumentsInSchemaNewDataTypes(schemaId, dataType)

	argName := "i"
	funcName := "dump"
	definition := c.SamplePythonDefinition(t, funcName, argName)
	dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType)
	returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
	argument := sdk.NewProcedureArgumentRequest(argName, dataType)
	packages := []sdk.ProcedurePackageRequest{
		*sdk.NewProcedurePackageRequest("snowflake-snowpark-python==1.14.0"),
	}
	request := sdk.NewCreateForPythonProcedureRequest(id.SchemaObjectId(), *returns, testvars.PythonRuntime, packages, funcName).
		WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
		WithProcedureDefinitionWrapped(definition)

	err := c.client().CreateForPython(ctx, request)
	require.NoError(t, err)

	procedure, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return procedure, c.DropProcedureFunc(t, id)
}

func (c *ProcedureClient) CreateScalaStaged(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments, dataType datatypes.DataType, importPath string, handler string) (*sdk.Procedure, func()) {
	t.Helper()
	ctx := context.Background()

	argName := "x"
	argument := sdk.NewProcedureArgumentRequest(argName, dataType)
	dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType)
	returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
	packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:1.14.0")}

	request := sdk.NewCreateForScalaProcedureRequest(id.SchemaObjectId(), *returns, "2.12", packages, handler).
		WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
		WithImports([]sdk.ProcedureImportRequest{*sdk.NewProcedureImportRequest(importPath)})

	err := c.client().CreateForScala(ctx, request)
	require.NoError(t, err)

	function, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return function, c.DropProcedureFunc(t, id)
}

func (c *ProcedureClient) Create(t *testing.T, arguments ...sdk.DataType) (*sdk.Procedure, func()) {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifierWithArguments(arguments...)
	return c.CreateWithIdentifier(t, id)
}

func (c *ProcedureClient) CreateInSchema(t *testing.T, schemaId sdk.DatabaseObjectIdentifier, arguments ...sdk.DataType) (*sdk.Procedure, func()) {
	t.Helper()
	return c.CreateWithIdentifier(t, c.ids.RandomSchemaObjectIdentifierWithArgumentsInSchema(schemaId, arguments...))
}

func (c *ProcedureClient) CreateWithIdentifier(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) (*sdk.Procedure, func()) {
	t.Helper()
	ctx := context.Background()
	argumentRequests := make([]sdk.ProcedureArgumentRequest, len(id.ArgumentDataTypes()))
	for i, argumentDataType := range id.ArgumentDataTypes() {
		argumentRequests[i] = *sdk.NewProcedureArgumentRequest(c.ids.Alpha(), nil).WithArgDataTypeOld(argumentDataType)
	}
	err := c.client().CreateForSQL(ctx,
		sdk.NewCreateForSQLProcedureRequestDefinitionWrapped(
			id.SchemaObjectId(),
			*sdk.NewProcedureSQLReturnsRequest().WithResultDataType(*sdk.NewProcedureReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeInt)),
			`BEGIN RETURN 1; END`).WithArguments(argumentRequests),
	)
	require.NoError(t, err)

	procedure, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return procedure, c.DropProcedureFunc(t, id)
}

func (c *ProcedureClient) DropProcedureFunc(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropProcedureRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *ProcedureClient) Show(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) (*sdk.Procedure, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *ProcedureClient) DescribeDetails(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) (*sdk.ProcedureDetails, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().DescribeDetails(ctx, id)
}

// formatProcedureDefinition removes the first newline and replaces all tabs with spaces.
func (c *ProcedureClient) formatProcedureDefinition(definition string) string {
	return strings.Replace(strings.ReplaceAll(definition, "\t", "  "), "\n", "", 1)
}

// Session argument is needed: https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-java#data-access-example
// More references: https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-java
func (c *ProcedureClient) SampleJavaDefinition(t *testing.T, className string, funcName string, argName string) string {
	t.Helper()

	return c.formatProcedureDefinition(fmt.Sprintf(`
	import com.snowflake.snowpark_java.*;
	class %[1]s {
		public static String %[2]s(Session session, String %[3]s) {
			return %[3]s;
		}
	}
`, className, funcName, argName))
}

func (c *ProcedureClient) SampleJavaDefinitionNoArgs(t *testing.T, className string, funcName string) string {
	t.Helper()

	return c.formatProcedureDefinition(fmt.Sprintf(`
	import com.snowflake.snowpark_java.*;
	class %[1]s {
		public static String %[2]s(Session session) {
			return null;
		}
	}
`, className, funcName))
}

// For more references: https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-javascript
func (c *ProcedureClient) SampleJavascriptDefinition(t *testing.T, argName string) string {
	t.Helper()

	return c.formatProcedureDefinition(fmt.Sprintf(`
	if (%[1]s == 0) {
		return 1;
	} else {
		return 2;
	}
`, argName))
}

func (c *ProcedureClient) SamplePythonDefinition(t *testing.T, funcName string, argName string) string {
	t.Helper()

	return c.formatProcedureDefinition(fmt.Sprintf(`
def %[1]s(%[2]s):
	result = ''
	for a in range(5):
		result += %[2]s
	return result
`, funcName, argName))
}

// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-scala
func (c *ProcedureClient) SampleScalaDefinition(t *testing.T, className string, funcName string, argName string) string {
	t.Helper()

	return c.formatProcedureDefinition(fmt.Sprintf(`
	import com.snowflake.snowpark_java.Session
	class %[1]s {
		def %[2]s(session : Session, %[3]s : String): String = {
			return %[3]s
		}
	}
`, className, funcName, argName))
}

func (c *ProcedureClient) SampleSqlDefinition(t *testing.T) string {
	t.Helper()

	return c.formatProcedureDefinition(`
BEGIN
	RETURN 3.141592654::FLOAT;
END;
`)
}

func (c *ProcedureClient) SampleSqlDefinitionWithArgument(t *testing.T) string {
	t.Helper()

	return c.formatProcedureDefinition(`
BEGIN
  RETURN message;
END;
`)
}

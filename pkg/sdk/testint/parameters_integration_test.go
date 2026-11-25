//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_Parameters_ObjectParameterWithGenericFallback(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	assertParameterValue := func(t *testing.T, object sdk.Object, parameter sdk.ObjectParameter, expected string) {
		t.Helper()
		p, err := client.Parameters.ShowObjectParameter(ctx, parameter, object)
		require.NoError(t, err)
		require.Equal(t, expected, p.Value)
	}
	assertParameterUnset := func(t *testing.T, object sdk.Object, parameter sdk.ObjectParameter) {
		t.Helper()
		p, err := client.Parameters.ShowObjectParameter(ctx, parameter, object)
		require.NoError(t, err)
		require.Equal(t, p.Default, p.Value)
		require.NotEqual(t, sdk.ParameterType(object.ObjectType), p.Level)
	}

	t.Run("supported object type and parameter: the catalog write path is used, and the kind is not", func(t *testing.T) {
		database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseCleanup)
		object := sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: database.ID()}
		parameter := sdk.ObjectParameterDataRetentionTimeInDays

		// The catalog parses the value as an integer, so STRING is not used.
		require.NoError(t, client.Parameters.SetObjectParameterOnObjectWithGenericFallback(ctx, object, parameter, sdk.ParameterValueKindString, "2"))
		assertParameterValue(t, object, parameter, "2")

		require.NoError(t, client.Parameters.UnsetObjectParameterOnObjectWithGenericFallback(ctx, object, parameter))
		assertParameterUnset(t, object, parameter)
	})

	t.Run("supported object type and unsupported parameter: the value is rendered by the kind", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)
		object := sdk.Object{ObjectType: sdk.ObjectTypeSchema, Name: schema.ID()}
		parameter := sdk.ObjectParameter("ICEBERG_DEFAULT_DDL_COLLATION")

		require.NoError(t, client.Parameters.SetObjectParameterOnObjectWithGenericFallback(ctx, object, parameter, sdk.ParameterValueKindString, "en-ci"))
		assertParameterValue(t, object, parameter, "en-ci")

		require.NoError(t, client.Parameters.UnsetObjectParameterOnObjectWithGenericFallback(ctx, object, parameter))
		assertParameterUnset(t, object, parameter)
	})

	t.Run("object type without a catalog writer: the value is rendered by the kind", func(t *testing.T) {
		table, tableCleanup := testClientHelper().Table.Create(t)
		t.Cleanup(tableCleanup)
		object := sdk.Object{ObjectType: sdk.ObjectTypeTable, Name: table.ID()}
		parameter := sdk.ObjectParameterDefaultDdlCollation

		require.NoError(t, client.Parameters.SetObjectParameterOnObjectWithGenericFallback(ctx, object, parameter, sdk.ParameterValueKindString, "en-ci"))
		assertParameterValue(t, object, parameter, "en-ci")

		require.NoError(t, client.Parameters.UnsetObjectParameterOnObjectWithGenericFallback(ctx, object, parameter))
		assertParameterUnset(t, object, parameter)
	})

	// ALERT has no dedicated field in ParametersIn, so SHOW PARAMETERS uses the generic IN <object_type> <name> form.
	t.Run("object type without a dedicated SHOW PARAMETERS field", func(t *testing.T) {
		alert, alertCleanup := testClientHelper().Alert.CreateAlert(t)
		t.Cleanup(alertCleanup)
		object := sdk.Object{ObjectType: sdk.ObjectTypeAlert, Name: alert.ID()}

		parameters, err := client.Parameters.ShowParameters(ctx, &sdk.ShowParametersOptions{In: &sdk.ParametersIn{Object: &sdk.ParametersInObject{ObjectType: object.ObjectType, Name: object.Name}}})
		require.NoError(t, err)
		require.NotEmpty(t, parameters)

		p, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameter(parameters[0].Key), object)
		require.NoError(t, err)
		require.Equal(t, parameters[0].Key, p.Key)
	})
}

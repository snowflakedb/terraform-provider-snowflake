//go:build non_account_level_tests

package testint

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_ExternalFunctions(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	defaultDataTypes := []sdk.DataType{sdk.DataTypeVARCHAR}

	integration, integrationCleanup := testClientHelper().ApiIntegration.CreateApiIntegration(t)
	t.Cleanup(integrationCleanup)

	cleanupExternalFunctionHandle := func(id sdk.SchemaObjectIdentifierWithArguments) func() {
		return func() {
			err := client.Functions.Drop(ctx, sdk.NewDropFunctionRequest(id).WithIfExists(true))
			require.NoError(t, err)
		}
	}

	createExternalFunction := func(t *testing.T) *sdk.ExternalFunction {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(defaultDataTypes[0])
		argument := sdk.NewExternalFunctionArgumentRequest("x", defaultDataTypes[0])
		as := "https://xyz.execute-api.us-west-2.amazonaws.com/production/remote_echo"
		request := sdk.NewCreateExternalFunctionRequest(id.SchemaObjectId(), sdk.DataTypeVariant, sdk.Pointer(integration.ID()), as).
			WithOrReplace(true).
			WithSecure(true).
			WithArguments([]sdk.ExternalFunctionArgumentRequest{*argument})
		err := client.ExternalFunctions.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupExternalFunctionHandle(id))

		e, err := client.ExternalFunctions.ShowByID(ctx, id)
		require.NoError(t, err)
		return e
	}

	assertExternalFunction := func(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments, secure bool) {
		t.Helper()

		e, err := client.ExternalFunctions.ShowByID(ctx, id)
		require.NoError(t, err)

		require.NotEmpty(t, e.CreatedOn)
		require.Equal(t, id, e.ID())
		require.False(t, e.IsBuiltin)
		require.False(t, e.IsAggregate)
		require.False(t, e.IsAnsi)
		if len(id.ArgumentDataTypes()) > 0 {
			require.Equal(t, 1, e.MinNumArguments)
			require.Equal(t, 1, e.MaxNumArguments)
		} else {
			require.Equal(t, 0, e.MinNumArguments)
			require.Equal(t, 0, e.MaxNumArguments)
		}
		require.NotEmpty(t, e.Description)
		require.NotEmpty(t, e.CatalogName)
		require.False(t, e.IsTableFunction)
		require.False(t, e.ValidForClustering)
		require.Equal(t, secure, e.IsSecure)
		require.True(t, e.IsExternalFunction)
		require.Equal(t, "EXTERNAL", e.Language)
		require.False(t, e.IsMemoizable)
		require.False(t, e.IsDataMetric)
	}

	t.Run("create external function", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(defaultDataTypes...)
		argument := sdk.NewExternalFunctionArgumentRequest("x", sdk.DataTypeVARCHAR)
		headers := []sdk.ExternalFunctionHeaderRequest{
			{
				Name:  "measure",
				Value: "kilometers",
			},
		}
		ch := []sdk.ExternalFunctionContextHeaderRequest{
			{
				ContextFunction: "CURRENT_DATE",
			},
			{
				ContextFunction: "CURRENT_TIMESTAMP",
			},
		}
		as := "https://xyz.execute-api.us-west-2.amazonaws.com/production/remote_echo"
		request := sdk.NewCreateExternalFunctionRequest(id.SchemaObjectId(), sdk.DataTypeVariant, sdk.Pointer(integration.ID()), as).
			WithOrReplace(true).
			WithSecure(true).
			WithArguments([]sdk.ExternalFunctionArgumentRequest{*argument}).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorCalledOnNullInput)).
			WithHeaders(headers).
			WithContextHeaders(ch).
			WithMaxBatchRows(10).
			WithCompression("GZIP")
		err := client.ExternalFunctions.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupExternalFunctionHandle(id))

		assertExternalFunction(t, id, true)
	})

	t.Run("create external function without arguments", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments()
		as := "https://xyz.execute-api.us-west-2.amazonaws.com/production/remote_echo"
		request := sdk.NewCreateExternalFunctionRequest(id.SchemaObjectId(), sdk.DataTypeVariant, sdk.Pointer(integration.ID()), as)
		err := client.ExternalFunctions.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupExternalFunctionHandle(id))

		assertExternalFunction(t, id, false)
	})

	t.Run("alter external function: set api integration", func(t *testing.T) {
		externalFunction := createExternalFunction(t)
		set := sdk.NewExternalFunctionSetRequest().
			WithApiIntegration(integration.ID())
		request := sdk.NewAlterExternalFunctionRequest(externalFunction.ID()).WithSet(*set)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)

		assertExternalFunction(t, externalFunction.ID(), true)
	})

	t.Run("alter external function: set headers", func(t *testing.T) {
		externalFunction := createExternalFunction(t)

		headers := []sdk.ExternalFunctionHeaderRequest{
			{
				Name:  "measure",
				Value: "kilometers",
			},
		}
		set := sdk.NewExternalFunctionSetRequest().WithHeaders(headers)
		request := sdk.NewAlterExternalFunctionRequest(externalFunction.ID()).WithSet(*set)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)
		assertExternalFunction(t, externalFunction.ID(), true)
	})

	t.Run("alter external function: set context headers", func(t *testing.T) {
		externalFunction := createExternalFunction(t)

		ch := []sdk.ExternalFunctionContextHeaderRequest{
			{
				ContextFunction: "CURRENT_DATE",
			},
			{
				ContextFunction: "CURRENT_TIMESTAMP",
			},
		}
		set := sdk.NewExternalFunctionSetRequest().WithContextHeaders(ch)
		request := sdk.NewAlterExternalFunctionRequest(externalFunction.ID()).WithSet(*set)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)
		assertExternalFunction(t, externalFunction.ID(), true)
	})

	t.Run("alter external function: set compression", func(t *testing.T) {
		externalFunction := createExternalFunction(t)

		set := sdk.NewExternalFunctionSetRequest().WithCompression("AUTO")
		request := sdk.NewAlterExternalFunctionRequest(externalFunction.ID()).WithSet(*set)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)
		assertExternalFunction(t, externalFunction.ID(), true)
	})

	t.Run("alter external function: set max batch rows", func(t *testing.T) {
		externalFunction := createExternalFunction(t)

		set := sdk.NewExternalFunctionSetRequest().WithMaxBatchRows(20)
		request := sdk.NewAlterExternalFunctionRequest(externalFunction.ID()).WithSet(*set)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)
		assertExternalFunction(t, externalFunction.ID(), true)
	})

	t.Run("alter external function: unset", func(t *testing.T) {
		externalFunction := createExternalFunction(t)

		unset := sdk.NewExternalFunctionUnsetRequest().
			WithComment(true).
			WithHeaders(true)
		request := sdk.NewAlterExternalFunctionRequest(externalFunction.ID()).WithUnset(*unset)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)

		assertExternalFunction(t, externalFunction.ID(), true)
	})

	t.Run("show external function: with like", func(t *testing.T) {
		e1 := createExternalFunction(t)
		e2 := createExternalFunction(t)

		es, err := client.ExternalFunctions.Show(ctx, sdk.NewShowExternalFunctionRequest().WithLike(sdk.Like{Pattern: sdk.String(e1.Name)}))
		require.NoError(t, err)

		require.Len(t, es, 1)
		require.Contains(t, es, *e1)
		require.NotContains(t, es, *e2)
	})

	t.Run("show external function: with in", func(t *testing.T) {
		otherDb, otherDbCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(otherDbCleanup)

		e1 := createExternalFunction(t)

		es, err := client.ExternalFunctions.Show(ctx, sdk.NewShowExternalFunctionRequest().WithIn(sdk.In{Schema: e1.ID().SchemaId()}))
		require.NoError(t, err)

		require.Contains(t, es, *e1)

		es, err = client.ExternalFunctions.Show(ctx, sdk.NewShowExternalFunctionRequest().WithIn(sdk.In{Database: testClientHelper().Ids.DatabaseId()}))
		require.NoError(t, err)

		require.Contains(t, es, *e1)

		es, err = client.ExternalFunctions.Show(ctx, sdk.NewShowExternalFunctionRequest().WithIn(sdk.In{Database: otherDb.ID()}))
		require.NoError(t, err)

		require.Empty(t, es)
	})

	t.Run("show external function: no matches", func(t *testing.T) {
		es, err := client.ExternalFunctions.Show(ctx, sdk.NewShowExternalFunctionRequest().WithLike(sdk.Like{Pattern: sdk.String("non-existing-id-pattern")}))
		require.NoError(t, err)
		require.Empty(t, es)
	})

	t.Run("show external function by id", func(t *testing.T) {
		e := createExternalFunction(t)

		es, err := client.ExternalFunctions.ShowByID(ctx, e.ID())
		require.NoError(t, err)
		require.Equal(t, *e, *es)
	})

	t.Run("show external function by id - different name, same arguments", func(t *testing.T) {
		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeInt, sdk.DataTypeFloat, sdk.DataTypeVARCHAR)
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeInt, sdk.DataTypeFloat, sdk.DataTypeVARCHAR)
		e := testClientHelper().ExternalFunction.CreateWithIdentifier(t, integration.ID(), id1)
		testClientHelper().ExternalFunction.CreateWithIdentifier(t, integration.ID(), id2)

		es, err := client.ExternalFunctions.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, *e, *es)
	})

	t.Run("show external function by id - same name, different arguments", func(t *testing.T) {
		name := testClientHelper().Ids.Alpha()
		id1 := testClientHelper().Ids.NewSchemaObjectIdentifierWithArgumentsInSchema(name, testClientHelper().Ids.SchemaId(), sdk.DataTypeInt, sdk.DataTypeFloat, sdk.DataTypeVARCHAR)
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierWithArgumentsInSchema(name, testClientHelper().Ids.SchemaId(), sdk.DataTypeInt, sdk.DataTypeVARCHAR)
		e := testClientHelper().ExternalFunction.CreateWithIdentifier(t, integration.ID(), id1)
		testClientHelper().ExternalFunction.CreateWithIdentifier(t, integration.ID(), id2)

		es, err := client.ExternalFunctions.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, *e, *es)
	})

	t.Run("describe external function", func(t *testing.T) {
		e := createExternalFunction(t)

		details, err := client.ExternalFunctions.Describe(ctx, e.ID())
		require.NoError(t, err)
		pairs := make(map[string]string)
		for _, detail := range details {
			pairs[detail.Property] = detail.Value
		}
		require.Equal(t, "EXTERNAL", pairs["language"])
		require.Equal(t, "VARIANT", pairs["returns"])
		require.Equal(t, "VOLATILE", pairs["volatility"])
		require.Equal(t, "AUTO", pairs["compression"])
		require.Equal(t, "(X VARCHAR)", pairs["signature"])
	})
}

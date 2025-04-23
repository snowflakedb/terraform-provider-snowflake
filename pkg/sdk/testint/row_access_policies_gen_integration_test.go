//go:build !account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_RowAccessPolicies(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	assertRowAccessPolicy := func(t *testing.T, rowAccessPolicy *sdk.RowAccessPolicy, id sdk.SchemaObjectIdentifier, comment string) {
		t.Helper()
		assert.NotEmpty(t, rowAccessPolicy.CreatedOn)
		assert.Equal(t, id.Name(), rowAccessPolicy.Name)
		assert.Equal(t, id.DatabaseName(), rowAccessPolicy.DatabaseName)
		assert.Equal(t, id.SchemaName(), rowAccessPolicy.SchemaName)
		assert.Equal(t, "ROW_ACCESS_POLICY", rowAccessPolicy.Kind)
		assert.Equal(t, "ACCOUNTADMIN", rowAccessPolicy.Owner)
		assert.Equal(t, comment, rowAccessPolicy.Comment)
		assert.Empty(t, rowAccessPolicy.Options)
		assert.Equal(t, "ROLE", rowAccessPolicy.OwnerRoleType)
	}

	assertRowAccessPolicyDescription := func(t *testing.T, rowAccessPolicyDescription *sdk.RowAccessPolicyDescription, id sdk.SchemaObjectIdentifier, signature []sdk.TableColumnSignature, expectedBody string) {
		t.Helper()
		assert.Equal(t, sdk.RowAccessPolicyDescription{
			Name:       id.Name(),
			Signature:  signature,
			ReturnType: "BOOLEAN",
			Body:       expectedBody,
		}, *rowAccessPolicyDescription)
	}

	createRowAccessPolicyRequest := func(t *testing.T, args []sdk.CreateRowAccessPolicyArgsRequest, body string) *sdk.CreateRowAccessPolicyRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		return sdk.NewCreateRowAccessPolicyRequest(id, args, body)
	}

	createRowAccessPolicyBasicRequest := func(t *testing.T) *sdk.CreateRowAccessPolicyRequest {
		t.Helper()

		argName := random.AlphaN(5)
		argType := testdatatypes.DataTypeVarchar
		args := sdk.NewCreateRowAccessPolicyArgsRequest(argName, argType)

		body := "true"

		return createRowAccessPolicyRequest(t, []sdk.CreateRowAccessPolicyArgsRequest{*args}, body)
	}

	t.Run("create row access policy: no optionals", func(t *testing.T) {
		request := createRowAccessPolicyBasicRequest(t)

		rowAccessPolicy, cleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicyWithRequest(t, *request)
		t.Cleanup(cleanup)

		assertRowAccessPolicy(t, rowAccessPolicy, request.GetName(), "")
	})

	t.Run("create row access policy: full", func(t *testing.T) {
		request := createRowAccessPolicyBasicRequest(t).WithComment(sdk.Pointer("some comment"))

		rowAccessPolicy, cleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicyWithRequest(t, *request)
		t.Cleanup(cleanup)

		assertRowAccessPolicy(t, rowAccessPolicy, request.GetName(), "some comment")
	})

	t.Run("drop row access policy: existing", func(t *testing.T) {
		request := createRowAccessPolicyBasicRequest(t)
		id := request.GetName()

		err := client.RowAccessPolicies.Create(ctx, request)
		require.NoError(t, err)

		err = client.RowAccessPolicies.Drop(ctx, sdk.NewDropRowAccessPolicyRequest(id))
		require.NoError(t, err)

		_, err = client.RowAccessPolicies.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop row access policy: non-existing", func(t *testing.T) {
		err := client.RowAccessPolicies.Drop(ctx, sdk.NewDropRowAccessPolicyRequest(NonExistingSchemaObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter row access policy: rename", func(t *testing.T) {
		createRequest := createRowAccessPolicyBasicRequest(t)
		id := createRequest.GetName()

		err := client.RowAccessPolicies.Create(ctx, createRequest)
		require.NoError(t, err)

		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		alterRequest := sdk.NewAlterRowAccessPolicyRequest(id).WithRenameTo(&newId)

		err = client.RowAccessPolicies.Alter(ctx, alterRequest)
		if err != nil {
			t.Cleanup(testClientHelper().RowAccessPolicy.DropRowAccessPolicyFunc(t, id))
		} else {
			t.Cleanup(testClientHelper().RowAccessPolicy.DropRowAccessPolicyFunc(t, newId))
		}
		require.NoError(t, err)

		_, err = client.RowAccessPolicies.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		rowAccessPolicy, err := client.RowAccessPolicies.ShowByID(ctx, newId)
		require.NoError(t, err)

		assertRowAccessPolicy(t, rowAccessPolicy, newId, "")
	})

	t.Run("alter row access policy: set and unset comment", func(t *testing.T) {
		rowAccessPolicy, cleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicy(t)
		t.Cleanup(cleanup)
		id := rowAccessPolicy.ID()

		alterRequest := sdk.NewAlterRowAccessPolicyRequest(id).WithSetComment(sdk.String("new comment"))
		err := client.RowAccessPolicies.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredRowAccessPolicy, err := client.RowAccessPolicies.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "new comment", alteredRowAccessPolicy.Comment)

		alterRequest = sdk.NewAlterRowAccessPolicyRequest(id).WithUnsetComment(sdk.Bool(true))
		err = client.RowAccessPolicies.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredRowAccessPolicy, err = client.RowAccessPolicies.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "", alteredRowAccessPolicy.Comment)
	})

	t.Run("alter row access policy: set body", func(t *testing.T) {
		rowAccessPolicy, cleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicy(t)
		t.Cleanup(cleanup)
		id := rowAccessPolicy.ID()

		alterRequest := sdk.NewAlterRowAccessPolicyRequest(id).WithSetBody(sdk.String("false"))
		err := client.RowAccessPolicies.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredRowAccessPolicyDescription, err := client.RowAccessPolicies.Describe(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "false", alteredRowAccessPolicyDescription.Body)

		alterRequest = sdk.NewAlterRowAccessPolicyRequest(id).WithSetBody(sdk.String("true"))
		err = client.RowAccessPolicies.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredRowAccessPolicyDescription, err = client.RowAccessPolicies.Describe(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "true", alteredRowAccessPolicyDescription.Body)
	})

	t.Run("show row access policy: default", func(t *testing.T) {
		rowAccessPolicy1, cleanup1 := testClientHelper().RowAccessPolicy.CreateRowAccessPolicy(t)
		t.Cleanup(cleanup1)
		rowAccessPolicy2, cleanup2 := testClientHelper().RowAccessPolicy.CreateRowAccessPolicy(t)
		t.Cleanup(cleanup2)

		showRequest := sdk.NewShowRowAccessPolicyRequest()
		returnedRowAccessPolicies, err := client.RowAccessPolicies.Show(ctx, showRequest)
		require.NoError(t, err)
		require.LessOrEqual(t, 2, len(returnedRowAccessPolicies))
		assert.Contains(t, returnedRowAccessPolicies, *rowAccessPolicy1)
		assert.Contains(t, returnedRowAccessPolicies, *rowAccessPolicy2)
	})

	t.Run("show row access policy: with options", func(t *testing.T) {
		rowAccessPolicy1, cleanup1 := testClientHelper().RowAccessPolicy.CreateRowAccessPolicy(t)
		t.Cleanup(cleanup1)
		rowAccessPolicy2, cleanup2 := testClientHelper().RowAccessPolicy.CreateRowAccessPolicy(t)
		t.Cleanup(cleanup2)

		showRequest := sdk.NewShowRowAccessPolicyRequest().
			WithLike(sdk.Like{Pattern: &rowAccessPolicy1.Name}).
			WithIn(sdk.ExtendedIn{In: sdk.In{Schema: testClientHelper().Ids.SchemaId()}}).
			WithLimit(&sdk.LimitFrom{Rows: sdk.Int(5)})
		returnedRowAccessPolicies, err := client.RowAccessPolicies.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Equal(t, 1, len(returnedRowAccessPolicies))
		assert.Contains(t, returnedRowAccessPolicies, *rowAccessPolicy1)
		assert.NotContains(t, returnedRowAccessPolicies, *rowAccessPolicy2)
	})

	t.Run("describe row access policy: existing", func(t *testing.T) {
		argName := random.AlphaN(5)
		argType := testdatatypes.DataTypeVarchar
		args := sdk.NewCreateRowAccessPolicyArgsRequest(argName, argType)
		body := "true"

		request := createRowAccessPolicyRequest(t, []sdk.CreateRowAccessPolicyArgsRequest{*args}, body)
		rowAccessPolicy, cleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicyWithRequest(t, *request)
		t.Cleanup(cleanup)

		returnedRowAccessPolicyDescription, err := client.RowAccessPolicies.Describe(ctx, rowAccessPolicy.ID())
		require.NoError(t, err)

		assertRowAccessPolicyDescription(t, returnedRowAccessPolicyDescription, rowAccessPolicy.ID(), []sdk.TableColumnSignature{{
			Name: argName,
			Type: argType,
		}}, body)
	})

	t.Run("describe row access policy: with timestamp data type normalization", func(t *testing.T) {
		argName := random.AlphaN(5)
		argType := testdatatypes.DataTypeTimestampLTZ
		args := sdk.NewCreateRowAccessPolicyArgsRequest(argName, argType)
		body := "true"

		request := createRowAccessPolicyRequest(t, []sdk.CreateRowAccessPolicyArgsRequest{*args}, body)
		rowAccessPolicy, cleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicyWithRequest(t, *request)
		t.Cleanup(cleanup)

		returnedRowAccessPolicyDescription, err := client.RowAccessPolicies.Describe(ctx, rowAccessPolicy.ID())
		require.NoError(t, err)

		assertRowAccessPolicyDescription(t, returnedRowAccessPolicyDescription, rowAccessPolicy.ID(), []sdk.TableColumnSignature{{
			Name: argName,
			Type: testdatatypes.DataTypeTimestampLTZ,
		}}, body)
	})

	t.Run("describe row access policy: with varchar data type normalization", func(t *testing.T) {
		argName := random.AlphaN(5)
		argType := testdatatypes.DataTypeVarchar_200
		args := sdk.NewCreateRowAccessPolicyArgsRequest(argName, argType)
		body := "true"

		request := createRowAccessPolicyRequest(t, []sdk.CreateRowAccessPolicyArgsRequest{*args}, body)
		rowAccessPolicy, cleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicyWithRequest(t, *request)
		t.Cleanup(cleanup)

		returnedRowAccessPolicyDescription, err := client.RowAccessPolicies.Describe(ctx, rowAccessPolicy.ID())
		require.NoError(t, err)

		assertRowAccessPolicyDescription(t, returnedRowAccessPolicyDescription, rowAccessPolicy.ID(), []sdk.TableColumnSignature{{
			Name: argName,
			Type: testdatatypes.DataTypeVarchar,
		}}, body)
	})

	t.Run("describe row access policy: non-existing", func(t *testing.T) {
		_, err := client.RowAccessPolicies.Describe(ctx, NonExistingSchemaObjectIdentifier)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_RowAccessPoliciesShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		body := "true"
		argName := random.AlphaN(5)
		argType := testdatatypes.DataTypeVarchar
		arg := sdk.NewCreateRowAccessPolicyArgsRequest(argName, argType)

		req1 := sdk.NewCreateRowAccessPolicyRequest(id1, []sdk.CreateRowAccessPolicyArgsRequest{*arg}, body)
		req2 := sdk.NewCreateRowAccessPolicyRequest(id2, []sdk.CreateRowAccessPolicyArgsRequest{*arg}, body)
		_, cleanup1 := testClientHelper().RowAccessPolicy.CreateRowAccessPolicyWithRequest(t, *req1)
		t.Cleanup(cleanup1)
		_, cleanup2 := testClientHelper().RowAccessPolicy.CreateRowAccessPolicyWithRequest(t, *req2)
		t.Cleanup(cleanup2)

		e1, err := client.RowAccessPolicies.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.RowAccessPolicies.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}

// TODO [next PR]: improve this test
func TestInt_RowAccessPoliciesDescribe(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("describe", func(t *testing.T) {
		args := []sdk.CreateRowAccessPolicyArgsRequest{
			*sdk.NewCreateRowAccessPolicyArgsRequest("A", testdatatypes.DataTypeNumber_2_0),
			*sdk.NewCreateRowAccessPolicyArgsRequest("B", testdatatypes.DataTypeDecimal),
			*sdk.NewCreateRowAccessPolicyArgsRequest("C", testdatatypes.DataTypeInteger),
			*sdk.NewCreateRowAccessPolicyArgsRequest("D", testdatatypes.DataTypeFloat),
			*sdk.NewCreateRowAccessPolicyArgsRequest("E", testdatatypes.DataTypeDouble),
			*sdk.NewCreateRowAccessPolicyArgsRequest("F", testdatatypes.DataTypeVarchar_100),
			*sdk.NewCreateRowAccessPolicyArgsRequest("G", testdatatypes.DataTypeChar),
			*sdk.NewCreateRowAccessPolicyArgsRequest("H", testdatatypes.DataTypeString),
			*sdk.NewCreateRowAccessPolicyArgsRequest("I", testdatatypes.DataTypeText),
			*sdk.NewCreateRowAccessPolicyArgsRequest("J", testdatatypes.DataTypeBinary),
			*sdk.NewCreateRowAccessPolicyArgsRequest("K", testdatatypes.DataTypeVarbinary),
			*sdk.NewCreateRowAccessPolicyArgsRequest("L", testdatatypes.DataTypeBoolean),
			*sdk.NewCreateRowAccessPolicyArgsRequest("M", testdatatypes.DataTypeDate),
			*sdk.NewCreateRowAccessPolicyArgsRequest("N", testdatatypes.DataTypeDatetime),
			*sdk.NewCreateRowAccessPolicyArgsRequest("O", testdatatypes.DataTypeTime),
			*sdk.NewCreateRowAccessPolicyArgsRequest("R", testdatatypes.DataTypeTimestampLTZ),
			*sdk.NewCreateRowAccessPolicyArgsRequest("S", testdatatypes.DataTypeTimestampNTZ),
			*sdk.NewCreateRowAccessPolicyArgsRequest("T", testdatatypes.DataTypeTimestampTZ),
			*sdk.NewCreateRowAccessPolicyArgsRequest("U", testdatatypes.DataTypeVariant),
			*sdk.NewCreateRowAccessPolicyArgsRequest("V", testdatatypes.DataTypeObject),
			*sdk.NewCreateRowAccessPolicyArgsRequest("W", testdatatypes.DataTypeArray),
			*sdk.NewCreateRowAccessPolicyArgsRequest("X", testdatatypes.DataTypeGeography),
			*sdk.NewCreateRowAccessPolicyArgsRequest("Y", testdatatypes.DataTypeGeometry),
			// TODO(SNOW-1596962): Fully support VECTOR data type sdk.ParseFunctionArgumentsFromString could be a base for another function that takes argument names into consideration.
			// *sdk.NewCreateRowAccessPolicyArgsRequest("Z", "VECTOR(INT, 16)"),
		}

		policy, cleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicyWithArguments(t, args)
		t.Cleanup(cleanup)

		id := policy.ID()
		policyDetails, err := client.RowAccessPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "true", policyDetails.Body)
		assert.Equal(t, id.Name(), policyDetails.Name)
		assert.Equal(t, "BOOLEAN", policyDetails.ReturnType)
		require.NoError(t, err)
		wantArgs := make([]sdk.TableColumnSignature, len(args))
		for i, arg := range args {
			wantType, err := datatypes.ParseDataType(arg.Type.ToLegacyDataTypeSql())
			require.NoError(t, err)
			wantArgs[i] = sdk.TableColumnSignature{
				Name: arg.Name,
				Type: wantType,
			}
		}
		assert.Equal(t, wantArgs, policyDetails.Signature)
	})
}

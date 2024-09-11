package testint

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

	assertRowAccessPolicyDescription := func(t *testing.T, rowAccessPolicyDescription *sdk.RowAccessPolicyDescription, id sdk.SchemaObjectIdentifier, expectedSignature string, expectedBody string) {
		t.Helper()
		assert.Equal(t, sdk.RowAccessPolicyDescription{
			Name:       id.Name(),
			Signature:  expectedSignature,
			ReturnType: "BOOLEAN",
			Body:       expectedBody,
		}, *rowAccessPolicyDescription)
	}

	cleanupRowAccessPolicyProvider := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.RowAccessPolicies.Drop(ctx, sdk.NewDropRowAccessPolicyRequest(id))
			require.NoError(t, err)
		}
	}

	createRowAccessPolicyRequest := func(t *testing.T, args []sdk.CreateRowAccessPolicyArgsRequest, body string) *sdk.CreateRowAccessPolicyRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		return sdk.NewCreateRowAccessPolicyRequest(id, args, body)
	}

	createRowAccessPolicyBasicRequest := func(t *testing.T) *sdk.CreateRowAccessPolicyRequest {
		t.Helper()

		argName := random.AlphaN(5)
		argType := sdk.DataTypeVARCHAR
		args := sdk.NewCreateRowAccessPolicyArgsRequest(argName, argType)

		body := "true"

		return createRowAccessPolicyRequest(t, []sdk.CreateRowAccessPolicyArgsRequest{*args}, body)
	}

	createRowAccessPolicyWithRequest := func(t *testing.T, request *sdk.CreateRowAccessPolicyRequest) *sdk.RowAccessPolicy {
		t.Helper()
		id := request.GetName()

		err := client.RowAccessPolicies.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupRowAccessPolicyProvider(id))

		rowAccessPolicy, err := client.RowAccessPolicies.ShowByID(ctx, id)
		require.NoError(t, err)

		return rowAccessPolicy
	}

	createRowAccessPolicy := func(t *testing.T) *sdk.RowAccessPolicy {
		t.Helper()
		return createRowAccessPolicyWithRequest(t, createRowAccessPolicyBasicRequest(t))
	}

	t.Run("create row access policy: no optionals", func(t *testing.T) {
		request := createRowAccessPolicyBasicRequest(t)

		rowAccessPolicy := createRowAccessPolicyWithRequest(t, request)

		assertRowAccessPolicy(t, rowAccessPolicy, request.GetName(), "")
	})

	t.Run("create row access policy: full", func(t *testing.T) {
		request := createRowAccessPolicyBasicRequest(t)
		request.Comment = sdk.String("some comment")

		rowAccessPolicy := createRowAccessPolicyWithRequest(t, request)

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
			t.Cleanup(cleanupRowAccessPolicyProvider(id))
		} else {
			t.Cleanup(cleanupRowAccessPolicyProvider(newId))
		}
		require.NoError(t, err)

		_, err = client.RowAccessPolicies.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		rowAccessPolicy, err := client.RowAccessPolicies.ShowByID(ctx, newId)
		require.NoError(t, err)

		assertRowAccessPolicy(t, rowAccessPolicy, newId, "")
	})

	t.Run("alter row access policy: set and unset comment", func(t *testing.T) {
		rowAccessPolicy := createRowAccessPolicy(t)
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
		rowAccessPolicy := createRowAccessPolicy(t)
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

	t.Run("alter row access policy: set and unset tags", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		rowAccessPolicy := createRowAccessPolicy(t)
		id := rowAccessPolicy.ID()

		tagValue := "abc"
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		alterRequestSetTags := sdk.NewAlterRowAccessPolicyRequest(id).WithSetTags(tags)

		err := client.RowAccessPolicies.Alter(ctx, alterRequestSetTags)
		require.NoError(t, err)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeRowAccessPolicy)
		require.NoError(t, err)

		assert.Equal(t, tagValue, returnedTagValue)

		unsetTags := []sdk.ObjectIdentifier{
			tag.ID(),
		}
		alterRequestUnsetTags := sdk.NewAlterRowAccessPolicyRequest(id).WithUnsetTags(unsetTags)

		err = client.RowAccessPolicies.Alter(ctx, alterRequestUnsetTags)
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeRowAccessPolicy)
		require.Error(t, err)
	})

	t.Run("show row access policy: default", func(t *testing.T) {
		rowAccessPolicy1 := createRowAccessPolicy(t)
		rowAccessPolicy2 := createRowAccessPolicy(t)

		showRequest := sdk.NewShowRowAccessPolicyRequest()
		returnedRowAccessPolicies, err := client.RowAccessPolicies.Show(ctx, showRequest)
		require.NoError(t, err)
		require.LessOrEqual(t, 2, len(returnedRowAccessPolicies))
		assert.Contains(t, returnedRowAccessPolicies, *rowAccessPolicy1)
		assert.Contains(t, returnedRowAccessPolicies, *rowAccessPolicy2)
	})

	t.Run("show row access policy: with options", func(t *testing.T) {
		rowAccessPolicy1 := createRowAccessPolicy(t)
		rowAccessPolicy2 := createRowAccessPolicy(t)

		showRequest := sdk.NewShowRowAccessPolicyRequest().
			WithLike(&sdk.Like{Pattern: &rowAccessPolicy1.Name}).
			WithIn(&sdk.In{Schema: testClientHelper().Ids.SchemaId()})
		returnedRowAccessPolicies, err := client.RowAccessPolicies.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Equal(t, 1, len(returnedRowAccessPolicies))
		assert.Contains(t, returnedRowAccessPolicies, *rowAccessPolicy1)
		assert.NotContains(t, returnedRowAccessPolicies, *rowAccessPolicy2)
	})

	t.Run("describe row access policy: existing", func(t *testing.T) {
		argName := random.AlphaN(5)
		argType := sdk.DataTypeVARCHAR
		args := sdk.NewCreateRowAccessPolicyArgsRequest(argName, argType)
		body := "true"

		request := createRowAccessPolicyRequest(t, []sdk.CreateRowAccessPolicyArgsRequest{*args}, body)
		rowAccessPolicy := createRowAccessPolicyWithRequest(t, request)

		returnedRowAccessPolicyDescription, err := client.RowAccessPolicies.Describe(ctx, rowAccessPolicy.ID())
		require.NoError(t, err)

		assertRowAccessPolicyDescription(t, returnedRowAccessPolicyDescription, rowAccessPolicy.ID(), fmt.Sprintf("(%s %s)", strings.ToUpper(argName), argType), body)
	})

	t.Run("describe row access policy: with data type normalization", func(t *testing.T) {
		argName := random.AlphaN(5)
		argType := sdk.DataTypeTimestamp
		args := sdk.NewCreateRowAccessPolicyArgsRequest(argName, argType)
		body := "true"

		request := createRowAccessPolicyRequest(t, []sdk.CreateRowAccessPolicyArgsRequest{*args}, body)
		rowAccessPolicy := createRowAccessPolicyWithRequest(t, request)

		returnedRowAccessPolicyDescription, err := client.RowAccessPolicies.Describe(ctx, rowAccessPolicy.ID())
		require.NoError(t, err)

		assertRowAccessPolicyDescription(t, returnedRowAccessPolicyDescription, rowAccessPolicy.ID(), fmt.Sprintf("(%s %s)", strings.ToUpper(argName), sdk.DataTypeTimestampNTZ), body)
	})

	t.Run("describe row access policy: with data type normalization", func(t *testing.T) {
		argName := random.AlphaN(5)
		argType := sdk.DataType("VARCHAR(200)")
		args := sdk.NewCreateRowAccessPolicyArgsRequest(argName, argType)
		body := "true"

		request := createRowAccessPolicyRequest(t, []sdk.CreateRowAccessPolicyArgsRequest{*args}, body)
		rowAccessPolicy := createRowAccessPolicyWithRequest(t, request)

		returnedRowAccessPolicyDescription, err := client.RowAccessPolicies.Describe(ctx, rowAccessPolicy.ID())
		require.NoError(t, err)

		assertRowAccessPolicyDescription(t, returnedRowAccessPolicyDescription, rowAccessPolicy.ID(), fmt.Sprintf("(%s %s)", strings.ToUpper(argName), sdk.DataTypeVARCHAR), body)
	})

	t.Run("describe row access policy: non-existing", func(t *testing.T) {
		_, err := client.RowAccessPolicies.Describe(ctx, NonExistingSchemaObjectIdentifier)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_RowAccessPoliciesShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupRowAccessPolicyHandle := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.RowAccessPolicies.Drop(ctx, sdk.NewDropRowAccessPolicyRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createRowAccessPolicyHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		args := sdk.NewCreateRowAccessPolicyArgsRequest(random.AlphaN(5), sdk.DataTypeVARCHAR)
		err := client.RowAccessPolicies.Create(ctx, sdk.NewCreateRowAccessPolicyRequest(id, []sdk.CreateRowAccessPolicyArgsRequest{*args}, "true"))
		require.NoError(t, err)
		t.Cleanup(cleanupRowAccessPolicyHandle(id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		createRowAccessPolicyHandle(t, id1)
		createRowAccessPolicyHandle(t, id2)

		e1, err := client.RowAccessPolicies.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.RowAccessPolicies.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}

func TestInt_RowAccessPoliciesDescribe(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupRowAccessPolicyHandle := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.RowAccessPolicies.Drop(ctx, sdk.NewDropRowAccessPolicyRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createRowAccessPolicyHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier, args []sdk.CreateRowAccessPolicyArgsRequest) {
		t.Helper()

		err := client.RowAccessPolicies.Create(ctx, sdk.NewCreateRowAccessPolicyRequest(id, args, "true"))
		require.NoError(t, err)
		t.Cleanup(cleanupRowAccessPolicyHandle(id))
	}

	t.Run("describe", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		arg1 := sdk.NewCreateRowAccessPolicyArgsRequest(random.AlphaN(5), sdk.DataTypeVARCHAR)
		arg2 := sdk.NewCreateRowAccessPolicyArgsRequest(random.AlphaN(5)+" foo", "TEXT")
		arg3 := sdk.NewCreateRowAccessPolicyArgsRequest(random.AlphaN(5), "NUMBER(10, 5)")
		args := []sdk.CreateRowAccessPolicyArgsRequest{*arg1, *arg2, *arg3}

		createRowAccessPolicyHandle(t, id, args)

		policy, err := client.RowAccessPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "true", policy.Body)
		assert.Equal(t, id.Name(), policy.Name)
		assert.Equal(t, "BOOLEAN", policy.ReturnType)
		assert.Equal(t, fmt.Sprintf("(%s VARCHAR, %s VARCHAR, %s NUMBER)", arg1.Name, arg2.Name, arg3.Name), policy.Signature)
		gotArgs, err := policy.Arguments()
		require.NoError(t, err)
		assert.Equal(t, []sdk.RowAccessPolicyArgument{
			{
				Name: arg1.Name,
				Type: "VARCHAR",
			},
			{
				Name: arg2.Name,
				Type: "VARCHAR",
			},
			{
				Name: arg3.Name,
				Type: "NUMBER",
			},
		}, gotArgs)
	})
}

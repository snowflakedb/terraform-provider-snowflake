//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_JoinPolicies(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	const defaultBody = "JOIN_CONSTRAINT(JOIN_REQUIRED => TRUE)"
	const otherBody = "JOIN_CONSTRAINT(JOIN_REQUIRED => FALSE)"

	createBasic := func(t *testing.T) sdk.SchemaObjectIdentifier {
		t.Helper()
		id, cleanup := testClientHelper().JoinPolicy.Create(t)
		t.Cleanup(cleanup)
		return id
	}

	t.Run("create: basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.JoinPolicies.Create(ctx, sdk.NewCreateJoinPolicyRequest(id, defaultBody))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().JoinPolicy.DropFunc(t, id))

		assertThatObject(t, objectassert.JoinPolicy(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasKind("JOIN_POLICY").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasOwnerRoleType("ROLE").
			HasOptions(""),
		)
		assertThatObject(t, objectassert.JoinPolicyDetails(t, id).
			HasName(id.Name()).
			HasSignature("()").
			HasReturnType("JOIN_CONSTRAINT").
			HasBody(defaultBody),
		)
	})

	t.Run("create: complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()
		body := "CASE WHEN CURRENT_ROLE() = 'ACCOUNTADMIN' THEN JOIN_CONSTRAINT(JOIN_REQUIRED => FALSE) ELSE JOIN_CONSTRAINT(JOIN_REQUIRED => TRUE) END"
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		err := client.JoinPolicies.Create(ctx, sdk.NewCreateJoinPolicyRequest(id, body).
			WithIfNotExists(true).
			WithComment(comment).
			WithTag([]sdk.TagAssociation{{Name: tag.ID(), Value: "finance"}}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().JoinPolicy.DropFunc(t, id))

		assertThatObject(t, objectassert.JoinPolicy(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasKind("JOIN_POLICY").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasOwnerRoleType("ROLE").
			HasOptions(""),
		)
		assertThatObject(t, objectassert.JoinPolicyDetails(t, id).
			HasName(id.Name()).
			HasSignature("()").
			HasReturnType("JOIN_CONSTRAINT").
			HasBody(body),
		)
		assertTagSet(t, tag.ID(), id, sdk.ObjectTypeJoinPolicy, "finance")

		err = client.JoinPolicies.Create(ctx, sdk.NewCreateJoinPolicyRequest(id, otherBody).
			WithIfNotExists(true))
		require.NoError(t, err)
		assertThatObject(t, objectassert.JoinPolicyDetails(t, id).
			HasBody(body),
		)
	})

	t.Run("create: or replace", func(t *testing.T) {
		id := createBasic(t)

		comment := random.Comment()
		err := client.JoinPolicies.Create(ctx, sdk.NewCreateJoinPolicyRequest(id, otherBody).
			WithOrReplace(true).
			WithComment(comment))
		require.NoError(t, err)

		assertThatObject(t, objectassert.JoinPolicy(t, id).
			HasComment(comment),
		)
		assertThatObject(t, objectassert.JoinPolicyDetails(t, id).
			HasBody(otherBody),
		)
	})

	t.Run("show: like, in, starts with, and limit", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		prefix := random.AlphaN(6)
		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, schema.ID())
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, schema.ID())
		id3 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		for _, id := range []sdk.SchemaObjectIdentifier{id1, id2, id3} {
			_, cleanup := testClientHelper().JoinPolicy.CreateWithRequest(t, sdk.NewCreateJoinPolicyRequest(id, defaultBody))
			t.Cleanup(cleanup)
		}

		byLike, err := client.JoinPolicies.Show(ctx, sdk.NewShowJoinPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id1.Name())}).
			WithIn(sdk.In{Schema: schema.ID()}))
		require.NoError(t, err)
		require.Len(t, byLike, 1)
		assertThatObject(t, objectassert.JoinPolicyFromObject(t, &byLike[0]).
			HasName(id1.Name()).
			HasDatabaseName(id1.DatabaseName()).
			HasSchemaName(id1.SchemaName()).
			HasKind("JOIN_POLICY"),
		)

		byStartsWith, err := client.JoinPolicies.Show(ctx, sdk.NewShowJoinPolicyRequest().
			WithIn(sdk.In{Schema: schema.ID()}).
			WithStartsWith(prefix))
		require.NoError(t, err)
		require.Len(t, byStartsWith, 2)

		byLimit, err := client.JoinPolicies.Show(ctx, sdk.NewShowJoinPolicyRequest().
			WithIn(sdk.In{Schema: schema.ID()}).
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
		require.NoError(t, err)
		require.Len(t, byLimit, 1)

		inSchema, err := client.JoinPolicies.Show(ctx, sdk.NewShowJoinPolicyRequest().
			WithIn(sdk.In{Schema: schema.ID()}))
		require.NoError(t, err)
		require.Len(t, inSchema, 3)

		inDatabase, err := client.JoinPolicies.Show(ctx, sdk.NewShowJoinPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id3.Name())}).
			WithIn(sdk.In{Database: id3.DatabaseId()}))
		require.NoError(t, err)
		require.Len(t, inDatabase, 1)

		inAccount, err := client.JoinPolicies.Show(ctx, sdk.NewShowJoinPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id3.Name())}).
			WithIn(sdk.In{Account: sdk.Bool(true)}))
		require.NoError(t, err)
		require.Len(t, inAccount, 1)
	})

	t.Run("show by id: same name in two schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())
		_, cleanup1 := testClientHelper().JoinPolicy.CreateWithRequest(t, sdk.NewCreateJoinPolicyRequest(id1, defaultBody))
		t.Cleanup(cleanup1)
		_, cleanup2 := testClientHelper().JoinPolicy.CreateWithRequest(t, sdk.NewCreateJoinPolicyRequest(id2, defaultBody).WithComment("second"))
		t.Cleanup(cleanup2)

		assertThatObject(t, objectassert.JoinPolicy(t, id1).
			HasName(id1.Name()).
			HasSchemaName(id1.SchemaName()).
			HasComment(""),
		)
		assertThatObject(t, objectassert.JoinPolicy(t, id2).
			HasName(id2.Name()).
			HasSchemaName(id2.SchemaName()).
			HasComment("second"),
		)
	})

	t.Run("show by id safely: existing and non-existing", func(t *testing.T) {
		id := createBasic(t)

		policy, err := client.JoinPolicies.ShowByIDSafely(ctx, id)
		require.NoError(t, err)
		assertThatObject(t, objectassert.JoinPolicyFromObject(t, policy).
			HasName(id.Name()),
		)

		_, err = client.JoinPolicies.ShowByIDSafely(ctx, NonExistingSchemaObjectIdentifier)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("describe: existing", func(t *testing.T) {
		id := createBasic(t)

		assertThatObject(t, objectassert.JoinPolicyDetails(t, id).
			HasName(id.Name()).
			HasSignature("()").
			HasReturnType("JOIN_CONSTRAINT").
			HasBody(defaultBody),
		)
	})

	t.Run("describe: non-existing", func(t *testing.T) {
		_, err := client.JoinPolicies.Describe(ctx, NonExistingSchemaObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter: set body", func(t *testing.T) {
		id := createBasic(t)

		err := client.JoinPolicies.Alter(ctx, sdk.NewAlterJoinPolicyRequest(id).
			WithIfExists(true).
			WithSet(*sdk.NewJoinPolicySetRequest().WithSetBody(otherBody)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.JoinPolicyDetails(t, id).
			HasBody(otherBody),
		)
	})

	t.Run("alter: set and unset comment", func(t *testing.T) {
		id := createBasic(t)
		comment := random.Comment()

		err := client.JoinPolicies.Alter(ctx, sdk.NewAlterJoinPolicyRequest(id).
			WithSet(*sdk.NewJoinPolicySetRequest().WithComment(comment)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.JoinPolicy(t, id).
			HasComment(comment),
		)

		err = client.JoinPolicies.Alter(ctx, sdk.NewAlterJoinPolicyRequest(id).
			WithUnset(*sdk.NewJoinPolicyUnsetRequest().WithComment(true)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.JoinPolicy(t, id).
			HasComment(""),
		)
	})

	t.Run("alter: rename", func(t *testing.T) {
		id := createBasic(t)
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.JoinPolicies.Alter(ctx, sdk.NewAlterJoinPolicyRequest(id).WithRenameTo(newId))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().JoinPolicy.DropFunc(t, newId))

		_, err = client.JoinPolicies.ShowByID(ctx, id)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)

		assertThatObject(t, objectassert.JoinPolicy(t, newId).
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()).
			HasKind("JOIN_POLICY"),
		)
	})

	t.Run("alter: rename to another schema", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)
		id := createBasic(t)
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

		err := client.JoinPolicies.Alter(ctx, sdk.NewAlterJoinPolicyRequest(id).WithRenameTo(newId))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().JoinPolicy.DropFunc(t, newId))

		assertThatObject(t, objectassert.JoinPolicy(t, newId).
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()),
		)
	})

	t.Run("drop: existing", func(t *testing.T) {
		id := createBasic(t)

		err := client.JoinPolicies.Drop(ctx, sdk.NewDropJoinPolicyRequest(id))
		require.NoError(t, err)
		_, err = client.JoinPolicies.ShowByID(ctx, id)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop: non-existing", func(t *testing.T) {
		err := client.JoinPolicies.Drop(ctx, sdk.NewDropJoinPolicyRequest(NonExistingSchemaObjectIdentifier))
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: if exists on non-existing", func(t *testing.T) {
		err := client.JoinPolicies.Drop(ctx, sdk.NewDropJoinPolicyRequest(NonExistingSchemaObjectIdentifier).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("drop safely: existing", func(t *testing.T) {
		id := createBasic(t)

		err := client.JoinPolicies.DropSafely(ctx, id)
		require.NoError(t, err)
		_, err = client.JoinPolicies.ShowByIDSafely(ctx, id)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop safely: non-existing", func(t *testing.T) {
		err := client.JoinPolicies.DropSafely(ctx, NonExistingSchemaObjectIdentifier)
		require.NoError(t, err)
	})
}

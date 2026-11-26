//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_AggregationPolicies(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	const (
		defaultBody = "AGGREGATION_CONSTRAINT(MIN_GROUP_SIZE => 5)"
		otherBody   = "AGGREGATION_CONSTRAINT(MIN_GROUP_SIZE => 2)"
		caseBody    = "CASE WHEN CURRENT_ROLE() = 'ADMIN' THEN NO_AGGREGATION_CONSTRAINT() ELSE AGGREGATION_CONSTRAINT(MIN_GROUP_SIZE => 3) END"
	)

	createBasic := func(t *testing.T) sdk.SchemaObjectIdentifier {
		t.Helper()
		id, cleanup := testClientHelper().AggregationPolicy.CreateAggregationPolicy(t)
		t.Cleanup(cleanup)
		return id
	}

	createInSchema := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		_, cleanup := testClientHelper().AggregationPolicy.CreateWithRequest(t, sdk.NewCreateAggregationPolicyRequest(id, defaultBody))
		t.Cleanup(cleanup)
	}

	t.Run("create: basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.AggregationPolicies.Create(ctx, sdk.NewCreateAggregationPolicyRequest(id, defaultBody))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().AggregationPolicy.DropAggregationPolicyFunc(t, id))

		assertThatObject(t, objectassert.AggregationPolicy(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasKind("AGGREGATION_POLICY").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasOwnerRoleType("ROLE").
			HasOptions(""),
		)
		assertThatObject(t, objectassert.AggregationPolicyDetails(t, id).
			HasName(id.Name()).
			HasSignature("()").
			HasReturnType("AGGREGATION_CONSTRAINT").
			HasBody(defaultBody),
		)
	})

	t.Run("create: complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		err := client.AggregationPolicies.Create(ctx, sdk.NewCreateAggregationPolicyRequest(id, caseBody).
			WithIfNotExists(true).
			WithComment(comment).
			WithTag([]sdk.TagAssociation{{Name: tag.ID(), Value: "finance"}}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().AggregationPolicy.DropAggregationPolicyFunc(t, id))

		assertThatObject(t, objectassert.AggregationPolicy(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasKind("AGGREGATION_POLICY").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasOwnerRoleType("ROLE").
			HasOptions(""),
		)
		assertThatObject(t, objectassert.AggregationPolicyDetails(t, id).
			HasName(id.Name()).
			HasSignature("()").
			HasReturnType("AGGREGATION_CONSTRAINT").
			HasBody(caseBody),
		)
		assertTagSet(t, tag.ID(), id, sdk.ObjectTypeAggregationPolicy, "finance")

		err = client.AggregationPolicies.Create(ctx, sdk.NewCreateAggregationPolicyRequest(id, otherBody).
			WithIfNotExists(true))
		require.NoError(t, err)
		assertThatObject(t, objectassert.AggregationPolicyDetails(t, id).
			HasBody(caseBody),
		)
	})

	t.Run("create: or replace", func(t *testing.T) {
		id := createBasic(t)

		comment := random.Comment()
		err := client.AggregationPolicies.Create(ctx, sdk.NewCreateAggregationPolicyRequest(id, otherBody).
			WithOrReplace(true).
			WithComment(comment))
		require.NoError(t, err)

		assertThatObject(t, objectassert.AggregationPolicy(t, id).
			HasComment(comment),
		)
		assertThatObject(t, objectassert.AggregationPolicyDetails(t, id).
			HasBody(otherBody),
		)
	})

	t.Run("alter: set body", func(t *testing.T) {
		id := createBasic(t)

		err := client.AggregationPolicies.Alter(ctx, sdk.NewAlterAggregationPolicyRequest(id).
			WithSet(*sdk.NewAggregationPolicySetRequest().WithSetBody(otherBody)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.AggregationPolicyDetails(t, id).
			HasBody(otherBody),
		)
	})

	t.Run("alter: set and unset comment", func(t *testing.T) {
		id := createBasic(t)
		comment := random.Comment()

		err := client.AggregationPolicies.Alter(ctx, sdk.NewAlterAggregationPolicyRequest(id).
			WithSet(*sdk.NewAggregationPolicySetRequest().WithComment(comment)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.AggregationPolicy(t, id).
			HasComment(comment),
		)

		err = client.AggregationPolicies.Alter(ctx, sdk.NewAlterAggregationPolicyRequest(id).
			WithUnset(*sdk.NewAggregationPolicyUnsetRequest().WithComment(true)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.AggregationPolicy(t, id).
			HasComment(""),
		)
	})

	t.Run("alter: rename", func(t *testing.T) {
		id := createBasic(t)
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.AggregationPolicies.Alter(ctx, sdk.NewAlterAggregationPolicyRequest(id).WithRenameTo(newId))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().AggregationPolicy.DropAggregationPolicyFunc(t, newId))

		_, err = client.AggregationPolicies.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)

		assertThatObject(t, objectassert.AggregationPolicy(t, newId).
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()).
			HasKind("AGGREGATION_POLICY"),
		)
	})

	t.Run("drop: existing", func(t *testing.T) {
		id := createBasic(t)

		err := client.AggregationPolicies.Drop(ctx, sdk.NewDropAggregationPolicyRequest(id))
		require.NoError(t, err)
		_, err = client.AggregationPolicies.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop: non-existing", func(t *testing.T) {
		err := client.AggregationPolicies.Drop(ctx, sdk.NewDropAggregationPolicyRequest(NonExistingSchemaObjectIdentifier))
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: if exists on non-existing", func(t *testing.T) {
		err := client.AggregationPolicies.Drop(ctx, sdk.NewDropAggregationPolicyRequest(NonExistingSchemaObjectIdentifier).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("drop safely: existing", func(t *testing.T) {
		id := createBasic(t)

		err := client.AggregationPolicies.DropSafely(ctx, id)
		require.NoError(t, err)
		_, err = client.AggregationPolicies.ShowByIDSafely(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop safely: non-existing", func(t *testing.T) {
		err := client.AggregationPolicies.DropSafely(ctx, NonExistingSchemaObjectIdentifier)
		require.NoError(t, err)
	})

	t.Run("show: like, in, starts with, and limit", func(t *testing.T) {
		db, dbCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(dbCleanup)
		schema, schemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
		t.Cleanup(schemaCleanup)
		otherSchema, otherSchemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
		t.Cleanup(otherSchemaCleanup)

		prefix := random.AlphaN(6)
		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, schema.ID())
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, schema.ID())
		id3 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(otherSchema.ID())
		createInSchema(t, id1)
		createInSchema(t, id2)
		createInSchema(t, id3)

		byLike, err := client.AggregationPolicies.Show(ctx, sdk.NewShowAggregationPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id1.Name())}).
			WithIn(sdk.In{Schema: schema.ID()}))
		require.NoError(t, err)
		require.Len(t, byLike, 1)
		assertThatObject(t, objectassert.AggregationPolicyFromObject(t, &byLike[0]).
			HasName(id1.Name()).
			HasDatabaseName(id1.DatabaseName()).
			HasSchemaName(id1.SchemaName()).
			HasKind("AGGREGATION_POLICY"),
		)

		inSchema, err := client.AggregationPolicies.Show(ctx, sdk.NewShowAggregationPolicyRequest().
			WithIn(sdk.In{Schema: schema.ID()}))
		require.NoError(t, err)
		require.Len(t, inSchema, 2)

		inDatabase, err := client.AggregationPolicies.Show(ctx, sdk.NewShowAggregationPolicyRequest().
			WithIn(sdk.In{Database: db.ID()}))
		require.NoError(t, err)
		require.Len(t, inDatabase, 3)

		inAccount, err := client.AggregationPolicies.Show(ctx, sdk.NewShowAggregationPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id3.Name())}).
			WithIn(sdk.In{Account: sdk.Bool(true)}))
		require.NoError(t, err)
		require.Len(t, inAccount, 1)

		startsWith, err := client.AggregationPolicies.Show(ctx, sdk.NewShowAggregationPolicyRequest().
			WithIn(sdk.In{Database: db.ID()}).
			WithStartsWith(prefix))
		require.NoError(t, err)
		require.Len(t, startsWith, 2)

		limited, err := client.AggregationPolicies.Show(ctx, sdk.NewShowAggregationPolicyRequest().
			WithIn(sdk.In{Schema: schema.ID()}).
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
		require.NoError(t, err)
		require.Len(t, limited, 1)
	})

	t.Run("show by id: same name in different schemas", func(t *testing.T) {
		otherSchema, otherSchemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(otherSchemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), otherSchema.ID())
		createInSchema(t, id1)
		createInSchema(t, id2)

		assertThatObject(t, objectassert.AggregationPolicy(t, id1).
			HasName(id1.Name()).
			HasSchemaName(id1.SchemaName()),
		)
		assertThatObject(t, objectassert.AggregationPolicy(t, id2).
			HasName(id2.Name()).
			HasSchemaName(id2.SchemaName()),
		)
	})

	t.Run("show by id safely: existing", func(t *testing.T) {
		id := createBasic(t)

		policy, err := client.AggregationPolicies.ShowByIDSafely(ctx, id)
		require.NoError(t, err)
		assertThatObject(t, objectassert.AggregationPolicyFromObject(t, policy).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()),
		)
	})

	t.Run("show by id safely: non-existing", func(t *testing.T) {
		_, err := client.AggregationPolicies.ShowByIDSafely(ctx, NonExistingSchemaObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("describe: existing", func(t *testing.T) {
		id := createBasic(t)

		assertThatObject(t, objectassert.AggregationPolicyDetails(t, id).
			HasName(id.Name()).
			HasSignature("()").
			HasReturnType("AGGREGATION_CONSTRAINT").
			HasBody(defaultBody),
		)
	})

	t.Run("describe: non-existing", func(t *testing.T) {
		_, err := client.AggregationPolicies.Describe(ctx, NonExistingSchemaObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

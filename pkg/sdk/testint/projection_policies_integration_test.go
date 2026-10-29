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

func TestInt_ProjectionPolicies(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	const defaultBody = "PROJECTION_CONSTRAINT(ALLOW => false)"

	createBasic := func(t *testing.T) sdk.SchemaObjectIdentifier {
		t.Helper()
		id, cleanup := testClientHelper().ProjectionPolicy.CreateProjectionPolicy(t)
		t.Cleanup(cleanup)
		return id
	}

	t.Run("create: basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.ProjectionPolicies.Create(ctx, sdk.NewCreateProjectionPolicyRequest(id, defaultBody))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ProjectionPolicy.DropProjectionPolicyFunc(t, id))

		assertThatObject(t, objectassert.ProjectionPolicy(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasKind("PROJECTION_POLICY").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasOwnerRoleType("ROLE").
			HasOptions(""),
		)
		assertThatObject(t, objectassert.ProjectionPolicyDetails(t, id).
			HasName(id.Name()).
			HasSignature("()").
			HasReturnType("PROJECTION_CONSTRAINT").
			HasBody(defaultBody),
		)
	})

	t.Run("create: complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		err := client.ProjectionPolicies.Create(ctx, sdk.NewCreateProjectionPolicyRequest(id, defaultBody).
			WithIfNotExists(true).
			WithComment(comment).
			WithTag([]sdk.TagAssociation{{Name: tag.ID(), Value: "finance"}}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ProjectionPolicy.DropProjectionPolicyFunc(t, id))

		assertThatObject(t, objectassert.ProjectionPolicy(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasKind("PROJECTION_POLICY").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasOwnerRoleType("ROLE").
			HasOptions(""),
		)
		assertThatObject(t, objectassert.ProjectionPolicyDetails(t, id).
			HasName(id.Name()).
			HasSignature("()").
			HasReturnType("PROJECTION_CONSTRAINT").
			HasBody(defaultBody),
		)
		assertTagSet(t, tag.ID(), id, sdk.ObjectTypeProjectionPolicy, "finance")

		err = client.ProjectionPolicies.Create(ctx, sdk.NewCreateProjectionPolicyRequest(id, "PROJECTION_CONSTRAINT(ALLOW => true)").
			WithIfNotExists(true))
		require.NoError(t, err)
		assertThatObject(t, objectassert.ProjectionPolicyDetails(t, id).
			HasBody(defaultBody),
		)
	})

	t.Run("create: or replace", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.ProjectionPolicies.Create(ctx, sdk.NewCreateProjectionPolicyRequest(id, defaultBody))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ProjectionPolicy.DropProjectionPolicyFunc(t, id))

		replacement := "PROJECTION_CONSTRAINT(ALLOW => true)"
		comment := random.Comment()
		err = client.ProjectionPolicies.Create(ctx, sdk.NewCreateProjectionPolicyRequest(id, replacement).
			WithOrReplace(true).
			WithComment(comment))
		require.NoError(t, err)

		assertThatObject(t, objectassert.ProjectionPolicy(t, id).
			HasComment(comment),
		)
		assertThatObject(t, objectassert.ProjectionPolicyDetails(t, id).
			HasBody(replacement),
		)
	})

	t.Run("show: like, in, starts with, and limit", func(t *testing.T) {
		id := createBasic(t)

		shown, err := client.ProjectionPolicies.Show(ctx, sdk.NewShowProjectionPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id.Name())}).
			WithIn(sdk.In{Schema: id.SchemaId()}).
			WithStartsWith(id.Name()).
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
		require.NoError(t, err)
		require.Len(t, shown, 1)
		assertThatObject(t, objectassert.ProjectionPolicyFromObject(t, &shown[0]).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasKind("PROJECTION_POLICY"),
		)

		inDatabase, err := client.ProjectionPolicies.Show(ctx, sdk.NewShowProjectionPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id.Name())}).
			WithIn(sdk.In{Database: id.DatabaseId()}))
		require.NoError(t, err)
		require.Len(t, inDatabase, 1)

		inAccount, err := client.ProjectionPolicies.Show(ctx, sdk.NewShowProjectionPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id.Name())}).
			WithIn(sdk.In{Account: sdk.Bool(true)}))
		require.NoError(t, err)
		require.Len(t, inAccount, 1)
	})

	t.Run("describe: existing", func(t *testing.T) {
		id := createBasic(t)

		assertThatObject(t, objectassert.ProjectionPolicyDetails(t, id).
			HasName(id.Name()).
			HasSignature("()").
			HasReturnType("PROJECTION_CONSTRAINT").
			HasBody(defaultBody),
		)
	})

	t.Run("describe: non-existing", func(t *testing.T) {
		_, err := client.ProjectionPolicies.Describe(ctx, NonExistingSchemaObjectIdentifier)
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter: set body", func(t *testing.T) {
		id := createBasic(t)
		body := "PROJECTION_CONSTRAINT(ALLOW => true)"

		err := client.ProjectionPolicies.Alter(ctx, sdk.NewAlterProjectionPolicyRequest(id).
			WithIfExists(true).
			WithSet(*sdk.NewProjectionPolicySetRequest().WithSetBody(body)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.ProjectionPolicyDetails(t, id).
			HasBody(body),
		)
	})

	t.Run("alter: set and unset comment", func(t *testing.T) {
		id := createBasic(t)
		comment := random.Comment()

		err := client.ProjectionPolicies.Alter(ctx, sdk.NewAlterProjectionPolicyRequest(id).
			WithSet(*sdk.NewProjectionPolicySetRequest().WithComment(comment)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.ProjectionPolicy(t, id).
			HasComment(comment),
		)

		err = client.ProjectionPolicies.Alter(ctx, sdk.NewAlterProjectionPolicyRequest(id).
			WithUnset(*sdk.NewProjectionPolicyUnsetRequest().WithComment(true)))
		require.NoError(t, err)
		assertThatObject(t, objectassert.ProjectionPolicy(t, id).
			HasComment(""),
		)
	})

	t.Run("alter: rename", func(t *testing.T) {
		id := createBasic(t)
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.ProjectionPolicies.Alter(ctx, sdk.NewAlterProjectionPolicyRequest(id).WithRenameTo(newId))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ProjectionPolicy.DropProjectionPolicyFunc(t, newId))

		_, err = client.ProjectionPolicies.ShowByID(ctx, id)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)

		assertThatObject(t, objectassert.ProjectionPolicy(t, newId).
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()).
			HasKind("PROJECTION_POLICY"),
		)
	})

	t.Run("drop: if exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.ProjectionPolicies.Drop(ctx, sdk.NewDropProjectionPolicyRequest(id).WithIfExists(true))
		require.NoError(t, err)

		created := createBasic(t)
		err = client.ProjectionPolicies.Drop(ctx, sdk.NewDropProjectionPolicyRequest(created))
		require.NoError(t, err)
		_, err = client.ProjectionPolicies.ShowByID(ctx, created)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})
}

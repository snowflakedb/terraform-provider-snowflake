//go:build non_account_level_tests

package testint

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_SessionPolicies(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	assertSessionPolicy := func(t *testing.T, sessionPolicy *sdk.SessionPolicy, id sdk.SchemaObjectIdentifier, expectedComment string) {
		t.Helper()
		assert.NotEmpty(t, sessionPolicy.CreatedOn)
		assert.Equal(t, id.Name(), sessionPolicy.Name)
		assert.Equal(t, id.SchemaName(), sessionPolicy.SchemaName)
		assert.Equal(t, id.DatabaseName(), sessionPolicy.DatabaseName)
		assert.Equal(t, "ACCOUNTADMIN", sessionPolicy.Owner)
		assert.Equal(t, expectedComment, sessionPolicy.Comment)
		assert.Equal(t, "SESSION_POLICY", sessionPolicy.Kind)
		assert.Equal(t, "ROLE", sessionPolicy.OwnerRoleType)
		assert.Equal(t, "", sessionPolicy.Options)
	}

	cleanupSessionPolicyProvider := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.SessionPolicies.Drop(ctx, sdk.NewDropSessionPolicyRequest(id).WithIfExists(true))
			require.NoError(t, err)
		}
	}

	createSessionPolicy := func(t *testing.T) *sdk.SessionPolicy {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.SessionPolicies.Create(ctx, sdk.NewCreateSessionPolicyRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupSessionPolicyProvider(id))

		sessionPolicy, err := client.SessionPolicies.ShowByID(ctx, id)
		require.NoError(t, err)

		return sessionPolicy
	}

	t.Run("create session_policy: complete case", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()

		role1, role1Cleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(role1Cleanup)

		role2, role2Cleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(role2Cleanup)

		role3, role3Cleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(role3Cleanup)

		request := sdk.NewCreateSessionPolicyRequest(id).
			WithSessionIdleTimeoutMins(5).
			WithSessionUiIdleTimeoutMins(34).
			WithAllowedSecondaryRoles(*sdk.NewSessionPolicySecondaryRolesRequest().
				WithRoles([]sdk.AccountObjectIdentifier{role1.ID(), role2.ID()})).
			WithBlockedSecondaryRoles(*sdk.NewSessionPolicySecondaryRolesRequest().
				WithRoles([]sdk.AccountObjectIdentifier{role3.ID()})).
			WithComment(comment).
			WithIfNotExists(true)

		err := client.SessionPolicies.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupSessionPolicyProvider(id))

		sessionPolicy, err := client.SessionPolicies.ShowByID(ctx, id)

		require.NoError(t, err)
		assertSessionPolicy(t, sessionPolicy, id, comment)
		assertThatObject(t, objectassert.SessionPolicyDetails(t, id).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(comment).
			HasSessionIdleTimeoutMins(5).
			HasSessionUiIdleTimeoutMins(34).
			HasAllowedSecondaryRolesUnordered(role1.Name, role2.Name).
			HasBlockedSecondaryRoles(role3.Name))
	})

	t.Run("create session_policy: no optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		request := sdk.NewCreateSessionPolicyRequest(id)

		err := client.SessionPolicies.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupSessionPolicyProvider(id))

		sessionPolicy, err := client.SessionPolicies.ShowByID(ctx, id)

		require.NoError(t, err)
		assertSessionPolicy(t, sessionPolicy, id, "")
		assertThatObject(t, objectassert.SessionPolicyDetails(t, id).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("").
			HasSessionIdleTimeoutMins(240).
			HasSessionUiIdleTimeoutMins(1080).
			HasAllowedSecondaryRolesUnordered("ALL").
			HasBlockedSecondaryRoles())
	})

	t.Run("drop session_policy: existing", func(t *testing.T) {
		id := createSessionPolicy(t).ID()

		err := client.SessionPolicies.Drop(ctx, sdk.NewDropSessionPolicyRequest(id))
		require.NoError(t, err)

		_, err = client.SessionPolicies.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop session_policy: non-existing", func(t *testing.T) {
		err := client.SessionPolicies.Drop(ctx, sdk.NewDropSessionPolicyRequest(NonExistingSchemaObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter session_policy: set value and unset value", func(t *testing.T) {
		id := createSessionPolicy(t).ID()

		alterRequest := sdk.NewAlterSessionPolicyRequest(id).WithSet(*sdk.NewSessionPolicySetRequest().
			WithSessionIdleTimeoutMins(60).
			WithSessionUiIdleTimeoutMins(60).
			WithAllowedSecondaryRoles(*sdk.NewSessionPolicySecondaryRolesRequest().WithNone(true)).
			WithBlockedSecondaryRoles(*sdk.NewSessionPolicySecondaryRolesRequest().WithAll(true)).
			WithComment("new comment"),
		)
		err := client.SessionPolicies.Alter(ctx, alterRequest)
		require.NoError(t, err)

		assertThatObject(t, objectassert.SessionPolicyDetails(t, id).
			HasComment("new comment").
			HasSessionIdleTimeoutMins(60).
			HasSessionUiIdleTimeoutMins(60).
			HasAllowedSecondaryRolesUnordered().
			HasBlockedSecondaryRoles("ALL"))

		alterRequest = sdk.NewAlterSessionPolicyRequest(id).WithUnset(*sdk.NewSessionPolicyUnsetRequest().
			WithSessionIdleTimeoutMins(true).
			WithSessionUiIdleTimeoutMins(true).
			WithAllowedSecondaryRoles(true).
			WithBlockedSecondaryRoles(true).
			WithComment(true),
		)
		err = client.SessionPolicies.Alter(ctx, alterRequest)
		require.NoError(t, err)

		assertThatObject(t, objectassert.SessionPolicyDetails(t, id).
			HasComment("").
			HasSessionIdleTimeoutMins(240).
			HasSessionUiIdleTimeoutMins(1080).
			HasAllowedSecondaryRolesUnordered("ALL").
			HasBlockedSecondaryRoles())
	})

	t.Run("alter session_policy: rename", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.SessionPolicies.Create(ctx, sdk.NewCreateSessionPolicyRequest(id))
		require.NoError(t, err)

		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		alterRequest := sdk.NewAlterSessionPolicyRequest(id).WithRenameTo(newId)

		err = client.SessionPolicies.Alter(ctx, alterRequest)
		if err != nil {
			t.Cleanup(cleanupSessionPolicyProvider(id))
		} else {
			t.Cleanup(cleanupSessionPolicyProvider(newId))
		}
		require.NoError(t, err)

		_, err = client.SessionPolicies.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		sessionPolicy, err := client.SessionPolicies.ShowByID(ctx, newId)
		require.NoError(t, err)

		assertSessionPolicy(t, sessionPolicy, newId, "")
	})

	t.Run("show session policies", func(t *testing.T) {
		db, dbCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(dbCleanup)

		user, userCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(userCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("test_session_policyzzz")
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("test_session_policy_2_")
		id3 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("test_session_policy_3_")
		id4 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(sdk.NewDatabaseObjectIdentifier(db.Name, "PUBLIC"))
		ids := []sdk.SchemaObjectIdentifier{id1, id2, id3, id4}
		for _, id := range ids {
			_, sessionPolicyCleanup := testClientHelper().SessionPolicy.CreateSessionPolicyWithOptions(t, id, sdk.NewCreateSessionPolicyRequest(id))
			t.Cleanup(sessionPolicyCleanup)
		}

		testClientHelper().User.Alter(t, user.ID(), &sdk.AlterUserOptions{Set: &sdk.UserSet{SessionPolicy: sdk.Pointer(id1)}})
		userSessionPolicyAttachmentCleanup := func() {
			testClientHelper().User.Alter(t, user.ID(), &sdk.AlterUserOptions{Unset: &sdk.UserUnset{SessionPolicy: sdk.Bool(true)}})
		}
		t.Cleanup(userSessionPolicyAttachmentCleanup)

		t.Run("like", func(t *testing.T) {
			sessionPolicies, err := client.SessionPolicies.Show(ctx, sdk.NewShowSessionPolicyRequest().
				WithLike(sdk.Like{Pattern: sdk.String("test_session_policy_2_%")}).
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id1.SchemaId()}}))
			require.NoError(t, err)
			assert.Len(t, sessionPolicies, 1)
		})

		t.Run("starts_with", func(t *testing.T) {
			sessionPolicies, err := client.SessionPolicies.Show(ctx, sdk.NewShowSessionPolicyRequest().
				WithStartsWith("test_session_policy_").
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id1.SchemaId()}}))
			require.NoError(t, err)
			assert.Len(t, sessionPolicies, 2)
		})

		t.Run("in_account", func(t *testing.T) {
			sessionPolicies, err := client.SessionPolicies.Show(ctx, sdk.NewShowSessionPolicyRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Account: sdk.Bool(true)}}))
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(sessionPolicies), 4)
		})

		t.Run("in_database", func(t *testing.T) {
			sessionPolicies, err := client.SessionPolicies.Show(ctx, sdk.NewShowSessionPolicyRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Database: id1.DatabaseId()}}))
			require.NoError(t, err)
			assert.Len(t, sessionPolicies, 3)
		})

		t.Run("in_schema", func(t *testing.T) {
			sessionPolicies, err := client.SessionPolicies.Show(ctx, sdk.NewShowSessionPolicyRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id1.SchemaId()}}))
			require.NoError(t, err)
			assert.Len(t, sessionPolicies, 3)
		})

		t.Run("on_account", func(t *testing.T) {
			sessionPolicies, err := client.SessionPolicies.Show(ctx, sdk.NewShowSessionPolicyRequest().WithOn(sdk.On{Account: sdk.Pointer(true)}))
			require.NoError(t, err)
			assert.Len(t, sessionPolicies, 0)
		})

		t.Run("on_user", func(t *testing.T) {
			sessionPolicies, err := client.SessionPolicies.Show(ctx, sdk.NewShowSessionPolicyRequest().WithOn(sdk.On{User: user.ID()}))
			require.NoError(t, err)
			assert.Len(t, sessionPolicies, 1)
		})

		t.Run("limit", func(t *testing.T) {
			sessionPolicies, err := client.SessionPolicies.Show(ctx, sdk.NewShowSessionPolicyRequest().
				WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}).
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id1.SchemaId()}}))
			require.NoError(t, err)
			assert.Len(t, sessionPolicies, 1)
		})

		t.Run("limit from", func(t *testing.T) {
			sessionPolicies, err := client.SessionPolicies.Show(ctx, sdk.NewShowSessionPolicyRequest().
				WithLimit(sdk.LimitFrom{Rows: sdk.Int(1), From: sdk.String("test_session_policy_")}).
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id1.SchemaId()}}))
			require.NoError(t, err)
			require.Len(t, sessionPolicies, 1)
			require.True(t, strings.HasPrefix(sessionPolicies[0].Name, "test_session_policy_2"))
		})
	})

	t.Run("describe session_policy", func(t *testing.T) {
		sessionPolicy := createSessionPolicy(t)

		details, err := client.SessionPolicies.DescribeDetails(ctx, sessionPolicy.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.SessionPolicyDetailsFromObject(t, details).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment("").
			HasSessionIdleTimeoutMins(240).
			HasSessionUiIdleTimeoutMins(1080).
			HasAllowedSecondaryRolesUnordered("ALL").
			HasBlockedSecondaryRoles())
	})

	t.Run("describe session policy: non-existing", func(t *testing.T) {
		id := NonExistingSchemaObjectIdentifier

		_, err := client.SessionPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

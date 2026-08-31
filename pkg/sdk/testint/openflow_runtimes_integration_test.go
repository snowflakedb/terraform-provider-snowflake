//go:build non_account_level_tests

package testint

import (
	"sync"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Runtimes require an ACTIVE deployment, named by TEST_SF_TF_OPENFLOW_DEPLOYMENT. A runtime takes four to
// five minutes to create, so subtests share one unless they mutate it.
func TestInt_OpenflowRuntimes(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	client := testClient(t)
	ctx := testContext(t)

	deploymentId := testClientHelper().OpenflowDeployment.ActiveDeploymentForRuntimes(t)

	currentRole := testClientHelper().Context.CurrentRole(t)
	executeAsRoleId := currentRole

	// Read-only fixture for the subtests that only need a runtime to exist: creating one takes four to five
	// minutes, so describe and the show cases share it rather than each paying that. Anything that mutates a
	// runtime creates its own. Teardown goes against the parent t so the closure outlives the subtest that
	// triggered it.
	parentT := t
	var (
		sharedRuntimeOnce sync.Once
		sharedRuntimeId   sdk.SchemaObjectIdentifier
		sharedRuntimeErr  error
	)
	sharedRuntime := func(t *testing.T) sdk.SchemaObjectIdentifier {
		t.Helper()
		sharedRuntimeOnce.Do(func() {
			// Asserted against the calling t: FailNow on a t this goroutine does not own exits silently.
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			if err := client.OpenflowRuntimes.Create(ctx, sdk.NewCreateOpenflowRuntimeRequest(id, deploymentId, sdk.OpenflowRuntimeNodeTypeSmall, 1, 1, executeAsRoleId)); err != nil {
				sharedRuntimeErr = err
				return
			}
			parentT.Cleanup(testClientHelper().OpenflowRuntime.DropFunc(parentT, id))
			testClientHelper().OpenflowRuntime.WaitForStatus(t, id, sdk.OpenflowRuntimeStatusActive, helpers.OpenflowRuntimeActiveTimeout)
			sharedRuntimeId = id
		})
		require.NoError(t, sharedRuntimeErr, "creating the shared runtime failed")
		require.NotEmpty(t, sharedRuntimeId.Name(), "the shared runtime is unavailable because creating it failed")
		return sharedRuntimeId
	}

	t.Run("create: basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateOpenflowRuntimeRequest(id, deploymentId, sdk.OpenflowRuntimeNodeTypeSmall, 1, 1, executeAsRoleId)

		err := client.OpenflowRuntimes.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().OpenflowRuntime.DropFunc(t, id))

		// Settle first so the status is deterministic. Teardown waits for the same thing before it can
		// terminate, so this mostly moves the wait earlier.
		testClientHelper().OpenflowRuntime.WaitForStatus(t, id, sdk.OpenflowRuntimeStatusActive, helpers.OpenflowRuntimeActiveTimeout)

		runtime, err := client.OpenflowRuntimes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.OpenflowRuntimeFromObject(t, runtime).
				HasName(id.Name()).
				HasStatus(sdk.OpenflowRuntimeStatusActive).
				HasDeployment(deploymentId.Name()).
				HasMinNodes(1).
				HasMaxNodes(1).
				HasNodeType(sdk.OpenflowRuntimeNodeTypeSmall).
				HasNoDisplayName().
				HasNoExternalAccessIntegrations().
				HasInitiallySuspended(false).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasExecuteAsRole(executeAsRoleId.Name()).
				HasKeyNotEmpty().
				HasOwner(currentRole.Name()).
				HasNoComment().
				HasCreatedOnNotEmpty().
				HasUpdatedOnNotEmpty(),
		)
	})

	t.Run("create: complete", func(t *testing.T) {
		networkRule, networkRuleCleanup := testClientHelper().NetworkRule.CreateEgressWithIdentifier(t, testClientHelper().Ids.RandomSchemaObjectIdentifier())
		t.Cleanup(networkRuleCleanup)
		externalAccessIntegration, eaiCleanup := testClientHelper().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
		t.Cleanup(eaiCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()
		displayName := random.AlphaN(12)

		request := sdk.NewCreateOpenflowRuntimeRequest(id, deploymentId, sdk.OpenflowRuntimeNodeTypeMedium, 1, 2, executeAsRoleId).
			WithIfNotExists(true).
			WithExternalAccessIntegrations(*sdk.NewOpenflowRuntimeExternalAccessIntegrationsRequest([]sdk.AccountObjectIdentifier{externalAccessIntegration})).
			WithDisplayName(displayName).
			WithComment(comment)

		err := client.OpenflowRuntimes.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().OpenflowRuntime.DropFunc(t, id))

		testClientHelper().OpenflowRuntime.WaitForStatus(t, id, sdk.OpenflowRuntimeStatusActive, helpers.OpenflowRuntimeActiveTimeout)

		runtime, err := client.OpenflowRuntimes.ShowByID(ctx, id)
		require.NoError(t, err)

		// external_access_integrations is a JSON array, so a correct conversion yields one element per
		// integration rather than a single bracketed string.
		assertThatObject(
			t, objectassert.OpenflowRuntimeFromObject(t, runtime).
				HasName(id.Name()).
				HasStatus(sdk.OpenflowRuntimeStatusActive).
				HasDeployment(deploymentId.Name()).
				HasMinNodes(1).
				HasMaxNodes(2).
				HasNodeType(sdk.OpenflowRuntimeNodeTypeMedium).
				HasDisplayName(displayName).
				HasExternalAccessIntegrations(externalAccessIntegration).
				HasInitiallySuspended(false).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasExecuteAsRole(executeAsRoleId.Name()).
				HasKeyNotEmpty().
				HasOwner(currentRole.Name()).
				HasComment(comment).
				HasCreatedOnNotEmpty().
				HasUpdatedOnNotEmpty(),
		)
	})

	t.Run("create: min nodes greater than max nodes is rejected", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		// Torn down regardless, in case the statement unexpectedly succeeds.
		t.Cleanup(testClientHelper().OpenflowRuntime.DropFunc(t, id))
		err := client.OpenflowRuntimes.Create(ctx, sdk.NewCreateOpenflowRuntimeRequest(id, deploymentId, sdk.OpenflowRuntimeNodeTypeSmall, 5, 2, executeAsRoleId))
		require.Error(t, err)
	})

	// Drives every field SET and UNSET accept. It runs on a runtime of its own rather than the shared one,
	// because changing node counts and the role would break the later subtests that assert the values the
	// shared runtime is created with.
	t.Run("alter: set and unset", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.OpenflowRuntimes.Create(ctx, sdk.NewCreateOpenflowRuntimeRequest(id, deploymentId, sdk.OpenflowRuntimeNodeTypeSmall, 1, 1, executeAsRoleId))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().OpenflowRuntime.DropFunc(t, id))
		testClientHelper().OpenflowRuntime.WaitForStatus(t, id, sdk.OpenflowRuntimeStatusActive, helpers.OpenflowRuntimeActiveTimeout)

		networkRule, networkRuleCleanup := testClientHelper().NetworkRule.CreateEgressWithIdentifier(t, testClientHelper().Ids.RandomSchemaObjectIdentifier())
		t.Cleanup(networkRuleCleanup)
		externalAccessIntegration, eaiCleanup := testClientHelper().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
		t.Cleanup(eaiCleanup)

		comment := random.Comment()
		displayName := random.AlphaN(12)
		// EXECUTE_AS_ROLE is set on its own further down. The server checks that the execute-as role can use
		// the runtime's integrations, so pairing a different role with these in one statement is refused.
		err = client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithSet(
			*sdk.NewOpenflowRuntimeSetRequest().
				WithDisplayName(displayName).
				WithMinNodes(1).
				WithMaxNodes(3).
				WithExternalAccessIntegrations(*sdk.NewOpenflowRuntimeExternalAccessIntegrationsRequest([]sdk.AccountObjectIdentifier{externalAccessIntegration})).
				WithComment(comment),
		))
		require.NoError(t, err)
		// Changing node counts drives the runtime through UPDATING before it settles back on ACTIVE.
		testClientHelper().OpenflowRuntime.WaitForStatus(t, id, sdk.OpenflowRuntimeStatusActive, helpers.OpenflowRuntimeActiveTimeout)

		assertThatObject(
			t, objectassert.OpenflowRuntime(t, id).
				HasDisplayName(displayName).
				HasMinNodes(1).
				HasMaxNodes(3).
				HasExternalAccessIntegrations(externalAccessIntegration).
				HasComment(comment),
		)

		// MIN_NODES and MAX_NODES have no UNSET; they always hold a value. EXECUTE_AS_ROLE is asserted to be
		// refused below rather than unset here, since one rejected property fails the whole statement.
		err = client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithUnset(
			*sdk.NewOpenflowRuntimeUnsetRequest().
				WithExternalAccessIntegrations(true).
				WithDisplayName(true).
				WithComment(true),
		))
		require.NoError(t, err)
		testClientHelper().OpenflowRuntime.WaitForStatus(t, id, sdk.OpenflowRuntimeStatusActive, helpers.OpenflowRuntimeActiveTimeout)

		assertThatObject(
			t, objectassert.OpenflowRuntime(t, id).
				HasNoDisplayName().
				HasNoComment().
				HasNoExternalAccessIntegrations(),
		)

		// The server rejects UNSET EXECUTE_AS_ROLE outright: a runtime always runs as some role. The SDK does
		// not model it, so raw SQL is the only way to reach it, as for a connector's DEFAULT_VERSION.
		_, err = client.ExecForTests(ctx, "ALTER OPENFLOW RUNTIME "+id.FullyQualifiedName()+" UNSET EXECUTE_AS_ROLE")
		require.ErrorContains(t, err, "cannot unset property")

		// The role can still be moved with SET. PUBLIC rather than a purpose-made role, since the test role
		// cannot CREATE ROLE, and it only works once the integrations are detached: the server checks that
		// the execute-as role can use whatever the runtime carries.
		err = client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithSet(
			*sdk.NewOpenflowRuntimeSetRequest().WithExecuteAsRole(snowflakeroles.Public),
		))
		require.NoError(t, err)
		testClientHelper().OpenflowRuntime.WaitForStatus(t, id, sdk.OpenflowRuntimeStatusActive, helpers.OpenflowRuntimeActiveTimeout)

		assertThatObject(t, objectassert.OpenflowRuntime(t, id).HasExecuteAsRole(snowflakeroles.Public.Name()))
	})

	t.Run("alter: rename", func(t *testing.T) {
		runtime, cleanup := testClientHelper().OpenflowRuntime.Create(t, deploymentId, executeAsRoleId)
		id := runtime.ID()
		// Explicitly a new name within the runtime's own schema, rather than relying on the default generator
		// resolving to the same one.
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(id.SchemaId())

		err := client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithRenameTo(newId))
		if err != nil {
			t.Cleanup(cleanup)
			require.NoError(t, err)
		}
		t.Cleanup(testClientHelper().OpenflowRuntime.DropFunc(t, newId))

		_, err = client.OpenflowRuntimes.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		assertThatObject(t, objectassert.OpenflowRuntime(t, newId).HasName(newId.Name()))
	})

	// Its own runtime, since suspending mutates it and the shared one is read-only.
	t.Run("alter: suspend and resume", func(t *testing.T) {
		runtime, cleanup := testClientHelper().OpenflowRuntime.Create(t, deploymentId, executeAsRoleId)
		t.Cleanup(cleanup)
		id := runtime.ID()
		testClientHelper().OpenflowRuntime.WaitForStatus(t, id, sdk.OpenflowRuntimeStatusActive, helpers.OpenflowRuntimeActiveTimeout)

		err := client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithSuspend(true))
		require.NoError(t, err)
		testClientHelper().OpenflowRuntime.WaitForStatus(t, id, sdk.OpenflowRuntimeStatusSuspended, helpers.OpenflowRuntimeSuspendTimeout)

		err = client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithResume(true))
		require.NoError(t, err)
		testClientHelper().OpenflowRuntime.WaitForStatus(t, id, sdk.OpenflowRuntimeStatusActive, helpers.OpenflowRuntimeActiveTimeout)
	})

	t.Run("describe", func(t *testing.T) {
		id := sharedRuntime(t)

		// DESCRIBE returns server_url and node_type_tier, which SHOW does not, and omits database_name,
		// schema_name, created_on and updated_on, which SHOW does return. Those two are Snowflake-assigned, so
		// they are asserted as populated rather than pinned; an unmapped column would leave them nil.
		//
		// The shared runtime is created without a display name, comment or integrations, and nothing mutates it.
		assertThatObject(
			t, objectassert.OpenflowRuntimeDetails(t, id).
				HasName(id.Name()).
				HasStatus(sdk.OpenflowRuntimeStatusActive).
				HasDeployment(deploymentId.Name()).
				HasMinNodes(1).
				HasMaxNodes(1).
				HasNodeType(sdk.OpenflowRuntimeNodeTypeSmall).
				HasNoDisplayName().
				HasNoExternalAccessIntegrations().
				HasInitiallySuspended(false).
				HasExecuteAsRole(executeAsRoleId.Name()).
				HasKeyNotEmpty().
				HasOwner(currentRole.Name()).
				HasNoComment().
				HasServerUrlNotEmpty().
				HasNodeTypeTierNotEmpty().
				// Id is threaded in by the caller, since DESCRIBE returns neither database_name nor schema_name.
				HasId(id),
		)
	})

	t.Run("show: with like", func(t *testing.T) {
		id := sharedRuntime(t)

		runtimes, err := client.OpenflowRuntimes.Show(ctx, sdk.NewShowOpenflowRuntimeRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id.Name())}))
		require.NoError(t, err)
		require.Len(t, runtimes, 1)
		assert.Equal(t, id.Name(), runtimes[0].Name)
	})

	t.Run("show: in account", func(t *testing.T) {
		id := sharedRuntime(t)

		// Without IN ACCOUNT, SHOW is scoped to the session's database and schema.
		runtimes, err := client.OpenflowRuntimes.Show(ctx, sdk.NewShowOpenflowRuntimeRequest().
			WithIn(sdk.In{Account: sdk.Bool(true)}))
		require.NoError(t, err)
		assert.Contains(t, collections.Map(runtimes, func(r sdk.OpenflowRuntime) string { return r.Name }), id.Name())
	})

	t.Run("show: in database and in schema", func(t *testing.T) {
		id := sharedRuntime(t)

		inDatabase, err := client.OpenflowRuntimes.Show(ctx, sdk.NewShowOpenflowRuntimeRequest().
			WithIn(sdk.In{Database: id.DatabaseId()}))
		require.NoError(t, err)
		assert.Contains(t, collections.Map(inDatabase, func(r sdk.OpenflowRuntime) string { return r.Name }), id.Name())

		inSchema, err := client.OpenflowRuntimes.Show(ctx, sdk.NewShowOpenflowRuntimeRequest().
			WithIn(sdk.In{Schema: id.SchemaId()}))
		require.NoError(t, err)
		assert.Contains(t, collections.Map(inSchema, func(r sdk.OpenflowRuntime) string { return r.Name }), id.Name())
	})

	// Both clauses together, in the order Snowflake requires. STARTS WITH on a unique name returns one row on
	// its own, so LIMIT is covered by "show: with limit" below.
	t.Run("show: with starts with and limit", func(t *testing.T) {
		id := sharedRuntime(t)

		runtimes, err := client.OpenflowRuntimes.Show(ctx, sdk.NewShowOpenflowRuntimeRequest().
			WithIn(sdk.In{Account: sdk.Bool(true)}).
			WithStartsWith(id.Name()).
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
		require.NoError(t, err)
		require.Len(t, runtimes, 1)
		assert.Equal(t, id.Name(), runtimes[0].Name)
	})

	// Asserts LIMIT actually limits. Creating a second runtime to guarantee the account holds more than one
	// costs minutes, so this compares against whatever an unlimited SHOW returns instead of skipping.
	t.Run("show: with limit", func(t *testing.T) {
		sharedRuntime(t)

		all, err := client.OpenflowRuntimes.Show(ctx, sdk.NewShowOpenflowRuntimeRequest().
			WithIn(sdk.In{Account: sdk.Bool(true)}))
		require.NoError(t, err)
		require.NotEmpty(t, all)

		limited, err := client.OpenflowRuntimes.Show(ctx, sdk.NewShowOpenflowRuntimeRequest().
			WithIn(sdk.In{Account: sdk.Bool(true)}).
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
		require.NoError(t, err)
		assert.Len(t, limited, 1)
		assert.LessOrEqual(t, len(limited), len(all), "LIMIT 1 returned more rows than an unlimited SHOW")
	})

	// ADD and REMOVE edit the list in place, where SET replaces it wholesale. The list runs
	// none -> one -> two -> one -> none, so later subtests find the shared runtime as they expect.
	// Its own runtime, since attaching integrations mutates it and the shared one is read-only.
	t.Run("alter: add and remove external access integrations", func(t *testing.T) {
		runtime, cleanup := testClientHelper().OpenflowRuntime.Create(t, deploymentId, executeAsRoleId)
		t.Cleanup(cleanup)
		id := runtime.ID()
		testClientHelper().OpenflowRuntime.WaitForStatus(t, id, sdk.OpenflowRuntimeStatusActive, helpers.OpenflowRuntimeActiveTimeout)

		firstRule, firstRuleCleanup := testClientHelper().NetworkRule.CreateEgressWithIdentifier(t, testClientHelper().Ids.RandomSchemaObjectIdentifier())
		t.Cleanup(firstRuleCleanup)
		firstEai, firstEaiCleanup := testClientHelper().ExternalAccessIntegration.CreateExternalAccessIntegration(t, firstRule.ID())
		t.Cleanup(firstEaiCleanup)
		secondRule, secondRuleCleanup := testClientHelper().NetworkRule.CreateEgressWithIdentifier(t, testClientHelper().Ids.RandomSchemaObjectIdentifier())
		t.Cleanup(secondRuleCleanup)
		secondEai, secondEaiCleanup := testClientHelper().ExternalAccessIntegration.CreateExternalAccessIntegration(t, secondRule.ID())
		t.Cleanup(secondEaiCleanup)

		// Detach both, or the cleanups above cannot drop integrations still referenced. After a green run
		// they are already detached, which Snowflake rejects, so tolerate only that error.
		t.Cleanup(func() {
			for _, eai := range []sdk.AccountObjectIdentifier{secondEai, firstEai} {
				err := client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).
					WithRemoveExternalAccessIntegrations(*sdk.NewOpenflowRuntimeExternalAccessIntegrationsRequest([]sdk.AccountObjectIdentifier{eai})))
				if err != nil {
					assert.ErrorContains(t, err, "001422 (22023): SQL compilation error:\ninvalid value")
					continue
				}
				testClientHelper().OpenflowRuntime.WaitUntilSettled(t, id, helpers.OpenflowRuntimeActiveTimeout)
			}
		})

		integrationNames := func() []string {
			t.Helper()
			runtime, err := client.OpenflowRuntimes.ShowByID(ctx, id)
			require.NoError(t, err)
			return collections.Map(runtime.ExternalAccessIntegrations, func(i sdk.AccountObjectIdentifier) string { return i.Name() })
		}
		// Editing the list drives the runtime through UPDATING; a second ALTER is refused until it settles.
		add := func(eais ...sdk.AccountObjectIdentifier) {
			t.Helper()
			require.NoError(t, client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).
				WithAddExternalAccessIntegrations(*sdk.NewOpenflowRuntimeExternalAccessIntegrationsRequest(eais))))
			testClientHelper().OpenflowRuntime.WaitForStatus(t, id, sdk.OpenflowRuntimeStatusActive, helpers.OpenflowRuntimeActiveTimeout)
		}
		remove := func(eais ...sdk.AccountObjectIdentifier) {
			t.Helper()
			require.NoError(t, client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).
				WithRemoveExternalAccessIntegrations(*sdk.NewOpenflowRuntimeExternalAccessIntegrationsRequest(eais))))
			testClientHelper().OpenflowRuntime.WaitForStatus(t, id, sdk.OpenflowRuntimeStatusActive, helpers.OpenflowRuntimeActiveTimeout)
		}

		// Two integrations at once, because the column is a comma-separated list and one element would not
		// exercise the split. Whether ADD appends rather than replaces is Snowflake's own behavior; what
		// belongs here is that the list converts, and that an emptied one converts to nothing rather than to
		// an integration named "null".
		require.Empty(t, integrationNames(), "the shared runtime should start with no external access integrations")

		add(firstEai, secondEai)
		assert.ElementsMatch(t, []string{firstEai.Name(), secondEai.Name()}, integrationNames())

		remove(firstEai, secondEai)
		assert.Empty(t, integrationNames())
	})

	t.Run("drop: requires terminate first", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.OpenflowRuntimes.Create(ctx, sdk.NewCreateOpenflowRuntimeRequest(id, deploymentId, sdk.OpenflowRuntimeNodeTypeSmall, 1, 1, executeAsRoleId))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().OpenflowRuntime.DropFunc(t, id))

		testClientHelper().OpenflowRuntime.WaitForStatus(t, id, sdk.OpenflowRuntimeStatusActive, helpers.OpenflowRuntimeActiveTimeout)

		// DROP is only permitted from DELETED, so it is refused while the runtime is ACTIVE.
		err = client.OpenflowRuntimes.Drop(ctx, sdk.NewDropOpenflowRuntimeRequest(id))
		require.ErrorContains(t, err, "DROP not allowed")

		// TERMINATE is accepted straight from ACTIVE, so the SUSPEND below covers that path rather than being
		// a precondition.
		err = client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithSuspend(true))
		require.NoError(t, err)
		testClientHelper().OpenflowRuntime.WaitForStatus(t, id, sdk.OpenflowRuntimeStatusSuspended, helpers.OpenflowRuntimeSuspendTimeout)

		err = client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithTerminate(true))
		require.NoError(t, err)
		testClientHelper().OpenflowRuntime.WaitForStatus(t, id, sdk.OpenflowRuntimeStatusDeleted, helpers.OpenflowRuntimeTerminateTimeout)

		err = client.OpenflowRuntimes.Drop(ctx, sdk.NewDropOpenflowRuntimeRequest(id))
		require.NoError(t, err)

		_, err = client.OpenflowRuntimes.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop: when a runtime does not exist", func(t *testing.T) {
		err := client.OpenflowRuntimes.Drop(ctx, sdk.NewDropOpenflowRuntimeRequest(NonExistingSchemaObjectIdentifier))
		assert.Error(t, err)
	})

	t.Run("drop: if exists on a non-existing runtime", func(t *testing.T) {
		err := client.OpenflowRuntimes.Drop(ctx, sdk.NewDropOpenflowRuntimeRequest(NonExistingSchemaObjectIdentifier).WithIfExists(true))
		require.NoError(t, err)
	})
}

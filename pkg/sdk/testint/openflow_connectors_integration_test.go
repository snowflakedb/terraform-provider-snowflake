//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Connectors require an ACTIVE runtime, so these tests share one. They use the OPENFLOW_POSTGRES_CDC
// definition. Its max_node_count is 1, so the runtime stays at one node.
//
// A new connector is a draft and settles on STOPPED, not RUNNING.
func TestInt_OpenflowConnectors(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	client := testClient(t)
	ctx := testContext(t)

	deploymentId := testClientHelper().OpenflowDeployment.ActiveDeploymentForRuntimes(t)
	definition := testClientHelper().OpenflowConnectorDefinition.ForTesting(t)

	currentRole := testClientHelper().Context.CurrentRole(t)
	runtime, runtimeCleanup := testClientHelper().OpenflowRuntime.Create(t, deploymentId, currentRole)
	t.Cleanup(runtimeCleanup)
	runtimeId := runtime.ID()
	testClientHelper().OpenflowRuntime.WaitForStatus(t, runtimeId, sdk.OpenflowRuntimeStatusActive, helpers.OpenflowRuntimeActiveTimeout)

	t.Run("create: from definition", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(runtimeId.SchemaId())
		request := sdk.NewCreateOpenflowConnectorRequest(id, runtimeId).WithFromDefinition(definition.Name)

		connector, cleanup := testClientHelper().OpenflowConnector.CreateWithRequest(t, request)
		t.Cleanup(cleanup)

		assertThatObject(
			t, objectassert.OpenflowConnectorFromObject(t, connector).
				HasName(id.Name()).
				HasStatus(sdk.OpenflowConnectorStatusStopped).
				HasRuntime(runtimeId.Name()).
				HasConnectorDefinition(definition.Name).
				HasNoDisplayName().
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner(currentRole.Name()).
				HasNoComment().
				HasCreatedOnNotEmpty().
				HasUpdatedOnNotEmpty().
				HasLiveVersionLocationUriNotEmpty().
				// default_version says which saved version the connector runs, and defaults to LAST. The
				// other default_version_* columns describe that version, and stay empty until something is
				// committed.
				HasDefaultVersion("LAST").
				HasNoDefaultVersionName().
				HasNoDefaultVersionAlias().
				HasNoDefaultVersionLocationUri().
				HasNoDefaultVersionSourceLocationUri().
				HasConnectorUrlNotEmpty(),
		)
	})

	t.Run("create: complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(runtimeId.SchemaId())
		comment := random.Comment()
		displayName := random.AlphaN(12)

		request := sdk.NewCreateOpenflowConnectorRequest(id, runtimeId).
			WithIfNotExists(true).
			WithFromDefinition(definition.Name).
			WithDisplayName(displayName).
			WithComment(comment)

		connector, cleanup := testClientHelper().OpenflowConnector.CreateWithRequest(t, request)
		t.Cleanup(cleanup)

		assertThatObject(
			t, objectassert.OpenflowConnectorFromObject(t, connector).
				HasName(id.Name()).
				HasStatus(sdk.OpenflowConnectorStatusStopped).
				HasRuntime(runtimeId.Name()).
				HasConnectorDefinition(definition.Name).
				HasDisplayName(displayName).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner(currentRole.Name()).
				HasDefaultVersion("LAST").
				HasNoDefaultVersionName().
				HasNoDefaultVersionAlias().
				HasNoDefaultVersionLocationUri().
				HasNoDefaultVersionSourceLocationUri().
				HasLiveVersionLocationUriNotEmpty().
				HasComment(comment).
				HasCreatedOnNotEmpty().
				HasUpdatedOnNotEmpty().
				HasConnectorUrlNotEmpty(),
		)
	})

	// Checks the version URI the SDK builds against the one Snowflake reports for the same version. The test
	// schema's identifiers are mixed case, so this covers the per-part quoting.
	t.Run("show versions: the URI matches the one Snowflake reports", func(t *testing.T) {
		// Committed, because only a committed version has a name; a draft's live version does not.
		source, sourceCleanup := testClientHelper().OpenflowConnector.CreateFromDefinitionAndCommit(t, runtimeId, definition.Name)
		t.Cleanup(sourceCleanup)
		sourceId := source.ID()

		versions, err := client.OpenflowConnectors.ShowVersions(ctx, sdk.NewShowVersionsOpenflowConnectorRequest(sourceId))
		require.NoError(t, err)
		committed, err := collections.FindFirst(versions, func(v sdk.OpenflowConnectorVersion) bool {
			return v.Name != nil && *v.Name != ""
		})
		require.NoError(t, err, "committing should leave a named version behind")

		assert.Equal(t, committed.LocationUri,
			sdk.NewOpenflowConnectorVersionLocation(sourceId, *committed.Name).ToSql())
	})

	t.Run("alter: set and unset", func(t *testing.T) {
		// Committed, because DEFAULT_VERSION can only name a committed version and a draft has none. FIRST and
		// LAST resolve against committed versions too, so no form of it works before the commit.
		connector, cleanup := testClientHelper().OpenflowConnector.CreateFromDefinitionAndCommit(t, runtimeId, definition.Name)
		t.Cleanup(cleanup)
		id := connector.ID()

		require.NotNil(t, connector.DefaultVersionName)
		committedVersion := *connector.DefaultVersionName

		comment := random.Comment()
		displayName := random.AlphaN(12)
		err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithSet(
			*sdk.NewOpenflowConnectorSetRequest().
				WithComment(comment).
				WithDisplayName(displayName).
				WithDefaultVersion(*sdk.NewOpenflowConnectorDefaultVersionRequest().WithVersion(committedVersion)),
		))
		require.NoError(t, err)

		// A draft serves LAST, so pinning a version is an observable change.
		assertThatObject(
			t, objectassert.OpenflowConnector(t, id).
				HasComment(comment).
				HasDisplayName(displayName).
				HasDefaultVersion(committedVersion),
		)

		// UNSET has no DEFAULT_VERSION: the property always holds a value and the server refuses to unset
		// it, which "alter: unset default version is not supported" covers.
		err = client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithUnset(
			*sdk.NewOpenflowConnectorUnsetRequest().WithComment(true).WithDisplayName(true),
		))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.OpenflowConnector(t, id).
				HasNoComment().
				HasNoDisplayName().
				HasDefaultVersion(committedVersion),
		)
	})

	t.Run("alter: rename", func(t *testing.T) {
		connector, cleanup := testClientHelper().OpenflowConnector.CreateFromDefinition(t, runtimeId, definition.Name)
		id := connector.ID()
		// Registered before the assertion below, which can fail. Teardown looks the connector up first, so it
		// is a no-op once the rename has succeeded.
		t.Cleanup(cleanup)

		// A connector cannot be renamed across schemas.
		otherSchemaId := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), "SOME_OTHER_SCHEMA", id.Name()+"_MOVED")
		err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithRenameTo(otherSchemaId))
		require.Error(t, err)

		newId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(id.SchemaId())
		err = client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithRenameTo(newId))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().OpenflowConnector.DropFunc(t, newId))

		_, err = client.OpenflowConnectors.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
		assertThatObject(t, objectassert.OpenflowConnector(t, newId).HasName(newId.Name()))
	})

	t.Run("alter: add live version and commit", func(t *testing.T) {
		connector, cleanup := testClientHelper().OpenflowConnector.CreateFromDefinition(t, runtimeId, definition.Name)
		t.Cleanup(cleanup)
		id := connector.ID()

		// COMMIT promotes the existing live version to the default. Configuring the files first is the
		// snowflake_execute flow the provider documents.
		err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).
			WithCommit(*sdk.NewOpenflowConnectorCommitRequest()))
		require.NoError(t, err)

		// Committing fills in the columns describing the resolved version, which the draft case only sees
		// empty. default_version is not asserted here, since it already holds LAST before the commit;
		// "alter: set and unset" pins it to a version.
		assertThatObject(
			t, objectassert.OpenflowConnectorDetails(t, id).
				HasDefaultVersionNameNotEmpty().
				HasDefaultVersionLocationUriNotEmpty().
				HasLastVersionNameNotEmpty(),
		)

		// Editing a committed connector starts by adding a fresh live version from the last one.
		err = client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).
			WithAddLiveVersion(*sdk.NewOpenflowConnectorAddLiveVersionRequest().WithVersionAlias("v2").WithComment(random.Comment())))
		require.NoError(t, err)

		assertThatObject(t, objectassert.OpenflowConnectorDetails(t, id).HasLiveVersionLocationUriNotEmpty())

		// ABORT discards the uncommitted live version.
		err = client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithAbort(true))
		require.NoError(t, err)
	})

	t.Run("alter: start and stop", func(t *testing.T) {
		// Committed, because a connector must have a committed default version before START is even accepted.
		connector, cleanup := testClientHelper().OpenflowConnector.CreateFromDefinitionAndCommit(t, runtimeId, definition.Name)
		t.Cleanup(cleanup)
		id := connector.ID()

		assertThatObject(t, objectassert.OpenflowConnector(t, id).HasDefaultVersionNameNotEmpty())

		// The committed version is unconfigured, so START is accepted but settles on START_FAILED. RUNNING
		// needs uploaded config files, which is a GET/PUT driver operation the provider cannot express.
		err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithStart(true))
		require.NoError(t, err)
		testClientHelper().OpenflowConnector.WaitForAnyStatus(t, id, []sdk.OpenflowConnectorStatus{
			sdk.OpenflowConnectorStatusRunning,
			sdk.OpenflowConnectorStatusStartFailed,
		}, helpers.OpenflowConnectorRunningTimeout)

		// STOP is accepted from either outcome and returns the connector to STOPPED.
		err = client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithStop(true))
		require.NoError(t, err)
		testClientHelper().OpenflowConnector.WaitForAnyStatus(t, id, []sdk.OpenflowConnectorStatus{
			sdk.OpenflowConnectorStatusStopped,
		}, helpers.OpenflowConnectorStopTimeout)
	})

	t.Run("describe", func(t *testing.T) {
		connector, cleanup := testClientHelper().OpenflowConnector.CreateFromDefinition(t, runtimeId, definition.Name)
		t.Cleanup(cleanup)
		id := connector.ID()

		// DESCRIBE carries last_version_* and the git commit hashes, and omits database_name, schema_name,
		// created_on and updated_on. Every version column is empty on a draft, default_version excepted.
		// "alter: add live version and commit" asserts the populated case, since an unmapped column reads the
		// same as an absent value.
		assertThatObject(
			t, objectassert.OpenflowConnectorDetails(t, id).
				HasName(id.Name()).
				HasRuntime(runtimeId.Name()).
				HasConnectorDefinition(definition.Name).
				HasNoDisplayName().
				HasOwner(currentRole.Name()).
				HasNoComment().
				HasStatus(sdk.OpenflowConnectorStatusStopped).
				// Id is threaded in by the caller, not returned by DESCRIBE.
				HasId(id).
				HasDefaultVersion("LAST").
				HasNoDefaultVersionName().
				HasNoDefaultVersionAlias().
				HasNoDefaultVersionLocationUri().
				HasNoDefaultVersionSourceLocationUri().
				HasNoDefaultVersionGitCommitHash().
				HasNoLastVersionName().
				HasNoLastVersionAlias().
				HasNoLastVersionLocationUri().
				HasNoLastVersionSourceLocationUri().
				HasNoLastVersionGitCommitHash().
				HasLiveVersionLocationUriNotEmpty().
				HasConnectorUrlNotEmpty(),
		)
	})

	t.Run("show: with like", func(t *testing.T) {
		connector, cleanup := testClientHelper().OpenflowConnector.CreateFromDefinition(t, runtimeId, definition.Name)
		t.Cleanup(cleanup)

		connectors, err := client.OpenflowConnectors.Show(ctx, sdk.NewShowOpenflowConnectorRequest().
			WithLike(sdk.Like{Pattern: sdk.String(connector.Name)}))
		require.NoError(t, err)
		require.Len(t, connectors, 1)
		assert.Equal(t, connector.Name, connectors[0].Name)
	})

	t.Run("show: in account, database and schema", func(t *testing.T) {
		connector, cleanup := testClientHelper().OpenflowConnector.CreateFromDefinition(t, runtimeId, definition.Name)
		t.Cleanup(cleanup)
		id := connector.ID()

		names := func(cs []sdk.OpenflowConnector) []string {
			return collections.Map(cs, func(c sdk.OpenflowConnector) string { return c.Name })
		}

		inAccount, err := client.OpenflowConnectors.Show(ctx, sdk.NewShowOpenflowConnectorRequest().
			WithIn(sdk.In{Account: sdk.Bool(true)}))
		require.NoError(t, err)
		assert.Contains(t, names(inAccount), connector.Name)

		inDatabase, err := client.OpenflowConnectors.Show(ctx, sdk.NewShowOpenflowConnectorRequest().
			WithIn(sdk.In{Database: id.DatabaseId()}))
		require.NoError(t, err)
		assert.Contains(t, names(inDatabase), connector.Name)

		inSchema, err := client.OpenflowConnectors.Show(ctx, sdk.NewShowOpenflowConnectorRequest().
			WithIn(sdk.In{Schema: id.SchemaId()}))
		require.NoError(t, err)
		assert.Contains(t, names(inSchema), connector.Name)
	})

	t.Run("show: with starts with and limit", func(t *testing.T) {
		connector, cleanup := testClientHelper().OpenflowConnector.CreateFromDefinition(t, runtimeId, definition.Name)
		t.Cleanup(cleanup)

		// Both clauses together, in the order Snowflake requires. STARTS WITH on a unique name returns one row
		// on its own, so this says nothing about LIMIT.
		connectors, err := client.OpenflowConnectors.Show(ctx, sdk.NewShowOpenflowConnectorRequest().
			WithIn(sdk.In{Account: sdk.Bool(true)}).
			WithStartsWith(connector.Name).
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
		require.NoError(t, err)
		require.Len(t, connectors, 1)
		assert.Equal(t, connector.Name, connectors[0].Name)

		// LIMIT is asserted separately, account-wide, and needs more than one connector to be observable.
		all, err := client.OpenflowConnectors.Show(ctx, sdk.NewShowOpenflowConnectorRequest().
			WithIn(sdk.In{Account: sdk.Bool(true)}))
		require.NoError(t, err)
		if len(all) < 2 {
			t.Skipf("LIMIT cannot be observed with %d connector(s) on the account", len(all))
		}

		limited, err := client.OpenflowConnectors.Show(ctx, sdk.NewShowOpenflowConnectorRequest().
			WithIn(sdk.In{Account: sdk.Bool(true)}).
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
		require.NoError(t, err)
		assert.Len(t, limited, 1, "LIMIT 1 should return exactly one row where an unlimited SHOW returns %d", len(all))
	})

	// SHOW VERSIONS has its own row shape, so every field is asserted. A new connector has one version,
	// which is at once the live, default, first and last.
	t.Run("show versions: a new connector has one live default version", func(t *testing.T) {
		connector, cleanup := testClientHelper().OpenflowConnector.CreateFromDefinition(t, runtimeId, definition.Name)
		t.Cleanup(cleanup)

		versions, err := client.OpenflowConnectors.ShowVersions(ctx, sdk.NewShowVersionsOpenflowConnectorRequest(connector.ID()))
		require.NoError(t, err)
		require.Len(t, versions, 1)

		// The only version of a new connector is its live one, which is uncommitted and so has no name. An
		// alias, source_location_uri and git_commit_hash only apply to versions from a stage or git remote.
		// location_uri is the snow:// URI the config lives at, asserted exactly against the server's value in
		// "show versions: the URI matches the one Snowflake reports".
		//
		// The comment is Snowflake's own, "Initial version" at the time of writing, so it is asserted as
		// populated rather than pinned to wording the provider does not control.
		version := versions[0]
		assertThatObject(
			t, objectassert.OpenflowConnectorVersionFromObject(t, &version).
				HasNoName().
				HasNoAlias().
				HasCommentNotEmpty().
				HasCreatedOnNotEmpty().
				HasIsDefault(true).
				HasIsFirst(true).
				HasIsLast(true).
				HasIsLive(true).
				HasLocationUriNotEmpty().
				HasNoSourceLocationUri().
				HasNoGitCommitHash(),
		)
	})

	// LIMIT is the only filter SHOW VERSIONS accepts, and needs two versions to be observable. COMMIT then
	// ADD LIVE VERSION arranges that, so the count is asserted rather than skipped on.
	t.Run("show versions: limit", func(t *testing.T) {
		connector, cleanup := testClientHelper().OpenflowConnector.CreateFromDefinitionAndCommit(t, runtimeId, definition.Name)
		t.Cleanup(cleanup)
		id := connector.ID()

		err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).
			WithAddLiveVersion(*sdk.NewOpenflowConnectorAddLiveVersionRequest().WithVersionAlias("v2")))
		require.NoError(t, err)

		all, err := client.OpenflowConnectors.ShowVersions(ctx, sdk.NewShowVersionsOpenflowConnectorRequest(id))
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(all), 2, "a committed version plus a new live version should be two rows")

		limited, err := client.OpenflowConnectors.ShowVersions(ctx, sdk.NewShowVersionsOpenflowConnectorRequest(id).WithLimit(1))
		require.NoError(t, err)
		assert.Len(t, limited, 1, "LIMIT 1 should return exactly one row where an unlimited SHOW VERSIONS returns %d", len(all))

		// Discard the uncommitted version so teardown sees the usual shape.
		require.NoError(t, client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithAbort(true)))
	})

	t.Run("show versions: on a non-existing connector", func(t *testing.T) {
		_, err := client.OpenflowConnectors.ShowVersions(ctx, sdk.NewShowVersionsOpenflowConnectorRequest(NonExistingSchemaObjectIdentifier))
		assert.ErrorContains(t, err, "does not exist")
	})

	// UNSET DEFAULT_VERSION parses but Snowflake always refuses it, so the SDK models only SET.
	t.Run("alter: unset default version is not supported", func(t *testing.T) {
		connector, cleanup := testClientHelper().OpenflowConnector.CreateFromDefinition(t, runtimeId, definition.Name)
		t.Cleanup(cleanup)

		_, err := client.ExecForTests(ctx, "ALTER OPENFLOW CONNECTOR "+connector.ID().FullyQualifiedName()+" UNSET DEFAULT_VERSION")
		require.ErrorContains(t, err, "cannot unset property")
	})

	t.Run("drop: requires terminate first", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(runtimeId.SchemaId())
		err := client.OpenflowConnectors.Create(ctx, sdk.NewCreateOpenflowConnectorRequest(id, runtimeId).WithFromDefinition(definition.Name))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().OpenflowConnector.DropFunc(t, id))

		// Snowflake refuses DROP and TERMINATE while the connector is CREATING, so wait first or the
		// assertion below passes for the wrong reason.
		testClientHelper().OpenflowConnector.WaitUntilSettled(t, id, helpers.OpenflowConnectorStopTimeout)
		assertThatObject(t, objectassert.OpenflowConnector(t, id).HasStatus(sdk.OpenflowConnectorStatusStopped))

		// DROP only succeeds once TERMINATE has driven the connector to DELETED.
		err = client.OpenflowConnectors.Drop(ctx, sdk.NewDropOpenflowConnectorRequest(id))
		require.ErrorContains(t, err, "DROP not allowed")

		err = client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithTerminate(true))
		require.NoError(t, err)
		testClientHelper().OpenflowConnector.WaitForStatus(t, id, sdk.OpenflowConnectorStatusDeleted, helpers.OpenflowConnectorTerminateTimeout)

		err = client.OpenflowConnectors.Drop(ctx, sdk.NewDropOpenflowConnectorRequest(id))
		require.NoError(t, err)

		_, err = client.OpenflowConnectors.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop: when a connector does not exist", func(t *testing.T) {
		err := client.OpenflowConnectors.Drop(ctx, sdk.NewDropOpenflowConnectorRequest(NonExistingSchemaObjectIdentifier))
		assert.Error(t, err)
	})

	t.Run("drop: if exists on a non-existing connector", func(t *testing.T) {
		err := client.OpenflowConnectors.Drop(ctx, sdk.NewDropOpenflowConnectorRequest(NonExistingSchemaObjectIdentifier).WithIfExists(true))
		require.NoError(t, err)
	})
}

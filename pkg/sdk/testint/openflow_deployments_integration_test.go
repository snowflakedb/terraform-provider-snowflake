//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These use BYOC deployments, which settle at INACTIVE without any customer infrastructure. Creating a
// Snowflake-managed deployment provisions it for real and takes a long time, so the acceptance tests cover
// that flow and these cover the SQL operations.
//
// Teardown is TERMINATE, wait for DELETED, then DROP; see helpers.OpenflowDeploymentClient.DropFunc.
func TestInt_OpenflowDeployments(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	client := testClient(t)
	ctx := testContext(t)

	currentRole := testClientHelper().Context.CurrentRole(t)

	t.Run("create: basic byoc", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreateOpenflowDeploymentRequest(id, sdk.OpenflowDeploymentTypeByoc).
			WithVpcType(sdk.OpenflowVpcTypeManaged)

		err := client.OpenflowDeployments.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().OpenflowDeployment.DropFunc(t, id))

		// A BYOC deployment moves CREATING -> INACTIVE on its own, so settle first to make the status
		// deterministic. Teardown waits for the same thing before it can terminate.
		testClientHelper().OpenflowDeployment.WaitUntilSettled(t, id, helpers.OpenflowDeploymentActiveTimeout)

		deployment, err := client.OpenflowDeployments.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.OpenflowDeploymentFromObject(t, deployment).
				HasName(id.Name()).
				HasType(sdk.OpenflowDeploymentTypeByoc).
				HasStatus(sdk.OpenflowDeploymentStatusInactive).
				HasVpcType(sdk.OpenflowVpcTypeManaged).
				HasNoDisplayName().
				HasUsePrivateLink(false).
				HasUseUserAuthOverPrivateLink(false).
				HasNoCustomIngressHostname().
				// key names the AWS resources in a BYOC deployment's CloudFormation template.
				HasKeyNotEmpty().
				HasOwner(currentRole.Name()).
				HasNoComment().
				HasCreatedOnNotEmpty().
				HasUpdatedOnNotEmpty(),
		)
	})

	t.Run("create: complete byoc", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		comment := random.Comment()
		displayName := random.AlphaN(12)
		hostname := "ingress.example.com"

		request := sdk.NewCreateOpenflowDeploymentRequest(id, sdk.OpenflowDeploymentTypeByoc).
			WithIfNotExists(true).
			WithVpcType(sdk.OpenflowVpcTypeProvided).
			WithCustomIngressHostname(hostname).
			WithUsePrivateLink(true).
			WithUseUserAuthOverPrivatelink(true).
			WithDisplayName(displayName).
			WithComment(comment)

		err := client.OpenflowDeployments.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().OpenflowDeployment.DropFunc(t, id))

		testClientHelper().OpenflowDeployment.WaitUntilSettled(t, id, helpers.OpenflowDeploymentActiveTimeout)

		deployment, err := client.OpenflowDeployments.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.OpenflowDeploymentFromObject(t, deployment).
				HasName(id.Name()).
				HasType(sdk.OpenflowDeploymentTypeByoc).
				HasStatus(sdk.OpenflowDeploymentStatusInactive).
				HasVpcType(sdk.OpenflowVpcTypeProvided).
				HasDisplayName(displayName).
				HasUsePrivateLink(true).
				HasUseUserAuthOverPrivateLink(true).
				HasCustomIngressHostname(hostname).
				HasKeyNotEmpty().
				HasOwner(currentRole.Name()).
				HasComment(comment).
				HasCreatedOnNotEmpty().
				HasUpdatedOnNotEmpty(),
		)
	})

	t.Run("alter: set and unset", func(t *testing.T) {
		eventTable, eventTableCleanup := testClientHelper().EventTable.Create(t)
		t.Cleanup(eventTableCleanup)

		deployment, cleanup := testClientHelper().OpenflowDeployment.CreateByoc(t)
		t.Cleanup(cleanup)
		id := deployment.ID()

		// EVENT_TABLE is in neither SHOW nor DESCRIBE output, so SHOW PARAMETERS is the only way to read it.
		// Unsetting hands the deployment back whatever it inherits, so the account's own value is what to
		// expect afterwards rather than an empty one.
		accountEventTable := helpers.FindParameter(t, testClientHelper().Parameter.ShowAccountParameters(t), sdk.AccountParameterEventTable).Value

		comment := random.Comment()
		displayName := random.AlphaN(12)
		err := client.OpenflowDeployments.Alter(ctx, sdk.NewAlterOpenflowDeploymentRequest(id).WithSet(
			*sdk.NewOpenflowDeploymentSetRequest().
				WithComment(comment).
				WithDisplayName(displayName).
				WithEventTable(*sdk.NewOpenflowDeploymentEventTableRequest().WithEventTable(eventTable.ID())),
		))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.OpenflowDeployment(t, id).
				HasComment(comment).
				HasDisplayName(displayName),
		)
		assertThatObject(
			t, objectparametersassert.OpenflowDeploymentParameters(t, id).
				HasEventTable(eventTable.ID().FullyQualifiedName()).
				HasEventTableLevel(sdk.ParameterTypeOpenflowDeployment),
		)

		err = client.OpenflowDeployments.Alter(ctx, sdk.NewAlterOpenflowDeploymentRequest(id).WithUnset(
			*sdk.NewOpenflowDeploymentUnsetRequest().WithComment(true).WithDisplayName(true).WithEventTable(true),
		))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.OpenflowDeployment(t, id).
				HasNoComment().
				HasNoDisplayName(),
		)
		assertThatObject(
			t, objectparametersassert.OpenflowDeploymentParameters(t, id).
				HasEventTable(accountEventTable),
		)
	})

	t.Run("alter: rename", func(t *testing.T) {
		deployment, cleanup := testClientHelper().OpenflowDeployment.CreateByoc(t)
		id := deployment.ID()
		newId := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.OpenflowDeployments.Alter(ctx, sdk.NewAlterOpenflowDeploymentRequest(id).WithRenameTo(newId))
		if err != nil {
			t.Cleanup(cleanup)
			require.NoError(t, err)
		}
		t.Cleanup(testClientHelper().OpenflowDeployment.DropFunc(t, newId))

		_, err = client.OpenflowDeployments.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		assertThatObject(t, objectassert.OpenflowDeployment(t, newId).HasName(newId.Name()))
	})

	t.Run("describe", func(t *testing.T) {
		deployment, cleanup := testClientHelper().OpenflowDeployment.CreateByoc(t)
		t.Cleanup(cleanup)
		id := deployment.ID()

		// DESCRIBE returns the same columns as SHOW minus created_on and updated_on. CreateByoc makes it
		// without a display name, comment or ingress hostname.
		assertThatObject(
			t, objectassert.OpenflowDeploymentDetails(t, id).
				HasName(id.Name()).
				HasType(sdk.OpenflowDeploymentTypeByoc).
				HasStatus(sdk.OpenflowDeploymentStatusInactive).
				HasVpcType(sdk.OpenflowVpcTypeManaged).
				HasNoDisplayName().
				HasUsePrivateLink(false).
				HasUseUserAuthOverPrivateLink(false).
				HasNoCustomIngressHostname().
				HasKeyNotEmpty().
				HasOwner(currentRole.Name()).
				HasNoComment(),
		)
	})

	t.Run("show: with like", func(t *testing.T) {
		deployment, cleanup := testClientHelper().OpenflowDeployment.CreateByoc(t)
		t.Cleanup(cleanup)
		id := deployment.ID()

		deployments, err := client.OpenflowDeployments.Show(ctx, sdk.NewShowOpenflowDeploymentRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id.Name())}))
		require.NoError(t, err)
		require.Len(t, deployments, 1)
		assert.Equal(t, id.Name(), deployments[0].Name)
	})

	t.Run("show: unfiltered returns the deployment", func(t *testing.T) {
		deployment, cleanup := testClientHelper().OpenflowDeployment.CreateByoc(t)
		t.Cleanup(cleanup)
		id := deployment.ID()

		deployments, err := client.OpenflowDeployments.Show(ctx, sdk.NewShowOpenflowDeploymentRequest())
		require.NoError(t, err)
		assert.Contains(t, collections.Map(deployments, func(d sdk.OpenflowDeployment) string { return d.Name }), id.Name())
	})

	t.Run("show: with like on a non-existing deployment", func(t *testing.T) {
		deployments, err := client.OpenflowDeployments.Show(ctx, sdk.NewShowOpenflowDeploymentRequest().
			WithLike(sdk.Like{Pattern: sdk.String(NonExistingAccountObjectIdentifier.Name())}))
		require.NoError(t, err)
		assert.Empty(t, deployments)
	})

	t.Run("show: with starts with", func(t *testing.T) {
		deployment, cleanup := testClientHelper().OpenflowDeployment.CreateByoc(t)
		t.Cleanup(cleanup)
		id := deployment.ID()

		deployments, err := client.OpenflowDeployments.Show(ctx, sdk.NewShowOpenflowDeploymentRequest().
			WithStartsWith(id.Name()))
		require.NoError(t, err)
		// Exactly one: the name is unique to this run, so Contains would pass with the clause dropped.
		require.Len(t, deployments, 1)
		assert.Equal(t, id.Name(), deployments[0].Name)
	})

	// Asserts LIMIT actually limits. Creating a second deployment to guarantee the account holds more than one
	// costs minutes, so this compares against whatever an unlimited SHOW returns instead of skipping.
	t.Run("show: with limit", func(t *testing.T) {
		_, cleanup := testClientHelper().OpenflowDeployment.CreateByoc(t)
		t.Cleanup(cleanup)

		all, err := client.OpenflowDeployments.Show(ctx, sdk.NewShowOpenflowDeploymentRequest())
		require.NoError(t, err)
		require.NotEmpty(t, all)

		limited, err := client.OpenflowDeployments.Show(ctx, sdk.NewShowOpenflowDeploymentRequest().
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
		require.NoError(t, err)
		assert.Len(t, limited, 1)
		assert.LessOrEqual(t, len(limited), len(all), "LIMIT 1 returned more rows than an unlimited SHOW")
	})

	// Both clauses together, in the order Snowflake requires. STARTS WITH on a unique name returns one row on
	// its own, so LIMIT is covered by "show: with limit" above.
	t.Run("show: starts with and limit combined", func(t *testing.T) {
		deployment, cleanup := testClientHelper().OpenflowDeployment.CreateByoc(t)
		t.Cleanup(cleanup)
		id := deployment.ID()

		deployments, err := client.OpenflowDeployments.Show(ctx, sdk.NewShowOpenflowDeploymentRequest().
			WithStartsWith(id.Name()).
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
		require.NoError(t, err)
		require.Len(t, deployments, 1)
		assert.Equal(t, id.Name(), deployments[0].Name)
	})

	t.Run("drop: requires terminate first", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.OpenflowDeployments.Create(ctx, sdk.NewCreateOpenflowDeploymentRequest(id, sdk.OpenflowDeploymentTypeByoc).
			WithVpcType(sdk.OpenflowVpcTypeManaged))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().OpenflowDeployment.DropFunc(t, id))

		// Snowflake refuses DROP and TERMINATE while the deployment is CREATING, so wait first or the
		// assertion below passes for the wrong reason.
		testClientHelper().OpenflowDeployment.WaitUntilSettled(t, id, helpers.OpenflowDeploymentActiveTimeout)
		assertThatObject(t, objectassert.OpenflowDeployment(t, id).HasStatus(sdk.OpenflowDeploymentStatusInactive))

		// DROP only succeeds once the deployment has reached DELETED, which TERMINATE drives it to.
		err = client.OpenflowDeployments.Drop(ctx, sdk.NewDropOpenflowDeploymentRequest(id))
		require.Error(t, err)

		err = client.OpenflowDeployments.Alter(ctx, sdk.NewAlterOpenflowDeploymentRequest(id).WithTerminate(true))
		require.NoError(t, err)
		testClientHelper().OpenflowDeployment.WaitForStatus(t, id, sdk.OpenflowDeploymentStatusDeleted, helpers.OpenflowDeploymentTerminateTimeout)

		err = client.OpenflowDeployments.Drop(ctx, sdk.NewDropOpenflowDeploymentRequest(id))
		require.NoError(t, err)

		_, err = client.OpenflowDeployments.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop: when a deployment does not exist", func(t *testing.T) {
		err := client.OpenflowDeployments.Drop(ctx, sdk.NewDropOpenflowDeploymentRequest(NonExistingAccountObjectIdentifier))
		assert.Error(t, err)
	})

	t.Run("drop: if exists on a non-existing deployment", func(t *testing.T) {
		err := client.OpenflowDeployments.Drop(ctx, sdk.NewDropOpenflowDeploymentRequest(NonExistingAccountObjectIdentifier).WithIfExists(true))
		require.NoError(t, err)
	})
}

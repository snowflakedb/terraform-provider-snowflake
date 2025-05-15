//go:build !account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_ComputePools(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("create - basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreateComputePoolRequest(id, 1, 2, sdk.ComputePoolInstanceFamilyCPUX64XS)

		err := client.ComputePools.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ComputePool.DropFunc(t, id))

		computePool, err := client.ComputePools.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ComputePoolFromObject(t, computePool).
			HasName(id.Name()).
			HasState(sdk.ComputePoolStateStarting).
			HasMinNodes(1).
			HasMaxNodes(2).
			HasInstanceFamily(sdk.ComputePoolInstanceFamilyCPUX64XS).
			HasNumServices(0).
			HasNumJobs(0).
			HasAutoSuspendSecs(3600).
			HasAutoResume(true).
			HasActiveNodes(0).
			HasIdleNodes(0).
			HasTargetNodes(1).
			HasCreatedOnNotEmpty().
			HasResumedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasNoComment().
			HasIsExclusive(false).
			HasNoApplication().
			HasNoBudget(),
		)
	})

	t.Run("create - complete", func(t *testing.T) {
		applicationPackage, applicationPackageCleanup := createApplicationPackage(t)
		t.Cleanup(applicationPackageCleanup)
		application, applicationCleanup := testClientHelper().Application.CreateApplication(t, applicationPackage.ID(), "V01")
		t.Cleanup(applicationCleanup)
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		comment := random.Comment()
		request := sdk.NewCreateComputePoolRequest(id, 1, 2, sdk.ComputePoolInstanceFamilyCPUX64XS).
			WithForApplication(application.ID()).
			WithAutoResume(true).
			WithInitiallySuspended(true).
			WithAutoSuspendSecs(6767).
			WithComment(comment)

		err := client.ComputePools.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ComputePool.DropFunc(t, id))

		computePool, err := client.ComputePools.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ComputePoolFromObject(t, computePool).
			HasName(id.Name()).
			HasState(sdk.ComputePoolStateSuspended).
			HasMinNodes(1).
			HasMaxNodes(2).
			HasInstanceFamily(sdk.ComputePoolInstanceFamilyCPUX64XS).
			HasNumServices(0).
			HasNumJobs(0).
			HasAutoSuspendSecs(6767).
			HasAutoResume(true).
			HasActiveNodes(0).
			HasIdleNodes(0).
			HasTargetNodes(0).
			HasCreatedOnNotEmpty().
			HasResumedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasIsExclusive(true).
			HasApplication(application.ID().FullyQualifiedName()).
			HasNoBudget(),
		)
	})

	t.Run("alter: set", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreateComputePoolRequest(id, 1, 2, sdk.ComputePoolInstanceFamilyCPUX64XS).
			WithAutoSuspendSecs(6767).
			WithAutoResume(false)

		_, cleanup := testClientHelper().ComputePool.CreateWithRequest(t, request)
		t.Cleanup(cleanup)

		comment := random.Comment()
		err := client.ComputePools.Alter(ctx, sdk.NewAlterComputePoolRequest(id).WithSet(
			*sdk.NewComputePoolSetRequest().
				WithMinNodes(4).
				WithMaxNodes(5).
				WithAutoResume(true).
				WithAutoSuspendSecs(3600).
				WithComment(comment),
		))
		require.NoError(t, err)

		computePool, err := client.ComputePools.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ComputePoolFromObject(t, computePool).
			HasName(id.Name()).
			HasState(sdk.ComputePoolStateStarting).
			HasMinNodes(4).
			HasMaxNodes(5).
			HasInstanceFamily(sdk.ComputePoolInstanceFamilyCPUX64XS).
			HasNumServices(0).
			HasNumJobs(0).
			HasAutoSuspendSecs(3600).
			HasAutoResume(true).
			HasActiveNodes(0).
			HasIdleNodes(0).
			HasTargetNodes(4).
			HasCreatedOnNotEmpty().
			HasResumedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasIsExclusive(false).
			HasNoApplication().
			HasNoBudget(),
		)
	})

	t.Run("alter: unset", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreateComputePoolRequest(id, 1, 2, sdk.ComputePoolInstanceFamilyCPUX64XS).
			WithAutoSuspendSecs(6767).
			WithAutoResume(false).
			WithComment(random.Comment())

		_, cleanup := testClientHelper().ComputePool.CreateWithRequest(t, request)
		t.Cleanup(cleanup)

		err := client.ComputePools.Alter(ctx, sdk.NewAlterComputePoolRequest(id).WithUnset(
			*sdk.NewComputePoolUnsetRequest().
				WithAutoSuspendSecs(true).
				WithComment(true).
				WithAutoResume(true),
		))
		require.NoError(t, err)

		computePool, err := client.ComputePools.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ComputePoolFromObject(t, computePool).
			HasName(id.Name()).
			HasState(sdk.ComputePoolStateStarting).
			HasMinNodes(1).
			HasMaxNodes(2).
			HasInstanceFamily(sdk.ComputePoolInstanceFamilyCPUX64XS).
			HasNumServices(0).
			HasNumJobs(0).
			HasAutoSuspendSecs(3600).
			HasAutoResume(true).
			HasActiveNodes(0).
			HasIdleNodes(0).
			HasTargetNodes(1).
			HasCreatedOnNotEmpty().
			HasResumedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasNoComment().
			HasIsExclusive(false).
			HasNoApplication().
			HasNoBudget(),
		)
	})

	t.Run("describe", func(t *testing.T) {
		computePool, funcCleanup := testClientHelper().ComputePool.Create(t)
		t.Cleanup(funcCleanup)

		assertThatObject(t, objectassert.ComputePoolDetails(t, computePool.ID()).
			HasName(computePool.ID().Name()).
			HasState(sdk.ComputePoolStateStarting).
			HasMinNodes(1).
			HasMaxNodes(1).
			HasInstanceFamily(sdk.ComputePoolInstanceFamilyCPUX64XS).
			HasNumServices(0).
			HasNumJobs(0).
			HasAutoSuspendSecs(3600).
			HasAutoResume(true).
			HasActiveNodes(0).
			HasIdleNodes(0).
			HasTargetNodes(1).
			HasCreatedOnNotEmpty().
			HasResumedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasNoComment().
			HasIsExclusive(false).
			HasNoApplication().
			HasNoBudget().
			HasErrorCode("").
			HasStatusMessage("Compute pool is starting for last 0 minutes"),
		)
	})

	t.Run("drop: when an object already exists", func(t *testing.T) {
		computePool, computePoolCleanup := testClientHelper().ComputePool.Create(t)
		t.Cleanup(computePoolCleanup)
		id := computePool.ID()
		err := client.ComputePools.Drop(ctx, sdk.NewDropComputePoolRequest(id))
		require.NoError(t, err)
		_, err = client.ComputePools.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: when an object does not exist", func(t *testing.T) {
		err := client.ComputePools.Drop(ctx, sdk.NewDropComputePoolRequest(NonExistingAccountObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("show: with like", func(t *testing.T) {
		computePool, funcCleanup := testClientHelper().ComputePool.Create(t)
		t.Cleanup(funcCleanup)

		computePools, err := client.ComputePools.Show(ctx, sdk.NewShowComputePoolRequest().WithLike(sdk.Like{Pattern: &computePool.Name}))
		require.NoError(t, err)
		require.Equal(t, 1, len(computePools))
		require.NotNil(t, computePool)
		require.Equal(t, computePool.Name, computePools[0].Name)
		// TODO(next PR): Add more assertions, based on generated builders.
		// Note that the value of updated_on may differ between both SHOW calls.
	})
}

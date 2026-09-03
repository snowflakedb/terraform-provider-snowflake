//go:build account_level_tests

package testint

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_PostgresInstances(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// Sequential Create: concurrent CREATE POSTGRES INSTANCE races the account keychain.
	sharedInstance, sharedCleanup := testClientHelper().PostgresInstance.Create(t)
	t.Cleanup(sharedCleanup)

	// Created one version below the default so "alter: set and unset properties" can
	// exercise a POSTGRES_VERSION upgrade; downgrades are rejected by Snowflake.
	alterId := testClientHelper().Ids.RandomAccountObjectIdentifier()
	alterInstance, alterCleanup := testClientHelper().PostgresInstance.CreateWithRequest(t,
		sdk.NewCreatePostgresInstanceRequest(alterId, "STANDARD_M", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres).
			WithPostgresVersion(17))
	t.Cleanup(alterCleanup)

	suspendInstance, suspendCleanup := testClientHelper().PostgresInstance.Create(t)
	t.Cleanup(suspendCleanup)

	// READY waits run in parallel. Fork is issued once sharedInstance is READY, then
	// waited on in the same goroutine so fork: basic only asserts.
	var preForkedId sdk.AccountObjectIdentifier
	var wg sync.WaitGroup
	wg.Go(func() {
		sharedInstance = testClientHelper().PostgresInstance.WaitForReady(t, sharedInstance.ID(), 6*time.Minute)
		preForkedId = testClientHelper().Ids.RandomAccountObjectIdentifier()
		postgresForkEventually(t, client, sdk.NewForkPostgresInstanceRequest(preForkedId, sharedInstance.ID()))
		testClientHelper().PostgresInstance.WaitForReady(t, preForkedId, 5*time.Minute)
	})
	wg.Go(func() {
		alterInstance = testClientHelper().PostgresInstance.WaitForReady(t, alterInstance.ID(), 6*time.Minute)
	})
	wg.Go(func() {
		suspendInstance = testClientHelper().PostgresInstance.WaitForReady(t, suspendInstance.ID(), 6*time.Minute)
	})
	wg.Wait()
	if preForkedId.Name() != "" {
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, preForkedId))
	}

	t.Run("create: basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_M", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres)

		err := client.PostgresInstances.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, id))

		postgresInstance, err := client.PostgresInstances.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.PostgresInstanceFromObject(t, postgresInstance).
				HasName(id.Name()).
				HasComputeFamily("STANDARD_M").
				HasStorageSize(10).
				HasAuthenticationAuthority("POSTGRES").
				HasIsHighlyAvailable(false).
				HasType("PRIMARY").
				HasOwner(snowflakeroles.Accountadmin.Name()).
				HasOwnerRoleType("ROLE").
				HasNoComment().
				// TODO: Origin is documented behavior but not currently observed.
				// HasNoOrigin().
				HasCreatedOnNotEmpty().
				HasUpdatedOnNotEmpty(),
		)
	})

	t.Run("create: complete", func(t *testing.T) {
		networkRule, networkRuleCleanup := testClientHelper().NetworkRule.CreateWithRequest(t, sdk.NewCreateNetworkRuleRequest(
			testClientHelper().Ids.RandomSchemaObjectIdentifier(),
			sdk.NetworkRuleTypeIpv4,
			[]sdk.NetworkRuleValue{},
			sdk.NetworkRuleModePostgresIngress,
		))
		t.Cleanup(networkRuleCleanup)

		networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t,
			sdk.NewCreateNetworkPolicyRequest(testClientHelper().Ids.RandomAccountObjectIdentifier()).
				WithAllowedNetworkRuleList([]sdk.SchemaObjectIdentifier{networkRule.ID()}))
		t.Cleanup(networkPolicyCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		comment := random.Comment()
		request := sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_M", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres).
			WithPostgresVersion(17).
			WithHighAvailability(true).
			WithNetworkPolicy(networkPolicy.ID()).
			WithPostgresSettings(`{"postgres:work_mem": "128MB"}`).
			WithComment(comment)

		err := client.PostgresInstances.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, id))

		postgresInstance, err := client.PostgresInstances.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.PostgresInstanceFromObject(t, postgresInstance).
				HasName(id.Name()).
				HasComputeFamily("STANDARD_M").
				HasStorageSize(10).
				HasAuthenticationAuthority("POSTGRES").
				HasPostgresVersion("17").
				HasIsHighlyAvailable(true).
				HasComment(comment).
				HasPostgresSettings(`{"postgres:work_mem":"128MB"}`).
				HasCreatedOnNotEmpty().
				HasUpdatedOnNotEmpty(),
		)
	})

	t.Run("create: with tags", func(t *testing.T) {
		t.Skip("tagging for POSTGRES INSTANCE is not yet supported")
		tag1, tag1Cleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tag1Cleanup)
		tag2, tag2Cleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tag2Cleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_M", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres).
			WithTag([]sdk.TagAssociation{
				{
					Name:  tag1.ID(),
					Value: "value1",
				},
				{
					Name:  tag2.ID(),
					Value: "value2",
				},
			})

		err := client.PostgresInstances.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, id))

		assertTagSet(t, tag1.ID(), id, sdk.ObjectTypePostgresInstance, "value1")
		assertTagSet(t, tag2.ID(), id, sdk.ObjectTypePostgresInstance, "value2")
	})

	t.Run("create: with authentication_authority postgres_or_snowflake", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_M", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgresOrSnowflake)

		err := client.PostgresInstances.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, id))

		assertThatObject(
			t, objectassert.PostgresInstance(t, id).
				HasName(id.Name()).
				HasAuthenticationAuthority("POSTGRES_OR_SNOWFLAKE"),
		)
	})

	// TODO(SNOW-3580377): Investigate and unskip:
	// 2026-05-26T09:17:46.8049911Z 2026/05/26 09:17:46 2026/05/26 09:17:46 [DEBUG] err: 001008 (22023): SQL compilation error:
	// 2026-05-26T09:17:46.8052734Z 2026/05/26 09:17:46 invalid value [HTMJHKIT_2374DB97D6E784A70C5FBE98C2CAE7F9201484A0AL] for parameter 'STORAGE_INTEGRATION (must be of type POSTGRES_EXTERNAL_STORAGE)'
	// 2026-05-26T09:17:46.8054876Z 2026/05/26 09:17:46     postgres_instances_gen_integration_test.go:240:
	// 2026-05-26T09:17:46.8057404Z 2026/05/26 09:17:46         	Error Trace:	/home/runner/work/terraform-provider-snowflake/terraform-provider-snowflake/pkg/sdk/testint/postgres_instances_gen_integration_test.go:240
	// 2026-05-26T09:17:46.8059166Z 2026/05/26 09:17:46         	Error:      	Received unexpected error:
	// 2026-05-26T09:17:46.8060390Z 2026/05/26 09:17:46         	            	001008 (22023): SQL compilation error:
	// 2026-05-26T09:17:46.8062098Z 2026/05/26 09:17:46         	            	invalid value [HTMJHKIT_2374DB97D6E784A70C5FBE98C2CAE7F9201484A0AL] for parameter 'STORAGE_INTEGRATION (must be of type POSTGRES_EXTERNAL_STORAGE)'
	// 2026-05-26T09:17:46.8063895Z 2026/05/26 09:17:46         	Test:       	TestInt_PostgresInstances/create_-_with_storage_integration
	t.Run("create: with storage_integration", func(t *testing.T) {
		t.Skip("TODO(SNOW-3580377): Investigate and unskip")
		awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
		awsRoleARN := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)

		storageIntegration, storageIntegrationCleanup := testClientHelper().StorageIntegration.CreateS3(t, awsBucketUrl, awsRoleARN)
		t.Cleanup(storageIntegrationCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_M", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres).
			WithStorageIntegration(storageIntegration.ID())

		err := client.PostgresInstances.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, id))

		postgresInstance, err := client.PostgresInstances.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), postgresInstance.Name)
	})

	t.Run("fork: basic", func(t *testing.T) {
		// TODO: Origin is documented behavior for forks but not currently observed.
		// Poll and assertion commented out until behavior is enabled.
		// var forkedInstance *sdk.PostgresInstance
		// require.Eventually(t, func() bool {
		//   var showErr error
		//   forkedInstance, showErr = client.PostgresInstances.ShowByID(ctx, preForkedId)
		//   require.NoError(t, showErr)
		//   return forkedInstance.Origin != nil
		// }, 5*time.Minute, 3*time.Second)
		//
		// assertThatObject(t, objectassert.PostgresInstanceFromObject(t, forkedInstance).
		//   HasName(preForkedId.Name()).
		//   HasOriginContaining(sharedInstance.Name),
		// )

		forkedInstance, showErr := client.PostgresInstances.ShowByID(ctx, preForkedId)
		require.NoError(t, showErr)
		assertThatObject(
			t, objectassert.PostgresInstanceFromObject(t, forkedInstance).
				HasName(preForkedId.Name()),
		)
	})

	t.Run("fork: with time travel options", func(t *testing.T) {
		// TODO(SNOW-3543815): Crunchy rejects the fork with 400 "There is no backup available."
		// because the source is younger than its first backup. Snowflake's pre-flight only
		// checks a flat MAX_BACKUP_RETENTION_DAYS window from now, so the statement is
		// accepted and ApiExceptionUtil maps the 400 to generic 604032 instead of
		// POSTGRES_TIMESTAMP_NOT_WITHIN_RETENTION_PERIOD (604028). Unskip once that mapping
		// exists (and/or the pre-flight compares against the source's actual backup horizon).
		t.Skip("TODO(SNOW-3543815): time-travel fork needs a backup on the source; 604028 not wired yet")

		// Compute a timestamp guaranteed within the retention window: 1 minute after
		// the source instance was created, so it is always in the past and within
		// the instance's RetentionTime-day window.
		validTimestamp := sharedInstance.CreatedOn.UTC().Add(time.Minute).Format("2006-01-02 15:04:05")

		// AT with timestamp
		forkId1 := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request1 := sdk.NewForkPostgresInstanceRequest(forkId1, sharedInstance.ID()).
			WithAt(*sdk.NewPostgresInstanceForkAtRequest().WithTimestamp(validTimestamp))

		postgresForkEventually(t, client, request1)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, forkId1))

		forkedInstance, showErr := client.PostgresInstances.ShowByID(ctx, forkId1)
		require.NoError(t, showErr)
		assert.Equal(t, forkId1.Name(), forkedInstance.Name)

		// AT with offset and compute overrides
		forkId2 := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request2 := sdk.NewForkPostgresInstanceRequest(forkId2, sharedInstance.ID()).
			WithAt(*sdk.NewPostgresInstanceForkAtRequest().WithOffset("-60")).
			WithComment("Fork with offset and compute override")

		if err := client.PostgresInstances.Fork(ctx, request2); err == nil {
			t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, forkId2))

			forkedInstance, showErr := client.PostgresInstances.ShowByID(ctx, forkId2)
			require.NoError(t, showErr)
			assert.Equal(t, forkId2.Name(), forkedInstance.Name)
		}

		// BEFORE with timestamp
		forkId3 := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request3 := sdk.NewForkPostgresInstanceRequest(forkId3, sharedInstance.ID()).
			WithBefore(*sdk.NewPostgresInstanceForkBeforeRequest().WithTimestamp(validTimestamp))

		postgresForkEventually(t, client, request3)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, forkId3))

		forkedInstance, showErr = client.PostgresInstances.ShowByID(ctx, forkId3)
		require.NoError(t, showErr)
		assert.Equal(t, forkId3.Name(), forkedInstance.Name)

		// BEFORE with offset
		forkId4 := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request4 := sdk.NewForkPostgresInstanceRequest(forkId4, sharedInstance.ID()).
			WithBefore(*sdk.NewPostgresInstanceForkBeforeRequest().WithOffset("-60"))

		if err := client.PostgresInstances.Fork(ctx, request4); err == nil {
			t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, forkId4))

			forkedInstance, showErr = client.PostgresInstances.ShowByID(ctx, forkId4)
			require.NoError(t, showErr)
			assert.Equal(t, forkId4.Name(), forkedInstance.Name)
		}
	})

	t.Run("fork: with all optional parameters", func(t *testing.T) {
		comment := random.Comment()
		forkId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewForkPostgresInstanceRequest(forkId, sharedInstance.ID()).
			WithComputeFamily("STANDARD_M").
			WithStorageSizeGb(20).
			WithComment(comment)

		postgresForkEventually(t, client, request)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, forkId))

		forkedInstance, err := client.PostgresInstances.ShowByID(ctx, forkId)
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.PostgresInstanceFromObject(t, forkedInstance).
				HasName(forkId.Name()).
				HasComment(comment),
		)
	})

	t.Run("fork: from non-existing source", func(t *testing.T) {
		forkId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewForkPostgresInstanceRequest(forkId, NonExistingAccountObjectIdentifier)

		err := client.PostgresInstances.Fork(ctx, request)
		assert.Error(t, err)
	})

	t.Run("fork: validation: At and Before are mutually exclusive", func(t *testing.T) {
		forkId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewForkPostgresInstanceRequest(forkId, NonExistingAccountObjectIdentifier).
			WithAt(*sdk.NewPostgresInstanceForkAtRequest().WithTimestamp("2025-01-15 12:00:00")).
			WithBefore(*sdk.NewPostgresInstanceForkBeforeRequest().WithTimestamp("2025-01-15 12:00:00"))

		err := client.PostgresInstances.Fork(ctx, request)
		require.ErrorContains(t, err, "ForkPostgresInstanceOptions")
		require.ErrorContains(t, err, "are incompatible and cannot be set at the same time")
	})

	t.Run("alter: set and unset properties", func(t *testing.T) {
		networkRule, networkRuleCleanup := testClientHelper().NetworkRule.CreateWithRequest(t, sdk.NewCreateNetworkRuleRequest(
			testClientHelper().Ids.RandomSchemaObjectIdentifier(),
			sdk.NetworkRuleTypeIpv4,
			[]sdk.NetworkRuleValue{},
			sdk.NetworkRuleModePostgresIngress,
		))
		t.Cleanup(networkRuleCleanup)

		networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t,
			sdk.NewCreateNetworkPolicyRequest(testClientHelper().Ids.RandomAccountObjectIdentifier()).
				WithAllowedNetworkRuleList([]sdk.SchemaObjectIdentifier{networkRule.ID()}))
		t.Cleanup(networkPolicyCleanup)

		comment := random.Comment()
		postgresAlterSafely(t, client, sdk.NewAlterPostgresInstanceRequest(alterInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().
				WithComment(comment).
				WithStorageSizeGb(20).
				WithComputeFamily("STANDARD_L").
				WithNetworkPolicy(networkPolicy.ID()).
				WithAuthenticationAuthority(sdk.PostgresInstanceAuthenticationAuthorityPostgresOrSnowflake).
				WithApply(*sdk.NewPostgresInstanceApplyRequest().WithImmediately(true))))

		postgresAlterSafely(t, client, sdk.NewAlterPostgresInstanceRequest(alterInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().
				WithMaintenanceWindowStart(3)))

		postgresAlterSafely(t, client, sdk.NewAlterPostgresInstanceRequest(alterInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().
				WithPostgresSettings(`{"postgres:work_mem": "128MB"}`)))

		postgresAlterSafely(t, client, sdk.NewAlterPostgresInstanceRequest(alterInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().
				WithHighAvailability(true)))

		assertThatObject(
			t, objectassert.PostgresInstance(t, alterInstance.ID()).
				HasComment(comment).
				HasStorageSize(20).
				HasComputeFamily("STANDARD_L").
				HasIsHighlyAvailable(true).
				HasAuthenticationAuthority("POSTGRES_OR_SNOWFLAKE"),
		)

		details, err := client.PostgresInstances.DescribeDetails(ctx, alterInstance.ID())
		require.NoError(t, err)
		require.NotNil(t, details.MaintenanceWindowStart)
		assert.Equal(t, 3, *details.MaintenanceWindowStart)
		assert.False(t, details.HasAnyRunningOperations)

		postgresAlterSafely(t, client, sdk.NewAlterPostgresInstanceRequest(alterInstance.ID()).
			WithUnset(*sdk.NewPostgresInstanceUnsetRequest().
				WithComment(true).
				WithPostgresSettings(true).
				WithMaintenanceWindowStart(true).
				WithNetworkPolicy(true)))

		// UNSET POSTGRES_SETTINGS clears the settings to an empty JSON object rather than to NULL.
		assertThatObject(
			t, objectassert.PostgresInstance(t, alterInstance.ID()).
				HasNoComment().
				HasPostgresSettings("{}"),
		)

		postgresAlterSafely(t, client, sdk.NewAlterPostgresInstanceRequest(alterInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().
				WithStorageSizeGb(30).
				WithApply(*sdk.NewPostgresInstanceApplyRequest().WithImmediately(true))))

		assertThatObject(
			t, objectassert.PostgresInstance(t, alterInstance.ID()).
				HasStorageSize(30),
		)

		// This instance is created at version 17 during setup; POSTGRES_VERSION only moves forward.
		postgresAlterSafely(t, client, sdk.NewAlterPostgresInstanceRequest(alterInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().
				WithPostgresVersion(18).
				WithApply(*sdk.NewPostgresInstanceApplyRequest().WithImmediately(true))))

		assertThatObject(
			t, objectassert.PostgresInstance(t, alterInstance.ID()).
				HasPostgresVersion("18"),
		)

		// APPLY ON a future timestamp queues work that will not drain in this test. Snowflake
		// rejects timestamps more than 72 hours out, so schedule it a day from now.
		scheduledFor := time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02 15:04:05")
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(alterInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().
				WithStorageSizeGb(40).
				WithApply(*sdk.NewPostgresInstanceApplyRequest().WithOn(scheduledFor))))
		require.NoError(t, err)
	})

	t.Run("alter: suspend and resume", func(t *testing.T) {
		postgresAlterSafely(t, client, sdk.NewAlterPostgresInstanceRequest(suspendInstance.ID()).
			WithSuspend(true))

		result, err := client.PostgresInstances.ShowByID(ctx, suspendInstance.ID())
		require.NoError(t, err)
		assertThatObject(
			t, objectassert.PostgresInstanceFromObject(t, result).
				HasStateOneOf(
					sdk.PostgresInstanceStateSuspending,
					sdk.PostgresInstanceStateSuspended,
				),
		)

		require.Eventually(t, func() bool {
			result, err = client.PostgresInstances.ShowByID(ctx, suspendInstance.ID())
			require.NoError(t, err)
			return result.State == sdk.PostgresInstanceStateSuspended
		}, 2*time.Minute, 5*time.Second)

		// Suspending an already-suspended instance is accepted as a no-op rather than rejected.
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(suspendInstance.ID()).
			WithSuspend(true))
		require.NoError(t, err)

		result, err = client.PostgresInstances.ShowByID(ctx, suspendInstance.ID())
		require.NoError(t, err)
		assertThatObject(
			t, objectassert.PostgresInstanceFromObject(t, result).
				HasState(sdk.PostgresInstanceStateSuspended),
		)

		postgresAlterSafely(t, client, sdk.NewAlterPostgresInstanceRequest(suspendInstance.ID()).
			WithResume(true))

		result, err = client.PostgresInstances.ShowByID(ctx, suspendInstance.ID())
		require.NoError(t, err)
		assertThatObject(
			t, objectassert.PostgresInstanceFromObject(t, result).
				HasState(sdk.PostgresInstanceStateReady),
		)
	})

	t.Run("alter: rename", func(t *testing.T) {
		t.Skip("RENAME TO not yet supported for POSTGRES INSTANCE")
		postgresInstance1, cleanup1 := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup1)
		postgresInstance2, cleanup2 := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup2)

		// Rename instance1 to a new name
		newId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, newId))

		err := client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance1.ID()).
			WithRenameTo(newId))
		require.NoError(t, err)

		// Old name should not exist
		_, err = client.PostgresInstances.ShowByID(ctx, postgresInstance1.ID())
		require.Error(t, err)

		// New name should exist
		result, err := client.PostgresInstances.ShowByID(ctx, newId)
		require.NoError(t, err)
		assert.Equal(t, newId.Name(), result.Name)

		// Try to rename instance2 to the new name (already taken) - should fail
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance2.ID()).
			WithRenameTo(newId))
		assert.Error(t, err)
	})

	t.Run("alter: reset access", func(t *testing.T) {
		// Reset access for snowflake_admin
		err := client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(sharedInstance.ID()).
			WithResetAccess(*sdk.NewPostgresInstanceResetAccessRequest(sdk.PostgresInstanceResetAccessRoleSnowflakeAdmin)))
		require.NoError(t, err)

		// Reset access for application
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(sharedInstance.ID()).
			WithResetAccess(*sdk.NewPostgresInstanceResetAccessRequest(sdk.PostgresInstanceResetAccessRoleApplication)))
		require.NoError(t, err)
	})

	t.Run("alter: set and unset tags", func(t *testing.T) {
		t.Skip("tagging for POSTGRES INSTANCE is not yet supported")
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		tagAssociation := sdk.TagAssociation{
			Name:  tag.ID(),
			Value: "tag_value",
		}

		err := client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSetTags([]sdk.TagAssociation{tagAssociation}))
		require.NoError(t, err)

		assertTagSet(t, tag.ID(), postgresInstance.ID(), sdk.ObjectTypePostgresInstance, "tag_value")

		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithUnsetTags([]sdk.ObjectIdentifier{tag.ID()}))
		require.NoError(t, err)

		assertTagUnset(t, tag.ID(), postgresInstance.ID(), sdk.ObjectTypePostgresInstance)
	})

	// TODO(SNOW-3580377): Investigate and unskip:
	// 2026-05-26T09:19:25.3894211Z 2026/05/26 09:19:25 2026/05/26 09:19:25 [DEBUG] err: 001008 (22023): SQL compilation error:
	// 2026-05-26T09:19:25.3897018Z 2026/05/26 09:19:25 invalid value [MTKMMDIT_2374DB97D6E784A70C5FBE98C2CAE7F9201484A0AL] for parameter 'STORAGE_INTEGRATION (must be of type POSTGRES_EXTERNAL_STORAGE)'
	// 2026-05-26T09:19:25.3899241Z 2026/05/26 09:19:25     postgres_instances_gen_integration_test.go:671:
	// 2026-05-26T09:19:25.3901032Z 2026/05/26 09:19:25         	Error Trace:	/home/runner/work/terraform-provider-snowflake/terraform-provider-snowflake/pkg/sdk/testint/postgres_instances_gen_integration_test.go:671
	// 2026-05-26T09:19:25.3902718Z 2026/05/26 09:19:25         	Error:      	Received unexpected error:
	// 2026-05-26T09:19:25.3903936Z 2026/05/26 09:19:25         	            	001008 (22023): SQL compilation error:
	// 2026-05-26T09:19:25.3906056Z 2026/05/26 09:19:25         	            	invalid value [MTKMMDIT_2374DB97D6E784A70C5FBE98C2CAE7F9201484A0AL] for parameter 'STORAGE_INTEGRATION (must be of type POSTGRES_EXTERNAL_STORAGE)'
	// 2026-05-26T09:19:25.3907907Z 2026/05/26 09:19:25         	Test:       	TestInt_PostgresInstances/alter:_set_and_unset_storage_integration
	t.Run("alter: set and unset storage_integration", func(t *testing.T) {
		t.Skip("TODO(SNOW-3580377): Investigate and unskip")
		awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
		awsRoleARN := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)

		storageIntegration, storageIntegrationCleanup := testClientHelper().StorageIntegration.CreateS3(t, awsBucketUrl, awsRoleARN)
		t.Cleanup(storageIntegrationCleanup)

		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		// Set storage_integration
		err := client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().WithStorageIntegration(storageIntegration.ID())))
		require.NoError(t, err)

		// Unset storage_integration
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithUnset(*sdk.NewPostgresInstanceUnsetRequest().WithStorageIntegration(true)))
		require.NoError(t, err)
	})

	t.Run("alter: non-existing object", func(t *testing.T) {
		err := client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(NonExistingAccountObjectIdentifier).
			WithSet(*sdk.NewPostgresInstanceSetRequest().WithComment("test")))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)

		// With IF EXISTS should succeed
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(NonExistingAccountObjectIdentifier).
			WithIfExists(true).
			WithSet(*sdk.NewPostgresInstanceSetRequest().WithComment("test")))
		require.NoError(t, err)
	})

	t.Run("show: all", func(t *testing.T) {
		instances, err := client.PostgresInstances.Show(ctx, sdk.NewShowPostgresInstanceRequest())
		require.NoError(t, err)
		require.NotEmpty(t, instances)

		_, err = collections.FindFirst(instances, func(inst sdk.PostgresInstance) bool {
			return inst.Name == sharedInstance.Name
		})
		require.NoError(t, err)
	})

	t.Run("show: like", func(t *testing.T) {
		instances, err := client.PostgresInstances.Show(ctx, sdk.NewShowPostgresInstanceRequest().
			WithLike(sdk.Like{Pattern: sdk.String(sharedInstance.Name)}))
		require.NoError(t, err)
		require.Len(t, instances, 1)
		assert.Equal(t, sharedInstance.Name, instances[0].Name)
	})

	t.Run("show: starts with", func(t *testing.T) {
		instances, err := client.PostgresInstances.Show(ctx, sdk.NewShowPostgresInstanceRequest().
			WithStartsWith(sharedInstance.Name))
		require.NoError(t, err)
		require.NotEmpty(t, instances)
		assert.Equal(t, sharedInstance.Name, instances[0].Name)
	})

	t.Run("ShowByID and ShowByIDSafely", func(t *testing.T) {
		result, err := client.PostgresInstances.ShowByID(ctx, sharedInstance.ID())
		require.NoError(t, err)
		assert.Equal(t, sharedInstance.Name, result.Name)

		assertThatObject(
			t, objectassert.PostgresInstanceFromObject(t, result).
				HasName(sharedInstance.Name).
				HasOwner(snowflakeroles.Accountadmin.Name()).
				HasOwnerRoleType("ROLE").
				HasType("PRIMARY").
				HasComputeFamily("STANDARD_M").
				HasStorageSize(10).
				HasAuthenticationAuthority("POSTGRES").
				HasIsHighlyAvailable(false).
				HasRetentionTime(0).
				HasPostgresVersionNotEmpty().
				HasStateOneOf(
					sdk.PostgresInstanceStateCreating,
					sdk.PostgresInstanceStateStarting,
					sdk.PostgresInstanceStateReady,
				).
				HasNoComment().
				// TODO: Origin is documented behavior but not currently observed.
				// HasNoOrigin().
				HasCreatedOnNotEmpty().
				HasUpdatedOnNotEmpty(),
		)

		result, err = client.PostgresInstances.ShowByIDSafely(ctx, sharedInstance.ID())
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, sharedInstance.Name, result.Name)
	})

	t.Run("ShowByID: missing object", func(t *testing.T) {
		_, err := client.PostgresInstances.ShowByID(ctx, testClientHelper().Ids.RandomAccountObjectIdentifier())
		require.Error(t, err)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("ShowByIDSafely: missing object", func(t *testing.T) {
		_, err := client.PostgresInstances.ShowByIDSafely(ctx, testClientHelper().Ids.RandomAccountObjectIdentifier())
		require.Error(t, err)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("describe", func(t *testing.T) {
		assertThatObject(
			t, objectassert.PostgresInstanceDetails(t, sharedInstance.ID()).
				HasName(sharedInstance.ID().Name()).
				HasOwner(snowflakeroles.Accountadmin.Name()).
				HasOwnerRoleType("ROLE").
				HasCreatedOnNotEmpty().
				HasUpdatedOnNotEmpty().
				HasType("PRIMARY").
				HasHostNotEmpty().
				HasComputeFamily("STANDARD_M").
				HasStorageSizeGb(10).
				HasPostgresVersionNotEmpty().
				HasHighAvailability(false).
				HasAuthenticationAuthority("POSTGRES").
				HasStateNotEmpty(),
		)
	})

	t.Run("drop: existing object", func(t *testing.T) {
		postgresInstance, _ := testClientHelper().PostgresInstance.Create(t)

		err := client.PostgresInstances.Drop(ctx, sdk.NewDropPostgresInstanceRequest(postgresInstance.ID()))
		require.NoError(t, err)

		_, err = client.PostgresInstances.ShowByID(ctx, postgresInstance.ID())
		assert.Error(t, err)
	})

	t.Run("drop: non-existing and already dropped", func(t *testing.T) {
		// Drop non-existing without IF EXISTS should error
		err := client.PostgresInstances.Drop(ctx, sdk.NewDropPostgresInstanceRequest(NonExistingAccountObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)

		// Drop non-existing with IF EXISTS should succeed
		err = client.PostgresInstances.Drop(ctx, sdk.NewDropPostgresInstanceRequest(NonExistingAccountObjectIdentifier).WithIfExists(true))
		require.NoError(t, err)

		// Create then drop, then try again
		postgresInstance, _ := testClientHelper().PostgresInstance.Create(t)

		// First drop succeeds
		err = client.PostgresInstances.Drop(ctx, sdk.NewDropPostgresInstanceRequest(postgresInstance.ID()))
		require.NoError(t, err)

		// Second drop without IF EXISTS should error
		err = client.PostgresInstances.Drop(ctx, sdk.NewDropPostgresInstanceRequest(postgresInstance.ID()))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)

		// Second drop with IF EXISTS should succeed
		err = client.PostgresInstances.Drop(ctx, sdk.NewDropPostgresInstanceRequest(postgresInstance.ID()).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("drop safely: existing object", func(t *testing.T) {
		postgresInstance, _ := testClientHelper().PostgresInstance.Create(t)

		err := client.PostgresInstances.DropSafely(ctx, postgresInstance.ID())
		require.NoError(t, err)

		_, err = client.PostgresInstances.ShowByID(ctx, postgresInstance.ID())
		assert.Error(t, err)
	})

	t.Run("drop safely: non-existing object", func(t *testing.T) {
		err := client.PostgresInstances.DropSafely(ctx, NonExistingAccountObjectIdentifier)
		require.NoError(t, err)
	})
}

func postgresAlterSafely(t *testing.T, client *sdk.Client, req *sdk.AlterPostgresInstanceRequest) {
	t.Helper()
	ctx, cancel := context.WithTimeout(testContext(t), 10*time.Minute)
	defer cancel()
	require.NoError(t, client.PostgresInstances.AlterSafely(ctx, req))
}

// postgresForkEventually retries a fork until it succeeds. A source instance in READY state is
// not necessarily fork-ready, and a fork issued too early fails with an internal error.
func postgresForkEventually(t *testing.T, client *sdk.Client, req *sdk.ForkPostgresInstanceRequest) {
	t.Helper()
	ctx := testContext(t)
	var err error
	require.Eventually(t, func() bool {
		err = client.PostgresInstances.Fork(ctx, req)
		return err == nil
	}, 2*time.Minute, 5*time.Second)
	require.NoError(t, err)
}

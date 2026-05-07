//go:build account_level_tests

package testint

import (
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_PostgresInstances(t *testing.T) {
	// Currently, almost all Postgres Instances tests are failing due to the following error:
	// 604001 (0A000): Compute Family STANDARD_1 is not a supported compute family for Postgres
	// This will be addressed in https://github.com/snowflakedb/terraform-provider-snowflake/pull/4704
	t.Skip("Skipping all Postgres Instances tests")
	client := testClient(t)
	ctx := testContext(t)

	// ==================
	// Create
	// ==================

	// Doc example: CREATE POSTGRES INSTANCE my_postgres COMPUTE_FAMILY = 'STANDARD_S' STORAGE_SIZE_GB = 50 AUTHENTICATION_AUTHORITY = POSTGRES;
	t.Run("create - basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_M", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres)

		err := client.PostgresInstances.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, id))

		postgresInstance, err := client.PostgresInstances.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstanceFromObject(t, postgresInstance).
			HasName(id.Name()).
			HasComputeFamily("STANDARD_M").
			HasStorageSize(10).
			HasAuthenticationAuthority("POSTGRES").
			HasIsHa(false).
			HasType("PRIMARY").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasNoComment().
			HasNoOrigin().
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty(),
		)
	})

	// Doc example: CREATE POSTGRES INSTANCE prod_postgres COMPUTE_FAMILY = 'STANDARD_M' STORAGE_SIZE_GB = 500
	//   AUTHENTICATION_AUTHORITY = POSTGRES POSTGRES_VERSION = 17 HIGH_AVAILABILITY = TRUE
	//   NETWORK_POLICY = 'my_network_policy' POSTGRES_SETTINGS = '{"postgres:work_mem": "128MB"}'
	//   COMMENT = 'Production Postgres instance';
	t.Run("create - complete", func(t *testing.T) {
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
			WithNetworkPolicy(networkPolicy.Name).
			WithPostgresSettings(`{"postgres:work_mem": "128MB"}`).
			WithComment(comment)

		err := client.PostgresInstances.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, id))

		postgresInstance, err := client.PostgresInstances.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstanceFromObject(t, postgresInstance).
			HasName(id.Name()).
			HasComputeFamily("STANDARD_M").
			HasStorageSize(10).
			HasAuthenticationAuthority("POSTGRES").
			HasPostgresVersion("17").
			HasIsHa(true).
			HasComment(comment).
			HasPostgresSettings(`{"postgres:work_mem":"128MB"}`).
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty(),
		)
	})

	// Doc example: TAG (tag1 = 'value1')
	t.Run("create - with tags", func(t *testing.T) {
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

	// Doc example: CREATE POSTGRES INSTANCE <name> COMPUTE_FAMILY = 'STANDARD_S' ... AUTHENTICATION_AUTHORITY = POSTGRES_OR_SNOWFLAKE
	t.Run("create - with authentication_authority postgres_or_snowflake", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_M", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgresOrSnowflake)

		err := client.PostgresInstances.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, id))

		assertThatObject(t, objectassert.PostgresInstance(t, id).
			HasName(id.Name()).
			HasAuthenticationAuthority("POSTGRES_OR_SNOWFLAKE"),
		)
	})

	t.Run("create - with storage_integration", func(t *testing.T) {
		awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
		awsRoleARN := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)

		storageIntegration, storageIntegrationCleanup := testClientHelper().StorageIntegration.CreateS3(t, awsBucketUrl, awsRoleARN)
		t.Cleanup(storageIntegrationCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_M", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres).
			WithStorageIntegration(storageIntegration.Name)

		err := client.PostgresInstances.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, id))

		postgresInstance, err := client.PostgresInstances.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), postgresInstance.Name)
	})

	// ==================
	// Fork
	// ==================

	// Doc example: CREATE POSTGRES INSTANCE my_fork FORK my_source_instance;
	t.Run("fork - basic", func(t *testing.T) {
		sourceInstance, sourceCleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(sourceCleanup)
		testClientHelper().PostgresInstance.WaitForReady(t, sourceInstance.ID(), 5*time.Minute)

		forkId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewForkPostgresInstanceRequest(forkId, sourceInstance.ID())

		// Retry fork operation — instance may not be internally fork-ready despite being in READY state
		var err error
		require.Eventually(t, func() bool {
			err = client.PostgresInstances.Fork(ctx, request)
			return err == nil
		}, 2*time.Minute, 5*time.Second)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, forkId))

		// Wait for fork to reach READY state — metadata (origin, type) populates after provisioning completes
		testClientHelper().PostgresInstance.WaitForReady(t, forkId, 5*time.Minute)

		// Poll until origin metadata is populated (fail if not received within timeout)
		var forkedInstance *sdk.PostgresInstance
		require.Eventually(t, func() bool {
			var showErr error
			forkedInstance, showErr = client.PostgresInstances.ShowByID(ctx, forkId)
			require.NoError(t, showErr)
			return forkedInstance.Origin != nil
		}, 5*time.Minute, 5*time.Second)

		assertThatObject(t, objectassert.PostgresInstanceFromObject(t, forkedInstance).
			HasName(forkId.Name()).
			HasOriginContaining(sourceInstance.Name),
		)
	})

	// Doc example: CREATE POSTGRES INSTANCE my_fork FORK my_source_instance AT (TIMESTAMP => '2025-01-15 12:00:00'::TIMESTAMP_NTZ);
	t.Run("fork - with time travel options", func(t *testing.T) {
		sourceInstance, sourceCleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(sourceCleanup)

		// AT with timestamp
		forkId1 := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request1 := sdk.NewForkPostgresInstanceRequest(forkId1, sourceInstance.ID()).
			WithAt(*sdk.NewPostgresInstanceForkAtRequest().WithTimestamp("2025-01-15 12:00:00"))

		// This may fail if the timestamp is outside the retention period; that's expected
		err := client.PostgresInstances.Fork(ctx, request1)
		if err == nil {
			t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, forkId1))

			forkedInstance, showErr := client.PostgresInstances.ShowByID(ctx, forkId1)
			require.NoError(t, showErr)
			assert.Equal(t, forkId1.Name(), forkedInstance.Name)
		}

		// AT with offset and compute overrides
		forkId2 := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request2 := sdk.NewForkPostgresInstanceRequest(forkId2, sourceInstance.ID()).
			WithAt(*sdk.NewPostgresInstanceForkAtRequest().WithOffset("-60")).
			WithComment("Fork with offset and compute override")

		err = client.PostgresInstances.Fork(ctx, request2)
		if err == nil {
			t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, forkId2))

			forkedInstance, showErr := client.PostgresInstances.ShowByID(ctx, forkId2)
			require.NoError(t, showErr)
			assert.Equal(t, forkId2.Name(), forkedInstance.Name)
		}

		// BEFORE with timestamp
		forkId3 := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request3 := sdk.NewForkPostgresInstanceRequest(forkId3, sourceInstance.ID()).
			WithBefore(*sdk.NewPostgresInstanceForkBeforeRequest().WithTimestamp("2025-01-15 12:00:00"))

		// This may fail if the timestamp is outside the retention period; that's expected
		err = client.PostgresInstances.Fork(ctx, request3)
		if err == nil {
			t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, forkId3))

			forkedInstance, showErr := client.PostgresInstances.ShowByID(ctx, forkId3)
			require.NoError(t, showErr)
			assert.Equal(t, forkId3.Name(), forkedInstance.Name)
		}

		// BEFORE with offset
		forkId4 := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request4 := sdk.NewForkPostgresInstanceRequest(forkId4, sourceInstance.ID()).
			WithBefore(*sdk.NewPostgresInstanceForkBeforeRequest().WithOffset("-60"))

		err = client.PostgresInstances.Fork(ctx, request4)
		if err == nil {
			t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, forkId4))

			forkedInstance, showErr := client.PostgresInstances.ShowByID(ctx, forkId4)
			require.NoError(t, showErr)
			assert.Equal(t, forkId4.Name(), forkedInstance.Name)
		}
	})

	t.Run("fork - with all optional parameters", func(t *testing.T) {
		sourceInstance, sourceCleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(sourceCleanup)
		testClientHelper().PostgresInstance.WaitForReady(t, sourceInstance.ID(), 5*time.Minute)

		comment := random.Comment()
		forkId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewForkPostgresInstanceRequest(forkId, sourceInstance.ID()).
			WithComputeFamily("STANDARD_M").
			WithStorageSizeGb(20).
			WithComment(comment)

		// Retry fork operation — instance may not be internally fork-ready despite being in READY state
		var err error
		require.Eventually(t, func() bool {
			err = client.PostgresInstances.Fork(ctx, request)
			return err == nil
		}, 2*time.Minute, 5*time.Second)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, forkId))

		forkedInstance, err := client.PostgresInstances.ShowByID(ctx, forkId)
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstanceFromObject(t, forkedInstance).
			HasName(forkId.Name()).
			HasComment(comment),
		)
	})

	t.Run("fork - from non-existing source", func(t *testing.T) {
		forkId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewForkPostgresInstanceRequest(forkId, NonExistingAccountObjectIdentifier)

		err := client.PostgresInstances.Fork(ctx, request)
		assert.Error(t, err)
	})

	// ==================
	// Alter
	// ==================

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

		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)
		testClientHelper().PostgresInstance.WaitForReady(t, postgresInstance.ID(), 5*time.Minute)

		// Set compute/storage properties with APPLY IMMEDIATELY to ensure the operation
		// completes before subsequent ALTERs that conflict with in-progress compute/storage changes
		comment := random.Comment()
		var err error
		require.Eventually(t, func() bool {
			err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
				WithSet(*sdk.NewPostgresInstanceSetRequest().
					WithComment(comment).
					WithStorageSizeGb(20).
					WithComputeFamily("STANDARD_L").
					WithNetworkPolicy(networkPolicy.Name).
					WithAuthenticationAuthority(sdk.PostgresInstanceAuthenticationAuthorityPostgresOrSnowflake).
					WithApply(*sdk.NewPostgresInstanceApplyRequest().WithImmediately(true))))
			return err == nil
		}, 2*time.Minute, 5*time.Second)
		require.NoError(t, err)

		// Wait for the instance to return to READY after compute/storage changes complete
		testClientHelper().PostgresInstance.WaitForReady(t, postgresInstance.ID(), 3*time.Minute)

		// Set MAINTENANCE_WINDOW_START separately (does not conflict with compute/storage)
		require.Eventually(t, func() bool {
			err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
				WithSet(*sdk.NewPostgresInstanceSetRequest().
					WithMaintenanceWindowStart(3)))
			return err == nil
		}, 3*time.Minute, 5*time.Second)
		require.NoError(t, err)

		// Set HIGH_AVAILABILITY separately (cannot be combined with COMPUTE_FAMILY/STORAGE_SIZE_GB)
		// Retry because the previous compute/storage change may still be in progress
		require.Eventually(t, func() bool {
			err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
				WithSet(*sdk.NewPostgresInstanceSetRequest().
					WithHighAvailability(true)))
			return err == nil
		}, 6*time.Minute, 5*time.Second)
		require.NoError(t, err)

		// Set postgres settings separately (cannot be combined with COMPUTE_FAMILY/STORAGE_SIZE_GB/HIGH_AVAILABILITY)
		// Retry because the previous HA operation may still be in progress
		require.Eventually(t, func() bool {
			err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
				WithSet(*sdk.NewPostgresInstanceSetRequest().
					WithPostgresSettings(`{"postgres:work_mem": "128MB"}`)))
			return err == nil
		}, 5*time.Minute, 5*time.Second)
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstance(t, postgresInstance.ID()).
			HasComment(comment).
			HasStorageSize(20).
			HasComputeFamily("STANDARD_L").
			HasIsHa(true).
			HasAuthenticationAuthority("POSTGRES_OR_SNOWFLAKE"),
		)

		properties, err := client.PostgresInstances.Describe(ctx, postgresInstance.ID())
		require.NoError(t, err)
		propertyMap := make(map[string]string)
		for _, p := range properties {
			propertyMap[p.Property] = p.Value
		}
		assert.Equal(t, "3", propertyMap["maintenance_window_start"])

		// Unset all unsettable properties in one call
		// Retry because the previous postgres_settings alter may still be in progress
		require.Eventually(t, func() bool {
			err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
				WithUnset(*sdk.NewPostgresInstanceUnsetRequest().
					WithComment(true).
					WithPostgresSettings(true).
					WithMaintenanceWindowStart(true).
					WithNetworkPolicy(true)))
			return err == nil
		}, 5*time.Minute, 5*time.Second)
		require.NoError(t, err)

		// Poll until unset properties propagate (HA may still be processing)
		require.Eventually(t, func() bool {
			instance, showErr := client.PostgresInstances.ShowByID(ctx, postgresInstance.ID())
			if showErr != nil {
				return false
			}
			return instance.Comment == nil && instance.PostgresSettings == nil
		}, 5*time.Minute, 5*time.Second)

		// Set with apply immediately
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().
				WithStorageSizeGb(30).
				WithApply(*sdk.NewPostgresInstanceApplyRequest().WithImmediately(true))))
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstance(t, postgresInstance.ID()).
			HasStorageSize(30),
		)

		// Set postgres_version
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().
				WithPostgresVersion(17).
				WithApply(*sdk.NewPostgresInstanceApplyRequest().WithImmediately(true))))
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstance(t, postgresInstance.ID()).
			HasPostgresVersion("17"),
		)

		// Set with apply on timestamp
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().
				WithStorageSizeGb(40).
				WithApply(*sdk.NewPostgresInstanceApplyRequest().WithOn("2099-01-01 00:00:00"))))
		require.NoError(t, err)
	})

	t.Run("alter: suspend and resume", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)
		testClientHelper().PostgresInstance.WaitForReady(t, postgresInstance.ID(), 5*time.Minute)

		// Suspend — retry because the instance may not be fully ready for suspend despite READY state
		var err error
		require.Eventually(t, func() bool {
			err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
				WithSuspend(true))
			return err == nil
		}, 2*time.Minute, 5*time.Second)
		require.NoError(t, err)

		// Wait for instance to reach SUSPENDING or SUSPENDED state
		var result *sdk.PostgresInstance
		require.Eventually(t, func() bool {
			result, err = client.PostgresInstances.ShowByID(ctx, postgresInstance.ID())
			require.NoError(t, err)
			return result.State == sdk.PostgresInstanceStateSuspending || result.State == sdk.PostgresInstanceStateSuspended
		}, 2*time.Minute, 5*time.Second)
		assertThatObject(t, objectassert.PostgresInstanceFromObject(t, result).
			HasStateOneOf(
				sdk.PostgresInstanceStateSuspending,
				sdk.PostgresInstanceStateSuspended,
			),
		)

		// Suspend again - expect error due to invalid state
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSuspend(true))
		assert.Error(t, err)

		// Wait for SUSPENDED state before resuming
		require.Eventually(t, func() bool {
			result, err = client.PostgresInstances.ShowByID(ctx, postgresInstance.ID())
			require.NoError(t, err)
			return result.State == sdk.PostgresInstanceStateSuspended
		}, 2*time.Minute, 5*time.Second)

		// Resume
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithResume(true))
		require.NoError(t, err)

		// Wait for state to transition from SUSPENDED after resume
		require.Eventually(t, func() bool {
			result, err = client.PostgresInstances.ShowByID(ctx, postgresInstance.ID())
			require.NoError(t, err)
			return result.State != sdk.PostgresInstanceStateSuspended
		}, 2*time.Minute, 5*time.Second)
		// State may be RESUMING, STARTING, CREATING, or READY depending on timing
		assertThatObject(t, objectassert.PostgresInstanceFromObject(t, result).
			HasStateOneOf(
				sdk.PostgresInstanceStateResuming,
				sdk.PostgresInstanceStateStarting,
				sdk.PostgresInstanceStateCreating,
				sdk.PostgresInstanceStateReady,
			),
		)
	})

	// Doc example: ALTER POSTGRES INSTANCE my_postgres RENAME TO prod_postgres;
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
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)
		testClientHelper().PostgresInstance.WaitForReady(t, postgresInstance.ID(), 5*time.Minute)

		// Reset access for snowflake_admin
		err := client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithResetAccess(*sdk.NewPostgresInstanceResetAccessRequest(sdk.PostgresInstanceResetAccessRoleSnowflakeAdmin)))
		require.NoError(t, err)

		// Reset access for application
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithResetAccess(*sdk.NewPostgresInstanceResetAccessRequest(sdk.PostgresInstanceResetAccessRoleApplication)))
		require.NoError(t, err)
	})

	// Doc example: ALTER POSTGRES INSTANCE <name> SET TAG <tag_name> = '<tag_value>'
	// Doc example: ALTER POSTGRES INSTANCE <name> UNSET TAG <tag_name>
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

	t.Run("alter: set and unset storage_integration", func(t *testing.T) {
		awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
		awsRoleARN := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)

		storageIntegration, storageIntegrationCleanup := testClientHelper().StorageIntegration.CreateS3(t, awsBucketUrl, awsRoleARN)
		t.Cleanup(storageIntegrationCleanup)

		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		// Set storage_integration
		err := client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().WithStorageIntegration(storageIntegration.Name)))
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

	// ==================
	// Show
	// ==================

	// Doc example: SHOW POSTGRES INSTANCES;
	t.Run("show: all, like, starts_with, and verify fields", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		// Show all
		instances, err := client.PostgresInstances.Show(ctx, sdk.NewShowPostgresInstanceRequest())
		require.NoError(t, err)
		require.NotEmpty(t, instances)

		found := false
		for _, inst := range instances {
			if inst.Name == postgresInstance.Name {
				found = true
				break
			}
		}
		assert.True(t, found, "expected to find created postgres instance in show results")

		// Show with LIKE
		instances, err = client.PostgresInstances.Show(ctx, sdk.NewShowPostgresInstanceRequest().
			WithLike(sdk.Like{Pattern: sdk.String(postgresInstance.Name)}))
		require.NoError(t, err)
		require.Len(t, instances, 1)
		assert.Equal(t, postgresInstance.Name, instances[0].Name)

		// Show with STARTS WITH
		instances, err = client.PostgresInstances.Show(ctx, sdk.NewShowPostgresInstanceRequest().
			WithStartsWith(postgresInstance.Name))
		require.NoError(t, err)
		require.NotEmpty(t, instances)
		assert.Equal(t, postgresInstance.Name, instances[0].Name)

		// Verify all result fields via ShowByID
		result, err := client.PostgresInstances.ShowByID(ctx, postgresInstance.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstanceFromObject(t, result).
			HasName(postgresInstance.Name).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasType("PRIMARY").
			HasComputeFamily("STANDARD_M").
			HasStorageSize(10).
			HasAuthenticationAuthority("POSTGRES").
			HasIsHa(false).
			HasRetentionTime(0).
			HasPostgresVersionNotEmpty().
			HasStateOneOf(
				sdk.PostgresInstanceStateCreating,
				sdk.PostgresInstanceStateStarting,
				sdk.PostgresInstanceStateReady,
			).
			HasNoComment().
			HasNoOrigin().
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty(),
		)
	})

	t.Run("ShowByID and ShowByIDSafely", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		// ShowByID
		result, err := client.PostgresInstances.ShowByID(ctx, postgresInstance.ID())
		require.NoError(t, err)
		assert.Equal(t, postgresInstance.Name, result.Name)

		// ShowByIDSafely
		result, err = client.PostgresInstances.ShowByIDSafely(ctx, postgresInstance.ID())
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, postgresInstance.Name, result.Name)
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

	// ==================
	// Describe
	// ==================

	// Doc example: DESCRIBE POSTGRES INSTANCE my_postgres;
	t.Run("describe", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		properties, err := client.PostgresInstances.Describe(ctx, postgresInstance.ID())
		require.NoError(t, err)
		require.NotEmpty(t, properties)

		// Verify all documented properties are present
		propertyMap := make(map[string]string)
		for _, p := range properties {
			propertyMap[p.Property] = p.Value
		}

		assert.Contains(t, propertyMap, "name")
		assert.Contains(t, propertyMap, "owner")
		assert.Contains(t, propertyMap, "owner_role_type")
		assert.Contains(t, propertyMap, "created_on")
		assert.Contains(t, propertyMap, "updated_on")
		assert.Contains(t, propertyMap, "type")
		assert.Contains(t, propertyMap, "host")
		assert.Contains(t, propertyMap, "compute_family")
		assert.Contains(t, propertyMap, "storage_size_gb")
		assert.Contains(t, propertyMap, "postgres_version")
		assert.Contains(t, propertyMap, "high_availability")
		assert.Contains(t, propertyMap, "authentication_authority")
		assert.Contains(t, propertyMap, "state")

		// Verify property values match expected defaults
		assert.Equal(t, postgresInstance.ID().Name(), propertyMap["name"])
		assert.Equal(t, snowflakeroles.Accountadmin.Name(), propertyMap["owner"])
		assert.Equal(t, "ROLE", propertyMap["owner_role_type"])
		assert.NotEmpty(t, propertyMap["created_on"])
		assert.NotEmpty(t, propertyMap["updated_on"])
		assert.Equal(t, "PRIMARY", propertyMap["type"])
		assert.NotEmpty(t, propertyMap["host"])
		assert.Equal(t, "STANDARD_M", propertyMap["compute_family"])
		assert.Equal(t, "10", propertyMap["storage_size_gb"])
		assert.NotEmpty(t, propertyMap["postgres_version"])
		assert.Equal(t, "false", propertyMap["high_availability"])
		assert.Equal(t, "POSTGRES", propertyMap["authentication_authority"])
		assert.NotEmpty(t, propertyMap["state"])
	})

	// ==================
	// Drop
	// ==================

	// Doc example: DROP POSTGRES INSTANCE my_postgres;
	t.Run("drop: existing object", func(t *testing.T) {
		postgresInstance, _ := testClientHelper().PostgresInstance.Create(t)

		err := client.PostgresInstances.Drop(ctx, sdk.NewDropPostgresInstanceRequest(postgresInstance.ID()))
		require.NoError(t, err)

		_, err = client.PostgresInstances.ShowByID(ctx, postgresInstance.ID())
		assert.Error(t, err)
	})

	// Doc example: DROP POSTGRES INSTANCE IF EXISTS my_postgres;
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

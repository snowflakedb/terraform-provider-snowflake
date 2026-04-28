//go:build account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_PostgresInstances(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// ==================
	// Create
	// ==================

	// Doc example: CREATE POSTGRES INSTANCE my_postgres COMPUTE_FAMILY = 'STANDARD_S' STORAGE_SIZE_GB = 50 AUTHENTICATION_AUTHORITY = POSTGRES;
	t.Run("create - basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_1", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres)

		err := client.PostgresInstances.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, id))

		postgresInstance, err := client.PostgresInstances.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstanceFromObject(t, postgresInstance).
			HasName(id.Name()).
			HasComputeFamily("STANDARD_1").
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
	//   NETWORK_POLICY = 'my_network_policy' POSTGRES_SETTINGS = '{"postgres:work_mem" = "128MB"}'
	//   COMMENT = 'Production Postgres instance';
	t.Run("create - complete", func(t *testing.T) {
		networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicy(t)
		t.Cleanup(networkPolicyCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		comment := random.Comment()
		request := sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_1", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres).
			WithPostgresVersion(17).
			WithHighAvailability(true).
			WithNetworkPolicy(networkPolicy.Name).
			WithPostgresSettings(`{"postgres:work_mem" = "128MB"}`).
			WithComment(comment)

		err := client.PostgresInstances.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, id))

		postgresInstance, err := client.PostgresInstances.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstanceFromObject(t, postgresInstance).
			HasName(id.Name()).
			HasComputeFamily("STANDARD_1").
			HasStorageSize(10).
			HasAuthenticationAuthority("POSTGRES").
			HasPostgresVersion("17").
			HasIsHa(true).
			HasComment(comment).
			HasPostgresSettings(`{"postgres:work_mem" = "128MB"}`).
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty(),
		)
	})

	// Doc example: TAG (tag1 = 'value1')
	t.Run("create - with tags", func(t *testing.T) {
		tag1, tag1Cleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tag1Cleanup)
		tag2, tag2Cleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tag2Cleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_1", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres).
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
		request := sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_1", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgresOrSnowflake)

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
		request := sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_1", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres).
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

		forkId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewForkPostgresInstanceRequest(forkId, sourceInstance.ID())

		err := client.PostgresInstances.Fork(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, forkId))

		forkedInstance, err := client.PostgresInstances.ShowByID(ctx, forkId)
		require.NoError(t, err)
		assert.Equal(t, forkId.Name(), forkedInstance.Name)

		// Forked instances should have an origin referencing the source
		assert.NotNil(t, forkedInstance.Origin)
		if forkedInstance.Origin != nil {
			assert.Contains(t, *forkedInstance.Origin, sourceInstance.Name)
		}

		// Verify type field is populated
		assert.NotEmpty(t, forkedInstance.Type)
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
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		sourceInstance, sourceCleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(sourceCleanup)

		comment := random.Comment()
		forkId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewForkPostgresInstanceRequest(forkId, sourceInstance.ID()).
			WithComputeFamily("STANDARD_1").
			WithStorageSizeGb(20).
			WithPostgresSettings(`{"postgres:work_mem" = "128MB"}`).
			WithComment(comment).
			WithTag([]sdk.TagAssociation{
				{
					Name:  tag.ID(),
					Value: "fork_tag_value",
				},
			})

		err := client.PostgresInstances.Fork(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, forkId))

		forkedInstance, err := client.PostgresInstances.ShowByID(ctx, forkId)
		require.NoError(t, err)
		assert.Equal(t, forkId.Name(), forkedInstance.Name)
		assert.Equal(t, comment, *forkedInstance.Comment)

		assertTagSet(t, tag.ID(), forkId, sdk.ObjectTypePostgresInstance, "fork_tag_value")
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
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		// Set and unset comment
		comment := random.Comment()
		err := client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().WithComment(comment)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstance(t, postgresInstance.ID()).
			HasComment(comment),
		)

		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithUnset(*sdk.NewPostgresInstanceUnsetRequest().WithComment(true)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstance(t, postgresInstance.ID()).
			HasNoComment(),
		)

		// Set and unset maintenance_window_start
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().WithMaintenanceWindowStart(3)))
		require.NoError(t, err)

		properties, err := client.PostgresInstances.Describe(ctx, postgresInstance.ID())
		require.NoError(t, err)
		propertyMap := make(map[string]string)
		for _, p := range properties {
			propertyMap[p.Property] = p.Value
		}
		assert.Equal(t, "3", propertyMap["maintenance_window_start"])

		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithUnset(*sdk.NewPostgresInstanceUnsetRequest().WithMaintenanceWindowStart(true)))
		require.NoError(t, err)

		// Set and unset postgres_settings
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().WithPostgresSettings(`{"postgres:work_mem" = "128MB"}`)))
		require.NoError(t, err)

		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithUnset(*sdk.NewPostgresInstanceUnsetRequest().WithPostgresSettings(true)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstance(t, postgresInstance.ID()).
			HasNoPostgresSettings(),
		)

		// Set storage_size
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().
				WithStorageSizeGb(20)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstance(t, postgresInstance.ID()).
			HasStorageSize(20),
		)

		// Set compute_family
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().
				WithComputeFamily("STANDARD_2")))
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstance(t, postgresInstance.ID()).
			HasComputeFamily("STANDARD_2"),
		)

		// Set high_availability
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().
				WithHighAvailability(true)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstance(t, postgresInstance.ID()).
			HasIsHa(true),
		)

		// Set authentication_authority
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().
				WithAuthenticationAuthority(sdk.PostgresInstanceAuthenticationAuthorityPostgresOrSnowflake)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstance(t, postgresInstance.ID()).
			HasAuthenticationAuthority("POSTGRES_OR_SNOWFLAKE"),
		)

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

		// Set multiple properties in one call
		comment = random.Comment()
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().
				WithComment(comment).
				WithStorageSizeGb(50).
				WithPostgresSettings(`{"postgres:work_mem" = "128MB"}`)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstance(t, postgresInstance.ID()).
			HasComment(comment).
			HasStorageSize(50),
		)
	})

	t.Run("alter: suspend and resume", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		// Resume without suspending - instance is already in a running state
		err := client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithResume(true))
		assert.Error(t, err)

		// Suspend
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSuspend(true))
		require.NoError(t, err)

		result, err := client.PostgresInstances.ShowByID(ctx, postgresInstance.ID())
		require.NoError(t, err)
		// State may be SUSPENDING or SUSPENDED depending on timing
		assert.Contains(t, []sdk.PostgresInstanceState{
			sdk.PostgresInstanceStateSuspending,
			sdk.PostgresInstanceStateSuspended,
		}, result.State)

		// Suspend again - expect error due to invalid state
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSuspend(true))
		assert.Error(t, err)

		// Resume
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithResume(true))
		require.NoError(t, err)

		result, err = client.PostgresInstances.ShowByID(ctx, postgresInstance.ID())
		require.NoError(t, err)
		// State may be RESUMING, STARTING, CREATING, or READY depending on timing
		assert.Contains(t, []sdk.PostgresInstanceState{
			sdk.PostgresInstanceStateResuming,
			sdk.PostgresInstanceStateStarting,
			sdk.PostgresInstanceStateCreating,
			sdk.PostgresInstanceStateReady,
		}, result.State)
	})

	// Doc example: ALTER POSTGRES INSTANCE my_postgres RENAME TO prod_postgres;
	t.Run("alter: rename", func(t *testing.T) {
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

	t.Run("alter: set and unset storage_integration and network_policy", func(t *testing.T) {
		awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
		awsRoleARN := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)

		storageIntegration, storageIntegrationCleanup := testClientHelper().StorageIntegration.CreateS3(t, awsBucketUrl, awsRoleARN)
		t.Cleanup(storageIntegrationCleanup)

		networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicy(t)
		t.Cleanup(networkPolicyCleanup)

		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		// Set and unset storage_integration
		err := client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().WithStorageIntegration(storageIntegration.Name)))
		require.NoError(t, err)

		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithUnset(*sdk.NewPostgresInstanceUnsetRequest().WithStorageIntegration(true)))
		require.NoError(t, err)

		// Set and unset network_policy
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().WithNetworkPolicy(networkPolicy.Name)))
		require.NoError(t, err)

		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithUnset(*sdk.NewPostgresInstanceUnsetRequest().WithNetworkPolicy(true)))
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
			HasComputeFamily("STANDARD_1").
			HasStorageSize(10).
			HasAuthenticationAuthority("POSTGRES").
			HasIsHa(false).
			HasNoComment().
			HasNoOrigin().
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty(),
		)

		// RetentionTime should have a default value
		assert.GreaterOrEqual(t, result.RetentionTime, 0)

		// PostgresVersion should be a non-empty string
		assert.NotEmpty(t, result.PostgresVersion)

		// State should be one of the valid states
		assert.Contains(t, sdk.AllPostgresInstanceStates, result.State)
	})

	t.Run("ShowByID and ShowByIDSafely", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		// ShowByID
		result, err := client.PostgresInstances.ShowByID(ctx, postgresInstance.ID())
		require.NoError(t, err)
		assert.Equal(t, postgresInstance.ID(), result.ID())
		assert.Equal(t, postgresInstance.Name, result.Name)

		// ShowByIDSafely
		result, err = client.PostgresInstances.ShowByIDSafely(ctx, postgresInstance.ID())
		assert.NotNil(t, result)
		require.NoError(t, err)
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
		assert.Contains(t, propertyMap, "retention_time")

		assert.Equal(t, postgresInstance.ID().Name(), propertyMap["name"])
		assert.Equal(t, "STANDARD_1", propertyMap["compute_family"])
		assert.Equal(t, "POSTGRES", propertyMap["authentication_authority"])
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

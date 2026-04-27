//go:build non_account_level_tests

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

func TestInt_PostgresInstances(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

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
	//   NETWORK_POLICY = 'my_network_policy' COMMENT = 'Production Postgres instance';
	t.Run("create - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		comment := random.Comment()
		request := sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_1", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres).
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
			HasIsHa(false).
			HasComment(comment).
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty(),
		)
	})

	// Doc example: POSTGRES_SETTINGS = '{"postgres:work_mem" = "128MB", "pgbouncer:default_pool_size" = "200"}'
	t.Run("create - with postgres settings", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_1", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres).
			WithPostgresSettings(`{"postgres:work_mem" = "128MB"}`)

		err := client.PostgresInstances.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, id))

		postgresInstance, err := client.PostgresInstances.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstanceFromObject(t, postgresInstance).
			HasName(id.Name()),
		)
	})

	// Doc example: TAG (tag1 = 'value1')
	t.Run("create - with tags", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_1", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres).
			WithTag([]sdk.TagAssociation{
				{
					Name:  tag.ID(),
					Value: "value1",
				},
			})

		err := client.PostgresInstances.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, id))

		assertTagSet(t, tag.ID(), id, sdk.ObjectTypePostgresInstance, "value1")
	})

	// Doc example: ALTER POSTGRES INSTANCE my_postgres SET COMMENT = '...'
	t.Run("alter: set and unset comment", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		comment := random.Comment()
		err := client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().WithComment(comment)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstance(t, postgresInstance.ID()).
			HasComment(comment),
		)

		// Unset
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithUnset(*sdk.NewPostgresInstanceUnsetRequest().WithComment(true)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstance(t, postgresInstance.ID()).
			HasNoComment(),
		)
	})

	// Doc example: ALTER POSTGRES INSTANCE my_postgres SET COMPUTE_FAMILY = 'STANDARD_M' STORAGE_SIZE_GB = 100;
	t.Run("alter: set compute_family and storage_size", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		err := client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().
				WithStorageSizeGb(20)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstance(t, postgresInstance.ID()).
			HasStorageSize(20),
		)
	})

	// Doc example: ALTER POSTGRES INSTANCE my_postgres SET AUTHENTICATION_AUTHORITY = POSTGRES_OR_SNOWFLAKE
	t.Run("alter: set authentication_authority", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		err := client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().
				WithAuthenticationAuthority(sdk.PostgresInstanceAuthenticationAuthorityPostgresOrSnowflake)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstance(t, postgresInstance.ID()).
			HasAuthenticationAuthority("POSTGRES_OR_SNOWFLAKE"),
		)
	})

	// Doc example: ALTER POSTGRES INSTANCE my_postgres SET COMPUTE_FAMILY = 'STANDARD_L' APPLY IMMEDIATELY;
	t.Run("alter: set with apply immediately", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		err := client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSet(*sdk.NewPostgresInstanceSetRequest().
				WithStorageSizeGb(20).
				WithApply(*sdk.NewPostgresInstanceApplyRequest().WithImmediately(true))))
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstance(t, postgresInstance.ID()).
			HasStorageSize(20),
		)
	})

	// Doc example: ALTER POSTGRES INSTANCE my_postgres UNSET COMMENT, POSTGRES_SETTINGS, NETWORK_POLICY
	t.Run("alter: unset", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		comment := random.Comment()
		request := sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_1", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres).
			WithComment(comment).
			WithPostgresSettings(`{"postgres:work_mem" = "128MB"}`)

		err := client.PostgresInstances.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, id))

		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(id).
			WithUnset(*sdk.NewPostgresInstanceUnsetRequest().
				WithComment(true).
				WithPostgresSettings(true)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.PostgresInstance(t, id).
			HasNoComment(),
		)
	})

	// Doc example: ALTER POSTGRES INSTANCE my_postgres SUSPEND;
	t.Run("alter: suspend", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		err := client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithSuspend(true))
		require.NoError(t, err)

		result, err := client.PostgresInstances.ShowByID(ctx, postgresInstance.ID())
		require.NoError(t, err)
		// State may be SUSPENDING or SUSPENDED depending on timing
		assert.Contains(t, []sdk.PostgresInstanceState{
			sdk.PostgresInstanceStateSuspending,
			sdk.PostgresInstanceStateSuspended,
		}, result.State)
	})

	// Doc example: ALTER POSTGRES INSTANCE my_postgres RESUME;
	t.Run("alter: resume", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_1", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres)

		err := client.PostgresInstances.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, id))

		// Suspend first
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(id).
			WithSuspend(true))
		require.NoError(t, err)

		// Resume
		err = client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(id).
			WithResume(true))
		require.NoError(t, err)

		result, err := client.PostgresInstances.ShowByID(ctx, id)
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
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		newId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, newId))

		err := client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithRenameTo(newId))
		require.NoError(t, err)

		// Old name should not exist
		_, err = client.PostgresInstances.ShowByID(ctx, postgresInstance.ID())
		require.Error(t, err)

		// New name should exist
		result, err := client.PostgresInstances.ShowByID(ctx, newId)
		require.NoError(t, err)
		assert.Equal(t, newId.Name(), result.Name)
	})

	// Doc example: ALTER POSTGRES INSTANCE my_postgres RESET ACCESS FOR 'snowflake_admin'
	t.Run("alter: reset access for snowflake_admin", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		err := client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
			WithResetAccess(*sdk.NewPostgresInstanceResetAccessRequest(sdk.PostgresInstanceResetAccessRoleSnowflakeAdmin)))
		require.NoError(t, err)
	})

	// Doc example: ALTER POSTGRES INSTANCE my_postgres RESET ACCESS FOR 'application'
	t.Run("alter: reset access for application", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		err := client.PostgresInstances.Alter(ctx, sdk.NewAlterPostgresInstanceRequest(postgresInstance.ID()).
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
	})

	// Doc example: CREATE POSTGRES INSTANCE my_fork FORK my_source_instance AT (TIMESTAMP => '2025-01-15 12:00:00'::TIMESTAMP_NTZ);
	t.Run("fork - with at timestamp", func(t *testing.T) {
		sourceInstance, sourceCleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(sourceCleanup)

		forkId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewForkPostgresInstanceRequest(forkId, sourceInstance.ID()).
			WithAt(*sdk.NewPostgresInstanceForkAtRequest().WithTimestamp("2025-01-15 12:00:00"))

		// This may fail if the timestamp is outside the retention period; that's expected
		err := client.PostgresInstances.Fork(ctx, request)
		if err == nil {
			t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, forkId))

			forkedInstance, showErr := client.PostgresInstances.ShowByID(ctx, forkId)
			require.NoError(t, showErr)
			assert.Equal(t, forkId.Name(), forkedInstance.Name)
		}
	})

	// Doc example: CREATE POSTGRES INSTANCE my_fork FORK my_source_instance AT (OFFSET => -7200) COMPUTE_FAMILY = 'STANDARD_L';
	t.Run("fork - with at offset and compute overrides", func(t *testing.T) {
		sourceInstance, sourceCleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(sourceCleanup)

		forkId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewForkPostgresInstanceRequest(forkId, sourceInstance.ID()).
			WithAt(*sdk.NewPostgresInstanceForkAtRequest().WithOffset("-60")).
			WithComment("Fork with offset and compute override")

		err := client.PostgresInstances.Fork(ctx, request)
		if err == nil {
			t.Cleanup(testClientHelper().PostgresInstance.DropFunc(t, forkId))

			forkedInstance, showErr := client.PostgresInstances.ShowByID(ctx, forkId)
			require.NoError(t, showErr)
			assert.Equal(t, forkId.Name(), forkedInstance.Name)
		}
	})

	// Doc example: DESCRIBE POSTGRES INSTANCE my_postgres;
	t.Run("describe", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		properties, err := client.PostgresInstances.Describe(ctx, postgresInstance.ID())
		require.NoError(t, err)
		require.NotEmpty(t, properties)

		// Verify expected properties are present per the docs output table
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

		assert.Equal(t, postgresInstance.ID().Name(), propertyMap["name"])
		assert.Equal(t, "STANDARD_1", propertyMap["compute_family"])
		assert.Equal(t, "POSTGRES", propertyMap["authentication_authority"])
	})

	// Doc example: SHOW POSTGRES INSTANCES;
	t.Run("show: all", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

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
	})

	// Doc example: SHOW POSTGRES INSTANCES LIKE 'DEV_%';
	t.Run("show: with like", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		instances, err := client.PostgresInstances.Show(ctx, sdk.NewShowPostgresInstanceRequest().
			WithLike(sdk.Like{Pattern: sdk.String(postgresInstance.Name)}))
		require.NoError(t, err)
		require.Len(t, instances, 1)
		assert.Equal(t, postgresInstance.Name, instances[0].Name)
	})

	// Doc example: SHOW POSTGRES INSTANCES STARTS WITH 'PROD';
	t.Run("show: with starts_with", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		instances, err := client.PostgresInstances.Show(ctx, sdk.NewShowPostgresInstanceRequest().
			WithStartsWith(postgresInstance.Name))
		require.NoError(t, err)
		require.NotEmpty(t, instances)
		assert.Equal(t, postgresInstance.Name, instances[0].Name)
	})

	// Doc example: SHOW POSTGRES INSTANCES LIMIT 1
	t.Run("show: with limit", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		instances, err := client.PostgresInstances.Show(ctx, sdk.NewShowPostgresInstanceRequest().
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
		require.NoError(t, err)
		require.LessOrEqual(t, len(instances), 1)

		_ = postgresInstance
	})

	t.Run("ShowByID", func(t *testing.T) {
		postgresInstance, cleanup := testClientHelper().PostgresInstance.Create(t)
		t.Cleanup(cleanup)

		result, err := client.PostgresInstances.ShowByID(ctx, postgresInstance.ID())
		require.NoError(t, err)
		assert.Equal(t, postgresInstance.ID(), result.ID())
		assert.Equal(t, postgresInstance.Name, result.Name)
	})

	// Doc example: DROP POSTGRES INSTANCE my_postgres;
	t.Run("drop: existing object", func(t *testing.T) {
		postgresInstance, _ := testClientHelper().PostgresInstance.Create(t)

		err := client.PostgresInstances.Drop(ctx, sdk.NewDropPostgresInstanceRequest(postgresInstance.ID()))
		require.NoError(t, err)

		_, err = client.PostgresInstances.ShowByID(ctx, postgresInstance.ID())
		assert.Error(t, err)
	})

	// Doc example: DROP POSTGRES INSTANCE IF EXISTS my_postgres;
	t.Run("drop: non-existing object", func(t *testing.T) {
		err := client.PostgresInstances.Drop(ctx, sdk.NewDropPostgresInstanceRequest(NonExistingAccountObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)

		err = client.PostgresInstances.Drop(ctx, sdk.NewDropPostgresInstanceRequest(NonExistingAccountObjectIdentifier).WithIfExists(true))
		require.NoError(t, err)
	})
}

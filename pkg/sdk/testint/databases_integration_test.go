//go:build account_level_tests

package testint

import (
	"fmt"
	"testing"
	"time"

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

func TestInt_DatabasesCreate(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("minimal", func(t *testing.T) {
		databaseID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Databases.Create(ctx, sdk.NewCreateDatabaseRequest(databaseID).WithOrReplace(true))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, databaseID))

		database, err := client.Databases.ShowByID(ctx, databaseID)
		require.NoError(t, err)
		assert.Equal(t, databaseID.Name(), database.Name)
		assertThatObject(t, objectassert.DatabaseFromObject(t, database).
			HasKind(sdk.DatabaseKindStandard))
	})

	t.Run("as clone", func(t *testing.T) {
		cloneDatabase, cloneDatabaseCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(cloneDatabaseCleanup)

		databaseID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Databases.Clone(ctx, sdk.NewCloneDatabaseRequest(databaseID).WithClone(sdk.Clone{
			SourceObject: cloneDatabase.ID(),
			At: &sdk.TimeTravel{
				Offset: sdk.Int(0),
			},
		}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, databaseID))

		database, err := client.Databases.ShowByID(ctx, databaseID)
		require.NoError(t, err)
		assert.Equal(t, databaseID.Name(), database.Name)
	})

	t.Run("complete", func(t *testing.T) {
		databaseId := testClientHelper().Ids.RandomAccountObjectIdentifier()

		// new database and schema created on purpose
		databaseTest, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseCleanup)

		schemaTest, schemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, databaseTest.ID())
		t.Cleanup(schemaCleanup)

		tagTest, tagCleanup := testClientHelper().Tag.CreateTagInSchema(t, schemaTest.ID())
		t.Cleanup(tagCleanup)

		tag2Test, tag2Cleanup := testClientHelper().Tag.CreateTagInSchema(t, schemaTest.ID())
		t.Cleanup(tag2Cleanup)

		externalVolume, externalVolumeCleanup := testClientHelper().ExternalVolume.Create(t)
		t.Cleanup(externalVolumeCleanup)

		catalog, catalogCleanup := testClientHelper().CatalogIntegration.Create(t)
		t.Cleanup(catalogCleanup)

		comment := random.Comment()
		err := client.Databases.Create(
			ctx, sdk.NewCreateDatabaseRequest(databaseId).
				WithTransient(true).
				WithIfNotExists(true).
				WithDataRetentionTimeInDays(0).
				WithMaxDataExtensionTimeInDays(10).
				WithExternalVolume(externalVolume).
				WithCatalog(catalog).
				WithReplaceInvalidCharacters(true).
				WithDefaultDdlCollation("en_US").
				WithDefaultNotebookComputePoolCpu("CPU_X64_S").
				WithDefaultNotebookComputePoolGpu("GPU_NV_S").
				WithStorageSerializationPolicy(sdk.StorageSerializationPolicyCompatible).
				WithLogLevel(sdk.LogLevelInfo).
				WithLogEventLevel(sdk.LogLevelInfo).
				WithTraceLevel(sdk.TraceLevelPropagate).
				WithSuspendTaskAfterNumFailures(10).
				WithTaskAutoRetryAttempts(10).
				WithUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeMedium).
				WithUserTaskTimeoutMs(12_000).
				WithUserTaskMinimumTriggerIntervalInSeconds(30).
				WithQuotedIdentifiersIgnoreCase(true).
				WithEnableConsoleOutput(true).
				WithComment(comment).
				WithTag([]sdk.TagAssociation{
					{
						Name:  tagTest.ID(),
						Value: "v1",
					},
					{
						Name:  tag2Test.ID(),
						Value: "v2",
					},
				}),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, databaseId))

		database, err := client.Databases.ShowByID(ctx, databaseId)
		require.NoError(t, err)
		assert.Equal(t, databaseId.Name(), database.Name)
		assert.Equal(t, comment, database.Comment)

		params, err := client.Databases.ShowParameters(ctx, databaseId)
		require.NoError(t, err)
		assertParameterEquals := func(t *testing.T, parameterName sdk.AccountParameter, expected string) {
			t.Helper()
			assert.Equal(t, expected, helpers.FindParameter(t, params, parameterName).Value)
		}

		assertParameterEquals(t, sdk.AccountParameterDataRetentionTimeInDays, "0")
		assertParameterEquals(t, sdk.AccountParameterMaxDataExtensionTimeInDays, "10")
		assertParameterEquals(t, sdk.AccountParameterDefaultDdlCollation, "en_US")
		assertParameterEquals(t, sdk.AccountParameterDefaultNotebookComputePoolCpu, "CPU_X64_S")
		assertParameterEquals(t, sdk.AccountParameterDefaultNotebookComputePoolGpu, "GPU_NV_S")
		assertParameterEquals(t, sdk.AccountParameterExternalVolume, externalVolume.Name())
		assertParameterEquals(t, sdk.AccountParameterCatalog, catalog.Name())
		assertParameterEquals(t, sdk.AccountParameterLogLevel, string(sdk.LogLevelInfo))
		assertParameterEquals(t, sdk.AccountParameterLogEventLevel, string(sdk.LogLevelInfo))
		assertParameterEquals(t, sdk.AccountParameterTraceLevel, string(sdk.TraceLevelPropagate))
		assertParameterEquals(t, sdk.AccountParameterReplaceInvalidCharacters, "true")
		assertParameterEquals(t, sdk.AccountParameterStorageSerializationPolicy, string(sdk.StorageSerializationPolicyCompatible))
		assertParameterEquals(t, sdk.AccountParameterSuspendTaskAfterNumFailures, "10")
		assertParameterEquals(t, sdk.AccountParameterTaskAutoRetryAttempts, "10")
		assertParameterEquals(t, sdk.AccountParameterUserTaskManagedInitialWarehouseSize, string(sdk.WarehouseSizeMedium))
		assertParameterEquals(t, sdk.AccountParameterUserTaskTimeoutMs, "12000")
		assertParameterEquals(t, sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds, "30")
		assertParameterEquals(t, sdk.AccountParameterQuotedIdentifiersIgnoreCase, "true")
		assertParameterEquals(t, sdk.AccountParameterEnableConsoleOutput, "true")

		tag1Value, err := client.SystemFunctions.GetTag(ctx, tagTest.ID(), database.ID(), sdk.ObjectTypeDatabase)
		require.NoError(t, err)
		assert.Equal(t, sdk.Pointer("v1"), tag1Value)

		tag2Value, err := client.SystemFunctions.GetTag(ctx, tag2Test.ID(), database.ID(), sdk.ObjectTypeDatabase)
		require.NoError(t, err)
		assert.Equal(t, sdk.Pointer("v2"), tag2Value)
	})
}

func TestInt_DatabasesCreateShared(t *testing.T) {
	client := testClient(t)
	secondaryClient := testSecondaryClient(t)
	ctx := testContext(t)

	testTag, testTagCleanup := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(testTagCleanup)

	externalVolume, externalVolumeCleanup := testClientHelper().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalog, catalogCleanup := testClientHelper().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	// prepare a database on the secondary account
	shareTest, shareCleanup := secondaryTestClientHelper().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	sharedDatabase, sharedDatabaseCleanup := secondaryTestClientHelper().Database.CreateDatabase(t)
	t.Cleanup(sharedDatabaseCleanup)

	databaseId := sharedDatabase.ID()

	err := secondaryClient.Grants.GrantPrivilegeToShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
		Database: sharedDatabase.ID(),
	}, shareTest.ID())
	require.NoError(t, err)
	t.Cleanup(func() {
		err := secondaryClient.Grants.RevokePrivilegeFromShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
			Database: sharedDatabase.ID(),
		}, shareTest.ID())
		require.NoError(t, err)
	})

	err = secondaryClient.Shares.Alter(ctx, sdk.NewAlterShareRequest(shareTest.ID()).WithIfExists(true).WithSet(sdk.ShareSetRequest{
		Accounts: []sdk.AccountIdentifier{
			testClientHelper().Account.GetAccountIdentifier(t),
		},
	}))
	require.NoError(t, err)

	comment := random.Comment()
	err = client.Databases.CreateShared(
		ctx, sdk.NewCreateSharedDatabaseRequest(databaseId, shareTest.ExternalID()).
			WithTransient(true).
			WithIfNotExists(true).
			WithExternalVolume(externalVolume).
			WithCatalog(catalog).
			WithLogLevel(sdk.LogLevelDebug).
			WithLogEventLevel(sdk.LogLevelDebug).
			WithTraceLevel(sdk.TraceLevelAlways).
			WithReplaceInvalidCharacters(true).
			WithDefaultDdlCollation("en_US").
			WithDefaultNotebookComputePoolCpu("CPU_X64_S").
			WithDefaultNotebookComputePoolGpu("GPU_NV_S").
			WithStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
			WithSuspendTaskAfterNumFailures(10).
			WithTaskAutoRetryAttempts(10).
			WithUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeMedium).
			WithUserTaskTimeoutMs(12_000).
			WithUserTaskMinimumTriggerIntervalInSeconds(30).
			WithQuotedIdentifiersIgnoreCase(true).
			WithEnableConsoleOutput(true).
			WithComment(comment).
			WithTag([]sdk.TagAssociation{
				{
					Name:  testTag.ID(),
					Value: "v1",
				},
			}),
	)
	require.NoError(t, err)
	t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, databaseId))

	database, err := client.Databases.ShowByID(ctx, databaseId)
	require.NoError(t, err)

	assert.Equal(t, databaseId.Name(), database.Name)
	assert.Equal(t, comment, database.Comment)
	assertThatObject(t, objectassert.DatabaseFromObject(t, database).
		HasKind(sdk.DatabaseKindImportedDatabase))

	params, err := client.Databases.ShowParameters(ctx, databaseId)
	require.NoError(t, err)
	assertParameterEquals := func(t *testing.T, parameterName sdk.AccountParameter, expected string) {
		t.Helper()
		assert.Equal(t, expected, helpers.FindParameter(t, params, parameterName).Value)
	}

	assertParameterEquals(t, sdk.AccountParameterDefaultDdlCollation, "en_US")
	assertParameterEquals(t, sdk.AccountParameterDefaultNotebookComputePoolCpu, "CPU_X64_S")
	assertParameterEquals(t, sdk.AccountParameterDefaultNotebookComputePoolGpu, "GPU_NV_S")
	assertParameterEquals(t, sdk.AccountParameterExternalVolume, externalVolume.Name())
	assertParameterEquals(t, sdk.AccountParameterCatalog, catalog.Name())
	assertParameterEquals(t, sdk.AccountParameterLogLevel, string(sdk.LogLevelDebug))
	assertParameterEquals(t, sdk.AccountParameterLogEventLevel, string(sdk.LogLevelDebug))
	assertParameterEquals(t, sdk.AccountParameterTraceLevel, string(sdk.TraceLevelAlways))
	assertParameterEquals(t, sdk.AccountParameterReplaceInvalidCharacters, "true")
	assertParameterEquals(t, sdk.AccountParameterStorageSerializationPolicy, string(sdk.StorageSerializationPolicyOptimized))
	assertParameterEquals(t, sdk.AccountParameterSuspendTaskAfterNumFailures, "10")
	assertParameterEquals(t, sdk.AccountParameterTaskAutoRetryAttempts, "10")
	assertParameterEquals(t, sdk.AccountParameterUserTaskManagedInitialWarehouseSize, string(sdk.WarehouseSizeMedium))
	assertParameterEquals(t, sdk.AccountParameterUserTaskTimeoutMs, "12000")
	assertParameterEquals(t, sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds, "30")
	assertParameterEquals(t, sdk.AccountParameterQuotedIdentifiersIgnoreCase, "true")
	assertParameterEquals(t, sdk.AccountParameterEnableConsoleOutput, "true")

	tag1Value, err := client.SystemFunctions.GetTag(ctx, testTag.ID(), database.ID(), sdk.ObjectTypeDatabase)
	require.NoError(t, err)
	assert.Equal(t, sdk.Pointer("v1"), tag1Value)
}

func TestInt_DatabasesCreateSecondary(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	primaryDatabase, externalDatabaseId := createPrimaryDatabase(t)
	databaseId := primaryDatabase.ID()

	externalVolume, externalVolumeCleanup := testClientHelper().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalog, catalogCleanup := testClientHelper().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	comment := random.Comment()
	err := client.Databases.CreateSecondary(
		ctx, sdk.NewCreateSecondaryDatabaseRequest(databaseId, externalDatabaseId).
			WithIfNotExists(true).
			WithDataRetentionTimeInDays(10).
			WithMaxDataExtensionTimeInDays(10).
			WithExternalVolume(externalVolume).
			WithCatalog(catalog).
			WithReplaceInvalidCharacters(true).
			WithDefaultDdlCollation("en_US").
			WithDefaultNotebookComputePoolCpu("CPU_X64_S").
			WithDefaultNotebookComputePoolGpu("GPU_NV_S").
			WithStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
			WithLogLevel(sdk.LogLevelDebug).
			WithLogEventLevel(sdk.LogLevelDebug).
			WithTraceLevel(sdk.TraceLevelAlways).
			WithSuspendTaskAfterNumFailures(10).
			WithTaskAutoRetryAttempts(10).
			WithUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeMedium).
			WithUserTaskTimeoutMs(12_000).
			WithUserTaskMinimumTriggerIntervalInSeconds(30).
			WithQuotedIdentifiersIgnoreCase(true).
			WithEnableConsoleOutput(true).
			WithComment(comment),
	)
	require.NoError(t, err)
	t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, databaseId))

	database, err := client.Databases.ShowByID(ctx, databaseId)
	require.NoError(t, err)

	assert.Equal(t, databaseId.Name(), database.Name)
	assert.Equal(t, comment, database.Comment)
	// a secondary database reports the same kind as a regular one
	assertThatObject(t, objectassert.DatabaseFromObject(t, database).
		HasKind(sdk.DatabaseKindStandard))

	params, err := client.Databases.ShowParameters(ctx, databaseId)
	require.NoError(t, err)
	assertParameterEquals := func(t *testing.T, parameterName sdk.AccountParameter, expected string) {
		t.Helper()
		assert.Equal(t, expected, helpers.FindParameter(t, params, parameterName).Value)
	}

	assertParameterEquals(t, sdk.AccountParameterDataRetentionTimeInDays, "10")
	assertParameterEquals(t, sdk.AccountParameterMaxDataExtensionTimeInDays, "10")
	assertParameterEquals(t, sdk.AccountParameterDefaultDdlCollation, "en_US")
	assertParameterEquals(t, sdk.AccountParameterDefaultNotebookComputePoolCpu, "CPU_X64_S")
	assertParameterEquals(t, sdk.AccountParameterDefaultNotebookComputePoolGpu, "GPU_NV_S")
	assertParameterEquals(t, sdk.AccountParameterExternalVolume, externalVolume.Name())
	assertParameterEquals(t, sdk.AccountParameterCatalog, catalog.Name())
	assertParameterEquals(t, sdk.AccountParameterLogLevel, string(sdk.LogLevelDebug))
	assertParameterEquals(t, sdk.AccountParameterLogEventLevel, string(sdk.LogLevelDebug))
	assertParameterEquals(t, sdk.AccountParameterTraceLevel, string(sdk.TraceLevelAlways))
	assertParameterEquals(t, sdk.AccountParameterReplaceInvalidCharacters, "true")
	assertParameterEquals(t, sdk.AccountParameterStorageSerializationPolicy, string(sdk.StorageSerializationPolicyOptimized))
	assertParameterEquals(t, sdk.AccountParameterSuspendTaskAfterNumFailures, "10")
	assertParameterEquals(t, sdk.AccountParameterTaskAutoRetryAttempts, "10")
	assertParameterEquals(t, sdk.AccountParameterUserTaskManagedInitialWarehouseSize, string(sdk.WarehouseSizeMedium))
	assertParameterEquals(t, sdk.AccountParameterUserTaskTimeoutMs, "12000")
	assertParameterEquals(t, sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds, "30")
	assertParameterEquals(t, sdk.AccountParameterQuotedIdentifiersIgnoreCase, "true")
	assertParameterEquals(t, sdk.AccountParameterEnableConsoleOutput, "true")
}

func TestInt_DatabasesCreateFromListing(t *testing.T) {
	t.Skip("TODO(SNOW-3556777): Use precreated listing")

	client := testClient(t)
	ctx := testContext(t)

	secondaryClient := testSecondaryClient(t)
	secondaryCtx := testSecondaryContext(t)

	share, shareCleanup := secondaryTestClientHelper().Share.CreateShare(t)
	t.Cleanup(shareCleanup)
	t.Cleanup(secondaryTestClientHelper().Grant.GrantPrivilegeOnDatabaseToShare(t, secondaryTestClientHelper().Ids.DatabaseId(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))

	primaryAccountId := testClientHelper().Context.CurrentAccountId(t)
	manifest, _ := secondaryTestClientHelper().Listing.BasicManifestWithTargetAccounts(t, primaryAccountId)

	listingId := secondaryTestClientHelper().Ids.RandomAccountObjectIdentifier()
	err := secondaryClient.Listings.Create(secondaryCtx, sdk.NewCreateListingRequest(listingId).
		WithAs(manifest).
		WithWith(*sdk.NewListingWithRequest().WithShare(share.ID())).
		WithReview(true).
		WithPublish(true))
	require.NoError(t, err)
	t.Cleanup(secondaryTestClientHelper().Listing.DropFunc(t, listingId))

	listing, err := secondaryClient.Listings.ShowByID(secondaryCtx, listingId)
	require.NoError(t, err)
	require.NotEmpty(t, listing.GlobalName)
	require.Equal(t, sdk.ListingStatePublished, listing.State)

	testClientHelper().Listing.AcceptLegalTermsWithRetry(t, listing.GlobalName, time.Minute, 5*time.Second)
	testClientHelper().Listing.RequestListingAndWaitForSuccess(t, listing.GlobalName, 10)

	t.Run("basic", func(t *testing.T) {
		databaseID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Databases.CreateFromListing(ctx, sdk.NewCreateFromListingDatabaseRequest(databaseID, listing.GlobalName))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, databaseID))

		database, err := client.Databases.ShowByID(ctx, databaseID)
		require.NoError(t, err)
		assert.Equal(t, databaseID.Name(), database.Name)
		assertThatObject(t, objectassert.DatabaseFromObject(t, database).
			HasKind(sdk.DatabaseKindImportedDatabase))
	})
}

func TestInt_DatabasesAlter(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	assertDatabaseParameterEquals := func(t *testing.T, params []*sdk.Parameter, parameterName sdk.AccountParameter, expected string) {
		t.Helper()
		assert.Equal(t, expected, helpers.FindParameter(t, params, parameterName).Value)
	}

	assertDatabaseParameterEqualsToDefaultValue := func(t *testing.T, params []*sdk.Parameter, parameterName sdk.ObjectParameter) {
		t.Helper()
		param, err := collections.FindFirst(params, func(param *sdk.Parameter) bool { return param.Key == string(parameterName) })
		assert.NoError(t, err)
		assert.NotNil(t, param)
		if param != nil && (*param).Level == "" {
			param := *param
			assert.Equal(t, param.Default, param.Value)
		}
	}

	testCases := []struct {
		DatabaseType string
		CreateFn     func(t *testing.T) (*sdk.Database, func())
		ExpectedKind sdk.DatabaseKind
	}{
		{
			DatabaseType: "Normal",
			CreateFn: func(t *testing.T) (*sdk.Database, func()) {
				t.Helper()
				return testClientHelper().Database.CreateDatabase(t)
			},
			ExpectedKind: sdk.DatabaseKindStandard,
		},
		{
			DatabaseType: "From Share",
			CreateFn:     createDatabaseFromShare,
			ExpectedKind: sdk.DatabaseKindImportedDatabase,
		},
		{
			DatabaseType: "Replica",
			CreateFn:     createDatabaseReplica,
			ExpectedKind: sdk.DatabaseKindStandard,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Database: %s - renaming", testCase.DatabaseType), func(t *testing.T) {
			databaseTest, databaseTestCleanup := testCase.CreateFn(t)
			t.Cleanup(databaseTestCleanup)
			newName := testClientHelper().Ids.RandomAccountObjectIdentifier()

			err := client.Databases.Alter(ctx, sdk.NewAlterDatabaseRequest(databaseTest.ID()).WithRenameTo(newName))
			require.NoError(t, err)
			t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, newName))

			database, err := client.Databases.ShowByID(ctx, newName)
			require.NoError(t, err)
			assert.Equal(t, newName.Name(), database.Name)
			assertThatObject(t, objectassert.DatabaseFromObject(t, database).
				HasKind(testCase.ExpectedKind))
		})

		t.Run(fmt.Sprintf("Database: %s - setting and unsetting parameters", testCase.DatabaseType), func(t *testing.T) {
			if testCase.DatabaseType == "From Share" {
				t.Skipf("Skipping database test because from share is not supported")
			}

			databaseTest, databaseTestCleanup := testCase.CreateFn(t)
			t.Cleanup(databaseTestCleanup)

			externalVolumeTest, externalVolumeTestCleanup := testClientHelper().ExternalVolume.Create(t)
			t.Cleanup(externalVolumeTestCleanup)

			catalogIntegrationTest, catalogIntegrationTestCleanup := testClientHelper().CatalogIntegration.Create(t)
			t.Cleanup(catalogIntegrationTestCleanup)

			err := client.Databases.Alter(ctx, sdk.NewAlterDatabaseRequest(databaseTest.ID()).WithSet(
				*sdk.NewDatabaseSetRequest().
					WithDataRetentionTimeInDays(42).
					WithMaxDataExtensionTimeInDays(42).
					WithExternalVolume(externalVolumeTest).
					WithCatalog(catalogIntegrationTest).
					WithReplaceInvalidCharacters(true).
					WithDefaultDdlCollation("en_US").
					WithDefaultNotebookComputePoolCpu("CPU_X64_S").
					WithDefaultNotebookComputePoolGpu("GPU_NV_S").
					WithStorageSerializationPolicy(sdk.StorageSerializationPolicyCompatible).
					WithLogLevel(sdk.LogLevelInfo).
					WithLogEventLevel(sdk.LogLevelInfo).
					WithTraceLevel(sdk.TraceLevelPropagate).
					WithSuspendTaskAfterNumFailures(10).
					WithTaskAutoRetryAttempts(10).
					WithUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeMedium).
					WithUserTaskTimeoutMs(12_000).
					WithUserTaskMinimumTriggerIntervalInSeconds(30).
					WithQuotedIdentifiersIgnoreCase(true).
					WithEnableConsoleOutput(true),
			))
			require.NoError(t, err)

			params, err := client.Databases.ShowParameters(ctx, databaseTest.ID())
			require.NoError(t, err)
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterDataRetentionTimeInDays, "42")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterMaxDataExtensionTimeInDays, "42")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterExternalVolume, externalVolumeTest.Name())
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterCatalog, catalogIntegrationTest.Name())
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterReplaceInvalidCharacters, "true")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterDefaultDdlCollation, "en_US")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterDefaultNotebookComputePoolCpu, "CPU_X64_S")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterDefaultNotebookComputePoolGpu, "GPU_NV_S")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterStorageSerializationPolicy, string(sdk.StorageSerializationPolicyCompatible))
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterLogLevel, string(sdk.LogLevelInfo))
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterLogEventLevel, string(sdk.LogLevelInfo))
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterTraceLevel, string(sdk.TraceLevelPropagate))
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterSuspendTaskAfterNumFailures, "10")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterTaskAutoRetryAttempts, "10")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterUserTaskManagedInitialWarehouseSize, string(sdk.WarehouseSizeMedium))
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterUserTaskTimeoutMs, "12000")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds, "30")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterQuotedIdentifiersIgnoreCase, "true")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterEnableConsoleOutput, "true")

			err = client.Databases.Alter(ctx, sdk.NewAlterDatabaseRequest(databaseTest.ID()).WithUnset(
				*sdk.NewDatabaseUnsetRequest().
					WithDataRetentionTimeInDays(true).
					WithMaxDataExtensionTimeInDays(true).
					WithExternalVolume(true).
					WithCatalog(true).
					WithReplaceInvalidCharacters(true).
					WithDefaultDdlCollation(true).
					WithDefaultNotebookComputePoolCpu(true).
					WithDefaultNotebookComputePoolGpu(true).
					WithStorageSerializationPolicy(true).
					WithLogLevel(true).
					WithLogEventLevel(true).
					WithTraceLevel(true).
					WithSuspendTaskAfterNumFailures(true).
					WithTaskAutoRetryAttempts(true).
					WithUserTaskManagedInitialWarehouseSize(true).
					WithUserTaskTimeoutMs(true).
					WithUserTaskMinimumTriggerIntervalInSeconds(true).
					WithQuotedIdentifiersIgnoreCase(true).
					WithEnableConsoleOutput(true),
			))
			require.NoError(t, err)

			params, err = client.Databases.ShowParameters(ctx, databaseTest.ID())
			require.NoError(t, err)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterDataRetentionTimeInDays)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterMaxDataExtensionTimeInDays)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterExternalVolume)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterCatalog)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterReplaceInvalidCharacters)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterDefaultDdlCollation)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterDefaultNotebookComputePoolCpu)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterDefaultNotebookComputePoolGpu)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterStorageSerializationPolicy)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterLogLevel)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterLogEventLevel)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterTraceLevel)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterSuspendTaskAfterNumFailures)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterTaskAutoRetryAttempts)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterUserTaskManagedInitialWarehouseSize)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterUserTaskTimeoutMs)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterQuotedIdentifiersIgnoreCase)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterEnableConsoleOutput)
		})

		t.Run(fmt.Sprintf("Database: %s - setting and unsetting comment", testCase.DatabaseType), func(t *testing.T) {
			databaseTest, databaseTestCleanup := testCase.CreateFn(t)
			t.Cleanup(databaseTestCleanup)

			err := client.Databases.Alter(ctx, sdk.NewAlterDatabaseRequest(databaseTest.ID()).WithSet(
				*sdk.NewDatabaseSetRequest().WithComment("test comment"),
			))
			require.NoError(t, err)

			database, err := client.Databases.ShowByID(ctx, databaseTest.ID())
			require.NoError(t, err)

			assert.Equal(t, "test comment", database.Comment)

			err = client.Databases.Alter(ctx, sdk.NewAlterDatabaseRequest(databaseTest.ID()).WithUnset(
				*sdk.NewDatabaseUnsetRequest().WithComment(true),
			))
			require.NoError(t, err)

			database, err = client.Databases.ShowByID(ctx, databaseTest.ID())
			require.NoError(t, err)
			assert.Equal(t, "", database.Comment)
		})

		t.Run(fmt.Sprintf("Database: %s - swap with another database", testCase.DatabaseType), func(t *testing.T) {
			databaseTest, databaseTestCleanup := testCase.CreateFn(t)
			t.Cleanup(databaseTestCleanup)

			databaseTest2, databaseCleanup2 := testClientHelper().Database.CreateDatabase(t)
			t.Cleanup(databaseCleanup2)

			err := client.Databases.Alter(ctx, sdk.NewAlterDatabaseRequest(databaseTest.ID()).WithSwapWith(databaseTest2.ID()))
			require.NoError(t, err)
		})
	}
}

func TestInt_DatabasesAlterReplication(t *testing.T) {
	t.Run("enable and disable replication", func(t *testing.T) {
		ctx := testContext(t)

		database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseCleanup)

		err := testClient(t).Databases.AlterReplication(
			ctx, sdk.NewAlterReplicationDatabaseRequest(database.ID()).
				WithEnableReplication(*sdk.NewEnableReplicationRequest().
					WithToAccounts([]sdk.AccountIdentifier{
						secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
					}).
					WithIgnoreEditionCheck(true)),
		)
		require.NoError(t, err)

		err = testClient(t).Databases.AlterReplication(
			ctx, sdk.NewAlterReplicationDatabaseRequest(database.ID()).
				WithDisableReplication(*sdk.NewDisableReplicationRequest().
					WithToAccounts([]sdk.AccountIdentifier{
						secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
					})),
		)
		require.NoError(t, err)
	})

	t.Run("refresh replicated database", func(t *testing.T) {
		client := testClient(t)
		secondaryClient := testSecondaryClient(t)
		ctx := testContext(t)

		sharedDatabase, externalDatabaseId := createPrimaryDatabase(t)

		externalVolume, externalVolumeCleanup := testClientHelper().ExternalVolume.Create(t)
		t.Cleanup(externalVolumeCleanup)

		catalog, catalogCleanup := testClientHelper().CatalogIntegration.Create(t)
		t.Cleanup(catalogCleanup)

		comment := random.Comment()
		err := client.Databases.CreateSecondary(
			ctx, sdk.NewCreateSecondaryDatabaseRequest(sharedDatabase.ID(), externalDatabaseId).
				WithIfNotExists(true).
				WithDataRetentionTimeInDays(1).
				WithMaxDataExtensionTimeInDays(10).
				WithExternalVolume(externalVolume).
				WithCatalog(catalog).
				WithDefaultDdlCollation("en_US").
				WithDefaultNotebookComputePoolCpu("CPU_X64_S").
				WithDefaultNotebookComputePoolGpu("GPU_NV_S").
				WithLogLevel(sdk.LogLevelDebug).
				WithTraceLevel(sdk.TraceLevelAlways).
				WithComment(comment),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, sharedDatabase.ID()))

		err = secondaryClient.Databases.Alter(ctx, sdk.NewAlterDatabaseRequest(sharedDatabase.ID()).WithSet(
			*sdk.NewDatabaseSetRequest().WithComment("some comment"),
		))
		require.NoError(t, err)

		database, err := client.Databases.ShowByID(ctx, sharedDatabase.ID())
		require.NoError(t, err)

		assert.Equal(t, sharedDatabase.ID().Name(), database.Name)
		assert.Equal(t, 1, database.RetentionTime)
		assert.Equal(t, comment, database.Comment)

		err = client.Databases.AlterReplication(ctx, sdk.NewAlterReplicationDatabaseRequest(sharedDatabase.ID()).WithRefresh(true))
		require.NoError(t, err)

		database, err = client.Databases.ShowByID(ctx, sharedDatabase.ID())
		require.NoError(t, err)

		assert.Equal(t, sharedDatabase.ID().Name(), database.Name)
		assert.Equal(t, 1, database.RetentionTime)
		assert.Equal(t, comment, database.Comment)
	})
}

func TestInt_DatabasesAlterFailover(t *testing.T) {
	t.Run("enable and disable failover", func(t *testing.T) {
		ctx := testContext(t)

		database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseCleanup)

		err := testClient(t).Databases.AlterReplication(
			ctx, sdk.NewAlterReplicationDatabaseRequest(database.ID()).
				WithEnableReplication(*sdk.NewEnableReplicationRequest().
					WithToAccounts([]sdk.AccountIdentifier{
						secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
					}).
					WithIgnoreEditionCheck(true)),
		)
		require.NoError(t, err)

		err = testClient(t).Databases.AlterFailover(
			ctx, sdk.NewAlterFailoverDatabaseRequest(database.ID()).
				WithEnableFailover(*sdk.NewEnableFailoverRequest().
					WithToAccounts([]sdk.AccountIdentifier{
						secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
					})),
		)
		require.NoError(t, err)

		err = testClient(t).Databases.AlterFailover(
			ctx, sdk.NewAlterFailoverDatabaseRequest(database.ID()).
				WithDisableFailover(*sdk.NewDisableFailoverRequest().
					WithToAccounts([]sdk.AccountIdentifier{
						secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
					})),
		)
		require.NoError(t, err)
	})

	t.Run("promote to primary", func(t *testing.T) {
		t.Skipf("Can be unskipped after [SNOW-1002023]. CI Snowflake Edition doesn't support this feature")

		ctx := testContext(t)

		database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseCleanup)

		err := testClient(t).Databases.AlterReplication(
			ctx, sdk.NewAlterReplicationDatabaseRequest(database.ID()).
				WithEnableReplication(*sdk.NewEnableReplicationRequest().
					WithToAccounts([]sdk.AccountIdentifier{
						secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
					}).
					WithIgnoreEditionCheck(true)),
		)
		require.NoError(t, err)

		err = testClient(t).Databases.AlterFailover(
			ctx, sdk.NewAlterFailoverDatabaseRequest(database.ID()).
				WithEnableFailover(*sdk.NewEnableFailoverRequest().
					WithToAccounts([]sdk.AccountIdentifier{
						secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
					})),
		)
		require.NoError(t, err)

		err = testClient(t).Databases.AlterFailover(ctx, sdk.NewAlterFailoverDatabaseRequest(database.ID()).WithPrimary(true))
		require.NoError(t, err)
	})
}

func TestInt_DatabasesDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("drop with nil options", func(t *testing.T) {
		databaseTest, databaseTestCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseTestCleanup)

		err := client.Databases.Drop(ctx, sdk.NewDropDatabaseRequest(databaseTest.ID()))
		require.NoError(t, err)
	})

	t.Run("drop if exists", func(t *testing.T) {
		databaseTest, databaseTestCleanup := testClientHelper().Database.CreateDatabase(t)
		databaseTestCleanup()

		err := client.Databases.Drop(ctx, sdk.NewDropDatabaseRequest(databaseTest.ID()).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("drop with cascade", func(t *testing.T) {
		databaseTest, databaseTestCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseTestCleanup)

		err := client.Databases.Drop(ctx, sdk.NewDropDatabaseRequest(databaseTest.ID()).
			WithIfExists(true).
			WithCascade(true))
		require.NoError(t, err)
	})

	t.Run("drop with restrict", func(t *testing.T) {
		databaseTest, databaseTestCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseTestCleanup)

		err := client.Databases.Drop(ctx, sdk.NewDropDatabaseRequest(databaseTest.ID()).
			WithIfExists(true).
			WithRestrict(true))
		require.NoError(t, err)
	})
}

func TestInt_DatabasesUndrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
	databaseCleanup()

	_, err := client.Databases.ShowByID(ctx, databaseTest.ID())
	require.Error(t, err)

	err = client.Databases.Undrop(ctx, sdk.NewUndropDatabaseRequest(databaseTest.ID()))
	require.NoError(t, err)

	database, err := client.Databases.ShowByID(ctx, databaseTest.ID())
	require.NoError(t, err)

	assert.Equal(t, databaseTest.Name, database.Name)
}

func TestInt_DatabasesShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	databaseTest2, databaseCleanup2 := testClientHelper().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup2)

	t.Run("without show options", func(t *testing.T) {
		databases, err := client.Databases.Show(ctx, sdk.NewShowDatabaseRequest())
		require.NoError(t, err)

		assert.GreaterOrEqual(t, len(databases), 2)
		databaseIDs := make([]sdk.AccountObjectIdentifier, len(databases))
		for i, database := range databases {
			databaseIDs[i] = database.ID()
		}
		assert.Contains(t, databaseIDs, databaseTest.ID())
		assert.Contains(t, databaseIDs, databaseTest2.ID())
		assert.Equal(t, "ROLE", databases[0].OwnerRoleType)

		database, err := collections.FindFirst(databases, func(database sdk.Database) bool { return database.Name == databaseTest.Name })
		require.NoError(t, err)
		assertThatObject(t, objectassert.DatabaseFromObject(t, database).
			HasKind(sdk.DatabaseKindStandard))
	})

	t.Run("with terse", func(t *testing.T) {
		databases, err := client.Databases.Show(
			ctx, sdk.NewShowDatabaseRequest().
				WithTerse(true).
				WithLike(sdk.Like{Pattern: sdk.String(databaseTest.Name)}),
		)
		require.NoError(t, err)

		database, err := collections.FindFirst(databases, func(database sdk.Database) bool { return database.Name == databaseTest.Name })
		require.NoError(t, err)

		assert.Equal(t, databaseTest.Name, database.Name)
		assert.NotEmpty(t, database.CreatedOn)
		assert.Empty(t, database.DroppedOn)
		assert.Empty(t, database.Owner)
		// kind is returned for TERSE, even though the docs claim it's always NULL there
		assertThatObject(t, objectassert.DatabaseFromObject(t, database).
			HasKind(sdk.DatabaseKindStandard))
	})

	t.Run("with history", func(t *testing.T) {
		databaseTest3, databaseCleanup3 := testClientHelper().Database.CreateDatabase(t)
		databaseCleanup3()

		databases, err := client.Databases.Show(
			ctx, sdk.NewShowDatabaseRequest().
				WithHistory(true).
				WithLike(sdk.Like{Pattern: sdk.String(databaseTest3.Name)}),
		)
		require.NoError(t, err)

		droppedDatabase, err := collections.FindFirst(databases, func(database sdk.Database) bool { return database.Name == databaseTest3.Name })
		require.NoError(t, err)

		assert.Equal(t, databaseTest3.Name, droppedDatabase.Name)
		assert.NotEmpty(t, droppedDatabase.DroppedOn)
	})

	t.Run("with like starts with", func(t *testing.T) {
		databases, err := client.Databases.Show(
			ctx, sdk.NewShowDatabaseRequest().
				WithStartsWith(databaseTest.Name).
				WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}),
		)
		require.NoError(t, err)

		database, err := collections.FindFirst(databases, func(database sdk.Database) bool { return database.Name == databaseTest.Name })
		require.NoError(t, err)

		assert.Equal(t, databaseTest.Name, database.Name)
	})

	t.Run("when searching a non-existent database", func(t *testing.T) {
		databases, err := client.Databases.Show(
			ctx, sdk.NewShowDatabaseRequest().
				WithLike(sdk.Like{Pattern: sdk.String("non-existent")}),
		)
		require.NoError(t, err)

		assert.Equal(t, 0, len(databases))
	})

	t.Run("show by id: missing database", func(t *testing.T) {
		_, err := client.Databases.ShowByID(ctx, testClientHelper().Ids.RandomAccountObjectIdentifier())
		require.Error(t, err)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("show by id safely", func(t *testing.T) {
		database, err := client.Databases.ShowByIDSafely(ctx, testClientHelper().Ids.DatabaseId())
		assert.NotNil(t, database)
		require.NoError(t, err)
	})

	t.Run("show by id safely: missing database", func(t *testing.T) {
		_, err := client.Databases.ShowByIDSafely(ctx, testClientHelper().Ids.RandomAccountObjectIdentifier())
		require.Error(t, err)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})
}

func TestInt_DatabasesDescribe(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	assertContainsSchema := func(details *sdk.DatabaseDetails, schemaName string) {
		_, err := collections.FindFirst(details.Rows, func(row sdk.DatabaseDetailsRow) bool { return row.Kind == "SCHEMA" && row.Name == schemaName })
		assert.NoError(t, err)
	}

	schemaTest, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	databaseDetails, err := client.Databases.Describe(ctx, schemaTest.ID().DatabaseId())
	require.NoError(t, err)

	assertContainsSchema(databaseDetails, schemaTest.ID().Name())
	assertContainsSchema(databaseDetails, "INFORMATION_SCHEMA")
	assertContainsSchema(databaseDetails, "PUBLIC")
}

// TestInt_DatabasesCatalogLinked_WithAdditionalDependencies requires a preconfigured external Iceberg
// REST catalog integration and external volume, which we can't provision dynamically for now.
// TODO(SNOW-3725859): Provide them dynamically and move these tests to the main test suite.
func TestInt_DatabasesCatalogLinked_WithAdditionalDependencies(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	restCatalogId := sdk.NewAccountObjectIdentifier(testenvs.GetOrSkipTest(t, testenvs.CatalogLinkedDatabaseCatalogIntegration))
	externalVolumeId := sdk.NewAccountObjectIdentifier(testenvs.GetOrSkipTest(t, testenvs.CatalogLinkedDatabaseExternalVolume))

	createCatalogLinkedDatabase := func(t *testing.T) *sdk.Database {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		database, databaseCleanup := testClientHelper().Database.CreateCatalogLinkedDatabaseWithRequest(t, sdk.NewCreateCatalogLinkedDatabaseRequest(id).
			WithLinkedCatalog(*sdk.NewLinkedCatalogRequest().WithCatalog(restCatalogId)).
			WithExternalVolume(externalVolumeId))
		t.Cleanup(databaseCleanup)
		return database
	}

	getConfig := func(t *testing.T, id sdk.AccountObjectIdentifier) *sdk.CatalogLinkedDatabaseConfig {
		t.Helper()
		config, err := client.SystemFunctions.GetCatalogLinkedDatabaseConfig(ctx, id)
		require.NoError(t, err)
		require.NotNil(t, config)
		return config
	}

	// Shared database for the subtests that only read it or fully revert their changes.
	sharedDatabase := createCatalogLinkedDatabase(t)

	t.Run("create: basic", func(t *testing.T) {
		assertThatObject(t, objectassert.DatabaseFromObject(t, sharedDatabase).
			HasCreatedOnNotEmpty().
			HasName(sharedDatabase.ID().Name()).
			HasIsDefault(false).
			HasOwner("ACCOUNTADMIN").
			HasComment("").
			HasTransient(false).
			HasKind(sdk.DatabaseKindCatalogLinkedDatabase).
			HasOwnerRoleType("ROLE"))

		assertThatObject(t, objectparametersassert.DatabaseParameters(t, sharedDatabase.ID()).
			HasExternalVolume(externalVolumeId.Name()))
	})

	t.Run("create: complete", func(t *testing.T) {
		// new database and schema created on purpose
		databaseTest, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseCleanup)

		schemaTest, schemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, databaseTest.ID())
		t.Cleanup(schemaCleanup)

		tagTest, tagCleanup := testClientHelper().Tag.CreateTagInSchema(t, schemaTest.ID())
		t.Cleanup(tagCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		comment := random.Comment()
		database, cleanup := testClientHelper().Database.CreateCatalogLinkedDatabaseWithRequest(t, sdk.NewCreateCatalogLinkedDatabaseRequest(id).
			WithLinkedCatalog(*sdk.NewLinkedCatalogRequest().
				WithCatalog(restCatalogId).
				WithAllowedNamespaces([]sdk.StringListItemWrapper{{Value: "ns1"}, {Value: "ns2"}}).
				WithBlockedNamespaces([]sdk.StringListItemWrapper{{Value: "ns3"}}).
				WithAllowedWriteOperations(sdk.CatalogLinkedDatabaseAllowedWriteOperationsAll).
				WithNamespaceMode(sdk.CatalogLinkedDatabaseNamespaceModeFlattenNestedNamespace).
				WithNamespaceFlattenDelimiter("_").
				WithSyncIntervalSeconds(60)).
			WithExternalVolume(externalVolumeId).
			WithComment(comment).
			WithTag([]sdk.TagAssociation{
				{
					Name:  tagTest.ID(),
					Value: "v1",
				},
			}).
			WithCatalogCaseSensitivity(sdk.DatabaseCatalogCaseSensitivityCaseInsensitive))
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.DatabaseFromObject(t, database).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasIsDefault(false).
			HasOwner("ACCOUNTADMIN").
			HasComment(comment).
			HasTransient(false).
			HasKind(sdk.DatabaseKindCatalogLinkedDatabase).
			HasOwnerRoleType("ROLE"))

		assertThatObject(t, objectassert.CatalogLinkedDatabaseConfigFromObject(t, getConfig(t, database.ID())).
			HasCatalogIntegration(restCatalogId).
			HasNoCatalogName().
			HasExternalVolume(externalVolumeId).
			HasSyncIntervalSeconds(60).
			HasNamespaceMode(sdk.CatalogLinkedDatabaseNamespaceModeFlattenNestedNamespace).
			HasNamespaceFlattenDelimiter("_").
			HasAllowedWriteOperations(sdk.CatalogLinkedDatabaseAllowedWriteOperationsAll).
			HasCatalogCaseSensitivity(sdk.DatabaseCatalogCaseSensitivityCaseInsensitive).
			HasIsSuspended(false).
			HasAllowedNamespaces("ns1", "ns2").
			HasBlockedNamespaces("ns3"))

		assertTagSet(t, tagTest.ID(), database.ID(), sdk.ObjectTypeDatabase, "v1")
	})

	t.Run("get catalog linked database config", func(t *testing.T) {
		// The values not set at creation are asserted with their documented defaults; they may need
		// adjustment after a run against a live external catalog.
		assertThatObject(t, objectassert.CatalogLinkedDatabaseConfigFromObject(t, getConfig(t, sharedDatabase.ID())).
			HasCatalogIntegration(restCatalogId).
			HasNoCatalogName().
			HasExternalVolume(externalVolumeId).
			HasSyncIntervalSeconds(30).
			HasNamespaceMode(sdk.CatalogLinkedDatabaseNamespaceModeIgnoreNestedNamespace).
			HasNoNamespaceFlattenDelimiter().
			HasAllowedWriteOperations(sdk.CatalogLinkedDatabaseAllowedWriteOperationsAll).
			HasCatalogCaseSensitivity(sdk.DatabaseCatalogCaseSensitivityCaseInsensitive).
			HasIsSuspended(false).
			HasNoAllowedNamespaces().
			HasNoBlockedNamespaces())
	})

	t.Run("alter: allowed namespaces", func(t *testing.T) {
		err := client.Databases.AlterCatalogLinked(ctx, sdk.NewAlterCatalogLinkedDatabaseRequest(sharedDatabase.ID()).
			WithIfExists(true).
			WithAddToAllowedNamespaces(*sdk.NewAddToAllowedNamespacesRequest([]sdk.StringListItemWrapper{{Value: "ns1"}, {Value: "ns2"}})))
		require.NoError(t, err)
		assertThatObject(t, objectassert.CatalogLinkedDatabaseConfigFromObject(t, getConfig(t, sharedDatabase.ID())).
			HasAllowedNamespaces("ns1", "ns2"))

		err = client.Databases.AlterCatalogLinked(ctx, sdk.NewAlterCatalogLinkedDatabaseRequest(sharedDatabase.ID()).
			WithRemoveFromAllowedNamespaces(*sdk.NewRemoveFromAllowedNamespacesRequest([]sdk.StringListItemWrapper{{Value: "ns2"}})))
		require.NoError(t, err)
		assertThatObject(t, objectassert.CatalogLinkedDatabaseConfigFromObject(t, getConfig(t, sharedDatabase.ID())).
			HasAllowedNamespaces("ns1"))

		err = client.Databases.AlterCatalogLinked(ctx, sdk.NewAlterCatalogLinkedDatabaseRequest(sharedDatabase.ID()).
			WithUnsetAllowedNamespaces(true))
		require.NoError(t, err)
		assertThatObject(t, objectassert.CatalogLinkedDatabaseConfigFromObject(t, getConfig(t, sharedDatabase.ID())).
			HasNoAllowedNamespaces())
	})

	t.Run("alter: blocked namespaces", func(t *testing.T) {
		err := client.Databases.AlterCatalogLinked(ctx, sdk.NewAlterCatalogLinkedDatabaseRequest(sharedDatabase.ID()).
			WithAddToBlockedNamespaces(*sdk.NewAddToBlockedNamespacesRequest([]sdk.StringListItemWrapper{{Value: "ns3"}, {Value: "ns4"}})))
		require.NoError(t, err)
		assertThatObject(t, objectassert.CatalogLinkedDatabaseConfigFromObject(t, getConfig(t, sharedDatabase.ID())).
			HasBlockedNamespaces("ns3", "ns4"))

		err = client.Databases.AlterCatalogLinked(ctx, sdk.NewAlterCatalogLinkedDatabaseRequest(sharedDatabase.ID()).
			WithRemoveFromBlockedNamespaces(*sdk.NewRemoveFromBlockedNamespacesRequest([]sdk.StringListItemWrapper{{Value: "ns4"}})))
		require.NoError(t, err)
		assertThatObject(t, objectassert.CatalogLinkedDatabaseConfigFromObject(t, getConfig(t, sharedDatabase.ID())).
			HasBlockedNamespaces("ns3"))

		err = client.Databases.AlterCatalogLinked(ctx, sdk.NewAlterCatalogLinkedDatabaseRequest(sharedDatabase.ID()).
			WithUnsetBlockedNamespaces(true))
		require.NoError(t, err)
		assertThatObject(t, objectassert.CatalogLinkedDatabaseConfigFromObject(t, getConfig(t, sharedDatabase.ID())).
			HasNoBlockedNamespaces())
	})

	t.Run("alter: set linked catalog parameters", func(t *testing.T) {
		database := createCatalogLinkedDatabase(t)

		err := client.Databases.AlterCatalogLinked(ctx, sdk.NewAlterCatalogLinkedDatabaseRequest(database.ID()).
			WithSet(*sdk.NewCatalogLinkedDatabaseSetRequest().
				WithSyncIntervalSeconds(120).
				WithAllowedWriteOperations(sdk.CatalogLinkedDatabaseAllowedWriteOperationsNone)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.CatalogLinkedDatabaseConfigFromObject(t, getConfig(t, database.ID())).
			HasSyncIntervalSeconds(120).
			HasAllowedWriteOperations(sdk.CatalogLinkedDatabaseAllowedWriteOperationsNone))
	})

	t.Run("alter: rename and comment through the generic alter", func(t *testing.T) {
		database := createCatalogLinkedDatabase(t)

		newName := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Databases.Alter(ctx, sdk.NewAlterDatabaseRequest(database.ID()).WithRenameTo(newName))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, newName))

		comment := random.Comment()
		err = client.Databases.Alter(ctx, sdk.NewAlterDatabaseRequest(newName).WithSet(*sdk.NewDatabaseSetRequest().WithComment(comment)))
		require.NoError(t, err)

		renamed, err := client.Databases.ShowByID(ctx, newName)
		require.NoError(t, err)
		assertThatObject(t, objectassert.DatabaseFromObject(t, renamed).
			HasName(newName.Name()).
			HasComment(comment).
			HasKind(sdk.DatabaseKindCatalogLinkedDatabase))

		err = client.Databases.Alter(ctx, sdk.NewAlterDatabaseRequest(newName).WithUnset(*sdk.NewDatabaseUnsetRequest().WithComment(true)))
		require.NoError(t, err)

		unset, err := client.Databases.ShowByID(ctx, newName)
		require.NoError(t, err)
		assertThatObject(t, objectassert.DatabaseFromObject(t, unset).
			HasComment(""))
	})

	t.Run("get catalog link status", func(t *testing.T) {
		// The execution state may not be reported in the short window right after creation, so poll briefly.
		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			status, err := client.SystemFunctions.GetCatalogLinkStatus(ctx, sharedDatabase.ID())
			assert.NoError(collect, err)
			if assert.NotNil(collect, status) {
				assert.NotEmpty(collect, status.ExecutionState)
			}
		}, 20*time.Second, 2*time.Second)

		// The failure-related fields are state- and timing-dependent (the docs show failureDetails
		// populated even for a RUNNING link), so only the execution state is asserted.
		status, err := client.SystemFunctions.GetCatalogLinkStatus(ctx, sharedDatabase.ID())
		require.NoError(t, err)
		assertThatObject(t, objectassert.CatalogLinkStatusFromObject(t, status).
			HasExecutionStateNotEmpty())
	})
}

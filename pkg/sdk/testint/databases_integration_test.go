//go:build account_level_tests

package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
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
		err := client.Databases.Create(ctx, databaseID, &sdk.CreateDatabaseOptions{
			OrReplace: sdk.Bool(true),
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, databaseID))

		database, err := client.Databases.ShowByID(ctx, databaseID)
		require.NoError(t, err)
		assert.Equal(t, databaseID.Name(), database.Name)
	})

	t.Run("as clone", func(t *testing.T) {
		cloneDatabase, cloneDatabaseCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(cloneDatabaseCleanup)

		databaseID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Databases.Create(ctx, databaseID, &sdk.CreateDatabaseOptions{
			Clone: &sdk.Clone{
				SourceObject: cloneDatabase.ID(),
				At: &sdk.TimeTravel{
					Offset: sdk.Int(0),
				},
			},
		})
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
		err := client.Databases.Create(ctx, databaseId, &sdk.CreateDatabaseOptions{
			Transient:                               sdk.Bool(true),
			IfNotExists:                             sdk.Bool(true),
			DataRetentionTimeInDays:                 sdk.Int(0),
			MaxDataExtensionTimeInDays:              sdk.Int(10),
			ExternalVolume:                          &externalVolume,
			Catalog:                                 &catalog,
			ReplaceInvalidCharacters:                sdk.Bool(true),
			DefaultDDLCollation:                     sdk.String("en_US"),
			StorageSerializationPolicy:              sdk.Pointer(sdk.StorageSerializationPolicyCompatible),
			LogLevel:                                sdk.Pointer(sdk.LogLevelInfo),
			TraceLevel:                              sdk.Pointer(sdk.TraceLevelPropagate),
			SuspendTaskAfterNumFailures:             sdk.Int(10),
			TaskAutoRetryAttempts:                   sdk.Int(10),
			UserTaskManagedInitialWarehouseSize:     sdk.Pointer(sdk.WarehouseSizeMedium),
			UserTaskTimeoutMs:                       sdk.Int(12_000),
			UserTaskMinimumTriggerIntervalInSeconds: sdk.Int(30),
			QuotedIdentifiersIgnoreCase:             sdk.Bool(true),
			EnableConsoleOutput:                     sdk.Bool(true),
			Comment:                                 sdk.String(comment),
			Tag: []sdk.TagAssociation{
				{
					Name:  tagTest.ID(),
					Value: "v1",
				},
				{
					Name:  tag2Test.ID(),
					Value: "v2",
				},
			},
		})
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
		assertParameterEquals(t, sdk.AccountParameterDefaultDDLCollation, "en_US")
		assertParameterEquals(t, sdk.AccountParameterExternalVolume, externalVolume.Name())
		assertParameterEquals(t, sdk.AccountParameterCatalog, catalog.Name())
		assertParameterEquals(t, sdk.AccountParameterLogLevel, string(sdk.LogLevelInfo))
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

	err = secondaryClient.Shares.Alter(ctx, shareTest.ID(), &sdk.AlterShareOptions{
		IfExists: sdk.Bool(true),
		Set: &sdk.ShareSet{
			Accounts: []sdk.AccountIdentifier{
				testClientHelper().Account.GetAccountIdentifier(t),
			},
		},
	})
	require.NoError(t, err)

	comment := random.Comment()
	err = client.Databases.CreateShared(ctx, databaseId, shareTest.ExternalID(), &sdk.CreateSharedDatabaseOptions{
		Transient:                               sdk.Bool(true),
		IfNotExists:                             sdk.Bool(true),
		ExternalVolume:                          &externalVolume,
		Catalog:                                 &catalog,
		LogLevel:                                sdk.Pointer(sdk.LogLevelDebug),
		TraceLevel:                              sdk.Pointer(sdk.TraceLevelAlways),
		ReplaceInvalidCharacters:                sdk.Bool(true),
		DefaultDDLCollation:                     sdk.String("en_US"),
		StorageSerializationPolicy:              sdk.Pointer(sdk.StorageSerializationPolicyOptimized),
		SuspendTaskAfterNumFailures:             sdk.Int(10),
		TaskAutoRetryAttempts:                   sdk.Int(10),
		UserTaskManagedInitialWarehouseSize:     sdk.Pointer(sdk.WarehouseSizeMedium),
		UserTaskTimeoutMs:                       sdk.Int(12_000),
		UserTaskMinimumTriggerIntervalInSeconds: sdk.Int(30),
		QuotedIdentifiersIgnoreCase:             sdk.Bool(true),
		EnableConsoleOutput:                     sdk.Bool(true),
		Comment:                                 sdk.String(comment),
		Tag: []sdk.TagAssociation{
			{
				Name:  testTag.ID(),
				Value: "v1",
			},
		},
	})
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

	assertParameterEquals(t, sdk.AccountParameterDefaultDDLCollation, "en_US")
	assertParameterEquals(t, sdk.AccountParameterExternalVolume, externalVolume.Name())
	assertParameterEquals(t, sdk.AccountParameterCatalog, catalog.Name())
	assertParameterEquals(t, sdk.AccountParameterLogLevel, string(sdk.LogLevelDebug))
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
	secondaryClient := testSecondaryClient(t)
	ctx := testContext(t)

	sharedDatabase, sharedDatabaseCleanup := secondaryTestClientHelper().Database.CreateDatabase(t)
	t.Cleanup(sharedDatabaseCleanup)

	databaseId := sharedDatabase.ID()

	err := secondaryClient.Databases.AlterReplication(ctx, sharedDatabase.ID(), &sdk.AlterDatabaseReplicationOptions{
		EnableReplication: &sdk.EnableReplication{
			ToAccounts: []sdk.AccountIdentifier{
				testClientHelper().Account.GetAccountIdentifier(t),
			},
			IgnoreEditionCheck: sdk.Bool(true),
		},
	})
	require.NoError(t, err)

	externalVolume, externalVolumeCleanup := testClientHelper().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalog, catalogCleanup := testClientHelper().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	externalDatabaseId := sdk.NewExternalObjectIdentifier(secondaryTestClientHelper().Account.GetAccountIdentifier(t), sharedDatabase.ID())

	comment := random.Comment()
	err = client.Databases.CreateSecondary(ctx, databaseId, externalDatabaseId, &sdk.CreateSecondaryDatabaseOptions{
		IfNotExists:                             sdk.Bool(true),
		DataRetentionTimeInDays:                 sdk.Int(10),
		MaxDataExtensionTimeInDays:              sdk.Int(10),
		ExternalVolume:                          &externalVolume,
		Catalog:                                 &catalog,
		ReplaceInvalidCharacters:                sdk.Bool(true),
		DefaultDDLCollation:                     sdk.String("en_US"),
		StorageSerializationPolicy:              sdk.Pointer(sdk.StorageSerializationPolicyOptimized),
		LogLevel:                                sdk.Pointer(sdk.LogLevelDebug),
		TraceLevel:                              sdk.Pointer(sdk.TraceLevelAlways),
		SuspendTaskAfterNumFailures:             sdk.Int(10),
		TaskAutoRetryAttempts:                   sdk.Int(10),
		UserTaskManagedInitialWarehouseSize:     sdk.Pointer(sdk.WarehouseSizeMedium),
		UserTaskTimeoutMs:                       sdk.Int(12_000),
		UserTaskMinimumTriggerIntervalInSeconds: sdk.Int(30),
		QuotedIdentifiersIgnoreCase:             sdk.Bool(true),
		EnableConsoleOutput:                     sdk.Bool(true),
		Comment:                                 sdk.String(comment),
	})
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

	assertParameterEquals(t, sdk.AccountParameterDataRetentionTimeInDays, "10")
	assertParameterEquals(t, sdk.AccountParameterMaxDataExtensionTimeInDays, "10")
	assertParameterEquals(t, sdk.AccountParameterDefaultDDLCollation, "en_US")
	assertParameterEquals(t, sdk.AccountParameterExternalVolume, externalVolume.Name())
	assertParameterEquals(t, sdk.AccountParameterCatalog, catalog.Name())
	assertParameterEquals(t, sdk.AccountParameterLogLevel, string(sdk.LogLevelDebug))
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
	}{
		{
			DatabaseType: "Normal",
			CreateFn: func(t *testing.T) (*sdk.Database, func()) {
				t.Helper()
				return testClientHelper().Database.CreateDatabase(t)
			},
		},
		{
			DatabaseType: "From Share",
			CreateFn:     createDatabaseFromShare,
		},
		{
			DatabaseType: "Replica",
			CreateFn:     createDatabaseReplica,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Database: %s - renaming", testCase.DatabaseType), func(t *testing.T) {
			databaseTest, databaseTestCleanup := testCase.CreateFn(t)
			t.Cleanup(databaseTestCleanup)
			newName := testClientHelper().Ids.RandomAccountObjectIdentifier()

			err := client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
				NewName: &newName,
			})
			require.NoError(t, err)
			t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, newName))

			database, err := client.Databases.ShowByID(ctx, newName)
			require.NoError(t, err)
			assert.Equal(t, newName.Name(), database.Name)
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

			err := client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
				Set: &sdk.DatabaseSet{
					DataRetentionTimeInDays:                 sdk.Int(42),
					MaxDataExtensionTimeInDays:              sdk.Int(42),
					ExternalVolume:                          &externalVolumeTest,
					Catalog:                                 &catalogIntegrationTest,
					ReplaceInvalidCharacters:                sdk.Bool(true),
					DefaultDDLCollation:                     sdk.String("en_US"),
					StorageSerializationPolicy:              sdk.Pointer(sdk.StorageSerializationPolicyCompatible),
					LogLevel:                                sdk.Pointer(sdk.LogLevelInfo),
					TraceLevel:                              sdk.Pointer(sdk.TraceLevelPropagate),
					SuspendTaskAfterNumFailures:             sdk.Int(10),
					TaskAutoRetryAttempts:                   sdk.Int(10),
					UserTaskManagedInitialWarehouseSize:     sdk.Pointer(sdk.WarehouseSizeMedium),
					UserTaskTimeoutMs:                       sdk.Int(12_000),
					UserTaskMinimumTriggerIntervalInSeconds: sdk.Int(30),
					QuotedIdentifiersIgnoreCase:             sdk.Bool(true),
					EnableConsoleOutput:                     sdk.Bool(true),
				},
			})
			require.NoError(t, err)

			params, err := client.Databases.ShowParameters(ctx, databaseTest.ID())
			require.NoError(t, err)
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterDataRetentionTimeInDays, "42")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterMaxDataExtensionTimeInDays, "42")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterExternalVolume, externalVolumeTest.Name())
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterCatalog, catalogIntegrationTest.Name())
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterReplaceInvalidCharacters, "true")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterDefaultDDLCollation, "en_US")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterStorageSerializationPolicy, string(sdk.StorageSerializationPolicyCompatible))
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterLogLevel, string(sdk.LogLevelInfo))
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterTraceLevel, string(sdk.TraceLevelPropagate))
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterSuspendTaskAfterNumFailures, "10")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterTaskAutoRetryAttempts, "10")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterUserTaskManagedInitialWarehouseSize, string(sdk.WarehouseSizeMedium))
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterUserTaskTimeoutMs, "12000")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds, "30")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterQuotedIdentifiersIgnoreCase, "true")
			assertDatabaseParameterEquals(t, params, sdk.AccountParameterEnableConsoleOutput, "true")

			err = client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
				Unset: &sdk.DatabaseUnset{
					DataRetentionTimeInDays:                 sdk.Bool(true),
					MaxDataExtensionTimeInDays:              sdk.Bool(true),
					ExternalVolume:                          sdk.Bool(true),
					Catalog:                                 sdk.Bool(true),
					ReplaceInvalidCharacters:                sdk.Bool(true),
					DefaultDDLCollation:                     sdk.Bool(true),
					StorageSerializationPolicy:              sdk.Bool(true),
					LogLevel:                                sdk.Bool(true),
					TraceLevel:                              sdk.Bool(true),
					SuspendTaskAfterNumFailures:             sdk.Bool(true),
					TaskAutoRetryAttempts:                   sdk.Bool(true),
					UserTaskManagedInitialWarehouseSize:     sdk.Bool(true),
					UserTaskTimeoutMs:                       sdk.Bool(true),
					UserTaskMinimumTriggerIntervalInSeconds: sdk.Bool(true),
					QuotedIdentifiersIgnoreCase:             sdk.Bool(true),
					EnableConsoleOutput:                     sdk.Bool(true),
				},
			})
			require.NoError(t, err)

			params, err = client.Databases.ShowParameters(ctx, databaseTest.ID())
			require.NoError(t, err)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterDataRetentionTimeInDays)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterMaxDataExtensionTimeInDays)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterExternalVolume)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterCatalog)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterReplaceInvalidCharacters)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterDefaultDDLCollation)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterStorageSerializationPolicy)
			assertDatabaseParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterLogLevel)
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

			err := client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
				Set: &sdk.DatabaseSet{
					Comment: sdk.String("test comment"),
				},
			})
			require.NoError(t, err)

			database, err := client.Databases.ShowByID(ctx, databaseTest.ID())
			require.NoError(t, err)

			assert.Equal(t, "test comment", database.Comment)

			err = client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
				Unset: &sdk.DatabaseUnset{
					Comment: sdk.Bool(true),
				},
			})
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

			err := client.Databases.Alter(ctx, databaseTest.ID(), &sdk.AlterDatabaseOptions{
				SwapWith: sdk.Pointer(databaseTest2.ID()),
			})
			require.NoError(t, err)
		})
	}
}

func TestInt_DatabasesAlterReplication(t *testing.T) {
	t.Run("enable and disable replication", func(t *testing.T) {
		ctx := testContext(t)

		database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseCleanup)

		err := testClient(t).Databases.AlterReplication(ctx, database.ID(), &sdk.AlterDatabaseReplicationOptions{
			EnableReplication: &sdk.EnableReplication{
				ToAccounts: []sdk.AccountIdentifier{
					secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
				},
				IgnoreEditionCheck: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		err = testClient(t).Databases.AlterReplication(ctx, database.ID(), &sdk.AlterDatabaseReplicationOptions{
			DisableReplication: &sdk.DisableReplication{
				ToAccounts: []sdk.AccountIdentifier{
					secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
				},
			},
		})
		require.NoError(t, err)
	})

	t.Run("refresh replicated database", func(t *testing.T) {
		client := testClient(t)
		secondaryClient := testSecondaryClient(t)
		ctx := testContext(t)

		sharedDatabase, sharedDatabaseCleanup := secondaryTestClientHelper().Database.CreateDatabase(t)
		t.Cleanup(sharedDatabaseCleanup)

		err := secondaryClient.Databases.AlterReplication(ctx, sharedDatabase.ID(), &sdk.AlterDatabaseReplicationOptions{
			EnableReplication: &sdk.EnableReplication{
				ToAccounts: []sdk.AccountIdentifier{
					testClientHelper().Account.GetAccountIdentifier(t),
				},
				IgnoreEditionCheck: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		externalVolume, externalVolumeCleanup := testClientHelper().ExternalVolume.Create(t)
		t.Cleanup(externalVolumeCleanup)

		catalog, catalogCleanup := testClientHelper().CatalogIntegration.Create(t)
		t.Cleanup(catalogCleanup)

		externalDatabaseId := sdk.NewExternalObjectIdentifier(secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(), sharedDatabase.ID())
		comment := random.Comment()
		err = client.Databases.CreateSecondary(ctx, sharedDatabase.ID(), externalDatabaseId, &sdk.CreateSecondaryDatabaseOptions{
			IfNotExists:                sdk.Bool(true),
			DataRetentionTimeInDays:    sdk.Int(1),
			MaxDataExtensionTimeInDays: sdk.Int(10),
			ExternalVolume:             &externalVolume,
			Catalog:                    &catalog,
			DefaultDDLCollation:        sdk.String("en_US"),
			LogLevel:                   sdk.Pointer(sdk.LogLevelDebug),
			TraceLevel:                 sdk.Pointer(sdk.TraceLevelAlways),
			Comment:                    sdk.String(comment),
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, sharedDatabase.ID()))

		err = secondaryClient.Databases.Alter(ctx, sharedDatabase.ID(), &sdk.AlterDatabaseOptions{
			Set: &sdk.DatabaseSet{
				Comment: sdk.String("some comment"),
			},
		})
		require.NoError(t, err)

		database, err := client.Databases.ShowByID(ctx, sharedDatabase.ID())
		require.NoError(t, err)

		assert.Equal(t, sharedDatabase.ID().Name(), database.Name)
		assert.Equal(t, 1, database.RetentionTime)
		assert.Equal(t, comment, database.Comment)

		err = client.Databases.AlterReplication(ctx, sharedDatabase.ID(), &sdk.AlterDatabaseReplicationOptions{
			Refresh: sdk.Bool(true),
		})
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

		err := testClient(t).Databases.AlterReplication(ctx, database.ID(), &sdk.AlterDatabaseReplicationOptions{
			EnableReplication: &sdk.EnableReplication{
				ToAccounts: []sdk.AccountIdentifier{
					secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
				},
				IgnoreEditionCheck: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		err = testClient(t).Databases.AlterFailover(ctx, database.ID(), &sdk.AlterDatabaseFailoverOptions{
			EnableFailover: &sdk.EnableFailover{
				ToAccounts: []sdk.AccountIdentifier{
					secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
				},
			},
		})
		require.NoError(t, err)

		err = testClient(t).Databases.AlterFailover(ctx, database.ID(), &sdk.AlterDatabaseFailoverOptions{
			DisableFailover: &sdk.DisableFailover{
				ToAccounts: []sdk.AccountIdentifier{
					secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
				},
			},
		})
		require.NoError(t, err)
	})

	t.Run("promote to primary", func(t *testing.T) {
		t.Skipf("Can be unskipped after [SNOW-1002023]. CI Snowflake Edition doesn't support this feature")

		ctx := testContext(t)

		database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseCleanup)

		err := testClient(t).Databases.AlterReplication(ctx, database.ID(), &sdk.AlterDatabaseReplicationOptions{
			EnableReplication: &sdk.EnableReplication{
				ToAccounts: []sdk.AccountIdentifier{
					secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
				},
				IgnoreEditionCheck: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		err = testClient(t).Databases.AlterFailover(ctx, database.ID(), &sdk.AlterDatabaseFailoverOptions{
			EnableFailover: &sdk.EnableFailover{
				ToAccounts: []sdk.AccountIdentifier{
					secondaryTestClientHelper().Ids.AccountIdentifierWithLocator(),
				},
			},
		})
		require.NoError(t, err)

		err = testClient(t).Databases.AlterFailover(ctx, database.ID(), &sdk.AlterDatabaseFailoverOptions{
			Primary: sdk.Bool(true),
		})
		require.NoError(t, err)
	})
}

func TestInt_DatabasesDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("drop with nil options", func(t *testing.T) {
		databaseTest, databaseTestCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseTestCleanup)

		err := client.Databases.Drop(ctx, databaseTest.ID(), nil)
		require.NoError(t, err)
	})

	t.Run("drop if exists", func(t *testing.T) {
		databaseTest, databaseTestCleanup := testClientHelper().Database.CreateDatabase(t)
		databaseTestCleanup()

		err := client.Databases.Drop(ctx, databaseTest.ID(), &sdk.DropDatabaseOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)
	})

	t.Run("drop with cascade", func(t *testing.T) {
		databaseTest, databaseTestCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseTestCleanup)

		err := client.Databases.Drop(ctx, databaseTest.ID(), &sdk.DropDatabaseOptions{
			IfExists: sdk.Bool(true),
			Cascade:  sdk.Bool(true),
		})
		require.NoError(t, err)
	})

	t.Run("drop with restrict", func(t *testing.T) {
		databaseTest, databaseTestCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseTestCleanup)

		err := client.Databases.Drop(ctx, databaseTest.ID(), &sdk.DropDatabaseOptions{
			IfExists: sdk.Bool(true),
			Restrict: sdk.Bool(true),
		})
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

	err = client.Databases.Undrop(ctx, databaseTest.ID())
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
		databases, err := client.Databases.Show(ctx, nil)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, len(databases), 2)
		databaseIDs := make([]sdk.AccountObjectIdentifier, len(databases))
		for i, database := range databases {
			databaseIDs[i] = database.ID()
		}
		assert.Contains(t, databaseIDs, databaseTest.ID())
		assert.Contains(t, databaseIDs, databaseTest2.ID())
		assert.Equal(t, "ROLE", databases[0].OwnerRoleType)
	})

	t.Run("with terse", func(t *testing.T) {
		databases, err := client.Databases.Show(ctx, &sdk.ShowDatabasesOptions{
			Terse: sdk.Bool(true),
			Like: &sdk.Like{
				Pattern: sdk.String(databaseTest.Name),
			},
		})
		require.NoError(t, err)

		database, err := collections.FindFirst(databases, func(database sdk.Database) bool { return database.Name == databaseTest.Name })
		require.NoError(t, err)

		assert.Equal(t, databaseTest.Name, database.Name)
		assert.NotEmpty(t, database.CreatedOn)
		assert.Empty(t, database.DroppedOn)
		assert.Empty(t, database.Owner)
	})

	t.Run("with history", func(t *testing.T) {
		databaseTest3, databaseCleanup3 := testClientHelper().Database.CreateDatabase(t)
		databaseCleanup3()

		databases, err := client.Databases.Show(ctx, &sdk.ShowDatabasesOptions{
			History: sdk.Bool(true),
			Like: &sdk.Like{
				Pattern: sdk.String(databaseTest3.Name),
			},
		})
		require.NoError(t, err)

		droppedDatabase, err := collections.FindFirst(databases, func(database sdk.Database) bool { return database.Name == databaseTest3.Name })
		require.NoError(t, err)

		assert.Equal(t, databaseTest3.Name, droppedDatabase.Name)
		assert.NotEmpty(t, droppedDatabase.DroppedOn)
	})

	t.Run("with like starts with", func(t *testing.T) {
		databases, err := client.Databases.Show(ctx, &sdk.ShowDatabasesOptions{
			StartsWith: sdk.String(databaseTest.Name),
			LimitFrom: &sdk.LimitFrom{
				Rows: sdk.Int(1),
			},
		})
		require.NoError(t, err)

		database, err := collections.FindFirst(databases, func(database sdk.Database) bool { return database.Name == databaseTest.Name })
		require.NoError(t, err)

		assert.Equal(t, databaseTest.Name, database.Name)
	})

	t.Run("when searching a non-existent database", func(t *testing.T) {
		databases, err := client.Databases.Show(ctx, &sdk.ShowDatabasesOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("non-existent"),
			},
		})
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

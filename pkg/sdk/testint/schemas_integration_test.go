//go:build !account_level_tests

package testint

import (
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Schemas(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	schema1, cleanupSchema1 := testClientHelper().Schema.CreateSchema(t)
	t.Cleanup(cleanupSchema1)
	schema2, cleanupSchema2 := testClientHelper().Schema.CreateSchema(t)
	t.Cleanup(cleanupSchema2)

	assertParameterEquals := func(t *testing.T, params []*sdk.Parameter, parameterName sdk.AccountParameter, expected string) {
		t.Helper()
		assert.Equal(t, expected, helpers.FindParameter(t, params, parameterName).Value)
	}

	assertParameterEqualsToDefaultValue := func(t *testing.T, params []*sdk.Parameter, parameterName sdk.ObjectParameter) {
		t.Helper()
		param, err := collections.FindFirst(params, func(param *sdk.Parameter) bool { return param.Key == string(parameterName) })
		assert.NoError(t, err)
		assert.NotNil(t, param)
		assert.Equal(t, (*param).Default, (*param).Value)
	}

	t.Run("create: minimal", func(t *testing.T) {
		schemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
		err := client.Schemas.Create(ctx, schemaId, &sdk.CreateSchemaOptions{
			OrReplace: sdk.Bool(true),
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Schema.DropSchemaFunc(t, schemaId))

		database, err := client.Schemas.ShowByID(ctx, schemaId)
		require.NoError(t, err)
		assert.Equal(t, schemaId.Name(), database.Name)

		params, err := client.Schemas.ShowParameters(ctx, schemaId)
		require.NoError(t, err)
		assertParameterEquals(t, params, sdk.AccountParameterDataRetentionTimeInDays, strconv.Itoa(testClientHelper().Database.TestDatabaseDataRetentionTimeInDays()))
		assertParameterEquals(t, params, sdk.AccountParameterMaxDataExtensionTimeInDays, strconv.Itoa(testClientHelper().Database.TestDatabaseMaxDataExtensionTimeInDays()))
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterExternalVolume)
		assertParameterEquals(t, params, sdk.AccountParameterCatalog, testClientHelper().Database.TestDatabaseCatalog().Name())
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterReplaceInvalidCharacters)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterDefaultDDLCollation)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterStorageSerializationPolicy)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterLogLevel)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterTraceLevel)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterSuspendTaskAfterNumFailures)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterTaskAutoRetryAttempts)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterUserTaskManagedInitialWarehouseSize)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterUserTaskTimeoutMs)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterQuotedIdentifiersIgnoreCase)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterEnableConsoleOutput)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterPipeExecutionPaused)
	})

	t.Run("create: or replace", func(t *testing.T) {
		schema, cleanupSchema := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(cleanupSchema)
		comment := "replaced"
		err := client.Schemas.Create(ctx, schema.ID(), &sdk.CreateSchemaOptions{
			OrReplace:                  sdk.Bool(true),
			DataRetentionTimeInDays:    sdk.Int(10),
			MaxDataExtensionTimeInDays: sdk.Int(10),
			DefaultDDLCollation:        sdk.Pointer(sdk.StringAllowEmpty{Value: "en_US-trim"}),
			WithManagedAccess:          sdk.Bool(true),
			Comment:                    sdk.String(comment),
		})
		require.NoError(t, err)
		s, err := client.Schemas.ShowByID(ctx, schema.ID())
		require.NoError(t, err)
		assert.Equal(t, schema.Name, s.Name)
		assert.Equal(t, "MANAGED ACCESS", *s.Options)
		assert.Equal(t, comment, s.Comment)
	})

	t.Run("create: if not exists", func(t *testing.T) {
		schema, cleanupSchema := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(cleanupSchema)
		comment := "some_comment"
		err := client.Schemas.Create(ctx, schema.ID(), &sdk.CreateSchemaOptions{
			IfNotExists: sdk.Bool(true),
			Comment:     sdk.String(comment),
		})
		require.NoError(t, err)
		s, err := client.Schemas.ShowByID(ctx, schema.ID())
		require.NoError(t, err)
		assert.NotEqual(t, comment, s.Comment)
	})

	t.Run("create: clone", func(t *testing.T) {
		comment := "some_comment"
		schemaID := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
		err := client.Schemas.Create(ctx, schemaID, &sdk.CreateSchemaOptions{
			Comment: sdk.String(comment),
		})
		require.NoError(t, err)

		clonedSchemaID := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
		err = client.Schemas.Create(ctx, clonedSchemaID, &sdk.CreateSchemaOptions{
			Comment: sdk.String(comment),
			Clone: &sdk.Clone{
				SourceObject: schemaID,
			},
		})
		require.NoError(t, err)

		s, err := client.Schemas.ShowByID(ctx, schemaID)
		require.NoError(t, err)

		cs, err := client.Schemas.ShowByID(ctx, clonedSchemaID)
		require.NoError(t, err)
		assert.Equal(t, s.Comment, cs.Comment)

		t.Cleanup(func() {
			err = client.Schemas.Drop(ctx, schemaID, nil)
			require.NoError(t, err)
			err = client.Schemas.Drop(ctx, clonedSchemaID, nil)
			require.NoError(t, err)
		})
	})

	t.Run("create: with tags", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		schemaID := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
		tagValue := random.String()
		err := client.Schemas.Create(ctx, schemaID, &sdk.CreateSchemaOptions{
			Tag: []sdk.TagAssociation{
				{
					Name:  tag.ID(),
					Value: tagValue,
				},
			},
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Schemas.Drop(ctx, schemaID, nil)
			require.NoError(t, err)
		})

		tv, err := client.SystemFunctions.GetTag(ctx, tag.ID(), schemaID, sdk.ObjectTypeSchema)
		require.NoError(t, err)
		assert.Equal(t, &tagValue, tv)
	})

	t.Run("create: complete", func(t *testing.T) {
		schemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifier()

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
		err := client.Schemas.Create(ctx, schemaId, &sdk.CreateSchemaOptions{
			Transient:                               sdk.Bool(true),
			IfNotExists:                             sdk.Bool(true),
			DataRetentionTimeInDays:                 sdk.Int(0),
			MaxDataExtensionTimeInDays:              sdk.Int(10),
			ExternalVolume:                          &externalVolume,
			Catalog:                                 &catalog,
			ReplaceInvalidCharacters:                sdk.Bool(true),
			DefaultDDLCollation:                     sdk.Pointer(sdk.StringAllowEmpty{Value: "en_US"}),
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
			PipeExecutionPaused:                     sdk.Bool(true),
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
		t.Cleanup(testClientHelper().Schema.DropSchemaFunc(t, schemaId))

		schema, err := client.Schemas.ShowByID(ctx, schemaId)
		require.NoError(t, err)
		assert.Equal(t, schemaId.Name(), schema.Name)
		assert.Equal(t, comment, schema.Comment)

		params, err := client.Schemas.ShowParameters(ctx, schemaId)
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
		assertParameterEquals(t, sdk.AccountParameterPipeExecutionPaused, "true")

		tag1Value, err := client.SystemFunctions.GetTag(ctx, tagTest.ID(), schema.ID(), sdk.ObjectTypeSchema)
		require.NoError(t, err)
		assert.Equal(t, sdk.Pointer("v1"), tag1Value)

		tag2Value, err := client.SystemFunctions.GetTag(ctx, tag2Test.ID(), schema.ID(), sdk.ObjectTypeSchema)
		require.NoError(t, err)
		assert.Equal(t, sdk.Pointer("v2"), tag2Value)
	})

	t.Run("alter: rename to", func(t *testing.T) {
		// new schema created on purpose
		schema, _ := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(func() {
			err := client.Sessions.UseSchema(ctx, testClientHelper().Ids.SchemaId())
			require.NoError(t, err)
		})
		newID := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
		err := client.Schemas.Alter(ctx, schema.ID(), &sdk.AlterSchemaOptions{
			NewName: sdk.Pointer(newID),
		})
		require.NoError(t, err)
		s, err := client.Schemas.ShowByID(ctx, newID)
		require.NoError(t, err)
		assert.Equal(t, newID, s.ID())
		t.Cleanup(func() {
			err = client.Schemas.Drop(ctx, newID, nil)
			require.NoError(t, err)
		})
	})

	t.Run("alter: swap with", func(t *testing.T) {
		// new schemas created on purpose
		schema, cleanupSchema := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(cleanupSchema)
		swapSchema, cleanupSwapSchema := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(cleanupSwapSchema)

		table, _ := testClientHelper().Table.CreateInSchema(t, schema.ID())
		t.Cleanup(func() {
			newId := sdk.NewSchemaObjectIdentifierInSchema(swapSchema.ID(), table.Name)
			err := client.Tables.Drop(ctx, sdk.NewDropTableRequest(newId))
			require.NoError(t, err)
		})

		err := client.Schemas.Alter(ctx, schema.ID(), &sdk.AlterSchemaOptions{
			SwapWith: sdk.Pointer(swapSchema.ID()),
		})
		require.NoError(t, err)

		schemaDetails, err := client.Schemas.Describe(ctx, swapSchema.ID())
		require.NoError(t, err)
		assert.Len(t, schemaDetails, 1)
		assert.Equal(t, "TABLE", schemaDetails[0].Kind)
		assert.Equal(t, table.Name, schemaDetails[0].Name)
	})

	t.Run("alter: set and unset parameters", func(t *testing.T) {
		// new schema created on purpose
		schemaTest, cleanupSchemaTest := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(cleanupSchemaTest)

		externalVolumeTest, externalVolumeTestCleanup := testClientHelper().ExternalVolume.Create(t)
		t.Cleanup(externalVolumeTestCleanup)

		catalogIntegrationTest, catalogIntegrationTestCleanup := testClientHelper().CatalogIntegration.Create(t)
		t.Cleanup(catalogIntegrationTestCleanup)

		err := client.Schemas.Alter(ctx, schemaTest.ID(), &sdk.AlterSchemaOptions{
			Set: &sdk.SchemaSet{
				DataRetentionTimeInDays:                 sdk.Int(42),
				MaxDataExtensionTimeInDays:              sdk.Int(42),
				ExternalVolume:                          &externalVolumeTest,
				Catalog:                                 &catalogIntegrationTest,
				ReplaceInvalidCharacters:                sdk.Bool(true),
				DefaultDDLCollation:                     sdk.Pointer(sdk.StringAllowEmpty{Value: "en_US"}),
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
				PipeExecutionPaused:                     sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		params, err := client.Schemas.ShowParameters(ctx, schemaTest.ID())
		require.NoError(t, err)
		assertParameterEquals(t, params, sdk.AccountParameterDataRetentionTimeInDays, "42")
		assertParameterEquals(t, params, sdk.AccountParameterMaxDataExtensionTimeInDays, "42")
		assertParameterEquals(t, params, sdk.AccountParameterExternalVolume, externalVolumeTest.Name())
		assertParameterEquals(t, params, sdk.AccountParameterCatalog, catalogIntegrationTest.Name())
		assertParameterEquals(t, params, sdk.AccountParameterReplaceInvalidCharacters, "true")
		assertParameterEquals(t, params, sdk.AccountParameterDefaultDDLCollation, "en_US")
		assertParameterEquals(t, params, sdk.AccountParameterStorageSerializationPolicy, string(sdk.StorageSerializationPolicyCompatible))
		assertParameterEquals(t, params, sdk.AccountParameterLogLevel, string(sdk.LogLevelInfo))
		assertParameterEquals(t, params, sdk.AccountParameterTraceLevel, string(sdk.TraceLevelPropagate))
		assertParameterEquals(t, params, sdk.AccountParameterSuspendTaskAfterNumFailures, "10")
		assertParameterEquals(t, params, sdk.AccountParameterTaskAutoRetryAttempts, "10")
		assertParameterEquals(t, params, sdk.AccountParameterUserTaskManagedInitialWarehouseSize, string(sdk.WarehouseSizeMedium))
		assertParameterEquals(t, params, sdk.AccountParameterUserTaskTimeoutMs, "12000")
		assertParameterEquals(t, params, sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds, "30")
		assertParameterEquals(t, params, sdk.AccountParameterQuotedIdentifiersIgnoreCase, "true")
		assertParameterEquals(t, params, sdk.AccountParameterEnableConsoleOutput, "true")
		assertParameterEquals(t, params, sdk.AccountParameterPipeExecutionPaused, "true")

		err = client.Schemas.Alter(ctx, schemaTest.ID(), &sdk.AlterSchemaOptions{
			Unset: &sdk.SchemaUnset{
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
				PipeExecutionPaused:                     sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		params, err = client.Schemas.ShowParameters(ctx, schemaTest.ID())
		require.NoError(t, err)
		assertParameterEquals(t, params, sdk.AccountParameterDataRetentionTimeInDays, strconv.Itoa(testClientHelper().Database.TestDatabaseDataRetentionTimeInDays()))
		assertParameterEquals(t, params, sdk.AccountParameterMaxDataExtensionTimeInDays, strconv.Itoa(testClientHelper().Database.TestDatabaseMaxDataExtensionTimeInDays()))
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterExternalVolume)
		assertParameterEquals(t, params, sdk.AccountParameterCatalog, testClientHelper().Database.TestDatabaseCatalog().Name())
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterReplaceInvalidCharacters)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterDefaultDDLCollation)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterStorageSerializationPolicy)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterLogLevel)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterTraceLevel)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterSuspendTaskAfterNumFailures)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterTaskAutoRetryAttempts)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterUserTaskManagedInitialWarehouseSize)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterUserTaskTimeoutMs)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterQuotedIdentifiersIgnoreCase)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterEnableConsoleOutput)
		assertParameterEqualsToDefaultValue(t, params, sdk.ObjectParameterPipeExecutionPaused)
	})

	t.Run("alter: set non-parameters", func(t *testing.T) {
		schemaID := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
		comment := random.Comment()
		err := client.Schemas.Create(ctx, schemaID, nil)
		require.NoError(t, err)

		err = client.Schemas.Alter(ctx, schemaID, &sdk.AlterSchemaOptions{
			Set: &sdk.SchemaSet{
				Comment: sdk.Pointer(comment),
			},
		})
		require.NoError(t, err)

		s, err := client.Schemas.ShowByID(ctx, schemaID)
		require.NoError(t, err)
		assert.Equal(t, comment, s.Comment)

		t.Cleanup(func() {
			err := client.Schemas.Drop(ctx, schemaID, nil)
			require.NoError(t, err)
		})
	})

	t.Run("alter: unset non-parameters", func(t *testing.T) {
		schemaID := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
		comment := random.Comment()
		err := client.Schemas.Create(ctx, schemaID, &sdk.CreateSchemaOptions{
			Comment: sdk.String(comment),
		})
		require.NoError(t, err)

		err = client.Schemas.Alter(ctx, schemaID, &sdk.AlterSchemaOptions{
			Unset: &sdk.SchemaUnset{
				Comment: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		s, err := client.Schemas.ShowByID(ctx, schemaID)
		require.NoError(t, err)
		assert.Empty(t, s.Comment)

		t.Cleanup(func() {
			err := client.Schemas.Drop(ctx, schemaID, nil)
			require.NoError(t, err)
		})
	})

	t.Run("alter: enable managed access", func(t *testing.T) {
		schema, cleanupSchema := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(cleanupSchema)

		err := client.Schemas.Alter(ctx, schema.ID(), &sdk.AlterSchemaOptions{
			EnableManagedAccess: sdk.Bool(true),
		})
		require.NoError(t, err)

		s, err := client.Schemas.ShowByID(ctx, schema.ID())
		require.NoError(t, err)
		assert.Equal(t, schema.Name, s.Name)
		assert.True(t, s.IsManagedAccess())
	})

	t.Run("alter: disable managed access", func(t *testing.T) {
		id := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
		schema, cleanupSchema := testClientHelper().Schema.CreateSchemaWithOpts(t, id, &sdk.CreateSchemaOptions{
			WithManagedAccess: sdk.Pointer(true),
		})
		t.Cleanup(cleanupSchema)

		err := client.Schemas.Alter(ctx, schema.ID(), &sdk.AlterSchemaOptions{
			DisableManagedAccess: sdk.Bool(true),
		})
		require.NoError(t, err)

		s, err := client.Schemas.ShowByID(ctx, schema.ID())
		require.NoError(t, err)
		assert.Equal(t, schema.Name, s.Name)
		assert.False(t, s.IsManagedAccess())
	})

	t.Run("show: no options", func(t *testing.T) {
		schemas, err := client.Schemas.Show(ctx, nil)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(schemas), 2)
		schemaIds := make([]sdk.DatabaseObjectIdentifier, len(schemas))
		for i, schema := range schemas {
			schemaIds[i] = schema.ID()
		}
		assert.Contains(t, schemaIds, schema1.ID())
		assert.Contains(t, schemaIds, schema2.ID())
	})

	t.Run("show: with terse", func(t *testing.T) {
		schemas, err := client.Schemas.Show(ctx, &sdk.ShowSchemaOptions{
			Terse: sdk.Bool(true),
			Like: &sdk.Like{
				Pattern: sdk.String(schema1.Name),
			},
		})
		require.NoError(t, err)

		schema, err := collections.FindFirst(schemas, func(schema sdk.Schema) bool { return schema.Name == schema1.Name })
		require.NoError(t, err)

		assert.Equal(t, schema1.Name, schema.Name)
		assert.NotEmpty(t, schema.CreatedOn)
		assert.Empty(t, schema.Owner)
	})

	t.Run("show: with history", func(t *testing.T) {
		schema3, cleanupSchema3 := testClientHelper().Schema.CreateSchema(t)
		cleanupSchema3()
		schemas, err := client.Schemas.Show(ctx, &sdk.ShowSchemaOptions{
			History: sdk.Bool(true),
			Like: &sdk.Like{
				Pattern: sdk.String(schema3.Name),
			},
		})
		require.NoError(t, err)

		droppedSchema, err := collections.FindFirst(schemas, func(schema sdk.Schema) bool { return schema.Name == schema3.Name })
		require.NoError(t, err)

		assert.Equal(t, schema3.Name, droppedSchema.Name)
		assert.NotEmpty(t, droppedSchema.DroppedOn)
	})

	t.Run("show: with options", func(t *testing.T) {
		schemas, err := client.Schemas.Show(ctx, &sdk.ShowSchemaOptions{
			Terse:   sdk.Bool(true),
			History: sdk.Bool(true),
			Like: &sdk.Like{
				Pattern: sdk.String(schema1.Name),
			},
			In: &sdk.SchemaIn{
				Account: sdk.Bool(true),
			},
			StartsWith: sdk.String(schema1.Name),
			LimitFrom: &sdk.LimitFrom{
				Rows: sdk.Int(1),
			},
		})
		require.NoError(t, err)
		schemaNames := make([]string, len(schemas))
		for i, s := range schemas {
			schemaNames[i] = s.Name
		}
		assert.Contains(t, schemaNames, schema1.Name)
		assert.Equal(t, "ROLE", schema1.OwnerRoleType)
	})

	t.Run("show by id: missing schema", func(t *testing.T) {
		schemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
		_, err := client.Schemas.ShowByID(ctx, schemaId)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("show by id: missing database", func(t *testing.T) {
		databaseId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		schemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
		_, err := client.Schemas.ShowByID(ctx, schemaId)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sdk.ErrDoesNotExistOrOperationCannotBePerformed)
	})

	t.Run("show by id safely", func(t *testing.T) {
		schema, err := client.Schemas.ShowByIDSafely(ctx, testClientHelper().Ids.SchemaId())
		assert.NotNil(t, schema)
		assert.NoError(t, err)
	})

	t.Run("show by id safely: missing schema", func(t *testing.T) {
		schemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
		_, err := client.Schemas.ShowByIDSafely(ctx, schemaId)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("show by id safely: missing database", func(t *testing.T) {
		databaseId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		schemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
		_, err := client.Schemas.ShowByIDSafely(ctx, schemaId)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
		assert.ErrorIs(t, err, sdk.ErrDoesNotExistOrOperationCannotBePerformed)
	})

	t.Run("drop", func(t *testing.T) {
		schema, cleanupSchema := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(cleanupSchema)

		s, err := client.Schemas.ShowByID(ctx, schema.ID())
		require.NoError(t, err)
		assert.Equal(t, schema.Name, s.Name)

		err = client.Schemas.Drop(ctx, schema.ID(), nil)
		require.NoError(t, err)

		_, err = client.Schemas.ShowByID(ctx, schema.ID())
		assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("undrop", func(t *testing.T) {
		schema, cleanupSchema := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(cleanupSchema)

		before, err := client.Schemas.ShowByID(ctx, schema.ID())
		require.NoError(t, err)
		assert.Equal(t, schema.Name, before.Name)

		err = client.Schemas.Drop(ctx, schema.ID(), nil)
		require.NoError(t, err)

		_, err = client.Schemas.ShowByID(ctx, schema.ID())
		assert.ErrorIs(t, err, sdk.ErrObjectNotFound)

		err = client.Schemas.Undrop(ctx, schema.ID())
		require.NoError(t, err)

		after, err := client.Schemas.ShowByID(ctx, schema.ID())
		require.NoError(t, err)
		assert.Equal(t, schema.Name, after.Name)
	})
}

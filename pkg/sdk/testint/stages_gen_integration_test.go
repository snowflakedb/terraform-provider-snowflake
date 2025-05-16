//go:build !account_level_tests

package testint

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/ids"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Stages(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)
	gcsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.GcsExternalBucketUrl)
	azureBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)
	azureSasToken := testenvs.GetOrSkipTest(t, testenvs.AzureExternalSasToken)

	s3StorageIntegration, err := client.StorageIntegrations.ShowByID(ctx, ids.PrecreatedS3StorageIntegration)
	require.NoError(t, err)
	gcpStorageIntegration, err := client.StorageIntegrations.ShowByID(ctx, ids.PrecreatedGcpStorageIntegration)
	require.NoError(t, err)
	azureStorageIntegration, err := client.StorageIntegrations.ShowByID(ctx, ids.PrecreatedAzureStorageIntegration)
	require.NoError(t, err)

	cleanupStage := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		t.Cleanup(func() {
			err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(id))
			require.NoError(t, err)
		})
	}

	createBasicS3Stage := func(t *testing.T, stageId sdk.SchemaObjectIdentifier) {
		t.Helper()
		s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
			WithCredentials(sdk.NewExternalStageS3CredentialsRequest().
				WithAwsKeyId(&awsKeyId).
				WithAwsSecretKey(&awsSecretKey))
		err := client.Stages.CreateOnS3(ctx, sdk.NewCreateOnS3StageRequest(stageId).
			WithFileFormat(sdk.NewStageFileFormatRequest().WithType(&sdk.FileFormatTypeJSON)).
			WithExternalStageParams(s3Req))
		require.NoError(t, err)
		cleanupStage(t, stageId)
	}

	createBasicGcsStage := func(t *testing.T, stageId sdk.SchemaObjectIdentifier) {
		t.Helper()
		err := client.Stages.CreateOnGCS(ctx, sdk.NewCreateOnGCSStageRequest(stageId).
			WithFileFormat(sdk.NewStageFileFormatRequest().WithType(&sdk.FileFormatTypeJSON)).
			WithExternalStageParams(sdk.NewExternalGCSStageParamsRequest(gcsBucketUrl).
				WithStorageIntegration(sdk.Pointer(ids.PrecreatedGcpStorageIntegration))))
		require.NoError(t, err)
		cleanupStage(t, stageId)
	}

	createBasicAzureStage := func(t *testing.T, stageId sdk.SchemaObjectIdentifier) {
		t.Helper()
		err := client.Stages.CreateOnAzure(ctx, sdk.NewCreateOnAzureStageRequest(stageId).
			WithFileFormat(sdk.NewStageFileFormatRequest().WithType(&sdk.FileFormatTypeJSON)).
			WithExternalStageParams(sdk.NewExternalAzureStageParamsRequest(azureBucketUrl).
				WithCredentials(sdk.NewExternalStageAzureCredentialsRequest(azureSasToken))))
		require.NoError(t, err)
		cleanupStage(t, stageId)
	}

	assertStage := func(t *testing.T, stage *sdk.Stage, id sdk.SchemaObjectIdentifier, stageType string, comment string, cloud string, url string, storageIntegration string) {
		t.Helper()
		assert.Equal(t, id.DatabaseName(), stage.DatabaseName)
		assert.Equal(t, id.SchemaName(), stage.SchemaName)
		assert.Equal(t, id.Name(), stage.Name)
		assert.Equal(t, comment, stage.Comment)
		if len(url) > 0 {
			assert.Equal(t, url, stage.Url)
		}
		assert.Equal(t, stageType, stage.Type)
		if len(cloud) > 0 {
			assert.Equal(t, cloud, *stage.Cloud)
		}
		if len(storageIntegration) > 0 {
			assert.Equal(t, storageIntegration, *stage.StorageIntegration)
		}
	}

	t.Run("CreateInternal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Stages.CreateInternal(ctx, sdk.NewCreateInternalStageRequest(id).
			WithFileFormat(sdk.NewStageFileFormatRequest().WithType(&sdk.FileFormatTypeJSON)).
			WithComment(sdk.String("some comment")))
		require.NoError(t, err)
		cleanupStage(t, id)

		stage, err := client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)
		assertStage(t, stage, id, "INTERNAL", "some comment", "", "", "")
	})

	t.Run("CreateInternal - temporary", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Stages.CreateInternal(ctx, sdk.NewCreateInternalStageRequest(id).
			WithTemporary(sdk.Bool(true)).
			WithFileFormat(sdk.NewStageFileFormatRequest().WithType(&sdk.FileFormatTypeJSON)).
			WithComment(sdk.String("some comment")))
		require.NoError(t, err)
		cleanupStage(t, id)

		stage, err := client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)
		assertStage(t, stage, id, "INTERNAL TEMPORARY", "some comment", "", "", "")
	})

	t.Run("CreateOnS3 - IAM User", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
			WithCredentials(sdk.NewExternalStageS3CredentialsRequest().
				WithAwsKeyId(&awsKeyId).
				WithAwsSecretKey(&awsSecretKey))
		err := client.Stages.CreateOnS3(ctx, sdk.NewCreateOnS3StageRequest(id).
			WithFileFormat(sdk.NewStageFileFormatRequest().WithType(&sdk.FileFormatTypeJSON)).
			WithExternalStageParams(s3Req).
			WithComment(sdk.String("some comment")))
		require.NoError(t, err)
		cleanupStage(t, id)

		stage, err := client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)
		assertStage(t, stage, id, "EXTERNAL", "some comment", "AWS", awsBucketUrl, "")
	})

	t.Run("CreateOnS3 - temporary - Storage Integration", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
			WithStorageIntegration(sdk.Pointer(ids.PrecreatedS3StorageIntegration))
		err := client.Stages.CreateOnS3(ctx, sdk.NewCreateOnS3StageRequest(id).
			WithTemporary(sdk.Bool(true)).
			WithFileFormat(sdk.NewStageFileFormatRequest().WithType(&sdk.FileFormatTypeJSON)).
			WithExternalStageParams(s3Req).
			WithComment(sdk.String("some comment")))
		require.NoError(t, err)
		cleanupStage(t, id)

		stage, err := client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)
		assertStage(t, stage, id, "EXTERNAL TEMPORARY", "some comment", "AWS", awsBucketUrl, s3StorageIntegration.Name)
	})

	t.Run("CreateOnGCS", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Stages.CreateOnGCS(ctx, sdk.NewCreateOnGCSStageRequest(id).
			WithFileFormat(sdk.NewStageFileFormatRequest().WithType(&sdk.FileFormatTypeJSON)).
			WithExternalStageParams(sdk.NewExternalGCSStageParamsRequest(gcsBucketUrl).
				WithStorageIntegration(sdk.Pointer(ids.PrecreatedGcpStorageIntegration))).
			WithComment(sdk.String("some comment")))
		require.NoError(t, err)
		cleanupStage(t, id)

		stage, err := client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)
		assertStage(t, stage, id, "EXTERNAL", "some comment", "GCP", gcsBucketUrl, gcpStorageIntegration.Name)
	})

	t.Run("CreateOnAzure - Storage Integration", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Stages.CreateOnAzure(ctx, sdk.NewCreateOnAzureStageRequest(id).
			WithFileFormat(sdk.NewStageFileFormatRequest().WithType(&sdk.FileFormatTypeJSON)).
			WithExternalStageParams(sdk.NewExternalAzureStageParamsRequest(azureBucketUrl).
				WithStorageIntegration(sdk.Pointer(ids.PrecreatedAzureStorageIntegration))).
			WithComment(sdk.String("some comment")))
		require.NoError(t, err)
		cleanupStage(t, id)

		stage, err := client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)
		assertStage(t, stage, id, "EXTERNAL", "some comment", "AZURE", azureBucketUrl, azureStorageIntegration.Name)
	})

	t.Run("CreateOnAzure - Shared Access Signature", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Stages.CreateOnAzure(ctx, sdk.NewCreateOnAzureStageRequest(id).
			WithFileFormat(sdk.NewStageFileFormatRequest().WithType(&sdk.FileFormatTypeJSON)).
			WithExternalStageParams(sdk.NewExternalAzureStageParamsRequest(azureBucketUrl).
				WithCredentials(sdk.NewExternalStageAzureCredentialsRequest(azureSasToken))).
			WithComment(sdk.String("some comment")))
		require.NoError(t, err)
		cleanupStage(t, id)

		stage, err := client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)
		assertStage(t, stage, id, "EXTERNAL", "some comment", "AZURE", azureBucketUrl, "")
	})

	t.Run("CreateOnS3Compatible", func(t *testing.T) {
		// TODO: (SNOW-1012064) create s3 compat service for tests
	})

	t.Run("Alter - rename", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		renamed := false

		err := client.Stages.CreateInternal(ctx, sdk.NewCreateInternalStageRequest(id))
		require.NoError(t, err)
		t.Cleanup(func() {
			if renamed {
				err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(newId))
				require.NoError(t, err)
			} else {
				err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(id))
				require.NoError(t, err)
			}
		})

		err = client.Stages.Alter(ctx, sdk.NewAlterStageRequest(id).
			WithIfExists(sdk.Bool(true)).
			WithRenameTo(&newId))
		require.NoError(t, err)
		renamed = true

		stage, err := client.Stages.ShowByID(ctx, newId)
		require.NotNil(t, stage)
		require.NoError(t, err)
	})

	t.Run("AlterInternalStage", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Stages.CreateInternal(ctx, sdk.NewCreateInternalStageRequest(id).
			WithCopyOptions(sdk.NewStageCopyOptionsRequest().WithSizeLimit(sdk.Int(100))).
			WithFileFormat(sdk.NewStageFileFormatRequest().WithType(&sdk.FileFormatTypeJSON)).
			WithComment(sdk.String("some comment")))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(id))
			require.NoError(t, err)
		})

		stage, err := client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, "some comment", stage.Comment)

		stageProperties, err := client.Stages.Describe(ctx, id)
		require.NoError(t, err)
		require.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "STAGE_COPY_OPTIONS",
			Name:    "SIZE_LIMIT",
			Type:    "Long",
			Value:   "100",
			Default: "",
		})
		require.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "STAGE_FILE_FORMAT",
			Name:    "TYPE",
			Type:    "String",
			Value:   "JSON",
			Default: "CSV",
		})

		err = client.Stages.AlterInternalStage(ctx, sdk.NewAlterInternalStageStageRequest(id).
			WithIfExists(sdk.Bool(true)).
			WithCopyOptions(sdk.NewStageCopyOptionsRequest().WithSizeLimit(sdk.Int(200))).
			WithFileFormat(sdk.NewStageFileFormatRequest().WithType(&sdk.FileFormatTypeCSV)).
			WithComment(sdk.String("altered comment")))
		require.NoError(t, err)

		stage, err = client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, "altered comment", stage.Comment)

		stageProperties, err = client.Stages.Describe(ctx, id)
		require.NoError(t, err)
		require.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "STAGE_COPY_OPTIONS",
			Name:    "SIZE_LIMIT",
			Type:    "Long",
			Value:   "200",
			Default: "",
		})
		require.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "STAGE_FILE_FORMAT",
			Name:    "TYPE",
			Type:    "String",
			Value:   "CSV",
			Default: "CSV",
		})
	})

	t.Run("AlterExternalS3Stage", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createBasicS3Stage(t, id)

		err := client.Stages.AlterExternalS3Stage(ctx, sdk.NewAlterExternalS3StageStageRequest(id).
			WithExternalStageParams(sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
				WithStorageIntegration(sdk.Pointer(ids.PrecreatedS3StorageIntegration))).
			WithComment(sdk.String("Updated comment")))
		require.NoError(t, err)

		stage, err := client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)
		assertStage(t, stage, id, "EXTERNAL", "Updated comment", "AWS", awsBucketUrl, s3StorageIntegration.Name)

		props, err := client.Stages.Describe(ctx, id)

		require.NoError(t, err)
		require.NotNil(t, props)
	})

	t.Run("AlterExternalGCSStage", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createBasicGcsStage(t, id)

		err := client.Stages.AlterExternalGCSStage(ctx, sdk.NewAlterExternalGCSStageStageRequest(id).
			WithExternalStageParams(sdk.NewExternalGCSStageParamsRequest(gcsBucketUrl).
				WithStorageIntegration(sdk.Pointer(ids.PrecreatedGcpStorageIntegration))).
			WithComment(sdk.String("Updated comment")))
		require.NoError(t, err)

		stage, err := client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)
		assertStage(t, stage, id, "EXTERNAL", "Updated comment", "GCP", gcsBucketUrl, gcpStorageIntegration.Name)
	})

	t.Run("AlterExternalAzureStage", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createBasicAzureStage(t, id)

		err := client.Stages.AlterExternalAzureStage(ctx, sdk.NewAlterExternalAzureStageStageRequest(id).
			WithExternalStageParams(sdk.NewExternalAzureStageParamsRequest(azureBucketUrl).
				WithStorageIntegration(sdk.Pointer(ids.PrecreatedAzureStorageIntegration))).
			WithComment(sdk.String("Updated comment")))
		require.NoError(t, err)

		stage, err := client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)
		assertStage(t, stage, id, "EXTERNAL", "Updated comment", "AZURE", azureBucketUrl, azureStorageIntegration.Name)
	})

	t.Run("AlterDirectoryTable", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createBasicS3Stage(t, id)

		stageProperties, err := client.Stages.Describe(ctx, id)
		require.NoError(t, err)
		assert.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "DIRECTORY",
			Name:    "ENABLE",
			Type:    "Boolean",
			Value:   "false",
			Default: "false",
		})

		err = client.Stages.AlterDirectoryTable(ctx, sdk.NewAlterDirectoryTableStageRequest(id).
			WithSetDirectory(sdk.NewDirectoryTableSetRequest(true)))
		require.NoError(t, err)

		err = client.Stages.AlterDirectoryTable(ctx, sdk.NewAlterDirectoryTableStageRequest(id).
			WithRefresh(sdk.NewDirectoryTableRefreshRequest().WithSubpath(sdk.String("/"))))
		require.NoError(t, err)

		stageProperties, err = client.Stages.Describe(ctx, id)
		require.NoError(t, err)
		assert.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "DIRECTORY",
			Name:    "ENABLE",
			Type:    "Boolean",
			Value:   "true",
			Default: "false",
		})
	})

	t.Run("Drop", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Stages.CreateInternal(ctx, sdk.NewCreateInternalStageRequest(id))
		require.NoError(t, err)

		stage, err := client.Stages.ShowByID(ctx, id)
		require.NotNil(t, stage)
		require.NoError(t, err)

		err = client.Stages.Drop(ctx, sdk.NewDropStageRequest(id))
		require.NoError(t, err)

		stage, err = client.Stages.ShowByID(ctx, id)
		require.Nil(t, stage)
		require.Error(t, err)
	})

	t.Run("Describe internal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Stages.CreateInternal(ctx, sdk.NewCreateInternalStageRequest(id))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(id))
			require.NoError(t, err)
		})

		stageProperties, err := client.Stages.Describe(ctx, id)
		require.NoError(t, err)
		require.NotEmpty(t, stageProperties)
		assert.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "DIRECTORY",
			Name:    "ENABLE",
			Type:    "Boolean",
			Value:   "false",
			Default: "false",
		})
		assert.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "STAGE_LOCATION",
			Name:    "URL",
			Type:    "String",
			Value:   "",
			Default: "",
		})
	})

	t.Run("Describe external s3", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createBasicS3Stage(t, id)

		stageProperties, err := client.Stages.Describe(ctx, id)
		require.NoError(t, err)
		require.NotEmpty(t, stageProperties)
		assert.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "STAGE_CREDENTIALS",
			Name:    "AWS_KEY_ID",
			Type:    "String",
			Value:   awsKeyId,
			Default: "",
		})
		assert.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "STAGE_LOCATION",
			Name:    "URL",
			Type:    "String",
			Value:   fmt.Sprintf("[\"%s\"]", awsBucketUrl),
			Default: "",
		})
	})

	t.Run("Describe external gcs", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createBasicGcsStage(t, id)

		stageProperties, err := client.Stages.Describe(ctx, id)
		require.NoError(t, err)
		require.NotEmpty(t, stageProperties)
		assert.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "STAGE_LOCATION",
			Name:    "URL",
			Type:    "String",
			Value:   fmt.Sprintf("[\"%s\"]", gcsBucketUrl),
			Default: "",
		})
	})

	t.Run("Describe external azure", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createBasicAzureStage(t, id)

		stageProperties, err := client.Stages.Describe(ctx, id)
		require.NoError(t, err)
		require.NotEmpty(t, stageProperties)
		assert.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "DIRECTORY",
			Name:    "ENABLE",
			Type:    "Boolean",
			Value:   "false",
			Default: "false",
		})
		assert.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "STAGE_LOCATION",
			Name:    "URL",
			Type:    "String",
			Value:   fmt.Sprintf("[\"%s\"]", azureBucketUrl),
			Default: "",
		})
	})

	t.Run("Show internal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Stages.CreateInternal(ctx, sdk.NewCreateInternalStageRequest(id).
			WithDirectoryTableOptions(sdk.NewInternalDirectoryTableOptionsRequest().WithEnable(sdk.Bool(true))).
			WithComment(sdk.String("some comment")))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(id))
			require.NoError(t, err)
		})

		stage, err := client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.DatabaseName(), stage.DatabaseName)
		assert.Equal(t, id.SchemaName(), stage.SchemaName)
		assert.Equal(t, id.Name(), stage.Name)
		assert.Empty(t, stage.Url)
		assert.False(t, stage.HasCredentials)
		assert.False(t, stage.HasEncryptionKey)
		assert.Equal(t, "some comment", stage.Comment)
		assert.Nil(t, stage.Region)
		assert.Equal(t, "INTERNAL", stage.Type)
		assert.Nil(t, stage.Cloud)
		assert.Nil(t, stage.StorageIntegration)
		assert.Nil(t, stage.Endpoint)
		assert.True(t, stage.DirectoryEnabled)
		assert.Equal(t, "ROLE", *stage.OwnerRoleType)
	})
}

func TestInt_StagesShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupStageHandle := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}
	createStageHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		err := client.Stages.CreateInternal(ctx, sdk.NewCreateInternalStageRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupStageHandle(id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		createStageHandle(t, id1)
		createStageHandle(t, id2)

		e1, err := client.Stages.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Stages.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}

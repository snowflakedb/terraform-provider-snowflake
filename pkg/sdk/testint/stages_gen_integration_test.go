//go:build non_account_level_tests

package testint

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/ids"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

	// ==================== INTERNAL STAGE TESTS ====================

	t.Run("CreateInternal - minimal", func(t *testing.T) {
		stage, cleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(stage.Name).
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasType("INTERNAL").
			HasComment("").
			HasUrl("").
			HasDirectoryEnabled(false).
			HasHasCredentials(false).
			HasHasEncryptionKey(false).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"))
	})

	t.Run("CreateInternal - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "test comment"

		request := sdk.NewCreateInternalStageRequest(id).
			WithEncryption(*sdk.NewInternalStageEncryptionRequest().
				WithSnowflakeFull(*sdk.NewInternalStageEncryptionSnowflakeFullRequest())).
			WithDirectoryTableOptions(*sdk.NewInternalDirectoryTableOptionsRequest().
				WithEnable(true).
				WithAutoRefresh(false)).
			WithComment(comment)

		stage, cleanup := testClientHelper().Stage.CreateStageWithRequest(t, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasType("INTERNAL").
			HasComment(comment).
			HasDirectoryEnabled(true).
			HasHasEncryptionKey(true).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"))
	})

	t.Run("CreateInternal - temporary", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "temporary stage"

		request := sdk.NewCreateInternalStageRequest(id).
			WithTemporary(true).
			WithComment(comment)

		stage, cleanup := testClientHelper().Stage.CreateStageWithRequest(t, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasType("INTERNAL TEMPORARY").
			HasComment(comment))
	})

	t.Run("AlterInternalStage - complete", func(t *testing.T) {
		stage, cleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanup)

		require.Equal(t, "", stage.Comment)

		err := client.Stages.AlterInternalStage(ctx, sdk.NewAlterInternalStageStageRequest(stage.ID()).
			WithIfExists(true).
			WithComment("altered comment"))
		require.NoError(t, err)

		stage, err = client.Stages.ShowByID(ctx, stage.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasComment("altered comment"))
	})

	// ==================== S3 STAGE TESTS ====================

	t.Run("CreateOnS3 - minimal with storage integration", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
			WithStorageIntegration(ids.PrecreatedS3StorageIntegration)

		request := sdk.NewCreateOnS3StageRequest(id, *s3Req)

		stage, cleanup := testClientHelper().Stage.CreateStageOnS3WithRequest(t, id, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasType("EXTERNAL").
			HasUrl(awsBucketUrl).
			HasCloud("AWS").
			HasStorageIntegration(s3StorageIntegration.Name).
			HasHasCredentials(false).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"))
	})

	t.Run("CreateOnS3 - minimal with credentials", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
			WithCredentials(*sdk.NewExternalStageS3CredentialsRequest().
				WithAwsKeyId(awsKeyId).
				WithAwsSecretKey(awsSecretKey))

		request := sdk.NewCreateOnS3StageRequest(id, *s3Req)

		stage, cleanup := testClientHelper().Stage.CreateStageOnS3WithRequest(t, id, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasType("EXTERNAL").
			HasUrl(awsBucketUrl).
			HasCloud("AWS").
			HasHasCredentials(true).
			HasOwner(snowflakeroles.Accountadmin.Name()))
	})

	t.Run("CreateOnS3 - complete with storage integration", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "complete s3 stage with storage integration"

		s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
			WithStorageIntegration(ids.PrecreatedS3StorageIntegration)

		request := sdk.NewCreateOnS3StageRequest(id, *s3Req).
			WithDirectoryTableOptions(*sdk.NewStageS3CommonDirectoryTableOptionsRequest().
				WithEnable(true).
				WithRefreshOnCreate(true).
				WithAutoRefresh(false)).
			WithComment(comment)

		stage, cleanup := testClientHelper().Stage.CreateStageOnS3WithRequest(t, id, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasType("EXTERNAL").
			HasUrl(awsBucketUrl).
			HasCloud("AWS").
			HasStorageIntegration(s3StorageIntegration.Name).
			HasDirectoryEnabled(true).
			HasComment(comment))
	})

	t.Run("CreateOnS3 - complete with credentials", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "complete s3 stage with credentials"

		s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
			WithAwsAccessPointArn("arn:aws:s3:us-west-2:123456789012:accesspoint/my-data-ap").
			WithUsePrivatelinkEndpoint(true).
			WithCredentials(*sdk.NewExternalStageS3CredentialsRequest().
				WithAwsKeyId(awsKeyId).
				WithAwsSecretKey(awsSecretKey)).
			WithEncryption(*sdk.NewExternalStageS3EncryptionRequest().
				WithAwsSseS3(*sdk.NewExternalStageS3EncryptionAwsSseS3Request()))

		request := sdk.NewCreateOnS3StageRequest(id, *s3Req).
			WithDirectoryTableOptions(*sdk.NewStageS3CommonDirectoryTableOptionsRequest().
				WithEnable(true).
				WithRefreshOnCreate(true).
				WithAutoRefresh(false)).
			WithComment(comment)

		stage, cleanup := testClientHelper().Stage.CreateStageOnS3WithRequest(t, id, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasType("EXTERNAL").
			HasUrl(awsBucketUrl).
			HasCloud("AWS").
			HasHasCredentials(true).
			HasHasEncryptionKey(false).
			HasDirectoryEnabled(true).
			HasComment(comment))

		stageProperties, err := client.Stages.Describe(ctx, id)
		require.NoError(t, err)
		require.NotEmpty(t, stageProperties)
		assert.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "STAGE_LOCATION",
			Name:    "AWS_ACCESS_POINT_ARN",
			Type:    "String",
			Value:   "arn:aws:s3:us-west-2:123456789012:accesspoint/my-data-ap",
			Default: "",
		})
		assert.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "PRIVATELINK",
			Name:    "USE_PRIVATELINK_ENDPOINT",
			Type:    "Boolean",
			Value:   "true",
			Default: "false",
		})
	})

	t.Run("CreateOnS3 - temporary", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "temporary s3 stage"

		s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
			WithStorageIntegration(ids.PrecreatedS3StorageIntegration)

		request := sdk.NewCreateOnS3StageRequest(id, *s3Req).
			WithTemporary(true).
			WithComment(comment)

		stage, cleanup := testClientHelper().Stage.CreateStageOnS3WithRequest(t, id, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasType("EXTERNAL TEMPORARY").
			HasUrl(awsBucketUrl).
			HasCloud("AWS").
			HasStorageIntegration(s3StorageIntegration.Name).
			HasComment(comment))
	})

	t.Run("AlterExternalS3Stage - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
			WithCredentials(*sdk.NewExternalStageS3CredentialsRequest().
				WithAwsKeyId(awsKeyId).
				WithAwsSecretKey(awsSecretKey))

		request := sdk.NewCreateOnS3StageRequest(id, *s3Req).
			WithComment("initial comment")

		stage, cleanup := testClientHelper().Stage.CreateStageOnS3WithRequest(t, id, request)
		t.Cleanup(cleanup)

		require.Equal(t, "initial comment", stage.Comment)

		err := client.Stages.AlterExternalS3Stage(ctx, sdk.NewAlterExternalS3StageStageRequest(id).
			WithExternalStageParams(*sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
				WithStorageIntegration(ids.PrecreatedS3StorageIntegration)).
			WithComment("Updated comment"))
		require.NoError(t, err)

		stage, err = client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasType("EXTERNAL").
			HasUrl(awsBucketUrl).
			HasCloud("AWS").
			HasStorageIntegration(s3StorageIntegration.Name).
			HasComment("Updated comment"))
	})

	// ==================== GCS STAGE TESTS ====================

	t.Run("CreateOnGCS - minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		gcsReq := sdk.NewExternalGCSStageParamsRequest(gcsBucketUrl).
			WithStorageIntegration(ids.PrecreatedGcpStorageIntegration)

		request := sdk.NewCreateOnGCSStageRequest(id, *gcsReq)

		stage, cleanup := testClientHelper().Stage.CreateStageOnGCSWithRequest(t, id, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasType("EXTERNAL").
			HasUrl(gcsBucketUrl).
			HasCloud("GCP").
			HasStorageIntegration(gcpStorageIntegration.Name).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"))
	})

	t.Run("CreateOnGCS - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "complete gcs stage"

		gcsReq := sdk.NewExternalGCSStageParamsRequest(gcsBucketUrl).
			WithStorageIntegration(ids.PrecreatedGcpStorageIntegration).
			WithEncryption(*sdk.NewExternalStageGCSEncryptionRequest().
				WithGcsSseKms(*sdk.NewExternalStageGCSEncryptionGcsSseKmsRequest()))

		request := sdk.NewCreateOnGCSStageRequest(id, *gcsReq).
			WithDirectoryTableOptions(*sdk.NewExternalGCSDirectoryTableOptionsRequest().
				WithEnable(true).
				WithRefreshOnCreate(true).
				WithAutoRefresh(false)).
			WithComment(comment)

		stage, cleanup := testClientHelper().Stage.CreateStageOnGCSWithRequest(t, id, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasType("EXTERNAL").
			HasUrl(gcsBucketUrl).
			HasCloud("GCP").
			HasStorageIntegration(gcpStorageIntegration.Name).
			HasDirectoryEnabled(true).
			HasComment(comment))
	})

	t.Run("CreateOnGCS - temporary", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "temporary gcs stage"

		gcsReq := sdk.NewExternalGCSStageParamsRequest(gcsBucketUrl).
			WithStorageIntegration(ids.PrecreatedGcpStorageIntegration)

		request := sdk.NewCreateOnGCSStageRequest(id, *gcsReq).
			WithTemporary(true).
			WithComment(comment)

		stage, cleanup := testClientHelper().Stage.CreateStageOnGCSWithRequest(t, id, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasType("EXTERNAL TEMPORARY").
			HasUrl(gcsBucketUrl).
			HasCloud("GCP").
			HasStorageIntegration(gcpStorageIntegration.Name).
			HasComment(comment))
	})

	t.Run("AlterExternalGCSStage - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		gcsReq := sdk.NewExternalGCSStageParamsRequest(gcsBucketUrl).
			WithStorageIntegration(ids.PrecreatedGcpStorageIntegration)

		request := sdk.NewCreateOnGCSStageRequest(id, *gcsReq).
			WithComment("initial comment")

		stage, cleanup := testClientHelper().Stage.CreateStageOnGCSWithRequest(t, id, request)
		t.Cleanup(cleanup)

		require.Equal(t, "initial comment", stage.Comment)

		err := client.Stages.AlterExternalGCSStage(ctx, sdk.NewAlterExternalGCSStageStageRequest(id).
			WithExternalStageParams(*sdk.NewExternalGCSStageParamsRequest(gcsBucketUrl).
				WithStorageIntegration(ids.PrecreatedGcpStorageIntegration)).
			WithComment("Updated comment"))
		require.NoError(t, err)

		stage, err = client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasType("EXTERNAL").
			HasUrl(gcsBucketUrl).
			HasCloud("GCP").
			HasStorageIntegration(gcpStorageIntegration.Name).
			HasComment("Updated comment"))
	})

	// ==================== AZURE STAGE TESTS ====================

	t.Run("CreateOnAzure - minimal with storage integration", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		azureReq := sdk.NewExternalAzureStageParamsRequest(azureBucketUrl).
			WithStorageIntegration(ids.PrecreatedAzureStorageIntegration)

		request := sdk.NewCreateOnAzureStageRequest(id, *azureReq)

		stage, cleanup := testClientHelper().Stage.CreateStageOnAzureWithRequest(t, id, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasType("EXTERNAL").
			HasUrl(azureBucketUrl).
			HasCloud("AZURE").
			HasStorageIntegration(azureStorageIntegration.Name).
			HasHasCredentials(false).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"))
	})

	t.Run("CreateOnAzure - minimal with credentials", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		azureReq := sdk.NewExternalAzureStageParamsRequest(azureBucketUrl).
			WithCredentials(*sdk.NewExternalStageAzureCredentialsRequest(azureSasToken))

		request := sdk.NewCreateOnAzureStageRequest(id, *azureReq)

		stage, cleanup := testClientHelper().Stage.CreateStageOnAzureWithRequest(t, id, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasType("EXTERNAL").
			HasUrl(azureBucketUrl).
			HasCloud("AZURE").
			HasHasCredentials(true).
			HasOwner(snowflakeroles.Accountadmin.Name()))
	})

	t.Run("CreateOnAzure - complete with storage integration", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "complete azure stage with storage integration"

		azureReq := sdk.NewExternalAzureStageParamsRequest(azureBucketUrl).
			WithStorageIntegration(ids.PrecreatedAzureStorageIntegration)

		request := sdk.NewCreateOnAzureStageRequest(id, *azureReq).
			WithDirectoryTableOptions(*sdk.NewExternalAzureDirectoryTableOptionsRequest().
				WithEnable(true).
				WithRefreshOnCreate(true).
				WithAutoRefresh(false)).
			WithComment(comment)

		stage, cleanup := testClientHelper().Stage.CreateStageOnAzureWithRequest(t, id, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasType("EXTERNAL").
			HasUrl(azureBucketUrl).
			HasCloud("AZURE").
			HasStorageIntegration(azureStorageIntegration.Name).
			HasDirectoryEnabled(true).
			HasComment(comment))
	})

	t.Run("CreateOnAzure - complete with credentials", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "complete azure stage with credentials"

		azureReq := sdk.NewExternalAzureStageParamsRequest(azureBucketUrl).
			WithCredentials(*sdk.NewExternalStageAzureCredentialsRequest(azureSasToken)).
			WithEncryption(*sdk.NewExternalStageAzureEncryptionRequest().
				WithAzureCse(*sdk.NewExternalStageAzureEncryptionAzureCseRequest().
					WithMasterKey("test-master-key")))

		request := sdk.NewCreateOnAzureStageRequest(id, *azureReq).
			WithDirectoryTableOptions(*sdk.NewExternalAzureDirectoryTableOptionsRequest().
				WithEnable(true).
				WithRefreshOnCreate(true).
				WithAutoRefresh(false)).
			WithComment(comment)

		stage, cleanup := testClientHelper().Stage.CreateStageOnAzureWithRequest(t, id, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasType("EXTERNAL").
			HasUrl(azureBucketUrl).
			HasCloud("AZURE").
			HasHasCredentials(true).
			HasHasEncryptionKey(true).
			HasDirectoryEnabled(true).
			HasComment(comment))
	})

	t.Run("CreateOnAzure - temporary", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "temporary azure stage"

		azureReq := sdk.NewExternalAzureStageParamsRequest(azureBucketUrl).
			WithStorageIntegration(ids.PrecreatedAzureStorageIntegration)

		request := sdk.NewCreateOnAzureStageRequest(id, *azureReq).
			WithTemporary(true).
			WithComment(comment)

		stage, cleanup := testClientHelper().Stage.CreateStageOnAzureWithRequest(t, id, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasType("EXTERNAL TEMPORARY").
			HasUrl(azureBucketUrl).
			HasCloud("AZURE").
			HasStorageIntegration(azureStorageIntegration.Name).
			HasComment(comment))
	})

	t.Run("AlterExternalAzureStage - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		azureReq := sdk.NewExternalAzureStageParamsRequest(azureBucketUrl).
			WithCredentials(*sdk.NewExternalStageAzureCredentialsRequest(azureSasToken))

		request := sdk.NewCreateOnAzureStageRequest(id, *azureReq).
			WithComment("initial comment")

		stage, cleanup := testClientHelper().Stage.CreateStageOnAzureWithRequest(t, id, request)
		t.Cleanup(cleanup)

		require.Equal(t, "initial comment", stage.Comment)

		err := client.Stages.AlterExternalAzureStage(ctx, sdk.NewAlterExternalAzureStageStageRequest(id).
			WithExternalStageParams(*sdk.NewExternalAzureStageParamsRequest(azureBucketUrl).
				WithStorageIntegration(ids.PrecreatedAzureStorageIntegration)).
			WithComment("Updated comment"))
		require.NoError(t, err)

		stage, err = client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasType("EXTERNAL").
			HasUrl(azureBucketUrl).
			HasCloud("AZURE").
			HasStorageIntegration(azureStorageIntegration.Name).
			HasComment("Updated comment"))
	})

	// ==================== S3-COMPATIBLE STAGE TESTS ====================

	t.Run("CreateOnS3Compatible - minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		compatibleBucketUrl := strings.Replace(awsBucketUrl, "s3://", "s3compat://", 1)
		endpoint := "s3.us-west-2.amazonaws.com"

		s3CompatReq := sdk.NewExternalS3CompatibleStageParamsRequest(compatibleBucketUrl, endpoint).
			WithCredentials(*sdk.NewExternalStageS3CompatibleCredentialsRequest(awsKeyId, awsSecretKey))

		request := sdk.NewCreateOnS3CompatibleStageRequest(id, *s3CompatReq)

		stage, cleanup := testClientHelper().Stage.CreateStageOnS3CompatibleWithRequest(t, id, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasType("EXTERNAL").
			HasUrl(compatibleBucketUrl).
			HasCloud("AWS").
			HasEndpoint(endpoint).
			HasHasCredentials(true).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE"))
	})

	t.Run("CreateOnS3Compatible - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		compatibleBucketUrl := strings.Replace(awsBucketUrl, "s3://", "s3compat://", 1)
		endpoint := "s3.us-west-2.amazonaws.com"
		comment := "complete s3 compatible stage"

		s3CompatReq := sdk.NewExternalS3CompatibleStageParamsRequest(compatibleBucketUrl, endpoint).
			WithCredentials(*sdk.NewExternalStageS3CompatibleCredentialsRequest(awsKeyId, awsSecretKey))

		request := sdk.NewCreateOnS3CompatibleStageRequest(id, *s3CompatReq).
			WithDirectoryTableOptions(*sdk.NewStageS3CommonDirectoryTableOptionsRequest().
				WithEnable(true).
				WithRefreshOnCreate(true).
				WithAutoRefresh(false)).
			WithComment(comment)

		stage, cleanup := testClientHelper().Stage.CreateStageOnS3CompatibleWithRequest(t, id, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasType("EXTERNAL").
			HasUrl(compatibleBucketUrl).
			HasCloud("AWS").
			HasEndpoint(endpoint).
			HasDirectoryEnabled(true).
			HasComment(comment))
	})

	t.Run("CreateOnS3Compatible - temporary", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		compatibleBucketUrl := strings.Replace(awsBucketUrl, "s3://", "s3compat://", 1)
		endpoint := "s3.us-west-2.amazonaws.com"
		comment := "temporary s3 compatible stage"

		s3CompatReq := sdk.NewExternalS3CompatibleStageParamsRequest(compatibleBucketUrl, endpoint).
			WithCredentials(*sdk.NewExternalStageS3CompatibleCredentialsRequest(awsKeyId, awsSecretKey))

		request := sdk.NewCreateOnS3CompatibleStageRequest(id, *s3CompatReq).
			WithTemporary(true).
			WithComment(comment)

		stage, cleanup := testClientHelper().Stage.CreateStageOnS3CompatibleWithRequest(t, id, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasType("EXTERNAL TEMPORARY").
			HasUrl(compatibleBucketUrl).
			HasCloud("AWS").
			HasEndpoint(endpoint).
			HasComment(comment))
	})

	// ==================== OTHER OPERATIONS ====================

	t.Run("Alter - rename", func(t *testing.T) {
		stage, cleanup := testClientHelper().Stage.CreateStage(t)
		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err = client.Stages.Alter(ctx, sdk.NewAlterStageRequest(stage.ID()).
			WithIfExists(true).
			WithRenameTo(newId))
		require.NoError(t, err)

		// Update cleanup to use new id
		t.Cleanup(func() {
			cleanup() // This will fail but we need to clean up with the new id
		})
		t.Cleanup(testClientHelper().Stage.DropStageFunc(t, newId))

		renamedStage, err := client.Stages.ShowByID(ctx, newId)
		require.NoError(t, err)
		require.NotNil(t, renamedStage)

		assertThatObject(t, objectassert.StageFromObject(t, renamedStage).
			HasName(newId.Name()))
	})

	t.Run("AlterDirectoryTable", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
			WithCredentials(*sdk.NewExternalStageS3CredentialsRequest().
				WithAwsKeyId(awsKeyId).
				WithAwsSecretKey(awsSecretKey))

		request := sdk.NewCreateOnS3StageRequest(id, *s3Req)

		_, cleanup := testClientHelper().Stage.CreateStageOnS3WithRequest(t, id, request)
		t.Cleanup(cleanup)

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
			WithSetDirectory(*sdk.NewDirectoryTableSetRequest(true)))
		require.NoError(t, err)

		err = client.Stages.AlterDirectoryTable(ctx, sdk.NewAlterDirectoryTableStageRequest(id).
			WithRefresh(*sdk.NewDirectoryTableRefreshRequest().WithSubpath("/")))
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
		stage, _ := testClientHelper().Stage.CreateStage(t)

		foundStage, err := client.Stages.ShowByID(ctx, stage.ID())
		require.NotNil(t, foundStage)
		require.NoError(t, err)

		err = client.Stages.Drop(ctx, sdk.NewDropStageRequest(stage.ID()))
		require.NoError(t, err)

		foundStage, err = client.Stages.ShowByID(ctx, stage.ID())
		require.Nil(t, foundStage)
		require.Error(t, err)
	})

	t.Run("Describe internal", func(t *testing.T) {
		stage, cleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanup)

		stageProperties, err := client.Stages.Describe(ctx, stage.ID())
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

		s3Req := sdk.NewExternalS3StageParamsRequest(awsBucketUrl).
			WithCredentials(*sdk.NewExternalStageS3CredentialsRequest().
				WithAwsKeyId(awsKeyId).
				WithAwsSecretKey(awsSecretKey))

		request := sdk.NewCreateOnS3StageRequest(id, *s3Req)

		_, cleanup := testClientHelper().Stage.CreateStageOnS3WithRequest(t, id, request)
		t.Cleanup(cleanup)

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

		gcsReq := sdk.NewExternalGCSStageParamsRequest(gcsBucketUrl).
			WithStorageIntegration(ids.PrecreatedGcpStorageIntegration)

		request := sdk.NewCreateOnGCSStageRequest(id, *gcsReq)

		_, cleanup := testClientHelper().Stage.CreateStageOnGCSWithRequest(t, id, request)
		t.Cleanup(cleanup)

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

		azureReq := sdk.NewExternalAzureStageParamsRequest(azureBucketUrl).
			WithCredentials(*sdk.NewExternalStageAzureCredentialsRequest(azureSasToken))

		request := sdk.NewCreateOnAzureStageRequest(id, *azureReq)

		_, cleanup := testClientHelper().Stage.CreateStageOnAzureWithRequest(t, id, request)
		t.Cleanup(cleanup)

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
		comment := "show internal test"

		request := sdk.NewCreateInternalStageRequest(id).
			WithDirectoryTableOptions(*sdk.NewInternalDirectoryTableOptionsRequest().WithEnable(true)).
			WithComment(comment)

		stage, cleanup := testClientHelper().Stage.CreateStageWithRequest(t, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.StageFromObject(t, stage).
			HasName(id.Name()).
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasUrl("").
			HasHasCredentials(false).
			HasHasEncryptionKey(false).
			HasComment(comment).
			HasType("INTERNAL").
			HasDirectoryEnabled(true).
			HasOwnerRoleType("ROLE"))
	})
}

func TestInt_StagesShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		stage1, cleanup1 := testClientHelper().Stage.CreateStageWithRequest(t, sdk.NewCreateInternalStageRequest(id1))
		t.Cleanup(cleanup1)
		stage2, cleanup2 := testClientHelper().Stage.CreateStageInSchema(t, schema.ID())
		t.Cleanup(cleanup2)

		// Re-create stage2 with the same name as stage1
		err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(stage2.ID()))
		require.NoError(t, err)
		stage2, cleanup2 = testClientHelper().Stage.CreateStageWithRequest(t, sdk.NewCreateInternalStageRequest(id2))
		t.Cleanup(cleanup2)

		e1, err := client.Stages.ShowByID(ctx, stage1.ID())
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Stages.ShowByID(ctx, stage2.ID())
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}

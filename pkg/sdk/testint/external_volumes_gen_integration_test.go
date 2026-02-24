//go:build non_account_level_tests

package testint

import (
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_ExternalVolumes(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	awsBaseUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsRoleARN := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)
	awsKmsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsExternalId := "123456789"

	gcsBaseUrl := testenvs.GetOrSkipTest(t, testenvs.GcsExternalBucketUrl)
	gcsKmsKeyId := "123456789"

	azureBaseUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl)
	azureTenantId := testenvs.GetOrSkipTest(t, testenvs.AzureExternalTenantId)

	// Storage location structs for testing
	// Note cannot test awsgov on non-gov Snowflake deployments

	s3StorageLocations := []sdk.ExternalVolumeStorageLocationItem{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: "s3_testing_storage_location",
			S3StorageLocationParams: &sdk.S3StorageLocationParams{
				StorageProvider:      sdk.S3StorageProviderS3,
				StorageAwsRoleArn:    awsRoleARN,
				StorageBaseUrl:       awsBaseUrl,
				StorageAwsExternalId: sdk.String(awsExternalId),
				Encryption: &sdk.ExternalVolumeS3Encryption{
					EncryptionType: sdk.S3EncryptionTypeSseKms,
					KmsKeyId:       &awsKmsKeyId,
				},
			},
		}},
	}

	s3StorageLocationsNoneEncryption := []sdk.ExternalVolumeStorageLocationItem{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: "s3_testing_storage_location_none_encryption",
			S3StorageLocationParams: &sdk.S3StorageLocationParams{
				StorageProvider:      sdk.S3StorageProviderS3,
				StorageAwsRoleArn:    awsRoleARN,
				StorageBaseUrl:       awsBaseUrl,
				StorageAwsExternalId: sdk.String(awsExternalId),
				Encryption: &sdk.ExternalVolumeS3Encryption{
					EncryptionType: sdk.S3EncryptionNone,
				},
			},
		}},
	}

	s3StorageLocationsNoEncryption := []sdk.ExternalVolumeStorageLocationItem{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: "s3_testing_storage_location_no_encryption",
			S3StorageLocationParams: &sdk.S3StorageLocationParams{
				StorageProvider:      sdk.S3StorageProviderS3,
				StorageAwsRoleArn:    awsRoleARN,
				StorageBaseUrl:       awsBaseUrl,
				StorageAwsExternalId: sdk.String(awsExternalId),
			},
		}},
	}

	gcsStorageLocationsNoneEncryption := []sdk.ExternalVolumeStorageLocationItem{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: "gcs_testing_storage_location_none_encryption",
			GCSStorageLocationParams: &sdk.GCSStorageLocationParams{
				StorageBaseUrl: gcsBaseUrl,
				Encryption: &sdk.ExternalVolumeGCSEncryption{
					EncryptionType: sdk.GCSEncryptionTypeNone,
				},
			},
		}},
	}

	gcsStorageLocationsNoEncryption := []sdk.ExternalVolumeStorageLocationItem{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: "gcs_testing_storage_location_no_encryption",
			GCSStorageLocationParams: &sdk.GCSStorageLocationParams{
				StorageBaseUrl: gcsBaseUrl,
			},
		}},
	}

	gcsStorageLocations := []sdk.ExternalVolumeStorageLocationItem{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: "gcs_testing_storage_location",
			GCSStorageLocationParams: &sdk.GCSStorageLocationParams{
				StorageBaseUrl: gcsBaseUrl,
				Encryption: &sdk.ExternalVolumeGCSEncryption{
					EncryptionType: sdk.GCSEncryptionTypeSseKms,
					KmsKeyId:       &gcsKmsKeyId,
				},
			},
		}},
	}

	azureStorageLocations := []sdk.ExternalVolumeStorageLocationItem{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: "azure_testing_storage_location",
			AzureStorageLocationParams: &sdk.AzureStorageLocationParams{
				AzureTenantId:  azureTenantId,
				StorageBaseUrl: azureBaseUrl,
			},
		}},
	}

	describeExternalVolume := func(t *testing.T, id sdk.AccountObjectIdentifier) sdk.ExternalVolumeDetails {
		t.Helper()
		externalVolumeProperties, err := client.ExternalVolumes.Describe(ctx, id)
		require.NoError(t, err)
		externalVolumeDetails, err := sdk.ParseExternalVolumeDescribed(externalVolumeProperties)
		require.NoError(t, err)
		return externalVolumeDetails
	}

	createExternalVolume := func(t *testing.T, storageLocations []sdk.ExternalVolumeStorageLocationItem, allowWrites bool, comment *string) sdk.AccountObjectIdentifier {
		t.Helper()

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		req := sdk.NewCreateExternalVolumeRequest(id, storageLocations).
			WithIfNotExists(true).
			WithAllowWrites(allowWrites)

		if comment != nil {
			req = req.WithComment(*comment)
		}

		err := client.ExternalVolumes.Create(ctx, req)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.ExternalVolumes.Drop(ctx, sdk.NewDropExternalVolumeRequest(id).WithIfExists(true))
			require.NoError(t, err)
		})

		return id
	}

	t.Run("Create - S3 Storage Location", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, s3StorageLocations, allowWrites, &comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, externalVolume).
			HasName(id.Name()).
			HasAllowWrites(allowWrites).
			HasComment(comment))
	})

	t.Run("Create - S3 Storage Location empty Comment", func(t *testing.T) {
		allowWrites := true
		emptyComment := ""
		id := createExternalVolume(t, s3StorageLocations, allowWrites, &emptyComment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, externalVolume).
			HasName(id.Name()).
			HasAllowWrites(allowWrites).
			HasComment(emptyComment))
	})

	t.Run("Create - S3 Storage Location No Comment", func(t *testing.T) {
		allowWrites := true
		id := createExternalVolume(t, s3StorageLocations, allowWrites, nil)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, externalVolume).
			HasName(id.Name()).
			HasAllowWrites(allowWrites).
			HasComment(""))
	})

	t.Run("Create - S3 Storage Location None Encryption", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, s3StorageLocationsNoneEncryption, allowWrites, &comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, externalVolume).
			HasName(id.Name()).
			HasAllowWrites(allowWrites).
			HasComment(comment))
	})

	t.Run("Create - S3 Storage Location No Encryption", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, s3StorageLocationsNoEncryption, allowWrites, &comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, externalVolume).
			HasName(id.Name()).
			HasAllowWrites(allowWrites).
			HasComment(comment))
	})

	t.Run("Create - GCS Storage Location", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, gcsStorageLocations, allowWrites, &comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, externalVolume).
			HasName(id.Name()).
			HasAllowWrites(allowWrites).
			HasComment(comment))
	})

	t.Run("Create - GCS Storage Location None Encryption", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, gcsStorageLocationsNoneEncryption, allowWrites, &comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, externalVolume).
			HasName(id.Name()).
			HasAllowWrites(allowWrites).
			HasComment(comment))

		externalVolumeDetails := describeExternalVolume(t, id)
		require.Len(t, externalVolumeDetails.StorageLocations, 1)

		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[0]).
			HasStorageProvider(string(sdk.StorageProviderGCS)).
			HasStorageBaseUrl(gcsStorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.GCSStorageLocationParams.StorageBaseUrl).
			HasStorageAllowedLocations(gcsBaseUrl+"*").
			HasEncryptionType(string(gcsStorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.GCSStorageLocationParams.Encryption.EncryptionType)))

		assertThatObject(t, objectassert.StorageLocationGcsDetailsFromObject(t, externalVolumeDetails.StorageLocations[0].GCSStorageLocation).
			HasEncryptionKmsKeyId("").
			HasStorageGcpServiceAccountNotEmpty())
	})

	t.Run("Create - GCS Storage Location No Encryption", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, gcsStorageLocationsNoEncryption, allowWrites, &comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, externalVolume).
			HasName(id.Name()).
			HasAllowWrites(allowWrites).
			HasComment(comment))
	})

	t.Run("Create - Azure Storage Location", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, azureStorageLocations, allowWrites, &comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, externalVolume).
			HasName(id.Name()).
			HasAllowWrites(allowWrites).
			HasComment(comment))
	})

	t.Run("Create - Multiple Storage Locations", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, append(append(s3StorageLocations, gcsStorageLocationsNoneEncryption...), azureStorageLocations...), allowWrites, &comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, externalVolume).
			HasName(id.Name()).
			HasAllowWrites(allowWrites).
			HasComment(comment))
	})

	t.Run("Alter - remove storage location", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, append(s3StorageLocationsNoneEncryption, gcsStorageLocationsNoneEncryption...), allowWrites, &comment)

		req := sdk.NewAlterExternalVolumeRequest(id).WithRemoveStorageLocation(gcsStorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.Name)

		err := client.ExternalVolumes.Alter(ctx, req)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeDetails(t, id).
			HasActive("").
			HasComment(comment).
			HasAllowWrites(strconv.FormatBool(allowWrites)))

		externalVolumeDetails := describeExternalVolume(t, id)
		require.Len(t, externalVolumeDetails.StorageLocations, 1)

		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[0]).
			HasName(s3StorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(s3StorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageProvider)).
			HasStorageAllowedLocations(awsBaseUrl+"*").
			HasStorageBaseUrl(s3StorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageBaseUrl).
			HasEncryptionType(string(s3StorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.Encryption.EncryptionType)))

		assertThatObject(t, objectassert.StorageLocationS3DetailsFromObject(t, externalVolumeDetails.StorageLocations[0].S3StorageLocation).
			HasStorageAwsRoleArn(s3StorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageAwsRoleArn).
			HasStorageAwsExternalId(*s3StorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageAwsExternalId).
			HasStorageAwsIamUserArnNotEmpty())
	})

	t.Run("Alter - set comment", func(t *testing.T) {
		allowWrites := true
		comment1 := "some comment"
		comment2 := ""
		id := createExternalVolume(t, s3StorageLocationsNoneEncryption, allowWrites, &comment1)

		req := sdk.NewAlterExternalVolumeRequest(id).WithSet(
			*sdk.NewAlterExternalVolumeSetRequest().WithComment(comment2),
		)

		err := client.ExternalVolumes.Alter(ctx, req)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeDetails(t, id).
			HasActive("").
			HasComment(comment2).
			HasAllowWrites(strconv.FormatBool(allowWrites)))

		externalVolumeDetails := describeExternalVolume(t, id)
		require.Len(t, externalVolumeDetails.StorageLocations, 1)

		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[0]).
			HasName(s3StorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(s3StorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageProvider)).
			HasStorageAllowedLocations(awsBaseUrl+"*").
			HasStorageBaseUrl(s3StorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageBaseUrl).
			HasEncryptionType(string(s3StorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.Encryption.EncryptionType)))

		assertThatObject(t, objectassert.StorageLocationS3DetailsFromObject(t, externalVolumeDetails.StorageLocations[0].S3StorageLocation).
			HasStorageAwsRoleArn(s3StorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageAwsRoleArn).
			HasStorageAwsExternalId(*s3StorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageAwsExternalId).
			HasStorageAwsIamUserArnNotEmpty())
	})

	t.Run("Alter - set allow writes", func(t *testing.T) {
		allowWrites := false
		comment := "some comment"
		id := createExternalVolume(t, s3StorageLocations, true, &comment)

		req := sdk.NewAlterExternalVolumeRequest(id).WithSet(
			*sdk.NewAlterExternalVolumeSetRequest().WithAllowWrites(allowWrites),
		)

		err := client.ExternalVolumes.Alter(ctx, req)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeDetails(t, id).
			HasActive("").
			HasComment(comment).
			HasAllowWrites(strconv.FormatBool(allowWrites)))

		externalVolumeDetails := describeExternalVolume(t, id)
		require.Len(t, externalVolumeDetails.StorageLocations, 1)

		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[0]).
			HasName(s3StorageLocations[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(s3StorageLocations[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageProvider)).
			HasStorageAllowedLocations(awsBaseUrl+"*").
			HasStorageBaseUrl(s3StorageLocations[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageBaseUrl).
			HasEncryptionType(string(s3StorageLocations[0].ExternalVolumeStorageLocation.S3StorageLocationParams.Encryption.EncryptionType)))

		assertThatObject(t, objectassert.StorageLocationS3DetailsFromObject(t, externalVolumeDetails.StorageLocations[0].S3StorageLocation).
			HasStorageAwsRoleArn(s3StorageLocations[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageAwsRoleArn).
			HasStorageAwsExternalId(*s3StorageLocations[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageAwsExternalId).
			HasEncryptionKmsKeyId(*s3StorageLocations[0].ExternalVolumeStorageLocation.S3StorageLocationParams.Encryption.KmsKeyId).
			HasStorageAwsIamUserArnNotEmpty())
	})

	t.Run("Alter - add s3 storage location to external volume", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, gcsStorageLocationsNoneEncryption, allowWrites, &comment)

		s3Loc := s3StorageLocations[0].ExternalVolumeStorageLocation
		req := sdk.NewAlterExternalVolumeRequest(id).WithAddStorageLocation(
			*sdk.NewExternalVolumeStorageLocationItemRequest(
				*sdk.NewExternalVolumeStorageLocationRequest(
					s3Loc.Name,
				).WithS3StorageLocationParams(
					*sdk.NewS3StorageLocationParamsRequest(
						s3Loc.S3StorageLocationParams.StorageProvider,
						s3Loc.S3StorageLocationParams.StorageAwsRoleArn,
						s3Loc.S3StorageLocationParams.StorageBaseUrl,
					).WithStorageAwsExternalId(*s3Loc.S3StorageLocationParams.StorageAwsExternalId).
						WithEncryption(
							*sdk.NewExternalVolumeS3EncryptionRequest(s3Loc.S3StorageLocationParams.Encryption.EncryptionType).
								WithKmsKeyId(*s3Loc.S3StorageLocationParams.Encryption.KmsKeyId),
						),
				),
			),
		)

		err := client.ExternalVolumes.Alter(ctx, req)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeDetails(t, id).
			HasActive("").
			HasComment(comment).
			HasAllowWrites(strconv.FormatBool(allowWrites)))

		externalVolumeDetails := describeExternalVolume(t, id)
		require.Len(t, externalVolumeDetails.StorageLocations, 2)

		// Location 0: GCS
		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[0]).
			HasName(gcsStorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderGCS)).
			HasStorageAllowedLocations(gcsBaseUrl+"*").
			HasStorageBaseUrl(gcsStorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.GCSStorageLocationParams.StorageBaseUrl).
			HasEncryptionType(string(gcsStorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.GCSStorageLocationParams.Encryption.EncryptionType)))

		assertThatObject(t, objectassert.StorageLocationGcsDetailsFromObject(t, externalVolumeDetails.StorageLocations[0].GCSStorageLocation).
			HasEncryptionKmsKeyId("").
			HasStorageGcpServiceAccountNotEmpty())

		// Location 1: S3
		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[1]).
			HasName(s3StorageLocations[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(s3StorageLocations[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageProvider)).
			HasStorageAllowedLocations(awsBaseUrl+"*").
			HasStorageBaseUrl(s3StorageLocations[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageBaseUrl).
			HasEncryptionType(string(s3StorageLocations[0].ExternalVolumeStorageLocation.S3StorageLocationParams.Encryption.EncryptionType)))

		assertThatObject(t, objectassert.StorageLocationS3DetailsFromObject(t, externalVolumeDetails.StorageLocations[1].S3StorageLocation).
			HasStorageAwsRoleArn(s3StorageLocations[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageAwsRoleArn).
			HasStorageAwsExternalId(*s3StorageLocations[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageAwsExternalId).
			HasEncryptionKmsKeyId(*s3StorageLocations[0].ExternalVolumeStorageLocation.S3StorageLocationParams.Encryption.KmsKeyId).
			HasStorageAwsIamUserArnNotEmpty())
	})

	t.Run("Describe - multiple storage locations", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(
			t,
			append(append(append(append(append(append(s3StorageLocations, gcsStorageLocationsNoneEncryption...), azureStorageLocations...), s3StorageLocationsNoneEncryption...), gcsStorageLocations...), s3StorageLocationsNoEncryption...), gcsStorageLocationsNoEncryption...),
			allowWrites,
			&comment,
		)

		assertThatObject(t, objectassert.ExternalVolumeDetails(t, id).
			HasActive("").
			HasComment(comment).
			HasAllowWrites(strconv.FormatBool(allowWrites)))

		externalVolumeDetails := describeExternalVolume(t, id)
		require.Len(t, externalVolumeDetails.StorageLocations, 7)

		// Location 0: S3 with KMS encryption
		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[0]).
			HasName(s3StorageLocations[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(s3StorageLocations[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageProvider)).
			HasStorageAllowedLocations(awsBaseUrl+"*").
			HasStorageBaseUrl(s3StorageLocations[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageBaseUrl).
			HasEncryptionType(string(s3StorageLocations[0].ExternalVolumeStorageLocation.S3StorageLocationParams.Encryption.EncryptionType)))

		assertThatObject(t, objectassert.StorageLocationS3DetailsFromObject(t, externalVolumeDetails.StorageLocations[0].S3StorageLocation).
			HasStorageAwsRoleArn(s3StorageLocations[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageAwsRoleArn).
			HasStorageAwsExternalId(*s3StorageLocations[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageAwsExternalId).
			HasEncryptionKmsKeyId(*s3StorageLocations[0].ExternalVolumeStorageLocation.S3StorageLocationParams.Encryption.KmsKeyId).
			HasStorageAwsIamUserArnNotEmpty())

		// Location 1: GCS with NONE encryption
		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[1]).
			HasName(gcsStorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderGCS)).
			HasStorageAllowedLocations(gcsBaseUrl+"*").
			HasStorageBaseUrl(gcsStorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.GCSStorageLocationParams.StorageBaseUrl).
			HasEncryptionType(string(gcsStorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.GCSStorageLocationParams.Encryption.EncryptionType)))

		assertThatObject(t, objectassert.StorageLocationGcsDetailsFromObject(t, externalVolumeDetails.StorageLocations[1].GCSStorageLocation).
			HasEncryptionKmsKeyId("").
			HasStorageGcpServiceAccountNotEmpty())

		// Location 2: Azure
		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[2]).
			HasName(azureStorageLocations[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderAzure)).
			HasStorageAllowedLocations(azureBaseUrl+"*").
			HasStorageBaseUrl(azureStorageLocations[0].ExternalVolumeStorageLocation.AzureStorageLocationParams.StorageBaseUrl).
			HasEncryptionType("NONE"))

		assertThatObject(t, objectassert.StorageLocationAzureDetailsFromObject(t, externalVolumeDetails.StorageLocations[2].AzureStorageLocation).
			HasAzureTenantId(azureStorageLocations[0].ExternalVolumeStorageLocation.AzureStorageLocationParams.AzureTenantId).
			HasAzureMultiTenantAppNameNotEmpty().
			HasAzureConsentUrlNotEmpty())

		// Location 3: S3 with NONE encryption
		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[3]).
			HasName(s3StorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(s3StorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageProvider)).
			HasStorageAllowedLocations(awsBaseUrl+"*").
			HasStorageBaseUrl(s3StorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageBaseUrl).
			HasEncryptionType(string(s3StorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.Encryption.EncryptionType)))

		assertThatObject(t, objectassert.StorageLocationS3DetailsFromObject(t, externalVolumeDetails.StorageLocations[3].S3StorageLocation).
			HasStorageAwsRoleArn(s3StorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageAwsRoleArn).
			HasStorageAwsExternalId(*s3StorageLocationsNoneEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageAwsExternalId).
			HasStorageAwsIamUserArnNotEmpty())

		// Location 4: GCS with KMS encryption
		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[4]).
			HasName(gcsStorageLocations[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderGCS)).
			HasStorageAllowedLocations(gcsBaseUrl+"*").
			HasStorageBaseUrl(gcsStorageLocations[0].ExternalVolumeStorageLocation.GCSStorageLocationParams.StorageBaseUrl).
			HasEncryptionType(string(gcsStorageLocations[0].ExternalVolumeStorageLocation.GCSStorageLocationParams.Encryption.EncryptionType)))

		assertThatObject(t, objectassert.StorageLocationGcsDetailsFromObject(t, externalVolumeDetails.StorageLocations[4].GCSStorageLocation).
			HasEncryptionKmsKeyId(*gcsStorageLocations[0].ExternalVolumeStorageLocation.GCSStorageLocationParams.Encryption.KmsKeyId).
			HasStorageGcpServiceAccountNotEmpty())

		// Location 5: S3 with no encryption specified
		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[5]).
			HasName(s3StorageLocationsNoEncryption[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(s3StorageLocationsNoEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageProvider)).
			HasStorageAllowedLocations(awsBaseUrl+"*").
			HasStorageBaseUrl(s3StorageLocationsNoEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageBaseUrl).
			HasEncryptionType("NONE"))

		assertThatObject(t, objectassert.StorageLocationS3DetailsFromObject(t, externalVolumeDetails.StorageLocations[5].S3StorageLocation).
			HasStorageAwsRoleArn(s3StorageLocationsNoEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageAwsRoleArn).
			HasStorageAwsExternalId(*s3StorageLocationsNoEncryption[0].ExternalVolumeStorageLocation.S3StorageLocationParams.StorageAwsExternalId).
			HasStorageAwsIamUserArnNotEmpty())

		// Location 6: GCS with no encryption specified
		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[6]).
			HasName(gcsStorageLocationsNoEncryption[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderGCS)).
			HasStorageAllowedLocations(gcsBaseUrl+"*").
			HasStorageBaseUrl(gcsStorageLocationsNoEncryption[0].ExternalVolumeStorageLocation.GCSStorageLocationParams.StorageBaseUrl).
			HasEncryptionType("NONE"))

		assertThatObject(t, objectassert.StorageLocationGcsDetailsFromObject(t, externalVolumeDetails.StorageLocations[6].GCSStorageLocation).
			HasEncryptionKmsKeyId("").
			HasStorageGcpServiceAccountNotEmpty())
	})

	t.Run("Show with like", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, s3StorageLocations, allowWrites, &comment)
		name := id.Name()
		req := sdk.NewShowExternalVolumeRequest().WithLike(sdk.Like{Pattern: &name})

		externalVolumes, err := client.ExternalVolumes.Show(ctx, req)
		require.NoError(t, err)

		assert.Len(t, externalVolumes, 1)
		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, &externalVolumes[0]).
			HasName(id.Name()).
			HasAllowWrites(allowWrites).
			HasComment(comment))
	})
}

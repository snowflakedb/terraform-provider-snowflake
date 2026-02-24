//go:build non_account_level_tests

package testint

import (
	"strconv"
	"strings"
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

	s3CompatBaseUrl := strings.Replace(awsBaseUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	// Storage location structs for testing
	// TODO(SNOW-2356128): Test use_privatelink_endpoint for azure

	s3StorageLocationsBasic := []sdk.ExternalVolumeStorageLocationItem{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: "s3_testing_storage_location_basic",
			S3StorageLocationParams: &sdk.S3StorageLocationParams{
				StorageProvider:   sdk.S3StorageProviderS3,
				StorageAwsRoleArn: awsRoleARN,
				StorageBaseUrl:    awsBaseUrl,
			},
		}},
	}

	s3StorageLocationsComplete := []sdk.ExternalVolumeStorageLocationItem{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: "s3_testing_storage_location_complete",
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

	s3GovBaseUrl := strings.Replace(awsBaseUrl, "s3://", "s3gov://", 1)
	s3GovStorageLocationsBasic := []sdk.ExternalVolumeStorageLocationItem{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: "s3gov_testing_storage_location_basic",
			S3StorageLocationParams: &sdk.S3StorageLocationParams{
				StorageProvider:   sdk.S3StorageProviderS3GOV,
				StorageAwsRoleArn: awsRoleARN,
				StorageBaseUrl:    s3GovBaseUrl,
			},
		}},
	}

	s3GovStorageLocationsComplete := []sdk.ExternalVolumeStorageLocationItem{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: "s3gov_testing_storage_location_complete",
			S3StorageLocationParams: &sdk.S3StorageLocationParams{
				StorageProvider:      sdk.S3StorageProviderS3GOV,
				StorageAwsRoleArn:    awsRoleARN,
				StorageBaseUrl:       s3GovBaseUrl,
				StorageAwsExternalId: sdk.String(awsExternalId),
				Encryption: &sdk.ExternalVolumeS3Encryption{
					EncryptionType: sdk.S3EncryptionTypeSseKms,
					KmsKeyId:       &awsKmsKeyId,
				},
			},
		}},
	}

	gcsStorageLocationsBasic := []sdk.ExternalVolumeStorageLocationItem{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: "gcs_testing_storage_location_basic",
			GCSStorageLocationParams: &sdk.GCSStorageLocationParams{
				StorageBaseUrl: gcsBaseUrl,
			},
		}},
	}

	gcsStorageLocationsComplete := []sdk.ExternalVolumeStorageLocationItem{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: "gcs_testing_storage_location_complete",
			GCSStorageLocationParams: &sdk.GCSStorageLocationParams{
				StorageBaseUrl: gcsBaseUrl,
				Encryption: &sdk.ExternalVolumeGCSEncryption{
					EncryptionType: sdk.GCSEncryptionTypeSseKms,
					KmsKeyId:       &gcsKmsKeyId,
				},
			},
		}},
	}

	s3CompatStorageLocationsBasic := []sdk.ExternalVolumeStorageLocationItem{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: "s3compat_testing_storage_location_basic",
			S3CompatStorageLocationParams: &sdk.S3CompatStorageLocationParams{
				StorageBaseUrl:  s3CompatBaseUrl,
				StorageEndpoint: s3CompatEndpoint,
				Credentials: &sdk.ExternalVolumeS3CompatCredentials{
					AwsKeyId:     awsKmsKeyId,
					AwsSecretKey: awsSecretKey,
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

	t.Run("Create - S3 - basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ExternalVolumes.Create(ctx, sdk.NewCreateExternalVolumeRequest(id, s3StorageLocationsBasic))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ExternalVolume.DropFunc(t, id))

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, externalVolume).
			HasName(id.Name()).
			HasAllowWrites(true).
			HasComment(""))

		externalVolumeDetails := describeExternalVolume(t, id)
		require.Len(t, externalVolumeDetails.StorageLocations, 1)

		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[0]).
			HasName(s3StorageLocationsBasic[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderS3)).
			HasStorageBaseUrl(awsBaseUrl).
			HasStorageAllowedLocations(awsBaseUrl+"*").
			HasEncryptionType("NONE"))

		assertThatObject(t, objectassert.StorageLocationS3DetailsFromObject(t, externalVolumeDetails.StorageLocations[0].S3StorageLocation).
			HasStorageAwsRoleArn(awsRoleARN).
			HasStorageAwsIamUserArnNotEmpty())
	})

	t.Run("Create - S3 - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ExternalVolumes.Create(ctx, sdk.NewCreateExternalVolumeRequest(id, s3StorageLocationsComplete).
			WithIfNotExists(true).
			WithAllowWrites(false).
			WithComment("some comment"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ExternalVolume.DropFunc(t, id))

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, externalVolume).
			HasName(id.Name()).
			HasAllowWrites(false).
			HasComment("some comment"))

		externalVolumeDetails := describeExternalVolume(t, id)
		require.Len(t, externalVolumeDetails.StorageLocations, 1)

		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[0]).
			HasName(s3StorageLocationsComplete[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderS3)).
			HasStorageBaseUrl(awsBaseUrl).
			HasStorageAllowedLocations(awsBaseUrl+"*").
			HasEncryptionType(string(sdk.S3EncryptionTypeSseKms)))

		assertThatObject(t, objectassert.StorageLocationS3DetailsFromObject(t, externalVolumeDetails.StorageLocations[0].S3StorageLocation).
			HasStorageAwsRoleArn(awsRoleARN).
			HasStorageAwsExternalId(awsExternalId).
			HasEncryptionKmsKeyId(awsKmsKeyId).
			HasStorageAwsIamUserArnNotEmpty())
	})

	t.Run("Create - S3Gov - basic", func(t *testing.T) {
		if testenvs.GetSnowflakeEnvironmentWithProdDefault() != testenvs.SnowflakePreProdGovEnvironment {
			t.Skip("Skipping S3Gov test, requires gov deployment")
		}

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ExternalVolumes.Create(ctx, sdk.NewCreateExternalVolumeRequest(id, s3GovStorageLocationsBasic).
			WithIfNotExists(true))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ExternalVolume.DropFunc(t, id))

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, externalVolume).
			HasName(id.Name()).
			HasAllowWrites(true).
			HasComment(""))

		externalVolumeDetails := describeExternalVolume(t, id)
		require.Len(t, externalVolumeDetails.StorageLocations, 1)

		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[0]).
			HasName(s3GovStorageLocationsBasic[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderS3GOV)).
			HasStorageBaseUrl(s3GovBaseUrl).
			HasStorageAllowedLocations(s3GovBaseUrl+"*").
			HasEncryptionType("NONE"))

		assertThatObject(t, objectassert.StorageLocationS3DetailsFromObject(t, externalVolumeDetails.StorageLocations[0].S3StorageLocation).
			HasStorageAwsRoleArn(awsRoleARN).
			HasStorageAwsIamUserArnNotEmpty())
	})

	t.Run("Create - S3Gov - complete", func(t *testing.T) {
		if testenvs.GetSnowflakeEnvironmentWithProdDefault() != testenvs.SnowflakePreProdGovEnvironment {
			t.Skip("Skipping S3Gov test, requires gov deployment")
		}

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ExternalVolumes.Create(ctx, sdk.NewCreateExternalVolumeRequest(id, s3GovStorageLocationsComplete).
			WithIfNotExists(true).
			WithAllowWrites(true).
			WithComment("some comment"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ExternalVolume.DropFunc(t, id))

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, externalVolume).
			HasName(id.Name()).
			HasAllowWrites(true).
			HasComment("some comment"))

		externalVolumeDetails := describeExternalVolume(t, id)
		require.Len(t, externalVolumeDetails.StorageLocations, 1)

		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[0]).
			HasName(s3GovStorageLocationsComplete[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderS3GOV)).
			HasStorageBaseUrl(s3GovBaseUrl).
			HasStorageAllowedLocations(s3GovBaseUrl+"*").
			HasEncryptionType(string(sdk.S3EncryptionTypeSseKms)))

		assertThatObject(t, objectassert.StorageLocationS3DetailsFromObject(t, externalVolumeDetails.StorageLocations[0].S3StorageLocation).
			HasStorageAwsRoleArn(awsRoleARN).
			HasStorageAwsExternalId(awsExternalId).
			HasEncryptionKmsKeyId(awsKmsKeyId).
			HasStorageAwsIamUserArnNotEmpty())
	})

	t.Run("Create - GCS - basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ExternalVolumes.Create(ctx, sdk.NewCreateExternalVolumeRequest(id, gcsStorageLocationsBasic))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ExternalVolume.DropFunc(t, id))

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, externalVolume).
			HasName(id.Name()).
			HasAllowWrites(true).
			HasComment(""))

		externalVolumeDetails := describeExternalVolume(t, id)
		require.Len(t, externalVolumeDetails.StorageLocations, 1)

		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[0]).
			HasName(gcsStorageLocationsBasic[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderGCS)).
			HasStorageBaseUrl(gcsBaseUrl).
			HasStorageAllowedLocations(gcsBaseUrl+"*").
			HasEncryptionType(string(sdk.GCSEncryptionTypeNone)))

		assertThatObject(t, objectassert.StorageLocationGcsDetailsFromObject(t, externalVolumeDetails.StorageLocations[0].GCSStorageLocation).
			HasEncryptionKmsKeyId("").
			HasStorageGcpServiceAccountNotEmpty())
	})

	t.Run("Create - GCS - all fields", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ExternalVolumes.Create(ctx, sdk.NewCreateExternalVolumeRequest(id, gcsStorageLocationsComplete).
			WithIfNotExists(true).
			WithAllowWrites(false).
			WithComment("some comment"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ExternalVolume.DropFunc(t, id))

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, externalVolume).
			HasName(id.Name()).
			HasAllowWrites(false).
			HasComment("some comment"))

		externalVolumeDetails := describeExternalVolume(t, id)
		require.Len(t, externalVolumeDetails.StorageLocations, 1)

		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[0]).
			HasName(gcsStorageLocationsComplete[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderGCS)).
			HasStorageBaseUrl(gcsBaseUrl).
			HasStorageAllowedLocations(gcsBaseUrl+"*").
			HasEncryptionType(string(sdk.GCSEncryptionTypeSseKms)))

		assertThatObject(t, objectassert.StorageLocationGcsDetailsFromObject(t, externalVolumeDetails.StorageLocations[0].GCSStorageLocation).
			HasEncryptionKmsKeyId(gcsKmsKeyId).
			HasStorageGcpServiceAccountNotEmpty())
	})

	t.Run("Create - Azure", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ExternalVolumes.Create(ctx, sdk.NewCreateExternalVolumeRequest(id, azureStorageLocations).
			WithIfNotExists(true).
			WithAllowWrites(true).
			WithComment("some comment"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ExternalVolume.DropFunc(t, id))

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, externalVolume).
			HasName(id.Name()).
			HasAllowWrites(true).
			HasComment("some comment"))

		externalVolumeDetails := describeExternalVolume(t, id)
		require.Len(t, externalVolumeDetails.StorageLocations, 1)

		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[0]).
			HasName(azureStorageLocations[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderAzure)).
			HasStorageBaseUrl(azureBaseUrl).
			HasStorageAllowedLocations(azureBaseUrl+"*").
			HasEncryptionType("NONE"))

		assertThatObject(t, objectassert.StorageLocationAzureDetailsFromObject(t, externalVolumeDetails.StorageLocations[0].AzureStorageLocation).
			HasAzureTenantId(azureTenantId).
			HasAzureMultiTenantAppNameNotEmpty().
			HasAzureConsentUrlNotEmpty())
	})

	t.Run("Create - S3Compat - basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ExternalVolumes.Create(ctx, sdk.NewCreateExternalVolumeRequest(id, s3CompatStorageLocationsBasic))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ExternalVolume.DropFunc(t, id))

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, externalVolume).
			HasName(id.Name()).
			HasAllowWrites(true).
			HasComment(""))

		externalVolumeDetails := describeExternalVolume(t, id)
		require.Len(t, externalVolumeDetails.StorageLocations, 1)

		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[0]).
			HasName(s3CompatStorageLocationsBasic[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderS3Compatible)).
			HasStorageBaseUrl(s3CompatBaseUrl).
			HasStorageAllowedLocations(s3CompatBaseUrl+"*").
			HasEncryptionType("NONE"))

		assertThatObject(t, objectassert.StorageLocationS3CompatDetailsFromObject(t, externalVolumeDetails.StorageLocations[0].S3CompatStorageLocation).
			HasStorageEndpoint(s3CompatEndpoint).
			HasAwsAccessKeyId(awsKmsKeyId))
	})

	t.Run("Alter - remove storage location", func(t *testing.T) {
		comment := "some comment"
		id, cleanup := testClientHelper().ExternalVolume.CreateWithRequest(t, sdk.NewCreateExternalVolumeRequest(
			testClientHelper().Ids.RandomAccountObjectIdentifier(),
			append(s3StorageLocationsBasic, gcsStorageLocationsBasic...),
		).WithIfNotExists(true).WithAllowWrites(true).WithComment(comment))
		t.Cleanup(cleanup)

		req := sdk.NewAlterExternalVolumeRequest(id).WithRemoveStorageLocation(gcsStorageLocationsBasic[0].ExternalVolumeStorageLocation.Name)

		err := client.ExternalVolumes.Alter(ctx, req)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeDetails(t, id).
			HasActive("").
			HasComment(comment).
			HasAllowWrites(strconv.FormatBool(true)))

		externalVolumeDetails := describeExternalVolume(t, id)
		// Only one storage location should be left (s3)
		require.Len(t, externalVolumeDetails.StorageLocations, 1)

		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[0]).
			HasName(s3StorageLocationsBasic[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderS3)).
			HasStorageAllowedLocations(awsBaseUrl+"*").
			HasStorageBaseUrl(awsBaseUrl).
			HasEncryptionType("NONE"))

		assertThatObject(t, objectassert.StorageLocationS3DetailsFromObject(t, externalVolumeDetails.StorageLocations[0].S3StorageLocation).
			HasStorageAwsRoleArn(awsRoleARN).
			HasStorageAwsIamUserArnNotEmpty())
	})

	t.Run("Alter - set comment", func(t *testing.T) {
		id, cleanup := testClientHelper().ExternalVolume.Create(t)
		t.Cleanup(cleanup)

		newComment := "comment"
		req := sdk.NewAlterExternalVolumeRequest(id).WithSet(
			*sdk.NewAlterExternalVolumeSetRequest().WithComment(newComment),
		)

		err := client.ExternalVolumes.Alter(ctx, req)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeDetails(t, id).
			HasActive("").
			HasComment(newComment).
			HasAllowWrites(strconv.FormatBool(true)))
	})

	t.Run("Alter - set allow writes", func(t *testing.T) {
		id, cleanup := testClientHelper().ExternalVolume.Create(t)
		t.Cleanup(cleanup)

		req := sdk.NewAlterExternalVolumeRequest(id).WithSet(
			*sdk.NewAlterExternalVolumeSetRequest().WithAllowWrites(false),
		)

		err := client.ExternalVolumes.Alter(ctx, req)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeDetails(t, id).
			HasActive("").
			HasComment("").
			HasAllowWrites(strconv.FormatBool(false)))
	})

	t.Run("Alter - add s3 storage location to external volume", func(t *testing.T) {
		comment := "some comment"
		id, cleanup := testClientHelper().ExternalVolume.CreateWithRequest(t, sdk.NewCreateExternalVolumeRequest(
			testClientHelper().Ids.RandomAccountObjectIdentifier(),
			gcsStorageLocationsBasic,
		).WithIfNotExists(true).WithAllowWrites(true).WithComment(comment))
		t.Cleanup(cleanup)

		s3Loc := s3StorageLocationsComplete[0].ExternalVolumeStorageLocation
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
			HasAllowWrites(strconv.FormatBool(true)))

		externalVolumeDetails := describeExternalVolume(t, id)
		require.Len(t, externalVolumeDetails.StorageLocations, 2)

		// Location 0: GCS
		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[0]).
			HasName(gcsStorageLocationsBasic[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderGCS)).
			HasStorageAllowedLocations(gcsBaseUrl+"*").
			HasStorageBaseUrl(gcsBaseUrl).
			HasEncryptionType("NONE"))

		assertThatObject(t, objectassert.StorageLocationGcsDetailsFromObject(t, externalVolumeDetails.StorageLocations[0].GCSStorageLocation).
			HasEncryptionKmsKeyId("").
			HasStorageGcpServiceAccountNotEmpty())

		// Location 1: S3
		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[1]).
			HasName(s3StorageLocationsComplete[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderS3)).
			HasStorageAllowedLocations(awsBaseUrl+"*").
			HasStorageBaseUrl(awsBaseUrl).
			HasEncryptionType(string(sdk.S3EncryptionTypeSseKms)))

		assertThatObject(t, objectassert.StorageLocationS3DetailsFromObject(t, externalVolumeDetails.StorageLocations[1].S3StorageLocation).
			HasStorageAwsRoleArn(awsRoleARN).
			HasStorageAwsExternalId(awsExternalId).
			HasEncryptionKmsKeyId(awsKmsKeyId).
			HasStorageAwsIamUserArnNotEmpty())
	})

	t.Run("Alter - add GCS storage location to external volume", func(t *testing.T) {
		comment := "some comment"
		id, cleanup := testClientHelper().ExternalVolume.CreateWithRequest(t, sdk.NewCreateExternalVolumeRequest(
			testClientHelper().Ids.RandomAccountObjectIdentifier(),
			s3StorageLocationsBasic,
		).WithIfNotExists(true).WithAllowWrites(true).WithComment(comment))
		t.Cleanup(cleanup)

		gcsLoc := gcsStorageLocationsComplete[0].ExternalVolumeStorageLocation
		req := sdk.NewAlterExternalVolumeRequest(id).WithAddStorageLocation(
			*sdk.NewExternalVolumeStorageLocationItemRequest(
				*sdk.NewExternalVolumeStorageLocationRequest(
					gcsLoc.Name,
				).WithGCSStorageLocationParams(
					*sdk.NewGCSStorageLocationParamsRequest(
						gcsLoc.GCSStorageLocationParams.StorageBaseUrl,
					).WithEncryption(
						*sdk.NewExternalVolumeGCSEncryptionRequest(gcsLoc.GCSStorageLocationParams.Encryption.EncryptionType).
							WithKmsKeyId(*gcsLoc.GCSStorageLocationParams.Encryption.KmsKeyId),
					),
				),
			),
		)

		err := client.ExternalVolumes.Alter(ctx, req)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeDetails(t, id).
			HasActive("").
			HasComment(comment).
			HasAllowWrites(strconv.FormatBool(true)))

		externalVolumeDetails := describeExternalVolume(t, id)
		require.Len(t, externalVolumeDetails.StorageLocations, 2)

		// Location 0: S3
		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[0]).
			HasName(s3StorageLocationsBasic[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderS3)).
			HasStorageAllowedLocations(awsBaseUrl+"*").
			HasStorageBaseUrl(awsBaseUrl).
			HasEncryptionType("NONE"))

		assertThatObject(t, objectassert.StorageLocationS3DetailsFromObject(t, externalVolumeDetails.StorageLocations[0].S3StorageLocation).
			HasStorageAwsRoleArn(awsRoleARN).
			HasStorageAwsIamUserArnNotEmpty())

		// Location 1: GCS
		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[1]).
			HasName(gcsStorageLocationsComplete[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderGCS)).
			HasStorageAllowedLocations(gcsBaseUrl+"*").
			HasStorageBaseUrl(gcsBaseUrl).
			HasEncryptionType(string(sdk.GCSEncryptionTypeSseKms)))

		assertThatObject(t, objectassert.StorageLocationGcsDetailsFromObject(t, externalVolumeDetails.StorageLocations[1].GCSStorageLocation).
			HasEncryptionKmsKeyId(gcsKmsKeyId).
			HasStorageGcpServiceAccountNotEmpty())
	})

	t.Run("Alter - add Azure storage location to external volume", func(t *testing.T) {
		id, cleanup := testClientHelper().ExternalVolume.CreateWithRequest(t, sdk.NewCreateExternalVolumeRequest(
			testClientHelper().Ids.RandomAccountObjectIdentifier(),
			s3StorageLocationsBasic,
		))
		t.Cleanup(cleanup)

		azureLoc := azureStorageLocations[0].ExternalVolumeStorageLocation
		req := sdk.NewAlterExternalVolumeRequest(id).WithAddStorageLocation(
			*sdk.NewExternalVolumeStorageLocationItemRequest(
				*sdk.NewExternalVolumeStorageLocationRequest(
					azureLoc.Name,
				).WithAzureStorageLocationParams(
					*sdk.NewAzureStorageLocationParamsRequest(
						azureLoc.AzureStorageLocationParams.AzureTenantId,
						azureLoc.AzureStorageLocationParams.StorageBaseUrl,
					),
				),
			),
		)

		err := client.ExternalVolumes.Alter(ctx, req)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ExternalVolumeDetails(t, id).
			HasActive("").
			HasComment("").
			HasAllowWrites(strconv.FormatBool(true)))

		externalVolumeDetails := describeExternalVolume(t, id)
		require.Len(t, externalVolumeDetails.StorageLocations, 2)

		// Location 0: S3
		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[0]).
			HasName(s3StorageLocationsBasic[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderS3)).
			HasStorageAllowedLocations(awsBaseUrl+"*").
			HasStorageBaseUrl(awsBaseUrl).
			HasEncryptionType("NONE"))

		assertThatObject(t, objectassert.StorageLocationS3DetailsFromObject(t, externalVolumeDetails.StorageLocations[0].S3StorageLocation).
			HasStorageAwsRoleArn(awsRoleARN).
			HasStorageAwsIamUserArnNotEmpty())

		// Location 1: Azure
		assertThatObject(t, objectassert.ExternalVolumeStorageLocationDetailsFromObject(t, &externalVolumeDetails.StorageLocations[1]).
			HasName(azureStorageLocations[0].ExternalVolumeStorageLocation.Name).
			HasStorageProvider(string(sdk.StorageProviderAzure)).
			HasStorageAllowedLocations(azureBaseUrl+"*").
			HasStorageBaseUrl(azureBaseUrl).
			HasEncryptionType("NONE"))

		assertThatObject(t, objectassert.StorageLocationAzureDetailsFromObject(t, externalVolumeDetails.StorageLocations[1].AzureStorageLocation).
			HasAzureTenantId(azureTenantId).
			HasAzureMultiTenantAppNameNotEmpty().
			HasAzureConsentUrlNotEmpty())
	})

	t.Run("Show with like", func(t *testing.T) {
		id, cleanup := testClientHelper().ExternalVolume.Create(t)
		t.Cleanup(cleanup)

		name := id.Name()
		req := sdk.NewShowExternalVolumeRequest().WithLike(sdk.Like{Pattern: &name})

		externalVolumes, err := client.ExternalVolumes.Show(ctx, req)
		require.NoError(t, err)

		assert.Len(t, externalVolumes, 1)
		assertThatObject(t, objectassert.ExternalVolumeFromObject(t, &externalVolumes[0]).
			HasAllowWrites(true).
			HasComment("").
			HasName(id.Name()))
	})
}

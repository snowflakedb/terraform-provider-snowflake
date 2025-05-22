//go:build !account_level_tests

package testint

import (
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
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

	assertExternalVolumeShowResult := func(t *testing.T, s *sdk.ExternalVolume, name sdk.AccountObjectIdentifier, allowWrites bool, comment string) {
		t.Helper()
		assert.Equal(t, name.Name(), s.Name)
		assert.Equal(t, allowWrites, s.AllowWrites)
		assert.Equal(t, s.Comment, comment)
	}

	// Storage location structs for testing
	// Note cannot test awsgov on non-gov Snowflake deployments

	s3StorageLocations := []sdk.ExternalVolumeStorageLocation{
		{
			S3StorageLocationParams: &sdk.S3StorageLocationParams{
				Name:                 "s3_testing_storage_location",
				StorageProvider:      sdk.S3StorageProviderS3,
				StorageAwsRoleArn:    awsRoleARN,
				StorageBaseUrl:       awsBaseUrl,
				StorageAwsExternalId: sdk.String(awsExternalId),
				Encryption: &sdk.ExternalVolumeS3Encryption{
					Type:     sdk.S3EncryptionTypeSseKms,
					KmsKeyId: &awsKmsKeyId,
				},
			},
		},
	}

	s3StorageLocationsNoneEncryption := []sdk.ExternalVolumeStorageLocation{
		{
			S3StorageLocationParams: &sdk.S3StorageLocationParams{
				Name:                 "s3_testing_storage_location_none_encryption",
				StorageProvider:      sdk.S3StorageProviderS3,
				StorageAwsRoleArn:    awsRoleARN,
				StorageBaseUrl:       awsBaseUrl,
				StorageAwsExternalId: sdk.String(awsExternalId),
				Encryption: &sdk.ExternalVolumeS3Encryption{
					Type: sdk.S3EncryptionNone,
				},
			},
		},
	}

	s3StorageLocationsNoEncryption := []sdk.ExternalVolumeStorageLocation{
		{
			S3StorageLocationParams: &sdk.S3StorageLocationParams{
				Name:                 "s3_testing_storage_location_no_encryption",
				StorageProvider:      sdk.S3StorageProviderS3,
				StorageAwsRoleArn:    awsRoleARN,
				StorageBaseUrl:       awsBaseUrl,
				StorageAwsExternalId: sdk.String(awsExternalId),
			},
		},
	}

	gcsStorageLocationsNoneEncryption := []sdk.ExternalVolumeStorageLocation{
		{
			GCSStorageLocationParams: &sdk.GCSStorageLocationParams{
				Name:           "gcs_testing_storage_location_none_encryption",
				StorageBaseUrl: gcsBaseUrl,
				Encryption: &sdk.ExternalVolumeGCSEncryption{
					Type: sdk.GCSEncryptionTypeNone,
				},
			},
		},
	}

	gcsStorageLocationsNoEncryption := []sdk.ExternalVolumeStorageLocation{
		{
			GCSStorageLocationParams: &sdk.GCSStorageLocationParams{
				Name:           "gcs_testing_storage_location_no_encryption",
				StorageBaseUrl: gcsBaseUrl,
			},
		},
	}

	gcsStorageLocations := []sdk.ExternalVolumeStorageLocation{
		{
			GCSStorageLocationParams: &sdk.GCSStorageLocationParams{
				Name:           "gcs_testing_storage_location",
				StorageBaseUrl: gcsBaseUrl,
				Encryption: &sdk.ExternalVolumeGCSEncryption{
					Type:     sdk.GCSEncryptionTypeSseKms,
					KmsKeyId: &gcsKmsKeyId,
				},
			},
		},
	}

	azureStorageLocations := []sdk.ExternalVolumeStorageLocation{
		{
			AzureStorageLocationParams: &sdk.AzureStorageLocationParams{
				Name:           "azure_testing_storage_location",
				AzureTenantId:  azureTenantId,
				StorageBaseUrl: azureBaseUrl,
			},
		},
	}

	createExternalVolume := func(t *testing.T, storageLocations []sdk.ExternalVolumeStorageLocation, allowWrites bool, comment *string) sdk.AccountObjectIdentifier {
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

		assertExternalVolumeShowResult(t, externalVolume, id, allowWrites, comment)
	})

	t.Run("Create - S3 Storage Location empty Comment", func(t *testing.T) {
		allowWrites := true
		emptyComment := ""
		id := createExternalVolume(t, s3StorageLocations, allowWrites, &emptyComment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertExternalVolumeShowResult(t, externalVolume, id, allowWrites, emptyComment)
	})

	t.Run("Create - S3 Storage Location No Comment", func(t *testing.T) {
		allowWrites := true
		id := createExternalVolume(t, s3StorageLocations, allowWrites, nil)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertExternalVolumeShowResult(t, externalVolume, id, allowWrites, "")
	})

	t.Run("Create - S3 Storage Location None Encryption", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, s3StorageLocationsNoneEncryption, allowWrites, &comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertExternalVolumeShowResult(t, externalVolume, id, allowWrites, comment)
	})

	t.Run("Create - S3 Storage Location No Encryption", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, s3StorageLocationsNoEncryption, allowWrites, &comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertExternalVolumeShowResult(t, externalVolume, id, allowWrites, comment)
	})

	t.Run("Create - GCS Storage Location", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, gcsStorageLocations, allowWrites, &comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertExternalVolumeShowResult(t, externalVolume, id, allowWrites, comment)
	})

	t.Run("Create - GCS Storage Location None Encryption", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, gcsStorageLocationsNoneEncryption, allowWrites, &comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertExternalVolumeShowResult(t, externalVolume, id, allowWrites, comment)
	})

	t.Run("Create - GCS Storage Location No Encryption", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, gcsStorageLocationsNoEncryption, allowWrites, &comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertExternalVolumeShowResult(t, externalVolume, id, allowWrites, comment)
	})

	t.Run("Create - Azure Storage Location", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, azureStorageLocations, allowWrites, &comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertExternalVolumeShowResult(t, externalVolume, id, allowWrites, comment)
	})

	t.Run("Create - Multiple Storage Locations", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, append(append(s3StorageLocations, gcsStorageLocationsNoneEncryption...), azureStorageLocations...), allowWrites, &comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertExternalVolumeShowResult(t, externalVolume, id, allowWrites, comment)
	})

	t.Run("Alter - remove storage location", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, append(s3StorageLocationsNoneEncryption, gcsStorageLocationsNoneEncryption...), allowWrites, &comment)

		req := sdk.NewAlterExternalVolumeRequest(id).WithRemoveStorageLocation(gcsStorageLocationsNoneEncryption[0].GCSStorageLocationParams.Name)

		err := client.ExternalVolumes.Alter(ctx, req)
		require.NoError(t, err)

		props, err := client.ExternalVolumes.Describe(ctx, id)
		require.NoError(t, err)

		parsedExternalVolumeDescribed, err := helpers.ParseExternalVolumeDescribed(props)
		require.NoError(t, err)
		expectedParsedExternalVolumeDescribed := helpers.ParsedExternalVolumeDescribed{
			StorageLocations: []helpers.StorageLocation{
				{
					Name:                 s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.Name,
					StorageProvider:      string(s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageProvider),
					StorageBaseUrl:       s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageBaseUrl,
					StorageAwsRoleArn:    s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageAwsRoleArn,
					StorageAwsExternalId: *s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageAwsExternalId,
					EncryptionType:       string(s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.Encryption.Type),
					EncryptionKmsKeyId:   "",
					AzureTenantId:        "",
				},
			},
			Active:      "",
			Comment:     comment,
			AllowWrites: strconv.FormatBool(allowWrites),
		}

		assert.Equal(t, expectedParsedExternalVolumeDescribed, parsedExternalVolumeDescribed)
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

		props, err := client.ExternalVolumes.Describe(ctx, id)
		require.NoError(t, err)

		parsedExternalVolumeDescribed, err := helpers.ParseExternalVolumeDescribed(props)
		require.NoError(t, err)
		expectedParsedExternalVolumeDescribed := helpers.ParsedExternalVolumeDescribed{
			StorageLocations: []helpers.StorageLocation{
				{
					Name:                 s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.Name,
					StorageProvider:      string(s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageProvider),
					StorageBaseUrl:       s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageBaseUrl,
					StorageAwsRoleArn:    s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageAwsRoleArn,
					StorageAwsExternalId: *s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageAwsExternalId,
					EncryptionType:       string(s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.Encryption.Type),
					EncryptionKmsKeyId:   "",
					AzureTenantId:        "",
				},
			},
			Active:      "",
			Comment:     comment2,
			AllowWrites: strconv.FormatBool(allowWrites),
		}

		assert.Equal(t, expectedParsedExternalVolumeDescribed, parsedExternalVolumeDescribed)
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

		props, err := client.ExternalVolumes.Describe(ctx, id)
		require.NoError(t, err)

		parsedExternalVolumeDescribed, err := helpers.ParseExternalVolumeDescribed(props)
		require.NoError(t, err)
		expectedParsedExternalVolumeDescribed := helpers.ParsedExternalVolumeDescribed{
			StorageLocations: []helpers.StorageLocation{
				{
					Name:                 s3StorageLocations[0].S3StorageLocationParams.Name,
					StorageProvider:      string(s3StorageLocations[0].S3StorageLocationParams.StorageProvider),
					StorageBaseUrl:       s3StorageLocations[0].S3StorageLocationParams.StorageBaseUrl,
					StorageAwsRoleArn:    s3StorageLocations[0].S3StorageLocationParams.StorageAwsRoleArn,
					StorageAwsExternalId: *s3StorageLocations[0].S3StorageLocationParams.StorageAwsExternalId,
					EncryptionType:       string(s3StorageLocations[0].S3StorageLocationParams.Encryption.Type),
					EncryptionKmsKeyId:   *s3StorageLocations[0].S3StorageLocationParams.Encryption.KmsKeyId,
					AzureTenantId:        "",
				},
			},
			Active:      "",
			Comment:     comment,
			AllowWrites: strconv.FormatBool(allowWrites),
		}

		assert.Equal(t, expectedParsedExternalVolumeDescribed, parsedExternalVolumeDescribed)
	})

	t.Run("Alter - add s3 storage location to external volume", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, gcsStorageLocationsNoneEncryption, allowWrites, &comment)

		req := sdk.NewAlterExternalVolumeRequest(id).WithAddStorageLocation(
			*sdk.NewExternalVolumeStorageLocationRequest().WithS3StorageLocationParams(
				*sdk.NewS3StorageLocationParamsRequest(
					s3StorageLocations[0].S3StorageLocationParams.Name,
					s3StorageLocations[0].S3StorageLocationParams.StorageProvider,
					s3StorageLocations[0].S3StorageLocationParams.StorageAwsRoleArn,
					s3StorageLocations[0].S3StorageLocationParams.StorageBaseUrl,
				).WithStorageAwsExternalId(*s3StorageLocations[0].S3StorageLocationParams.StorageAwsExternalId).
					WithEncryption(
						*sdk.NewExternalVolumeS3EncryptionRequest(s3StorageLocations[0].S3StorageLocationParams.Encryption.Type).
							WithKmsKeyId(*s3StorageLocations[0].S3StorageLocationParams.Encryption.KmsKeyId),
					),
			),
		)

		err := client.ExternalVolumes.Alter(ctx, req)
		require.NoError(t, err)

		props, err := client.ExternalVolumes.Describe(ctx, id)
		require.NoError(t, err)

		parsedExternalVolumeDescribed, err := helpers.ParseExternalVolumeDescribed(props)
		require.NoError(t, err)
		expectedParsedExternalVolumeDescribed := helpers.ParsedExternalVolumeDescribed{
			StorageLocations: []helpers.StorageLocation{
				{
					Name:                 gcsStorageLocationsNoneEncryption[0].GCSStorageLocationParams.Name,
					StorageProvider:      string(sdk.StorageProviderGCS),
					StorageBaseUrl:       gcsStorageLocationsNoneEncryption[0].GCSStorageLocationParams.StorageBaseUrl,
					StorageAwsRoleArn:    "",
					StorageAwsExternalId: "",
					EncryptionType:       string(gcsStorageLocationsNoneEncryption[0].GCSStorageLocationParams.Encryption.Type),
					EncryptionKmsKeyId:   "",
					AzureTenantId:        "",
				},
				{
					Name:                 s3StorageLocations[0].S3StorageLocationParams.Name,
					StorageProvider:      string(s3StorageLocations[0].S3StorageLocationParams.StorageProvider),
					StorageBaseUrl:       s3StorageLocations[0].S3StorageLocationParams.StorageBaseUrl,
					StorageAwsRoleArn:    s3StorageLocations[0].S3StorageLocationParams.StorageAwsRoleArn,
					StorageAwsExternalId: *s3StorageLocations[0].S3StorageLocationParams.StorageAwsExternalId,
					EncryptionType:       string(s3StorageLocations[0].S3StorageLocationParams.Encryption.Type),
					EncryptionKmsKeyId:   *s3StorageLocations[0].S3StorageLocationParams.Encryption.KmsKeyId,
					AzureTenantId:        "",
				},
			},
			Active:      "",
			Comment:     comment,
			AllowWrites: strconv.FormatBool(allowWrites),
		}

		assert.Equal(t, expectedParsedExternalVolumeDescribed, parsedExternalVolumeDescribed)
	})

	t.Run("Describe", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(
			t,
			append(append(append(append(append(append(s3StorageLocations, gcsStorageLocationsNoneEncryption...), azureStorageLocations...), s3StorageLocationsNoneEncryption...), gcsStorageLocations...), s3StorageLocationsNoEncryption...), gcsStorageLocationsNoEncryption...),
			allowWrites,
			&comment,
		)

		props, err := client.ExternalVolumes.Describe(ctx, id)
		require.NoError(t, err)

		parsedExternalVolumeDescribed, err := helpers.ParseExternalVolumeDescribed(props)
		require.NoError(t, err)
		expectedParsedExternalVolumeDescribed := helpers.ParsedExternalVolumeDescribed{
			StorageLocations: []helpers.StorageLocation{
				{
					Name:                 s3StorageLocations[0].S3StorageLocationParams.Name,
					StorageProvider:      string(s3StorageLocations[0].S3StorageLocationParams.StorageProvider),
					StorageBaseUrl:       s3StorageLocations[0].S3StorageLocationParams.StorageBaseUrl,
					StorageAwsRoleArn:    s3StorageLocations[0].S3StorageLocationParams.StorageAwsRoleArn,
					StorageAwsExternalId: *s3StorageLocations[0].S3StorageLocationParams.StorageAwsExternalId,
					EncryptionType:       string(s3StorageLocations[0].S3StorageLocationParams.Encryption.Type),
					EncryptionKmsKeyId:   *s3StorageLocations[0].S3StorageLocationParams.Encryption.KmsKeyId,
					AzureTenantId:        "",
				},
				{
					Name:                 gcsStorageLocationsNoneEncryption[0].GCSStorageLocationParams.Name,
					StorageProvider:      string(sdk.StorageProviderGCS),
					StorageBaseUrl:       gcsStorageLocationsNoneEncryption[0].GCSStorageLocationParams.StorageBaseUrl,
					StorageAwsRoleArn:    "",
					StorageAwsExternalId: "",
					EncryptionType:       string(gcsStorageLocationsNoneEncryption[0].GCSStorageLocationParams.Encryption.Type),
					EncryptionKmsKeyId:   "",
					AzureTenantId:        "",
				},
				{
					Name:                 azureStorageLocations[0].AzureStorageLocationParams.Name,
					StorageProvider:      string(sdk.StorageProviderAzure),
					StorageBaseUrl:       azureStorageLocations[0].AzureStorageLocationParams.StorageBaseUrl,
					StorageAwsRoleArn:    "",
					StorageAwsExternalId: "",
					EncryptionType:       "NONE",
					EncryptionKmsKeyId:   "",
					AzureTenantId:        azureStorageLocations[0].AzureStorageLocationParams.AzureTenantId,
				},
				{
					Name:                 s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.Name,
					StorageProvider:      string(s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageProvider),
					StorageBaseUrl:       s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageBaseUrl,
					StorageAwsRoleArn:    s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageAwsRoleArn,
					StorageAwsExternalId: *s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageAwsExternalId,
					EncryptionType:       string(s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.Encryption.Type),
					EncryptionKmsKeyId:   "",
					AzureTenantId:        "",
				},
				{
					Name:                 gcsStorageLocations[0].GCSStorageLocationParams.Name,
					StorageProvider:      string(sdk.StorageProviderGCS),
					StorageBaseUrl:       gcsStorageLocations[0].GCSStorageLocationParams.StorageBaseUrl,
					StorageAwsRoleArn:    "",
					StorageAwsExternalId: "",
					EncryptionType:       string(gcsStorageLocations[0].GCSStorageLocationParams.Encryption.Type),
					EncryptionKmsKeyId:   *gcsStorageLocations[0].GCSStorageLocationParams.Encryption.KmsKeyId,
					AzureTenantId:        "",
				},
				{
					Name:                 s3StorageLocationsNoEncryption[0].S3StorageLocationParams.Name,
					StorageProvider:      string(s3StorageLocationsNoEncryption[0].S3StorageLocationParams.StorageProvider),
					StorageBaseUrl:       s3StorageLocationsNoEncryption[0].S3StorageLocationParams.StorageBaseUrl,
					StorageAwsRoleArn:    s3StorageLocationsNoEncryption[0].S3StorageLocationParams.StorageAwsRoleArn,
					StorageAwsExternalId: *s3StorageLocationsNoEncryption[0].S3StorageLocationParams.StorageAwsExternalId,
					EncryptionType:       "NONE",
					EncryptionKmsKeyId:   "",
					AzureTenantId:        "",
				},
				{
					Name:                 gcsStorageLocationsNoEncryption[0].GCSStorageLocationParams.Name,
					StorageProvider:      string(sdk.StorageProviderGCS),
					StorageBaseUrl:       gcsStorageLocationsNoEncryption[0].GCSStorageLocationParams.StorageBaseUrl,
					StorageAwsRoleArn:    "",
					StorageAwsExternalId: "",
					EncryptionType:       "NONE",
					EncryptionKmsKeyId:   "",
					AzureTenantId:        "",
				},
			},
			Active:      "",
			Comment:     comment,
			AllowWrites: strconv.FormatBool(allowWrites),
		}

		assert.Equal(t, expectedParsedExternalVolumeDescribed, parsedExternalVolumeDescribed)
	})

	t.Run("Show with like", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, s3StorageLocations, allowWrites, &comment)
		name := id.Name()
		req := sdk.NewShowExternalVolumeRequest().WithLike(sdk.Like{Pattern: &name})

		externalVolumes, err := client.ExternalVolumes.Show(ctx, req)
		require.NoError(t, err)

		assert.Equal(t, 1, len(externalVolumes))
		assertExternalVolumeShowResult(t, &externalVolumes[0], id, allowWrites, comment)
	})
}

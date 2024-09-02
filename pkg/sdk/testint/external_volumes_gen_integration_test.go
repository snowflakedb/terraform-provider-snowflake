package testint

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Cannot test awsgov on non-gov Snowflake deployments

func TestInt_ExternalVolumes(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	awsBaseUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsRoleARN := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)
	awsKmsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsExternalId := "123456789"

	gcsBaseUrl := testenvs.GetOrSkipTest(t, testenvs.GcsExternalBuckerUrl)
	// //TODO - To test external volumes wth gcs storage locations that have encryption TYPE=GCS_SSE_KMS a gcs KMS_KEY_ID is required, which is not currently available in testenvs

	azureBaseUrl := testenvs.GetOrSkipTest(t, testenvs.AzureExternalBucketUrl) // TODO verify assumption on form
	azureTenantId := testenvs.GetOrSkipTest(t, testenvs.AzureExternalTenantId)

	assertExternalVolumeShowResult := func(t *testing.T, s *sdk.ExternalVolume, name sdk.AccountObjectIdentifier, allowWrites bool, comment string) {
		t.Helper()
		assert.Equal(t, name.Name(), s.Name)
		assert.Equal(t, strconv.FormatBool(allowWrites), s.AllowWrites)
		assert.Equal(t, comment, s.Comment)
	}

	type ExternalVolumePropNameValue struct {
		Name  string
		Value string
	}

	type S3StorageLocation struct {
		Name                    string   `json:"NAME"`
		StorageProvider         string   `json:"STORAGE_PROVIDER"`
		StorageBaseUrl          string   `json:"STORAGE_BASE_URL"`
		StorageAllowedLocations []string `json:"STORAGE_ALLOWED_LOCATIONS"`
		StorageAwsRoleArn       string   `json:"STORAGE_AWS_ROLE_ARN"`
		StroageAwsIamUserArn    string   `json:"STORAGE_AWS_IAM_USER_ARN"`
		StorageAwsExternalId    string   `json:"STORAGE_AWS_EXTERNAL_ID"`
		EncryptionType          string   `json:"ENCRYPTION_TYPE"`
		EncryptionKmsId         string   `json:"ENCRYPTION_KMS_KEY_ID"`
	}

	type S3StorageLocationTrimmed struct {
		Name                 string `json:"NAME"`
		StorageProvider      string `json:"STORAGE_PROVIDER"`
		StorageBaseUrl       string `json:"STORAGE_BASE_URL"`
		StorageAwsRoleArn    string `json:"STORAGE_AWS_ROLE_ARN"`
		StorageAwsExternalId string `json:"STORAGE_AWS_EXTERNAL_ID"`
		EncryptionType       string `json:"ENCRYPTION_TYPE"`
		EncryptionKmsId      string `json:"ENCRYPTION_KMS_KEY_ID,omitempty"`
	}

	type GCSStorageLocation struct {
		Name                     string   `json:"NAME"`
		StorageProvider          string   `json:"STORAGE_PROVIDER"`
		StorageBaseUrl           string   `json:"STORAGE_BASE_URL"`
		StorageAllowedLocations  []string `json:"STORAGE_ALLOWED_LOCATIONS"`
		StorageGcpServiceAccount string   `json:"STORAGE_GCP_SERVICE_ACCOUNT"`
		EncryptionType           string   `json:"ENCRYPTION_TYPE"`
		EncryptionKmsId          string   `json:"ENCRYPTION_KMS_KEY_ID"`
	}

	type GCSStorageLocationTrimmed struct {
		Name            string `json:"NAME"`
		StorageProvider string `json:"STORAGE_PROVIDER"`
		StorageBaseUrl  string `json:"STORAGE_BASE_URL"`
		EncryptionType  string `json:"ENCRYPTION_TYPE"`
		EncryptionKmsId string `json:"ENCRYPTION_KMS_KEY_ID,omitempty"`
	}

	type AzureStorageLocation struct {
		Name                    string   `json:"NAME"`
		StorageProvider         string   `json:"STORAGE_PROVIDER"`
		StorageBaseUrl          string   `json:"STORAGE_BASE_URL"`
		StorageAllowedLocations []string `json:"STORAGE_ALLOWED_LOCATIONS"`
		AzureTenantId           string   `json:"AZURE_TENANT_ID"`
		AzureMultiTenantAppName string   `json:"AZURE_MULTI_TENANT_APP_NAME"`
		AzureConsentUrl         string   `json:"AZURE_CONSENT_URL"`
		EncryptionType          string   `json:"ENCRYPTION_TYPE"`
		EncryptionKmsId         string   `json:"ENCRYPTION_KMS_KEY_ID"`
	}

	type AzureStorageLocationTrimmed struct {
		Name            string `json:"NAME"`
		StorageProvider string `json:"STORAGE_PROVIDER"`
		StorageBaseUrl  string `json:"STORAGE_BASE_URL"`
		AzureTenantId   string `json:"AZURE_TENANT_ID"`
	}

	// Enforce only property names and values in tests, not parent_property, type and property_default
	// In addition the storage location properties are trimmed to only contain values that we set
	trimProperties := func(t *testing.T, props []sdk.ExternalVolumeProperty) []ExternalVolumePropNameValue {
		t.Helper()
		var externalVolumePropNameValue []ExternalVolumePropNameValue
		for _, p := range props {
			if strings.Contains(p.Name, "STORAGE_LOCATION_") {
				if strings.Contains(p.Value, `"STORAGE_PROVIDER":"S3"`) {
					s3StorageLocation := S3StorageLocation{}
					json.Unmarshal([]byte(p.Value), &s3StorageLocation)
					s3StorageLocationTrimmed := S3StorageLocationTrimmed{
						Name:                 s3StorageLocation.Name,
						StorageProvider:      s3StorageLocation.StorageProvider,
						StorageBaseUrl:       s3StorageLocation.StorageBaseUrl,
						StorageAwsRoleArn:    s3StorageLocation.StorageAwsRoleArn,
						StorageAwsExternalId: s3StorageLocation.StorageAwsExternalId,
						EncryptionType:       s3StorageLocation.EncryptionType,
						EncryptionKmsId:      s3StorageLocation.EncryptionKmsId,
					}
					s3StorageLocationTrimmedMarshaled, err := json.Marshal(s3StorageLocationTrimmed)
					require.NoError(t, err)
					externalVolumePropNameValue = append(
						externalVolumePropNameValue,
						ExternalVolumePropNameValue{Name: p.Name, Value: string(s3StorageLocationTrimmedMarshaled)},
					)
				} else if strings.Contains(p.Value, `"STORAGE_PROVIDER":"GCS"`) {
					gcsStorageLocation := GCSStorageLocation{}
					json.Unmarshal([]byte(p.Value), &gcsStorageLocation)
					gcsStorageLocationTrimmed := GCSStorageLocationTrimmed{
						Name:            gcsStorageLocation.Name,
						StorageProvider: gcsStorageLocation.StorageProvider,
						StorageBaseUrl:  gcsStorageLocation.StorageBaseUrl,
						EncryptionType:  gcsStorageLocation.EncryptionType,
						EncryptionKmsId: gcsStorageLocation.EncryptionKmsId,
					}
					gcsStorageLocationTrimmedMarshaled, err := json.Marshal(gcsStorageLocationTrimmed)
					require.NoError(t, err)
					externalVolumePropNameValue = append(
						externalVolumePropNameValue,
						ExternalVolumePropNameValue{Name: p.Name, Value: string(gcsStorageLocationTrimmedMarshaled)},
					)
				} else if strings.Contains(p.Value, `"STORAGE_PROVIDER":"AZURE"`) {
					azureStorageLocation := AzureStorageLocation{}
					json.Unmarshal([]byte(p.Value), &azureStorageLocation)
					azureStorageLocationTrimmed := AzureStorageLocationTrimmed{
						Name:            azureStorageLocation.Name,
						StorageProvider: azureStorageLocation.StorageProvider,
						StorageBaseUrl:  azureStorageLocation.StorageBaseUrl,
						AzureTenantId:   azureStorageLocation.AzureTenantId,
					}
					azureStorageLocationTrimmedMarshaled, err := json.Marshal(azureStorageLocationTrimmed)
					require.NoError(t, err)
					externalVolumePropNameValue = append(
						externalVolumePropNameValue,
						ExternalVolumePropNameValue{Name: p.Name, Value: string(azureStorageLocationTrimmedMarshaled)},
					)
				} else {
					panic("Unrecognised storage provider in storage location property")
				}
			} else {
				externalVolumePropNameValue = append(externalVolumePropNameValue, ExternalVolumePropNameValue{Name: p.Name, Value: p.Value})
			}
		}

		return externalVolumePropNameValue
	}

	s3StorageLocations := []sdk.ExternalVolumeStorageLocation{
		{
			S3StorageLocationParams: &sdk.S3StorageLocationParams{
				Name:                 "s3_testing_storage_location",
				StorageProvider:      sdk.S3StorageProviderS3,
				StorageAwsRoleArn:    awsRoleARN,
				StorageBaseUrl:       awsBaseUrl,
				StorageAwsExternalId: sdk.String(awsExternalId),
				Encryption: sdk.ExternalVolumeS3Encryption{
					Type:     sdk.S3EncryptionTypeSseKms,
					KmsKeyId: &awsKmsKeyId,
				},
			},
		},
	}

	s3StorageLocationsNoneEncryption := []sdk.ExternalVolumeStorageLocation{
		{
			S3StorageLocationParams: &sdk.S3StorageLocationParams{
				Name:                 "s3_testing_storage_location_no_encryption",
				StorageProvider:      sdk.S3StorageProviderS3,
				StorageAwsRoleArn:    awsRoleARN,
				StorageBaseUrl:       awsBaseUrl,
				StorageAwsExternalId: sdk.String(awsExternalId),
				Encryption: sdk.ExternalVolumeS3Encryption{
					Type: sdk.S3EncryptionNone,
				},
			},
		},
	}

	gcsStorageLocations := []sdk.ExternalVolumeStorageLocation{
		{
			GCSStorageLocationParams: &sdk.GCSStorageLocationParams{
				Name:           "gcs_testing_storage_location",
				StorageBaseUrl: gcsBaseUrl,
				Encryption: sdk.ExternalVolumeGCSEncryption{
					Type: sdk.GCSEncryptionTypeNone,
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

	createExternalVolume := func(t *testing.T, storageLocations []sdk.ExternalVolumeStorageLocation, allowWrites bool, comment string) sdk.AccountObjectIdentifier {
		t.Helper()

		// id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		// TODO remove this
		id := sdk.NewAccountObjectIdentifier("test_external_volume")
		req := sdk.NewCreateExternalVolumeRequest(id, storageLocations).
			WithIfNotExists(true).
			WithAllowWrites(allowWrites).
			WithComment(comment)

		err := client.ExternalVolumes.Create(ctx, req)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.ExternalVolumes.Drop(ctx, sdk.NewDropExternalVolumeRequest(id))
			require.NoError(t, err)
		})

		return id
	}

	t.Run("Create - S3 Storage Location", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, s3StorageLocations, allowWrites, comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertExternalVolumeShowResult(t, externalVolume, id, allowWrites, comment)
	})

	t.Run("Create - S3 Storage Location No Encryption", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, s3StorageLocationsNoneEncryption, allowWrites, comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertExternalVolumeShowResult(t, externalVolume, id, allowWrites, comment)
	})

	t.Run("Create - GCS Storage Location", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, gcsStorageLocations, allowWrites, comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertExternalVolumeShowResult(t, externalVolume, id, allowWrites, comment)
	})

	t.Run("Create - Azure Storage Location", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, azureStorageLocations, allowWrites, comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertExternalVolumeShowResult(t, externalVolume, id, allowWrites, comment)
	})

	t.Run("Create - Multiple Storage Locations", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, append(append(s3StorageLocations, gcsStorageLocations...), azureStorageLocations...), allowWrites, comment)

		externalVolume, err := client.ExternalVolumes.ShowByID(ctx, id)
		require.NoError(t, err)

		assertExternalVolumeShowResult(t, externalVolume, id, allowWrites, comment)
	})

	t.Run("Alter - remove storage location", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, append(s3StorageLocationsNoneEncryption, gcsStorageLocations...), allowWrites, comment)

		req := sdk.NewAlterExternalVolumeRequest(id).WithRemoveStorageLocation(gcsStorageLocations[0].GCSStorageLocationParams.Name)

		err := client.ExternalVolumes.Alter(ctx, req)
		require.NoError(t, err)

		props, err := client.ExternalVolumes.Describe(ctx, id)
		require.NoError(t, err)

		trimmedProperties := trimProperties(t, props)
		assert.Equal(t, 4, len(trimmedProperties))
		assert.Contains(t, trimmedProperties, ExternalVolumePropNameValue{Name: "COMMENT", Value: comment})
		assert.Contains(t, trimmedProperties, ExternalVolumePropNameValue{Name: "ALLOW_WRITES", Value: strconv.FormatBool(allowWrites)})
		assert.Contains(t, trimmedProperties, ExternalVolumePropNameValue{Name: "ACTIVE", Value: ""})
		assert.Contains(
			t,
			trimmedProperties,
			ExternalVolumePropNameValue{
				Name: "STORAGE_LOCATION_1",
				Value: fmt.Sprintf(
					`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_AWS_ROLE_ARN":"%s","STORAGE_AWS_EXTERNAL_ID":"%s","ENCRYPTION_TYPE":"%s"}`,
					s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.Name,
					s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageProvider,
					s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageBaseUrl,
					s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageAwsRoleArn,
					*s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageAwsExternalId,
					s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.Encryption.Type,
				),
			},
		)
	})

	t.Run("Alter - set comment", func(t *testing.T) {
		allowWrites := true
		comment := ""
		id := createExternalVolume(t, s3StorageLocationsNoneEncryption, allowWrites, "some comment")

		req := sdk.NewAlterExternalVolumeRequest(id).WithSet(
			*sdk.NewAlterExternalVolumeSetRequest().WithComment(comment),
		)

		err := client.ExternalVolumes.Alter(ctx, req)
		require.NoError(t, err)

		props, err := client.ExternalVolumes.Describe(ctx, id)
		require.NoError(t, err)

		trimmedProperties := trimProperties(t, props)
		assert.Equal(t, 3, len(trimmedProperties))
		assert.Contains(t, trimmedProperties, ExternalVolumePropNameValue{Name: "ALLOW_WRITES", Value: strconv.FormatBool(allowWrites)})
		assert.Contains(t, trimmedProperties, ExternalVolumePropNameValue{Name: "ACTIVE", Value: ""})
		assert.Contains(
			t,
			trimmedProperties,
			ExternalVolumePropNameValue{
				Name: "STORAGE_LOCATION_1",
				Value: fmt.Sprintf(
					`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_AWS_ROLE_ARN":"%s","STORAGE_AWS_EXTERNAL_ID":"%s","ENCRYPTION_TYPE":"%s"}`,
					s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.Name,
					s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageProvider,
					s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageBaseUrl,
					s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageAwsRoleArn,
					*s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageAwsExternalId,
					s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.Encryption.Type,
				),
			},
		)
	})

	t.Run("Alter - set allow writes", func(t *testing.T) {
		allowWrites := false
		comment := "some comment"
		id := createExternalVolume(t, s3StorageLocations, true, comment)

		req := sdk.NewAlterExternalVolumeRequest(id).WithSet(
			*sdk.NewAlterExternalVolumeSetRequest().WithAllowWrites(allowWrites),
		)

		err := client.ExternalVolumes.Alter(ctx, req)
		require.NoError(t, err)

		props, err := client.ExternalVolumes.Describe(ctx, id)
		require.NoError(t, err)

		trimmedProperties := trimProperties(t, props)
		assert.Equal(t, 4, len(trimmedProperties))
		assert.Contains(t, trimmedProperties, ExternalVolumePropNameValue{Name: "COMMENT", Value: comment})
		assert.Contains(t, trimmedProperties, ExternalVolumePropNameValue{Name: "ALLOW_WRITES", Value: strconv.FormatBool(allowWrites)})
		assert.Contains(t, trimmedProperties, ExternalVolumePropNameValue{Name: "ACTIVE", Value: ""})
		assert.Contains(
			t,
			trimmedProperties,
			ExternalVolumePropNameValue{
				Name: "STORAGE_LOCATION_1",
				Value: fmt.Sprintf(
					`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_AWS_ROLE_ARN":"%s","STORAGE_AWS_EXTERNAL_ID":"%s","ENCRYPTION_TYPE":"%s","ENCRYPTION_KMS_KEY_ID":"%s"}`,
					s3StorageLocations[0].S3StorageLocationParams.Name,
					s3StorageLocations[0].S3StorageLocationParams.StorageProvider,
					s3StorageLocations[0].S3StorageLocationParams.StorageBaseUrl,
					s3StorageLocations[0].S3StorageLocationParams.StorageAwsRoleArn,
					*s3StorageLocations[0].S3StorageLocationParams.StorageAwsExternalId,
					s3StorageLocations[0].S3StorageLocationParams.Encryption.Type,
					*s3StorageLocations[0].S3StorageLocationParams.Encryption.KmsKeyId,
				),
			},
		)
	})

	t.Run("Alter - add s3 storage location to external volume", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, gcsStorageLocations, allowWrites, comment)

		req := sdk.NewAlterExternalVolumeRequest(id).WithAddStorageLocation(
			*sdk.NewExternalVolumeStorageLocationRequest().WithS3StorageLocationParams(
				*sdk.NewS3StorageLocationParamsRequest(
					s3StorageLocations[0].S3StorageLocationParams.Name,
					s3StorageLocations[0].S3StorageLocationParams.StorageProvider,
					s3StorageLocations[0].S3StorageLocationParams.StorageAwsRoleArn,
					s3StorageLocations[0].S3StorageLocationParams.StorageBaseUrl,
					*sdk.NewExternalVolumeS3EncryptionRequest(
						s3StorageLocations[0].S3StorageLocationParams.Encryption.Type,
					).WithKmsKeyId(*s3StorageLocations[0].S3StorageLocationParams.Encryption.KmsKeyId),
				).WithStorageAwsExternalId(*s3StorageLocations[0].S3StorageLocationParams.StorageAwsExternalId),
			),
		)

		err := client.ExternalVolumes.Alter(ctx, req)
		require.NoError(t, err)

		props, err := client.ExternalVolumes.Describe(ctx, id)
		require.NoError(t, err)

		trimmedProperties := trimProperties(t, props)
		assert.Equal(t, 5, len(trimmedProperties))
		assert.Contains(t, trimmedProperties, ExternalVolumePropNameValue{Name: "COMMENT", Value: comment})
		assert.Contains(t, trimmedProperties, ExternalVolumePropNameValue{Name: "ALLOW_WRITES", Value: strconv.FormatBool(allowWrites)})
		assert.Contains(t, trimmedProperties, ExternalVolumePropNameValue{Name: "ACTIVE", Value: ""})
		assert.Contains(
			t,
			trimmedProperties,
			ExternalVolumePropNameValue{
				Name: "STORAGE_LOCATION_1",
				Value: fmt.Sprintf(
					`{"NAME":"%s","STORAGE_PROVIDER":"GCS","STORAGE_BASE_URL":"%s","ENCRYPTION_TYPE":"%s"}`,
					gcsStorageLocations[0].GCSStorageLocationParams.Name,
					gcsStorageLocations[0].GCSStorageLocationParams.StorageBaseUrl,
					gcsStorageLocations[0].GCSStorageLocationParams.Encryption.Type,
				),
			},
		)
		assert.Contains(
			t,
			trimmedProperties,
			ExternalVolumePropNameValue{
				Name: "STORAGE_LOCATION_2",
				Value: fmt.Sprintf(
					`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_AWS_ROLE_ARN":"%s","STORAGE_AWS_EXTERNAL_ID":"%s","ENCRYPTION_TYPE":"%s","ENCRYPTION_KMS_KEY_ID":"%s"}`,
					s3StorageLocations[0].S3StorageLocationParams.Name,
					s3StorageLocations[0].S3StorageLocationParams.StorageProvider,
					s3StorageLocations[0].S3StorageLocationParams.StorageBaseUrl,
					s3StorageLocations[0].S3StorageLocationParams.StorageAwsRoleArn,
					*s3StorageLocations[0].S3StorageLocationParams.StorageAwsExternalId,
					s3StorageLocations[0].S3StorageLocationParams.Encryption.Type,
					*s3StorageLocations[0].S3StorageLocationParams.Encryption.KmsKeyId,
				),
			},
		)
	})

	t.Run("Describe", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(
			t,
			append(append(append(s3StorageLocations, gcsStorageLocations...), azureStorageLocations...), s3StorageLocationsNoneEncryption...),
			allowWrites,
			comment,
		)

		props, err := client.ExternalVolumes.Describe(ctx, id)
		require.NoError(t, err)

		trimmedProperties := trimProperties(t, props)
		assert.Equal(t, 7, len(trimmedProperties))
		assert.Contains(t, trimmedProperties, ExternalVolumePropNameValue{Name: "COMMENT", Value: comment})
		assert.Contains(t, trimmedProperties, ExternalVolumePropNameValue{Name: "ALLOW_WRITES", Value: strconv.FormatBool(allowWrites)})
		assert.Contains(t, trimmedProperties, ExternalVolumePropNameValue{Name: "ACTIVE", Value: ""})
		assert.Contains(
			t,
			trimmedProperties,
			ExternalVolumePropNameValue{
				Name: "STORAGE_LOCATION_1",
				Value: fmt.Sprintf(
					`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_AWS_ROLE_ARN":"%s","STORAGE_AWS_EXTERNAL_ID":"%s","ENCRYPTION_TYPE":"%s","ENCRYPTION_KMS_KEY_ID":"%s"}`,
					s3StorageLocations[0].S3StorageLocationParams.Name,
					s3StorageLocations[0].S3StorageLocationParams.StorageProvider,
					s3StorageLocations[0].S3StorageLocationParams.StorageBaseUrl,
					s3StorageLocations[0].S3StorageLocationParams.StorageAwsRoleArn,
					*s3StorageLocations[0].S3StorageLocationParams.StorageAwsExternalId,
					s3StorageLocations[0].S3StorageLocationParams.Encryption.Type,
					*s3StorageLocations[0].S3StorageLocationParams.Encryption.KmsKeyId,
				),
			},
		)
		assert.Contains(
			t,
			trimmedProperties,
			ExternalVolumePropNameValue{
				Name: "STORAGE_LOCATION_2",
				Value: fmt.Sprintf(
					`{"NAME":"%s","STORAGE_PROVIDER":"GCS","STORAGE_BASE_URL":"%s","ENCRYPTION_TYPE":"%s"}`,
					gcsStorageLocations[0].GCSStorageLocationParams.Name,
					gcsStorageLocations[0].GCSStorageLocationParams.StorageBaseUrl,
					gcsStorageLocations[0].GCSStorageLocationParams.Encryption.Type,
				),
			},
		)
		assert.Contains(
			t,
			trimmedProperties,
			ExternalVolumePropNameValue{
				Name: "STORAGE_LOCATION_3",
				Value: fmt.Sprintf(
					`{"NAME":"%s","STORAGE_PROVIDER":"AZURE","STORAGE_BASE_URL":"%s","AZURE_TENANT_ID":"%s"}`,
					azureStorageLocations[0].AzureStorageLocationParams.Name,
					azureStorageLocations[0].AzureStorageLocationParams.StorageBaseUrl,
					azureStorageLocations[0].AzureStorageLocationParams.AzureTenantId,
				),
			},
		)
		assert.Contains(
			t,
			trimmedProperties,
			ExternalVolumePropNameValue{
				Name: "STORAGE_LOCATION_4",
				Value: fmt.Sprintf(
					`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_AWS_ROLE_ARN":"%s","STORAGE_AWS_EXTERNAL_ID":"%s","ENCRYPTION_TYPE":"%s"}`,
					s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.Name,
					s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageProvider,
					s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageBaseUrl,
					s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageAwsRoleArn,
					*s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.StorageAwsExternalId,
					s3StorageLocationsNoneEncryption[0].S3StorageLocationParams.Encryption.Type,
				),
			},
		)
	})

	t.Run("Show with like", func(t *testing.T) {
		allowWrites := true
		comment := "some comment"
		id := createExternalVolume(t, s3StorageLocations, allowWrites, comment)
		name := id.Name()
		req := sdk.NewShowExternalVolumeRequest().WithLike(sdk.Like{Pattern: &name})

		externalVolumes, err := client.ExternalVolumes.Show(ctx, req)
		require.NoError(t, err)

		assert.Equal(t, 1, len(externalVolumes))
		assertExternalVolumeShowResult(t, &externalVolumes[0], id, allowWrites, comment)
	})
}

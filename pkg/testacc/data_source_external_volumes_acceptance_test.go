//go:build non_account_level_tests

package testacc

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ExternalVolumes_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	idOne := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := testClient().Ids.RandomAccountObjectIdentifier()

	storageLocations := []sdk.ExternalVolumeStorageLocationItem{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: "my-s3-us-west-2",
			S3StorageLocationParams: &sdk.S3StorageLocationParams{
				StorageProvider:   sdk.S3StorageProviderS3,
				StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole",
				StorageBaseUrl:    "s3://my-example-bucket/",
			},
		}},
	}

	_, cleanupOne := testClient().ExternalVolume.CreateWithRequest(t, sdk.NewCreateExternalVolumeRequest(idOne, storageLocations))
	t.Cleanup(cleanupOne)
	_, cleanupTwo := testClient().ExternalVolume.CreateWithRequest(t, sdk.NewCreateExternalVolumeRequest(idTwo, storageLocations))
	t.Cleanup(cleanupTwo)
	_, cleanupThree := testClient().ExternalVolume.CreateWithRequest(t, sdk.NewCreateExternalVolumeRequest(idThree, storageLocations))
	t.Cleanup(cleanupThree)

	externalVolumesModelLikeFirst := datasourcemodel.ExternalVolumes("test").
		WithLike(idOne.Name()).
		WithWithDescribe(false)

	externalVolumesModelLikePrefix := datasourcemodel.ExternalVolumes("test").
		WithLike(prefix + "%").
		WithWithDescribe(false)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// like (exact)
			{
				Config: accconfig.FromModels(t, externalVolumesModelLikeFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(externalVolumesModelLikeFirst.DatasourceReference(), "external_volumes.#", "1"),
				),
			},
			// like (prefix)
			{
				Config: accconfig.FromModels(t, externalVolumesModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(externalVolumesModelLikePrefix.DatasourceReference(), "external_volumes.#", "2"),
				),
			},
		},
	})
}

func TestAcc_ExternalVolumes_CompleteUseCase(t *testing.T) {
	if testenvs.GetSnowflakeEnvironmentWithProdDefault() == testenvs.SnowflakePreProdGovEnvironment {
		t.Skip("Skipping test - Snowflake error 393962 (42601): External volumes in government deployments cannot use non-government S3 storage locations. Storage type S3 is not allowed. Please use S3GOV instead.")
	}
	s3CompatAwsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	s3LocationName := "s3_loc"
	s3BaseUrl := "s3://my-example-bucket/"
	s3RoleArn := "arn:aws:iam::123456789012:role/myrole"

	gcsLocationName := "gcs_loc"
	gcsBaseUrl := "gcs://my-example-bucket/"

	azureLocationName := "azure_loc"
	azureBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container/"
	azureTenantId := "123456789"

	s3CompatLocationName := "s3compat_loc"
	s3CompatBaseUrl := strings.Replace(s3BaseUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"
	s3CompatAwsKeyId := "AKIAIOSFODNN7EXAMPLE"

	id := testClient().Ids.RandomAccountObjectIdentifier()
	storageLocations := []sdk.ExternalVolumeStorageLocationItem{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: s3LocationName,
			S3StorageLocationParams: &sdk.S3StorageLocationParams{
				StorageProvider:   sdk.S3StorageProviderS3,
				StorageAwsRoleArn: s3RoleArn,
				StorageBaseUrl:    s3BaseUrl,
			},
		}},
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: gcsLocationName,
			GCSStorageLocationParams: &sdk.GCSStorageLocationParams{
				StorageBaseUrl: gcsBaseUrl,
			},
		}},
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: azureLocationName,
			AzureStorageLocationParams: &sdk.AzureStorageLocationParams{
				AzureTenantId:  azureTenantId,
				StorageBaseUrl: azureBaseUrl,
			},
		}},
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: s3CompatLocationName,
			S3CompatStorageLocationParams: &sdk.S3CompatStorageLocationParams{
				StorageBaseUrl:  s3CompatBaseUrl,
				StorageEndpoint: s3CompatEndpoint,
				Credentials: sdk.ExternalVolumeS3CompatCredentials{
					AwsKeyId:     s3CompatAwsKeyId,
					AwsSecretKey: s3CompatAwsSecretKey,
				},
			},
		}},
	}

	_, cleanup := testClient().ExternalVolume.CreateWithRequest(t, sdk.NewCreateExternalVolumeRequest(id, storageLocations))
	t.Cleanup(cleanup)

	externalVolumesModel := datasourcemodel.ExternalVolumes("test").
		WithLike(id.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, externalVolumesModel),
				Check: assertThat(t,
					resourceshowoutputassert.ExternalVolumesDatasourceShowOutput(t, externalVolumesModel.DatasourceReference()).
						HasName(id.Name()).
						HasAllowWrites(true).
						HasCommentEmpty(),
					resourceshowoutputassert.ExternalVolumesDatasourceDescribeOutput(t, externalVolumesModel.DatasourceReference()).
						HasAllowWrites("true").
						HasActive("").
						HasCommentEmpty().
						HasStorageLocations([]sdk.ExternalVolumeStorageLocationDetails{
							{
								Name:            s3LocationName,
								StorageProvider: string(sdk.StorageProviderS3),
								StorageBaseUrl:  s3BaseUrl,
								EncryptionType:  "NONE",
								S3StorageLocation: &sdk.StorageLocationS3Details{
									StorageAwsRoleArn: s3RoleArn,
								},
							},
							{
								Name:               gcsLocationName,
								StorageProvider:    string(sdk.StorageProviderGCS),
								StorageBaseUrl:     gcsBaseUrl,
								EncryptionType:     "NONE",
								GCSStorageLocation: &sdk.StorageLocationGcsDetails{},
							},
							{
								Name:            azureLocationName,
								StorageProvider: string(sdk.StorageProviderAzure),
								StorageBaseUrl:  azureBaseUrl,
								EncryptionType:  "NONE",
								AzureStorageLocation: &sdk.StorageLocationAzureDetails{
									AzureTenantId: azureTenantId,
								},
							},
							{
								Name:            s3CompatLocationName,
								StorageProvider: string(sdk.StorageProviderS3Compatible),
								StorageBaseUrl:  s3CompatBaseUrl,
								EncryptionType:  "NONE",
								S3CompatStorageLocation: &sdk.StorageLocationS3CompatDetails{
									Endpoint:       s3CompatEndpoint,
									AwsAccessKeyId: s3CompatAwsKeyId,
								},
							},
						}),
				),
			},
		},
	})
}

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

	storageLocations := []sdk.ExternalVolumeStorageLocationItemRequest{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocationRequest{
			Name: "my-s3-us-west-2",
			S3StorageLocationParams: &sdk.S3StorageLocationParamsRequest{
				StorageProvider:   sdk.S3StorageProviderS3,
				StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole",
				StorageBaseUrl:    "s3://my-example-bucket/",
			},
		}},
	}

	commentOne, commentTwo := random.Comment(), random.Comment()
	_, cleanupOne := testClient().ExternalVolume.CreateWithRequest(t, sdk.NewCreateExternalVolumeRequest(idOne, storageLocations).WithComment(commentOne))
	t.Cleanup(cleanupOne)
	_, cleanupTwo := testClient().ExternalVolume.CreateWithRequest(t, sdk.NewCreateExternalVolumeRequest(idTwo, storageLocations).WithComment(commentTwo))
	t.Cleanup(cleanupTwo)
	_, cleanupThree := testClient().ExternalVolume.CreateWithRequest(t, sdk.NewCreateExternalVolumeRequest(idThree, storageLocations))
	t.Cleanup(cleanupThree)

	externalVolumesModelLikeFirst := datasourcemodel.ExternalVolumes("test").
		WithLike(idOne.Name()).
		WithWithDescribe(false)

	externalVolumesModelLikePrefix := datasourcemodel.ExternalVolumes("test").
		WithLike(prefix + "%").
		WithWithDescribe(false)

	externalVolumesModelLikePrefixWithDescribe := datasourcemodel.ExternalVolumes("test").
		WithLike(prefix + "%")

	otherCasePrefix := strings.ToUpper(prefix)
	if otherCasePrefix == prefix {
		otherCasePrefix = strings.ToLower(prefix)
	}
	externalVolumesModelLikeOtherCase := datasourcemodel.ExternalVolumes("test").
		WithLike(otherCasePrefix + "%").
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
					resource.TestCheckResourceAttr(externalVolumesModelLikePrefix.DatasourceReference(), "external_volumes.0.describe_output.#", "0"),
					resource.TestCheckResourceAttr(externalVolumesModelLikePrefix.DatasourceReference(), "external_volumes.1.describe_output.#", "0"),
				),
			},
			// with more than one result, every item holds its own show and describe output
			{
				Config: accconfig.FromModels(t, externalVolumesModelLikePrefixWithDescribe),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(externalVolumesModelLikePrefixWithDescribe.DatasourceReference(), "id", "external_volumes_read"),
					resource.TestCheckResourceAttr(externalVolumesModelLikePrefixWithDescribe.DatasourceReference(), "external_volumes.#", "2"),
					resource.TestCheckResourceAttr(externalVolumesModelLikePrefixWithDescribe.DatasourceReference(), "external_volumes.0.describe_output.#", "1"),
					resource.TestCheckResourceAttrPair(
						externalVolumesModelLikePrefixWithDescribe.DatasourceReference(), "external_volumes.0.show_output.0.comment",
						externalVolumesModelLikePrefixWithDescribe.DatasourceReference(), "external_volumes.0.describe_output.0.comment",
					),
					resource.TestCheckResourceAttr(externalVolumesModelLikePrefixWithDescribe.DatasourceReference(), "external_volumes.1.describe_output.#", "1"),
					resource.TestCheckResourceAttrPair(
						externalVolumesModelLikePrefixWithDescribe.DatasourceReference(), "external_volumes.1.show_output.0.comment",
						externalVolumesModelLikePrefixWithDescribe.DatasourceReference(), "external_volumes.1.describe_output.0.comment",
					),
					resource.TestCheckResourceAttr(externalVolumesModelLikePrefixWithDescribe.DatasourceReference(), "external_volumes.0.describe_output.0.storage_locations.#", "1"),
					resource.TestCheckResourceAttr(externalVolumesModelLikePrefixWithDescribe.DatasourceReference(), "external_volumes.0.describe_output.0.storage_locations.0.name", "my-s3-us-west-2"),
					resource.TestCheckResourceAttr(externalVolumesModelLikePrefixWithDescribe.DatasourceReference(), "external_volumes.0.describe_output.0.storage_locations.0.storage_provider", string(sdk.S3StorageProviderS3)),
					resource.TestCheckResourceAttr(externalVolumesModelLikePrefixWithDescribe.DatasourceReference(), "external_volumes.0.describe_output.0.storage_locations.0.storage_base_url", "s3://my-example-bucket/"),
					resource.TestCheckResourceAttr(externalVolumesModelLikePrefixWithDescribe.DatasourceReference(), "external_volumes.0.describe_output.0.storage_locations.0.s3_storage_location.0.storage_aws_role_arn", "arn:aws:iam::123456789012:role/myrole"),
					resource.TestCheckResourceAttr(externalVolumesModelLikePrefixWithDescribe.DatasourceReference(), "external_volumes.1.describe_output.0.storage_locations.#", "1"),
					resource.TestCheckResourceAttr(externalVolumesModelLikePrefixWithDescribe.DatasourceReference(), "external_volumes.1.describe_output.0.storage_locations.0.name", "my-s3-us-west-2"),
					resource.TestCheckResourceAttr(externalVolumesModelLikePrefixWithDescribe.DatasourceReference(), "external_volumes.1.describe_output.0.storage_locations.0.storage_provider", string(sdk.S3StorageProviderS3)),
					resource.TestCheckResourceAttr(externalVolumesModelLikePrefixWithDescribe.DatasourceReference(), "external_volumes.1.describe_output.0.storage_locations.0.storage_base_url", "s3://my-example-bucket/"),
					resource.TestCheckResourceAttr(externalVolumesModelLikePrefixWithDescribe.DatasourceReference(), "external_volumes.1.describe_output.0.storage_locations.0.s3_storage_location.0.storage_aws_role_arn", "arn:aws:iam::123456789012:role/myrole"),
				),
			},
			// like is case-insensitive
			{
				Config: accconfig.FromModels(t, externalVolumesModelLikeOtherCase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(externalVolumesModelLikeOtherCase.DatasourceReference(), "external_volumes.#", "2"),
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
	storageLocations := []sdk.ExternalVolumeStorageLocationItemRequest{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocationRequest{
			Name: s3LocationName,
			S3StorageLocationParams: &sdk.S3StorageLocationParamsRequest{
				StorageProvider:   sdk.S3StorageProviderS3,
				StorageAwsRoleArn: s3RoleArn,
				StorageBaseUrl:    s3BaseUrl,
			},
		}},
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocationRequest{
			Name: gcsLocationName,
			GCSStorageLocationParams: &sdk.GCSStorageLocationParamsRequest{
				StorageBaseUrl: gcsBaseUrl,
			},
		}},
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocationRequest{
			Name: azureLocationName,
			AzureStorageLocationParams: &sdk.AzureStorageLocationParamsRequest{
				AzureTenantId:  azureTenantId,
				StorageBaseUrl: azureBaseUrl,
			},
		}},
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocationRequest{
			Name: s3CompatLocationName,
			S3CompatStorageLocationParams: &sdk.S3CompatStorageLocationParamsRequest{
				StorageBaseUrl:  s3CompatBaseUrl,
				StorageEndpoint: s3CompatEndpoint,
				Credentials: sdk.ExternalVolumeS3CompatCredentialsRequest{
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
				Check: assertThat(
					t,
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
								StorageProvider:    string(sdk.StorageProviderGcs),
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
								StorageProvider: string(sdk.StorageProviderS3compat),
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
